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

	JSONV2Marshal       = jen.Qual(ImportJSONV2, "Marshal").Clone
	JSONV2MarshalEncode = jen.Qual(ImportJSONV2, "MarshalEncode").Clone
	JSONV2MarshalerTo   = jen.Qual(ImportJSONV2, "MarshalerTo").Clone
)
