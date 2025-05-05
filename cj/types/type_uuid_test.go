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
	"fmt"
	"testing"

	"github.com/palantir/conjure-go/v6/cj/types"
	"github.com/palantir/pkg/uuid"
)

func TestUUID(t *testing.T) {
	for i, test := range []typeTest{
		typeTestCase[[16]byte, types.UUID[[16]byte], types.UUID[[16]byte]]{
			Value: uuid.UUID{},
			JSON:  "\"00000000-0000-0000-0000-000000000000\"",
		},
		typeTestCase[uuid.UUID, types.UUID[uuid.UUID], types.UUID[uuid.UUID]]{
			Value: must(uuid.ParseUUID("10101010-1010-1010-1010-101010101010")),
			JSON:  "\"10101010-1010-1010-1010-101010101010\"",
		},
		typeTestCase[uuid.UUID, types.UUID[uuid.UUID], types.UUID[uuid.UUID]]{
			Value:                must(uuid.ParseUUID("10101010-1010-1010-1010-101010101010")),
			JSON:                 "null",
			SkipTestMarshal:      true,
			ErrUnmarshalJSONFrom: "KindMismatchError at 4: want json string, got null",
		},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}
