// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"

	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
	werror "github.com/palantir/witchcraft-go-error"
	gjson "github.com/tidwall/gjson"
)

type Basic struct {
	Data string `json:"data"`
}

func (o Basic) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o Basic) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"data\":"...)
		out = safejson.AppendQuotedString(out, o.Data)
	}
	out = append(out, '}')
	if !gjson.ValidBytes(out) {
		return nil, werror.ErrorWithContextParams(context.TODO(), "generated invalid json: please report this as a bug on github.com/palantir/conjure-go/issues")
	}
	return out, nil
}

func (o *Basic) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Basic")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (o *Basic) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Basic")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (o *Basic) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Basic")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (o *Basic) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Basic")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (o *Basic) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type Basic expected JSON object")
	}
	var seenData bool
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "data":
			if seenData {
				err = werror.ErrorWithContextParams(ctx, "type Basic encountered duplicate \"data\" field")
				return false
			} else {
				seenData = true
			}
			if value.Type != gjson.String {
				err = werror.ErrorWithContextParams(ctx, "field Basic[\"data\"] expected JSON string")
				return false
			}
			o.Data = value.Str
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
		return werror.ErrorWithContextParams(ctx, "type Basic missing required JSON fields", werror.SafeParam("missingFields", missingFields))
	}
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type Basic encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (o Basic) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *Basic) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
