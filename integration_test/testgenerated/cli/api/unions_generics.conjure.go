// This file was generated by Conjure and should not be manually edited.

//go:build go1.18

package api

import (
	"context"
	"fmt"
)

type CustomUnionWithT[T any] CustomUnion

func (u *CustomUnionWithT[T]) Accept(ctx context.Context, v CustomUnionVisitorWithT[T]) (T, error) {
	var result T
	switch u.typ {
	default:
		if u.typ == "" {
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "asString":
		if u.asString == nil {
			return result, fmt.Errorf("field asString is required")
		}
		return v.VisitAsString(ctx, *u.asString)
	case "asInteger":
		if u.asInteger == nil {
			return result, fmt.Errorf("field asInteger is required")
		}
		return v.VisitAsInteger(ctx, *u.asInteger)
	}
}

type CustomUnionVisitorWithT[T any] interface {
	VisitAsString(ctx context.Context, v string) (T, error)
	VisitAsInteger(ctx context.Context, v int) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}
