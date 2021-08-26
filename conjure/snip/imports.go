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
)

// A set of imported references included in generated code.
// Each entry is a func() *jen.Statement, typically the Clone method.
// This ensures there are no side effects caused by mutating the global variables.
var (
	Context             = jen.Qual("context", "Context").Clone
	ContextVar          = jen.Id("ctx").Qual("context", "Context").Clone
	Base64Encode        = jen.Qual("encoding/base64", "StdEncoding").Dot("Encode").Clone
	Base64EncodedLen    = jen.Qual("encoding/base64", "StdEncoding").Dot("EncodedLen").Clone
	FmtErrorf           = jen.Qual("fmt", "Errorf").Clone
	FmtSprint           = jen.Qual("fmt", "Sprint").Clone
	FmtSprintf          = jen.Qual("fmt", "Sprintf").Clone
	IOReadCloser        = jen.Qual("io", "ReadCloser").Clone
	MathIsInf           = jen.Qual("math", "IsInf").Clone
	MathIsNaN           = jen.Qual("math", "IsNaN").Clone
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
	CGRClientWithRequestMethod          = jen.Qual(cgr+"conjure-go-client/httpclient", "WithRequestMethod").Clone
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
	CGRErrorsWrapWithInvalidArgument    = jen.Qual(cgr+"conjure-go-contract/errors", "WrapWithInvalidArgument").Clone
	CGRErrorsWrapWithPermissionDenied   = jen.Qual(cgr+"conjure-go-contract/errors", "WrapWithPermissionDenied").Clone
	CGRHTTPServerErrHandler             = jen.Qual(cgr+"conjure-go-server/httpserver", "ErrHandler").Clone
	CGRHTTPServerNewJSONHandler         = jen.Qual(cgr+"conjure-go-server/httpserver", "NewJSONHandler").Clone
	CGRHTTPServerParseBearerTokenHeader = jen.Qual(cgr+"conjure-go-server/httpserver", "ParseBearerTokenHeader").Clone
	CGRHTTPServerStatusCodeMapper       = jen.Qual(cgr+"conjure-go-server/httpserver", "StatusCodeMapper").Clone

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
	SafeJSONAppendFunc             = jen.Qual(pal+"pkg/safejson", "AppendFunc").Clone
	SafeJSONAppendQuotedString     = jen.Qual(pal+"pkg/safejson", "AppendQuotedString").Clone
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
	WerrorWrap            = jen.Qual(pal+"witchcraft-go-error", "Wrap").Clone
	WerrorWrapContext     = jen.Qual(pal+"witchcraft-go-error", "WrapWithContextParams").Clone

	WresourceNew            = jen.Qual(pal+"witchcraft-go-server/v2/witchcraft/wresource", "New").Clone
	WrouterPathParams       = jen.Qual(pal+"witchcraft-go-server/v2/wrouter", "PathParams").Clone
	WrouterRouter           = jen.Qual(pal+"witchcraft-go-server/v2/wrouter", "Router").Clone
	WrouterSafeHeaderParams = jen.Qual(pal+"witchcraft-go-server/v2/wrouter", "SafeHeaderParams").Clone
	WrouterSafePathParams   = jen.Qual(pal+"witchcraft-go-server/v2/wrouter", "SafePathParams").Clone
	WrouterSafeQueryParams  = jen.Qual(pal+"witchcraft-go-server/v2/wrouter", "SafeQueryParams").Clone
)
