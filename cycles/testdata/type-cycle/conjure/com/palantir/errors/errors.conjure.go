// This file was generated by Conjure and should not be manually edited.

package errors

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	"github.com/palantir/conjure-go/v6/cycles/testdata/type-cycle/conjure/com/palantir/bar"
	barfoo "github.com/palantir/conjure-go/v6/cycles/testdata/type-cycle/conjure/com/palantir/bar_foo"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
	"github.com/palantir/pkg/uuid"
	werror "github.com/palantir/witchcraft-go-error"
)

type myError struct {
	SafeArg1   barfoo.Type1    `json:"safeArg1"`
	SafeArg2   bar.Type2       `json:"safeArg2"`
	UnsafeArg3 barfoo.BarType3 `json:"unsafeArg3"`
}

func (o myError) MarshalJSON() ([]byte, error) {
	if o.SafeArg1 == nil {
		o.SafeArg1 = make([]barfoo.BarType3, 0)
	}
	type myErrorAlias myError
	return safejson.Marshal(myErrorAlias(o))
}

func (o *myError) UnmarshalJSON(data []byte) error {
	type myErrorAlias myError
	var rawmyError myErrorAlias
	if err := safejson.Unmarshal(data, &rawmyError); err != nil {
		return err
	}
	if rawmyError.SafeArg1 == nil {
		rawmyError.SafeArg1 = make([]barfoo.BarType3, 0)
	}
	*o = myError(rawmyError)
	return nil
}

func (o myError) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *myError) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

// NewMyError returns new instance of MyError error.
func NewMyError(safeArg1Arg barfoo.Type1, safeArg2Arg bar.Type2, unsafeArg3Arg barfoo.BarType3) *MyError {
	return &MyError{errorInstanceID: uuid.NewUUID(), stack: werror.NewStackTrace(), myError: myError{SafeArg1: safeArg1Arg, SafeArg2: safeArg2Arg, UnsafeArg3: unsafeArg3Arg}}
}

// WrapWithMyError returns new instance of MyError error wrapping an existing error.
func WrapWithMyError(err error, safeArg1Arg barfoo.Type1, safeArg2Arg bar.Type2, unsafeArg3Arg barfoo.BarType3) *MyError {
	return &MyError{errorInstanceID: uuid.NewUUID(), stack: werror.NewStackTrace(), cause: err, myError: myError{SafeArg1: safeArg1Arg, SafeArg2: safeArg2Arg, UnsafeArg3: unsafeArg3Arg}}
}

// MyError is an error type.
type MyError struct {
	errorInstanceID uuid.UUID
	myError
	cause error
	stack werror.StackTrace
}

// IsMyError returns true if err is an instance of MyError.
func IsMyError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := errors.GetConjureError(err).(*MyError)
	return ok
}

func (e *MyError) Error() string {
	return fmt.Sprintf("INTERNAL Namespace:MyError (%s)", e.errorInstanceID)
}

// Cause returns the underlying cause of the error, or nil if none.
// Note that cause is not serialized and sent over the wire.
func (e *MyError) Cause() error {
	return e.cause
}

// StackTrace returns the StackTrace for the error, or nil if none.
// Note that stack traces are not serialized and sent over the wire.
func (e *MyError) StackTrace() werror.StackTrace {
	return e.stack
}

// Message returns the message body for the error.
func (e *MyError) Message() string {
	return "INTERNAL Namespace:MyError"
}

// Format implements fmt.Formatter, a requirement of werror.Werror.
func (e *MyError) Format(state fmt.State, verb rune) {
	werror.Format(e, e.safeParams(), state, verb)
}

// Code returns an enum describing error category.
func (e *MyError) Code() errors.ErrorCode {
	return errors.Internal
}

// Name returns an error name identifying error type.
func (e *MyError) Name() string {
	return "Namespace:MyError"
}

// InstanceID returns unique identifier of this particular error instance.
func (e *MyError) InstanceID() uuid.UUID {
	return e.errorInstanceID
}

// Parameters returns a set of named parameters detailing this particular error instance.
func (e *MyError) Parameters() map[string]interface{} {
	return map[string]interface{}{"safeArg1": e.SafeArg1, "safeArg2": e.SafeArg2, "unsafeArg3": e.UnsafeArg3}
}

// safeParams returns a set of named safe parameters detailing this particular error instance.
func (e *MyError) safeParams() map[string]interface{} {
	return map[string]interface{}{"safeArg1": e.SafeArg1, "safeArg2": e.SafeArg2, "errorInstanceId": e.errorInstanceID, "errorName": e.Name()}
}

// SafeParams returns a set of named safe parameters detailing this particular error instance and
// any underlying causes.
func (e *MyError) SafeParams() map[string]interface{} {
	safeParams, _ := werror.ParamsFromError(e.cause)
	for k, v := range e.safeParams() {
		if _, exists := safeParams[k]; !exists {
			safeParams[k] = v
		}
	}
	return safeParams
}

// unsafeParams returns a set of named unsafe parameters detailing this particular error instance.
func (e *MyError) unsafeParams() map[string]interface{} {
	return map[string]interface{}{"unsafeArg3": e.UnsafeArg3}
}

// UnsafeParams returns a set of named unsafe parameters detailing this particular error instance and
// any underlying causes.
func (e *MyError) UnsafeParams() map[string]interface{} {
	_, unsafeParams := werror.ParamsFromError(e.cause)
	for k, v := range e.unsafeParams() {
		if _, exists := unsafeParams[k]; !exists {
			unsafeParams[k] = v
		}
	}
	return unsafeParams
}

func (e MyError) MarshalJSON() ([]byte, error) {
	parameters, err := safejson.Marshal(e.myError)
	if err != nil {
		return nil, err
	}
	return safejson.Marshal(errors.SerializableError{ErrorCode: errors.Internal, ErrorName: "Namespace:MyError", ErrorInstanceID: e.errorInstanceID, Parameters: json.RawMessage(parameters)})
}

func (e *MyError) UnmarshalJSON(data []byte) error {
	var serializableError errors.SerializableError
	if err := safejson.Unmarshal(data, &serializableError); err != nil {
		return err
	}
	var parameters myError
	if err := safejson.Unmarshal([]byte(serializableError.Parameters), &parameters); err != nil {
		return err
	}
	e.errorInstanceID = serializableError.ErrorInstanceID
	e.myError = parameters
	return nil
}

func init() {
	errors.RegisterErrorType("Namespace:MyError", reflect.TypeOf(MyError{}))
}
