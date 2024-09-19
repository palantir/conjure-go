// This file was generated by Conjure and should not be manually edited.

package spec

import (
	"strings"
)

type ErrorCode struct {
	val ErrorCode_Value
}

type ErrorCode_Value string

const (
	ErrorCode_PERMISSION_DENIED        ErrorCode_Value = "PERMISSION_DENIED"
	ErrorCode_INVALID_ARGUMENT         ErrorCode_Value = "INVALID_ARGUMENT"
	ErrorCode_NOT_FOUND                ErrorCode_Value = "NOT_FOUND"
	ErrorCode_CONFLICT                 ErrorCode_Value = "CONFLICT"
	ErrorCode_REQUEST_ENTITY_TOO_LARGE ErrorCode_Value = "REQUEST_ENTITY_TOO_LARGE"
	ErrorCode_FAILED_PRECONDITION      ErrorCode_Value = "FAILED_PRECONDITION"
	ErrorCode_INTERNAL                 ErrorCode_Value = "INTERNAL"
	ErrorCode_TIMEOUT                  ErrorCode_Value = "TIMEOUT"
	ErrorCode_CUSTOM_CLIENT            ErrorCode_Value = "CUSTOM_CLIENT"
	ErrorCode_CUSTOM_SERVER            ErrorCode_Value = "CUSTOM_SERVER"
	ErrorCode_UNKNOWN                  ErrorCode_Value = "UNKNOWN"
)

func (e ErrorCode_Value) New() ErrorCode {
	return ErrorCode{val: e}
}

// ErrorCode_Values returns all known variants of ErrorCode.
func ErrorCode_Values() []ErrorCode_Value {
	return []ErrorCode_Value{ErrorCode_PERMISSION_DENIED, ErrorCode_INVALID_ARGUMENT, ErrorCode_NOT_FOUND, ErrorCode_CONFLICT, ErrorCode_REQUEST_ENTITY_TOO_LARGE, ErrorCode_FAILED_PRECONDITION, ErrorCode_INTERNAL, ErrorCode_TIMEOUT, ErrorCode_CUSTOM_CLIENT, ErrorCode_CUSTOM_SERVER}
}

func New_ErrorCode(value ErrorCode_Value) ErrorCode {
	return ErrorCode{val: value}
}

// IsUnknown returns false for all known variants of ErrorCode and true otherwise.
func (e ErrorCode) IsUnknown() bool {
	switch e.val {
	case ErrorCode_PERMISSION_DENIED, ErrorCode_INVALID_ARGUMENT, ErrorCode_NOT_FOUND, ErrorCode_CONFLICT, ErrorCode_REQUEST_ENTITY_TOO_LARGE, ErrorCode_FAILED_PRECONDITION, ErrorCode_INTERNAL, ErrorCode_TIMEOUT, ErrorCode_CUSTOM_CLIENT, ErrorCode_CUSTOM_SERVER:
		return false
	}
	return true
}

func (e ErrorCode) Value() ErrorCode_Value {
	if e.IsUnknown() {
		return ErrorCode_UNKNOWN
	}
	return e.val
}

func (e ErrorCode) String() string {
	return string(e.val)
}

func (e ErrorCode) MarshalText() ([]byte, error) {
	return []byte(e.val), nil
}

func (e *ErrorCode) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		*e = ErrorCode_Value(v).New()
	case "PERMISSION_DENIED":
		*e = ErrorCode_PERMISSION_DENIED.New()
	case "INVALID_ARGUMENT":
		*e = ErrorCode_INVALID_ARGUMENT.New()
	case "NOT_FOUND":
		*e = ErrorCode_NOT_FOUND.New()
	case "CONFLICT":
		*e = ErrorCode_CONFLICT.New()
	case "REQUEST_ENTITY_TOO_LARGE":
		*e = ErrorCode_REQUEST_ENTITY_TOO_LARGE.New()
	case "FAILED_PRECONDITION":
		*e = ErrorCode_FAILED_PRECONDITION.New()
	case "INTERNAL":
		*e = ErrorCode_INTERNAL.New()
	case "TIMEOUT":
		*e = ErrorCode_TIMEOUT.New()
	case "CUSTOM_CLIENT":
		*e = ErrorCode_CUSTOM_CLIENT.New()
	case "CUSTOM_SERVER":
		*e = ErrorCode_CUSTOM_SERVER.New()
	}
	return nil
}

type HttpMethod struct {
	val HttpMethod_Value
}

type HttpMethod_Value string

const (
	HttpMethod_GET     HttpMethod_Value = "GET"
	HttpMethod_POST    HttpMethod_Value = "POST"
	HttpMethod_PUT     HttpMethod_Value = "PUT"
	HttpMethod_DELETE  HttpMethod_Value = "DELETE"
	HttpMethod_UNKNOWN HttpMethod_Value = "UNKNOWN"
)

func (e HttpMethod_Value) New() HttpMethod {
	return HttpMethod{val: e}
}

// HttpMethod_Values returns all known variants of HttpMethod.
func HttpMethod_Values() []HttpMethod_Value {
	return []HttpMethod_Value{HttpMethod_GET, HttpMethod_POST, HttpMethod_PUT, HttpMethod_DELETE}
}

func New_HttpMethod(value HttpMethod_Value) HttpMethod {
	return HttpMethod{val: value}
}

// IsUnknown returns false for all known variants of HttpMethod and true otherwise.
func (e HttpMethod) IsUnknown() bool {
	switch e.val {
	case HttpMethod_GET, HttpMethod_POST, HttpMethod_PUT, HttpMethod_DELETE:
		return false
	}
	return true
}

func (e HttpMethod) Value() HttpMethod_Value {
	if e.IsUnknown() {
		return HttpMethod_UNKNOWN
	}
	return e.val
}

func (e HttpMethod) String() string {
	return string(e.val)
}

func (e HttpMethod) MarshalText() ([]byte, error) {
	return []byte(e.val), nil
}

func (e *HttpMethod) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		*e = HttpMethod_Value(v).New()
	case "GET":
		*e = HttpMethod_GET.New()
	case "POST":
		*e = HttpMethod_POST.New()
	case "PUT":
		*e = HttpMethod_PUT.New()
	case "DELETE":
		*e = HttpMethod_DELETE.New()
	}
	return nil
}

// Safety with regards to logging based on [safe-logging](https://github.com/palantir/safe-logging) concepts.
type LogSafety struct {
	val LogSafety_Value
}

type LogSafety_Value string

const (
	// Explicitly marks an element as safe.
	LogSafety_SAFE LogSafety_Value = "SAFE"
	// Explicitly marks an element as unsafe, diallowing contents from being logged as `SAFE`.
	LogSafety_UNSAFE LogSafety_Value = "UNSAFE"
	// Marks elements that must never be logged. For example, credentials, keys, and other secrets cannot be logged because such an action would compromise security.
	LogSafety_DO_NOT_LOG LogSafety_Value = "DO_NOT_LOG"
	LogSafety_UNKNOWN    LogSafety_Value = "UNKNOWN"
)

func (e LogSafety_Value) New() LogSafety {
	return LogSafety{val: e}
}

// LogSafety_Values returns all known variants of LogSafety.
func LogSafety_Values() []LogSafety_Value {
	return []LogSafety_Value{LogSafety_SAFE, LogSafety_UNSAFE, LogSafety_DO_NOT_LOG}
}

func New_LogSafety(value LogSafety_Value) LogSafety {
	return LogSafety{val: value}
}

// IsUnknown returns false for all known variants of LogSafety and true otherwise.
func (e LogSafety) IsUnknown() bool {
	switch e.val {
	case LogSafety_SAFE, LogSafety_UNSAFE, LogSafety_DO_NOT_LOG:
		return false
	}
	return true
}

func (e LogSafety) Value() LogSafety_Value {
	if e.IsUnknown() {
		return LogSafety_UNKNOWN
	}
	return e.val
}

func (e LogSafety) String() string {
	return string(e.val)
}

func (e LogSafety) MarshalText() ([]byte, error) {
	return []byte(e.val), nil
}

func (e *LogSafety) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		*e = LogSafety_Value(v).New()
	case "SAFE":
		*e = LogSafety_SAFE.New()
	case "UNSAFE":
		*e = LogSafety_UNSAFE.New()
	case "DO_NOT_LOG":
		*e = LogSafety_DO_NOT_LOG.New()
	}
	return nil
}

type PrimitiveType struct {
	val PrimitiveType_Value
}

type PrimitiveType_Value string

const (
	PrimitiveType_STRING      PrimitiveType_Value = "STRING"
	PrimitiveType_DATETIME    PrimitiveType_Value = "DATETIME"
	PrimitiveType_INTEGER     PrimitiveType_Value = "INTEGER"
	PrimitiveType_DOUBLE      PrimitiveType_Value = "DOUBLE"
	PrimitiveType_SAFELONG    PrimitiveType_Value = "SAFELONG"
	PrimitiveType_BINARY      PrimitiveType_Value = "BINARY"
	PrimitiveType_ANY         PrimitiveType_Value = "ANY"
	PrimitiveType_BOOLEAN     PrimitiveType_Value = "BOOLEAN"
	PrimitiveType_UUID        PrimitiveType_Value = "UUID"
	PrimitiveType_RID         PrimitiveType_Value = "RID"
	PrimitiveType_BEARERTOKEN PrimitiveType_Value = "BEARERTOKEN"
	PrimitiveType_UNKNOWN     PrimitiveType_Value = "UNKNOWN"
)

func (e PrimitiveType_Value) New() PrimitiveType {
	return PrimitiveType{val: e}
}

// PrimitiveType_Values returns all known variants of PrimitiveType.
func PrimitiveType_Values() []PrimitiveType_Value {
	return []PrimitiveType_Value{PrimitiveType_STRING, PrimitiveType_DATETIME, PrimitiveType_INTEGER, PrimitiveType_DOUBLE, PrimitiveType_SAFELONG, PrimitiveType_BINARY, PrimitiveType_ANY, PrimitiveType_BOOLEAN, PrimitiveType_UUID, PrimitiveType_RID, PrimitiveType_BEARERTOKEN}
}

func New_PrimitiveType(value PrimitiveType_Value) PrimitiveType {
	return PrimitiveType{val: value}
}

// IsUnknown returns false for all known variants of PrimitiveType and true otherwise.
func (e PrimitiveType) IsUnknown() bool {
	switch e.val {
	case PrimitiveType_STRING, PrimitiveType_DATETIME, PrimitiveType_INTEGER, PrimitiveType_DOUBLE, PrimitiveType_SAFELONG, PrimitiveType_BINARY, PrimitiveType_ANY, PrimitiveType_BOOLEAN, PrimitiveType_UUID, PrimitiveType_RID, PrimitiveType_BEARERTOKEN:
		return false
	}
	return true
}

func (e PrimitiveType) Value() PrimitiveType_Value {
	if e.IsUnknown() {
		return PrimitiveType_UNKNOWN
	}
	return e.val
}

func (e PrimitiveType) String() string {
	return string(e.val)
}

func (e PrimitiveType) MarshalText() ([]byte, error) {
	return []byte(e.val), nil
}

func (e *PrimitiveType) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		*e = PrimitiveType_Value(v).New()
	case "STRING":
		*e = PrimitiveType_STRING.New()
	case "DATETIME":
		*e = PrimitiveType_DATETIME.New()
	case "INTEGER":
		*e = PrimitiveType_INTEGER.New()
	case "DOUBLE":
		*e = PrimitiveType_DOUBLE.New()
	case "SAFELONG":
		*e = PrimitiveType_SAFELONG.New()
	case "BINARY":
		*e = PrimitiveType_BINARY.New()
	case "ANY":
		*e = PrimitiveType_ANY.New()
	case "BOOLEAN":
		*e = PrimitiveType_BOOLEAN.New()
	case "UUID":
		*e = PrimitiveType_UUID.New()
	case "RID":
		*e = PrimitiveType_RID.New()
	case "BEARERTOKEN":
		*e = PrimitiveType_BEARERTOKEN.New()
	}
	return nil
}
