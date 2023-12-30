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
	"os"
	"testing"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	specold "github.com/palantir/conjure-go/v6/conjure-api/conjure/spec_old"
	"github.com/palantir/pkg/safejson"
	"github.com/stretchr/testify/require"
)

func BenchmarkUnmarshal(b *testing.B) {
	irFileBytes, err := os.ReadFile("conjure-api-4.35.0.conjure.json")
	require.NoError(b, err)

	b.Run("empty IR", func(b *testing.B) {
		//b.Skip("profile")
		irBytes := []byte(`{"version":1}`)
		doBenchUnmarshalOld(b, irBytes)
		doBenchUnmarshalNew(b, irBytes)
	})
	b.Run("small IR", func(b *testing.B) {
		//b.Skip("profile")
		irBytes := []byte(`{"version":1,"errors":[],"types":[{"type":"object","object":{"typeName":{"name":"AliasDefinition","package":"com.palantir.conjure.spec"},"fields":[{"fieldName":"typeName","type":{"type":"reference","reference":{"name":"TypeName","package":"com.palantir.conjure.spec"}}}]}}],"services":[],"extensions":{"recommended-product-dependencies":[]}}`)
		doBenchUnmarshalOld(b, irBytes)
		doBenchUnmarshalNew(b, irBytes)
	})
	b.Run("large IR", func(b *testing.B) {
		doBenchUnmarshalOld(b, irFileBytes)
		doBenchUnmarshalNew(b, irFileBytes)
	})
}

func doBenchUnmarshalOld(b *testing.B, irBytes []byte) {
	b.Run("old", func(b *testing.B) {
		//b.Skip("profile")
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if err := (&specold.ConjureDefinition{}).UnmarshalJSON(irBytes); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func doBenchUnmarshalNew(b *testing.B, irBytes []byte) {
	b.Run("new", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if err := (&spec.ConjureDefinition{}).UnmarshalJSON(irBytes); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkMarshal(b *testing.B) {
	irFileBytes, err := os.ReadFile("conjure-api-4.35.0.conjure.json")
	require.NoError(b, err)
	b.Run("empty IR", func(b *testing.B) {
		//b.Skip("profile")
		irBytes := []byte(`{"version":1}`)
		doBenchMarshalOld(b, irBytes)
		doBenchMarshalNew(b, irBytes)
	})
	b.Run("small IR", func(b *testing.B) {
		//b.Skip("profile")
		b.ReportAllocs()
		irBytes := []byte(`{"version":1,"errors":[],"types":[{"type":"object","object":{"typeName":{"name":"AliasDefinition","package":"com.palantir.conjure.spec"},"fields":[{"fieldName":"typeName","type":{"type":"reference","reference":{"name":"TypeName","package":"com.palantir.conjure.spec"}}}]}}],"services":[],"extensions":{"recommended-product-dependencies":[]}}`)
		doBenchMarshalOld(b, irBytes)
		doBenchMarshalNew(b, irBytes)
	})
	b.Run("large IR", func(b *testing.B) {
		doBenchMarshalOld(b, irFileBytes)
		doBenchMarshalNew(b, irFileBytes)
	})
}

func doBenchMarshalOld(b *testing.B, irBytes []byte) {
	var ir specold.ConjureDefinition
	require.NoError(b, ir.UnmarshalJSON(irBytes))
	b.Run("old", func(b *testing.B) {
		//b.Skip("profile")
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf := new(bytes.Buffer)
			enc := safejson.Encoder(buf)
			enc.SetEscapeHTML(false)
			if err := enc.Encode(ir); err != nil {
				b.Fatal(err)
			}
			_ = buf.Bytes()
		}
	})
}

func doBenchMarshalNew(b *testing.B, irBytes []byte) {
	var ir spec.ConjureDefinition
	require.NoError(b, ir.UnmarshalJSON(irBytes))
	b.Run("new", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf := new(bytes.Buffer)
			if _, err := ir.WriteJSON(buf); err != nil {
				b.Fatal(err)
			}
			_ = buf.Bytes()
		}
	})
}
