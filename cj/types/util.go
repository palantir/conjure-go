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
	werror "github.com/palantir/witchcraft-go-error"
)

func readStringToken(dec *jsontext.Decoder) (jsontext.Token, error) {
	tok, err := dec.ReadToken()
	if err != nil {
		return tok, werror.Convert(err)
	}
	if kind := tok.Kind(); kind != '"' {
		return tok, cj.NewKindMismatchError(dec, kind, "json string")
	}
	return tok, nil
}
