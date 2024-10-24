// This file was generated by Conjure and should not be manually edited.

package foo1

import (
	"context"
	"fmt"

	"github.com/palantir/conjure-go/v6/cycles/testdata/pkg-cycle/conjure/com/palantir/bar"
	"github.com/palantir/conjure-go/v6/cycles/testdata/pkg-cycle/conjure/com/palantir/foo"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type Type3 struct {
	typ    string
	field1 *foo.Type2
	field2 *foo.Type4
	field3 *bar.Type1
}

type type3Deserializer struct {
	Type   string     `json:"type"`
	Field1 *foo.Type2 `json:"field1"`
	Field2 *foo.Type4 `json:"field2"`
	Field3 *bar.Type1 `json:"field3"`
}

func (u *type3Deserializer) toStruct() Type3 {
	return Type3{typ: u.Type, field1: u.Field1, field2: u.Field2, field3: u.Field3}
}

func (u *Type3) toSerializer() (interface{}, error) {
	switch u.typ {
	default:
		return nil, fmt.Errorf("unknown type %q", u.typ)
	case "field1":
		if u.field1 == nil {
			return nil, fmt.Errorf("field \"field1\" is required")
		}
		return struct {
			Type   string    `json:"type"`
			Field1 foo.Type2 `json:"field1"`
		}{Type: "field1", Field1: *u.field1}, nil
	case "field2":
		if u.field2 == nil {
			return nil, fmt.Errorf("field \"field2\" is required")
		}
		return struct {
			Type   string    `json:"type"`
			Field2 foo.Type4 `json:"field2"`
		}{Type: "field2", Field2: *u.field2}, nil
	case "field3":
		if u.field3 == nil {
			return nil, fmt.Errorf("field \"field3\" is required")
		}
		return struct {
			Type   string    `json:"type"`
			Field3 bar.Type1 `json:"field3"`
		}{Type: "field3", Field3: *u.field3}, nil
	}
}

func (u Type3) MarshalJSON() ([]byte, error) {
	ser, err := u.toSerializer()
	if err != nil {
		return nil, err
	}
	return safejson.Marshal(ser)
}

func (u *Type3) UnmarshalJSON(data []byte) error {
	var deser type3Deserializer
	if err := safejson.Unmarshal(data, &deser); err != nil {
		return err
	}
	*u = deser.toStruct()
	switch u.typ {
	case "field1":
		if u.field1 == nil {
			return fmt.Errorf("field \"field1\" is required")
		}
	case "field2":
		if u.field2 == nil {
			return fmt.Errorf("field \"field2\" is required")
		}
	case "field3":
		if u.field3 == nil {
			return fmt.Errorf("field \"field3\" is required")
		}
	}
	return nil
}

func (u Type3) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(u)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (u *Type3) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&u)
}

func (u *Type3) AcceptFuncs(field1Func func(foo.Type2) error, field2Func func(foo.Type4) error, field3Func func(bar.Type1) error, unknownFunc func(string) error) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return unknownFunc(u.typ)
	case "field1":
		if u.field1 == nil {
			return fmt.Errorf("field \"field1\" is required")
		}
		return field1Func(*u.field1)
	case "field2":
		if u.field2 == nil {
			return fmt.Errorf("field \"field2\" is required")
		}
		return field2Func(*u.field2)
	case "field3":
		if u.field3 == nil {
			return fmt.Errorf("field \"field3\" is required")
		}
		return field3Func(*u.field3)
	}
}

func (u *Type3) Field1NoopSuccess(foo.Type2) error {
	return nil
}

func (u *Type3) Field2NoopSuccess(foo.Type4) error {
	return nil
}

func (u *Type3) Field3NoopSuccess(bar.Type1) error {
	return nil
}

func (u *Type3) ErrorOnUnknown(typeName string) error {
	return fmt.Errorf("invalid value in union type. Type name: %s", typeName)
}

func (u *Type3) Accept(v Type3Visitor) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(u.typ)
	case "field1":
		if u.field1 == nil {
			return fmt.Errorf("field \"field1\" is required")
		}
		return v.VisitField1(*u.field1)
	case "field2":
		if u.field2 == nil {
			return fmt.Errorf("field \"field2\" is required")
		}
		return v.VisitField2(*u.field2)
	case "field3":
		if u.field3 == nil {
			return fmt.Errorf("field \"field3\" is required")
		}
		return v.VisitField3(*u.field3)
	}
}

type Type3Visitor interface {
	VisitField1(v foo.Type2) error
	VisitField2(v foo.Type4) error
	VisitField3(v bar.Type1) error
	VisitUnknown(typeName string) error
}

func (u *Type3) AcceptWithContext(ctx context.Context, v Type3VisitorWithContext) error {
	switch u.typ {
	default:
		if u.typ == "" {
			return fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknownWithContext(ctx, u.typ)
	case "field1":
		if u.field1 == nil {
			return fmt.Errorf("field \"field1\" is required")
		}
		return v.VisitField1WithContext(ctx, *u.field1)
	case "field2":
		if u.field2 == nil {
			return fmt.Errorf("field \"field2\" is required")
		}
		return v.VisitField2WithContext(ctx, *u.field2)
	case "field3":
		if u.field3 == nil {
			return fmt.Errorf("field \"field3\" is required")
		}
		return v.VisitField3WithContext(ctx, *u.field3)
	}
}

type Type3VisitorWithContext interface {
	VisitField1WithContext(ctx context.Context, v foo.Type2) error
	VisitField2WithContext(ctx context.Context, v foo.Type4) error
	VisitField3WithContext(ctx context.Context, v bar.Type1) error
	VisitUnknownWithContext(ctx context.Context, typeName string) error
}

func NewType3FromField1(v foo.Type2) Type3 {
	return Type3{typ: "field1", field1: &v}
}

func NewType3FromField2(v foo.Type4) Type3 {
	return Type3{typ: "field2", field2: &v}
}

func NewType3FromField3(v bar.Type1) Type3 {
	return Type3{typ: "field3", field3: &v}
}
