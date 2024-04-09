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
func Parse[DATA string | []byte](data DATA) (ResultImpl, error) {
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
func ParseNext[DATA string | []byte](data DATA, i int) (int, ResultImpl, error) {
	return validPayload(data, i)
}

func validPayload[DATA string | []byte](data DATA, i int) (outi int, res ResultImpl, err error) {
	i = validSpace(data, i)
	if i >= len(data) {
		return 0, ResultImpl{}, NewSyntaxError(i, "no content found")
	}
	res.index = i
	i, res.typ, err = validAny(data, i)
	if err != nil {
		return 0, ResultImpl{}, err
	}
	res.Raw = string(data[res.index:i])
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
