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

package types_test

import (
	"math"
	"testing"

	"github.com/palantir/conjure-go/v6/cj/types"
)

func TestFloat(t *testing.T) {
	for name, test := range map[string]typeTest{
		"zero": typeTestCase[float64, types.Float[float64], types.Float[float64]]{
			Value: 0.0, JSON: "0",
		},
		"positive": typeTestCase[float64, types.Float[float64], types.Float[float64]]{
			Value: 123.456, JSON: "123.456",
		},
		"negative": typeTestCase[float64, types.Float[float64], types.Float[float64]]{
			Value: -42.5, JSON: "-42.5",
		},
		"large": typeTestCase[float64, types.Float[float64], types.Float[float64]]{
			Value: 1e30, JSON: "1e+30",
		},
		"small": typeTestCase[float64, types.Float[float64], types.Float[float64]]{
			Value: 1e-18, JSON: "1e-18",
		},
		"nan": typeTestCase[float64, types.Float[float64], types.Float[float64]]{
			Value: math.NaN(), JSON: "\"NaN\"",
		},
		"+inf": typeTestCase[float64, types.Float[float64], types.Float[float64]]{
			Value: math.Inf(1), JSON: "\"Infinity\"",
		},
		"-inf": typeTestCase[float64, types.Float[float64], types.Float[float64]]{
			Value: math.Inf(-1), JSON: "\"-Infinity\"",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}
