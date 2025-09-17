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
	"testing"
	"time"

	"github.com/palantir/conjure-go/v6/cj"
	"github.com/palantir/pkg/datetime"
)

func TestDateTime(t *testing.T) {
	tests := []struct {
		Name string
		Test typeTest
	}{
		{
			Name: "zero_time",
			Test: typeTestCase[time.Time, cj.DateTime[time.Time], cj.DateTime[time.Time]]{
				Value: time.Time{}, JSON: "\"0001-01-01T00:00:00Z\"",
			},
		},
		{
			Name: "iso8601_time",
			Test: typeTestCase[time.Time, cj.DateTime[time.Time], cj.DateTime[time.Time]]{
				Value: time.Date(2025, 5, 12, 19, 26, 0, 0, time.UTC), JSON: "\"2025-05-12T19:26:00Z\"",
			},
		},
		{
			Name: "zero_datetime",
			Test: typeTestCase[datetime.DateTime, cj.DateTime[datetime.DateTime], cj.DateTime[datetime.DateTime]]{
				Value: datetime.DateTime{}, JSON: "\"0001-01-01T00:00:00Z\"",
			},
		},
		{
			Name: "iso8601_datetime",
			Test: typeTestCase[datetime.DateTime, cj.DateTime[datetime.DateTime], cj.DateTime[datetime.DateTime]]{
				Value: must(datetime.ParseDateTime("2025-05-12T19:26:00Z")), JSON: "\"2025-05-12T19:26:00Z\"",
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
