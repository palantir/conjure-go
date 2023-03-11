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

type typeTransformFn func(typeName spec.TypeName) (spec.TypeName, error)

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
