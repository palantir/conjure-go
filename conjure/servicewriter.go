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

func writeServiceType(cfg OutputConfiguration, file *jen.Group, serviceDef *types.ServiceDefinition) {
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
		InterfaceFunc(func(methods *jen.Group) {
			for _, endpointDef := range serviceDef.Endpoints {
				methods.Add(endpointDef.CommentLineWithDeprecation(endpointDef.Deprecated)).
					Id(transforms.Export(endpointDef.EndpointName)).
					ParamsFunc(func(args *jen.Group) {
						astForEndpointArgsFunc(args, endpointDef, withAuth, isServer)
					}).
					ParamsFunc(func(args *jen.Group) {
						astForEndpointReturnsFunc(args, endpointDef)
					})
			}
		})
}

func astForEndpointArgsFunc(args *jen.Group, endpointDef *types.EndpointDefinition, withAuth, isServer bool) {
	args.Id(ctxName).Add(snip.Context())
	if !withAuth {
		if endpointDef.HeaderAuth {
			args.Id(authHeaderVar).Add(types.Bearertoken{}.Code())
		} else if endpointDef.CookieAuth != nil {
			args.Id(cookieTokenVar).Add(types.Bearertoken{}.Code())
		}
	}
	for _, paramDef := range endpointDef.Params {
		args.Add(astForEndpointParameterArg(paramDef, isServer))
	}
}

func astForEndpointParameterArg(argDef *types.EndpointArgumentDefinition, isServer bool) *jen.Statement {
	argType := argDef.Type.Code()
	if argDef.Type.IsBinary() {
		// special case: "binary" types resolve to []byte, but this indicates a streaming parameter when
		// specified as the request argument of a service, so use "io.ReadCloser".
		// If the type is optional<binary>, use "*io.ReadCloser".
		if isServer {
			if argDef.Type.IsOptional() {
				argType = jen.Op("*").Add(snip.IOReadCloser())
			} else {
				argType = snip.IOReadCloser()
			}
		} else {
			// special case: the client provides "func() io.ReadCloser" instead of "io.ReadCloser" so
			// that a fresh "io.ReadCloser" can be retrieved for retries.
			argType = snip.FuncIOReadCloser()
		}
	}
	return jen.Id(transforms.ArgName(argDef.Name)).Add(argType)
}

func astForEndpointReturnsFunc(args *jen.Group, endpointDef *types.EndpointDefinition) {
	if endpointDef.Returns != nil {
		r := *endpointDef.Returns
		if !r.IsBinary() {
			args.Add(r.Code())
		} else {
			// special case: "binary" type resolves to []byte in structs, but indicates a streaming response when
			// specified as the return type of a service, so replace all nested references with "io.ReadCloser".
			if r.IsOptional() {
				args.Op("*").Add(snip.IOReadCloser())
			} else {
				args.Add(snip.IOReadCloser())
			}
		}
	}
	args.Error()
}

func astForClientStructDecl(serviceName string) *jen.Statement {
	return jen.Type().Id(clientStructTypeName(serviceName)).Struct(
		jen.Id(clientStructFieldName).Add(snip.CGRClientClient()),
	)
}

func astForClientStructDeclWithAuth(serviceDef *types.ServiceDefinition) *jen.Statement {
	return jen.Type().Id(withAuthName(clientStructTypeName(serviceDef.Name))).StructFunc(func(structDecls *jen.Group) {
		structDecls.Id(clientStructFieldName).Id(clientInterfaceTypeName(serviceDef.Name))
		if serviceDef.HasHeaderAuth() {
			structDecls.Id(authHeaderVar).Add(types.Bearertoken{}.Code())
		}
		if serviceDef.HasCookieAuth() {
			structDecls.Id(cookieTokenVar).Add(types.Bearertoken{}.Code())
		}
	})
}

func astForNewClientFunc(serviceName string) *jen.Statement {
	return jen.Func().Id("New" + clientInterfaceTypeName(serviceName)).
		Params(jen.Id(wrappedClientVar).Add(snip.CGRClientClient())).
		Params(jen.Id(clientInterfaceTypeName(serviceName))).
		Block(jen.Return(
			jen.Op("&").Id(clientStructTypeName(serviceName)).ValuesFunc(func(values *jen.Group) {
				values.Id(clientStructFieldName).Op(":").Id(wrappedClientVar)
			}),
		))
}

func astForNewServiceFuncWithAuth(serviceDef *types.ServiceDefinition) *jen.Statement {
	return jen.Func().Id(withAuthName("New" + clientInterfaceTypeName(serviceDef.Name))).
		ParamsFunc(func(args *jen.Group) {
			args.Id(wrappedClientVar).Id(clientInterfaceTypeName(serviceDef.Name))
			if serviceDef.HasHeaderAuth() {
				args.Id(authHeaderVar).Add(types.Bearertoken{}.Code())
			}
			if serviceDef.HasCookieAuth() {
				args.Id(cookieTokenVar).Add(types.Bearertoken{}.Code())
			}
		}).
		Params(jen.Id(withAuthName(clientInterfaceTypeName(serviceDef.Name)))).
		Block(jen.Return(
			jen.Op("&").Id(withAuthName(clientStructTypeName(serviceDef.Name))).ValuesFunc(func(values *jen.Group) {
				values.Id(clientStructFieldName).Op(":").Id(wrappedClientVar)
				if serviceDef.HasHeaderAuth() {
					values.Id(authHeaderVar).Op(":").Id(authHeaderVar)
				}
				if serviceDef.HasCookieAuth() {
					values.Id(cookieTokenVar).Op(":").Id(cookieTokenVar)
				}
			}),
		))
}

func astForEndpointMethod(serviceName string, endpointDef *types.EndpointDefinition, withAuth bool) *jen.Statement {
	return jen.Func().
		ParamsFunc(func(receiver *jen.Group) {
			if withAuth {
				receiver.Id(clientReceiverName).Op("*").Id(withAuthName(clientStructTypeName(serviceName)))
			} else {
				receiver.Id(clientReceiverName).Op("*").Id(clientStructTypeName(serviceName))
			}
		}).
		Id(transforms.Export(endpointDef.EndpointName)).
		ParamsFunc(func(args *jen.Group) {
			astForEndpointArgsFunc(args, endpointDef, withAuth, false)
		}).
		ParamsFunc(func(args *jen.Group) {
			astForEndpointReturnsFunc(args, endpointDef)
		}).
		BlockFunc(func(methodBody *jen.Group) {
			if withAuth {
				astForEndpointAuthMethodBodyFunc(methodBody, endpointDef)
			} else {
				astForEndpointMethodBodyFunc(methodBody, endpointDef)
			}
		})
}

func astForEndpointMethodBodyFunc(methodBody *jen.Group, endpointDef *types.EndpointDefinition) {
	var (
		hasReturnVal         = endpointDef.Returns != nil
		returnsBinary        = hasReturnVal && (*endpointDef.Returns).IsBinary()
		returnsCollection    = hasReturnVal && (*endpointDef.Returns).IsCollection()
		returnsOptional      = hasReturnVal && (*endpointDef.Returns).IsOptional()
		returnsNamedOptional = returnsOptional && (*endpointDef.Returns).IsNamed()
		// If return can not be nil, we'll declare a zero-value variable to return in case of error
		returnDefaultValue = hasReturnVal && !returnsBinary && !returnsCollection && !returnsOptional
	)

	// if endpoint returns a value, declare variables for value and store default return type for error returns.
	returnVar := func(returnVals *jen.Group) {}
	switch {
	case returnsBinary:
		returnVar = func(returnVals *jen.Group) { returnVals.Nil() }
	case returnDefaultValue:
		methodBody.Var().Id(defaultReturnValVar).Add((*endpointDef.Returns).Code())
		methodBody.Var().Id(returnValVar).Op("*").Add((*endpointDef.Returns).Code())

		returnVar = func(returnVals *jen.Group) { returnVals.Id(defaultReturnValVar) }
	case returnsNamedOptional:
		// alias<optional<T>> creates a struct with pointer field, so return default empty struct
		methodBody.Var().Id(defaultReturnValVar).Add((*endpointDef.Returns).Code())
		methodBody.Var().Id(returnValVar).Add((*endpointDef.Returns).Code())

		returnVar = func(returnVals *jen.Group) { returnVals.Id(defaultReturnValVar) }
	case hasReturnVal:
		methodBody.Var().Id(returnValVar).Add((*endpointDef.Returns).Code())

		returnVar = func(returnVals *jen.Group) { returnVals.Nil() }
	}

	// build requestParams
	astForEndpointMethodBodyRequestParams(methodBody, endpointDef)

	// execute request
	callStmt := jen.Id(clientReceiverName).Dot(clientStructFieldName).Dot("Do").Call(
		jen.Id("ctx"),
		jen.Id(requestParamsVar).Op("..."))
	returnErr := jen.ReturnFunc(func(returnVals *jen.Group) {
		returnVar(returnVals)
		returnVals.Add(snip.WerrorWrapContext()).Call(
			jen.Id("ctx"),
			jen.Err(),
			jen.Lit(fmt.Sprintf("%s failed", endpointDef.EndpointName)),
		)
	})

	if returnsBinary {
		methodBody.List(jen.Id(respVar), jen.Err()).Op(":=").Add(callStmt)
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErr)
		if returnsOptional {
			// If an endpoint with a return type of optional<binary> provides a response with a code of StatusNoContent
			// then the return value is empty and nil is returned.
			methodBody.If(jen.Id(respVar).Dot("StatusCode").Op("==").Add(snip.HTTPStatusNoContent())).Block(
				jen.Return(jen.Nil(), jen.Nil()),
			)
			methodBody.Return(jen.Op("&").Id(respVar).Dot("Body"), jen.Nil())
		} else {
			// if endpoint returns binary, return body of response directly
			methodBody.Return(jen.Id(respVar).Dot("Body"), jen.Nil())
		}
		return
	}

	methodBody.If(
		jen.List(jen.Id("_"), jen.Err()).Op(":=").Add(callStmt),
		jen.Err().Op("!=").Nil(),
	).Block(returnErr)

	if !returnsOptional && (returnDefaultValue || returnsCollection) {
		// verify that return value is non-nil and dereference
		methodBody.If(jen.Id(returnValVar).Op("==").Nil()).Block(jen.ReturnFunc(func(returnVals *jen.Group) {
			returnVar(returnVals)
			returnVals.Add(snip.WerrorErrorContext()).Call(
				jen.Id("ctx"),
				jen.Lit(fmt.Sprintf("%s response cannot be nil", endpointDef.EndpointName)),
			)
		}))
	}

	if returnDefaultValue {
		methodBody.Return(jen.Op("*").Id(returnValVar), jen.Nil())
	} else if hasReturnVal {
		methodBody.Return(jen.Id(returnValVar), jen.Nil())
	} else {
		methodBody.Return(jen.Nil())
	}
}

func astForEndpointMethodBodyRequestParams(methodBody *jen.Group, endpointDef *types.EndpointDefinition) {
	methodBody.Var().Id(requestParamsVar).Op("[]").Add(snip.CGRClientRequestParam())

	// helper for the statement "requestParams = append(requestParams, {code})"
	appendRequestParams := func(methodBody *jen.Group, code jen.Code) {
		methodBody.Id(requestParamsVar).Op("=").Append(jen.Id(requestParamsVar), code)
	}

	appendRequestParams(methodBody, snip.CGRClientWithRPCMethodName().Call(jen.Lit(transforms.Export(endpointDef.EndpointName))))
	appendRequestParams(methodBody, snip.CGRClientWithRequestMethod().Call(jen.Lit(endpointDef.HTTPMethod.String())))
	// auth params
	if endpointDef.HeaderAuth {
		appendRequestParams(methodBody, snip.CGRClientWithHeader().Call(jen.Lit("Authorization"),
			snip.FmtSprint().Call(jen.Lit("Bearer "), jen.Id(authHeaderVar)),
		))
	} else if endpointDef.CookieAuth != nil {
		appendRequestParams(methodBody, snip.CGRClientWithHeader().Call(jen.Lit("Cookie"),
			snip.FmtSprint().Call(jen.Lit(*endpointDef.CookieAuth+"="), jen.Id(cookieTokenVar)),
		))
	}
	// path params
	appendRequestParams(methodBody, snip.CGRClientWithPathf().CallFunc(func(args *jen.Group) {
		args.Lit(pathParamRegexp.ReplaceAllString(endpointDef.HTTPPath, regexp.QuoteMeta(`%s`)))
		for _, param := range endpointDef.PathParams() {
			args.Add(snip.URLPathEscape()).Call(snip.FmtSprint().Call(jen.Id(transforms.ArgName(param.ParamID))))
		}
	}))
	// body params
	if body := endpointDef.BodyParam(); body != nil {
		bodyArg := transforms.ArgName(body.Name)
		if body.Type.IsOptional() {
			bodyVal := jen.Id(bodyArg)
			if body.Type.IsNamed() && !body.Type.IsBinary() {
				// If the response type is named (i.e. an alias), check the inner Value field for absence.
				bodyVal = bodyVal.Dot("Value")
			}
			methodBody.If(bodyVal.Clone().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
				if body.Type.IsBinary() {
					appendRequestParams(ifBody, snip.CGRClientWithRawRequestBodyProvider().Call(jen.Id(bodyArg)))
				} else {
					appendRequestParams(ifBody, snip.CGRClientWithJSONRequest().Call(jen.Id(bodyArg)))
				}
			})
		} else if body.Type.IsBinary() {
			appendRequestParams(methodBody, snip.CGRClientWithRawRequestBodyProvider().Call(jen.Id(bodyArg)))
		} else {
			appendRequestParams(methodBody, snip.CGRClientWithJSONRequest().Call(jen.Id(bodyArg)))
		}
	}
	// header params
	for _, param := range endpointDef.HeaderParams() {
		argName := transforms.ArgName(param.Name)
		if param.Type.IsOptional() {
			selector := jen.Id(argName).Clone
			if _, isAlias := param.Type.(*types.AliasType); isAlias {
				selector = selector().Dot("Value").Clone
			}
			// if header parameter type is an optional, append dereferenced value if it is non-nil
			methodBody.If(selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
				appendRequestParams(ifBody, snip.CGRClientWithHeader().Call(
					jen.Lit(param.ParamID),
					snip.FmtSprint().Call(jen.Op("*").Add(selector())),
				))
			})
		} else {
			appendRequestParams(methodBody, snip.CGRClientWithHeader().Call(
				jen.Lit(param.ParamID),
				snip.FmtSprint().Call(jen.Id(argName)),
			))
		}
	}
	// query params
	if queryParams := endpointDef.QueryParams(); len(queryParams) > 0 {
		methodBody.Id(queryParamsVar).Op(":=").Make(snip.URLValues())
		for _, param := range endpointDef.QueryParams() {
			argName := transforms.ArgName(param.Name)
			if param.Type.IsOptional() {
				selector := jen.Id(argName).Clone
				if _, isAlias := param.Type.(*types.AliasType); isAlias {
					selector = selector().Dot("Value").Clone
				}
				methodBody.If(selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
					ifBody.Id(queryParamsVar).Dot("Set").Call(jen.Lit(param.ParamID),
						snip.FmtSprint().Call(jen.Op("*").Add(selector())))
				})
			} else if param.Type.IsList() {
				methodBody.For(jen.List(jen.Id("_"), jen.Id("v")).Op(":=").Range().Id(argName)).Block(
					jen.Id(queryParamsVar).Dot("Add").Call(jen.Lit(param.ParamID), snip.FmtSprint().Call(jen.Id("v"))),
				)
			} else {
				methodBody.Id(queryParamsVar).Dot("Set").Call(jen.Lit(param.ParamID),
					snip.FmtSprint().Call(jen.Id(argName)))
			}
		}
		appendRequestParams(methodBody, snip.CGRClientWithQueryValues().Call(jen.Id(queryParamsVar)))
	}
	// response
	if endpointDef.Returns != nil {
		if (*endpointDef.Returns).IsBinary() {
			appendRequestParams(methodBody, snip.CGRClientWithRawResponseBody().Call())
		} else {
			appendRequestParams(methodBody, snip.CGRClientWithJSONResponse().Call(jen.Op("&").Id(returnValVar)))
		}
	}
}

func astForEndpointAuthMethodBodyFunc(methodBody *jen.Group, endpointDef *types.EndpointDefinition) {
	methodBody.Return(
		jen.Id(clientReceiverName).
			Dot(clientStructFieldName).
			Dot(transforms.Export(endpointDef.EndpointName)).
			CallFunc(func(args *jen.Group) {
				args.Id("ctx")
				if endpointDef.HeaderAuth {
					args.Id(clientReceiverName).Dot(authHeaderVar)
				}
				if endpointDef.CookieAuth != nil {
					args.Id(clientReceiverName).Dot(cookieTokenVar)
				}
				for _, param := range endpointDef.Params {
					args.Id(transforms.ArgName(param.Name))
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

func astForTokenServiceEndpointMethodBody(methodBody *jen.Group, endpointDef *types.EndpointDefinition, hasAuth bool) {
	if hasAuth {
		if endpointDef.Returns != nil {
			returnsType := *endpointDef.Returns
			argType := returnsType.Code()
			if returnsType.IsBinary() {
				// special case: "binary" types resolve to []byte, but this indicates a streaming parameter when
				// specified as the request argument of a service, so use "io.ReadCloser".
				// If the type is optional<binary>, use "*io.ReadCloser".
				if returnsType.IsOptional() {
					argType = jen.Op("*").Add(snip.IOReadCloser())
				} else {
					argType = snip.IOReadCloser()
				}
			}
			methodBody.Var().Id(defaultReturnValVar).Add(argType)
		}
		methodBody.List(jen.Id("token"), jen.Err()).Op(":=").Id(clientReceiverName).Dot(tokenProviderVar).Call(jen.Id("ctx"))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(jen.ReturnFunc(func(returns *jen.Group) {
			if endpointDef.Returns != nil {
				returns.Id(defaultReturnValVar)
			}
			returns.Err()
		}))
	}
	methodBody.Return(jen.Id(clientReceiverName).Dot(clientStructFieldName).Dot(transforms.Export(endpointDef.EndpointName)).
		CallFunc(func(args *jen.Group) {
			args.Id("ctx")
			if hasAuth {
				args.Add(types.Bearertoken{}.Code()).Call(jen.Id("token"))
			}
			for _, paramDef := range endpointDef.Params {
				args.Id(transforms.ArgName(paramDef.Name))
			}
		}),
	)
}

func astForTokenServiceEndpointMethod(serviceName string, endpointDef *types.EndpointDefinition) *jen.Statement {
	hasAuth := endpointDef.HeaderAuth || endpointDef.CookieAuth != nil
	return jen.Func().
		Params(jen.Id(clientReceiverName).Op("*").Id(withTokenProviderName(clientStructTypeName(serviceName)))).
		Id(transforms.Export(endpointDef.EndpointName)).
		ParamsFunc(func(args *jen.Group) {
			astForEndpointArgsFunc(args, endpointDef, hasAuth, false)
		}).
		ParamsFunc(func(args *jen.Group) {
			astForEndpointReturnsFunc(args, endpointDef)
		}).
		BlockFunc(func(methodBody *jen.Group) {
			astForTokenServiceEndpointMethodBody(methodBody, endpointDef, hasAuth)
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
