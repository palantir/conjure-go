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
	"time"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
	"github.com/palantir/pkg/datetime"
)

type DateTime[T time.Time | datetime.DateTime] struct{}

func (DateTime[T]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	return (StringerMarshaler[datetime.DateTime]{}).MarshalJSONTo(enc, datetime.DateTime(receiver))
}

func (DateTime[T]) Compare(a, b T) int {
	aTime, bTime := time.Time(a), time.Time(b)
	if aTime.After(bTime) {
		return 1
	}
	if aTime.Before(bTime) {
		return -1
	}
	return 0
}

func (DateTime[T]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	parse, err := time.Parse(time.RFC3339Nano, tok.String())
	if err != nil {
		return cj.WrapSyntaxError(dec, "invalid datetime", err)
	}
	*receiver = T(parse)
	return nil
}
