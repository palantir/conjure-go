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

package transforms_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/palantir/conjure-go/v5/conjure/transforms"
)

func TestFieldNames(t *testing.T) {
	assert.Equal(t, "Test", transforms.ExportedFieldName("test"))
	assert.Equal(t, "KebabTest", transforms.ExportedFieldName("kebab-test"))
	assert.Equal(t, "KebabTestTest", transforms.ExportedFieldName("kebab-test-test"))
	assert.Equal(t, "KebabTest", transforms.ExportedFieldName("kebab-test-"))
	assert.Equal(t, "KebabTestTest", transforms.ExportedFieldName("kebab---test--test"))

	assert.Equal(t, "SnakeTest", transforms.ExportedFieldName("snake_test"))
	assert.Equal(t, "SnakeTestTest", transforms.ExportedFieldName("snake_test_test"))
	assert.Equal(t, "SnakeTest", transforms.ExportedFieldName("snake_test_"))
	assert.Equal(t, "SnakeTestTest", transforms.ExportedFieldName("snake___test__test"))
}
