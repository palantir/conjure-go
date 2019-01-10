// This file was generated by Conjure and should not be manually edited.

package api

import (
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/uuid"
)

type OptionalUuidAlias *uuid.UUID
type UuidAlias uuid.UUID

func (a UuidAlias) MarshalText() ([]byte, error) {
	return uuid.UUID(a).MarshalText()
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
	return safejson.Marshal(Compound(a))
}

func (a *UuidAlias2) UnmarshalJSON(data []byte) error {
	var rawUuidAlias2 Compound
	if err := safejson.Unmarshal(data, &rawUuidAlias2); err != nil {
		return err
	}
	*a = UuidAlias2(rawUuidAlias2)
	return nil
}

func (a UuidAlias2) MarshalYAML() (interface{}, error) {
	return Compound(a), nil
}

func (a *UuidAlias2) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var rawUuidAlias2 Compound
	if err := unmarshal(&rawUuidAlias2); err != nil {
		return err
	}
	*a = UuidAlias2(rawUuidAlias2)
	return nil
}

type BinaryAlias []byte
