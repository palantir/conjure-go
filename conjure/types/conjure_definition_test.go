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
	"io/ioutil"
	"testing"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConjureDefinition(t *testing.T) {
	outputDir := "./test"
	for _, test := range []struct {
		Name string
		In   spec.ConjureDefinition
		Out  *ConjureDefinition
	}{
		{
			Name: "full-featured object definition",
			In: spec.ConjureDefinition{
				Version: 1,
				Errors:  nil,
				Types: []spec.TypeDefinition{
					spec.NewTypeDefinitionFromObject(spec.ObjectDefinition{
						TypeName: spec.TypeName{
							Name:    "BackingFileSystem",
							Package: "com.palantir.foundry.catalog.api.datasets",
						},
						Fields: []spec.FieldDefinition{
							{
								FieldName: "fileSystemId",
								Type:      newPrimitive(spec.PrimitiveType_STRING),
								Docs:      docsPtr("The name by which this file system is identified."),
							},
							{
								FieldName: "baseUri",
								Type:      newPrimitive(spec.PrimitiveType_STRING),
							},
							{
								FieldName: "exenum",
								Type:      spec.NewTypeFromReference(spec.TypeName{Name: "ExampleEnumeration", Package: "example.api"}),
							},
							{
								FieldName: "client",
								Type: spec.NewTypeFromExternal(spec.ExternalReference{
									ExternalReference: spec.TypeName{
										Name:    "com/palantir/go-palantir/httpclient:RESTClient",
										Package: "github",
									},
									Fallback: newPrimitive(spec.PrimitiveType_STRING),
								}),
							},
						},
						Docs: docsPtr("Optional Docs"),
					}),
					spec.NewTypeDefinitionFromObject(spec.ObjectDefinition{
						TypeName: spec.TypeName{
							Name:    "TestType",
							Package: "com.palantir.foundry.catalog.api.datasets",
						},
						Fields: []spec.FieldDefinition{
							{
								FieldName: "alias",
								Type:      spec.NewTypeFromReference(spec.TypeName{Name: "ExampleAlias", Package: "com.palantir.test.api"}),
							},
							{
								FieldName: "rid",
								Type:      newPrimitive(spec.PrimitiveType_RID),
							},
							{
								FieldName: "large_int",
								Type:      newPrimitive(spec.PrimitiveType_SAFELONG),
							},
							{
								FieldName: "time",
								Type:      newPrimitive(spec.PrimitiveType_DATETIME),
							},
							{
								FieldName: "bytes",
								Type:      newPrimitive(spec.PrimitiveType_BINARY),
							},
						},
					}),
					spec.NewTypeDefinitionFromEnum(spec.EnumDefinition{
						TypeName: spec.TypeName{
							Name:    "ExampleEnumeration",
							Package: "example.api",
						},
						Values: []spec.EnumValueDefinition{{Value: "A"}, {Value: "B"}},
					}),
					spec.NewTypeDefinitionFromEnum(spec.EnumDefinition{
						TypeName: spec.TypeName{
							Name:    "Months",
							Package: "com.palantir.test.api",
						},
						Values: []spec.EnumValueDefinition{{Value: "JANUARY"}, {Value: "MULTI_MONTHS"}},
					}),
					spec.NewTypeDefinitionFromEnum(spec.EnumDefinition{
						TypeName: spec.TypeName{
							Name:    "Days",
							Package: "com.palantir.test.api",
						},
						Values: []spec.EnumValueDefinition{{Value: "FRIDAY"}, {Value: "SATURDAY"}},
					}),
					spec.NewTypeDefinitionFromAlias(spec.AliasDefinition{
						TypeName: spec.TypeName{
							Name:    "ExampleAlias",
							Package: "com.palantir.test.api",
						},
						Alias: newPrimitive(spec.PrimitiveType_STRING),
					}),
					spec.NewTypeDefinitionFromAlias(spec.AliasDefinition{
						TypeName: spec.TypeName{
							Name:    "LongAlias",
							Package: "com.palantir.test.api",
						},
						Alias: newPrimitive(spec.PrimitiveType_SAFELONG),
					}),
					spec.NewTypeDefinitionFromAlias(spec.AliasDefinition{
						TypeName: spec.TypeName{
							Name:    "Status",
							Package: "com.palantir.test.api",
						},
						Alias: newPrimitive(spec.PrimitiveType_INTEGER),
					}),
					spec.NewTypeDefinitionFromAlias(spec.AliasDefinition{
						TypeName: spec.TypeName{
							Name:    "ObjectAlias",
							Package: "com.palantir.test.api",
						},
						Alias: spec.NewTypeFromReference(spec.TypeName{Name: "TestType", Package: "com.palantir.foundry.catalog.api.datasets"}),
					}),
					spec.NewTypeDefinitionFromAlias(spec.AliasDefinition{
						TypeName: spec.TypeName{
							Name:    "AliasAlias",
							Package: "com.palantir.test.api",
						},
						Alias: spec.NewTypeFromReference(spec.TypeName{Name: "Status", Package: "com.palantir.test.api"}),
					}),
					spec.NewTypeDefinitionFromUnion(spec.UnionDefinition{
						TypeName: spec.TypeName{
							Name:    "ExampleUnion",
							Package: "com.palantir.test.api",
						},
						Union: []spec.FieldDefinition{
							{
								FieldName: "str",
								Type:      newPrimitive(spec.PrimitiveType_STRING),
							},
							{
								FieldName: "other",
								Type:      newPrimitive(spec.PrimitiveType_STRING),
								Docs:      docsPtr("Another string"),
							},
							{
								FieldName: "myMap",
								Type: spec.NewTypeFromMap(spec.MapType{
									KeyType: newPrimitive(spec.PrimitiveType_STRING),
									ValueType: spec.NewTypeFromList(spec.ListType{
										ItemType: newPrimitive(spec.PrimitiveType_INTEGER),
									}),
								}),
							},
							{
								FieldName: "tester",
								Type:      spec.NewTypeFromReference(spec.TypeName{Name: "TestType", Package: "com.palantir.foundry.catalog.api.datasets"}),
							},
							{
								FieldName: "recursive",
								Type:      spec.NewTypeFromReference(spec.TypeName{Name: "ExampleUnion", Package: "com.palantir.test.api"}),
							},
						},
					}),
				},
				Services:   nil,
				Extensions: nil,
			},
			Out: &ConjureDefinition{
				Version: 1,
				Packages: map[string]ConjurePackage{
					"com.palantir.foundry.catalog.api.datasets": {
						ConjurePackage: "com.palantir.foundry.catalog.api.datasets",
						ImportPath:     "github.com/palantir/conjure-go/v6/conjure/types/test/foundry/catalog/api/datasets",
						OutputDir:      "test/foundry/catalog/api/datasets",
						PackageName:    "datasets",
						Objects: []*ObjectType{
							{
								Name:       "BackingFileSystem",
								Docs:       "Optional Docs",
								conjurePkg: "com.palantir.foundry.catalog.api.datasets",
								importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/foundry/catalog/api/datasets",
								Fields: []*Field{
									{
										Name: "fileSystemId",
										Docs: "The name by which this file system is identified.",
										Type: String{},
									},
									{
										Name: "baseUri",
										Type: String{},
									},
									{
										Name: "exenum",
										Type: &EnumType{
											Name:       "ExampleEnumeration",
											conjurePkg: "example.api",
											importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/example/api",
											Values: []*Field{
												{Name: "A", Type: String{}},
												{Name: "B", Type: String{}},
											},
										},
									},
									{
										Name: "client",
										Type: &External{
											Spec: spec.TypeName{
												Name:    "com/palantir/go-palantir/httpclient:RESTClient",
												Package: "github",
											},
											Fallback: String{},
										},
									},
								},
							},
							{
								Name:       "TestType",
								conjurePkg: "com.palantir.foundry.catalog.api.datasets",
								importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/foundry/catalog/api/datasets",
								Fields: []*Field{
									{
										Name: "alias",
										Type: &AliasType{
											Name:       "ExampleAlias",
											Item:       String{},
											conjurePkg: "com.palantir.test.api",
											importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
										},
									},
									{Name: "rid", Type: RID{}},
									{Name: "large_int", Type: Safelong{}},
									{Name: "time", Type: DateTime{}},
									{Name: "bytes", Type: Binary{}},
								},
							},
						},
					},
					"com.palantir.test.api": {
						ConjurePackage: "com.palantir.test.api",
						ImportPath:     "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
						OutputDir:      "test/test/api",
						PackageName:    "api",
						Aliases: []*AliasType{
							{
								Name:       "ExampleAlias",
								Item:       String{},
								conjurePkg: "com.palantir.test.api",
								importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
							},
							{
								Name:       "LongAlias",
								Item:       Safelong{},
								conjurePkg: "com.palantir.test.api",
								importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
							},
							{
								Name:       "Status",
								Item:       Integer{},
								conjurePkg: "com.palantir.test.api",
								importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
							},
							{
								Name: "ObjectAlias",
								Item: &ObjectType{
									Name: "TestType",
									Fields: []*Field{
										{
											Name: "alias",
											Type: &AliasType{
												Name:       "ExampleAlias",
												Item:       String{},
												conjurePkg: "com.palantir.test.api",
												importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
											},
										},
										{Name: "rid", Type: RID{}},
										{Name: "large_int", Type: Safelong{}},
										{Name: "time", Type: DateTime{}},
										{Name: "bytes", Type: Binary{}},
									},
									conjurePkg: "com.palantir.foundry.catalog.api.datasets",
									importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/foundry/catalog/api/datasets",
								},
								conjurePkg: "com.palantir.test.api",
								importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
							},
							{
								Name: "AliasAlias",
								Item: &AliasType{
									Name:       "Status",
									Item:       Integer{},
									conjurePkg: "com.palantir.test.api",
									importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
								},
								conjurePkg: "com.palantir.test.api",
								importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
							},
						},
						Enums: []*EnumType{
							{
								Name:       "Months",
								Values:     []*Field{{Name: "JANUARY", Type: String{}}, {Name: "MULTI_MONTHS", Type: String{}}},
								conjurePkg: "com.palantir.test.api",
								importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
							},
							{
								Name:       "Days",
								Values:     []*Field{{Name: "FRIDAY", Type: String{}}, {Name: "SATURDAY", Type: String{}}},
								conjurePkg: "com.palantir.test.api",
								importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
							},
						},
						Unions: []*UnionType{
							func() *UnionType {
								// Use function so we can set up recursive field
								u := &UnionType{
									Name: "ExampleUnion",
									Fields: []*Field{
										{
											Name: "str",
											Type: String{},
										},
										{
											Docs: "Another string",
											Name: "other",
											Type: String{},
										},
										{
											Name: "myMap",
											Type: &Map{Key: String{}, Val: &List{Item: Integer{}}},
										},
										{
											Name: "tester",
											Type: &ObjectType{
												Name: "TestType",
												Fields: []*Field{
													{
														Name: "alias",
														Type: &AliasType{
															Name:       "ExampleAlias",
															Item:       String{},
															conjurePkg: "com.palantir.test.api",
															importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
														},
													},
													{Name: "rid", Type: RID{}},
													{Name: "large_int", Type: Safelong{}},
													{Name: "time", Type: DateTime{}},
													{Name: "bytes", Type: Binary{}},
												},
												conjurePkg: "com.palantir.foundry.catalog.api.datasets",
												importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/foundry/catalog/api/datasets",
											},
										},
									},
									conjurePkg: "com.palantir.test.api",
									importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
								}
								u.Fields = append(u.Fields, &Field{
									Name: "recursive",
									Type: u,
								})
								return u
							}(),
						},
					},
					"example.api": {
						ConjurePackage: "example.api",
						ImportPath:     "github.com/palantir/conjure-go/v6/conjure/types/test/example/api",
						OutputDir:      "test/example/api",
						PackageName:    "api",
						Enums: []*EnumType{{
							Name:       "ExampleEnumeration",
							Values:     []*Field{{Name: "A", Type: String{}}, {Name: "B", Type: String{}}},
							conjurePkg: "example.api",
							importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/example/api",
						}},
					},
				},
				Extensions: nil,
			},
		},
		{
			Name: "full-featured service definition",
			In: spec.ConjureDefinition{
				Version: 1,
				Services: []spec.ServiceDefinition{{
					ServiceName: spec.TypeName{
						Name:    "TestService",
						Package: "test.api",
					},
					Endpoints: []spec.EndpointDefinition{
						{
							EndpointName: "getFileSystems",
							HttpMethod:   spec.HttpMethod_GET.New(),
							HttpPath:     "/catalog/fileSystems",
							Auth:         authPtr(spec.NewAuthTypeFromHeader(spec.HeaderAuthType{})),
							Returns: specTypePtr(spec.NewTypeFromMap(spec.MapType{
								KeyType:   newPrimitive(spec.PrimitiveType_STRING),
								ValueType: newPrimitive(spec.PrimitiveType_INTEGER),
							})),
							Docs: docsPtr("Returns a mapping from file system id to backing file system configuration."),
						},
						{
							EndpointName: "createDataset",
							HttpMethod:   spec.HttpMethod_POST.New(),
							HttpPath:     "/catalog/datasets",
							Auth:         authPtr(spec.NewAuthTypeFromCookie(spec.CookieAuthType{CookieName: "PALANTIR_TOKEN"})),
							Args: []spec.ArgumentDefinition{
								{
									ArgName:   "request",
									Type:      newPrimitive(spec.PrimitiveType_STRING),
									ParamType: spec.NewParameterTypeFromBody(spec.BodyParameterType{}),
								},
							},
						},
						{
							EndpointName: "streamResponse",
							HttpMethod:   spec.HttpMethod_GET.New(),
							HttpPath:     "/catalog/streamResponse",
							Auth:         authPtr(spec.NewAuthTypeFromHeader(spec.HeaderAuthType{})),
							Returns:      specTypePtr(newPrimitive(spec.PrimitiveType_BINARY)),
						},
						{
							EndpointName: "maybeStreamResponse",
							HttpMethod:   spec.HttpMethod_GET.New(),
							HttpPath:     "/catalog/maybe/streamResponse",
							Auth:         authPtr(spec.NewAuthTypeFromHeader(spec.HeaderAuthType{})),
							Returns: specTypePtr(spec.NewTypeFromOptional(spec.OptionalType{
								ItemType: newPrimitive(spec.PrimitiveType_BINARY),
							})),
						},
						{
							EndpointName: "queryParams",
							HttpMethod:   spec.HttpMethod_GET.New(),
							HttpPath:     "/catalog/echo",
							Args: []spec.ArgumentDefinition{
								{
									ArgName:   "input",
									Type:      newPrimitive(spec.PrimitiveType_STRING),
									ParamType: spec.NewParameterTypeFromQuery(spec.QueryParameterType{ParamId: "input"}),
								},
								{
									ArgName:   "reps",
									Type:      newPrimitive(spec.PrimitiveType_INTEGER),
									ParamType: spec.NewParameterTypeFromQuery(spec.QueryParameterType{ParamId: "reps"}),
								},
							},
						},
					},
					Docs: docsPtr("A Markdown description of the service.\n"),
				}},
			},
			Out: &ConjureDefinition{
				Version: 1,
				Packages: map[string]ConjurePackage{
					"test.api": {
						ConjurePackage: "test.api",
						ImportPath:     "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
						OutputDir:      "test/test/api",
						PackageName:    "api",
						Services: []*ServiceDefinition{
							{
								Docs: "A Markdown description of the service.",
								Name: "TestService",
								Endpoints: []*EndpointDefinition{
									{
										Docs:         "Returns a mapping from file system id to backing file system configuration.",
										EndpointName: "getFileSystems",
										HTTPMethod:   spec.HttpMethod_GET.New(),
										HTTPPath:     "/catalog/fileSystems",
										HeaderAuth:   true,
										Returns:      typePtr(&Map{Key: String{}, Val: Integer{}}),
									},
									{
										EndpointName: "createDataset",
										HTTPMethod:   spec.HttpMethod_POST.New(),
										HTTPPath:     "/catalog/datasets",
										CookieAuth:   stringPtr("PALANTIR_TOKEN"),
										Params: []*EndpointArgumentDefinition{
											{
												Name:      "request",
												Type:      String{},
												ParamType: BodyParam,
												ParamID:   "request",
											},
										},
									},
									{
										EndpointName: "streamResponse",
										HTTPMethod:   spec.HttpMethod_GET.New(),
										HTTPPath:     "/catalog/streamResponse",
										HeaderAuth:   true,
										Returns:      typePtr(Binary{}),
									},
									{
										EndpointName: "maybeStreamResponse",
										HTTPMethod:   spec.HttpMethod_GET.New(),
										HTTPPath:     "/catalog/maybe/streamResponse",
										HeaderAuth:   true,
										Returns:      typePtr(&Optional{Item: Binary{}}),
									},
									{
										EndpointName: "queryParams",
										HTTPMethod:   spec.HttpMethod_GET.New(),
										HTTPPath:     "/catalog/echo",
										Params: []*EndpointArgumentDefinition{
											{
												Name:      "input",
												Type:      String{},
												ParamType: QueryParam,
												ParamID:   "input",
											},
											{
												Name:      "reps",
												Type:      Integer{},
												ParamType: QueryParam,
												ParamID:   "reps",
											},
										},
									},
								},
								conjurePkg: "test.api",
								importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
							},
						},
					},
				},
			},
		},
		{
			Name: "type and service definition",
			In: spec.ConjureDefinition{
				Version: 1,
				Types: []spec.TypeDefinition{
					spec.NewTypeDefinitionFromObject(spec.ObjectDefinition{
						TypeName: spec.TypeName{
							Name:    "BackingFileSystem",
							Package: "com.palantir.foundry.catalog.api.datasets",
						},
						Fields: []spec.FieldDefinition{
							{
								FieldName: "fileSystemId",
								Type:      newPrimitive(spec.PrimitiveType_STRING),
								Docs:      docsPtr("The name by which this file system is identified."),
							},
							{
								FieldName: "baseUri",
								Type:      newPrimitive(spec.PrimitiveType_STRING),
							},
						},
						Docs: docsPtr("Optional Docs"),
					}),
				},
				Services: []spec.ServiceDefinition{
					{
						ServiceName: spec.TypeName{
							Name:    "TestService",
							Package: "test.api",
						},
						Endpoints: []spec.EndpointDefinition{
							{
								EndpointName: "getFileSystems",
								HttpMethod:   spec.HttpMethod_GET.New(),
								HttpPath:     "/catalog/fileSystems",
								Auth:         authPtr(spec.NewAuthTypeFromHeader(spec.HeaderAuthType{})),
								Returns: specTypePtr(spec.NewTypeFromMap(spec.MapType{
									KeyType: newPrimitive(spec.PrimitiveType_STRING),
									ValueType: spec.NewTypeFromReference(spec.TypeName{
										Name:    "BackingFileSystem",
										Package: "com.palantir.foundry.catalog.api.datasets",
									}),
								})),
								Docs: docsPtr("Returns a mapping from file system id to backing file system configuration."),
							},
						},
						Docs: docsPtr("A Markdown description of the service.\n"),
					},
				},
			},
			Out: &ConjureDefinition{
				Version: 1,
				Packages: map[string]ConjurePackage{
					"com.palantir.foundry.catalog.api.datasets": {
						ConjurePackage: "com.palantir.foundry.catalog.api.datasets",
						ImportPath:     "github.com/palantir/conjure-go/v6/conjure/types/test/foundry/catalog/api/datasets",
						OutputDir:      "test/foundry/catalog/api/datasets",
						PackageName:    "datasets",
						Objects: []*ObjectType{
							{
								Name:       "BackingFileSystem",
								Docs:       "Optional Docs",
								conjurePkg: "com.palantir.foundry.catalog.api.datasets",
								importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/foundry/catalog/api/datasets",
								Fields: []*Field{
									{
										Name: "fileSystemId",
										Docs: "The name by which this file system is identified.",
										Type: String{},
									},
									{
										Name: "baseUri",
										Type: String{},
									},
								},
							},
						},
					},
					"test.api": {
						ConjurePackage: "test.api",
						ImportPath:     "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
						OutputDir:      "test/test/api",
						PackageName:    "api",
						Services: []*ServiceDefinition{
							{
								Docs: "A Markdown description of the service.",
								Name: "TestService",
								Endpoints: []*EndpointDefinition{
									{
										Docs:         "Returns a mapping from file system id to backing file system configuration.",
										EndpointName: "getFileSystems",
										HTTPMethod:   spec.HttpMethod_GET.New(),
										HTTPPath:     "/catalog/fileSystems",
										HeaderAuth:   true,
										Returns: typePtr(&Map{
											Key: String{},
											Val: &ObjectType{
												Name:       "BackingFileSystem",
												Docs:       "Optional Docs",
												conjurePkg: "com.palantir.foundry.catalog.api.datasets",
												importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/foundry/catalog/api/datasets",
												Fields: []*Field{
													{
														Name: "fileSystemId",
														Docs: "The name by which this file system is identified.",
														Type: String{},
													},
													{
														Name: "baseUri",
														Type: String{},
													},
												},
											},
										}),
									},
								},
								conjurePkg: "test.api",
								importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
							},
						},
					},
				},
				Extensions: nil,
			},
		},
		{
			Name: "type definition with kebab cases",
			In: spec.ConjureDefinition{
				Version: 1,
				Types: []spec.TypeDefinition{
					spec.NewTypeDefinitionFromObject(spec.ObjectDefinition{
						TypeName: spec.TypeName{
							Name:    "ServiceLogV1",
							Package: "com.palantir.sls.spec.logging",
						},
						Fields: []spec.FieldDefinition{
							{
								FieldName: "kebab-case",
								Type:      newPrimitive(spec.PrimitiveType_STRING),
							},
						},
					}),
				},
			},
			Out: &ConjureDefinition{
				Version: 1,
				Packages: map[string]ConjurePackage{"com.palantir.sls.spec.logging": {
					ConjurePackage: "com.palantir.sls.spec.logging",
					ImportPath:     "github.com/palantir/conjure-go/v6/conjure/types/test/sls/spec/logging",
					OutputDir:      "test/sls/spec/logging",
					PackageName:    "logging",
					Objects: []*ObjectType{{
						Name: "ServiceLogV1",
						Fields: []*Field{{
							Name: "kebab-case",
							Type: String{},
						}},
						conjurePkg: "com.palantir.sls.spec.logging",
						importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/sls/spec/logging",
					}},
				}},
			},
		},
		{
			Name: "full-featured error definition",
			In: spec.ConjureDefinition{
				Version: 1,
				Errors: []spec.ErrorDefinition{{
					ErrorName: spec.TypeName{
						Name:    "MyNotFound",
						Package: "com.palantir.test.another.api",
					},
					Docs:      docsPtr("This is documentation of MyNotFound error."),
					Namespace: "MyNamespace",
					Code:      spec.ErrorCode_NOT_FOUND.New(),
					SafeArgs: []spec.FieldDefinition{
						{
							FieldName: "safeArgA",
							Type:      spec.NewTypeFromReference(spec.TypeName{Name: "SimpleObject", Package: "com.palantir.test.api"}),
							Docs:      docsPtr("This is safeArgA doc."),
						},
						{
							FieldName: "safeArgB",
							Type:      newPrimitive(spec.PrimitiveType_INTEGER),
						},
					},
					UnsafeArgs: []spec.FieldDefinition{
						{
							FieldName: "unsafeArgA",
							Type:      newPrimitive(spec.PrimitiveType_STRING),
							Docs:      docsPtr("This is unsafeArgA doc."),
						},
					},
				}},
				Types: []spec.TypeDefinition{
					spec.NewTypeDefinitionFromObject(spec.ObjectDefinition{
						TypeName: spec.TypeName{
							Name:    "SimpleObject",
							Package: "com.palantir.test.api",
						},
						Fields: []spec.FieldDefinition{{
							FieldName: "someField",
							Type:      newPrimitive(spec.PrimitiveType_STRING),
						}},
					}),
				},
			},
			Out: &ConjureDefinition{
				Version: 1,
				Packages: map[string]ConjurePackage{
					"com.palantir.test.api": {
						ConjurePackage: "com.palantir.test.api",
						ImportPath:     "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
						OutputDir:      "test/test/api",
						PackageName:    "api",
						Objects: []*ObjectType{{
							Name: "SimpleObject",
							Fields: []*Field{{
								Name: "someField",
								Type: String{},
							}},
							conjurePkg: "com.palantir.test.api",
							importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
						}},
					},
					"com.palantir.test.another.api": {
						ConjurePackage: "com.palantir.test.another.api",
						ImportPath:     "github.com/palantir/conjure-go/v6/conjure/types/test/test/another/api",
						OutputDir:      "test/test/another/api",
						PackageName:    "api",
						Errors: []*ErrorDefinition{{
							Docs:           "This is documentation of MyNotFound error.",
							Name:           "MyNotFound",
							ErrorNamespace: "MyNamespace",
							ErrorCode:      spec.ErrorCode_NOT_FOUND.New(),
							SafeArgs: []*Field{
								{
									Docs: "This is safeArgA doc.",
									Name: "safeArgA",
									Type: &ObjectType{
										Name: "SimpleObject",
										Fields: []*Field{{
											Name: "someField",
											Type: String{},
										}},
										conjurePkg: "com.palantir.test.api",
										importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/api",
									},
								},
								{
									Name: "safeArgB",
									Type: Integer{},
								},
							},
							UnsafeArgs: []*Field{
								{
									Name: "unsafeArgA",
									Docs: "This is unsafeArgA doc.",
									Type: String{},
								},
							},
							conjurePkg: "com.palantir.test.another.api",
							importPath: "github.com/palantir/conjure-go/v6/conjure/types/test/test/another/api",
						}},
					},
				},
			},
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			out, err := NewConjureDefinition(outputDir, test.In)
			require.NoError(t, err)
			require.NotNil(t, out)
			require.Equal(t, test.Out, out)
		})
	}
}

func TestRecursiveTypeDefinition(t *testing.T) {
	def := spec.ConjureDefinition{
		Version: 1,
		Types: []spec.TypeDefinition{
			spec.NewTypeDefinitionFromAlias(spec.AliasDefinition{
				TypeName: spec.TypeName{
					Package: "com.palantir.test",
					Name:    "RecursiveMap",
				},
				Alias: spec.NewTypeFromMap(spec.MapType{
					KeyType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_STRING)),
					ValueType: spec.NewTypeFromReference(spec.TypeName{
						Package: "com.palantir.test",
						Name:    "RecursiveMap",
					}),
				}),
			}),
			spec.NewTypeDefinitionFromUnion(spec.UnionDefinition{
				TypeName: spec.TypeName{
					Package: "com.palantir.test",
					Name:    "RecursiveUnion",
				},
				Union: []spec.FieldDefinition{
					{
						FieldName: "list",
						Type: spec.NewTypeFromList(spec.ListType{
							ItemType: spec.NewTypeFromReference(spec.TypeName{
								Package: "com.palantir.test",
								Name:    "RecursiveUnion",
							}),
						}),
					},
				},
			}),
		},
	}

	mapAlias := &AliasType{
		Name:       "RecursiveMap",
		conjurePkg: "com.palantir.test",
		importPath: "github.com/palantir/conjure-go/v6/conjure/types/com/palantir/test",
	}
	mapAlias.Item = &Map{
		Key: String{},
		Val: mapAlias,
	}

	unionType := &UnionType{
		Name:       "RecursiveUnion",
		Fields:     nil, // populated below
		conjurePkg: "com.palantir.test",
		importPath: "github.com/palantir/conjure-go/v6/conjure/types/com/palantir/test",
	}
	unionType.Fields = []*Field{
		{
			Name: "list",
			Type: &List{Item: unionType},
		},
	}

	want := &ConjureDefinition{
		Version: 1,
		Packages: map[string]ConjurePackage{
			"com.palantir.test": {
				ImportPath:     "github.com/palantir/conjure-go/v6/conjure/types/com/palantir/test",
				OutputDir:      "com/palantir/test",
				PackageName:    "test",
				ConjurePackage: "com.palantir.test",
				Aliases:        []*AliasType{mapAlias},
				Unions:         []*UnionType{unionType},
			},
		},
	}

	out, err := NewConjureDefinition(".", def)
	require.NoError(t, err)
	require.Equal(t, want, out)
}

func TestNewConjureDefinition_ConjureAPI(t *testing.T) {
	apiBody, err := ioutil.ReadFile("../../conjure-api/conjure-api-4.35.0.conjure.json")
	require.NoError(t, err)
	var inputDef spec.ConjureDefinition
	require.NoError(t, inputDef.UnmarshalJSON(apiBody))
	out, err := NewConjureDefinition("./test", inputDef)
	require.NoError(t, err)
	require.NotNil(t, out)
	t.Logf("%#v", out)
}

func TestNewConjureDefinition_Verifier(t *testing.T) {
	apiBody, err := ioutil.ReadFile("../../conjure-go-verifier/verification-server-api.conjure.json")
	require.NoError(t, err)
	var inputDef spec.ConjureDefinition
	require.NoError(t, inputDef.UnmarshalJSON(apiBody))
	out, err := NewConjureDefinition("./test", inputDef)
	require.NoError(t, err)
	require.NotNil(t, out)
	t.Logf("%#v", out)
}

func TestSanitizePackageName(t *testing.T) {
	for _, test := range []struct {
		Import, Name string
	}{
		{"foo", "foo"},
		{"foo/bar", "bar"},
		{"foo/bar.2", "bar2"},
		{"foo/2bar", "bar"},
		{"foo/2_bar", "bar"},
		{"foo/123", "pkg"},
	} {
		t.Run(test.Import, func(t *testing.T) {
			assert.Equal(t, test.Name, sanitizePackageName(test.Import))
		})
	}
}

func newPrimitive(kind spec.PrimitiveType_Value) spec.Type {
	return spec.NewTypeFromPrimitive(kind.New())
}

func stringPtr(s string) *string             { return &s }
func authPtr(a spec.AuthType) *spec.AuthType { return &a }
func specTypePtr(t spec.Type) *spec.Type     { return &t }
func typePtr(t Type) *Type                   { return &t }
func docsPtr(s string) *spec.Documentation   { return (*spec.Documentation)(&s) }
