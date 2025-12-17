// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package snip

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

const (
	pal = "github.com/palantir/"
	wgl = pal + "witchcraft-go-logging/"
)

var (
	// cgrModuleVersion controls which conjure-go-runtime module version to use in generated code.
	// Defaults to 2. Set via SetCGRVersion before generation.
	cgrModuleVersion = 2
	// wgsModuleVersion controls which witchcraft-go-server module version to use in generated code.
	// Defaults to 2. Set via SetWGSVersion before generation.
	wgsModuleVersion = 2
)

// SetCGRModuleVersion sets the version of conjure-go-runtime to use in generated imports.
// Valid values are 2 or 3.
func SetCGRModuleVersion(version int) {
	if version != 2 && version != 3 {
		version = cgrModuleVersion
	}
	cgrModuleVersion = version
}

// SetWGSModuleVersion sets the version of witchcraft-go-server to use in generated imports.
// Valid values are 2 or 3.
func SetWGSModuleVersion(version int) {
	if version != 2 && version != 3 {
		version = wgsModuleVersion
	}
	wgsModuleVersion = version
}

// cgr returns the conjure-go-runtime import path prefix for the configured version.
func cgr() string {
	return fmt.Sprintf("%sconjure-go-runtime/v%d/", pal, cgrModuleVersion)
}

// wgs returns the witchcraft-go-server import path prefix for the configured version.
func wgs() string {
	return fmt.Sprintf("%switchcraft-go-server/v%d/", pal, wgsModuleVersion)
}

// ImportsToPackageNames returns a map of import paths to package names for use with jennifer.
// This must be called after SetCGRVersion/SetWGSVersion to get the correct version-specific paths.
func ImportsToPackageNames() map[string]string {
	return map[string]string{
		cgr() + "conjure-go-client/httpclient": "httpclient",
		cgr() + "conjure-go-contract/codecs":   "codecs",
		cgr() + "conjure-go-contract/errors":   "errors",
		cgr() + "conjure-go-server/httpserver": "httpserver",
		pal + "pkg/binary":                     "binary",
		pal + "pkg/bearertoken":                "bearertoken",
		pal + "pkg/boolean":                    "boolean",
		pal + "pkg/datetime":                   "datetime",
		pal + "pkg/rid":                        "rid",
		pal + "pkg/safelong":                   "safelong",
		pal + "pkg/safejson":                   "safejson",
		pal + "pkg/safeyaml":                   "safeyaml",
		pal + "pkg/uuid":                       "uuid",
		pal + "witchcraft-go-error":            "werror",
		pal + "witchcraft-go-params":           "wparams",
		wgl + "wlog":                           "wlog",
		wgl + "wlog-zap":                       "wlogzap",
		wgl + "wlog/evtlog/evt2log":            "evt2log",
		wgl + "wlog/svclog/svc1log":            "svc1log",
		wgl + "wlog/trclog/trc1log":            "trc1log",
		pal + "witchcraft-go-tracing/wtracing": "wtracing",
		pal + "witchcraft-go-tracing/wzipkin":  "wzipkin",
		wgs() + "witchcraft/wresource":         "wresource",
		wgs() + "wrouter":                      "wrouter",
		"gopkg.in/yaml.v3":                     "yaml",
		"github.com/spf13/cobra":               "cobra",
		"github.com/spf13/pflag":               "pflag",
	}
}

// A set of imported references included in generated code.
// Each entry is a func() *jen.Statement, typically the Clone method.
// This ensures there are no side effects caused by mutating the global variables.
var (
	ByteReader          = jen.Qual("bytes", "NewReader").Clone
	Context             = jen.Qual("context", "Context").Clone
	ContextTODO         = jen.Qual("context", "TODO").Clone
	ContextBackground   = jen.Qual("context", "Background").Clone
	ContextVar          = jen.Id("ctx").Qual("context", "Context").Clone
	Base64StdEncoding   = jen.Qual("encoding/base64", "StdEncoding").Clone
	JSONMarshalIndent   = jen.Qual("encoding/json", "MarshalIndent").Clone
	FmtErrorf           = jen.Qual("fmt", "Errorf").Clone
	FmtPrintf           = jen.Qual("fmt", "Printf").Clone
	FmtFprintf          = jen.Qual("fmt", "Fprintf").Clone
	FmtSprint           = jen.Qual("fmt", "Sprint").Clone
	FmtSprintf          = jen.Qual("fmt", "Sprintf").Clone
	IOReadCloser        = jen.Qual("io", "ReadCloser").Clone
	IONopCloser         = jen.Qual("io", "NopCloser").Clone
	IOCopy              = jen.Qual("io", "Copy").Clone
	IODiscard           = jen.Qual("io", "Discard").Clone
	HTTPNoBody          = jen.Qual("net/http", "NoBody").Clone
	HTTPStatusNoContent = jen.Qual("net/http", "StatusNoContent").Clone
	HTTPRequest         = jen.Qual("net/http", "Request").Clone
	HTTPResponseWriter  = jen.Qual("net/http", "ResponseWriter").Clone
	URLPathEscape       = jen.Qual("net/url", "PathEscape").Clone
	URLValues           = jen.Qual("net/url", "Values").Clone
	OSStdout            = jen.Qual("os", "Stdout").Clone
	OSReadFile          = jen.Qual("os", "ReadFile").Clone
	OSOpen              = jen.Qual("os", "Open").Clone
	ReflectTypeOf       = jen.Qual("reflect", "TypeOf").Clone
	StringsToUpper      = jen.Qual("strings", "ToUpper").Clone
	StringsHasPrefix    = jen.Qual("strings", "HasPrefix").Clone
	StringsTrimSpace    = jen.Qual("strings", "TrimSpace").Clone
	StrconvAtoi         = jen.Qual("strconv", "Atoi").Clone
	StrconvParseBool    = jen.Qual("strconv", "ParseBool").Clone
	StrconvParseFloat   = jen.Qual("strconv", "ParseFloat").Clone
)

// conjure-go-runtime imports (version-dependent)

func CGRClientClient() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "Client")
}
func CGRClientNewClient() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "NewClient")
}
func CGRClientClientConfig() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "ClientConfig")
}
func CGRClientWithConfig() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "WithConfig")
}
func CGRClientRequestBody() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "RequestBody")
}
func CGRClientRequestBodyInMemory() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "RequestBodyInMemory")
}
func CGRClientRequestBodyStreamOnce() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "RequestBodyStreamOnce")
}
func CGRClientRequestBodyStreamWithReplay() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "RequestBodyStreamWithReplay")
}
func CGRClientRequestParam() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "RequestParam")
}
func CGRClientTokenProvider() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "TokenProvider")
}
func CGRClientWithBinaryRequestBody() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "WithBinaryRequestBody")
}
func CGRClientWithHeader() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "WithHeader")
}
func CGRClientWithJSONRequest() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "WithJSONRequest")
}
func CGRClientWithJSONResponse() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "WithJSONResponse")
}
func CGRClientWithPathf() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "WithPathf")
}
func CGRClientWithQueryValues() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "WithQueryValues")
}
func CGRClientWithRPCMethodName() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "WithRPCMethodName")
}
func CGRClientWithRawResponseBody() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "WithRawResponseBody")
}
func CGRClientWithRequestConjureErrorDecoder() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "WithRequestConjureErrorDecoder")
}
func CGRClientWithRequestMethod() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-client/httpclient", "WithRequestMethod")
}
func CGRCodecsBinary() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/codecs", "Binary")
}
func CGRCodecsJSON() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/codecs", "JSON")
}
func CGRErrorsPermissionDenied() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "PermissionDenied")
}
func CGRErrorsInvalidArgument() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "InvalidArgument")
}
func CGRErrorsNotFound() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "NotFound")
}
func CGRErrorsConflict() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "Conflict")
}
func CGRErrorsRequestEntityTooLarge() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "RequestEntityTooLarge")
}
func CGRErrorsFailedPrecondition() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "FailedPrecondition")
}
func CGRErrorsInternal() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "Internal")
}
func CGRErrorsTimeout() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "Timeout")
}
func CGRErrorsCustomClient() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "CustomClient")
}
func CGRErrorsCustomServer() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "CustomServer")
}
func CGRErrorsErrorCode() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "ErrorCode")
}
func CGRErrorsGetConjureError() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "GetConjureError")
}
func CGRErrorsNewInternal() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "NewInternal")
}
func CGRErrorsNewInvalidArgument() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "NewInvalidArgument")
}
func CGRErrorsNewReflectTypeConjureErrorDecoder() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "NewReflectTypeConjureErrorDecoder")
}
func CGRErrorsConjureErrorDecoder() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "ConjureErrorDecoder")
}
func CGRErrorsSerializableError() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "SerializableError")
}
func CGRErrorsWrapWithInvalidArgument() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "WrapWithInvalidArgument")
}
func CGRErrorsWrapWithPermissionDenied() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-contract/errors", "WrapWithPermissionDenied")
}
func CGRHTTPServerErrHandler() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-server/httpserver", "ErrHandler")
}
func CGRHTTPServerNewJSONHandler() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-server/httpserver", "NewJSONHandler")
}
func CGRHTTPServerParseBearerTokenHeader() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-server/httpserver", "ParseBearerTokenHeader")
}
func CGRHTTPServerStatusCodeMapper() *jen.Statement {
	return jen.Qual(cgr()+"conjure-go-server/httpserver", "StatusCodeMapper")
}

// witchcraft-go-server imports (version-dependent)

func WresourceNew() *jen.Statement {
	return jen.Qual(wgs()+"witchcraft/wresource", "New")
}
func WrouterPathParams() *jen.Statement {
	return jen.Qual(wgs()+"wrouter", "PathParams")
}
func WrouterRouteParam() *jen.Statement {
	return jen.Qual(wgs()+"wrouter", "RouteParam")
}
func WrouterRouter() *jen.Statement {
	return jen.Qual(wgs()+"wrouter", "Router")
}
func WrouterForbiddenHeaderParams() *jen.Statement {
	return jen.Qual(wgs()+"wrouter", "ForbiddenHeaderParams")
}
func WrouterForbiddenPathParams() *jen.Statement {
	return jen.Qual(wgs()+"wrouter", "ForbiddenPathParams")
}
func WrouterForbiddenQueryParams() *jen.Statement {
	return jen.Qual(wgs()+"wrouter", "ForbiddenQueryParams")
}
func WrouterSafeHeaderParams() *jen.Statement {
	return jen.Qual(wgs()+"wrouter", "SafeHeaderParams")
}
func WrouterSafePathParams() *jen.Statement {
	return jen.Qual(wgs()+"wrouter", "SafePathParams")
}
func WrouterSafeQueryParams() *jen.Statement {
	return jen.Qual(wgs()+"wrouter", "SafeQueryParams")
}

var (
	BinaryBinary                   = jen.Qual(pal+"pkg/binary", "Binary").Clone
	BinaryNew                      = jen.Qual(pal+"pkg/binary", "New").Clone
	BearerTokenToken               = jen.Qual(pal+"pkg/bearertoken", "Token").Clone
	BooleanBoolean                 = jen.Qual(pal+"pkg/boolean", "Boolean").Clone
	DateTimeDateTime               = jen.Qual(pal+"pkg/datetime", "DateTime").Clone
	DateTimeParseDateTime          = jen.Qual(pal+"pkg/datetime", "ParseDateTime").Clone
	RIDParseRID                    = jen.Qual(pal+"pkg/rid", "ParseRID").Clone
	RIDResourceIdentifier          = jen.Qual(pal+"pkg/rid", "ResourceIdentifier").Clone
	SafeLongParseSafeLong          = jen.Qual(pal+"pkg/safelong", "ParseSafeLong").Clone
	SafeLongSafeLong               = jen.Qual(pal+"pkg/safelong", "SafeLong").Clone
	SafeJSONMarshal                = jen.Qual(pal+"pkg/safejson", "Marshal").Clone
	SafeJSONUnmarshal              = jen.Qual(pal+"pkg/safejson", "Unmarshal").Clone
	SafeYAMLJSONtoYAMLMapSlice     = jen.Qual(pal+"pkg/safeyaml", "JSONtoYAMLMapSlice").Clone
	SafeYAMLUnmarshalerToJSONBytes = jen.Qual(pal+"pkg/safeyaml", "UnmarshalerToJSONBytes").Clone
	UUIDUUID                       = jen.Qual(pal+"pkg/uuid", "UUID").Clone
	UUIDNewUUID                    = jen.Qual(pal+"pkg/uuid", "NewUUID").Clone
	UUIDParseUUID                  = jen.Qual(pal+"pkg/uuid", "ParseUUID").Clone

	WerrorErrorContext    = jen.Qual(pal+"witchcraft-go-error", "ErrorWithContextParams").Clone
	WerrorFormat          = jen.Qual(pal+"witchcraft-go-error", "Format").Clone
	WerrorNewStackTrace   = jen.Qual(pal+"witchcraft-go-error", "NewStackTrace").Clone
	WerrorParamsFromError = jen.Qual(pal+"witchcraft-go-error", "ParamsFromError").Clone
	WerrorStackTrace      = jen.Qual(pal+"witchcraft-go-error", "StackTrace").Clone
	WerrorWrapContext     = jen.Qual(pal+"witchcraft-go-error", "WrapWithContextParams").Clone

	WGLLogSetDefaultLoggerProvider = jen.Qual(wgl+"wlog", "SetDefaultLoggerProvider").Clone
	WGLLogNoopLoggerProvider       = jen.Qual(wgl+"wlog", "NewNoopLoggerProvider").Clone
	WGLLogDebugLevel               = jen.Qual(wgl+"wlog", "DebugLevel").Clone
	WGLWlogZapLoggerProvider       = jen.Qual(wgl+"wlog-zap", "LoggerProvider").Clone
	WGLSvc1logWithLogger           = jen.Qual(wgl+"wlog/svclog/svc1log", "WithLogger").Clone
	WGLSvc1logNew                  = jen.Qual(wgl+"wlog/svclog/svc1log", "New").Clone
	WGLTrc1logWithLogger           = jen.Qual(wgl+"wlog/trclog/trc1log", "WithLogger").Clone
	WGLTrc1logNewLogger            = jen.Qual(wgl+"wlog/trclog/trc1log", "New").Clone
	WGLEvt2logWithLogger           = jen.Qual(wgl+"wlog/evtlog/evt2log", "WithLogger").Clone
	WGLEvt2logNew                  = jen.Qual(wgl+"wlog/evtlog/evt2log", "New").Clone
	WGTContextWithTracer           = jen.Qual(pal+"witchcraft-go-tracing/wtracing", "ContextWithTracer").Clone
	WGTZipkinNewTracer             = jen.Qual(pal+"witchcraft-go-tracing/wzipkin", "NewTracer").Clone

	TAny          = jen.Op("[").Id("T").Id("any").Op("]").Clone
	YamlUnmarshal = jen.Qual("gopkg.in/yaml.v3", "Unmarshal").Clone

	CobraCommand = jen.Qual("github.com/spf13/cobra", "Command").Clone

	PflagsFlagset = jen.Qual("github.com/spf13/pflag", "FlagSet").Clone
)
