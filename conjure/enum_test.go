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
	"strings"
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
type Months string

const (
	MonthsJanuary  Months = "JANUARY"
	MonthsFebruary Months = "FEBRUARY"
)

// Months_Values returns all known variants of Months.
func Months_Values() []Months {
	return []Months{MonthsJanuary, MonthsFebruary}
}

// IsUnknown returns false for all known variants of Months and true otherwise.
func (e Months) IsUnknown() bool {
	switch e {
	case MonthsJanuary, MonthsFebruary:
		return false
	}
	return true
}
func (e *Months) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		if !enumValuePattern.MatchString(v) {
			return werror.Convert(errors.NewInvalidArgument(wparams.NewSafeAndUnsafeParamStorer(map[string]interface{}{"enumType": "Months", "message": "enum value must match pattern ^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$"}, map[string]interface{}{"enumValue": string(data)})))
		}
		*e = Months(v)
	case "JANUARY":
		*e = MonthsJanuary
	case "FEBRUARY":
		*e = MonthsFebruary
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
type Months string

const (
	MonthsJanuary  Months = "JANUARY"
	MonthsFebruary Months = "FEBRUARY"
)

// Months_Values returns all known variants of Months.
func Months_Values() []Months {
	return []Months{MonthsJanuary, MonthsFebruary}
}

// IsUnknown returns false for all known variants of Months and true otherwise.
func (e Months) IsUnknown() bool {
	switch e {
	case MonthsJanuary, MonthsFebruary:
		return false
	}
	return true
}
func (e *Months) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		if !enumValuePattern.MatchString(v) {
			return werror.Convert(errors.NewInvalidArgument(wparams.NewSafeAndUnsafeParamStorer(map[string]interface{}{"enumType": "Months", "message": "enum value must match pattern ^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$"}, map[string]interface{}{"enumValue": string(data)})))
		}
		*e = Months(v)
	case "JANUARY":
		*e = MonthsJanuary
	case "FEBRUARY":
		*e = MonthsFebruary
	}
	return nil
}

// These represent values
type Values string

const (
	ValuesNullValue  Values = "NULL_VALUE"
	ValuesValidValue Values = "VALID_VALUE"
)

// Values_Values returns all known variants of Values.
func Values_Values() []Values {
	return []Values{ValuesNullValue, ValuesValidValue}
}

// IsUnknown returns false for all known variants of Values and true otherwise.
func (e Values) IsUnknown() bool {
	switch e {
	case ValuesNullValue, ValuesValidValue:
		return false
	}
	return true
}
func (e *Values) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		if !enumValuePattern.MatchString(v) {
			return werror.Convert(errors.NewInvalidArgument(wparams.NewSafeAndUnsafeParamStorer(map[string]interface{}{"enumType": "Values", "message": "enum value must match pattern ^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$"}, map[string]interface{}{"enumValue": string(data)})))
		}
		*e = Values(v)
	case "NULL_VALUE":
		*e = ValuesNullValue
	case "VALID_VALUE":
		*e = ValuesValidValue
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
type Months string

const (
	// Docs for JANUARY
	MonthsJanuary Months = "JANUARY"
	// Docs for FEBRUARY
	MonthsFebruary Months = "FEBRUARY"
)

// Months_Values returns all known variants of Months.
func Months_Values() []Months {
	return []Months{MonthsJanuary, MonthsFebruary}
}

// IsUnknown returns false for all known variants of Months and true otherwise.
func (e Months) IsUnknown() bool {
	switch e {
	case MonthsJanuary, MonthsFebruary:
		return false
	}
	return true
}
func (e *Months) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		if !enumValuePattern.MatchString(v) {
			return werror.Convert(errors.NewInvalidArgument(wparams.NewSafeAndUnsafeParamStorer(map[string]interface{}{"enumType": "Months", "message": "enum value must match pattern ^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$"}, map[string]interface{}{"enumValue": string(data)})))
		}
		*e = Months(v)
	case "JANUARY":
		*e = MonthsJanuary
	case "FEBRUARY":
		*e = MonthsFebruary
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

			assert.Equal(t, strings.Split(currCase.want, "\n"), strings.Split(string(got), "\n"))
		})
	}
}

func docPtr(doc string) *spec.Documentation {
	return (*spec.Documentation)(&doc)
}
