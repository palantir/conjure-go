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

package visitors

import (
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
)

// GetParamID returns the parameter ID for the provided argument definition. If the provided definition is a header or
// query parameter and its ParamId field is non-empty, it is returned; otherwise, the argName is returned.
func GetParamID(argDef spec.ArgumentDefinition) (paramID string) {
	_ = argDef.ParamType.AcceptFuncs(
		argDef.ParamType.BodyNoopSuccess,
		func(h spec.HeaderParameterType) error {
			paramID = string(h.ParamId)
			return nil
		},
		argDef.ParamType.PathNoopSuccess,
		func(q spec.QueryParameterType) error {
			paramID = string(q.ParamId)
			return nil
		},
		func(string) error {
			return nil
		},
	)
	if paramID == "" {
		paramID = string(argDef.ArgName)
	}
	return paramID
}
