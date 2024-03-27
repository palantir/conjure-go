package dj

import (
	"math"
	"strconv"
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

// Result represents a json value that is returned from Get().
type Result struct {
	// Type is the json type
	Type Type
	// Index of raw value in original json, zero means index unknown
	Index int
	// Raw is the raw json
	Raw string
}

func (t Result) IsNull() bool {
	return t.Type == Null
}

// String returns a string representation of the value.
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

// Bool returns a boolean representation.
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
func (t Result) Int() (int64, error) {
	if t.Type != Number {
		return 0, NewTypeMismatchError(t, Number.String())
	}
	// now try to parse the raw string
	i, err := parseIntResult(t.Raw)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func parseIntResult(s string) (n int64, err error) {
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

// Float returns an float64 representation.
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

// TODO: move this
// ForEach iterates through values.
// If the result represents a non-existent value, then no values will be
// iterated. If the result is an Object, the iterator will pass the key and
// value of each item. If the result is an Array, the iterator will only pass
// the value of each item. If the result is not a JSON array or object, the
// iterator will pass back one value equal to the result.

func (t Result) ObjectIterator(i int) (ObjectIterator, int, error) {
	if t.Type != Object {
		return ObjectIterator{}, 0, NewTypeMismatchError(t, Object.String())
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
	for ; i < len(json); i++ {
		if json[i] != '"' {
			continue
		}
		keyOffset := i
		i, str, err = parseString(json, i+1)
		if err != nil {
			return Result{}, Result{}, i, err
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
			return Result{}, Result{}, i, err
		}
		value.Index = valOffset + t.Index

		return key, value, i, nil
	}
	return Result{}, Result{}, i, NewSyntaxError(i, "expected object entry")
}

func (t Result) ArrayIterator(i int) (ArrayIterator, int, error) {
	if t.Type != Array {
		return ArrayIterator{}, 0, NewTypeMismatchError(t, Array.String())
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
		return nil, NewSyntaxError(t.Index, "unrecognized result type "+t.Type.String())
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
