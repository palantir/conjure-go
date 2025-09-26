// Copyright (c) 2025 Palantir Technologies. All rights reserved.
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

package cj_test

import (
	"math"
	"testing"

	"github.com/palantir/conjure-go/v6/cj"
)

func TestFloat(t *testing.T) {
	tests := []struct {
		Name string
		Test typeTest
	}{
		{
			Name: "zero",
			Test: typeTestCase[float64, cj.Float[float64], cj.Float[float64]]{
				Value: 0.0, JSON: "0",
			},
		},
		{
			Name: "positive",
			Test: typeTestCase[float64, cj.Float[float64], cj.Float[float64]]{
				Value: 123.456, JSON: "123.456",
			},
		},
		{
			Name: "negative",
			Test: typeTestCase[float64, cj.Float[float64], cj.Float[float64]]{
				Value: -42.5, JSON: "-42.5",
			},
		},
		{
			Name: "large",
			Test: typeTestCase[float64, cj.Float[float64], cj.Float[float64]]{
				Value: 1e30, JSON: "1e+30",
			},
		},
		{
			Name: "small",
			Test: typeTestCase[float64, cj.Float[float64], cj.Float[float64]]{
				Value: 1e-18, JSON: "1e-18",
			},
		},
		{
			Name: "nan",
			Test: typeTestCase[float64, cj.Float[float64], cj.Float[float64]]{
				Value: math.NaN(), JSON: "\"NaN\"",
			},
		},
		{
			Name: "+inf",
			Test: typeTestCase[float64, cj.Float[float64], cj.Float[float64]]{
				Value: math.Inf(1), JSON: "\"Infinity\"",
			},
		},
		{
			Name: "-inf",
			Test: typeTestCase[float64, cj.Float[float64], cj.Float[float64]]{
				Value: math.Inf(-1), JSON: "\"-Infinity\"",
			},
		},
		{
			Name: "map",
			Test: typeTestCase[map[float64]float64, cj.OrderedMapMarshaler[map[float64]float64, float64, float64, cj.FloatMapKey[float64], cj.Float[float64]], cj.MapUnmarshaler[map[float64]float64, float64, float64, cj.FloatMapKey[float64], cj.Float[float64]]]{
				Value: map[float64]float64{0.0: 0.0, 123.456: 123.456, -42.5: -42.5, 1e30: 1e30, 1e-18: 1e-18, math.Inf(1): math.Inf(1), math.Inf(-1): math.Inf(-1)},
				JSON:  `{"-Infinity":"-Infinity","-42.5":-42.5,"0":0,"0.000000000000000001":1e-18,"123.456":123.456,"1000000000000000000000000000000":1e+30,"Infinity":"Infinity"}`,
			},
		},
		{
			Name: "nan_as_map_key_rejected",
			Test: typeTestCase[float64, cj.FloatMapKey[float64], cj.FloatMapKey[float64]]{
				JSON:                 `"NaN"`,
				SkipTestMarshal:      true,
				ErrUnmarshalJSONFrom: "InvalidValueError at offset 5: cannot use NaN as map key",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			t.Run("Marshal", tc.Test.TestMarshal)
			t.Run("Unmarshal", tc.Test.TestUnmarshal)
		})
	}
}
