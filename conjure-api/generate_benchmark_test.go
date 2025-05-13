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

package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	specold "github.com/palantir/conjure-go/v6/conjure-api/conjure/spec_old"
	"github.com/stretchr/testify/require"
)

const (
	newLargeOnly = true
)

func BenchmarkUnmarshal(b *testing.B) {
	b.Run("empty IR", func(b *testing.B) {
		if newLargeOnly {
			b.Skip("profile")
		}
		irBytes := []byte(`{"version":1}`)
		doBenchUnmarshal(b, irBytes)
	})
	b.Run("small IR", func(b *testing.B) {
		if newLargeOnly {
			b.Skip("profile")
		}
		irBytes := []byte(`{"version":1,"errors":[],"types":[{"type":"object","object":{"typeName":{"name":"AliasDefinition","package":"com.palantir.conjure.spec"},"fields":[{"fieldName":"typeName","type":{"type":"reference","reference":{"name":"TypeName","package":"com.palantir.conjure.spec"}}}]}}],"services":[],"extensions":{"recommended-product-dependencies":[]}}`)
		doBenchUnmarshal(b, irBytes)
	})
	b.Run("large IR", func(b *testing.B) {
		irFileBytes, err := os.ReadFile("conjure-api-4.35.0.conjure.json")
		require.NoError(b, err)
		doBenchUnmarshal(b, irFileBytes)
	})
}

func doBenchUnmarshal(b *testing.B, irBytes []byte) {
	b.Helper()
	b.Run("generated", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if err := (&spec.ConjureDefinition{}).UnmarshalJSON(irBytes); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("develop", func(b *testing.B) {
		if newLargeOnly {
			b.Skip("profile")
		}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if err := (&specold.ConjureDefinition{}).UnmarshalJSON(irBytes); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkMarshal(b *testing.B) {
	b.Run("empty IR", func(b *testing.B) {
		if newLargeOnly {
			b.Skip("profile")
		}
		irBytes := []byte(`{"version":1}`)
		doBenchMarshal(b, irBytes)
	})
	b.Run("small IR", func(b *testing.B) {
		if newLargeOnly {
			b.Skip("profile")
		}
		irBytes := []byte(`{"version":1,"errors":[],"types":[{"type":"object","object":{"typeName":{"name":"AliasDefinition","package":"com.palantir.conjure.spec"},"fields":[{"fieldName":"typeName","type":{"type":"reference","reference":{"name":"TypeName","package":"com.palantir.conjure.spec"}}}]}}],"services":[],"extensions":{"recommended-product-dependencies":[]}}`)
		doBenchMarshal(b, irBytes)
	})
	b.Run("large IR", func(b *testing.B) {
		irBytes, err := os.ReadFile("conjure-api-4.35.0.conjure.json")
		require.NoError(b, err)
		doBenchMarshal(b, irBytes)
	})
}

func doBenchMarshal(b *testing.B, irBytes []byte) {
	var irGenerated spec.ConjureDefinition
	require.NoError(b, irGenerated.UnmarshalJSON(irBytes))
	doBenchMarshalJSON(b, "generated", irGenerated)
	if newLargeOnly {
		return
	}
	var irDevelop specold.ConjureDefinition
	require.NoError(b, irDevelop.UnmarshalJSON(irBytes))
	doBenchMarshalJSON(b, "develop", irDevelop)
}

func doBenchMarshalJSON(b *testing.B, name string, irObj json.Marshaler) {
	b.Helper()
	b.Run(name, func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := irObj.MarshalJSON()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
