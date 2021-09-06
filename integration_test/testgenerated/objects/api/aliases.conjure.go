// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"encoding/base64"

	errors "github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	binary "github.com/palantir/pkg/binary"
	rid "github.com/palantir/pkg/rid"
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
	uuid "github.com/palantir/pkg/uuid"
	werror "github.com/palantir/witchcraft-go-error"
	gjson "github.com/tidwall/gjson"
)

type BinaryAlias []byte

func (a BinaryAlias) String() string {
	return binary.New(a).String()
}

func (a *BinaryAlias) UnmarshalString(data string) error {
	rawBinaryAlias := []byte(data)
	*a = BinaryAlias(rawBinaryAlias)
	return nil
}

func (a BinaryAlias) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a BinaryAlias) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '"')
	if len([]byte(a)) > 0 {
		b64out := make([]byte, base64.StdEncoding.EncodedLen(len([]byte(a))))
		base64.StdEncoding.Encode(b64out, []byte(a))
		out = append(out, b64out...)
	}
	out = append(out, '"')
	if !gjson.ValidBytes(out) {
		return nil, werror.ErrorWithContextParams(context.TODO(), "generated invalid json: please report this as a bug on github.com/palantir/conjure-go/issues")
	}
	return out, nil
}

func (a *BinaryAlias) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for BinaryAlias")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *BinaryAlias) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for BinaryAlias")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *BinaryAlias) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawBinaryAlias []byte
	var err error
	if value.Type != gjson.String {
		err = werror.ErrorWithContextParams(ctx, "BinaryAlias expected JSON string")
		return err
	}
	rawBinaryAlias, err = binary.Binary(value.Str).Bytes()
	if err != nil {
		err = werror.WrapWithContextParams(ctx, err, "BinaryAlias")
		return err
	}
	*a = BinaryAlias(rawBinaryAlias)
	return nil
}

func (a BinaryAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *BinaryAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type NestedAlias1 NestedAlias2

func (a NestedAlias1) String() string {
	return NestedAlias2(a).String()
}

func (a *NestedAlias1) UnmarshalString(data string) error {
	var rawNestedAlias1 NestedAlias2
	if err := rawNestedAlias1.UnmarshalString(data); err != nil {
		return werror.WrapWithContextParams(context.TODO(), errors.WrapWithInvalidArgument(err), "unmarshal string as NestedAlias2(NestedAlias3(optional<string>))")
	}
	*a = NestedAlias1(rawNestedAlias1)
	return nil
}

func (a NestedAlias1) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a NestedAlias1) AppendJSON(out []byte) ([]byte, error) {
	var err error
	out, err = NestedAlias2(a).AppendJSON(out)
	if err != nil {
		return nil, err
	}
	if !gjson.ValidBytes(out) {
		return nil, werror.ErrorWithContextParams(context.TODO(), "generated invalid json: please report this as a bug on github.com/palantir/conjure-go/issues")
	}
	return out, nil
}

func (a *NestedAlias1) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for NestedAlias1")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *NestedAlias1) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for NestedAlias1")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *NestedAlias1) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawNestedAlias1 NestedAlias2
	var err error
	if err = rawNestedAlias1.UnmarshalJSONString(value.Raw); err != nil {
		err = werror.WrapWithContextParams(ctx, err, "NestedAlias1")
		return err
	}
	*a = NestedAlias1(rawNestedAlias1)
	return nil
}

func (a NestedAlias1) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *NestedAlias1) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type NestedAlias2 NestedAlias3

func (a NestedAlias2) String() string {
	return NestedAlias3(a).String()
}

func (a *NestedAlias2) UnmarshalString(data string) error {
	var rawNestedAlias2 NestedAlias3
	if err := rawNestedAlias2.UnmarshalString(data); err != nil {
		return werror.WrapWithContextParams(context.TODO(), errors.WrapWithInvalidArgument(err), "unmarshal string as NestedAlias3(optional<string>)")
	}
	*a = NestedAlias2(rawNestedAlias2)
	return nil
}

func (a NestedAlias2) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a NestedAlias2) AppendJSON(out []byte) ([]byte, error) {
	var err error
	out, err = NestedAlias3(a).AppendJSON(out)
	if err != nil {
		return nil, err
	}
	if !gjson.ValidBytes(out) {
		return nil, werror.ErrorWithContextParams(context.TODO(), "generated invalid json: please report this as a bug on github.com/palantir/conjure-go/issues")
	}
	return out, nil
}

func (a *NestedAlias2) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for NestedAlias2")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *NestedAlias2) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for NestedAlias2")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *NestedAlias2) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawNestedAlias2 NestedAlias3
	var err error
	if err = rawNestedAlias2.UnmarshalJSONString(value.Raw); err != nil {
		err = werror.WrapWithContextParams(ctx, err, "NestedAlias2")
		return err
	}
	*a = NestedAlias2(rawNestedAlias2)
	return nil
}

func (a NestedAlias2) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *NestedAlias2) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type NestedAlias3 struct {
	Value *string
}

func (a NestedAlias3) String() string {
	if a.Value == nil {
		return ""
	}
	return *a.Value
}

func (a *NestedAlias3) UnmarshalString(data string) error {
	rawNestedAlias3 := data
	a.Value = &rawNestedAlias3
	return nil
}

func (a NestedAlias3) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a NestedAlias3) AppendJSON(out []byte) ([]byte, error) {
	if a.Value != nil {
		optVal := *a.Value
		out = safejson.AppendQuotedString(out, optVal)
	} else {
		out = append(out, "null"...)
	}
	if !gjson.ValidBytes(out) {
		return nil, werror.ErrorWithContextParams(context.TODO(), "generated invalid json: please report this as a bug on github.com/palantir/conjure-go/issues")
	}
	return out, nil
}

func (a *NestedAlias3) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for NestedAlias3")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *NestedAlias3) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for NestedAlias3")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *NestedAlias3) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawNestedAlias3 *string
	var err error
	if value.Type != gjson.Null {
		var optVal string
		if value.Type != gjson.String {
			err = werror.ErrorWithContextParams(ctx, "NestedAlias3 expected JSON string")
			return err
		}
		optVal = value.Str
		rawNestedAlias3 = &optVal
	}
	a.Value = rawNestedAlias3
	return nil
}

func (a NestedAlias3) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *NestedAlias3) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type OptionalUuidAlias struct {
	Value *uuid.UUID
}

func (a OptionalUuidAlias) String() string {
	if a.Value == nil {
		return ""
	}
	return a.Value.String()
}

func (a *OptionalUuidAlias) UnmarshalString(data string) error {
	rawOptionalUuidAlias, err := uuid.ParseUUID(data)
	if err != nil {
		return werror.WrapWithContextParams(context.TODO(), errors.WrapWithInvalidArgument(err), "unmarshal string as uuid")
	}
	a.Value = &rawOptionalUuidAlias
	return nil
}

func (a OptionalUuidAlias) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a OptionalUuidAlias) AppendJSON(out []byte) ([]byte, error) {
	if a.Value != nil {
		optVal := *a.Value
		out = safejson.AppendQuotedString(out, optVal.String())
	} else {
		out = append(out, "null"...)
	}
	if !gjson.ValidBytes(out) {
		return nil, werror.ErrorWithContextParams(context.TODO(), "generated invalid json: please report this as a bug on github.com/palantir/conjure-go/issues")
	}
	return out, nil
}

func (a *OptionalUuidAlias) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for OptionalUuidAlias")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *OptionalUuidAlias) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for OptionalUuidAlias")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *OptionalUuidAlias) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawOptionalUuidAlias *uuid.UUID
	var err error
	if value.Type != gjson.Null {
		var optVal uuid.UUID
		if value.Type != gjson.String {
			err = werror.ErrorWithContextParams(ctx, "OptionalUuidAlias expected JSON string")
			return err
		}
		optVal, err = uuid.ParseUUID(value.Str)
		if err != nil {
			err = werror.WrapWithContextParams(ctx, err, "OptionalUuidAlias")
			return err
		}
		rawOptionalUuidAlias = &optVal
	}
	a.Value = rawOptionalUuidAlias
	return nil
}

func (a OptionalUuidAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *OptionalUuidAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type RidAlias rid.ResourceIdentifier

func (a RidAlias) String() string {
	return rid.ResourceIdentifier(a).String()
}

func (a *RidAlias) UnmarshalString(data string) error {
	rawRidAlias, err := rid.ParseRID(data)
	if err != nil {
		return werror.WrapWithContextParams(context.TODO(), errors.WrapWithInvalidArgument(err), "unmarshal string as rid")
	}
	*a = RidAlias(rawRidAlias)
	return nil
}

func (a RidAlias) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a RidAlias) AppendJSON(out []byte) ([]byte, error) {
	out = safejson.AppendQuotedString(out, rid.ResourceIdentifier(a).String())
	if !gjson.ValidBytes(out) {
		return nil, werror.ErrorWithContextParams(context.TODO(), "generated invalid json: please report this as a bug on github.com/palantir/conjure-go/issues")
	}
	return out, nil
}

func (a *RidAlias) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for RidAlias")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *RidAlias) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for RidAlias")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *RidAlias) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawRidAlias rid.ResourceIdentifier
	var err error
	if value.Type != gjson.String {
		err = werror.ErrorWithContextParams(ctx, "RidAlias expected JSON string")
		return err
	}
	rawRidAlias, err = rid.ParseRID(value.Str)
	if err != nil {
		err = werror.WrapWithContextParams(ctx, err, "RidAlias")
		return err
	}
	*a = RidAlias(rawRidAlias)
	return nil
}

func (a RidAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *RidAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type UuidAlias uuid.UUID

func (a UuidAlias) String() string {
	return uuid.UUID(a).String()
}

func (a *UuidAlias) UnmarshalString(data string) error {
	rawUuidAlias, err := uuid.ParseUUID(data)
	if err != nil {
		return werror.WrapWithContextParams(context.TODO(), errors.WrapWithInvalidArgument(err), "unmarshal string as uuid")
	}
	*a = UuidAlias(rawUuidAlias)
	return nil
}

func (a UuidAlias) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a UuidAlias) AppendJSON(out []byte) ([]byte, error) {
	out = safejson.AppendQuotedString(out, uuid.UUID(a).String())
	if !gjson.ValidBytes(out) {
		return nil, werror.ErrorWithContextParams(context.TODO(), "generated invalid json: please report this as a bug on github.com/palantir/conjure-go/issues")
	}
	return out, nil
}

func (a *UuidAlias) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for UuidAlias")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *UuidAlias) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for UuidAlias")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *UuidAlias) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawUuidAlias uuid.UUID
	var err error
	if value.Type != gjson.String {
		err = werror.ErrorWithContextParams(ctx, "UuidAlias expected JSON string")
		return err
	}
	rawUuidAlias, err = uuid.ParseUUID(value.Str)
	if err != nil {
		err = werror.WrapWithContextParams(ctx, err, "UuidAlias")
		return err
	}
	*a = UuidAlias(rawUuidAlias)
	return nil
}

func (a UuidAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *UuidAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type UuidAlias2 Compound

func (a UuidAlias2) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a UuidAlias2) AppendJSON(out []byte) ([]byte, error) {
	var err error
	out, err = Compound(a).AppendJSON(out)
	if err != nil {
		return nil, err
	}
	if !gjson.ValidBytes(out) {
		return nil, werror.ErrorWithContextParams(context.TODO(), "generated invalid json: please report this as a bug on github.com/palantir/conjure-go/issues")
	}
	return out, nil
}

func (a *UuidAlias2) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for UuidAlias2")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (a *UuidAlias2) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for UuidAlias2")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (a *UuidAlias2) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for UuidAlias2")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (a *UuidAlias2) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for UuidAlias2")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (a *UuidAlias2) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	var rawUuidAlias2 Compound
	var err error
	if strict {
		if err = rawUuidAlias2.UnmarshalJSONStringStrict(value.Raw); err != nil {
			err = werror.WrapWithContextParams(ctx, err, "UuidAlias2")
			return err
		}
	} else {
		if err = rawUuidAlias2.UnmarshalJSONString(value.Raw); err != nil {
			err = werror.WrapWithContextParams(ctx, err, "UuidAlias2")
			return err
		}
	}
	*a = UuidAlias2(rawUuidAlias2)
	return nil
}

func (a UuidAlias2) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *UuidAlias2) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}
