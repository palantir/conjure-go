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
	"testing"
	"time"

	"github.com/palantir/conjure-go/v6/cj/types"
	"github.com/palantir/pkg/datetime"
)

func TestDateTime(t *testing.T) {
	type TimeAlias time.Time
	for name, test := range map[string]typeTest{
		"zero_time": typeTestCase[time.Time, types.DateTime[time.Time], types.DateTime[time.Time]]{
			Value: time.Time{}, JSON: "\"0001-01-01T00:00:00Z\"",
		},
		"iso8601_time": typeTestCase[time.Time, types.DateTime[time.Time], types.DateTime[time.Time]]{
			Value: time.Date(2025, 5, 12, 19, 26, 0, 0, time.UTC), JSON: "\"2025-05-12T19:26:00Z\"",
		},
		"zero_datetime": typeTestCase[datetime.DateTime, types.DateTime[datetime.DateTime], types.DateTime[datetime.DateTime]]{
			Value: datetime.DateTime{}, JSON: "\"0001-01-01T00:00:00Z\"",
		},
		"iso8601_datetime": typeTestCase[datetime.DateTime, types.DateTime[datetime.DateTime], types.DateTime[datetime.DateTime]]{
			Value: must(datetime.ParseDateTime("2025-05-12T19:26:00Z")), JSON: "\"2025-05-12T19:26:00Z\"",
		},
		"alias": typeTestCase[TimeAlias, types.DateTime[TimeAlias], types.DateTime[TimeAlias]]{
			Value: TimeAlias(time.Date(2025, 5, 12, 19, 26, 0, 0, time.UTC)), JSON: "\"2025-05-12T19:26:00Z\"",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}
