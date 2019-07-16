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

package objects_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/palantir/pkg/rid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/palantir/conjure-go/integration_test/testgenerated/objects/api"
)

type FuncType int

func (t FuncType) String() string {
	switch t {
	case JSON:
		return "JSON"
	case YAML:
		return "YAML"
	default:
		return strconv.Itoa(int(t))
	}
}

const (
	JSON FuncType = iota
	YAML
)

var unmarshalFuncs = []func([]byte, interface{}) (err error){
	JSON: json.Unmarshal,
	YAML: yaml.Unmarshal,
}

func TestStrOptionalNil(t *testing.T) {
	union := api.NewExampleUnionFromStrOptional(nil)
	bytes, err := union.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `{"type":"strOptional","strOptional":null}`, string(bytes))

	union = api.ExampleUnion{}
	if err := json.Unmarshal(bytes, &union); err != nil {
		require.NoError(t, err)
	}
	v := &visitor{}
	err = union.Accept(v)
	require.NoError(t, err)
	assert.Nil(t, v.visitedStrOptional)

	bytes, err = union.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `{"type":"strOptional","strOptional":null}`, string(bytes))
}

func TestStrOptionalNonNil(t *testing.T) {
	strVal := "hello"
	union := api.NewExampleUnionFromStrOptional(&strVal)
	bytes, err := union.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `{"type":"strOptional","strOptional":"hello"}`, string(bytes))

	union = api.ExampleUnion{}
	if err := json.Unmarshal(bytes, &union); err != nil {
		require.NoError(t, err)
	}
	v := &visitor{}
	err = union.Accept(v)
	require.NoError(t, err)
	assert.Equal(t, "hello", *v.visitedStrOptional)

	bytes, err = union.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `{"type":"strOptional","strOptional":"hello"}`, string(bytes))
}

func TestRidAliasString(t *testing.T) {
	parsedRID, err := rid.ParseRID("ri.a1p2p3.south-west.data-set.my-hello_WORLD-123")
	require.NoError(t, err)

	ridAlias := api.RidAlias(parsedRID)
	assert.Equal(t, "ri.a1p2p3.south-west.data-set.my-hello_WORLD-123", fmt.Sprint(ridAlias))
}

func TestMarshal(t *testing.T) {
	for i, tc := range []struct {
		obj      interface{}
		wantJSON string
		wantYAML string
	}{
		{api.Collections{}, `{"mapVar":{},"listVar":[],"multiDim":[]}`, "mapVar: {}\nlistVar: []\nmultiDim: []\n"},
		{&api.Collections{}, `{"mapVar":{},"listVar":[],"multiDim":[]}`, "mapVar: {}\nlistVar: []\nmultiDim: []\n"},
		{api.Compound{}, `{"obj":{"mapVar":{},"listVar":[],"multiDim":[]}}`, "obj:\n  mapVar: {}\n  listVar: []\n  multiDim: []\n"},
	} {
		bytes, err := json.Marshal(tc.obj)
		require.NoError(t, err)
		assert.Equal(t, tc.wantJSON, string(bytes), "Case %d (JSON)", i)

		bytes, err = yaml.Marshal(tc.obj)
		require.NoError(t, err)
		assert.Equal(t, tc.wantYAML, string(bytes), "Case %d (YAML)", i)
	}
}

func TestUnmarshal(t *testing.T) {
	for idx, unmarshalFunc := range unmarshalFuncs {
		var test1 api.Collections
		err := unmarshalFunc([]byte(`{}`), &test1)
		require.NoError(t, err)
		assert.Equal(t, api.Collections{
			MapVar:   make(map[string][]int, 0),
			ListVar:  make([]string, 0),
			MultiDim: make([][]map[string]int, 0),
		}, test1)
		assert.NotNil(t, test1.MapVar, "Case %s", FuncType(idx).String())
		assert.NotNil(t, test1.ListVar, "Case %s", FuncType(idx).String())

		var test2 api.Compound
		err = unmarshalFunc([]byte(`{"obj":{}}`), &test2)
		require.NoError(t, err)
		assert.Equal(t, api.Compound{
			Obj: api.Collections{
				MapVar:   make(map[string][]int, 0),
				ListVar:  make([]string, 0),
				MultiDim: make([][]map[string]int, 0),
			},
		}, test2)
		assert.NotNil(t, test2.Obj.MapVar, "Case %s", FuncType(idx).String())
		assert.NotNil(t, test2.Obj.ListVar, "Case %s", FuncType(idx).String())
		assert.NotNil(t, test2.Obj.MultiDim, "Case %s", FuncType(idx).String())
	}
}

func TestUnions(t *testing.T) {
	for i, tc := range []struct {
		creator       func() api.ExampleUnion
		wantJSONBytes []byte
		wantYAMLBytes []byte
		wantVisitor   visitor
	}{
		{
			func() api.ExampleUnion {
				return api.NewExampleUnionFromOther(5)
			},
			[]byte(`{"type":"other","other":5}`),
			[]byte("type: other\nother: 5\n"),
			visitor{
				visitedInt: 5,
			},
		},
		{
			func() api.ExampleUnion {
				return api.NewExampleUnionFromStr("foo")
			},
			[]byte(`{"type":"str","str":"foo"}`),
			[]byte("type: str\nstr: foo\n"),
			visitor{
				visitedStr: "foo",
			},
		},
	} {
		union := tc.creator()
		bytes, err := json.Marshal(union)
		require.NoError(t, err)
		assert.Equal(t, string(tc.wantJSONBytes), string(bytes), "Case %d (JSON)", i)
		v := &visitor{}
		err = union.Accept(v)
		require.NoError(t, err, "Case %d (JSON)", i)
		assert.Equal(t, tc.wantVisitor, *v, "Case %d (JSON)", i)

		bytes, err = yaml.Marshal(union)
		require.NoError(t, err)
		assert.Equal(t, string(tc.wantYAMLBytes), string(bytes), "Case %d (YAML)", i)
		v = &visitor{}
		err = union.Accept(v)
		require.NoError(t, err, "Case %d (YAML)", i)
		assert.Equal(t, tc.wantVisitor, *v, "Case %d (YAML)", i)
	}
}

func TestUnknownUnions(t *testing.T) {
	for idx, unmarshalFunc := range unmarshalFuncs {
		var unknownUnion *api.ExampleUnion
		err := unmarshalFunc([]byte(`{"type":"notAValidType","notAValidType":"foo"}`), &unknownUnion)
		require.NoError(t, err, "Case %s", FuncType(idx).String())
		v := &visitor{}
		err = unknownUnion.Accept(v)
		require.NoError(t, err, "Case %s", FuncType(idx).String())
		assert.Equal(t, "notAValidType", v.unknownType, "Case %s", FuncType(idx).String())
	}
}

type visitor struct {
	visitedStr         string
	visitedStrOptional *string
	visitedInt         int
	unknownType        string
}

func (v *visitor) VisitStr(val string) error {
	v.visitedStr = val
	return nil
}

func (v *visitor) VisitStrOptional(val *string) error {
	v.visitedStrOptional = val
	return nil
}

func (v *visitor) VisitOther(val int) error {
	v.visitedInt = val
	return nil
}

func (v *visitor) VisitUnknown(typeName string) error {
	v.unknownType = typeName
	return nil
}
