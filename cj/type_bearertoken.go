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
	"sync"
	"unicode/utf8"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/pkg/bearertoken"
)

type BearerToken[T ~string] struct{}

func (BearerToken[T]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	return String[bearertoken.Token]{}.MarshalJSONTo(enc, bearertoken.Token(receiver))
}

func (BearerToken[T]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	str := tok.String()
	if len(str) == 0 || str[0] == '=' {
		return NewInvalidValueError(dec, "empty bearer token", nil)
	}
	chars := validBearerTokenChars()
	for i := 0; i < len(str); i++ {
		if !chars[str[i]] {
			return NewInvalidValueError(dec, "invalid character for bearer token", nil)
		}
	}
	*receiver = T(tok.String())
	return nil
}

var validBearerTokenChars = sync.OnceValue(func() [utf8.RuneSelf]bool {
	var chars [utf8.RuneSelf]bool
	for i := '0'; i <= '9'; i++ {
		chars[i] = true
	}
	for i := 'A'; i <= 'Z'; i++ {
		chars[i] = true
	}
	for i := 'a'; i <= 'z'; i++ {
		chars[i] = true
	}
	chars['+'] = true
	chars['-'] = true
	chars['.'] = true
	chars['/'] = true
	chars['='] = true
	chars['_'] = true
	chars['~'] = true
	return chars
})
