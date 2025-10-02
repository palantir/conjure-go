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
	"bytes"
	stdjson "encoding/json"
	"io"
	"os"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	spec_old "github.com/palantir/conjure-go/v6/conjure-api/spec_old"
	"github.com/stretchr/testify/require"
)

const (
	newLargeOnly = false
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
	b.ReportMetric(float64(len(irBytes)), "B/payload")
	b.Run("jsonv2_NewDecoder_UnmarshalJSONFrom", func(b *testing.B) {
		if newLargeOnly {
			b.Skip("profile")
		}
		b.ReportAllocs()
		b.ReportMetric(float64(len(irBytes)), "B/payload")
		for i := 0; i < b.N; i++ {
			reader := bytes.NewReader(irBytes)
			dec := jsontext.NewDecoder(reader)
			err := (&spec.ConjureDefinition{}).UnmarshalJSONFrom(dec)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("jsonv2_NewUnmarshalerFrom", func(b *testing.B) {
		if newLargeOnly {
			b.Skip("profile")
		}
		b.ReportAllocs()
		b.ReportMetric(float64(len(irBytes)), "B/payload")
		for i := 0; i < b.N; i++ {
			err := json.Unmarshal(irBytes, cj.NewUnmarshalerFrom[spec.ConjureDefinition, cj.StructUnmarshaler[*spec.ConjureDefinition]](&spec.ConjureDefinition{}))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("jsonv2", func(b *testing.B) {
		b.ReportAllocs()
		b.ReportMetric(float64(len(irBytes)), "B/payload")
		for i := 0; i < b.N; i++ {
			err := json.Unmarshal(irBytes, &spec.ConjureDefinition{})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("develop", func(b *testing.B) {
		if newLargeOnly {
			b.Skip("profile")
		}
		b.ReportAllocs()
		b.ReportMetric(float64(len(irBytes)), "B/payload")
		for i := 0; i < b.N; i++ {
			if err := stdjson.Unmarshal(irBytes, &spec_old.ConjureDefinition{}); err != nil {
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
	b.Run("jsontext_NewEncoder_MarshalJSONTo", func(b *testing.B) {
		b.ReportAllocs()
		b.ReportMetric(float64(len(irBytes)), "B/payload")
		for i := 0; i < b.N; i++ {
			enc := jsontext.NewEncoder(io.Discard)
			err := irGenerated.MarshalJSONTo(enc)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("jsonv2_NewMarshalerTo", func(b *testing.B) {
		b.ReportAllocs()
		b.ReportMetric(float64(len(irBytes)), "B/payload")
		for i := 0; i < b.N; i++ {
			err := json.MarshalWrite(io.Discard, cj.NewMarshalerTo[spec.ConjureDefinition, cj.StructMarshaler[spec.ConjureDefinition]](irGenerated))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("jsonv2", func(b *testing.B) {
		b.ReportAllocs()
		b.ReportMetric(float64(len(irBytes)), "B/payload")
		for i := 0; i < b.N; i++ {
			err := json.MarshalWrite(io.Discard, irGenerated)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	var irDevelop spec_old.ConjureDefinition
	require.NoError(b, irDevelop.UnmarshalJSON(irBytes))
	b.Run("develop", func(b *testing.B) {
		if newLargeOnly {
			b.Skip("profile")
		}
		b.ReportAllocs()
		b.ReportMetric(float64(len(irBytes)), "B/payload")
		for i := 0; i < b.N; i++ {
			enc := stdjson.NewEncoder(io.Discard)
			err := enc.Encode(irDevelop)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
