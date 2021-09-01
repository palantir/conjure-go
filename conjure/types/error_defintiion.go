package types

import (
	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
)

type ErrorDefinition struct {
	Docs
	Name           string
	ErrorNamespace spec.ErrorNamespace
	ErrorCode      spec.ErrorCode
	SafeArgs       []*Field
	UnsafeArgs     []*Field
	conjurePkg     string
	importPath     string
}

func (t *ErrorDefinition) Code() *jen.Statement {
	return jen.Qual(t.importPath, t.Name)
}
