// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"encoding/base64"

	binary "github.com/palantir/pkg/binary"
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
	werror "github.com/palantir/witchcraft-go-error"
	gjson "github.com/tidwall/gjson"
)

type CustomObject struct {
	Data []byte `json:"data"`
}

func (o CustomObject) MarshalJSON() ([]byte, error) {
	size, err := o.JSONSize()
	if err != nil {
		return nil, err
	}
	return o.AppendJSON(make([]byte, 0, size))
}

func (o CustomObject) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"data\":"...)
		out = append(out, '"')
		if len(o.Data) > 0 {
			b64out := make([]byte, base64.StdEncoding.EncodedLen(len(o.Data)))
			base64.StdEncoding.Encode(b64out, o.Data)
			out = append(out, b64out...)
		}
		out = append(out, '"')
	}
	out = append(out, '}')
	return out, nil
}

func (o CustomObject) JSONSize() (int, error) {
	var out int
	out += 1 // '{'
	{
		out += 7 // "data":
		out += 1 // '"'
		if len(o.Data) > 0 {
			b64out := make([]byte, base64.StdEncoding.EncodedLen(len(o.Data)))
			base64.StdEncoding.Encode(b64out, o.Data)
			out += len(b64out)
		}
		out += 1 // '"'
	}
	out += 1 // '}'
	return out, nil
}

func (o *CustomObject) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for CustomObject")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (o *CustomObject) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for CustomObject")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (o *CustomObject) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for CustomObject")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (o *CustomObject) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for CustomObject")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (o *CustomObject) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type CustomObject expected JSON object")
	}
	var seenData bool
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "data":
			if seenData {
				err = werror.ErrorWithContextParams(ctx, "type CustomObject encountered duplicate \"data\" field")
				return false
			}
			seenData = true
			if value.Type != gjson.String {
				err = werror.ErrorWithContextParams(ctx, "field CustomObject[\"data\"] expected JSON string")
				return false
			}
			o.Data, err = binary.Binary(value.Str).Bytes()
			if err != nil {
				err = werror.WrapWithContextParams(ctx, err, "field CustomObject[\"data\"]")
				return false
			}
		default:
			if strict {
				unrecognizedFields = append(unrecognizedFields, key.Str)
			}
		}
		return err == nil
	})
	if err != nil {
		return err
	}
	var missingFields []string
	if !seenData {
		missingFields = append(missingFields, "data")
	}
	if len(missingFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type CustomObject missing required JSON fields", werror.SafeParam("missingFields", missingFields))
	}
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type CustomObject encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (o CustomObject) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *CustomObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
