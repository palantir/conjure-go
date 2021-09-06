// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"

	errors "github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	rid "github.com/palantir/pkg/rid"
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
	werror "github.com/palantir/witchcraft-go-error"
	gjson "github.com/tidwall/gjson"
)

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

type StringAlias string

func (a StringAlias) String() string {
	return string(a)
}

func (a *StringAlias) UnmarshalString(data string) error {
	rawStringAlias := data
	*a = StringAlias(rawStringAlias)
	return nil
}

func (a StringAlias) MarshalJSON() ([]byte, error) {
	return a.AppendJSON(nil)
}

func (a StringAlias) AppendJSON(out []byte) ([]byte, error) {
	out = safejson.AppendQuotedString(out, string(a))
	if !gjson.ValidBytes(out) {
		return nil, werror.ErrorWithContextParams(context.TODO(), "generated invalid json: please report this as a bug on github.com/palantir/conjure-go/issues")
	}
	return out, nil
}

func (a *StringAlias) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for StringAlias")
	}
	return a.unmarshalJSONResult(ctx, gjson.ParseBytes(data))
}

func (a *StringAlias) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for StringAlias")
	}
	return a.unmarshalJSONResult(ctx, gjson.Parse(data))
}

func (a *StringAlias) unmarshalJSONResult(ctx context.Context, value gjson.Result) error {
	var rawStringAlias string
	var err error
	if value.Type != gjson.String {
		err = werror.ErrorWithContextParams(ctx, "StringAlias expected JSON string")
		return err
	}
	rawStringAlias = value.Str
	*a = StringAlias(rawStringAlias)
	return nil
}

func (a StringAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *StringAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}
