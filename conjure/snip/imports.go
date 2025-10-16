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
	"github.com/dave/jennifer/jen"
)

const (
	pal = "github.com/palantir/"
	cgr = pal + "conjure-go-runtime/v2/"
	wgs = pal + "witchcraft-go-server/v2/"
	wgl = pal + "witchcraft-go-logging/"
)

var DefaultImportsToPackageNames = map[string]string{
	cgr + "conjure-go-client/httpclient":   "httpclient",
	cgr + "conjure-go-contract/codecs":     "codecs",
	cgr + "conjure-go-contract/errors":     "errors",
	cgr + "conjure-go-server/httpserver":   "httpserver",
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
	wgs + "witchcraft/wresource":           "wresource",
	wgs + "wrouter":                        "wrouter",
	"gopkg.in/yaml.v3":                     "yaml",
	"github.com/spf13/cobra":               "cobra",
	"github.com/spf13/pflag":               "pflag",
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

	CGRClientClient                            = jen.Qual(cgr+"conjure-go-client/httpclient", "Client").Clone
	CGRClientNewClient                         = jen.Qual(cgr+"conjure-go-client/httpclient", "NewClient").Clone
	CGRClientClientConfig                      = jen.Qual(cgr+"conjure-go-client/httpclient", "ClientConfig").Clone
	CGRClientWithConfig                        = jen.Qual(cgr+"conjure-go-client/httpclient", "WithConfig").Clone
	CGRClientRequestBody                       = jen.Qual(cgr+"conjure-go-client/httpclient", "RequestBody").Clone
	CGRClientRequestBodyInMemory               = jen.Qual(cgr+"conjure-go-client/httpclient", "RequestBodyInMemory").Clone
	CGRClientRequestBodyStreamOnce             = jen.Qual(cgr+"conjure-go-client/httpclient", "RequestBodyStreamOnce").Clone
	CGRClientRequestBodyStreamWithReplay       = jen.Qual(cgr+"conjure-go-client/httpclient", "RequestBodyStreamWithReplay").Clone
	CGRClientRequestParam                      = jen.Qual(cgr+"conjure-go-client/httpclient", "RequestParam").Clone
	CGRClientTokenProvider                     = jen.Qual(cgr+"conjure-go-client/httpclient", "TokenProvider").Clone
	CGRClientWithBinaryRequestBody             = jen.Qual(cgr+"conjure-go-client/httpclient", "WithBinaryRequestBody").Clone
	CGRClientWithHeader                        = jen.Qual(cgr+"conjure-go-client/httpclient", "WithHeader").Clone
	CGRClientWithJSONRequest                   = jen.Qual(cgr+"conjure-go-client/httpclient", "WithJSONRequest").Clone
	CGRClientWithJSONResponse                  = jen.Qual(cgr+"conjure-go-client/httpclient", "WithJSONResponse").Clone
	CGRClientWithPathf                         = jen.Qual(cgr+"conjure-go-client/httpclient", "WithPathf").Clone
	CGRClientWithQueryValues                   = jen.Qual(cgr+"conjure-go-client/httpclient", "WithQueryValues").Clone
	CGRClientWithRPCMethodName                 = jen.Qual(cgr+"conjure-go-client/httpclient", "WithRPCMethodName").Clone
	CGRClientWithRawResponseBody               = jen.Qual(cgr+"conjure-go-client/httpclient", "WithRawResponseBody").Clone
	CGRClientWithRequestConjureErrorDecoder    = jen.Qual(cgr+"conjure-go-client/httpclient", "WithRequestConjureErrorDecoder").Clone
	CGRClientWithRequestMethod                 = jen.Qual(cgr+"conjure-go-client/httpclient", "WithRequestMethod").Clone
	CGRCodecsBinary                            = jen.Qual(cgr+"conjure-go-contract/codecs", "Binary").Clone
	CGRCodecsJSON                              = jen.Qual(cgr+"conjure-go-contract/codecs", "JSON").Clone
	CGRErrorsPermissionDenied                  = jen.Qual(cgr+"conjure-go-contract/errors", "PermissionDenied").Clone
	CGRErrorsInvalidArgument                   = jen.Qual(cgr+"conjure-go-contract/errors", "InvalidArgument").Clone
	CGRErrorsNotFound                          = jen.Qual(cgr+"conjure-go-contract/errors", "NotFound").Clone
	CGRErrorsConflict                          = jen.Qual(cgr+"conjure-go-contract/errors", "Conflict").Clone
	CGRErrorsRequestEntityTooLarge             = jen.Qual(cgr+"conjure-go-contract/errors", "RequestEntityTooLarge").Clone
	CGRErrorsFailedPrecondition                = jen.Qual(cgr+"conjure-go-contract/errors", "FailedPrecondition").Clone
	CGRErrorsInternal                          = jen.Qual(cgr+"conjure-go-contract/errors", "Internal").Clone
	CGRErrorsTimeout                           = jen.Qual(cgr+"conjure-go-contract/errors", "Timeout").Clone
	CGRErrorsCustomClient                      = jen.Qual(cgr+"conjure-go-contract/errors", "CustomClient").Clone
	CGRErrorsCustomServer                      = jen.Qual(cgr+"conjure-go-contract/errors", "CustomServer").Clone
	CGRErrorsErrorCode                         = jen.Qual(cgr+"conjure-go-contract/errors", "ErrorCode").Clone
	CGRErrorsGetConjureError                   = jen.Qual(cgr+"conjure-go-contract/errors", "GetConjureError").Clone
	CGRErrorsNewInternal                       = jen.Qual(cgr+"conjure-go-contract/errors", "NewInternal").Clone
	CGRErrorsNewInvalidArgument                = jen.Qual(cgr+"conjure-go-contract/errors", "NewInvalidArgument").Clone
	CGRErrorsNewReflectTypeConjureErrorDecoder = jen.Qual(cgr+"conjure-go-contract/errors", "NewReflectTypeConjureErrorDecoder").Clone
	CGRErrorsConjureErrorDecoder               = jen.Qual(cgr+"conjure-go-contract/errors", "ConjureErrorDecoder").Clone
	CGRErrorsSerializableError                 = jen.Qual(cgr+"conjure-go-contract/errors", "SerializableError").Clone
	CGRErrorsWrapWithInvalidArgument           = jen.Qual(cgr+"conjure-go-contract/errors", "WrapWithInvalidArgument").Clone
	CGRErrorsWrapWithPermissionDenied          = jen.Qual(cgr+"conjure-go-contract/errors", "WrapWithPermissionDenied").Clone
	CGRHTTPServerErrHandler                    = jen.Qual(cgr+"conjure-go-server/httpserver", "ErrHandler").Clone
	CGRHTTPServerNewJSONHandler                = jen.Qual(cgr+"conjure-go-server/httpserver", "NewJSONHandler").Clone
	CGRHTTPServerParseBearerTokenHeader        = jen.Qual(cgr+"conjure-go-server/httpserver", "ParseBearerTokenHeader").Clone
	CGRHTTPServerStatusCodeMapper              = jen.Qual(cgr+"conjure-go-server/httpserver", "StatusCodeMapper").Clone

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

	WresourceNew                 = jen.Qual(wgs+"witchcraft/wresource", "New").Clone
	WrouterPathParams            = jen.Qual(wgs+"wrouter", "PathParams").Clone
	WrouterRouteParam            = jen.Qual(wgs+"wrouter", "RouteParam").Clone
	WrouterRouter                = jen.Qual(wgs+"wrouter", "Router").Clone
	WrouterForbiddenHeaderParams = jen.Qual(wgs+"wrouter", "ForbiddenHeaderParams").Clone
	WrouterForbiddenPathParams   = jen.Qual(wgs+"wrouter", "ForbiddenPathParams").Clone
	WrouterForbiddenQueryParams  = jen.Qual(wgs+"wrouter", "ForbiddenQueryParams").Clone
	WrouterSafeHeaderParams      = jen.Qual(wgs+"wrouter", "SafeHeaderParams").Clone
	WrouterSafePathParams        = jen.Qual(wgs+"wrouter", "SafePathParams").Clone
	WrouterSafeQueryParams       = jen.Qual(wgs+"wrouter", "SafeQueryParams").Clone

	TAny          = jen.Op("[").Id("T").Id("any").Op("]").Clone
	YamlUnmarshal = jen.Qual("gopkg.in/yaml.v3", "Unmarshal").Clone

	CobraCommand = jen.Qual("github.com/spf13/cobra", "Command").Clone

	PflagsFlagset = jen.Qual("github.com/spf13/pflag", "FlagSet").Clone
)
