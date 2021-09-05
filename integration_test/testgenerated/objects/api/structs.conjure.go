// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strconv"

	binary "github.com/palantir/pkg/binary"
	boolean "github.com/palantir/pkg/boolean"
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
	uuid "github.com/palantir/pkg/uuid"
	werror "github.com/palantir/witchcraft-go-error"
	gjson "github.com/tidwall/gjson"
)

type AnyValue struct {
	Value interface{} `json:"value"`
}

func (o AnyValue) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o AnyValue) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		if o.Value == nil {
			out = append(out, "null"...)
		} else if appender, ok := o.Value.(interface {
			AppendJSON([]byte) ([]byte, error)
		}); ok {
			var err error
			out, err = appender.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		} else if marshaler, ok := o.Value.(json.Marshaler); ok {
			data, err := marshaler.MarshalJSON()
			if err != nil {
				return nil, err
			}
			out = append(out, data...)
		} else if data, err := safejson.Marshal(o.Value); err != nil {
			return nil, err
		} else {
			out = append(out, data...)
		}
	}
	out = append(out, '}')
	return out, nil
}

func (o *AnyValue) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for AnyValue")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (o *AnyValue) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for AnyValue")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (o *AnyValue) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for AnyValue")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (o *AnyValue) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for AnyValue")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (o *AnyValue) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type AnyValue expected JSON object")
	}
	var seenValue bool
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "value":
			if value.Type != gjson.JSON && value.Type != gjson.String && value.Type != gjson.Number && value.Type != gjson.True && value.Type != gjson.False {
				err = werror.ErrorWithContextParams(ctx, "field AnyValue[\"value\"] expected JSON non-null value")
				return false
			}
			o.Value = value.Value()
			seenValue = true
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
	if !seenValue {
		missingFields = append(missingFields, "value")
	}
	if len(missingFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type AnyValue missing required JSON fields", werror.SafeParam("missingFields", missingFields))
	}
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type AnyValue encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (o AnyValue) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *AnyValue) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

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
			if value.Type != gjson.String {
				err = werror.ErrorWithContextParams(ctx, "field Basic[\"data\"] expected JSON string")
				return false
			}
			o.Data = value.Str
			seenData = true
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

type BinaryMap struct {
	Map map[binary.Binary][]byte `json:"map"`
}

func (o BinaryMap) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o BinaryMap) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"map\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.Map {
				{
					out = safejson.AppendQuotedString(out, string(k))
				}
				out = append(out, ':')
				{
					out = append(out, '"')
					if len(v) > 0 {
						b64out := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
						base64.StdEncoding.Encode(b64out, v)
						out = append(out, b64out...)
					}
					out = append(out, '"')
				}
				i++
				if i < len(o.Map) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
	}
	out = append(out, '}')
	return out, nil
}

func (o *BinaryMap) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for BinaryMap")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (o *BinaryMap) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for BinaryMap")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (o *BinaryMap) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for BinaryMap")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (o *BinaryMap) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for BinaryMap")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (o *BinaryMap) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type BinaryMap expected JSON object")
	}
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "map":
			if !value.IsObject() {
				err = werror.ErrorWithContextParams(ctx, "field BinaryMap[\"map\"] expected JSON object")
				return false
			}
			if o.Map == nil {
				o.Map = make(map[binary.Binary][]byte, 0)
			}
			value.ForEach(func(key, value gjson.Result) bool {
				var mapKey binary.Binary
				{
					if key.Type != gjson.String {
						err = werror.ErrorWithContextParams(ctx, "field BinaryMap[\"map\"] map key expected JSON string")
						return false
					}
					mapKey = binary.Binary(key.Str)
				}
				if _, exists := o.Map[mapKey]; exists {
					err = werror.ErrorWithContextParams(ctx, "field BinaryMap[\"map\"] encountered duplicate map key")
					return false
				}
				var mapVal []byte
				{
					if value.Type != gjson.String {
						err = werror.ErrorWithContextParams(ctx, "field BinaryMap[\"map\"] map value expected JSON string")
						return false
					}
					mapVal, err = binary.Binary(value.Str).Bytes()
					if err != nil {
						err = werror.WrapWithContextParams(ctx, err, "field BinaryMap[\"map\"] map value")
						return false
					}
				}
				o.Map[mapKey] = mapVal
				return err == nil
			})
			if err != nil {
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
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type BinaryMap encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (o BinaryMap) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *BinaryMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type BooleanIntegerMap struct {
	Map map[boolean.Boolean]int `json:"map"`
}

func (o BooleanIntegerMap) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o BooleanIntegerMap) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"map\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.Map {
				{
					if k {
						out = append(out, "\"true\""...)
					} else {
						out = append(out, "\"false\""...)
					}
				}
				out = append(out, ':')
				{
					out = strconv.AppendInt(out, int64(v), 10)
				}
				i++
				if i < len(o.Map) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
	}
	out = append(out, '}')
	return out, nil
}

func (o *BooleanIntegerMap) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for BooleanIntegerMap")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (o *BooleanIntegerMap) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for BooleanIntegerMap")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (o *BooleanIntegerMap) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for BooleanIntegerMap")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (o *BooleanIntegerMap) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for BooleanIntegerMap")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (o *BooleanIntegerMap) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type BooleanIntegerMap expected JSON object")
	}
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "map":
			if !value.IsObject() {
				err = werror.ErrorWithContextParams(ctx, "field BooleanIntegerMap[\"map\"] expected JSON object")
				return false
			}
			if o.Map == nil {
				o.Map = make(map[boolean.Boolean]int, 0)
			}
			value.ForEach(func(key, value gjson.Result) bool {
				var mapKey boolean.Boolean
				{
					if key.Type != gjson.String {
						err = werror.ErrorWithContextParams(ctx, "field BooleanIntegerMap[\"map\"] map key expected JSON string")
						return false
					}
					var boolVal bool
					boolVal, err = strconv.ParseBool(key.Str)
					if err != nil {
						err = werror.WrapWithContextParams(ctx, err, "field BooleanIntegerMap[\"map\"] map key")
						return false
					}
					mapKey = boolean.Boolean(boolVal)
				}
				if _, exists := o.Map[mapKey]; exists {
					err = werror.ErrorWithContextParams(ctx, "field BooleanIntegerMap[\"map\"] encountered duplicate map key")
					return false
				}
				var mapVal int
				{
					if value.Type != gjson.Number {
						err = werror.ErrorWithContextParams(ctx, "field BooleanIntegerMap[\"map\"] map value expected JSON number")
						return false
					}
					mapVal, err = strconv.Atoi(value.Raw)
					if err != nil {
						err = werror.WrapWithContextParams(ctx, err, "field BooleanIntegerMap[\"map\"] map value")
						return false
					}
				}
				o.Map[mapKey] = mapVal
				return err == nil
			})
			if err != nil {
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
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type BooleanIntegerMap encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (o BooleanIntegerMap) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *BooleanIntegerMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type Collections struct {
	MapVar   map[string][]int   `json:"mapVar"`
	ListVar  []string           `json:"listVar"`
	MultiDim [][]map[string]int `json:"multiDim"`
}

func (o Collections) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o Collections) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"mapVar\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.MapVar {
				{
					out = safejson.AppendQuotedString(out, k)
				}
				out = append(out, ':')
				{
					out = append(out, '[')
					for i1 := range v {
						out = strconv.AppendInt(out, int64(v[i1]), 10)
						if i1 < len(v)-1 {
							out = append(out, ',')
						}
					}
					out = append(out, ']')
				}
				i++
				if i < len(o.MapVar) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
		out = append(out, ',')
	}
	{
		out = append(out, "\"listVar\":"...)
		out = append(out, '[')
		for i := range o.ListVar {
			out = safejson.AppendQuotedString(out, o.ListVar[i])
			if i < len(o.ListVar)-1 {
				out = append(out, ',')
			}
		}
		out = append(out, ']')
		out = append(out, ',')
	}
	{
		out = append(out, "\"multiDim\":"...)
		out = append(out, '[')
		for i := range o.MultiDim {
			out = append(out, '[')
			for i1 := range o.MultiDim[i] {
				out = append(out, '{')
				{
					var i2 int
					for k, v := range o.MultiDim[i][i1] {
						{
							out = safejson.AppendQuotedString(out, k)
						}
						out = append(out, ':')
						{
							out = strconv.AppendInt(out, int64(v), 10)
						}
						i2++
						if i2 < len(o.MultiDim[i][i1]) {
							out = append(out, ',')
						}
					}
				}
				out = append(out, '}')
				if i1 < len(o.MultiDim[i])-1 {
					out = append(out, ',')
				}
			}
			out = append(out, ']')
			if i < len(o.MultiDim)-1 {
				out = append(out, ',')
			}
		}
		out = append(out, ']')
	}
	out = append(out, '}')
	return out, nil
}

func (o *Collections) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Collections")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (o *Collections) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Collections")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (o *Collections) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Collections")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (o *Collections) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Collections")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (o *Collections) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type Collections expected JSON object")
	}
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "mapVar":
			if !value.IsObject() {
				err = werror.ErrorWithContextParams(ctx, "field Collections[\"mapVar\"] expected JSON object")
				return false
			}
			if o.MapVar == nil {
				o.MapVar = make(map[string][]int, 0)
			}
			value.ForEach(func(key, value gjson.Result) bool {
				var mapKey string
				{
					if key.Type != gjson.String {
						err = werror.ErrorWithContextParams(ctx, "field Collections[\"mapVar\"] map key expected JSON string")
						return false
					}
					mapKey = key.Str
				}
				if _, exists := o.MapVar[mapKey]; exists {
					err = werror.ErrorWithContextParams(ctx, "field Collections[\"mapVar\"] encountered duplicate map key")
					return false
				}
				var mapVal []int
				{
					if !value.IsArray() {
						err = werror.ErrorWithContextParams(ctx, "field Collections[\"mapVar\"] map value expected JSON array")
						return false
					}
					value.ForEach(func(_, value gjson.Result) bool {
						var listElement1 int
						if value.Type != gjson.Number {
							err = werror.ErrorWithContextParams(ctx, "field Collections[\"mapVar\"] map value list element expected JSON number")
							return false
						}
						listElement1, err = strconv.Atoi(value.Raw)
						if err != nil {
							err = werror.WrapWithContextParams(ctx, err, "field Collections[\"mapVar\"] map value list element")
							return false
						}
						mapVal = append(mapVal, listElement1)
						return err == nil
					})
					if err != nil {
						return false
					}
				}
				o.MapVar[mapKey] = mapVal
				return err == nil
			})
			if err != nil {
				return false
			}
		case "listVar":
			if !value.IsArray() {
				err = werror.ErrorWithContextParams(ctx, "field Collections[\"listVar\"] expected JSON array")
				return false
			}
			value.ForEach(func(_, value gjson.Result) bool {
				var listElement string
				if value.Type != gjson.String {
					err = werror.ErrorWithContextParams(ctx, "field Collections[\"listVar\"] list element expected JSON string")
					return false
				}
				listElement = value.Str
				o.ListVar = append(o.ListVar, listElement)
				return err == nil
			})
			if err != nil {
				return false
			}
		case "multiDim":
			if !value.IsArray() {
				err = werror.ErrorWithContextParams(ctx, "field Collections[\"multiDim\"] expected JSON array")
				return false
			}
			value.ForEach(func(_, value gjson.Result) bool {
				var listElement []map[string]int
				if !value.IsArray() {
					err = werror.ErrorWithContextParams(ctx, "field Collections[\"multiDim\"] list element expected JSON array")
					return false
				}
				value.ForEach(func(_, value gjson.Result) bool {
					var listElement1 map[string]int
					if !value.IsObject() {
						err = werror.ErrorWithContextParams(ctx, "field Collections[\"multiDim\"] list element list element expected JSON object")
						return false
					}
					if listElement1 == nil {
						listElement1 = make(map[string]int, 0)
					}
					value.ForEach(func(key, value gjson.Result) bool {
						var mapKey2 string
						{
							if key.Type != gjson.String {
								err = werror.ErrorWithContextParams(ctx, "field Collections[\"multiDim\"] list element list element map key expected JSON string")
								return false
							}
							mapKey2 = key.Str
						}
						if _, exists := listElement1[mapKey2]; exists {
							err = werror.ErrorWithContextParams(ctx, "field Collections[\"multiDim\"] list element list element encountered duplicate map key")
							return false
						}
						var mapVal2 int
						{
							if value.Type != gjson.Number {
								err = werror.ErrorWithContextParams(ctx, "field Collections[\"multiDim\"] list element list element map value expected JSON number")
								return false
							}
							mapVal2, err = strconv.Atoi(value.Raw)
							if err != nil {
								err = werror.WrapWithContextParams(ctx, err, "field Collections[\"multiDim\"] list element list element map value")
								return false
							}
						}
						listElement1[mapKey2] = mapVal2
						return err == nil
					})
					if err != nil {
						return false
					}
					listElement = append(listElement, listElement1)
					return err == nil
				})
				if err != nil {
					return false
				}
				o.MultiDim = append(o.MultiDim, listElement)
				return err == nil
			})
			if err != nil {
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
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type Collections encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (o Collections) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *Collections) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type Compound struct {
	Obj Collections `json:"obj"`
}

func (o Compound) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o Compound) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"obj\":"...)
		var err error
		out, err = o.Obj.AppendJSON(out)
		if err != nil {
			return nil, err
		}
	}
	out = append(out, '}')
	return out, nil
}

func (o *Compound) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Compound")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (o *Compound) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Compound")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (o *Compound) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Compound")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (o *Compound) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Compound")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (o *Compound) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type Compound expected JSON object")
	}
	var seenObj bool
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "obj":
			if strict {
				if err = o.Obj.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Compound[\"obj\"]")
					return false
				}
			} else {
				if err = o.Obj.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Compound[\"obj\"]")
					return false
				}
			}
			seenObj = true
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
	if !seenObj {
		missingFields = append(missingFields, "obj")
	}
	if len(missingFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type Compound missing required JSON fields", werror.SafeParam("missingFields", missingFields))
	}
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type Compound encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (o Compound) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *Compound) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type ExampleUuid struct {
	Uid uuid.UUID `json:"uid"`
}

func (o ExampleUuid) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o ExampleUuid) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"uid\":"...)
		out = safejson.AppendQuotedString(out, o.Uid.String())
	}
	out = append(out, '}')
	return out, nil
}

func (o *ExampleUuid) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ExampleUuid")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (o *ExampleUuid) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ExampleUuid")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (o *ExampleUuid) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ExampleUuid")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (o *ExampleUuid) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ExampleUuid")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (o *ExampleUuid) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type ExampleUuid expected JSON object")
	}
	var seenUid bool
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "uid":
			if value.Type != gjson.String {
				err = werror.ErrorWithContextParams(ctx, "field ExampleUuid[\"uid\"] expected JSON string")
				return false
			}
			o.Uid, err = uuid.ParseUUID(value.Str)
			if err != nil {
				err = werror.WrapWithContextParams(ctx, err, "field ExampleUuid[\"uid\"]")
				return false
			}
			seenUid = true
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
	if !seenUid {
		missingFields = append(missingFields, "uid")
	}
	if len(missingFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type ExampleUuid missing required JSON fields", werror.SafeParam("missingFields", missingFields))
	}
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type ExampleUuid encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (o ExampleUuid) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *ExampleUuid) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type MapOptional struct {
	Map map[OptionalUuidAlias]string `json:"map"`
}

func (o MapOptional) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o MapOptional) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"map\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.Map {
				{
					var err error
					out, err = k.AppendJSON(out)
					if err != nil {
						return nil, err
					}
				}
				out = append(out, ':')
				{
					out = safejson.AppendQuotedString(out, v)
				}
				i++
				if i < len(o.Map) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
	}
	out = append(out, '}')
	return out, nil
}

func (o *MapOptional) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for MapOptional")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (o *MapOptional) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for MapOptional")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (o *MapOptional) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for MapOptional")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (o *MapOptional) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for MapOptional")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (o *MapOptional) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type MapOptional expected JSON object")
	}
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "map":
			if !value.IsObject() {
				err = werror.ErrorWithContextParams(ctx, "field MapOptional[\"map\"] expected JSON object")
				return false
			}
			if o.Map == nil {
				o.Map = make(map[OptionalUuidAlias]string, 0)
			}
			value.ForEach(func(key, value gjson.Result) bool {
				var mapKey OptionalUuidAlias
				{
					if err = mapKey.UnmarshalJSONString(key.Raw); err != nil {
						err = werror.WrapWithContextParams(ctx, err, "field MapOptional[\"map\"] map key")
						return false
					}
				}
				if _, exists := o.Map[mapKey]; exists {
					err = werror.ErrorWithContextParams(ctx, "field MapOptional[\"map\"] encountered duplicate map key")
					return false
				}
				var mapVal string
				{
					if value.Type != gjson.String {
						err = werror.ErrorWithContextParams(ctx, "field MapOptional[\"map\"] map value expected JSON string")
						return false
					}
					mapVal = value.Str
				}
				o.Map[mapKey] = mapVal
				return err == nil
			})
			if err != nil {
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
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type MapOptional encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (o MapOptional) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *MapOptional) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

// A type using go keywords
type Type struct {
	Type []string          `json:"type"`
	Chan map[string]string `json:"chan"`
}

func (o Type) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o Type) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"type\":"...)
		out = append(out, '[')
		for i := range o.Type {
			out = safejson.AppendQuotedString(out, o.Type[i])
			if i < len(o.Type)-1 {
				out = append(out, ',')
			}
		}
		out = append(out, ']')
		out = append(out, ',')
	}
	{
		out = append(out, "\"chan\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.Chan {
				{
					out = safejson.AppendQuotedString(out, k)
				}
				out = append(out, ':')
				{
					out = safejson.AppendQuotedString(out, v)
				}
				i++
				if i < len(o.Chan) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
	}
	out = append(out, '}')
	return out, nil
}

func (o *Type) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Type")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (o *Type) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Type")
	}
	return o.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (o *Type) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Type")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (o *Type) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Type")
	}
	return o.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (o *Type) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type Type expected JSON object")
	}
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "type":
			if !value.IsArray() {
				err = werror.ErrorWithContextParams(ctx, "field Type[\"type\"] expected JSON array")
				return false
			}
			value.ForEach(func(_, value gjson.Result) bool {
				var listElement string
				if value.Type != gjson.String {
					err = werror.ErrorWithContextParams(ctx, "field Type[\"type\"] list element expected JSON string")
					return false
				}
				listElement = value.Str
				o.Type = append(o.Type, listElement)
				return err == nil
			})
			if err != nil {
				return false
			}
		case "chan":
			if !value.IsObject() {
				err = werror.ErrorWithContextParams(ctx, "field Type[\"chan\"] expected JSON object")
				return false
			}
			if o.Chan == nil {
				o.Chan = make(map[string]string, 0)
			}
			value.ForEach(func(key, value gjson.Result) bool {
				var mapKey string
				{
					if key.Type != gjson.String {
						err = werror.ErrorWithContextParams(ctx, "field Type[\"chan\"] map key expected JSON string")
						return false
					}
					mapKey = key.Str
				}
				if _, exists := o.Chan[mapKey]; exists {
					err = werror.ErrorWithContextParams(ctx, "field Type[\"chan\"] encountered duplicate map key")
					return false
				}
				var mapVal string
				{
					if value.Type != gjson.String {
						err = werror.ErrorWithContextParams(ctx, "field Type[\"chan\"] map value expected JSON string")
						return false
					}
					mapVal = value.Str
				}
				o.Chan[mapKey] = mapVal
				return err == nil
			})
			if err != nil {
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
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type Type encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (o Type) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
