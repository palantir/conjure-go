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
	"encoding/base64"
	"encoding/json"
	"io"
	"strconv"
	"unsafe"

	"github.com/tidwall/gjson"
)

// JSONWriter is implemented by types that can write themselves as JSON.
type JSONWriter interface {
	WriteJSON(writer io.Writer) (int, error)
}

// stringConst is a string that is known to be constant.
// External callers can only instantiate this type with a string literal.
type stringConst string

var (
	bOpenObject  = []byte{'{'}
	bCloseObject = []byte{'}'}
	bOpenArray   = []byte{'['}
	bCloseArray  = []byte{']'}
	bColon       = []byte{':'}
	bComma       = []byte{','}
	bQuote       = []byte{'"'}
	bNull        = []byte("null")
	bTrue        = []byte("true")
	bFalse       = []byte("false")
)

// WriteOpenObject writes the opening brace of a JSON object.
func WriteOpenObject(w io.Writer) (int, error) {
	return w.Write(bOpenObject)
}

// WriteCloseObject writes the closing brace of a JSON object.
func WriteCloseObject(w io.Writer) (int, error) {
	return w.Write(bCloseObject)
}

// WriteOpenArray writes the opening bracket of a JSON array.
func WriteOpenArray(w io.Writer) (int, error) {
	return w.Write(bOpenArray)
}

// WriteCloseArray writes the closing bracket of a JSON array.
func WriteCloseArray(w io.Writer) (int, error) {
	return w.Write(bCloseArray)
}

// WriteColon writes the colon that separates a JSON object key from its value.
func WriteColon(w io.Writer) (int, error) {
	return w.Write(bColon)
}

// WriteComma writes the comma that separates JSON array and object elements.
func WriteComma(w io.Writer) (int, error) {
	return w.Write(bComma)
}

// WriteQuote writes one double-quote character.
func WriteQuote(w io.Writer) (int, error) {
	return w.Write(bQuote)
}

// WriteNull writes the JSON null value.
func WriteNull(w io.Writer) (int, error) {
	return w.Write(bNull)
}

// WriteLiteral writes a string literal that is known to be constant.
func WriteLiteral(w io.Writer, s stringConst) (int, error) {
	b := unsafe.Slice(unsafe.StringData(string(s)), len(s)) // convert to []byte without allocation
	return w.Write(b)
}

// WriteString quotes, escapes, and writes a JSON string.
func WriteString(w io.Writer, s string) (int, error) {
	if app, ok := w.(*AppendWriter); ok {
		*app = gjson.AppendJSONString(*app, s)
		return len(s), nil
	}
	return w.Write(gjson.AppendJSONString(nil, s))
}

// WriteInt writes an integer as a JSON number.
func WriteInt(w io.Writer, i int64) (int, error) {
	if app, ok := w.(*AppendWriter); ok {
		preLen := len(*app)
		*app = strconv.AppendInt(*app, i, 10)
		return len(*app) - preLen, nil
	}
	return w.Write(strconv.AppendInt(nil, i, 10))
}

// WriteIntString writes an integer as a JSON number, surrounded by quotes.
func WriteIntString(w io.Writer, i int64) (int, error) {
	if _, err := WriteQuote(w); err != nil {
		return 0, err
	}
	n, err := WriteInt(w, i)
	if err != nil {
		return 0, err
	}
	if _, err := WriteQuote(w); err != nil {
		return 0, err
	}
	return n + 2, nil
}

// WriteFloat writes a float64 as a JSON number.
func WriteFloat(w io.Writer, f float64) (int, error) {
	if app, ok := w.(*AppendWriter); ok {
		*app = strconv.AppendFloat(*app, f, 'f', -1, 64)
	}
	return w.Write(strconv.AppendFloat(nil, f, 'f', -1, 64))
}

// WriteFloatString writes a float64 as a JSON number, surrounded by quotes.
func WriteFloatString(w io.Writer, f float64) (int, error) {
	if _, err := WriteQuote(w); err != nil {
		return 0, err
	}
	n, err := WriteFloat(w, f)
	if err != nil {
		return 0, err
	}
	if _, err := WriteQuote(w); err != nil {
		return 0, err
	}
	return n + 2, nil
}

// WriteBool writes a JSON boolean.
func WriteBool(w io.Writer, b bool) (int, error) {
	if app, ok := w.(*AppendWriter); ok {
		if b {
			*app = append(*app, bTrue...)
			return 4, nil
		}
		*app = append(*app, bFalse...)
		return 5, nil
	}
	if b {
		return w.Write(bTrue)
	}
	return w.Write(bFalse)
}

// WriteBoolString writes a JSON boolean, surrounded by quotes.
func WriteBoolString(w io.Writer, b bool) (int, error) {
	if _, err := WriteQuote(w); err != nil {
		return 0, err
	}
	n, err := WriteBool(w, b)
	if err != nil {
		return 0, err
	}
	if _, err := WriteQuote(w); err != nil {
		return 0, err
	}
	return n + 2, nil
}

// WriteBase64 writes a byte slice as a JSON string containing the base64 encoding of the bytes.
func WriteBase64(w io.Writer, data []byte) (int, error) {
	if w == io.Discard {
		return 2 + base64.StdEncoding.EncodedLen(len(data)), nil
	}
	if _, err := WriteQuote(w); err != nil {
		return 0, err
	}
	// todo: can avoid this allocation if we can get to the raw bytes of the writer
	b64out := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(b64out, data)
	n, err := w.Write(b64out)
	if err != nil {
		return 0, err
	}
	if _, err := WriteQuote(w); err != nil {
		return 0, err
	}
	return n + 2, nil
}

// WriteObject writes a JSON representation of the given object.
func WriteObject(w io.Writer, obj any) (int, error) {
	switch v := obj.(type) {
	case JSONWriter:
		return v.WriteJSON(w)
	case json.Marshaler:
		jsonBytes, err := v.MarshalJSON()
		if err != nil {
			return 0, err
		}
		return w.Write(jsonBytes)
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return 0, err
		}
		return w.Write(jsonBytes)
	}
}

// NewAppender returns a Writer that appends to the given byte slice.
func NewAppender(buf *[]byte) *AppendWriter {
	return (*AppendWriter)(buf)
}

// AppendWriter is an io.Writer that appends to a byte slice.
// Some Write* methods check for this type to avoid allocating intermediate buffers.
type AppendWriter []byte

func (w *AppendWriter) Write(p []byte) (int, error) {
	*w = append(*w, p...)
	return len(p), nil
}

func (w *AppendWriter) WriteString(s string) (int, error) {
	*w = append(*w, s...)
	return len(s), nil
}

func (w *AppendWriter) String() string {
	return string(*w)
}
