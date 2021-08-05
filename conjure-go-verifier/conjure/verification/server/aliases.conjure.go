// This file was generated by Conjure and should not be manually edited.

package server

import (
	"context"
	"encoding/json"

	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/tidwall/gjson"
)

type EndpointName string

func (a EndpointName) String() string {
	return string(a)
}

func (a EndpointName) MarshalJSON() ([]byte, error) {
	return a.MarshalJSONBuffer(nil)
}

func (a EndpointName) MarshalJSONBuffer(buf []byte) ([]byte, error) {
	buf = safejson.AppendQuotedString(buf, string(a))
	return nil, nil
}

func (a *EndpointName) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.ParseBytes(data), false)
}

func (a *EndpointName) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.Parse(data), false)
}

func (a *EndpointName) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.ParseBytes(data), true)
}

func (a *EndpointName) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid json")
	}
	return a.unmarshalGJSON(ctx, gjson.Parse(data), true)
}

func (a *EndpointName) unmarshalGJSON(ctx context.Context, value gjson.Result, strict bool) error {
	var err error
	var objectValue string
	if value.Type != gjson.String {
		err = werror.ErrorWithContextParams(ctx, "type EndpointName expected json type String")
		return err
	}
	objectValue = value.Str
	*a = EndpointName(objectValue)
	return err
}

func (a *EndpointName) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

func (a EndpointName) MarshalYAML() (interface{}, error) {
	jsonBytes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
