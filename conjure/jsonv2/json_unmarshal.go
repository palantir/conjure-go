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

package jsonv2

const (
	decName = "dec"
)

//func UnmarshalJSONFromMethod(receiverName string, receiverTypeName string, receiverType types.Type) *jen.Statement {
//	return jen.Func().
//		Params(jen.Id(receiverName).Op("*").Id(receiverTypeName)).
//		Id("UnmarshalJSONFrom").
//		Params(jen.Id(decName).Op("*").Add(snip.JSONV2Decoder())).
//		Error().
//		BlockFunc(func(methodBody *jen.Group) {
//			switch v := receiverType.(type) {
//			default:
//				marshalJSONValue(methodBody, jen.Id(receiverName), v, 0, false)
//				methodBody.Return(jen.Nil())
//			case *types.ObjectType:
//				var fields []jsonStructField
//				for _, field := range v.Fields {
//					fields = append(fields, jsonStructField{
//						Key:      field.Name,
//						Type:     field.Type,
//						Selector: jen.Id(receiverName).Dot(transforms.ExportedFieldName(field.Name)).Clone,
//					})
//				}
//				methodBody.Add(mustWriteToken(snip.JSONV2BeginObject()))
//				for fieldIdx := range fields {
//					marshalJSONStructField(methodBody, fields, fieldIdx)
//				}
//				methodBody.Add(mustWriteToken(snip.JSONV2EndObject()))
//				methodBody.Return(jen.Nil())
//			case *types.UnionType:
//				var fields []jsonStructField
//				for _, field := range v.Fields {
//					fields = append(fields, jsonStructField{
//						Key:      field.Name,
//						Type:     field.Type,
//						Selector: jen.Id(receiverName).Dot(transforms.PrivateFieldName(field.Name)).Clone,
//					})
//				}
//				methodBody.Add(mustWriteToken(snip.JSONV2BeginObject()))
//				methodBody.Switch(jen.Id(receiverName).Dot("typ")).BlockFunc(func(cases *jen.Group) {
//					for _, field := range fields {
//						cases.Case(jen.Lit(field.Key)).BlockFunc(func(caseBody *jen.Group) {
//							caseBody.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit("type"))))
//							caseBody.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit(field.Key))))
//							caseBody.If(field.Selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
//								ifBody.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit(field.Key))))
//								ifBody.Id("unionVal").Op(":=").Op("*").Add(field.Selector())
//								marshalJSONValue(ifBody, jen.Id("unionVal"), field.Type, 0, false)
//							})
//						})
//					}
//					cases.Default().Block(
//						mustWriteToken(snip.JSONV2String().Call(jen.Lit("type"))),
//						mustWriteToken(snip.JSONV2String().Call(jen.Id(receiverName).Dot("typ"))),
//					)
//				})
//				methodBody.Add(mustWriteToken(snip.JSONV2EndObject()))
//				methodBody.Return(jen.Nil())
//			}
//		})
//}
