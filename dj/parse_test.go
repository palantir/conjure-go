// Copyright (c) 2023 Palantir Technologies. All rights reserved.
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

package dj_test

import (
	"testing"

	"github.com/palantir/conjure-go/v6/dj"
	"github.com/stretchr/testify/require"
)

// this json block is poorly formed on purpose.
var basicJSON = `  {"age":100, "name":{"here":"B\\\"R"},
	"noop":{"what is a wren?":"a bird"},
	"happy":true,"immortal":false,
	"items":[1,2,3,{"tags":[1,2,3],"points":[[1,2],[3,4]]},4,5,6,7],
	"arr":["1",2,"3",{"hello":"world"},"4",5],
	"vals":[1,2,3,{"sadf":"sdf\"asdf"}],"name":{"first":"tom","last":null},
	"created":"2014-05-16T08:28:06.989Z",
	"loggy":{
		"programmers": [
    	    {
    	        "firstName": "Brett",
    	        "lastName": "McLaughlin",
    	        "email": "aaaa",
				"tag": "good"
    	    },
    	    {
    	        "firstName": "Jason",
    	        "lastName": "Hunter",
    	        "email": "bbbb",
				"tag": "bad"
    	    },
    	    {
    	        "firstName": "Elliotte",
    	        "lastName": "Harold",
    	        "email": "cccc",
				"tag": "good"
    	    },
			{
				"firstName": 1002.3,
				"age": 101
			}
    	]
	},
	"lastly":{"end...ing":"soon","yay":"final"}
}`

var basicJSONObj = map[string]any{
	"age":      float64(100),
	"name":     map[string]any{"first": "tom", "last": nil}, // overridden
	"noop":     map[string]any{"what is a wren?": "a bird"},
	"happy":    true,
	"immortal": false,
	"items": []any{
		float64(1),
		float64(2),
		float64(3),
		map[string]any{
			"tags": []any{
				float64(1),
				float64(2),
				float64(3),
			},
			"points": []any{
				[]any{float64(1), float64(2)},
				[]any{float64(3), float64(4)},
			},
		},
		float64(4),
		float64(5),
		float64(6),
		float64(7),
	},
	"arr": []any{
		"1",
		float64(2),
		"3",
		map[string]any{"hello": "world"},
		"4",
		float64(5),
	},
	"vals": []any{
		float64(1),
		float64(2),
		float64(3),
		map[string]any{"sadf": `sdf"asdf`},
	},
	"created": "2014-05-16T08:28:06.989Z",
	"loggy": map[string]any{
		"programmers": []any{
			map[string]any{
				"firstName": "Brett",
				"lastName":  "McLaughlin",
				"email":     "aaaa",
				"tag":       "good",
			},
			map[string]any{
				"firstName": "Jason",
				"lastName":  "Hunter",
				"email":     "bbbb",
				"tag":       "bad",
			},
			map[string]any{
				"firstName": "Elliotte",
				"lastName":  "Harold",
				"email":     "cccc",
				"tag":       "good",
			},
			map[string]any{
				"firstName": float64(1002.3),
				"age":       float64(101),
			},
		},
	},
	"lastly": map[string]any{
		"end...ing": "soon",
		"yay":       "final",
	},
}

var complicatedJSON = `
{
	"tagged": "OK",
	"Tagged": "KO",
	"NotTagged": true,
	"unsettable": 101,
	"Nested": {
		"Yellow": "Green",
		"yellow": "yellow"
	},
	"nestedTagged": {
		"Green": "Green",
		"Map": {
			"this": "that",
			"and": "the other thing"
		},
		"Ints": {
			"Uint": 99,
			"Uint16": 16,
			"Uint32": 32,
			"Uint64": 65
		},
		"Uints": {
			"int": -99,
			"Int": -98,
			"Int16": -16,
			"Int32": -32,
			"int64": -64,
			"Int64": -65
		},
		"Uints": {
			"Float32": 32.32,
			"Float64": 64.64
		},
		"Byte": 254,
		"Bool": true
	},
	"LeftOut": "you shouldn't be here",
	"SelfPtr": {"tagged":"OK","nestedTagged":{"Ints":{"Uint32":32}}},
	"SelfSlice": [{"tagged":"OK","nestedTagged":{"Ints":{"Uint32":32}}}],
	"SelfSlicePtr": [{"tagged":"OK","nestedTagged":{"Ints":{"Uint32":32}}}],
	"SelfPtrSlice": [{"tagged":"OK","nestedTagged":{"Ints":{"Uint32":32}}}],
	"interface": "Tile38 Rocks!",
	"Interface": "Please Download",
	"Array": [0,2,3,4,5],
	"time": "2017-05-07T13:24:43-07:00",
	"Binary": "R0lGODlhPQBEAPeo",
	"NonBinary": [9,3,100,115]
}
`

var complicatedJSONObj = map[string]any{
	"tagged":     "OK",
	"Tagged":     "KO",
	"NotTagged":  true,
	"unsettable": float64(101),
	"Nested": map[string]any{
		"Yellow": "Green",
		"yellow": "yellow",
	},
	"nestedTagged": map[string]any{
		"Green": "Green",
		"Map": map[string]any{
			"this": "that",
			"and":  "the other thing",
		},
		"Ints": map[string]any{
			"Uint":   float64(99),
			"Uint16": float64(16),
			"Uint32": float64(32),
			"Uint64": float64(65),
		},
		"Uints": map[string]any{
			"Float32": float64(32.32),
			"Float64": float64(64.64),
		},
		"Byte": float64(254),
		"Bool": true,
	},
	"Array":     []any{float64(0), float64(2), float64(3), float64(4), float64(5)},
	"Binary":    "R0lGODlhPQBEAPeo",
	"LeftOut":   "you shouldn't be here",
	"NonBinary": []any{float64(9), float64(3), float64(100), float64(115)},
	"SelfPtr": map[string]any{
		"tagged": "OK",
		"nestedTagged": map[string]any{
			"Ints": map[string]any{
				"Uint32": float64(32),
			},
		},
	},
	"SelfSlice": []any{
		map[string]any{
			"tagged": "OK",
			"nestedTagged": map[string]any{
				"Ints": map[string]any{
					"Uint32": float64(32),
				},
			},
		},
	},
	"SelfSlicePtr": []any{
		map[string]any{
			"tagged": "OK",
			"nestedTagged": map[string]any{
				"Ints": map[string]any{
					"Uint32": float64(32),
				},
			},
		},
	},
	"SelfPtrSlice": []any{
		map[string]any{
			"tagged": "OK",
			"nestedTagged": map[string]any{
				"Ints": map[string]any{
					"Uint32": float64(32),
				},
			},
		},
	},
	"interface": "Tile38 Rocks!",
	"Interface": "Please Download",
	"time":      "2017-05-07T13:24:43-07:00",
}

var exampleJSON = `{
	"widget": {
		"debug": "on",
		"window": {
			"title": "Sample Konfabulator Widget",
			"name": "main_window",
			"width": 500,
			"height": 500
		},
		"image": {
			"src": "Images/Sun.png",
			"hOffset": 250,
			"vOffset": 250,
			"alignment": "center"
		},
		"text": {
			"data": "Click Here",
			"size": 36,
			"style": "bold",
			"vOffset": 100,
			"alignment": "center",
			"onMouseUp": "sun1.opacity = (sun1.opacity / 100) * 90;"
		}
	}
}`

var exampleJSONObj = map[string]any{
	"widget": map[string]any{
		"debug": "on",
		"window": map[string]any{
			"title":  "Sample Konfabulator Widget",
			"name":   "main_window",
			"width":  float64(500),
			"height": float64(500),
		},
		"image": map[string]any{
			"src":       "Images/Sun.png",
			"hOffset":   float64(250),
			"vOffset":   float64(250),
			"alignment": "center",
		},
		"text": map[string]any{
			"data":      "Click Here",
			"size":      float64(36),
			"style":     "bold",
			"vOffset":   float64(100),
			"alignment": "center",
			"onMouseUp": "sun1.opacity = (sun1.opacity / 100) * 90;",
		},
	},
}

func TestInvalidJSON(t *testing.T) {
	for _, test := range []struct {
		JSON string
		Err  dj.SyntaxError
	}{
		{
			JSON: "",
			Err:  dj.NewSyntaxError(0, "invalid character before JSON", nil),
		},
		{
			JSON: "bad string",
			Err:  dj.NewSyntaxError(0, "invalid character beginning JSON", nil),
		},
		{
			JSON: `"open string`,
			Err:  dj.NewSyntaxError(12, "string not closed", nil),
		},
		{
			JSON: ` a""`,
			Err:  dj.NewSyntaxError(1, "invalid character beginning JSON", nil),
		},
		{
			JSON: `""a`,
			Err:  dj.NewSyntaxError(2, "invalid character after JSON", nil),
		},
		{
			JSON: "[1,2,3",
			Err:  dj.NewSyntaxError(6, "expected comma", nil),
		},
		{
			JSON: "[1 2 3]",
			Err:  dj.NewSyntaxError(3, "invalid character for comma", nil),
		},
		{
			JSON: "[1,2,3,]",
			Err:  dj.NewSyntaxError(7, "invalid character beginning JSON", nil),
		},
		{
			JSON: `{"a":1`,
			Err:  dj.NewSyntaxError(6, "expected comma", nil),
		},
		{
			JSON: `{"a":"b,}`,
			Err:  dj.NewSyntaxError(9, "string not closed", nil),
		},
		{
			JSON: `{"a":1 "b":2}`,
			Err:  dj.NewSyntaxError(7, "invalid character for comma", nil),
		},
		{
			JSON: `{"a":1,"b":2,}`,
			Err:  dj.NewSyntaxError(13, "invalid character between object entries", nil),
		},
		{
			JSON: `{"a":[1,2,3}`,
			Err:  dj.NewSyntaxError(11, "invalid character for comma", nil),
		},
		{
			JSON: "00",
			Err:  dj.NewSyntaxError(1, "invalid character after JSON", nil),
		},
		{
			JSON: "-00",
			Err:  dj.NewSyntaxError(2, "invalid character after JSON", nil),
		},
		{
			JSON: "-.",
			Err:  dj.NewSyntaxError(1, "expected digit after sign", nil),
		},
		{
			JSON: "-.123",
			Err:  dj.NewSyntaxError(1, "expected digit after sign", nil),
		},
		{
			JSON: "10EE",
			Err:  dj.NewSyntaxError(3, "expected valid digit in exp number", nil),
		},
		{
			JSON: "10E-",
			Err:  dj.NewSyntaxError(4, "expected digit following sign in exp number", nil),
		},
		{
			JSON: "10E+",
			Err:  dj.NewSyntaxError(4, "expected digit following sign in exp number", nil),
		},
		{
			JSON: " ",
			Err:  dj.NewSyntaxError(1, "invalid character before JSON", nil),
		},
		{
			JSON: "{",
			Err:  dj.NewSyntaxError(1, "object not closed", nil),
		},
		{
			JSON: "-",
			Err:  dj.NewSyntaxError(1, "sign character at end of data", nil),
		},
		{
			JSON: "-1.",
			Err:  dj.NewSyntaxError(3, "expected digit following dot", nil),
		},
		{
			JSON: "-1.0 i",
			Err:  dj.NewSyntaxError(5, "invalid character after JSON", nil),
		},
		{
			JSON: " True ",
			Err:  dj.NewSyntaxError(1, "invalid character beginning JSON", nil),
		},
		{
			JSON: " tru",
			Err:  dj.NewSyntaxError(2, "expected 'true'", nil),
		},
		{
			JSON: " False ",
			Err:  dj.NewSyntaxError(1, "invalid character beginning JSON", nil),
		},
		{
			JSON: " fals",
			Err:  dj.NewSyntaxError(2, "expected 'false'", nil),
		},
		{
			JSON: " Null ",
			Err:  dj.NewSyntaxError(1, "invalid character beginning JSON", nil),
		},
		{
			JSON: " nul",
			Err:  dj.NewSyntaxError(2, "expected 'null'", nil),
		},
		{
			JSON: " [ true,]",
			Err:  dj.NewSyntaxError(8, "invalid character beginning JSON", nil),
		},
		{
			JSON: `{ "hello": "world", }`,
			Err:  dj.NewSyntaxError(20, "invalid character between object entries", nil),
		},
		{
			JSON: `{"a":"b",}`,
			Err:  dj.NewSyntaxError(9, "invalid character between object entries", nil),
		},
		{
			JSON: `{"a":"b","a"}`,
			Err:  dj.NewSyntaxError(12, "invalid character for colon", nil),
		},
		{
			JSON: `{"a":"b","a":}`,
			Err:  dj.NewSyntaxError(13, "invalid character beginning JSON", nil),
		},
		{
			JSON: `{"a":"b",2"1":2}`,
			Err:  dj.NewSyntaxError(9, "invalid character between object entries", nil),
		},
		{
			JSON: `"`,
			Err:  dj.NewSyntaxError(1, "string not closed", nil),
		},
		{
			JSON: `"\"`,
			Err:  dj.NewSyntaxError(3, "string not closed", nil),
		},
		{
			JSON: `"a\\b\\\uFFAZa"`,
			Err:  dj.NewSyntaxError(12, "invalid unicode character", nil),
		},
		{
			JSON: `"a\\b\\\uFFA"`,
			Err:  dj.NewSyntaxError(12, "invalid unicode character", nil),
		},
		{
			JSON: "[-]",
			Err:  dj.NewSyntaxError(2, "expected digit after sign", nil),
		},
		{
			JSON: "[-.123]",
			Err:  dj.NewSyntaxError(2, "expected digit after sign", nil),
		},
	} {
		t.Run(test.JSON, func(t *testing.T) {
			require.EqualError(t, dj.Valid(test.JSON), test.Err.Error())
		})
	}
}

func TestJSON(t *testing.T) {
	for _, test := range []struct {
		Name     string
		JSON     string
		Value    any
		Err      error
		ValueErr error
	}{
		{
			JSON:  "0",
			Value: float64(0),
		},
		{
			JSON:  "0.0",
			Value: float64(0.0),
		},
		{
			JSON:  "10.0",
			Value: float64(10.0),
		},
		{
			JSON:  "10e1",
			Value: float64(10e1),
		},
		{
			JSON:  "10E123",
			Value: float64(10e123),
		},
		{
			JSON:  "10E-123",
			Value: float64(10e-123),
		},
		{
			JSON:  "10E-0123",
			Value: float64(10e-0123),
		},
		{
			JSON:  "{}",
			Value: map[string]any{},
		},
		{
			JSON:  "-1",
			Value: float64(-1),
		},
		{
			JSON:  "-1.0",
			Value: float64(-1.0),
		},
		{
			JSON:  " -1.0",
			Value: float64(-1.0),
		},
		{
			JSON:  " -1.0 ",
			Value: float64(-1.0),
		},
		{
			JSON:  "-1.0 ",
			Value: float64(-1.0),
		},
		{
			JSON:  "true",
			Value: true,
		},
		{
			JSON:  " true",
			Value: true,
		},
		{
			JSON:  " true ",
			Value: true,
		},
		{
			JSON:  "false",
			Value: false,
		},
		{
			JSON:  " false",
			Value: false,
		},
		{
			JSON:  " false ",
			Value: false,
		},
		{
			JSON:  "null",
			Value: nil,
		},
		{
			JSON:  " null",
			Value: nil,
		},
		{
			JSON:  " null ",
			Value: nil,
		},
		{
			JSON:  " []",
			Value: []any{},
		},
		{
			JSON:  " [true]",
			Value: []any{true},
		},
		{
			JSON:  " [ true, null ]",
			Value: []any{true, nil},
		},
		{
			JSON:  `{"hello":"world"}`,
			Value: map[string]any{"hello": "world"},
		},
		{
			JSON:  `{ "hello": "world" }`,
			Value: map[string]any{"hello": "world"},
		},
		{
			JSON:  `{"a":"b","a":1}`,
			Value: map[string]any{"a": float64(1)},
		},
		{
			JSON:  `{"a":"b","a": 1, "c":{"hi":"there"} }`,
			Value: map[string]any{"a": float64(1), "c": map[string]any{"hi": "there"}},
		},
		{
			JSON:  `{"a":"b","a": 1, "c":{"hi":"there", "easy":["going",{"mixed":"bag"}]} }`,
			Value: map[string]any{"a": float64(1), "c": map[string]any{"hi": "there", "easy": []any{"going", map[string]any{"mixed": "bag"}}}},
		},
		{
			JSON:  `""`,
			Value: "",
		},
		{
			JSON:  `"\n"`,
			Value: "\n",
		},
		{
			JSON:  `"\\"`,
			Value: "\\",
		},
		{
			JSON:  `"a\\b"`,
			Value: "a\\b",
		},
		{
			JSON:  `"a\\b\\\"a"`,
			Value: "a\\b\\\"a",
		},
		{
			JSON:  `"a\\b\\\uFFAAa"`,
			Value: "a\\b\\\uFFAAa",
		},
		{
			JSON:  complicatedJSON,
			Value: complicatedJSONObj,
		},
		{
			JSON:  exampleJSON,
			Value: exampleJSONObj,
		},
		{
			JSON:  basicJSON,
			Value: basicJSONObj,
		},
	} {
		t.Run(test.JSON, func(t *testing.T) {
			result, parseErr := dj.Parse(test.JSON)
			require.Equal(t, test.Err, parseErr)
			if parseErr != nil {
				return
			}
			resultValue, valueErr := result.Value()
			require.Equal(t, test.ValueErr, valueErr)
			require.Equal(t, test.Value, resultValue)
		})
	}
}

func TestResult_String(t *testing.T) {
	for _, test := range []struct {
		Name  string
		JSON  string
		Value string
		Err   string
	}{
		{
			Name:  "empty",
			JSON:  `""`,
			Value: "",
		},
		{
			Name:  "simple",
			JSON:  `"hello"`,
			Value: "hello",
		},
		{
			Name: "invalid",
			JSON: `123`,
			Err:  "type mismatch at index 0: want String got Number",
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			res, err := dj.Parse(test.JSON)
			require.NoError(t, err)
			str, err := res.String()
			if test.Err != "" {
				require.EqualError(t, err, test.Err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.Value, str)
		})
	}
}

func TestResult_Int(t *testing.T) {
	for _, test := range []struct {
		Name  string
		JSON  string
		Value int64
		Err   string
	}{
		{
			Name:  "0",
			JSON:  `0`,
			Value: 0,
		},
		{
			Name:  "simple",
			JSON:  `-123`,
			Value: -123,
		},
		{
			Name: "invalid type",
			JSON: `"123"`,
			Err:  "type mismatch at index 0: want Number got String",
		},
		{
			Name: "floating point",
			JSON: `1.23`,
			Err:  "invalid json at index 1: invalid character for int",
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			res, err := dj.Parse(test.JSON)
			require.NoError(t, err)
			str, err := res.Int()
			if test.Err != "" {
				require.EqualError(t, err, test.Err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.Value, str)
		})
	}
}
