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

package conjuretype_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/conjure/types/conjuretype"
)

func TestConjureTypes(t *testing.T) {
	type allTypes struct {
		Bearertoken conjuretype.Bearertoken        `json:"bearertoken"`
		DateTime    conjuretype.DateTime           `json:"datetime"`
		RID         conjuretype.ResourceIdentifier `json:"rid"`
		SafeLong    conjuretype.SafeLong           `json:"safelong"`
		UUID        conjuretype.UUID               `json:"uuid"`
	}

	for _, test := range []struct {
		Name   string
		Object allTypes
		JSON   string
	}{
		{
			Name: "struct with all types",
			Object: allTypes{
				Bearertoken: "so-secret",
				DateTime:    conjuretype.DateTime(time.Date(2018, 12, 27, 14, 20, 30, 0, time.UTC)),
				RID: conjuretype.ResourceIdentifier{
					Service:  "my-service",
					Instance: "my-instance",
					Type:     "my-type",
					Locator:  "my-locator",
				},
				SafeLong: 1234567890,
				UUID:     testUUID,
			},
			JSON: `{
  "bearertoken":"so-secret",
  "datetime":"2018-12-27T14:20:30Z",
  "rid":"my-service.my-instance.my-type.my-locator",
  "safelong":1234567890,
  "uuid":"00010203-0405-0607-0809-0a0b0c0d0e0f"
}`,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(test.Object)
			require.NoError(t, err)
			require.JSONEq(t, test.JSON, string(jsonBytes))
			var unmarshaled allTypes
			require.NoError(t, json.Unmarshal(jsonBytes, &unmarshaled))
			require.Equal(t, test.Object, unmarshaled)
		})
	}
}
