// This file was generated by Conjure and should not be manually edited.

//go:build go1.18

package api

import (
	"context"
	"fmt"
)

type ExampleUnionWithT[T any] ExampleUnion

func (u *ExampleUnionWithT[T]) Accept(ctx context.Context, v ExampleUnionVisitorWithT[T]) (T, error) {
	switch u.typ {
	default:
		if u.typ == "" {
			var result T
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "str":
		return v.VisitStr(ctx, *u.str)
	case "strOptional":
		var strOptional *string
		if u.strOptional != nil {
			strOptional = *u.strOptional
		}
		return v.VisitStrOptional(ctx, strOptional)
	case "other":
		return v.VisitOther(ctx, *u.other)
	}
}

type ExampleUnionVisitorWithT[T any] interface {
	VisitStr(ctx context.Context, v string) (T, error)
	VisitStrOptional(ctx context.Context, v *string) (T, error)
	VisitOther(ctx context.Context, v int) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}