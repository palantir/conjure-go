// This file was generated by Conjure and should not be manually edited.

package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/codecs"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wlogzap "github.com/palantir/witchcraft-go-logging/wlog-zap"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
	"github.com/palantir/witchcraft-go-tracing/wzipkin"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	testService_Echo_Cmd.Flags().String("input", "", "Required. ")
	testService_Echo_Cmd.Flags().String("reps", "", "Required. ")
	testService_Echo_Cmd.Flags().String("optional", "", "Optional. ")
	testService_Echo_Cmd.Flags().String("listParam", "", "Required. ")
	testService_Echo_Cmd.Flags().String("lastParam", "", "Optional. ")

	return rootCmd
}

func (c TestServiceCLICommand) testService_Echo_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
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

	repsRaw, err := flags.GetString("reps")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument reps")
	}
	if repsRaw == "" {
		return werror.ErrorWithContextParams(ctx, "reps is a required argument")
	}
	repsArg, err := strconv.Atoi(repsRaw)
	if err != nil {
		return werror.WrapWithContextParams(ctx, errors.WrapWithInvalidArgument(err), "failed to parse \"reps\" as integer")
	}

	optionalRaw, err := flags.GetString("optional")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument optional")
	}
	var optionalArg *string
	if optionalArgStr := optionalRaw; optionalArgStr != "" {
		optionalArgInternal := optionalArgStr
		optionalArg = &optionalArgInternal
	}

	listParamRaw, err := flags.GetString("listParam")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument listParam")
	}
	if listParamRaw == "" {
		return werror.ErrorWithContextParams(ctx, "listParam is a required argument")
	}
	var listParamArg []int
	var listParamArgReader io.ReadCloser
	switch {
	case listParamRaw == "@-":
		listParamArgReader = io.NopCloser(cmd.InOrStdin())
	case strings.HasPrefix(listParamRaw, "@"):
		listParamArgReader, err = os.Open(strings.TrimSpace(listParamRaw[1:]))
		if err != nil {
			return werror.WrapWithContextParams(ctx, err, "failed to open file for argument listParam")
		}
	default:
		listParamArgReader = io.NopCloser(bytes.NewReader([]byte(listParamRaw)))
	}
	defer listParamArgReader.Close()
	if err := codecs.JSON.Decode(listParamArgReader, &listParamArg); err != nil {
		return werror.WrapWithContextParams(ctx, err, "invalid value for listParam argument")
	}

	lastParamRaw, err := flags.GetString("lastParam")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument lastParam")
	}
	var lastParamArg *string
	if lastParamArgStr := lastParamRaw; lastParamArgStr != "" {
		lastParamArgInternal := lastParamArgStr
		lastParamArg = &lastParamArgInternal
	}

	result, err := client.Echo(ctx, inputArg, repsArg, optionalArg, listParamArg, lastParamArg)
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
