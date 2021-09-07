// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"encoding/base64"

	errors "github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	binary "github.com/palantir/pkg/binary"
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
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

type BinaryAliasAlias struct {
	Value *BinaryAlias
}

func (a BinaryAliasAlias) String() string {
	if a.Value == nil {
		return ""
	}
	return binary.New(*a.Value).String()
}

func (a *BinaryAliasAlias) UnmarshalString(data string) error {
	var rawBinaryAliasAlias BinaryAlias
	if err := rawBinaryAliasAlias.UnmarshalString(data); err != nil {
		return werror.WrapWithContextParams(context.TODO(), errors.WrapWithInvalidArgument(err), "unmarshal string as BinaryAlias(binary)")
	}
	a.Value = &rawBinaryAliasAlias
	return nil
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

func (a *BinaryAliasAlias) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for BinaryAliasAlias")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *BinaryAliasAlias) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for BinaryAliasAlias")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *BinaryAliasAlias) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawBinaryAliasAlias *BinaryAlias
	var err error
	if value.Type != gjson.Null {
		var optVal BinaryAlias
		if err = optVal.UnmarshalJSONString(value.Raw); err != nil {
			err = werror.WrapWithContextParams(ctx, err, "BinaryAliasAlias")
			return err
		}
		rawBinaryAliasAlias = &optVal
	}
	a.Value = rawBinaryAliasAlias
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

func (a BinaryAliasOptional) String() string {
	if a.Value == nil {
		return ""
	}
	return binary.New(*a.Value).String()
}

func (a *BinaryAliasOptional) UnmarshalString(data string) error {
	rawBinaryAliasOptional := []byte(data)
	a.Value = &rawBinaryAliasOptional
	return nil
}

func (a BinaryAliasOptional) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a BinaryAliasOptional) AppendJSON(out []byte) ([]byte, error) {
	if a.Value != nil {
		optVal := *a.Value
		out = append(out, '"')
		if len(optVal) > 0 {
			b64out := make([]byte, base64.StdEncoding.EncodedLen(len(optVal)))
			base64.StdEncoding.Encode(b64out, optVal)
			out = append(out, b64out...)
		}
		out = append(out, '"')
	} else {
		out = append(out, "null"...)
	}
	return out, nil
}

func (a *BinaryAliasOptional) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for BinaryAliasOptional")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *BinaryAliasOptional) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for BinaryAliasOptional")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *BinaryAliasOptional) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawBinaryAliasOptional *[]byte
	var err error
	if value.Type != gjson.Null {
		var optVal []byte
		if value.Type != gjson.String {
			err = werror.ErrorWithContextParams(ctx, "BinaryAliasOptional expected JSON string")
			return err
		}
		optVal, err = binary.Binary(value.Str).Bytes()
		if err != nil {
			err = werror.WrapWithContextParams(ctx, err, "BinaryAliasOptional")
			return err
		}
		rawBinaryAliasOptional = &optVal
	}
	a.Value = rawBinaryAliasOptional
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
