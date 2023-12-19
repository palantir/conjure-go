package dj_test

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"testing"

	"github.com/palantir/conjure-go/v6/dj"
	"github.com/palantir/pkg/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func BenchmarkUnmarshalJSON(b *testing.B) {
	obj := newBenchmarkOuter(b, 5)
	jsonBytes, err := json.Marshal(obj)
	require.NoError(b, err)
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
	b.Run("dj direct iterator", func(b *testing.B) {
		b.ReportAllocs()
		for bN := 0; bN < b.N; bN++ {
			var out benchmarkOuter
			value, err := dj.Parse(jsonBytes)
			if err != nil {
				b.Fatal(err)
			}
			err = out.djIteratorUnmarshalJSON(value, 0)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("dj func visitor", func(b *testing.B) {
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
		obj := newBenchmarkOuter(b, 5)
		jsonBytes, err := json.Marshal(obj)
		require.NoError(b, err)
		runBench(b, jsonBytes)
	})
	b.Run("basic json", func(b *testing.B) {
		runBench(b, []byte(basicJSON))
	})
}

func BenchmarkMarshalJSON(b *testing.B) {
	obj := newBenchmarkOuter(b, 5)
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

func newBenchmarkOuter(tb testing.TB, count int) benchmarkOuter {
	var out benchmarkOuter
	for i := 0; i < count; i++ {
		out.Inners = append(out.Inners, newBenchmarkInner(tb))
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
		out = dj.AppendJSONString(out, bo.Inners[i].Field0)
		out = append(out, ',')
		out = append(out, "\"field1\":"...)
		out = dj.AppendJSONString(out, bo.Inners[i].Field1)
		out = append(out, ',')
		out = append(out, "\"field2\":"...)
		out = dj.AppendJSONString(out, bo.Inners[i].Field2)
		out = append(out, ',')
		out = append(out, "\"field3\":"...)
		out = dj.AppendJSONString(out, bo.Inners[i].Field3)
		out = append(out, ',')
		out = append(out, "\"field4\":"...)
		out = dj.AppendJSONString(out, bo.Inners[i].Field4)
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
		switch key.Str {
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

func (bo *benchmarkOuter) djIteratorUnmarshalJSON(t dj.Result, i int) error {
	iter, i, err := t.ObjectIterator(i)
	if err != nil {
		return err
	}
	for iter.HasNext(t, i) {
		var key, value dj.Result
		key, value, i, err = iter.Next(t, i)
		if err != nil {
			return err
		}
		switch key.Str {
		case "inner":
			iter1, i1, err := value.ArrayIterator(i)
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
				err = inner.djIteratorUnmarshalJSON(value1, i1)
				if err != nil {
					return err
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

func newBenchmarkInner(tb testing.TB) benchmarkInner {
	return benchmarkInner{
		Field0: newUUID(tb),
		Field1: newUUID(tb),
		Field2: newUUID(tb),
		Field3: newUUID(tb),
		Field4: newUUID(tb),
	}
}

func (bi *benchmarkInner) stdlibUnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, bi)
}

func (bi *benchmarkInner) djVisitorUnmarshalJSON(value dj.Result) error {
	return value.VisitObject(func(key, value dj.Result) error {
		switch key.Str {
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

func (bi *benchmarkInner) djIteratorUnmarshalJSON(t dj.Result, i int) error {
	iter, i, err := t.ObjectIterator(i)
	if err != nil {
		return err
	}
	for iter.HasNext(t, i) {
		var key, value dj.Result
		key, value, i, err = iter.Next(t, i)
		if err != nil {
			return err
		}
		switch key.Str {
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

func newUUID(tb testing.TB) string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	require.NoError(tb, err)
	return uuid.UUID(bytes).String()
}
