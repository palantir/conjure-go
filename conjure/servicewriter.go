// Copyright (c) 2018 Palantir Technologies. All rights reserved.
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
	"go/token"
	"regexp"

	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	"github.com/palantir/goastwriter/statement"
	"github.com/pkg/errors"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/conjure-go/v6/conjure/visitors"
)

const (
	clientStructFieldName = "client"

	receiverName     = "c"
	ctxName          = "ctx"
	wrappedClientVar = "client"
	authHeaderVar    = "authHeader"
	cookieTokenVar   = "cookieToken"

	httpClientImportPath = "github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	httpClientClientType = expression.Type("httpclient.Client")
	httpClientPkgName    = "httpclient"

	tokenProviderVar        = "tokenProvider"
	tokenProviderType       = expression.Type("httpclient.TokenProvider")
	tokenProviderImportPath = "github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
)

type serviceASTConfig struct {
	withAuth          bool
	withTokenProvider bool
}

func astForService(serviceDefinition spec.ServiceDefinition, info types.PkgInfo) ([]astgen.ASTDecl, StringSet, error) {
	allImports := NewStringSet()
	serviceName := serviceDefinition.ServiceName.Name
	interfaceAST, imports, err := clientServiceInterfaceAST(serviceDefinition, info, serviceASTConfig{})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to generate interface for service %q", serviceName)
	}
	allImports.AddAll(imports)

	serviceNewFunc, imports := serviceNewFuncAST(serviceName)
	allImports.AddAll(imports)

	serviceStruct := decl.NewStruct(clientStructTypeName(serviceName), []*expression.StructField{
		{
			Name: clientStructFieldName,
			Type: httpClientClientType,
		},
	}, "")
	allImports[httpClientImportPath] = struct{}{}

	methodsAST, imports, err := serviceStructMethodsAST(serviceDefinition, info)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to generate methods for service %q", serviceName)
	}
	allImports.AddAll(imports)
	components := []astgen.ASTDecl{
		interfaceAST,
		serviceStruct,
		serviceNewFunc,
	}
	components = append(components, methodsAST...)

	hasHeaderAuth, hasCookieAuth, err := hasAuth(serviceDefinition.Endpoints)
	if err != nil {
		return nil, nil, err
	}

	if hasHeaderAuth || hasCookieAuth {
		// at least one endpoint uses authentication: define decorator structures
		withAuthInterfaceAST, imports, err := clientServiceInterfaceAST(serviceDefinition, info, serviceASTConfig{withAuth: true})
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to generate interface with auth for service %q", serviceName)
		}
		components = append(components, withAuthInterfaceAST)
		allImports.AddAll(imports)

		withAuthServiceNewFunc, authServiceNewFuncImports := withAuthServiceNewFuncAST(serviceName, hasHeaderAuth, hasCookieAuth, info)
		components = append(components, withAuthServiceNewFunc)
		allImports.AddAll(authServiceNewFuncImports)

		withAuthServiceStruct, authServiceStructImports := withAuthServiceStructAST(serviceName, hasHeaderAuth, hasCookieAuth, info)
		components = append(components, withAuthServiceStruct)
		allImports.AddAll(authServiceStructImports)

		withAuthMethodsAST, imports, err := withAuthServiceStructMethodsAST(serviceDefinition, info)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to generate methods with auth for service %q", serviceName)
		}
		components = append(components, withAuthMethodsAST...)
		allImports.AddAll(imports)
	}

	if canAddTokenInterface(serviceDefinition.Endpoints) {
		withTokenServiceNewFunc, tokenServiceNewFuncImports := withTokenServiceNewFuncAST(serviceName, info)
		components = append(components, withTokenServiceNewFunc)
		allImports.AddAll(tokenServiceNewFuncImports)

		withTokenServiceStruct, tokenServiceStructImports := withTokenServiceStructAST(serviceName, info)
		components = append(components, withTokenServiceStruct)
		allImports.AddAll(tokenServiceStructImports)

		withTokenMethodsAST, imports, err := withTokenServiceStructMethodsAST(serviceDefinition, info)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Failed to generate methods with token provider for service %q", serviceName)
		}
		components = append(components, withTokenMethodsAST...)
		allImports.AddAll(imports)
	}

	return components, allImports, nil
}

func hasAuth(endpoints []spec.EndpointDefinition) (hasHeaderAuth, hasCookieAuth bool, err error) {
	for _, endpointDefinition := range endpoints {
		if endpointDefinition.Auth == nil {
			continue
		}
		if possibleHeaderAuth, err := visitors.GetPossibleHeaderAuth(*endpointDefinition.Auth); err != nil {
			return false, false, err
		} else if possibleHeaderAuth != nil {
			hasHeaderAuth = true
		}
		if possibleCookieAuth, err := visitors.GetPossibleCookieAuth(*endpointDefinition.Auth); err != nil {
			return false, false, err
		} else if possibleCookieAuth != nil {
			hasCookieAuth = true
		}
	}
	return
}

// Return true if all endpoints that require authentication are of the same auth type (header or cookie) and at least
// one endpoint has auth. The same auth type is required because a single token provider will likely not be useful for
// both auth types.
func canAddTokenInterface(endpoints []spec.EndpointDefinition) bool {
	var hasHeaderAuth, hasCookieAuth bool
	for _, endpointDefinition := range endpoints {
		if endpointDefinition.Auth == nil {
			continue
		}
		possibleHeaderAuth, err := visitors.GetPossibleHeaderAuth(*endpointDefinition.Auth)
		if err != nil {
			return false
		}
		hasHeaderAuth = hasHeaderAuth || possibleHeaderAuth != nil
		possibleCookieAuth, err := visitors.GetPossibleCookieAuth(*endpointDefinition.Auth)
		if err != nil {
			return false
		}
		hasCookieAuth = hasCookieAuth || possibleCookieAuth != nil
		if hasHeaderAuth && hasCookieAuth {
			return false
		}
	}
	return hasHeaderAuth || hasCookieAuth
}

type generatorType bool

const (
	generatorTypeClient generatorType = true
	generatorTypeServer generatorType = false
)

func clientServiceInterfaceAST(serviceDefinition spec.ServiceDefinition, info types.PkgInfo, config serviceASTConfig) (astgen.ASTDecl, StringSet, error) {
	return serviceInterfaceAST(serviceDefinition, info, config, generatorTypeClient)
}

func serverServiceInterfaceAST(serviceDefinition spec.ServiceDefinition, info types.PkgInfo, config serviceASTConfig) (astgen.ASTDecl, StringSet, error) {
	return serviceInterfaceAST(serviceDefinition, info, config, generatorTypeServer)
}

func serviceInterfaceAST(serviceDefinition spec.ServiceDefinition, info types.PkgInfo, config serviceASTConfig, generatorType generatorType) (astgen.ASTDecl, StringSet, error) {
	allImports := make(StringSet)
	var interfaceFuncs []*expression.InterfaceFunctionDecl
	serviceName := serviceDefinition.ServiceName.Name
	for _, endpointDefinition := range serviceDefinition.Endpoints {
		endpointName := string(endpointDefinition.EndpointName)
		params, imports, err := paramsForEndpoint(endpointDefinition, info, config, generatorType)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to generate parameters for endpoint %q", endpointName)
		}
		allImports.AddAll(imports)

		returnTypes, imports, err := returnTypesForEndpoint(endpointDefinition, info)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to generate return types for endpoint %q", endpointName)
		}
		allImports.AddAll(imports)

		interfaceFuncs = append(interfaceFuncs, &expression.InterfaceFunctionDecl{
			Name:        transforms.Export(endpointName),
			Params:      params,
			ReturnTypes: returnTypes,
			Comment:     transforms.Documentation(endpointDefinition.Docs),
		})
	}

	name := interfaceTypeName(serviceName)
	if generatorType == generatorTypeClient {
		name = clientInterfaceTypeName(name)
	}
	if config.withAuth {
		name = withAuthName(name)
	}
	return &decl.Interface{
		Name: name,
		InterfaceType: expression.InterfaceType{
			Functions: interfaceFuncs,
		},
		Comment: transforms.Documentation(serviceDefinition.Docs),
	}, allImports, nil
}

func withAuthServiceStructAST(serviceName string, hasHeaderAuth, hasCookieAuth bool, info types.PkgInfo) (astgen.ASTDecl, StringSet) {
	imports := NewStringSet()
	fields := []*expression.StructField{
		{
			Name: wrappedClientVar,
			Type: expression.Type(clientInterfaceTypeName(serviceName)),
		},
	}
	if hasHeaderAuth {
		fields = append(fields, &expression.StructField{
			Name: authHeaderVar,
			Type: expression.Type(types.Bearertoken.GoType(info)),
		})
		imports.Add(types.Bearertoken.ImportPaths()...)
	}
	if hasCookieAuth {
		fields = append(fields, &expression.StructField{
			Name: cookieTokenVar,
			Type: expression.Type(types.Bearertoken.GoType(info)),
		})
		imports.Add(types.Bearertoken.ImportPaths()...)
	}
	return decl.NewStruct(withAuthName(clientStructTypeName(serviceName)), fields, ""), imports
}

func withTokenServiceStructAST(serviceName string, info types.PkgInfo) (astgen.ASTDecl, StringSet) {
	imports := NewStringSet(types.Bearertoken.ImportPaths()...)
	imports.Add(tokenProviderImportPath)
	fields := []*expression.StructField{
		{
			Name: wrappedClientVar,
			Type: expression.Type(clientInterfaceTypeName(serviceName)),
		},
		{
			Name: tokenProviderVar,
			Type: tokenProviderType,
		},
	}
	return decl.NewStruct(withTokenProviderName(clientStructTypeName(serviceName)), fields, ""), imports
}

func serviceNewFuncAST(serviceName string) (astgen.ASTDecl, StringSet) {
	return &decl.Function{
		Name: "New" + clientInterfaceTypeName(serviceName),
		FuncType: expression.FuncType{
			Params: []*expression.FuncParam{
				expression.NewFuncParam(clientStructFieldName, httpClientClientType),
			},
			ReturnTypes: []expression.Type{
				expression.Type(clientInterfaceTypeName(serviceName)),
			},
		},
		Body: []astgen.ASTStmt{
			statement.NewReturn(
				expression.NewUnary(token.AND, expression.NewCompositeLit(
					expression.Type(clientStructTypeName(serviceName)),
					expression.NewKeyValue(clientStructFieldName, expression.VariableVal(clientStructFieldName)),
				)),
			),
		},
	}, NewStringSet(httpClientImportPath)
}

func withAuthServiceNewFuncAST(serviceName string, hasHeaderAuth, hasCookieAuth bool, info types.PkgInfo) (astgen.ASTDecl, StringSet) {
	funcParams := []*expression.FuncParam{
		expression.NewFuncParam(wrappedClientVar, expression.Type(clientInterfaceTypeName(serviceName))),
	}
	imports := NewStringSet()
	if hasHeaderAuth {
		funcParams = append(
			funcParams,
			expression.NewFuncParam(authHeaderVar, expression.Type(types.Bearertoken.GoType(info))),
		)
		imports.AddAll(NewStringSet(types.Bearertoken.ImportPaths()...))
	}
	if hasCookieAuth {
		funcParams = append(
			funcParams,
			expression.NewFuncParam(cookieTokenVar, expression.Type(types.Bearertoken.GoType(info))),
		)
		imports.AddAll(NewStringSet(types.Bearertoken.ImportPaths()...))
	}

	structElems := []astgen.ASTExpr{
		expression.NewKeyValue(wrappedClientVar, expression.VariableVal(wrappedClientVar)),
	}
	if hasHeaderAuth {
		structElems = append(
			structElems,
			expression.NewKeyValue(authHeaderVar, expression.VariableVal(authHeaderVar)),
		)
	}
	if hasCookieAuth {
		structElems = append(
			structElems,
			expression.NewKeyValue(cookieTokenVar, expression.VariableVal(cookieTokenVar)),
		)
	}

	return &decl.Function{
		Name: withAuthName("New" + clientInterfaceTypeName(serviceName)),
		FuncType: expression.FuncType{
			Params: funcParams,
			ReturnTypes: []expression.Type{
				expression.Type(withAuthName(clientInterfaceTypeName(serviceName))),
			},
		},
		Body: []astgen.ASTStmt{
			&statement.Return{
				Values: []astgen.ASTExpr{
					&expression.Unary{
						Op: token.AND,
						Receiver: &expression.CompositeLit{
							Type:     expression.Type(withAuthName(clientStructTypeName(serviceName))),
							Elements: structElems,
						},
					},
				},
			},
		},
	}, imports
}

func withTokenServiceNewFuncAST(serviceName string, info types.PkgInfo) (astgen.ASTDecl, StringSet) {
	funcParams := []*expression.FuncParam{
		expression.NewFuncParam(wrappedClientVar, expression.Type(clientInterfaceTypeName(serviceName))),
		expression.NewFuncParam(tokenProviderVar, tokenProviderType),
	}
	imports := NewStringSet()

	structElems := []astgen.ASTExpr{
		expression.NewKeyValue(wrappedClientVar, expression.VariableVal(wrappedClientVar)),
		expression.NewKeyValue(tokenProviderVar, expression.VariableVal(tokenProviderVar)),
	}

	return &decl.Function{
		Name: withTokenProviderName("New" + clientInterfaceTypeName(serviceName)),
		FuncType: expression.FuncType{
			Params: funcParams,
			ReturnTypes: []expression.Type{
				expression.Type(withAuthName(clientInterfaceTypeName(serviceName))),
			},
		},
		Body: []astgen.ASTStmt{
			&statement.Return{
				Values: []astgen.ASTExpr{
					&expression.Unary{
						Op: token.AND,
						Receiver: &expression.CompositeLit{
							Type:     expression.Type(withTokenProviderName(clientStructTypeName(serviceName))),
							Elements: structElems,
						},
					},
				},
			},
		},
	}, imports
}

func serviceStructMethodsAST(serviceDefinition spec.ServiceDefinition, info types.PkgInfo) ([]astgen.ASTDecl, StringSet, error) {
	allImports := make(StringSet)
	var methods []astgen.ASTDecl
	serviceName := serviceDefinition.ServiceName.Name
	for _, endpointDefinition := range serviceDefinition.Endpoints {
		endpointName := string(endpointDefinition.EndpointName)
		params, imports, err := paramsForEndpoint(endpointDefinition, info, serviceASTConfig{}, generatorTypeClient)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to generate parameters for endpoint %q", endpointName)
		}
		allImports.AddAll(imports)
		returnTypes, imports, err := returnTypesForEndpoint(endpointDefinition, info)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to generate return types for endpoint %q", endpointName)
		}
		allImports.AddAll(imports)
		returnBinary, err := isReturnTypeSpecificType(endpointDefinition.Returns, visitors.IsBinary)
		if err != nil {
			return nil, nil, err
		}
		body, err := serviceStructMethodBodyAST(endpointDefinition, returnTypes, returnBinary, info)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to generate method body for endpoint %q", endpointName)
		}

		methods = append(methods, &decl.Method{
			ReceiverName: receiverName,
			ReceiverType: expression.Type(clientStructTypeName(serviceName)).Pointer(),
			Function: decl.Function{
				Name: transforms.Export(endpointName),
				FuncType: expression.FuncType{
					Params:      params,
					ReturnTypes: returnTypes,
				},
				Body: body,
			},
		})
	}
	return methods, allImports, nil
}

func withAuthServiceStructMethodsAST(serviceDefinition spec.ServiceDefinition, info types.PkgInfo) ([]astgen.ASTDecl, StringSet, error) {
	allImports := make(StringSet)
	var methods []astgen.ASTDecl
	serviceName := serviceDefinition.ServiceName.Name
	for _, endpointDefinition := range serviceDefinition.Endpoints {
		endpointName := string(endpointDefinition.EndpointName)
		params, imports, err := paramsForEndpoint(endpointDefinition, info, serviceASTConfig{withAuth: true}, generatorTypeClient)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to generate parameters for endpoint %q", endpointName)
		}
		allImports.AddAll(imports)

		returnTypes, imports, err := returnTypesForEndpoint(endpointDefinition, info)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to generate return types for endpoint %q", endpointName)
		}
		allImports.AddAll(imports)
		body, err := serviceWithAuthStructMethodBodyAST(endpointDefinition, params)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to generate return auth structs for endpoint %q", endpointName)
		}

		methods = append(methods, &decl.Method{
			ReceiverName: receiverName,
			ReceiverType: expression.Type(withAuthName(clientStructTypeName(serviceName))).Pointer(),
			Function: decl.Function{
				Name: transforms.Export(endpointName),
				FuncType: expression.FuncType{
					Params:      params,
					ReturnTypes: returnTypes,
				},
				Body: body,
			},
		})
	}
	return methods, allImports, nil
}

func withTokenServiceStructMethodsAST(serviceDefinition spec.ServiceDefinition, info types.PkgInfo) ([]astgen.ASTDecl, StringSet, error) {
	allImports := make(StringSet)
	var methods []astgen.ASTDecl
	serviceName := serviceDefinition.ServiceName.Name
	for _, endpointDefinition := range serviceDefinition.Endpoints {
		endpointName := string(endpointDefinition.EndpointName)
		params, imports, err := paramsForEndpoint(endpointDefinition, info, serviceASTConfig{withTokenProvider: true}, generatorTypeClient)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Failed to generate parameters for endpoint %q", endpointName)
		}
		allImports.AddAll(imports)

		returnTypes, imports, err := returnTypesForEndpoint(endpointDefinition, info)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Failed to generate return types for endpoint %q", endpointName)
		}
		allImports.AddAll(imports)

		body, err := serviceWithTokenStructMethodBodyAST(endpointDefinition, params, returnTypes)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Failed to generate token provider structs for endpoint %q", endpointName)
		}

		methods = append(methods, &decl.Method{
			ReceiverName: receiverName,
			ReceiverType: expression.Type(withTokenProviderName(clientStructTypeName(serviceName))).Pointer(),
			Function: decl.Function{
				Name: transforms.Export(endpointName),
				FuncType: expression.FuncType{
					Params:      params,
					ReturnTypes: returnTypes,
				},
				Body: body,
			},
		})
	}
	return methods, allImports, nil
}

func isReturnTypeCollectionType(inType *spec.Type) (bool, error) {
	if inType == nil {
		return false, nil
	}
	isListType, err := visitors.IsSpecificConjureType(*inType, visitors.IsList)
	if err != nil {
		return false, err
	}
	if isListType {
		return true, nil
	}
	isMapType, err := visitors.IsSpecificConjureType(*inType, visitors.IsMap)
	if err != nil {
		return false, err
	}
	if isMapType {
		return true, nil
	}
	return visitors.IsSpecificConjureType(*inType, visitors.IsSet)
}

func isReturnTypeSpecificType(returnType *spec.Type, typeCheck visitors.TypeCheck) (bool, error) {
	if returnType == nil {
		return false, nil
	}
	isType, err := visitors.IsSpecificConjureType(*returnType, typeCheck)
	if err != nil {
		return false, err
	}
	return isType, nil
}

func isBinaryType(specType spec.Type) (bool, error) {
	return visitors.IsSpecificConjureType(specType, visitors.IsBinary)
}

var pathParamRegexp = regexp.MustCompile(regexp.QuoteMeta("{") + "[^}]+" + regexp.QuoteMeta("}"))

func serviceStructMethodBodyAST(endpointDefinition spec.EndpointDefinition, returnTypes expression.Types, returnsBinary bool, info types.PkgInfo) ([]astgen.ASTStmt, error) {
	var body []astgen.ASTStmt

	const (
		defaultReturnValVar = "defaultReturnVal"
		returnValVar        = "returnVal"
		respVar             = "resp"

		requestParamsVar = "requestParams"
		queryParamsVar   = "queryParams"
		errVar           = "err"
	)

	hasReturnVal := len(returnTypes) > 1
	returnsCollection := false
	returnsOptional := false

	// if endpoint returns a value, declare variables for value
	if hasReturnVal && !returnsBinary {
		isCollection, err := isReturnTypeCollectionType(endpointDefinition.Returns)
		if err != nil {
			return nil, err
		}
		returnsCollection = isCollection

		isOptional, err := isReturnTypeSpecificType(endpointDefinition.Returns, visitors.IsOptional)
		if err != nil {
			return nil, err
		}
		returnsOptional = isOptional

		if !returnsCollection && !returnsOptional {
			// return value cannot be nil: create an indirected version of the variable to unmarshal into to verify it is non-nil
			body = append(body, statement.NewDecl(decl.NewVar(defaultReturnValVar, returnTypes[0])))
			body = append(body, statement.NewDecl(decl.NewVar(returnValVar, returnTypes[0].Pointer())))
		} else {
			body = append(body, statement.NewDecl(decl.NewVar(returnValVar, returnTypes[0])))
		}
	}
	body = append(body, statement.NewDecl(decl.NewVar(requestParamsVar, expression.Type(fmt.Sprintf("[]%s.RequestParam", httpClientPkgName)))))

	// function that creates the statement "requestParams = append(requestParams, httpclient.{httpClientFuncName}({args}))"
	appendToRequestParamsFn := func(httpClientFuncName string, args ...astgen.ASTExpr) astgen.ASTStmt {
		return statement.NewAssignment(
			expression.VariableVal(requestParamsVar),
			token.ASSIGN,
			expression.NewCallExpression(expression.AppendBuiltIn,
				expression.VariableVal(requestParamsVar),
				expression.NewCallFunction(httpClientPkgName, httpClientFuncName, args...),
			),
		)
	}

	// function that creates the statement "requestParams = append(requestParams, httpclient.{httpClientFuncName}({args}))" and appends it to the "body" variable
	bodyAppendToRequestParamsFn := func(httpClientFuncName string, args ...astgen.ASTExpr) {
		body = append(body, appendToRequestParamsFn(httpClientFuncName, args...))
	}

	bodyAppendToRequestParamsFn("WithRPCMethodName", expression.StringVal(transforms.Export(string(endpointDefinition.EndpointName))))
	bodyAppendToRequestParamsFn("WithRequestMethod", expression.StringVal(endpointDefinition.HttpMethod.Value()))

	// auth
	if endpointDefinition.Auth != nil {
		info.AddImports("fmt")
		if authHeader, err := visitors.GetPossibleHeaderAuth(*endpointDefinition.Auth); err != nil {
			return nil, err
		} else if authHeader != nil {
			bodyAppendToRequestParamsFn("WithHeader",
				expression.StringVal("Authorization"),
				expression.NewCallFunction("fmt", "Sprint", expression.StringVal("Bearer "), expression.VariableVal(authHeaderVar)),
			)
		}
		if authCookie, err := visitors.GetPossibleCookieAuth(*endpointDefinition.Auth); err != nil {
			return nil, err
		} else if authCookie != nil {
			bodyAppendToRequestParamsFn("WithHeader",
				expression.StringVal("Cookie"),
				expression.NewCallFunction("fmt", "Sprint", expression.StringVal(authCookie.CookieName+"="), expression.VariableVal(cookieTokenVar)),
			)
		}
	}

	pathParamArgs := []astgen.ASTExpr{
		expression.StringVal(pathParamRegexp.ReplaceAllString(string(endpointDefinition.HttpPath), regexp.QuoteMeta(`%s`))),
	}
	pathParams, err := visitors.GetPathParams(endpointDefinition.Args)
	if err != nil {
		return nil, err
	}
	for _, pathParam := range pathParams {
		info.AddImports("fmt", "net/url")
		pathParamArgs = append(pathParamArgs,
			expression.NewCallFunction("url", "PathEscape",
				expression.NewCallFunction("fmt", "Sprint", expression.VariableVal(argNameTransform(string(pathParam.ArgumentDefinition.ArgName))))),
		)
	}
	// path params
	bodyAppendToRequestParamsFn("WithPathf", pathParamArgs...)

	// body params
	bodyParams, err := visitors.GetBodyParams(endpointDefinition.Args)
	if err != nil {
		return nil, err
	}
	if len(bodyParams) > 0 {
		if len(bodyParams) != 1 {
			return nil, errors.Errorf("more than 1 body param exists: %v", bodyParams)
		}
		requestFn := "WithJSONRequest"
		bodyArgDef := bodyParams[0].ArgumentDefinition
		if isBinaryParam, err := isBinaryType(bodyArgDef.Type); err != nil {
			return nil, err
		} else if isBinaryParam {
			requestFn = "WithRawRequestBodyProvider"
		}
		bodyAppendToRequestParamsFn(requestFn, expression.VariableVal(argNameTransform(string(bodyArgDef.ArgName))))
	}

	// header params
	headerParams, err := visitors.GetHeaderParams(endpointDefinition.Args)
	if err != nil {
		return nil, err
	}
	for _, headerParam := range headerParams {
		argName := argNameTransform(string(headerParam.ArgumentDefinition.ArgName))

		isOptional, err := isReturnTypeSpecificType(&headerParam.ArgumentDefinition.Type, visitors.IsOptional)
		if err != nil {
			return nil, err
		}

		var variableVar astgen.ASTExpr = expression.VariableVal(argName)
		if isOptional {
			variableVar = expression.NewStar(variableVar)
		}
		appendExpr := appendToRequestParamsFn("WithHeader", expression.StringVal(visitors.GetParamID(headerParam.ArgumentDefinition)), expression.NewCallFunction("fmt", "Sprint", variableVar))

		// if header parameter type is an optional, append dereferenced value if it is non-nil
		if isOptional {
			appendExpr = &statement.If{
				Cond: &expression.Binary{
					LHS: expression.VariableVal(argName),
					Op:  token.NEQ,
					RHS: expression.Nil,
				},
				Body: []astgen.ASTStmt{
					appendExpr,
				},
			}
		}

		body = append(body, appendExpr)
		info.AddImports("fmt")
	}

	// query params
	queryParams, err := visitors.GetQueryParams(endpointDefinition.Args)
	if err != nil {
		return nil, err
	}
	if len(queryParams) > 0 {
		body = append(body, &statement.Assignment{
			LHS: []astgen.ASTExpr{
				expression.VariableVal(queryParamsVar),
			},
			Tok: token.DEFINE,
			RHS: expression.NewCallExpression(expression.MakeBuiltIn, expression.Type("url.Values")),
		})
		info.AddImports("net/url")

		for _, queryParam := range queryParams {
			currQueryParamVarName := argNameTransform(string(queryParam.ArgumentDefinition.ArgName))
			currQueryParamKeyName := visitors.GetParamID(queryParam.ArgumentDefinition)
			currQueryTypeProvider, err := visitors.NewConjureTypeProvider(queryParam.ArgumentDefinition.Type)
			if err != nil {
				return nil, err
			}

			isOptional := currQueryTypeProvider.IsSpecificType(visitors.IsOptional)
			isList := currQueryTypeProvider.IsSpecificType(visitors.IsList) || currQueryTypeProvider.IsSpecificType(visitors.IsSet)

			var accessVarContentExpr astgen.ASTExpr = expression.VariableVal(currQueryParamVarName)
			if isOptional {
				accessVarContentExpr = expression.NewUnary(token.MUL, accessVarContentExpr)
			}

			var addQueryParamStmt astgen.ASTStmt
			if isList {
				addQueryParamStmt = &statement.Range{
					Key:   expression.VariableVal("_"),
					Value: expression.VariableVal("v"),
					Tok:   token.DEFINE,
					Expr:  accessVarContentExpr,
					Body: []astgen.ASTStmt{
						statement.NewExpression(expression.NewCallFunction(
							queryParamsVar,
							"Add",
							expression.StringVal(currQueryParamKeyName),
							expression.NewCallFunction("fmt", "Sprint", expression.VariableVal("v"))),
						),
					},
				}
			} else {
				addQueryParamStmt = statement.NewExpression(expression.NewCallFunction(
					queryParamsVar,
					"Set",
					expression.StringVal(currQueryParamKeyName),
					expression.NewCallFunction("fmt", "Sprint", accessVarContentExpr)),
				)
			}
			info.AddImports("fmt")

			if isOptional {
				addQueryParamStmt = &statement.If{
					Cond: &expression.Binary{
						LHS: expression.VariableVal(currQueryParamVarName),
						Op:  token.NEQ,
						RHS: expression.Nil,
					},
					Body: []astgen.ASTStmt{
						addQueryParamStmt,
					},
				}
			}
			body = append(body, addQueryParamStmt)
		}
		bodyAppendToRequestParamsFn("WithQueryValues", expression.VariableVal(queryParamsVar))
	}

	// return val
	switch {
	case returnsBinary:
		bodyAppendToRequestParamsFn("WithRawResponseBody")
	case hasReturnVal:
		bodyAppendToRequestParamsFn("WithJSONResponse", expression.NewUnary(token.AND, expression.VariableVal(returnValVar)))
	}

	body = append(body, &statement.Assignment{
		LHS: []astgen.ASTExpr{
			expression.VariableVal(respVar),
			expression.VariableVal(errVar),
		},
		Tok: token.DEFINE,
		RHS: expression.NewCallFunction(fmt.Sprintf("%s.%s", receiverName, clientStructFieldName), "Do",
			expression.VariableVal(ctxName),
			expression.VariableVal(requestParamsVar+"...")),
	})

	valVarToReturnInErr := returnValVar
	if returnsBinary {
		valVarToReturnInErr = "nil"
	} else if !returnsCollection && !returnsOptional {
		valVarToReturnInErr = defaultReturnValVar
	}
	body = append(body, ifErrNotNilReturnHelper(hasReturnVal, valVarToReturnInErr, errVar, nil))

	if returnsBinary {
		// if endpoint returns binary, return body of response directly
		body = append(body, statement.NewReturn(
			expression.NewSelector(
				expression.VariableVal(respVar),
				"Body",
			),
			expression.Nil,
		))
		return body, nil
	}

	// otherwise, return values
	body = append(body, &statement.Assignment{
		LHS: []astgen.ASTExpr{
			expression.Blank,
		},
		Tok: token.ASSIGN,
		RHS: expression.VariableVal(respVar),
	})

	var returnExp astgen.ASTExpr = expression.VariableVal(returnValVar)
	if hasReturnVal {
		switch {
		case !returnsCollection && !returnsOptional:
			// verify that return value is non-nil and dereference
			body = append(body, &statement.If{
				Cond: &expression.Binary{
					LHS: expression.VariableVal(returnValVar),
					Op:  token.EQL,
					RHS: expression.Nil,
				},
				Body: []astgen.ASTStmt{
					statement.NewReturn(
						expression.VariableVal(defaultReturnValVar),
						expression.NewCallFunction("fmt", "Errorf", expression.StringVal(returnValVar+" cannot be nil")),
					),
				},
			})
			info.AddImports("fmt")
			returnExp = expression.NewUnary(token.MUL, returnExp)
		case returnsCollection:
			// if returned value is nil, initialize to empty instead
			body = append(body, &statement.If{
				Cond: &expression.Binary{
					LHS: expression.VariableVal(returnValVar),
					Op:  token.EQL,
					RHS: expression.Nil,
				},
				Body: []astgen.ASTStmt{
					statement.NewAssignment(
						expression.VariableVal(returnValVar),
						token.ASSIGN,
						expression.NewCallExpression(expression.MakeBuiltIn, expression.VariableVal(returnTypes[0]), expression.IntVal(0)),
					),
				},
			})
		}
	}

	body = append(body, &statement.Return{
		Values: returnVals(hasReturnVal,
			returnExp,
			expression.Nil),
	})
	return body, nil
}

func serviceWithAuthStructMethodBodyAST(endpointDefinition spec.EndpointDefinition, params expression.FuncParams) ([]astgen.ASTStmt, error) {
	endpointName := string(endpointDefinition.EndpointName)
	args := []astgen.ASTExpr{expression.VariableVal(ctxName)}
	if endpointDefinition.Auth != nil {
		possibleHeader, err := visitors.GetPossibleHeaderAuth(*endpointDefinition.Auth)
		if err != nil {
			return nil, err
		}
		if possibleHeader != nil {
			args = append(args, expression.NewSelector(expression.VariableVal(receiverName), authHeaderVar))
		}
		possibleCookie, err := visitors.GetPossibleCookieAuth(*endpointDefinition.Auth)
		if err != nil {
			return nil, err
		}
		if possibleCookie != nil {
			args = append(args, expression.NewSelector(expression.VariableVal(receiverName), cookieTokenVar))
		}
	}
	for _, param := range params {
		if param.Type == "context.Context" {
			// We already added ctx as the first argument.
			continue
		}
		for _, curr := range param.Names {
			args = append(args, expression.VariableVal(curr))
		}
	}

	// Invoke wrapped client with updated arguments
	// return c.client.DoThing(ctx, authHeader, arg1)
	return []astgen.ASTStmt{
		statement.NewReturn(
			expression.NewCallExpression(
				expression.NewSelector(
					expression.NewSelector(
						expression.VariableVal(receiverName),
						wrappedClientVar,
					),
					transforms.Export(endpointName),
				),
				args...,
			),
		),
	}, nil

}

func serviceWithTokenStructMethodBodyAST(endpointDefinition spec.EndpointDefinition, params expression.FuncParams, returnTypes expression.Types) ([]astgen.ASTStmt, error) {
	endpointName := string(endpointDefinition.EndpointName)
	args := []astgen.ASTExpr{expression.VariableVal(ctxName)}
	statements := []astgen.ASTStmt{}
	if endpointDefinition.Auth != nil {
		possibleHeader, err := visitors.GetPossibleHeaderAuth(*endpointDefinition.Auth)
		if err != nil {
			return nil, err
		}
		possibleCookie, err := visitors.GetPossibleCookieAuth(*endpointDefinition.Auth)
		if err != nil {
			return nil, err
		}
		if possibleHeader != nil || possibleCookie != nil {
			tokenInit := &statement.Assignment{
				LHS: []astgen.ASTExpr{
					expression.VariableVal("token"),
					expression.VariableVal("err"),
				},
				Tok: token.DEFINE,
				RHS: expression.NewCallExpression(expression.NewSelector(expression.VariableVal(receiverName), tokenProviderVar), expression.VariableVal(ctxName)),
			}
			var errReturn *statement.If
			if len(returnTypes) > 1 {
				statements = append(statements, statement.NewDecl(decl.NewVar("defaultReturnVal", returnTypes[0])))
				errReturn = ifErrNotNilReturnHelper(true, "defaultReturnVal", "err", nil)
			} else {
				errReturn = ifErrNotNilReturnErrStatement("err", nil)
			}
			args = append(args, expression.NewCallExpression(expression.Type("bearertoken.Token"), expression.VariableVal("token")))
			statements = append(statements, tokenInit, errReturn)
		}
	}

	for _, param := range params {
		if param.Type == "context.Context" {
			// We already added ctx as the first argument.
			continue
		}
		for _, curr := range param.Names {
			args = append(args, expression.VariableVal(curr))
		}
	}

	return append(statements,
		statement.NewReturn(
			expression.NewCallExpression(
				expression.NewSelector(
					expression.NewSelector(expression.VariableVal(receiverName),
						wrappedClientVar,
					),
					transforms.Export(endpointName),
				),
				args...,
			),
		),
	), nil
}

func returnVals(hasReturnVal bool, optional, required astgen.ASTExpr) []astgen.ASTExpr {
	var rvals []astgen.ASTExpr
	if hasReturnVal {
		rvals = append(rvals, optional)
	}
	return append(rvals, required)
}

func ifErrNotNilReturnErrStatement(errVarName string, init astgen.ASTStmt) *statement.If {
	return ifErrNotNilReturnHelper(false, "", errVarName, init)
}

func ifErrNotNilReturnHelper(hasReturnVal bool, valVarName, errVarName string, init astgen.ASTStmt) *statement.If {
	return &statement.If{
		Init: init,
		Cond: &expression.Binary{
			LHS: expression.VariableVal(errVarName),
			Op:  token.NEQ,
			RHS: expression.Nil,
		},
		Body: []astgen.ASTStmt{
			&statement.Return{
				Values: returnVals(hasReturnVal,
					expression.VariableVal(valVarName),
					expression.VariableVal(errVarName)),
			},
		},
	}
}

func returnTypesForEndpoint(endpointDefinition spec.EndpointDefinition, info types.PkgInfo) (expression.Types, StringSet, error) {
	var returnTypes []expression.Type
	imports := make(StringSet)
	if endpointDefinition.Returns != nil {
		var goType string
		returnBinary, err := isReturnTypeSpecificType(endpointDefinition.Returns, visitors.IsBinary)
		if err != nil {
			return nil, nil, err
		}
		if returnBinary {
			// special case: "binary" type resolves to []byte in structs, but indicates a streaming response when
			// specified as the return type of a service, so use "io.ReadCloser".
			goType = types.IOReadCloserType.GoType(info)
			imports.AddAll(NewStringSet(types.IOReadCloserType.ImportPaths()...))
		} else {
			typer, err := visitors.NewConjureTypeProviderTyper(*endpointDefinition.Returns, info)
			if err != nil {
				return nil, nil, err
			}
			goType = typer.GoType(info)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "failed to process return type %q", goType)
			}
		}
		returnTypes = append(returnTypes, expression.Type(goType))
	}
	return append(returnTypes, expression.ErrorType), imports, nil
}

func paramsForEndpoint(endpointDefinition spec.EndpointDefinition, info types.PkgInfo, config serviceASTConfig, generatorType generatorType) (expression.FuncParams, StringSet, error) {
	imports := NewStringSet("context")
	params := []*expression.FuncParam{expression.NewFuncParam(ctxName, expression.Type("context.Context"))}
	if endpointDefinition.Auth != nil && !config.withAuth && !config.withTokenProvider {
		if authHeader, err := visitors.GetPossibleHeaderAuth(*endpointDefinition.Auth); err != nil {
			return nil, nil, err
		} else if authHeader != nil {
			params = append(params, expression.NewFuncParam(authHeaderVar, expression.Type(types.Bearertoken.GoType(info))))
		}
		if authCookie, err := visitors.GetPossibleCookieAuth(*endpointDefinition.Auth); err != nil {
			return nil, nil, err
		} else if authCookie != nil {
			params = append(params, expression.NewFuncParam(cookieTokenVar, expression.Type(types.Bearertoken.GoType(info))))
		}
		imports.AddAll(NewStringSet(types.Bearertoken.ImportPaths()...))
	}
	for _, arg := range endpointDefinition.Args {
		binaryParam, err := isBinaryType(arg.Type)
		if err != nil {
			return nil, nil, err
		}

		var goType string
		argName := string(arg.ArgName)
		if binaryParam {
			// special case: "binary" types resolve to []byte, but this indicates a streaming parameter when
			// specified as the request argument of a service, so use "io.ReadCloser".
			if generatorType == generatorTypeClient {
				// special case: the client provides "func() io.ReadCloser" instead of "io.ReadCloser" so
				// that a fresh "io.ReadCloser" can be retrieved for retries.
				goType = types.GetBodyType.GoType(info)
				imports.AddAll(NewStringSet(types.GetBodyType.ImportPaths()...))
			} else {
				goType = types.IOReadCloserType.GoType(info)
				imports.AddAll(NewStringSet(types.IOReadCloserType.ImportPaths()...))
			}
		} else {
			typer, err := visitors.NewConjureTypeProviderTyper(arg.Type, info)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "failed to process param %q", argName)
			}
			goType = typer.GoType(info)
		}
		params = append(params, expression.NewFuncParam(argNameTransform(argName), expression.Type(goType)))
		imports.AddAll(NewStringSet(types.Bearertoken.ImportPaths()...))
	}
	return params, imports, nil
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
