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

package conjure

import (
	"encoding/json"
	"io/ioutil"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/pkg/errors"
)

func FromIRFile(file string) (spec.ConjureDefinition, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return spec.ConjureDefinition{}, errors.Wrapf(err, "failed to read IR from file %s", file)
	}
	return FromIRBytes(bytes)
}

func FromIRBytes(irJSONBytes []byte) (spec.ConjureDefinition, error) {
	var conjureDefinition spec.ConjureDefinition
	if err := json.Unmarshal(irJSONBytes, &conjureDefinition); err != nil {
		return spec.ConjureDefinition{}, errors.Wrapf(err, "failed to unmarshal JSON IR for ConjureDefinition")
	}
	return conjureDefinition, nil
}
