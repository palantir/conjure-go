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

package cj_test

import (
	"fmt"
	"testing"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
)

func TestStructs(t *testing.T) {
	for name, test := range map[string]typeTest{
		"simpleStruct": typeTestCase[simpleStruct, cj.StructMarshaler[simpleStruct], cj.StructUnmarshaler[*simpleStruct]]{
			Value: simpleStruct{Name: "foo", Num: 42}, JSON: `{"name":"foo","num":42}`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}

type simpleStruct struct {
	Name string
	Num  int
}

func (s simpleStruct) MarshalJSONTo(enc *jsontext.Encoder) error {
	return enc.WriteValue(jsontext.Value(fmt.Sprintf(`{"name":"%s","num":%d}`, s.Name, s.Num)))
}

func (s *simpleStruct) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if kind := tok.Kind(); kind != '{' {
		return cj.NewKindMismatchError(dec, kind, "json object")
	}
	for {
		key, err := dec.ReadToken()
		if err != nil {
			return err
		}
		if kind := key.Kind(); kind == '}' {
			break // End of object
		} else if kind != '"' {
			return cj.NewKindMismatchError(dec, kind, "next key or closing brace for EnumValueDefinition")
		}
		switch key.String() {
		case "name":
			if err := (cj.String[string]{}).UnmarshalJSONFrom(dec, &s.Name); err != nil {
				return err
			}
		case "num":
			if err := (cj.Int32[int]{}).UnmarshalJSONFrom(dec, &s.Num); err != nil {
				return err
			}
		}
	}
	return nil
}
