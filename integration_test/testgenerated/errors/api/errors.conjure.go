// This file was generated by Conjure and should not be manually edited.

package api

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/palantir/conjure-go-runtime/conjure-go-contract/errors"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
	"github.com/palantir/pkg/uuid"
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
func NewMyInternal(safeArgA Basic, safeArgB []int, type_ string, unsafeArgA string, unsafeArgB *string) *MyInternal {
	return &MyInternal{errorInstanceID: uuid.NewUUID(), myInternal: myInternal{SafeArgA: safeArgA, SafeArgB: safeArgB, Type: type_, UnsafeArgA: unsafeArgA, UnsafeArgB: unsafeArgB}}
}

// MyInternal is an error type.
//
// Internal server error.
type MyInternal struct {
	errorInstanceID uuid.UUID
	myInternal
}

func (e *MyInternal) Error() string {
	return fmt.Sprintf("INTERNAL MyNamespace:MyInternal (%s)", e.errorInstanceID)
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
	return map[string]interface{}{"safeArgA": e.SafeArgA, "safeArgB": e.SafeArgB, "type": e.Type, "unsafeArgA": e.UnsafeArgA, "unsafeArgB": e.UnsafeArgB}
}

// SafeParams returns a set of named safe parameters detailing this particular error instance.
func (e *MyInternal) SafeParams() map[string]interface{} {
	return map[string]interface{}{"safeArgA": e.SafeArgA, "safeArgB": e.SafeArgB, "type": e.Type, "errorInstanceId": e.errorInstanceID}
}

// UnsafeParams returns a set of named unsafe parameters detailing this particular error instance.
func (e *MyInternal) UnsafeParams() map[string]interface{} {
	return map[string]interface{}{"unsafeArgA": e.UnsafeArgA, "unsafeArgB": e.UnsafeArgB}
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
func NewMyNotFound(safeArgA Basic, safeArgB []int, type_ string, unsafeArgA string, unsafeArgB *string) *MyNotFound {
	return &MyNotFound{errorInstanceID: uuid.NewUUID(), myNotFound: myNotFound{SafeArgA: safeArgA, SafeArgB: safeArgB, Type: type_, UnsafeArgA: unsafeArgA, UnsafeArgB: unsafeArgB}}
}

// MyNotFound is an error type.
//
// Something was not found.
type MyNotFound struct {
	errorInstanceID uuid.UUID
	myNotFound
}

func (e *MyNotFound) Error() string {
	return fmt.Sprintf("NOT_FOUND MyNamespace:MyNotFound (%s)", e.errorInstanceID)
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

// SafeParams returns a set of named safe parameters detailing this particular error instance.
func (e *MyNotFound) SafeParams() map[string]interface{} {
	return map[string]interface{}{"safeArgA": e.SafeArgA, "safeArgB": e.SafeArgB, "type": e.Type, "errorInstanceId": e.errorInstanceID}
}

// UnsafeParams returns a set of named unsafe parameters detailing this particular error instance.
func (e *MyNotFound) UnsafeParams() map[string]interface{} {
	return map[string]interface{}{"unsafeArgA": e.UnsafeArgA, "unsafeArgB": e.UnsafeArgB}
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
