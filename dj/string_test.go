package dj_test

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/palantir/conjure-go/v6/dj"
	"github.com/stretchr/testify/require"
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
		outJSONEncoder := bytes.Clone(bytes.TrimRight(buf.Bytes(), "\n"))
		buf.Reset()

		outAppendString := dj.QuoteString(string(b[:n]))
		require.Len(t, outAppendString, dj.QuotedLength(string(b[:n])))

		_, err = dj.WriteQuotedString(buf, string(b[:n]))
		require.NoError(t, err)
		outWriteString := bytes.Clone(bytes.TrimRight(buf.Bytes(), "\n"))
		require.Len(t, outWriteString, dj.QuotedLength(string(b[:n])))

		require.Equal(t, outAppendString, string(outWriteString))
		require.Equal(t, outAppendString, string(outJSONEncoder))
	}
}
