// Copyright (c) 2018 Palantir Technologies. All rights reserved.
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

package conjure_test

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
	"testing"

	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/expression"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/conjure"
	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/conjure/types"
)

func TestParseStringParam(t *testing.T) {
	customTypes := types.NewCustomConjureTypes()
	err := customTypes.Add("Foo", "com.example.foo", types.SafeLong)
	require.NoError(t, err)
	for _, test := range []struct {
		Name            string
		ArgName         spec.ArgumentName
		ArgType         spec.Type
		ExpectedImports []string
		ExpectedSrc     string
		ExpectedErr     string
	}{
		{
			Name:        "Primitive string param",
			ArgName:     spec.ArgumentName("myArg"),
			ArgType:     spec.NewTypeFromPrimitive(spec.PrimitiveTypeString),
			ExpectedSrc: `myArg := myString`,
		},
		{
			Name:            "Primitive datetime param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.PrimitiveTypeDatetime),
			ExpectedImports: []string{"github.com/palantir/pkg/datetime"},
			ExpectedSrc: `myArg, err := datetime.ParseDateTime(myString)
if err != nil {
	return err
}`,
		},
		{
			Name:            "Primitive int param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.PrimitiveTypeInteger),
			ExpectedImports: []string{"strconv"},
			ExpectedSrc: `myArg, err := strconv.Atoi(myString)
if err != nil {
	return err
}`,
		},
		{
			Name:            "Primitive double param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.PrimitiveTypeDouble),
			ExpectedImports: []string{"strconv"},
			ExpectedSrc: `myArg, err := strconv.ParseFloat(myString, 64)
if err != nil {
	return err
}`,
		},
		{
			Name:            "Primitive safelong param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.PrimitiveTypeSafelong),
			ExpectedImports: []string{"github.com/palantir/pkg/safelong"},
			ExpectedSrc: `myArg, err := safelong.ParseSafeLong(myString)
if err != nil {
	return err
}`,
		},
		{
			Name:        "Primitive binary param",
			ArgName:     spec.ArgumentName("myArg"),
			ArgType:     spec.NewTypeFromPrimitive(spec.PrimitiveTypeBinary),
			ExpectedSrc: `myArg := []byte(myString)`,
		},
		{
			Name:        "Primitive any param",
			ArgName:     spec.ArgumentName("myArg"),
			ArgType:     spec.NewTypeFromPrimitive(spec.PrimitiveTypeAny),
			ExpectedSrc: `myArg := myString`,
		},
		{
			Name:            "Primitive boolean param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.PrimitiveTypeBoolean),
			ExpectedImports: []string{"strconv"},
			ExpectedSrc: `myArg, err := strconv.ParseBool(myString)
if err != nil {
	return err
}`,
		},
		{
			Name:            "Primitive uuid param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.PrimitiveTypeUuid),
			ExpectedImports: []string{"github.com/palantir/pkg/uuid"},
			ExpectedSrc: `myArg, err := uuid.ParseUUID(myString)
if err != nil {
	return err
}`,
		},
		{
			Name:            "Primitive rid param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.PrimitiveTypeRid),
			ExpectedImports: []string{"github.com/palantir/pkg/rid"},
			ExpectedSrc: `myArg, err := rid.ParseRID(myString)
if err != nil {
	return err
}`,
		},
		{
			Name:            "Primitive bearertoken param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.PrimitiveTypeBearertoken),
			ExpectedImports: []string{"github.com/palantir/pkg/bearertoken"},
			ExpectedSrc:     `myArg := bearertoken.Token(myString)`,
		},
		{
			Name:        "Primitive unknown param",
			ArgName:     spec.ArgumentName("myArg"),
			ArgType:     spec.NewTypeFromPrimitive(spec.PrimitiveTypeUnknown),
			ExpectedErr: "Unsupported primitive type UNKNOWN",
		},
		{
			Name:    "Optional string param",
			ArgName: spec.ArgumentName("myArg"),
			ArgType: spec.NewTypeFromOptional(spec.OptionalType{ItemType: spec.NewTypeFromPrimitive(spec.PrimitiveTypeString)}),
			ExpectedSrc: `var myArg *string
if myArgStr := myString; myArgStr != "" {
	myArgInternal := myArgStr
	myArg = &myArgInternal
}`,
		},
		{
			Name:            "Optional int param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromOptional(spec.OptionalType{ItemType: spec.NewTypeFromPrimitive(spec.PrimitiveTypeInteger)}),
			ExpectedImports: []string{"strconv"},
			ExpectedSrc: `var myArg *int
if myArgStr := myString; myArgStr != "" {
	myArgInternal, err := strconv.Atoi(myArgStr)
	if err != nil {
		return err
	}
	myArg = &myArgInternal
}`,
		},
		{
			Name:    "Reference param",
			ArgName: spec.ArgumentName("myArg"),
			ArgType: spec.NewTypeFromReference(spec.TypeName{
				Name:    "Foo",
				Package: "com.example.foo",
			}),
			ExpectedImports: []string{"strconv"},
			//TODO(bmoylan) This output is wrong - how are reference imports supposed to work?
			ExpectedSrc: `var myArg *int
if myArgStr := myString; myArgStr != "" {
	myArgInternal, err := strconv.Atoi(myArgStr)
	if err != nil {
		return err
	}
	myArg = &myArgInternal
}`,
		},
		{
			Name:    "External param",
			ArgName: spec.ArgumentName("myArg"),
			ArgType: spec.NewTypeFromExternal(spec.ExternalReference{
				ExternalReference: spec.TypeName{
					Name:    "foo:Foo",
					Package: "com.example.foo",
				},
				Fallback: spec.NewTypeFromPrimitive(spec.PrimitiveTypeString),
			}),
			ExpectedImports: []string{"com.example.foo.foo"},
			//TODO(bmoylan) This output is wrong - how are external imports supposed to work?
			ExpectedSrc: `myArgInternal := myString
myArg = com.example.foo.foo.Foo(myArgInternal)`,
		},
		{
			Name:        "List param",
			ArgName:     spec.ArgumentName("myArg"),
			ArgType:     spec.NewTypeFromList(spec.ListType{ItemType: spec.NewTypeFromPrimitive(spec.PrimitiveTypeInteger)}),
			ExpectedErr: "can not assign string expression to list type",
		},
		{
			Name:    "Map param",
			ArgName: spec.ArgumentName("myArg"),
			ArgType: spec.NewTypeFromMap(spec.MapType{
				KeyType:   spec.NewTypeFromPrimitive(spec.PrimitiveTypeString),
				ValueType: spec.NewTypeFromPrimitive(spec.PrimitiveTypeInteger)}),
			ExpectedErr: "can not assign string expression to map type",
		},
		{
			Name:        "Set param",
			ArgName:     spec.ArgumentName("myArg"),
			ArgType:     spec.NewTypeFromSet(spec.SetType{ItemType: spec.NewTypeFromPrimitive(spec.PrimitiveTypeInteger)}),
			ExpectedErr: "can not assign string expression to set type",
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			info := types.NewPkgInfo("", customTypes)
			stmts, err := conjure.ParseStringParam(
				test.ArgName,
				test.ArgType,
				expression.VariableVal("myString"),
				info,
			)
			if test.ExpectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.ExpectedErr)
				return
			}
			require.NoError(t, err)

			var buf bytes.Buffer
			fset := token.NewFileSet()
			require.NoError(t, format.Node(&buf, fset, createStmts(stmts)))
			assert.Equal(t, strings.Split(test.ExpectedSrc, "\n"), strings.Split(buf.String(), "\n"))
			importMap := info.ImportAliases()
			var imports []string
			for k := range importMap {
				imports = append(imports, k)
			}
			assert.Equal(t, test.ExpectedImports, imports)
		})
	}
}

func createStmts(in []astgen.ASTStmt) (out []ast.Stmt) {
	for _, stmt := range in {
		out = append(out, stmt.ASTStmt())
	}
	return out
}
