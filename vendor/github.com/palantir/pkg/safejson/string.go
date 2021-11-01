package safejson

import (
	"reflect"
	"unicode/utf8"
	"unsafe"
)

// constants we can preallocate
var (
	// jsonReplace holds the values below 128 which require replacement in JSON strings.
	// If an entry is nil, the rune can be used as-is.
	// All values are nil except for the ASCII control characters (0-31), the
	// double quote ("), and the backslash character ("\").
	jsonReplace = [utf8.RuneSelf][]byte{
		'\\': []byte(`\\`),
		'"':  []byte(`\"`),
		'\n': []byte(`\n`),
		'\r': []byte(`\r`),
		'\t': []byte(`\t`),
	}
)

func init() {
	const hex = "0123456789abcdef"
	for i := 0; i < ' '; i++ {
		switch i {
		case '\n', '\r', '\t':
		default:
			// This encodes bytes < 0x20 except for \t, \n and \r.
			jsonReplace[i] = append([]byte(`\u00`), hex[i>>4], hex[i&0xF])
		}
	}
}

func QuoteString(s string) string {
	return string(AppendQuotedString(nil, s))
}

func QuotedStringLength(s string) int {
	// Create an unsafe copy of s as a []byte.
	// This is safe because we do not mutate b in WriteQuotedBytes.
	b := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (*reflect.StringHeader)(unsafe.Pointer(&s)).Data,
		Len:  len(s),
		Cap:  len(s),
	}))
	return QuotedBytesLength(b)
}

// AppendQuotedString quotes and JSON-escapes s and appends the result to dst.
// The resulting slice is returned in case it was resized by append().
func AppendQuotedString(dst []byte, s string) []byte {
	// Create an unsafe copy of s as a []byte.
	// This is safe because we do not mutate b in WriteQuotedBytes.
	b := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (*reflect.StringHeader)(unsafe.Pointer(&s)).Data,
		Len:  len(s),
		Cap:  len(s),
	}))
	return AppendQuotedBytes(dst, b)
}

func QuotedBytesLength(b []byte) int {
	out := 2 // open/close quotes
	for i := 0; i < len(b); {
		if b[i] < utf8.RuneSelf {
			repl := jsonReplace[b[i]]
			if repl == nil {
				out++
			} else {
				out += len(repl)
			}
			i++
			continue
		}
		c, size := utf8.DecodeRune(b[i:])
		i += size
		if c == utf8.RuneError && size == 1 {
			out += len(`\ufffd`)
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			out += len(`\u2028`)
			continue
		}
		out += size
	}
	return out
}

// AppendQuotedBytes quotes and JSON-escapes b and appends the result to dst.
// The resulting slice is returned in case it was resized by append().
func AppendQuotedBytes(dst []byte, b []byte) []byte {
	dst = append(dst, '"')
	start := 0
	for i := 0; i < len(b); {
		if b[i] < utf8.RuneSelf {
			repl := jsonReplace[b[i]]
			if repl == nil {
				i++
				continue
			}
			if start < i {
				dst = append(dst, b[start:i]...)
			}
			dst = append(dst, repl...)
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRune(b[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				dst = append(dst, b[start:i]...)
			}
			dst = append(dst, `\ufffd`...)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				dst = append(dst, b[start:i]...)
			}
			if c == '\u2028' {
				dst = append(dst, `\u2028`...)
			} else {
				dst = append(dst, `\u2029`...)
			}
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(b) {
		dst = append(dst, b[start:]...)
	}
	dst = append(dst, '"')
	return dst
}
