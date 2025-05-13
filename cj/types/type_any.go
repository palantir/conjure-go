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

package types

import (
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
)

// Any provides generic JSON marshaling and unmarshaling for any Go type T.
// It is a fallback encoder/decoder for types not otherwise handled by more specific
// implementations. Use this when you want to delegate to the default Go JSON logic,
// but still participate in the MarshalerTo/UnmarshalerFrom interfaces.
type Any[T any] struct{}

func (Any[T]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	return json.MarshalEncode(enc, receiver, DefaultOptions)
}

func (Any[T]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	if kind := dec.PeekKind(); kind == 'n' {
		return cj.NewKindMismatchError(dec, kind, "non-optional value")
	}
	return json.UnmarshalDecode(dec, receiver, DefaultOptions)
}
