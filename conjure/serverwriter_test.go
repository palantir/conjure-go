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

package conjure

import (
	"bytes"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerASTDecodeHTTPParam(t *testing.T) {
	for _, test := range []struct {
		Name string
		Arg  types.EndpointArgumentDefinition
		Out  string
	}{
		{
			Name: "string path param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      types.String{},
				ParamType: types.PathParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParam, ok := pathParams["myParam"]
	if !ok {
		return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myParam\" not present")
	}
}`,
		},
		{
			Name: "bearertoken query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      types.Bearertoken{},
				ParamType: types.QueryParam,
				ParamID:   "myParamCustomParamId",
			},
			Out: `{
	myParam := bearertoken.Token(req.URL.Query().Get("myParamCustomParamId"))
}`,
		},
		{
			Name: "optional bearertoken path param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      &types.Optional{Item: types.Bearertoken{}},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	var myParam *bearertoken.Token
	if myParamStr := req.URL.Query().Get("myParam"); myParamStr != "" {
		myParamInternal := bearertoken.Token(myParamStr)
		myParam = &myParamInternal
	}
}`,
		},
		{
			Name: "list bearertoken query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      &types.List{Item: types.Bearertoken{}},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	var myParam []bearertoken.Token
	for _, v := range req.URL.Query()["myParam"] {
		convertedVal := bearertoken.Token(v)
		myParam = append(myParam, convertedVal)
	}
}`,
		},
		{
			Name: "binary query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      types.Binary{},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParam := []byte(req.URL.Query().Get("myParam"))
}`,
		},
		{
			Name: "boolean query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      types.Boolean{},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParam, err := strconv.ParseBool(req.URL.Query().Get("myParam"))
	if err != nil {
		return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as boolean")
	}
}`,
		},
		{
			Name: "optional boolean query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      &types.Optional{Item: types.Boolean{}},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	var myParam *bool
	if myParamStr := req.URL.Query().Get("myParam"); myParamStr != "" {
		myParamInternal, err := strconv.ParseBool(myParamStr)
		if err != nil {
			return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as boolean")
		}
		myParam = &myParamInternal
	}
}`,
		},
		{
			Name: "list boolean query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      &types.List{Item: types.Boolean{}},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	var myParam []bool
	for _, v := range req.URL.Query()["myParam"] {
		convertedVal, err := strconv.ParseBool(v)
		if err != nil {
			return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as boolean")
		}
		myParam = append(myParam, convertedVal)
	}
}`,
		},
		{
			Name: "datetime query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      types.DateTime{},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParam, err := datetime.ParseDateTime(req.URL.Query().Get("myParam"))
	if err != nil {
		return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as datetime")
	}
}`,
		},
		{
			Name: "optional datetime query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      &types.Optional{Item: types.DateTime{}},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	var myParam *datetime.DateTime
	if myParamStr := req.URL.Query().Get("myParam"); myParamStr != "" {
		myParamInternal, err := datetime.ParseDateTime(myParamStr)
		if err != nil {
			return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as datetime")
		}
		myParam = &myParamInternal
	}
}`,
		},
		{
			Name: "list datetime query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      &types.List{Item: types.DateTime{}},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	var myParam []datetime.DateTime
	for _, v := range req.URL.Query()["myParam"] {
		convertedVal, err := datetime.ParseDateTime(v)
		if err != nil {
			return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as datetime")
		}
		myParam = append(myParam, convertedVal)
	}
}`,
		},
		{
			Name: "binary query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      types.Binary{},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParam := []byte(req.URL.Query().Get("myParam"))
}`,
		},
		{
			Name: "double query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      types.Double{},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParam, err := strconv.ParseFloat(req.URL.Query().Get("myParam"), 64)
	if err != nil {
		return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as double")
	}
}`,
		},
		{
			Name: "optional double query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      &types.Optional{Item: types.Double{}},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	var myParam *float64
	if myParamStr := req.URL.Query().Get("myParam"); myParamStr != "" {
		myParamInternal, err := strconv.ParseFloat(myParamStr, 64)
		if err != nil {
			return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as double")
		}
		myParam = &myParamInternal
	}
}`,
		},
		{
			Name: "list double query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      &types.List{Item: types.Double{}},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	var myParam []float64
	for _, v := range req.URL.Query()["myParam"] {
		convertedVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as double")
		}
		myParam = append(myParam, convertedVal)
	}
}`,
		},
		{
			Name: "integer query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      types.Integer{},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParam, err := strconv.Atoi(req.URL.Query().Get("myParam"))
	if err != nil {
		return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as integer")
	}
}`,
		},
		{
			Name: "rid query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      types.RID{},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParam, err := rid.ParseRID(req.URL.Query().Get("myParam"))
	if err != nil {
		return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as rid")
	}
}`,
		},
		{
			Name: "safelong query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      types.Safelong{},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParam, err := safelong.ParseSafeLong(req.URL.Query().Get("myParam"))
	if err != nil {
		return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as safelong")
	}
}`,
		},
		{
			Name: "string query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      types.String{},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParam := req.URL.Query().Get("myParam")
}`,
		},
		{
			Name: "optional string query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      &types.Optional{Item: types.String{}},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	var myParam *string
	if myParamStr := req.URL.Query().Get("myParam"); myParamStr != "" {
		myParamInternal := myParamStr
		myParam = &myParamInternal
	}
}`,
		},
		{
			Name: "list string query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      &types.List{Item: types.String{}},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParam := req.URL.Query()["myParam"]
}`,
		},
		{
			Name: "uuid query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      types.UUID{},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParam, err := uuid.ParseUUID(req.URL.Query().Get("myParam"))
	if err != nil {
		return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as uuid")
	}
}`,
		},
		{
			Name: "any query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      types.Any{},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParam := req.URL.Query().Get("myParam")
}`,
		},
		{
			Name: "optional any query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      &types.Optional{Item: types.Any{}},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	var myParam *interface{}
	if myParamStr := req.URL.Query().Get("myParam"); myParamStr != "" {
		myParamInternal := myParamStr
		myParam = &myParamInternal
	}
}`,
		},
		{
			Name: "list any query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      &types.List{Item: types.Any{}},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	var myParam []interface{}
	for _, v := range req.URL.Query()["myParam"] {
		convertedVal := v
		myParam = append(myParam, convertedVal)
	}
}`,
		},
		{
			Name: "enum query param",
			Arg: types.EndpointArgumentDefinition{
				Name:      "myParam",
				Type:      &types.EnumType{Name: "MyEnum"},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	var myParam MyEnum
	if err := myParam.UnmarshalText([]byte(req.URL.Query().Get("myParam"))); err != nil {
		return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as MyEnum")
	}
}`,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			out := jen.BlockFunc(func(g *jen.Group) {
				switch test.Arg.ParamType {
				case types.PathParam:
					astForHandlerMethodPathParam(g, &test.Arg)
				case types.HeaderParam:
					astForHandlerMethodHeaderParam(g, &test.Arg)
				case types.QueryParam:
					astForHandlerMethodQueryParam(g, &test.Arg)
				case types.BodyParam:
					astForHandlerMethodDecodeBody(g, &test.Arg)
				}
			})
			var buf bytes.Buffer
			require.NoError(t, out.Render(&buf))
			assert.Equal(t, test.Out, buf.String())
		})
	}
}
