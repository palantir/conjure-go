// Copyright (c) 2023 Palantir Technologies. All rights reserved.
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

package dj

import (
	"io"
)

var CODEC = codec{}

type codec struct{}

func (codec) Accept() string {
	return "application/json"
}

func (codec) Decode(r io.Reader, v interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (codec) ContentType() string {
	return "application/json"
}

func (codec) Encode(w io.Writer, v interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (codec) Marshal(v interface{}) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
