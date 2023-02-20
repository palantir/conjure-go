package graph

import (
	"sort"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
)

func buildTypeGraph(def spec.ConjureDefinition) (*Graph[spec.TypeName], error) {
	g := newGraph[spec.TypeName](0)

	getNode := func(typeName spec.TypeName) *Node[spec.TypeName] {
		if _, ok := g.NodesByID[typeName]; !ok {
			g.addNode(typeName)
		}
		return g.NodesByID[typeName]
	}
	getNodes := func(typeNames []spec.TypeName) []*Node[spec.TypeName] {
		ret := make([]*Node[spec.TypeName], 0, len(typeNames))
		for _, typeName := range typeNames {
			ret = append(ret, getNode(typeName))
		}
		return ret
	}

	for _, errorDef := range def.Errors {
		u := getNode(errorDef.ErrorName)
		for _, field := range errorDef.SafeArgs {
			deps, err := typeNamesWithinType(field.Type)
			if err != nil {
				return nil, err
			}
			g.addEdges(u, getNodes(deps)...)
		}
		for _, field := range errorDef.UnsafeArgs {
			deps, err := typeNamesWithinType(field.Type)
			if err != nil {
				return nil, err
			}
			g.addEdges(u, getNodes(deps)...)
		}
	}

	for _, typeDef := range def.Types {
		if err := typeDef.AcceptFuncs(
			func(def spec.AliasDefinition) error {
				u := getNode(def.TypeName)
				deps, err := typeNamesWithinType(def.Alias)
				if err != nil {
					return err
				}
				g.addEdges(u, getNodes(deps)...)
				return nil
			},
			func(def spec.EnumDefinition) error {
				_ = getNode(def.TypeName)
				return nil
			},
			func(def spec.ObjectDefinition) error {
				u := getNode(def.TypeName)
				for _, field := range def.Fields {
					deps, err := typeNamesWithinType(field.Type)
					if err != nil {
						return err
					}
					g.addEdges(u, getNodes(deps)...)
				}
				return nil
			},
			func(def spec.UnionDefinition) error {
				u := getNode(def.TypeName)
				for _, field := range def.Union {
					deps, err := typeNamesWithinType(field.Type)
					if err != nil {
						return err
					}
					g.addEdges(u, getNodes(deps)...)
				}
				return nil
			},
			typeDef.ErrorOnUnknown,
		); err != nil {
			return nil, err
		}
	}

	for _, serviceDef := range def.Services {
		u := getNode(serviceDef.ServiceName)
		for _, endpointDef := range serviceDef.Endpoints {
			if endpointDef.Returns != nil {
				deps, err := typeNamesWithinType(*endpointDef.Returns)
				if err != nil {
					return nil, err
				}
				g.addEdges(u, getNodes(deps)...)
			}
			for _, arg := range endpointDef.Args {
				deps, err := typeNamesWithinType(arg.Type)
				if err != nil {
					return nil, err
				}
				g.addEdges(u, getNodes(deps)...)
				for _, marker := range arg.Markers {
					deps, err := typeNamesWithinType(marker)
					if err != nil {
						return nil, err
					}
					g.addEdges(u, getNodes(deps)...)
				}
			}
			for _, marker := range endpointDef.Markers {
				deps, err := typeNamesWithinType(marker)
				if err != nil {
					return nil, err
				}
				g.addEdges(u, getNodes(deps)...)
			}
		}
	}

	// Sort graph to keep it stable
	sort.SliceStable(g.Nodes, func(i, j int) bool {
		return compareTypes(g.Nodes[i].ID, g.Nodes[j].ID)
	})
	return g, nil
}

func compareTypes(t1, t2 spec.TypeName) bool {
	if t1.Package != t2.Package {
		return t1.Package < t2.Package
	}
	if t1.Name != t2.Name {
		return t1.Name < t2.Name
	}
	return false
}

func typeNamesWithinType(typ spec.Type) ([]spec.TypeName, error) {
	var ret []spec.TypeName
	if err := typ.AcceptFuncs(
		func(primitiveType spec.PrimitiveType) error {
			return nil
		},
		func(optionalType spec.OptionalType) error {
			names, err := typeNamesWithinType(optionalType.ItemType)
			if err != nil {
				return err
			}
			ret = append(ret, names...)
			return nil
		},
		func(listType spec.ListType) error {
			names, err := typeNamesWithinType(listType.ItemType)
			if err != nil {
				return err
			}
			ret = append(ret, names...)
			return nil
		},
		func(setType spec.SetType) error {
			names, err := typeNamesWithinType(setType.ItemType)
			if err != nil {
				return err
			}
			ret = append(ret, names...)
			return nil
		},
		func(mapType spec.MapType) error {
			names, err := typeNamesWithinType(mapType.KeyType)
			if err != nil {
				return err
			}
			ret = append(ret, names...)
			names, err = typeNamesWithinType(mapType.ValueType)
			if err != nil {
				return err
			}
			ret = append(ret, names...)
			return nil
		},
		func(name spec.TypeName) error {
			ret = append(ret, name)
			return nil
		},
		func(reference spec.ExternalReference) error {
			ret = append(ret, reference.ExternalReference)
			names, err := typeNamesWithinType(reference.Fallback)
			if err != nil {
				return err
			}
			ret = append(ret, names...)
			return nil
		},
		typ.ErrorOnUnknown,
	); err != nil {
		return nil, err
	}
	return dedup(ret), nil
}

func dedup(names []spec.TypeName) []spec.TypeName {
	nameSet := make(map[spec.TypeName]struct{})
	for _, t := range names {
		nameSet[t] = struct{}{}
	}
	ret := make([]spec.TypeName, 0, len(nameSet))
	for t := range nameSet {
		ret = append(ret, t)
	}
	sort.SliceStable(ret, func(i, j int) bool {
		return compareTypes(ret[i], ret[j])
	})
	return ret
}
