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
	"encoding/hex"
	"encoding/json"
	"io"
	"math/rand"
	"testing"

	"github.com/palantir/conjure-go/v6/dj"
	"github.com/tidwall/gjson"
)

func BenchmarkQuoteString(b *testing.B) {
	for _, bb := range []struct {
		name  string
		input string
	}{
		{
			name:  "10B hex",
			input: newBenchInput(10, true),
		},
		{
			name:  "10B utf",
			input: newBenchInput(10, false),
		},
		{
			name:  "10K hex",
			input: newBenchInput(10000, true),
		},
		{
			name:  "10K utf",
			input: newBenchInput(10000, false),
		},
	} {
		input := bb.input
		b.Run(bb.name, func(b *testing.B) {
			doBench := func(doBenchFn func(b *testing.B, input string, heuristicLen bool, calcLen bool)) func(*testing.B) {
				return func(b *testing.B) {
					b.Run("precalculated", func(b *testing.B) {
						doBenchFn(b, input, false, true)
					})
					b.Run("heuristic", func(b *testing.B) {
						doBenchFn(b, input, true, false)
					})
					b.Run("unallocated", func(b *testing.B) {
						doBenchFn(b, input, false, false)
					})
				}
			}
			b.Run("gjson.AppendJSONString", doBench(doBenchGJSONAppendJSONString))
			b.Run("dj.AppendQuotedString", doBench(doBenchAppendQuotedString))
			b.Run("json.Encoder", doBench(doBenchStringJSONEncoder))
			b.Run("dj.WriteQuotedString", doBench(doBenchWriteQuotedString))

			b.Run("json.Marshal", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					data, err := json.Marshal(input)
					if err != nil {
						panic(err)
					}
					_ = data
				}
			})

			b.Run("dj.WriteQuotedString Discard", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					data, err := dj.WriteQuotedString(io.Discard, input)
					if err != nil {
						b.Fatal(err)
					}
					_ = data
				}
			})

			b.Run("dj.QuotedLength", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					data := dj.QuotedLength(input)
					_ = data
				}
			})
		})
	}
}

func newBenchInput(strLen int, ascii bool) string {
	randData := make([]byte, strLen)
	fixedRand := rand.New(rand.NewSource(0))
	_, err := fixedRand.Read(randData)
	if err != nil {
		panic(err)
	}
	if ascii {
		return hex.EncodeToString(randData[:hex.DecodedLen(strLen)])
	}
	return string(randData)
}

func doBenchAppendQuotedString(b *testing.B, input string, heuristicLen bool, calcLen bool) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var buf []byte
		if heuristicLen {
			buf = make([]byte, 0, 2+len(input))
		} else if calcLen {
			buf = make([]byte, 0, dj.QuotedLength(input))
		}
		data := dj.AppendQuotedString(buf, input)
		_ = data
	}
}

func doBenchWriteQuotedString(b *testing.B, input string, heuristicLen bool, calcLen bool) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var buf *bytes.Buffer
		if heuristicLen {
			buf = bytes.NewBuffer(make([]byte, 0, 2+len(input)))
		} else if calcLen {
			buf = bytes.NewBuffer(make([]byte, 0, dj.QuotedLength(input)))
		} else {
			buf = new(bytes.Buffer)
		}
		_, err := dj.WriteQuotedString(buf, input)
		if err != nil {
			b.Fatal(err)
		}
		data := buf.Bytes()
		_ = data
	}
}

func doBenchStringJSONEncoder(b *testing.B, input string, heuristicLen bool, calcLen bool) {
	b.Run("json.Encoder", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var buf *bytes.Buffer
			if heuristicLen {
				buf = bytes.NewBuffer(make([]byte, 0, 2+len(input)))
			} else if calcLen {
				buf = bytes.NewBuffer(make([]byte, 0, dj.QuotedLength(input)))
			} else {
				buf = new(bytes.Buffer)
			}
			enc := json.NewEncoder(buf)
			enc.SetEscapeHTML(false)
			err := enc.Encode(input)
			if err != nil {
				panic(err)
			}
			data := buf.Bytes()
			_ = data
		}
	})
}

func doBenchGJSONAppendJSONString(b *testing.B, input string, heuristicLen bool, calcLen bool) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var buf []byte
		if heuristicLen {
			buf = make([]byte, 0, 2+len(input))
		} else if calcLen {
			buf = make([]byte, 0, dj.QuotedLength(input))
		}
		data := gjson.AppendJSONString(buf, input)
		_ = data
	}
}
