// This file was generated by Conjure and should not be manually edited.

package types

import (
	"regexp"
	"strings"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	werror "github.com/palantir/witchcraft-go-error"
	wparams "github.com/palantir/witchcraft-go-params"
)

type EnumExample string

const (
	EnumExampleOne EnumExample = "ONE"
	EnumExampleTwo EnumExample = "TWO"
)

func (e *EnumExample) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		if !regexp.MustCompile("^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$").MatchString(v) {
			return werror.Convert(errors.NewInvalidArgument(wparams.NewSafeAndUnsafeParamStorer(map[string]interface{}{"enumType": "EnumExample", "message": "enum value must match pattern [A-Z][A-Z0-9]*(_[A-Z0-9]+)*"}, map[string]interface{}{"enumValue": string(data)})))
		}
		*e = EnumExample(v)
	case "ONE":
		*e = EnumExampleOne
	case "TWO":
		*e = EnumExampleTwo
	}
	return nil
}

type Enum string

const (
	EnumOne Enum = "ONE"
	EnumTwo Enum = "TWO"
)

func (e *Enum) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		if !regexp.MustCompile("^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$").MatchString(v) {
			return werror.Convert(errors.NewInvalidArgument(wparams.NewSafeAndUnsafeParamStorer(map[string]interface{}{"enumType": "Enum", "message": "enum value must match pattern [A-Z][A-Z0-9]*(_[A-Z0-9]+)*"}, map[string]interface{}{"enumValue": string(data)})))
		}
		*e = Enum(v)
	case "ONE":
		*e = EnumOne
	case "TWO":
		*e = EnumTwo
	}
	return nil
}
