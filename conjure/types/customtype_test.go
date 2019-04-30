// Copyright (c) 2018 Palantir Technologies. All rights reserved.
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
	"testing"

	"github.com/palantir/conjure-go/conjure/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomConjureTypes(t *testing.T) {
	customTypes := types.NewCustomConjureTypes()

	err := customTypes.Add("foo", "", types.Integer)
	require.NoError(t, err)

	val, ok := customTypes.Get("foo")
	require.True(t, ok)
	require.Equal(t, types.CustomConjureType{
		Name:  "foo",
		Pkg:   "",
		Typer: types.Integer,
	}, val)

	err = customTypes.Add("bar", "", types.String)
	require.NoError(t, err)

	val, ok = customTypes.Get("bar")
	require.True(t, ok)
	require.Equal(t, types.CustomConjureType{
		Name:  "bar",
		Pkg:   "",
		Typer: types.String,
	}, val)

	err = customTypes.Add("foo", "", types.Double)
	require.EqualError(t, err, `"foo" has already been defined as a custom Conjure type`)

	val, ok = customTypes.Get("foo")
	require.True(t, ok)
	require.Equal(t, types.CustomConjureType{
		Name:  "foo",
		Pkg:   "",
		Typer: types.Integer,
	}, val)
}

func TestAddCaseInsensitive(t *testing.T) {
	customTypes := types.NewCustomConjureTypes()

	err := customTypes.Add("FooBar", "", types.Integer)
	require.NoError(t, err)

	val, ok := customTypes.Get("fooBAR")
	require.True(t, ok)
	require.Equal(t, types.CustomConjureType{
		Name:  "FooBar",
		Pkg:   "",
		Typer: types.Integer,
	}, val)

	err = customTypes.Add("FOOBAR", "", types.Integer)
	assert.EqualError(t, err, "\"FooBar\" has already been defined as a custom Conjure type")
}
