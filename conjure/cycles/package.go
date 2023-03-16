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
	"sort"
	"strings"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
)

type packageSet map[string]struct{}
type packageSetStr string

func (s *packageSet) toString() packageSetStr {
	strs := make([]string, 0, len(*s))
	for str := range *s {
		strs = append(strs, str)
	}
	sort.Strings(strs)
	return packageSetStr(strings.Join(strs, ";"))
}

// RemovePackageCycles modifies the conjure definition in order to remove package cycles in the compiled Go code.
// Please check the README.md file of this package for information on how this is done.
func RemovePackageCycles(def spec.ConjureDefinition) (spec.ConjureDefinition, error) {
	// Step 1: build the type graph for all errors, objects and services of the conjure def
	typeGraph, err := buildTypeGraph(def)
	if err != nil {
		return spec.ConjureDefinition{}, err
	}

	// Step 2: calculate the strongly connected components (SCCs) of the type graph
	sccs := calculateStronglyConnectedComponents(typeGraph)

	// Step 3: merge strongly connected components as much as possible
	packageSetByComponent := make(map[componentID]packageSetStr, len(sccs.componentGraph.nodes))
	for compID, types := range sccs.components {
		packages := make(packageSet)
		for _, typ := range types {
			packages[typ.Package] = struct{}{}
		}
		packageSetByComponent[compID] = packages.toString()
	}
	_ = partition(sccs.componentGraph, packageSetByComponent)

	// TODO: remove this once bitset is no longer deadcode
	_ = newBitset(1)

	return def, nil
}
