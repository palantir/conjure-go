// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"fmt"
	"strconv"

	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
	werror "github.com/palantir/witchcraft-go-error"
	gjson "github.com/tidwall/gjson"
)

type ExampleUnion struct {
	typ         string
	str         *string
	strOptional **string
	other       *int
}

func NewExampleUnionFromStr(v string) ExampleUnion {
	return ExampleUnion{typ: "str", str: &v}
}

func NewExampleUnionFromStrOptional(v *string) ExampleUnion {
	return ExampleUnion{typ: "strOptional", strOptional: &v}
}

func NewExampleUnionFromOther(v int) ExampleUnion {
	return ExampleUnion{typ: "other", other: &v}
}

type ExampleUnionVisitor interface {
	VisitStr(string) error
	VisitStrOptional(*string) error
	VisitOther(int) error
	VisitUnknown(typeName string) error
}

func (u *ExampleUnion) Accept(v ExampleUnionVisitor) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(u.typ)
	case "str":
		return v.VisitStr(*u.str)
	case "strOptional":
		var strOptional *string
		if u.strOptional != nil {
			strOptional = *u.strOptional
		}
		return v.VisitStrOptional(strOptional)
	case "other":
		return v.VisitOther(*u.other)
	}
}

type ExampleUnionVisitorWithContext interface {
	VisitStrWithContext(context.Context, string) error
	VisitStrOptionalWithContext(context.Context, *string) error
	VisitOtherWithContext(context.Context, int) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
}

func (u *ExampleUnion) AcceptWithContext(ctx context.Context, v ExampleUnionVisitorWithContext) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknownWithContext(ctx, u.typ)
	case "str":
		return v.VisitStrWithContext(ctx, *u.str)
	case "strOptional":
		var strOptional *string
		if u.strOptional != nil {
			strOptional = *u.strOptional
		}
		return v.VisitStrOptionalWithContext(ctx, strOptional)
	case "other":
		return v.VisitOtherWithContext(ctx, *u.other)
	}
}

func (u *ExampleUnion) AcceptFuncs(strFunc func(string) error, strOptionalFunc func(*string) error, otherFunc func(int) error, unknownFunc func(string) error) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return unknownFunc(u.typ)
	case "str":
		return strFunc(*u.str)
	case "strOptional":
		var strOptional *string
		if u.strOptional != nil {
			strOptional = *u.strOptional
		}
		return strOptionalFunc(strOptional)
	case "other":
		return otherFunc(*u.other)
	}
}

func (u *ExampleUnion) StrNoopSuccess(string) error {
	return nil
}

func (u *ExampleUnion) StrOptionalNoopSuccess(*string) error {
	return nil
}

func (u *ExampleUnion) OtherNoopSuccess(int) error {
	return nil
}

func (u *ExampleUnion) ErrorOnUnknown(typeName string) error {
	return fmt.Errorf("invalid value in union type. Type name: %s", typeName)
}

func (u ExampleUnion) MarshalJSON() ([]byte, error) {
	size, err := u.JSONSize()
	if err != nil {
		return nil, err
	}
	return u.AppendJSON(make([]byte, 0, size))
}

func (u ExampleUnion) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	switch u.typ {
	case "str":
		out = append(out, "\"type\":\"str\""...)
		if u.str != nil {
			out = append(out, ',')
			out = append(out, "\"str\""...)
			out = append(out, ':')
			unionVal := *u.str
			out = safejson.AppendQuotedString(out, unionVal)
		}
	case "strOptional":
		out = append(out, "\"type\":\"strOptional\""...)
		if u.strOptional != nil {
			out = append(out, ',')
			out = append(out, "\"strOptional\""...)
			out = append(out, ':')
			unionVal := *u.strOptional
			if unionVal != nil {
				optVal := *unionVal
				out = safejson.AppendQuotedString(out, optVal)
			} else {
				out = append(out, "null"...)
			}
		}
	case "other":
		out = append(out, "\"type\":\"other\""...)
		if u.other != nil {
			out = append(out, ',')
			out = append(out, "\"other\""...)
			out = append(out, ':')
			unionVal := *u.other
			out = strconv.AppendInt(out, int64(unionVal), 10)
		}
	default:
		out = append(out, "\"type\":"...)
		out = safejson.AppendQuotedString(out, u.typ)
	}
	out = append(out, '}')
	return out, nil
}

func (u ExampleUnion) JSONSize() (int, error) {
	var out int
	out += 1 // '{'
	switch u.typ {
	case "str":
		out += 12 // "type":"str"
		if u.str != nil {
			out += 1 // ','
			out += 5 // "str"
			out += 1 // ':'
			unionVal := *u.str
			out += safejson.QuotedStringLength(unionVal)
		}
	case "strOptional":
		out += 20 // "type":"strOptional"
		if u.strOptional != nil {
			out += 1  // ','
			out += 13 // "strOptional"
			out += 1  // ':'
			unionVal := *u.strOptional
			if unionVal != nil {
				optVal := *unionVal
				out += safejson.QuotedStringLength(optVal)
			} else {
				out += 4 // null
			}
		}
	case "other":
		out += 14 // "type":"other"
		if u.other != nil {
			out += 1 // ','
			out += 7 // "other"
			out += 1 // ':'
			unionVal := *u.other
			out += len(strconv.AppendInt(nil, int64(unionVal), 10))
		}
	default:
		out += 7 // "type":
		out += safejson.QuotedStringLength(u.typ)
	}
	out += 1 // '}'
	return out, nil
}

func (u *ExampleUnion) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ExampleUnion")
	}
	return u.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (u *ExampleUnion) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ExampleUnion")
	}
	return u.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (u *ExampleUnion) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ExampleUnion")
	}
	return u.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (u *ExampleUnion) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ExampleUnion")
	}
	return u.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (u *ExampleUnion) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type ExampleUnion expected JSON object")
	}
	var seenType bool
	var seenStr bool
	var seenStrOptional bool
	var seenOther bool
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "type":
			if seenType {
				err = werror.ErrorWithContextParams(ctx, "type ExampleUnion encountered duplicate \"type\" field")
				return false
			}
			seenType = true
			if value.Type != gjson.String {
				err = werror.ErrorWithContextParams(ctx, "field ExampleUnion[\"type\"] expected JSON string")
				return false
			}
			u.typ = value.Str
		case "str":
			if seenStr {
				err = werror.ErrorWithContextParams(ctx, "type ExampleUnion encountered duplicate \"str\" field")
				return false
			}
			seenStr = true
			var unionVal string
			if value.Type != gjson.String {
				err = werror.ErrorWithContextParams(ctx, "field ExampleUnion[\"str\"] expected JSON string")
				return false
			}
			unionVal = value.Str
			u.str = &unionVal
		case "strOptional":
			if seenStrOptional {
				err = werror.ErrorWithContextParams(ctx, "type ExampleUnion encountered duplicate \"strOptional\" field")
				return false
			}
			seenStrOptional = true
			var unionVal *string
			if value.Type != gjson.Null {
				var optVal string
				if value.Type != gjson.String {
					err = werror.ErrorWithContextParams(ctx, "field ExampleUnion[\"strOptional\"] expected JSON string")
					return false
				}
				optVal = value.Str
				unionVal = &optVal
			}
			u.strOptional = &unionVal
		case "other":
			if seenOther {
				err = werror.ErrorWithContextParams(ctx, "type ExampleUnion encountered duplicate \"other\" field")
				return false
			}
			seenOther = true
			var unionVal int
			if value.Type != gjson.Number {
				err = werror.ErrorWithContextParams(ctx, "field ExampleUnion[\"other\"] expected JSON number")
				return false
			}
			unionVal, err = strconv.Atoi(value.Raw)
			if err != nil {
				err = werror.WrapWithContextParams(ctx, err, "field ExampleUnion[\"other\"]")
				return false
			}
			u.other = &unionVal
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
	if !seenType {
		missingFields = append(missingFields, "type")
	}
	if u.typ == "str" && !seenStr {
		missingFields = append(missingFields, "str")
	}
	if u.typ == "other" && !seenOther {
		missingFields = append(missingFields, "other")
	}
	if len(missingFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type ExampleUnion missing required JSON fields", werror.SafeParam("missingFields", missingFields))
	}
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type ExampleUnion encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (u ExampleUnion) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(u)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (u *ExampleUnion) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&u)
}
