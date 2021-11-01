package server

import (
	"context"
	"fmt"
	types "github.com/palantir/conjure-go/v6/conjure-go-verifier/conjure/verification/types"
)

type AutoDeserializeServiceImpl struct {
	TestCases map[EndpointName]PositiveAndNegativeTestCases
}

func (a AutoDeserializeServiceImpl) ReceiveBearerTokenExample(_ context.Context, indexArg int) (types.BearerTokenExample, error) {
	var value types.BearerTokenExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveBearerTokenExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveBinaryExample(_ context.Context, indexArg int) (types.BinaryExample, error) {
	var value types.BinaryExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveBinaryExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveBooleanExample(_ context.Context, indexArg int) (types.BooleanExample, error) {
	var value types.BooleanExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveBooleanExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveDateTimeExample(_ context.Context, indexArg int) (types.DateTimeExample, error) {
	var value types.DateTimeExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveDateTimeExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveDoubleExample(_ context.Context, indexArg int) (types.DoubleExample, error) {
	var value types.DoubleExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveDoubleExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveIntegerExample(_ context.Context, indexArg int) (types.IntegerExample, error) {
	var value types.IntegerExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveIntegerExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveRidExample(_ context.Context, indexArg int) (types.RidExample, error) {
	var value types.RidExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveRidExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSafeLongExample(_ context.Context, indexArg int) (types.SafeLongExample, error) {
	var value types.SafeLongExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSafeLongExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveStringExample(_ context.Context, indexArg int) (types.StringExample, error) {
	var value types.StringExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveStringExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveUuidExample(_ context.Context, indexArg int) (types.UuidExample, error) {
	var value types.UuidExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveUuidExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveAnyExample(_ context.Context, indexArg int) (types.AnyExample, error) {
	var value types.AnyExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveAnyExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveEnumExample(_ context.Context, indexArg int) (types.EnumExample, error) {
	var value types.EnumExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveEnumExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveListExample(_ context.Context, indexArg int) (types.ListExample, error) {
	var value types.ListExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveListExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSetStringExample(_ context.Context, indexArg int) (types.SetStringExample, error) {
	var value types.SetStringExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSetStringExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSetDoubleExample(_ context.Context, indexArg int) (types.SetDoubleExample, error) {
	var value types.SetDoubleExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSetDoubleExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveMapExample(_ context.Context, indexArg int) (types.MapExample, error) {
	var value types.MapExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveMapExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveOptionalExample(_ context.Context, indexArg int) (types.OptionalExample, error) {
	var value types.OptionalExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveOptionalExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveOptionalBooleanExample(_ context.Context, indexArg int) (types.OptionalBooleanExample, error) {
	var value types.OptionalBooleanExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveOptionalBooleanExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveOptionalIntegerExample(_ context.Context, indexArg int) (types.OptionalIntegerExample, error) {
	var value types.OptionalIntegerExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveOptionalIntegerExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveLongFieldNameOptionalExample(_ context.Context, indexArg int) (types.LongFieldNameOptionalExample, error) {
	var value types.LongFieldNameOptionalExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveLongFieldNameOptionalExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveRawOptionalExample(_ context.Context, indexArg int) (types.RawOptionalExample, error) {
	var value types.RawOptionalExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveRawOptionalExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveStringAliasExample(_ context.Context, indexArg int) (types.StringAliasExample, error) {
	var value types.StringAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveStringAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveDoubleAliasExample(_ context.Context, indexArg int) (types.DoubleAliasExample, error) {
	var value types.DoubleAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveDoubleAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveIntegerAliasExample(_ context.Context, indexArg int) (types.IntegerAliasExample, error) {
	var value types.IntegerAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveIntegerAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveBooleanAliasExample(_ context.Context, indexArg int) (types.BooleanAliasExample, error) {
	var value types.BooleanAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveBooleanAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSafeLongAliasExample(_ context.Context, indexArg int) (types.SafeLongAliasExample, error) {
	var value types.SafeLongAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSafeLongAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveRidAliasExample(_ context.Context, indexArg int) (types.RidAliasExample, error) {
	var value types.RidAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveRidAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveBearerTokenAliasExample(_ context.Context, indexArg int) (types.BearerTokenAliasExample, error) {
	var value types.BearerTokenAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveBearerTokenAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveUuidAliasExample(_ context.Context, indexArg int) (types.UuidAliasExample, error) {
	var value types.UuidAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveUuidAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveReferenceAliasExample(_ context.Context, indexArg int) (types.ReferenceAliasExample, error) {
	var value types.ReferenceAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveReferenceAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveDateTimeAliasExample(_ context.Context, indexArg int) (types.DateTimeAliasExample, error) {
	var value types.DateTimeAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveDateTimeAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveBinaryAliasExample(_ context.Context, indexArg int) (types.BinaryAliasExample, error) {
	var value types.BinaryAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveBinaryAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveKebabCaseObjectExample(_ context.Context, indexArg int) (types.KebabCaseObjectExample, error) {
	var value types.KebabCaseObjectExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveKebabCaseObjectExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSnakeCaseObjectExample(_ context.Context, indexArg int) (types.SnakeCaseObjectExample, error) {
	var value types.SnakeCaseObjectExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSnakeCaseObjectExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveOptionalBearerTokenAliasExample(_ context.Context, indexArg int) (types.OptionalBearerTokenAliasExample, error) {
	var value types.OptionalBearerTokenAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveOptionalBearerTokenAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveOptionalBooleanAliasExample(_ context.Context, indexArg int) (types.OptionalBooleanAliasExample, error) {
	var value types.OptionalBooleanAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveOptionalBooleanAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveOptionalDateTimeAliasExample(_ context.Context, indexArg int) (types.OptionalDateTimeAliasExample, error) {
	var value types.OptionalDateTimeAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveOptionalDateTimeAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveOptionalDoubleAliasExample(_ context.Context, indexArg int) (types.OptionalDoubleAliasExample, error) {
	var value types.OptionalDoubleAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveOptionalDoubleAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveOptionalIntegerAliasExample(_ context.Context, indexArg int) (types.OptionalIntegerAliasExample, error) {
	var value types.OptionalIntegerAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveOptionalIntegerAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveOptionalRidAliasExample(_ context.Context, indexArg int) (types.OptionalRidAliasExample, error) {
	var value types.OptionalRidAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveOptionalRidAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveOptionalSafeLongAliasExample(_ context.Context, indexArg int) (types.OptionalSafeLongAliasExample, error) {
	var value types.OptionalSafeLongAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveOptionalSafeLongAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveOptionalStringAliasExample(_ context.Context, indexArg int) (types.OptionalStringAliasExample, error) {
	var value types.OptionalStringAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveOptionalStringAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveOptionalUuidAliasExample(_ context.Context, indexArg int) (types.OptionalUuidAliasExample, error) {
	var value types.OptionalUuidAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveOptionalUuidAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveOptionalAnyAliasExample(_ context.Context, indexArg int) (types.OptionalAnyAliasExample, error) {
	var value types.OptionalAnyAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveOptionalAnyAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveListBearerTokenAliasExample(_ context.Context, indexArg int) (types.ListBearerTokenAliasExample, error) {
	var value types.ListBearerTokenAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveListBearerTokenAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveListBinaryAliasExample(_ context.Context, indexArg int) (types.ListBinaryAliasExample, error) {
	var value types.ListBinaryAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveListBinaryAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveListBooleanAliasExample(_ context.Context, indexArg int) (types.ListBooleanAliasExample, error) {
	var value types.ListBooleanAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveListBooleanAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveListDateTimeAliasExample(_ context.Context, indexArg int) (types.ListDateTimeAliasExample, error) {
	var value types.ListDateTimeAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveListDateTimeAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveListDoubleAliasExample(_ context.Context, indexArg int) (types.ListDoubleAliasExample, error) {
	var value types.ListDoubleAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveListDoubleAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveListIntegerAliasExample(_ context.Context, indexArg int) (types.ListIntegerAliasExample, error) {
	var value types.ListIntegerAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveListIntegerAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveListRidAliasExample(_ context.Context, indexArg int) (types.ListRidAliasExample, error) {
	var value types.ListRidAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveListRidAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveListSafeLongAliasExample(_ context.Context, indexArg int) (types.ListSafeLongAliasExample, error) {
	var value types.ListSafeLongAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveListSafeLongAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveListStringAliasExample(_ context.Context, indexArg int) (types.ListStringAliasExample, error) {
	var value types.ListStringAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveListStringAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveListUuidAliasExample(_ context.Context, indexArg int) (types.ListUuidAliasExample, error) {
	var value types.ListUuidAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveListUuidAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveListAnyAliasExample(_ context.Context, indexArg int) (types.ListAnyAliasExample, error) {
	var value types.ListAnyAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveListAnyAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveListOptionalAnyAliasExample(_ context.Context, indexArg int) (types.ListOptionalAnyAliasExample, error) {
	var value types.ListOptionalAnyAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveListOptionalAnyAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSetBearerTokenAliasExample(_ context.Context, indexArg int) (types.SetBearerTokenAliasExample, error) {
	var value types.SetBearerTokenAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSetBearerTokenAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSetBinaryAliasExample(_ context.Context, indexArg int) (types.SetBinaryAliasExample, error) {
	var value types.SetBinaryAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSetBinaryAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSetBooleanAliasExample(_ context.Context, indexArg int) (types.SetBooleanAliasExample, error) {
	var value types.SetBooleanAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSetBooleanAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSetDateTimeAliasExample(_ context.Context, indexArg int) (types.SetDateTimeAliasExample, error) {
	var value types.SetDateTimeAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSetDateTimeAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSetDoubleAliasExample(_ context.Context, indexArg int) (types.SetDoubleAliasExample, error) {
	var value types.SetDoubleAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSetDoubleAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSetIntegerAliasExample(_ context.Context, indexArg int) (types.SetIntegerAliasExample, error) {
	var value types.SetIntegerAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSetIntegerAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSetRidAliasExample(_ context.Context, indexArg int) (types.SetRidAliasExample, error) {
	var value types.SetRidAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSetRidAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSetSafeLongAliasExample(_ context.Context, indexArg int) (types.SetSafeLongAliasExample, error) {
	var value types.SetSafeLongAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSetSafeLongAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSetStringAliasExample(_ context.Context, indexArg int) (types.SetStringAliasExample, error) {
	var value types.SetStringAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSetStringAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSetUuidAliasExample(_ context.Context, indexArg int) (types.SetUuidAliasExample, error) {
	var value types.SetUuidAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSetUuidAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSetAnyAliasExample(_ context.Context, indexArg int) (types.SetAnyAliasExample, error) {
	var value types.SetAnyAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSetAnyAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveSetOptionalAnyAliasExample(_ context.Context, indexArg int) (types.SetOptionalAnyAliasExample, error) {
	var value types.SetOptionalAnyAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveSetOptionalAnyAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveMapBearerTokenAliasExample(_ context.Context, indexArg int) (types.MapBearerTokenAliasExample, error) {
	var value types.MapBearerTokenAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveMapBearerTokenAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveMapBinaryAliasExample(_ context.Context, indexArg int) (types.MapBinaryAliasExample, error) {
	var value types.MapBinaryAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveMapBinaryAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveMapBooleanAliasExample(_ context.Context, indexArg int) (types.MapBooleanAliasExample, error) {
	var value types.MapBooleanAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveMapBooleanAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveMapDateTimeAliasExample(_ context.Context, indexArg int) (types.MapDateTimeAliasExample, error) {
	var value types.MapDateTimeAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveMapDateTimeAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveMapDoubleAliasExample(_ context.Context, indexArg int) (types.MapDoubleAliasExample, error) {
	var value types.MapDoubleAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveMapDoubleAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveMapIntegerAliasExample(_ context.Context, indexArg int) (types.MapIntegerAliasExample, error) {
	var value types.MapIntegerAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveMapIntegerAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveMapRidAliasExample(_ context.Context, indexArg int) (types.MapRidAliasExample, error) {
	var value types.MapRidAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveMapRidAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveMapSafeLongAliasExample(_ context.Context, indexArg int) (types.MapSafeLongAliasExample, error) {
	var value types.MapSafeLongAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveMapSafeLongAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveMapStringAliasExample(_ context.Context, indexArg int) (types.MapStringAliasExample, error) {
	var value types.MapStringAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveMapStringAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveMapUuidAliasExample(_ context.Context, indexArg int) (types.MapUuidAliasExample, error) {
	var value types.MapUuidAliasExample
	err := value.UnmarshalJSON(a.testCaseBytes("receiveMapUuidAliasExample", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) ReceiveMapEnumExampleAlias(_ context.Context, indexArg int) (types.MapEnumExampleAlias, error) {
	var value types.MapEnumExampleAlias
	err := value.UnmarshalJSON(a.testCaseBytes("receiveMapEnumExampleAlias", indexArg))
	return value, err
}
func (a AutoDeserializeServiceImpl) testCaseBytes(endpointName EndpointName, i int) []byte {
	cases := a.TestCases[endpointName]
	posLen, negLen := len(cases.Positive), len(cases.Negative)
	if i < posLen {
		return []byte(cases.Positive[i])
	} else if i < posLen+negLen {
		return []byte(cases.Negative[i-posLen])
	}
	panic(fmt.Sprintf("invalid test case index %s[%d]", endpointName, i))
}
