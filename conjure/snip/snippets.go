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

// MethodUnmarshalText returns 'func (o *Foo) UnmarshalText(data []byte) error'
func MethodUnmarshalText(receiverName, receiverType string) *jen.Statement {
	return jen.Func().Params(jen.Id(receiverName).Op("*").Id(receiverType)).
		Id("UnmarshalText").Params(jen.Id("data").Index().Byte()).Params(jen.Error())
}

func MethodMarshalJSONV2(receiverName string, receiverTypeName string) *jen.Statement {
	return MethodMarshalJSON(receiverName, receiverTypeName).Block(
		jen.Return(JSONV2Marshal().Call(jen.Id(receiverName), JSONV2AllowDuplicateNames().Call(jen.True()))),
	)
}

// MethodMarshalJSONTo returns 'func (o Foo) MarshalJSONTo(enc jsontext.JSONEncoder) error'
func MethodMarshalJSONTo(receiverName, receiverType string) *jen.Statement {
	return jen.Func().Params(jen.Id(receiverName).Id(receiverType)).
		Id("MarshalJSONTo").Params(jen.Id("enc").Op("*").Add(JSONV2Encoder())).Error()
}

func MethodUnmarshalJSONV2(receiverName string, receiverTypeName string) *jen.Statement {
	return MethodUnmarshalJSON(receiverName, receiverTypeName).Block(
		jen.Return(JSONV2Unmarshal().Call(jen.Id("data"), jen.Id(receiverName))),
	)
}

// MethodUnmarshalJSONFrom returns 'func (o *Foo) UnmarshalJSONFrom(dec jsontext.JSONDecoder) error'
func MethodUnmarshalJSONFrom(receiverName, receiverType string) *jen.Statement {
	return jen.Func().Params(jen.Id(receiverName).Op("*").Id(receiverType)).
		Id("UnmarshalJSONFrom").Params(jen.Id("dec").Op("*").Add(JSONV2Decoder())).Error()
}

// MethodMarshalYAML returns:
//
//	func (o Foo) MarshalYAML() (interface{}, error) {
//		jsonBytes, err := safejson.Marshal(o)
//		if err != nil {
//			return nil, err
//		}
//		return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
//	}
func MethodMarshalYAML(receiverName, receiverType string) *jen.Statement {
	return MethodMarshalYAMLSig(receiverName, receiverType).Block(
		jen.List(jen.Id("jsonBytes"), jen.Err()).Op(":=").Add(SafeJSONMarshal()).Params(jen.Id(receiverName)),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Err()),
		),
		jen.Return(SafeYAMLJSONtoYAMLMapSlice().Call(jen.Id("jsonBytes"))),
	)
}

// MethodJSONV2MarshalYAML returns:
//
//	func (o Foo) MarshalYAML() (interface{}, error) {
//		return cj.MarshalYAML(o)
//	}
func MethodJSONV2MarshalYAML(receiverName, receiverType string) *jen.Statement {
	return MethodMarshalYAMLSig(receiverName, receiverType).Block(
		jen.Return(CJMarshalYAML().Call(jen.Id(receiverName))),
	)
}

// MethodMarshalYAMLSig returns:
//
//	func (o Foo) MarshalYAML() (interface{}, error)
func MethodMarshalYAMLSig(receiverName, receiverType string) *jen.Statement {
	return jen.Func().Params(jen.Id(receiverName).Id(receiverType)).
		Id("MarshalYAML").Params().Params(jen.Interface(), jen.Id("error"))
}

// MethodUnmarshalYAML returns:
//
//	func (o *Foo) UnmarshalYAML(unmarshal func(any) error) error {
//	  jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
//	  if err != nil {
//	    return err
//	  }
//	  return safejson.Unmarshal(jsonBytes, *&o)
//	}
func MethodUnmarshalYAML(receiverName, receiverType string) *jen.Statement {
	return MethodUnmarshalYAMLSig(receiverName, receiverType).Block(
		jen.List(jen.Id("jsonBytes"), jen.Err()).Op(":=").Add(SafeYAMLUnmarshalerToJSONBytes()).Params(jen.Id("unmarshal")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Err()),
		),
		jen.Return(SafeJSONUnmarshal().Call(jen.Id("jsonBytes"), jen.Op("*").Op("&").Id(receiverName))),
	)
}

// MethodJSONV2UnmarshalYAML returns:
//
//	func (o *Foo) UnmarshalYAML(unmarshal func(any) error) error {
//		return cj.UnmarshalYAML(o, unmarshal)
//	}
func MethodJSONV2UnmarshalYAML(receiverName, receiverType string) *jen.Statement {
	return MethodUnmarshalYAMLSig(receiverName, receiverType).Block(
		jen.Return(CJUnmarshalYAML().Call(jen.Id(receiverName), jen.Id("unmarshal"))),
	)
}

// MethodUnmarshalYAMLSig returns:
//
//	func (o *Foo) UnmarshalYAML(unmarshal func(any) error) error
func MethodUnmarshalYAMLSig(receiverName, receiverType string) *jen.Statement {
	return jen.Func().Params(jen.Id(receiverName).Op("*").Id(receiverType)).
		Id("UnmarshalYAML").Params(jen.Id("unmarshal").Func().Params(jen.Interface()).Error()).Error()
}
