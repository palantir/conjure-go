// This file was generated by Conjure and should not be manually edited.

package api

import (
	"encoding/base64"

	binary "github.com/palantir/pkg/binary"
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

type BinaryAliasAlias struct {
	Value *BinaryAlias
}

func (a BinaryAliasAlias) String() string {
	if a.Value == nil {
		return ""
	}
	return binary.New(*a.Value).String()
}

func (a BinaryAliasAlias) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a BinaryAliasAlias) AppendJSON(out []byte) ([]byte, error) {
	if a.Value != nil {
		optVal := *a.Value
		var err error
		out, err = optVal.AppendJSON(out)
		if err != nil {
			return nil, err
		}
	} else {
		out = append(out, "null"...)
	}
	return out, nil
}

func (a *BinaryAliasAlias) UnmarshalText(data []byte) error {
	rawBinaryAliasAlias, err := binary.Binary(data).Bytes()
	if err != nil {
		return err
	}
	*a.Value = rawBinaryAliasAlias
	return nil
}

type BinaryAliasOptional struct {
	Value *[]byte
}

func (a BinaryAliasOptional) String() string {
	if a.Value == nil {
		return ""
	}
	return binary.New(*a.Value).String()
}

func (a BinaryAliasOptional) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a BinaryAliasOptional) AppendJSON(out []byte) ([]byte, error) {
	if a.Value != nil {
		optVal := *a.Value
		out = append(out, '"')
		if len(optVal) > 0 {
			b64out := make([]byte, 0, base64.StdEncoding.EncodedLen(len(optVal)))
			base64.StdEncoding.Encode(b64out, optVal)
			out = append(out, b64out...)
		}
		out = append(out, '"')
	} else {
		out = append(out, "null"...)
	}
	return out, nil
}

func (a *BinaryAliasOptional) UnmarshalText(data []byte) error {
	rawBinaryAliasOptional, err := binary.Binary(data).Bytes()
	if err != nil {
		return err
	}
	*a.Value = rawBinaryAliasOptional
	return nil
}
