package dj

import (
	"bytes"
	stdjson "encoding/json"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// Type is Result type
type Type uint8

const (
	// Null is a null json value
	Null Type = iota
	// False is a json false boolean
	False
	// True is a json true boolean
	True
	// Number is json number
	Number
	// String is a json string
	String
	// Object is a json key-value mapping.
	Object
	// Array is a json array.
	Array
)

// String returns a string representation of the type.
func (t Type) String() string {
	switch t {
	default:
		return ""
	case Null:
		return "Null"
	case False:
		return "False"
	case Number:
		return "Number"
	case String:
		return "String"
	case True:
		return "True"
	case Object:
		return "Object"
	case Array:
		return "Array"
	}
}

// Result represents a json value that is returned from the Parse functions,
// or returned as a child element of an outer collection Result. Its Type
// method can be used to determine the actual type of the value for deserialization.
type Result interface {

	// Type will be one of Null, False, True, Number, String, Object, or Array.
	Type() Type

	// Index returns the index of the raw value in the original json.
	// If the index is unknown, it will return 0.
	// Generally this should only be used for debugging and error reporting purposes.
	Index() int

	// IsNull returns true if the value is a json null.
	IsNull() bool
	// String returns a decoded string representation of the value.
	// If the value is not a json string, it will return an error.
	String() (string, error)
	// Text returns a decoded string as a byte slice, suitable for encoding.TextUnmarshaler.
	Text() ([]byte, error)
	// Bool returns a boolean representation.
	// If the value is not a json boolean, it will return an error.
	Bool() (bool, error)
	// Int returns an integer representation.
	// If the value is not a json number, it will return an error.
	Int() (int64, error)
	// Float returns a float64 representation.
	// If the value is not a json number, it will return an error.
	// Recognized exceptions are "NaN", "Infinity", "+Inf", "-Infinity", or "-Inf",
	// which will return the corresponding float64 value.
	Float() (float64, error)

	// NextObjectEntry returns the next key and value in an object.
	// If the value is not a json object, it will return an error.
	// The i param is the index of the last value returned by NextObjectEntry, or 0 to start.
	// The key always has a type of String.
	// If there are no more entries, it will return ok=false.
	NextObjectEntry(i int) (key Result, value Result, iOut int, ok bool, err error)
	// VisitObject iterates through key-value pairs in an object.
	// If the value is not a json object, it will return an error.
	// The iterator will be called for each key-value pair in the object.
	// If the iterator returns an error, the iteration will stop and return the error.
	//
	// VisitObject is a convenient wrapper around NextObjectEntry. Hot code paths
	// may avoid the overhead of the iterator function by using NextObjectEntry directly.
	VisitObject(iterator func(key Result, value Result) error) error

	// NextArrayEntry returns the next value in an array.
	// If the value is not a json array, it will return an error.
	// The i param is the index of the last value returned by NextArrayEntry, or 0 to start.
	// If there are no more entries, it will return ok=false.
	NextArrayEntry(i int) (value Result, iOut int, ok bool, err error)
	// VisitArray iterates through values in an array.
	// If the value is not a json array, it will return an error.
	// The iterator will be called for each value in the array.
	//
	// VisitArray is a convenient wrapper around NextArrayEntry. Hot code paths
	// may avoid the overhead of the iterator function by using NextArrayEntry directly.
	VisitArray(iterator func(value Result) error) error

	// Unmarshal will decode the value into a struct or pointer.
	Unmarshal(unmarshaler stdjson.Unmarshaler) error
	// Value returns the value as a native Go type.
	// The return value will be one of nil, bool, float64, string, map[string]any, or []any.
	Value() (any, error)
}

// Result represents a json value that is returned from Get().
type ResultImpl[DATA string | []byte] struct {
	// Type is the json type
	typ Type
	// Index of raw value in original json, zero means index unknown
	index int
	// Raw is the raw json
	Raw string
}

func (t ResultImpl[DATA]) Type() Type {
	return t.typ
}

func (t ResultImpl[DATA]) Index() int {
	return t.index
}

// IsNull returns true if the value is a json null.
func (t ResultImpl[DATA]) IsNull() bool {
	return t.typ == Null
}

// String returns a string representation of the value.
func (t ResultImpl[DATA]) String() (string, error) {
	if t.typ != String {
		return "", NewTypeMismatchError(t.index, t.typ, String.String())
	}
	if len(t.Raw) < 2 || t.Raw[0] != '"' || t.Raw[len(t.Raw)-1] != '"' {
		return "", NewSyntaxError(t.index, "invalid string")
	}
	// unescape on first \\ found
	for i := 1; i < len(t.Raw); i++ {
		if t.Raw[i] == '\\' {
			sb := new(strings.Builder)
			// trim quotes
			if err := unescape(t.Raw[1:len(t.Raw)-1], sb); err != nil {
				return "", err
			}
			return sb.String(), nil
		}
	}
	// trim quotes
	return string(t.Raw[1 : len(t.Raw)-1]), nil
}

func (t ResultImpl[DATA]) Text() ([]byte, error) {
	if t.typ != String {
		return nil, NewTypeMismatchError(t.index, t.typ, String.String())
	}
	if len(t.Raw) < 2 || t.Raw[0] != '"' || t.Raw[len(t.Raw)-1] != '"' {
		return nil, NewSyntaxError(t.index, "invalid string")
	}
	// unescape on first \\ found
	for i := 1; i < len(t.Raw); i++ {
		if t.Raw[i] == '\\' {
			bb := new(bytes.Buffer)
			// trim quotes
			if err := unescape(t.Raw[1:len(t.Raw)-1], bb); err != nil {
				return nil, err
			}
			return bb.Bytes(), nil
		}
	}
	// trim quotes
	return []byte(t.Raw[1 : len(t.Raw)-1]), nil
}

// Bool returns a boolean representation.
func (t ResultImpl[DATA]) Bool() (bool, error) {
	switch t.typ {
	default:
		return false, NewTypeMismatchError(t.index, t.typ, "boolean")
	case True:
		return true, nil
	case False:
		return false, nil
	}
}

// Int returns an integer representation.
func (t ResultImpl[DATA]) Int() (int64, error) {
	if t.typ != Number {
		return 0, NewTypeMismatchError(t.index, t.typ, Number.String())
	}
	// now try to parse the raw string
	n, err := strconv.ParseInt(string(t.Raw), 10, 64)
	if err != nil {
		return 0, NewInvalidValueError(t.index, "invalid integer", err)
	}
	return n, nil
}

// Float returns an float64 representation.
func (t ResultImpl[DATA]) Float() (float64, error) {
	s := string(t.Raw)
	switch s {
	case `"NaN"`:
		return math.NaN(), nil
	case `"Inf"`, `"Infinity"`:
		return math.Inf(1), nil
	case `"-Inf","-Infinity"`:
		return math.Inf(-1), nil
	}
	if t.typ != Number {
		return 0, NewTypeMismatchError(t.index, t.typ, Number.String())
	}
	return strconv.ParseFloat(s, 64)
}

func (t ResultImpl[DATA]) Unmarshal(unmarshaler stdjson.Unmarshaler) error {
	return unmarshaler.UnmarshalJSON([]byte(t.Raw))
}

func (t ResultImpl[DATA]) NextObjectEntry(i int) (key, value Result, iOut int, ok bool, err error) {
	if i >= len(t.Raw) {
		return nil, nil, 0, false, NewSyntaxError(i, "object index out of bounds")
	}
	if t.typ != Object {
		return nil, nil, 0, false, NewTypeMismatchError(t.index, t.typ, Object.String())
	}
	json := t.Raw

	i = validSpace(json, i)
	switch json[i] {
	case '}':
		return nil, nil, i + 1, false, nil
	case ',':
		i++
	case '{':
		i++
		i = validSpace(json, i)
		if json[i] == '}' {
			return nil, nil, i + 1, false, nil
		}
	default:
		return nil, nil, 0, false, NewSyntaxError(i, "invalid character preceding object entry")
	}

	i, key, err = validPayload(json, i)
	if err != nil {
		return nil, nil, 0, false, err
	}
	if key.Type() != String {
		return nil, nil, 0, false, NewTypeMismatchError(i, t.typ, String.String())
	}
	i, err = validColon(json, i)
	if err != nil {
		return nil, nil, 0, false, err
	}
	i, value, err = validPayload(json, i)
	if err != nil {
		return nil, nil, 0, false, err
	}
	return key, value, i, true, nil
}

func (t ResultImpl[DATA]) VisitObject(iterator func(key, value Result) error) error {
	var i int
	for {
		var key, value Result
		var ok bool
		var err error
		key, value, i, ok, err = t.NextObjectEntry(i)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if err := iterator(key, value); err != nil {
			return err
		}
	}
}

func (t ResultImpl[DATA]) NextArrayEntry(i int) (value Result, iOut int, ok bool, err error) {
	if i >= len(t.Raw) {
		return nil, 0, false, NewSyntaxError(i, "array index out of bounds")
	}
	if t.typ != Array {
		return nil, 0, false, NewTypeMismatchError(t.index, t.typ, Array.String())
	}
	json := t.Raw
	i = validSpace(json, i)
	switch json[i] {
	case ']':
		return nil, i + 1, false, nil
	case ',':
		i++
	case '[':
		i++
		i = validSpace(json, i)
		if json[i] == ']' {
			return nil, i + 1, false, nil
		}
	default:
		return nil, 0, false, NewSyntaxError(i, "invalid character preceding array entry")
	}
	i, value, err = validPayload(json, i)
	if err != nil {
		return nil, 0, false, err
	}
	return value, i, true, nil
}

func (t ResultImpl[DATA]) VisitArray(iterator func(value Result) error) error {
	var i int
	for {
		var value Result
		var ok bool
		var err error
		value, i, ok, err = t.NextArrayEntry(i)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if err := iterator(value); err != nil {
			return err
		}
	}
}

//type ObjectIterator[DATA string | []byte] struct{}
//
//// HasNext returns true if there are more values to iterate.
//// The i param is the index of the last value returned by Next().
//func (ObjectIterator[DATA]) HasNext(t Result[DATA], i int) bool {
//	json := t.Raw
//	if i == 0 {
//		i++ // skip the first '{'
//	}
//	for ; i < len(json); i++ {
//		if json[i] <= ' ' || json[i] == ',' || json[i] == ':' {
//			continue
//		}
//		return json[i] != '}'
//	}
//	return false
//}
//
//// Next returns the next key, value, and index to pass to HasNext().
//func (ObjectIterator[DATA]) Next(t Result[DATA], i int) (key Result[DATA], value Result[DATA], iOut int, err error) {
//	json := t.Raw
//	if i == 0 {
//		i++ // skip the first '{'
//	}
//	for ; i < len(json); i++ {
//		if json[i] != '"' {
//			continue
//		}
//		keyOffset := i
//		i, key, err = ParseNext(json, i)
//		if err != nil {
//			return Result[DATA]{}, Result[DATA]{}, 0, err
//		}
//		if key.Type != String {
//			return Result[DATA]{}, Result[DATA]{}, 0, NewSyntaxError(key.Index, "expected string key")
//		}
//		key.Index = keyOffset + t.Index
//		for ; i < len(json); i++ {
//			if json[i] <= ' ' || json[i] == ',' || json[i] == ':' {
//				continue
//			}
//			break
//		}
//		i, value, err = ParseNext(json, i)
//		if err != nil {
//			return Result[DATA]{}, Result[DATA]{}, 0, err
//		}
//
//		return key, value, i, nil
//	}
//	return Result[DATA]{}, Result[DATA]{}, 0, NewSyntaxError(i, "expected object entry")
//}
//
//func (t Result[DATA]) ArrayIterator(i int) (ArrayIterator[DATA], int, error) {
//	if t.Type != Array {
//		return ArrayIterator[DATA]{}, 0, NewTypeMismatchError(t.Index, t.Type, Array.String())
//	}
//	return ArrayIterator[DATA]{}, i + 1, nil
//}
//
//type ArrayIterator[DATA string | []byte] struct{}
//
//// HasNext returns true if there are more values to iterate.
//// The i param is the index of the last value returned by Next().
//func (ArrayIterator[DATA]) HasNext(t Result[DATA], i int) bool {
//	json := t.Raw
//	for ; i < len(json); i++ {
//		if json[i] <= ' ' || json[i] == ',' || json[i] == ':' {
//			continue
//		}
//		return json[i] != ']'
//	}
//	return false
//}
//
//// Next returns the next value and index to pass to HasNext().
//func (ArrayIterator[DATA]) Next(t Result[DATA], i int) (value Result[DATA], iOut int, err error) {
//	json := t.Raw
//	for ; i < len(json); i++ {
//		if json[i] <= ' ' || json[i] == ',' {
//			continue
//		}
//		i, value, err = ParseNext(json, i)
//		if err != nil {
//			return Result[DATA]{}, i, err
//		}
//		return value, i, nil
//	}
//	return Result[DATA]{}, i, NewSyntaxError(i, "expected array element")
//}

// Value returns one of these types:
//
//	bool, for JSON booleans
//	float64, for JSON numbers
//	string, for JSON string literals
//	nil, for JSON null
//	map[string]any, for JSON objects
//	[]any, for JSON arrays
func (t ResultImpl[DATA]) Value() (any, error) {
	switch t.typ {
	default:
		return nil, NewTypeMismatchError(t.index, t.typ, "any")
	case Null:
		return nil, nil
	case False:
		return false, nil
	case True:
		return true, nil
	case Number:
		return t.Float()
	case String:
		return t.String()
	case Object:
		mapValue := make(map[string]any)
		var i int
		for {
			var key, value Result
			var ok bool
			var err error
			key, value, i, ok, err = t.NextObjectEntry(i)
			if err != nil {
				return nil, err
			}
			if !ok {
				return mapValue, nil
			}
			k, err := key.String()
			if err != nil {
				return nil, err
			}
			v, err := value.Value()
			if err != nil {
				return nil, err
			}
			mapValue[k] = v
		}
	case Array:
		arrayValue := make([]any, 0)
		var i int
		for {
			var value Result
			var ok bool
			var err error
			value, i, ok, err = t.NextArrayEntry(i)
			if !ok {
				return arrayValue, nil
			}
			if err != nil {
				return nil, err
			}
			v, err := value.Value()
			if err != nil {
				return nil, err
			}
			arrayValue = append(arrayValue, v)
		}
	}
}

// runeit returns the rune from the the \uXXXX
func runeit[DATA string | []byte](json DATA) rune {
	n, _ := strconv.ParseUint(string(json[:4]), 16, 64)
	return rune(n)
}

// unescape unescapes a string
func unescape[DATA string | []byte, OUT interface {
	io.Writer
	io.ByteWriter
	Grow(int)
}](json DATA, out OUT) error {
	out.Grow(len(json))
	ln := len(json)
	for i := 0; i < ln; i++ {
		switch {
		default:
			err := out.WriteByte(json[i])
			if err != nil {
				return err
			}
		case json[i] < ' ':
			return NewSyntaxError(i, "invalid character for encoded string")
		case json[i] == '\\':
			i++
			if i >= len(json) {
				return NewSyntaxError(i, "incomplete escape sequence in encoded string")
			}
			switch json[i] {
			default:
				return NewSyntaxError(i, "invalid escape sequence in encoded string")
			case '\\':
				if err := out.WriteByte('\\'); err != nil {
					return err
				}
			case '/':
				if err := out.WriteByte('/'); err != nil {
					return err
				}
			case 'b':
				if err := out.WriteByte('\b'); err != nil {
					return err
				}
			case 'f':
				if err := out.WriteByte('\f'); err != nil {
					return err
				}
			case 'n':
				if err := out.WriteByte('\n'); err != nil {
					return err
				}
			case 'r':
				if err := out.WriteByte('\r'); err != nil {
					return err
				}
			case 't':
				if err := out.WriteByte('\t'); err != nil {
					return err
				}
			case '"':
				if err := out.WriteByte('"'); err != nil {
					return err
				}
			case 'u':
				if i+5 > len(json) {
					return NewSyntaxError(i, "incomplete unicode sequence in encoded string")
				}
				r := runeit(json[i+1:])
				i += 5
				if utf16.IsSurrogate(r) {
					// need another code
					if len(json[i:]) >= 6 && json[i] == '\\' &&
						json[i+1] == 'u' {
						// we expect it to be correct so just consume it
						r = utf16.DecodeRune(r, runeit(json[i+2:]))
						i += 6
					}
				}
				// provide enough space to encode the largest utf8 possible
				buf := make([]byte, 8)
				n := utf8.EncodeRune(buf, r)
				if _, err := out.Write(buf[:n]); err != nil {
					return err
				}
				i-- // backtrack index by one
			}
		}
	}
	return nil
}
