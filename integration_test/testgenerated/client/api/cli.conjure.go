// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	"github.com/palantir/pkg/rid"
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wlogzap "github.com/palantir/witchcraft-go-logging/wlog-zap"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
	"github.com/palantir/witchcraft-go-tracing/wzipkin"
	"github.com/spf13/cobra"
	pflag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

type CLIConfig struct {
	Client httpclient.ClientConfig `yaml:",inline"`
}

// Commands for TestService

type CLITestServiceClientProvider interface {
	Get(ctx context.Context, flags *pflag.FlagSet) (TestServiceClient, error)
}

type defaultCLITestServiceClientProvider struct{}

func NewDefaultCLITestServiceClientProvider() CLITestServiceClientProvider {
	return defaultCLITestServiceClientProvider{}
}

func (d defaultCLITestServiceClientProvider) Get(ctx context.Context, flags *pflag.FlagSet) (TestServiceClient, error) {
	conf, err := loadCLIConfig(ctx, flags)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "failed to load CLI configuration file")
	}
	client, err := httpclient.NewClient(httpclient.WithConfig(conf.Client))
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "failed to create client with provided config")
	}
	return NewTestServiceClient(client), nil
}

type TestServiceCLICommand struct {
	clientProvider CLITestServiceClientProvider
}

func NewTestServiceCLICommand() *cobra.Command {
	return NewTestServiceCLICommandWithClientProvider(NewDefaultCLITestServiceClientProvider())
}

func NewTestServiceCLICommandWithClientProvider(clientProvider CLITestServiceClientProvider) *cobra.Command {
	rootCmd := &cobra.Command{
		Short: "Runs commands on the TestService",
		Use:   "testService",
	}
	rootCmd.PersistentFlags().String("conf", "var/conf/configuration.yml", "The configuration file is optional. The default path is ./var/conf/configuration.yml.")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enables verbose mode for debugging client connections.")

	cliCommand := TestServiceCLICommand{clientProvider: clientProvider}

	testService_Echo_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_Echo_CmdRun,
		Short: "Calls the echo endpoint.",
		Use:   "echo",
	}
	rootCmd.AddCommand(testService_Echo_Cmd)

	testService_PathParam_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_PathParam_CmdRun,
		Short: "Calls the pathParam endpoint.",
		Use:   "pathParam",
	}
	rootCmd.AddCommand(testService_PathParam_Cmd)
	testService_PathParam_Cmd.Flags().String("param", "", "Required. ")

	testService_PathParamAlias_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_PathParamAlias_CmdRun,
		Short: "Calls the pathParamAlias endpoint.",
		Use:   "pathParamAlias",
	}
	rootCmd.AddCommand(testService_PathParamAlias_Cmd)
	testService_PathParamAlias_Cmd.Flags().String("param", "", "Required. ")

	testService_PathParamRid_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_PathParamRid_CmdRun,
		Short: "Calls the pathParamRid endpoint.",
		Use:   "pathParamRid",
	}
	rootCmd.AddCommand(testService_PathParamRid_Cmd)
	testService_PathParamRid_Cmd.Flags().String("param", "", "Required. ")

	testService_PathParamRidAlias_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_PathParamRidAlias_CmdRun,
		Short: "Calls the pathParamRidAlias endpoint.",
		Use:   "pathParamRidAlias",
	}
	rootCmd.AddCommand(testService_PathParamRidAlias_Cmd)
	testService_PathParamRidAlias_Cmd.Flags().String("param", "", "Required. ")

	testService_Bytes_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_Bytes_CmdRun,
		Short: "Calls the bytes endpoint.",
		Use:   "bytes",
	}
	rootCmd.AddCommand(testService_Bytes_Cmd)

	testService_Binary_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_Binary_CmdRun,
		Short: "Calls the binary endpoint.",
		Use:   "binary",
	}
	rootCmd.AddCommand(testService_Binary_Cmd)

	testService_MaybeBinary_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_MaybeBinary_CmdRun,
		Short: "Calls the maybeBinary endpoint.",
		Use:   "maybeBinary",
	}
	rootCmd.AddCommand(testService_MaybeBinary_Cmd)

	testService_Query_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_Query_CmdRun,
		Short: "Calls the query endpoint.",
		Use:   "query",
	}
	rootCmd.AddCommand(testService_Query_Cmd)
	testService_Query_Cmd.Flags().String("query", "", "Optional. ")

	return rootCmd
}

func (c TestServiceCLICommand) testService_Echo_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	return client.Echo(ctx)
}

func (c TestServiceCLICommand) testService_PathParam_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	paramRaw, err := flags.GetString("param")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument param")
	}
	if paramRaw == "" {
		return werror.ErrorWithContextParams(ctx, "param is a required argument")
	}
	paramArg := paramRaw

	return client.PathParam(ctx, paramArg)
}

func (c TestServiceCLICommand) testService_PathParamAlias_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	paramRaw, err := flags.GetString("param")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument param")
	}
	if paramRaw == "" {
		return werror.ErrorWithContextParams(ctx, "param is a required argument")
	}
	paramArg := StringAlias(paramRaw)

	return client.PathParamAlias(ctx, paramArg)
}

func (c TestServiceCLICommand) testService_PathParamRid_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	paramRaw, err := flags.GetString("param")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument param")
	}
	if paramRaw == "" {
		return werror.ErrorWithContextParams(ctx, "param is a required argument")
	}
	paramArg, err := rid.ParseRID(paramRaw)
	if err != nil {
		return werror.WrapWithContextParams(ctx, errors.WrapWithInvalidArgument(err), "failed to parse \"param\" as rid")
	}

	return client.PathParamRid(ctx, paramArg)
}

func (c TestServiceCLICommand) testService_PathParamRidAlias_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	paramRaw, err := flags.GetString("param")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument param")
	}
	if paramRaw == "" {
		return werror.ErrorWithContextParams(ctx, "param is a required argument")
	}
	paramArgValue, err := rid.ParseRID(paramRaw)
	if err != nil {
		return werror.WrapWithContextParams(ctx, errors.WrapWithInvalidArgument(err), "failed to parse \"param\" as rid")
	}
	paramArg := RidAlias(paramArgValue)

	return client.PathParamRidAlias(ctx, paramArg)
}

func (c TestServiceCLICommand) testService_Bytes_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	result, err := client.Bytes(ctx)
	if err != nil {
		return err
	}
	resultBytes, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		fmt.Printf("Failed to marshal to json with err: %v\n\nPrinting as string:\n%v\n", err, result)
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%v\n", string(resultBytes))
	return nil
}

func (c TestServiceCLICommand) testService_Binary_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	result, err := client.Binary(ctx)
	if err != nil {
		return err
	}
	_, err = io.Copy(cmd.OutOrStdout(), result)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to write result bytes to stdout")
	}
	return result.Close()
}

func (c TestServiceCLICommand) testService_MaybeBinary_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	result, err := client.MaybeBinary(ctx)
	if err != nil {
		return err
	}
	if result == nil {
		return nil
	}
	resultDeref := *result
	_, err = io.Copy(cmd.OutOrStdout(), resultDeref)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to write result bytes to stdout")
	}
	return resultDeref.Close()
}

func (c TestServiceCLICommand) testService_Query_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	queryRaw, err := flags.GetString("query")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument query")
	}
	var queryArg *StringAlias
	if queryArgStr := queryRaw; queryArgStr != "" {
		queryArgInternal := StringAlias(queryArgStr)
		queryArg = &queryArgInternal
	}

	return client.Query(ctx, queryArg)
}

func loadCLIConfig(ctx context.Context, flags *pflag.FlagSet) (CLIConfig, error) {
	var emptyConfig CLIConfig
	configPath, err := flags.GetString("conf")
	if err != nil || configPath == "" {
		return emptyConfig, werror.WrapWithContextParams(ctx, err, "config file location must be specified")
	}
	confBytes, err := os.ReadFile(configPath)
	if err != nil {
		return emptyConfig, err
	}
	var conf CLIConfig
	err = yaml.Unmarshal(confBytes, &conf)
	if err != nil {
		return emptyConfig, err
	}
	return conf, nil
}

func getCLIContext(flags *pflag.FlagSet) context.Context {
	ctx := context.Background()
	logProvider := wlog.NewNoopLoggerProvider()
	logWriter := io.Discard
	verbose, err := flags.GetBool("verbose")
	if verbose && err == nil {
		logProvider = wlogzap.LoggerProvider()
		logWriter = os.Stdout
	}
	wlog.SetDefaultLoggerProvider(logProvider)
	ctx = svc1log.WithLogger(ctx, svc1log.New(logWriter, wlog.DebugLevel))
	traceLogger := trc1log.New(logWriter)
	ctx = trc1log.WithLogger(ctx, traceLogger)
	ctx = evt2log.WithLogger(ctx, evt2log.New(logWriter))
	tracer, err := wzipkin.NewTracer(traceLogger)
	if err != nil {
		return ctx
	}
	return wtracing.ContextWithTracer(ctx, tracer)
}
