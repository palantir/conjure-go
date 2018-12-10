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

package conjure_test

import (
	"testing"

	"github.com/palantir/goastwriter"
	"github.com/palantir/goastwriter/astgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/conjure"
)

func TestEnum(t *testing.T) {
	for caseNum, currCase := range []struct {
		pkg   string
		name  string
		enums []*conjure.Enum
		want  string
	}{
		{
			pkg:  "testpkg",
			name: "single enum",
			enums: []*conjure.Enum{
				{
					Name:    "Months",
					Comment: "These represent months",
					Values: []conjure.Value{
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
	MonthsUnknown  Months = "UNKNOWN"
)

func (e *Months) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch strings.ToUpper(s) {
	default:
		*e = MonthsUnknown
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
			enums: []*conjure.Enum{
				{
					Name:    "Months",
					Comment: "These represent months",
					Values: []conjure.Value{
						{Value: "JANUARY"},
						{Value: "FEBRUARY"},
					},
				},
				{
					Name:    "Values",
					Comment: "These represent values",
					Values: []conjure.Value{
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
	MonthsUnknown  Months = "UNKNOWN"
)

func (e *Months) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch strings.ToUpper(s) {
	default:
		*e = MonthsUnknown
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
	ValuesUnknown    Values = "UNKNOWN"
)

func (e *Values) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch strings.ToUpper(s) {
	default:
		*e = ValuesUnknown
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
			enums: []*conjure.Enum{
				{
					Name:    "Months",
					Comment: "These represent months",
					Values: []conjure.Value{
						{
							Value: "JANUARY",
							Docs:  "Docs for JANUARY",
						},
						{
							Value: "FEBRUARY",
							Docs:  "Docs for FEBRUARY",
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
	MonthsUnknown  Months = "UNKNOWN"
)

func (e *Months) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch strings.ToUpper(s) {
	default:
		*e = MonthsUnknown
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
		var components []astgen.ASTDecl
		for _, e := range currCase.enums {
			declers, _ := e.ASTDeclers()
			components = append(components, declers...)
		}

		got, err := goastwriter.Write(currCase.pkg, components...)
		require.NoError(t, err, "Case %d: %s", caseNum, currCase.name)
		assert.Equal(t, currCase.want, string(got), "Case %d: %s\n%s", caseNum, currCase.name, string(got))
	}
}
