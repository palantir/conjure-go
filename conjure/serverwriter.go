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
	"strings"

	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	"github.com/palantir/goastwriter/statement"
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/pkg/errors"

	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/conjure/transforms"
	"github.com/palantir/conjure-go/conjure/types"
	"github.com/palantir/conjure-go/conjure/visitors"
	"github.com/palantir/conjure-go/conjure/werrorexpressions"
)

const (
	registerPrefix = "RegisterRoutes"
	errorName      = "err"
	okName         = "ok"
	implName       = "impl"

	// Handler
	handlerName = "handler"

	// Router
	routerVarName           = "router"
	routerImportPackage     = "wrouter"
	routerImportClass       = "Router"
	routerPathParamsMapFunc = "PathParams"
	resourceName            = "resource"

	// Server
	serverResourceImportPackage = "wresource"
	serverResourceFunctionName  = "New"
	restImportPackage           = "rest"

	// Handler
	handlerStructNameSuffix   = "Handler"
	handlerFunctionNamePrefix = "Handle"

	// Auth
	funcParseBearerTokenHeader = "ParseBearerTokenHeader"
	authCookieVar              = "authCookie"

	// ResponseWriter
	responseWriterVarName = "rw"
	responseArgVarName    = "respArg"
	httpPackageName       = "http"
	responseWriterType    = "ResponseWriter"

	// Request
	requestVarName    = "req"
	requestVarType    = "*" + httpPackageName + ".Request"
	requestHeaderFunc = "Header"
	requestURLField   = "URL"
	urlQueryFunc      = "Query"

	// Codecs
	codecsJSON           = "codecs.JSON"
	codecEncodeFunc      = "Encode"
	codecDecodeFunc      = "Decode"
	codecContentTypeFunc = "ContentType"
)

func ASTForServerRouteRegistration(serviceDefinition spec.ServiceDefinition, info types.PkgInfo) ([]astgen.ASTDecl, error) {
	info.AddImports(
		"github.com/palantir/conjure-go-runtime/conjure-go-contract/codecs",
		"github.com/palantir/witchcraft-go-server/rest",
		"github.com/palantir/witchcraft-go-server/witchcraft",
		"github.com/palantir/witchcraft-go-server/witchcraft/wresource",
		"github.com/palantir/witchcraft-go-server/wrouter")
	info.SetImports("werror", "github.com/palantir/witchcraft-go-error")
	serviceName := serviceDefinition.ServiceName.Name
	funcName := registerPrefix + strings.Title(serviceName)
	serviceImplName := transforms.Export(serviceName)
	body, err := getRegisterRoutesBody(serviceDefinition)
	if err != nil {
		return nil, err
	}
	registerRoutesFunc := &decl.Function{
		Comment: funcName + " registers handlers for the " + serviceName + " endpoints with a witchcraft wrouter.\n" +
			"This should typically be called in a witchcraft server's InitFunc.\n" +
			"impl provides an implementation of each endpoint, which can assume the request parameters have been parsed\n" +
			"in accordance with the Conjure specification.",
		Name: funcName,
		FuncType: expression.FuncType{
			Params: []*expression.FuncParam{
				{
					Names: []string{routerVarName},
					Type:  expression.Type(fmt.Sprintf("%s.%s", routerImportPackage, routerImportClass)),
				},
				{
					Names: []string{implName},
					Type:  expression.Type(serviceImplName),
				},
			},
			ReturnTypes: []expression.Type{
				expression.ErrorType,
			},
		},
		Body: body,
	}
	components := []astgen.ASTDecl{
		registerRoutesFunc,
	}
	return components, nil
}

func getRegisterRoutesBody(serviceDefinition spec.ServiceDefinition) ([]astgen.ASTStmt, error) {
	var body []astgen.ASTStmt
	// Create the handler struct
	body = append(body, &statement.Assignment{
		LHS: []astgen.ASTExpr{
			expression.VariableVal(handlerName),
		},
		Tok: token.DEFINE,
		RHS: createHandlerSpec(serviceDefinition),
	})
	// Create the witchcraft resource
	body = append(body, &statement.Assignment{
		LHS: []astgen.ASTExpr{
			expression.VariableVal(resourceName),
		},
		Tok: token.DEFINE,
		RHS: expression.NewCallFunction(serverResourceImportPackage, serverResourceFunctionName, []astgen.ASTExpr{
			expression.StringVal(strings.ToLower(serviceDefinition.ServiceName.Name)),
			expression.VariableVal(routerVarName),
		}...),
	})
	// For each endpoint, register a route on the provided router
	// if err := resource.Get(...); err != nil {
	//     return werror.Wrap(err, ...)
	// }
	// TODO(bmoylan): Register safe params - Nothing in the conjure def tells us what is safe or unsafe. How do we know? Kevin says this is markers
	for _, endpoint := range serviceDefinition.Endpoints {
		endpointTitleName := strings.Title(string(endpoint.EndpointName))
		stmt := statement.If{
			Init: statement.NewAssignment(
				expression.VariableVal(errorName),
				token.DEFINE,
				expression.NewCallFunction(resourceName, getResourceFunction(endpoint), []astgen.ASTExpr{
					expression.StringVal(endpointTitleName),
					expression.StringVal(getPathToRegister(endpoint)),
					astForRestJSONHandler(expression.NewSelector(expression.VariableVal(handlerName), "Handle"+endpointTitleName)),
				}...),
			),
			Cond: &expression.Binary{
				LHS: expression.VariableVal(errorName),
				Op:  token.NEQ,
				RHS: expression.Nil,
			},
			Body: []astgen.ASTStmt{
				&statement.Return{
					Values: []astgen.ASTExpr{
						werrorexpressions.CreateWrapWErrorExpression(errorName, "failed to add route", map[string]string{"routeName": endpointTitleName}),
					},
				},
			},
		}
		body = append(body, &stmt)
	}
	// Return nil if everything registered
	body = append(body, &statement.Return{
		Values: []astgen.ASTExpr{expression.Nil},
	})
	return body, nil
}

func createHandlerSpec(serviceDefinition spec.ServiceDefinition) astgen.ASTExpr {
	return expression.NewCompositeLit(
		expression.Type(getHandlerStuctName(serviceDefinition)),
		expression.NewKeyValue(implName, expression.VariableVal(implName)),
	)
}

func getPathToRegister(endpointDefinition spec.EndpointDefinition) string {
	return string(endpointDefinition.HttpPath)
}

func getResourceFunction(endpointDefinition spec.EndpointDefinition) string {
	switch endpointDefinition.HttpMethod {
	case spec.HttpMethodGet:
		return "Get"
	case spec.HttpMethodPost:
		return "Post"
	case spec.HttpMethodPut:
		return "Put"
	case spec.HttpMethodDelete:
		return "Delete"
	default:
		return "Unknown"
	}
}

func AstForServerInterface(serviceDefinition spec.ServiceDefinition, info types.PkgInfo) ([]astgen.ASTDecl, error) {
	serviceName := serviceDefinition.ServiceName.Name
	isClient := false
	interfaceAST, _, err := serviceInterfaceAST(serviceDefinition, info, false, isClient)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate interface for service %q", serviceName)
	}
	components := []astgen.ASTDecl{
		interfaceAST,
	}
	return components, nil
}

func AstForServerFunctionHandler(serviceDefinition spec.ServiceDefinition, info types.PkgInfo) ([]astgen.ASTDecl, error) {
	var components []astgen.ASTDecl
	implStructs := getHandlerStruct(serviceDefinition)
	components = append(components, implStructs)
	methods, err := getHandleMethods(serviceDefinition, info)
	if err != nil {
		return nil, err
	}
	for _, method := range methods {
		components = append(components, method)
	}
	return components, nil
}

func getHandleMethods(serviceDefinition spec.ServiceDefinition, info types.PkgInfo) ([]*decl.Method, error) {
	var methods []*decl.Method
	for _, endpoint := range serviceDefinition.Endpoints {
		method, err := getHandleMethod(serviceDefinition, endpoint, info)
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}
	return methods, nil
}

func getHandleMethod(serviceDefinition spec.ServiceDefinition, endpoint spec.EndpointDefinition, info types.PkgInfo) (*decl.Method, error) {
	info.AddImports("net/http")
	body, err := getHandleMethodBody(serviceDefinition, endpoint, info)
	if err != nil {
		return nil, err
	}
	receiverName := getReceiverName(serviceDefinition)
	titleEndpoint := strings.Title(string(endpoint.EndpointName))
	methods := &decl.Method{
		ReceiverName: receiverName,
		ReceiverType: expression.Type(getHandlerStuctName(serviceDefinition)).Pointer(),
		Function: decl.Function{
			Name: handlerFunctionNamePrefix + titleEndpoint,
			FuncType: expression.FuncType{
				Params: []*expression.FuncParam{
					{
						Names: []string{responseWriterVarName},
						Type:  expression.Type(strings.Join([]string{httpPackageName, responseWriterType}, ".")),
					},
					{
						Names: []string{requestVarName},
						Type:  expression.Type(requestVarType),
					},
				},
				ReturnTypes: []expression.Type{expression.ErrorType},
			},
			Body: body,
		},
	}
	return methods, nil
}

func getHandleMethodBody(serviceDefinition spec.ServiceDefinition, endpoint spec.EndpointDefinition, info types.PkgInfo) ([]astgen.ASTStmt, error) {
	var body []astgen.ASTStmt

	pathParams, err := visitors.GetPathParams(endpoint.Args)
	if err != nil {
		return nil, err
	}
	headerParams, err := visitors.GetHeaderParams(endpoint.Args)
	if err != nil {
		return nil, err
	}
	queryParams, err := visitors.GetQueryParams(endpoint.Args)
	if err != nil {
		return nil, err
	}
	bodyParams, err := visitors.GetBodyParams(endpoint.Args)
	if err != nil {
		return nil, err
	}
	var bodyParam *visitors.ArgumentDefinitionBodyParam
	switch len(bodyParams) {
	case 0:
	case 1:
		bodyParam = &bodyParams[0]
	default:
		return nil, errors.New("only 1 body param is supported: Conjure IR generator should have caught this")
	}

	authStatements, err := getAuthStatements(endpoint.Auth, info)
	if err != nil {
		return nil, err
	}
	body = append(body, authStatements...)

	pathParamStatements, err := getPathParamStatements(pathParams, info)
	if err != nil {
		return nil, err
	}
	body = append(body, pathParamStatements...)

	queryParamStatements, err := getQueryParamStatements(queryParams, info)
	if err != nil {
		return nil, err
	}
	body = append(body, queryParamStatements...)

	headerParamStatements, err := getHeaderParamStatements(headerParams, info)
	if err != nil {
		return nil, err
	}
	body = append(body, headerParamStatements...)

	bodyParamStatements, err := getBodyParamStatements(bodyParam, info)
	if err != nil {
		return nil, err
	}
	body = append(body, bodyParamStatements...)

	varsToPassIntoImpl := []astgen.ASTExpr{expression.NewCallFunction(requestVarName, "Context")}

	if endpoint.Auth != nil {
		if headerAuth, err := visitors.GetPossibleHeaderAuth(*endpoint.Auth); err != nil {
			return nil, err
		} else if headerAuth != nil {
			varsToPassIntoImpl = append(varsToPassIntoImpl, expression.NewCallExpression(
				expression.Type(types.Bearertoken.GoType(info)),
				expression.VariableVal(authHeaderVar),
			))
		}
		if cookieAuth, err := visitors.GetPossibleCookieAuth(*endpoint.Auth); err != nil {
			return nil, err
		} else if cookieAuth != nil {
			varsToPassIntoImpl = append(varsToPassIntoImpl, expression.VariableVal(cookieTokenVar))
		}
	}

	for _, arg := range endpoint.Args {
		varsToPassIntoImpl = append(varsToPassIntoImpl, expression.VariableVal(transforms.SafeName(string(arg.ArgName))))
	}

	returnStatements, err := getReturnStatements(serviceDefinition, endpoint, varsToPassIntoImpl, info)
	if err != nil {
		return nil, err
	}
	body = append(body, returnStatements...)

	return body, nil
}

func getReturnStatements(
	serviceDefinition spec.ServiceDefinition,
	endpoint spec.EndpointDefinition,
	varsToPassIntoImpl []astgen.ASTExpr,
	info types.PkgInfo,
) ([]astgen.ASTStmt, error) {
	var body []astgen.ASTStmt
	receiverName := getReceiverName(serviceDefinition)
	endpointName := string(endpoint.EndpointName)
	endpointNameFirstLetterUpper := strings.Title(endpointName)
	// This is make the call to the interface
	makeFunctionCall := expression.NewCallFunction(receiverName+"."+implName, endpointNameFirstLetterUpper, varsToPassIntoImpl...)

	if endpoint.Returns == nil {
		// The endpoint doesn't return anything, just return the interface call
		body = append(body, &statement.Return{
			Values: []astgen.ASTExpr{makeFunctionCall},
		})
		return body, nil
	}

	// Make the call
	body = append(body, &statement.Assignment{
		LHS: []astgen.ASTExpr{
			expression.VariableVal(responseArgVarName),
			expression.VariableVal(errorName),
		},
		Tok: token.DEFINE,
		RHS: makeFunctionCall,
	})
	// Return an error if present
	body = append(body, getIfErrNotNilReturnErrExpression())

	var codec types.Typer
	if isBinary, err := isBinaryType(*endpoint.Returns); err != nil {
		return nil, err
	} else if isBinary {
		codec = types.CodecBinary
	} else {
		codec = types.CodecJSON
	}
	info.AddImports(codec.ImportPaths()...)

	body = append(body, statement.NewExpression(&expression.CallExpression{
		Function: &expression.Selector{
			Receiver: expression.NewCallFunction(responseWriterVarName, "Header"),
			Selector: "Add",
		},
		Args: []astgen.ASTExpr{
			expression.StringVal("Content-Type"),
			expression.NewCallFunction(codec.GoType(info), codecContentTypeFunc),
		},
	}))

	// Return error from writing object into response
	body = append(body, &statement.Return{
		Values: []astgen.ASTExpr{
			expression.NewCallFunction(codec.GoType(info), codecEncodeFunc,
				expression.VariableVal(responseWriterVarName),
				expression.VariableVal(responseArgVarName),
			),
		},
	})

	return body, nil
}

func getBodyParamStatements(bodyParam *visitors.ArgumentDefinitionBodyParam, info types.PkgInfo) ([]astgen.ASTStmt, error) {
	if bodyParam == nil {
		return nil, nil
	}
	var body []astgen.ASTStmt
	argName := transforms.SafeName(string(bodyParam.ArgumentDefinition.ArgName))
	typer, err := visitors.NewConjureTypeProviderTyper(bodyParam.ArgumentDefinition.Type, info)
	if err != nil {
		typJSON, _ := bodyParam.ArgumentDefinition.Type.MarshalJSON()
		return nil, errors.Wrapf(err, "failed to process return type %s", string(typJSON))
	}
	info.AddImports(typer.ImportPaths()...)

	if isBinary, err := isBinaryType(bodyParam.ArgumentDefinition.Type); err != nil {
		return nil, err
	} else if isBinary {
		// If the body argument is binary, pass req.Body directly to the impl.
		body = append(body, &statement.Assignment{
			LHS: []astgen.ASTExpr{expression.VariableVal(argName)},
			Tok: token.DEFINE,
			RHS: expression.NewSelector(expression.VariableVal(requestVarName), "Body"),
		})
	} else {
		// If the request is not binary, it is JSON. Unmarshal the req.Body.

		// Create the empty type of this object
		body = append(body, statement.NewDecl(decl.NewVar(argName, expression.Type(typer.GoType(info)))))
		// Decode request
		body = append(body, &statement.If{
			Init: &statement.Assignment{
				LHS: []astgen.ASTExpr{expression.VariableVal(errorName)},
				Tok: token.DEFINE,
				RHS: expression.NewCallFunction(
					codecsJSON,
					codecDecodeFunc,
					expression.NewSelector(expression.VariableVal(requestVarName), "Body"),
					expression.NewUnary(token.AND, expression.VariableVal(argName))),
			},
			Cond: getIfErrNotNilExpression(),
			Body: []astgen.ASTStmt{statement.NewReturn(generateNewRestError("StatusBadRequest"))},
		})
	}
	return body, nil
}

// rest.NewError(err, rest.StatusCode(http.$statusCode))
func generateNewRestError(statusCode string) *expression.CallExpression {
	return expression.NewCallFunction(restImportPackage, "NewError",
		expression.VariableVal(errorName),
		expression.NewCallFunction(restImportPackage, "StatusCode",
			expression.NewSelector(
				expression.VariableVal(httpPackageName),
				statusCode,
			),
		),
	)
}

func getAuthStatements(auth *spec.AuthType, info types.PkgInfo) ([]astgen.ASTStmt, error) {
	var body []astgen.ASTStmt
	if auth == nil {
		return body, nil
	}

	if headerAuth, err := visitors.GetPossibleHeaderAuth(*auth); err != nil {
		return nil, err
	} else if headerAuth != nil {
		body = append(body,
			//	authHeader, err := rest.ParseBearerTokenHeader(req)
			//	if err != nil {
			//		return rest.NewError(err, rest.StatusCode(http.StatusForbidden))
			//	}
			&statement.Assignment{
				LHS: []astgen.ASTExpr{
					expression.VariableVal(authHeaderVar),
					expression.VariableVal(errorName),
				},
				Tok: token.DEFINE,
				RHS: expression.NewCallFunction(restImportPackage, funcParseBearerTokenHeader, expression.VariableVal(requestVarName)),
			},
			&statement.If{
				Cond: getIfErrNotNilExpression(),
				Body: []astgen.ASTStmt{statement.NewReturn(generateNewRestError("StatusForbidden"))},
			},
		)
		return body, nil
	}

	if cookieAuth, err := visitors.GetPossibleCookieAuth(*auth); err != nil {
		return nil, err
	} else if cookieAuth != nil {
		//	authCookie, err := req.Cookie("P_TOKEN")
		//	if err != nil {
		//		return rest.NewError(err, rest.StatusCode(http.StatusForbidden))
		//	}
		//	cookieToken := bearertoken.Token(authCookie.Value)
		body = append(body,
			&statement.Assignment{
				LHS: []astgen.ASTExpr{
					expression.VariableVal(authCookieVar),
					expression.VariableVal(errorName),
				},
				Tok: token.DEFINE,
				RHS: expression.NewCallFunction(requestVarName, "Cookie", expression.StringVal(cookieAuth.CookieName)),
			},
			&statement.If{
				Cond: getIfErrNotNilExpression(),
				Body: []astgen.ASTStmt{statement.NewReturn(
					// rest.NewError(err, rest.StatusCode(http.StatusForbidden))
					expression.NewCallFunction(restImportPackage, "NewError",
						expression.VariableVal(errorName),
						expression.NewCallFunction(restImportPackage, "StatusCode",
							expression.NewSelector(
								expression.VariableVal(httpPackageName),
								"StatusForbidden",
							),
						),
					),
				)},
			},
			statement.NewAssignment(
				expression.VariableVal(cookieTokenVar),
				token.DEFINE,
				expression.NewCallExpression(expression.Type(types.Bearertoken.GoType(info)),
					expression.NewSelector(expression.VariableVal(authCookieVar), "Value"),
				),
			),
		)

		return body, nil
	}

	return nil, werror.Error("Unrecognized auth type", werror.SafeParam("authType", auth))
}

func getPathParamStatements(pathParams []visitors.ArgumentDefinitionPathParam, info types.PkgInfo) ([]astgen.ASTStmt, error) {
	if len(pathParams) == 0 {
		return nil, nil
	}
	var body []astgen.ASTStmt
	// Validate path params
	pathParamVar := "pathParams"
	// Use call back to get the path params for this request
	body = append(body, &statement.Assignment{
		LHS: []astgen.ASTExpr{
			expression.VariableVal(pathParamVar),
		},
		Tok: token.DEFINE,
		RHS: expression.NewCallFunction(routerImportPackage, routerPathParamsMapFunc, expression.VariableVal(requestVarName)),
	}, &statement.If{
		Cond: &expression.Binary{
			LHS: expression.VariableVal(pathParamVar),
			Op:  token.EQL,
			RHS: expression.Nil,
		},
		Body: []astgen.ASTStmt{&statement.Return{Values: []astgen.ASTExpr{
			werrorexpressions.CreateWErrorExpression("path params not found on request: ensure this endpoint is registered with wrouter", nil),
		}}},
	})

	for _, pathParam := range pathParams {
		arg := pathParam.ArgumentDefinition

		isString, err := visitors.IsSpecificConjureType(arg.Type, visitors.IsString)
		if err != nil {
			return nil, err
		}

		var strVar expression.VariableVal
		if isString {
			strVar = expression.VariableVal(transforms.SafeName(string(arg.ArgName)))
		} else {
			strVar = expression.VariableVal(arg.ArgName + "Str")
		}

		// For each path param, pull out the value and if it is present in the map
		// argNameStr, ok := pathParams["argName"]
		body = append(body, &statement.Assignment{
			LHS: []astgen.ASTExpr{
				strVar,
				expression.VariableVal("ok"),
			},
			Tok: token.DEFINE,
			RHS: &expression.Index{
				Receiver: expression.VariableVal(pathParamVar),
				Index:    expression.StringVal(visitors.GetParamID(arg)),
			},
		})

		// Check if the param does not exist
		// if !ok { return werror... }
		errorIfNotPresent := werrorexpressions.CreateWErrorExpression("path param not present", map[string]string{"pathParamName": string(arg.ArgName)})
		createWError := &statement.Assignment{
			LHS: []astgen.ASTExpr{expression.VariableVal(errorName)},
			Tok: token.DEFINE,
			RHS: errorIfNotPresent,
		}
		body = append(body, &statement.If{
			Cond: expression.NewUnary(token.NOT, expression.VariableVal(okName)),
			Body: []astgen.ASTStmt{
				createWError,
				&statement.Return{Values: []astgen.ASTExpr{generateNewRestError("StatusBadRequest")}},
			},
		})

		// type-specific unmarshal behavior
		if !isString {
			argName := spec.ArgumentName(transforms.SafeName(string(arg.ArgName)))
			paramStmts, err := visitors.ParseStringParam(argName, arg.Type, strVar, info)
			if err != nil {
				return nil, err
			}
			body = append(body, paramStmts...)
		}
	}
	return body, nil
}

func getHeaderParamStatements(headerParams []visitors.ArgumentDefinitionHeaderParam, info types.PkgInfo) ([]astgen.ASTStmt, error) {
	var body []astgen.ASTStmt
	for _, headerParam := range headerParams {
		arg := headerParam.ArgumentDefinition
		// Pull out the header from the request
		// req.Header.Get("paramID")
		getHeader := &expression.CallExpression{
			Function: &expression.Selector{
				Receiver: &expression.Selector{
					Receiver: expression.VariableVal(requestVarName),
					Selector: requestHeaderFunc,
				},
				Selector: "Get",
			},
			Args: []astgen.ASTExpr{
				expression.StringVal(visitors.GetParamID(headerParam.ArgumentDefinition)),
			},
		}
		// type-specific unmarshal behavior
		// TODO (bmoylan): lists are unimplemented right now, but we _could_ iterate through the raw map and pull them out.
		argName := spec.ArgumentName(transforms.SafeName(string(arg.ArgName)))
		paramStmts, err := visitors.ParseStringParam(argName, arg.Type, getHeader, info)
		if err != nil {
			return nil, err
		}
		body = append(body, paramStmts...)
	}
	return body, nil
}

func getQueryParamStatements(queryParams []visitors.ArgumentDefinitionQueryParam, info types.PkgInfo) ([]astgen.ASTStmt, error) {
	var body []astgen.ASTStmt
	for _, queryParam := range queryParams {
		arg := queryParam.ArgumentDefinition
		// Pull out the query param from the request URL
		// req.URL.Query.Get("paramID")
		getQuery := &expression.CallExpression{
			Function: &expression.Selector{
				Receiver: &expression.CallExpression{
					Function: &expression.Selector{
						Receiver: &expression.Selector{
							Receiver: expression.VariableVal(requestVarName),
							Selector: requestURLField,
						},
						Selector: urlQueryFunc,
					},
				},
				Selector: "Get",
			},
			Args: []astgen.ASTExpr{
				expression.StringVal(visitors.GetParamID(queryParam.ArgumentDefinition)),
			},
		}
		ifErrNotNilReturnErrStatement("err", nil)
		// type-specific unmarshal behavior
		// TODO(bmoylan): lists are unimplemented right now, but we _could_ iterate through the raw map and pull them out.
		argName := spec.ArgumentName(transforms.SafeName(string(arg.ArgName)))
		paramStmts, err := visitors.ParseStringParam(argName, arg.Type, getQuery, info)
		if err != nil {
			return nil, err
		}
		body = append(body, paramStmts...)
	}
	return body, nil
}

func getHandlerStruct(serviceDefinition spec.ServiceDefinition) *decl.Struct {
	return &decl.Struct{
		Name: getHandlerStuctName(serviceDefinition),
		StructType: expression.StructType{
			Fields: []*expression.StructField{
				{
					Name: implName,
					Type: expression.Type(serviceDefinition.ServiceName.Name),
				},
			},
		},
	}
}

func getIfErrNotNilReturnErrExpression() astgen.ASTStmt {
	return &statement.If{
		Cond: getIfErrNotNilExpression(),
		Body: []astgen.ASTStmt{&statement.Return{Values: []astgen.ASTExpr{expression.VariableVal(errorName)}}},
	}
}

func getIfErrNotNilExpression() astgen.ASTExpr {
	return &expression.Binary{
		LHS: expression.VariableVal(errorName),
		Op:  token.NEQ,
		RHS: expression.Nil,
	}
}

func getHandlerStuctName(serviceDefinition spec.ServiceDefinition) string {
	name := serviceDefinition.ServiceName.Name
	firstCharLower := strings.ToLower(string(name[0]))
	return strings.Join([]string{firstCharLower, name[1:], handlerStructNameSuffix}, "")
}

func getReceiverName(serviceDefinition spec.ServiceDefinition) string {
	return string(getHandlerStuctName(serviceDefinition)[0])
}

// rest.NewJSONHandler(funcExpr, rest.StatusCodeMapper, rest.ErrHandler)
func astForRestJSONHandler(funcExpr astgen.ASTExpr) astgen.ASTExpr {
	return expression.NewCallFunction(restImportPackage, "NewJSONHandler",
		funcExpr,
		expression.NewSelector(expression.VariableVal(restImportPackage), "StatusCodeMapper"),
		expression.NewSelector(expression.VariableVal(restImportPackage), "ErrHandler"),
	)
}
