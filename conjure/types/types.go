// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package types

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/snip"
)

type Type interface {
	// Code returns the fully-qualified go type name.
	// If using exported names within the same package, jen.Qual will handle omitting the import.
	Code() *jen.Statement
	// Make returns an expression to `make` a collection type, if required.
	// If the type does not require initialization, Make returns nil.
	Make() *jen.Statement

	IsString() bool
	IsText() bool
	IsBinary() bool
	IsBoolean() bool
	IsOptional() bool
	IsCollection() bool
	IsList() bool
	ContainsStrictFields() bool

	typ() // block external implementations
}

// Primitive Types

type Any struct{ base }

func (Any) Code() *jen.Statement { return jen.Interface() }

type Bearertoken struct{ base }

func (Bearertoken) Code() *jen.Statement { return snip.BearerTokenToken() }
func (Bearertoken) IsText() bool         { return true }

type Binary struct{ base }

func (Binary) Code() *jen.Statement { return jen.Op("[]").Byte() }
func (Binary) IsText() bool         { return true }
func (Binary) IsBinary() bool       { return true }

type Boolean struct{ base }

func (Boolean) Code() *jen.Statement { return jen.Bool() }
func (Boolean) IsBoolean() bool      { return true }

type DateTime struct{ base }

func (DateTime) Code() *jen.Statement { return snip.DateTimeDateTime() }
func (DateTime) IsText() bool         { return true }

type Double struct{ base }

func (Double) Code() *jen.Statement { return jen.Float64() }

type Integer struct{ base }

func (Integer) Code() *jen.Statement { return jen.Int() }

type RID struct{ base }

func (RID) Code() *jen.Statement { return snip.RIDResourceIdentifier() }
func (RID) IsText() bool         { return true }

type Safelong struct{ base }

func (Safelong) Code() *jen.Statement { return snip.SafeLongSafeLong() }

type String struct{ base }

func (String) Code() *jen.Statement { return jen.String() }
func (String) IsString() bool       { return true }
func (String) IsText() bool         { return true }

type UUID struct{ base }

func (UUID) Code() *jen.Statement { return snip.UUIDUUID() }
func (UUID) IsText() bool         { return true }

// Composite Types

type Optional struct {
	Item Type
	base
}

func (t *Optional) Code() *jen.Statement {
	return jen.Op("*").Add(t.Item.Code())
}

func (t *Optional) Make() *jen.Statement {
	// Optionals never get default initialization, even if the underlying does.
	return nil
}

func (t *Optional) IsString() bool             { return t.Item.IsString() }
func (t *Optional) IsText() bool               { return t.Item.IsText() }
func (t *Optional) IsBinary() bool             { return t.Item.IsBinary() }
func (t *Optional) IsBoolean() bool            { return t.Item.IsBoolean() }
func (t *Optional) IsOptional() bool           { return true }
func (t *Optional) IsCollection() bool         { return t.Item.IsCollection() }
func (t *Optional) IsList() bool               { return t.Item.IsList() }
func (t *Optional) ContainsStrictFields() bool { return t.Item.ContainsStrictFields() }

type List struct {
	Item Type
	base
}

func (t *List) Code() *jen.Statement {
	return jen.Op("[]").Add(t.Item.Code())
}

func (*List) IsCollection() bool { return true }
func (*List) IsList() bool       { return true }

func (t *List) Make() *jen.Statement {
	return jen.Make(t.Code(), jen.Lit(0))
}

type Map struct {
	Key Type
	Val Type
	base
}

func (t *Map) Code() *jen.Statement {
	var mapKey *jen.Statement
	switch {
	case t.Key.IsBinary():
		mapKey = jen.Map(snip.BinaryBinary())
	case t.Key.IsBoolean():
		mapKey = jen.Map(snip.BooleanBoolean())
	default:
		mapKey = jen.Map(t.Key.Code())
	}
	return mapKey.Add(t.Val.Code())
}

func (t *Map) IsCollection() bool { return true }

func (t *Map) Make() *jen.Statement {
	return jen.Make(t.Code(), jen.Lit(0))
}

// Named Types

type AliasType struct {
	Docs
	Name       string
	Item       Type
	conjurePkg string
	importPath string
	base
}

func (t *AliasType) Code() *jen.Statement {
	return jen.Qual(t.importPath, t.Name)
}

func (t *AliasType) Make() *jen.Statement {
	switch t.Item.(type) {
	case *Map, *List:
		return t.Item.Code()
	}
	if m := t.Item.Make(); m != nil {
		return t.Code().Call(m)
	}
	return nil
}

func (t *AliasType) IsString() bool {
	_, isString := t.Item.(String)
	return isString
}
func (t *AliasType) IsText() bool    { return t.Item.IsText() }
func (t *AliasType) IsBinary() bool  { return t.Item.IsBinary() }
func (t *AliasType) IsBoolean() bool { return t.Item.IsBoolean() }
func (t *AliasType) IsOptional() bool {
	_, isOptional := t.Item.(*Optional)
	return isOptional
}
func (t *AliasType) IsCollection() bool         { return t.Item.IsCollection() }
func (t *AliasType) IsList() bool               { return t.Item.IsList() }
func (t *AliasType) ContainsStrictFields() bool { return t.Item.ContainsStrictFields() }

type EnumType struct {
	Docs
	Name       string
	Values     []*Field
	conjurePkg string
	importPath string
	base
}

func (t *EnumType) Code() *jen.Statement {
	return jen.Qual(t.importPath, t.Name)
}

func (t *EnumType) IsText() bool { return true }

type ObjectType struct {
	Docs
	Name       string
	Fields     []*Field
	conjurePkg string
	importPath string
	base
}

func (t *ObjectType) Code() *jen.Statement {
	return jen.Qual(t.importPath, t.Name)
}

func (*ObjectType) ContainsStrictFields() bool { return true }

type UnionType struct {
	Docs
	Name       string
	Fields     []*Field
	conjurePkg string
	importPath string
	base
}

func (t *UnionType) Code() *jen.Statement {
	return jen.Qual(t.importPath, t.Name)
}

func (*UnionType) ContainsStrictFields() bool { return true }

type ErrorType struct {
	Docs
	Name           string
	ErrorNamespace spec.ErrorNamespace
	ErrorCode      spec.ErrorCode
	SafeArgs       []*Field
	UnsafeArgs     []*Field
	conjurePkg     string
	importPath     string
	base
}

func (t *ErrorType) Code() *jen.Statement {
	return jen.Qual(t.importPath, t.Name)
}

type External struct {
	Spec     spec.TypeName
	fallback Type
	base
}

func (t *External) Code() *jen.Statement {
	if !strings.Contains(t.Spec.Name, ":") {
		// did not find expected delimiter in type name, use fallback
		return t.fallback.Code()
	}
	pathAndName := strings.Split(t.Spec.Name, ":")
	importPath := t.Spec.Package + "." + pathAndName[0]
	goTypeName := pathAndName[1]
	return jen.Qual(importPath, goTypeName)
}

// Public member types

type Docs string

func (c Docs) CommentLine() *jen.Statement {
	if c != "" {
		return jen.Comment(string(c)).Line()
	}
	return jen.Empty()
}

type EnumValue struct {
	Docs
	Value string
}

type Field struct {
	Docs
	Deprecated Docs
	Name       string // JSON key or enum value
	Type       Type   // string for enum value
}

// private utility types

type base struct{}

func (base) Make() *jen.Statement       { return nil }
func (base) IsString() bool             { return false }
func (base) IsText() bool               { return false }
func (base) IsBinary() bool             { return false }
func (base) IsBoolean() bool            { return false }
func (base) IsOptional() bool           { return false }
func (base) IsCollection() bool         { return false }
func (base) IsList() bool               { return false }
func (base) ContainsStrictFields() bool { return false }
func (base) typ()                       {}
