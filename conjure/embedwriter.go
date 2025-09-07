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

package conjure

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

// embedFileAsBlankIdentifierString adds code to the provided jen.File to
// embed the file at the specified path into a package-level blank identifier
// variable (`_`) as a string using the `//go:embed` directive. Also adds an
// anonymous import for the "embed" package.
func embedFileAsBlankIdentifierString(file *jen.File, filePath string) {
	file.Anon("embed")
	file.Comment(fmt.Sprintf("//go:embed %s", filePath))
	file.Var().Id("_").String()
}
