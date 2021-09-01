// This file was generated by Conjure and should not be manually edited.

package server

import (
	"context"

	safejson "github.com/palantir/pkg/safejson"
	werror "github.com/palantir/witchcraft-go-error"
	gjson "github.com/tidwall/gjson"
)

type ClientTestCases struct {
	AutoDeserialize         map[EndpointName]PositiveAndNegativeTestCases `json:"autoDeserialize"`
	SingleHeaderService     map[EndpointName][]string                     `json:"singleHeaderService"`
	SinglePathParamService  map[EndpointName][]string                     `json:"singlePathParamService"`
	SingleQueryParamService map[EndpointName][]string                     `json:"singleQueryParamService"`
}

func (o ClientTestCases) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o ClientTestCases) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"autoDeserialize\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.AutoDeserialize {
				{
					var err error
					out, err = k.AppendJSON(out)
					if err != nil {
						return nil, err
					}
				}
				out = append(out, ':')
				{
					var err error
					out, err = v.AppendJSON(out)
					if err != nil {
						return nil, err
					}
				}
				i++
				if i < len(o.AutoDeserialize) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
		out = append(out, ',')
	}
	{
		out = append(out, "\"singleHeaderService\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.SingleHeaderService {
				{
					var err error
					out, err = k.AppendJSON(out)
					if err != nil {
						return nil, err
					}
				}
				out = append(out, ':')
				{
					out = append(out, '[')
					for i1 := range v {
						out = safejson.AppendQuotedString(out, v[i1])
						if i1 < len(v)-1 {
							out = append(out, ',')
						}
					}
					out = append(out, ']')
				}
				i++
				if i < len(o.SingleHeaderService) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
		out = append(out, ',')
	}
	{
		out = append(out, "\"singlePathParamService\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.SinglePathParamService {
				{
					var err error
					out, err = k.AppendJSON(out)
					if err != nil {
						return nil, err
					}
				}
				out = append(out, ':')
				{
					out = append(out, '[')
					for i1 := range v {
						out = safejson.AppendQuotedString(out, v[i1])
						if i1 < len(v)-1 {
							out = append(out, ',')
						}
					}
					out = append(out, ']')
				}
				i++
				if i < len(o.SinglePathParamService) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
		out = append(out, ',')
	}
	{
		out = append(out, "\"singleQueryParamService\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.SingleQueryParamService {
				{
					var err error
					out, err = k.AppendJSON(out)
					if err != nil {
						return nil, err
					}
				}
				out = append(out, ':')
				{
					out = append(out, '[')
					for i1 := range v {
						out = safejson.AppendQuotedString(out, v[i1])
						if i1 < len(v)-1 {
							out = append(out, ',')
						}
					}
					out = append(out, ']')
				}
				i++
				if i < len(o.SingleQueryParamService) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
	}
	out = append(out, '}')
	return out, nil
}

func (o *ClientTestCases) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ClientTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (o *ClientTestCases) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ClientTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (o *ClientTestCases) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ClientTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (o *ClientTestCases) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ClientTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (o *ClientTestCases) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type ClientTestCases expected JSON object")
	}
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "autoDeserialize":
			if !value.IsObject() {
				err = werror.ErrorWithContextParams(ctx, "field ClientTestCases[\"autoDeserialize\"] expected JSON object")
				return false
			}
			if o.AutoDeserialize == nil {
				o.AutoDeserialize = make(map[EndpointName]PositiveAndNegativeTestCases, 0)
			}
			value.ForEach(func(key, value gjson.Result) bool {
				var mapKey EndpointName
				var mapVal PositiveAndNegativeTestCases
				{
					if err = mapKey.UnmarshalJSONString(key.Raw); err != nil {
						err = werror.WrapWithContextParams(ctx, err, "field ClientTestCases[\"autoDeserialize\"] map key")
						return false
					}
				}
				{
					if strict {
						if err = mapVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
							err = werror.WrapWithContextParams(ctx, err, "field ClientTestCases[\"autoDeserialize\"] map value")
							return false
						}
					} else {
						if err = mapVal.UnmarshalJSONString(value.Raw); err != nil {
							err = werror.WrapWithContextParams(ctx, err, "field ClientTestCases[\"autoDeserialize\"] map value")
							return false
						}
					}
				}
				o.AutoDeserialize[mapKey] = mapVal
				return err == nil
			})
		case "singleHeaderService":
			if !value.IsObject() {
				err = werror.ErrorWithContextParams(ctx, "field ClientTestCases[\"singleHeaderService\"] expected JSON object")
				return false
			}
			if o.SingleHeaderService == nil {
				o.SingleHeaderService = make(map[EndpointName][]string, 0)
			}
			value.ForEach(func(key, value gjson.Result) bool {
				var mapKey EndpointName
				var mapVal []string
				{
					if err = mapKey.UnmarshalJSONString(key.Raw); err != nil {
						err = werror.WrapWithContextParams(ctx, err, "field ClientTestCases[\"singleHeaderService\"] map key")
						return false
					}
				}
				{
					if !value.IsArray() {
						err = werror.ErrorWithContextParams(ctx, "field ClientTestCases[\"singleHeaderService\"] map value expected JSON array")
						return false
					}
					value.ForEach(func(_, value gjson.Result) bool {
						var listElement1 string
						if value.Type != gjson.String {
							err = werror.ErrorWithContextParams(ctx, "field ClientTestCases[\"singleHeaderService\"] map value list element expected JSON string")
							return false
						}
						listElement1 = value.Str
						mapVal = append(mapVal, listElement1)
						return err == nil
					})
				}
				o.SingleHeaderService[mapKey] = mapVal
				return err == nil
			})
		case "singlePathParamService":
			if !value.IsObject() {
				err = werror.ErrorWithContextParams(ctx, "field ClientTestCases[\"singlePathParamService\"] expected JSON object")
				return false
			}
			if o.SinglePathParamService == nil {
				o.SinglePathParamService = make(map[EndpointName][]string, 0)
			}
			value.ForEach(func(key, value gjson.Result) bool {
				var mapKey EndpointName
				var mapVal []string
				{
					if err = mapKey.UnmarshalJSONString(key.Raw); err != nil {
						err = werror.WrapWithContextParams(ctx, err, "field ClientTestCases[\"singlePathParamService\"] map key")
						return false
					}
				}
				{
					if !value.IsArray() {
						err = werror.ErrorWithContextParams(ctx, "field ClientTestCases[\"singlePathParamService\"] map value expected JSON array")
						return false
					}
					value.ForEach(func(_, value gjson.Result) bool {
						var listElement1 string
						if value.Type != gjson.String {
							err = werror.ErrorWithContextParams(ctx, "field ClientTestCases[\"singlePathParamService\"] map value list element expected JSON string")
							return false
						}
						listElement1 = value.Str
						mapVal = append(mapVal, listElement1)
						return err == nil
					})
				}
				o.SinglePathParamService[mapKey] = mapVal
				return err == nil
			})
		case "singleQueryParamService":
			if !value.IsObject() {
				err = werror.ErrorWithContextParams(ctx, "field ClientTestCases[\"singleQueryParamService\"] expected JSON object")
				return false
			}
			if o.SingleQueryParamService == nil {
				o.SingleQueryParamService = make(map[EndpointName][]string, 0)
			}
			value.ForEach(func(key, value gjson.Result) bool {
				var mapKey EndpointName
				var mapVal []string
				{
					if err = mapKey.UnmarshalJSONString(key.Raw); err != nil {
						err = werror.WrapWithContextParams(ctx, err, "field ClientTestCases[\"singleQueryParamService\"] map key")
						return false
					}
				}
				{
					if !value.IsArray() {
						err = werror.ErrorWithContextParams(ctx, "field ClientTestCases[\"singleQueryParamService\"] map value expected JSON array")
						return false
					}
					value.ForEach(func(_, value gjson.Result) bool {
						var listElement1 string
						if value.Type != gjson.String {
							err = werror.ErrorWithContextParams(ctx, "field ClientTestCases[\"singleQueryParamService\"] map value list element expected JSON string")
							return false
						}
						listElement1 = value.Str
						mapVal = append(mapVal, listElement1)
						return err == nil
					})
				}
				o.SingleQueryParamService[mapKey] = mapVal
				return err == nil
			})
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
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type ClientTestCases encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

type IgnoredClientTestCases struct {
	AutoDeserialize         map[EndpointName][]string `json:"autoDeserialize"`
	SingleHeaderService     map[EndpointName][]string `json:"singleHeaderService"`
	SinglePathParamService  map[EndpointName][]string `json:"singlePathParamService"`
	SingleQueryParamService map[EndpointName][]string `json:"singleQueryParamService"`
}

func (o IgnoredClientTestCases) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o IgnoredClientTestCases) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"autoDeserialize\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.AutoDeserialize {
				{
					var err error
					out, err = k.AppendJSON(out)
					if err != nil {
						return nil, err
					}
				}
				out = append(out, ':')
				{
					out = append(out, '[')
					for i1 := range v {
						out = safejson.AppendQuotedString(out, v[i1])
						if i1 < len(v)-1 {
							out = append(out, ',')
						}
					}
					out = append(out, ']')
				}
				i++
				if i < len(o.AutoDeserialize) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
		out = append(out, ',')
	}
	{
		out = append(out, "\"singleHeaderService\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.SingleHeaderService {
				{
					var err error
					out, err = k.AppendJSON(out)
					if err != nil {
						return nil, err
					}
				}
				out = append(out, ':')
				{
					out = append(out, '[')
					for i1 := range v {
						out = safejson.AppendQuotedString(out, v[i1])
						if i1 < len(v)-1 {
							out = append(out, ',')
						}
					}
					out = append(out, ']')
				}
				i++
				if i < len(o.SingleHeaderService) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
		out = append(out, ',')
	}
	{
		out = append(out, "\"singlePathParamService\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.SinglePathParamService {
				{
					var err error
					out, err = k.AppendJSON(out)
					if err != nil {
						return nil, err
					}
				}
				out = append(out, ':')
				{
					out = append(out, '[')
					for i1 := range v {
						out = safejson.AppendQuotedString(out, v[i1])
						if i1 < len(v)-1 {
							out = append(out, ',')
						}
					}
					out = append(out, ']')
				}
				i++
				if i < len(o.SinglePathParamService) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
		out = append(out, ',')
	}
	{
		out = append(out, "\"singleQueryParamService\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.SingleQueryParamService {
				{
					var err error
					out, err = k.AppendJSON(out)
					if err != nil {
						return nil, err
					}
				}
				out = append(out, ':')
				{
					out = append(out, '[')
					for i1 := range v {
						out = safejson.AppendQuotedString(out, v[i1])
						if i1 < len(v)-1 {
							out = append(out, ',')
						}
					}
					out = append(out, ']')
				}
				i++
				if i < len(o.SingleQueryParamService) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
	}
	out = append(out, '}')
	return out, nil
}

func (o *IgnoredClientTestCases) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for IgnoredClientTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (o *IgnoredClientTestCases) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for IgnoredClientTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (o *IgnoredClientTestCases) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for IgnoredClientTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (o *IgnoredClientTestCases) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for IgnoredClientTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (o *IgnoredClientTestCases) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type IgnoredClientTestCases expected JSON object")
	}
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "autoDeserialize":
			if !value.IsObject() {
				err = werror.ErrorWithContextParams(ctx, "field IgnoredClientTestCases[\"autoDeserialize\"] expected JSON object")
				return false
			}
			if o.AutoDeserialize == nil {
				o.AutoDeserialize = make(map[EndpointName][]string, 0)
			}
			value.ForEach(func(key, value gjson.Result) bool {
				var mapKey EndpointName
				var mapVal []string
				{
					if err = mapKey.UnmarshalJSONString(key.Raw); err != nil {
						err = werror.WrapWithContextParams(ctx, err, "field IgnoredClientTestCases[\"autoDeserialize\"] map key")
						return false
					}
				}
				{
					if !value.IsArray() {
						err = werror.ErrorWithContextParams(ctx, "field IgnoredClientTestCases[\"autoDeserialize\"] map value expected JSON array")
						return false
					}
					value.ForEach(func(_, value gjson.Result) bool {
						var listElement1 string
						if value.Type != gjson.String {
							err = werror.ErrorWithContextParams(ctx, "field IgnoredClientTestCases[\"autoDeserialize\"] map value list element expected JSON string")
							return false
						}
						listElement1 = value.Str
						mapVal = append(mapVal, listElement1)
						return err == nil
					})
				}
				o.AutoDeserialize[mapKey] = mapVal
				return err == nil
			})
		case "singleHeaderService":
			if !value.IsObject() {
				err = werror.ErrorWithContextParams(ctx, "field IgnoredClientTestCases[\"singleHeaderService\"] expected JSON object")
				return false
			}
			if o.SingleHeaderService == nil {
				o.SingleHeaderService = make(map[EndpointName][]string, 0)
			}
			value.ForEach(func(key, value gjson.Result) bool {
				var mapKey EndpointName
				var mapVal []string
				{
					if err = mapKey.UnmarshalJSONString(key.Raw); err != nil {
						err = werror.WrapWithContextParams(ctx, err, "field IgnoredClientTestCases[\"singleHeaderService\"] map key")
						return false
					}
				}
				{
					if !value.IsArray() {
						err = werror.ErrorWithContextParams(ctx, "field IgnoredClientTestCases[\"singleHeaderService\"] map value expected JSON array")
						return false
					}
					value.ForEach(func(_, value gjson.Result) bool {
						var listElement1 string
						if value.Type != gjson.String {
							err = werror.ErrorWithContextParams(ctx, "field IgnoredClientTestCases[\"singleHeaderService\"] map value list element expected JSON string")
							return false
						}
						listElement1 = value.Str
						mapVal = append(mapVal, listElement1)
						return err == nil
					})
				}
				o.SingleHeaderService[mapKey] = mapVal
				return err == nil
			})
		case "singlePathParamService":
			if !value.IsObject() {
				err = werror.ErrorWithContextParams(ctx, "field IgnoredClientTestCases[\"singlePathParamService\"] expected JSON object")
				return false
			}
			if o.SinglePathParamService == nil {
				o.SinglePathParamService = make(map[EndpointName][]string, 0)
			}
			value.ForEach(func(key, value gjson.Result) bool {
				var mapKey EndpointName
				var mapVal []string
				{
					if err = mapKey.UnmarshalJSONString(key.Raw); err != nil {
						err = werror.WrapWithContextParams(ctx, err, "field IgnoredClientTestCases[\"singlePathParamService\"] map key")
						return false
					}
				}
				{
					if !value.IsArray() {
						err = werror.ErrorWithContextParams(ctx, "field IgnoredClientTestCases[\"singlePathParamService\"] map value expected JSON array")
						return false
					}
					value.ForEach(func(_, value gjson.Result) bool {
						var listElement1 string
						if value.Type != gjson.String {
							err = werror.ErrorWithContextParams(ctx, "field IgnoredClientTestCases[\"singlePathParamService\"] map value list element expected JSON string")
							return false
						}
						listElement1 = value.Str
						mapVal = append(mapVal, listElement1)
						return err == nil
					})
				}
				o.SinglePathParamService[mapKey] = mapVal
				return err == nil
			})
		case "singleQueryParamService":
			if !value.IsObject() {
				err = werror.ErrorWithContextParams(ctx, "field IgnoredClientTestCases[\"singleQueryParamService\"] expected JSON object")
				return false
			}
			if o.SingleQueryParamService == nil {
				o.SingleQueryParamService = make(map[EndpointName][]string, 0)
			}
			value.ForEach(func(key, value gjson.Result) bool {
				var mapKey EndpointName
				var mapVal []string
				{
					if err = mapKey.UnmarshalJSONString(key.Raw); err != nil {
						err = werror.WrapWithContextParams(ctx, err, "field IgnoredClientTestCases[\"singleQueryParamService\"] map key")
						return false
					}
				}
				{
					if !value.IsArray() {
						err = werror.ErrorWithContextParams(ctx, "field IgnoredClientTestCases[\"singleQueryParamService\"] map value expected JSON array")
						return false
					}
					value.ForEach(func(_, value gjson.Result) bool {
						var listElement1 string
						if value.Type != gjson.String {
							err = werror.ErrorWithContextParams(ctx, "field IgnoredClientTestCases[\"singleQueryParamService\"] map value list element expected JSON string")
							return false
						}
						listElement1 = value.Str
						mapVal = append(mapVal, listElement1)
						return err == nil
					})
				}
				o.SingleQueryParamService[mapKey] = mapVal
				return err == nil
			})
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
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type IgnoredClientTestCases encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

type IgnoredTestCases struct {
	Client IgnoredClientTestCases `json:"client"`
}

func (o IgnoredTestCases) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o IgnoredTestCases) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"client\":"...)
		var err error
		out, err = o.Client.AppendJSON(out)
		if err != nil {
			return nil, err
		}
	}
	out = append(out, '}')
	return out, nil
}

func (o *IgnoredTestCases) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for IgnoredTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (o *IgnoredTestCases) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for IgnoredTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (o *IgnoredTestCases) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for IgnoredTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (o *IgnoredTestCases) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for IgnoredTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (o *IgnoredTestCases) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type IgnoredTestCases expected JSON object")
	}
	var seenClient bool
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "client":
			if strict {
				if err = o.Client.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field IgnoredTestCases[\"client\"]")
					return false
				}
			} else {
				if err = o.Client.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field IgnoredTestCases[\"client\"]")
					return false
				}
			}
			seenClient = true
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
	if !seenClient {
		missingFields = append(missingFields, "client")
	}
	if len(missingFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type IgnoredTestCases missing required JSON fields", werror.SafeParam("missingFields", missingFields))
	}
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type IgnoredTestCases encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

type PositiveAndNegativeTestCases struct {
	Positive []string `json:"positive"`
	Negative []string `json:"negative"`
}

func (o PositiveAndNegativeTestCases) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o PositiveAndNegativeTestCases) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"positive\":"...)
		out = append(out, '[')
		for i := range o.Positive {
			out = safejson.AppendQuotedString(out, o.Positive[i])
			if i < len(o.Positive)-1 {
				out = append(out, ',')
			}
		}
		out = append(out, ']')
		out = append(out, ',')
	}
	{
		out = append(out, "\"negative\":"...)
		out = append(out, '[')
		for i := range o.Negative {
			out = safejson.AppendQuotedString(out, o.Negative[i])
			if i < len(o.Negative)-1 {
				out = append(out, ',')
			}
		}
		out = append(out, ']')
	}
	out = append(out, '}')
	return out, nil
}

func (o *PositiveAndNegativeTestCases) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for PositiveAndNegativeTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (o *PositiveAndNegativeTestCases) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for PositiveAndNegativeTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (o *PositiveAndNegativeTestCases) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for PositiveAndNegativeTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (o *PositiveAndNegativeTestCases) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for PositiveAndNegativeTestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (o *PositiveAndNegativeTestCases) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type PositiveAndNegativeTestCases expected JSON object")
	}
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "positive":
			if !value.IsArray() {
				err = werror.ErrorWithContextParams(ctx, "field PositiveAndNegativeTestCases[\"positive\"] expected JSON array")
				return false
			}
			value.ForEach(func(_, value gjson.Result) bool {
				var listElement string
				if value.Type != gjson.String {
					err = werror.ErrorWithContextParams(ctx, "field PositiveAndNegativeTestCases[\"positive\"] list element expected JSON string")
					return false
				}
				listElement = value.Str
				o.Positive = append(o.Positive, listElement)
				return err == nil
			})
		case "negative":
			if !value.IsArray() {
				err = werror.ErrorWithContextParams(ctx, "field PositiveAndNegativeTestCases[\"negative\"] expected JSON array")
				return false
			}
			value.ForEach(func(_, value gjson.Result) bool {
				var listElement string
				if value.Type != gjson.String {
					err = werror.ErrorWithContextParams(ctx, "field PositiveAndNegativeTestCases[\"negative\"] list element expected JSON string")
					return false
				}
				listElement = value.Str
				o.Negative = append(o.Negative, listElement)
				return err == nil
			})
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
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type PositiveAndNegativeTestCases encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

type TestCases struct {
	Client ClientTestCases `json:"client"`
}

func (o TestCases) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o TestCases) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"client\":"...)
		var err error
		out, err = o.Client.AppendJSON(out)
		if err != nil {
			return nil, err
		}
	}
	out = append(out, '}')
	return out, nil
}

func (o *TestCases) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for TestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (o *TestCases) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for TestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (o *TestCases) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for TestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (o *TestCases) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for TestCases")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (o *TestCases) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type TestCases expected JSON object")
	}
	var seenClient bool
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "client":
			if strict {
				if err = o.Client.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field TestCases[\"client\"]")
					return false
				}
			} else {
				if err = o.Client.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field TestCases[\"client\"]")
					return false
				}
			}
			seenClient = true
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
	if !seenClient {
		missingFields = append(missingFields, "client")
	}
	if len(missingFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type TestCases missing required JSON fields", werror.SafeParam("missingFields", missingFields))
	}
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type TestCases encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}
