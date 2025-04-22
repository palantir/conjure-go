package cj

import (
	"bytes"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/pkg/safeyaml"
)

func MarshalYAMLWithGoYAMLV3(obj json.MarshalerTo) (any, error) {
	buf := bytes.NewBuffer(nil)
	enc := jsontext.NewEncoder(buf)
	if err := obj.MarshalJSONTo(enc); err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(buf.Bytes())
}
