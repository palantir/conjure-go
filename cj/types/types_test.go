package types_test

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
	"github.com/palantir/conjure-go/v6/cj/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type typeTest interface {
	TestMarshal(t *testing.T)
	TestUnmarshal(t *testing.T)
}

type typeTestCase[T any, ENC cj.TypeEncoder[T], DEC cj.TypeDecoder[T]] struct {
	// Value is the value to encode/decode.
	Value T
	// JSON is the JSON representation of the value.
	JSON string
	// Options (optional)
	json.Options

	SkipTestMarshal      bool
	ErrMarshalJSONTo     string // if nonempty, expect MarshalJSONTo to fail
	SkipTestUnmarshal    bool
	ErrUnmarshalJSONFrom string // if nonempty, expect UnmarshalJSONFrom to fail
}

func (tc typeTestCase[T, ENC, DEC]) TestMarshal(t *testing.T) {
	t.Helper()
	if tc.SkipTestMarshal {
		t.SkipNow()
	}
	buf := bytes.NewBuffer(nil)
	enc := jsontext.NewEncoder(buf, json.JoinOptions(types.DefaultOptions, tc.Options))
	err := (*new(ENC)).MarshalJSONTo(tc.Value, enc)
	if tc.ErrMarshalJSONTo != "" {
		require.EqualErrorf(t, err, tc.ErrMarshalJSONTo, "expected error from (%T)(%v).MarshalJSON", tc.Value, tc.Value)
		return
	}
	require.NoErrorf(t, err, "unexpected error from (%T)(%v).MarshalJSON", tc.Value, tc.Value)
	got := strings.TrimSpace(buf.String())
	if assert.JSONEqf(t, tc.JSON, got, "unexpected JSON from (%T)(%v).MarshalJSON", tc.Value, tc.Value) {
		// If values are json-equivalent, assert JSON formatting/spacing
		assert.EqualValuesf(t, tc.JSON, got, "unexpected JSON formatting/spacing from (%T)(%v).MarshalJSON", tc.Value, tc.Value)
	}
}

func (tc typeTestCase[T, ENC, DEC]) TestUnmarshal(t *testing.T) {
	t.Helper()
	if tc.SkipTestUnmarshal {
		t.SkipNow()
	}
	dec := jsontext.NewDecoder(strings.NewReader(tc.JSON), json.JoinOptions(types.DefaultOptions, tc.Options))
	var got T
	err := (*new(DEC)).UnmarshalJSONFrom(&got, dec)
	if tc.ErrUnmarshalJSONFrom != "" {
		require.EqualErrorf(t, err, tc.ErrUnmarshalJSONFrom, "expected error from (%T).UnmarshalJSON(%q)", tc.Value, tc.JSON)
		return
	}
	require.NoErrorf(t, err, "unexpected error from (%T).UnmarshalJSON(%q)", tc.Value, tc.JSON)
	if isNaN(tc.Value) {
		assert.Truef(t, isNaN(got), "unexpected value from (%T).UnmarshalJSON(%q)", tc.Value, tc.JSON)
	} else {
		assert.EqualValuesf(t, tc.Value, got, "unexpected value from (%T).UnmarshalJSON(%q)", tc.Value, tc.JSON)
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func isNaN(v any) bool {
	f, ok := v.(float64)
	return ok && math.IsNaN(f)
}
