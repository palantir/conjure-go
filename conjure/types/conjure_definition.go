// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package types

import (
	"fmt"
	"path"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/pkg/errors"
)

type ConjureDefinition struct {
	Version    int
	Packages   map[string]ConjurePackage
	Extensions map[string]interface{}
}

type ConjurePackage struct {
	ConjurePackage string
	ImportPath     string
	OutputDir      string
	PackageName    string

	Aliases  []*AliasType
	Enums    []*EnumType
	Objects  []*ObjectType
	Unions   []*UnionType
	Errors   []*ErrorDefinition
	Services []*ServiceDefinition
}

func NewConjureDefinition(outputBaseDir string, def spec.ConjureDefinition) (*ConjureDefinition, error) {
	paths, err := newPathTranslator(outputBaseDir)
	if err != nil {
		return nil, err
	}
	names := &namedTypes{}
	packages := map[string]ConjurePackage{}
	// Add all named types to the registry. If a field/member type is a not-yet-processed Reference type,
	// names.TypeFromSpec will return a unresolvedReferencePlaceholder we will resolve later.
	for _, typeDef := range def.Types {
		if err := typeDef.AcceptFuncs(
			func(def spec.AliasDefinition) error {
				if def.Safety != nil {
					logSafetyWarning()
				}
				alias := &AliasType{
					Docs:       Docs(transforms.Documentation(def.Docs)),
					Item:       names.GetBySpec(def.Alias),
					conjurePkg: def.TypeName.Package,
					importPath: paths.conjurePkgToGoPkg(def.TypeName.Package),
					Name:       def.TypeName.Name,
				}
				names.put(def.TypeName, alias)
				pkgTypes := packages[def.TypeName.Package]
				pkgTypes.Aliases = append(pkgTypes.Aliases, alias)
				packages[def.TypeName.Package] = pkgTypes
				return nil
			},
			func(def spec.EnumDefinition) error {
				enum := &EnumType{
					Docs:       Docs(transforms.Documentation(def.Docs)),
					Values:     newFields(names, nil, def.Values),
					conjurePkg: def.TypeName.Package,
					importPath: paths.conjurePkgToGoPkg(def.TypeName.Package),
					Name:       def.TypeName.Name,
				}
				names.put(def.TypeName, enum)
				pkgTypes := packages[def.TypeName.Package]
				pkgTypes.Enums = append(pkgTypes.Enums, enum)
				packages[def.TypeName.Package] = pkgTypes
				return nil
			},
			func(def spec.ObjectDefinition) error {
				object := &ObjectType{
					Docs:       Docs(transforms.Documentation(def.Docs)),
					Fields:     newFields(names, def.Fields, nil),
					conjurePkg: def.TypeName.Package,
					importPath: paths.conjurePkgToGoPkg(def.TypeName.Package),
					Name:       def.TypeName.Name,
				}
				names.put(def.TypeName, object)
				pkgTypes := packages[def.TypeName.Package]
				pkgTypes.Objects = append(pkgTypes.Objects, object)
				packages[def.TypeName.Package] = pkgTypes
				return nil
			},
			func(def spec.UnionDefinition) error {
				union := &UnionType{
					Docs:       Docs(transforms.Documentation(def.Docs)),
					Fields:     newFields(names, def.Union, nil),
					conjurePkg: def.TypeName.Package,
					importPath: paths.conjurePkgToGoPkg(def.TypeName.Package),
					Name:       def.TypeName.Name,
				}
				names.put(def.TypeName, union)
				pkgTypes := packages[def.TypeName.Package]
				pkgTypes.Unions = append(pkgTypes.Unions, union)
				packages[def.TypeName.Package] = pkgTypes
				return nil
			},
			typeDef.ErrorOnUnknown,
		); err != nil {
			return nil, err
		}
	}

	// Now that we have all named types in the registry, go back looking
	// for unresolvedReferencePlaceholder and resolve them.
	for pkgName, pkgTypes := range names.pkgNameType {
		for typeName, typeI := range pkgTypes {
			if err := names.resolveType(typeI); err != nil {
				return nil, err
			}
			names.markComplete(pkgName, typeName)
		}
	}

	// Types are finished, move on to errors and services

	for _, def := range def.Errors {
		pkgTypes := packages[def.ErrorName.Package]
		pkgTypes.Errors = append(pkgTypes.Errors, &ErrorDefinition{
			Docs:           Docs(transforms.Documentation(def.Docs)),
			Name:           def.ErrorName.Name,
			ErrorNamespace: def.Namespace,
			ErrorCode:      def.Code,
			SafeArgs:       newFields(names, def.SafeArgs, nil),
			UnsafeArgs:     newFields(names, def.UnsafeArgs, nil),
			conjurePkg:     def.ErrorName.Package,
			importPath:     paths.conjurePkgToGoPkg(def.ErrorName.Package),
		})
		packages[def.ErrorName.Package] = pkgTypes
	}

	for _, def := range def.Services {
		var endpoints []*EndpointDefinition
		for _, endpointDef := range def.Endpoints {
			endpoint, err := newEndpointDefinition(names, endpointDef)
			if err != nil {
				return nil, err
			}
			endpoints = append(endpoints, endpoint)
		}
		pkgTypes := packages[def.ServiceName.Package]
		pkgTypes.Services = append(pkgTypes.Services, &ServiceDefinition{
			Docs:       Docs(transforms.Documentation(def.Docs)),
			Name:       def.ServiceName.Name,
			Endpoints:  endpoints,
			conjurePkg: def.ServiceName.Package,
			importPath: paths.conjurePkgToGoPkg(def.ServiceName.Package),
		})
		packages[def.ServiceName.Package] = pkgTypes
	}
	// Populate package-wide config
	for pkgName, pkg := range packages {
		pkg.ConjurePackage = pkgName
		pkg.ImportPath = paths.conjurePkgToGoPkg(pkgName)
		pkg.OutputDir = paths.conjurePkgToFilePath(pkgName)
		pkg.PackageName = sanitizePackageName(pkg.ImportPath)
		packages[pkgName] = pkg
	}
	return &ConjureDefinition{
		Version:    def.Version,
		Packages:   packages,
		Extensions: def.Extensions,
	}, nil
}

type namedTypes struct {
	pkgNameType map[string]map[string]Type
	complete    map[string]map[string]bool
}

func (t *namedTypes) put(name spec.TypeName, typ Type) {
	if t.pkgNameType == nil {
		t.pkgNameType = map[string]map[string]Type{}
	}
	if t.pkgNameType[name.Package] == nil {
		t.pkgNameType[name.Package] = map[string]Type{}
	}
	t.pkgNameType[name.Package][name.Name] = typ
}

func (t *namedTypes) GetByName(name spec.TypeName) Type {
	if pkg := t.pkgNameType[name.Package]; len(pkg) > 0 {
		return pkg[name.Name]
	}
	return nil
}

func (t *namedTypes) GetBySpec(typ spec.Type) (out Type) {
	if err := typ.AcceptFuncs(
		func(p spec.PrimitiveType) error {
			switch p.Value() {
			case spec.PrimitiveType_STRING:
				out = String{}
			case spec.PrimitiveType_DATETIME:
				out = DateTime{}
			case spec.PrimitiveType_INTEGER:
				out = Integer{}
			case spec.PrimitiveType_DOUBLE:
				out = Double{}
			case spec.PrimitiveType_SAFELONG:
				out = Safelong{}
			case spec.PrimitiveType_BINARY:
				out = Binary{}
			case spec.PrimitiveType_ANY:
				out = Any{}
			case spec.PrimitiveType_BOOLEAN:
				out = Boolean{}
			case spec.PrimitiveType_UUID:
				out = UUID{}
			case spec.PrimitiveType_RID:
				out = RID{}
			case spec.PrimitiveType_BEARERTOKEN:
				out = Bearertoken{}
			default:
				return fmt.Errorf("unknown type %s", p)
			}
			return nil
		},
		func(optional spec.OptionalType) error {
			out = &Optional{Item: t.GetBySpec(optional.ItemType)}
			return nil
		},
		func(list spec.ListType) error {
			out = &List{Item: t.GetBySpec(list.ItemType)}
			return nil
		},
		func(set spec.SetType) error {
			out = &List{Item: t.GetBySpec(set.ItemType)}
			return nil
		},
		func(map_ spec.MapType) error {
			out = &Map{Key: t.GetBySpec(map_.KeyType), Val: t.GetBySpec(map_.ValueType)}
			return nil
		},
		func(reference spec.TypeName) error {
			if ref := t.GetByName(reference); ref != nil {
				out = ref
			} else {
				out = unresolvedReferencePlaceholder{Ref: reference}
			}
			return nil
		},
		func(external spec.ExternalReference) error {
			out = &External{Spec: external.ExternalReference, Fallback: t.GetBySpec(external.Fallback)}
			return nil
		}, typ.ErrorOnUnknown); err != nil {
		return nil
	}
	return
}

// resolveType recursively searches typeI for unresolvedReferencePlaceholder in types containing member types
// and replaces them with resolved references.
func (t *namedTypes) resolveType(typeI Type) error {
	switch v := typeI.(type) {
	case *Optional:
		if unresolved, ok := v.Item.(unresolvedReferencePlaceholder); ok {
			if resolved := t.GetByName(unresolved.Ref); resolved != nil {
				v.Item = resolved
			} else {
				return errors.Errorf("Unresolved optional type reference %s %s", unresolved.Ref.Package, unresolved.Ref.Name)
			}
		} else if err := t.resolveType(v.Item); err != nil {
			return err
		}
	case *List:
		if unresolved, ok := v.Item.(unresolvedReferencePlaceholder); ok {
			if resolved := t.GetByName(unresolved.Ref); resolved != nil {
				v.Item = resolved
			} else {
				return errors.Errorf("Unresolved list item type reference %s %s", unresolved.Ref.Package, unresolved.Ref.Name)
			}
		} else if err := t.resolveType(v.Item); err != nil {
			return err
		}
	case *Map:
		if unresolved, ok := v.Key.(unresolvedReferencePlaceholder); ok {
			if resolved := t.GetByName(unresolved.Ref); resolved != nil {
				v.Key = resolved
			} else {
				return errors.Errorf("Unresolved map key type reference %s %s", unresolved.Ref.Package, unresolved.Ref.Name)
			}
		} else {
			if err := t.resolveType(v.Key); err != nil {
				return err
			}
		}
		if unresolved, ok := v.Val.(unresolvedReferencePlaceholder); ok {
			if resolved := t.GetByName(unresolved.Ref); resolved != nil {
				v.Val = resolved
			} else {
				return errors.Errorf("Unresolved map value type reference %s %s", unresolved.Ref.Package, unresolved.Ref.Name)
			}
		} else if err := t.resolveType(v.Val); err != nil {
			return err
		}
	case *AliasType:
		if !t.isComplete(v.conjurePkg, v.Name) {
			if unresolved, ok := v.Item.(unresolvedReferencePlaceholder); ok {
				if resolved := t.GetByName(unresolved.Ref); resolved != nil {
					v.Item = resolved
				} else {
					return errors.Errorf("Unresolved alias type reference %s %s", unresolved.Ref.Package, unresolved.Ref.Name)
				}
			} else if err := t.resolveType(v.Item); err != nil {
				return err
			}
			t.markComplete(v.conjurePkg, v.Name)
		}
	case *ObjectType:
		if !t.isComplete(v.conjurePkg, v.Name) {
			for i := range v.Fields {
				field := v.Fields[i]
				if unresolved, ok := field.Type.(unresolvedReferencePlaceholder); ok {
					if resolved := t.GetByName(unresolved.Ref); resolved != nil {
						field.Type = resolved
					} else {
						return errors.Errorf("Unresolved object field type reference %s %s", unresolved.Ref.Package, unresolved.Ref.Name)
					}
				} else if err := t.resolveType(field.Type); err != nil {
					return err
				}
			}
			t.markComplete(v.conjurePkg, v.Name)
		}
	case *UnionType:
		if !t.isComplete(v.conjurePkg, v.Name) {
			for i := range v.Fields {
				field := v.Fields[i]
				if unresolved, ok := field.Type.(unresolvedReferencePlaceholder); ok {
					if resolved := t.GetByName(unresolved.Ref); resolved != nil {
						field.Type = resolved
					} else {
						return errors.Errorf("Unresolved union field type reference %s %s", unresolved.Ref.Package, unresolved.Ref.Name)
					}
				} else if err := t.resolveType(field.Type); err != nil {
					return err
				}
			}
			t.markComplete(v.conjurePkg, v.Name)
		}
	}
	return nil
}

func (t *namedTypes) isComplete(pkg, name string) bool {
	if t.complete != nil {
		if pkgMap, ok := t.complete[pkg]; ok {
			return pkgMap[name]
		}
	}
	return false
}

func (t *namedTypes) markComplete(pkg, name string) {
	if t.complete == nil {
		t.complete = map[string]map[string]bool{}
	}
	if _, ok := t.complete[pkg]; !ok {
		t.complete[pkg] = map[string]bool{}
	}
	t.complete[pkg][name] = true
}

func newFields(names *namedTypes, structDefs []spec.FieldDefinition, enumDefs []spec.EnumValueDefinition) []*Field {
	var fields []*Field
	for _, value := range structDefs {
		if value.Safety != nil {
			logSafetyWarning()
		}
		fields = append(fields, &Field{
			Docs: Docs(transforms.Documentation(value.Docs)),
			Name: string(value.FieldName),
			Type: names.GetBySpec(value.Type),
		})
	}
	for _, value := range enumDefs {
		fields = append(fields, &Field{
			Docs: Docs(transforms.Documentation(value.Docs)),
			Name: value.Value,
			Type: String{},
		})
	}
	return fields
}

func newEndpointDefinition(names *namedTypes, def spec.EndpointDefinition) (*EndpointDefinition, error) {
	endpoint := &EndpointDefinition{
		Docs:         Docs(transforms.Documentation(def.Docs)),
		Deprecated:   Docs(transforms.Documentation(def.Deprecated)),
		EndpointName: string(def.EndpointName),
		HTTPMethod:   def.HttpMethod,
		HTTPPath:     string(def.HttpPath),
		Markers:      newMarkers(names, def.Markers),
		Tags:         def.Tags,
	}
	if def.Auth != nil {
		if err := def.Auth.AcceptFuncs(
			func(spec.HeaderAuthType) error {
				endpoint.HeaderAuth = true
				return nil
			},
			func(cookie spec.CookieAuthType) error {
				endpoint.CookieAuth = &cookie.CookieName
				return nil
			}, def.Auth.ErrorOnUnknown); err != nil {
			return nil, err
		}
	}
	for _, argDef := range def.Args {
		arg := &EndpointArgumentDefinition{
			Docs:    Docs(transforms.Documentation(argDef.Docs)),
			Name:    string(argDef.ArgName),
			ParamID: string(argDef.ArgName),
			Type:    names.GetBySpec(argDef.Type),
			Markers: newMarkers(names, argDef.Markers),
			Safety:  argDef.Safety,
			Tags:    argDef.Tags,
		}
		if err := argDef.ParamType.AcceptFuncs(
			func(spec.BodyParameterType) error {
				if endpoint.BodyParam() != nil {
					return fmt.Errorf("only one body parameter allowed")
				}
				arg.ParamType = BodyParam
				endpoint.Params = append(endpoint.Params, arg)
				return nil
			},
			func(p spec.HeaderParameterType) error {
				arg.ParamType = HeaderParam
				if p.ParamId != "" {
					arg.ParamID = string(p.ParamId)
				}
				endpoint.Params = append(endpoint.Params, arg)
				return nil
			},
			func(spec.PathParameterType) error {
				arg.ParamType = PathParam
				endpoint.Params = append(endpoint.Params, arg)
				return nil
			},
			func(p spec.QueryParameterType) error {
				arg.ParamType = QueryParam
				if p.ParamId != "" {
					arg.ParamID = string(p.ParamId)
				}
				endpoint.Params = append(endpoint.Params, arg)
				return nil
			},
			argDef.ParamType.ErrorOnUnknown); err != nil {
			return nil, err
		}
	}
	if def.Returns != nil {
		returns := names.GetBySpec(*def.Returns)
		endpoint.Returns = &returns
	}
	return endpoint, nil
}

func newMarkers(names *namedTypes, markers []spec.Type) (out []Type) {
	for i := range markers {
		out = append(out, names.GetBySpec(markers[i]))
	}
	return out
}

// unresolvedReferencePlaceholder takes the place of a Reference type when iterating through a definition.
// These will be replaced before a type is returned out of this package.
type unresolvedReferencePlaceholder struct {
	Ref spec.TypeName
	base
}

func (unresolvedReferencePlaceholder) Code() *jen.Statement {
	panic("unresolvedReferencePlaceholder does not implement methods")
}

func (unresolvedReferencePlaceholder) String() string {
	panic("unresolvedReferencePlaceholder does not implement methods")
}

// sanitizePackageName is based on `guessAlias` from jen/file.go: https://github.com/dave/jennifer/blob/45cc0b7eb71a469771aa486323bc53189030b60e/jen/file.go#L224
// It uses path.Base then removes non-alphanumerics and leading digits.
func sanitizePackageName(importPath string) string {
	alias := path.Base(importPath)

	// alias should be lower case
	alias = strings.ToLower(alias)

	// alias should now only contain alphanumerics
	importsRegex := regexp.MustCompile(`[^a-z0-9]`)
	alias = importsRegex.ReplaceAllString(alias, "")

	// can't have a first digit, per Go identifier rules, so just skip them
	alias = strings.TrimLeftFunc(alias, unicode.IsDigit)

	// If path part was all digits, we may be left with an empty string. In this case use "pkg" as the alias.
	if alias == "" {
		alias = "pkg"
	}

	return alias
}

var logSafetyWarningOnce sync.Once

func logSafetyWarning() {
	logSafetyWarningOnce.Do(func() {
		fmt.Println("Warning: Object definition(s) use 'safety' fields unimplemented by conjure-go.")
	})
}
