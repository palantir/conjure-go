// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/pkg/bearertoken"
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
	Client httpclient.ClientConfig
}

// Commands for BothAuthService

var RootBothAuthServiceCmd = &cobra.Command{
	Short: "Runs commands on the BothAuthService",
	Use:   "bothAuthService",
}

func getBothAuthServiceClient(ctx context.Context, flags *pflag.FlagSet) (BothAuthServiceClient, error) {
	conf, err := loadConfig(ctx, flags)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "failed to load CLI configuration file")
	}
	client, err := httpclient.NewClient(httpclient.WithConfig(conf.Client))
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "failed to create client with provided config")
	}
	return NewBothAuthServiceClient(client), nil
}

var BothAuthServicedefaultCmd = &cobra.Command{
	RunE:  bothAuthServicedefaultCmdRun,
	Short: "Calls the default endpoint",
	Use:   "default",
}

func bothAuthServicedefaultCmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := getBothAuthServiceClient(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	return bothAuthServicedefaultCmdRunInternal(ctx, flags, client)
}

func bothAuthServicedefaultCmdRunInternal(ctx context.Context, flags *pflag.FlagSet, client BothAuthServiceClient) error {
	var err error

	bearer_tokenRaw, err := flags.GetString("bearer_token")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument __authVar")
	}
	if bearer_tokenRaw == "" {
		return werror.ErrorWithContextParams(ctx, "__authVarArg is a required argument")
	}
	__authVarArg := bearertoken.Token(bearer_tokenRaw)
	result, err := client.Default(ctx, __authVarArg)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", result)
	return nil
}

var BothAuthServicecookieCmd = &cobra.Command{
	RunE:  bothAuthServicecookieCmdRun,
	Short: "Calls the cookie endpoint",
	Use:   "cookie",
}

func bothAuthServicecookieCmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := getBothAuthServiceClient(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	return bothAuthServicecookieCmdRunInternal(ctx, flags, client)
}

func bothAuthServicecookieCmdRunInternal(ctx context.Context, flags *pflag.FlagSet, client BothAuthServiceClient) error {
	var err error

	bearer_tokenRaw, err := flags.GetString("bearer_token")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument __authVar")
	}
	if bearer_tokenRaw == "" {
		return werror.ErrorWithContextParams(ctx, "__authVarArg is a required argument")
	}
	__authVarArg := bearertoken.Token(bearer_tokenRaw)
	err = client.Cookie(ctx, __authVarArg)
	return err
}

var BothAuthServicenoneCmd = &cobra.Command{
	RunE:  bothAuthServicenoneCmdRun,
	Short: "Calls the none endpoint",
	Use:   "none",
}

func bothAuthServicenoneCmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := getBothAuthServiceClient(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	return bothAuthServicenoneCmdRunInternal(ctx, flags, client)
}

func bothAuthServicenoneCmdRunInternal(ctx context.Context, flags *pflag.FlagSet, client BothAuthServiceClient) error {
	var err error

	err = client.None(ctx)
	return err
}

var BothAuthServicewithArgCmd = &cobra.Command{
	RunE:  bothAuthServicewithArgCmdRun,
	Short: "Calls the withArg endpoint",
	Use:   "withArg",
}

func bothAuthServicewithArgCmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := getBothAuthServiceClient(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	return bothAuthServicewithArgCmdRunInternal(ctx, flags, client)
}

func bothAuthServicewithArgCmdRunInternal(ctx context.Context, flags *pflag.FlagSet, client BothAuthServiceClient) error {
	var err error

	bearer_tokenRaw, err := flags.GetString("bearer_token")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument __authVar")
	}
	if bearer_tokenRaw == "" {
		return werror.ErrorWithContextParams(ctx, "__authVarArg is a required argument")
	}
	__authVarArg := bearertoken.Token(bearer_tokenRaw)
	argRaw, err := flags.GetString("arg")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument arg")
	}
	if argRaw == "" {
		return werror.ErrorWithContextParams(ctx, "argArg is a required argument")
	}
	argArg := argRaw

	err = client.WithArg(ctx, __authVarArg, argArg)
	return err
}

// Commands for CookieAuthService

var RootCookieAuthServiceCmd = &cobra.Command{
	Short: "Runs commands on the CookieAuthService",
	Use:   "cookieAuthService",
}

func getCookieAuthServiceClient(ctx context.Context, flags *pflag.FlagSet) (CookieAuthServiceClient, error) {
	conf, err := loadConfig(ctx, flags)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "failed to load CLI configuration file")
	}
	client, err := httpclient.NewClient(httpclient.WithConfig(conf.Client))
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "failed to create client with provided config")
	}
	return NewCookieAuthServiceClient(client), nil
}

var CookieAuthServicecookieCmd = &cobra.Command{
	RunE:  cookieAuthServicecookieCmdRun,
	Short: "Calls the cookie endpoint",
	Use:   "cookie",
}

func cookieAuthServicecookieCmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := getCookieAuthServiceClient(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	return cookieAuthServicecookieCmdRunInternal(ctx, flags, client)
}

func cookieAuthServicecookieCmdRunInternal(ctx context.Context, flags *pflag.FlagSet, client CookieAuthServiceClient) error {
	var err error

	bearer_tokenRaw, err := flags.GetString("bearer_token")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument __authVar")
	}
	if bearer_tokenRaw == "" {
		return werror.ErrorWithContextParams(ctx, "__authVarArg is a required argument")
	}
	__authVarArg := bearertoken.Token(bearer_tokenRaw)
	err = client.Cookie(ctx, __authVarArg)
	return err
}

// Commands for HeaderAuthService

var RootHeaderAuthServiceCmd = &cobra.Command{
	Short: "Runs commands on the HeaderAuthService",
	Use:   "headerAuthService",
}

func getHeaderAuthServiceClient(ctx context.Context, flags *pflag.FlagSet) (HeaderAuthServiceClient, error) {
	conf, err := loadConfig(ctx, flags)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "failed to load CLI configuration file")
	}
	client, err := httpclient.NewClient(httpclient.WithConfig(conf.Client))
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "failed to create client with provided config")
	}
	return NewHeaderAuthServiceClient(client), nil
}

var HeaderAuthServicedefaultCmd = &cobra.Command{
	RunE:  headerAuthServicedefaultCmdRun,
	Short: "Calls the default endpoint",
	Use:   "default",
}

func headerAuthServicedefaultCmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := getHeaderAuthServiceClient(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	return headerAuthServicedefaultCmdRunInternal(ctx, flags, client)
}

func headerAuthServicedefaultCmdRunInternal(ctx context.Context, flags *pflag.FlagSet, client HeaderAuthServiceClient) error {
	var err error

	bearer_tokenRaw, err := flags.GetString("bearer_token")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument __authVar")
	}
	if bearer_tokenRaw == "" {
		return werror.ErrorWithContextParams(ctx, "__authVarArg is a required argument")
	}
	__authVarArg := bearertoken.Token(bearer_tokenRaw)
	result, err := client.Default(ctx, __authVarArg)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", result)
	return nil
}

var HeaderAuthServicebinaryCmd = &cobra.Command{
	RunE:  headerAuthServicebinaryCmdRun,
	Short: "Calls the binary endpoint",
	Use:   "binary",
}

func headerAuthServicebinaryCmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := getHeaderAuthServiceClient(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	return headerAuthServicebinaryCmdRunInternal(ctx, flags, client)
}

func headerAuthServicebinaryCmdRunInternal(ctx context.Context, flags *pflag.FlagSet, client HeaderAuthServiceClient) error {
	var err error

	bearer_tokenRaw, err := flags.GetString("bearer_token")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument __authVar")
	}
	if bearer_tokenRaw == "" {
		return werror.ErrorWithContextParams(ctx, "__authVarArg is a required argument")
	}
	__authVarArg := bearertoken.Token(bearer_tokenRaw)
	result, err := client.Binary(ctx, __authVarArg)
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, result)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to write result bytes to stdout")
	}
	return result.Close()
}

var HeaderAuthServicebinaryOptionalCmd = &cobra.Command{
	RunE:  headerAuthServicebinaryOptionalCmdRun,
	Short: "Calls the binaryOptional endpoint",
	Use:   "binaryOptional",
}

func headerAuthServicebinaryOptionalCmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := getHeaderAuthServiceClient(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	return headerAuthServicebinaryOptionalCmdRunInternal(ctx, flags, client)
}

func headerAuthServicebinaryOptionalCmdRunInternal(ctx context.Context, flags *pflag.FlagSet, client HeaderAuthServiceClient) error {
	var err error

	bearer_tokenRaw, err := flags.GetString("bearer_token")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument __authVar")
	}
	if bearer_tokenRaw == "" {
		return werror.ErrorWithContextParams(ctx, "__authVarArg is a required argument")
	}
	__authVarArg := bearertoken.Token(bearer_tokenRaw)
	result, err := client.BinaryOptional(ctx, __authVarArg)
	if err != nil {
		return err
	}
	if result == nil {
		return nil
	}
	resultDeref := *result
	_, err = io.Copy(os.Stdout, resultDeref)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to write result bytes to stdout")
	}
	return resultDeref.Close()
}

// Commands for SomeHeaderAuthService

var RootSomeHeaderAuthServiceCmd = &cobra.Command{
	Short: "Runs commands on the SomeHeaderAuthService",
	Use:   "someHeaderAuthService",
}

func getSomeHeaderAuthServiceClient(ctx context.Context, flags *pflag.FlagSet) (SomeHeaderAuthServiceClient, error) {
	conf, err := loadConfig(ctx, flags)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "failed to load CLI configuration file")
	}
	client, err := httpclient.NewClient(httpclient.WithConfig(conf.Client))
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "failed to create client with provided config")
	}
	return NewSomeHeaderAuthServiceClient(client), nil
}

var SomeHeaderAuthServicedefaultCmd = &cobra.Command{
	RunE:  someHeaderAuthServicedefaultCmdRun,
	Short: "Calls the default endpoint",
	Use:   "default",
}

func someHeaderAuthServicedefaultCmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := getSomeHeaderAuthServiceClient(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	return someHeaderAuthServicedefaultCmdRunInternal(ctx, flags, client)
}

func someHeaderAuthServicedefaultCmdRunInternal(ctx context.Context, flags *pflag.FlagSet, client SomeHeaderAuthServiceClient) error {
	var err error

	bearer_tokenRaw, err := flags.GetString("bearer_token")
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to parse argument __authVar")
	}
	if bearer_tokenRaw == "" {
		return werror.ErrorWithContextParams(ctx, "__authVarArg is a required argument")
	}
	__authVarArg := bearertoken.Token(bearer_tokenRaw)
	result, err := client.Default(ctx, __authVarArg)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", result)
	return nil
}

var SomeHeaderAuthServicenoneCmd = &cobra.Command{
	RunE:  someHeaderAuthServicenoneCmdRun,
	Short: "Calls the none endpoint",
	Use:   "none",
}

func someHeaderAuthServicenoneCmdRun(cmd *cobra.Command, _ []string) error {
	ctx := getCLIContext()
	flags := cmd.Flags()
	client, err := getSomeHeaderAuthServiceClient(ctx, flags)
	if err != nil {
		return werror.WrapWithContextParams(ctx, err, "failed to initialize client")
	}
	return someHeaderAuthServicenoneCmdRunInternal(ctx, flags, client)
}

func someHeaderAuthServicenoneCmdRunInternal(ctx context.Context, flags *pflag.FlagSet, client SomeHeaderAuthServiceClient) error {
	var err error

	err = client.None(ctx)
	return err
}

func loadConfig(ctx context.Context, flags *pflag.FlagSet) (CLIConfig, error) {
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

func RegisterCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(RootBothAuthServiceCmd)
	rootCmd.AddCommand(RootCookieAuthServiceCmd)
	rootCmd.AddCommand(RootHeaderAuthServiceCmd)
	rootCmd.AddCommand(RootSomeHeaderAuthServiceCmd)
}

func init() {
	// BothAuthService commands and flags
	RootBothAuthServiceCmd.PersistentFlags().String("conf", "../var/conf/configuration.yml", "The configuration file is optional. The default path is ./var/conf/configuration.yml.")
	RootBothAuthServiceCmd.AddCommand(BothAuthServicedefaultCmd)
	BothAuthServicedefaultCmd.Flags().String("bearer_token", "", "bearer_token is a required field.")
	RootBothAuthServiceCmd.AddCommand(BothAuthServicecookieCmd)
	BothAuthServicecookieCmd.Flags().String("bearer_token", "", "bearer_token is a required field.")
	RootBothAuthServiceCmd.AddCommand(BothAuthServicenoneCmd)
	RootBothAuthServiceCmd.AddCommand(BothAuthServicewithArgCmd)
	BothAuthServicewithArgCmd.Flags().String("arg", "", "arg is a required param.")
	BothAuthServicewithArgCmd.Flags().String("bearer_token", "", "bearer_token is a required field.")

	// CookieAuthService commands and flags
	RootCookieAuthServiceCmd.PersistentFlags().String("conf", "../var/conf/configuration.yml", "The configuration file is optional. The default path is ./var/conf/configuration.yml.")
	RootCookieAuthServiceCmd.AddCommand(CookieAuthServicecookieCmd)
	CookieAuthServicecookieCmd.Flags().String("bearer_token", "", "bearer_token is a required field.")

	// HeaderAuthService commands and flags
	RootHeaderAuthServiceCmd.PersistentFlags().String("conf", "../var/conf/configuration.yml", "The configuration file is optional. The default path is ./var/conf/configuration.yml.")
	RootHeaderAuthServiceCmd.AddCommand(HeaderAuthServicedefaultCmd)
	HeaderAuthServicedefaultCmd.Flags().String("bearer_token", "", "bearer_token is a required field.")
	RootHeaderAuthServiceCmd.AddCommand(HeaderAuthServicebinaryCmd)
	HeaderAuthServicebinaryCmd.Flags().String("bearer_token", "", "bearer_token is a required field.")
	RootHeaderAuthServiceCmd.AddCommand(HeaderAuthServicebinaryOptionalCmd)
	HeaderAuthServicebinaryOptionalCmd.Flags().String("bearer_token", "", "bearer_token is a required field.")

	// SomeHeaderAuthService commands and flags
	RootSomeHeaderAuthServiceCmd.PersistentFlags().String("conf", "../var/conf/configuration.yml", "The configuration file is optional. The default path is ./var/conf/configuration.yml.")
	RootSomeHeaderAuthServiceCmd.AddCommand(SomeHeaderAuthServicedefaultCmd)
	SomeHeaderAuthServicedefaultCmd.Flags().String("bearer_token", "", "bearer_token is a required field.")
	RootSomeHeaderAuthServiceCmd.AddCommand(SomeHeaderAuthServicenoneCmd)
}
