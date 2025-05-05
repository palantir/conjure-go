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
)

// StructMarshaler provides JSON marshaling for types that implement json.MarshalerTo.
// Delegates marshaling to the type's MarshalJSONTo method.
type StructMarshaler[T json.MarshalerTo] struct{}

func (StructMarshaler[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return receiver.MarshalJSONTo(enc)
}

// StructUnmarshaler provides JSON unmarshaling for types that implement json.UnmarshalerFrom.
// Delegates unmarshaling to the type's UnmarshalJSONFrom method.
type StructUnmarshaler[T json.UnmarshalerFrom] struct{}

func (StructUnmarshaler[T]) UnmarshalJSONFrom(receiver T, dec *jsontext.Decoder) error {
	return receiver.UnmarshalJSONFrom(dec)
}
