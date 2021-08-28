// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package snip

import (
	"github.com/dave/jennifer/jen"
)

// MethodString returns 'func (o Foo) String() string'
func MethodString(receiverName, receiverType string) *jen.Statement {
	return jen.Func().Params(jen.Id(receiverName).Id(receiverType)).
		Id("String").Params().String()
}

// MethodAppendJSON returns 'func (o Foo) AppendJSON(out []byte) ([]byte, error)'
func MethodAppendJSON(receiverName, receiverType string) *jen.Statement {
	return jen.Func().Params(jen.Id(receiverName).Id(receiverType)).
		Id("AppendJSON").Params(jen.Id("out").Id("[]byte")).Params(jen.Id("[]byte"), jen.Error())
}

// MethodMarshalJSON returns 'func (o Foo) MarshalJSON() ([]byte, error)'
func MethodMarshalJSON(receiverName, receiverType string) *jen.Statement {
	return jen.Func().Params(jen.Id(receiverName).Id(receiverType)).
		Id("MarshalJSON").Params().Params(jen.Id("[]byte"), jen.Error())
}

// MethodMarshalText returns 'func (o Foo) MarshalText() ([]byte, error)'
func MethodMarshalText(receiverName, receiverType string) *jen.Statement {
	return jen.Func().Params(jen.Id(receiverName).Id(receiverType)).
		Id("MarshalText").Params().Params(jen.Id("[]byte"), jen.Error())
}

// MethodUnmarshalJSON returns 'func (o *Foo) UnmarshalJSON(data []byte) error'
func MethodUnmarshalJSON(receiverName, receiverType string) *jen.Statement {
	return jen.Func().Params(jen.Id(receiverName).Op("*").Id(receiverType)).
		Id("UnmarshalJSON").Params(jen.Id("data").Id("[]byte")).Params(jen.Error())
}

// MethodUnmarshalJSONStrict returns 'func (o *Foo) UnmarshalJSONStrict(data []byte) error'
func MethodUnmarshalJSONStrict(receiverName, receiverType string) *jen.Statement {
	return jen.Func().Params(jen.Id(receiverName).Op("*").Id(receiverType)).
		Id("UnmarshalJSONStrict").Params(jen.Id("data").Id("[]byte")).Params(jen.Error())
}

// MethodUnmarshalText returns 'func (o *Foo) UnmarshalText(data []byte) error'
func MethodUnmarshalText(receiverName, receiverType string) *jen.Statement {
	return jen.Func().Params(jen.Id(receiverName).Op("*").Id(receiverType)).
		Id("UnmarshalText").Params(jen.Id("data").Id("[]byte")).Params(jen.Error())
}

// MethodMarshalYAML returns:
//	func (o Foo) MarshalYAML() (interface{}, error) {
//		jsonBytes, err := safejson.Marshal(o)
//		if err != nil {
//			return nil, err
//		}
//		return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
//	}
func MethodMarshalYAML(receiverName, receiverType string) *jen.Statement {
	return jen.Func().Params(jen.Id(receiverName).Id(receiverType)).
		Id("MarshalYAML").Params().Params(jen.Interface(), jen.Id("error")).Block(
		jen.List(jen.Id("jsonBytes"), jen.Err()).Op(":=").Add(SafeJSONMarshal()).Params(jen.Id(receiverName)),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Err()),
		),
		jen.Return(SafeYAMLJSONtoYAMLMapSlice().Call(jen.Id("jsonBytes"))),
	)
}

// MethodUnmarshalYAML returns:
//  func (o *Foo) UnmarshalYAML(unmarshal func(interface{}) error) error {
//    jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
//    if err != nil {
//      return err
//	  }
//	  return safejson.Unmarshal(jsonBytes, *&o)
//  }
func MethodUnmarshalYAML(receiverName, receiverType string) *jen.Statement {
	return jen.Func().Params(jen.Id(receiverName).Op("*").Id(receiverType)).
		Id("UnmarshalYAML").Params(jen.Id("unmarshal").Func().Params(jen.Interface()).Error()).Error().Block(
		jen.List(jen.Id("jsonBytes"), jen.Err()).Op(":=").Add(SafeYAMLUnmarshalerToJSONBytes()).Params(jen.Id("unmarshal")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Err()),
		),
		jen.Return(SafeJSONUnmarshal().Call(jen.Id("jsonBytes"), jen.Op("*").Op("&").Id(receiverName))),
	)
}
