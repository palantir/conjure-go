// Copyright (c) 2018 Palantir Technologies. All rights reserved.
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

package conjure

import (
	"testing"

	"github.com/palantir/goastwriter"
	"github.com/palantir/goastwriter/astgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/v5/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v5/conjure/types"
)

func TestEnum(t *testing.T) {
	for caseNum, currCase := range []struct {
		pkg   string
		name  string
		enums []spec.EnumDefinition
		want  string
	}{
		{
			pkg:  "testpkg",
			name: "single enum",
			enums: []spec.EnumDefinition{
				{
					TypeName: spec.TypeName{
						Name:    "Months",
						Package: "api",
					},
					Docs: docPtr("These represent months"),
					Values: []spec.EnumValueDefinition{
						{Value: "JANUARY"},
						{Value: "FEBRUARY"},
					},
				},
			},
			want: `package testpkg

// These represent months
type Months struct {
	val MonthsValue
}

// These represent months
type MonthsValue string

const (
	MonthsJanuary  MonthsValue = "JANUARY"
	MonthsFebruary MonthsValue = "FEBRUARY"
	MonthsUnknown  MonthsValue = "UNKNOWN"
)

// Months_Values returns all known variants of Months.
func Months_Values() []MonthsValue {
	return []MonthsValue{MonthsJanuary, MonthsFebruary}
}
func NewMonths(value MonthsValue) Months {
	return Months{val: value}
}

// IsUnknown returns false for all known variants of Months and true otherwise.
func (e Months) IsUnknown() bool {
	switch e.val {
	case MonthsJanuary, MonthsFebruary:
		return false
	}
	return true
}
func (e Months) Value() MonthsValue {
	if e.IsUnknown() {
		return MonthsUnknown
	}
	return e.val
}
func (e Months) String() string {
	return string(e.val)
}
func (e Months) MarshalText() ([]byte, error) {
	return []byte(e.val), nil
}
func (e *Months) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		if !enumValuePattern.MatchString(v) {
			return werror.Convert(errors.NewInvalidArgument(wparams.NewSafeAndUnsafeParamStorer(map[string]interface{}{"enumType": "Months", "message": "enum value must match pattern ^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$"}, map[string]interface{}{"enumValue": string(data)})))
		}
		*e = NewMonths(MonthsValue(v))
	case "JANUARY":
		*e = NewMonths(MonthsJanuary)
	case "FEBRUARY":
		*e = NewMonths(MonthsFebruary)
	}
	return nil
}
`,
		},
		{
			pkg:  "testpkg",
			name: "multiple enums",
			enums: []spec.EnumDefinition{
				{
					TypeName: spec.TypeName{
						Name:    "Months",
						Package: "api",
					},
					Docs: docPtr("These represent months"),
					Values: []spec.EnumValueDefinition{
						{Value: "JANUARY"},
						{Value: "FEBRUARY"},
					},
				},
				{

					TypeName: spec.TypeName{
						Name:    "Values",
						Package: "api",
					},
					Docs: docPtr("These represent values"),
					Values: []spec.EnumValueDefinition{
						{Value: "NULL_VALUE"},
						{Value: "VALID_VALUE"},
					},
				},
			},
			want: `package testpkg

// These represent months
type Months struct {
	val MonthsValue
}

// These represent months
type MonthsValue string

const (
	MonthsJanuary  MonthsValue = "JANUARY"
	MonthsFebruary MonthsValue = "FEBRUARY"
	MonthsUnknown  MonthsValue = "UNKNOWN"
)

// Months_Values returns all known variants of Months.
func Months_Values() []MonthsValue {
	return []MonthsValue{MonthsJanuary, MonthsFebruary}
}
func NewMonths(value MonthsValue) Months {
	return Months{val: value}
}

// IsUnknown returns false for all known variants of Months and true otherwise.
func (e Months) IsUnknown() bool {
	switch e.val {
	case MonthsJanuary, MonthsFebruary:
		return false
	}
	return true
}
func (e Months) Value() MonthsValue {
	if e.IsUnknown() {
		return MonthsUnknown
	}
	return e.val
}
func (e Months) String() string {
	return string(e.val)
}
func (e Months) MarshalText() ([]byte, error) {
	return []byte(e.val), nil
}
func (e *Months) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		if !enumValuePattern.MatchString(v) {
			return werror.Convert(errors.NewInvalidArgument(wparams.NewSafeAndUnsafeParamStorer(map[string]interface{}{"enumType": "Months", "message": "enum value must match pattern ^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$"}, map[string]interface{}{"enumValue": string(data)})))
		}
		*e = NewMonths(MonthsValue(v))
	case "JANUARY":
		*e = NewMonths(MonthsJanuary)
	case "FEBRUARY":
		*e = NewMonths(MonthsFebruary)
	}
	return nil
}

// These represent values
type Values struct {
	val ValuesValue
}

// These represent values
type ValuesValue string

const (
	ValuesNullValue  ValuesValue = "NULL_VALUE"
	ValuesValidValue ValuesValue = "VALID_VALUE"
	ValuesUnknown    ValuesValue = "UNKNOWN"
)

// Values_Values returns all known variants of Values.
func Values_Values() []ValuesValue {
	return []ValuesValue{ValuesNullValue, ValuesValidValue}
}
func NewValues(value ValuesValue) Values {
	return Values{val: value}
}

// IsUnknown returns false for all known variants of Values and true otherwise.
func (e Values) IsUnknown() bool {
	switch e.val {
	case ValuesNullValue, ValuesValidValue:
		return false
	}
	return true
}
func (e Values) Value() ValuesValue {
	if e.IsUnknown() {
		return ValuesUnknown
	}
	return e.val
}
func (e Values) String() string {
	return string(e.val)
}
func (e Values) MarshalText() ([]byte, error) {
	return []byte(e.val), nil
}
func (e *Values) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		if !enumValuePattern.MatchString(v) {
			return werror.Convert(errors.NewInvalidArgument(wparams.NewSafeAndUnsafeParamStorer(map[string]interface{}{"enumType": "Values", "message": "enum value must match pattern ^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$"}, map[string]interface{}{"enumValue": string(data)})))
		}
		*e = NewValues(ValuesValue(v))
	case "NULL_VALUE":
		*e = NewValues(ValuesNullValue)
	case "VALID_VALUE":
		*e = NewValues(ValuesValidValue)
	}
	return nil
}
`,
		},
		{
			pkg:  "testpkg",
			name: "enum with comments",
			enums: []spec.EnumDefinition{
				{
					TypeName: spec.TypeName{
						Name:    "Months",
						Package: "api",
					},
					Docs: docPtr("These represent months"),
					Values: []spec.EnumValueDefinition{
						{
							Value: "JANUARY",
							Docs:  docPtr("Docs for JANUARY"),
						},
						{
							Value: "FEBRUARY",
							Docs:  docPtr("Docs for FEBRUARY"),
						},
					},
				},
			},
			want: `package testpkg

// These represent months
type Months struct {
	val MonthsValue
}

// These represent months
type MonthsValue string

const (
	// Docs for JANUARY
	MonthsJanuary MonthsValue = "JANUARY"
	// Docs for FEBRUARY
	MonthsFebruary MonthsValue = "FEBRUARY"
	MonthsUnknown  MonthsValue = "UNKNOWN"
)

// Months_Values returns all known variants of Months.
func Months_Values() []MonthsValue {
	return []MonthsValue{MonthsJanuary, MonthsFebruary}
}
func NewMonths(value MonthsValue) Months {
	return Months{val: value}
}

// IsUnknown returns false for all known variants of Months and true otherwise.
func (e Months) IsUnknown() bool {
	switch e.val {
	case MonthsJanuary, MonthsFebruary:
		return false
	}
	return true
}
func (e Months) Value() MonthsValue {
	if e.IsUnknown() {
		return MonthsUnknown
	}
	return e.val
}
func (e Months) String() string {
	return string(e.val)
}
func (e Months) MarshalText() ([]byte, error) {
	return []byte(e.val), nil
}
func (e *Months) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		if !enumValuePattern.MatchString(v) {
			return werror.Convert(errors.NewInvalidArgument(wparams.NewSafeAndUnsafeParamStorer(map[string]interface{}{"enumType": "Months", "message": "enum value must match pattern ^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$"}, map[string]interface{}{"enumValue": string(data)})))
		}
		*e = NewMonths(MonthsValue(v))
	case "JANUARY":
		*e = NewMonths(MonthsJanuary)
	case "FEBRUARY":
		*e = NewMonths(MonthsFebruary)
	}
	return nil
}
`,
		},
	} {
		t.Run(currCase.name, func(t *testing.T) {
			info := types.NewPkgInfo("foo", nil)
			var components []astgen.ASTDecl
			for _, e := range currCase.enums {
				declers := astForEnum(e, info)
				components = append(components, declers...)
			}

			got, err := goastwriter.Write(currCase.pkg, components...)
			require.NoError(t, err, "Case %d: %s", caseNum, currCase.name)

			assert.Equal(t, currCase.want, string(got))
		})
	}
}

func docPtr(doc string) *spec.Documentation {
	return (*spec.Documentation)(&doc)
}
