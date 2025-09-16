# `cj`: Conjure JSON

Package cj provides high-performance JSON marshaling and unmarshaling utilities specifically designed for conjure-go generated code.

# Overview

The cj package leverages Go's type system and generics to provide compile-time type safety and zero-allocation JSON
processing. It implements the encoding/json/v2 MarshalerTo and UnmarshalerFrom interfaces for all Conjure types,
enabling efficient serialization without runtime reflection or interface boxing.

## Zero-Allocation Design

The package's core innovation is the use of generic struct{} types as encoders and decoders. Each type encoder/decoder
is a zero-sized struct that implements TypeEncoder[T] or TypeDecoder[T]:

	type String[T ~string] struct{}

These struct{} types have several key advantages:
  - Zero memory allocation: struct{} has no fields and takes zero bytes
  - No heap allocation: created with (*new(EncoderType)) at compile time
  - No interface boxing: generic type parameters eliminate runtime type assertions
  - Compile-time specialization: each type combination gets its own optimized code path without runtime checking or reflection.
  - Freedom to deviate from upstream semantics: the conjure specification conflicts with default go behavior in a few situations.
    In these cases, cj's encoding and decoding logic strive to match the conjure-specified behavior with more control than reflection
    and interface boxing can provide.

When conjure-go generates code, it creates specific encoder/decoder chains like:

	cj.ListMarshaler[[]string, string, cj.String[string]]

This compiles to specialized code with no runtime overhead, unlike traditional interface{}-based JSON libraries
that require reflection and heap allocation for type discovery.

## Core Interfaces

TypeEncoder[T] provides JSON marshaling for type T:

	type TypeEncoder[T any] interface {
		MarshalJSONTo(enc *jsontext.Encoder, receiver T) error
	}

TypeDecoder[T] provides JSON unmarshaling for type T:

	type TypeDecoder[T any] interface {
		UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error
	}

MapKeyEncoder[K] extends TypeEncoder for types that can be map keys with custom ordering:

	type MapKeyEncoder[K comparable] interface {
		TypeEncoder[K]
		Compare(K, K) int  // For deterministic key ordering
	}

## Type System

The package provides encoders/decoders for all Conjure primitive and container types:

Primitive Types:
  - String[T ~string]: string-like types
  - Int32[T ~int|~int8|~int16|~int32|~int64]: signed integers
  - SafeLong[T ~int|~int8|~int16|~int32|~int64]: JavaScript-safe integers
  - Float[T ~float64]: floating-point numbers
  - Boolean[T ~bool]: boolean types
  - BearerToken[T ~string]: bearer tokens with validation
  - DateTime[T time.Time|datetime.DateTime]: RFC3339 timestamps
  - UUID[T ~[16]byte]: UUID types
  - RID[T ridConstraint]: resource identifiers
  - Binary: []byte handling
  - Any[T any]: fallback for arbitrary types

Map Key Variants:
  - Int32MapKey, SafeLongMapKey, FloatMapKey, BooleanMapKey, BinaryMapKey: primitive types as string keys
  - Text-like keys (Boolean, RID, UUID, TextMarshaler, StringerMarshaler, etc.) implement MapKeyEncoder for deterministic ordering.

Container Types:
  - ListMarshaler/ListUnmarshaler[T ~[]U, U, ITEM]: slices using ITEM encoder/decoder
  - OrderedMapMarshaler/MapUnmarshaler[T ~map[K]V, K, V, KEY, VAL]: maps with keys that support > < = operators
  - ComparableMapMarshaler: maps with custom key comparison
  - OptionalMarshaler/OptionalUnmarshaler[T *U, U, ITEM]: pointer types (nullable)

Types with interface constraints implemented by generated and library types:
  - StructMarshaler[T]: delegates to T's MarshalJSONTo method
  - StructUnmarshaler[T]: delegates to T's UnmarshalJSONFrom method
  - StringerMarshaler, TextMarshaler, TextUnmarshaler: for types with text conversion methods

# Usage in Generated Code

Conjure-go generates type-specific encoder/decoder chains. For example, a service method that accepts
[]string generates:

	httpclient.WithRequestBody(bodyArg,
		cj.ClientEncoder[[]string, cj.ListMarshaler[[]string, string, cj.String[string]]]{})

The type chain reads as: "ClientEncoder for []string using ListMarshaler for []string where each string
element uses String encoder". This provides:
  - Compile-time type checking: mismatched types cause build errors
  - Zero runtime overhead: no reflection or interface boxing
  - Memory efficiency: no intermediate allocations during encoding/decoding

For complex nested types like map[string]*CustomObject:

	cj.OrderedMapMarshaler[map[string]*CustomObject, string, *CustomObject,
		cj.String[string],
		cj.OptionalMarshaler[*CustomObject, CustomObject, cj.StructMarshaler[CustomObject]]]

## Integration with conjure-go-runtime

The package provides ClientEncoder and ClientDecoder types that implement conjure-go-runtime's
codecs.Encoder and codecs.Decoder interfaces, enabling seamless integration with HTTP clients:

	encoder := cj.ClientEncoder[MyType, cj.StructMarshaler[MyType]]{}
	decoder := cj.ClientDecoder[MyType, cj.StructUnmarshaler[*MyType]]{}

These adapt the cj type system to the runtime's interface-based codec system while maintaining
the performance benefits internally.

## Error Handling

The package provides a rich error taxonomy that preserves JSON parsing context:

Error Types:
- SyntaxError: malformed JSON syntax
- KindMismatchError: wrong JSON type (e.g., string when number expected)
- InvalidValueError: correct JSON type but invalid value (e.g., NaN as map key)
- UnmarshalFieldError: field-specific unmarshaling failure
- MissingFieldsError: required struct fields missing
- UnknownFieldsError: unknown struct fields (when RejectUnknownMembers enabled)
- DuplicateFieldKeyError/DuplicateMapKeyError: duplicate key detection

All errors include:
- Byte offset position in the JSON stream
- JSON Pointer path to the error location (RFC 6901)
- Stack trace for debugging
- Wrapped underlying cause errors

### Required Field Validation

Generated code tracks required struct fields using local `seen*` boolean variables during
unmarshaling. For each required field, the generated UnmarshalJSONFrom method maintains
a corresponding boolean flag:

	var seenName bool
	var seenAge bool
	// ... unmarshaling logic sets these to true when fields are encountered

	var missingFields []string
	if !seenName {
		missingFields = append(missingFields, "name")
	}
	if !seenAge {
		missingFields = append(missingFields, "age")
	}
	if len(missingFields) > 0 {
		return cj.NewMissingFieldsError(dec, "Person", missingFields)
	}

This pattern ensures that API contracts are enforced at the JSON layer, preventing partially
initialized structs from propagating through the system. The approach is more explicit and
debuggable than reflection-based validation performed after unmarshaling.

### Duplicate Key Detection

Both struct field and map key decoders detect and reject duplicate keys, which while legal
JSON, often indicates data corruption or API misuse. This provides early error detection
rather than silent last-wins behavior.

_Disclosure: This document is primarily the work of Claude. The code is not._


# Other Edge Cases and Specific Features

## Strict vs. Lenient Parsing

The package supports both strict server-side parsing and lenient client-side parsing through the
json.RejectUnknownMembers option. Generated server code uses:

	json.RejectUnknownMembers(true)

This ensures API compatibility - servers reject payloads with unknown fields, preventing silent
data loss during API evolution. Client decoders are lenient by default, ignoring unknown fields
to enable forward compatibility when servers add new optional fields.

## IEEE 754 Special Values

Float encoders handle IEEE 754 special cases correctly:
  - NaN: encoded as JSON string "NaN"
  - +Infinity: encoded as JSON string "Infinity"
  - -Infinity: encoded as JSON string "-Infinity"

However, NaN cannot be used as map keys since it violates mathematical ordering requirements,
resulting in a specific InvalidValueError during decoding.

## Nil Slice and Map Handling

Container types support configurable nil handling via json/v2 options to provide semantic
distinction between "no values" (null) and "empty collection" ([]/{}):
  - json.FormatNilSliceAsNull(true): nil slices encode as JSON null instead of empty arrays
  - json.FormatNilMapAsNull(true): nil maps encode as JSON null instead of empty objects

By default, nil slices and maps are encoded as empty arrays and objects respectively.

## Map Key Serialization and Ordering

JSON requires object keys to be strings, but Go maps can have non-string keys. The package
handles this through specialized MapKey encoders:
  - `cj.Int32MapKey`: converts integers to string keys
  - `cj.FloatMapKey`: converts floats to string keys (with special value handling)
  - `cj.TextMarshaler`: uses MarshalText for custom key types

Go's map iteration order is intentionally randomized, but deterministic JSON output is helpful for
testing, caching, reproducibility, etc. The package provides two distinct map marshalers based on
Go's type system capabilities:

**OrderedMapMarshaler** for `cmp.Ordered` types (strings and numbers):

	OrderedMapMarshaler[map[string]Value, string, Value, cj.String[string], ValueEncoder]

Uses Go 1.21+'s `cmp.Ordered` constraint and `slices.Sorted()` for built-in comparison operators.
This is the fast path for types that support `<`, `>`, `==` operators natively.

**ComparableMapMarshaler** for custom comparable types (binary, boolean, datetime, rid, uuid):

	ComparableMapMarshaler[map[CustomKey]Value, CustomKey, Value, CustomKeyEncoder, ValueEncoder]

Uses the `MapKeyEncoder.Compare(K, K) int` method for types that are comparable (`==`, `!=`) but
don't support ordering operators (`<`, `>`). Examples include:
  - UUID types: lexicographic comparison of string representation
  - Custom time types: chronological ordering
  - Composite keys: multi-field comparison logic

This design leverages Go's type system to choose the most efficient sorting algorithm at compile
time while supporting both standard types and domain-specific ordering requirements.

## Integer and SafeLong Precision Handling

The Int32 type encoder validates that integers fit within 32-bit signed integer range even when
the backing integer type is larger.

The SafeLong type encoder validates that integers fit within JavaScript's safe integer range
(-(2^53-1) to 2^53-1), preventing precision loss when JSON is consumed by JavaScript clients.
