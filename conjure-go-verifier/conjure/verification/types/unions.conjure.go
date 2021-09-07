// This file was generated by Conjure and should not be manually edited.

package types

import (
	"context"
	"fmt"
	"strconv"

	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
	werror "github.com/palantir/witchcraft-go-error"
	gjson "github.com/tidwall/gjson"
)

// A type which can either be a StringExample, a set of strings, or an integer.
type Union struct {
	typ                  string
	stringExample        *StringExample
	set                  *[]string
	thisFieldIsAnInteger *int
	alsoAnInteger        *int
	if_                  *int
	new                  *int
	interface_           *int
}

func NewUnionFromStringExample(v StringExample) Union {
	return Union{typ: "stringExample", stringExample: &v}
}

func NewUnionFromSet(v []string) Union {
	return Union{typ: "set", set: &v}
}

func NewUnionFromThisFieldIsAnInteger(v int) Union {
	return Union{typ: "thisFieldIsAnInteger", thisFieldIsAnInteger: &v}
}

func NewUnionFromAlsoAnInteger(v int) Union {
	return Union{typ: "alsoAnInteger", alsoAnInteger: &v}
}

func NewUnionFromIf(v int) Union {
	return Union{typ: "if", if_: &v}
}

func NewUnionFromNew(v int) Union {
	return Union{typ: "new", new: &v}
}

func NewUnionFromInterface(v int) Union {
	return Union{typ: "interface", interface_: &v}
}

type UnionVisitor interface {
	VisitStringExample(StringExample) error
	VisitSet([]string) error
	VisitThisFieldIsAnInteger(int) error
	VisitAlsoAnInteger(int) error
	VisitIf(int) error
	VisitNew(int) error
	VisitInterface(int) error
	VisitUnknown(typeName string) error
}

func (u *Union) Accept(v UnionVisitor) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(u.typ)
	case "stringExample":
		return v.VisitStringExample(*u.stringExample)
	case "set":
		return v.VisitSet(*u.set)
	case "thisFieldIsAnInteger":
		return v.VisitThisFieldIsAnInteger(*u.thisFieldIsAnInteger)
	case "alsoAnInteger":
		return v.VisitAlsoAnInteger(*u.alsoAnInteger)
	case "if":
		return v.VisitIf(*u.if_)
	case "new":
		return v.VisitNew(*u.new)
	case "interface":
		return v.VisitInterface(*u.interface_)
	}
}

type UnionVisitorWithContext interface {
	VisitStringExampleWithContext(context.Context, StringExample) error
	VisitSetWithContext(context.Context, []string) error
	VisitThisFieldIsAnIntegerWithContext(context.Context, int) error
	VisitAlsoAnIntegerWithContext(context.Context, int) error
	VisitIfWithContext(context.Context, int) error
	VisitNewWithContext(context.Context, int) error
	VisitInterfaceWithContext(context.Context, int) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
}

func (u *Union) AcceptWithContext(ctx context.Context, v UnionVisitorWithContext) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknownWithContext(ctx, u.typ)
	case "stringExample":
		return v.VisitStringExampleWithContext(ctx, *u.stringExample)
	case "set":
		return v.VisitSetWithContext(ctx, *u.set)
	case "thisFieldIsAnInteger":
		return v.VisitThisFieldIsAnIntegerWithContext(ctx, *u.thisFieldIsAnInteger)
	case "alsoAnInteger":
		return v.VisitAlsoAnIntegerWithContext(ctx, *u.alsoAnInteger)
	case "if":
		return v.VisitIfWithContext(ctx, *u.if_)
	case "new":
		return v.VisitNewWithContext(ctx, *u.new)
	case "interface":
		return v.VisitInterfaceWithContext(ctx, *u.interface_)
	}
}

func (u Union) MarshalJSON() ([]byte, error) {
	return u.AppendJSON(nil)
}

func (u Union) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	switch u.typ {
	case "stringExample":
		out = append(out, "\"type\":\"stringExample\""...)
		if u.stringExample != nil {
			out = append(out, ',')
			out = append(out, "\"stringExample\""...)
			out = append(out, ':')
			unionVal := *u.stringExample
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	case "set":
		out = append(out, "\"type\":\"set\""...)
		if u.set != nil {
			out = append(out, ',')
			out = append(out, "\"set\""...)
			out = append(out, ':')
			unionVal := *u.set
			out = append(out, '[')
			for i := range unionVal {
				out = safejson.AppendQuotedString(out, unionVal[i])
				if i < len(unionVal)-1 {
					out = append(out, ',')
				}
			}
			out = append(out, ']')
		}
	case "thisFieldIsAnInteger":
		out = append(out, "\"type\":\"thisFieldIsAnInteger\""...)
		if u.thisFieldIsAnInteger != nil {
			out = append(out, ',')
			out = append(out, "\"thisFieldIsAnInteger\""...)
			out = append(out, ':')
			unionVal := *u.thisFieldIsAnInteger
			out = strconv.AppendInt(out, int64(unionVal), 10)
		}
	case "alsoAnInteger":
		out = append(out, "\"type\":\"alsoAnInteger\""...)
		if u.alsoAnInteger != nil {
			out = append(out, ',')
			out = append(out, "\"alsoAnInteger\""...)
			out = append(out, ':')
			unionVal := *u.alsoAnInteger
			out = strconv.AppendInt(out, int64(unionVal), 10)
		}
	case "if":
		out = append(out, "\"type\":\"if\""...)
		if u.if_ != nil {
			out = append(out, ',')
			out = append(out, "\"if\""...)
			out = append(out, ':')
			unionVal := *u.if_
			out = strconv.AppendInt(out, int64(unionVal), 10)
		}
	case "new":
		out = append(out, "\"type\":\"new\""...)
		if u.new != nil {
			out = append(out, ',')
			out = append(out, "\"new\""...)
			out = append(out, ':')
			unionVal := *u.new
			out = strconv.AppendInt(out, int64(unionVal), 10)
		}
	case "interface":
		out = append(out, "\"type\":\"interface\""...)
		if u.interface_ != nil {
			out = append(out, ',')
			out = append(out, "\"interface\""...)
			out = append(out, ':')
			unionVal := *u.interface_
			out = strconv.AppendInt(out, int64(unionVal), 10)
		}
	default:
		out = append(out, "\"type\":"...)
		out = safejson.AppendQuotedString(out, u.typ)
	}
	out = append(out, '}')
	return out, nil
}

func (u *Union) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Union")
	}
	return u.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (u *Union) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Union")
	}
	return u.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (u *Union) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Union")
	}
	return u.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (u *Union) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Union")
	}
	return u.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (u *Union) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type Union expected JSON object")
	}
	var seenType bool
	var seenStringExample bool
	var seenSet bool
	var seenThisFieldIsAnInteger bool
	var seenAlsoAnInteger bool
	var seenIf bool
	var seenNew bool
	var seenInterface bool
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "type":
			if seenType {
				err = werror.ErrorWithContextParams(ctx, "type Union encountered duplicate \"type\" field")
				return false
			} else {
				seenType = true
			}
			if value.Type != gjson.String {
				err = werror.ErrorWithContextParams(ctx, "field Union[\"type\"] expected JSON string")
				return false
			}
			u.typ = value.Str
		case "stringExample":
			if seenStringExample {
				err = werror.ErrorWithContextParams(ctx, "type Union encountered duplicate \"stringExample\" field")
				return false
			} else {
				seenStringExample = true
			}
			var unionVal StringExample
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Union[\"stringExample\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Union[\"stringExample\"]")
					return false
				}
			}
			u.stringExample = &unionVal
		case "set":
			if seenSet {
				err = werror.ErrorWithContextParams(ctx, "type Union encountered duplicate \"set\" field")
				return false
			} else {
				seenSet = true
			}
			var unionVal []string
			if !value.IsArray() {
				err = werror.ErrorWithContextParams(ctx, "field Union[\"set\"] expected JSON array")
				return false
			}
			value.ForEach(func(_, value gjson.Result) bool {
				var listElement string
				if value.Type != gjson.String {
					err = werror.ErrorWithContextParams(ctx, "field Union[\"set\"] list element expected JSON string")
					return false
				}
				listElement = value.Str
				unionVal = append(unionVal, listElement)
				return err == nil
			})
			if err != nil {
				return false
			}
			u.set = &unionVal
		case "thisFieldIsAnInteger":
			if seenThisFieldIsAnInteger {
				err = werror.ErrorWithContextParams(ctx, "type Union encountered duplicate \"thisFieldIsAnInteger\" field")
				return false
			} else {
				seenThisFieldIsAnInteger = true
			}
			var unionVal int
			if value.Type != gjson.Number {
				err = werror.ErrorWithContextParams(ctx, "field Union[\"thisFieldIsAnInteger\"] expected JSON number")
				return false
			}
			unionVal, err = strconv.Atoi(value.Raw)
			if err != nil {
				err = werror.WrapWithContextParams(ctx, err, "field Union[\"thisFieldIsAnInteger\"]")
				return false
			}
			u.thisFieldIsAnInteger = &unionVal
		case "alsoAnInteger":
			if seenAlsoAnInteger {
				err = werror.ErrorWithContextParams(ctx, "type Union encountered duplicate \"alsoAnInteger\" field")
				return false
			} else {
				seenAlsoAnInteger = true
			}
			var unionVal int
			if value.Type != gjson.Number {
				err = werror.ErrorWithContextParams(ctx, "field Union[\"alsoAnInteger\"] expected JSON number")
				return false
			}
			unionVal, err = strconv.Atoi(value.Raw)
			if err != nil {
				err = werror.WrapWithContextParams(ctx, err, "field Union[\"alsoAnInteger\"]")
				return false
			}
			u.alsoAnInteger = &unionVal
		case "if":
			if seenIf {
				err = werror.ErrorWithContextParams(ctx, "type Union encountered duplicate \"if\" field")
				return false
			} else {
				seenIf = true
			}
			var unionVal int
			if value.Type != gjson.Number {
				err = werror.ErrorWithContextParams(ctx, "field Union[\"if\"] expected JSON number")
				return false
			}
			unionVal, err = strconv.Atoi(value.Raw)
			if err != nil {
				err = werror.WrapWithContextParams(ctx, err, "field Union[\"if\"]")
				return false
			}
			u.if_ = &unionVal
		case "new":
			if seenNew {
				err = werror.ErrorWithContextParams(ctx, "type Union encountered duplicate \"new\" field")
				return false
			} else {
				seenNew = true
			}
			var unionVal int
			if value.Type != gjson.Number {
				err = werror.ErrorWithContextParams(ctx, "field Union[\"new\"] expected JSON number")
				return false
			}
			unionVal, err = strconv.Atoi(value.Raw)
			if err != nil {
				err = werror.WrapWithContextParams(ctx, err, "field Union[\"new\"]")
				return false
			}
			u.new = &unionVal
		case "interface":
			if seenInterface {
				err = werror.ErrorWithContextParams(ctx, "type Union encountered duplicate \"interface\" field")
				return false
			} else {
				seenInterface = true
			}
			var unionVal int
			if value.Type != gjson.Number {
				err = werror.ErrorWithContextParams(ctx, "field Union[\"interface\"] expected JSON number")
				return false
			}
			unionVal, err = strconv.Atoi(value.Raw)
			if err != nil {
				err = werror.WrapWithContextParams(ctx, err, "field Union[\"interface\"]")
				return false
			}
			u.interface_ = &unionVal
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
	if u.typ == "stringExample" && !seenStringExample {
		missingFields = append(missingFields, "stringExample")
	}
	if u.typ == "thisFieldIsAnInteger" && !seenThisFieldIsAnInteger {
		missingFields = append(missingFields, "thisFieldIsAnInteger")
	}
	if u.typ == "alsoAnInteger" && !seenAlsoAnInteger {
		missingFields = append(missingFields, "alsoAnInteger")
	}
	if u.typ == "if" && !seenIf {
		missingFields = append(missingFields, "if")
	}
	if u.typ == "new" && !seenNew {
		missingFields = append(missingFields, "new")
	}
	if u.typ == "interface" && !seenInterface {
		missingFields = append(missingFields, "interface")
	}
	if len(missingFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type Union missing required JSON fields", werror.SafeParam("missingFields", missingFields))
	}
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type Union encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (u Union) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(u)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (u *Union) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&u)
}
