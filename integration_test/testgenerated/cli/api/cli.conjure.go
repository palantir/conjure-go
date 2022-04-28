// This file was generated by Conjure and should not be manually edited.

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/codecs"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	"github.com/palantir/pkg/bearertoken"
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/pkg/uuid"
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
	rootCmd        *cobra.Command
}

func NewTestServiceCLICommand() TestServiceCLICommand {
	return NewTestServiceCLICommandWithClientProvider(NewDefaultCLITestServiceClientProvider())
}

func NewTestServiceCLICommandWithClientProvider(clientProvider CLITestServiceClientProvider) TestServiceCLICommand {
	rootCmd := &cobra.Command{
		Short: "Runs commands on the TestService",
		Use:   "testService",
	}
	rootCmd.PersistentFlags().String("conf", "../var/conf/configuration.yml", "The configuration file is optional. The default path is ./var/conf/configuration.yml.")

	cliCommand := TestServiceCLICommand{
		clientProvider: clientProvider,
		rootCmd:        rootCmd,
	}

	testService_Echo_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_Echo_CmdRun,
		Short: "Calls the echo endpoint",
		Use:   "echo",
	}
	rootCmd.AddCommand(testService_Echo_Cmd)
	testService_Echo_Cmd.Flags().String("bearer_token", "", "bearer_token is a required field.")

	testService_EchoStrings_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_EchoStrings_CmdRun,
		Short: "Calls the echoStrings endpoint",
		Use:   "echoStrings",
	}
	rootCmd.AddCommand(testService_EchoStrings_Cmd)
	testService_EchoStrings_Cmd.Flags().String("body", "", "body is a required param.")

	testService_EchoCustomObject_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_EchoCustomObject_CmdRun,
		Short: "Calls the echoCustomObject endpoint",
		Use:   "echoCustomObject",
	}
	rootCmd.AddCommand(testService_EchoCustomObject_Cmd)
	testService_EchoCustomObject_Cmd.Flags().String("body", "", "body is an optional param.")

	testService_EchoOptionalAlias_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_EchoOptionalAlias_CmdRun,
		Short: "Calls the echoOptionalAlias endpoint",
		Use:   "echoOptionalAlias",
	}
	rootCmd.AddCommand(testService_EchoOptionalAlias_Cmd)
	testService_EchoOptionalAlias_Cmd.Flags().String("body", "", "body is an optional param.")

	testService_EchoOptionalListAlias_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_EchoOptionalListAlias_CmdRun,
		Short: "Calls the echoOptionalListAlias endpoint",
		Use:   "echoOptionalListAlias",
	}
	rootCmd.AddCommand(testService_EchoOptionalListAlias_Cmd)
	testService_EchoOptionalListAlias_Cmd.Flags().String("body", "", "body is an optional param.")

	testService_GetPathParam_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_GetPathParam_CmdRun,
		Short: "Calls the getPathParam endpoint",
		Use:   "getPathParam",
	}
	rootCmd.AddCommand(testService_GetPathParam_Cmd)
	testService_GetPathParam_Cmd.Flags().String("myPathParam", "", "myPathParam is a required param.")
	testService_GetPathParam_Cmd.Flags().String("bearer_token", "", "bearer_token is a required field.")

	testService_GetListBoolean_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_GetListBoolean_CmdRun,
		Short: "Calls the getListBoolean endpoint",
		Use:   "getListBoolean",
	}
	rootCmd.AddCommand(testService_GetListBoolean_Cmd)
	testService_GetListBoolean_Cmd.Flags().String("myQueryParam1", "", "myQueryParam1 is a required param.")

	testService_PutMapStringString_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_PutMapStringString_CmdRun,
		Short: "Calls the putMapStringString endpoint",
		Use:   "putMapStringString",
	}
	rootCmd.AddCommand(testService_PutMapStringString_Cmd)
	testService_PutMapStringString_Cmd.Flags().String("myParam", "", "myParam is a required param.")

	testService_PutMapStringAny_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_PutMapStringAny_CmdRun,
		Short: "Calls the putMapStringAny endpoint",
		Use:   "putMapStringAny",
	}
	rootCmd.AddCommand(testService_PutMapStringAny_Cmd)
	testService_PutMapStringAny_Cmd.Flags().String("myParam", "", "myParam is a required param.")

	testService_GetDateTime_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_GetDateTime_CmdRun,
		Short: "Calls the getDateTime endpoint",
		Use:   "getDateTime",
	}
	rootCmd.AddCommand(testService_GetDateTime_Cmd)
	testService_GetDateTime_Cmd.Flags().String("myParam", "", "myParam is a required param.")

	testService_GetDouble_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_GetDouble_CmdRun,
		Short: "Calls the getDouble endpoint",
		Use:   "getDouble",
	}
	rootCmd.AddCommand(testService_GetDouble_Cmd)
	testService_GetDouble_Cmd.Flags().String("myParam", "", "myParam is a required param.")

	testService_GetRid_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_GetRid_CmdRun,
		Short: "Calls the getRid endpoint",
		Use:   "getRid",
	}
	rootCmd.AddCommand(testService_GetRid_Cmd)
	testService_GetRid_Cmd.Flags().String("myParam", "", "myParam is a required param.")

	testService_GetSafeLong_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_GetSafeLong_CmdRun,
		Short: "Calls the getSafeLong endpoint",
		Use:   "getSafeLong",
	}
	rootCmd.AddCommand(testService_GetSafeLong_Cmd)
	testService_GetSafeLong_Cmd.Flags().String("myParam", "", "myParam is a required param.")

	testService_GetUuid_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_GetUuid_CmdRun,
		Short: "Calls the getUuid endpoint",
		Use:   "getUuid",
	}
	rootCmd.AddCommand(testService_GetUuid_Cmd)
	testService_GetUuid_Cmd.Flags().String("myParam", "", "myParam is a required param.")

	testService_GetBinary_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_GetBinary_CmdRun,
		Short: "Calls the getBinary endpoint",
		Use:   "getBinary",
	}
	rootCmd.AddCommand(testService_GetBinary_Cmd)

	testService_GetOptionalBinary_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_GetOptionalBinary_CmdRun,
		Short: "Calls the getOptionalBinary endpoint",
		Use:   "getOptionalBinary",
	}
	rootCmd.AddCommand(testService_GetOptionalBinary_Cmd)

	testService_GetReserved_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_GetReserved_CmdRun,
		Short: "Calls the getReserved endpoint",
		Use:   "getReserved",
	}
	rootCmd.AddCommand(testService_GetReserved_Cmd)
	testService_GetReserved_Cmd.Flags().String("conf_Arg", "", "conf is a required param.")
	testService_GetReserved_Cmd.Flags().String("bearertoken", "", "bearertoken is a required param.")

	testService_Chan_Cmd := &cobra.Command{
		RunE:  cliCommand.testService_Chan_CmdRun,
		Short: "Calls the chan endpoint",
		Use:   "chan",
	}
	rootCmd.AddCommand(testService_Chan_Cmd)
	testService_Chan_Cmd.Flags().String("var", "", "var is a required param.")
	testService_Chan_Cmd.Flags().String("import", "", "import is a required param.")
	testService_Chan_Cmd.Flags().String("type", "", "type is a required param.")
	testService_Chan_Cmd.Flags().String("return", "", "return is a required param.")
	testService_Chan_Cmd.Flags().String("http", "", "http is a required param.")
	testService_Chan_Cmd.Flags().String("json", "", "json is a required param.")
	testService_Chan_Cmd.Flags().String("req", "", "req is a required param.")
	testService_Chan_Cmd.Flags().String("rw", "", "rw is a required param.")

	return cliCommand
}

func (c TestServiceCLICommand) Command() *cobra.Command {
	return c.rootCmd
}

func (c TestServiceCLICommand) testService_Echo_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	bearer_tokenRaw, err := flags.GetString("bearer_token")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument __authVar")
	}
	if bearer_tokenRaw == "" {
		return werror.ErrorWithContextParams(ctx, "bearer_token is a required argument")
	}
	__authVarArg := bearertoken.Token(bearer_tokenRaw)
	return client.Echo(ctx, __authVarArg)
}

func (c TestServiceCLICommand) testService_EchoStrings_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
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
	var bodyArg []string
	bodyArgBytes := []byte(bodyRaw)
	if err := codecs.JSON.Decode(bytes.NewReader(bodyArgBytes), &bodyArg); err != nil {
		return werror.WrapWithContextParams(ctx, err, "invalid value for body argument")
	}

	result, err := client.EchoStrings(ctx, bodyArg)
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

func (c TestServiceCLICommand) testService_EchoCustomObject_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	bodyRaw, err := flags.GetString("body")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument body")
	}
	var bodyArg *CustomObject
	if bodyRaw != "" {
		bodyArgBytes := []byte(bodyRaw)
		if err := codecs.JSON.Decode(bytes.NewReader(bodyArgBytes), &bodyArg); err != nil {
			return werror.WrapWithContextParams(ctx, err, "invalid value for body argument")
		}
	}

	result, err := client.EchoCustomObject(ctx, bodyArg)
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

func (c TestServiceCLICommand) testService_EchoOptionalAlias_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	bodyRaw, err := flags.GetString("body")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument body")
	}
	var bodyArgValue *int
	if bodyArgValueStr1 := bodyRaw; bodyArgValueStr1 != "" {
		bodyArgValueInternal1, err := strconv.Atoi(bodyArgValueStr1)
		if err != nil {
			return werror.WrapWithContextParams(ctx, errors.WrapWithInvalidArgument(err), "failed to parse \"body\" as integer")
		}
		bodyArgValue = &bodyArgValueInternal1
	}
	bodyArg := OptionalIntegerAlias{Value: bodyArgValue}

	result, err := client.EchoOptionalAlias(ctx, bodyArg)
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

func (c TestServiceCLICommand) testService_EchoOptionalListAlias_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	bodyRaw, err := flags.GetString("body")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument body")
	}
	var bodyArg OptionalListAlias
	if bodyRaw != "" {
		bodyArgBytes := []byte(bodyRaw)
		if err := codecs.JSON.Decode(bytes.NewReader(bodyArgBytes), &bodyArg); err != nil {
			return werror.WrapWithContextParams(ctx, err, "invalid value for body argument")
		}
	}

	result, err := client.EchoOptionalListAlias(ctx, bodyArg)
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

func (c TestServiceCLICommand) testService_GetPathParam_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	bearer_tokenRaw, err := flags.GetString("bearer_token")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument __authVar")
	}
	if bearer_tokenRaw == "" {
		return werror.ErrorWithContextParams(ctx, "bearer_token is a required argument")
	}
	__authVarArg := bearertoken.Token(bearer_tokenRaw)
	myPathParamRaw, err := flags.GetString("myPathParam")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument myPathParam")
	}
	if myPathParamRaw == "" {
		return werror.ErrorWithContextParams(ctx, "myPathParam is a required argument")
	}
	myPathParamArg := myPathParamRaw

	return client.GetPathParam(ctx, __authVarArg, myPathParamArg)
}

func (c TestServiceCLICommand) testService_GetListBoolean_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	myQueryParam1Raw, err := flags.GetString("myQueryParam1")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument myQueryParam1")
	}
	if myQueryParam1Raw == "" {
		return werror.ErrorWithContextParams(ctx, "myQueryParam1 is a required argument")
	}
	var myQueryParam1Arg []bool
	myQueryParam1ArgBytes := []byte(myQueryParam1Raw)
	if err := codecs.JSON.Decode(bytes.NewReader(myQueryParam1ArgBytes), &myQueryParam1Arg); err != nil {
		return werror.WrapWithContextParams(ctx, err, "invalid value for myQueryParam1 argument")
	}

	result, err := client.GetListBoolean(ctx, myQueryParam1Arg)
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

func (c TestServiceCLICommand) testService_PutMapStringString_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	myParamRaw, err := flags.GetString("myParam")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument myParam")
	}
	if myParamRaw == "" {
		return werror.ErrorWithContextParams(ctx, "myParam is a required argument")
	}
	var myParamArg map[string]string
	myParamArgBytes := []byte(myParamRaw)
	if err := codecs.JSON.Decode(bytes.NewReader(myParamArgBytes), &myParamArg); err != nil {
		return werror.WrapWithContextParams(ctx, err, "invalid value for myParam argument")
	}

	result, err := client.PutMapStringString(ctx, myParamArg)
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

func (c TestServiceCLICommand) testService_PutMapStringAny_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	myParamRaw, err := flags.GetString("myParam")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument myParam")
	}
	if myParamRaw == "" {
		return werror.ErrorWithContextParams(ctx, "myParam is a required argument")
	}
	var myParamArg map[string]interface{}
	myParamArgBytes := []byte(myParamRaw)
	if err := codecs.JSON.Decode(bytes.NewReader(myParamArgBytes), &myParamArg); err != nil {
		return werror.WrapWithContextParams(ctx, err, "invalid value for myParam argument")
	}

	result, err := client.PutMapStringAny(ctx, myParamArg)
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

func (c TestServiceCLICommand) testService_GetDateTime_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	myParamRaw, err := flags.GetString("myParam")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument myParam")
	}
	if myParamRaw == "" {
		return werror.ErrorWithContextParams(ctx, "myParam is a required argument")
	}
	myParamArg, err := datetime.ParseDateTime(myParamRaw)
	if err != nil {
		return werror.WrapWithContextParams(ctx, errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as datetime")
	}

	result, err := client.GetDateTime(ctx, myParamArg)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%v\n", result)
	return nil
}

func (c TestServiceCLICommand) testService_GetDouble_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	myParamRaw, err := flags.GetString("myParam")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument myParam")
	}
	if myParamRaw == "" {
		return werror.ErrorWithContextParams(ctx, "myParam is a required argument")
	}
	myParamArg, err := strconv.ParseFloat(myParamRaw, 64)
	if err != nil {
		return werror.WrapWithContextParams(ctx, errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as double")
	}

	result, err := client.GetDouble(ctx, myParamArg)
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

func (c TestServiceCLICommand) testService_GetRid_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	myParamRaw, err := flags.GetString("myParam")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument myParam")
	}
	if myParamRaw == "" {
		return werror.ErrorWithContextParams(ctx, "myParam is a required argument")
	}
	myParamArg, err := rid.ParseRID(myParamRaw)
	if err != nil {
		return werror.WrapWithContextParams(ctx, errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as rid")
	}

	result, err := client.GetRid(ctx, myParamArg)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%v\n", result)
	return nil
}

func (c TestServiceCLICommand) testService_GetSafeLong_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	myParamRaw, err := flags.GetString("myParam")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument myParam")
	}
	if myParamRaw == "" {
		return werror.ErrorWithContextParams(ctx, "myParam is a required argument")
	}
	myParamArg, err := safelong.ParseSafeLong(myParamRaw)
	if err != nil {
		return werror.WrapWithContextParams(ctx, errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as safelong")
	}

	result, err := client.GetSafeLong(ctx, myParamArg)
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

func (c TestServiceCLICommand) testService_GetUuid_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	myParamRaw, err := flags.GetString("myParam")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument myParam")
	}
	if myParamRaw == "" {
		return werror.ErrorWithContextParams(ctx, "myParam is a required argument")
	}
	myParamArg, err := uuid.ParseUUID(myParamRaw)
	if err != nil {
		return werror.WrapWithContextParams(ctx, errors.WrapWithInvalidArgument(err), "failed to parse \"myParam\" as uuid")
	}

	result, err := client.GetUuid(ctx, myParamArg)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%v\n", result)
	return nil
}

func (c TestServiceCLICommand) testService_GetBinary_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	result, err := client.GetBinary(ctx)
	if err != nil {
		return err
	}
	_, err = io.Copy(cmd.OutOrStdout(), result)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to write result bytes to stdout")
	}
	return result.Close()
}

func (c TestServiceCLICommand) testService_GetOptionalBinary_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	result, err := client.GetOptionalBinary(ctx)
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

func (c TestServiceCLICommand) testService_GetReserved_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	confRaw, err := flags.GetString("conf")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument conf")
	}
	if confRaw == "" {
		return werror.ErrorWithContextParams(ctx, "conf is a required argument")
	}
	confArg := confRaw

	bearertokenRaw, err := flags.GetString("bearertoken")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument bearertoken")
	}
	if bearertokenRaw == "" {
		return werror.ErrorWithContextParams(ctx, "bearertoken is a required argument")
	}
	bearertokenArg := bearertokenRaw

	return client.GetReserved(ctx, confArg, bearertokenArg)
}

func (c TestServiceCLICommand) testService_Chan_CmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := c.clientProvider.Get(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	varRaw, err := flags.GetString("var")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument var")
	}
	if varRaw == "" {
		return werror.ErrorWithContextParams(ctx, "var is a required argument")
	}
	varArg := varRaw

	importRaw, err := flags.GetString("import")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument import")
	}
	if importRaw == "" {
		return werror.ErrorWithContextParams(ctx, "import is a required argument")
	}
	var importArg map[string]string
	importArgBytes := []byte(importRaw)
	if err := codecs.JSON.Decode(bytes.NewReader(importArgBytes), &importArg); err != nil {
		return werror.WrapWithContextParams(ctx, err, "invalid value for import argument")
	}

	typeRaw, err := flags.GetString("type")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument type")
	}
	if typeRaw == "" {
		return werror.ErrorWithContextParams(ctx, "type is a required argument")
	}
	typeArg := typeRaw

	returnRaw, err := flags.GetString("return")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument return")
	}
	if returnRaw == "" {
		return werror.ErrorWithContextParams(ctx, "return is a required argument")
	}
	returnArg, err := safelong.ParseSafeLong(returnRaw)
	if err != nil {
		return werror.WrapWithContextParams(ctx, errors.WrapWithInvalidArgument(err), "failed to parse \"return\" as safelong")
	}

	httpRaw, err := flags.GetString("http")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument http")
	}
	if httpRaw == "" {
		return werror.ErrorWithContextParams(ctx, "http is a required argument")
	}
	httpArg := httpRaw

	jsonRaw, err := flags.GetString("json")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument json")
	}
	if jsonRaw == "" {
		return werror.ErrorWithContextParams(ctx, "json is a required argument")
	}
	jsonArg := jsonRaw

	reqRaw, err := flags.GetString("req")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument req")
	}
	if reqRaw == "" {
		return werror.ErrorWithContextParams(ctx, "req is a required argument")
	}
	reqArg := reqRaw

	rwRaw, err := flags.GetString("rw")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument rw")
	}
	if rwRaw == "" {
		return werror.ErrorWithContextParams(ctx, "rw is a required argument")
	}
	rwArg := rwRaw

	return client.Chan(ctx, varArg, importArg, typeArg, returnArg, httpArg, jsonArg, reqArg, rwArg)
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
