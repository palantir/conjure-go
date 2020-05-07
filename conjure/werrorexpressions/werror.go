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

package werrorexpressions

import (
	"sort"

	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/expression"
)

const (
	werrorPackage         = "werror"
	wrapFunctionName      = "Wrap"
	safeParamFunctionName = "SafeParam"
)

func CreateWrapWErrorExpression(errorExpr astgen.ASTExpr, message string, safeParams map[string]string) astgen.ASTExpr {
	params := []astgen.ASTExpr{
		errorExpr,
		expression.StringVal(message),
	}
	for _, s := range turnParamsIntoStatements(safeParams) {
		params = append(params, s)
	}
	return expression.NewCallFunction(werrorPackage, wrapFunctionName, params...)
}

func turnParamsIntoStatements(safeParams map[string]string) []astgen.ASTExpr {
	var params []astgen.ASTExpr
	paramKeys := make([]string, 0, len(safeParams))
	for k := range safeParams {
		paramKeys = append(paramKeys, k)
	}
	sort.Strings(paramKeys)

	for _, k := range paramKeys {
		v := safeParams[k]
		safeParamExpression := expression.NewCallFunction(werrorPackage, safeParamFunctionName, expression.StringVal(k), expression.StringVal(v))
		params = append(params, safeParamExpression)
	}
	return params
}
