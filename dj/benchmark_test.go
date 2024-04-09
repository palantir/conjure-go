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

package dj_test

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"math/rand"
	"testing"

	"github.com/palantir/conjure-go/v6/dj"
	"github.com/palantir/pkg/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func BenchmarkWriteConstant(b *testing.B) {
	const (
		openBraceString = "{"
	)
	var (
		openBraceBytes = []byte("{")
	)
	buf := bytes.NewBuffer(make([]byte, 1024))
	buf.Reset()
	w := io.Writer(buf)
	b.Run("constantStringCastToBytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = w.Write([]byte(openBraceString))
			buf.Reset()
		}
	})
	b.Run("literalStringCastToBytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = w.Write([]byte("{"))
			buf.Reset()
		}
	})
	b.Run("constantStringWriteString", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = io.WriteString(w, openBraceString)
			buf.Reset()
		}
	})
	b.Run("literalStringBufWriteString", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = buf.WriteString("{")
			buf.Reset()
		}
	})
	b.Run("literalStringWriteString", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = io.WriteString(w, "{")
			buf.Reset()
		}
	})
	b.Run("constantBytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = w.Write(openBraceBytes)
			buf.Reset()
		}
	})
	b.Run("literalBytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = w.Write([]byte{'{'})
			buf.Reset()
		}
	})
	b.Run("dj.WriteLiteral", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = dj.WriteLiteral(w, "{")
			buf.Reset()
		}
	})
	b.Run("dj.WriteOpenObject", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = dj.WriteOpenObject(w)
			buf.Reset()
		}
	})
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	obj := newBenchmarkOuter(5)
	jsonBytes, err := json.Marshal(obj)
	require.NoError(b, err)
	jsonString := string(jsonBytes)
	b.Run("standard library encoding/json", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var out benchmarkOuter
			err := json.Unmarshal(jsonBytes, &out)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("dj direct iterator []byte", func(b *testing.B) {
		b.ReportAllocs()
		for bN := 0; bN < b.N; bN++ {
			var out benchmarkOuter
			value, err := dj.Parse(jsonBytes)
			if err != nil {
				b.Fatal(err)
			}
			err = out.djIteratorUnmarshalJSON(value)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("dj direct iterator string", func(b *testing.B) {
		b.ReportAllocs()
		for bN := 0; bN < b.N; bN++ {
			var out benchmarkOuter
			value, err := dj.Parse(jsonString)
			if err != nil {
				b.Fatal(err)
			}
			err = out.djIteratorUnmarshalJSON(value)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("dj func visitor string", func(b *testing.B) {
		b.ReportAllocs()
		for bN := 0; bN < b.N; bN++ {
			var out benchmarkOuter
			value, err := dj.Parse(jsonString)
			if err != nil {
				b.Fatal(err)
			}
			err = out.djVisitorUnmarshalJSON(value)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("dj func visitor []byte", func(b *testing.B) {
		b.ReportAllocs()
		for bN := 0; bN < b.N; bN++ {
			var out benchmarkOuter
			_, value, err := dj.ParseNext(jsonBytes, 0)
			if err != nil {
				b.Fatal(err)
			}
			err = out.djVisitorUnmarshalJSON(value)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("gjson", func(b *testing.B) {
		b.ReportAllocs()
		for bN := 0; bN < b.N; bN++ {
			var out benchmarkOuter
			if !gjson.ValidBytes(jsonBytes) {
				b.Fatal("invalid json")
			}
			value := gjson.ParseBytes(jsonBytes)
			if !value.IsObject() {
				b.Fatal("expected object")
			}
			var err error
			value.ForEach(func(key, value gjson.Result) bool {
				switch key.Str {
				case "inner":
					if !value.IsArray() {
						err = dj.SyntaxError{}
						return false
					}
					value.ForEach(func(_, value gjson.Result) bool {
						var inner benchmarkInner
						value.ForEach(func(key, value gjson.Result) bool {
							if value.Type != gjson.String {
								err = dj.SyntaxError{}
								return false
							}
							switch key.Str {
							case "field0":
								inner.Field0 = value.String()
							case "field1":
								inner.Field1 = value.String()
							case "field2":
								inner.Field2 = value.String()
							case "field3":
								inner.Field3 = value.String()
							case "field4":
								inner.Field4 = value.String()
							}
							return true
						})
						if err != nil {
							return false
						}
						out.Inners = append(out.Inners, inner)
						return true
					})
				}
				return true
			})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkValidJSON(b *testing.B) {
	runBench := func(b *testing.B, jsonBytes []byte) {
		b.Run("standard library encoding/json", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if !json.Valid(jsonBytes) {
					b.Fatal("invalid json")
				}
			}
		})
		b.Run("gjson", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if ok := gjson.ValidBytes(jsonBytes); !ok {
					b.Fatal("invalid json")
				}
			}
		})
		b.Run("dj", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if err := dj.Valid(jsonBytes); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
	b.Run("benchmark struct", func(b *testing.B) {
		obj := newBenchmarkOuter(5)
		jsonBytes, err := json.Marshal(obj)
		require.NoError(b, err)
		runBench(b, jsonBytes)
	})
	b.Run("basic json", func(b *testing.B) {
		runBench(b, []byte(basicJSON))
	})
}

func BenchmarkMarshalJSON(b *testing.B) {
	obj := newBenchmarkOuter(5)
	b.Run("standard library encoding/json", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := json.Marshal(obj)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("append", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			out, err := obj.djAppendJSON(nil)
			if err != nil {
				b.Fatal(err)
			}
			_ = out
		}
	})
	b.Run("dj *bytes.Buffer", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf := bytes.NewBuffer(nil)
			_, err := obj.djMarshalJSON(buf)
			if err != nil {
				b.Fatal(err)
			}
			var _ = buf.Bytes()
		}
	})
	b.Run("dj empty appender", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf := make([]byte, 0)
			w := dj.NewAppender(&buf)
			_, err := obj.djMarshalJSON(w)
			if err != nil {
				b.Fatal(err)
			}
			var _ []byte = *w
		}
	})
	b.Run("dj discard-prealloc appender", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, err := obj.djMarshalJSON(io.Discard)
			if err != nil {
				b.Fatal(err)
			}
			buf := make([]byte, 0, n)
			w := dj.NewAppender(&buf)
			_, err = obj.djMarshalJSON(w)
			if err != nil {
				b.Fatal(err)
			}
			var _ []byte = *w
		}
	})
}

type benchmarkOuter struct {
	Inners []benchmarkInner `json:"inner"`
}

func newBenchmarkOuter(count int) benchmarkOuter {
	var out benchmarkOuter
	for i := 0; i < count; i++ {
		out.Inners = append(out.Inners, newBenchmarkInner())
	}
	return out
}

func (bo benchmarkOuter) djAppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	out = append(out, "\"inner\":"...)
	out = append(out, '[')
	for i := range bo.Inners {
		out = append(out, '{')
		out = append(out, "\"field0\":"...)
		out = dj.AppendQuotedString(out, bo.Inners[i].Field0)
		out = append(out, ',')
		out = append(out, "\"field1\":"...)
		out = dj.AppendQuotedString(out, bo.Inners[i].Field1)
		out = append(out, ',')
		out = append(out, "\"field2\":"...)
		out = dj.AppendQuotedString(out, bo.Inners[i].Field2)
		out = append(out, ',')
		out = append(out, "\"field3\":"...)
		out = dj.AppendQuotedString(out, bo.Inners[i].Field3)
		out = append(out, ',')
		out = append(out, "\"field4\":"...)
		out = dj.AppendQuotedString(out, bo.Inners[i].Field4)
		out = append(out, '}')
	}
	out = append(out, ']')
	out = append(out, '}')
	return out, nil
}

func (bo benchmarkOuter) djMarshalJSON(w io.Writer) (out int, err error) {
	if n, err := dj.WriteOpenObject(w); err != nil {
		return out, err
	} else {
		out += n
	}
	if n, err := dj.WriteLiteral(w, "\"inner\":"); err != nil {
		return out, err
	} else {
		out += n
	}
	if n, err := dj.WriteOpenArray(w); err != nil {
		return out, err
	} else {
		out += n
	}
	for i := range bo.Inners {
		if n, err := dj.WriteOpenObject(w); err != nil {
			return out, err
		} else {
			out += n
		}
		if n, err := dj.WriteLiteral(w, "\"field0\":"); err != nil {
			return out, err
		} else {
			out += n
		}
		if n, err := dj.WriteString(w, bo.Inners[i].Field0); err != nil {
			return out, err
		} else {
			out += n
		}
		if n, err := dj.WriteCloseObject(w); err != nil {
			return out, err
		} else {
			out += n
		}
	}
	if n, err := dj.WriteCloseArray(w); err != nil {
		return out, err
	} else {
		out += n
	}
	if n, err := dj.WriteCloseObject(w); err != nil {
		return out, err
	} else {
		out += n
	}
	return out, nil
}

func (bo *benchmarkOuter) djVisitorUnmarshalJSON(value dj.Result) error {
	return value.VisitObject(func(key, value dj.Result) error {
		keyString, err := key.String()
		if err != nil {
			return err
		}
		switch keyString {
		case "inner":
			if err := value.VisitArray(func(value dj.Result) error {
				var inner benchmarkInner
				err := inner.djVisitorUnmarshalJSON(value)
				if err != nil {
					return err
				}
				bo.Inners = append(bo.Inners, inner)
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (bo *benchmarkOuter) djIteratorUnmarshalJSON(t dj.Result) error {
	var objectIndex int
	for {
		var key, value dj.Result
		var ok bool
		var err error
		key, value, objectIndex, ok, err = t.NextObjectEntry(objectIndex)
		if err != nil {
			return err
		}
		if !ok {
			break
		}
		keyString, err := key.String()
		if err != nil {
			return err
		}
		switch keyString {
		case "inner":
			arrayIndex := 0
			for {
				var value1 dj.Result
				var ok1 bool
				var err1 error
				value1, arrayIndex, ok1, err1 = value.NextArrayEntry(arrayIndex)
				if err1 != nil {
					return err1
				}
				if !ok1 {
					break
				}
				var inner benchmarkInner
				err1 = inner.djIteratorUnmarshalJSON(value1)
				if err1 != nil {
					return err1
				}
				bo.Inners = append(bo.Inners, inner)
			}
		}
	}
	return nil
}

type benchmarkInner struct {
	Field0 string `json:"field0"`
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
	Field3 string `json:"field3"`
	Field4 string `json:"field4"`
}

func newBenchmarkInner() benchmarkInner {
	return benchmarkInner{
		Field0: newUUID(),
		Field1: newUUID(),
		Field2: newUUID(),
		Field3: newUUID(),
		Field4: newUUID(),
	}
}

func (bi *benchmarkInner) stdlibUnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, *&bi)
}

func (bi *benchmarkInner) djVisitorUnmarshalJSON(value dj.Result) error {
	return value.VisitObject(func(key, value dj.Result) error {
		keyString, err := key.String()
		if err != nil {
			return err
		}
		switch keyString {
		case "field0":
			stringVal, err := value.String()
			if err != nil {
				return err
			}
			bi.Field0 = stringVal
		case "field1":
			stringVal, err := value.String()
			if err != nil {
				return err
			}
			bi.Field1 = stringVal
		case "field2":
			stringVal, err := value.String()
			if err != nil {
				return err
			}
			bi.Field2 = stringVal
		case "field3":
			stringVal, err := value.String()
			if err != nil {
				return err
			}
			bi.Field3 = stringVal
		case "field4":
			stringVal, err := value.String()
			if err != nil {
				return err
			}
			bi.Field4 = stringVal
		default:
		}
		return nil
	})
}

func (bi *benchmarkInner) djIteratorUnmarshalJSON(t dj.Result) error {
	var objectIndex int
	for {
		var key, value dj.Result
		var ok bool
		var err error
		key, value, objectIndex, ok, err = t.NextObjectEntry(objectIndex)
		if err != nil {
			return err
		}
		if !ok {
			break
		}
		keyString, err := key.String()
		if err != nil {
			return err
		}
		switch keyString {
		case "field0":
			stringVal, err := value.String()
			if err != nil {
				return err
			}
			bi.Field0 = stringVal
		case "field1":
			stringVal, err := value.String()
			if err != nil {
				return err
			}
			bi.Field1 = stringVal
		case "field2":
			stringVal, err := value.String()
			if err != nil {
				return err
			}
			bi.Field2 = stringVal
		case "field3":
			stringVal, err := value.String()
			if err != nil {
				return err
			}
			bi.Field3 = stringVal
		case "field4":
			stringVal, err := value.String()
			if err != nil {
				return err
			}
			bi.Field4 = stringVal
		default:
		}
	}
	return nil
}

func newUUID() string {
	bytes := make([]byte, 0, 16)
	bytes = binary.AppendUvarint(bytes, rand.Uint64())
	bytes = binary.AppendUvarint(bytes, rand.Uint64())
	return uuid.UUID(bytes).String()
}

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
