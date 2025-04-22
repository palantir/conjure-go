package cj

import (
	"cmp"
	"encoding"
	"encoding/base64"
	"fmt"
	"maps"
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/uuid"
)

type TypeArshaler[T any] interface {
	TypeEncoder[T]
	TypeDecoder[T]
}

type TypeEncoder[T any] interface {
	MarshalJSONTo(receiver T, enc *jsontext.Encoder) error
}

type TypeDecoder[T any] interface {
	UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error
}

type MapKeyArshaler[K comparable] interface {
	TypeArshaler[K]
	MapKeyComparator[K]
}

type MapKeyComparator[K comparable] interface {
	// Compare returns -1 if a < b, 0 if a == b, and 1 if a > b.
	// This is used to sort keys in a deterministic order.
	Compare(K, K) int
}

type TypeAny[T any] struct{}

func (TypeAny[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, receiver)
}

func (TypeAny[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	return json.UnmarshalDecode(dec, receiver)
}

type TypeBinary[T ~[]byte] struct{}

func (TypeBinary[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	b64out := enc.UnusedBuffer()
	b64out = append(b64out, '"')
	b64out = base64.StdEncoding.AppendEncode(b64out, receiver)
	b64out = append(b64out, '"')
	return enc.WriteValue(b64out)
}

func (TypeBinary[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	val, err := readValueOfKind(dec, '"')
	if err != nil {
		return err
	}
	if val[len(val)-1] != '"' {
		return NewSyntaxError(dec, "expected closing quote", nil)
	}
	*receiver, err = base64.StdEncoding.AppendDecode(*receiver, val[1:len(val)-1])
	if err != nil {
		return NewSyntaxError(dec, "invalid base64", err)
	}
	return nil
}

type TypeBoolean[T ~bool] struct{}

func (TypeBoolean[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	if receiver {
		return enc.WriteToken(jsontext.True)
	}
	return enc.WriteToken(jsontext.False)
}

func (TypeBoolean[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, 't', 'f')
	if err != nil {
		return err
	}
	if tok.Kind() == 't' {
		*receiver = true
	} else {
		*receiver = false
	}
	return nil
}

type TypeBooleanMapKey[T ~bool] struct{}

func (TypeBooleanMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	if receiver {
		return enc.WriteValue(jsontext.Value(`"true"`))
	}
	return enc.WriteValue(jsontext.Value(`"false"`))
}

func (TypeBooleanMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	b, err := strconv.ParseBool(tok.String())
	if err != nil {
		return NewSyntaxError(dec, "invalid boolean", err)
	}
	*receiver = T(b)
	return nil
}

type TypeDateTime[T time.Time | datetime.DateTime] struct{}

func (TypeDateTime[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	out = append(out, '"')
	out = time.Time(receiver).AppendFormat(out, time.RFC3339Nano)
	out = append(out, '"')
	return enc.WriteValue(out)
}

func (TypeDateTime[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	parse, err := time.Parse(tok.String(), time.RFC3339Nano)
	if err != nil {
		return NewSyntaxError(dec, "invalid datetime", err)
	}
	*receiver = T(parse)
	return nil
}

type TypeFloat[T ~float64] struct{}

func (TypeFloat[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	out = append(out, '"')
	n := float64(receiver)
	switch {
	case math.Float64bits(n) == 0:
		out = append(out, '0')
	case math.IsNaN(n):
		out = append(out, "\"NaN\""...)
	case math.IsInf(n, +1):
		out = append(out, "\"Infinity\""...)
	case math.IsInf(n, -1):
		out = append(out, "\"-Infinity\""...)
	default:
		out = strconv.AppendFloat(out, n, 'f', -1, 64)
	}
	return enc.WriteValue(out)
}

func (TypeFloat[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '0')
	if err != nil {
		return err
	}
	*receiver = T(tok.Float())
	return nil
}

type TypeFloatMapKey[T ~float64] struct{}

func (TypeFloatMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	n := float64(receiver)
	switch {
	case math.Float64bits(n) == 0:
		out = append(out, "\"0\""...)
	case math.IsNaN(n):
		out = append(out, "\"NaN\""...)
	case math.IsInf(n, +1):
		out = append(out, "\"Infinity\""...)
	case math.IsInf(n, -1):
		out = append(out, "\"-Infinity\""...)
	default:
		out = append(out, '"')
		out = strconv.AppendFloat(out, n, 'f', -1, 64)
		out = append(out, '"')
	}
	return enc.WriteValue(out)
}

func (TypeFloatMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	f, err := strconv.ParseFloat(tok.String(), 64)
	if err != nil {
		return NewSyntaxError(dec, "invalid float", err)
	}
	*receiver = T(f)
	return nil
}

type TypeInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64] struct{}

func (TypeInt[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	out = strconv.AppendInt(out, int64(receiver), 10)
	return enc.WriteValue(out)
}

func (TypeInt[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '0')
	if err != nil {
		return err
	}
	*receiver = T(tok.Int())
	return nil
}

type TypeIntMapKey[T ~int | ~int8 | ~int16 | ~int32 | ~int64] struct{}

func (TypeIntMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	out = append(out, '"')
	out = strconv.AppendInt(out, int64(receiver), 10)
	out = append(out, '"')
	return enc.WriteValue(out)
}

func (TypeIntMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	i, err := strconv.ParseInt(tok.String(), 10, 64)
	if err != nil {
		return NewSyntaxError(dec, "invalid int", err)
	}
	*receiver = T(i)
	return nil
}

type TypeUint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64] struct{}

func (TypeUint[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	out = strconv.AppendUint(out, uint64(receiver), 10)
	return enc.WriteValue(out)
}

func (TypeUint[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '0')
	if err != nil {
		return err
	}
	*receiver = T(tok.Uint())
	return nil
}

type TypeUintMapKey[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64] struct{}

func (TypeUintMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	out = append(out, '"')
	out = strconv.AppendUint(out, uint64(receiver), 10)
	out = append(out, '"')
	return enc.WriteValue(out)
}

func (TypeUintMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	i, err := strconv.ParseUint(tok.String(), 10, 64)
	if err != nil {
		return NewSyntaxError(dec, "invalid uint", err)
	}
	*receiver = T(i)
	return nil
}

type TypeRID[T rid.ResourceIdentifier] struct{}

func (TypeRID[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(rid.ResourceIdentifier(receiver).String()))
}

func (TypeRID[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	parsed, err := rid.ParseRID(tok.String())
	if err != nil {
		return NewSyntaxError(dec, "invalid resource identifier", err)
	}
	*receiver = T(parsed)
	return nil
}

type TypeString[T ~string] struct{}

func (TypeString[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(string(receiver)))
}

func (TypeString[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	*receiver = T(tok.String())
	return nil
}

type TypeUUID[T ~[16]byte] struct{}

func (TypeUUID[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(uuid.UUID(receiver).String()))
}

func (TypeUUID[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	parsed, err := uuid.ParseUUID(tok.String())
	if err != nil {
		return NewSyntaxError(dec, "invalid UUID", err)
	}
	*receiver = T(parsed)
	return nil
}

type TypeOptional[T any, ITEM TypeArshaler[T]] struct{}

func (TypeOptional[T, ITEM]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	if receiver == nil {
		return enc.WriteToken(jsontext.Null)
	}
	return (*new(ITEM)).MarshalJSONTo(receiver, enc)
}

func (TypeOptional[T, ITEM]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	if dec.PeekKind() == 'n' {
		// still have to consume token
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		*receiver = *new(T)
		return nil
	}
	return (*new(ITEM)).UnmarshalJSONFrom(receiver, dec)
}

type TypeList[T any, ITEM TypeArshaler[T]] struct{}

func (TypeList[T, ITEM]) MarshalJSONTo(receiver []T, enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginArray); err != nil {
		return err
	}
	for _, item := range receiver {
		err := (*new(ITEM)).MarshalJSONTo(item, enc)
		if err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndArray)
}

func (TypeList[T, ITEM]) UnmarshalJSONFrom(receiver *[]T, dec *jsontext.Decoder) error {
	if dec.PeekKind() == 'n' {
		// still have to consume token
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		*receiver = make([]T, 0)
		return nil
	}
	if _, err := readTokenOfKind(dec, '['); err != nil {
		return err
	}
	if *receiver == nil {
		*receiver = make([]T, 0)
	} else {
		*receiver = (*receiver)[:0]
	}
	for {
		item := *new(T)
		if err := (*new(ITEM)).UnmarshalJSONFrom(&item, dec); err != nil {
			return err
		}
		*receiver = append(*receiver, item)

		if dec.PeekKind() == ']' {
			break
		}
	}
	if _, err := readTokenOfKind(dec, ']'); err != nil {
		return err
	}
	return nil
}

type TypeSortedMap[K comparable, V any, KEY MapKeyArshaler[K], VAL TypeArshaler[V]] struct{}

func (TypeSortedMap[K, V, KEY, VAL]) MarshalJSONTo(receiver map[K]V, enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	keys := make([]K, 0, len(receiver))
	for k := range receiver {
		keys = append(keys, k)
	}

	for _, k := range slices.SortedFunc(maps.Keys(receiver), (*new(KEY)).Compare) {
		if err := (*new(KEY)).MarshalJSONTo(k, enc); err != nil {
			return err
		}
		if err := (*new(VAL)).MarshalJSONTo(receiver[k], enc); err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndObject)
}

type TypeOrderedMap[K cmp.Ordered, V any, KEY TypeArshaler[K], VAL TypeArshaler[V]] struct{}

func (TypeOrderedMap[K, V, KEY, VAL]) MarshalJSONTo(receiver map[K]V, enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	for _, k := range slices.Sorted(maps.Keys(receiver)) {
		err := (*new(KEY)).MarshalJSONTo(k, enc)
		if err != nil {
			return err
		}
		err = (*new(VAL)).MarshalJSONTo(receiver[k], enc)
		if err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndObject)
}

type stringerTokenMarshaler[T fmt.Stringer] struct{}

func (stringerTokenMarshaler[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(receiver.String()))
}

// use only when string is known to be ASCII only (for example, hex or base64 encoded).
// this is a performance optimization that skips the UTF-8 validation step.
type stringerValueMarshaler[T fmt.Stringer] struct{}

func (stringerValueMarshaler[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	out = append(out, '"')
	out = append(out, receiver.String()...)
	out = append(out, '"')
	return enc.WriteValue(out)
}

type TextArshaler[TM encoding.TextMarshaler, TU encoding.TextUnmarshaler] struct{}

func (TextArshaler[TM, TU]) MarshalJSONTo(receiver TM, enc *jsontext.Encoder) error {
	text, err := receiver.MarshalText()
	if err != nil {
		return err
	}
	out := enc.UnusedBuffer()
	out, err = jsontext.AppendQuote(out, text)
	if err != nil {
		return err
	}
	return enc.WriteValue(out)
}

func (TextArshaler[TM, TU]) UnmarshalJSONFrom(receiver *TM, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	unmarshaler := any(receiver).(TU)
	return unmarshaler.UnmarshalText([]byte(tok.String()))
}
