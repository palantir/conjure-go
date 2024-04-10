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
	stdjson "encoding/json"
	"math"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// Type is Result type
type Type int

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

// Result represents a json value that is returned from Get().
type Result struct {
	// Type is the json type, one of Null, False, True, Number, String, Object, or Array.
	Type Type
	// Raw is the raw json
	Raw string
	// Index returns the index of the raw value in the original json.
	// If the index is unknown, it will return 0.
	// Generally this should only be used for debugging and error reporting purposes.
	Index int
}

// String returns a string representation of the value.
// If the value is not a json string, it will return an error.
func (t Result) String() (string, error) {
	switch t.Type {
	default:
		return "", NewTypeMismatchError(t, String.String())
	case String:
		if len(t.Raw) < 2 || t.Raw[0] != '"' || t.Raw[len(t.Raw)-1] != '"' {
			return "", NewSyntaxError(t.Index, "invalid string")
		}
		// unescape on first \\ found
		for i := 1; i < len(t.Raw); i++ {
			if t.Raw[i] == '\\' {
				// trim quotes
				return unescape(t.Raw[1 : len(t.Raw)-1])
			}
		}
		// trim quotes
		return t.Raw[1 : len(t.Raw)-1], nil
	}
}

// Text returns a decoded string as a byte slice, suitable for encoding.TextUnmarshaler.
// If the value is not a json string, it will return an error.
func (t Result) Text() ([]byte, error) {
	switch t.Type {
	default:
		return nil, NewTypeMismatchError(t, String.String())
	case String:
		if len(t.Raw) < 2 || t.Raw[0] != '"' || t.Raw[len(t.Raw)-1] != '"' {
			return nil, NewSyntaxError(t.Index, "invalid string")
		}
		// unescape on first \\ found
		for i := 1; i < len(t.Raw); i++ {
			if t.Raw[i] == '\\' {
				// trim quotes
				s, err := unescape(t.Raw[1 : len(t.Raw)-1])
				if err != nil {
					return nil, err
				}
				return []byte(s), nil
			}
		}
		// trim quotes
		return []byte(t.Raw[1 : len(t.Raw)-1]), nil
	}
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
	n, err := parseInt(t.Raw)
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
	case `"Infinity"`:
		return math.Inf(1), nil
	case `"-Infinity"`:
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
	var key, value Result
	for iter.HasNext(t, i) {
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
	json := t.Raw
	if i == 0 {
		i++ // skip the first '{'
	}
	for ; i < len(json); i++ {
		if json[i] <= ' ' || json[i] == ',' || json[i] == ':' {
			continue
		}
		return json[i] != '}'
	}
	return false
}

// Next returns the next key, value, and index to pass to HasNext().
func (ObjectIterator) Next(t Result, i int) (key Result, value Result, iOut int, err error) {
	json := t.Raw
	if i == 0 {
		i++ // skip the first '{'
	}
	var str string
	for ; i < len(json); i++ {
		if json[i] != '"' {
			continue
		}
		keyOffset := i
		i, str, _, err = parseString(json, i+1)
		if err != nil {
			return Result{}, Result{}, 0, err
		}
		key.Type = String
		key.Raw = str
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
	case False, True:
		return t.Bool()
	case Number:
		return t.Float()
	case String:
		return t.String()
	case Object:
		mapValue := make(map[string]any)
		iter, i, err := t.ObjectIterator(0)
		if err != nil {
			return nil, err
		}
		var key, value Result
		for iter.HasNext(t, i) {
			key, value, i, err = iter.Next(t, i)
			if err != nil {
				return nil, err
			}
			v, err := value.Value()
			if err != nil {
				return nil, err
			}
			k, err := key.String()
			if err != nil {
				return nil, err
			}
			mapValue[k] = v
		}
		return mapValue, nil
	case Array:
		arrayValue := make([]any, 0)
		iter, i, err := t.ArrayIterator(0)
		if err != nil {
			return nil, err
		}
		for iter.HasNext(t, i) {
			var value Result
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

// Parse parses the json and returns a result. It returns a SyntaxError if the data is invalid JSON.
// The input can be a string or []byte.
// The returned Result's Index field is the starting index of the value.
// To parse multiple JSON values in a single string, use ParseNext.
func Parse[DATA string | []byte](json DATA) (Result, error) {
	if err := Valid(json); err != nil {
		return Result{}, err
	}
	_, res, err := parseAny(string(json), 0)
	return res, err
}

// ParseNext parses the next value from a json string.
// This function is useful when you have multiple json values in a single string.
// The return values are (i int, res Result, err error)
// The i is the index of the next character after the parsed value.
// When the returned i is equal to len(data), there are no more values to parse and the next call will error.
func ParseNext[DATA string | []byte](json DATA, i int) (int, Result, error) {
	i, res, err := parseAny(string(json), i)
	if err != nil {
		return 0, Result{}, err
	}
	if err := Valid(res.Raw); err != nil {
		return 0, Result{}, err
	}
	return i, res, nil
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

func parseInt(s string) (n int64, err error) {
	var i int
	var sign bool
	if len(s) > 0 && s[0] == '-' {
		sign = true
		i++
	}
	if i == len(s) {
		return 0, NewSyntaxError(i, "short data for int")
	}
	for ; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			n = n*10 + int64(s[i]-'0')
		} else {
			return 0, NewSyntaxError(i, "invalid character for int")
		}
	}
	if sign {
		return n * -1, nil
	}
	return n, nil
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
		if json[i] <= ' ' {
			continue
		}
		var num bool
		switch json[i] {
		case '{':
			i, val = parseSquash(json, i)
			res.Raw = val
			res.Type = Object
			res.Index = i
			return i, res, nil
		case '[':
			i, val = parseSquash(json, i)
			res.Raw = val
			res.Type = Array
			res.Index = i
			return i, res, nil
		case '"':
			i++
			var err error
			i, val, _, err = parseString(json, i)
			if err != nil {
				return i, res, err
			}
			res.Type = String
			res.Raw = val
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

var hexchars = [...]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f',
}

func appendHex16(dst []byte, x uint16) []byte {
	return append(dst,
		hexchars[x>>12&0xF], hexchars[x>>8&0xF],
		hexchars[x>>4&0xF], hexchars[x>>0&0xF],
	)
}

// AppendJSONString is a convenience function that converts the provided string
// to a valid JSON string and appends it to dst.
func AppendJSONString(dst []byte, s string) []byte {
	dst = append(dst, make([]byte, len(s)+2)...)
	dst = append(dst[:len(dst)-len(s)-2], '"')
	for i := 0; i < len(s); i++ {
		if s[i] < ' ' {
			switch s[i] {
			case '\n':
				dst = append(dst, '\\', 'n')
			case '\r':
				dst = append(dst, '\\', 'r')
			case '\t':
				dst = append(dst, '\\', 't')
			default:
				dst = append(dst, '\\', 'u')
				dst = appendHex16(dst, uint16(s[i]))
			}
		} else if s[i] == '>' || s[i] == '<' || s[i] == '&' {
			dst = append(dst, '\\', 'u')
			dst = appendHex16(dst, uint16(s[i]))
		} else if s[i] == '\\' {
			dst = append(dst, '\\', '\\')
		} else if s[i] == '"' {
			dst = append(dst, '\\', '"')
		} else if s[i] > 127 {
			// read utf8 character
			r, n := utf8.DecodeRuneInString(s[i:])
			if n == 0 {
				break
			}
			if r == utf8.RuneError && n == 1 {
				dst = append(dst, `\ufffd`...)
			} else if r == '\u2028' || r == '\u2029' {
				dst = append(dst, `\u202`...)
				dst = append(dst, hexchars[r&0xF])
			} else {
				dst = append(dst, s[i:i+n]...)
			}
			i = i + n - 1
		} else {
			dst = append(dst, s[i])
		}
	}
	return append(dst, '"')
}

// runeit returns the rune from the the \uXXXX
func runeit[DATA string | []byte](json DATA) rune {
	n, _ := strconv.ParseUint(string(json[:4]), 16, 64)
	return rune(n)
}

// unescape unescapes a string
func unescape[DATA string | []byte](json DATA) (string, error) {
	var sb strings.Builder
	sb.Grow(len(json))
	ln := len(json)
	for i := 0; i < ln; i++ {
		switch {
		default:
			sb.WriteByte(json[i])
		case json[i] < ' ':
			return "", NewSyntaxError(i, "invalid character for encoded string")
		case json[i] == '\\':
			i++
			if i >= len(json) {
				return "", NewSyntaxError(i, "incomplete escape sequence in encoded string")
			}
			switch json[i] {
			default:
				return "", NewSyntaxError(i, "invalid escape sequence in encoded string")
			case '\\':
				sb.WriteByte('\\')
			case '/':
				sb.WriteByte('/')
			case 'b':
				sb.WriteByte('\b')
			case 'f':
				sb.WriteByte('\f')
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			case '"':
				sb.WriteByte('"')
			case 'u':
				if i+5 > len(json) {
					return "", NewSyntaxError(i, "incomplete unicode sequence in encoded string")
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
				sb.Write(buf[:n])
				i-- // backtrack index by one
			}
		}
	}
	return sb.String(), nil
}
