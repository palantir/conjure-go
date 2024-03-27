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
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// Parse parses the json and returns a result.
//
// This function expects that the json is well-formed, and does not validate.
// Invalid json will not panic, but it may return back unexpected results.
// If you are consuming JSON from an unpredictable source then you may want to
// use the Valid function first.
func Parse[DATA string | []byte](json DATA) (Result, error) {
	_, res, err := parseAny(string(json), 0)
	return res, err
}

// parseAny parses the next value from a json string.
// A Result is returned when the hit param is set.
// The return values are (i int, res Result, err error)
func parseAny(json string, i int) (int, Result, error) {
	var res Result
	var val string
	var err error
	switch json[i] {
	case '{':
		i, val, err = parseObject(json, i)
		if err != nil {
			return i, res, err
		}
		res.Raw = val
		res.Type = Object
		res.Index = i
		return i, res, nil
	case '[':
		i, val, err = parseArray(json, i)
		if err != nil {
			return i, res, err
		}
		res.Raw = val
		res.Type = Array
		res.Index = i
		return i, res, nil
	case '"':
		i++
		i, val, err = parseString(json, i)
		if err != nil {
			return i, res, err
		}
		res.Type = String
		res.Raw = val
		res.Index = i
		return i, res, nil
	case 'n':
		if i+3 < len(json) && json[i+1] == 'u' && json[i+2] == 'l' && json[i+3] == 'l' {
			i += 4
			res.Type = Null
			res.Index = i
			return i, res, nil
		}
		return i, res, NewSyntaxError(i, "expected 'null'")
	case 't':
		if i+3 < len(json) && json[i+1] == 'r' && json[i+2] == 'u' && json[i+3] == 'e' {
			i += 4
			res.Type = True
			res.Index = i
			return i, res, nil
		}
		return i, res, NewSyntaxError(i, "expected 'true'")
	case 'f':
		if i+4 < len(json) && json[i+1] == 'a' && json[i+2] == 'l' && json[i+3] == 's' && json[i+4] == 'e' {
			i += 5
			res.Type = False
			res.Index = i
			return i, res, nil
		}
		return i, res, NewSyntaxError(i, "expected 'false'")
	case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		i, val, err = parseNumber(json, i)
		if err != nil {
			return i, res, err
		}
		res.Raw = val
		res.Type = Number
		res.Index = i
		return i, res, nil
	default:
		return i, res, NewSyntaxError(i, "invalid character for json")
	}
}

func parseString(json string, i int) (iOut int, val string, err error) {
	var s = i
	ln := len(json)
	for ; i < ln; i++ {
		ji := json[i]
		if ji < ' ' || ji >= utf8.RuneSelf {
			return i, json[s-1:], NewSyntaxError(i, "invalid character for string")
		}
		if ji == '"' {
			return i + 1, json[s-1 : i+1], nil // String closed
		}
		if i == ln-1 {
			return i, json[s-1:], NewSyntaxError(i, "string not closed")
		}
		if ji == '\\' {
			// escape character: advance through full escape sequence
			i++
			ji = json[i]
			if ji < ' ' || ji >= utf8.RuneSelf {
				return i, json[s-1:], NewSyntaxError(i, "invalid character after escape character")
			}
			if i == ln-1 {
				return i, json[s-1:], NewSyntaxError(i, "escape character at end of data")
			}
			switch ji {
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				// valid escape
			case 'u':
				for j := 0; j < 4; j++ {
					i++
					if i >= ln {
						return i, json[s-1:], NewSyntaxError(i, "too short unicode character")
					}
					if !((ji >= '0' && ji <= '9') || (ji >= 'a' && ji <= 'f') || (ji >= 'A' && ji <= 'F')) {
						return i, json[s-1:], NewSyntaxError(i, "invalid unicode character")
					}
				}
			default:
				return i, json[s-1:], NewSyntaxError(i, "invalid escape character")
			}
			// continue
		}
	}
	return i, json[s-1:], NewSyntaxError(i, "string not closed")
}

func parseNumber(json string, i int) (int, string, error) {
	var s = i
	if i == len(json) {
		return i, json[s:i], NewSyntaxError(i, "short data for number")
	}
	for ; i < len(json); i++ {
		switch json[i] {
		case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', 'e', 'E':
			continue
		case ' ', '\t', '\n', '\r', ',', '}', ']':
			return i, json[s:i], nil
		default:
			return i, json[s:i], NewSyntaxError(i, "invalid character for number")
		}
	}
	return i, json[s:i], nil
}

// returns the substring containing the json value up to the closing brace.
func parseObject(json string, i int) (int, string, error) {
	s := i
	if json[i] != '{' {
		return i, json[s : i+1], NewSyntaxError(i, "object expected opening bracket")
	}
	i++
	for ; i < len(json); i++ {
		i = parseSpace(json, i)
		if i >= len(json) {
			return i, json[s:], NewSyntaxError(i, "object not closed")
		}
		if json[i] == '"' {
			// parse key
			i++
			if i >= len(json) {
				return i, json[s:], NewSyntaxError(i, "object missing string for key")
			}
			var err error
			i, _, err = parseString(json, i)
			if err != nil {
				return i, json[s:i], err
			}

			// parse colon
			i = parseSpace(json, i)
			if i >= len(json) {
				return i, json[s:], NewSyntaxError(i, "object key not closed")
			}
			if json[i] != ':' {
				return i, json[s : i+1], NewSyntaxError(i, "object key-value pair expected ':'")
			}

			// parse value
			i = parseSpace(json, i)
			if i >= len(json) {
				return i, json[s:], NewSyntaxError(i, "object value not closed")
			}
			i, _, err = parseAny(json, i)
			if err != nil {
				return i, json[s:i], err
			}

			// parse comma or closing brace
			i = parseSpace(json, i)
			if i >= len(json) {
				return i, json[s:], NewSyntaxError(i, "object value missing comma or closing brace")
			}
			if json[i] == ',' {
				continue
			}
			if json[i] == '}' {
				return i + 1, json[s : i+1], nil
			}
			return i, json[s:i], NewSyntaxError(i, "object value expected comma or closing brace")
		}
		if json[i] == '}' {
			return i + 1, json[s : i+1], nil
		}
		return i, json[s : i+1], NewSyntaxError(i, "object expected key or closing brace")
	}
	return i, json[s : i+1], NewSyntaxError(i, "invalid character for object")
}

func parseArray(json string, i int) (int, string, error) {
	s := i
	if json[i] != '[' {
		return i, json[s : i+1], NewSyntaxError(i, "array expected opening bracket")
	}
	i++
	for ; i < len(json); i++ {
		// parse value
		i = parseSpace(json, i)
		if i >= len(json) {
			return i, json[s:], NewSyntaxError(i, "array not closed")
		}
		var err error
		i, _, err = parseAny(json, i)
		if err != nil {
			return i, json[s:i], err
		}

		// parse comma or closing bracket
		i = parseSpace(json, i)
		if i >= len(json) {
			return i, json[s:], NewSyntaxError(i, "array missing comma or closing bracket")
		}
		if json[i] == ',' {
			continue
		}
		if json[i] == ']' {
			return i + 1, json[s : i+1], nil
		}
		return i, json[s:i], NewSyntaxError(i, "array expected comma or closing bracket")
	}
	return i, json[s : i+1], NewSyntaxError(i, "invalid character for array")
}

// returns the index of the first non-space character
func parseSpace(json string, i int) int {
	for ; i < len(json); i++ {
		switch json[i] {
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
func unescape[DATA string | []byte](json DATA) (string, error) {
	sb := new(strings.Builder)
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
