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
