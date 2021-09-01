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

func writeServiceType(file *jen.Group, serviceDef *types.ServiceDefinition, cfg OutputConfiguration) {
	file.Add(astForServiceInterface(serviceDef, false, false))
	file.Add(astForClientStructDecl(serviceDef.Name))
	file.Add(astForNewClientFunc(serviceDef.Name))
	for _, endpointDef := range serviceDef.Endpoints {
		file.Add(astForEndpointMethod(serviceDef.Name, endpointDef, false, cfg.LiteralJSONMethods))
	}
	if serviceDef.HasHeaderAuth() || serviceDef.HasCookieAuth() {
		// at least one endpoint uses authentication: define decorator structures
		file.Add(astForServiceInterface(serviceDef, true, false))
		file.Add(astForNewServiceFuncWithAuth(serviceDef))
		file.Add(astForClientStructDeclWithAuth(serviceDef))
		for _, endpointDef := range serviceDef.Endpoints {
			file.Add(astForEndpointMethod(serviceDef.Name, endpointDef, true, cfg.LiteralJSONMethods))
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
				methods.Add(endpointDef.Docs.CommentLine()).Id(transforms.Export(endpointDef.EndpointName)).
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

func astForEndpointMethod(serviceName string, endpointDef *types.EndpointDefinition, withAuth bool, litJSON bool) *jen.Statement {
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
				astForEndpointMethodBodyFunc(methodBody, endpointDef, litJSON)
			}
		})
}

func astForEndpointMethodBodyFunc(methodBody *jen.Group, endpointDef *types.EndpointDefinition, litJSON bool) {
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
	returnVar := func(returns *jen.Group) {
		switch {
		case returnDefaultValue:
			returns.Id(defaultReturnValVar)
		case hasReturnVal:
			returns.Nil()
		}
	}

	// if endpoint returns a value, declare variables for value
	if hasReturnVal && !returnsBinary {
		if returnDefaultValue {
			methodBody.Var().Id(defaultReturnValVar).Add((*endpointDef.Returns).Code())
			methodBody.Var().Id(returnValVar).Op("*").Add((*endpointDef.Returns).Code())
		} else {
			methodBody.Var().Id(returnValVar).Add((*endpointDef.Returns).Code())
		}
	}
	// build requestParams
	astForEndpointMethodBodyRequestParams(methodBody, endpointDef, litJSON)

	// execute request
	callStmt := jen.Id(clientReceiverName).Dot(clientStructFieldName).Dot("Do").Call(
		jen.Id("ctx"),
		jen.Id(requestParamsVar).Op("..."))
	returnErr := jen.ReturnFunc(func(returns *jen.Group) {
		returnVar(returns)
		returns.Add(snip.WerrorWrapContext()).Call(
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
	if returnDefaultValue {
		if _, returnsAliasType := (*endpointDef.Returns).(*types.AliasType); !(returnsOptional && returnsAliasType) {
			// verify that return value is non-nil and dereference
			methodBody.If(jen.Id(returnValVar).Op("==").Nil()).Block(jen.ReturnFunc(func(returns *jen.Group) {
				returnVar(returns)
				returns.Add(snip.WerrorErrorContext()).Call(
					jen.Id("ctx"),
					jen.Lit(fmt.Sprintf("%s response cannot be nil", endpointDef.EndpointName)),
				)
			}))
		}
	}

	if returnDefaultValue {
		methodBody.Return(jen.Op("*").Id(returnValVar), jen.Nil())
	} else if hasReturnVal {
		methodBody.Return(jen.Id(returnValVar), jen.Nil())
	} else {
		methodBody.Return(jen.Nil())
	}
}

func astForEndpointMethodBodyRequestParams(methodBody *jen.Group, endpointDef *types.EndpointDefinition, litJSON bool) {
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
			args.Add(snip.URLPathEscape()).Call(snip.FmtSprint().Call(jen.Id(argNameTransform(param.ParamID))))
		}
	}))
	// body params
	if body := endpointDef.BodyParam(); body != nil {
		bodyArg := argNameTransform(body.Name)
		if body.Type.IsBinary() {
			appendRequestParams(methodBody, snip.CGRClientWithRawRequestBodyProvider().Call(jen.Id(bodyArg)))
		} else {
			if litJSON {
				appendRequestParams(methodBody, snip.CGRClientWithRequestAppendFunc().CallFunc(func(args *jen.Group) {
					args.Add(snip.CGRCodecsJSON()).Dot("ContentType").Call()
					if body.Type.IsNamed() {
						args.Id(bodyArg).Dot("AppendJSON")
					} else {
						args.Func().
							Params(jen.Id("out").Op("[]").Byte()).
							Params(jen.Op("[]").Byte(), jen.Error()).
							BlockFunc(func(funcBody *jen.Group) {
								encoding.AnonFuncBodyAppendJSON(funcBody, jen.Id(bodyArg).Clone, body.Type)
							})
					}
				}))
			} else {
				appendRequestParams(methodBody, snip.CGRClientWithJSONRequest().Call(jen.Id(bodyArg)))
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
			argName := argNameTransform(param.Name)
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
			if litJSON {
				appendRequestParams(methodBody, snip.CGRClientWithResponseUnmarshalFunc().CallFunc(func(args *jen.Group) {
					args.Add(snip.CGRCodecsJSON()).Dot("Accept").Call()
					if (*endpointDef.Returns).IsNamed() {
						args.Id(returnValVar).Dot("UnmarshalJSON")
					} else {
						args.Func().
							Params(jen.Id("data").Op("[]").Byte()).
							Params(jen.Op("[]").Byte(), jen.Error()).
							BlockFunc(func(funcBody *jen.Group) {
								encoding.AnonFuncBodyUnmarshalJSON(funcBody, jen.Id(returnValVar).Clone, *endpointDef.Returns, false)
							})
					}
				}))
			} else {
				appendRequestParams(methodBody, snip.CGRClientWithJSONResponse().Call(jen.Op("&").Id(returnValVar)))
			}
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
					args.Id(argNameTransform(param.Name))
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
		ParamsFunc(func(args *jen.Group) {
			astForEndpointArgsFunc(args, endpointDef, true, false)
		}).
		ParamsFunc(func(args *jen.Group) {
			astForEndpointReturnsFunc(args, endpointDef)
		}).
		BlockFunc(func(methodBody *jen.Group) {
			if hasAuth {
				if endpointDef.Returns != nil {
					methodBody.Var().Id(defaultReturnValVar).Add((*endpointDef.Returns).Code())
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
						args.Id(argNameTransform(paramDef.Name))
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
