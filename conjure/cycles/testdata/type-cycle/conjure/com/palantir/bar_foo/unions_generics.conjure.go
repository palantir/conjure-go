// This file was generated by Conjure and should not be manually edited.

//go:build go1.18

package barfoo

import (
	"context"
	"fmt"
)

type FooType3WithT[T any] FooType3

func (u *FooType3WithT[T]) Accept(ctx context.Context, v FooType3VisitorWithT[T]) (T, error) {
	var result T
	switch u.typ {
	default:
		if u.typ == "" {
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "field1":
		if u.field1 == nil {
			return result, fmt.Errorf("field \"field1\" is required")
		}
		return v.VisitField1(ctx, *u.field1)
	case "field3":
		if u.field3 == nil {
			return result, fmt.Errorf("field \"field3\" is required")
		}
		return v.VisitField3(ctx, *u.field3)
	}
}

type FooType3VisitorWithT[T any] interface {
	VisitField1(ctx context.Context, v Type2) (T, error)
	VisitField3(ctx context.Context, v Type1) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}
