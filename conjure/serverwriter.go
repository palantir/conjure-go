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
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/encoding"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	implName = "impl"

	// Handler
	handlerName = "handler"

	// Router
	routerVarName     = "router"
	resourceName      = "resource"
	pathParamsVarName = "pathParams"

	// Handler
	handlerStructNameSuffix = "Handler"

	// ResponseWriter
	responseWriterVarName = "rw"
	responseArgVarName    = "respArg"

	// Request
	reqName = "req"
)

func writeServerType(file *jen.Group, serviceDef *types.ServiceDefinition) {
	file.Add(astForServiceInterface(serviceDef, false, true))
	file.Add(astForRouteRegistration(serviceDef))
	file.Add(astForHandlerStructDecl(serviceDef.Name))
	file.Add(astForHandlerMethods(serviceDef))
}

func astForRouteRegistration(serviceDef *types.ServiceDefinition) *jen.Statement {
	funcName := routeRegistrationFuncName(serviceDef.Name)
	ifaceType := transforms.Export(serviceDef.Name)
	return jen.
		Commentf("%s registers handlers for the %s endpoints with a witchcraft wrouter.", funcName, serviceDef.Name).Line().
		Comment("This should typically be called in a witchcraft server's InitFunc.").Line().
		Comment("impl provides an implementation of each endpoint, which can assume the request parameters have been parsed").Line().
		Comment("in accordance with the Conjure specification.").Line().
		Func().Id(funcName).
		Params(jen.Id(routerVarName).Add(snip.WrouterRouter()), jen.Id(implName).Id(ifaceType)).
		Params(jen.Error()).
		BlockFunc(func(methodBody *jen.Group) {
			// Create the handler struct
			methodBody.Id(handlerName).Op(":=").Id(handlerStuctName(serviceDef.Name)).Values(jen.Id(implName).Op(":").Id(implName))
			// Create the witchcraft resource
			methodBody.Id(resourceName).Op(":=").Add(snip.WresourceNew()).Call(jen.Lit(strings.ToLower(serviceDef.Name)), jen.Id(routerVarName))
			// For each endpoint, register a route on the provided router
			// if err := resource.Get(...); err != nil {
			//     return werror.Wrap(err, ...)
			// }
			for _, endpointDef := range serviceDef.Endpoints {
				methodBody.If(
					jen.Err().Op(":=").Id(resourceName).Dot(wresourceMethod(endpointDef.HTTPMethod)).CallFunc(func(args *jen.Group) {
						astForWrouterRegisterArgsFunc(args, endpointDef)
						args.Line()
					}),
					jen.Err().Op("!=").Nil(),
				).Block(
					jen.Return(snip.WerrorWrapContext()).Call(
						jen.Add(snip.ContextTODO()).Call(), jen.Err(), jen.Lit(fmt.Sprintf("failed to add %s route", endpointDef.EndpointName))),
				)
			}
			// Return nil if everything registered
			methodBody.Return(jen.Nil())
		})
}

func astForWrouterRegisterArgsFunc(args *jen.Group, endpointDef *types.EndpointDefinition) {
	args.Line().Lit(strings.Title(endpointDef.EndpointName))
	args.Line().Lit(endpointDef.HTTPPath)
	args.Line().Add(snip.CGRHTTPServerNewJSONHandler()).Call(
		jen.Id(handlerName).Dot(handleFuncName(endpointDef.EndpointName)),
		snip.CGRHTTPServerStatusCodeMapper(),
		snip.CGRHTTPServerErrHandler(),
	)
	for _, argDef := range endpointDef.PathParams() {
		for _, marker := range argDef.Markers {
			if isSafeMarker(marker) {
				args.Line().Add(snip.WrouterSafePathParams()).Call(jen.Lit(argDef.ParamID))
			}
		}
	}
	for _, argDef := range endpointDef.HeaderParams() {
		for _, marker := range argDef.Markers {
			if isSafeMarker(marker) {
				args.Line().Add(snip.WrouterSafeHeaderParams()).Call(jen.Lit(argDef.ParamID))
			}
		}
	}
	for _, argDef := range endpointDef.QueryParams() {
		for _, marker := range argDef.Markers {
			if isSafeMarker(marker) {
				args.Line().Add(snip.WrouterSafeQueryParams()).Call(jen.Lit(argDef.ParamID))
			}
		}
	}
}

func astForHandlerStructDecl(serviceName string) *jen.Statement {
	return jen.Type().Id(handlerStuctName(serviceName)).Struct(jen.Id(implName).Id(serviceName))
}

func astForHandlerMethods(serviceDef *types.ServiceDefinition) *jen.Statement {
	stmt := jen.Empty()
	for _, endpointDef := range serviceDef.Endpoints {
		stmt = stmt.Func().
			Params(jen.Id(handlerReceiverName(serviceDef.Name)).Op("*").Id(handlerStuctName(serviceDef.Name))).
			Id(handleFuncName(endpointDef.EndpointName)).
			ParamsFunc(func(params *jen.Group) {
				if endpointDef.Returns != nil {
					params.Id(responseWriterVarName).Add(snip.HTTPResponseWriter())
				} else {
					params.Id("_").Add(snip.HTTPResponseWriter())
				}
				params.Id(reqName).Op("*").Add(snip.HTTPRequest())
			}).
			Params(jen.Error()).
			BlockFunc(func(methodBody *jen.Group) {
				astForHandlerMethodBody(methodBody, serviceDef.Name, endpointDef)
			}).
			Line()
	}
	return stmt
}

func astForHandlerMethodBody(methodBody *jen.Group, serviceName string, endpointDef *types.EndpointDefinition) {
	// decode auth header
	astForHandlerMethodAuthParams(methodBody, endpointDef)
	// decode arguments
	astForHandlerMethodPathParams(methodBody, endpointDef.PathParams())
	astForHandlerMethodQueryParams(methodBody, endpointDef.QueryParams())
	astForHandlerMethodHeaderParams(methodBody, endpointDef.HeaderParams())
	astForHandlerMethodDecodeBody(methodBody, endpointDef.BodyParam())
	// call impl handler & return
	astForHandlerExecImplAndReturn(methodBody, serviceName, endpointDef)
}

func astForHandlerMethodAuthParams(methodBody *jen.Group, endpointDef *types.EndpointDefinition) {
	switch {
	case endpointDef.HeaderAuth:
		//	authHeader, err := httpserver.ParseBearerTokenHeader(req)
		//	if err != nil {
		//		return errors.WrapWithPermissionDenied(err)
		//	}
		methodBody.List(jen.Id(authHeaderVar), jen.Err()).Op(":=").
			Add(snip.CGRHTTPServerParseBearerTokenHeader()).Call(jen.Id(reqName))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(jen.Return(snip.CGRErrorsWrapWithPermissionDenied().Call(jen.Err())))
	case endpointDef.CookieAuth != nil:
		//	authCookie, err := req.Cookie("PALANTIR_TOKEN")
		//	if err != nil {
		//		return errors.WrapWithPermissionDenied(err)
		//	}
		//	cookieToken := bearertoken.Token(authCookie.Value)
		methodBody.List(jen.Id("authCookie"), jen.Err()).Op(":=").Id(reqName).Dot("Cookie").Call(jen.Lit(*endpointDef.CookieAuth))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(jen.Return(snip.CGRErrorsWrapWithPermissionDenied().Call(jen.Err())))
		methodBody.Id(cookieTokenVar).Op(":=").Add(types.Bearertoken{}.Code()).Call(jen.Id("authCookie").Dot("Value"))
	}
}

func astForHandlerMethodPathParams(methodBody *jen.Group, pathParams []*types.EndpointArgumentDefinition) {
	if len(pathParams) == 0 {
		return
	}
	methodBody.Id(pathParamsVarName).Op(":=").Add(snip.WrouterPathParams()).Call(jen.Id(reqName))
	methodBody.If(jen.Id(pathParamsVarName).Op("==").Nil()).Block(jen.Return(snip.WerrorWrapContext().Call(
		reqContext(),
		snip.CGRErrorsNewInternal().Call(),
		jen.Lit("path params not found on request: ensure this endpoint is registered with wrouter"),
	)))
	for _, argDef := range pathParams {
		astForHandlerMethodPathParam(methodBody, argDef)
	}
}

func astForHandlerMethodPathParam(methodBody *jen.Group, argDef *types.EndpointArgumentDefinition) {
	strVar := transforms.SafeName(argDef.ParamID) + "Str"
	switch argDef.Type.(type) {
	case types.Any, types.String:
		strVar = transforms.SafeName(argDef.ParamID)
	}
	// For each path param, pull out the value and check if it is present in the map
	// argNameStr, ok := pathParams["argName"]; if !ok { werror... }
	methodBody.List(jen.Id(strVar), jen.Id("ok")).Op(":=").Id(pathParamsVarName).Index(jen.Lit(argDef.ParamID))
	methodBody.If(jen.Op("!").Id("ok")).Block(jen.Return(
		snip.WerrorWrapContext().Call(
			reqContext(),
			snip.CGRErrorsNewInvalidArgument().Call(),
			jen.Lit(fmt.Sprintf("path parameter %q not present", argDef.ParamID))),
	))
	// type-specific unmarshal behavior
	switch argDef.Type.(type) {
	case types.Any, types.String:
	default:
		argDescriptor := fmt.Sprintf("path[%q]", argDef.ParamID)
		astForDecodeHTTPParam(methodBody, argDef.Type, transforms.SafeName(argDef.Name), jen.Id(strVar).Clone, argDescriptor)
	}
}

func astForHandlerMethodHeaderParams(methodBody *jen.Group, headerParams []*types.EndpointArgumentDefinition) {
	for _, arg := range headerParams {
		astForHandlerMethodHeaderParam(methodBody, arg)
	}
}

func astForHandlerMethodHeaderParam(methodBody *jen.Group, argDef *types.EndpointArgumentDefinition) {
	var headerVar func() *jen.Statement
	switch argDef.Type.(type) {
	case *types.List:
		headerVar = jen.Id(reqName).Dot("Header").Dot("Values").Call(jen.Lit(argDef.ParamID)).Clone
	default:
		headerVar = jen.Id(reqName).Dot("Header").Dot("Get").Call(jen.Lit(argDef.ParamID)).Clone
	}
	argDescriptor := fmt.Sprintf("header[%q]", argDef.ParamID)
	astForDecodeHTTPParam(methodBody, argDef.Type, transforms.SafeName(argDef.Name), headerVar, argDescriptor)
}

func astForHandlerMethodQueryParams(methodBody *jen.Group, queryParams []*types.EndpointArgumentDefinition) {
	for _, arg := range queryParams {
		astForHandlerMethodQueryParam(methodBody, arg)
	}
}

func astForHandlerMethodQueryParam(methodBody *jen.Group, argDef *types.EndpointArgumentDefinition) {
	var queryVar func() *jen.Statement
	switch argDef.Type.(type) {
	case *types.List:
		queryVar = jen.Id(reqName).Dot("URL").Dot("Query").Call().Index(jen.Lit(argDef.ParamID)).Clone
	default:
		queryVar = jen.Id(reqName).Dot("URL").Dot("Query").Call().Dot("Get").Call(jen.Lit(argDef.ParamID)).Clone
	}
	argDescriptor := fmt.Sprintf("query[%q]", argDef.ParamID)
	astForDecodeHTTPParam(methodBody, argDef.Type, transforms.SafeName(argDef.Name), queryVar, argDescriptor)
}

func astForHandlerMethodDecodeBody(methodBody *jen.Group, argDef *types.EndpointArgumentDefinition) {
	if argDef == nil {
		return
	}
	if argDef.Type.IsBinary() {
		// If the body argument is binary, pass req.Body directly to the impl.
		methodBody.Id(transforms.SafeName(argDef.Name)).Op(":=").Id(reqName).Dot("Body")
		return
	}
	// If the request is not binary, it is JSON. Unmarshal the req.Body.
	methodBody.Var().Id(transforms.SafeName(argDef.Name)).Add(argDef.Type.Code())
	methodBody.If(
		jen.Err().Op(":=").Add(snip.CGRCodecsJSON().Dot("Decode")).Call(
			jen.Id(reqName).Dot("Body"),
			jen.Op("&").Id(transforms.SafeName(argDef.Name)),
		),
		jen.Err().Op("!=").Nil(),
	).Block(jen.Return(snip.CGRErrorsWrapWithInvalidArgument().Call(jen.Err())))
}

func astForDecodeHTTPParam(
	methodBody *jen.Group,
	argType types.Type,
	outVarName string,
	inStr func() *jen.Statement,
	argDescriptor string,
) {
	if argType.IsList() {
		encoding.UnmarshalStringListStatements(methodBody, argType, outVarName, inStr, argDescriptor, reqContext)
	} else {
		encoding.UnmarshalStringStatements(methodBody, argType, outVarName, inStr, argDescriptor, reqContext)
	}
}

func astForHandlerExecImplAndReturn(methodBody *jen.Group, serviceName string, endpointDef *types.EndpointDefinition) {
	callFunc := jen.Id(handlerReceiverName(serviceName)).Dot(implName).Dot(strings.Title(endpointDef.EndpointName)).CallFunc(func(args *jen.Group) {
		args.Add(reqContext())
		if endpointDef.HeaderAuth {
			args.Add(snip.BearerTokenToken()).Call(jen.Id(authHeaderVar))
		} else if endpointDef.CookieAuth != nil {
			args.Id(cookieTokenVar)
		}
		for _, paramDef := range endpointDef.Params {
			args.Id(transforms.SafeName(paramDef.Name))
		}
	})

	if endpointDef.Returns == nil {
		// The endpoint doesn't return anything, just return the interface call
		methodBody.Return(callFunc)
		return
	}

	methodBody.List(jen.Id(responseArgVarName), jen.Err()).Op(":=").Add(callFunc)
	methodBody.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))

	respArg := jen.Id(responseArgVarName)
	codec := snip.CGRCodecsJSON()
	if (*endpointDef.Returns).IsBinary() {
		if (*endpointDef.Returns).IsOptional() {
			// Empty binaries return a 204 (No Content) response
			methodBody.If(respArg.Clone().Op("==").Nil()).Block(
				jen.Id(responseWriterVarName).Dot("WriteHeader").Call(snip.HTTPStatusNoContent()),
				jen.Return(jen.Nil()),
			)
			respArg = jen.Op("*").Add(respArg.Clone())
		}
		codec = snip.CGRCodecsBinary()
	} else {
		if !(*endpointDef.Returns).IsNamed() {
			// If we have an unnamed type, wrap marshal logic in safejson.AppendFunc.
			respArg = snip.SafeJSONAppendFunc().Call(jen.Func().
				Params(jen.Id("out").Op("[]").Byte()).
				Params(jen.Op("[]").Byte(), jen.Error()).
				BlockFunc(func(funcBody *jen.Group) {
					encoding.AnonFuncBodyAppendJSON(funcBody, respArg.Clone, *endpointDef.Returns)
				}))
		}
	}

	methodBody.Id(responseWriterVarName).Dot("Header").Call().Dot("Add").Call(
		jen.Lit("Content-Type"),
		codec.Clone().Dot("ContentType").Call(),
	)
	methodBody.Return(codec.Clone().Dot("Encode").Call(jen.Id(responseWriterVarName), respArg.Clone()))
}

func routeRegistrationFuncName(serviceName string) string {
	return "RegisterRoutes" + strings.Title(serviceName)
}

func handlerStuctName(serviceName string) string {
	firstCharLower := strings.ToLower(string(serviceName[0]))
	return strings.Join([]string{firstCharLower, serviceName[1:], handlerStructNameSuffix}, "")
}

func handlerReceiverName(serviceName string) string {
	return strings.ToLower(string(serviceName[0]))
}

func handleFuncName(endpointName string) string {
	return "Handle" + strings.Title(endpointName)
}

func wresourceMethod(method spec.HttpMethod) string {
	switch method.Value() {
	case spec.HttpMethod_GET:
		return "Get"
	case spec.HttpMethod_POST:
		return "Post"
	case spec.HttpMethod_PUT:
		return "Put"
	case spec.HttpMethod_DELETE:
		return "Delete"
	default:
		panic("Unexpected http method " + method.String())
	}
}

func isSafeMarker(marker types.Type) bool {
	ext, ok := marker.(*types.External)
	if !ok {
		return false
	}
	return ext.Spec.Package == "com.palantir.logsafe" && ext.Spec.Name == "Safe"
}

func reqContext() *jen.Statement {
	return jen.Id(reqName).Dot("Context").Call()
}
