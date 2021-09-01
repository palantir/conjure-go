// This file was generated by Conjure and should not be manually edited.

package spec

import (
	"context"
	"fmt"

	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
	werror "github.com/palantir/witchcraft-go-error"
	gjson "github.com/tidwall/gjson"
)

type AuthType struct {
	typ    string
	header *HeaderAuthType
	cookie *CookieAuthType
}

func NewAuthTypeFromHeader(v HeaderAuthType) AuthType {
	return AuthType{typ: "header", header: &v}
}

func NewAuthTypeFromCookie(v CookieAuthType) AuthType {
	return AuthType{typ: "cookie", cookie: &v}
}

type AuthTypeVisitor interface {
	VisitHeader(HeaderAuthType) error
	VisitCookie(CookieAuthType) error
	VisitUnknown(typeName string) error
}

func (u *AuthType) Accept(v AuthTypeVisitor) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(u.typ)
	case "header":
		return v.VisitHeader(*u.header)
	case "cookie":
		return v.VisitCookie(*u.cookie)
	}
}

type AuthTypeVisitorWithContext interface {
	VisitHeaderWithContext(context.Context, HeaderAuthType) error
	VisitCookieWithContext(context.Context, CookieAuthType) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
}

func (u *AuthType) AcceptWithContext(ctx context.Context, v AuthTypeVisitorWithContext) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknownWithContext(ctx, u.typ)
	case "header":
		return v.VisitHeaderWithContext(ctx, *u.header)
	case "cookie":
		return v.VisitCookieWithContext(ctx, *u.cookie)
	}
}

func (u *AuthType) AcceptFuncs(headerFunc func(HeaderAuthType) error, cookieFunc func(CookieAuthType) error, unknownFunc func(string) error) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return unknownFunc(u.typ)
	case "header":
		return headerFunc(*u.header)
	case "cookie":
		return cookieFunc(*u.cookie)
	}
}

func (u *AuthType) HeaderNoopSuccess(HeaderAuthType) error {
	return nil
}

func (u *AuthType) CookieNoopSuccess(CookieAuthType) error {
	return nil
}

func (u *AuthType) ErrorOnUnknown(typeName string) error {
	return fmt.Errorf("invalid value in union type. Type name: %s", typeName)
}

func (u AuthType) MarshalJSON() ([]byte, error) {
	return u.AppendJSON(nil)
}

func (u AuthType) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	switch u.typ {
	default:
		out = append(out, "\"type\":"...)
		out = safejson.AppendQuotedString(out, u.typ)
	case "header":
		out = append(out, "\"type\":\"header\""...)
		if u.header != nil {
			out = append(out, ',')
			out = append(out, "\"header\""...)
			out = append(out, ':')
			unionVal := *u.header
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	case "cookie":
		out = append(out, "\"type\":\"cookie\""...)
		if u.cookie != nil {
			out = append(out, ',')
			out = append(out, "\"cookie\""...)
			out = append(out, ':')
			unionVal := *u.cookie
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	}
	out = append(out, '}')
	return out, nil
}

func (u *AuthType) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for AuthType")
	}
	return u.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (u *AuthType) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for AuthType")
	}
	return u.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (u *AuthType) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for AuthType")
	}
	return u.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (u *AuthType) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for AuthType")
	}
	return u.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (u *AuthType) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type AuthType expected JSON object")
	}
	var seenType bool
	var seenHeader bool
	var seenCookie bool
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "type":
			if value.Type != gjson.String {
				err = werror.ErrorWithContextParams(ctx, "field AuthType[\"type\"] expected JSON string")
				return false
			}
			u.typ = value.Str
			seenType = true
		case "header":
			var unionVal HeaderAuthType
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field AuthType[\"header\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field AuthType[\"header\"]")
					return false
				}
			}
			u.header = &unionVal
			seenHeader = true
		case "cookie":
			var unionVal CookieAuthType
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field AuthType[\"cookie\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field AuthType[\"cookie\"]")
					return false
				}
			}
			u.cookie = &unionVal
			seenCookie = true
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
	if u.typ == "header" && !seenHeader {
		missingFields = append(missingFields, "header")
	}
	if u.typ == "cookie" && !seenCookie {
		missingFields = append(missingFields, "cookie")
	}
	if len(missingFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type AuthType missing required JSON fields", werror.SafeParam("missingFields", missingFields))
	}
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type AuthType encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (u AuthType) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(u)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (u *AuthType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&u)
}

type ParameterType struct {
	typ    string
	body   *BodyParameterType
	header *HeaderParameterType
	path   *PathParameterType
	query  *QueryParameterType
}

func NewParameterTypeFromBody(v BodyParameterType) ParameterType {
	return ParameterType{typ: "body", body: &v}
}

func NewParameterTypeFromHeader(v HeaderParameterType) ParameterType {
	return ParameterType{typ: "header", header: &v}
}

func NewParameterTypeFromPath(v PathParameterType) ParameterType {
	return ParameterType{typ: "path", path: &v}
}

func NewParameterTypeFromQuery(v QueryParameterType) ParameterType {
	return ParameterType{typ: "query", query: &v}
}

type ParameterTypeVisitor interface {
	VisitBody(BodyParameterType) error
	VisitHeader(HeaderParameterType) error
	VisitPath(PathParameterType) error
	VisitQuery(QueryParameterType) error
	VisitUnknown(typeName string) error
}

func (u *ParameterType) Accept(v ParameterTypeVisitor) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(u.typ)
	case "body":
		return v.VisitBody(*u.body)
	case "header":
		return v.VisitHeader(*u.header)
	case "path":
		return v.VisitPath(*u.path)
	case "query":
		return v.VisitQuery(*u.query)
	}
}

type ParameterTypeVisitorWithContext interface {
	VisitBodyWithContext(context.Context, BodyParameterType) error
	VisitHeaderWithContext(context.Context, HeaderParameterType) error
	VisitPathWithContext(context.Context, PathParameterType) error
	VisitQueryWithContext(context.Context, QueryParameterType) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
}

func (u *ParameterType) AcceptWithContext(ctx context.Context, v ParameterTypeVisitorWithContext) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknownWithContext(ctx, u.typ)
	case "body":
		return v.VisitBodyWithContext(ctx, *u.body)
	case "header":
		return v.VisitHeaderWithContext(ctx, *u.header)
	case "path":
		return v.VisitPathWithContext(ctx, *u.path)
	case "query":
		return v.VisitQueryWithContext(ctx, *u.query)
	}
}

func (u *ParameterType) AcceptFuncs(bodyFunc func(BodyParameterType) error, headerFunc func(HeaderParameterType) error, pathFunc func(PathParameterType) error, queryFunc func(QueryParameterType) error, unknownFunc func(string) error) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return unknownFunc(u.typ)
	case "body":
		return bodyFunc(*u.body)
	case "header":
		return headerFunc(*u.header)
	case "path":
		return pathFunc(*u.path)
	case "query":
		return queryFunc(*u.query)
	}
}

func (u *ParameterType) BodyNoopSuccess(BodyParameterType) error {
	return nil
}

func (u *ParameterType) HeaderNoopSuccess(HeaderParameterType) error {
	return nil
}

func (u *ParameterType) PathNoopSuccess(PathParameterType) error {
	return nil
}

func (u *ParameterType) QueryNoopSuccess(QueryParameterType) error {
	return nil
}

func (u *ParameterType) ErrorOnUnknown(typeName string) error {
	return fmt.Errorf("invalid value in union type. Type name: %s", typeName)
}

func (u ParameterType) MarshalJSON() ([]byte, error) {
	return u.AppendJSON(nil)
}

func (u ParameterType) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	switch u.typ {
	default:
		out = append(out, "\"type\":"...)
		out = safejson.AppendQuotedString(out, u.typ)
	case "body":
		out = append(out, "\"type\":\"body\""...)
		if u.body != nil {
			out = append(out, ',')
			out = append(out, "\"body\""...)
			out = append(out, ':')
			unionVal := *u.body
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	case "header":
		out = append(out, "\"type\":\"header\""...)
		if u.header != nil {
			out = append(out, ',')
			out = append(out, "\"header\""...)
			out = append(out, ':')
			unionVal := *u.header
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	case "path":
		out = append(out, "\"type\":\"path\""...)
		if u.path != nil {
			out = append(out, ',')
			out = append(out, "\"path\""...)
			out = append(out, ':')
			unionVal := *u.path
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	case "query":
		out = append(out, "\"type\":\"query\""...)
		if u.query != nil {
			out = append(out, ',')
			out = append(out, "\"query\""...)
			out = append(out, ':')
			unionVal := *u.query
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	}
	out = append(out, '}')
	return out, nil
}

func (u *ParameterType) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ParameterType")
	}
	return u.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (u *ParameterType) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ParameterType")
	}
	return u.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (u *ParameterType) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ParameterType")
	}
	return u.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (u *ParameterType) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for ParameterType")
	}
	return u.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (u *ParameterType) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type ParameterType expected JSON object")
	}
	var seenType bool
	var seenBody bool
	var seenHeader bool
	var seenPath bool
	var seenQuery bool
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "type":
			if value.Type != gjson.String {
				err = werror.ErrorWithContextParams(ctx, "field ParameterType[\"type\"] expected JSON string")
				return false
			}
			u.typ = value.Str
			seenType = true
		case "body":
			var unionVal BodyParameterType
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field ParameterType[\"body\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field ParameterType[\"body\"]")
					return false
				}
			}
			u.body = &unionVal
			seenBody = true
		case "header":
			var unionVal HeaderParameterType
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field ParameterType[\"header\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field ParameterType[\"header\"]")
					return false
				}
			}
			u.header = &unionVal
			seenHeader = true
		case "path":
			var unionVal PathParameterType
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field ParameterType[\"path\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field ParameterType[\"path\"]")
					return false
				}
			}
			u.path = &unionVal
			seenPath = true
		case "query":
			var unionVal QueryParameterType
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field ParameterType[\"query\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field ParameterType[\"query\"]")
					return false
				}
			}
			u.query = &unionVal
			seenQuery = true
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
	if u.typ == "body" && !seenBody {
		missingFields = append(missingFields, "body")
	}
	if u.typ == "header" && !seenHeader {
		missingFields = append(missingFields, "header")
	}
	if u.typ == "path" && !seenPath {
		missingFields = append(missingFields, "path")
	}
	if u.typ == "query" && !seenQuery {
		missingFields = append(missingFields, "query")
	}
	if len(missingFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type ParameterType missing required JSON fields", werror.SafeParam("missingFields", missingFields))
	}
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type ParameterType encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (u ParameterType) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(u)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (u *ParameterType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&u)
}

type Type struct {
	typ       string
	primitive *PrimitiveType
	optional  *OptionalType
	list      *ListType
	set       *SetType
	map_      *MapType
	reference *TypeName
	external  *ExternalReference
}

func NewTypeFromPrimitive(v PrimitiveType) Type {
	return Type{typ: "primitive", primitive: &v}
}

func NewTypeFromOptional(v OptionalType) Type {
	return Type{typ: "optional", optional: &v}
}

func NewTypeFromList(v ListType) Type {
	return Type{typ: "list", list: &v}
}

func NewTypeFromSet(v SetType) Type {
	return Type{typ: "set", set: &v}
}

func NewTypeFromMap(v MapType) Type {
	return Type{typ: "map", map_: &v}
}

func NewTypeFromReference(v TypeName) Type {
	return Type{typ: "reference", reference: &v}
}

func NewTypeFromExternal(v ExternalReference) Type {
	return Type{typ: "external", external: &v}
}

type TypeVisitor interface {
	VisitPrimitive(PrimitiveType) error
	VisitOptional(OptionalType) error
	VisitList(ListType) error
	VisitSet(SetType) error
	VisitMap(MapType) error
	VisitReference(TypeName) error
	VisitExternal(ExternalReference) error
	VisitUnknown(typeName string) error
}

func (u *Type) Accept(v TypeVisitor) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(u.typ)
	case "primitive":
		return v.VisitPrimitive(*u.primitive)
	case "optional":
		return v.VisitOptional(*u.optional)
	case "list":
		return v.VisitList(*u.list)
	case "set":
		return v.VisitSet(*u.set)
	case "map":
		return v.VisitMap(*u.map_)
	case "reference":
		return v.VisitReference(*u.reference)
	case "external":
		return v.VisitExternal(*u.external)
	}
}

type TypeVisitorWithContext interface {
	VisitPrimitiveWithContext(context.Context, PrimitiveType) error
	VisitOptionalWithContext(context.Context, OptionalType) error
	VisitListWithContext(context.Context, ListType) error
	VisitSetWithContext(context.Context, SetType) error
	VisitMapWithContext(context.Context, MapType) error
	VisitReferenceWithContext(context.Context, TypeName) error
	VisitExternalWithContext(context.Context, ExternalReference) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
}

func (u *Type) AcceptWithContext(ctx context.Context, v TypeVisitorWithContext) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknownWithContext(ctx, u.typ)
	case "primitive":
		return v.VisitPrimitiveWithContext(ctx, *u.primitive)
	case "optional":
		return v.VisitOptionalWithContext(ctx, *u.optional)
	case "list":
		return v.VisitListWithContext(ctx, *u.list)
	case "set":
		return v.VisitSetWithContext(ctx, *u.set)
	case "map":
		return v.VisitMapWithContext(ctx, *u.map_)
	case "reference":
		return v.VisitReferenceWithContext(ctx, *u.reference)
	case "external":
		return v.VisitExternalWithContext(ctx, *u.external)
	}
}

func (u *Type) AcceptFuncs(primitiveFunc func(PrimitiveType) error, optionalFunc func(OptionalType) error, listFunc func(ListType) error, setFunc func(SetType) error, map_Func func(MapType) error, referenceFunc func(TypeName) error, externalFunc func(ExternalReference) error, unknownFunc func(string) error) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return unknownFunc(u.typ)
	case "primitive":
		return primitiveFunc(*u.primitive)
	case "optional":
		return optionalFunc(*u.optional)
	case "list":
		return listFunc(*u.list)
	case "set":
		return setFunc(*u.set)
	case "map":
		return map_Func(*u.map_)
	case "reference":
		return referenceFunc(*u.reference)
	case "external":
		return externalFunc(*u.external)
	}
}

func (u *Type) PrimitiveNoopSuccess(PrimitiveType) error {
	return nil
}

func (u *Type) OptionalNoopSuccess(OptionalType) error {
	return nil
}

func (u *Type) ListNoopSuccess(ListType) error {
	return nil
}

func (u *Type) SetNoopSuccess(SetType) error {
	return nil
}

func (u *Type) MapNoopSuccess(MapType) error {
	return nil
}

func (u *Type) ReferenceNoopSuccess(TypeName) error {
	return nil
}

func (u *Type) ExternalNoopSuccess(ExternalReference) error {
	return nil
}

func (u *Type) ErrorOnUnknown(typeName string) error {
	return fmt.Errorf("invalid value in union type. Type name: %s", typeName)
}

func (u Type) MarshalJSON() ([]byte, error) {
	return u.AppendJSON(nil)
}

func (u Type) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	switch u.typ {
	default:
		out = append(out, "\"type\":"...)
		out = safejson.AppendQuotedString(out, u.typ)
	case "primitive":
		out = append(out, "\"type\":\"primitive\""...)
		if u.primitive != nil {
			out = append(out, ',')
			out = append(out, "\"primitive\""...)
			out = append(out, ':')
			unionVal := *u.primitive
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	case "optional":
		out = append(out, "\"type\":\"optional\""...)
		if u.optional != nil {
			out = append(out, ',')
			out = append(out, "\"optional\""...)
			out = append(out, ':')
			unionVal := *u.optional
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	case "list":
		out = append(out, "\"type\":\"list\""...)
		if u.list != nil {
			out = append(out, ',')
			out = append(out, "\"list\""...)
			out = append(out, ':')
			unionVal := *u.list
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
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	case "map":
		out = append(out, "\"type\":\"map\""...)
		if u.map_ != nil {
			out = append(out, ',')
			out = append(out, "\"map\""...)
			out = append(out, ':')
			unionVal := *u.map_
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	case "reference":
		out = append(out, "\"type\":\"reference\""...)
		if u.reference != nil {
			out = append(out, ',')
			out = append(out, "\"reference\""...)
			out = append(out, ':')
			unionVal := *u.reference
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	case "external":
		out = append(out, "\"type\":\"external\""...)
		if u.external != nil {
			out = append(out, ',')
			out = append(out, "\"external\""...)
			out = append(out, ':')
			unionVal := *u.external
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	}
	out = append(out, '}')
	return out, nil
}

func (u *Type) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Type")
	}
	return u.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (u *Type) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Type")
	}
	return u.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (u *Type) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Type")
	}
	return u.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (u *Type) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for Type")
	}
	return u.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (u *Type) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type Type expected JSON object")
	}
	var seenType bool
	var seenPrimitive bool
	var seenOptional bool
	var seenList bool
	var seenSet bool
	var seenMap bool
	var seenReference bool
	var seenExternal bool
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "type":
			if value.Type != gjson.String {
				err = werror.ErrorWithContextParams(ctx, "field Type[\"type\"] expected JSON string")
				return false
			}
			u.typ = value.Str
			seenType = true
		case "primitive":
			var unionVal PrimitiveType
			if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
				err = werror.WrapWithContextParams(ctx, err, "field Type[\"primitive\"]")
				return false
			}
			u.primitive = &unionVal
			seenPrimitive = true
		case "optional":
			var unionVal OptionalType
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Type[\"optional\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Type[\"optional\"]")
					return false
				}
			}
			u.optional = &unionVal
			seenOptional = true
		case "list":
			var unionVal ListType
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Type[\"list\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Type[\"list\"]")
					return false
				}
			}
			u.list = &unionVal
			seenList = true
		case "set":
			var unionVal SetType
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Type[\"set\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Type[\"set\"]")
					return false
				}
			}
			u.set = &unionVal
			seenSet = true
		case "map":
			var unionVal MapType
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Type[\"map\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Type[\"map\"]")
					return false
				}
			}
			u.map_ = &unionVal
			seenMap = true
		case "reference":
			var unionVal TypeName
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Type[\"reference\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Type[\"reference\"]")
					return false
				}
			}
			u.reference = &unionVal
			seenReference = true
		case "external":
			var unionVal ExternalReference
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Type[\"external\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field Type[\"external\"]")
					return false
				}
			}
			u.external = &unionVal
			seenExternal = true
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
	if u.typ == "primitive" && !seenPrimitive {
		missingFields = append(missingFields, "primitive")
	}
	if u.typ == "optional" && !seenOptional {
		missingFields = append(missingFields, "optional")
	}
	if u.typ == "list" && !seenList {
		missingFields = append(missingFields, "list")
	}
	if u.typ == "set" && !seenSet {
		missingFields = append(missingFields, "set")
	}
	if u.typ == "map" && !seenMap {
		missingFields = append(missingFields, "map")
	}
	if u.typ == "reference" && !seenReference {
		missingFields = append(missingFields, "reference")
	}
	if u.typ == "external" && !seenExternal {
		missingFields = append(missingFields, "external")
	}
	if len(missingFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type Type missing required JSON fields", werror.SafeParam("missingFields", missingFields))
	}
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type Type encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (u Type) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(u)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (u *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&u)
}

type TypeDefinition struct {
	typ    string
	alias  *AliasDefinition
	enum   *EnumDefinition
	object *ObjectDefinition
	union  *UnionDefinition
}

func NewTypeDefinitionFromAlias(v AliasDefinition) TypeDefinition {
	return TypeDefinition{typ: "alias", alias: &v}
}

func NewTypeDefinitionFromEnum(v EnumDefinition) TypeDefinition {
	return TypeDefinition{typ: "enum", enum: &v}
}

func NewTypeDefinitionFromObject(v ObjectDefinition) TypeDefinition {
	return TypeDefinition{typ: "object", object: &v}
}

func NewTypeDefinitionFromUnion(v UnionDefinition) TypeDefinition {
	return TypeDefinition{typ: "union", union: &v}
}

type TypeDefinitionVisitor interface {
	VisitAlias(AliasDefinition) error
	VisitEnum(EnumDefinition) error
	VisitObject(ObjectDefinition) error
	VisitUnion(UnionDefinition) error
	VisitUnknown(typeName string) error
}

func (u *TypeDefinition) Accept(v TypeDefinitionVisitor) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(u.typ)
	case "alias":
		return v.VisitAlias(*u.alias)
	case "enum":
		return v.VisitEnum(*u.enum)
	case "object":
		return v.VisitObject(*u.object)
	case "union":
		return v.VisitUnion(*u.union)
	}
}

type TypeDefinitionVisitorWithContext interface {
	VisitAliasWithContext(context.Context, AliasDefinition) error
	VisitEnumWithContext(context.Context, EnumDefinition) error
	VisitObjectWithContext(context.Context, ObjectDefinition) error
	VisitUnionWithContext(context.Context, UnionDefinition) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
}

func (u *TypeDefinition) AcceptWithContext(ctx context.Context, v TypeDefinitionVisitorWithContext) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknownWithContext(ctx, u.typ)
	case "alias":
		return v.VisitAliasWithContext(ctx, *u.alias)
	case "enum":
		return v.VisitEnumWithContext(ctx, *u.enum)
	case "object":
		return v.VisitObjectWithContext(ctx, *u.object)
	case "union":
		return v.VisitUnionWithContext(ctx, *u.union)
	}
}

func (u *TypeDefinition) AcceptFuncs(aliasFunc func(AliasDefinition) error, enumFunc func(EnumDefinition) error, objectFunc func(ObjectDefinition) error, unionFunc func(UnionDefinition) error, unknownFunc func(string) error) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return unknownFunc(u.typ)
	case "alias":
		return aliasFunc(*u.alias)
	case "enum":
		return enumFunc(*u.enum)
	case "object":
		return objectFunc(*u.object)
	case "union":
		return unionFunc(*u.union)
	}
}

func (u *TypeDefinition) AliasNoopSuccess(AliasDefinition) error {
	return nil
}

func (u *TypeDefinition) EnumNoopSuccess(EnumDefinition) error {
	return nil
}

func (u *TypeDefinition) ObjectNoopSuccess(ObjectDefinition) error {
	return nil
}

func (u *TypeDefinition) UnionNoopSuccess(UnionDefinition) error {
	return nil
}

func (u *TypeDefinition) ErrorOnUnknown(typeName string) error {
	return fmt.Errorf("invalid value in union type. Type name: %s", typeName)
}

func (u TypeDefinition) MarshalJSON() ([]byte, error) {
	return u.AppendJSON(nil)
}

func (u TypeDefinition) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	switch u.typ {
	default:
		out = append(out, "\"type\":"...)
		out = safejson.AppendQuotedString(out, u.typ)
	case "alias":
		out = append(out, "\"type\":\"alias\""...)
		if u.alias != nil {
			out = append(out, ',')
			out = append(out, "\"alias\""...)
			out = append(out, ':')
			unionVal := *u.alias
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	case "enum":
		out = append(out, "\"type\":\"enum\""...)
		if u.enum != nil {
			out = append(out, ',')
			out = append(out, "\"enum\""...)
			out = append(out, ':')
			unionVal := *u.enum
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	case "object":
		out = append(out, "\"type\":\"object\""...)
		if u.object != nil {
			out = append(out, ',')
			out = append(out, "\"object\""...)
			out = append(out, ':')
			unionVal := *u.object
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	case "union":
		out = append(out, "\"type\":\"union\""...)
		if u.union != nil {
			out = append(out, ',')
			out = append(out, "\"union\""...)
			out = append(out, ':')
			unionVal := *u.union
			var err error
			out, err = unionVal.AppendJSON(out)
			if err != nil {
				return nil, err
			}
		}
	}
	out = append(out, '}')
	return out, nil
}

func (u *TypeDefinition) UnmarshalJSON(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for TypeDefinition")
	}
	return u.unmarshalJSONResult(ctx, gjson.ParseBytes(data), false)
}

func (u *TypeDefinition) UnmarshalJSONStrict(data []byte) error {
	ctx := context.TODO()
	if !gjson.ValidBytes(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for TypeDefinition")
	}
	return u.unmarshalJSONResult(ctx, gjson.ParseBytes(data), true)
}

func (u *TypeDefinition) UnmarshalJSONString(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for TypeDefinition")
	}
	return u.unmarshalJSONResult(ctx, gjson.Parse(data), false)
}

func (u *TypeDefinition) UnmarshalJSONStringStrict(data string) error {
	ctx := context.TODO()
	if !gjson.Valid(data) {
		return werror.ErrorWithContextParams(ctx, "invalid JSON for TypeDefinition")
	}
	return u.unmarshalJSONResult(ctx, gjson.Parse(data), true)
}

func (u *TypeDefinition) unmarshalJSONResult(ctx context.Context, value gjson.Result, strict bool) error {
	if !value.IsObject() {
		return werror.ErrorWithContextParams(ctx, "type TypeDefinition expected JSON object")
	}
	var seenType bool
	var seenAlias bool
	var seenEnum bool
	var seenObject bool
	var seenUnion bool
	var unrecognizedFields []string
	var err error
	value.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "type":
			if value.Type != gjson.String {
				err = werror.ErrorWithContextParams(ctx, "field TypeDefinition[\"type\"] expected JSON string")
				return false
			}
			u.typ = value.Str
			seenType = true
		case "alias":
			var unionVal AliasDefinition
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field TypeDefinition[\"alias\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field TypeDefinition[\"alias\"]")
					return false
				}
			}
			u.alias = &unionVal
			seenAlias = true
		case "enum":
			var unionVal EnumDefinition
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field TypeDefinition[\"enum\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field TypeDefinition[\"enum\"]")
					return false
				}
			}
			u.enum = &unionVal
			seenEnum = true
		case "object":
			var unionVal ObjectDefinition
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field TypeDefinition[\"object\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field TypeDefinition[\"object\"]")
					return false
				}
			}
			u.object = &unionVal
			seenObject = true
		case "union":
			var unionVal UnionDefinition
			if strict {
				if err = unionVal.UnmarshalJSONStringStrict(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field TypeDefinition[\"union\"]")
					return false
				}
			} else {
				if err = unionVal.UnmarshalJSONString(value.Raw); err != nil {
					err = werror.WrapWithContextParams(ctx, err, "field TypeDefinition[\"union\"]")
					return false
				}
			}
			u.union = &unionVal
			seenUnion = true
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
	if u.typ == "alias" && !seenAlias {
		missingFields = append(missingFields, "alias")
	}
	if u.typ == "enum" && !seenEnum {
		missingFields = append(missingFields, "enum")
	}
	if u.typ == "object" && !seenObject {
		missingFields = append(missingFields, "object")
	}
	if u.typ == "union" && !seenUnion {
		missingFields = append(missingFields, "union")
	}
	if len(missingFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type TypeDefinition missing required JSON fields", werror.SafeParam("missingFields", missingFields))
	}
	if strict && len(unrecognizedFields) > 0 {
		return werror.ErrorWithContextParams(ctx, "type TypeDefinition encountered unrecognized JSON fields", werror.UnsafeParam("unrecognizedFields", unrecognizedFields))
	}
	return nil
}

func (u TypeDefinition) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(u)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (u *TypeDefinition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&u)
}
