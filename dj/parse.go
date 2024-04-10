// Copyright (c) 2023 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// ResultAPI represents a json value that is returned from the Parse functions,
// or returned as a child element of an outer collection Result. Its Type
// method can be used to determine the actual type of the value for deserialization.
type ResultAPI interface {

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
type Result struct {
	// Type is the json type, one of Null, False, True, Number, String, Object, or Array.
	Type Type
	// Index returns the index of the raw value in the original json.
	// If the index is unknown, it will return 0.
	// Generally this should only be used for debugging and error reporting purposes.
	Index int
	// Raw is the raw json
	Raw string
}

// String returns a string representation of the value.
// If the value is not a json string, it will return an error.
func (t Result) String() (string, error) {
	if t.Type != String {
		return "", NewTypeMismatchError(t, String.String())
	}
	if len(t.Raw) < 2 || t.Raw[0] != '"' || t.Raw[len(t.Raw)-1] != '"' {
		return "", NewSyntaxError(t.Index, "invalid string")
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
	return t.Raw[1 : len(t.Raw)-1], nil
}

// Text returns a decoded string as a byte slice, suitable for encoding.TextUnmarshaler.
// If the value is not a json string, it will return an error.
func (t Result) Text() ([]byte, error) {
	if t.Type != String {
		return nil, NewTypeMismatchError(t, String.String())
	}
	if len(t.Raw) < 2 || t.Raw[0] != '"' || t.Raw[len(t.Raw)-1] != '"' {
		return nil, NewSyntaxError(t.Index, "invalid string")
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
// If the value is not a json boolean, it will return an error.
func (t Result) Bool() (bool, error) {
	switch t.Type {
	default:
		return false, NewTypeMismatchError(t, "boolean")
	case True:
		return true, nil
	case False:
		return false, nil
	}
}

// Int returns an integer representation.
// If the value is not a json number, it will return an error.
func (t Result) Int() (int64, error) {
	if t.Type != Number {
		return 0, NewTypeMismatchError(t, Number.String())
	}
	// now try to parse the raw string
	n, err := strconv.ParseInt(t.Raw, 10, 64)
	if err != nil {
		return 0, NewInvalidValueError(t, "invalid integer", err)
	}
	return n, nil
}

// Float returns a float64 representation.
// If the value is not a json number, it will return an error.
// Recognized exceptions are "NaN", "Infinity", "+Inf", "-Infinity", or "-Inf",
// which will return the corresponding float64 value.
func (t Result) Float() (float64, error) {
	switch t.Raw {
	case `"NaN"`:
		return math.NaN(), nil
	case `"Inf"`, `"Infinity"`:
		return math.Inf(1), nil
	case `"-Inf","-Infinity"`:
		return math.Inf(-1), nil
	}
	if t.Type != Number {
		return 0, NewTypeMismatchError(t, Number.String())
	}
	return strconv.ParseFloat(t.Raw, 64)
}

// Unmarshal will decode the value into a struct or pointer.
func (t Result) Unmarshal(unmarshaler stdjson.Unmarshaler) error {
	return unmarshaler.UnmarshalJSON([]byte(t.Raw))
}

//// NextObjectEntry returns the next key and value in an object.
//// If the value is not a json object, it will return an error.
//// The i param is the index of the last value returned by NextObjectEntry, or 0 to start.
//// The key always has a type of String.
//// If there are no more entries, it will return ok=false.
//func (t Result) NextObjectEntry(i int) (key, value Result, iOut int, ok bool, err error) {
//	if i >= len(t.Raw) {
//		return Result{}, Result{}, 0, false, NewSyntaxError(i, "object index out of bounds")
//	}
//	if t.Type != Object {
//		return Result{}, Result{}, 0, false, NewTypeMismatchError(t, Object.String())
//	}
//	json := t.Raw
//
//	i = validSpace(json, i)
//	switch json[i] {
//	case '}':
//		return Result{}, Result{}, i + 1, false, nil
//	case ',':
//		i++
//	case '{':
//		i++
//		i = validSpace(json, i)
//		if json[i] == '}' {
//			return Result{}, Result{}, i + 1, false, nil
//		}
//	default:
//		return Result{}, Result{}, 0, false, NewSyntaxError(i, "invalid character preceding object entry")
//	}
//
//	i, key, err = validPayload(json, i)
//	if err != nil {
//		return Result{}, Result{}, 0, false, err
//	}
//	if key.Type != String {
//		return Result{}, Result{}, 0, false, NewTypeMismatchError(t, String.String())
//	}
//	i, err = validColon(json, i)
//	if err != nil {
//		return Result{}, Result{}, 0, false, err
//	}
//	i, value, err = validPayload(json, i)
//	if err != nil {
//		return Result{}, Result{}, 0, false, err
//	}
//	return key, value, i, true, nil
//}

func (t Result) ObjectIterator(i int) (ObjectIterator, int, error) {
	if t.Type != Object {
		return ObjectIterator{}, 0, NewTypeMismatchError(t, Object.String())
	}
	return ObjectIterator{}, i + 1, nil
}

// VisitObject iterates through key-value pairs in an object.
// If the value is not a json object, it will return an error.
// The iterator will be called for each key-value pair in the object.
// If the iterator returns an error, the iteration will stop and return the error.
//
// VisitObject is a convenient wrapper around NextObjectEntry. Hot code paths
// may avoid the overhead of the iterator function by using NextObjectEntry directly.
func (t Result) VisitObject(iterator func(key, value Result) error) error {
	iter, i, err := t.ObjectIterator(0)
	if err != nil {
		return err
	}
	for iter.HasNext(t, i) {
		var key, value Result
		var err error
		key, value, i, err = iter.Next(t, i)
		if err != nil {
			return err
		}
		if err := iterator(key, value); err != nil {
			return err
		}
	}
	return nil
}

type ObjectIterator struct{}

// HasNext returns true if there are more values to iterate.
// The i param is the index of the last value returned by Next().
func (ObjectIterator) HasNext(t Result, i int) bool {
	if i == 0 {
		i++ // skip the first '{'
	}
	for ; i < len(t.Raw); i++ {
		ji := t.Raw[i]
		if ji <= ' ' || ji == ',' || ji == ':' {
			continue
		}
		return ji != '}'
	}
	return false
}

// Next returns the next key, value, and index to pass to HasNext().
func (ObjectIterator) Next(t Result, i int) (key Result, value Result, iOut int, err error) {
	json := t.Raw
	if i == 0 {
		i++ // skip the first '{'
	}
	for ; i < len(json); i++ {
		if json[i] != '"' {
			continue
		}
		keyOffset := i
		var str string
		var esc bool
		i, str, esc, err = parseString(json, i+1)
		if err != nil {
			return Result{}, Result{}, 0, err
		}
		key.Type = String
		if esc {
			sb := new(strings.Builder)
			if err := unescape(str, sb); err != nil {
				return Result{}, Result{}, 0, err
			}
			key.Raw = sb.String()
		} else {
			key.Raw = str
		}
		key.Index = keyOffset + t.Index
		for ; i < len(json); i++ {
			if json[i] <= ' ' || json[i] == ',' || json[i] == ':' {
				continue
			}
			break
		}
		valOffset := i
		i, value, err = parseAny(json, i)
		if err != nil {
			return Result{}, Result{}, 0, err
		}
		value.Index = valOffset + t.Index

		return key, value, i, nil
	}
	return Result{}, Result{}, 0, NewSyntaxError(i, "expected object entry")
}

func (t Result) ArrayIterator(i int) (ArrayIterator, int, error) {
	if t.Type != Array {
		return ArrayIterator{}, 0, NewTypeMismatchError(t, Array.String())
	}
	return ArrayIterator{}, i + 1, nil
}

//// NextArrayEntry returns the next value in an array.
//// If the value is not a json array, it will return an error.
//// The i param is the index of the last value returned by NextArrayEntry, or 0 to start.
//// If there are no more entries, it will return ok=false.
//func (t Result) NextArrayEntry(i int) (value Result, iOut int, ok bool, err error) {
//	if i >= len(t.Raw) {
//		return Result{}, 0, false, NewSyntaxError(i, "array index out of bounds")
//	}
//	if t.Type != Array {
//		return Result{}, 0, false, NewTypeMismatchError(t, Array.String())
//	}
//	json := t.Raw
//	i = validSpace(json, i)
//	switch json[i] {
//	case ']':
//		return Result{}, i + 1, false, nil
//	case ',':
//		i++
//	case '[':
//		i++
//		i = validSpace(json, i)
//		if json[i] == ']' {
//			return Result{}, i + 1, false, nil
//		}
//	default:
//		return Result{}, 0, false, NewSyntaxError(i, "invalid character preceding array entry")
//	}
//	i, value, err = validPayload(json, i)
//	if err != nil {
//		return Result{}, 0, false, err
//	}
//	return value, i, true, nil
//}

// VisitArray iterates through values in an array.
// If the value is not a json array, it will return an error.
// The iterator will be called for each value in the array.
//
// VisitArray is a convenient wrapper around usage of ArrayIterator. Hot code paths
// may avoid the overhead of the iterator function by using ArrayIterator directly.
func (t Result) VisitArray(iterator func(value Result) error) error {
	iter, i, err := t.ArrayIterator(0)
	if err != nil {
		return err
	}
	for iter.HasNext(t, i) {
		var value Result
		var err error
		value, i, err = iter.Next(t, i)
		if err != nil {
			return err
		}
		if err := iterator(value); err != nil {
			return err
		}
	}
	return nil
}

type ArrayIterator struct{}

// HasNext returns true if there are more values to iterate.
// The i param is the index of the last value returned by Next().
func (ArrayIterator) HasNext(t Result, i int) bool {
	json := t.Raw
	for ; i < len(json); i++ {
		if json[i] <= ' ' || json[i] == ',' || json[i] == ':' {
			continue
		}
		return json[i] != ']'
	}
	return false
}

// Next returns the next value and index to pass to HasNext().
func (ArrayIterator) Next(t Result, i int) (value Result, iOut int, err error) {
	json := t.Raw
	for ; i < len(json); i++ {
		if json[i] <= ' ' || json[i] == ',' {
			continue
		}
		valOffset := i
		i, value, err = parseAny(json, i)
		if err != nil {
			return Result{}, i, err
		}
		value.Index = valOffset + t.Index
		return value, i, nil
	}
	return Result{}, i, NewSyntaxError(i, "expected array element")
}

// Value returns one of these types:
//
//	bool, for JSON booleans
//	float64, for JSON numbers
//	string, for JSON string literals
//	nil, for JSON null
//	map[string]any, for JSON objects
//	[]any, for JSON arrays
func (t Result) Value() (any, error) {
	switch t.Type {
	default:
		return nil, NewTypeMismatchError(t, "any")
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
		iter, i, err := t.ObjectIterator(0)
		if err != nil {
			return nil, err
		}
		mapValue := make(map[string]any)
		for iter.HasNext(t, i) {
			var key, value Result
			var err error
			key, value, i, err = iter.Next(t, i)
			if err != nil {
				return nil, err
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
		return mapValue, nil
	case Array:
		iter, i, err := t.ArrayIterator(0)
		if err != nil {
			return nil, err
		}
		arrayValue := make([]any, 0)
		for iter.HasNext(t, i) {
			var value Result
			var err error
			value, i, err = iter.Next(t, i)
			if err != nil {
				return nil, err
			}
			v, err := value.Value()
			if err != nil {
				return nil, err
			}
			arrayValue = append(arrayValue, v)
		}
		return arrayValue, nil
	}
}

// Valid returns true if the input is valid json.
// The input can be a string or []byte.
func Valid[DATA string | []byte](data DATA) error {
	// Use validAny directly to avoid allocating the Result.Raw field.
	i, _, err := validAny(data, 0)
	if err != nil {
		return err
	}
	i = validSpace(data, i)
	if i < len(data) {
		return NewSyntaxError(i, "invalid character after JSON")
	}
	return nil
}

// Parse parses the json and returns a result. It returns a SyntaxError if the data is invalid JSON.
// The input can be a string or []byte.
// The returned Result's Index field is the starting index of the value.
// To parse multiple JSON values in a single string, use ParseNext.
func Parse[DATA string | []byte](data DATA) (Result, error) {
	i, res, err := ParseNext(data, 0)
	if err != nil {
		return res, err
	}
	if i < len(data) {
		return res, NewSyntaxError(i, "invalid character after JSON")
	}
	return res, nil
}

// ParseNext parses the next value from a json string.
// This function is useful when you have multiple json values in a single string.
// The return values are (i int, res Result, err error)
// The i is the index of the next character after the parsed value.
// When the returned i is equal to len(data), there are no more values to parse and the next call will error.
func ParseNext[DATA string | []byte](data DATA, i int) (int, Result, error) {
	return validPayload(data, i)
}

func parseString(json string, i int) (iOut int, val string, vesc bool, err error) {
	var s = i
	ln := len(json)
	for ; i < ln; i++ {
		ji := json[i]
		if ji > '\\' {
			continue
		}
		if ji == '"' {
			return i + 1, json[s-1 : i+1], false, nil
		}
		if ji == '\\' {
			i++
			for ; i < ln; i++ {
				ji := json[i]
				if ji > '\\' {
					continue
				}
				if ji == '"' {
					// look for an escaped slash
					if json[i-1] == '\\' {
						n := 0
						for j := i - 2; j > 0; j-- {
							if json[j] != '\\' {
								break
							}
							n++
						}
						if n%2 == 0 {
							continue
						}
					}
					return i + 1, json[s-1 : i+1], true, nil
				}
			}
			break
		}
	}
	return i, json[s-1:], false, NewSyntaxError(i, "invalid character for string")
}

// parse until any number-terminating token: space, comma, bracket, brace
func parseNumber(json string, i int) (int, string) {
	var s = i
	i++
	for ; i < len(json); i++ {
		if json[i] <= ' ' || json[i] == ',' || json[i] == ']' || json[i] == '}' {
			return i, json[s:i]
		}
	}
	return i, json[s:]
}

// parse unquoted values (true, false, null)
func parseLiteral(json string, i int) (int, string) {
	var s = i
	i++
	ln := len(json)
	for ; i < ln; i++ {
		if json[i] < 'a' || json[i] > 'z' {
			return i, json[s:i]
		}
	}
	return i, json[s:]
}

// returns the substring containing the json value up to the closing brace/bracket.
func parseSquash(json string, i int) (int, string) {
	// expects that the lead character is a '[' or '{' or '('
	// squash the value, ignoring all nested arrays and objects.
	// the first '[' or '{' or '(' has already been read
	s := i
	depth := 1
	ln := len(json)
	for {
		i++
		if i >= ln {
			break
		}
		ji := json[i]
		if ji >= '"' && ji <= '}' {
			switch ji {
			case '"':
				// parse string and skip escaped quotes
				s2 := i
				for {
					i++
					if i >= ln {
						break
					}
					ji := json[i]
					if ji > '\\' {
						continue
					}
					if ji == '"' {
						// look for an escaped slash
						if json[i-1] == '\\' {
							n := 0
							for j := i - 2; j > s2; j-- {
								if json[j] != '\\' {
									break
								}
								n++
							}
							if n%2 == 0 {
								continue
							}
						}
						break
					}
				}
			case '{', '[', '(':
				depth++
			case '}', ']', ')':
				depth--
				if depth == 0 {
					i++
					return i, json[s:i]
				}
			}
		}
	}
	return i, json[s:]
}

// parseAny parses the next value from a json string.
// A Result is returned when the hit param is set.
// The return values are (i int, res Result, err error)
func parseAny(json string, i int) (int, Result, error) {

	var res Result
	var val string
	for ; i < len(json); i++ {
		debugC := string([]byte{json[i]})
		debugS := json[i:]
		_, _ = debugC, debugS
		if json[i] <= ' ' {
			continue
		}
		var num bool
		switch json[i] {
		case '{':
			i, val = parseSquash(json, i)
			debugC1 := string([]byte{json[i]})
			debugS1 := json[i:]
			_, _ = debugC1, debugS1
			res.Raw = val
			res.Type = Object
			res.Index = i
			return i, res, nil
		case '[':
			i, val = parseSquash(json, i)
			debugC1 := string([]byte{json[i]})
			debugS1 := json[i:]
			_, _ = debugC1, debugS1
			res.Raw = val
			res.Type = Array
			res.Index = i
			return i, res, nil
		case '"':
			i, val, esc, err := parseString(json, i+1)
			if err != nil {
				return i, res, err
			}
			res.Type = String
			if esc {
				sb := new(strings.Builder)
				if err := unescape(val, sb); err != nil {
					return 0, Result{}, err
				}
				res.Raw = sb.String()
			} else {
				res.Raw = val
			}
			res.Index = i
			return i, res, nil
		case 'n':
			if i+1 < len(json) && json[i+1] != 'u' {
				num = true
				break
			}
			fallthrough
		case 't', 'f':
			vc := json[i]
			i, val = parseLiteral(json, i)
			res.Raw = val
			switch vc {
			case 't':
				res.Type = True
			case 'f':
				res.Type = False
			}
			res.Index = i
			return i, res, nil
		case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'i', 'I', 'N':
			num = true
		}
		if num {
			i, val = parseNumber(json, i)
			res.Raw = val
			res.Type = Number
			res.Index = i
			return i, res, nil
		}
	}
	return i, res, NewSyntaxError(i, "invalid character for json")
}

// valid* functions are slower than parse* equivalents because they check for invalid JSON.
// Once the JSON is validated, the parse* functions can be used to unmarshal the JSON faster.

func validPayload[DATA string | []byte](data DATA, i int) (outi int, res Result, err error) {
	i = validSpace(data, i)
	if i >= len(data) {
		return 0, Result{}, NewSyntaxError(i, "no content found")
	}
	res.Index = i
	i, res.Type, err = validAny(data, i)
	if err != nil {
		return 0, Result{}, err
	}
	res.Raw = string(data[res.Index:i])
	i = validSpace(data, i)
	return i, res, nil
}

func validAny[DATA string | []byte](data DATA, i int) (outi int, typ Type, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return 0, 0, NewSyntaxError(i, "invalid character beginning JSON")
		case ' ', '\t', '\n', '\r':
			continue
		case '{':
			i, err = validObject(data, i+1)
			return i, Object, err
		case '[':
			i, err = validArray(data, i+1)
			return i, Array, err
		case '"':
			i, err = validString(data, i+1)
			return i, String, err
		case 't':
			i, err = validTrue(data, i+1)
			return i, True, err
		case 'f':
			i, err = validFalse(data, i+1)
			return i, False, err
		case 'n':
			i, err = validNull(data, i+1)
			return i, Null, err
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			i, err = validNumber(data, i+1)
			return i, Number, err
		}
	}
	return 0, 0, NewSyntaxError(i, "empty content")
}

func validObject[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return i, NewSyntaxError(i, "expected object key or closing brace")
		case ' ', '\t', '\n', '\r':
			continue
		case '}':
			return i + 1, nil
		case '"':
		key:
			if i, err = validString(data, i+1); err != nil {
				return 0, err
			}
			if i, err = validColon(data, i); err != nil {
				return 0, err
			}
			if i, _, err = validAny(data, i); err != nil {
				return 0, err
			}
			if i, err = validComma(data, i, '}'); err != nil {
				return 0, err
			}
			if data[i] == '}' {
				return i + 1, nil
			}
			i++
			for ; i < len(data); i++ {
				switch data[i] {
				default:
					return i, NewSyntaxError(i, "invalid character between object entries")
				case ' ', '\t', '\n', '\r':
					continue
				case '"':
					goto key
				}
			}
			return i, NewSyntaxError(i, "object not closed after entry")
		}
	}
	return i, NewSyntaxError(i, "object not closed")
}

func validColon[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return i, NewSyntaxError(i, "invalid character for colon")
		case ' ', '\t', '\n', '\r':
			continue
		case ':':
			return i + 1, nil
		}
	}
	return i, NewSyntaxError(i, "expected colon")
}

func validComma[DATA string | []byte](data DATA, i int, end byte) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			return i, NewSyntaxError(i, "invalid character for comma")
		case ' ', '\t', '\n', '\r':
			continue
		case ',', end:
			return i, nil
		}
	}
	return i, NewSyntaxError(i, "expected comma")
}

func validArray[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		switch data[i] {
		default:
			for ; i < len(data); i++ {
				if i, _, err = validAny(data, i); err != nil {
					return 0, err
				}
				if i, err = validComma(data, i, ']'); err != nil {
					return 0, err
				}
				if data[i] == ']' {
					return i + 1, nil
				}
			}
		case ' ', '\t', '\n', '\r':
			continue
		case ']':
			return i + 1, nil
		}
	}
	return i, NewSyntaxError(i, "array not closed")
}

func validString[DATA string | []byte](data DATA, i int) (outi int, err error) {
	for ; i < len(data); i++ {
		if data[i] < ' ' {
			return i, NewSyntaxError(i, "invalid character for string")
		} else if data[i] == '\\' {
			i++
			if i == len(data) {
				return i, NewSyntaxError(i, "escape character at end of data")
			}
			switch data[i] {
			default:
				return i, NewSyntaxError(i, "invalid escape character "+string(data[i:i+1]))
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
			case 'u':
				for j := 0; j < 4; j++ {
					i++
					if i >= len(data) {
						return i, NewSyntaxError(i, "too short unicode character")
					}
					if !((data[i] >= '0' && data[i] <= '9') ||
						(data[i] >= 'a' && data[i] <= 'f') ||
						(data[i] >= 'A' && data[i] <= 'F')) {
						return i, NewSyntaxError(i, "invalid unicode character")
					}
				}
			}
		} else if data[i] == '"' {
			return i + 1, nil
		}
	}
	return i, NewSyntaxError(i, "string not closed")
}

func validNumber[DATA string | []byte](data DATA, i int) (outi int, err error) {
	i--
	// sign
	if data[i] == '-' {
		i++
		if i == len(data) {
			return i, NewSyntaxError(i, "sign character at end of data")
		}
		if data[i] < '0' || data[i] > '9' {
			return i, NewSyntaxError(i, "expected digit after sign")
		}
	}
	// int
	if i == len(data) {
		return i, NewSyntaxError(i, "short data for number")
	}
	if data[i] == '0' {
		i++
	} else {
		for ; i < len(data); i++ {
			if data[i] >= '0' && data[i] <= '9' {
				continue
			}
			break
		}
	}
	// frac
	if i == len(data) {
		return i, nil
	}
	if data[i] == '.' {
		i++
		if i == len(data) {
			return i, NewSyntaxError(i, "expected digit following dot")
		}
		if data[i] < '0' || data[i] > '9' {
			return i, NewSyntaxError(i, "expected digit following dot")
		}
		i++
		for ; i < len(data); i++ {
			if data[i] >= '0' && data[i] <= '9' {
				continue
			}
			break
		}
	}
	// exp
	if i == len(data) {
		return i, nil
	}
	if data[i] == 'e' || data[i] == 'E' {
		i++
		if i == len(data) {
			return i, NewSyntaxError(i, "expected digit following exponent in exp number")
		}
		if data[i] == '+' || data[i] == '-' {
			i++
		}
		if i == len(data) {
			return i, NewSyntaxError(i, "expected digit following sign in exp number")
		}
		if data[i] < '0' || data[i] > '9' {
			return i, NewSyntaxError(i, "expected valid digit in exp number")
		}
		i++
		for ; i < len(data); i++ {
			if data[i] >= '0' && data[i] <= '9' {
				continue
			}
			break
		}
	}
	return i, nil
}

func validTrue[DATA string | []byte](data DATA, i int) (outi int, err error) {
	if i+3 <= len(data) && data[i] == 'r' && data[i+1] == 'u' &&
		data[i+2] == 'e' {
		return i + 3, nil
	}
	return 0, NewSyntaxError(i, "expected 'true'")
}

func validFalse[DATA string | []byte](data DATA, i int) (outi int, err error) {
	if i+4 <= len(data) && data[i] == 'a' && data[i+1] == 'l' &&
		data[i+2] == 's' && data[i+3] == 'e' {
		return i + 4, nil
	}
	return 0, NewSyntaxError(i, "expected 'false'")
}

func validNull[DATA string | []byte](data DATA, i int) (outi int, err error) {
	if i+3 <= len(data) && data[i] == 'u' && data[i+1] == 'l' && data[i+2] == 'l' {
		return i + 3, nil
	}
	return 0, NewSyntaxError(i, "expected 'null'")
}

func validSpace[DATA string | []byte](data DATA, i int) int {
	for ; i < len(data); i++ {
		switch data[i] {
		case ' ', '\t', '\n', '\r':
			continue
		default:
			return i
		}
	}
	return i
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
