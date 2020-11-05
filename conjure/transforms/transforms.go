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

// Package transforms provides common transformations from Conjure spec types and to Go specific types.
package transforms

import (
	"path"
	"strings"
	"unicode"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
)

func ExportedFieldName(input string) string {
	return Export(getFieldName(input))
}

func PrivateFieldName(input string) string {
	return SafeName(Private(getFieldName(input)))
}

func getFieldName(input string) string {
	if !strings.Contains(input, "-") && !strings.Contains(input, "_") {
		return input
	}
	var s []string
	i := 0
	for i < len(input) {
		charString := string(input[i])
		i = i + 1
		if charString != "-" && charString != "_" {
			s = append(s, charString)
			continue
		}
		if i == len(input) {
			continue
		}
		nextCharString := string(input[i])
		if nextCharString == "-" || nextCharString == "_" {
			continue
		}
		s = append(s, strings.ToUpper(nextCharString))
		i = i + 1
	}
	return strings.Join(s, "")
}

func Export(input string) string {
	return firstCharTransform(input, unicode.ToUpper)
}

func Private(input string) string {
	return firstCharTransform(input, unicode.ToLower)
}

func firstCharTransform(input string, t func(rune) rune) string {
	if len(input) == 0 {
		return input
	}
	return string([]rune{t(rune(input[0]))}) + input[1:]
}

func Documentation(documentation *spec.Documentation) string {
	var docs string
	if documentation != nil {
		docs = string(*documentation)
	}
	return docs
}

// PackagePath takes a period-delimited Conjure package path and converts it to a slash-delimited one. If the input path has
// 3 or more segments, the first two segments are omitted. For example, "com.example.folder1.folder2" -> "folder1/folder2".
func PackagePath(conjurePkgName string) string {
	parts := strings.Split(conjurePkgName, ".")
	if len(parts) > 3 {
		// if package has more than 3 parts, trim first two (typically "com.palantir")
		parts = parts[2:]
	}
	return path.Join(parts...)
}

var keywords = map[string]struct{}{
	"break":       {},
	"default":     {},
	"func":        {},
	"interface":   {},
	"select":      {},
	"case":        {},
	"defer":       {},
	"go":          {},
	"map":         {},
	"struct":      {},
	"chan":        {},
	"else":        {},
	"goto":        {},
	"package":     {},
	"switch":      {},
	"const":       {},
	"fallthrough": {},
	"if":          {},
	"range":       {},
	"type":        {},
	"continue":    {},
	"for":         {},
	"import":      {},
	"return":      {},
	"var":         {},
}

func SafeName(in string) string {
	if _, ok := keywords[in]; ok {
		return in + "_"
	}
	return in
}
