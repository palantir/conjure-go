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
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
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
	myParamArg, ok := pathParams["myParam"]
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
	myParamArg := bearertoken.Token(req.URL.Query().Get("myParamCustomParamId"))
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
	var myParamArg *bearertoken.Token
	if myParamArgStr := req.URL.Query().Get("myParam"); myParamArgStr != "" {
		myParamArgInternal := bearertoken.Token(myParamArgStr)
		myParamArg = &myParamArgInternal
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
	var myParamArg []bearertoken.Token
	for _, v := range req.URL.Query()["myParam"] {
		convertedVal := bearertoken.Token(v)
		myParamArg = append(myParamArg, convertedVal)
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
	myParamArg := []byte(req.URL.Query().Get("myParam"))
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
	myParamArg, err := strconv.ParseBool(req.URL.Query().Get("myParam"))
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
	var myParamArg *bool
	if myParamArgStr := req.URL.Query().Get("myParam"); myParamArgStr != "" {
		myParamArgInternal, err := strconv.ParseBool(myParamArgStr)
		if err != nil {
			return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as boolean")
		}
		myParamArg = &myParamArgInternal
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
	var myParamArg []bool
	for _, v := range req.URL.Query()["myParam"] {
		convertedVal, err := strconv.ParseBool(v)
		if err != nil {
			return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as boolean")
		}
		myParamArg = append(myParamArg, convertedVal)
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
	myParamArg, err := datetime.ParseDateTime(req.URL.Query().Get("myParam"))
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
	var myParamArg *datetime.DateTime
	if myParamArgStr := req.URL.Query().Get("myParam"); myParamArgStr != "" {
		myParamArgInternal, err := datetime.ParseDateTime(myParamArgStr)
		if err != nil {
			return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as datetime")
		}
		myParamArg = &myParamArgInternal
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
	var myParamArg []datetime.DateTime
	for _, v := range req.URL.Query()["myParam"] {
		convertedVal, err := datetime.ParseDateTime(v)
		if err != nil {
			return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as datetime")
		}
		myParamArg = append(myParamArg, convertedVal)
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
	myParamArg := []byte(req.URL.Query().Get("myParam"))
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
	myParamArg, err := strconv.ParseFloat(req.URL.Query().Get("myParam"), 64)
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
	var myParamArg *float64
	if myParamArgStr := req.URL.Query().Get("myParam"); myParamArgStr != "" {
		myParamArgInternal, err := strconv.ParseFloat(myParamArgStr, 64)
		if err != nil {
			return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as double")
		}
		myParamArg = &myParamArgInternal
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
	var myParamArg []float64
	for _, v := range req.URL.Query()["myParam"] {
		convertedVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as double")
		}
		myParamArg = append(myParamArg, convertedVal)
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
	myParamArg, err := strconv.Atoi(req.URL.Query().Get("myParam"))
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
	myParamArg, err := rid.ParseRID(req.URL.Query().Get("myParam"))
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
	myParamArg, err := safelong.ParseSafeLong(req.URL.Query().Get("myParam"))
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
	myParamArg := req.URL.Query().Get("myParam")
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
	var myParamArg *string
	if myParamArgStr := req.URL.Query().Get("myParam"); myParamArgStr != "" {
		myParamArgInternal := myParamArgStr
		myParamArg = &myParamArgInternal
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
	myParamArg := req.URL.Query()["myParam"]
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
	myParamArg, err := uuid.ParseUUID(req.URL.Query().Get("myParam"))
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
	myParamArg := req.URL.Query().Get("myParam")
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
	var myParamArg *interface{}
	if myParamArgStr := req.URL.Query().Get("myParam"); myParamArgStr != "" {
		myParamArgInternal := myParamArgStr
		myParamArg = &myParamArgInternal
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
	var myParamArg []interface{}
	for _, v := range req.URL.Query()["myParam"] {
		convertedVal := v
		myParamArg = append(myParamArg, convertedVal)
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
	var myParamArg MyEnum
	if err := myParamArg.UnmarshalText([]byte(req.URL.Query().Get("myParam"))); err != nil {
		return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as MyEnum")
	}
}`,
		},
		{
			Name: "external string query param",
			Arg: types.EndpointArgumentDefinition{
				Name: "myParam",
				Type: &types.External{
					Spec: spec.TypeName{
						Name:    "myParam",
						Package: "com.palantir.service.api.MyParam",
					},
					Fallback: types.String{},
				},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParamArg := req.URL.Query().Get("myParam")
}`,
		},
		{
			Name: "external string path param",
			Arg: types.EndpointArgumentDefinition{
				Name: "myParam",
				Type: &types.External{
					Spec: spec.TypeName{
						Name:    "myParam",
						Package: "com.palantir.service.api.MyParam",
					},
					Fallback: types.String{},
				},
				ParamType: types.PathParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParamArgStr, ok := pathParams["myParam"]
	if !ok {
		return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myParam\" not present")
	}
	myParamArg := myParamArgStr
}`,
		},
		{
			Name: "External integer query param",
			Arg: types.EndpointArgumentDefinition{
				Name: "myParam",
				Type: &types.External{
					Spec: spec.TypeName{
						Name:    "myParam",
						Package: "com.palantir.service.api.MyParam",
					},
					Fallback: types.Integer{},
				},
				ParamType: types.QueryParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParamArg, err := strconv.Atoi(req.URL.Query().Get("myParam"))
	if err != nil {
		return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as integer")
	}
}`,
		},
		{
			Name: "External integer path param",
			Arg: types.EndpointArgumentDefinition{
				Name: "myParam",
				Type: &types.External{
					Spec: spec.TypeName{
						Name:    "myParam",
						Package: "com.palantir.service.api.MyParam",
					},
					Fallback: types.Integer{},
				},
				ParamType: types.PathParam,
				ParamID:   "myParam",
			},
			Out: `{
	myParamArgStr, ok := pathParams["myParam"]
	if !ok {
		return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myParam\" not present")
	}
	myParamArg, err := strconv.Atoi(myParamArgStr)
	if err != nil {
		return witchcraftgoerror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as integer")
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
					astForHandlerMethodDecodeBody(g, &test.Arg, "MyService", "MyEndpoint")
				}
			})
			var buf bytes.Buffer
			require.NoError(t, out.Render(&buf))
			assert.Equal(t, test.Out, buf.String())
		})
	}
}
