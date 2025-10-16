// Copyright (c) 2025 Palantir Technologies. All rights reserved.
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

package cj_test

import (
	"encoding/base64"
	"testing"

	"github.com/palantir/conjure-go/v6/cj"
	"github.com/palantir/pkg/binary"
)

func TestBinary(t *testing.T) {
	tests := []struct {
		Name string
		Test typeTest
	}{
		{
			Name: "nil",
			Test: typeTestCase[[]byte, cj.Binary[[]byte], cj.Binary[[]byte]]{
				Value: nil, JSON: `""`, SkipTestUnmarshal: true,
			},
		},
		{
			Name: "empty_marshal",
			Test: typeTestCase[[]byte, cj.Binary[[]byte], cj.Binary[[]byte]]{
				Value: []byte{}, JSON: `""`,
			},
		},
		{
			Name: "null_unmarshal",
			Test: typeTestCase[[]byte, cj.Binary[[]byte], cj.Binary[[]byte]]{
				JSON: "null", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "KindMismatchError at offset 0: want json string, got null",
			},
		},
		{
			Name: "bytes",
			Test: typeTestCase[[]byte, cj.Binary[[]byte], cj.Binary[[]byte]]{
				Value: []byte("hello 👋"), JSON: "\"aGVsbG8g8J+Riw==\"",
			},
		},
		{
			Name: "invalid base64",
			Test: typeTestCase[[]byte, cj.Binary[[]byte], cj.Binary[[]byte]]{
				JSON: "\"not_base64!@#\"", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "InvalidValueError at offset 15: illegal base64 data at input byte 3",
			},
		},
		{
			Name: "not a string",
			Test: typeTestCase[[]byte, cj.Binary[[]byte], cj.Binary[[]byte]]{
				JSON: "123", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "KindMismatchError at offset 0: want json string, got number",
			},
		},
		{
			Name: "map",
			Test: typeTestCase[map[binary.Binary]string, cj.OrderedMapMarshaler[map[binary.Binary]string, binary.Binary, string, cj.BinaryMapKey[binary.Binary], cj.String[string]], cj.MapUnmarshaler[map[binary.Binary]string, binary.Binary, string, cj.BinaryMapKey[binary.Binary], cj.String[string]]]{
				Value: map[binary.Binary]string{
					binary.Binary(base64.StdEncoding.EncodeToString([]byte("a"))): "a",
					binary.Binary(base64.StdEncoding.EncodeToString([]byte("b"))): "b",
					binary.Binary(base64.StdEncoding.EncodeToString([]byte("c"))): "c",
				},
				JSON: `{"YQ==":"a","Yg==":"b","Yw==":"c"}`,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			t.Run("Marshal", tc.Test.TestMarshal)
			t.Run("Unmarshal", tc.Test.TestUnmarshal)
		})
	}
}
