package types

import (
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
	werror "github.com/palantir/witchcraft-go-error"
)

func readStringToken(dec *jsontext.Decoder) (jsontext.Token, error) {
	tok, err := dec.ReadToken()
	if err != nil {
		return tok, werror.Convert(err)
	}
	if kind := tok.Kind(); kind != '"' {
		return tok, cj.NewKindMismatchError(dec, kind, "json string")
	}
	return tok, nil
}
