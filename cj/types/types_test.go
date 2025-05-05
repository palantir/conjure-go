package types_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
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
	if tc.SkipTestMarshal {
		t.SkipNow()
	}
	buf := bytes.NewBuffer(nil)
	enc := jsontext.NewEncoder(buf, tc.Options)
	err := (*new(ENC)).MarshalJSONTo(tc.Value, enc)
	if tc.ErrMarshalJSONTo != "" {
		require.EqualError(t, err, tc.ErrMarshalJSONTo)
		return
	}
	require.NoError(t, err)
	got := strings.TrimSpace(buf.String())
	if assert.JSONEq(t, tc.JSON, got) {
		assert.Equal(t, tc.JSON, got)
	}
}

func (tc typeTestCase[T, ENC, DEC]) TestUnmarshal(t *testing.T) {
	if tc.SkipTestUnmarshal {
		t.SkipNow()
	}
	dec := jsontext.NewDecoder(strings.NewReader(tc.JSON))
	var got T
	err := (*new(DEC)).UnmarshalJSONFrom(&got, dec)
	if tc.ErrUnmarshalJSONFrom != "" {
		require.EqualError(t, err, tc.ErrUnmarshalJSONFrom)
		return
	}
	require.NoError(t, err)
	assert.Equal(t, tc.Value, got)
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
