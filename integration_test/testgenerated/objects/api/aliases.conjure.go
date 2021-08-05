// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/palantir/pkg/binary"
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
	"github.com/palantir/pkg/uuid"
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/tidwall/gjson"
)

type BinaryAlias []byte

func (a BinaryAlias) MarshalJSON() ([]byte, error) {
	return a.MarshalJSONBuffer(nil)
}

func (a BinaryAlias) MarshalJSONBuffer(buf []byte) ([]byte, error) {
	buf = append(buf, '"')
	if len([]byte(a)) > 0 {
		b64out := make([]byte, 0, base64.StdEncoding.EncodedLen(len([]byte(a))))
		base64.StdEncoding.Encode(b64out, []byte(a))
		buf = append(buf, b64out...)
	}
	buf = append(buf, '"')
	return nil, nil
}

func (a *BinaryAlias) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.ParseBytes(data), false)
}

func (a *BinaryAlias) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.Parse(data), false)
}

func (a *BinaryAlias) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.ParseBytes(data), true)
}

func (a *BinaryAlias) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.Parse(data), true)
}

func (a *BinaryAlias) unmarshalGJSON(ctx context.Context, value gjson.Result, strict bool) error {
	var err error
	var objectValue []byte
	objectValue = make([]byte, 0)
	if value.Type != gjson.String {
		err = werror.ErrorWithContextParams(ctx, "type BinaryAlias expected json type String")
		return err
	}
	objectValue, err = binary.Binary(value.Str).Bytes()
	*a = BinaryAlias(objectValue)
	return err
}

func (a *BinaryAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

func (a BinaryAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

type CompoundAlias Compound

func (a CompoundAlias) MarshalJSON() ([]byte, error) {
	return a.MarshalJSONBuffer(nil)
}

func (a CompoundAlias) MarshalJSONBuffer(buf []byte) ([]byte, error) {
	if out, err := Compound(a).MarshalJSONBuffer(buf); err != nil {
		return nil, err
	} else {
		buf = out
	}
	return nil, nil
}

func (a *CompoundAlias) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.ParseBytes(data), false)
}

func (a *CompoundAlias) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.Parse(data), false)
}

func (a *CompoundAlias) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.ParseBytes(data), true)
}

func (a *CompoundAlias) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.Parse(data), true)
}

func (a *CompoundAlias) unmarshalGJSON(ctx context.Context, value gjson.Result, strict bool) error {
	var err error
	var objectValue Compound
	if strict {
		err = objectValue.UnmarshalJSONStringStrict(value.Raw)
	} else {
		err = objectValue.UnmarshalJSONString(value.Raw)
	}
	err = werror.WrapWithContextParams(ctx, err, "type CompoundAlias")
	*a = CompoundAlias(objectValue)
	return err
}

func (a *CompoundAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

func (a CompoundAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

type OptionalCompound struct {
	Value *Compound
}

func (a OptionalCompound) MarshalJSON() ([]byte, error) {
	return a.MarshalJSONBuffer(nil)
}

func (a OptionalCompound) MarshalJSONBuffer(buf []byte) ([]byte, error) {
	if a.Value != nil {
		return nil, err
	} else {
		buf = out
	}
	if a.Value != nil {
		if out, err := (*a.Value).MarshalJSONBuffer(buf); err != nil {
			return nil, err
		} else {
			buf = out
		}
	} else {
		buf = append(buf, "null"...)
	}
	return nil, nil
}

func (a *OptionalCompound) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.ParseBytes(data), false)
}

func (a *OptionalCompound) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.Parse(data), false)
}

func (a *OptionalCompound) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.ParseBytes(data), true)
}

func (a *OptionalCompound) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.Parse(data), true)
}

func (a *OptionalCompound) unmarshalGJSON(ctx context.Context, value gjson.Result, strict bool) error {
	var err error
	if value.Type != gjson.Null {
		var optionalValue Compound
		if strict {
			err = optionalValue.UnmarshalJSONStringStrict(value.Raw)
		} else {
			err = optionalValue.UnmarshalJSONString(value.Raw)
		}
		err = werror.WrapWithContextParams(ctx, err, "type OptionalCompound")
		a.Value = &optionalValue
	}
	return err
}

func (a *OptionalCompound) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

func (a OptionalCompound) MarshalYAML() (interface{}, error) {
	jsonBytes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
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

func (a OptionalUuidAlias) MarshalJSON() ([]byte, error) {
	return a.MarshalJSONBuffer(nil)
}

func (a OptionalUuidAlias) MarshalJSONBuffer(buf []byte) ([]byte, error) {
	if a.Value != nil {
		return nil, err
	} else {
		buf = out
	}
	if a.Value != nil {
		buf = safejson.AppendQuotedString(buf, (*a.Value).String())
	} else {
		buf = append(buf, "null"...)
	}
	return nil, nil
}

func (a *OptionalUuidAlias) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.ParseBytes(data), false)
}

func (a *OptionalUuidAlias) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.Parse(data), false)
}

func (a *OptionalUuidAlias) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.ParseBytes(data), true)
}

func (a *OptionalUuidAlias) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.Parse(data), true)
}

func (a *OptionalUuidAlias) unmarshalGJSON(ctx context.Context, value gjson.Result, strict bool) error {
	var err error
	if value.Type != gjson.Null {
		if value.Type != gjson.String {
			err = werror.ErrorWithContextParams(ctx, "type OptionalUuidAlias expected json type String")
			return err
		}
		var optionalValue uuid.UUID
		optionalValue, err = uuid.ParseUUID(value.Str)
		a.Value = &optionalValue
	}
	return err
}

func (a *OptionalUuidAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

func (a OptionalUuidAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

type RidAlias rid.ResourceIdentifier

func (a RidAlias) String() string {
	return rid.ResourceIdentifier(a).String()
}

func (a RidAlias) MarshalJSON() ([]byte, error) {
	return a.MarshalJSONBuffer(nil)
}

func (a RidAlias) MarshalJSONBuffer(buf []byte) ([]byte, error) {
	buf = safejson.AppendQuotedString(buf, rid.ResourceIdentifier(a).String())
	return nil, nil
}

func (a *RidAlias) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.ParseBytes(data), false)
}

func (a *RidAlias) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.Parse(data), false)
}

func (a *RidAlias) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.ParseBytes(data), true)
}

func (a *RidAlias) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.Parse(data), true)
}

func (a *RidAlias) unmarshalGJSON(ctx context.Context, value gjson.Result, strict bool) error {
	var err error
	var objectValue rid.ResourceIdentifier
	if value.Type != gjson.String {
		err = werror.ErrorWithContextParams(ctx, "type RidAlias expected json type String")
		return err
	}
	objectValue, err = rid.ParseRID(value.Str)
	*a = RidAlias(objectValue)
	return err
}

func (a *RidAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

func (a RidAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

type UuidAlias uuid.UUID

func (a UuidAlias) String() string {
	return uuid.UUID(a).String()
}

func (a UuidAlias) MarshalJSON() ([]byte, error) {
	return a.MarshalJSONBuffer(nil)
}

func (a UuidAlias) MarshalJSONBuffer(buf []byte) ([]byte, error) {
	buf = safejson.AppendQuotedString(buf, uuid.UUID(a).String())
	return nil, nil
}

func (a *UuidAlias) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.ParseBytes(data), false)
}

func (a *UuidAlias) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.Parse(data), false)
}

func (a *UuidAlias) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.ParseBytes(data), true)
}

func (a *UuidAlias) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.Parse(data), true)
}

func (a *UuidAlias) unmarshalGJSON(ctx context.Context, value gjson.Result, strict bool) error {
	var err error
	var objectValue uuid.UUID
	if value.Type != gjson.String {
		err = werror.ErrorWithContextParams(ctx, "type UuidAlias expected json type String")
		return err
	}
	objectValue, err = uuid.ParseUUID(value.Str)
	*a = UuidAlias(objectValue)
	return err
}

func (a *UuidAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

func (a UuidAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
