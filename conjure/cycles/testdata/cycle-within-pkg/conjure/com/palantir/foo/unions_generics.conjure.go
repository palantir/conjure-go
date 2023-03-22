// This file was generated by Conjure and should not be manually edited.

//go:build go1.18

package foo

import (
	"context"
	"fmt"

	"github.com/palantir/conjure-go/v6/conjure/cycles/testdata/cycle-within-pkg/conjure/com/palantir/bar"
)

type Type3WithT[T any] Type3

func (u *Type3WithT[T]) Accept(ctx context.Context, v Type3VisitorWithT[T]) (T, error) {
	switch u.typ {
	default:
		if u.typ == "" {
			var result T
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "field1":
		return v.VisitField1(ctx, *u.field1)
	case "field2":
		return v.VisitField2(ctx, *u.field2)
	case "field3":
		return v.VisitField3(ctx, *u.field3)
	}
}

type Type3VisitorWithT[T any] interface {
	VisitField1(ctx context.Context, v Type2) (T, error)
	VisitField2(ctx context.Context, v Type4) (T, error)
	VisitField3(ctx context.Context, v bar.Type3) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}