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
	"fmt"
	"regexp"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/encoding"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	clientStructFieldName = "client"
	wrappedClientVar      = "client"

	clientReceiverName = "c"
	ctxName            = "ctx"
	authHeaderVar      = "authHeader"
	cookieTokenVar     = "cookieToken"
	tokenProviderVar   = "tokenProvider"

	defaultReturnValVar = "defaultReturnVal"
	returnValVar        = "returnVal"
	respVar             = "resp"

	requestParamsVar = "requestParams"
	queryParamsVar   = "queryParams"
)

var (
	pathParamRegexp = regexp.MustCompile(regexp.QuoteMeta("{") + "[^}]+" + regexp.QuoteMeta("}"))
)

func writeServiceType(file *jen.Group, serviceDef *types.ServiceDefinition) {
	file.Add(astForServiceInterface(serviceDef, false, false))
	file.Add(astForClientStructDecl(serviceDef.Name))
	file.Add(astForNewClientFunc(serviceDef.Name))
	for _, endpointDef := range serviceDef.Endpoints {
		file.Add(astForEndpointMethod(serviceDef.Name, endpointDef, false))
	}
	if serviceDef.HasHeaderAuth() || serviceDef.HasCookieAuth() {
		// at least one endpoint uses authentication: define decorator structures
		file.Add(astForServiceInterface(serviceDef, true, false))
		file.Add(astForNewServiceFuncWithAuth(serviceDef))
		file.Add(astForClientStructDeclWithAuth(serviceDef))
		for _, endpointDef := range serviceDef.Endpoints {
			file.Add(astForEndpointMethod(serviceDef.Name, endpointDef, true))
		}

		// Return true if all endpoints that require authentication are of the same auth type (header or cookie) and at least
		// one endpoint has auth. The same auth type is required because a single token provider will likely not be useful for
		// both auth types.
		if serviceDef.HasHeaderAuth() != serviceDef.HasCookieAuth() {
			file.Add(astForNewTokenServiceFunc(serviceDef.Name))
			file.Add(astForTokenServiceStructDecl(serviceDef.Name))
			for _, endpointDef := range serviceDef.Endpoints {
				file.Add(astForTokenServiceEndpointMethod(serviceDef.Name, endpointDef))
			}
		}
	}
}

func astForServiceInterface(serviceDef *types.ServiceDefinition, withAuth, isServer bool) *jen.Statement {
	name := interfaceTypeName(serviceDef.Name)
	if !isServer {
		name = clientInterfaceTypeName(serviceDef.Name)
	}
	if withAuth {
		name = withAuthName(name)
	}
	return serviceDef.CommentLine().
		Type().
		Id(name).
		InterfaceFunc(func(g *jen.Group) {
			for _, endpointDef := range serviceDef.Endpoints {
				g.Add(endpointDef.Docs.CommentLine()).Id(transforms.Export(endpointDef.EndpointName)).
					ParamsFunc(func(g *jen.Group) {
						astForEndpointArgsFunc(g, endpointDef, withAuth, isServer)
					}).
					ParamsFunc(func(g *jen.Group) {
						astForEndpointReturnsFunc(g, endpointDef)
					})
			}
		})
}

func astForEndpointArgsFunc(g *jen.Group, endpointDef *types.EndpointDefinition, withAuth, isServer bool) {
	g.Id(ctxName).Add(snip.Context())
	if !withAuth {
		if endpointDef.HeaderAuth {
			g.Id(authHeaderVar).Add(types.Bearertoken{}.Code())
		} else if endpointDef.CookieAuth != nil {
			g.Id(cookieTokenVar).Add(types.Bearertoken{}.Code())
		}
	}
	for _, paramDef := range endpointDef.Params {
		g.Add(astForEndpointParameterArg(paramDef, isServer))
	}
}

func astForEndpointParameterArg(argDef *types.EndpointArgumentDefinition, isServer bool) *jen.Statement {
	argType := argDef.Type.Code()
	if argDef.Type.IsBinary() {
		// special case: "binary" types resolve to []byte, but this indicates a streaming parameter when
		// specified as the request argument of a service, so use "io.ReadCloser".
		if isServer {
			argType = snip.IOReadCloser()
		} else {
			// special case: the client provides "func() io.ReadCloser" instead of "io.ReadCloser" so
			// that a fresh "io.ReadCloser" can be retrieved for retries.
			argType = snip.FuncIOReadCloser()
		}
	}
	return jen.Id(argNameTransform(argDef.Name)).Add(argType)
}

func astForEndpointReturnsFunc(g *jen.Group, endpointDef *types.EndpointDefinition) {
	if endpointDef.Returns != nil {
		r := *endpointDef.Returns
		if !r.IsBinary() {
			g.Add(r.Code())
		} else {
			// special case: "binary" type resolves to []byte in structs, but indicates a streaming response when
			// specified as the return type of a service, so replace all nested references with "io.ReadCloser".
			if r.IsOptional() {
				g.Op("*").Add(snip.IOReadCloser())
			} else {
				g.Add(snip.IOReadCloser())
			}
		}
	}
	g.Error()
}

func astForClientStructDecl(serviceName string) *jen.Statement {
	return jen.Type().Id(clientStructTypeName(serviceName)).Struct(
		jen.Id(clientStructFieldName).Add(snip.CGRClientClient()),
	)
}

func astForClientStructDeclWithAuth(serviceDef *types.ServiceDefinition) *jen.Statement {
	return jen.Type().Id(withAuthName(clientStructTypeName(serviceDef.Name))).StructFunc(func(g *jen.Group) {
		g.Id(clientStructFieldName).Id(clientInterfaceTypeName(serviceDef.Name))
		if serviceDef.HasHeaderAuth() {
			g.Id(authHeaderVar).Add(types.Bearertoken{}.Code())
		}
		if serviceDef.HasCookieAuth() {
			g.Id(cookieTokenVar).Add(types.Bearertoken{}.Code())
		}
	})
}

func astForNewClientFunc(serviceName string) *jen.Statement {
	return jen.Func().Id("New" + clientInterfaceTypeName(serviceName)).
		Params(jen.Id(wrappedClientVar).Add(snip.CGRClientClient())).
		Params(jen.Id(clientInterfaceTypeName(serviceName))).
		Block(jen.Return(
			jen.Op("&").Id(clientStructTypeName(serviceName)).ValuesFunc(func(g *jen.Group) {
				g.Id(clientStructFieldName).Op(":").Id(wrappedClientVar)
			}),
		))
}

func astForNewServiceFuncWithAuth(serviceDef *types.ServiceDefinition) *jen.Statement {
	return jen.Func().Id(withAuthName("New" + clientInterfaceTypeName(serviceDef.Name))).
		ParamsFunc(func(g *jen.Group) {
			g.Id(wrappedClientVar).Id(clientInterfaceTypeName(serviceDef.Name))
			if serviceDef.HasHeaderAuth() {
				g.Id(authHeaderVar).Add(types.Bearertoken{}.Code())
			}
			if serviceDef.HasCookieAuth() {
				g.Id(cookieTokenVar).Add(types.Bearertoken{}.Code())
			}
		}).
		Params(jen.Id(withAuthName(clientInterfaceTypeName(serviceDef.Name)))).
		Block(jen.Return(
			jen.Op("&").Id(withAuthName(clientStructTypeName(serviceDef.Name))).ValuesFunc(func(g *jen.Group) {
				g.Id(clientStructFieldName).Op(":").Id(wrappedClientVar)
				if serviceDef.HasHeaderAuth() {
					g.Id(authHeaderVar).Op(":").Id(authHeaderVar)
				}
				if serviceDef.HasCookieAuth() {
					g.Id(cookieTokenVar).Op(":").Id(cookieTokenVar)
				}
			}),
		))
}

func astForEndpointMethod(serviceName string, endpointDef *types.EndpointDefinition, withAuth bool) *jen.Statement {
	return jen.Func().
		ParamsFunc(func(g *jen.Group) {
			if withAuth {
				g.Id(clientReceiverName).Op("*").Id(withAuthName(clientStructTypeName(serviceName)))
			} else {
				g.Id(clientReceiverName).Op("*").Id(clientStructTypeName(serviceName))
			}
		}).
		Id(transforms.Export(endpointDef.EndpointName)).
		ParamsFunc(func(g *jen.Group) {
			astForEndpointArgsFunc(g, endpointDef, withAuth, false)
		}).
		ParamsFunc(func(g *jen.Group) {
			astForEndpointReturnsFunc(g, endpointDef)
		}).
		BlockFunc(func(g *jen.Group) {
			if withAuth {
				astForEndpointAuthMethodBodyFunc(g, endpointDef)
			} else {
				astForEndpointMethodBodyFunc(g, endpointDef)
			}
		})
}

func astForEndpointMethodBodyFunc(g *jen.Group, endpointDef *types.EndpointDefinition) {
	var (
		hasReturnVal      = endpointDef.Returns != nil
		returnsBinary     = hasReturnVal && (*endpointDef.Returns).IsBinary()
		returnsCollection = hasReturnVal && (*endpointDef.Returns).IsCollection()
		returnsOptional   = hasReturnVal && (*endpointDef.Returns).IsOptional()
		// If return can not be nil, we'll declare a zero-value variable to return in case of error
		returnDefaultValue = false
	)
	if hasReturnVal {
		_, returnsAliasType := (*endpointDef.Returns).(*types.AliasType)
		returnDefaultValue = hasReturnVal &&
			!returnsBinary &&
			!returnsCollection &&
			(!returnsOptional || returnsAliasType) // alias<optional<>> creates a struct with pointer field, return default empty struct
	}
	returnVar := func(g *jen.Group) {
		switch {
		case returnDefaultValue:
			g.Id(defaultReturnValVar)
		case hasReturnVal:
			g.Nil()
		}
	}

	// if endpoint returns a value, declare variables for value
	if hasReturnVal && !returnsBinary {
		if returnDefaultValue {
			g.Var().Id(defaultReturnValVar).Add((*endpointDef.Returns).Code())
			g.Var().Id(returnValVar).Op("*").Add((*endpointDef.Returns).Code())
		} else {
			g.Var().Id(returnValVar).Add((*endpointDef.Returns).Code())
		}
	}
	// build requestParams
	astForEndpointMethodBodyRequestParams(g, endpointDef)

	// execute request
	callStmt := jen.Id(clientReceiverName).Dot(clientStructFieldName).Dot("Do").Call(
		jen.Id("ctx"),
		jen.Id(requestParamsVar).Op("..."))
	returnErr := jen.ReturnFunc(func(g *jen.Group) {
		returnVar(g)
		g.Add(snip.WerrorWrapContext()).Call(
			jen.Id("ctx"),
			jen.Err(),
			jen.Lit(fmt.Sprintf("%s failed", endpointDef.EndpointName)),
		)
	})

	if returnsBinary {
		g.List(jen.Id(respVar), jen.Err()).Op(":=").Add(callStmt)
		g.If(jen.Err().Op("!=").Nil()).Block(returnErr)
		if returnsOptional {
			// If an endpoint with a return type of optional<binary> provides a response with a code of StatusNoContent
			// then the return value is empty and nil is returned.
			g.If(jen.Id(respVar).Dot("StatusCode").Op("==").Add(snip.HTTPStatusNoContent())).Block(
				jen.Return(jen.Nil(), jen.Nil()),
			)
			g.Return(jen.Op("&").Id(respVar).Dot("Body"), jen.Nil())
		} else {
			// if endpoint returns binary, return body of response directly
			g.Return(jen.Id(respVar).Dot("Body"), jen.Nil())
		}
		return
	}

	g.If(
		jen.List(jen.Id("_"), jen.Err()).Op(":=").Add(callStmt),
		jen.Err().Op("!=").Nil(),
	).Block(returnErr)

	if returnDefaultValue || returnsCollection {
		// verify that return value is non-nil and dereference
		g.If(jen.Id(returnValVar).Op("==").Nil()).Block(jen.ReturnFunc(func(g *jen.Group) {
			returnVar(g)
			g.Add(snip.WerrorErrorContext()).Call(
				jen.Id("ctx"),
				jen.Lit(fmt.Sprintf("%s response cannot be nil", endpointDef.EndpointName)),
			)
		}))
	}

	if returnDefaultValue {
		g.Return(jen.Op("*").Id(returnValVar), jen.Nil())
	} else if hasReturnVal {
		g.Return(jen.Id(returnValVar), jen.Nil())
	} else {
		g.Return(jen.Nil())
	}
}

func astForEndpointMethodBodyRequestParams(g *jen.Group, endpointDef *types.EndpointDefinition) {
	g.Var().Id(requestParamsVar).Op("[]").Add(snip.CGRClientRequestParam())

	// helper for the statement "requestParams = append(requestParams, {code})"
	appendRequestParams := func(g *jen.Group, code jen.Code) {
		g.Id(requestParamsVar).Op("=").Append(jen.Id(requestParamsVar), code)
	}

	appendRequestParams(g, snip.CGRClientWithRPCMethodName().Call(jen.Lit(transforms.Export(endpointDef.EndpointName))))
	appendRequestParams(g, snip.CGRClientWithRequestMethod().Call(jen.Lit(endpointDef.HTTPMethod.String())))
	// auth params
	if endpointDef.HeaderAuth {
		appendRequestParams(g, snip.CGRClientWithHeader().Call(jen.Lit("Authorization"),
			snip.FmtSprint().Call(jen.Lit("Bearer "), jen.Id(authHeaderVar)),
		))
	} else if endpointDef.CookieAuth != nil {
		appendRequestParams(g, snip.CGRClientWithHeader().Call(jen.Lit("Cookie"),
			snip.FmtSprint().Call(jen.Lit(*endpointDef.CookieAuth+"="), jen.Id(cookieTokenVar)),
		))
	}
	// path params
	appendRequestParams(g, snip.CGRClientWithPathf().CallFunc(func(g *jen.Group) {
		g.Lit(pathParamRegexp.ReplaceAllString(endpointDef.HTTPPath, regexp.QuoteMeta(`%s`)))
		for _, param := range endpointDef.PathParams() {
			g.Add(snip.URLPathEscape()).Call(snip.FmtSprint().Call(jen.Id(argNameTransform(param.ParamID))))
		}
	}))
	// body params
	if body := endpointDef.BodyParam(); body != nil {
		if body.Type.IsBinary() {
			appendRequestParams(g, snip.CGRClientWithRawRequestBodyProvider().Call(jen.Id(argNameTransform(body.Name))))
		} else {
			switch body.Type.(type) {
			case *types.AliasType, *types.EnumType, *types.ObjectType, *types.UnionType:
				appendRequestParams(g, snip.CGRClientWithJSONRequest().Call(jen.Id(argNameTransform(body.Name))))
			default:
				appendRequestParams(g, snip.CGRClientWithJSONRequest().Call(snip.SafeJSONAppendFunc().Call(jen.Func().
					Params(jen.Id("out").Op("[]").Byte()).
					Params(jen.Op("[]").Byte(), jen.Error()).
					BlockFunc(func(g *jen.Group) {
						encoding.AnonFuncBodyAppendJSON(g, jen.Id(argNameTransform(body.Name)).Clone, body.Type)
					}))),
				)
			}
		}
	}
	// header params
	for _, param := range endpointDef.HeaderParams() {
		argName := argNameTransform(param.Name)
		if param.Type.IsOptional() {
			selector := jen.Id(argName).Clone
			if _, isAlias := param.Type.(*types.AliasType); isAlias {
				selector = selector().Dot("Value").Clone
			}
			// if header parameter type is an optional, append dereferenced value if it is non-nil
			g.If(selector().Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
				appendRequestParams(g, snip.CGRClientWithHeader().Call(
					jen.Lit(param.ParamID),
					snip.FmtSprint().Call(jen.Op("*").Add(selector())),
				))
			})
		} else {
			appendRequestParams(g, snip.CGRClientWithHeader().Call(
				jen.Lit(param.ParamID),
				snip.FmtSprint().Call(jen.Id(argName)),
			))
		}
	}
	// query params
	if queryParams := endpointDef.QueryParams(); len(queryParams) > 0 {
		g.Id(queryParamsVar).Op(":=").Make(snip.URLValues())
		for _, param := range endpointDef.QueryParams() {
			argName := argNameTransform(param.Name)
			if param.Type.IsOptional() {
				selector := jen.Id(argName).Clone
				if _, isAlias := param.Type.(*types.AliasType); isAlias {
					selector = selector().Dot("Value").Clone
				}
				g.If(selector().Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
					g.Id(queryParamsVar).Dot("Set").Call(jen.Lit(param.ParamID),
						snip.FmtSprint().Call(jen.Op("*").Add(selector())))
				})
			} else if param.Type.IsList() {
				g.For(jen.List(jen.Id("_"), jen.Id("v")).Op(":=").Range().Id(argName)).Block(
					jen.Id(queryParamsVar).Dot("Add").Call(jen.Lit(param.ParamID), snip.FmtSprint().Call(jen.Id("v"))),
				)
			} else {
				g.Id(queryParamsVar).Dot("Set").Call(jen.Lit(param.ParamID),
					snip.FmtSprint().Call(jen.Id(argName)))
			}
		}
		appendRequestParams(g, snip.CGRClientWithQueryValues().Call(jen.Id(queryParamsVar)))
	}
	// response
	if endpointDef.Returns != nil {
		if (*endpointDef.Returns).IsBinary() {
			appendRequestParams(g, snip.CGRClientWithRawResponseBody().Call())
		} else {
			appendRequestParams(g, snip.CGRClientWithJSONResponse().Call(jen.Op("&").Id(returnValVar)))
		}
	}
}

func astForEndpointAuthMethodBodyFunc(g *jen.Group, endpointDef *types.EndpointDefinition) {
	g.Return(
		jen.Id(clientReceiverName).
			Dot(clientStructFieldName).
			Dot(transforms.Export(endpointDef.EndpointName)).
			CallFunc(func(g *jen.Group) {
				g.Id("ctx")
				if endpointDef.HeaderAuth {
					g.Id(clientReceiverName).Dot(authHeaderVar)
				}
				if endpointDef.CookieAuth != nil {
					g.Id(clientReceiverName).Dot(cookieTokenVar)
				}
				for _, param := range endpointDef.Params {
					g.Id(argNameTransform(param.Name))
				}
			}),
	)
}

func astForNewTokenServiceFunc(serviceName string) *jen.Statement {
	return jen.Func().Id(withTokenProviderName("New"+clientInterfaceTypeName(serviceName))).
		Params(
			jen.Id(wrappedClientVar).Id(clientInterfaceTypeName(serviceName)),
			jen.Id(tokenProviderVar).Add(snip.CGRClientTokenProvider()),
		).
		Params(jen.Id(withAuthName(clientInterfaceTypeName(serviceName)))).
		Block(jen.Return(
			jen.Op("&").Id(withTokenProviderName(clientStructTypeName(serviceName))).Values(
				jen.Id(clientStructFieldName).Op(":").Id(wrappedClientVar),
				jen.Id(tokenProviderVar).Op(":").Id(tokenProviderVar),
			),
		))
}

func astForTokenServiceStructDecl(serviceName string) *jen.Statement {
	return jen.Type().Id(withTokenProviderName(clientStructTypeName(serviceName))).Struct(
		jen.Id(clientStructFieldName).Id(clientInterfaceTypeName(serviceName)),
		jen.Id(tokenProviderVar).Add(snip.CGRClientTokenProvider()),
	)
}

func astForTokenServiceEndpointMethod(serviceName string, endpointDef *types.EndpointDefinition) *jen.Statement {
	hasAuth := endpointDef.HeaderAuth || endpointDef.CookieAuth != nil
	return jen.Func().
		Params(jen.Id(clientReceiverName).Op("*").Id(withTokenProviderName(clientStructTypeName(serviceName)))).
		Id(transforms.Export(endpointDef.EndpointName)).
		ParamsFunc(func(g *jen.Group) {
			astForEndpointArgsFunc(g, endpointDef, true, false)
		}).
		ParamsFunc(func(g *jen.Group) {
			astForEndpointReturnsFunc(g, endpointDef)
		}).
		BlockFunc(func(g *jen.Group) {
			if hasAuth {
				if endpointDef.Returns != nil {
					g.Var().Id(defaultReturnValVar).Add((*endpointDef.Returns).Code())
				}
				g.List(jen.Id("token"), jen.Err()).Op(":=").Id(clientReceiverName).Dot(tokenProviderVar).Call(jen.Id("ctx"))
				g.If(jen.Err().Op("!=").Nil()).Block(jen.ReturnFunc(func(g *jen.Group) {
					if endpointDef.Returns != nil {
						g.Id(defaultReturnValVar)
					}
					g.Err()
				}))
			}
			g.Return(jen.Id(clientReceiverName).Dot(clientStructFieldName).Dot(transforms.Export(endpointDef.EndpointName)).
				CallFunc(func(g *jen.Group) {
					g.Id("ctx")
					if hasAuth {
						g.Add(types.Bearertoken{}.Code()).Call(jen.Id("token"))
					}
					for _, paramDef := range endpointDef.Params {
						g.Id(argNameTransform(paramDef.Name))
					}
				}),
			)
		})
}

func interfaceTypeName(serviceName string) string {
	return transforms.Export(serviceName)
}

func clientInterfaceTypeName(serviceName string) string {
	return interfaceTypeName(serviceName) + "Client"
}

func clientStructTypeName(serviceName string) string {
	return transforms.Private(serviceName) + "Client"
}

func withAuthName(name string) string {
	return name + "WithAuth"
}

func withTokenProviderName(name string) string {
	return name + "WithTokenProvider"
}

// argNameTransform returns the input string with "Arg" appended to it. This transformation is done to ensure that
// argument variable names do not shadow any package names.
func argNameTransform(input string) string {
	return input + "Arg"
}
