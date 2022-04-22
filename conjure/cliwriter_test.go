// Copyright (c) 2022 Palantir Technologies. All rights reserved.
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
	"os"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/stretchr/testify/assert"
)

var typeBinary types.Type = types.Integer{}

var fakeTestService = &types.ServiceDefinition{
	Name: "MyTestService",
	Endpoints: []*types.EndpointDefinition{
		{
			EndpointName: "GetResultByCustomId",
			Returns:      &typeBinary,
			Params: []*types.EndpointArgumentDefinition{
				{
					Name: "customId",
					Type: &types.AliasType{
						Docs: "",
						Name: "CustomId",
						Item: &types.String{},
					},
				},
				{
					Name: "date",
					Type: types.DateTime{},
				},
				{
					Name: "myOptional",
					Type: &types.Optional{
						Item: types.RID{},
					},
				},
				{
					Name: "myList",
					Type: &types.List{
						Item: types.String{},
					},
				},
			},
		},
	},
}

func TestAstCLIRoot(t *testing.T) {
	file := jen.NewFile("cli")
	writeCLIType(file.Group, []*types.ServiceDefinition{fakeTestService})
	err := file.Render(os.Stdout)
	assert.NoError(t, err)
}
