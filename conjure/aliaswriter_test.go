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
	"testing"

	"github.com/dave/jennifer/jen"
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
			Name: "astForAliasTextStringer",
			In:   astForAliasTextStringer("Foo", types.DateTime{}.Code()),
			Out: `func (a Foo) String() string {
	return datetime.DateTime(a).String()
}`,
		},
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
			In:   astForAliasOptionalStringTextUnmarshal("Foo"),
			Out: `func (a *Foo) UnmarshalText(data []byte) error {
	rawFoo := string(data)
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
		return nil, nil
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
