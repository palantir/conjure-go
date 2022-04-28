// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
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
	rootCmd.PersistentFlags().String("conf", "../var/conf/configuration.yml", "The configuration file is optional. The default path is ./var/conf/configuration.yml.")

	cliCommand := TestServiceCLICommand{clientProvider: clientProvider}

	testService_Echo_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_Echo_CmdRun,
		Short: "Calls the echo endpoint.",
		Use:   "echo",
	}
	rootCmd.AddCommand(testService_Echo_Cmd)
	testService_Echo_Cmd.Flags().String("input", "", "Required. ")

	return rootCmd
}

func (c TestServiceCLICommand) testService_Echo_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	inputRaw, err := flags.GetString("input")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument input")
	}
	if inputRaw == "" {
		return werror.ErrorWithContextParams(ctx, "input is a required argument")
	}
	inputArg := inputRaw

	result, err := client.Echo(ctx, inputArg)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%v\n", result)
	return nil
}

func loadCLIConfig(ctx context.Context, flags *pflag.FlagSet) (CLIConfig, error) {
	var emptyConfig CLIConfig
	configPath, err := flags.GetString("conf")
	if err != nil || configPath == "" {
		return emptyConfig, werror.WrapWithContextParams(ctx, err, "config file location must be specified")
	}
	confBytes, err := ioutil.ReadFile(configPath)
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
