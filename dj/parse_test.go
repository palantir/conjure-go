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
			Err:  dj.SyntaxError{Index: 0, Msg: "invalid character before JSON"},
		},
		{
			JSON: "bad string",
			Err:  dj.SyntaxError{Index: 0, Msg: "invalid character beginning JSON"},
		},
		{
			JSON: `"open string`,
			Err:  dj.SyntaxError{Index: 12, Msg: "string not closed"},
		},
		{
			JSON: ` a""`,
			Err:  dj.SyntaxError{Index: 1, Msg: "invalid character beginning JSON"},
		},
		{
			JSON: `""a`,
			Err:  dj.SyntaxError{Index: 2, Msg: "invalid character after JSON"},
		},
		{
			JSON: "[1,2,3",
			Err:  dj.SyntaxError{Index: 6, Msg: "expected comma"},
		},
		{
			JSON: "[1 2 3]",
			Err:  dj.SyntaxError{Index: 3, Msg: "invalid character for comma"},
		},
		{
			JSON: "[1,2,3,]",
			Err:  dj.SyntaxError{Index: 7, Msg: "invalid character beginning JSON"},
		},
		{
			JSON: `{"a":1`,
			Err:  dj.SyntaxError{Index: 6, Msg: "expected comma"},
		},
		{
			JSON: `{"a":"b,}`,
			Err:  dj.SyntaxError{Index: 9, Msg: "string not closed"},
		},
		{
			JSON: `{"a":1 "b":2}`,
			Err:  dj.SyntaxError{Index: 7, Msg: "invalid character for comma"},
		},
		{
			JSON: `{"a":1,"b":2,}`,
			Err:  dj.SyntaxError{Index: 13, Msg: "invalid character between object entries"},
		},
		{
			JSON: `{"a":[1,2,3}`,
			Err:  dj.SyntaxError{Index: 11, Msg: "invalid character for comma"},
		},
		{
			JSON: "00",
			Err:  dj.SyntaxError{Index: 1, Msg: "invalid character after JSON"},
		},
		{
			JSON: "-00",
			Err:  dj.SyntaxError{Index: 2, Msg: "invalid character after JSON"},
		},
		{
			JSON: "-.",
			Err:  dj.SyntaxError{Index: 1, Msg: "expected digit after sign"},
		},
		{
			JSON: "-.123",
			Err:  dj.SyntaxError{Index: 1, Msg: "expected digit after sign"},
		},
		{
			JSON: "10EE",
			Err:  dj.SyntaxError{Index: 3, Msg: "expected valid digit in exp number"},
		},
		{
			JSON: "10E-",
			Err:  dj.SyntaxError{Index: 4, Msg: "expected digit following sign in exp number"},
		},
		{
			JSON: "10E+",
			Err:  dj.SyntaxError{Index: 4, Msg: "expected digit following sign in exp number"},
		},
		{
			JSON: " ",
			Err:  dj.SyntaxError{Index: 1, Msg: "invalid character before JSON"},
		},
		{
			JSON: "{",
			Err:  dj.SyntaxError{Index: 1, Msg: "object not closed"},
		},
		{
			JSON: "-",
			Err:  dj.SyntaxError{Index: 1, Msg: "sign character at end of data"},
		},
		{
			JSON: "-1.",
			Err:  dj.SyntaxError{Index: 3, Msg: "expected digit following dot"},
		},
		{
			JSON: "-1.0 i",
			Err:  dj.SyntaxError{Index: 5, Msg: "invalid character after JSON"},
		},
		{
			JSON: " True ",
			Err:  dj.SyntaxError{Index: 1, Msg: "invalid character beginning JSON"},
		},
		{
			JSON: " tru",
			Err:  dj.SyntaxError{Index: 2, Msg: "expected 'true'"},
		},
		{
			JSON: " False ",
			Err:  dj.SyntaxError{Index: 1, Msg: "invalid character beginning JSON"},
		},
		{
			JSON: " fals",
			Err:  dj.SyntaxError{Index: 2, Msg: "expected 'false'"},
		},
		{
			JSON: " Null ",
			Err:  dj.SyntaxError{Index: 1, Msg: "invalid character beginning JSON"},
		},
		{
			JSON: " nul",
			Err:  dj.SyntaxError{Index: 2, Msg: "expected 'null'"},
		},
		{
			JSON: " [ true,]",
			Err:  dj.SyntaxError{Index: 8, Msg: "invalid character beginning JSON"},
		},
		{
			JSON: `{ "hello": "world", }`,
			Err:  dj.SyntaxError{Index: 20, Msg: "invalid character between object entries"},
		},
		{
			JSON: `{"a":"b",}`,
			Err:  dj.SyntaxError{Index: 9, Msg: "invalid character between object entries"},
		},
		{
			JSON: `{"a":"b","a"}`,
			Err:  dj.SyntaxError{Index: 12, Msg: "invalid character for colon"},
		},
		{
			JSON: `{"a":"b","a":}`,
			Err:  dj.SyntaxError{Index: 13, Msg: "invalid character beginning JSON"},
		},
		{
			JSON: `{"a":"b",2"1":2}`,
			Err:  dj.SyntaxError{Index: 9, Msg: "invalid character between object entries"},
		},
		{
			JSON: `"`,
			Err:  dj.SyntaxError{Index: 1, Msg: "string not closed"},
		},
		{
			JSON: `"\"`,
			Err:  dj.SyntaxError{Index: 3, Msg: "string not closed"},
		},
		{
			JSON: `"a\\b\\\uFFAZa"`,
			Err:  dj.SyntaxError{Index: 12, Msg: "invalid unicode character"},
		},
		{
			JSON: `"a\\b\\\uFFA"`,
			Err:  dj.SyntaxError{Index: 12, Msg: "invalid unicode character"},
		},
		{
			JSON: "[-]",
			Err:  dj.SyntaxError{Index: 2, Msg: "expected digit after sign"},
		},
		{
			JSON: "[-.123]",
			Err:  dj.SyntaxError{Index: 2, Msg: "expected digit after sign"},
		},
	} {
		t.Run(test.JSON, func(t *testing.T) {
			require.Equal(t, test.Err, dj.Valid(test.JSON))
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
