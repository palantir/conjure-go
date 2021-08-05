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

package visitors_test

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
	"testing"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/conjure-go/v6/conjure/visitors"
	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/expression"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatementsForHTTPParam(t *testing.T) {
	customTypes := types.NewCustomConjureTypes()
	err := customTypes.Add("Foo", "com.example.foo", types.SafeLong, nil)
	require.NoError(t, err)
	err = customTypes.Add("com.example.foo.FooId", "com.example.foo", types.SafeLong, nil)
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
			Name:            "Primitive bearertoken param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_BEARERTOKEN)),
			ExpectedImports: []string{"github.com/palantir/pkg/bearertoken"},
			ExpectedSrc:     `myArg := bearertoken.Token(myString)`,
		},
		{
			Name:            "Optional bearertoken param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromOptional(spec.OptionalType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_BEARERTOKEN))}),
			ExpectedImports: []string{"github.com/palantir/pkg/bearertoken"},
			ExpectedSrc: `var myArg *bearertoken.Token
if myArgStr := myString; myArgStr != "" {
	myArgInternal := bearertoken.Token(myArgStr)
	myArg = &myArgInternal
}`,
		},
		{
			Name:            "List bearertoken param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromList(spec.ListType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_BEARERTOKEN))}),
			ExpectedImports: []string{"github.com/palantir/pkg/bearertoken"},
			ExpectedSrc: `var myArg []bearertoken.Token
for _, v := range myString {
	convertedVal := bearertoken.Token(v)
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:            "Set bearertoken param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromSet(spec.SetType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_BEARERTOKEN))}),
			ExpectedImports: []string{"github.com/palantir/pkg/bearertoken"},
			ExpectedSrc: `var myArg []bearertoken.Token
for _, v := range myString {
	convertedVal := bearertoken.Token(v)
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:        "Primitive binary param",
			ArgName:     spec.ArgumentName("myArg"),
			ArgType:     spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_BINARY)),
			ExpectedSrc: `myArg := []byte(myString)`,
		},
		{
			Name:            "Primitive boolean param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_BOOLEAN)),
			ExpectedImports: []string{"strconv"},
			ExpectedSrc: `myArg, err := strconv.ParseBool(myString)
if err != nil {
	return err
}`,
		},
		{
			Name:            "Optional boolean param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromOptional(spec.OptionalType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_BOOLEAN))}),
			ExpectedImports: []string{"strconv"},
			ExpectedSrc: `var myArg *bool
if myArgStr := myString; myArgStr != "" {
	myArgInternal, err := strconv.ParseBool(myArgStr)
	if err != nil {
		return err
	}
	myArg = &myArgInternal
}`,
		},
		{
			Name:            "List boolean param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromList(spec.ListType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_BOOLEAN))}),
			ExpectedImports: []string{"strconv"},
			ExpectedSrc: `var myArg []bool
for _, v := range myString {
	convertedVal, err := strconv.ParseBool(v)
	if err != nil {
		return err
	}
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:            "Set bool param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromSet(spec.SetType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_BOOLEAN))}),
			ExpectedImports: []string{"strconv"},
			ExpectedSrc: `var myArg []bool
for _, v := range myString {
	convertedVal, err := strconv.ParseBool(v)
	if err != nil {
		return err
	}
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:            "Primitive datetime param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_DATETIME)),
			ExpectedImports: []string{"github.com/palantir/pkg/datetime"},
			ExpectedSrc: `myArg, err := datetime.ParseDateTime(myString)
if err != nil {
	return err
}`,
		},
		{
			Name:            "Optional datetime param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromOptional(spec.OptionalType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_DATETIME))}),
			ExpectedImports: []string{"github.com/palantir/pkg/datetime"},
			ExpectedSrc: `var myArg *datetime.DateTime
if myArgStr := myString; myArgStr != "" {
	myArgInternal, err := datetime.ParseDateTime(myArgStr)
	if err != nil {
		return err
	}
	myArg = &myArgInternal
}`,
		},
		{
			Name:            "List datetime param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromList(spec.ListType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_DATETIME))}),
			ExpectedImports: []string{"github.com/palantir/pkg/datetime"},
			ExpectedSrc: `var myArg []datetime.DateTime
for _, v := range myString {
	convertedVal, err := datetime.ParseDateTime(v)
	if err != nil {
		return err
	}
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:            "Set datetime param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromSet(spec.SetType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_DATETIME))}),
			ExpectedImports: []string{"github.com/palantir/pkg/datetime"},
			ExpectedSrc: `var myArg []datetime.DateTime
for _, v := range myString {
	convertedVal, err := datetime.ParseDateTime(v)
	if err != nil {
		return err
	}
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:            "Primitive double param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_DOUBLE)),
			ExpectedImports: []string{"strconv"},
			ExpectedSrc: `myArg, err := strconv.ParseFloat(myString, 64)
if err != nil {
	return err
}`,
		},
		{
			Name:            "Optional double param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromOptional(spec.OptionalType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_DOUBLE))}),
			ExpectedImports: []string{"strconv"},
			ExpectedSrc: `var myArg *float64
if myArgStr := myString; myArgStr != "" {
	myArgInternal, err := strconv.ParseFloat(myArgStr, 64)
	if err != nil {
		return err
	}
	myArg = &myArgInternal
}`,
		},
		{
			Name:            "List double param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromList(spec.ListType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_DOUBLE))}),
			ExpectedImports: []string{"strconv"},
			ExpectedSrc: `var myArg []float64
for _, v := range myString {
	convertedVal, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return err
	}
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:            "Set double param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromSet(spec.SetType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_DOUBLE))}),
			ExpectedImports: []string{"strconv"},
			ExpectedSrc: `var myArg []float64
for _, v := range myString {
	convertedVal, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return err
	}
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:            "Primitive int param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_INTEGER)),
			ExpectedImports: []string{"strconv"},
			ExpectedSrc: `myArg, err := strconv.Atoi(myString)
if err != nil {
	return err
}`,
		},
		{
			Name:            "Optional int param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromOptional(spec.OptionalType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_INTEGER))}),
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
			Name:            "List int param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromList(spec.ListType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_INTEGER))}),
			ExpectedImports: []string{"strconv"},
			ExpectedSrc: `var myArg []int
for _, v := range myString {
	convertedVal, err := strconv.Atoi(v)
	if err != nil {
		return err
	}
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:            "Set int param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromSet(spec.SetType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_INTEGER))}),
			ExpectedImports: []string{"strconv"},
			ExpectedSrc: `var myArg []int
for _, v := range myString {
	convertedVal, err := strconv.Atoi(v)
	if err != nil {
		return err
	}
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:            "Primitive rid param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_RID)),
			ExpectedImports: []string{"github.com/palantir/pkg/rid"},
			ExpectedSrc: `myArg, err := rid.ParseRID(myString)
if err != nil {
	return err
}`,
		},
		{
			Name:            "Optional rid param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromOptional(spec.OptionalType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_RID))}),
			ExpectedImports: []string{"github.com/palantir/pkg/rid"},
			ExpectedSrc: `var myArg *rid.ResourceIdentifier
if myArgStr := myString; myArgStr != "" {
	myArgInternal, err := rid.ParseRID(myArgStr)
	if err != nil {
		return err
	}
	myArg = &myArgInternal
}`,
		},
		{
			Name:            "List rid param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromList(spec.ListType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_RID))}),
			ExpectedImports: []string{"github.com/palantir/pkg/rid"},
			ExpectedSrc: `var myArg []rid.ResourceIdentifier
for _, v := range myString {
	convertedVal, err := rid.ParseRID(v)
	if err != nil {
		return err
	}
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:            "Set rid param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromSet(spec.SetType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_RID))}),
			ExpectedImports: []string{"github.com/palantir/pkg/rid"},
			ExpectedSrc: `var myArg []rid.ResourceIdentifier
for _, v := range myString {
	convertedVal, err := rid.ParseRID(v)
	if err != nil {
		return err
	}
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:            "Primitive safelong param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_SAFELONG)),
			ExpectedImports: []string{"github.com/palantir/pkg/safelong"},
			ExpectedSrc: `myArg, err := safelong.ParseSafeLong(myString)
if err != nil {
	return err
}`,
		},
		{
			Name:            "Optional safelong param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromOptional(spec.OptionalType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_SAFELONG))}),
			ExpectedImports: []string{"github.com/palantir/pkg/safelong"},
			ExpectedSrc: `var myArg *safelong.SafeLong
if myArgStr := myString; myArgStr != "" {
	myArgInternal, err := safelong.ParseSafeLong(myArgStr)
	if err != nil {
		return err
	}
	myArg = &myArgInternal
}`,
		},
		{
			Name:            "List safelong param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromList(spec.ListType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_SAFELONG))}),
			ExpectedImports: []string{"github.com/palantir/pkg/safelong"},
			ExpectedSrc: `var myArg []safelong.SafeLong
for _, v := range myString {
	convertedVal, err := safelong.ParseSafeLong(v)
	if err != nil {
		return err
	}
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:            "Set safelong param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromSet(spec.SetType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_SAFELONG))}),
			ExpectedImports: []string{"github.com/palantir/pkg/safelong"},
			ExpectedSrc: `var myArg []safelong.SafeLong
for _, v := range myString {
	convertedVal, err := safelong.ParseSafeLong(v)
	if err != nil {
		return err
	}
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:        "Primitive string param",
			ArgName:     spec.ArgumentName("myArg"),
			ArgType:     spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_STRING)),
			ExpectedSrc: `myArg := myString`,
		},
		{
			Name:    "Optional string param",
			ArgName: spec.ArgumentName("myArg"),
			ArgType: spec.NewTypeFromOptional(spec.OptionalType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_STRING))}),
			ExpectedSrc: `var myArg *string
if myArgStr := myString; myArgStr != "" {
	myArgInternal := myArgStr
	myArg = &myArgInternal
}`,
		},
		{
			Name:        "List string param",
			ArgName:     spec.ArgumentName("myArg"),
			ArgType:     spec.NewTypeFromList(spec.ListType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_STRING))}),
			ExpectedSrc: `myArg := myString`,
		},
		{
			Name:        "Set string param",
			ArgName:     spec.ArgumentName("myArg"),
			ArgType:     spec.NewTypeFromSet(spec.SetType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_STRING))}),
			ExpectedSrc: `myArg := myString`,
		},
		{
			Name:            "Primitive uuid param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_UUID)),
			ExpectedImports: []string{"github.com/palantir/pkg/uuid"},
			ExpectedSrc: `myArg, err := uuid.ParseUUID(myString)
if err != nil {
	return err
}`,
		},
		{
			Name:            "Optional uuid param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromOptional(spec.OptionalType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_UUID))}),
			ExpectedImports: []string{"github.com/palantir/pkg/uuid"},
			ExpectedSrc: `var myArg *uuid.UUID
if myArgStr := myString; myArgStr != "" {
	myArgInternal, err := uuid.ParseUUID(myArgStr)
	if err != nil {
		return err
	}
	myArg = &myArgInternal
}`,
		},
		{
			Name:            "List uuid param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromList(spec.ListType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_UUID))}),
			ExpectedImports: []string{"github.com/palantir/pkg/uuid"},
			ExpectedSrc: `var myArg []uuid.UUID
for _, v := range myString {
	convertedVal, err := uuid.ParseUUID(v)
	if err != nil {
		return err
	}
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:            "Set uuid param",
			ArgName:         spec.ArgumentName("myArg"),
			ArgType:         spec.NewTypeFromSet(spec.SetType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_UUID))}),
			ExpectedImports: []string{"github.com/palantir/pkg/uuid"},
			ExpectedSrc: `var myArg []uuid.UUID
for _, v := range myString {
	convertedVal, err := uuid.ParseUUID(v)
	if err != nil {
		return err
	}
	myArg = append(myArg, convertedVal)
}`,
		},
		{
			Name:        "Primitive any param",
			ArgName:     spec.ArgumentName("myArg"),
			ArgType:     spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_ANY)),
			ExpectedSrc: `myArg := myString`,
		},
		{
			Name:        "Primitive unknown param",
			ArgName:     spec.ArgumentName("myArg"),
			ArgType:     spec.NewTypeFromPrimitive(spec.New_PrimitiveType("unknown")),
			ExpectedErr: "Unsupported primitive type unknown",
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
			Name:    "Reference string param",
			ArgName: spec.ArgumentName("myArg"),
			ArgType: spec.NewTypeFromReference(spec.TypeName{
				Name:    "FooId",
				Package: "com.example.foo",
			}),
			ExpectedImports: []string{"github.com/palantir/pkg/safejson", "github.com/palantir/pkg/safelong", "strconv"},
			ExpectedSrc: `var myArg safelong.SafeLong
myArgQuote := strconv.Quote(myString)
if err := safejson.Unmarshal([]byte(myArgQuote), &myArg); err != nil {
	return werror.Wrap(err, "failed to unmarshal argument", werror.SafeParam("argName", "myArg"), werror.SafeParam("argType", "safelong.SafeLong"))
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
				Fallback: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_STRING)),
			}),
			ExpectedImports: []string{"com.example.foo.foo"},
			//TODO(bmoylan) This output is wrong - how are external imports supposed to work?
			ExpectedSrc: `myArgInternal := myString
myArg := com.example.foo.foo.Foo(myArgInternal)`,
		},
		{
			Name:    "Map param",
			ArgName: spec.ArgumentName("myArg"),
			ArgType: spec.NewTypeFromMap(spec.MapType{
				KeyType:   spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_STRING)),
				ValueType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_INTEGER))}),
			ExpectedErr: "can not assign string expression to map type",
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			if test.Name == "Reference param" {
				t.Skip()
			}
			info := types.NewPkgInfo("", customTypes)
			stmts, err := visitors.StatementsForHTTPParam(
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
			assert.Equal(t, len(test.ExpectedImports), len(imports))
			for _, actualImport := range imports {
				assert.Contains(t, test.ExpectedImports, actualImport)
			}
		})
	}
}

func createStmts(in []astgen.ASTStmt) (out []ast.Stmt) {
	for _, stmt := range in {
		out = append(out, stmt.ASTStmt())
	}
	return out
}
