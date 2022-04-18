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

func astCLIRoot(file *jen.Group) {
	file.Type().Id(cliConfigTypeName).Struct(
		jen.Id("Client").Add(snip.CGRClientClientConfig()))
	file.Var().Id("configFile").Op("*").String()
}

func astInitFunc(file *jen.Group, services []*types.ServiceDefinition) {
	astForLoadCLIConfig(file)
	astForGetCLIContext(file)
	astForRegisterCommands(file, services)
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
			jen.Qual("", cliConfigTypeName),
			jen.Error()).
		BlockFunc(func(g *jen.Group) {
			astForLoadCLIConfigBody(g)
		}))
}

func astForLoadCLIConfigBody(file *jen.Group) {
	returnErrBlock := jen.If(jen.Id("err").Op("!=").Nil()).Block(
		jen.Return(jen.Id("emptyConfig"), jen.Id("err"))).Clone

	file.Var().Id("emptyConfig").Qual("", cliConfigTypeName)
	file.If(jen.Id(cliConfigFileVarName).Op("==").Nil()).Block(
		jen.Return(jen.Id("emptyConfig"),
			snip.WerrorErrorContext().Call(jen.Id("ctx"), jen.Lit("config file location must be specified"))))
	file.List(jen.Id("confBytes"), jen.Id("err")).Op(":=").
		Add(snip.IOUtilReadFile()).Call(jen.Op("*").Id(cliConfigFileVarName))
	file.Add(returnErrBlock())
	file.Var().Id("conf").Qual("", cliConfigTypeName)
	file.Id("err").Op("=").Add(snip.YamlUnmarshal()).Call(jen.Id("confBytes"), jen.Op("&").Id("conf"))
	file.Add(returnErrBlock())
	file.Return(jen.Id("conf"), jen.Nil())
}

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

func astForInitFunc(file *jen.Group, services []*types.ServiceDefinition) {
	for _, service := range services {
		file.Comment(fmt.Sprintf("%s commands and flags", service.Name))

		file.Line()
	}
}

//func astForInitServiceCommand(file *jen.Group, service *types.ServiceDefinition) {
//	file.Id(getRootServiceCommandName(service.Name)).
//		Dot("PersistentFlags").Call().
//		Dot("StringVarP").Call(
//		jen.Id(cliConfigFileVarName), jen.Lit("conf"), jen.Lit(""), jen.Lit(defaultConfigFilePath), jen.Lit("The configuration file is optional. The default path is ./var/conf/configuration.yml."))
//
//	for _, endpoint := range service.Endpoints {
//		astForInitEndpointDefinition(file, service, endpoint)
//	}
//}

//func astForInitEndpointDefinition(file *jen.Group, service *types.ServiceDefinition, endpoint *types.EndpointDefinition) {
//	rootCmd := jen.Id(getRootServiceCommandName(service.Name)).Clone
//	endpointCmd := jen.Id(getEndpointCommandName(service.Name, endpoint.EndpointName)).Clone
//
//	file.Add(rootCmd()).Dot("AddCommand").Call(endpointCmd())
//	for _, param := range endpoint.Params {
//
//	}
//	file.Add(endpointCmd())
//}
//
//func astForInitEndpointParam(file *jen.Group, param *types.EndpointArgumentDefinition, endpointCmd *jen.Statement) {
//	if param.ParamType == types.BodyParam {
//		if param.Type.Code()
//	}
//}

func writeCLIType(file *jen.Group, serviceDef *types.ServiceDefinition) {
	file.Comment(fmt.Sprintf("// Commands for %s", serviceDef.Name)).Line()
	file.Var().Id(getRootServiceCommandName(serviceDef.Name)).Op("=").Op("&").Add(snip.CobraCommand()).Values(jen.Dict{
		jen.Id("Use"):   jen.Lit(transforms.Private(serviceDef.Name)),
		jen.Id("Short"): jen.Lit(fmt.Sprintf("Runs commands on the %s", serviceDef.Name)),
	}).Line()
	astForServiceClient(file, serviceDef)
	for _, endpoint := range serviceDef.Endpoints {
		astForEndpointCommand(file, serviceDef, endpoint)
	}
	file.Line().Line()
}

func astForServiceClient(file *jen.Group, service *types.ServiceDefinition) {
	file.Func().Id(getServiceClientFuncName(service.Name)).Params(
		snip.ContextVar()).
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
	file.Func().Id(endpointCmdRun).Params(
		jen.Id("_").Op("*").Add(snip.CobraCommand()),
		jen.Id("_").Index().String()).
		Params(jen.Error()).
		BlockFunc(func(g *jen.Group) {
			astForEndpointCommandBody(g, service, endpoint)
		})
}

func astForEndpointCommandBody(file *jen.Group, service *types.ServiceDefinition, endpoint *types.EndpointDefinition) {
	file.Id("ctx").Op(":=").Id(getCLIContextFuncName).Call()
	file.List(jen.Id("client"), jen.Err()).Op(":=").Id(getServiceClientFuncName(service.Name)).
		Call(jen.Id("ctx"))
	file.If(jen.Err().Op("!=").Nil()).Block(
		jen.Return(snip.WerrorWrapContext().
			Call(jen.Id("ctx"), jen.Err(), jen.Lit("failed to initialize client"))))
	for _, param := range endpoint.Params {
		astForEndpointParam(file, service, endpoint, param)
	}
}

func astForEndpointParam(file *jen.Group, service *types.ServiceDefinition, endpoint *types.EndpointDefinition, param *types.EndpointArgumentDefinition) {
	argName := getArgName(param)
	flagVarName := getFlagVarName(service.Name, endpoint.EndpointName, param.Name)
	if !param.Type.IsOptional() {
		file.Var().Id(argName).Add(param.Type.Code())
		file.If(jen.Id(flagVarName).Op("==").Nil()).Block(
			jen.Return(snip.WerrorErrorContext().Call(
				jen.Id("ctx"), jen.Lit(fmt.Sprintf("%s is a required argument", argName)))))
		// TODO: Get argument value
	} else {

	}
	//assignment := getAssignmentForParam(param)
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

func getServiceClientName(serviceName string) string {
	return serviceName + "Client"
}

func getServiceClientFuncName(serviceName string) string {
	return "get" + getServiceClientName(serviceName)
}

func getArgName(param *types.EndpointArgumentDefinition) string {
	return transforms.Private(param.Name + "Arg")
}

func getFlagVarName(serviceName, endpointName, paramName string) string {
	return transforms.Private(serviceName) + "_" +
		transforms.Private(endpointName) + "_" +
		transforms.Private(paramName)
}

func getAssignmentForParam(param *types.EndpointArgumentDefinition) *jen.Statement {
	switch {
	case param.Type.IsString():

	}
	return nil
}
