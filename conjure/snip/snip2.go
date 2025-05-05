// Copyright (c) 2025 Palantir Technologies. All rights reserved.
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

package snip

import (
	"github.com/dave/jennifer/jen"
)

const (
	ImportJSONV2     = "github.com/go-json-experiment/json"
	ImportJSONV2Text = "github.com/go-json-experiment/json/jsontext"
)

var (
	JSONV2Encoder     = jen.Qual(ImportJSONV2Text, "Encoder").Clone
	JSONV2Value       = jen.Qual(ImportJSONV2Text, "Value").Clone
	JSONV2Float       = jen.Qual(ImportJSONV2Text, "Float").Clone
	JSONV2Int         = jen.Qual(ImportJSONV2Text, "Int").Clone
	JSONV2Null        = jen.Qual(ImportJSONV2Text, "Null").Clone
	JSONV2False       = jen.Qual(ImportJSONV2Text, "False").Clone
	JSONV2True        = jen.Qual(ImportJSONV2Text, "True").Clone
	JSONV2String      = jen.Qual(ImportJSONV2Text, "String").Clone
	JSONV2BeginObject = jen.Qual(ImportJSONV2Text, "BeginObject").Clone
	JSONV2EndObject   = jen.Qual(ImportJSONV2Text, "EndObject").Clone
	JSONV2BeginArray  = jen.Qual(ImportJSONV2Text, "BeginArray").Clone
	JSONV2EndArray    = jen.Qual(ImportJSONV2Text, "EndArray").Clone

	JSONV2Decoder = jen.Qual(ImportJSONV2Text, "Decoder").Clone

	JSONV2Marshal              = jen.Qual(ImportJSONV2, "Marshal").Clone
	JSONV2MarshalEncode        = jen.Qual(ImportJSONV2, "MarshalEncode").Clone
	JSONV2MarshalerTo          = jen.Qual(ImportJSONV2, "MarshalerTo").Clone
	JSONV2Unmarshal            = jen.Qual(ImportJSONV2, "Unmarshal").Clone
	JSONV2UnmarshalDecode      = jen.Qual(ImportJSONV2, "UnmarshalDecode").Clone
	JSONV2UnmarshalerFrom      = jen.Qual(ImportJSONV2, "UnmarshalerFrom").Clone
	JSONV2GetOption            = jen.Qual(ImportJSONV2, "GetOption").Clone
	JSONV2RejectUnknownMembers = jen.Qual(ImportJSONV2, "RejectUnknownMembers").Clone
)

const (
	cjImport      = "github.com/palantir/conjure-go/v6/cj"
	cjTypesImport = "github.com/palantir/conjure-go/v6/cj/types"
)

var (
	CJTypeAny                    = jen.Qual(cjTypesImport, "Any").Clone
	CJTypeBinary                 = jen.Qual(cjTypesImport, "Binary").Clone
	CJTypeBoolean                = jen.Qual(cjTypesImport, "Boolean").Clone
	CJTypeBooleanMapKey          = jen.Qual(cjTypesImport, "BooleanMapKey").Clone
	CJTypeDateTime               = jen.Qual(cjTypesImport, "DateTime").Clone
	CJTypeFloat                  = jen.Qual(cjTypesImport, "Float").Clone
	CJTypeFloatMapKey            = jen.Qual(cjTypesImport, "FloatMapKey").Clone
	CJTypeInt                    = jen.Qual(cjTypesImport, "Int").Clone
	CJTypeIntMapKey              = jen.Qual(cjTypesImport, "IntMapKey").Clone
	CJTypeUint                   = jen.Qual(cjTypesImport, "Uint").Clone
	CJTypeUintMapKey             = jen.Qual(cjTypesImport, "UintMapKey").Clone
	CJTypeRID                    = jen.Qual(cjTypesImport, "RID").Clone
	CJTypeString                 = jen.Qual(cjTypesImport, "String").Clone
	CJTypeUUID                   = jen.Qual(cjTypesImport, "UUID").Clone
	CJTypeOptionalMarshaler      = jen.Qual(cjTypesImport, "OptionalMarshaler").Clone
	CJTypeOptionalUnmarshaler    = jen.Qual(cjTypesImport, "OptionalUnmarshaler").Clone
	CJTypeListMarshaler          = jen.Qual(cjTypesImport, "ListMarshaler").Clone
	CJTypeListUnmarshaler        = jen.Qual(cjTypesImport, "ListUnmarshaler").Clone
	CJTypeOrderedMapMarshaler    = jen.Qual(cjTypesImport, "OrderedMapMarshaler").Clone
	CJTypeComparableMapMarshaler = jen.Qual(cjTypesImport, "ComparableMapMarshaler").Clone
	CJTypeMapUnmarshaler         = jen.Qual(cjTypesImport, "MapUnmarshaler").Clone
	CJTypeTextMarshaler          = jen.Qual(cjTypesImport, "TextMarshaler").Clone
	CJTypeTextUnmarshaler        = jen.Qual(cjTypesImport, "TextUnmarshaler").Clone
	CJTypeStructMarshaler        = jen.Qual(cjTypesImport, "StructMarshaler").Clone
	CJTypeStructUnmarshaler      = jen.Qual(cjTypesImport, "StructUnmarshaler").Clone
	CJYAMLV3MarshalerFromJSON    = jen.Qual(cjImport, "YAMLV3MarshalerFromJSON").Clone
	CJYAMLV3UnmarshalerToJSON    = jen.Qual(cjImport, "YAMLV3UnmarshalerToJSON").Clone

	CJNewSyntaxError                   = jen.Qual(cjImport, "NewSyntaxError").Clone
	CJNewKindMismatchError             = jen.Qual(cjImport, "NewKindMismatchError").Clone
	CJNewInvalidValueError             = jen.Qual(cjImport, "NewInvalidValueError").Clone
	CJNewUnmarshalFieldError           = jen.Qual(cjImport, "NewUnmarshalFieldError").Clone
	CJNewUnmarshalMissingFieldsError   = jen.Qual(cjImport, "NewUnmarshalMissingFieldsError").Clone
	CJNewUnmarshalUnknownFieldsError   = jen.Qual(cjImport, "NewUnmarshalUnknownFieldsError").Clone
	CJNewUnmarshalDuplicateFieldError  = jen.Qual(cjImport, "NewUnmarshalDuplicateFieldError").Clone
	CJNewUnmarshalDuplicateMapKeyError = jen.Qual(cjImport, "NewUnmarshalDuplicateMapKeyError").Clone
)
