// This file was generated by Conjure and should not be manually edited.

//go:build go1.18

package foo1

import (
	"context"
	"fmt"

	"github.com/palantir/conjure-go/v6/conjure/cycles/testdata/pkg-cycle-disconnected/conjure/com/palantir/bar"
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
	case "field3":
		if u.field3 == nil {
			return result, fmt.Errorf("field \"field3\" is required")
		}
		return v.VisitField3(ctx, *u.field3)
	}
}

type Type3VisitorWithT[T any] interface {
	VisitField3(ctx context.Context, v bar.Type1) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}
