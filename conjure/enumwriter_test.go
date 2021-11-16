// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package conjure

import (
	"testing"

	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/stretchr/testify/assert"
)

func Test_EmptyEnum(t *testing.T) {
	for _, tc := range []struct {
		name     string
		typeName string
		values   []*types.Field
		expected string
	}{
		{
			name:     "Empty list",
			typeName: "EmptyValuesEnum",
			values:   nil,
			expected: `// IsUnknown returns false for all known variants of EmptyValuesEnum and true otherwise.
func (e EmptyValuesEnum) IsUnknown() bool {
	switch e.val {
	default:
		return false
	}
	return true
}`,
		},
		{
			name:     "String values",
			typeName: "Enum",
			values: []*types.Field{
				{
					Name: "SATURDAY",
					Type: types.String{},
				},
				{
					Name: "SUNDAY",
					Type: types.String{},
				},
			},
			expected: `// IsUnknown returns false for all known variants of Enum and true otherwise.
func (e Enum) IsUnknown() bool {
	switch e.val {
	case Enum_SATURDAY, Enum_SUNDAY:
		return false
	}
	return true
}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stmt := astForEnumIsUnknown(tc.typeName, tc.values)
			assert.Equal(t, tc.expected, stmt.GoString())
		})
	}

}
