// Copyright (c) 2022 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conjure

import (
	"bytes"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/stretchr/testify/assert"
)

var testUnionType = &types.UnionType{
	Name: "MyUnion",
	Fields: []*types.Field{
		{
			Name: "stringVal", Type: types.String{},
		},
		{
			Name: "boolVal", Type: types.Boolean{},
		},
	},
}

func TestUnionWriter_unionVisitorInterfaceT(t *testing.T) {
	f := jen.NewFile("testpkg")
	unionVisitorWithT(f.Group, testUnionType)
	var buf bytes.Buffer
	assert.NoError(t, f.Render(&buf))
	assert.Equal(t, `package testpkg

import "context"

type MyUnionVisitorWithT[T any] interface {
	VisitStringVal(ctx context.Context, v string) (T, error)
	VisitBoolVal(ctx context.Context, v bool) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}
`, buf.String())
}

func TestUnionWriter_unionTypeWithT(t *testing.T) {
	f := jen.NewFile("testpkg")
	unionTypeWithT(f.Group, testUnionType)
	var buf bytes.Buffer
	assert.NoError(t, f.Render(&buf))
	assert.Equal(t, `package testpkg

type MyUnionWithT[T any] MyUnion
`, buf.String())
}

func TestUnionWriter_unionTypeWithTAccept(t *testing.T) {
	f := jen.NewFile("testpkg")
	unionTypeWithTAccept(f.Group, testUnionType)
	var buf bytes.Buffer
	assert.NoError(t, f.Render(&buf))
	assert.Equal(t, `package testpkg

import (
	"context"
	"fmt"
)

func (u *MyUnionWithT[T]) Accept(ctx context.Context, v MyUnionVisitorWithT[T]) (T, error) {
	var result T
	switch u.typ {
	default:
		if u.typ == "" {
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "stringVal":
		if u.stringVal == nil {
			return result, fmt.Errorf("field stringVal is required")
		}
		return v.VisitStringVal(ctx, *u.stringVal)
	case "boolVal":
		if u.boolVal == nil {
			return result, fmt.Errorf("field boolVal is required")
		}
		return v.VisitBoolVal(ctx, *u.boolVal)
	}
}
`, buf.String())
}
