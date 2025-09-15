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

package types_test

import (
	"encoding/base64"
	"testing"

	"github.com/palantir/conjure-go/v6/cj/types"
	"github.com/palantir/pkg/binary"
	"github.com/palantir/pkg/uuid"
)

func TestBinary(t *testing.T) {
	for name, test := range map[string]typeTest{
		"nil": typeTestCase[[]byte, types.Binary[[]byte], types.Binary[[]byte]]{
			Value: nil, JSON: `""`, SkipTestUnmarshal: true,
		},
		"empty_marshal": typeTestCase[[]byte, types.Binary[[]byte], types.Binary[[]byte]]{
			Value: []byte{}, JSON: `""`,
		},
		"null_unmarshal": typeTestCase[[]byte, types.Binary[[]byte], types.Binary[[]byte]]{
			JSON: "null", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "KindMismatchError at 0: want json string, got null",
		},
		"bytes": typeTestCase[[]byte, types.Binary[[]byte], types.Binary[[]byte]]{
			Value: []byte("hello 👋"), JSON: "\"aGVsbG8g8J+Riw==\"",
		},
		"invalid base64": typeTestCase[[]byte, types.Binary[[]byte], types.Binary[[]byte]]{
			JSON: "\"not_base64!@#\"", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "InvalidValueError at 15: illegal base64 data at input byte 3",
		},
		"not a string": typeTestCase[[]byte, types.Binary[[]byte], types.Binary[[]byte]]{
			JSON: "123", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "KindMismatchError at 0: want json string, got number",
		},
		"map": typeTestCase[map[binary.Binary]string, types.OrderedMapMarshaler[map[binary.Binary]string, binary.Binary, string, types.BinaryMapKey[binary.Binary], types.String[string]], types.MapUnmarshaler[map[binary.Binary]string, binary.Binary, string, types.BinaryMapKey[binary.Binary], types.String[string]]]{
			Value: map[binary.Binary]string{
				binary.Binary(base64.StdEncoding.EncodeToString([]byte("a"))): "a",
				binary.Binary(base64.StdEncoding.EncodeToString([]byte("b"))): "b",
				binary.Binary(base64.StdEncoding.EncodeToString([]byte("c"))): "c",
			},
			JSON: `{"YQ==":"a","Yg==":"b","Yw==":"c"}`,
		},
		"BinaryMarshaler": typeTestCase[uuid.UUID, types.BinaryMarshaler[uuid.UUID], types.BinaryUnmarshaler[*uuid.UUID]]{
			Value: must(uuid.ParseUUID("10101010-1010-1010-1010-101010101010")), JSON: "\"EBAQEBAQEBAQEBAQEBAQEA==\"",
		},
		"BinaryMarshaler map": typeTestCase[map[uuid.UUID]string, types.ComparableMapMarshaler[map[uuid.UUID]string, uuid.UUID, string, types.BinaryMarshaler[uuid.UUID], types.String[string]], types.MapUnmarshaler[map[uuid.UUID]string, uuid.UUID, string, types.BinaryUnmarshaler[*uuid.UUID], types.String[string]]]{
			Value: map[uuid.UUID]string{
				must(uuid.ParseUUID("10101010-1010-1010-1010-101010101010")): "foo",
				must(uuid.ParseUUID("10202020-2020-2020-2020-202020202020")): "bar",
			},
			JSON: `{"EBAQEBAQEBAQEBAQEBAQEA==":"foo","ECAgICAgICAgICAgICAgIA==":"bar"}`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}
