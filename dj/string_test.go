package dj_test

import (
	"bytes"
	"encoding/json"
	"github.com/palantir/conjure-go/v6/dj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strings"
	"testing"
	"time"
	"unicode"
)

func TestQuote(t *testing.T) {
	for _, test := range []struct {
		in  string
		out string
	}{
		{"plain", `"plain"`},
		{"new\nline", `"new\nline"`},
		{"\n❤️\t", `"\n❤️\t"`},
		{"I❤️NY", `"I❤️NY"`},
		{"I❤️", `"I❤️"`},
		{"\u2028", `"\u2028"`},
		{"\u2029", `"\u2029"`},
	} {
		t.Run(test.in, func(t *testing.T) {
			out := dj.QuoteString(test.in)
			require.Equal(t, test.out, out)

			ref, err := json.Marshal(test.in)
			require.NoError(t, err)
			require.Equal(t, string(ref), out)
			require.Len(t, out, dj.QuotedLength(test.in))
		})
	}
}

// TestEncodeString_RandomData is a fuzzing test that throws random data at the QuoteString
// function looking for panics.
func TestQuoteString_RandomData(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 200)
	for i := 0; i < 100000; i++ {
		n, err := rnd.Read(b[:rand.Int()%len(b)])
		require.NoError(t, err)

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err = enc.Encode(string(b[:n]))
		require.NoError(t, err)
		marshaled := strings.TrimSpace(buf.String())

		quoted := dj.QuoteString(string(b[:n]))
		require.Len(t, quoted, dj.QuotedLength(string(b[:n])))

		buf.Reset()
		_, err = dj.WriteQuotedString(buf, string(b[:n]))
		require.NoError(t, err)
		written := buf.String()
		require.Len(t, written, dj.QuotedLength(string(b[:n])))

		require.Equal(t, written, quoted, "QuoteString should produce the same output as WriteQuotedString")

		require.Equal(t, marshaled, quoted, "QuoteString should produce the same output as json.Encoder")
	}
}

func TestQuoteAllUnicode(t *testing.T) {
	var sb strings.Builder
	for _, table := range []*unicode.RangeTable{
		unicode.C, unicode.M, unicode.N, unicode.P, unicode.S, unicode.Z,
	} {
		for _, r16 := range table.R16 {
			for i := r16.Lo; i <= r16.Hi; i += r16.Stride {
				sb.WriteRune(rune(i))
			}
		}
		for _, r32 := range table.R32 {
			for i := r32.Lo; i <= r32.Hi; i += r32.Stride {
				sb.WriteRune(rune(i))
			}
		}
	}
	str := sb.String()

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(str)
	require.NoError(t, err)
	expected := string(bytes.TrimSpace(buf.Bytes()))

	quoted := dj.QuoteString(str)

	buf.Reset()
	_, err = dj.WriteQuotedString(buf, str)
	require.NoError(t, err)
	written := buf.String()
	require.Equal(t, quoted, written, "WriteQuotedString should produce the same output as QuoteString")

	if !assert.Equal(t, expected, quoted) {
		for i := 0; i < len(expected) && i < len(quoted); i++ {
			if expected[i] != quoted[i] {
				t.Logf("first difference at index %d:\n want: %s\n  got: %s", i, expected[i:], quoted[i:])
				break
			}
		}
	}
}
