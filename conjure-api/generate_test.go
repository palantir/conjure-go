// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package main_test

import (
	"testing"

	"github.com/palantir/conjure-go/v6/cmd"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	err := cmd.Generate("conjure-api-4.14.1.conjure.json", "")
	require.NoError(t, err)
}

//func BenchmarkUnmarshal(b *testing.B) {
//	irFileBytes, err := ioutil.ReadFile("conjure-api-4.14.1.conjure.json")
//	require.NoError(b, err)
//	b.Run("empty IR", func(b *testing.B) {
//		b.ReportAllocs()
//		irBytes := []byte(`{"version":1}`)
//		for i := 0; i < b.N; i++ {
//			if err := (&spec.ConjureDefinition{}).UnmarshalJSON(irBytes); err != nil {
//				b.Fatal(err)
//			}
//		}
//	})
//	b.Run("small IR", func(b *testing.B) {
//		b.ReportAllocs()
//		irBytes := []byte(`{"version":1,"errors":[],"types":[{"type":"object","object":{"typeName":{"name":"AliasDefinition","package":"com.palantir.conjure.spec"},"fields":[{"fieldName":"typeName","type":{"type":"reference","reference":{"name":"TypeName","package":"com.palantir.conjure.spec"}}}]}}],"services":[],"extensions":{"recommended-product-dependencies":[]}}`)
//		for i := 0; i < b.N; i++ {
//			if err := (&spec.ConjureDefinition{}).UnmarshalJSON(irBytes); err != nil {
//				b.Fatal(err)
//			}
//		}
//	})
//	b.Run("large IR", func(b *testing.B) {
//		b.ReportAllocs()
//		for i := 0; i < b.N; i++ {
//			if err := (&spec.ConjureDefinition{}).UnmarshalJSON(irFileBytes); err != nil {
//				b.Fatal(err)
//			}
//		}
//	})
//}
//
//func BenchmarkMarshal(b *testing.B) {
//	irFileBytes, err := ioutil.ReadFile("conjure-api-4.14.1.conjure.json")
//	require.NoError(b, err)
//	var irFileObj spec.ConjureDefinition
//	require.NoError(b, irFileObj.UnmarshalJSON(irFileBytes))
//	b.Run("empty IR", func(b *testing.B) {
//		b.ReportAllocs()
//		ir := spec.ConjureDefinition{Version: 1}
//		for i := 0; i < b.N; i++ {
//			if _, err := ir.MarshalJSON(); err != nil {
//				b.Fatal(err)
//			}
//		}
//	})
//	b.Run("small IR", func(b *testing.B) {
//		b.ReportAllocs()
//		irBytes := []byte(`{"version":1,"errors":[],"types":[{"type":"object","object":{"typeName":{"name":"AliasDefinition","package":"com.palantir.conjure.spec"},"fields":[{"fieldName":"typeName","type":{"type":"reference","reference":{"name":"TypeName","package":"com.palantir.conjure.spec"}}}]}}],"services":[],"extensions":{"recommended-product-dependencies":[]}}`)
//		var ir spec.ConjureDefinition
//		require.NoError(b, ir.UnmarshalJSON(irBytes))
//		for i := 0; i < b.N; i++ {
//			if _, err := ir.MarshalJSON(); err != nil {
//				b.Fatal(err)
//			}
//		}
//	})
//	b.Run("large IR", func(b *testing.B) {
//		b.ReportAllocs()
//		for i := 0; i < b.N; i++ {
//			if _, err := irFileObj.MarshalJSON(); err != nil {
//				b.Fatal(err)
//			}
//		}
//	})
//}
