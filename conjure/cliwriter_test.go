package conjure

import (
	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var fakeTestService = &types.ServiceDefinition{
	Name: "MyTestService",
	Endpoints: []*types.EndpointDefinition{
		{
			EndpointName: "GetResultByCustomId",
			Params: []*types.EndpointArgumentDefinition{
				{
					Name: "customId",
					Type: &types.AliasType{
						Docs: "",
						Name: "CustomId",
						Item: &types.String{},
					},
				},
			},
		},
	},
}

func TestAstCLIRoot(t *testing.T) {
	file := jen.NewFile("cli")
	astCLIRoot(file.Group)
	writeCLIType(file.Group, fakeTestService)
	astInitFunc(file.Group, []*types.ServiceDefinition{fakeTestService})
	err := file.Render(os.Stdout)
	assert.NoError(t, err)
}
