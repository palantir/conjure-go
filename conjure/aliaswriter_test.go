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

package conjure

import (
	"bytes"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/stretchr/testify/assert"
)

func TestAliasWriter(t *testing.T) {
	for _, test := range []struct {
		Name string
		In   *jen.Statement
		Out  string
	}{
		{
			Name: "astForAliasTextMarshal",
			In:   astForAliasTextMarshal("Foo", types.DateTime{}.Code()),
			Out: `func (a Foo) MarshalText() ([]byte, error) {
	return datetime.DateTime(a).MarshalText()
}`,
		},
		{
			Name: "astForAliasOptionalTextMarshal",
			In:   astForAliasOptionalTextMarshal("Foo"),
			Out: `func (a Foo) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return a.Value.MarshalText()
}`,
		},
		{
			Name: "astForAliasOptionalStringTextMarshal",
			In:   astForAliasOptionalStringTextMarshal("Foo"),
			Out: `func (a Foo) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return []byte(*a.Value), nil
}`,
		},
		{
			Name: "astForAliasOptionalBinaryTextMarshal",
			In:   astForAliasOptionalBinaryTextMarshal("Foo"),
			Out: `func (a Foo) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return binary.New(*a.Value).MarshalText()
}`,
		},
		{
			Name: "astForAliasTextUnmarshal",
			In:   astForAliasTextUnmarshal("Foo", types.DateTime{}.Code()),
			Out: `func (a *Foo) UnmarshalText(data []byte) error {
	var rawFoo datetime.DateTime
	if err := rawFoo.UnmarshalText(data); err != nil {
		return err
	}
	*a = Foo(rawFoo)
	return nil
}`,
		},
		{
			Name: "astForAliasBinaryTextUnmarshal",
			In:   astForAliasBinaryTextUnmarshal("Foo"),
			Out: `func (a *Foo) UnmarshalText(data []byte) error {
	rawFoo, err := binary.Binary(data).Bytes()
	if err != nil {
		return err
	}
	*a = Foo(rawFoo)
	return nil
}`,
		},
		{
			Name: "astForAliasOptionalTextUnmarshal",
			In:   astForAliasOptionalTextUnmarshal("Foo", jen.New(types.DateTime{}.Code())),
			Out: `func (a *Foo) UnmarshalText(data []byte) error {
	if a.Value == nil {
		a.Value = new(datetime.DateTime)
	}
	return a.Value.UnmarshalText(data)
}`,
		},
		{
			Name: "astForAliasOptionalStringTextUnmarshal",
			In:   astForAliasOptionalStringTextUnmarshal("Foo", types.String{}.Code()),
			Out: `func (a *Foo) UnmarshalText(data []byte) error {
	rawFoo := string(data)
	a.Value = &rawFoo
	return nil
}`,
		},
		{
			Name: "astForAliasOptionalStringTextUnmarshal_Alias",
			In:   astForAliasOptionalStringTextUnmarshal("Foo", (&types.AliasType{Name: "FooAlias", Item: types.String{}}).Code()),
			Out: `func (a *Foo) UnmarshalText(data []byte) error {
	rawFoo := FooAlias(data)
	a.Value = &rawFoo
	return nil
}`,
		},
		{
			Name: "astForAliasOptionalBinaryTextUnmarshal",
			In:   astForAliasOptionalBinaryTextUnmarshal("Foo"),
			Out: `func (a *Foo) UnmarshalText(data []byte) error {
	rawFoo, err := binary.Binary(data).Bytes()
	if err != nil {
		return err
	}
	*a.Value = rawFoo
	return nil
}`,
		},
		{
			Name: "astForAliasJSONMarshal",
			In:   astForAliasJSONMarshal("Foo", types.DateTime{}.Code()),
			Out: `func (a Foo) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(datetime.DateTime(a))
}`,
		},
		{
			Name: "astForAliasOptionalJSONMarshal",
			In:   astForAliasOptionalJSONMarshal("Foo"),
			Out: `func (a Foo) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}`,
		},
		{
			Name: "astForAliasJSONUnmarshal",
			In:   astForAliasJSONUnmarshal("Foo", types.DateTime{}.Code()),
			Out: `func (a *Foo) UnmarshalJSON(data []byte) error {
	var rawFoo datetime.DateTime
	if err := safejson.Unmarshal(data, &rawFoo); err != nil {
		return err
	}
	*a = Foo(rawFoo)
	return nil
}`,
		},
		{
			Name: "astForAliasOptionalJSONUnmarshal",
			In:   astForAliasOptionalJSONUnmarshal("Foo", jen.New(types.DateTime{}.Code())),
			Out: `func (a *Foo) UnmarshalJSON(data []byte) error {
	if a.Value == nil {
		a.Value = new(datetime.DateTime)
	}
	return safejson.Unmarshal(data, a.Value)
}`,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.Out, test.In.GoString())
		})
	}
}

func TestWriteAliasTypeWithSafety(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   *types.AliasType
		out  string
	}{
		{
			name: "Simple alias without safety annotation",
			in: &types.AliasType{
				Name: "StringAlias",
				Item: types.String{},
			},
			out: `package testpkg

type StringAlias string
`,
		},
		{
			name: "Simple alias with unsafe safety annotation",
			in: types.NewAliasTypeWithSafety("UnsafeStringAlias", types.String{},
				func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_UNSAFE); return &s }()),
			out: `package testpkg

type UnsafeStringAlias string // safelogging:@Unsafe
`,
		},
		{
			name: "Simple alias with safe safety annotation",
			in: types.NewAliasTypeWithSafety("SafeStringAlias", types.String{},
				func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_SAFE); return &s }()),
			out: `package testpkg

type SafeStringAlias string // safelogging:@Safe
`,
		},
		{
			name: "Simple alias with do-not-log safety annotation",
			in: types.NewAliasTypeWithSafety("SecretStringAlias", types.String{},
				func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_DO_NOT_LOG); return &s }()),
			out: `package testpkg

type SecretStringAlias string // safelogging:@DoNotLog
`,
		},
		{
			name: "Alias with documentation and safety annotation",
			in: func() *types.AliasType {
				alias := types.NewAliasTypeWithSafety("DocumentedUnsafeAlias", types.String{},
					func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_UNSAFE); return &s }())
				alias.Docs = types.Docs("This is a documented unsafe alias.")
				return alias
			}(),
			out: `package testpkg

// This is a documented unsafe alias.
type DocumentedUnsafeAlias string // safelogging:@Unsafe
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := jen.NewFile("testpkg")
			writeAliasType(f.Group, tc.in)
			var buf bytes.Buffer
			assert.NoError(t, f.Render(&buf))
			assert.Equal(t, tc.out, buf.String())
		})
	}
}

func TestWriteOptionalAliasTypeWithSafety(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   *types.AliasType
		out  string
	}{
		{
			name: "Optional alias with documentation only (no safety)",
			in: func() *types.AliasType {
				alias := &types.AliasType{
					Name: "OptionalDocumentedAlias",
					Item: &types.Optional{Item: types.String{}},
				}
				alias.Docs = types.Docs("This is an optional documented alias without safety annotation.")
				return alias
			}(),
			out: `package testpkg

import (
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

// This is an optional documented alias without safety annotation.
type OptionalDocumentedAlias struct {
	Value *string
}

func (a OptionalDocumentedAlias) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return []byte(*a.Value), nil
}
func (a OptionalDocumentedAlias) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}
func (a *OptionalDocumentedAlias) UnmarshalText(data []byte) error {
	rawOptionalDocumentedAlias := string(data)
	a.Value = &rawOptionalDocumentedAlias
	return nil
}
func (a OptionalDocumentedAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (a *OptionalDocumentedAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}
`,
		},
		{
			name: "Optional alias with safety only (no documentation)",
			in: func() *types.AliasType {
				alias := types.NewAliasTypeWithSafety("OptionalUnsafeAlias", &types.Optional{Item: types.String{}},
					func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_UNSAFE); return &s }())
				return alias
			}(),
			out: `package testpkg

import (
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

// safelogging:@Unsafe
type OptionalUnsafeAlias struct {
	Value *string
}

func (a OptionalUnsafeAlias) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return []byte(*a.Value), nil
}
func (a OptionalUnsafeAlias) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}
func (a *OptionalUnsafeAlias) UnmarshalText(data []byte) error {
	rawOptionalUnsafeAlias := string(data)
	a.Value = &rawOptionalUnsafeAlias
	return nil
}
func (a OptionalUnsafeAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (a *OptionalUnsafeAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}
`,
		},
		{
			name: "Optional alias with both documentation and safety",
			in: func() *types.AliasType {
				alias := types.NewAliasTypeWithSafety("OptionalDocumentedUnsafeAlias", &types.Optional{Item: types.String{}},
					func() *spec.LogSafety { s := spec.New_LogSafety(spec.LogSafety_UNSAFE); return &s }())
				alias.Docs = types.Docs("This is an optional documented unsafe alias.")
				return alias
			}(),
			out: `package testpkg

import (
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

// This is an optional documented unsafe alias.
// safelogging:@Unsafe
type OptionalDocumentedUnsafeAlias struct {
	Value *string
}

func (a OptionalDocumentedUnsafeAlias) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return []byte(*a.Value), nil
}
func (a OptionalDocumentedUnsafeAlias) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}
func (a *OptionalDocumentedUnsafeAlias) UnmarshalText(data []byte) error {
	rawOptionalDocumentedUnsafeAlias := string(data)
	a.Value = &rawOptionalDocumentedUnsafeAlias
	return nil
}
func (a OptionalDocumentedUnsafeAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (a *OptionalDocumentedUnsafeAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := jen.NewFile("testpkg")
			writeOptionalAliasType(f.Group, tc.in)
			var buf bytes.Buffer
			assert.NoError(t, f.Render(&buf))
			assert.Equal(t, tc.out, buf.String())
		})
	}
}
