// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package jsonencoding

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	out, err := testMarshalJSONBuffer(nil)
	require.NoError(t, err)
	t.Log(string(out))
}

func testMarshalJSONBuffer(b []byte) ([]byte, error) {
	out := bytes.NewBuffer(b)
	_ = out.WriteByte(byte('{'))
	_, _ = out.WriteString(`"value":`)
	_, _ = out.WriteString(strconv.FormatBool(true))
	_ = out.WriteByte(byte('}'))
	return out.Bytes(), nil
}
