package spec

import "testing"

type v struct {
	name string
}

func (jv *v) VisitHeader(v HeaderAuthType) error {
	println(jv.name)
	return nil
}

func (jv *v) VisitCookie(v CookieAuthType) error {
	println(jv.name)
	return nil
}

func (jv *v) VisitUnknown(typeName string) error {
	panic("implement me")
}

func Test(t *testing.T) {
	auth := NewAuthTypeFromCookie(CookieAuthType{})

	v := v{name: "bob"}

	_ = auth.Accept(&v)
	_ = auth.AcceptFuncs(v.VisitHeader, v.VisitCookie, v.VisitUnknown)
	return
}

