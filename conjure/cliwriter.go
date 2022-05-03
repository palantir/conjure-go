// Copyright (c) 2022 Palantir Technologies. All rights reserved.
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

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	cliConfigTypeName     = "CLIConfig"
	defaultConfigFilePath = "../var/conf/configuration.yml"

	loadConfigFuncName    = "loadCLIConfig"
	getCLIContextFuncName = "getCLIContext"

	bearerTokenFlagName = "bearer_token"
	confFlagName        = "conf"
)

var (
	// cliReservedArgNames tracks the list of flag names that may not be used directly as they collide
	// with global flags
	cliReservedArgNames = []string{
		confFlagName,
		bearerTokenFlagName,
	}
)

// writeCLIType is the entry point for generating CLI commands from a set of service definitions
func writeCLIType(file *jen.Group, services []*types.ServiceDefinition) {
	writeCLIConfigStruct(file)
	for _, service := range services {
		writeCommandsForService(file, service)
	}
	writeSharedCLIFuncs(file)
}

// writeCLIConfigStruct generates a struct for unmarshaling a client config file
func writeCLIConfigStruct(file *jen.Group) {
	file.Type().Id(cliConfigTypeName).Struct(
		jen.Id("Client").Add(snip.CGRClientClientConfig()).Tag(map[string]string{"yaml": ",inline"})).Line()
}

// writeSharedCLIFuncs writes a set of shared functions used across all services
func writeSharedCLIFuncs(file *jen.Group) {
	astForLoadCLIConfig(file)
	astForGetCLIContext(file)
}

// astForLoadCLIConfig writes a function for getting the config file for configuring service clients
func astForLoadCLIConfig(file *jen.Group) {
	file.Add(jen.Func().Id(loadConfigFuncName).
		Params(snip.ContextVar(), jen.Id("flags").Op("*").Add(snip.PflagsFlagset())).
		Params(
			jen.Id(cliConfigTypeName),
			jen.Error()).
		BlockFunc(func(g *jen.Group) {
			astForLoadCLIConfigBody(g)
		})).Line()
}

// astForLoadCLIConfigBody implements the getCLIConfig function
func astForLoadCLIConfigBody(file *jen.Group) {
	configPathVar := "configPath"
	returnErrBlock := jen.If(jen.Id("err").Op("!=").Nil()).Block(
		jen.Return(jen.Id("emptyConfig"), jen.Err())).Clone

	// Get config path from command flag
	file.Var().Id("emptyConfig").Id(cliConfigTypeName)
	file.List(jen.Id(configPathVar), jen.Err()).Op(":=").Id("flags").Dot("GetString").Call(jen.Lit(confFlagName))
	file.If(jen.Err().Op("!=").Nil().Op("||").Id(configPathVar).Op("==").Lit("")).Block(
		jen.Return(jen.Id("emptyConfig"),
			snip.WerrorWrapContext().Call(jen.Id("ctx"), jen.Err(), jen.Lit("config file location must be specified"))))

	// Read config bytes from disk
	file.List(jen.Id("confBytes"), jen.Err()).Op(":=").
		Add(snip.OSReadFile()).Call(jen.Id(configPathVar))
	file.Add(returnErrBlock())

	// Unmarshal client config and return
	file.Var().Id("conf").Id(cliConfigTypeName)
	file.Id("err").Op("=").Add(snip.YamlUnmarshal()).Call(jen.Id("confBytes"), jen.Op("&").Id("conf"))
	file.Add(returnErrBlock())
	file.Return(jen.Id("conf"), jen.Nil())
}

// astForGetCLIContext implements getCLIContext, which returns a context with initialized loggers
func astForGetCLIContext(file *jen.Group) {
	file.Add(jen.Func().Id(getCLIContextFuncName).
		Params().
		Params(snip.Context()).
		BlockFunc(func(g *jen.Group) {
			astForGetCLIContextBody(g)
		})).Line()
}

// astForGetCLIContextBody implements the getCLIContext function
func astForGetCLIContextBody(file *jen.Group) {
	stdout := jen.Qual("os", "Stdout").Clone
	file.Id("ctx").Op(":=").Add(snip.ContextBackground()).Call()
	file.Add(snip.WGLLogSetDefaultLoggerProvider()).Call(snip.WGLWlogZapLoggerProvider().Call())
	file.Id("ctx").Op("=").Add(snip.WGLSvc1logWithLogger()).Call(
		jen.Id("ctx"), snip.WGLSvc1logNew().Call(stdout(), snip.WGLLogDebugLevel()))
	file.Id("traceLogger").Op(":=").Add(snip.WGLTrc1logDefaultLogger()).Call()
	file.Id("ctx").Op("=").Add(snip.WGLTrc1logWithLogger()).Call(
		jen.Id("ctx"), jen.Id("traceLogger"))
	file.Id("ctx").Op("=").Add(snip.WGLEvt2logWithLogger()).Call(
		jen.Id("ctx"), snip.WGLEvt2logNew().Call(stdout()))
	file.List(jen.Id("tracer"), jen.Id("err")).Op(":=").Add(snip.WGTZipkinNewTracer()).Call(jen.Id("traceLogger"))
	file.If(jen.Id("err").Op("!=").Nil()).Block(
		jen.Return(jen.Id("ctx")))
	file.Return(snip.WGTContextWithTracer().Call(
		jen.Id("ctx"), jen.Id("tracer")))
}

// writeCommandsForService creates a root command for a service. to which subcommands will be attached for each
// endpoint definition.
func writeCommandsForService(file *jen.Group, serviceDef *types.ServiceDefinition) {
	serviceName := serviceDef.Name
	file.Comment(fmt.Sprintf("// Commands for %s", serviceName)).Line()

	// Generate client provider type and default implementation
	astForServiceClientProvider(file, serviceDef)

	// Generate root service command type
	astForRootServiceCommand(file, serviceDef)

	// For each endpoint defined on the service, generate a subcommand
	for _, endpoint := range serviceDef.Endpoints {
		// TODO: Add support for endpoints with binary type params
		if hasBinaryArgs(endpoint) {
			continue
		}
		astForEndpointCommand(file, serviceDef, endpoint)
	}
	file.Line().Line()
}

// astForServiceClientProvider writes a client provider interface and a default implementation
func astForServiceClientProvider(file *jen.Group, service *types.ServiceDefinition) {
	serviceName := service.Name

	// Define client provider interface. This simplifies testing by making it possible to provide an
	// alternate implementation that injects a mocked client.
	file.Type().Id(getCLIServiceClientProviderName(serviceName)).Interface(
		jen.Id("Get").Params(snip.ContextVar(), jen.Id("flags").Op("*").Add(snip.PflagsFlagset())).
			Params(jen.Id(getServiceClientName(serviceName)), jen.Error())).Line()

	// Generate default implementation
	file.Type().Id(getDefaultClientProviderImplName(serviceName)).Struct().Line()
	file.Func().Id(getNewDefaultClientProviderName(serviceName)).Params().Params(jen.Id(getCLIServiceClientProviderName(serviceName))).Block(
		jen.Return(jen.Id(getDefaultClientProviderImplName(serviceName)).Values())).Line()
	file.Func().Params(jen.Id("d").Id(getDefaultClientProviderImplName(serviceName))).
		Id("Get").
		Params(snip.ContextVar(), jen.Id("flags").Op("*").Add(snip.PflagsFlagset())).
		Params(jen.Id(getServiceClientName(serviceName)), jen.Error()).
		BlockFunc(func(f *jen.Group) {
			astForServiceClientBody(f, service)
		})
}

// astForServiceClientBody is the body of the get client function
func astForServiceClientBody(file *jen.Group, service *types.ServiceDefinition) {
	// retrieve client configuration
	file.List(jen.Id("conf"), jen.Err()).Op(":=").Id(loadConfigFuncName).
		Call(jen.Id("ctx"), jen.Id("flags"))
	file.If(jen.Err().Op("!=").Nil()).Block(
		jen.Return(jen.Nil(), snip.WerrorWrapContext().
			Call(jen.Id("ctx"), jen.Err(), jen.Lit("failed to load CLI configuration file"))))

	// create client based on configuration
	file.List(jen.Id("client"), jen.Err()).Op(":=").
		Add(snip.CGRClientNewClient()).
		Call(snip.CGRClientWithConfig().
			Call(jen.Id("conf").Dot("Client")))
	file.If(jen.Err().Op("!=").Nil()).Block(
		jen.Return(jen.Nil(), snip.WerrorWrapContext().
			Call(jen.Id("ctx"), jen.Err(), jen.Lit("failed to create client with provided config"))))
	file.Return(jen.Id("New"+getServiceClientName(service.Name)).Call(jen.Id("client")), jen.Nil())
}

func astForRootServiceCommand(file *jen.Group, service *types.ServiceDefinition) {
	serviceName := service.Name

	// Define root service command type
	file.Type().Id(getRootServiceCommandName(serviceName)).Struct(
		jen.Id("clientProvider").Id(getCLIServiceClientProviderName(serviceName))).Line()

	// Generate default constructor
	file.Func().Id(getNewRootServiceCommandName(serviceName)).Params().Params(jen.Op("*").Add(snip.CobraCommand())).
		Block(
			jen.Return(
				jen.Id(getNewRootServiceCommandWithClientProviderName(serviceName)).
					Call(jen.Id(getNewDefaultClientProviderName(serviceName)).Call()))).Line()

	// Generate constructor that accepts a client provider argument, enabling injection of a client for testing
	file.Func().Id(getNewRootServiceCommandWithClientProviderName(serviceName)).
		Params(jen.Id("clientProvider").Id(getCLIServiceClientProviderName(serviceName))).
		Params(jen.Op("*").Add(snip.CobraCommand())).
		BlockFunc(func(f *jen.Group) {
			astForRootServiceCommandConstructorBody(f, service)
		}).Line()
}

func astForRootServiceCommandConstructorBody(file *jen.Group, service *types.ServiceDefinition) {
	serviceName := service.Name

	// Initialize root cobra command for service, as well as persistent conf flag
	file.Id("rootCmd").Op(":=").Op("&").Add(snip.CobraCommand()).Values(jen.Dict{
		jen.Id("Use"):   jen.Lit(transforms.Private(serviceName)),
		jen.Id("Short"): jen.Lit(fmt.Sprintf("Runs commands on the %s", serviceName)),
	})
	file.Id("rootCmd").
		Dot("PersistentFlags").Call().
		Dot("String").Call(
		jen.Lit(confFlagName), jen.Lit(defaultConfigFilePath), jen.Lit("The configuration file is optional. The default path is ./var/conf/configuration.yml.")).Line()

	// Initialize service command
	file.Id("cliCommand").Op(":=").Id(getRootServiceCommandName(serviceName)).Values(jen.Dict{
		jen.Id("clientProvider"): jen.Id("clientProvider"),
	}).Line()

	// Initialize and configure endpoint commands
	for _, endpoint := range service.Endpoints {
		// TODO: Add support for endpoints with binary type params
		if hasBinaryArgs(endpoint) {
			continue
		}
		astForEndpointFlags(file, service, endpoint)
	}

	file.Return(jen.Id("rootCmd"))
}

// astForEndpointFlags registers each endpoint subcommand and associated flags
func astForEndpointFlags(file *jen.Group, service *types.ServiceDefinition, endpoint *types.EndpointDefinition) {
	// Initialize endpoint command
	endpointCmd := getEndpointCommandName(service.Name, endpoint.EndpointName)
	endpointCmdRun := jen.Id("cliCommand").Dot(getEndpointCommandRunName(service.Name, endpoint.EndpointName))
	endpointDocs := fmt.Sprintf("Calls the %s endpoint.", endpoint.EndpointName)
	if len(endpoint.Docs) > 0 {
		endpointDocs = string(endpoint.Docs)
	}
	file.Id(endpointCmd).
		Op(":=").
		Op("&").Add(snip.CobraCommand()).Values(jen.Dict{
		jen.Id("Use"):   jen.Lit(transforms.Private(endpoint.EndpointName)),
		jen.Id("Short"): jen.Lit(endpointDocs),
		jen.Id("RunE"):  endpointCmdRun,
	})

	// Register endpoint subcommand on root service command
	file.Id("rootCmd").
		Dot("AddCommand").
		Call(jen.Id(endpointCmd))

	// Register a flag for each endpoint param
	for _, param := range endpoint.Params {
		optionality := "Required"
		if param.Type.IsOptional() {
			optionality = "Optional"
		}
		argDocs := ""
		if len(param.Docs) > 0 {
			argDocs = string(param.Docs)
		}
		file.Id(endpointCmd).Dot("Flags").Call().
			Dot("String").Call(
			jen.Lit(getFlagName(param.Name)),
			jen.Lit(""),
			jen.Lit(fmt.Sprintf("%s. %s", optionality, argDocs)))
	}

	// Register an additional bearer token flag if auth is enabled for the endpoint
	if endpoint.CookieAuth != nil || endpoint.HeaderAuth {
		file.Id(endpointCmd).Dot("Flags").Call().
			Dot("String").Call(
			jen.Lit(bearerTokenFlagName),
			jen.Lit(""),
			jen.Lit(fmt.Sprintf("bearer_token is a required field.")))
	}

	file.Line()
}

// astForEndpointCommand creates a subcommand for a service endpoint and its associated function called on execution
func astForEndpointCommand(file *jen.Group, service *types.ServiceDefinition, endpoint *types.EndpointDefinition) {
	file.Func().
		Params(jen.Id("c").Id(getRootServiceCommandName(service.Name))).
		Id(getEndpointCommandRunName(service.Name, endpoint.EndpointName)).
		Params(
			jen.Id("cmd").Op("*").Add(snip.CobraCommand()),
			jen.Id("_").Index().String()).
		Params(jen.Error()).
		BlockFunc(func(g *jen.Group) {
			astForEndpointCommandBody(g, service, endpoint)
		}).Line()
}

// astForEndpointCommandBody generates the command function, which initializes a client before calling the internal
// command function
func astForEndpointCommandBody(file *jen.Group, service *types.ServiceDefinition, endpoint *types.EndpointDefinition) {
	// Get CLI with logging
	file.Id("ctx").Op(":=").Id(getCLIContextFuncName).Call()

	// Get flags from command
	file.Id("flags").Op(":=").Id("cmd").Dot("Flags").Call()

	// Get client for service
	file.List(jen.Id("client"), jen.Err()).
		Op(":=").
		Id("c").
		Dot("clientProvider").
		Dot("Get").
		Call(jen.Id("ctx"), jen.Id("flags"))
	file.If(jen.Err().Op("!=").Nil()).Block(
		jen.Return(snip.WerrorWrapContext().
			Call(jen.Id("ctx"), jen.Err(), jen.Lit("failed to initialize client"))))

	// create arg list, which will be used to call the client once all arguments are parsed
	clientArgList := make([]jen.Code, 0, len(endpoint.Params)+1)
	clientArgList = append(clientArgList, jen.Id("ctx"))

	// If auth is enabled, we must inject and handle an additional token param, which is always passed as the second
	// argument when present
	if endpoint.CookieAuth != nil || endpoint.HeaderAuth {
		param := &types.EndpointArgumentDefinition{
			Name: "__authVar",
			Type: types.Bearertoken{},
		}
		astForEndpointParam(file, bearerTokenFlagName, param)
		clientArgList = append(clientArgList, jen.Id("__authVarArg"))
	}

	// Parse each endpoint param into an argument for calling the client
	for _, param := range endpoint.Params {
		astForEndpointParam(file, param.Name, param)
		file.Line()
		clientArgList = append(clientArgList, jen.Id(getArgName(param)))
	}

	clientCallCode := jen.Id("client").Dot(transforms.Export(endpoint.EndpointName)).
		Call(clientArgList...)
	// If an endpoint has no return value, we handle only any returned error
	if endpoint.Returns == nil {
		file.Return(clientCallCode)
		return
	}

	// For endpoints with a return value, call the client and print the result unless it returns an error
	file.List(jen.Id("result"), jen.Err()).Op(":=").Add(clientCallCode)
	file.If().Err().Op("!=").Nil().Block(
		jen.Return(jen.Err()))
	astForPrintResult(file, endpoint)
}

// astForEndpointParam handles getting a param value from a flag and parsing it into the type expected by the client
func astForEndpointParam(file *jen.Group, flagName string, param *types.EndpointArgumentDefinition) {
	argName := getArgName(param)

	// Get the param value from the flag
	flagVarNameRaw := flagName + "Raw"
	file.List(jen.Id(flagVarNameRaw), jen.Err()).Op(":=").
		Id("flags").Dot("GetString").Call(jen.Lit(flagName))
	file.If(jen.Err().Op("!=").Nil()).Block(
		jen.Return(
			snip.WerrorWrapContext().Call(jen.Id("ctx"), jen.Err(), jen.Lit(fmt.Sprintf("failed to parse argument %s", param.Name)))),
	)

	// For optional params, always add handling code
	if param.Type.IsOptional() {
		astForEndpointParamInner(file, argName, jen.Id(flagVarNameRaw), param)
		return
	}

	// For required params, return an error if no param value specified
	file.If(jen.Id(flagVarNameRaw).Op("==").Lit("")).Block(
		jen.Return(snip.WerrorErrorContext().Call(
			jen.Id("ctx"), jen.Lit(fmt.Sprintf("%s is a required argument", flagName)))))
	astForEndpointParamInner(file, argName, jen.Id(flagVarNameRaw), param)
}

// astForEndpointParamInner delegates param parsing based on param type
func astForEndpointParamInner(file *jen.Group, argName string, flagVar jen.Code, param *types.EndpointArgumentDefinition) {
	// Collection types are handled via json decoding
	if param.Type.IsCollection() || param.Type.ContainsStrictFields() {
		astForEndpointCollectionParam(file, argName, flagVar, param)
		return
	}
	// All non-complex types delegate to the serverwriter http param decoding
	astForDecodeHTTPParam(file, param.Name, param.Type, argName, jen.Id("ctx"), flagVar)
}

// astForEndpointCollectionParam applies json decoding to handle collection param values
func astForEndpointCollectionParam(file *jen.Group, argName string, flagVar jen.Code, param *types.EndpointArgumentDefinition) {
	// Declare argument var
	file.Var().Id(argName).Add(param.Type.Code())

	// If param is optional, only attempt to decode if value is not empty
	if param.Type.IsOptional() {
		file.If(jen.Add(flagVar).Op("!=").Lit("")).BlockFunc(func(g *jen.Group) {
			astForEndpointCollectionParamDecode(g, argName, flagVar, param)
		})
		return
	}
	astForEndpointCollectionParamDecode(file, argName, flagVar, param)
}

func astForEndpointCollectionParamDecode(file *jen.Group, argName string, flagVar jen.Code, param *types.EndpointArgumentDefinition) {
	argBytesName := argName + "Bytes"

	file.Id(argBytesName).Op(":=").Index().Byte().Parens(flagVar)
	file.If(
		jen.Err().Op(":=").Add(snip.CGRCodecsJSON().Dot("Decode")).Call(
			snip.ByteReader().Call(jen.Id(argBytesName)),
			jen.Op("&").Id(argName),
		),
		jen.Err().Op("!=").Nil(),
	).Block(jen.Return(snip.WerrorWrapContext().
		Call(jen.Id("ctx"), jen.Err(), jen.Lit(fmt.Sprintf("invalid value for %s argument", param.Name)))))
}

// astForPrintResult prints a client result based on return type
func astForPrintResult(file *jen.Group, endpoint *types.EndpointDefinition) {
	returnType := *endpoint.Returns
	switch {
	// Write binary output to STDOUT, enabling it to be piped to a file
	case returnType.IsBinary():
		resultVar := jen.Id("result").Clone
		if returnType.IsOptional() {
			file.If().Id("result").Op("==").Nil().Block(
				jen.Return(jen.Nil()))
			file.Id("resultDeref").Op(":=").Op("*").Id("result")
			resultVar = jen.Id("resultDeref").Clone
		}
		file.List(jen.Id("_"), jen.Err()).Op("=").Add(snip.IOCopy).Call(jen.Id("cmd").Dot("OutOrStdout").Call(), resultVar())
		file.If(jen.Err().Op("!=").Nil().Block(
			jen.Return(jen.Add(snip.WerrorWrapContext().Call(jen.Id("ctx"), jen.Err(), jen.Lit("failed to write result bytes to stdout"))))),
		)
		file.Return(resultVar().Dot("Close").Call())
	// Write simple text results as formatted text
	case returnType.IsText():
		file.Add(snip.FmtFprintf()).Call(jen.Id("cmd").Dot("OutOrStdout").Call(), jen.Lit("%v\n"), jen.Id("result"))
		file.Return(jen.Nil())
	// For any remaining types, including objects, marshal to json and pretty print
	default:
		file.List(jen.Id("resultBytes"), jen.Err()).Op(":=").Add(snip.JSONMarshalIndent()).Call(
			jen.Id("result"), jen.Lit(""), jen.Lit("    "))
		file.If(jen.Err().Op("!=").Nil()).Block(
			jen.Add(snip.FmtPrintf()).Call(jen.Lit("Failed to marshal to json with err: %v\n\nPrinting as string:\n%v\n"), jen.Err(), jen.Id("result")),
			jen.Return(jen.Nil()))
		file.Add(snip.FmtFprintf().Call(jen.Id("cmd").Dot("OutOrStdout").Call(), jen.Lit("%v\n"), jen.String().Parens(jen.Id("resultBytes"))))
		file.Return(jen.Nil())
	}
}

func getRootServiceCommandName(serviceName string) string {
	return fmt.Sprintf("%sCLICommand", serviceName)
}

func getNewRootServiceCommandName(serviceName string) string {
	return "New" + getRootServiceCommandName(serviceName)
}

func getNewRootServiceCommandWithClientProviderName(serviceName string) string {
	return getNewRootServiceCommandName(serviceName) + "WithClientProvider"
}

func getCLIServiceClientProviderName(serviceName string) string {
	return fmt.Sprintf("CLI%sClientProvider", serviceName)
}

func getDefaultClientProviderImplName(serviceName string) string {
	return fmt.Sprintf("default%s", getCLIServiceClientProviderName(serviceName))
}

func getNewDefaultClientProviderName(serviceName string) string {
	return "New" + transforms.Export(getDefaultClientProviderImplName(serviceName))
}

func getEndpointCommandName(serviceName, endpointName string) string {
	return transforms.Private(
		fmt.Sprintf("%s_%s_Cmd", serviceName, transforms.Export(endpointName)))
}

func getEndpointCommandRunName(serviceName, endpointName string) string {
	return transforms.Private(getEndpointCommandName(serviceName, endpointName) + "Run")
}

func getServiceClientName(serviceName string) string {
	return serviceName + "Client"
}

func getArgName(param *types.EndpointArgumentDefinition) string {
	return transforms.Private(param.Name + "Arg")
}

func getFlagName(paramName string) string {
	name := transforms.Private(paramName)
	for _, reservedName := range cliReservedArgNames {
		if name == reservedName {
			return name + "_Arg"
		}
	}
	return name
}

func hasBinaryArgs(endpoint *types.EndpointDefinition) bool {
	for _, param := range endpoint.Params {
		if param.Type.IsBinary() {
			return true
		}
	}
	return false
}
