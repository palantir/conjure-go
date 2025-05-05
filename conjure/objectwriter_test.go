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

import cj "github.com/palantir/conjure-go/v6/cj"

type User struct {
	UserName string ` + "`json:\"UserName\"`" + `
}

func (o User) MarshalYAML() (any, error) {
	return cj.YAMLV3MarshalerFromJSON(o)
}
func (o *User) UnmarshalYAML(unmarshal func(any) error) error {
	return cj.YAMLV3UnmarshalerToJSON(o, unmarshal)
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
	cj "github.com/palantir/conjure-go/v6/cj"
	safejson "github.com/palantir/pkg/safejson"
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
func (o User) MarshalYAML() (any, error) {
	return cj.YAMLV3MarshalerFromJSON(o)
}
func (o *User) UnmarshalYAML(unmarshal func(any) error) error {
	return cj.YAMLV3UnmarshalerToJSON(o, unmarshal)
}
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := jen.NewFile("testpkg")
			writeObjectType(OutputConfiguration{}, f.Group, tc.in)
			var buf bytes.Buffer
			assert.NoError(t, f.Render(&buf))
			assert.Equal(t, tc.Out, buf.String())
		})
	}
}
