// This file was generated by Conjure and should not be manually edited.

package api

import (
	"encoding/json"
	"fmt"

	"github.com/palantir/conjure-go-runtime/conjure-go-contract/errors"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
	"github.com/palantir/pkg/uuid"
)

type myNotFound struct {
	// This is safeArgA doc.
	SafeArgA Basic `json:"safeArgA" conjure-docs:"This is safeArgA doc."`
	// This is safeArgB doc.
	SafeArgB   []int   `json:"safeArgB" conjure-docs:"This is safeArgB doc."`
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
func NewMyNotFound(safeArgA Basic, safeArgB []int, unsafeArgA string, unsafeArgB *string) *MyNotFound {
	return &MyNotFound{errorInstanceID: uuid.NewUUID(), myNotFound: myNotFound{SafeArgA: safeArgA, SafeArgB: safeArgB, UnsafeArgA: unsafeArgA, UnsafeArgB: unsafeArgB}}
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
	return map[string]interface{}{"safeArgA": e.SafeArgA, "safeArgB": e.SafeArgB, "unsafeArgA": e.UnsafeArgA, "unsafeArgB": e.UnsafeArgB}
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
