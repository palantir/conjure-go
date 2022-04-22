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
)

var DefaultImportsToPackageNames = map[string]string{
	cgr + "conjure-go-client/httpclient": "httpclient",
	cgr + "conjure-go-contract/codecs":   "codecs",
	cgr + "conjure-go-contract/errors":   "errors",
	cgr + "conjure-go-server/httpserver": "httpserver",
	pal + "pkg/binary":                   "binary",
	pal + "pkg/bearertoken":              "bearertoken",
	pal + "pkg/boolean":                  "boolean",
	pal + "pkg/datetime":                 "datetime",
	pal + "pkg/rid":                      "rid",
	pal + "pkg/safelong":                 "safelong",
	pal + "pkg/safejson":                 "safejson",
	pal + "pkg/safeyaml":                 "safeyaml",
	pal + "pkg/uuid":                     "uuid",
	pal + "witchcraft-go-error":          "werror",
	pal + "witchcraft-go-params":         "wparams",
	wgs + "witchcraft/wresource":         "wresource",
	wgs + "wrouter":                      "wrouter",
	"github.com/tidwall/gjson":           "gjson",
}

// A set of imported references included in generated code.
// Each entry is a func() *jen.Statement, typically the Clone method.
// This ensures there are no side effects caused by mutating the global variables.
var (
	Context             = jen.Qual("context", "Context").Clone
	ContextTODO         = jen.Qual("context", "TODO").Clone
	ContextVar          = jen.Id("ctx").Qual("context", "Context").Clone
	Base64Encode        = jen.Qual("encoding/base64", "StdEncoding").Dot("Encode").Clone
	Base64EncodedLen    = jen.Qual("encoding/base64", "StdEncoding").Dot("EncodedLen").Clone
	FmtErrorf           = jen.Qual("fmt", "Errorf").Clone
	FmtSprint           = jen.Qual("fmt", "Sprint").Clone
	FmtSprintf          = jen.Qual("fmt", "Sprintf").Clone
	IOReadCloser        = jen.Qual("io", "ReadCloser").Clone
	IOUtilReadAll       = jen.Qual("io/ioutil", "ReadAll").Clone
	JSONMarshaler       = jen.Qual("encoding/json", "Marshaler").Clone
	JSONUnmarshaler     = jen.Qual("encoding/json", "Unmarshaler").Clone
	MathIsInf           = jen.Qual("math", "IsInf").Clone
	MathIsNaN           = jen.Qual("math", "IsNaN").Clone
	MathInf             = jen.Qual("math", "Inf").Clone
	MathNaN             = jen.Qual("math", "NaN").Clone
	HTTPNoBody          = jen.Qual("net/http", "NoBody").Clone
	HTTPStatusNoContent = jen.Qual("net/http", "StatusNoContent").Clone
	HTTPRequest         = jen.Qual("net/http", "Request").Clone
	HTTPResponseWriter  = jen.Qual("net/http", "ResponseWriter").Clone
	URLPathEscape       = jen.Qual("net/url", "PathEscape").Clone
	URLValues           = jen.Qual("net/url", "Values").Clone
	ReflectTypeOf       = jen.Qual("reflect", "TypeOf").Clone
	StringsToUpper      = jen.Qual("strings", "ToUpper").Clone
	StrconvAppendFloat  = jen.Qual("strconv", "AppendFloat").Clone
	StrconvAppendInt    = jen.Qual("strconv", "AppendInt").Clone
	StrconvAtoi         = jen.Qual("strconv", "Atoi").Clone
	StrconvItoa         = jen.Qual("strconv", "Itoa").Clone
	StrconvParseBool    = jen.Qual("strconv", "ParseBool").Clone
	StrconvParseFloat   = jen.Qual("strconv", "ParseFloat").Clone
	StrconvQuote        = jen.Qual("strconv", "Quote").Clone
	FuncIOReadCloser    = jen.Func().Params().Params(IOReadCloser()).Clone // 'func() io.ReadCloser', the type of to http.Request.GetBody.

	CGRClientClient                     = jen.Qual(cgr+"conjure-go-client/httpclient", "Client").Clone
	CGRClientRequestParam               = jen.Qual(cgr+"conjure-go-client/httpclient", "RequestParam").Clone
	CGRClientTokenProvider              = jen.Qual(cgr+"conjure-go-client/httpclient", "TokenProvider").Clone
	CGRClientWithHeader                 = jen.Qual(cgr+"conjure-go-client/httpclient", "WithHeader").Clone
	CGRClientWithJSONRequest            = jen.Qual(cgr+"conjure-go-client/httpclient", "WithJSONRequest").Clone
	CGRClientWithJSONResponse           = jen.Qual(cgr+"conjure-go-client/httpclient", "WithJSONResponse").Clone
	CGRClientWithPathf                  = jen.Qual(cgr+"conjure-go-client/httpclient", "WithPathf").Clone
	CGRClientWithQueryValues            = jen.Qual(cgr+"conjure-go-client/httpclient", "WithQueryValues").Clone
	CGRClientWithRPCMethodName          = jen.Qual(cgr+"conjure-go-client/httpclient", "WithRPCMethodName").Clone
	CGRClientWithRawRequestBodyProvider = jen.Qual(cgr+"conjure-go-client/httpclient", "WithRawRequestBodyProvider").Clone
	CGRClientWithRawResponseBody        = jen.Qual(cgr+"conjure-go-client/httpclient", "WithRawResponseBody").Clone
	CGRClientWithRequiredResponse       = jen.Qual(cgr+"conjure-go-client/httpclient", "WithRequiredResponse").Clone
	CGRClientWithRequestMethod          = jen.Qual(cgr+"conjure-go-client/httpclient", "WithRequestMethod").Clone
	CGRClientWithRequestAppendFunc      = jen.Qual(cgr+"conjure-go-client/httpclient", "WithRequestAppendFunc").Clone
	CGRClientWithResponseUnmarshalFunc  = jen.Qual(cgr+"conjure-go-client/httpclient", "WithResponseUnmarshalFunc").Clone
	CGRCodecsBinary                     = jen.Qual(cgr+"conjure-go-contract/codecs", "Binary").Clone
	CGRCodecsJSON                       = jen.Qual(cgr+"conjure-go-contract/codecs", "JSON").Clone
	CGRErrorsPermissionDenied           = jen.Qual(cgr+"conjure-go-contract/errors", "PermissionDenied").Clone
	CGRErrorsInvalidArgument            = jen.Qual(cgr+"conjure-go-contract/errors", "InvalidArgument").Clone
	CGRErrorsNotFound                   = jen.Qual(cgr+"conjure-go-contract/errors", "NotFound").Clone
	CGRErrorsConflict                   = jen.Qual(cgr+"conjure-go-contract/errors", "Conflict").Clone
	CGRErrorsRequestEntityTooLarge      = jen.Qual(cgr+"conjure-go-contract/errors", "RequestEntityTooLarge").Clone
	CGRErrorsFailedPrecondition         = jen.Qual(cgr+"conjure-go-contract/errors", "FailedPrecondition").Clone
	CGRErrorsInternal                   = jen.Qual(cgr+"conjure-go-contract/errors", "Internal").Clone
	CGRErrorsTimeout                    = jen.Qual(cgr+"conjure-go-contract/errors", "Timeout").Clone
	CGRErrorsCustomClient               = jen.Qual(cgr+"conjure-go-contract/errors", "CustomClient").Clone
	CGRErrorsCustomServer               = jen.Qual(cgr+"conjure-go-contract/errors", "CustomServer").Clone
	CGRErrorsErrorCode                  = jen.Qual(cgr+"conjure-go-contract/errors", "ErrorCode").Clone
	CGRErrorsGetConjureError            = jen.Qual(cgr+"conjure-go-contract/errors", "GetConjureError").Clone
	CGRErrorsNewInternal                = jen.Qual(cgr+"conjure-go-contract/errors", "NewInternal").Clone
	CGRErrorsNewInvalidArgument         = jen.Qual(cgr+"conjure-go-contract/errors", "NewInvalidArgument").Clone
	CGRErrorsRegisterErrorType          = jen.Qual(cgr+"conjure-go-contract/errors", "RegisterErrorType").Clone
	CGRErrorsSerializableError          = jen.Qual(cgr+"conjure-go-contract/errors", "SerializableError").Clone
	CGRErrorsWrapWithInternal           = jen.Qual(cgr+"conjure-go-contract/errors", "WrapWithInternal").Clone
	CGRErrorsWrapWithInvalidArgument    = jen.Qual(cgr+"conjure-go-contract/errors", "WrapWithInvalidArgument").Clone
	CGRErrorsWrapWithPermissionDenied   = jen.Qual(cgr+"conjure-go-contract/errors", "WrapWithPermissionDenied").Clone
	CGRHTTPServerErrHandler             = jen.Qual(cgr+"conjure-go-server/httpserver", "ErrHandler").Clone
	CGRHTTPServerNewJSONHandler         = jen.Qual(cgr+"conjure-go-server/httpserver", "NewJSONHandler").Clone
	CGRHTTPServerParseBearerTokenHeader = jen.Qual(cgr+"conjure-go-server/httpserver", "ParseBearerTokenHeader").Clone
	CGRHTTPServerStatusCodeMapper       = jen.Qual(cgr+"conjure-go-server/httpserver", "StatusCodeMapper").Clone

	BinaryBinary                   = jen.Qual(pal+"pkg/binary", "Binary").Clone
	BinaryNew                      = jen.Qual(pal+"pkg/binary", "New").Clone
	BearerTokenNew                 = jen.Qual(pal+"pkg/bearertoken", "New").Clone
	BearerTokenToken               = jen.Qual(pal+"pkg/bearertoken", "Token").Clone
	BooleanBoolean                 = jen.Qual(pal+"pkg/boolean", "Boolean").Clone
	DateTimeDateTime               = jen.Qual(pal+"pkg/datetime", "DateTime").Clone
	DateTimeParseDateTime          = jen.Qual(pal+"pkg/datetime", "ParseDateTime").Clone
	RIDParseRID                    = jen.Qual(pal+"pkg/rid", "ParseRID").Clone
	RIDResourceIdentifier          = jen.Qual(pal+"pkg/rid", "ResourceIdentifier").Clone
	SafeLongNewSafeLong            = jen.Qual(pal+"pkg/safelong", "NewSafeLong").Clone
	SafeLongParseSafeLong          = jen.Qual(pal+"pkg/safelong", "ParseSafeLong").Clone
	SafeLongSafeLong               = jen.Qual(pal+"pkg/safelong", "SafeLong").Clone
	SafeJSONAppendFunc             = jen.Qual(pal+"pkg/safejson", "AppendFunc").Clone
	SafeJSONAppendQuotedString     = jen.Qual(pal+"pkg/safejson", "AppendQuotedString").Clone
	SafeJSONMarshal                = jen.Qual(pal+"pkg/safejson", "Marshal").Clone
	SafeJSONQuotedStringLength     = jen.Qual(pal+"pkg/safejson", "QuotedStringLength").Clone
	SafeJSONQuoteString            = jen.Qual(pal+"pkg/safejson", "QuoteString").Clone
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
	WerrorSafeParam       = jen.Qual(pal+"witchcraft-go-error", "SafeParam").Clone
	WerrorStackTrace      = jen.Qual(pal+"witchcraft-go-error", "StackTrace").Clone
	WerrorUnsafeParam     = jen.Qual(pal+"witchcraft-go-error", "UnsafeParam").Clone
	WerrorWrap            = jen.Qual(pal+"witchcraft-go-error", "Wrap").Clone
	WerrorWrapContext     = jen.Qual(pal+"witchcraft-go-error", "WrapWithContextParams").Clone

	WresourceNew            = jen.Qual(wgs+"witchcraft/wresource", "New").Clone
	WrouterPathParams       = jen.Qual(wgs+"wrouter", "PathParams").Clone
	WrouterRouteParam       = jen.Qual(wgs+"wrouter", "RouteParam").Clone
	WrouterRouter           = jen.Qual(wgs+"wrouter", "Router").Clone
	WrouterSafeHeaderParams = jen.Qual(wgs+"wrouter", "SafeHeaderParams").Clone
	WrouterSafePathParams   = jen.Qual(wgs+"wrouter", "SafePathParams").Clone
	WrouterSafeQueryParams  = jen.Qual(wgs+"wrouter", "SafeQueryParams").Clone

	GJSONNull       = jen.Qual("github.com/tidwall/gjson", "Null").Clone
	GJSONFalse      = jen.Qual("github.com/tidwall/gjson", "False").Clone
	GJSONNumber     = jen.Qual("github.com/tidwall/gjson", "Number").Clone
	GJSONString     = jen.Qual("github.com/tidwall/gjson", "String").Clone
	GJSONTrue       = jen.Qual("github.com/tidwall/gjson", "True").Clone
	GJSONJSON       = jen.Qual("github.com/tidwall/gjson", "JSON").Clone
	GJSONParse      = jen.Qual("github.com/tidwall/gjson", "Parse").Clone
	GJSONParseBytes = jen.Qual("github.com/tidwall/gjson", "ParseBytes").Clone
	GJSONResult     = jen.Qual("github.com/tidwall/gjson", "Result").Clone
	GJSONValid      = jen.Qual("github.com/tidwall/gjson", "Valid").Clone
	GJSONValidBytes = jen.Qual("github.com/tidwall/gjson", "ValidBytes").Clone

	TAny = jen.Op("[").Id("T").Id("any").Op("]").Clone
)
