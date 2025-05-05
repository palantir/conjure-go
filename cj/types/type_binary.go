package types

import (
	"bytes"
	"encoding"
	"encoding/base64"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
)

// Binary provides JSON marshaling and unmarshaling for byte slice types (e.g., []byte).
// Encodes values as base64-encoded JSON strings using base64.StdEncoding.
// Implements comparison for equality and ordering based on byte content
// so values can be used as map keys.
type Binary[T ~[]byte] struct{}

func (Binary[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(base64.StdEncoding.EncodeToString(receiver)))
}

func (Binary[T]) Compare(a, b T) int {
	return bytes.Compare(a, b)
}

func (Binary[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	unquoted, err := jsontext.AppendUnquote(nil, tok.String())
	if err != nil {
		return err
	}
	decodedLen := base64.StdEncoding.DecodedLen(len(unquoted))
	if cap(*receiver) < decodedLen {
		*receiver = make([]byte, 0, decodedLen)
	} else {
		*receiver = (*receiver)[:0]
	}
	*receiver, err = base64.StdEncoding.AppendDecode(*receiver, unquoted)
	if err != nil {
		return cj.WrapSyntaxError(dec, "invalid base64", err)
	}
	return nil
}

// BinaryMarshaler provides JSON marshaling for types implementing encoding.BinaryMarshaler.
// Values are marshaled as base64-encoded JSON strings using the MarshalBinary method.
type BinaryMarshaler[T encoding.BinaryMarshaler] struct{}

func (BinaryMarshaler[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	binary, err := receiver.MarshalBinary()
	if err != nil {
		return err
	}
	return enc.WriteToken(jsontext.String(base64.StdEncoding.EncodeToString(binary)))
}

func (t BinaryMarshaler[T]) Compare(a, b T) int {
	aBinary, errA := a.MarshalBinary()
	bBinary, errB := b.MarshalBinary()
	if errA != nil || errB != nil {
		// If either fails, treat as equal (could log or handle differently)
		return 0
	}
	return bytes.Compare(aBinary, bBinary)
}

// BinaryUnmarshaler provides JSON unmarshaling for types implementing encoding.BinaryUnmarshaler.
// Expects base64-encoded JSON strings, which are decoded and passed to UnmarshalBinary.
type BinaryUnmarshaler[T encoding.BinaryUnmarshaler] struct{}

func (BinaryUnmarshaler[T]) UnmarshalJSONFrom(receiver T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	binary, err := base64.StdEncoding.DecodeString(tok.String())
	if err != nil {
		return cj.WrapSyntaxError(dec, "invalid base64", err)
	}
	return receiver.UnmarshalBinary(binary)
}
