// Copyright (c) 2019 Palantir Technologies. All rights reserved.
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

package visitors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/v5/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v5/conjure/types"
)

func TestExternalTypeFallback(t *testing.T) {

	t.Run("External Fallback", func(t *testing.T) {
		def := spec.ExternalReference{
			ExternalReference: spec.TypeName{
				Name:    "Foo",
				Package: "com.example.foo",
			},
			Fallback: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_STRING)),
		}

		provider := newExternalVisitor(def)

		info := types.NewPkgInfo("", nil)
		typ, err := provider.ParseType(info)

		assert.Equal(t, err, nil)
		assert.Equal(t, types.String, typ)
	})
}

func TestExternalType_IsSafeMarker(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		ref := spec.ExternalReference{
			ExternalReference: spec.TypeName{
				Name:    "Safe",
				Package: "com.palantir.logsafe",
			},
		}
		result, err := IsSpecificConjureType(spec.NewTypeFromExternal(ref), IsSafeMarker)
		require.NoError(t, err)
		assert.True(t, result)
	})
	t.Run("bad name", func(t *testing.T) {
		ref := spec.ExternalReference{
			ExternalReference: spec.TypeName{
				Name:    "Unsafe",
				Package: "com.palantir.logsafe",
			},
		}
		result, err := IsSpecificConjureType(spec.NewTypeFromExternal(ref), IsSafeMarker)
		require.NoError(t, err)
		assert.False(t, result)
	})
	t.Run("bad pkg", func(t *testing.T) {
		ref := spec.ExternalReference{
			ExternalReference: spec.TypeName{
				Name:    "Safe",
				Package: "com.palantir.conjure",
			},
		}
		result, err := IsSpecificConjureType(spec.NewTypeFromExternal(ref), IsSafeMarker)
		require.NoError(t, err)
		assert.False(t, result)
	})
}
