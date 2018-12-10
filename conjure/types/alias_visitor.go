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

package types

import (
	"github.com/pkg/errors"

	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
)

type AliasVisitor struct {
	AliasValue string
}

var _ spec.TypeVisitor = &AliasVisitor{}

func NewAliasVisitor() *AliasVisitor {
	return &AliasVisitor{
		AliasValue: "",
	}
}

func (a *AliasVisitor) VisitPrimitive(v spec.PrimitiveType) error {
	a.AliasValue = string(v)
	return nil
}

func (a *AliasVisitor) VisitOptional(v spec.OptionalType) error {
	return errors.New("Golang only supports primitive aliases")
}

func (a *AliasVisitor) VisitList(v spec.ListType) error {
	return errors.New("Golang only supports primitive aliases")
}

func (a *AliasVisitor) VisitSet(v spec.SetType) error {
	return errors.New("Golang only supports primitive aliases")
}

func (a *AliasVisitor) VisitMap(v spec.MapType) error {
	return errors.New("Golang only supports primitive aliases")
}

func (a *AliasVisitor) VisitReference(v spec.TypeName) error {
	return errors.New("Golang only supports primitive aliases")
}

func (a *AliasVisitor) VisitExternal(v spec.ExternalReference) error {
	return errors.New("Golang only supports primitive aliases")
}

func (a *AliasVisitor) VisitUnknown(typeName string) error {
	return errors.New("Golang only supports primitive aliases")
}
