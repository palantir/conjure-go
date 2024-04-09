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

package dj_test

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/palantir/conjure-go/v6/dj"
	"github.com/stretchr/testify/assert"
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
	b := make([]byte, 20)
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
	testRune := func(t *testing.T, n int, r rune) {
		var sb strings.Builder
		sb.WriteRune(r)
		str := sb.String()

		quoted := dj.QuoteString(str)

		writeBuf := &bytes.Buffer{}
		_, err := dj.WriteQuotedString(writeBuf, str)
		require.NoError(t, err)
		written := writeBuf.String()

		require.Equal(t, quoted, written, "WriteQuotedString should produce the same output as QuoteString")

		quotedBytes := string(dj.AppendQuotedBytes(nil, []byte(str)))
		require.Equal(t, quoted, quotedBytes, "AppendQuotedBytes should produce the same output as QuoteString")

		jsonBuf := &bytes.Buffer{}
		enc := json.NewEncoder(jsonBuf)
		enc.SetEscapeHTML(false)
		require.NoError(t, enc.Encode(str))
		stdlibOutput := strings.TrimSpace(jsonBuf.String())

		if stdlibOutput == `"�"` {
			assert.Equal(t, `"\ufffd"`, quoted, "case %d unexpected quoted string for rune 0x%x", n, r)
		} else {
			assert.Equal(t, stdlibOutput, quoted, "case %d unexpected quoted string for rune 0x%x", n, r)
		}
	}

	for _, test := range []struct {
		Name  string
		Table *unicode.RangeTable
	}{
		{Name: "Numbers", Table: unicode.N},
		{Name: "Punctuation", Table: unicode.P},
		{Name: "Symbols", Table: unicode.S},
		{Name: "Spaces", Table: unicode.Z},
		{Name: "Control characters", Table: unicode.C},
		{Name: "Mark characters", Table: unicode.M},
	} {
		t.Run(test.Name, func(t *testing.T) {
			t.Run("16-bit", func(t *testing.T) {
				for n, r16 := range test.Table.R16 {

					for i := r16.Lo; i <= r16.Hi; i += r16.Stride {
						testRune(t, n, rune(i))
					}
				}
			})
			t.Run("32-bit", func(t *testing.T) {
				for n, r32 := range test.Table.R32 {
					for i := r32.Lo; i <= r32.Hi; i += r32.Stride {
						testRune(t, n, rune(i))
					}
				}
			})
		})
	}
}
