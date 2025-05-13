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
	"io"

	"github.com/go-json-experiment/json"
)

// Codec implements conjure-go-runtime's Codec interface
var Codec codecJSONV2

type codecJSONV2 struct{}

func (codecJSONV2) Accept() string {
	return "application/json"
}

func (codecJSONV2) Decode(r io.Reader, v interface{}) error {
	return json.UnmarshalRead(r, *&v)
}

func (codecJSONV2) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, *&v)
}

func (codecJSONV2) ContentType() string {
	return "application/json"
}

func (codecJSONV2) Encode(w io.Writer, v interface{}) error {
	return json.MarshalWrite(w, v)
}

func (codecJSONV2) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
