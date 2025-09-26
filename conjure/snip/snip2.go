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

const (
	cjImport = "github.com/palantir/conjure-go/v6/cj" // TODO(bmoylan) move to CGR or pkg
)

var (
	JSONV2Encoder     = jen.Qual(ImportJSONV2Text, "Encoder").Clone
	JSONV2NewEncoder  = jen.Qual(ImportJSONV2Text, "NewEncoder").Clone
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
	JSONV2Decoder     = jen.Qual(ImportJSONV2Text, "Decoder").Clone
	JSONV2NewDecoder  = jen.Qual(ImportJSONV2Text, "NewDecoder").Clone

	JSONV2Marshal              = jen.Qual(ImportJSONV2, "Marshal").Clone
	JSONV2MarshalEncode        = jen.Qual(ImportJSONV2, "MarshalEncode").Clone
	JSONV2MarshalerTo          = jen.Qual(ImportJSONV2, "MarshalerTo").Clone
	JSONV2Unmarshal            = jen.Qual(ImportJSONV2, "Unmarshal").Clone
	JSONV2UnmarshalDecode      = jen.Qual(ImportJSONV2, "UnmarshalDecode").Clone
	JSONV2UnmarshalerFrom      = jen.Qual(ImportJSONV2, "UnmarshalerFrom").Clone
	JSONV2GetOption            = jen.Qual(ImportJSONV2, "GetOption").Clone
	JSONV2RejectUnknownMembers = jen.Qual(ImportJSONV2, "RejectUnknownMembers").Clone

	CJAny                    = jen.Qual(cjImport, "Any").Clone
	CJBearerToken            = jen.Qual(cjImport, "BearerToken").Clone
	CJBinary                 = jen.Qual(cjImport, "Binary").Clone
	CJBinaryMapKey           = jen.Qual(cjImport, "BinaryMapKey").Clone
	CJBoolean                = jen.Qual(cjImport, "Boolean").Clone
	CJBooleanMapKey          = jen.Qual(cjImport, "BooleanMapKey").Clone
	CJFloat                  = jen.Qual(cjImport, "Float").Clone
	CJFloatMapKey            = jen.Qual(cjImport, "FloatMapKey").Clone
	CJInt32                  = jen.Qual(cjImport, "Int32").Clone
	CJInt32MapKey            = jen.Qual(cjImport, "Int32MapKey").Clone
	CJSafeLong               = jen.Qual(cjImport, "SafeLong").Clone
	CJSafeLongMapKey         = jen.Qual(cjImport, "SafeLongMapKey").Clone
	CJRID                    = jen.Qual(cjImport, "RID").Clone
	CJString                 = jen.Qual(cjImport, "String").Clone
	CJUUID                   = jen.Qual(cjImport, "UUID").Clone
	CJOptionalMarshaler      = jen.Qual(cjImport, "OptionalMarshaler").Clone
	CJOptionalUnmarshaler    = jen.Qual(cjImport, "OptionalUnmarshaler").Clone
	CJListMarshaler          = jen.Qual(cjImport, "ListMarshaler").Clone
	CJListUnmarshaler        = jen.Qual(cjImport, "ListUnmarshaler").Clone
	CJSetMarshaler           = jen.Qual(cjImport, "SetMarshaler").Clone
	CJSetUnmarshaler         = jen.Qual(cjImport, "SetUnmarshaler").Clone
	CJOrderedMapMarshaler    = jen.Qual(cjImport, "OrderedMapMarshaler").Clone
	CJComparableMapMarshaler = jen.Qual(cjImport, "ComparableMapMarshaler").Clone
	CJMapUnmarshaler         = jen.Qual(cjImport, "MapUnmarshaler").Clone
	CJStringerMarshaler      = jen.Qual(cjImport, "StringerMarshaler").Clone
	CJTextMarshaler          = jen.Qual(cjImport, "TextMarshaler").Clone
	CJTextUnmarshaler        = jen.Qual(cjImport, "TextUnmarshaler").Clone
	CJStructMarshaler        = jen.Qual(cjImport, "StructMarshaler").Clone
	CJStructUnmarshaler      = jen.Qual(cjImport, "StructUnmarshaler").Clone
	CJVisitObjectFields      = jen.Qual(cjImport, "VisitObjectFields").Clone

	CJClientDecoder             = jen.Qual(cjImport, "ClientDecoder").Clone
	CJClientEncoder             = jen.Qual(cjImport, "ClientEncoder").Clone
	CJMarshal                   = jen.Qual(cjImport, "Marshal").Clone
	CJMarshalEncode             = jen.Qual(cjImport, "MarshalEncode").Clone
	CJUnmarshal                 = jen.Qual(cjImport, "Unmarshal").Clone
	CJUnmarshalDecode           = jen.Qual(cjImport, "UnmarshalDecode").Clone
	CJUnmarshalRead             = jen.Qual(cjImport, "UnmarshalRead").Clone
	CJNewSyntaxError            = jen.Qual(cjImport, "NewSyntaxError").Clone
	CJWrapSyntaxError           = jen.Qual(cjImport, "WrapSyntaxError").Clone
	CJNewKindMismatchError      = jen.Qual(cjImport, "NewKindMismatchError").Clone
	CJNewInvalidValueError      = jen.Qual(cjImport, "NewInvalidValueError").Clone
	CJNewUnmarshalFieldError    = jen.Qual(cjImport, "NewUnmarshalFieldError").Clone
	CJNewMissingFieldsError     = jen.Qual(cjImport, "NewMissingFieldsError").Clone
	CJNewUnknownFieldsError     = jen.Qual(cjImport, "NewUnknownFieldsError").Clone
	CJNewDuplicateFieldKeyError = jen.Qual(cjImport, "NewDuplicateFieldKeyError").Clone
	CJMarshalYAML               = jen.Qual(cjImport, "MarshalYAML").Clone
	CJUnmarshalYAML             = jen.Qual(cjImport, "UnmarshalYAML").Clone
)
