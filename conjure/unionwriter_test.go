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
	"github.com/palantir/conjure-go/v7/conjure/types"
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

var testUnionTypeWithOptional = &types.UnionType{
	Name: "MyUnionWithOptional",
	Fields: []*types.Field{
		{
			Name: "stringVal", Type: types.String{},
		},
		{
			Name: "optionalVal", Type: &types.Optional{Item: types.String{}},
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
			return result, fmt.Errorf("field \"stringVal\" is required")
		}
		return v.VisitStringVal(ctx, *u.stringVal)
	case "boolVal":
		if u.boolVal == nil {
			return result, fmt.Errorf("field \"boolVal\" is required")
		}
		return v.VisitBoolVal(ctx, *u.boolVal)
	}
}
`, buf.String())
}

func TestUnionWriter_unionTypeWithTAcceptFuncs(t *testing.T) {
	f := jen.NewFile("testpkg")
	unionTypeWithTAcceptFuncs(f.Group, testUnionType)
	var buf bytes.Buffer
	assert.NoError(t, f.Render(&buf))
	assert.Equal(t, `package testpkg

import "fmt"

func (u *MyUnionWithT[T]) AcceptFuncs(stringValFunc func(string) (T, error), boolValFunc func(bool) (T, error), unknownFunc func(string) (T, error)) (T, error) {
	var result T
	switch u.typ {
	default:
		if u.typ == "" {
			return result, fmt.Errorf("invalid value in union type")
		}
		return unknownFunc(u.typ)
	case "stringVal":
		if u.stringVal == nil {
			return result, fmt.Errorf("field \"stringVal\" is required")
		}
		return stringValFunc(*u.stringVal)
	case "boolVal":
		if u.boolVal == nil {
			return result, fmt.Errorf("field \"boolVal\" is required")
		}
		return boolValFunc(*u.boolVal)
	}
}
func (u *MyUnionWithT[T]) StringValNoopSuccess(string) (T, error) {
	var result T
	return result, nil
}
func (u *MyUnionWithT[T]) BoolValNoopSuccess(bool) (T, error) {
	var result T
	return result, nil
}
func (u *MyUnionWithT[T]) ErrorOnUnknown(typeName string) (T, error) {
	var result T
	return result, fmt.Errorf("invalid value in union type. Type name: %s", typeName)
}
`, buf.String())
}

func TestUnionWriter_unionTypeWithTAcceptFuncs_withOptional(t *testing.T) {
	f := jen.NewFile("testpkg")
	unionTypeWithTAcceptFuncs(f.Group, testUnionTypeWithOptional)
	var buf bytes.Buffer
	assert.NoError(t, f.Render(&buf))
	assert.Equal(t, `package testpkg

import "fmt"

func (u *MyUnionWithOptionalWithT[T]) AcceptFuncs(stringValFunc func(string) (T, error), optionalValFunc func(*string) (T, error), unknownFunc func(string) (T, error)) (T, error) {
	var result T
	switch u.typ {
	default:
		if u.typ == "" {
			return result, fmt.Errorf("invalid value in union type")
		}
		return unknownFunc(u.typ)
	case "stringVal":
		if u.stringVal == nil {
			return result, fmt.Errorf("field \"stringVal\" is required")
		}
		return stringValFunc(*u.stringVal)
	case "optionalVal":
		var optionalVal *string
		if u.optionalVal != nil {
			optionalVal = *u.optionalVal
		}
		return optionalValFunc(optionalVal)
	}
}
func (u *MyUnionWithOptionalWithT[T]) StringValNoopSuccess(string) (T, error) {
	var result T
	return result, nil
}
func (u *MyUnionWithOptionalWithT[T]) OptionalValNoopSuccess(*string) (T, error) {
	var result T
	return result, nil
}
func (u *MyUnionWithOptionalWithT[T]) ErrorOnUnknown(typeName string) (T, error) {
	var result T
	return result, fmt.Errorf("invalid value in union type. Type name: %s", typeName)
}
`, buf.String())
}

func TestUnionWriter_writeUnionTypeWithGenerics_withAcceptFuncs(t *testing.T) {
	f := jen.NewFile("testpkg")
	writeUnionTypeWithGenerics(f.Group, testUnionType, true)
	var buf bytes.Buffer
	assert.NoError(t, f.Render(&buf))
	// Should contain AcceptFuncs
	assert.Contains(t, buf.String(), "AcceptFuncs")
	assert.Contains(t, buf.String(), "StringValNoopSuccess")
	assert.Contains(t, buf.String(), "ErrorOnUnknown")
}

func TestUnionWriter_writeUnionTypeWithGenerics_withoutAcceptFuncs(t *testing.T) {
	f := jen.NewFile("testpkg")
	writeUnionTypeWithGenerics(f.Group, testUnionType, false)
	var buf bytes.Buffer
	assert.NoError(t, f.Render(&buf))
	// Should NOT contain AcceptFuncs
	assert.NotContains(t, buf.String(), "AcceptFuncs")
	assert.NotContains(t, buf.String(), "StringValNoopSuccess")
	assert.NotContains(t, buf.String(), "ErrorOnUnknown")
	// Should still contain Accept and visitor interface
	assert.Contains(t, buf.String(), "Accept(ctx context.Context")
	assert.Contains(t, buf.String(), "MyUnionVisitorWithT[T any]")
}
