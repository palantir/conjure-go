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
	cliConfigFileVarName  = "configFile"
	defaultConfigFilePath = "../var/conf/configuration.yml"

	loadConfigFuncName    = "loadConfig"
	getCLIContextFuncName = "getCLIContext"
)

func writeCLIType(file *jen.Group, services []*types.ServiceDefinition) {
	writeCLIRoot(file)
	for _, service := range services {
		writeCommandsForService(file, service)
	}
	writeInitAndSharedFuncs(file, services)
}

func writeCLIRoot(file *jen.Group) {
	file.Type().Id(cliConfigTypeName).Struct(
		jen.Id("Client").Add(snip.CGRClientClientConfig())).Line()
	file.Var().Id("configFile").Op("*").String().Line()
}

func writeInitAndSharedFuncs(file *jen.Group, services []*types.ServiceDefinition) {
	astForLoadCLIConfig(file)
	astForGetCLIContext(file)
	astForRegisterCommands(file, services)
	astInitFunc(file, services)
}

func astInitFunc(file *jen.Group, services []*types.ServiceDefinition) {
	file.Func().Id("init").Params().BlockFunc(func(g *jen.Group) {
		astForInitFuncBody(g, services)
	})
}

func astForGetCLIContext(file *jen.Group) {
	file.Add(jen.Func().Id(getCLIContextFuncName).
		Params().
		Params(snip.Context()).
		BlockFunc(func(g *jen.Group) {
			astForGetCLIContextBody(g)
		}))
}

func astForLoadCLIConfig(file *jen.Group) {
	file.Add(jen.Func().Id(loadConfigFuncName).
		Params(snip.ContextVar()).
		Params(
			jen.Id(cliConfigTypeName),
			jen.Error()).
		BlockFunc(func(g *jen.Group) {
			astForLoadCLIConfigBody(g)
		}))
}

func astForLoadCLIConfigBody(file *jen.Group) {
	returnErrBlock := jen.If(jen.Id("err").Op("!=").Nil()).Block(
		jen.Return(jen.Id("emptyConfig"), jen.Id("err"))).Clone

	file.Var().Id("emptyConfig").Id(cliConfigTypeName)
	file.If(jen.Id(cliConfigFileVarName).Op("==").Nil()).Block(
		jen.Return(jen.Id("emptyConfig"),
			snip.WerrorErrorContext().Call(jen.Id("ctx"), jen.Lit("config file location must be specified"))))
	file.List(jen.Id("confBytes"), jen.Id("err")).Op(":=").
		Add(snip.IOUtilReadFile()).Call(jen.Op("*").Id(cliConfigFileVarName))
	file.Add(returnErrBlock())
	file.Var().Id("conf").Id(cliConfigTypeName)
	file.Id("err").Op("=").Add(snip.YamlUnmarshal()).Call(jen.Id("confBytes"), jen.Op("&").Id("conf"))
	file.Add(returnErrBlock())
	file.Return(jen.Id("conf"), jen.Nil())
}

/*
	astForGetCLIContextBody returns the body of the getCLIContext func, used to initialize a context for logging on
	each command invocation.

	func getCLIContext() context.Context {
		ctx := context.Background()
		wlog.SetDefaultLoggerProvider(wlogzap.LoggerProvider())
		ctx = svc1log.WithLogger(ctx, svc1log.New(os.Stdout, wlog.DebugLevel))
		traceLogger := trc1log.DefaultLogger()
		ctx = trc1log.WithLogger(ctx, traceLogger)
		ctx = evt2log.WithLogger(ctx, evt2log.New(os.Stdout))
		tracer, err := wzipkin.NewTracer(traceLogger)
		if err != nil {
			return ctx
		}
		return wtracing.ContextWithTracer(ctx, tracer)
	}
*/
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

// astForRegisterCommands renders the RegisterCommand function, which provides the binding for registering each
// command root (corresponding to a service definition) to the root cobra CLI command.
func astForRegisterCommands(file *jen.Group, services []*types.ServiceDefinition) {
	file.Func().Id("RegisterCommands").
		Params(jen.Id("rootCmd").Op("*").Add(snip.CobraCommand())).
		BlockFunc(func(g *jen.Group) {
			astForRegisterCommandsBody(g, services)
		})
}

func astForRegisterCommandsBody(file *jen.Group, services []*types.ServiceDefinition) {
	for _, service := range services {
		file.Id("rootCmd").Dot("AddCommand").
			Call(jen.Id(getRootServiceCommandName(service.Name)))
	}
}

// astForInitFuncBody renders the init func that builds both the hierarchy of subcommands for each service and
// binds argument flags to each subcommand.
func astForInitFuncBody(file *jen.Group, services []*types.ServiceDefinition) {
	lastIdx := len(services) - 1
	for idx, service := range services {
		file.Comment(fmt.Sprintf("%s commands and flags", service.Name))
		astForInitServiceCommand(file, service)
		if idx != lastIdx {
			file.Line()
		}
	}
}

func astForInitServiceCommand(file *jen.Group, service *types.ServiceDefinition) {
	file.Id(getRootServiceCommandName(service.Name)).
		Dot("PersistentFlags").Call().
		Dot("StringVarP").Call(
		jen.Id(cliConfigFileVarName), jen.Lit("conf"), jen.Lit(""), jen.Lit(defaultConfigFilePath), jen.Lit("The configuration file is optional. The default path is ./var/conf/configuration.yml."))

	for _, endpoint := range service.Endpoints {
		if hasBinaryArgs(endpoint) {
			continue
		}
		astForInitEndpointDefinition(file, service, endpoint)
	}
}

func astForInitEndpointDefinition(file *jen.Group, service *types.ServiceDefinition, endpoint *types.EndpointDefinition) {
	rootCmd := jen.Id(getRootServiceCommandName(service.Name)).Clone
	endpointCmd := jen.Id(getEndpointCommandName(service.Name, endpoint.EndpointName)).Clone

	file.Add(rootCmd()).Dot("AddCommand").Call(endpointCmd())
	for _, param := range endpoint.Params {
		optionality := "a required param"
		if param.Type.IsOptional() {
			optionality = "an optional param"
		}
		argDocs := ""
		if len(param.Docs) > 0 {
			argDocs = fmt.Sprintf(" Argument docs: %s", param.Docs)
		}
		file.Add(endpointCmd()).Dot("PersistentFlags").Call().
			Dot("StringVarP").Call(
			jen.Id(getFlagVarName(service.Name, endpoint.EndpointName, param.Name)),
			jen.Lit(param.Name),
			jen.Lit(""),
			jen.Lit(""),
			jen.Lit(fmt.Sprintf("%s is %s.%s", param.Name, optionality, argDocs)))
	}
	if endpoint.CookieAuth != nil || endpoint.HeaderAuth {
		file.Add(endpointCmd()).Dot("PersistentFlags").Call().
			Dot("StringVarP").Call(
			jen.Id(getFlagAuthVarName(service.Name, endpoint.EndpointName)),
			jen.Lit("bearer_token"),
			jen.Lit(""),
			jen.Lit(""),
			jen.Lit(fmt.Sprintf("bearer_token is a required field.")))
	}
}

func writeCommandsForService(file *jen.Group, serviceDef *types.ServiceDefinition) {
	file.Comment(fmt.Sprintf("// Commands for %s", serviceDef.Name)).Line()
	file.Var().Id(getRootServiceCommandName(serviceDef.Name)).Op("=").Op("&").Add(snip.CobraCommand()).Values(jen.Dict{
		jen.Id("Use"):   jen.Lit(transforms.Private(serviceDef.Name)),
		jen.Id("Short"): jen.Lit(fmt.Sprintf("Runs commands on the %s", serviceDef.Name)),
	}).Line()
	astForServiceClient(file, serviceDef)
	for _, endpoint := range serviceDef.Endpoints {
		if hasBinaryArgs(endpoint) {
			continue
		}
		astForEndpointCommand(file, serviceDef, endpoint)
	}
	file.Line().Line()
}

func astForServiceClient(file *jen.Group, service *types.ServiceDefinition) {
	file.Func().Id(getServiceClientFuncName(service.Name)).
		Params(snip.ContextVar()).
		Params(
			jen.Id(getServiceClientName(service.Name)),
			jen.Error()).
		BlockFunc(func(g *jen.Group) {
			astForServiceClientBody(g, service)
		})
}

func astForServiceClientBody(file *jen.Group, service *types.ServiceDefinition) {
	file.List(jen.Id("conf"), jen.Err()).Op(":=").Id(loadConfigFuncName).
		Call(jen.Id("ctx"))
	file.If(jen.Err().Op("!=").Nil()).Block(
		jen.Return(jen.Nil(), snip.WerrorWrapContext().
			Call(jen.Id("ctx"), jen.Err(), jen.Lit("failed to load CLI configuration file"))))
	file.List(jen.Id("client"), jen.Err()).Op(":=").
		Add(snip.CGRClientNewClient()).
		Call(snip.CGRClientWithConfig().
			Call(jen.Id("conf").Dot("Client")))
	file.If(jen.Err().Op("!=").Nil()).Block(
		jen.Return(jen.Nil(), snip.WerrorWrapContext().
			Call(jen.Id("ctx"), jen.Err(), jen.Lit("failed to create client with provided config"))))
	file.Return(jen.Id("New"+getServiceClientName(service.Name)).Call(jen.Id("client")), jen.Nil())
}

func astForEndpointCommand(file *jen.Group, service *types.ServiceDefinition, endpoint *types.EndpointDefinition) {
	endpointCmdRun := getEndpointCommandRunName(service.Name, endpoint.EndpointName)
	file.Var().Id(getEndpointCommandName(service.Name, endpoint.EndpointName)).
		Op("=").
		Op("&").Add(snip.CobraCommand()).Values(jen.Dict{
		jen.Id("Use"):   jen.Lit(transforms.Private(endpoint.EndpointName)),
		jen.Id("Short"): jen.Lit(fmt.Sprintf("Calls the %s endpoint", endpoint.EndpointName)),
		jen.Id("RunE"):  jen.Id(endpointCmdRun),
	}).Line()
	for _, param := range endpoint.Params {
		file.Var().Id(getFlagVarName(service.Name, endpoint.EndpointName, param.Name)).Op("*").String()
	}
	if endpoint.CookieAuth != nil || endpoint.HeaderAuth {
		file.Var().Id(getFlagAuthVarName(service.Name, endpoint.EndpointName)).Op("*").String()
	}
	file.Line()
	file.Func().Id(endpointCmdRun).Params(
		jen.Id("_").Op("*").Add(snip.CobraCommand()),
		jen.Id("_").Index().String()).
		Params(jen.Error()).
		BlockFunc(func(g *jen.Group) {
			astForEndpointCommandBody(g, service, endpoint)
		}).Line()
	file.Func().Id(getEndpointCommandRunInternalName(service.Name, endpoint.EndpointName)).
		Params(jen.Add(snip.ContextVar()), jen.Id("client").Id(getServiceClientName(service.Name))).
		Params(jen.Error()).
		BlockFunc(func(g *jen.Group) {
			astForEndpointCommandInternalBody(g, service, endpoint)
		})
}

func astForEndpointCommandBody(file *jen.Group, service *types.ServiceDefinition, endpoint *types.EndpointDefinition) {
	file.Id("ctx").Op(":=").Id(getCLIContextFuncName).Call()
	file.List(jen.Id("client"), jen.Err()).Op(":=").Id(getServiceClientFuncName(service.Name)).
		Call(jen.Id("ctx"))
	file.If(jen.Err().Op("!=").Nil()).Block(
		jen.Return(snip.WerrorWrapContext().
			Call(jen.Id("ctx"), jen.Err(), jen.Lit("failed to initialize client"))))
	file.Return(
		jen.Id(getEndpointCommandRunInternalName(service.Name, endpoint.EndpointName)).Call(
			jen.Id("ctx"),
			jen.Id("client"),
		),
	)
}

// astForEndpointCommandInternalBody renders the internal implementation of each command, taking both a context and a client
// as an argument. This enables injecting a mocked client for unit testing.
func astForEndpointCommandInternalBody(file *jen.Group, service *types.ServiceDefinition, endpoint *types.EndpointDefinition) {
	clientArgList := make([]jen.Code, 0, len(endpoint.Params)+1)
	clientArgList = append(clientArgList, jen.Id("ctx"))

	file.Var().Err().Error().Line()

	if endpoint.CookieAuth != nil || endpoint.HeaderAuth {
		authFlagVar := getFlagAuthVarName(service.Name, endpoint.EndpointName)
		param := &types.EndpointArgumentDefinition{
			Name: "__authVar",
			Type: types.Bearertoken{},
		}
		astForEndpointParam(file, authFlagVar, param)
		clientArgList = append(clientArgList, jen.Id("__authVarArg"))
	}

	for _, param := range endpoint.Params {
		flagVarName := getFlagVarName(service.Name, endpoint.EndpointName, param.Name)
		astForEndpointParam(file, flagVarName, param)
		file.Line()
		clientArgList = append(clientArgList, jen.Id(getArgName(param)))
	}

	clientCallCode := jen.Id("client").Dot(transforms.Export(endpoint.EndpointName)).
		Call(clientArgList...)
	if endpoint.Returns == nil {
		file.Err().Op("=").Add(clientCallCode)
		file.Return(jen.Err())
		return
	}

	file.List(jen.Id("result"), jen.Err()).Op(":=").Add(clientCallCode)
	file.If().Err().Op("!=").Nil().Block(
		jen.Return(jen.Err()))
	astForPrintResult(file, endpoint)
}

func astForEndpointParam(file *jen.Group, flagVarName string, param *types.EndpointArgumentDefinition) {
	argName := getArgName(param)

	// TODO: Add support for reading file from path as binary input
	if param.Type.IsBinary() {
		file.Var().Id(argName).Index().Byte()
		file.Id("panic").Call(jen.Lit("Commands with binary arguments are not yet supported."))
		return
	}

	// For optional params, handle only if flag not nil
	if param.Type.IsOptional() {
		flagVarNameDeref := flagVarName + "Deref"
		file.Var().Id(flagVarNameDeref).String()
		file.If(jen.Id(flagVarName).Op("!=").Nil()).Block(
			jen.Id(flagVarNameDeref).Op("=").Op("*").Id(flagVarName))
		astForEndpointParamInner(file, argName, jen.Id(flagVarNameDeref), param)
		return
	}

	// otherwise, error if nil, handle otherwise
	file.If(jen.Id(flagVarName).Op("==").Nil()).Block(
		jen.Return(snip.WerrorErrorContext().Call(
			jen.Id("ctx"), jen.Lit(fmt.Sprintf("%s is a required argument", argName)))))
	astForEndpointParamInner(file, argName, jen.Op("*").Id(flagVarName), param)
}

func astForEndpointParamInner(file *jen.Group, argName string, flagVar jen.Code, param *types.EndpointArgumentDefinition) {
	ctxExpr := jen.Id("ctx")
	if param.Type.IsCollection() || param.Type.ContainsStrictFields() {
		astForEndpointCollectionParam(file, argName, flagVar, param)
		return
	}
	astForDecodeHTTPParam(file, param.Name, param.Type, argName, ctxExpr, flagVar)
}

func astForEndpointCollectionParam(file *jen.Group, argName string, flagVar jen.Code, param *types.EndpointArgumentDefinition) {
	argBytesName := argName + "Bytes"

	file.Var().Id(argName).Add(param.Type.Code())
	file.Id(argBytesName).Op(":=").Index().Byte().Parens(flagVar)
	file.If(
		jen.Err().Op(":=").Add(snip.CGRCodecsJSON().Dot("Decode")).Call(
			jen.Add(snip.ByteReader).Call(jen.Id(argBytesName)),
			jen.Op("&").Id(argName),
		),
		jen.Err().Op("!=").Nil(),
	).Block(jen.Return(snip.CGRErrorsWrapWithInvalidArgument().Call(jen.Err())))
}

func astForPrintResult(file *jen.Group, endpoint *types.EndpointDefinition) {
	returnType := *endpoint.Returns
	switch {
	case returnType.IsBinary():
		resultVar := jen.Id("result").Clone
		if returnType.IsOptional() {
			file.If().Id("result").Op("==").Nil().Block(
				jen.Return(jen.Nil()))
			file.Id("resultDeref").Op(":=").Op("*").Id("result")
			resultVar = jen.Id("resultDeref").Clone
		}
		file.List(jen.Id("_"), jen.Err()).Op("=").Add(snip.IOCopy).Call(snip.OSStdout(), resultVar())
		file.If(jen.Err().Op("!=").Nil().Block(
			jen.Return(jen.Add(snip.WerrorWrapContext().Call(jen.Id("ctx"), jen.Err(), jen.Lit("failed to write result bytes to stdout"))))),
		)
		file.Return(resultVar().Dot("Close").Call())
	case returnType.IsText():
		file.Add(snip.FmtPrintf()).Call(jen.Lit("%v\n"), jen.Id("result"))
		file.Return(jen.Nil())
	default:
		file.List(jen.Id("resultBytes"), jen.Err()).Op(":=").Add(snip.JSONMarshalIndent()).Call(
			jen.Id("result"), jen.Lit(""), jen.Lit("    "))
		file.If(jen.Err().Op("!=").Nil()).Block(
			jen.Add(snip.FmtPrintf()).Call(jen.Lit("Failed to marshal to json with err: %v\n\nPrinting as string:\n%v\n"), jen.Err(), jen.Id("result")),
			jen.Return(jen.Nil()))
		file.Add(snip.FmtPrintf().Call(jen.Lit("%v\n"), jen.String().Parens(jen.Id("resultBytes"))))
		file.Return(jen.Nil())
	}
}

func getRootServiceCommandName(serviceName string) string {
	return "Root" + serviceName + "Cmd"
}

func getEndpointCommandName(serviceName, endpointName string) string {
	return serviceName + endpointName + "Cmd"
}

func getEndpointCommandRunName(serviceName, endpointName string) string {
	return transforms.Private(getEndpointCommandName(serviceName, endpointName) + "Run")
}

func getEndpointCommandRunInternalName(serviceName, endpointName string) string {
	return getEndpointCommandRunName(serviceName, endpointName) + "Internal"
}

func getServiceClientName(serviceName string) string {
	return serviceName + "Client"
}

func getServiceClientFuncName(serviceName string) string {
	return "get" + getServiceClientName(serviceName)
}

func getArgName(param *types.EndpointArgumentDefinition) string {
	return transforms.Private(param.Name + "Arg")
}

func getFlagAuthVarName(serviceName, endpointName string) string {
	return transforms.Private(serviceName) + "_" +
		transforms.Private(endpointName) + "__auth"
}

func getFlagVarName(serviceName, endpointName, paramName string) string {
	return transforms.Private(serviceName) + "_" +
		transforms.Private(endpointName) + "_" +
		transforms.Private(paramName)
}

func hasBinaryArgs(endpoint *types.EndpointDefinition) bool {
	for _, param := range endpoint.Params {
		if param.Type.IsBinary() {
			return true
		}
	}
	return false
}
