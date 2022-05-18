// This file was generated by Conjure and should not be manually edited.

package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/codecs"
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

	testService_BinaryAlias_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_BinaryAlias_CmdRun,
		Short: "Calls the binaryAlias endpoint.",
		Use:   "binaryAlias",
	}
	rootCmd.AddCommand(testService_BinaryAlias_Cmd)
	testService_BinaryAlias_Cmd.Flags().String("body", "", "Required. ")

	testService_BinaryAliasOptional_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_BinaryAliasOptional_CmdRun,
		Short: "Calls the binaryAliasOptional endpoint.",
		Use:   "binaryAliasOptional",
	}
	rootCmd.AddCommand(testService_BinaryAliasOptional_Cmd)

	testService_BinaryAliasAlias_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_BinaryAliasAlias_CmdRun,
		Short: "Calls the binaryAliasAlias endpoint.",
		Use:   "binaryAliasAlias",
	}
	rootCmd.AddCommand(testService_BinaryAliasAlias_Cmd)
	testService_BinaryAliasAlias_Cmd.Flags().String("body", "", "Optional. ")

	testService_Binary_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_Binary_CmdRun,
		Short: "Calls the binary endpoint.",
		Use:   "binary",
	}
	rootCmd.AddCommand(testService_Binary_Cmd)
	testService_Binary_Cmd.Flags().String("body", "", "Required. ")

	testService_BinaryOptional_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_BinaryOptional_CmdRun,
		Short: "Calls the binaryOptional endpoint.",
		Use:   "binaryOptional",
	}
	rootCmd.AddCommand(testService_BinaryOptional_Cmd)

	testService_BinaryOptionalAlias_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_BinaryOptionalAlias_CmdRun,
		Short: "Calls the binaryOptionalAlias endpoint.",
		Use:   "binaryOptionalAlias",
	}
	rootCmd.AddCommand(testService_BinaryOptionalAlias_Cmd)
	testService_BinaryOptionalAlias_Cmd.Flags().String("body", "", "Optional. ")

	testService_BinaryList_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_BinaryList_CmdRun,
		Short: "Calls the binaryList endpoint.",
		Use:   "binaryList",
	}
	rootCmd.AddCommand(testService_BinaryList_Cmd)
	testService_BinaryList_Cmd.Flags().String("body", "", "Required. ")

	testService_Bytes_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_Bytes_CmdRun,
		Short: "Calls the bytes endpoint.",
		Use:   "bytes",
	}
	rootCmd.AddCommand(testService_Bytes_Cmd)
	testService_Bytes_Cmd.Flags().String("body", "", "Required. ")

	return rootCmd
}

func (c TestServiceCLICommand) testService_BinaryAlias_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	bodyRaw, err := flags.GetString("body")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument body")
	}
	if bodyRaw == "" {
		return werror.ErrorWithContextParams(ctx, "body is a required argument")
	}
	var bodyArg func() io.ReadCloser
	var bodyArgReader io.ReadCloser
	switch {
	case bodyRaw == "@-":
		bodyArgReader = io.NopCloser(cmd.InOrStdin())
	case strings.HasPrefix(bodyRaw, "@"):
		bodyArgReader, err = os.Open(strings.TrimSpace(bodyRaw[1:]))
		if err != nil {
			return werror.WrapWithContextParams(ctx, err, "failed to open file for argument body")
		}
	default:
		bodyArgReader = io.NopCloser(base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(bodyRaw))))
	}
	bodyArg = func() io.ReadCloser {
		return bodyArgReader
	}

	result, err := client.BinaryAlias(ctx, bodyArg)
	if err != nil {
		return err
	}
	_, err = io.Copy(cmd.OutOrStdout(), result)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to write result bytes to stdout")
	}
	return result.Close()
}

func (c TestServiceCLICommand) testService_BinaryAliasOptional_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	result, err := client.BinaryAliasOptional(ctx)
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

func (c TestServiceCLICommand) testService_BinaryAliasAlias_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	bodyRaw, err := flags.GetString("body")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument body")
	}
	var bodyArg func() io.ReadCloser
	if bodyRaw != "" {
		var bodyArgReader io.ReadCloser
		switch {
		case bodyRaw == "@-":
			bodyArgReader = io.NopCloser(cmd.InOrStdin())
		case strings.HasPrefix(bodyRaw, "@"):
			bodyArgReader, err = os.Open(strings.TrimSpace(bodyRaw[1:]))
			if err != nil {
				return werror.WrapWithContextParams(ctx, err, "failed to open file for argument body")
			}
		default:
			bodyArgReader = io.NopCloser(base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(bodyRaw))))
		}
		bodyArg = func() io.ReadCloser {
			return bodyArgReader
		}
	}

	result, err := client.BinaryAliasAlias(ctx, bodyArg)
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

func (c TestServiceCLICommand) testService_Binary_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	bodyRaw, err := flags.GetString("body")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument body")
	}
	if bodyRaw == "" {
		return werror.ErrorWithContextParams(ctx, "body is a required argument")
	}
	var bodyArg func() io.ReadCloser
	var bodyArgReader io.ReadCloser
	switch {
	case bodyRaw == "@-":
		bodyArgReader = io.NopCloser(cmd.InOrStdin())
	case strings.HasPrefix(bodyRaw, "@"):
		bodyArgReader, err = os.Open(strings.TrimSpace(bodyRaw[1:]))
		if err != nil {
			return werror.WrapWithContextParams(ctx, err, "failed to open file for argument body")
		}
	default:
		bodyArgReader = io.NopCloser(base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(bodyRaw))))
	}
	bodyArg = func() io.ReadCloser {
		return bodyArgReader
	}

	result, err := client.Binary(ctx, bodyArg)
	if err != nil {
		return err
	}
	_, err = io.Copy(cmd.OutOrStdout(), result)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to write result bytes to stdout")
	}
	return result.Close()
}

func (c TestServiceCLICommand) testService_BinaryOptional_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	result, err := client.BinaryOptional(ctx)
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

func (c TestServiceCLICommand) testService_BinaryOptionalAlias_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	bodyRaw, err := flags.GetString("body")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument body")
	}
	var bodyArg func() io.ReadCloser
	if bodyRaw != "" {
		var bodyArgReader io.ReadCloser
		switch {
		case bodyRaw == "@-":
			bodyArgReader = io.NopCloser(cmd.InOrStdin())
		case strings.HasPrefix(bodyRaw, "@"):
			bodyArgReader, err = os.Open(strings.TrimSpace(bodyRaw[1:]))
			if err != nil {
				return werror.WrapWithContextParams(ctx, err, "failed to open file for argument body")
			}
		default:
			bodyArgReader = io.NopCloser(base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(bodyRaw))))
		}
		bodyArg = func() io.ReadCloser {
			return bodyArgReader
		}
	}

	result, err := client.BinaryOptionalAlias(ctx, bodyArg)
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

func (c TestServiceCLICommand) testService_BinaryList_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	bodyRaw, err := flags.GetString("body")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument body")
	}
	if bodyRaw == "" {
		return werror.ErrorWithContextParams(ctx, "body is a required argument")
	}
	var bodyArg [][]byte
	var bodyArgReader io.ReadCloser
	switch {
	case bodyRaw == "@-":
		bodyArgReader = io.NopCloser(cmd.InOrStdin())
	case strings.HasPrefix(bodyRaw, "@"):
		bodyArgReader, err = os.Open(strings.TrimSpace(bodyRaw[1:]))
		if err != nil {
			return werror.WrapWithContextParams(ctx, err, "failed to open file for argument body")
		}
	default:
		bodyArgReader = io.NopCloser(bytes.NewReader([]byte(bodyRaw)))
	}
	defer bodyArgReader.Close()
	if err := codecs.JSON.Decode(bodyArgReader, &bodyArg); err != nil {
		return werror.WrapWithContextParams(ctx, err, "invalid value for body argument")
	}

	result, err := client.BinaryList(ctx, bodyArg)
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

func (c TestServiceCLICommand) testService_Bytes_CmdRun(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := getCLIContext(flags)
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	bodyRaw, err := flags.GetString("body")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument body")
	}
	if bodyRaw == "" {
		return werror.ErrorWithContextParams(ctx, "body is a required argument")
	}
	var bodyArg CustomObject
	var bodyArgReader io.ReadCloser
	switch {
	case bodyRaw == "@-":
		bodyArgReader = io.NopCloser(cmd.InOrStdin())
	case strings.HasPrefix(bodyRaw, "@"):
		bodyArgReader, err = os.Open(strings.TrimSpace(bodyRaw[1:]))
		if err != nil {
			return werror.WrapWithContextParams(ctx, err, "failed to open file for argument body")
		}
	default:
		bodyArgReader = io.NopCloser(bytes.NewReader([]byte(bodyRaw)))
	}
	defer bodyArgReader.Close()
	if err := codecs.JSON.Decode(bodyArgReader, &bodyArg); err != nil {
		return werror.WrapWithContextParams(ctx, err, "invalid value for body argument")
	}

	result, err := client.Bytes(ctx, bodyArg)
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
