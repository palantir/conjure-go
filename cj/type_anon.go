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
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

func NewMarshalerTo[T any, E TypeEncoder[T]](receiver T) json.MarshalerTo {
	return anonMarshalerTo[T, E]{receiver: receiver}
}

type anonMarshalerTo[T any, E TypeEncoder[T]] struct {
	receiver T
}

func (f anonMarshalerTo[T, E]) MarshalJSONTo(enc *jsontext.Encoder) error {
	return (*new(E)).MarshalJSONTo(enc, f.receiver)
}

func NewUnmarshalerFrom[T any, D TypeDecoder[T]](receiver *T) json.UnmarshalerFrom {
	return &anonUnmarshalerFrom[T, D]{receiver: receiver}
}

type anonUnmarshalerFrom[T any, D TypeDecoder[T]] struct {
	receiver *T
}

func (f *anonUnmarshalerFrom[T, D]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return (*new(D)).UnmarshalJSONFrom(dec, f.receiver)
}
