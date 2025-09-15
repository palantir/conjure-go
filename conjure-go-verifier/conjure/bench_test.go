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

package conjure

import (
	stdjson "encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/palantir/conjure-go/v6/cj"
	types1 "github.com/palantir/conjure-go/v6/cj/types"
	"github.com/palantir/conjure-go/v6/conjure-go-verifier/conjure/verification/types"
	"github.com/palantir/pkg/datetime"
)

func BenchmarkDateTime(b *testing.B) {
	date := types.DateTimeAliasExample(time.Date(2025, 5, 12, 19, 26, 0, 0, time.UTC))
	dateText := []byte(strconv.Quote(date.String()))
	b.Run("Marshal", func(b *testing.B) {
		b.Run("text", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := cj.Marshal[types.DateTimeAliasExample, types1.StringerMarshaler[types.DateTimeAliasExample]](date)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("cj", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := date.MarshalJSON()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("v1", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := stdjson.Marshal(datetime.DateTime(date))
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})
	b.Run("Unmarshal", func(b *testing.B) {
		b.Run("text", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var obj types.DateTimeAliasExample
				err := cj.Unmarshal[types.DateTimeAliasExample, types1.TextUnmarshaler[*types.DateTimeAliasExample]](dateText, &obj)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("cj", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var obj types.DateTimeAliasExample
				if err := obj.UnmarshalJSON(dateText); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("v1", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var out datetime.DateTime
				if err := stdjson.Unmarshal(dateText, &out); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}
