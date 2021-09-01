package encoding

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

func UnmarshalStringStatements(
	methodBody *jen.Group,
	argType types.Type,
	outVarName string,
	inStr func() *jen.Statement,
	argName string,
	ctxVar func() *jen.Statement,
) {
	unmarshalStringStatements(methodBody, argType, outVarName, inStr, argName, ctxVar, 0)
}

func unmarshalStringStatements(
	methodBody *jen.Group,
	argType types.Type,
	outVarName string,
	inStr func() *jen.Statement,
	argName string,
	ctxVar func() *jen.Statement,
	nestDepth int,
) {
	var (
		// Simple types can reuse the assignment logic at the end of this function by setting these variables
		expr       func() *jen.Statement
		returnsVal bool
		returnsErr bool
	)
	switch typVal := argType.(type) {
	case types.Any, types.String:
		expr = inStr
		returnsVal = true
	case types.Bearertoken:
		expr = snip.BearerTokenToken().Call(inStr()).Clone
		returnsVal = true
	case types.Binary:
		expr = jen.Id("[]byte").Call(inStr()).Clone
		returnsVal = true
	case types.Boolean:
		expr = snip.StrconvParseBool().Call(inStr()).Clone
		returnsVal = true
		returnsErr = true
	case types.DateTime:
		expr = snip.DateTimeParseDateTime().Call(inStr()).Clone
		returnsVal = true
		returnsErr = true
	case types.Double:
		expr = snip.StrconvParseFloat().Call(inStr(), jen.Lit(64)).Clone
		returnsVal = true
		returnsErr = true
	case types.Integer:
		expr = snip.StrconvAtoi().Call(inStr()).Clone
		returnsVal = true
		returnsErr = true
	case types.RID:
		expr = snip.RIDParseRID().Call(inStr()).Clone
		returnsVal = true
		returnsErr = true
	case types.Safelong:
		expr = snip.SafeLongParseSafeLong().Call(inStr()).Clone
		returnsVal = true
		returnsErr = true
	case types.UUID:
		expr = snip.UUIDParseUUID().Call(inStr()).Clone
		returnsVal = true
		returnsErr = true

	case *types.Optional:
		// declare output variable
		strVar := tmpVarName(outVarName+"Str", nestDepth)
		valVar := tmpVarName(outVarName+"Internal", nestDepth)
		methodBody.Var().Id(outVarName).Add(typVal.Code())
		methodBody.If(
			jen.Id(strVar).Op(":=").Add(inStr()),
			jen.Id(strVar).Op("!=").Lit(""),
		).BlockFunc(func(ifBody *jen.Group) {
			unmarshalStringStatements(ifBody, typVal.Item, valVar, jen.Id(strVar).Clone, argName, ctxVar, nestDepth+1)
			ifBody.Id(outVarName).Op("=").Op("&").Id(valVar)
		})
	case *types.AliasType:
		if typVal.IsText() {
			expr = jen.Id(outVarName).Dot("UnmarshalString").Call(inStr()).Clone
		} else {
			expr = jen.Id(outVarName).Dot("UnmarshalJSONString").Call(snip.SafeJSONQuoteString().Call(inStr())).Clone
		}
		returnsErr = true
	case *types.EnumType:
		expr = jen.Id(outVarName).Dot("UnmarshalString").Call(inStr()).Clone
		returnsErr = true
	case *types.List, *types.Map, *types.ObjectType, *types.UnionType:
		panic(fmt.Sprintf("unsupported complex type for string param %v", argType))
	default:
		panic(fmt.Sprintf("unrecognized type %v", argType))
	}

	if expr != nil {
		if !returnsErr {
			methodBody.Id(outVarName).Op(":=").Add(expr())
		} else {
			if !returnsVal {
				methodBody.Var().Id(outVarName).Add(argType.Code())
			} else {
				methodBody.List(jen.Id(outVarName), jen.Err()).Op(":=").Add(expr())
			}
			methodBody.IfFunc(func(conds *jen.Group) {
				if !returnsVal {
					conds.Err().Op(":=").Add(expr())
				}
				conds.Err().Op("!=").Nil()
			}).Block(
				jen.Return(snip.WerrorWrapContext().Call(
					ctxVar(),
					snip.CGRErrorsWrapWithInvalidArgument().Call(jen.Err()),
					jen.Lit(fmt.Sprintf("unmarshal %s as %s", argName, argType)),
				)),
			)
		}
	}
}

func UnmarshalStringListStatements(
	methodBody *jen.Group,
	listType types.Type,
	outVarName string,
	inStrs func() *jen.Statement,
	argName string,
	ctxVar func() *jen.Statement,
) {
	methodBody.Var().Id(outVarName).Add(listType.Code())
	unmarshalStringListStatements(methodBody, listType, outVarName, inStrs, argName, ctxVar)
}

func unmarshalStringListStatements(
	methodBody *jen.Group,
	listType types.Type,
	outVarName string,
	inStrs func() *jen.Statement,
	argName string,
	ctxVar func() *jen.Statement,
) {
	switch typVal := listType.(type) {
	case *types.Optional:
		unmarshalStringListStatements(methodBody, typVal.Item, outVarName, inStrs, argName, ctxVar)
	case *types.AliasType:
		unmarshalStringListStatements(methodBody, typVal.Item, outVarName, inStrs, argName, ctxVar)
	case *types.List:
		if _, isString := typVal.Item.(types.String); isString {
			methodBody.Id(outVarName).Op("=").Add(inStrs())
			return
		}

		methodBody.For(jen.List(jen.Id("_"), jen.Id("v")).Op(":=").Range().Add(inStrs())).BlockFunc(func(rangeBody *jen.Group) {
			UnmarshalStringStatements(rangeBody, typVal.Item, "convertedVal", jen.Id("v").Clone, argName, ctxVar)
			rangeBody.Id(outVarName).Op("=").Append(jen.Id(outVarName), jen.Id("convertedVal"))
		})
	default:
		panic(fmt.Sprintf("unsupported complex type for string list param %v", listType))
	}
}
