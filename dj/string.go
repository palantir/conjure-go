// Copyright (c) 2024 Palantir Technologies. All rights reserved.
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
	"io"
	"unicode/utf8"
	"unsafe"
)

// QuoteString quotes and JSON-escapes s.
func QuoteString(s string) string {
	b := appendQuoted(nil, nil, s)
	// convert to string without allocation; safe because b is freshly allocated
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// AppendQuotedString quotes and JSON-escapes s and appends the result to dst.
// The resulting slice is returned in case it was resized by append().
func AppendQuotedString(dst []byte, s string) []byte {
	return appendQuoted(dst, nil, s)
}

// AppendQuotedBytes quotes and JSON-escapes b and appends the result to dst.
// The resulting slice is returned in case it was resized by append().
func AppendQuotedBytes(dst []byte, b []byte) []byte {
	return appendQuoted(dst, b, "")
}

// WriteQuotedString writes s as a quoted and JSON-escaped string to w.
func WriteQuotedString(w io.Writer, s string) (int, error) {
	return writeQuoted(w, nil, s)
}

// WriteQuotedBytes writes b as a quoted and JSON-escaped string to w.
func WriteQuotedBytes(w io.Writer, b []byte) (int, error) {
	return writeQuoted(w, b, "")
}

// QuotedLength returns the length in bytes of the s when quoted and escaped.
// This is useful for pre-allocating memory before AppendQuotedString.
func QuotedLength(s string) int {
	return lengthQuoted(nil, s)
}

// QuotedBytesLength returns the length in bytes of the b when quoted and escaped.
// This is useful for pre-allocating memory before AppendQuotedBytes.
func QuotedBytesLength(b []byte) int {
	return lengthQuoted(b, "")
}

func escapeNextRune(b []byte) (int, []byte) {
	b0 := b[0]
	if b0 < utf8.RuneSelf {
		repl := jsonReplace[b0]
		return 1, repl
	}
	c, size := utf8.DecodeRune(b)
	switch {
	// U+2028 is LINE SEPARATOR.
	// U+2029 is PARAGRAPH SEPARATOR.
	// They are both technically valid characters in JSON strings,
	// but don't work in JSONP, which has to be evaluated as JavaScript,
	// and can lead to security holes there. It is valid JSON to
	// escape them, so we do so unconditionally.
	// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
	case c == '\u2028':
		return size, bu2028
	case c == '\u2029':
		return size, bu2029
	case c == utf8.RuneError:
		return size, buFFFD
	default:
		return size, nil // no replacement
	}
}

// appendQuoted is inspired by json.Marshal's private implementation: https://github.com/golang/go/blob/go1.19.1/src/encoding/json/encode.go#L1102-L1171
func appendQuoted(dst []byte, b []byte, s string) []byte {
	if b == nil && s != "" {
		b = unsafe.Slice(unsafe.StringData(s), len(s)) // convert to []byte without allocation
	}
	dst = append(dst, '"')
	start := 0
	for i := 0; i < len(b); {
		size, repl := escapeNextRune(b[i:])
		if repl == nil {
			i += size
			continue
		}
		// append the bytes before the replacement
		if start < i {
			dst = append(dst, b[start:i]...)
		}
		// append the replacement
		dst = append(dst, repl...)
		i += size
		start = i
	}
	if start < len(b) {
		// append the bytes after the last replacement
		dst = append(dst, b[start:]...)
	}
	dst = append(dst, '"')
	return dst
}

func writeQuoted(w io.Writer, b []byte, s string) (int, error) {
	if b == nil && s != "" {
		b = unsafe.Slice(unsafe.StringData(s), len(s)) // convert to []byte without allocation
	}
	var out int
	n, err := WriteDoubleQuote(w)
	if err != nil {
		return n, err
	}
	out += n
	start := 0
	for i := 0; i < len(b); {
		size, repl := escapeNextRune(b[i:])
		if repl == nil {
			i += size
			continue
		}
		// write the bytes before the replacement
		if start < i {
			n, err := w.Write(b[start:i])
			if err != nil {
				return 0, err
			}
			out += n
		}
		// write the replacement
		n, err := w.Write(repl)
		if err != nil {
			return 0, err
		}
		out += n
		i += size
		start = i
	}
	if start < len(b) {
		// write the bytes after the last replacement
		n, err := w.Write(b[start:])
		if err != nil {
			return 0, err
		}
		out += n
	}
	n, err = WriteDoubleQuote(w)
	if err != nil {
		return n, err
	}
	out += n
	return out, nil
}

func lengthQuoted(b []byte, s string) int {
	if len(b) == 0 && s == "" {
		return 2
	}
	if b == nil && s != "" {
		b = unsafe.Slice(unsafe.StringData(s), len(s)) // convert to []byte without allocation
	}
	out := 2 // open/close quotes
	for i := 0; i < len(b); {
		size, repl := escapeNextRune(b[i:])
		if repl == nil {
			out += size
		} else {
			out += len(repl)
		}
		i += size
	}
	return out
}

// jsonReplace holds the values below 128 which require replacement in JSON strings.
// If an entry is nil, the rune can be used as-is.
// All values are nil except for the ASCII control characters (0-31), the
// double quote ("), and the backslash character ("\").
var jsonReplace = [utf8.RuneSelf][]byte{}

func init() {
	const hex = "0123456789abcdef"
	for i := uint8(0); i < utf8.RuneSelf; i++ {
		switch {
		case i == '\\':
			jsonReplace[i] = []byte(`\\`)
		case i == '"':
			jsonReplace[i] = []byte(`\"`)
		case i == '\n':
			jsonReplace[i] = []byte(`\n`)
		case i == '\r':
			jsonReplace[i] = []byte(`\r`)
		case i == '\t':
			jsonReplace[i] = []byte(`\t`)
		case i < ' ':
			jsonReplace[i] = []byte{'\\', 'u', '0', '0', hex[i>>4], hex[i&0xF]}
		default:
			// leave as nil
		}
	}
}
