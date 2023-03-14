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

package cycles

import (
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
)

// RemovePackageCycles modifies the conjure definition in order to remove package cycles in the compiled Go code.
// Please check the README.md file of this package for information on how this is done.
func RemovePackageCycles(def spec.ConjureDefinition) (spec.ConjureDefinition, error) {
	// Step 1: build the type graph for all errors, objects and services of the conjure def
	_, err := buildTypeGraph(def)
	if err != nil {
		return spec.ConjureDefinition{}, err
	}

	return def, nil
}
