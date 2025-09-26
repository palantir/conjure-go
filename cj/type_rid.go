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

package cj

import (
	"strings"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/pkg/rid"
)

// ridConstraint provides a common interface for types based on rid.ResourceIdentifier.
type ridConstraint interface {
	~struct {
		Service  string
		Instance string
		Type     string
		Locator  string
	}
}

// RID provides JSON marshaling and unmarshaling for types based on rid.ResourceIdentifier.
// Encodes values as JSON strings using the canonical string representation of the resource identifier.
// Implements comparison based on all RID fields for use as a map key.
type RID[T ridConstraint] struct{}

func (RID[T]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	return (StringerMarshaler[rid.ResourceIdentifier]{}).MarshalJSONTo(enc, rid.ResourceIdentifier(receiver))
}

func (RID[T]) Compare(a, b T) int {
	ra, rb := rid.ResourceIdentifier(a), rid.ResourceIdentifier(b)
	if cmp := strings.Compare(ra.Service, rb.Service); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(ra.Instance, rb.Instance); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(ra.Type, rb.Type); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(ra.Locator, rb.Locator); cmp != 0 {
		return cmp
	}
	return 0
}

func (RID[T]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return WrapSyntaxError(dec, "", err)
	}
	if kind := tok.Kind(); kind != '"' {
		return NewKindMismatchError(dec, kind, "json string")
	}
	parsed, err := rid.ParseRID(tok.String())
	if err != nil {
		return WrapSyntaxError(dec, "invalid resource identifier", err)
	}
	*receiver = T(parsed)
	return nil
}
