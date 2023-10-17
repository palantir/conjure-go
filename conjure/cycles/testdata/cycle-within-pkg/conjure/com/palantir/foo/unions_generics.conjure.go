// This file was generated by Conjure and should not be manually edited.

//go:build go1.18

package foo

import (
	"context"
	"fmt"

	"github.com/palantir/conjure-go/v6/conjure/cycles/testdata/cycle-within-pkg/conjure/com/palantir/bar"
	werror "github.com/palantir/witchcraft-go-error"
)

type Type3WithT[T any] Type3

func (u *Type3WithT[T]) Accept(ctx context.Context, v Type3VisitorWithT[T]) (T, error) {
	var result T
	switch u.typ {
	default:
		if u.typ == "" {
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "field1":
		if u.field1 == nil {
			return result, werror.Error("field field1 is required")
		}
		return v.VisitField1(ctx, *u.field1)
	case "field2":
		if u.field2 == nil {
			return result, werror.Error("field field2 is required")
		}
		return v.VisitField2(ctx, *u.field2)
	case "field3":
		if u.field3 == nil {
			return result, werror.Error("field field3 is required")
		}
		return v.VisitField3(ctx, *u.field3)
	}
}

type Type3VisitorWithT[T any] interface {
	VisitField1(ctx context.Context, v Type2) (T, error)
	VisitField2(ctx context.Context, v Type4) (T, error)
	VisitField3(ctx context.Context, v bar.Type3) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}
