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
	"fmt"
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
	// String returns a human-friendly name for log messages
	String() string

	IsNamed() bool
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
func (Any) String() string       { return "any" }

type Bearertoken struct{ base }

func (Bearertoken) Code() *jen.Statement { return snip.BearerTokenToken() }
func (Bearertoken) String() string       { return "bearertoken" }
func (Bearertoken) IsText() bool         { return true }

type Binary struct{ base }

func (Binary) Code() *jen.Statement { return jen.Op("[]").Byte() }
func (Binary) String() string       { return "binary" }
func (Binary) IsText() bool         { return true }
func (Binary) IsBinary() bool       { return true }

type Boolean struct{ base }

func (Boolean) Code() *jen.Statement { return jen.Bool() }
func (Boolean) String() string       { return "boolean" }
func (Boolean) IsBoolean() bool      { return true }

type DateTime struct{ base }

func (DateTime) Code() *jen.Statement { return snip.DateTimeDateTime() }
func (DateTime) String() string       { return "datetime" }
func (DateTime) IsText() bool         { return true }

type Double struct{ base }

func (Double) Code() *jen.Statement { return jen.Float64() }
func (Double) String() string       { return "double" }

type Integer struct{ base }

func (Integer) Code() *jen.Statement { return jen.Int() }
func (Integer) String() string       { return "integer" }

type RID struct{ base }

func (RID) Code() *jen.Statement { return snip.RIDResourceIdentifier() }
func (RID) String() string       { return "rid" }
func (RID) IsText() bool         { return true }

type Safelong struct{ base }

func (Safelong) Code() *jen.Statement { return snip.SafeLongSafeLong() }
func (Safelong) String() string       { return "safelong" }

type String struct{ base }

func (String) Code() *jen.Statement { return jen.String() }
func (String) String() string       { return "string" }
func (String) IsString() bool       { return true }
func (String) IsText() bool         { return true }

type UUID struct{ base }

func (UUID) Code() *jen.Statement { return snip.UUIDUUID() }
func (UUID) String() string       { return "uuid" }
func (UUID) IsText() bool         { return true }

// Composite Types

type Optional struct {
	Item Type
	base
}

func (t *Optional) Code() *jen.Statement {
	return jen.Op("*").Add(t.Item.Code())
}
func (t *Optional) String() string { return fmt.Sprintf("optional<%s>", t.Item.String()) }

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
func (t *List) String() string { return fmt.Sprintf("list<%s>", t.Item.String()) }

func (*List) IsCollection() bool { return true }
func (*List) IsList() bool       { return true }

func (t *List) Make() *jen.Statement {
	return jen.Make(t.Code(), jen.Lit(0))
}

type Set struct {
	Item Type
	base
}

func (t *Set) Code() *jen.Statement {
	return jen.Op("[]").Add(t.Item.Code())
}
func (t *Set) String() string { return fmt.Sprintf("set<%s>", t.Item) }

func (*Set) IsCollection() bool { return true }

func (t *Set) Make() *jen.Statement {
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

func (t *Map) String() string { return fmt.Sprintf("map<%s, %s>", t.Key, t.Val) }

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
func (t *AliasType) String() string { return fmt.Sprintf("%s (%s)", t.Name, t.Item) }

func (t *AliasType) Make() *jen.Statement {
	switch t.Item.(type) {
	case *Map, *List:
		return t.Item.Make()
	}
	if m := t.Item.Make(); m != nil {
		return t.Code().Call(m)
	}
	return nil
}

func (*AliasType) IsNamed() bool { return true }
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

func (t *EnumType) String() string { return t.Name }

func (*EnumType) IsNamed() bool { return true }
func (*EnumType) IsText() bool  { return true }

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

func (t *ObjectType) String() string { return t.Name }

func (*ObjectType) IsNamed() bool              { return true }
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

func (t *UnionType) String() string { return t.Name }

func (*UnionType) IsNamed() bool              { return true }
func (*UnionType) ContainsStrictFields() bool { return true }

type External struct {
	Spec     spec.TypeName
	Fallback Type
	base
}

func (t *External) Code() *jen.Statement {
	if !t.ExternalHasGoType() {
		// did not find expected delimiter in type name, use fallback
		return t.Fallback.Code()
	}
	pathAndName := strings.Split(t.Spec.Name, ":")
	importPath := t.Spec.Package + "." + pathAndName[0]
	goTypeName := pathAndName[1]
	return jen.Qual(importPath, goTypeName)
}

func (t *External) String() string {
	if !t.ExternalHasGoType() {
		// did not find expected delimiter in type name, use fallback
		return fmt.Sprintf("%s (%s)", t.Spec.Name, t.Fallback)
	}
	return t.Spec.Name
}

func (t *External) ExternalHasGoType() bool {
	return strings.Contains(t.Spec.Name, ":")
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
func (base) IsNamed() bool              { return false }
func (base) IsString() bool             { return false }
func (base) IsText() bool               { return false }
func (base) IsBinary() bool             { return false }
func (base) IsBoolean() bool            { return false }
func (base) IsOptional() bool           { return false }
func (base) IsCollection() bool         { return false }
func (base) IsList() bool               { return false }
func (base) ContainsStrictFields() bool { return false }
func (base) typ()                       {}
