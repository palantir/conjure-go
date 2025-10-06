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

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	EncName = "enc"
	DecName = "dec"
)

func MethodMarshalJSONTo(receiverName string, receiverTypeName string, receiverType types.Type) *jen.Statement {
	return snip.MethodMarshalJSONTo(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
		switch typ := receiverType.(type) {
		default:
			panic(fmt.Sprintf("unexpected type %T", receiverType))
		case *types.AliasType:
			methodBody.Return(MarshalJSONValue(jen.Id(EncName), aliasTypeItemSelector(typ, jen.Id(receiverName), false), typ.Item))
		case *types.EnumType:
			methodBody.Return(MarshalJSONValue(jen.Id(EncName), jen.Id(receiverName), typ))
		case *types.ObjectType:
			var fields []jsonStructField
			for _, field := range typ.Fields {
				fields = append(fields, jsonStructField{
					Key:      field.Name,
					Type:     field.Type,
					Selector: jen.Id(receiverName).Dot(transforms.ExportedFieldName(field.Name)).Clone,
				})
			}
			marshalJSONStructFields(methodBody, receiverName, fields, false)
		case *types.UnionType:
			var fields []jsonStructField
			for _, field := range typ.Fields {
				fields = append(fields, jsonStructField{
					Key:      field.Name,
					Type:     field.Type,
					Selector: jen.Id(receiverName).Dot(transforms.PrivateFieldName(field.Name)).Clone,
				})
			}
			marshalJSONStructFields(methodBody, receiverName, fields, true)
		}
	})
}

func marshalJSONStructFields(methodBody *jen.Group, receiverName string, fields []jsonStructField, isUnion bool) {
	methodBody.If(
		jen.Err().Op(":=").Id(EncName).Dot("WriteToken").Call(snip.JSONV2BeginObject()),
		jen.Err().Op("!=").Nil(),
	).Block(
		jen.Return(jen.Err()),
	)

	if isUnion {
		methodBody.If(
			jen.Err().Op(":=").Id(EncName).Dot("WriteToken").Call(snip.JSONV2String().Call(jen.Lit("type"))),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Err()),
		)
		methodBody.If(
			jen.Err().Op(":=").Id(EncName).Dot("WriteToken").Call(snip.JSONV2String().Call(jen.Id(receiverName).Dot("typ"))),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Err()),
		)
		methodBody.Switch(jen.Id(receiverName).Dot("typ")).BlockFunc(func(cases *jen.Group) {
			for _, field := range fields {
				cases.Case(jen.Lit(field.Key)).BlockFunc(func(caseBody *jen.Group) {
					caseBody.If(
						jen.Err().Op(":=").Id(EncName).Dot("WriteToken").Call(snip.JSONV2String().Call(jen.Lit(field.Key))),
						jen.Err().Op("!=").Nil(),
					).Block(
						jen.Return(jen.Err()),
					)
					caseBody.If(field.Selector().Op("!=").Nil()).Block(
						jen.If(
							jen.Err().Op(":=").Add(MarshalJSONValue(jen.Id(EncName), jen.Op("*").Add(field.Selector()), field.Type)),
							jen.Err().Op("!=").Nil(),
						).Block(
							jen.Return(jen.Err()),
						),
					).Else().Block(
						jen.If(
							jen.Err().Op(":=").Id(EncName).Dot("WriteToken").Call(snip.JSONV2Null()),
							jen.Err().Op("!=").Nil(),
						).Block(
							jen.Return(jen.Err()),
						),
					)
				})
			}
		})
	} else {
		marshalField := func(condition *jen.Statement, key string, marshalJSONValue *jen.Statement) {
			methodBody.Add(condition).Block(
				jen.If(
					jen.Err().Op(":=").Id(EncName).Dot("WriteToken").Call(snip.JSONV2String().Call(jen.Lit(key))),
					jen.Err().Op("!=").Nil(),
				).Block(
					jen.Return(jen.Err()),
				),
				jen.If(
					jen.Err().Op(":=").Add(marshalJSONValue),
					jen.Err().Op("!=").Nil(),
				).Block(
					jen.Return(jen.Err()),
				),
			)
		}
		for _, field := range fields {
			if !field.Type.IsOptional() {
				marshalField(
					jen.Empty(),
					field.Key,
					MarshalJSONValue(jen.Id(EncName), field.Selector(), field.Type))
			} else {
				switch typ := field.Type.(type) {
				case *types.Optional:
					marshalField(
						jen.If(field.Selector().Op("!=").Nil()),
						field.Key,
						MarshalJSONValue(jen.Id(EncName), jen.Op("*").Add(field.Selector()), typ.Item))
				case *types.AliasType:
					marshalField(
						jen.If(field.Selector().Dot("Value").Op("!=").Nil()),
						field.Key,
						MarshalJSONValue(jen.Id(EncName), field.Selector().Dot("Value"), typ.Item))
				default:
					panic(fmt.Sprintf("unexpected optional type %T", field.Type))
				}
			}
		}
	}
	methodBody.If(
		jen.Err().Op(":=").Id(EncName).Dot("WriteToken").Call(snip.JSONV2EndObject()),
		jen.Err().Op("!=").Nil(),
	).Block(
		jen.Return(jen.Err()),
	)
	methodBody.Return(jen.Nil())
}

func GetCJMarshalerType(valueType types.Type) *jen.Statement {
	return getTypeArshaler(valueType, valueType.Code, false, false)
}

func MarshalJSONValue(
	encoder *jen.Statement,
	selector *jen.Statement,
	valueType types.Type,
) *jen.Statement {
	return snip.CJMarshalEncode().Types(valueType.Code(), GetCJMarshalerType(valueType)).Call(encoder, selector)
}

type jsonStructField struct {
	Key      string
	Type     types.Type
	Selector func() *jen.Statement
}

func GetCJUnmarshalerType(valueType types.Type) *jen.Statement {
	return getTypeArshaler(valueType, valueType.Code, false, true)
}

func UnmarshalJSONValue(
	decoder *jen.Statement,
	selector *jen.Statement,
	valueType types.Type,
) *jen.Statement {
	return snip.CJUnmarshalDecode().Types(valueType.Code(), GetCJUnmarshalerType(valueType)).Call(decoder, selector)
}

func MethodUnmarshalJSONFrom(receiverName string, receiverTypeName string, receiverType types.Type) *jen.Statement {
	return snip.MethodUnmarshalJSONFrom(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
		switch typ := receiverType.(type) {
		default:
			methodBody.Return(UnmarshalJSONValue(jen.Id(DecName), jen.Id(receiverName), receiverType))
		case *types.AliasType:
			methodBody.Return(UnmarshalJSONValue(jen.Id(DecName), aliasTypeItemSelector(typ, jen.Id(receiverName), true), typ.Item))
		case *types.EnumType:
			methodBody.If(
				jen.Err().Op(":=").Add(UnmarshalJSONValue(jen.Id(DecName), jen.Id(receiverName), typ)),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Err()),
			)
			methodBody.If(jen.Id(receiverName).Dot("IsUnknown").Call()).Block(
				jen.If(
					jen.List(jen.Id("strict"), jen.Op("_")).Op(":=").Add(snip.JSONV2GetOption()).Call(jen.Id(DecName).Dot("Options").Call(), snip.JSONV2RejectUnknownMembers()),
					jen.Id("strict"),
				).Block(
					jen.Return(snip.CJNewInvalidValueError().Call(
						jen.Id(DecName),
						jen.Lit(fmt.Sprintf("unknown %s value", receiverName)),
						jen.Nil(),
					)),
				),
			)
			methodBody.Return(jen.Nil())
		case *types.ObjectType:
			var fields []jsonStructField
			for _, field := range typ.Fields {
				fields = append(fields, jsonStructField{
					Key:      field.Name,
					Type:     field.Type,
					Selector: jen.Id(receiverName).Dot(transforms.ExportedFieldName(field.Name)).Clone,
				})
			}
			unmarshalJSONStructFields(methodBody, receiverName, receiverTypeName, fields, false)
			methodBody.Return(jen.Nil())
		case *types.UnionType:
			var fields []jsonStructField
			for _, field := range typ.Fields {
				fields = append(fields, jsonStructField{
					Key:      field.Name,
					Type:     field.Type,
					Selector: jen.Id(receiverName).Dot(transforms.PrivateFieldName(field.Name)).Clone,
				})
			}
			unmarshalJSONStructFields(methodBody, receiverName, receiverTypeName, fields, true)
			methodBody.Return(jen.Nil())
		}
	})
}

func unmarshalJSONStructFields(methodBody *jen.Group, receiverName string, receiverType string, fields []jsonStructField, isUnion bool) {
	var fieldResults []unmarshalJSONStructFieldResult
	hasRequiredFields := false
	hasCollections := false
	if isUnion {
		hasRequiredFields = true
		field := jsonStructField{
			Key:      "type",
			Type:     types.String{},
			Selector: jen.Id(receiverName).Dot("typ").Clone,
		}
		result := unmarshalJSONStructField(receiverName, receiverType, field, false, jen.Id(DecName))
		fieldResults = append(fieldResults, result)
	}
	for _, field := range fields {
		result := unmarshalJSONStructField(receiverName, receiverType, field, isUnion, jen.Id(DecName))
		if result.Validate != nil {
			hasRequiredFields = true
		}
		if result.DefaultCollection != nil {
			hasCollections = true
		}
		fieldResults = append(fieldResults, result)
	}

	methodBody.List(jen.Id("tok"), jen.Err()).Op(":=").Id(DecName).Dot("ReadToken").Call()
	methodBody.If(jen.Err().Op("!=").Nil()).Block(
		jen.Return(snip.CJWrapSyntaxError().Call(jen.Id(DecName), jen.Lit(""), jen.Err())),
	)
	methodBody.If(
		jen.Id("kind").Op(":=").Id("tok").Dot("Kind").Call(),
		jen.Id("kind").Op("!=").LitRune('{'),
	).Block(
		jen.Return(snip.CJNewKindMismatchError().Call(jen.Id(DecName), jen.Id("kind"), jen.Lit(receiverType+" opening brace"))),
	)

	for _, result := range fieldResults {
		if result.Init != nil {
			result.Init(methodBody)
		}
	}
	methodBody.Var().Id("unknownMembers").Index().String()

	methodBody.For().BlockFunc(func(forBody *jen.Group) {
		forBody.List(jen.Id("key"), jen.Err()).Op(":=").Id(DecName).Dot("ReadToken").Call()
		forBody.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(snip.CJWrapSyntaxError().Call(jen.Id(DecName), jen.Lit(""), jen.Err())),
		)
		forBody.Id("kind").Op(":=").Id("key").Dot("Kind").Call()
		forBody.If(jen.Id("kind").Op("==").LitRune('}')).Block(
			jen.Break(),
		)
		forBody.If(jen.Id("kind").Op("!=").LitRune('"')).Block(
			jen.Return(snip.CJNewKindMismatchError().Call(jen.Id(DecName), jen.Id("kind"), jen.Lit(receiverType+" closing brace or next key"))),
		)
		forBody.Switch(jen.Id("key").Dot("String").Call()).BlockFunc(func(cases *jen.Group) {
			for _, result := range fieldResults {
				if result.Unmarshal != nil {
					result.Unmarshal(cases)
				}
			}
			cases.Default().Block(
				jen.Id("unknownMembers").
					Op("=").
					Append(jen.Id("unknownMembers"), jen.Id("key").Dot("String").Call()),
				jen.If(
					jen.Err().Op(":=").Id(DecName).Dot("SkipValue").Call(),
					jen.Err().Op("!=").Nil(),
				).Block(
					jen.Return(jen.Err()),
				),
			)
		})
	})
	if hasRequiredFields {
		methodBody.Var().Id("missingFields").Index().String()
		for _, result := range fieldResults {
			if result.Validate != nil {
				result.Validate(methodBody)
			}
			if result.DefaultCollection != nil {
				result.DefaultCollection(methodBody)
			}
		}
		methodBody.If(jen.Len(jen.Id("missingFields")).Op(">").Lit(0)).Block(
			jen.Return(snip.CJNewMissingFieldsError().Call(
				jen.Id(DecName),
				jen.Lit(receiverType),
				jen.Id("missingFields"),
			)),
		)
	} else if hasCollections {
		for _, result := range fieldResults {
			if result.DefaultCollection != nil {
				result.DefaultCollection(methodBody)
			}
		}
	}
	methodBody.If(jen.Len(jen.Id("unknownMembers")).Op(">").Lit(0)).Block(
		jen.If(
			jen.List(jen.Id("strict"), jen.Op("_")).Op(":=").Add(snip.JSONV2GetOption()).Call(jen.Id(DecName).Dot("Options").Call(), snip.JSONV2RejectUnknownMembers()),
			jen.Id("strict"),
		).Block(
			jen.Return(snip.CJNewUnknownFieldsError().Call(
				jen.Id(DecName),
				jen.Lit(receiverType),
				jen.Id("unknownMembers"),
			)),
		),
	)
}

type unmarshalJSONStructFieldResult struct {
	Init              func(*jen.Group)
	Unmarshal         func(*jen.Group)
	Validate          func(*jen.Group)
	DefaultCollection func(*jen.Group)
}

func unmarshalJSONStructField(
	receiverName string,
	receiverType string,
	field jsonStructField,
	isUnionField bool,
	decoder *jen.Statement,
) (result unmarshalJSONStructFieldResult) {
	seenVar := "seen" + transforms.ExportedFieldName(field.Key)
	result.Init = func(methodBody *jen.Group) {
		methodBody.Var().Id(seenVar).Bool()
	}
	requiredField := !(field.Type.IsCollection() || field.Type.IsOptional())
	if requiredField {
		result.Validate = func(methodBody *jen.Group) {
			methodBody.IfFunc(func(conds *jen.Group) {
				if isUnionField {
					conds.Id(receiverName).Dot("typ").Op("==").Lit(field.Key).
						Op("&&").
						Op("!").Id(seenVar)
				} else {
					conds.Op("!").Id(seenVar)
				}
			}).Block(
				jen.Id("missingFields").Op("=").Append(jen.Id("missingFields"), jen.Lit(field.Key)),
			)
		}
	} else if mk := field.Type.Make(); mk != nil && !isUnionField {
		result.DefaultCollection = func(methodBody *jen.Group) {
			methodBody.If(jen.Op("!").Id(seenVar)).Block(
				field.Selector().Op("=").Add(mk),
			)
		}
	}
	result.Unmarshal = func(cases *jen.Group) {
		cases.Case(jen.Lit(field.Key)).BlockFunc(func(caseBody *jen.Group) {
			caseBody.If(jen.Id(seenVar)).Block(
				jen.Return(snip.CJNewDuplicateFieldKeyError().Call(
					decoder,
					jen.Lit(fmt.Sprintf("%s[%q]", receiverType, field.Key)),
				)),
			)

			var selector *jen.Statement
			if isUnionField {
				caseBody.Add(field.Selector()).Op("=").New(field.Type.Code())
				selector = field.Selector()
			} else {
				selector = jen.Op("&").Add(field.Selector())
			}
			caseBody.If(
				jen.Err().Op(":=").Add(UnmarshalJSONValue(decoder, selector, field.Type)),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(snip.CJNewUnmarshalFieldError().Call(jen.Id(DecName), jen.Lit(fmt.Sprintf("%s[%q]", receiverType, field.Key)), jen.Err())),
			)

			caseBody.Id(seenVar).Op("=").True()
		})
	}

	return result
}

func aliasTypeItemSelector(typ *types.AliasType, selector *jen.Statement, isUnmarshal bool) *jen.Statement {
	selector = selector.Clone()
	if typ.IsOptional() {
		if isUnmarshal {
			return jen.Op("&").Add(selector).Dot("Value")
		}
		return selector.Dot("Value")
	}
	if isUnmarshal {
		return jen.Parens(jen.Op("*").Add(typ.Item.Code())).Call(selector)
	}
	if typ.IsCollection() {
		return selector
	}
	return typ.Item.Code().Call(selector)
}

func getTypeArshaler(valueType types.Type, declType func() *jen.Statement, isMapKey bool, isUnmarshal bool) *jen.Statement {
	switch typ := valueType.(type) {
	case types.Any:
		return snip.CJAny().Types(declType())
	case types.String:
		return snip.CJString().Types(declType())
	case types.Bearertoken:
		return snip.CJBearerToken().Types(declType())
	case types.DateTime:
		// TODO: It is not possible to use a generic constraint ~time.Time, so just use the text methods instead.
		// This will make sorting map keys to marshal a little more expensive since it has to compare the string representations.
		if isUnmarshal {
			return snip.CJTextUnmarshaler().Types(jen.Op("*").Add(declType()))
		}
		return snip.CJStringerMarshaler().Types(declType())
	case types.RID:
		return snip.CJRID().Types(declType())
	case types.UUID:
		return snip.CJUUID().Types(declType())
	case types.Boolean:
		if isMapKey {
			return snip.CJBooleanMapKey().Types(declType())
		}
		return snip.CJBoolean().Types(declType())
	case types.Double:
		if isMapKey {
			return snip.CJFloatMapKey().Types(declType())
		}
		return snip.CJFloat().Types(declType())
	case types.Integer:
		if isMapKey {
			return snip.CJInt32MapKey().Types(declType())
		}
		return snip.CJInt32().Types(declType())
	case types.Safelong:
		if isMapKey {
			return snip.CJSafeLongMapKey().Types(declType())
		}
		return snip.CJSafeLong().Types(declType())
	case types.Binary:
		if isMapKey {
			return snip.CJBinaryMapKey().Types(declType())
		}
		return snip.CJBinary().Types(declType())
	case *types.Optional:
		if isUnmarshal {
			return snip.CJOptionalUnmarshaler().Types(declType(), typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code, false, isUnmarshal))
		}
		return snip.CJOptionalMarshaler().Types(declType(), typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code, false, isUnmarshal))
	case *types.List:
		if isUnmarshal {
			return snip.CJListUnmarshaler().Types(declType(), typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code, false, isUnmarshal))
		}
		return snip.CJListMarshaler().Types(declType(), typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code, false, isUnmarshal))
	case *types.Set:
		if typ.Item.IsComparable() {
			if isUnmarshal {
				return snip.CJSetUnmarshaler().Types(declType(), typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code, false, isUnmarshal))
			}
			return snip.CJSetMarshaler().Types(declType(), typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code, false, isUnmarshal))
		}
		// TODO: Add an UncomparableSet type that takes some kind of Comparator type parameter to do the collision checks.
		if isUnmarshal {
			return snip.CJListUnmarshaler().Types(declType(), typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code, false, isUnmarshal))
		}
		return snip.CJListMarshaler().Types(declType(), typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code, false, isUnmarshal))
	case *types.Map:
		var keyType *jen.Statement
		if typ.Key.IsBinary() {
			keyType = snip.BinaryBinary()
		} else if typ.Key.IsBoolean() {
			keyType = snip.BooleanBoolean()
		} else {
			keyType = typ.Key.Code()
		}
		key := getTypeArshaler(typ.Key, keyType.Clone, true, isUnmarshal)
		val := getTypeArshaler(typ.Val, typ.Val.Code, false, isUnmarshal)
		typeArgs := jen.Types(declType(), keyType, typ.Val.Code(), key, val)
		switch {
		case isUnmarshal:
			return snip.CJMapUnmarshaler().Add(typeArgs)
		case typ.Key.IsOrdered():
			return snip.CJOrderedMapMarshaler().Add(typeArgs)
		default:
			return snip.CJComparableMapMarshaler().Add(typeArgs)
		}
	case *types.External:
		if typ.ExternalHasGoType() {
			return snip.CJAny().Types(declType())
		}
		return getTypeArshaler(typ.Fallback, declType, isMapKey, isUnmarshal)
	case *types.AliasType:
		if typ.Item.IsOptional() {
			if isUnmarshal {
				return snip.CJStructUnmarshaler().Types(jen.Op("*").Add(declType()))
			}
			return snip.CJStructMarshaler().Types(declType())
		}
		return getTypeArshaler(typ.Item, declType, isMapKey, isUnmarshal)
	case *types.EnumType:
		if isUnmarshal {
			return snip.CJTextUnmarshaler().Types(jen.Op("*").Add(declType()))
		}
		return snip.CJStringerMarshaler().Types(declType())
	case *types.ObjectType, *types.UnionType:
		if isUnmarshal {
			return snip.CJStructUnmarshaler().Types(jen.Op("*").Add(declType()))
		}
		return snip.CJStructMarshaler().Types(declType())
	default:
		panic(fmt.Sprintf("unknown type %T", typ))
	}
}
