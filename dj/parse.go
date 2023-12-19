// Package ggjson provides searching for json strings.
// It is forked from github.com/tidwall/gjson, a great project!
package dj

import (
	"math"
	"strconv"
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
	// Type is the json type
	Type Type
	// Raw is the raw json
	Raw string
	// Str is the json string
	Str string
	// Num is the json number
	Num float64
	// Index of raw value in original json, zero means index unknown
	Index int
}

// Exists returns true if value exists.
//
//	 if gjson.Get(json, "name.last").Exists(){
//			println("value exists")
//	 }
func (t Result) Exists() bool {
	return t.Type != Null || len(t.Raw) != 0
}

// String returns a string representation of the value.
func (t Result) String() (string, error) {
	switch t.Type {
	default:
		return "", TypeMismatchError{Index: t.Index, Want: String.String(), Got: t.Type}
	case String:
		return t.Str, nil
	}
}

// Bool returns a boolean representation.
func (t Result) Bool() (bool, error) {
	switch t.Type {
	default:
		return false, TypeMismatchError{Index: t.Index, Want: "boolean", Got: t.Type}
	case True:
		return true, nil
	case False:
		return false, nil
	}
}

// Int returns an integer representation.
func (t Result) Int() (int64, error) {
	if t.Type != Number {
		return 0, TypeMismatchError{Index: t.Index, Want: Number.String(), Got: t.Type}
	}
	// try to directly convert the float64 to int64
	i, ok := safeInt(t.Num)
	if ok && i >= 0 {
		return i, nil
	}
	// now try to parse the raw string
	i, err := parseInt(t.Raw)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// Float returns an float64 representation.
func (t Result) Float() (float64, error) {
	switch t.Str {
	case "NaN":
		return math.NaN(), nil
	case "Infinity":
		return math.Inf(1), nil
	case "-Infinity":
		return math.Inf(-1), nil
	}
	if t.Type != Number {
		return 0, TypeMismatchError{Index: t.Index, Want: Number.String(), Got: t.Type}
	}
	return t.Num, nil
}

// TODO: move this
// ForEach iterates through values.
// If the result represents a non-existent value, then no values will be
// iterated. If the result is an Object, the iterator will pass the key and
// value of each item. If the result is an Array, the iterator will only pass
// the value of each item. If the result is not a JSON array or object, the
// iterator will pass back one value equal to the result.

func (t Result) ObjectIterator(i int) (ObjectIterator, int, error) {
	if !t.Exists() || t.Type != Object {
		return ObjectIterator{}, 0, TypeMismatchError{Index: t.Index, Want: Object.String(), Got: t.Type}
	}
	return ObjectIterator{}, i + 1, nil
}

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
	var vesc bool
	for ; i < len(json); i++ {
		if json[i] != '"' {
			continue
		}
		keyOffset := i
		i, str, vesc, err = parseString(json, i+1)
		if err != nil {
			return Result{}, Result{}, i, err
		}
		key.Type = String
		if vesc {
			key.Str = unescape(str[1 : len(str)-1])
		} else {
			key.Str = str[1 : len(str)-1]
		}
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
			return Result{}, Result{}, i, err
		}
		value.Index = valOffset + t.Index

		return key, value, i, nil
	}
	return Result{}, Result{}, i, SyntaxError{Index: i, Msg: "expected object entry"}
}

func (t Result) ArrayIterator(i int) (ArrayIterator, int, error) {
	if !t.Exists() || t.Type != Array {
		return ArrayIterator{}, 0, TypeMismatchError{Index: t.Index, Want: Array.String(), Got: t.Type}
	}
	return ArrayIterator{}, i + 1, nil
}

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
	return Result{}, i, SyntaxError{Index: i, Msg: "expected array element"}
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
		return nil, SyntaxError{Index: t.Index, Msg: "unrecognized result type"}
	case Null:
		return nil, nil
	case False:
		return false, nil
	case True:
		return true, nil
	case Number:
		return t.Num, nil
	case String:
		return t.Str, nil
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
			mapValue[key.Str] = v
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

// Parse parses the json and returns a result.
//
// This function expects that the json is well-formed, and does not validate.
// Invalid json will not panic, but it may return back unexpected results.
// If you are consuming JSON from an unpredictable source then you may want to
// use the Valid function first.
func Parse[DATA string | []byte](json DATA) (Result, error) {
	if err := Valid(json); err != nil {
		return Result{}, err
	}
	_, res, err := parseAny(string(json), 0)
	return res, err
}

func parseString(json string, i int) (int, string, bool, error) {
	var s = i
	for ; i < len(json); i++ {
		if json[i] > '\\' {
			continue
		}
		if json[i] == '"' {
			return i + 1, json[s-1 : i+1], false, nil
		}
		if json[i] == '\\' {
			i++
			for ; i < len(json); i++ {
				if json[i] > '\\' {
					continue
				}
				if json[i] == '"' {
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
	return i, json[s-1:], false, SyntaxError{Index: i, Msg: "invalid character for string"}
}

func parseNumber(json string, i int) (int, string) {
	var s = i
	i++
	for ; i < len(json); i++ {
		if json[i] <= ' ' || json[i] == ',' || json[i] == ']' ||
			json[i] == '}' {
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
		return 0, SyntaxError{Index: i, Msg: "short data for int"}
	}
	for ; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			n = n*10 + int64(s[i]-'0')
		} else {
			return 0, SyntaxError{Index: i, Msg: "invalid character for int"}
		}
	}
	if sign {
		return n * -1, nil
	}
	return n, nil
}

// safeInt validates a given JSON number
// ensures it lies within the minimum and maximum representable JSON numbers
func safeInt(f float64) (n int64, ok bool) {
	// https://tc39.es/ecma262/#sec-number.min_safe_integer
	// https://tc39.es/ecma262/#sec-number.max_safe_integer
	if f < -9007199254740991 || f > 9007199254740991 {
		return 0, false
	}
	return int64(f), true
}

// parse unquoted values (true, false, null)
func parseLiteral(json string, i int) (int, string) {
	var s = i
	i++
	for ; i < len(json); i++ {
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
	i++
	depth := 1
	for ; i < len(json); i++ {
		if json[i] >= '"' && json[i] <= '}' {
			switch json[i] {
			case '"':
				i++
				s2 := i
				for ; i < len(json); i++ {
					if json[i] > '\\' {
						continue
					}
					if json[i] == '"' {
						// look for an escaped slash
						if json[i-1] == '\\' {
							n := 0
							for j := i - 2; j > s2-1; j-- {
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
			var vesc bool
			var err error
			i, val, vesc, err = parseString(json, i)
			if err != nil {
				return i, res, err
			}
			res.Type = String
			res.Raw = val
			if vesc {
				res.Str = unescape(val[1 : len(val)-1])
			} else {
				res.Str = val[1 : len(val)-1]
			}
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
			return i, res, nil
		case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'i', 'I', 'N':
			num = true
		}
		if num {
			i, val = parseNumber(json, i)
			res.Raw = val
			res.Type = Number
			res.Num, _ = strconv.ParseFloat(val, 64)
			return i, res, nil
		}
	}
	return i, res, SyntaxError{Index: i, Msg: "invalid character for json"}
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
			dst = append(dst, '\\')
			switch s[i] {
			case '\n':
				dst = append(dst, 'n')
			case '\r':
				dst = append(dst, 'r')
			case '\t':
				dst = append(dst, 't')
			default:
				dst = append(dst, 'u')
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
func runeit(json string) rune {
	n, _ := strconv.ParseUint(json[:4], 16, 64)
	return rune(n)
}

// unescape unescapes a string
func unescape(json string) string {
	var str = make([]byte, 0, len(json))
	for i := 0; i < len(json); i++ {
		switch {
		default:
			str = append(str, json[i])
		case json[i] < ' ':
			return string(str)
		case json[i] == '\\':
			i++
			if i >= len(json) {
				return string(str)
			}
			switch json[i] {
			default:
				return string(str)
			case '\\':
				str = append(str, '\\')
			case '/':
				str = append(str, '/')
			case 'b':
				str = append(str, '\b')
			case 'f':
				str = append(str, '\f')
			case 'n':
				str = append(str, '\n')
			case 'r':
				str = append(str, '\r')
			case 't':
				str = append(str, '\t')
			case '"':
				str = append(str, '"')
			case 'u':
				if i+5 > len(json) {
					return string(str)
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
				str = append(str, 0, 0, 0, 0, 0, 0, 0, 0)
				n := utf8.EncodeRune(str[len(str)-8:], r)
				str = str[:len(str)-8+n]
				i-- // backtrack index by one
			}
		}
	}
	return string(str)
}
