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
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
	"github.com/palantir/pkg/uuid"
)

// UUID provides JSON marshaling and unmarshaling for uuid.UUID.
type UUID[T ~[16]byte] struct{}

func (UUID[T]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	return (StringerMarshaler[uuid.UUID]{}).MarshalJSONTo(enc, uuid.UUID(receiver))
}

func (t UUID[T]) Compare(a, b T) int {
	// UUIDs are 16 bytes, so we can compare them directly
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return int(a[i]) - int(b[i])
		}
	}
	return 0
}

func (UUID[T]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	parsed, err := uuid.ParseUUID(tok.String())
	if err != nil {
		return cj.NewInvalidValueError(dec, "", err)
	}
	*receiver = T(parsed)
	return nil
}
