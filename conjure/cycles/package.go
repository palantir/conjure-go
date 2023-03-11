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
	"strconv"
	"strings"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	werror "github.com/palantir/witchcraft-go-error"
)

type packageSet map[string]struct{}
type packageSetStr string

type typeTransformFn func(typeName spec.TypeName) (spec.TypeName, error)

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
	mergedComponents := partition(sccs.componentGraph, packageSetByComponent)

	// Step 4: merge types within a set of merged components and create a type transformation map
	typeTransform := make(map[spec.TypeName]spec.TypeName)
	for _, componentGroups := range mergedComponents {
		var typeGroups [][]spec.TypeName

		for _, componentGroup := range componentGroups {
			var types []spec.TypeName
			for _, compID := range componentGroup {
				types = append(types, sccs.components[compID]...)
			}
			typeGroups = append(typeGroups, dedup(types))
		}

		for idx, typeGroup := range typeGroups {
			mergeTypesIntoSamePackage(typeGroup, typeTransform, idx)
		}
	}

	// Step 5: transform types in conjure def with the transformation map
	return applyTypeTransformToDef(def, func(typeName spec.TypeName) (spec.TypeName, error) {
		newName, ok := typeTransform[typeName]
		if !ok {
			return spec.TypeName{}, werror.Error("found type not originally found in definition", werror.SafeParam("type", typeName))
		}
		return newName, nil
	})
}

func mergeTypesIntoSamePackage(types []spec.TypeName, typeTransform map[spec.TypeName]spec.TypeName, numSimilarPackageSet int) {
	sort.SliceStable(types, func(i, j int) bool {
		return compareTypes(types[i], types[j])
	})

	packageSet := make(map[string]struct{})
	for _, typ := range types {
		packageSet[typ.Package] = struct{}{}
	}
	newPackage := mergePackages(packageSet, numSimilarPackageSet)

	nameCount := make(map[string]int)
	for _, typ := range types {
		count := nameCount[typ.Name]
		nameCount[typ.Name] = count + 1

		newName := typ.Name
		if count > 0 {
			newName = newName + strconv.Itoa(count)
		}

		typeTransform[typ] = spec.TypeName{
			Package: newPackage,
			Name:    newName,
		}
	}
}

// mergePackages resolves the name of a new package that will contain all types in the packages.
func mergePackages(packageSet packageSet, numSimilarPackageSet int) (ret string) {
	// If multiple Go packages are the result of the merge of a same set of conjure packages, we keep the merged name
	// of the first one, second one is suffixed with 1, third one with 2 and so on.
	defer func() {
		if numSimilarPackageSet > 0 {
			ret = ret + strconv.Itoa(numSimilarPackageSet)
		}
	}()

	packages := make([]string, 0, len(packageSet))
	wordTable := make(map[string][]string, len(packageSet))
	for pkg := range packageSet {
		packages = append(packages, pkg)
		for _, word := range strings.Split(pkg, ".") {
			wordTable[pkg] = append(wordTable[pkg], word)
		}
	}
	sort.Strings(packages)

	// Corner case: do not merge if there is a single package
	if len(packages) == 1 {
		ret = packages[0]
		return
	}

	longestCommonPrefix := 0
	foundDiff := false
	for !foundDiff {
		firstWords := wordTable[packages[0]]
		if longestCommonPrefix+1 >= len(firstWords) {
			foundDiff = true
		}

		for i := range packages {
			words := wordTable[packages[i]]
			if longestCommonPrefix+1 >= len(words) {
				foundDiff = true
				break
			}
			if firstWords[longestCommonPrefix] != words[longestCommonPrefix] {
				foundDiff = true
				break
			}
		}
		if !foundDiff {
			longestCommonPrefix++
		}
	}

	prefix := strings.Join(wordTable[packages[0]][:longestCommonPrefix], ".")
	suffixes := make([]string, 0, len(packages))
	for _, pkg := range packages {
		suffix := strings.Join(wordTable[pkg][longestCommonPrefix:], "")
		suffixes = append(suffixes, suffix)
	}
	sort.Strings(suffixes)
	ret = prefix + "." + strings.Join(suffixes, "")
	return
}

func applyTypeTransformToDef(def spec.ConjureDefinition, typeTransform typeTransformFn) (spec.ConjureDefinition, error) {
	for i, errorDef := range def.Errors {
		for j, field := range errorDef.SafeArgs {
			newType, err := applyTypeTransformToType(field.Type, typeTransform)
			if err != nil {
				return spec.ConjureDefinition{}, err
			}
			def.Errors[i].SafeArgs[j].Type = newType
		}
		for j, field := range errorDef.UnsafeArgs {
			newType, err := applyTypeTransformToType(field.Type, typeTransform)
			if err != nil {
				return spec.ConjureDefinition{}, err
			}
			def.Errors[i].UnsafeArgs[j].Type = newType
		}
	}

	for i, typeDef := range def.Types {
		var newTypeDef spec.TypeDefinition
		if err := typeDef.AcceptFuncs(
			func(def spec.AliasDefinition) error {
				newAlias, err := applyTypeTransformToType(def.Alias, typeTransform)
				if err != nil {
					return err
				}
				def.TypeName, err = typeTransform(def.TypeName)
				if err != nil {
					return err
				}
				def.Alias = newAlias
				newTypeDef = spec.NewTypeDefinitionFromAlias(def)
				return nil
			},
			func(def spec.EnumDefinition) error {
				var err error
				def.TypeName, err = typeTransform(def.TypeName)
				if err != nil {
					return err
				}
				newTypeDef = spec.NewTypeDefinitionFromEnum(def)
				return nil
			},
			func(def spec.ObjectDefinition) error {
				var err error
				def.TypeName, err = typeTransform(def.TypeName)
				if err != nil {
					return err
				}
				for j, field := range def.Fields {
					newType, err := applyTypeTransformToType(field.Type, typeTransform)
					if err != nil {
						return err
					}
					def.Fields[j].Type = newType
				}
				newTypeDef = spec.NewTypeDefinitionFromObject(def)
				return nil
			},
			func(def spec.UnionDefinition) error {
				var err error
				def.TypeName, err = typeTransform(def.TypeName)
				if err != nil {
					return err
				}
				for j, field := range def.Union {
					newType, err := applyTypeTransformToType(field.Type, typeTransform)
					if err != nil {
						return err
					}
					def.Union[j].Type = newType
				}
				newTypeDef = spec.NewTypeDefinitionFromUnion(def)
				return nil
			},
			typeDef.ErrorOnUnknown,
		); err != nil {
			return spec.ConjureDefinition{}, err
		}
		def.Types[i] = newTypeDef
	}

	for i, serviceDef := range def.Services {
		for j, endpointDef := range serviceDef.Endpoints {
			if endpointDef.Returns != nil {
				newType, err := applyTypeTransformToType(*endpointDef.Returns, typeTransform)
				if err != nil {
					return spec.ConjureDefinition{}, err
				}
				*def.Services[i].Endpoints[j].Returns = newType
			}
			for k, arg := range endpointDef.Args {
				newType, err := applyTypeTransformToType(arg.Type, typeTransform)
				if err != nil {
					return spec.ConjureDefinition{}, err
				}
				def.Services[i].Endpoints[j].Args[k].Type = newType
				for m, marker := range arg.Markers {
					newType, err := applyTypeTransformToType(marker, typeTransform)
					if err != nil {
						return spec.ConjureDefinition{}, err
					}
					def.Services[i].Endpoints[j].Args[k].Markers[m] = newType
				}
			}
			for m, marker := range endpointDef.Markers {
				newType, err := applyTypeTransformToType(marker, typeTransform)
				if err != nil {
					return spec.ConjureDefinition{}, err
				}
				def.Services[i].Endpoints[j].Markers[m] = newType
			}
		}
	}

	return def, nil
}

func applyTypeTransformToType(typ spec.Type, typeTransform typeTransformFn) (spec.Type, error) {
	var newType spec.Type
	if err := typ.AcceptFuncs(
		func(primitiveType spec.PrimitiveType) error {
			newType = spec.NewTypeFromPrimitive(primitiveType)
			return nil
		},
		func(optionalType spec.OptionalType) error {
			newItemType, err := applyTypeTransformToType(optionalType.ItemType, typeTransform)
			if err != nil {
				return err
			}
			optionalType.ItemType = newItemType
			newType = spec.NewTypeFromOptional(optionalType)
			return nil
		},
		func(listType spec.ListType) error {
			newItemType, err := applyTypeTransformToType(listType.ItemType, typeTransform)
			if err != nil {
				return err
			}
			listType.ItemType = newItemType
			newType = spec.NewTypeFromList(listType)
			return nil
		},
		func(setType spec.SetType) error {
			newItemType, err := applyTypeTransformToType(setType.ItemType, typeTransform)
			if err != nil {
				return err
			}
			setType.ItemType = newItemType
			newType = spec.NewTypeFromSet(setType)
			return nil
		},
		func(mapType spec.MapType) error {
			newKeyType, err := applyTypeTransformToType(mapType.KeyType, typeTransform)
			if err != nil {
				return err
			}
			mapType.KeyType = newKeyType
			newValueType, err := applyTypeTransformToType(mapType.ValueType, typeTransform)
			if err != nil {
				return err
			}
			mapType.ValueType = newValueType
			newType = spec.NewTypeFromMap(mapType)
			return nil
		},
		func(name spec.TypeName) error {
			newTypeName, err := typeTransform(name)
			if err != nil {
				return err
			}
			newType = spec.NewTypeFromReference(newTypeName)
			return nil
		},
		func(reference spec.ExternalReference) error {
			var err error
			reference.ExternalReference, err = typeTransform(reference.ExternalReference)
			if err != nil {
				return err
			}
			newFallback, err := applyTypeTransformToType(reference.Fallback, typeTransform)
			if err != nil {
				return err
			}
			reference.Fallback = newFallback
			newType = spec.NewTypeFromExternal(reference)
			return nil
		},
		typ.ErrorOnUnknown,
	); err != nil {
		return spec.Type{}, err
	}
	return newType, nil
}
