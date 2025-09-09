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

package snip

import (
	"fmt"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/assert"
)

func TestSnippets(t *testing.T) {
	for _, test := range []struct {
		name     string
		stmt     jen.Code
		expected string
	}{
		{"MethodString", MethodString("x", "Foo"), `func (x Foo) String() string`},
		{"MethodMarshalJSON", MethodMarshalJSON("x", "Foo"), `func (x Foo) MarshalJSON() ([]byte, error)`},
		{"MethodUnmarshalJSON", MethodUnmarshalJSON("x", "Foo"), `func (x *Foo) UnmarshalJSON(data []byte) error`},
		{"MethodMarshalYAML", MethodMarshalYAML("x", "Foo"), `func (x Foo) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(x)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}`},
		{"MethodUnmarshalYAML", MethodUnmarshalYAML("x", "Foo"), `func (x *Foo) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&x)
}`},
		{"MethodStringImpl", MethodString("x", "Foo").Block(
			jen.Return(jen.String().Call(jen.Id("x")).Dot("String").Call()),
		), `func (x Foo) String() string {
	return string(x).String()
}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, fmt.Sprintf("%#v", test.stmt))
		})
	}
}
