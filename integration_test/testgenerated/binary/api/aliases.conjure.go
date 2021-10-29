// This file was generated by Conjure and should not be manually edited.

package api

import (
	binary "github.com/palantir/pkg/binary"
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
)

type BinaryAlias []byte

func (a BinaryAlias) String() string {
	return binary.New(a).String()
}

func (a BinaryAlias) MarshalText() ([]byte, error) {
	return binary.New(a).MarshalText()
}

func (a *BinaryAlias) UnmarshalText(data []byte) error {
	rawBinaryAlias, err := binary.Binary(data).Bytes()
	if err != nil {
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

type BinaryAliasAlias struct {
	Value *BinaryAlias
}

func (a BinaryAliasAlias) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return binary.New(*a.Value).MarshalText()
}

func (a *BinaryAliasAlias) UnmarshalText(data []byte) error {
	rawBinaryAliasAlias, err := binary.Binary(data).Bytes()
	if err != nil {
		return err
	}
	*a.Value = rawBinaryAliasAlias
	return nil
}

func (a BinaryAliasAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *BinaryAliasAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type BinaryAliasOptional struct {
	Value *[]byte
}

func (a BinaryAliasOptional) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return binary.New(*a.Value).MarshalText()
}

func (a *BinaryAliasOptional) UnmarshalText(data []byte) error {
	rawBinaryAliasOptional, err := binary.Binary(data).Bytes()
	if err != nil {
		return err
	}
	*a.Value = rawBinaryAliasOptional
	return nil
}

func (a BinaryAliasOptional) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *BinaryAliasOptional) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}
