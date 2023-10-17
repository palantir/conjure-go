// This file was generated by Conjure and should not be manually edited.

//go:build go1.18

package spec

import (
	"context"
	"fmt"

	werror "github.com/palantir/witchcraft-go-error"
)

type AuthTypeWithT[T any] AuthType

func (u *AuthTypeWithT[T]) Accept(ctx context.Context, v AuthTypeVisitorWithT[T]) (T, error) {
	var result T
	switch u.typ {
	default:
		if u.typ == "" {
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "header":
		if u.header == nil {
			return result, werror.Error("field header is required")
		}
		return v.VisitHeader(ctx, *u.header)
	case "cookie":
		if u.cookie == nil {
			return result, werror.Error("field cookie is required")
		}
		return v.VisitCookie(ctx, *u.cookie)
	}
}

type AuthTypeVisitorWithT[T any] interface {
	VisitHeader(ctx context.Context, v HeaderAuthType) (T, error)
	VisitCookie(ctx context.Context, v CookieAuthType) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}

type ParameterTypeWithT[T any] ParameterType

func (u *ParameterTypeWithT[T]) Accept(ctx context.Context, v ParameterTypeVisitorWithT[T]) (T, error) {
	var result T
	switch u.typ {
	default:
		if u.typ == "" {
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "body":
		if u.body == nil {
			return result, werror.Error("field body is required")
		}
		return v.VisitBody(ctx, *u.body)
	case "header":
		if u.header == nil {
			return result, werror.Error("field header is required")
		}
		return v.VisitHeader(ctx, *u.header)
	case "path":
		if u.path == nil {
			return result, werror.Error("field path is required")
		}
		return v.VisitPath(ctx, *u.path)
	case "query":
		if u.query == nil {
			return result, werror.Error("field query is required")
		}
		return v.VisitQuery(ctx, *u.query)
	}
}

type ParameterTypeVisitorWithT[T any] interface {
	VisitBody(ctx context.Context, v BodyParameterType) (T, error)
	VisitHeader(ctx context.Context, v HeaderParameterType) (T, error)
	VisitPath(ctx context.Context, v PathParameterType) (T, error)
	VisitQuery(ctx context.Context, v QueryParameterType) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}

type TypeWithT[T any] Type

func (u *TypeWithT[T]) Accept(ctx context.Context, v TypeVisitorWithT[T]) (T, error) {
	var result T
	switch u.typ {
	default:
		if u.typ == "" {
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "primitive":
		if u.primitive == nil {
			return result, werror.Error("field primitive is required")
		}
		return v.VisitPrimitive(ctx, *u.primitive)
	case "optional":
		if u.optional == nil {
			return result, werror.Error("field optional is required")
		}
		return v.VisitOptional(ctx, *u.optional)
	case "list":
		if u.list == nil {
			return result, werror.Error("field list is required")
		}
		return v.VisitList(ctx, *u.list)
	case "set":
		if u.set == nil {
			return result, werror.Error("field set is required")
		}
		return v.VisitSet(ctx, *u.set)
	case "map":
		if u.map_ == nil {
			return result, werror.Error("field map is required")
		}
		return v.VisitMap(ctx, *u.map_)
	case "reference":
		if u.reference == nil {
			return result, werror.Error("field reference is required")
		}
		return v.VisitReference(ctx, *u.reference)
	case "external":
		if u.external == nil {
			return result, werror.Error("field external is required")
		}
		return v.VisitExternal(ctx, *u.external)
	}
}

type TypeVisitorWithT[T any] interface {
	VisitPrimitive(ctx context.Context, v PrimitiveType) (T, error)
	VisitOptional(ctx context.Context, v OptionalType) (T, error)
	VisitList(ctx context.Context, v ListType) (T, error)
	VisitSet(ctx context.Context, v SetType) (T, error)
	VisitMap(ctx context.Context, v MapType) (T, error)
	VisitReference(ctx context.Context, v TypeName) (T, error)
	VisitExternal(ctx context.Context, v ExternalReference) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}

type TypeDefinitionWithT[T any] TypeDefinition

func (u *TypeDefinitionWithT[T]) Accept(ctx context.Context, v TypeDefinitionVisitorWithT[T]) (T, error) {
	var result T
	switch u.typ {
	default:
		if u.typ == "" {
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "alias":
		if u.alias == nil {
			return result, werror.Error("field alias is required")
		}
		return v.VisitAlias(ctx, *u.alias)
	case "enum":
		if u.enum == nil {
			return result, werror.Error("field enum is required")
		}
		return v.VisitEnum(ctx, *u.enum)
	case "object":
		if u.object == nil {
			return result, werror.Error("field object is required")
		}
		return v.VisitObject(ctx, *u.object)
	case "union":
		if u.union == nil {
			return result, werror.Error("field union is required")
		}
		return v.VisitUnion(ctx, *u.union)
	}
}

type TypeDefinitionVisitorWithT[T any] interface {
	VisitAlias(ctx context.Context, v AliasDefinition) (T, error)
	VisitEnum(ctx context.Context, v EnumDefinition) (T, error)
	VisitObject(ctx context.Context, v ObjectDefinition) (T, error)
	VisitUnion(ctx context.Context, v UnionDefinition) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}
