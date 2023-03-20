// This file was generated by Conjure and should not be manually edited.

package barfoo

import (
	"context"
	"fmt"

	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type FooType3 struct {
	typ    string
	field1 *Type2
	field3 *Type1
}

type fooType3Deserializer struct {
	Type   string `json:"type"`
	Field1 *Type2 `json:"field1"`
	Field3 *Type1 `json:"field3"`
}

func (u *fooType3Deserializer) toStruct() FooType3 {
	return FooType3{typ: u.Type, field1: u.Field1, field3: u.Field3}
}

func (u *FooType3) toSerializer() (interface{}, error) {
	switch u.typ {
	default:
		return nil, fmt.Errorf("unknown type %s", u.typ)
	case "field1":
		return struct {
			Type   string `json:"type"`
			Field1 Type2  `json:"field1"`
		}{Type: "field1", Field1: *u.field1}, nil
	case "field3":
		return struct {
			Type   string `json:"type"`
			Field3 Type1  `json:"field3"`
		}{Type: "field3", Field3: *u.field3}, nil
	}
}

func (u FooType3) MarshalJSON() ([]byte, error) {
	ser, err := u.toSerializer()
	if err != nil {
		return nil, err
	}
	return safejson.Marshal(ser)
}

func (u *FooType3) UnmarshalJSON(data []byte) error {
	var deser fooType3Deserializer
	if err := safejson.Unmarshal(data, &deser); err != nil {
		return err
	}
	*u = deser.toStruct()
	return nil
}

func (u FooType3) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(u)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (u *FooType3) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&u)
}

func (u *FooType3) AcceptFuncs(field1Func func(Type2) error, field3Func func(Type1) error, unknownFunc func(string) error) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return unknownFunc(u.typ)
	case "field1":
		return field1Func(*u.field1)
	case "field3":
		return field3Func(*u.field3)
	}
}

func (u *FooType3) Field1NoopSuccess(Type2) error {
	return nil
}

func (u *FooType3) Field3NoopSuccess(Type1) error {
	return nil
}

func (u *FooType3) ErrorOnUnknown(typeName string) error {
	return fmt.Errorf("invalid value in union type. Type name: %s", typeName)
}

func (u *FooType3) Accept(v FooType3Visitor) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(u.typ)
	case "field1":
		return v.VisitField1(*u.field1)
	case "field3":
		return v.VisitField3(*u.field3)
	}
}

type FooType3Visitor interface {
	VisitField1(v Type2) error
	VisitField3(v Type1) error
	VisitUnknown(typeName string) error
}

func (u *FooType3) AcceptWithContext(ctx context.Context, v FooType3VisitorWithContext) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknownWithContext(ctx, u.typ)
	case "field1":
		return v.VisitField1WithContext(ctx, *u.field1)
	case "field3":
		return v.VisitField3WithContext(ctx, *u.field3)
	}
}

type FooType3VisitorWithContext interface {
	VisitField1WithContext(ctx context.Context, v Type2) error
	VisitField3WithContext(ctx context.Context, v Type1) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
}

func NewFooType3FromField1(v Type2) FooType3 {
	return FooType3{typ: "field1", field1: &v}
}

func NewFooType3FromField3(v Type1) FooType3 {
	return FooType3{typ: "field3", field3: &v}
}
