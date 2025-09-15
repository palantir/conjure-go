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
	"github.com/go-json-experiment/json/jsontext"
)

// OptionalMarshaler provides JSON marshaling for optional (pointer) values of type T.
// Encodes nil pointers as JSON null, otherwise delegates encoding to ITEM.
type OptionalMarshaler[T *U, U any, ITEM TypeEncoder[U]] struct{}

func (OptionalMarshaler[T, U, ITEM]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	if receiver == nil {
		return enc.WriteToken(jsontext.Null)
	}
	return (*new(ITEM)).MarshalJSONTo(enc, *receiver)
}

// OptionalUnmarshaler provides JSON unmarshaling for optional (pointer) values of type T.
// Decodes JSON null as nil, otherwise delegates decoding to ITEM.
type OptionalUnmarshaler[T *U, U any, ITEM TypeDecoder[U]] struct{}

func (OptionalUnmarshaler[T, U, ITEM]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	if dec.PeekKind() == 'n' {
		// still have to consume 'null' token
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		*receiver = nil
		return nil
	}
	*receiver = new(U)
	return (*new(ITEM)).UnmarshalJSONFrom(dec, *receiver)
}
