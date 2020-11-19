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

package conjure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
)

var testCases = []struct {
	name      string
	src       string
	wantFiles map[string]string
}{
	{
		name: "full-featured object definition",
		src: `
			{
				"version" : 1,
				"errors" : [ ],
				"types" : [ {
				  "type" : "object",
				  "object" : {
					"typeName" : {
					  "name" : "BackingFileSystem",
					  "package" : "com.palantir.foundry.catalog.api.datasets"
					},
					"fields" : [ {
					  "fieldName" : "fileSystemId",
					  "type" : {
						"type" : "primitive",
						"primitive" : "STRING"
					  },
					  "docs" : "The name by which this file system is identified."
					}, {
					  "fieldName" : "baseUri",
					  "type" : {
						"type" : "primitive",
						"primitive" : "STRING"
					  }
					}, {
					  "fieldName" : "exenum",
					  "type" : {
						"type" : "reference",
						"reference" : {
						  "name" : "ExampleEnumeration",
						  "package" : "example.api"
						}
					  }
					}, {
					  "fieldName" : "client",
					  "type" : {
						"type" : "external",
						"external" : {
						  "externalReference" : {
							"name" : "com/palantir/go-palantir/httpclient:RESTClient",
							"package" : "github"
						  },
						  "fallback" : {
							"type" : "primitive",
							"primitive" : "STRING"
						  }
						}
					  }
					} ],
					"docs" : "Optional Docs"
				  }
				}, {
				  "type" : "object",
				  "object" : {
					"typeName" : {
					  "name" : "TestType",
					  "package" : "com.palantir.foundry.catalog.api.datasets"
					},
					"fields" : [ {
					  "fieldName" : "alias",
					  "type" : {
						"type" : "reference",
						"reference" : {
						  "name" : "ExampleAlias",
						  "package" : "com.palantir.test.api"
						}
					  }
					}, {
					  "fieldName" : "rid",
					  "type" : {
						"type" : "primitive",
						"primitive" : "RID"
					  }
					}, {
					  "fieldName" : "large_int",
					  "type" : {
						"type" : "primitive",
						"primitive" : "SAFELONG"
					  }
					}, {
					  "fieldName" : "time",
					  "type" : {
						"type" : "primitive",
						"primitive" : "DATETIME"
					  }
					}, {
					  "fieldName" : "bytes",
					  "type" : {
						"type" : "primitive",
						"primitive" : "BINARY"
					  }
					} ]
				  }
				}, {
				  "type" : "enum",
				  "enum" : {
					"typeName" : {
					  "name" : "ExampleEnumeration",
					  "package" : "example.api"
					},
					"values" : [ {
					  "value" : "A"
					}, {
					  "value" : "B"
					} ]
				  }
				}, {
				  "type" : "enum",
				  "enum" : {
					"typeName" : {
					  "name" : "Months",
					  "package" : "com.palantir.test.api"
					},
					"values" : [ {
					  "value" : "JANUARY"
					}, {
					  "value" : "MULTI_MONTHS"
					} ]
				  }
				}, {
				  "type" : "enum",
				  "enum" : {
					"typeName" : {
					  "name" : "Days",
					  "package" : "com.palantir.test.api"
					},
					"values" : [ {
					  "value" : "FRIDAY"
					}, {
					  "value" : "SATURDAY"
					} ]
				  }
				}, {
				  "type" : "alias",
				  "alias" : {
					"typeName" : {
					  "name" : "ExampleAlias",
					  "package" : "com.palantir.test.api"
					},
					"alias" : {
					  "type" : "primitive",
					  "primitive" : "STRING"
					}
				  }
				}, {
				  "type" : "alias",
				  "alias" : {
					"typeName" : {
					  "name" : "LongAlias",
					  "package" : "com.palantir.test.api"
					},
					"alias" : {
					  "type" : "primitive",
					  "primitive" : "SAFELONG"
					}
				  }
				}, {
				  "type" : "alias",
				  "alias" : {
					"typeName" : {
					  "name" : "Status",
					  "package" : "com.palantir.test.api"
					},
					"alias" : {
					  "type" : "primitive",
					  "primitive" : "INTEGER"
					}
				  }
				}, {
				  "type" : "alias",
				  "alias" : {
					"typeName" : {
					  "name" : "ObjectAlias",
					  "package" : "com.palantir.test.api"
					},
					"alias" : {
					  "type" : "reference",
					  "reference" : {
						"name" : "TestType",
						"package" : "com.palantir.foundry.catalog.api.datasets"
					  }
					}
				  }
				}, {
				  "type" : "alias",
				  "alias" : {
					"typeName" : {
					  "name" : "MapAlias",
					  "package" : "com.palantir.test.api"
					},
					"alias" : {
					  "type" : "map",
					  "map" : {
						"keyType" : {
						  "type" : "primitive",
						  "primitive" : "STRING"
						},
						"valueType" : {
						  "type" : "reference",
						  "reference" : {
							"name" : "Status",
							"package" : "com.palantir.test.api"
						  }
						}
					  }
					}
				  }
				}, {
				  "type" : "alias",
				  "alias" : {
					"typeName" : {
					  "name" : "AliasAlias",
					  "package" : "com.palantir.test.api"
					},
					"alias" : {
					  "type" : "reference",
					  "reference" : {
						"name" : "Status",
						"package" : "com.palantir.test.api"
					  }
					}
				  }
				}, {
				  "type" : "union",
				  "union" : {
					"typeName" : {
					  "name" : "ExampleUnion",
					  "package" : "com.palantir.test.api"
					},
					"union" : [ {
					  "fieldName" : "str",
					  "type" : {
						"type" : "primitive",
						"primitive" : "STRING"
					  }
					}, {
					  "fieldName" : "other",
					  "type" : {
						"type" : "primitive",
						"primitive" : "STRING"
					  },
					  "docs" : "Another string"
					}, {
					  "fieldName" : "myMap",
					  "type" : {
						"type" : "map",
						"map" : {
						  "keyType" : {
							"type" : "primitive",
							"primitive" : "STRING"
						  },
						  "valueType" : {
							"type" : "list",
							"list" : {
							  "itemType" : {
								"type" : "primitive",
								"primitive" : "INTEGER"
							  }
							}
						  }
						}
					  }
					}, {
					  "fieldName" : "tester",
					  "type" : {
						"type" : "reference",
						"reference" : {
						  "name" : "TestType",
						  "package" : "com.palantir.foundry.catalog.api.datasets"
						}
					  }
					}, {
					  "fieldName" : "recursive",
					  "type" : {
						"type" : "reference",
						"reference" : {
						  "name" : "ExampleUnion",
						  "package" : "com.palantir.test.api"
						}
					  }
					} ]
				  }
				} ],
				"services" : [ ]
}
`,
		wantFiles: map[string]string{
			"example/api/enums.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"strings"
)

type ExampleEnumeration struct {
	val ExampleEnumeration_Value
}

type ExampleEnumeration_Value string

const (
	ExampleEnumeration_A       ExampleEnumeration_Value = "A"
	ExampleEnumeration_B       ExampleEnumeration_Value = "B"
	ExampleEnumeration_UNKNOWN ExampleEnumeration_Value = "UNKNOWN"
)

// ExampleEnumeration_Values returns all known variants of ExampleEnumeration.
func ExampleEnumeration_Values() []ExampleEnumeration_Value {
	return []ExampleEnumeration_Value{ExampleEnumeration_A, ExampleEnumeration_B}
}

func New_ExampleEnumeration(value ExampleEnumeration_Value) ExampleEnumeration {
	return ExampleEnumeration{val: value}
}

// IsUnknown returns false for all known variants of ExampleEnumeration and true otherwise.
func (e ExampleEnumeration) IsUnknown() bool {
	switch e.val {
	case ExampleEnumeration_A, ExampleEnumeration_B:
		return false
	}
	return true
}

func (e ExampleEnumeration) Value() ExampleEnumeration_Value {
	if e.IsUnknown() {
		return ExampleEnumeration_UNKNOWN
	}
	return e.val
}

func (e ExampleEnumeration) String() string {
	return string(e.val)
}

func (e ExampleEnumeration) MarshalText() ([]byte, error) {
	return []byte(e.val), nil
}

func (e *ExampleEnumeration) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		*e = New_ExampleEnumeration(ExampleEnumeration_Value(v))
	case "A":
		*e = New_ExampleEnumeration(ExampleEnumeration_A)
	case "B":
		*e = New_ExampleEnumeration(ExampleEnumeration_B)
	}
	return nil
}
`,
			"foundry/catalog/api/datasets/structs.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package datasets

import (
	"github.com/palantir/conjure-go/v6/conjure/{{currCaseTmpDir}}/example/api"
	api_1 "github.com/palantir/conjure-go/v6/conjure/{{currCaseTmpDir}}/test/api"
	"github.com/palantir/go-palantir/httpclient"
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/pkg/safeyaml"
)

// Optional Docs
type BackingFileSystem struct {
	// The name by which this file system is identified.
	FileSystemId string                 ` + "`json:\"fileSystemId\" conjure-docs:\"The name by which this file system is identified.\"`" + `
	BaseUri      string                 ` + "`json:\"baseUri\"`" + `
	Exenum       api.ExampleEnumeration ` + "`json:\"exenum\"`" + `
	Client       httpclient.RESTClient  ` + "`json:\"client\"`" + `
}

func (o BackingFileSystem) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *BackingFileSystem) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type TestType struct {
	Alias    api_1.ExampleAlias     ` + "`json:\"alias\"`" + `
	Rid      rid.ResourceIdentifier ` + "`json:\"rid\"`" + `
	LargeInt safelong.SafeLong      ` + "`json:\"large_int\"`" + `
	Time     datetime.DateTime      ` + "`json:\"time\"`" + `
	Bytes    []byte                 ` + "`json:\"bytes\"`" + `
}

func (o TestType) MarshalJSON() ([]byte, error) {
	if o.Bytes == nil {
		o.Bytes = make([]byte, 0)
	}
	type TestTypeAlias TestType
	return safejson.Marshal(TestTypeAlias(o))
}

func (o *TestType) UnmarshalJSON(data []byte) error {
	type TestTypeAlias TestType
	var rawTestType TestTypeAlias
	if err := safejson.Unmarshal(data, &rawTestType); err != nil {
		return err
	}
	if rawTestType.Bytes == nil {
		rawTestType.Bytes = make([]byte, 0)
	}
	*o = TestType(rawTestType)
	return nil
}

func (o TestType) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *TestType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
			"test/api/aliases.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"github.com/palantir/conjure-go/v6/conjure/{{currCaseTmpDir}}/foundry/catalog/api/datasets"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/pkg/safeyaml"
)

type ExampleAlias string
type LongAlias safelong.SafeLong

func (a LongAlias) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(safelong.SafeLong(a))
}

func (a *LongAlias) UnmarshalJSON(data []byte) error {
	var rawLongAlias safelong.SafeLong
	if err := safejson.Unmarshal(data, &rawLongAlias); err != nil {
		return err
	}
	*a = LongAlias(rawLongAlias)
	return nil
}

func (a LongAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *LongAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type Status int
type ObjectAlias datasets.TestType

func (a ObjectAlias) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(datasets.TestType(a))
}

func (a *ObjectAlias) UnmarshalJSON(data []byte) error {
	var rawObjectAlias datasets.TestType
	if err := safejson.Unmarshal(data, &rawObjectAlias); err != nil {
		return err
	}
	*a = ObjectAlias(rawObjectAlias)
	return nil
}

func (a ObjectAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *ObjectAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type MapAlias map[string]Status

func (a MapAlias) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(map[string]Status(a))
}

func (a *MapAlias) UnmarshalJSON(data []byte) error {
	var rawMapAlias map[string]Status
	if err := safejson.Unmarshal(data, &rawMapAlias); err != nil {
		return err
	}
	*a = MapAlias(rawMapAlias)
	return nil
}

func (a MapAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *MapAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type AliasAlias Status

func (a AliasAlias) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(Status(a))
}

func (a *AliasAlias) UnmarshalJSON(data []byte) error {
	var rawAliasAlias Status
	if err := safejson.Unmarshal(data, &rawAliasAlias); err != nil {
		return err
	}
	*a = AliasAlias(rawAliasAlias)
	return nil
}

func (a AliasAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *AliasAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}
`,
			"test/api/enums.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"strings"
)

type Months struct {
	val Months_Value
}

type Months_Value string

const (
	Months_JANUARY      Months_Value = "JANUARY"
	Months_MULTI_MONTHS Months_Value = "MULTI_MONTHS"
	Months_UNKNOWN      Months_Value = "UNKNOWN"
)

// Months_Values returns all known variants of Months.
func Months_Values() []Months_Value {
	return []Months_Value{Months_JANUARY, Months_MULTI_MONTHS}
}

func New_Months(value Months_Value) Months {
	return Months{val: value}
}

// IsUnknown returns false for all known variants of Months and true otherwise.
func (e Months) IsUnknown() bool {
	switch e.val {
	case Months_JANUARY, Months_MULTI_MONTHS:
		return false
	}
	return true
}

func (e Months) Value() Months_Value {
	if e.IsUnknown() {
		return Months_UNKNOWN
	}
	return e.val
}

func (e Months) String() string {
	return string(e.val)
}

func (e Months) MarshalText() ([]byte, error) {
	return []byte(e.val), nil
}

func (e *Months) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		*e = New_Months(Months_Value(v))
	case "JANUARY":
		*e = New_Months(Months_JANUARY)
	case "MULTI_MONTHS":
		*e = New_Months(Months_MULTI_MONTHS)
	}
	return nil
}

type Days struct {
	val Days_Value
}

type Days_Value string

const (
	Days_FRIDAY   Days_Value = "FRIDAY"
	Days_SATURDAY Days_Value = "SATURDAY"
	Days_UNKNOWN  Days_Value = "UNKNOWN"
)

// Days_Values returns all known variants of Days.
func Days_Values() []Days_Value {
	return []Days_Value{Days_FRIDAY, Days_SATURDAY}
}

func New_Days(value Days_Value) Days {
	return Days{val: value}
}

// IsUnknown returns false for all known variants of Days and true otherwise.
func (e Days) IsUnknown() bool {
	switch e.val {
	case Days_FRIDAY, Days_SATURDAY:
		return false
	}
	return true
}

func (e Days) Value() Days_Value {
	if e.IsUnknown() {
		return Days_UNKNOWN
	}
	return e.val
}

func (e Days) String() string {
	return string(e.val)
}

func (e Days) MarshalText() ([]byte, error) {
	return []byte(e.val), nil
}

func (e *Days) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		*e = New_Days(Days_Value(v))
	case "FRIDAY":
		*e = New_Days(Days_FRIDAY)
	case "SATURDAY":
		*e = New_Days(Days_SATURDAY)
	}
	return nil
}
`,
			"test/api/unions.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"fmt"

	"github.com/palantir/conjure-go/v6/conjure/{{currCaseTmpDir}}/foundry/catalog/api/datasets"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type ExampleUnion struct {
	typ       string
	str       *string
	other     *string
	myMap     *map[string][]int
	tester    *datasets.TestType
	recursive *ExampleUnion
}

type exampleUnionDeserializer struct {
	Type      string             ` + "`" + `json:"type"` + "`" + `
	Str       *string            ` + "`" + `json:"str"` + "`" + `
	Other     *string            ` + "`" + `json:"other"` + "`" + `
	MyMap     *map[string][]int  ` + "`" + `json:"myMap"` + "`" + `
	Tester    *datasets.TestType ` + "`" + `json:"tester"` + "`" + `
	Recursive *ExampleUnion      ` + "`" + `json:"recursive"` + "`" + `
}

func (u *exampleUnionDeserializer) toStruct() ExampleUnion {
	return ExampleUnion{typ: u.Type, str: u.Str, other: u.Other, myMap: u.MyMap, tester: u.Tester, recursive: u.Recursive}
}

func (u *ExampleUnion) toSerializer() (interface{}, error) {
	switch u.typ {
	default:
		return nil, fmt.Errorf("unknown type %s", u.typ)
	case "str":
		return struct {
			Type string ` + "`" + `json:"type"` + "`" + `
			Str  string ` + "`" + `json:"str"` + "`" + `
		}{Type: "str", Str: *u.str}, nil
	case "other":
		return struct {
			Type  string ` + "`" + `json:"type"` + "`" + `
			Other string ` + "`" + `json:"other"` + "`" + `
		}{Type: "other", Other: *u.other}, nil
	case "myMap":
		return struct {
			Type  string           ` + "`" + `json:"type"` + "`" + `
			MyMap map[string][]int ` + "`" + `json:"myMap"` + "`" + `
		}{Type: "myMap", MyMap: *u.myMap}, nil
	case "tester":
		return struct {
			Type   string            ` + "`" + `json:"type"` + "`" + `
			Tester datasets.TestType ` + "`" + `json:"tester"` + "`" + `
		}{Type: "tester", Tester: *u.tester}, nil
	case "recursive":
		return struct {
			Type      string       ` + "`" + `json:"type"` + "`" + `
			Recursive ExampleUnion ` + "`" + `json:"recursive"` + "`" + `
		}{Type: "recursive", Recursive: *u.recursive}, nil
	}
}

func (u ExampleUnion) MarshalJSON() ([]byte, error) {
	ser, err := u.toSerializer()
	if err != nil {
		return nil, err
	}
	return safejson.Marshal(ser)
}

func (u *ExampleUnion) UnmarshalJSON(data []byte) error {
	var deser exampleUnionDeserializer
	if err := safejson.Unmarshal(data, &deser); err != nil {
		return err
	}
	*u = deser.toStruct()
	return nil
}

func (u ExampleUnion) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(u)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (u *ExampleUnion) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&u)
}

func (u *ExampleUnion) Accept(v ExampleUnionVisitor) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(u.typ)
	case "str":
		return v.VisitStr(*u.str)
	case "other":
		return v.VisitOther(*u.other)
	case "myMap":
		return v.VisitMyMap(*u.myMap)
	case "tester":
		return v.VisitTester(*u.tester)
	case "recursive":
		return v.VisitRecursive(*u.recursive)
	}
}

type ExampleUnionVisitor interface {
	VisitStr(v string) error
	VisitOther(v string) error
	VisitMyMap(v map[string][]int) error
	VisitTester(v datasets.TestType) error
	VisitRecursive(v ExampleUnion) error
	VisitUnknown(typeName string) error
}

func (u *ExampleUnion) AcceptWithContext(ctx context.Context, v ExampleUnionVisitorWithContext) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknownWithContext(ctx, u.typ)
	case "str":
		return v.VisitStrWithContext(ctx, *u.str)
	case "other":
		return v.VisitOtherWithContext(ctx, *u.other)
	case "myMap":
		return v.VisitMyMapWithContext(ctx, *u.myMap)
	case "tester":
		return v.VisitTesterWithContext(ctx, *u.tester)
	case "recursive":
		return v.VisitRecursiveWithContext(ctx, *u.recursive)
	}
}

type ExampleUnionVisitorWithContext interface {
	VisitStrWithContext(ctx context.Context, v string) error
	VisitOtherWithContext(ctx context.Context, v string) error
	VisitMyMapWithContext(ctx context.Context, v map[string][]int) error
	VisitTesterWithContext(ctx context.Context, v datasets.TestType) error
	VisitRecursiveWithContext(ctx context.Context, v ExampleUnion) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
}

func NewExampleUnionFromStr(v string) ExampleUnion {
	return ExampleUnion{typ: "str", str: &v}
}

func NewExampleUnionFromOther(v string) ExampleUnion {
	return ExampleUnion{typ: "other", other: &v}
}

func NewExampleUnionFromMyMap(v map[string][]int) ExampleUnion {
	return ExampleUnion{typ: "myMap", myMap: &v}
}

func NewExampleUnionFromTester(v datasets.TestType) ExampleUnion {
	return ExampleUnion{typ: "tester", tester: &v}
}

func NewExampleUnionFromRecursive(v ExampleUnion) ExampleUnion {
	return ExampleUnion{typ: "recursive", recursive: &v}
}
`,
		},
	},
	{
		name: "full-featured service definition",
		src: `
{
	"version" : 1,
	"errors" : [ ],
	"types" : [ ],
	"services" : [ {
		"serviceName" : {
		"name" : "TestService",
		"package" : "test.api"
		},
		"endpoints" : [ {
		"endpointName" : "getFileSystems",
		"httpMethod" : "GET",
		"httpPath" : "/catalog/fileSystems",
		"auth" : {
			"type" : "header",
			"header" : { }
		},
		"args" : [ ],
		"returns" : {
			"type" : "map",
			"map" : {
			"keyType" : {
				"type" : "primitive",
				"primitive" : "STRING"
			},
			"valueType" : {
				"type" : "primitive",
				"primitive" : "INTEGER"
			}
			}
		},
		"docs" : "Returns a mapping from file system id to backing file system configuration.\n",
		"markers" : [ ]
		}, {
		"endpointName" : "createDataset",
		"httpMethod" : "POST",
		"httpPath" : "/catalog/datasets",
		"auth" : {
			"type" : "cookie",
			"cookie" : {
			"cookieName" : "PALANTIR_TOKEN"
			}
		},
		"args" : [ {
			"argName" : "request",
			"type" : {
			"type" : "primitive",
			"primitive" : "STRING"
			},
			"paramType" : {
			"type" : "body",
			"body" : { }
			},
			"markers" : [ ]
		} ],
		"markers" : [ ]
		}, {
		"endpointName" : "streamResponse",
		"httpMethod" : "GET",
		"httpPath" : "/catalog/streamResponse",
		"auth" : {
			"type" : "header",
			"header" : { }
		},
		"args" : [ ],
		"returns" : {
			"type" : "primitive",
			"primitive" : "BINARY"
		},
		"markers" : [ ]
		}, {
		"endpointName" : "queryParams",
		"httpMethod" : "GET",
		"httpPath" : "/catalog/echo",
		"args" : [ {
			"argName" : "input",
			"type" : {
			"type" : "primitive",
			"primitive" : "STRING"
			},
			"paramType" : {
			"type" : "query",
			"query" : {
				"paramId" : "input"
			}
			},
			"markers" : [ ]
		}, {
			"argName" : "reps",
			"type" : {
			"type" : "primitive",
			"primitive" : "INTEGER"
			},
			"paramType" : {
			"type" : "query",
			"query" : {
				"paramId" : "reps"
			}
			},
			"markers" : [ ]
		} ],
		"markers" : [ ]
		} ],
		"docs" : "A Markdown description of the service.\n"
	} ]
	}
`,
		wantFiles: map[string]string{
			"test/api/services.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/pkg/bearertoken"
)

// A Markdown description of the service.
type TestServiceClient interface {
	// Returns a mapping from file system id to backing file system configuration.
	GetFileSystems(ctx context.Context, authHeader bearertoken.Token) (map[string]int, error)
	CreateDataset(ctx context.Context, cookieToken bearertoken.Token, requestArg string) error
	StreamResponse(ctx context.Context, authHeader bearertoken.Token) (io.ReadCloser, error)
	QueryParams(ctx context.Context, inputArg string, repsArg int) error
}

type testServiceClient struct {
	client httpclient.Client
}

func NewTestServiceClient(client httpclient.Client) TestServiceClient {
	return &testServiceClient{client: client}
}

func (c *testServiceClient) GetFileSystems(ctx context.Context, authHeader bearertoken.Token) (map[string]int, error) {
	var returnVal map[string]int
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetFileSystems"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/fileSystems"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return returnVal, err
	}
	_ = resp
	if returnVal == nil {
		returnVal = make(map[string]int, 0)
	}
	return returnVal, nil
}

func (c *testServiceClient) CreateDataset(ctx context.Context, cookieToken bearertoken.Token, requestArg string) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("CreateDataset"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithHeader("Cookie", fmt.Sprint("PALANTIR_TOKEN=", cookieToken)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/datasets"))
	requestParams = append(requestParams, httpclient.WithJSONRequest(requestArg))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return err
	}
	_ = resp
	return nil
}

func (c *testServiceClient) StreamResponse(ctx context.Context, authHeader bearertoken.Token) (io.ReadCloser, error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("StreamResponse"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/streamResponse"))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *testServiceClient) QueryParams(ctx context.Context, inputArg string, repsArg int) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParams"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/echo"))
	queryParams := make(url.Values)
	queryParams.Set("input", fmt.Sprint(inputArg))
	queryParams.Set("reps", fmt.Sprint(repsArg))
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return err
	}
	_ = resp
	return nil
}

// A Markdown description of the service.
type TestServiceClientWithAuth interface {
	// Returns a mapping from file system id to backing file system configuration.
	GetFileSystems(ctx context.Context) (map[string]int, error)
	CreateDataset(ctx context.Context, requestArg string) error
	StreamResponse(ctx context.Context) (io.ReadCloser, error)
	QueryParams(ctx context.Context, inputArg string, repsArg int) error
}

func NewTestServiceClientWithAuth(client TestServiceClient, authHeader bearertoken.Token, cookieToken bearertoken.Token) TestServiceClientWithAuth {
	return &testServiceClientWithAuth{client: client, authHeader: authHeader, cookieToken: cookieToken}
}

type testServiceClientWithAuth struct {
	client      TestServiceClient
	authHeader  bearertoken.Token
	cookieToken bearertoken.Token
}

func (c *testServiceClientWithAuth) GetFileSystems(ctx context.Context) (map[string]int, error) {
	return c.client.GetFileSystems(ctx, c.authHeader)
}

func (c *testServiceClientWithAuth) CreateDataset(ctx context.Context, requestArg string) error {
	return c.client.CreateDataset(ctx, c.cookieToken, requestArg)
}

func (c *testServiceClientWithAuth) StreamResponse(ctx context.Context) (io.ReadCloser, error) {
	return c.client.StreamResponse(ctx, c.authHeader)
}

func (c *testServiceClientWithAuth) QueryParams(ctx context.Context, inputArg string, repsArg int) error {
	return c.client.QueryParams(ctx, inputArg, repsArg)
}
`,
		},
	},
	{
		name: "service definition without auth",
		src: `
{
	"version" : 1,
	"errors" : [ ],
	"types" : [ ],
	"services" : [ {
		"serviceName" : {
		"name" : "TestService",
		"package" : "test.api"
		},
		"endpoints" : [ {
		"endpointName" : "getFileSystems",
		"httpMethod" : "GET",
		"httpPath" : "/catalog/fileSystems",
		"args" : [ ],
		"returns" : {
			"type" : "map",
			"map" : {
			"keyType" : {
				"type" : "primitive",
				"primitive" : "STRING"
			},
			"valueType" : {
				"type" : "primitive",
				"primitive" : "INTEGER"
			}
			}
		},
		"docs" : "Returns a mapping from file system id to backing file system configuration.",
		"markers" : [ ]
		} ],
		"docs" : "A Markdown description of the service.\n"
	} ]
}
`,
		wantFiles: map[string]string{
			"test/api/services.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
)

// A Markdown description of the service.
type TestServiceClient interface {
	// Returns a mapping from file system id to backing file system configuration.
	GetFileSystems(ctx context.Context) (map[string]int, error)
}

type testServiceClient struct {
	client httpclient.Client
}

func NewTestServiceClient(client httpclient.Client) TestServiceClient {
	return &testServiceClient{client: client}
}

func (c *testServiceClient) GetFileSystems(ctx context.Context) (map[string]int, error) {
	var returnVal map[string]int
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetFileSystems"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/fileSystems"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return returnVal, err
	}
	_ = resp
	if returnVal == nil {
		returnVal = make(map[string]int, 0)
	}
	return returnVal, nil
}
`,
		},
	},
	{
		name: "type and service definition",
		src: `
{
	"version" : 1,
	"errors" : [ ],
	"types" : [ {
		"type" : "object",
		"object" : {
		"typeName" : {
			"name" : "BackingFileSystem",
			"package" : "com.palantir.foundry.catalog.api.datasets"
		},
		"fields" : [ {
			"fieldName" : "fileSystemId",
			"type" : {
			"type" : "primitive",
			"primitive" : "STRING"
			},
			"docs" : "The name by which this file system is identified."
		}, {
			"fieldName" : "baseUri",
			"type" : {
			"type" : "primitive",
			"primitive" : "STRING"
			}
		} ],
		"docs" : "Optional Docs"
		}
	} ],
	"services" : [ {
		"serviceName" : {
		"name" : "TestService",
		"package" : "test.api"
		},
		"endpoints" : [ {
		"endpointName" : "getFileSystems",
		"httpMethod" : "GET",
		"httpPath" : "/catalog/fileSystems",
		"auth" : {
			"type" : "header",
			"header" : { }
		},
		"args" : [ ],
		"returns" : {
			"type" : "map",
			"map" : {
			"keyType" : {
				"type" : "primitive",
				"primitive" : "STRING"
			},
			"valueType" : {
				"type" : "reference",
				"reference" : {
				"name" : "BackingFileSystem",
				"package" : "com.palantir.foundry.catalog.api.datasets"
				}
			}
			}
		},
		"docs" : "Returns a mapping from file system id to backing file system configuration.",
		"markers" : [ ]
		} ],
		"docs" : "A Markdown description of the service.\n"
	} ]
	}
`,
		wantFiles: map[string]string{
			"foundry/catalog/api/datasets/structs.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package datasets

import (
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

// Optional Docs
type BackingFileSystem struct {
	// The name by which this file system is identified.
	FileSystemId string ` + "`json:\"fileSystemId\" conjure-docs:\"The name by which this file system is identified.\"`" + `
	BaseUri      string ` + "`json:\"baseUri\"`" + `
}

func (o BackingFileSystem) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *BackingFileSystem) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
			"test/api/services.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"fmt"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/conjure-go/v6/conjure/{{currCaseTmpDir}}/foundry/catalog/api/datasets"
	"github.com/palantir/pkg/bearertoken"
)

// A Markdown description of the service.
type TestServiceClient interface {
	// Returns a mapping from file system id to backing file system configuration.
	GetFileSystems(ctx context.Context, authHeader bearertoken.Token) (map[string]datasets.BackingFileSystem, error)
}

type testServiceClient struct {
	client httpclient.Client
}

func NewTestServiceClient(client httpclient.Client) TestServiceClient {
	return &testServiceClient{client: client}
}

func (c *testServiceClient) GetFileSystems(ctx context.Context, authHeader bearertoken.Token) (map[string]datasets.BackingFileSystem, error) {
	var returnVal map[string]datasets.BackingFileSystem
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetFileSystems"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/catalog/fileSystems"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return returnVal, err
	}
	_ = resp
	if returnVal == nil {
		returnVal = make(map[string]datasets.BackingFileSystem, 0)
	}
	return returnVal, nil
}

// A Markdown description of the service.
type TestServiceClientWithAuth interface {
	// Returns a mapping from file system id to backing file system configuration.
	GetFileSystems(ctx context.Context) (map[string]datasets.BackingFileSystem, error)
}

func NewTestServiceClientWithAuth(client TestServiceClient, authHeader bearertoken.Token) TestServiceClientWithAuth {
	return &testServiceClientWithAuth{client: client, authHeader: authHeader}
}

type testServiceClientWithAuth struct {
	client     TestServiceClient
	authHeader bearertoken.Token
}

func (c *testServiceClientWithAuth) GetFileSystems(ctx context.Context) (map[string]datasets.BackingFileSystem, error) {
	return c.client.GetFileSystems(ctx, c.authHeader)
}

func NewTestServiceClientWithTokenProvider(client TestServiceClient, tokenProvider httpclient.TokenProvider) TestServiceClientWithAuth {
	return &testServiceClientWithTokenProvider{client: client, tokenProvider: tokenProvider}
}

type testServiceClientWithTokenProvider struct {
	client        TestServiceClient
	tokenProvider httpclient.TokenProvider
}

func (c *testServiceClientWithTokenProvider) GetFileSystems(ctx context.Context) (map[string]datasets.BackingFileSystem, error) {
	var defaultReturnVal map[string]datasets.BackingFileSystem
	token, err := c.tokenProvider(ctx)
	if err != nil {
		return defaultReturnVal, err
	}
	return c.client.GetFileSystems(ctx, bearertoken.Token(token))
}
`,
		},
	},
	{
		name: "type definition with multi-line comment",
		src: `
{
	"version" : 1,
	"errors" : [ ],
	"types" : [ {
		"type" : "object",
		"object" : {
		"typeName" : {
			"name" : "ServiceLogV1",
			"package" : "com.palantir.spec.logging"
		},
		"fields" : [ {
			"fieldName" : "type",
			"type" : {
			"type" : "primitive",
			"primitive" : "STRING"
			},
			"docs" : "Docs for the type field.\nMulti-line comment on a struct field."
		} ],
		"docs" : "Definition of the service.1 format.\nFor more information, refer to the logging specification.\n"
		}
	} ],
	"services" : [ ]
}
`,
		wantFiles: map[string]string{
			"spec/logging/structs.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package logging

import (
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

// Definition of the service.1 format.
// For more information, refer to the logging specification.
type ServiceLogV1 struct {
	// Docs for the type field.
	// Multi-line comment on a struct field.
	Type string ` + "`" + `json:"type" conjure-docs:"Docs for the type field.\nMulti-line comment on a struct field."` + "`" + `
}

func (o ServiceLogV1) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *ServiceLogV1) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
	},
	{
		name: "collection types",
		src: `
			{
				"version" : 1,
				"errors" : [ ],
				"types" : [ {
				  "type" : "object",
				  "object" : {
					"typeName" : {
					  "name" : "BackingFileSystem",
					  "package" : "com.palantir.sls.spec.logging"
					},
					"fields" : [ {
					  "fieldName" : "baseUri",
					  "type" : {
						"type" : "primitive",
						"primitive" : "STRING"
					  }
					}, {
					  "fieldName" : "configuration",
					  "type" : {
						"type" : "map",
						"map" : {
						  "keyType" : {
							"type" : "primitive",
							"primitive" : "STRING"
						  },
						  "valueType" : {
							"type" : "primitive",
							"primitive" : "STRING"
						  }
						}
					  }
					}, {
					  "fieldName" : "configurationList",
					  "type" : {
						"type" : "list",
						"list" : {
						  "itemType" : {
							"type" : "primitive",
							"primitive" : "STRING"
						  }
						}
					  }
					}, {
					  "fieldName" : "configurationSet",
					  "type" : {
						"type" : "set",
						"set" : {
						  "itemType" : {
							"type" : "primitive",
							"primitive" : "STRING"
						  }
						}
					  }
					} ]
				  }
				} ],
				"services" : [ ]
			  }
`,
		wantFiles: map[string]string{
			"sls/spec/logging/structs.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package logging

import (
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type BackingFileSystem struct {
	BaseUri           string            ` + "`" + `json:"baseUri"` + "`" + `
	Configuration     map[string]string ` + "`" + `json:"configuration"` + "`" + `
	ConfigurationList []string          ` + "`" + `json:"configurationList"` + "`" + `
	ConfigurationSet  []string          ` + "`" + `json:"configurationSet"` + "`" + `
}

func (o BackingFileSystem) MarshalJSON() ([]byte, error) {
	if o.Configuration == nil {
		o.Configuration = make(map[string]string, 0)
	}
	if o.ConfigurationList == nil {
		o.ConfigurationList = make([]string, 0)
	}
	if o.ConfigurationSet == nil {
		o.ConfigurationSet = make([]string, 0)
	}
	type BackingFileSystemAlias BackingFileSystem
	return safejson.Marshal(BackingFileSystemAlias(o))
}

func (o *BackingFileSystem) UnmarshalJSON(data []byte) error {
	type BackingFileSystemAlias BackingFileSystem
	var rawBackingFileSystem BackingFileSystemAlias
	if err := safejson.Unmarshal(data, &rawBackingFileSystem); err != nil {
		return err
	}
	if rawBackingFileSystem.Configuration == nil {
		rawBackingFileSystem.Configuration = make(map[string]string, 0)
	}
	if rawBackingFileSystem.ConfigurationList == nil {
		rawBackingFileSystem.ConfigurationList = make([]string, 0)
	}
	if rawBackingFileSystem.ConfigurationSet == nil {
		rawBackingFileSystem.ConfigurationSet = make([]string, 0)
	}
	*o = BackingFileSystem(rawBackingFileSystem)
	return nil
}

func (o BackingFileSystem) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *BackingFileSystem) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
	},
	{
		name: "type definition with comment containing backtick",
		src: `
			{
				"version" : 1,
				"errors" : [ ],
				"types" : [ {
				  "type" : "object",
				  "object" : {
					"typeName" : {
					  "name" : "ServiceLogV1",
					  "package" : "com.palantir.sls.spec.logging"
					},
					"fields" : [ {
					  "fieldName" : "type",
					  "type" : {
						"type" : "primitive",
						"primitive" : "STRING"
					  },
					  "docs" : "Docs for the ` + "`" + `type` + "`" + ` field."
					} ],
					"docs" : "Definition of the ` + "`" + `service.1` + "`" + ` format.\n"
				  }
				} ],
				"services" : [ ]
			  }
`,
		wantFiles: map[string]string{
			"sls/spec/logging/structs.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package logging

import (
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

// Definition of the ` + "`" + `service.1` + "`" + ` format.
type ServiceLogV1 struct {
	// Docs for the ` + "`" + `type` + "`" + ` field.
	Type string ` + "`" + `json:"type" conjure-docs:"Docs for the \"type\" field."` + "`" + `
}

func (o ServiceLogV1) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *ServiceLogV1) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
	},
	{
		name: "type definition with comment containing double quotes",
		src: `
{
	"version" : 1,
	"errors" : [ ],
	"types" : [ {
		"type" : "object",
		"object" : {
		"typeName" : {
			"name" : "ServiceLogV1",
			"package" : "com.palantir.sls.spec.logging"
		},
		"fields" : [ {
			"fieldName" : "type",
			"type" : {
			"type" : "primitive",
			"primitive" : "STRING"
			},
			"docs" : "Docs for the \"type\" field."
		} ],
		"docs" : "Definition of the \"service.1\" format.\n"
		}
	} ],
	"services" : [ ]
}
`,
		wantFiles: map[string]string{
			"sls/spec/logging/structs.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package logging

import (
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

// Definition of the "service.1" format.
type ServiceLogV1 struct {
	// Docs for the "type" field.
	Type string ` + "`" + `json:"type" conjure-docs:"Docs for the \"type\" field."` + "`" + `
}

func (o ServiceLogV1) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *ServiceLogV1) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
	},
	{
		name: "type definition with comment containing backslashes",
		src: `
{
	"version" : 1,
	"errors" : [ ],
	"types" : [ {
		"type" : "object",
		"object" : {
		"typeName" : {
			"name" : "ServiceLogV1",
			"package" : "com.palantir.sls.spec.logging"
		},
		"fields" : [ {
			"fieldName" : "type",
			"type" : {
			"type" : "primitive",
			"primitive" : "STRING"
			},
			"docs" : "Docs for the \\\"type\\\" \\\\ field."
		} ],
		"docs" : "Definition of the \\\"service.1\\\" \\\\ format.\n"
		}
	} ],
	"services" : [ ]
	}

`,
		wantFiles: map[string]string{
			"sls/spec/logging/structs.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package logging

import (
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

// Definition of the \"service.1\" \\ format.
type ServiceLogV1 struct {
	// Docs for the \"type\" \\ field.
	Type string ` + "`" + `json:"type" conjure-docs:"Docs for the \\\"type\\\" \\\\ field."` + "`" + `
}

func (o ServiceLogV1) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *ServiceLogV1) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
	},
	{
		name: "full-featured union definition",
		src: `
{
	"version" : 1,
	"errors" : [ ],
	"types" : [ {
		"type" : "union",
		"union" : {
		"typeName" : {
			"name" : "ExampleUnion",
			"package" : "com.palantir.test.api"
		},
		"union" : [ {
			"fieldName" : "str",
			"type" : {
			"type" : "primitive",
			"primitive" : "STRING"
			}
		}, {
			"fieldName" : "other",
			"type" : {
			"type" : "primitive",
			"primitive" : "STRING"
			},
			"docs" : "Another string"
		}, {
			"fieldName" : "myMap",
			"type" : {
			"type" : "map",
			"map" : {
				"keyType" : {
				"type" : "primitive",
				"primitive" : "STRING"
				},
				"valueType" : {
				"type" : "list",
				"list" : {
					"itemType" : {
					"type" : "primitive",
					"primitive" : "INTEGER"
					}
				}
				}
			}
			}
		}, {
			"fieldName" : "recursive",
			"type" : {
			"type" : "reference",
			"reference" : {
				"name" : "ExampleUnion",
				"package" : "com.palantir.test.api"
			}
			}
		} ]
		}
	}, {
		"type" : "union",
		"union" : {
		"typeName" : {
			"name" : "OtherUnion",
			"package" : "com.palantir.test.api"
		},
		"union" : [ {
			"fieldName" : "str",
			"type" : {
			"type" : "primitive",
			"primitive" : "STRING"
			}
		}, {
			"fieldName" : "myMap",
			"type" : {
			"type" : "map",
			"map" : {
				"keyType" : {
				"type" : "primitive",
				"primitive" : "STRING"
				},
				"valueType" : {
				"type" : "primitive",
				"primitive" : "INTEGER"
				}
			}
			}
		} ]
		}
	} ],
	"services" : [ ]
	}
`,
		wantFiles: map[string]string{
			"test/api/unions.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"fmt"

	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type ExampleUnion struct {
	typ       string
	str       *string
	other     *string
	myMap     *map[string][]int
	recursive *ExampleUnion
}

type exampleUnionDeserializer struct {
	Type      string            ` + "`" + `json:"type"` + "`" + `
	Str       *string           ` + "`" + `json:"str"` + "`" + `
	Other     *string           ` + "`" + `json:"other"` + "`" + `
	MyMap     *map[string][]int ` + "`" + `json:"myMap"` + "`" + `
	Recursive *ExampleUnion     ` + "`" + `json:"recursive"` + "`" + `
}

func (u *exampleUnionDeserializer) toStruct() ExampleUnion {
	return ExampleUnion{typ: u.Type, str: u.Str, other: u.Other, myMap: u.MyMap, recursive: u.Recursive}
}

func (u *ExampleUnion) toSerializer() (interface{}, error) {
	switch u.typ {
	default:
		return nil, fmt.Errorf("unknown type %s", u.typ)
	case "str":
		return struct {
			Type string ` + "`" + `json:"type"` + "`" + `
			Str  string ` + "`" + `json:"str"` + "`" + `
		}{Type: "str", Str: *u.str}, nil
	case "other":
		return struct {
			Type  string ` + "`" + `json:"type"` + "`" + `
			Other string ` + "`" + `json:"other"` + "`" + `
		}{Type: "other", Other: *u.other}, nil
	case "myMap":
		return struct {
			Type  string           ` + "`" + `json:"type"` + "`" + `
			MyMap map[string][]int ` + "`" + `json:"myMap"` + "`" + `
		}{Type: "myMap", MyMap: *u.myMap}, nil
	case "recursive":
		return struct {
			Type      string       ` + "`" + `json:"type"` + "`" + `
			Recursive ExampleUnion ` + "`" + `json:"recursive"` + "`" + `
		}{Type: "recursive", Recursive: *u.recursive}, nil
	}
}

func (u ExampleUnion) MarshalJSON() ([]byte, error) {
	ser, err := u.toSerializer()
	if err != nil {
		return nil, err
	}
	return safejson.Marshal(ser)
}

func (u *ExampleUnion) UnmarshalJSON(data []byte) error {
	var deser exampleUnionDeserializer
	if err := safejson.Unmarshal(data, &deser); err != nil {
		return err
	}
	*u = deser.toStruct()
	return nil
}

func (u ExampleUnion) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(u)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (u *ExampleUnion) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&u)
}

func (u *ExampleUnion) Accept(v ExampleUnionVisitor) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(u.typ)
	case "str":
		return v.VisitStr(*u.str)
	case "other":
		return v.VisitOther(*u.other)
	case "myMap":
		return v.VisitMyMap(*u.myMap)
	case "recursive":
		return v.VisitRecursive(*u.recursive)
	}
}

type ExampleUnionVisitor interface {
	VisitStr(v string) error
	VisitOther(v string) error
	VisitMyMap(v map[string][]int) error
	VisitRecursive(v ExampleUnion) error
	VisitUnknown(typeName string) error
}

func (u *ExampleUnion) AcceptWithContext(ctx context.Context, v ExampleUnionVisitorWithContext) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknownWithContext(ctx, u.typ)
	case "str":
		return v.VisitStrWithContext(ctx, *u.str)
	case "other":
		return v.VisitOtherWithContext(ctx, *u.other)
	case "myMap":
		return v.VisitMyMapWithContext(ctx, *u.myMap)
	case "recursive":
		return v.VisitRecursiveWithContext(ctx, *u.recursive)
	}
}

type ExampleUnionVisitorWithContext interface {
	VisitStrWithContext(ctx context.Context, v string) error
	VisitOtherWithContext(ctx context.Context, v string) error
	VisitMyMapWithContext(ctx context.Context, v map[string][]int) error
	VisitRecursiveWithContext(ctx context.Context, v ExampleUnion) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
}

func NewExampleUnionFromStr(v string) ExampleUnion {
	return ExampleUnion{typ: "str", str: &v}
}

func NewExampleUnionFromOther(v string) ExampleUnion {
	return ExampleUnion{typ: "other", other: &v}
}

func NewExampleUnionFromMyMap(v map[string][]int) ExampleUnion {
	return ExampleUnion{typ: "myMap", myMap: &v}
}

func NewExampleUnionFromRecursive(v ExampleUnion) ExampleUnion {
	return ExampleUnion{typ: "recursive", recursive: &v}
}

type OtherUnion struct {
	typ   string
	str   *string
	myMap *map[string]int
}

type otherUnionDeserializer struct {
	Type  string          ` + "`" + `json:"type"` + "`" + `
	Str   *string         ` + "`" + `json:"str"` + "`" + `
	MyMap *map[string]int ` + "`" + `json:"myMap"` + "`" + `
}

func (u *otherUnionDeserializer) toStruct() OtherUnion {
	return OtherUnion{typ: u.Type, str: u.Str, myMap: u.MyMap}
}

func (u *OtherUnion) toSerializer() (interface{}, error) {
	switch u.typ {
	default:
		return nil, fmt.Errorf("unknown type %s", u.typ)
	case "str":
		return struct {
			Type string ` + "`" + `json:"type"` + "`" + `
			Str  string ` + "`" + `json:"str"` + "`" + `
		}{Type: "str", Str: *u.str}, nil
	case "myMap":
		return struct {
			Type  string         ` + "`" + `json:"type"` + "`" + `
			MyMap map[string]int ` + "`" + `json:"myMap"` + "`" + `
		}{Type: "myMap", MyMap: *u.myMap}, nil
	}
}

func (u OtherUnion) MarshalJSON() ([]byte, error) {
	ser, err := u.toSerializer()
	if err != nil {
		return nil, err
	}
	return safejson.Marshal(ser)
}

func (u *OtherUnion) UnmarshalJSON(data []byte) error {
	var deser otherUnionDeserializer
	if err := safejson.Unmarshal(data, &deser); err != nil {
		return err
	}
	*u = deser.toStruct()
	return nil
}

func (u OtherUnion) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(u)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (u *OtherUnion) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&u)
}

func (u *OtherUnion) Accept(v OtherUnionVisitor) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(u.typ)
	case "str":
		return v.VisitStr(*u.str)
	case "myMap":
		return v.VisitMyMap(*u.myMap)
	}
}

type OtherUnionVisitor interface {
	VisitStr(v string) error
	VisitMyMap(v map[string]int) error
	VisitUnknown(typeName string) error
}

func (u *OtherUnion) AcceptWithContext(ctx context.Context, v OtherUnionVisitorWithContext) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknownWithContext(ctx, u.typ)
	case "str":
		return v.VisitStrWithContext(ctx, *u.str)
	case "myMap":
		return v.VisitMyMapWithContext(ctx, *u.myMap)
	}
}

type OtherUnionVisitorWithContext interface {
	VisitStrWithContext(ctx context.Context, v string) error
	VisitMyMapWithContext(ctx context.Context, v map[string]int) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
}

func NewOtherUnionFromStr(v string) OtherUnion {
	return OtherUnion{typ: "str", str: &v}
}

func NewOtherUnionFromMyMap(v map[string]int) OtherUnion {
	return OtherUnion{typ: "myMap", myMap: &v}
}
`,
		},
	},
	{
		name: "type definition with kebab cases",
		src: `
{
	"version" : 1,
	"errors" : [ ],
	"types" : [ {
		"type" : "object",
		"object" : {
		"typeName" : {
			"name" : "ServiceLogV1",
			"package" : "com.palantir.sls.spec.logging"
		},
		"fields" : [ {
			"fieldName" : "kebab-case",
			"type" : {
				"type" : "primitive",
				"primitive" : "STRING"
			}
		} ]
		}
	} ],
	"services" : [ ]
}
`,
		wantFiles: map[string]string{
			"sls/spec/logging/structs.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package logging

import (
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type ServiceLogV1 struct {
	KebabCase string ` + "`" + `json:"kebab-case"` + "`" + `
}

func (o ServiceLogV1) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *ServiceLogV1) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
	},
	{
		name: "full-featured error definition",
		src: `
{
	"version" : 1,
	"errors" : [
		{
			"code" : "NOT_FOUND",
			"namespace" : "MyNamespace",
			"errorName" : {
				"name" : "MyNotFound",
				"package" : "com.palantir.test.another.api"
			},
			"docs" : "This is documentation of MyNotFound error.",
			"safeArgs" : [
				{
					"fieldName" : "safeArgA",
					"type" : {
						"type" : "reference",
						"reference" : {
							"name" : "SimpleObject",
							"package" : "com.palantir.test.api"
						}
					},
					"docs" : "This is safeArgA doc."
				},
				{
					"fieldName" : "safeArgB",
					"type" : {
						"type" : "primitive",
						"primitive" : "INTEGER"
					}
				}
			],
			"unsafeArgs" : [
				{
					"fieldName" : "unsafeArgA",
					"type" : {
						"type" : "primitive",
						"primitive" : "STRING"
					},
					"docs" : "This is unsafeArgA doc."
				}
			] 
		}
	],
	"types" : [ {
		"type" : "object",
		"object" : {
			"typeName" : {
				"name" : "SimpleObject",
				"package" : "com.palantir.test.api"
			},
			"fields" : [ {
				"fieldName" : "someField",
				"type" : {
					"type" : "primitive",
					"primitive" : "STRING"
				}
			} ]
		}
	} ],
	"services" : [ ]
}
`,
		wantFiles: map[string]string{
			"test/another/api/errors.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	"github.com/palantir/conjure-go/v6/conjure/{{currCaseTmpDir}}/test/api"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
	"github.com/palantir/pkg/uuid"
	werror "github.com/palantir/witchcraft-go-error"
)

type myNotFound struct {
	// This is safeArgA doc.
	SafeArgA api.SimpleObject ` + "`" + `json:"safeArgA" conjure-docs:"This is safeArgA doc."` + "`" + `
	SafeArgB int              ` + "`" + `json:"safeArgB"` + "`" + `
	// This is unsafeArgA doc.
	UnsafeArgA string ` + "`" + `json:"unsafeArgA" conjure-docs:"This is unsafeArgA doc."` + "`" + `
}

func (o myNotFound) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *myNotFound) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

// NewMyNotFound returns new instance of MyNotFound error.
func NewMyNotFound(safeArgAArg api.SimpleObject, safeArgBArg int, unsafeArgAArg string) *MyNotFound {
	return &MyNotFound{errorInstanceID: uuid.NewUUID(), stack: werror.NewStackTrace(), myNotFound: myNotFound{SafeArgA: safeArgAArg, SafeArgB: safeArgBArg, UnsafeArgA: unsafeArgAArg}}
}

// WrapWithMyNotFound returns new instance of MyNotFound error wrapping an existing error.
func WrapWithMyNotFound(err error, safeArgAArg api.SimpleObject, safeArgBArg int, unsafeArgAArg string) *MyNotFound {
	return &MyNotFound{errorInstanceID: uuid.NewUUID(), stack: werror.NewStackTrace(), cause: err, myNotFound: myNotFound{SafeArgA: safeArgAArg, SafeArgB: safeArgBArg, UnsafeArgA: unsafeArgAArg}}
}

// MyNotFound is an error type.
//
// This is documentation of MyNotFound error.
type MyNotFound struct {
	errorInstanceID uuid.UUID
	myNotFound
	cause error
	stack werror.StackTrace
}

// IsMyNotFound returns true if err is an instance of MyNotFound.
func IsMyNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := errors.GetConjureError(err).(*MyNotFound)
	return ok
}

func (e *MyNotFound) Error() string {
	return fmt.Sprintf("NOT_FOUND MyNamespace:MyNotFound (%s)", e.errorInstanceID)
}

// Cause returns the underlying cause of the error, or nil if none.
// Note that cause is not serialized and sent over the wire.
func (e *MyNotFound) Cause() error {
	return e.cause
}

// StackTrace returns the StackTrace for the error, or nil if none.
// Note that stack traces are not serialized and sent over the wire.
func (e *MyNotFound) StackTrace() werror.StackTrace {
	return e.stack
}

// Message returns the message body for the error.
func (e *MyNotFound) Message() string {
	return "NOT_FOUND MyNamespace:MyNotFound"
}

// Format implements fmt.Formatter, a requirement of werror.Werror.
func (e *MyNotFound) Format(state fmt.State, verb rune) {
	werror.Format(e, e.safeParams(), state, verb)
}

// Code returns an enum describing error category.
func (e *MyNotFound) Code() errors.ErrorCode {
	return errors.NotFound
}

// Name returns an error name identifying error type.
func (e *MyNotFound) Name() string {
	return "MyNamespace:MyNotFound"
}

// InstanceID returns unique identifier of this particular error instance.
func (e *MyNotFound) InstanceID() uuid.UUID {
	return e.errorInstanceID
}

// Parameters returns a set of named parameters detailing this particular error instance.
func (e *MyNotFound) Parameters() map[string]interface{} {
	return map[string]interface{}{"safeArgA": e.SafeArgA, "safeArgB": e.SafeArgB, "unsafeArgA": e.UnsafeArgA}
}

// safeParams returns a set of named safe parameters detailing this particular error instance.
func (e *MyNotFound) safeParams() map[string]interface{} {
	return map[string]interface{}{"safeArgA": e.SafeArgA, "safeArgB": e.SafeArgB, "errorInstanceId": e.errorInstanceID}
}

// SafeParams returns a set of named safe parameters detailing this particular error instance and
// any underlying causes.
func (e *MyNotFound) SafeParams() map[string]interface{} {
	safeParams, _ := werror.ParamsFromError(e.cause)
	for k, v := range e.safeParams() {
		if _, exists := safeParams[k]; !exists {
			safeParams[k] = v
		}
	}
	return safeParams
}

// unsafeParams returns a set of named unsafe parameters detailing this particular error instance.
func (e *MyNotFound) unsafeParams() map[string]interface{} {
	return map[string]interface{}{"unsafeArgA": e.UnsafeArgA}
}

// UnsafeParams returns a set of named unsafe parameters detailing this particular error instance and
// any underlying causes.
func (e *MyNotFound) UnsafeParams() map[string]interface{} {
	_, unsafeParams := werror.ParamsFromError(e.cause)
	for k, v := range e.unsafeParams() {
		if _, exists := unsafeParams[k]; !exists {
			unsafeParams[k] = v
		}
	}
	return unsafeParams
}

func (e MyNotFound) MarshalJSON() ([]byte, error) {
	parameters, err := safejson.Marshal(e.myNotFound)
	if err != nil {
		return nil, err
	}
	return safejson.Marshal(errors.SerializableError{ErrorCode: errors.NotFound, ErrorName: "MyNamespace:MyNotFound", ErrorInstanceID: e.errorInstanceID, Parameters: json.RawMessage(parameters)})
}

func (e *MyNotFound) UnmarshalJSON(data []byte) error {
	var serializableError errors.SerializableError
	if err := safejson.Unmarshal(data, &serializableError); err != nil {
		return err
	}
	var parameters myNotFound
	if err := safejson.Unmarshal([]byte(serializableError.Parameters), &parameters); err != nil {
		return err
	}
	e.errorInstanceID = serializableError.ErrorInstanceID
	e.myNotFound = parameters
	return nil
}

func init() {
	errors.RegisterErrorType("MyNamespace:MyNotFound", reflect.TypeOf(MyNotFound{}))
}
`,
		},
	},
	{
		name: "service definition with binary request and response",
		src: `
{
  "version" : 1,
  "errors" : [ ],
  "types" : [ ],
  "services" : [ {
    "serviceName" : {
      "name" : "TestService",
      "package" : "test.api"
    },
    "endpoints" : [ {
      "endpointName" : "putStatus",
      "httpMethod" : "PUT",
      "httpPath" : "/status",
      "returns" : {
        "type" : "primitive",
		"primitive" : "BINARY"
      },
      "args" : [ {
        "argName" : "request",
        "type" : {
          "type" : "primitive",
          "primitive" : "BINARY"
        },
        "paramType" : {
          "type" : "body",
          "body" : { }
        },
        "markers" : [ ]
      } ],
      "markers" : [ ]
    } ],
    "docs" : "A Markdown description of the service.\n"
  } ]
}
`,
		wantFiles: map[string]string{
			"test/api/services.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"io"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
)

// A Markdown description of the service.
type TestServiceClient interface {
	PutStatus(ctx context.Context, requestArg func() io.ReadCloser) (io.ReadCloser, error)
}

type testServiceClient struct {
	client httpclient.Client
}

func NewTestServiceClient(client httpclient.Client) TestServiceClient {
	return &testServiceClient{client: client}
}

func (c *testServiceClient) PutStatus(ctx context.Context, requestArg func() io.ReadCloser) (io.ReadCloser, error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PutStatus"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("PUT"))
	requestParams = append(requestParams, httpclient.WithPathf("/status"))
	requestParams = append(requestParams, httpclient.WithRawRequestBodyProvider(requestArg))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
`,
		},
	},
	{
		name: "service definition with optional header parameter that is nil",
		src: `
{
  "version" : 1,
  "errors" : [ ],
  "types" : [ ],
  "services" : [ {
    "serviceName" : {
      "name" : "TestService",
      "package" : "test.api"
    },
    "endpoints" : [ {
      "endpointName" : "putStatus",
      "httpMethod" : "PUT",
      "httpPath" : "/status",
      "args" : [ {
        "argName" : "onBehalfOf",
        "type" : {
          "type" : "optional",
          "optional" : {
            "itemType" : {
              "type" : "primitive",
              "primitive" : "UUID"
            }
          }
        },
        "paramType" : {
          "type" : "header",
          "header" : {
            "paramId" : "On-Behalf-Of"
          }
        },
        "markers" : [ ]
      } ],
      "markers" : [ ]
    } ],
    "docs" : "A Markdown description of the service.\n"
  } ]
}
`,
		wantFiles: map[string]string{
			"test/api/services.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"fmt"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/pkg/uuid"
)

// A Markdown description of the service.
type TestServiceClient interface {
	PutStatus(ctx context.Context, onBehalfOfArg *uuid.UUID) error
}

type testServiceClient struct {
	client httpclient.Client
}

func NewTestServiceClient(client httpclient.Client) TestServiceClient {
	return &testServiceClient{client: client}
}

func (c *testServiceClient) PutStatus(ctx context.Context, onBehalfOfArg *uuid.UUID) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PutStatus"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("PUT"))
	requestParams = append(requestParams, httpclient.WithPathf("/status"))
	if onBehalfOfArg != nil {
		requestParams = append(requestParams, httpclient.WithHeader("On-Behalf-Of", fmt.Sprint(*onBehalfOfArg)))
	}
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return err
	}
	_ = resp
	return nil
}
`,
		},
	},

	{
		name: "Object definition with boolean map key",
		src: `
			{
				"version" : 1,
				"errors" : [ ],
				"types" : [ {
				  "type" : "object",
				  "object" : {
					"typeName" : {
					  "name" : "Result",
					  "package" : "com.palantir.test.api"
					},
					"fields" : [ {
					  "fieldName" : "counts",
					  "type" : {
						"type" : "map",
						"map" : {
						  "keyType" : {
							"type" : "primitive",
							"primitive" : "BOOLEAN"
						  },
						  "valueType" : {
							"type" : "primitive",
							"primitive" : "INTEGER"
						  }
						}
					  }
					} ]
				  }
				} ],
				"services" : [ ]
			  }
`,
		wantFiles: map[string]string{
			"test/api/structs.conjure.go": `// This file was generated by Conjure and should not be manually edited.

package api

import (
	"github.com/palantir/pkg/boolean"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type Result struct {
	Counts map[boolean.Boolean]int ` + "`" + `json:"counts"` + "`" + `
}

func (o Result) MarshalJSON() ([]byte, error) {
	if o.Counts == nil {
		o.Counts = make(map[boolean.Boolean]int, 0)
	}
	type ResultAlias Result
	return safejson.Marshal(ResultAlias(o))
}

func (o *Result) UnmarshalJSON(data []byte) error {
	type ResultAlias Result
	var rawResult ResultAlias
	if err := safejson.Unmarshal(data, &rawResult); err != nil {
		return err
	}
	if rawResult.Counts == nil {
		rawResult.Counts = make(map[boolean.Boolean]int, 0)
	}
	*o = Result(rawResult)
	return nil
}

func (o Result) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *Result) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
`,
		},
	},
}

func TestGenerate(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir(".", "conjure-go-TestGenerate")
	defer cleanup()
	require.NoError(t, err)

	for currCaseNum, currCase := range testCases {
		t.Run(currCase.name, func(t *testing.T) {
			currCaseTmpDir, err := ioutil.TempDir(tmpDir, fmt.Sprintf("case-%d-", currCaseNum))
			require.NoError(t, err, "Case %d: %s", currCaseNum, currCase.name)

			ir, err := readConjureIRFromJSON([]byte(currCase.src))
			require.NoError(t, err, "Case %d: %s", currCaseNum, currCase.name)

			err = Generate(ir, OutputConfiguration{OutputDir: currCaseTmpDir})
			require.NoError(t, err, "Case %d: %s", currCaseNum, currCase.name)

			for k, wantSrc := range currCase.wantFiles {
				t.Run(k, func(t *testing.T) {
					wantSrc = strings.Replace(wantSrc, "{{currCaseTmpDir}}", currCaseTmpDir, -1)
					filename := path.Join(currCaseTmpDir, k)
					bytes, err := ioutil.ReadFile(filename)
					require.NoError(t, err, "Case %d: %s", currCaseNum, currCase.name)
					gotSrc := string(bytes)
					assert.Equal(t, strings.Split(wantSrc, "\n"), strings.Split(gotSrc, "\n"), "Case %d: %s\nUnexpected content for file %s:\n%s", currCaseNum, currCase.name, k, gotSrc)
				})
			}
		})
	}
}

func readConjureIRFromJSON(jsonBytes []byte) (spec.ConjureDefinition, error) {
	var conjureDefinition spec.ConjureDefinition
	if err := json.Unmarshal(jsonBytes, &conjureDefinition); err != nil {
		return spec.ConjureDefinition{}, errors.Wrapf(err, "failed to unmarshal JSON for configuration")
	}
	return conjureDefinition, nil
}
