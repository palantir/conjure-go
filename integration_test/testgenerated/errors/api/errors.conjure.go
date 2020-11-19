// This file was generated by Conjure and should not be manually edited.

package api

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
	"github.com/palantir/pkg/uuid"
	werror "github.com/palantir/witchcraft-go-error"
)

type myInternal struct {
	// This is safeArgA doc.
	SafeArgA Basic `json:"safeArgA" conjure-docs:"This is safeArgA doc."`
	// This is safeArgB doc.
	SafeArgB []int `json:"safeArgB" conjure-docs:"This is safeArgB doc."`
	// A field named with a go keyword
	Type       string  `json:"type" conjure-docs:"A field named with a go keyword"`
	UnsafeArgA string  `json:"unsafeArgA"`
	UnsafeArgB *string `json:"unsafeArgB"`
	MyInternal string  `json:"myInternal"`
}

func (o myInternal) MarshalJSON() ([]byte, error) {
	if o.SafeArgB == nil {
		o.SafeArgB = make([]int, 0)
	}
	type myInternalAlias myInternal
	return safejson.Marshal(myInternalAlias(o))
}

func (o *myInternal) UnmarshalJSON(data []byte) error {
	type myInternalAlias myInternal
	var rawmyInternal myInternalAlias
	if err := safejson.Unmarshal(data, &rawmyInternal); err != nil {
		return err
	}
	if rawmyInternal.SafeArgB == nil {
		rawmyInternal.SafeArgB = make([]int, 0)
	}
	*o = myInternal(rawmyInternal)
	return nil
}

func (o myInternal) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *myInternal) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

// NewMyInternal returns new instance of MyInternal error.
func NewMyInternal(safeArgAArg Basic, safeArgBArg []int, typeArg string, unsafeArgAArg string, unsafeArgBArg *string, myInternalArg string) *MyInternal {
	return &MyInternal{errorInstanceID: uuid.NewUUID(), stack: werror.NewStackTrace(), myInternal: myInternal{SafeArgA: safeArgAArg, SafeArgB: safeArgBArg, Type: typeArg, UnsafeArgA: unsafeArgAArg, UnsafeArgB: unsafeArgBArg, MyInternal: myInternalArg}}
}

// WrapWithMyInternal returns new instance of MyInternal error wrapping an existing error.
func WrapWithMyInternal(err error, safeArgAArg Basic, safeArgBArg []int, typeArg string, unsafeArgAArg string, unsafeArgBArg *string, myInternalArg string) *MyInternal {
	return &MyInternal{errorInstanceID: uuid.NewUUID(), stack: werror.NewStackTrace(), cause: err, myInternal: myInternal{SafeArgA: safeArgAArg, SafeArgB: safeArgBArg, Type: typeArg, UnsafeArgA: unsafeArgAArg, UnsafeArgB: unsafeArgBArg, MyInternal: myInternalArg}}
}

// MyInternal is an error type.
//
// Internal server error.
type MyInternal struct {
	errorInstanceID uuid.UUID
	myInternal
	cause error
	stack werror.StackTrace
}

// IsMyInternal returns true if err is an instance of MyInternal.
func IsMyInternal(err error) bool {
	if err == nil {
		return false
	}
	_, ok := errors.GetConjureError(err).(*MyInternal)
	return ok
}

func (e *MyInternal) Error() string {
	return fmt.Sprintf("INTERNAL MyNamespace:MyInternal (%s)", e.errorInstanceID)
}

// Cause returns the underlying cause of the error, or nil if none.
// Note that cause is not serialized and sent over the wire.
func (e *MyInternal) Cause() error {
	return e.cause
}

// StackTrace returns the StackTrace for the error, or nil if none.
// Note that stack traces are not serialized and sent over the wire.
func (e *MyInternal) StackTrace() werror.StackTrace {
	return e.stack
}

// Message returns the message body for the error.
func (e *MyInternal) Message() string {
	return "INTERNAL MyNamespace:MyInternal"
}

// Format implements fmt.Formatter, a requirement of werror.Werror.
func (e *MyInternal) Format(state fmt.State, verb rune) {
	werror.Format(e, e.safeParams(), state, verb)
}

// Code returns an enum describing error category.
func (e *MyInternal) Code() errors.ErrorCode {
	return errors.Internal
}

// Name returns an error name identifying error type.
func (e *MyInternal) Name() string {
	return "MyNamespace:MyInternal"
}

// InstanceID returns unique identifier of this particular error instance.
func (e *MyInternal) InstanceID() uuid.UUID {
	return e.errorInstanceID
}

// Parameters returns a set of named parameters detailing this particular error instance.
func (e *MyInternal) Parameters() map[string]interface{} {
	return map[string]interface{}{"safeArgA": e.SafeArgA, "safeArgB": e.SafeArgB, "type": e.Type, "unsafeArgA": e.UnsafeArgA, "unsafeArgB": e.UnsafeArgB, "myInternal": e.MyInternal}
}

// safeParams returns a set of named safe parameters detailing this particular error instance.
func (e *MyInternal) safeParams() map[string]interface{} {
	return map[string]interface{}{"safeArgA": e.SafeArgA, "safeArgB": e.SafeArgB, "type": e.Type, "errorInstanceId": e.errorInstanceID}
}

// SafeParams returns a set of named safe parameters detailing this particular error instance and
// any underlying causes.
func (e *MyInternal) SafeParams() map[string]interface{} {
	safeParams, _ := werror.ParamsFromError(e.cause)
	for k, v := range e.safeParams() {
		if _, exists := safeParams[k]; !exists {
			safeParams[k] = v
		}
	}
	return safeParams
}

// unsafeParams returns a set of named unsafe parameters detailing this particular error instance.
func (e *MyInternal) unsafeParams() map[string]interface{} {
	return map[string]interface{}{"unsafeArgA": e.UnsafeArgA, "unsafeArgB": e.UnsafeArgB, "myInternal": e.MyInternal}
}

// UnsafeParams returns a set of named unsafe parameters detailing this particular error instance and
// any underlying causes.
func (e *MyInternal) UnsafeParams() map[string]interface{} {
	_, unsafeParams := werror.ParamsFromError(e.cause)
	for k, v := range e.unsafeParams() {
		if _, exists := unsafeParams[k]; !exists {
			unsafeParams[k] = v
		}
	}
	return unsafeParams
}

func (e MyInternal) MarshalJSON() ([]byte, error) {
	parameters, err := safejson.Marshal(e.myInternal)
	if err != nil {
		return nil, err
	}
	return safejson.Marshal(errors.SerializableError{ErrorCode: errors.Internal, ErrorName: "MyNamespace:MyInternal", ErrorInstanceID: e.errorInstanceID, Parameters: json.RawMessage(parameters)})
}

func (e *MyInternal) UnmarshalJSON(data []byte) error {
	var serializableError errors.SerializableError
	if err := safejson.Unmarshal(data, &serializableError); err != nil {
		return err
	}
	var parameters myInternal
	if err := safejson.Unmarshal([]byte(serializableError.Parameters), &parameters); err != nil {
		return err
	}
	e.errorInstanceID = serializableError.ErrorInstanceID
	e.myInternal = parameters
	return nil
}

type myNotFound struct {
	// This is safeArgA doc.
	SafeArgA Basic `json:"safeArgA" conjure-docs:"This is safeArgA doc."`
	// This is safeArgB doc.
	SafeArgB []int `json:"safeArgB" conjure-docs:"This is safeArgB doc."`
	// A field named with a go keyword
	Type       string  `json:"type" conjure-docs:"A field named with a go keyword"`
	UnsafeArgA string  `json:"unsafeArgA"`
	UnsafeArgB *string `json:"unsafeArgB"`
}

func (o myNotFound) MarshalJSON() ([]byte, error) {
	if o.SafeArgB == nil {
		o.SafeArgB = make([]int, 0)
	}
	type myNotFoundAlias myNotFound
	return safejson.Marshal(myNotFoundAlias(o))
}

func (o *myNotFound) UnmarshalJSON(data []byte) error {
	type myNotFoundAlias myNotFound
	var rawmyNotFound myNotFoundAlias
	if err := safejson.Unmarshal(data, &rawmyNotFound); err != nil {
		return err
	}
	if rawmyNotFound.SafeArgB == nil {
		rawmyNotFound.SafeArgB = make([]int, 0)
	}
	*o = myNotFound(rawmyNotFound)
	return nil
}

func (o myNotFound) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *myNotFound) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

// NewMyNotFound returns new instance of MyNotFound error.
func NewMyNotFound(safeArgAArg Basic, safeArgBArg []int, typeArg string, unsafeArgAArg string, unsafeArgBArg *string) *MyNotFound {
	return &MyNotFound{errorInstanceID: uuid.NewUUID(), stack: werror.NewStackTrace(), myNotFound: myNotFound{SafeArgA: safeArgAArg, SafeArgB: safeArgBArg, Type: typeArg, UnsafeArgA: unsafeArgAArg, UnsafeArgB: unsafeArgBArg}}
}

// WrapWithMyNotFound returns new instance of MyNotFound error wrapping an existing error.
func WrapWithMyNotFound(err error, safeArgAArg Basic, safeArgBArg []int, typeArg string, unsafeArgAArg string, unsafeArgBArg *string) *MyNotFound {
	return &MyNotFound{errorInstanceID: uuid.NewUUID(), stack: werror.NewStackTrace(), cause: err, myNotFound: myNotFound{SafeArgA: safeArgAArg, SafeArgB: safeArgBArg, Type: typeArg, UnsafeArgA: unsafeArgAArg, UnsafeArgB: unsafeArgBArg}}
}

// MyNotFound is an error type.
//
// Something was not found.
type MyNotFound struct {
	errorInstanceID uuid.UUID
	myNotFound
	cause error
	stack werror.StackTrace
}

// IsMyNotFound returns true if err is an instance of MyNotFound.
func IsMyNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := errors.GetConjureError(err).(*MyNotFound)
	return ok
}

func (e *MyNotFound) Error() string {
	return fmt.Sprintf("NOT_FOUND MyNamespace:MyNotFound (%s)", e.errorInstanceID)
}

// Cause returns the underlying cause of the error, or nil if none.
// Note that cause is not serialized and sent over the wire.
func (e *MyNotFound) Cause() error {
	return e.cause
}

// StackTrace returns the StackTrace for the error, or nil if none.
// Note that stack traces are not serialized and sent over the wire.
func (e *MyNotFound) StackTrace() werror.StackTrace {
	return e.stack
}

// Message returns the message body for the error.
func (e *MyNotFound) Message() string {
	return "NOT_FOUND MyNamespace:MyNotFound"
}

// Format implements fmt.Formatter, a requirement of werror.Werror.
func (e *MyNotFound) Format(state fmt.State, verb rune) {
	werror.Format(e, e.safeParams(), state, verb)
}

// Code returns an enum describing error category.
func (e *MyNotFound) Code() errors.ErrorCode {
	return errors.NotFound
}

// Name returns an error name identifying error type.
func (e *MyNotFound) Name() string {
	return "MyNamespace:MyNotFound"
}

// InstanceID returns unique identifier of this particular error instance.
func (e *MyNotFound) InstanceID() uuid.UUID {
	return e.errorInstanceID
}

// Parameters returns a set of named parameters detailing this particular error instance.
func (e *MyNotFound) Parameters() map[string]interface{} {
	return map[string]interface{}{"safeArgA": e.SafeArgA, "safeArgB": e.SafeArgB, "type": e.Type, "unsafeArgA": e.UnsafeArgA, "unsafeArgB": e.UnsafeArgB}
}

// safeParams returns a set of named safe parameters detailing this particular error instance.
func (e *MyNotFound) safeParams() map[string]interface{} {
	return map[string]interface{}{"safeArgA": e.SafeArgA, "safeArgB": e.SafeArgB, "type": e.Type, "errorInstanceId": e.errorInstanceID}
}

// SafeParams returns a set of named safe parameters detailing this particular error instance and
// any underlying causes.
func (e *MyNotFound) SafeParams() map[string]interface{} {
	safeParams, _ := werror.ParamsFromError(e.cause)
	for k, v := range e.safeParams() {
		if _, exists := safeParams[k]; !exists {
			safeParams[k] = v
		}
	}
	return safeParams
}

// unsafeParams returns a set of named unsafe parameters detailing this particular error instance.
func (e *MyNotFound) unsafeParams() map[string]interface{} {
	return map[string]interface{}{"unsafeArgA": e.UnsafeArgA, "unsafeArgB": e.UnsafeArgB}
}

// UnsafeParams returns a set of named unsafe parameters detailing this particular error instance and
// any underlying causes.
func (e *MyNotFound) UnsafeParams() map[string]interface{} {
	_, unsafeParams := werror.ParamsFromError(e.cause)
	for k, v := range e.unsafeParams() {
		if _, exists := unsafeParams[k]; !exists {
			unsafeParams[k] = v
		}
	}
	return unsafeParams
}

func (e MyNotFound) MarshalJSON() ([]byte, error) {
	parameters, err := safejson.Marshal(e.myNotFound)
	if err != nil {
		return nil, err
	}
	return safejson.Marshal(errors.SerializableError{ErrorCode: errors.NotFound, ErrorName: "MyNamespace:MyNotFound", ErrorInstanceID: e.errorInstanceID, Parameters: json.RawMessage(parameters)})
}

func (e *MyNotFound) UnmarshalJSON(data []byte) error {
	var serializableError errors.SerializableError
	if err := safejson.Unmarshal(data, &serializableError); err != nil {
		return err
	}
	var parameters myNotFound
	if err := safejson.Unmarshal([]byte(serializableError.Parameters), &parameters); err != nil {
		return err
	}
	e.errorInstanceID = serializableError.ErrorInstanceID
	e.myNotFound = parameters
	return nil
}

func init() {
	errors.RegisterErrorType("MyNamespace:MyInternal", reflect.TypeOf(MyInternal{}))
	errors.RegisterErrorType("MyNamespace:MyNotFound", reflect.TypeOf(MyNotFound{}))
}
