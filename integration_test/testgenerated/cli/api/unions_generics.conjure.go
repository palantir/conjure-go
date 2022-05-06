// This file was generated by Conjure and should not be manually edited.

//go:build go1.18

package api

import (
	"context"
	"fmt"
)

type CustomUnionWithT[T any] CustomUnion

func (u *CustomUnionWithT[T]) Accept(ctx context.Context, v CustomUnionVisitorWithT[T]) (T, error) {
	switch u.typ {
	default:
		if u.typ == "" {
			var result T
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "asString":
		return v.VisitAsString(ctx, *u.asString)
	case "asInteger":
		return v.VisitAsInteger(ctx, *u.asInteger)
	}
}

type CustomUnionVisitorWithT[T any] interface {
	VisitAsString(ctx context.Context, v string) (T, error)
	VisitAsInteger(ctx context.Context, v int) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}
