package jsonv2

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

func getTypeArshaler(valueType types.Type, declType *jen.Statement, isMapKey bool, isUnmarshal bool) *jen.Statement {
	switch typ := valueType.(type) {
	case types.Any:
		return snip.CJTypeAny().Types(declType)
	case types.String, types.Bearertoken:
		return snip.CJTypeString().Types(declType)
	case types.DateTime:
		return snip.CJTypeDateTime().Types(declType)
	case types.RID:
		return snip.CJTypeRID().Types(declType)
	case types.UUID:
		return snip.CJTypeUUID().Types(declType)
	case types.Boolean:
		if isMapKey {
			return snip.CJTypeBooleanMapKey().Types(declType)
		}
		return snip.CJTypeBoolean().Types(declType)
	case types.Double:
		if isMapKey {
			return snip.CJTypeFloatMapKey().Types(declType)
		}
		return snip.CJTypeFloat().Types(declType)
	case types.Integer, types.Safelong:
		if isMapKey {
			return snip.CJTypeIntMapKey().Types(declType)
		}
		return snip.CJTypeInt().Types(declType)
	case types.Binary:
		if isMapKey {
			return snip.CJTypeString().Types(snip.BinaryBinary())
		}
		return snip.CJTypeBinary().Types(declType)
	case *types.Optional:
		if isUnmarshal {
			return snip.CJTypeOptionalUnmarshaler().Types(declType, typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code(), isMapKey, isUnmarshal))
		}
		return snip.CJTypeOptionalMarshaler().Types(declType, typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code(), isMapKey, isUnmarshal))
	case *types.List:
		if isUnmarshal {
			return snip.CJTypeListUnmarshaler().Types(declType, typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code(), false, isUnmarshal))
		}
		return snip.CJTypeListMarshaler().Types(declType, typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code(), false, isUnmarshal))
	case *types.Set:
		if isUnmarshal {
			return snip.CJTypeListUnmarshaler().Types(declType, typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code(), false, isUnmarshal))
		}
		return snip.CJTypeListMarshaler().Types(declType, typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code(), false, isUnmarshal))
	case *types.Map:
		var keyType *jen.Statement
		if typ.Key.IsBinary() {
			keyType = snip.BinaryBinary()
		} else if typ.Key.IsBoolean() {
			keyType = snip.BooleanBoolean()
		} else {
			keyType = typ.Key.Code()
		}
		key := getTypeArshaler(typ.Key, keyType, true, isUnmarshal)
		val := getTypeArshaler(typ.Val, typ.Val.Code(), false, isUnmarshal)
		typeArgs := jen.Types(declType, keyType, typ.Val.Code(), key, val)
		switch {
		case isUnmarshal:
			return snip.CJTypeMapUnmarshaler().Add(typeArgs)
		case typ.Key.IsOrdered():
			return snip.CJTypeOrderedMapMarshaler().Add(typeArgs)
		default:
			return snip.CJTypeComparableMapMarshaler().Add(typeArgs)
		}
	case *types.External:
		if typ.ExternalHasGoType() {
			return snip.CJTypeAny().Types(declType)
		}
		return getTypeArshaler(typ.Fallback, declType, isMapKey, isUnmarshal)
	case *types.AliasType:
		if typ.IsOptional() {
			if isUnmarshal {
				return snip.CJTypeStructUnmarshaler().Types(declType)
			}
			return snip.CJTypeStructMarshaler().Types(declType)
		}
		return getTypeArshaler(typ.Item, declType, isMapKey, isUnmarshal)
	case *types.EnumType:
		if isUnmarshal {
			return snip.CJTypeTextUnmarshaler().Types(jen.Op("*").Add(declType))
		}
		return snip.CJTypeStringerMarshaler().Types(declType)
	case *types.ObjectType, *types.UnionType:
		if isUnmarshal {
			return snip.CJTypeStructUnmarshaler().Types(jen.Op("*").Add(declType))
		}
		return snip.CJTypeStructMarshaler().Types(declType)
	default:
		panic(fmt.Sprintf("unknown type %T", typ))
	}
}
