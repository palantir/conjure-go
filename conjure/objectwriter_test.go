// Copyright (c) 2024 Palantir Technologies. All rights reserved.
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
	"bytes"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/stretchr/testify/assert"
)

func TestObjectWriter(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   *types.ObjectType
		Out  string
	}{
		{
			name: "Object without collections",
			in: &types.ObjectType{
				Name: "User",
				Fields: []*types.Field{
					{
						Name: "UserName",
						Type: types.String{},
					},
				},
			},
			Out: `package testpkg

import (
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

type User struct {
	UserName string ` + "`json:\"UserName\"`" + `
}

func (o User) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (o *User) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
		{
			name: "Object with collections",
			in: &types.ObjectType{
				Name: "User",
				Fields: []*types.Field{
					{
						Name: "UserName",
						Type: types.String{},
					},
					{
						Name: "UserAliases",
						Type: &types.List{Item: &types.AliasType{
							Item: types.String{},
							Name: "UserAlias",
						}},
					},
				},
			},
			Out: `package testpkg

import (
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

type User struct {
	UserName    string      ` + "`json:\"UserName\"`" + `
	UserAliases []UserAlias ` + "`json:\"UserAliases\"`" + `
}

func (o User) MarshalJSON() ([]byte, error) {
	if o.UserAliases == nil {
		o.UserAliases = make([]UserAlias, 0)
	}
	type _tmpUser User
	return safejson.Marshal(_tmpUser(o))
}
func (o *User) UnmarshalJSON(data []byte) error {
	type _tmpUser User
	var rawUser _tmpUser
	if err := safejson.Unmarshal(data, &rawUser); err != nil {
		return err
	}
	if rawUser.UserAliases == nil {
		rawUser.UserAliases = make([]UserAlias, 0)
	}
	*o = User(rawUser)
	return nil
}
func (o User) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (o *User) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
		{
			name: "Object with safety annotations",
			in: &types.ObjectType{
				Name: "SafetyTestObject",
				Fields: []*types.Field{
					{
						Name: "safeField",
						Type: types.NewAliasTypeWithSafety("SafeString", types.String{},
							func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_SAFE); return &s }()),
					},
					{
						Name: "unsafeField",
						Type: types.NewAliasTypeWithSafety("UnsafeString", types.String{},
							func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_UNSAFE); return &s }()),
					},
					{
						Name: "doNotLogField",
						Type: types.NewAliasTypeWithSafety("SecretString", types.String{},
							func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_DO_NOT_LOG); return &s }()),
					},
					{
						Name: "normalField",
						Type: types.String{},
					},
					{
						Name:   "fieldLevelSafeString",
						Type:   types.String{},
						Safety: func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_SAFE); return &s }(),
					},
					{
						Name:   "fieldLevelUnsafeString",
						Type:   types.String{},
						Safety: func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_UNSAFE); return &s }(),
					},
				},
			},
			Out: `package testpkg

import (
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

// safelogging:@DoNotLog
type SafetyTestObject struct {
	SafeField              SafeString   ` + "`json:\"safeField\" safelogging:\"@Safe\"`" + `
	UnsafeField            UnsafeString ` + "`json:\"unsafeField\" safelogging:\"@Unsafe\"`" + `
	DoNotLogField          SecretString ` + "`json:\"doNotLogField\" safelogging:\"@DoNotLog\"`" + `
	NormalField            string       ` + "`json:\"normalField\"`" + `
	FieldLevelSafeString   string       ` + "`json:\"fieldLevelSafeString\" safelogging:\"@Safe\"`" + `
	FieldLevelUnsafeString string       ` + "`json:\"fieldLevelUnsafeString\" safelogging:\"@Unsafe\"`" + `
}

func (o SafetyTestObject) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (o *SafetyTestObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
		{
			name: "Object with only safe fields",
			in: &types.ObjectType{
				Name: "OnlySafeObject",
				Fields: []*types.Field{
					{
						Name: "safeField1",
						Type: types.NewAliasTypeWithSafety("SafeString", types.String{},
							func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_SAFE); return &s }()),
					},
					{
						Name: "safeField2",
						Type: types.NewAliasTypeWithSafety("SafeString", types.String{},
							func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_SAFE); return &s }()),
					},
					{
						Name: "normalField",
						Type: types.String{},
					},
				},
			},
			Out: `package testpkg

import (
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

type OnlySafeObject struct {
	SafeField1  SafeString ` + "`json:\"safeField1\" safelogging:\"@Safe\"`" + `
	SafeField2  SafeString ` + "`json:\"safeField2\" safelogging:\"@Safe\"`" + `
	NormalField string     ` + "`json:\"normalField\"`" + `
}

func (o OnlySafeObject) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (o *OnlySafeObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
		{
			name: "Object with unsafe but no do-not-log fields",
			in: &types.ObjectType{
				Name: "UnsafeObject",
				Fields: []*types.Field{
					{
						Name: "safeField",
						Type: types.NewAliasTypeWithSafety("SafeString", types.String{},
							func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_SAFE); return &s }()),
					},
					{
						Name: "unsafeField",
						Type: types.NewAliasTypeWithSafety("UnsafeString", types.String{},
							func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_UNSAFE); return &s }()),
					},
				},
			},
			Out: `package testpkg

import (
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

// safelogging:@Unsafe
type UnsafeObject struct {
	SafeField   SafeString   ` + "`json:\"safeField\" safelogging:\"@Safe\"`" + `
	UnsafeField UnsafeString ` + "`json:\"unsafeField\" safelogging:\"@Unsafe\"`" + `
}

func (o UnsafeObject) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (o *UnsafeObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
		{
			name: "Nested objects with different safety levels",
			in: &types.ObjectType{
				Name: "RootObject",
				Fields: []*types.Field{
					{
						Name: "foo",
						Type: &types.ObjectType{
							Name: "Foo",
							Fields: []*types.Field{
								{
									Name: "safeField1",
									Type: types.NewAliasTypeWithSafety("SafeString", types.String{},
										func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_SAFE); return &s }()),
								},
								{
									Name: "safeField2",
									Type: types.NewAliasTypeWithSafety("SafeString", types.String{},
										func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_SAFE); return &s }()),
								},
							},
						},
					},
					{
						Name: "bar",
						Type: &types.ObjectType{
							Name: "Bar",
							Fields: []*types.Field{
								{
									Name: "safeField",
									Type: types.NewAliasTypeWithSafety("SafeString", types.String{},
										func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_SAFE); return &s }()),
								},
								{
									Name: "unsafeField",
									Type: types.NewAliasTypeWithSafety("UnsafeString", types.String{},
										func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_UNSAFE); return &s }()),
								},
							},
						},
					},
					{
						Name: "qux",
						Type: &types.ObjectType{
							Name: "Qux",
							Fields: []*types.Field{
								{
									Name: "safeField",
									Type: types.NewAliasTypeWithSafety("SafeString", types.String{},
										func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_SAFE); return &s }()),
								},
								{
									Name: "unannotatedField",
									Type: types.String{}, // No safety annotation
								},
							},
						},
					},
				},
			},
			Out: `package testpkg

import (
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

// safelogging:@Unsafe
type RootObject struct {
	Foo Foo ` + "`json:\"foo\" safelogging:\"@Safe\"`" + `
	Bar Bar ` + "`json:\"bar\" safelogging:\"@Unsafe\"`" + `
	Qux Qux ` + "`json:\"qux\"`" + `
}

func (o RootObject) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (o *RootObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
		{
			name: "Object with mix of safe and unannotated fields",
			in: &types.ObjectType{
				Name: "MixedSafetyObject",
				Fields: []*types.Field{
					{
						Name: "safeField",
						Type: types.NewAliasTypeWithSafety("SafeString", types.String{},
							func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_SAFE); return &s }()),
					},
					{
						Name: "unannotatedField",
						Type: types.String{}, // No safety annotation
					},
				},
			},
			Out: `package testpkg

import (
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

type MixedSafetyObject struct {
	SafeField        SafeString ` + "`json:\"safeField\" safelogging:\"@Safe\"`" + `
	UnannotatedField string     ` + "`json:\"unannotatedField\"`" + `
}

func (o MixedSafetyObject) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (o *MixedSafetyObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
		{
			name: "Empty struct should have unknown safety",
			in: &types.ObjectType{
				Name:   "EmptyStruct",
				Fields: []*types.Field{}, // No fields
			},
			Out: `package testpkg

import (
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

type EmptyStruct struct{}

func (o EmptyStruct) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (o *EmptyStruct) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
		{
			name: "Object with docs and safety annotation",
			in: &types.ObjectType{
				Name: "DocumentedSafeObject",
				Docs: types.Docs("DocumentedSafeObject represents a safe object with documentation."),
				Fields: []*types.Field{
					{
						Name: "safeField",
						Type: types.NewAliasTypeWithSafety("SafeString", types.String{},
							func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_SAFE); return &s }()),
					},
				},
			},
			Out: `package testpkg

import (
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

// DocumentedSafeObject represents a safe object with documentation.
// safelogging:@Safe
type DocumentedSafeObject struct {
	SafeField SafeString ` + "`json:\"safeField\" safelogging:\"@Safe\"`" + `
}

func (o DocumentedSafeObject) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (o *DocumentedSafeObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
		{
			name: "Object with docs but no safety annotation",
			in: &types.ObjectType{
				Name: "DocumentedMixedObject",
				Docs: types.Docs("DocumentedMixedObject has mixed field safety levels."),
				Fields: []*types.Field{
					{
						Name: "safeField",
						Type: types.NewAliasTypeWithSafety("SafeString", types.String{},
							func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_SAFE); return &s }()),
					},
					{
						Name: "normalField",
						Type: types.String{}, // No safety annotation - makes struct unknown
					},
				},
			},
			Out: `package testpkg

import (
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

// DocumentedMixedObject has mixed field safety levels.
type DocumentedMixedObject struct {
	SafeField   SafeString ` + "`json:\"safeField\" safelogging:\"@Safe\"`" + `
	NormalField string     ` + "`json:\"normalField\"`" + `
}

func (o DocumentedMixedObject) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (o *DocumentedMixedObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
		{
			name: "Object with External type with safety",
			in: &types.ObjectType{
				Name: "MyObject",
				Fields: []*types.Field{
					{
						Name: "externalField",
						Type: &types.External{
							Spec: spec.TypeName{
								Name:    "com/palantir/apollo/deployment:ApolloEnvironmentId",
								Package: "github.com/palantir/apollo-deployment-api",
							},
							Fallback: types.String{},
						},
					},
				},
			},
			Out: `package testpkg

import (
	deployment "github.com/palantir/apollo-deployment-api.com/palantir/apollo/deployment"
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

type MyObject struct {
	ExternalField deployment.ApolloEnvironmentId ` + "`json:\"externalField\"`" + `
}

func (o MyObject) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (o *MyObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
		{
			name: "Object with External type with SAFE safety",
			in: &types.ObjectType{
				Name: "MyObject",
				Fields: []*types.Field{
					{
						Name: "safeExternalField",
						Type: types.NewExternalWithSafety(
							spec.TypeName{
								Name:    "com/palantir/apollo/deployment:ApolloEnvironmentId",
								Package: "github.com/palantir/apollo-deployment-api",
							},
							types.String{},
							func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_SAFE); return &s }(),
						),
					},
				},
			},
			Out: `package testpkg

import (
	deployment "github.com/palantir/apollo-deployment-api.com/palantir/apollo/deployment"
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

// safelogging:@Safe
type MyObject struct {
	SafeExternalField deployment.ApolloEnvironmentId ` + "`json:\"safeExternalField\" safelogging:\"@Safe\"`" + `
}

func (o MyObject) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (o *MyObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
		{
			name: "Object with External type with UNSAFE safety overridden by field SAFE safety",
			in: &types.ObjectType{
				Name: "MyObject",
				Fields: []*types.Field{
					{
						Name: "overriddenField",
						Type: types.NewExternalWithSafety(
							spec.TypeName{
								Name:    "com/palantir/apollo/deployment:ApolloEnvironmentId",
								Package: "github.com/palantir/apollo-deployment-api",
							},
							types.String{},
							func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_UNSAFE); return &s }(),
						),
						Safety: func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_SAFE); return &s }(),
					},
				},
			},
			Out: `package testpkg

import (
	deployment "github.com/palantir/apollo-deployment-api.com/palantir/apollo/deployment"
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

// safelogging:@Safe
type MyObject struct {
	OverriddenField deployment.ApolloEnvironmentId ` + "`json:\"overriddenField\" safelogging:\"@Safe\"`" + `
}

func (o MyObject) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (o *MyObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := jen.NewFile("testpkg")
			safetyCache := make(map[types.Type]spec.LogSafety)
			writeObjectType(f.Group, tc.in, safetyCache, OutputConfiguration{})
			var buf bytes.Buffer
			assert.NoError(t, f.Render(&buf))
			assert.Equal(t, tc.Out, buf.String())
		})
	}
}
