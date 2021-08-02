// This file was generated by Conjure and should not be manually edited.

package spec

import (
	"context"
	"fmt"

	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type AuthType struct {
	typ    string
	header *HeaderAuthType
	cookie *CookieAuthType
}

type authTypeDeserializer struct {
	Type   string          `json:"type"`
	Header *HeaderAuthType `json:"header"`
	Cookie *CookieAuthType `json:"cookie"`
}

func (u *authTypeDeserializer) toStruct() AuthType {
	return AuthType{typ: u.Type, header: u.Header, cookie: u.Cookie}
}

func (u *AuthType) toSerializer() (interface{}, error) {
	switch u.typ {
	default:
		return nil, fmt.Errorf("unknown type %s", u.typ)
	case "header":
		return struct {
			Type   string         `json:"type"`
			Header HeaderAuthType `json:"header"`
		}{Type: "header", Header: *u.header}, nil
	case "cookie":
		return struct {
			Type   string         `json:"type"`
			Cookie CookieAuthType `json:"cookie"`
		}{Type: "cookie", Cookie: *u.cookie}, nil
	}
}

func (u AuthType) MarshalJSON() ([]byte, error) {
	ser, err := u.toSerializer()
	if err != nil {
		return nil, err
	}
	return safejson.Marshal(ser)
}

func (u *AuthType) UnmarshalJSON(data []byte) error {
	var deser authTypeDeserializer
	if err := safejson.Unmarshal(data, &deser); err != nil {
		return err
	}
	*u = deser.toStruct()
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

type AuthTypeVisitor interface {
	VisitHeader(v HeaderAuthType) error
	VisitCookie(v CookieAuthType) error
	VisitUnknown(typeName string) error
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

type AuthTypeVisitorWithContext interface {
	VisitHeaderWithContext(ctx context.Context, v HeaderAuthType) error
	VisitCookieWithContext(ctx context.Context, v CookieAuthType) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
}

func NewAuthTypeFromHeader(v HeaderAuthType) AuthType {
	return AuthType{typ: "header", header: &v}
}

func NewAuthTypeFromCookie(v CookieAuthType) AuthType {
	return AuthType{typ: "cookie", cookie: &v}
}

type ParameterType struct {
	typ    string
	body   *BodyParameterType
	header *HeaderParameterType
	path   *PathParameterType
	query  *QueryParameterType
}

type parameterTypeDeserializer struct {
	Type   string               `json:"type"`
	Body   *BodyParameterType   `json:"body"`
	Header *HeaderParameterType `json:"header"`
	Path   *PathParameterType   `json:"path"`
	Query  *QueryParameterType  `json:"query"`
}

func (u *parameterTypeDeserializer) toStruct() ParameterType {
	return ParameterType{typ: u.Type, body: u.Body, header: u.Header, path: u.Path, query: u.Query}
}

func (u *ParameterType) toSerializer() (interface{}, error) {
	switch u.typ {
	default:
		return nil, fmt.Errorf("unknown type %s", u.typ)
	case "body":
		return struct {
			Type string            `json:"type"`
			Body BodyParameterType `json:"body"`
		}{Type: "body", Body: *u.body}, nil
	case "header":
		return struct {
			Type   string              `json:"type"`
			Header HeaderParameterType `json:"header"`
		}{Type: "header", Header: *u.header}, nil
	case "path":
		return struct {
			Type string            `json:"type"`
			Path PathParameterType `json:"path"`
		}{Type: "path", Path: *u.path}, nil
	case "query":
		return struct {
			Type  string             `json:"type"`
			Query QueryParameterType `json:"query"`
		}{Type: "query", Query: *u.query}, nil
	}
}

func (u ParameterType) MarshalJSON() ([]byte, error) {
	ser, err := u.toSerializer()
	if err != nil {
		return nil, err
	}
	return safejson.Marshal(ser)
}

func (u *ParameterType) UnmarshalJSON(data []byte) error {
	var deser parameterTypeDeserializer
	if err := safejson.Unmarshal(data, &deser); err != nil {
		return err
	}
	*u = deser.toStruct()
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

type ParameterTypeVisitor interface {
	VisitBody(v BodyParameterType) error
	VisitHeader(v HeaderParameterType) error
	VisitPath(v PathParameterType) error
	VisitQuery(v QueryParameterType) error
	VisitUnknown(typeName string) error
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

type ParameterTypeVisitorWithContext interface {
	VisitBodyWithContext(ctx context.Context, v BodyParameterType) error
	VisitHeaderWithContext(ctx context.Context, v HeaderParameterType) error
	VisitPathWithContext(ctx context.Context, v PathParameterType) error
	VisitQueryWithContext(ctx context.Context, v QueryParameterType) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
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

type typeDeserializer struct {
	Type      string             `json:"type"`
	Primitive *PrimitiveType     `json:"primitive"`
	Optional  *OptionalType      `json:"optional"`
	List      *ListType          `json:"list"`
	Set       *SetType           `json:"set"`
	Map       *MapType           `json:"map"`
	Reference *TypeName          `json:"reference"`
	External  *ExternalReference `json:"external"`
}

func (u *typeDeserializer) toStruct() Type {
	return Type{typ: u.Type, primitive: u.Primitive, optional: u.Optional, list: u.List, set: u.Set, map_: u.Map, reference: u.Reference, external: u.External}
}

func (u *Type) toSerializer() (interface{}, error) {
	switch u.typ {
	default:
		return nil, fmt.Errorf("unknown type %s", u.typ)
	case "primitive":
		return struct {
			Type      string        `json:"type"`
			Primitive PrimitiveType `json:"primitive"`
		}{Type: "primitive", Primitive: *u.primitive}, nil
	case "optional":
		return struct {
			Type     string       `json:"type"`
			Optional OptionalType `json:"optional"`
		}{Type: "optional", Optional: *u.optional}, nil
	case "list":
		return struct {
			Type string   `json:"type"`
			List ListType `json:"list"`
		}{Type: "list", List: *u.list}, nil
	case "set":
		return struct {
			Type string  `json:"type"`
			Set  SetType `json:"set"`
		}{Type: "set", Set: *u.set}, nil
	case "map":
		return struct {
			Type string  `json:"type"`
			Map  MapType `json:"map"`
		}{Type: "map", Map: *u.map_}, nil
	case "reference":
		return struct {
			Type      string   `json:"type"`
			Reference TypeName `json:"reference"`
		}{Type: "reference", Reference: *u.reference}, nil
	case "external":
		return struct {
			Type     string            `json:"type"`
			External ExternalReference `json:"external"`
		}{Type: "external", External: *u.external}, nil
	}
}

func (u Type) MarshalJSON() ([]byte, error) {
	ser, err := u.toSerializer()
	if err != nil {
		return nil, err
	}
	return safejson.Marshal(ser)
}

func (u *Type) UnmarshalJSON(data []byte) error {
	var deser typeDeserializer
	if err := safejson.Unmarshal(data, &deser); err != nil {
		return err
	}
	*u = deser.toStruct()
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

func (u *Type) AcceptFuncs(primitiveFunc func(PrimitiveType) error, optionalFunc func(OptionalType) error, listFunc func(ListType) error, setFunc func(SetType) error, mapFunc func(MapType) error, referenceFunc func(TypeName) error, externalFunc func(ExternalReference) error, unknownFunc func(string) error) error {
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
		return mapFunc(*u.map_)
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

type TypeVisitor interface {
	VisitPrimitive(v PrimitiveType) error
	VisitOptional(v OptionalType) error
	VisitList(v ListType) error
	VisitSet(v SetType) error
	VisitMap(v MapType) error
	VisitReference(v TypeName) error
	VisitExternal(v ExternalReference) error
	VisitUnknown(typeName string) error
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

type TypeVisitorWithContext interface {
	VisitPrimitiveWithContext(ctx context.Context, v PrimitiveType) error
	VisitOptionalWithContext(ctx context.Context, v OptionalType) error
	VisitListWithContext(ctx context.Context, v ListType) error
	VisitSetWithContext(ctx context.Context, v SetType) error
	VisitMapWithContext(ctx context.Context, v MapType) error
	VisitReferenceWithContext(ctx context.Context, v TypeName) error
	VisitExternalWithContext(ctx context.Context, v ExternalReference) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
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

type TypeDefinition struct {
	typ    string
	alias  *AliasDefinition
	enum   *EnumDefinition
	object *ObjectDefinition
	union  *UnionDefinition
}

type typeDefinitionDeserializer struct {
	Type   string            `json:"type"`
	Alias  *AliasDefinition  `json:"alias"`
	Enum   *EnumDefinition   `json:"enum"`
	Object *ObjectDefinition `json:"object"`
	Union  *UnionDefinition  `json:"union"`
}

func (u *typeDefinitionDeserializer) toStruct() TypeDefinition {
	return TypeDefinition{typ: u.Type, alias: u.Alias, enum: u.Enum, object: u.Object, union: u.Union}
}

func (u *TypeDefinition) toSerializer() (interface{}, error) {
	switch u.typ {
	default:
		return nil, fmt.Errorf("unknown type %s", u.typ)
	case "alias":
		return struct {
			Type  string          `json:"type"`
			Alias AliasDefinition `json:"alias"`
		}{Type: "alias", Alias: *u.alias}, nil
	case "enum":
		return struct {
			Type string         `json:"type"`
			Enum EnumDefinition `json:"enum"`
		}{Type: "enum", Enum: *u.enum}, nil
	case "object":
		return struct {
			Type   string           `json:"type"`
			Object ObjectDefinition `json:"object"`
		}{Type: "object", Object: *u.object}, nil
	case "union":
		return struct {
			Type  string          `json:"type"`
			Union UnionDefinition `json:"union"`
		}{Type: "union", Union: *u.union}, nil
	}
}

func (u TypeDefinition) MarshalJSON() ([]byte, error) {
	ser, err := u.toSerializer()
	if err != nil {
		return nil, err
	}
	return safejson.Marshal(ser)
}

func (u *TypeDefinition) UnmarshalJSON(data []byte) error {
	var deser typeDefinitionDeserializer
	if err := safejson.Unmarshal(data, &deser); err != nil {
		return err
	}
	*u = deser.toStruct()
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

type TypeDefinitionVisitor interface {
	VisitAlias(v AliasDefinition) error
	VisitEnum(v EnumDefinition) error
	VisitObject(v ObjectDefinition) error
	VisitUnion(v UnionDefinition) error
	VisitUnknown(typeName string) error
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

type TypeDefinitionVisitorWithContext interface {
	VisitAliasWithContext(ctx context.Context, v AliasDefinition) error
	VisitEnumWithContext(ctx context.Context, v EnumDefinition) error
	VisitObjectWithContext(ctx context.Context, v ObjectDefinition) error
	VisitUnionWithContext(ctx context.Context, v UnionDefinition) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
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
