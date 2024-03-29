// This file was generated by Conjure and should not be manually edited.

//go:build go1.18

package api

import (
	"context"
	"fmt"
)

type ExampleUnionWithT[T any] ExampleUnion

func (u *ExampleUnionWithT[T]) Accept(ctx context.Context, v ExampleUnionVisitorWithT[T]) (T, error) {
	var result T
	switch u.typ {
	default:
		if u.typ == "" {
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "str":
		if u.str == nil {
			return result, fmt.Errorf("field \"str\" is required")
		}
		return v.VisitStr(ctx, *u.str)
	case "strOptional":
		var strOptional *string
		if u.strOptional != nil {
			strOptional = *u.strOptional
		}
		return v.VisitStrOptional(ctx, strOptional)
	case "other":
		if u.other == nil {
			return result, fmt.Errorf("field \"other\" is required")
		}
		return v.VisitOther(ctx, *u.other)
	}
}

type ExampleUnionVisitorWithT[T any] interface {
	VisitStr(ctx context.Context, v string) (T, error)
	VisitStrOptional(ctx context.Context, v *string) (T, error)
	VisitOther(ctx context.Context, v int) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}
