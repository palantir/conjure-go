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
	"encoding/json"
	"io"
	"math/rand"
	"testing"

	"github.com/palantir/conjure-go/v6/dj"
	"github.com/palantir/pkg/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func BenchmarkUnmarshalJSON_iter(b *testing.B) {
	obj := newBenchmarkOuter(5)
	jsonBytes, err := json.Marshal(obj)
	require.NoError(b, err)
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
		if len(out.Inners) != 5 {
			b.Fatal("invalid length")
		}
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	obj := newBenchmarkOuter(5)
	jsonBytes, err := json.Marshal(obj)
	require.NoError(b, err)
	jsonString := string(jsonBytes)
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
	b.Run("dj func visitor []byte", func(b *testing.B) {
		b.ReportAllocs()
		for bN := 0; bN < b.N; bN++ {
			var out benchmarkOuter
			value, err := dj.Parse(jsonBytes)
			if err != nil {
				b.Fatal(err)
			}
			err = out.djVisitorUnmarshalJSON(value)
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
	b.Run("gjson []byte", func(b *testing.B) {
		b.ReportAllocs()
		for bN := 0; bN < b.N; bN++ {
			if !gjson.ValidBytes(jsonBytes) {
				b.Fatal("invalid json")
			}
			value := gjson.ParseBytes(jsonBytes)
			if !value.IsObject() {
				b.Fatal("expected object")
			}
			var out benchmarkOuter
			if err := out.gjsonUnmarshalJSON(value); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("gjson string", func(b *testing.B) {
		b.ReportAllocs()
		for bN := 0; bN < b.N; bN++ {
			if !gjson.Valid(jsonString) {
				b.Fatal("invalid json")
			}
			value := gjson.Parse(jsonString)
			if !value.IsObject() {
				b.Fatal("expected object")
			}
			var out benchmarkOuter
			if err := out.gjsonUnmarshalJSON(value); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("encoding/json", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var out benchmarkOuter
			err := json.Unmarshal(jsonBytes, &out)
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
	iter, i, err := t.ObjectIterator(0)
	if err != nil {
		return err
	}
	for iter.HasNext(t, i) {
		var key, value dj.Result
		key, value, i, err = iter.Next(t, i)
		if err != nil {
			return err
		}
		keyString, err := key.String()
		if err != nil {
			return err
		}
		switch keyString {
		case "inner":
			iter1, i1, err := value.ArrayIterator(0)
			if err != nil {
				return err
			}
			for iter1.HasNext(value, i1) {
				var value1 dj.Result
				value1, i1, err = iter1.Next(value, i1)
				if err != nil {
					return err
				}
				var inner benchmarkInner
				err = inner.djIteratorUnmarshalJSON(value1)
				if err != nil {
					return err
				}
				bo.Inners = append(bo.Inners, inner)
			}
		}
	}
	return nil
}

func (bo *benchmarkOuter) gjsonUnmarshalJSON(value gjson.Result) error {
	var out benchmarkOuter
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
		return err
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
	iter, i, err := t.ObjectIterator(0)
	if err != nil {
		return err
	}
	for iter.HasNext(t, i) {
		var key, value dj.Result
		key, value, i, err = iter.Next(t, i)
		if err != nil {
			return err
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
