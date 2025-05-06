package types

import (
	"sync"
	"unicode/utf8"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
)

type BearerToken[T ~string] struct{}

func (BearerToken[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(string(receiver)))
}

func (BearerToken[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	str := tok.String()
	if len(str) == 0 || str[0] == '=' {
		return cj.NewInvalidValueError(dec, "empty bearer token", nil)
	}
	chars := validBearerTokenChars()
	for i := 0; i < len(str); i++ {
		if !chars[str[i]] {
			return cj.NewInvalidValueError(dec, "invalid character for bearer token", nil)
		}
	}
	*receiver = T(tok.String())
	return nil
}

var validBearerTokenChars = sync.OnceValue(func() [utf8.RuneSelf]bool {
	var chars [utf8.RuneSelf]bool
	for i := '0'; i <= '9'; i++ {
		chars[i] = true
	}
	for i := 'A'; i <= 'Z'; i++ {
		chars[i] = true
	}
	for i := 'a'; i <= 'z'; i++ {
		chars[i] = true
	}
	chars['+'] = true
	chars['-'] = true
	chars['.'] = true
	chars['/'] = true
	chars['='] = true
	chars['_'] = true
	chars['~'] = true
	return chars
})
