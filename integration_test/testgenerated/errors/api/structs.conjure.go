// This file was generated by Conjure and should not be manually edited.

package api

import (
	safejson "github.com/palantir/pkg/safejson"
)

type Basic struct {
	Data string `json:"data"`
}

func (o Basic) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o Basic) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"data\":"...)
		out = safejson.AppendQuotedString(out, o.Data)
	}
	out = append(out, '}')
	return out, nil
}
