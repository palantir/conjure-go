package server

import (
	"context"
	"fmt"
	"io"

	"github.com/palantir/conjure-go/v6/conjure-go-verifier/conjure/verification/server"
	"github.com/palantir/conjure-go/v6/conjure-go-verifier/conjure/verification/types"
)

type AutoDeserializeImpl struct {
	TestCases map[server.EndpointName]server.PositiveAndNegativeTestCases
}

func (a AutoDeserializeImpl) testCaseBytes(endpointName server.EndpointName, indexArg int) []byte {
	cases := a.TestCases[endpointName]
	posLen, negLen := len(cases.Positive), len(cases.Negative)
	if indexArg < posLen {
		return []byte(cases.Positive[indexArg])
	} else if indexArg < posLen+negLen {
		return []byte(cases.Negative[indexArg-posLen])
	}
	panic(fmt.Sprintf("invalid test case index %s[%d]", endpointName, indexArg))
}

func (a AutoDeserializeImpl) ReceiveBearerTokenExample(ctx context.Context, indexArg int) (retVal types.BearerTokenExample, err error) {
	err = retVal.UnmarshalJSON(a.testCaseBytes("receiveBearerTokenExample", indexArg))
	return
}

func (AutoDeserializeImpl) ReceiveBinaryExample(ctx context.Context, indexArg int) (types.BinaryExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveBooleanExample(ctx context.Context, indexArg int) (types.BooleanExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveDateTimeExample(ctx context.Context, indexArg int) (types.DateTimeExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveDoubleExample(ctx context.Context, indexArg int) (types.DoubleExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveIntegerExample(ctx context.Context, indexArg int) (types.IntegerExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveRidExample(ctx context.Context, indexArg int) (types.RidExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSafeLongExample(ctx context.Context, indexArg int) (types.SafeLongExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveStringExample(ctx context.Context, indexArg int) (types.StringExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveUuidExample(ctx context.Context, indexArg int) (types.UuidExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveAnyExample(ctx context.Context, indexArg int) (types.AnyExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveEnumExample(ctx context.Context, indexArg int) (types.EnumExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveListExample(ctx context.Context, indexArg int) (types.ListExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSetStringExample(ctx context.Context, indexArg int) (types.SetStringExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSetDoubleExample(ctx context.Context, indexArg int) (types.SetDoubleExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveMapExample(ctx context.Context, indexArg int) (types.MapExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveOptionalExample(ctx context.Context, indexArg int) (types.OptionalExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveOptionalBooleanExample(ctx context.Context, indexArg int) (types.OptionalBooleanExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveOptionalIntegerExample(ctx context.Context, indexArg int) (types.OptionalIntegerExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveLongFieldNameOptionalExample(ctx context.Context, indexArg int) (types.LongFieldNameOptionalExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveRawOptionalExample(ctx context.Context, indexArg int) (types.RawOptionalExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveStringAliasExample(ctx context.Context, indexArg int) (types.StringAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveDoubleAliasExample(ctx context.Context, indexArg int) (types.DoubleAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveIntegerAliasExample(ctx context.Context, indexArg int) (types.IntegerAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveBooleanAliasExample(ctx context.Context, indexArg int) (types.BooleanAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSafeLongAliasExample(ctx context.Context, indexArg int) (types.SafeLongAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveRidAliasExample(ctx context.Context, indexArg int) (types.RidAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveBearerTokenAliasExample(ctx context.Context, indexArg int) (types.BearerTokenAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveUuidAliasExample(ctx context.Context, indexArg int) (types.UuidAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveReferenceAliasExample(ctx context.Context, indexArg int) (types.ReferenceAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveDateTimeAliasExample(ctx context.Context, indexArg int) (types.DateTimeAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveBinaryAliasExample(ctx context.Context, indexArg int) (io.ReadCloser, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveKebabCaseObjectExample(ctx context.Context, indexArg int) (types.KebabCaseObjectExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSnakeCaseObjectExample(ctx context.Context, indexArg int) (types.SnakeCaseObjectExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveOptionalBearerTokenAliasExample(ctx context.Context, indexArg int) (types.OptionalBearerTokenAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveOptionalBooleanAliasExample(ctx context.Context, indexArg int) (types.OptionalBooleanAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveOptionalDateTimeAliasExample(ctx context.Context, indexArg int) (types.OptionalDateTimeAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveOptionalDoubleAliasExample(ctx context.Context, indexArg int) (types.OptionalDoubleAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveOptionalIntegerAliasExample(ctx context.Context, indexArg int) (types.OptionalIntegerAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveOptionalRidAliasExample(ctx context.Context, indexArg int) (types.OptionalRidAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveOptionalSafeLongAliasExample(ctx context.Context, indexArg int) (types.OptionalSafeLongAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveOptionalStringAliasExample(ctx context.Context, indexArg int) (types.OptionalStringAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveOptionalUuidAliasExample(ctx context.Context, indexArg int) (types.OptionalUuidAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveOptionalAnyAliasExample(ctx context.Context, indexArg int) (types.OptionalAnyAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveListBearerTokenAliasExample(ctx context.Context, indexArg int) (types.ListBearerTokenAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveListBinaryAliasExample(ctx context.Context, indexArg int) (types.ListBinaryAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveListBooleanAliasExample(ctx context.Context, indexArg int) (types.ListBooleanAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveListDateTimeAliasExample(ctx context.Context, indexArg int) (types.ListDateTimeAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveListDoubleAliasExample(ctx context.Context, indexArg int) (types.ListDoubleAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveListIntegerAliasExample(ctx context.Context, indexArg int) (types.ListIntegerAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveListRidAliasExample(ctx context.Context, indexArg int) (types.ListRidAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveListSafeLongAliasExample(ctx context.Context, indexArg int) (types.ListSafeLongAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveListStringAliasExample(ctx context.Context, indexArg int) (types.ListStringAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveListUuidAliasExample(ctx context.Context, indexArg int) (types.ListUuidAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveListAnyAliasExample(ctx context.Context, indexArg int) (types.ListAnyAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveListOptionalAnyAliasExample(ctx context.Context, indexArg int) (types.ListOptionalAnyAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSetBearerTokenAliasExample(ctx context.Context, indexArg int) (types.SetBearerTokenAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSetBinaryAliasExample(ctx context.Context, indexArg int) (types.SetBinaryAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSetBooleanAliasExample(ctx context.Context, indexArg int) (types.SetBooleanAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSetDateTimeAliasExample(ctx context.Context, indexArg int) (types.SetDateTimeAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSetDoubleAliasExample(ctx context.Context, indexArg int) (types.SetDoubleAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSetIntegerAliasExample(ctx context.Context, indexArg int) (types.SetIntegerAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSetRidAliasExample(ctx context.Context, indexArg int) (types.SetRidAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSetSafeLongAliasExample(ctx context.Context, indexArg int) (types.SetSafeLongAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSetStringAliasExample(ctx context.Context, indexArg int) (types.SetStringAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSetUuidAliasExample(ctx context.Context, indexArg int) (types.SetUuidAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSetAnyAliasExample(ctx context.Context, indexArg int) (types.SetAnyAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveSetOptionalAnyAliasExample(ctx context.Context, indexArg int) (types.SetOptionalAnyAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveMapBearerTokenAliasExample(ctx context.Context, indexArg int) (types.MapBearerTokenAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveMapBinaryAliasExample(ctx context.Context, indexArg int) (types.MapBinaryAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveMapBooleanAliasExample(ctx context.Context, indexArg int) (types.MapBooleanAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveMapDateTimeAliasExample(ctx context.Context, indexArg int) (types.MapDateTimeAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveMapDoubleAliasExample(ctx context.Context, indexArg int) (types.MapDoubleAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveMapIntegerAliasExample(ctx context.Context, indexArg int) (types.MapIntegerAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveMapRidAliasExample(ctx context.Context, indexArg int) (types.MapRidAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveMapSafeLongAliasExample(ctx context.Context, indexArg int) (types.MapSafeLongAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveMapStringAliasExample(ctx context.Context, indexArg int) (types.MapStringAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveMapUuidAliasExample(ctx context.Context, indexArg int) (types.MapUuidAliasExample, error) {
	panic("implement me")
}

func (AutoDeserializeImpl) ReceiveMapEnumExampleAlias(ctx context.Context, indexArg int) (types.MapEnumExampleAlias, error) {
	panic("implement me")
}
