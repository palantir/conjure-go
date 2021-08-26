// This file was generated by Conjure and should not be manually edited.

package api

import (
	"encoding/base64"

	binary "github.com/palantir/pkg/binary"
	rid "github.com/palantir/pkg/rid"
	safejson "github.com/palantir/pkg/safejson"
	uuid "github.com/palantir/pkg/uuid"
)

type BinaryAlias []byte

func (a BinaryAlias) String() string {
	return binary.New(a).String()
}

func (a BinaryAlias) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a BinaryAlias) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '"')
	if len([]byte(a)) > 0 {
		b64out := make([]byte, 0, base64.StdEncoding.EncodedLen(len([]byte(a))))
		base64.StdEncoding.Encode(b64out, []byte(a))
		out = append(out, b64out...)
	}
	out = append(out, '"')
	return out, nil
}

func (a *BinaryAlias) UnmarshalText(data []byte) error {
	rawBinaryAlias, err := binary.Binary(data).Bytes()
	if err != nil {
		return err
	}
	*a = BinaryAlias(rawBinaryAlias)
	return nil
}

type NestedAlias1 struct {
	Value NestedAlias2
}

func (a NestedAlias1) String() string {
	if a.Value == nil {
		return ""
	}
	return string(*a.Value)
}

func (a NestedAlias1) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a NestedAlias1) AppendJSON(out []byte) ([]byte, error) {
	if tmpOut, err := a.Value.AppendJSON(out); err != nil {
		return nil, err
	} else {
		out = tmpOut
	}
	return out, nil
}

func (a *NestedAlias1) UnmarshalText(data []byte) error {
	rawNestedAlias1 := string(data)
	a.Value = &rawNestedAlias1
	return nil
}

type NestedAlias2 struct {
	Value NestedAlias3
}

func (a NestedAlias2) String() string {
	if a.Value == nil {
		return ""
	}
	return string(*a.Value)
}

func (a NestedAlias2) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a NestedAlias2) AppendJSON(out []byte) ([]byte, error) {
	if tmpOut, err := a.Value.AppendJSON(out); err != nil {
		return nil, err
	} else {
		out = tmpOut
	}
	return out, nil
}

func (a *NestedAlias2) UnmarshalText(data []byte) error {
	rawNestedAlias2 := string(data)
	a.Value = &rawNestedAlias2
	return nil
}

type NestedAlias3 struct {
	Value *string
}

func (a NestedAlias3) String() string {
	if a.Value == nil {
		return ""
	}
	return string(*a.Value)
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
	return out, nil
}

func (a *NestedAlias3) UnmarshalText(data []byte) error {
	rawNestedAlias3 := string(data)
	a.Value = &rawNestedAlias3
	return nil
}

type OptionalUuidAlias struct {
	Value *uuid.UUID
}

func (a OptionalUuidAlias) String() string {
	if a.Value == nil {
		return ""
	}
	return string(*a.Value)
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
	return out, nil
}

func (a *OptionalUuidAlias) UnmarshalText(data []byte) error {
	if a.Value == nil {
		a.Value = new(uuid.UUID)
	}
	return a.Value.UnmarshalText(data)
}

type RidAlias rid.ResourceIdentifier

func (a RidAlias) String() string {
	return rid.ResourceIdentifier(a).String()
}

func (a RidAlias) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a RidAlias) AppendJSON(out []byte) ([]byte, error) {
	out = safejson.AppendQuotedString(out, rid.ResourceIdentifier(a).String())
	return out, nil
}

func (a *RidAlias) UnmarshalText(data []byte) error {
	var rawRidAlias rid.ResourceIdentifier
	if err := rawRidAlias.UnmarshalText(data); err != nil {
		return err
	}
	*a = RidAlias(rawRidAlias)
	return nil
}

type UuidAlias uuid.UUID

func (a UuidAlias) String() string {
	return uuid.UUID(a).String()
}

func (a UuidAlias) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a UuidAlias) AppendJSON(out []byte) ([]byte, error) {
	out = safejson.AppendQuotedString(out, uuid.UUID(a).String())
	return out, nil
}

func (a *UuidAlias) UnmarshalText(data []byte) error {
	var rawUuidAlias uuid.UUID
	if err := rawUuidAlias.UnmarshalText(data); err != nil {
		return err
	}
	*a = UuidAlias(rawUuidAlias)
	return nil
}

type UuidAlias2 Compound

func (a UuidAlias2) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a UuidAlias2) AppendJSON(out []byte) ([]byte, error) {
	if tmpOut, err := Compound(a).AppendJSON(out); err != nil {
		return nil, err
	} else {
		out = tmpOut
	}
	return out, nil
}

func (a *UuidAlias2) UnmarshalJSON(data []byte) error {
	var rawUuidAlias2 Compound
	if err := safejson.Unmarshal(data, &rawUuidAlias2); err != nil {
		return err
	}
	*a = UuidAlias2(rawUuidAlias2)
	return nil
}
