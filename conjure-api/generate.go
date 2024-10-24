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

// This directory contains the IR representation of a conjure definition. It also contains the generated conjure go files objects

package main

import (
	"github.com/palantir/conjure-go/v6/conjure"
)

//go:generate go run $GOFILE

func main() {
	ir, err := conjure.FromIRFile("conjure-api-4.35.0.conjure.json")
	if err != nil {
		panic(err)
	}
	if err := conjure.Generate(ir, conjure.OutputConfiguration{
		GenerateFuncsVisitor: true,
		OutputDir:            ".",
	}); err != nil {
		panic(err)
	}
}
