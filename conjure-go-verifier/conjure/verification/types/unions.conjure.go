// This file was generated by Conjure and should not be manually edited.

package types

import (
	"context"
	"fmt"

	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

// A type which can either be a StringExample, a set of strings, or an integer.
type Union struct {
	typ                  string
	stringExample        *StringExample
	set                  *[]string
	thisFieldIsAnInteger *int
	alsoAnInteger        *int
	if_                  *int
	new                  *int
	interface_           *int
}

type unionDeserializer struct {
	Type                 string         `json:"type"`
	StringExample        *StringExample `json:"stringExample"`
	Set                  *[]string      `json:"set"`
	ThisFieldIsAnInteger *int           `json:"thisFieldIsAnInteger"`
	AlsoAnInteger        *int           `json:"alsoAnInteger"`
	If                   *int           `json:"if"`
	New                  *int           `json:"new"`
	Interface            *int           `json:"interface"`
}

func (u *unionDeserializer) toStruct() Union {
	return Union{typ: u.Type, stringExample: u.StringExample, set: u.Set, thisFieldIsAnInteger: u.ThisFieldIsAnInteger, alsoAnInteger: u.AlsoAnInteger, if_: u.If, new: u.New, interface_: u.Interface}
}

func (u *Union) toSerializer() (interface{}, error) {
	switch u.typ {
	default:
		return nil, fmt.Errorf("unknown type %q", u.typ)
	case "stringExample":
		if u.stringExample == nil {
			return nil, fmt.Errorf("field \"stringExample\" is required")
		}
		return struct {
			Type          string        `json:"type"`
			StringExample StringExample `json:"stringExample"`
		}{Type: "stringExample", StringExample: *u.stringExample}, nil
	case "set":
		if u.set == nil {
			return nil, fmt.Errorf("field \"set\" is required")
		}
		return struct {
			Type string   `json:"type"`
			Set  []string `json:"set"`
		}{Type: "set", Set: *u.set}, nil
	case "thisFieldIsAnInteger":
		if u.thisFieldIsAnInteger == nil {
			return nil, fmt.Errorf("field \"thisFieldIsAnInteger\" is required")
		}
		return struct {
			Type                 string `json:"type"`
			ThisFieldIsAnInteger int    `json:"thisFieldIsAnInteger"`
		}{Type: "thisFieldIsAnInteger", ThisFieldIsAnInteger: *u.thisFieldIsAnInteger}, nil
	case "alsoAnInteger":
		if u.alsoAnInteger == nil {
			return nil, fmt.Errorf("field \"alsoAnInteger\" is required")
		}
		return struct {
			Type          string `json:"type"`
			AlsoAnInteger int    `json:"alsoAnInteger"`
		}{Type: "alsoAnInteger", AlsoAnInteger: *u.alsoAnInteger}, nil
	case "if":
		if u.if_ == nil {
			return nil, fmt.Errorf("field \"if\" is required")
		}
		return struct {
			Type string `json:"type"`
			If   int    `json:"if"`
		}{Type: "if", If: *u.if_}, nil
	case "new":
		if u.new == nil {
			return nil, fmt.Errorf("field \"new\" is required")
		}
		return struct {
			Type string `json:"type"`
			New  int    `json:"new"`
		}{Type: "new", New: *u.new}, nil
	case "interface":
		if u.interface_ == nil {
			return nil, fmt.Errorf("field \"interface\" is required")
		}
		return struct {
			Type      string `json:"type"`
			Interface int    `json:"interface"`
		}{Type: "interface", Interface: *u.interface_}, nil
	}
}

func (u Union) MarshalJSON() ([]byte, error) {
	ser, err := u.toSerializer()
	if err != nil {
		return nil, err
	}
	return safejson.Marshal(ser)
}

func (u *Union) UnmarshalJSON(data []byte) error {
	var deser unionDeserializer
	if err := safejson.Unmarshal(data, &deser); err != nil {
		return err
	}
	*u = deser.toStruct()
	switch u.typ {
	case "stringExample":
		if u.stringExample == nil {
			return fmt.Errorf("field \"stringExample\" is required")
		}
	case "set":
		if u.set == nil {
			return fmt.Errorf("field \"set\" is required")
		}
	case "thisFieldIsAnInteger":
		if u.thisFieldIsAnInteger == nil {
			return fmt.Errorf("field \"thisFieldIsAnInteger\" is required")
		}
	case "alsoAnInteger":
		if u.alsoAnInteger == nil {
			return fmt.Errorf("field \"alsoAnInteger\" is required")
		}
	case "if":
		if u.if_ == nil {
			return fmt.Errorf("field \"if\" is required")
		}
	case "new":
		if u.new == nil {
			return fmt.Errorf("field \"new\" is required")
		}
	case "interface":
		if u.interface_ == nil {
			return fmt.Errorf("field \"interface\" is required")
		}
	}
	return nil
}

func (u Union) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(u)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (u *Union) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&u)
}

func (u *Union) Accept(v UnionVisitor) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(u.typ)
	case "stringExample":
		if u.stringExample == nil {
			return fmt.Errorf("field \"stringExample\" is required")
		}
		return v.VisitStringExample(*u.stringExample)
	case "set":
		if u.set == nil {
			return fmt.Errorf("field \"set\" is required")
		}
		return v.VisitSet(*u.set)
	case "thisFieldIsAnInteger":
		if u.thisFieldIsAnInteger == nil {
			return fmt.Errorf("field \"thisFieldIsAnInteger\" is required")
		}
		return v.VisitThisFieldIsAnInteger(*u.thisFieldIsAnInteger)
	case "alsoAnInteger":
		if u.alsoAnInteger == nil {
			return fmt.Errorf("field \"alsoAnInteger\" is required")
		}
		return v.VisitAlsoAnInteger(*u.alsoAnInteger)
	case "if":
		if u.if_ == nil {
			return fmt.Errorf("field \"if\" is required")
		}
		return v.VisitIf(*u.if_)
	case "new":
		if u.new == nil {
			return fmt.Errorf("field \"new\" is required")
		}
		return v.VisitNew(*u.new)
	case "interface":
		if u.interface_ == nil {
			return fmt.Errorf("field \"interface\" is required")
		}
		return v.VisitInterface(*u.interface_)
	}
}

type UnionVisitor interface {
	VisitStringExample(v StringExample) error
	VisitSet(v []string) error
	VisitThisFieldIsAnInteger(v int) error
	VisitAlsoAnInteger(v int) error
	VisitIf(v int) error
	VisitNew(v int) error
	VisitInterface(v int) error
	VisitUnknown(typeName string) error
}

func (u *Union) AcceptWithContext(ctx context.Context, v UnionVisitorWithContext) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknownWithContext(ctx, u.typ)
	case "stringExample":
		if u.stringExample == nil {
			return fmt.Errorf("field \"stringExample\" is required")
		}
		return v.VisitStringExampleWithContext(ctx, *u.stringExample)
	case "set":
		if u.set == nil {
			return fmt.Errorf("field \"set\" is required")
		}
		return v.VisitSetWithContext(ctx, *u.set)
	case "thisFieldIsAnInteger":
		if u.thisFieldIsAnInteger == nil {
			return fmt.Errorf("field \"thisFieldIsAnInteger\" is required")
		}
		return v.VisitThisFieldIsAnIntegerWithContext(ctx, *u.thisFieldIsAnInteger)
	case "alsoAnInteger":
		if u.alsoAnInteger == nil {
			return fmt.Errorf("field \"alsoAnInteger\" is required")
		}
		return v.VisitAlsoAnIntegerWithContext(ctx, *u.alsoAnInteger)
	case "if":
		if u.if_ == nil {
			return fmt.Errorf("field \"if\" is required")
		}
		return v.VisitIfWithContext(ctx, *u.if_)
	case "new":
		if u.new == nil {
			return fmt.Errorf("field \"new\" is required")
		}
		return v.VisitNewWithContext(ctx, *u.new)
	case "interface":
		if u.interface_ == nil {
			return fmt.Errorf("field \"interface\" is required")
		}
		return v.VisitInterfaceWithContext(ctx, *u.interface_)
	}
}

type UnionVisitorWithContext interface {
	VisitStringExampleWithContext(ctx context.Context, v StringExample) error
	VisitSetWithContext(ctx context.Context, v []string) error
	VisitThisFieldIsAnIntegerWithContext(ctx context.Context, v int) error
	VisitAlsoAnIntegerWithContext(ctx context.Context, v int) error
	VisitIfWithContext(ctx context.Context, v int) error
	VisitNewWithContext(ctx context.Context, v int) error
	VisitInterfaceWithContext(ctx context.Context, v int) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
}

func NewUnionFromStringExample(v StringExample) Union {
	return Union{typ: "stringExample", stringExample: &v}
}

func NewUnionFromSet(v []string) Union {
	return Union{typ: "set", set: &v}
}

func NewUnionFromThisFieldIsAnInteger(v int) Union {
	return Union{typ: "thisFieldIsAnInteger", thisFieldIsAnInteger: &v}
}

func NewUnionFromAlsoAnInteger(v int) Union {
	return Union{typ: "alsoAnInteger", alsoAnInteger: &v}
}

func NewUnionFromIf(v int) Union {
	return Union{typ: "if", if_: &v}
}

func NewUnionFromNew(v int) Union {
	return Union{typ: "new", new: &v}
}

func NewUnionFromInterface(v int) Union {
	return Union{typ: "interface", interface_: &v}
}
