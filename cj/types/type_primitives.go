package types

import (
	"bytes"
	"encoding"
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/uuid"
)

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

type Boolean[T ~bool] struct{}

func (Boolean[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	if receiver {
		return enc.WriteToken(jsontext.True)
	}
	return enc.WriteToken(jsontext.False)
}

func (Boolean[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, 't', 'f')
	if err != nil {
		return err
	}
	switch tok.Kind() {
	case 't':
		*receiver = true
	case 'f':
		*receiver = false
	default:
		return fmt.Errorf("unexpected token kind: %v", tok.Kind())
	}
	return nil
}

type BooleanMapKey[T ~bool] struct{}

func (BooleanMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	if receiver {
		return enc.WriteValue(jsontext.Value(`"true"`))
	}
	return enc.WriteValue(jsontext.Value(`"false"`))
}

func (BooleanMapKey[T]) Compare(a, b T) int {
	if a == b {
		return 0
	}
	if a {
		return 1
	}
	return -1
}

func (BooleanMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	b, err := strconv.ParseBool(tok.String())
	if err != nil {
		return cj.WrapSyntaxError(dec, "invalid boolean", err)
	}
	*receiver = T(b)
	return nil
}

type DateTime[T time.Time | datetime.DateTime] struct{}

func (DateTime[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(time.Time(receiver).Format(time.RFC3339Nano)))
}

func (DateTime[T]) Compare(a, b T) int {
	aTime, bTime := time.Time(a), time.Time(b)
	if aTime.After(bTime) {
		return 1
	}
	if aTime.Before(bTime) {
		return -1
	}
	return 0
}

func (DateTime[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	parse, err := time.Parse(tok.String(), time.RFC3339Nano)
	if err != nil {
		return cj.WrapSyntaxError(dec, "invalid datetime", err)
	}
	*receiver = T(parse)
	return nil
}

type Float[T ~float64] struct{}

func (Float[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.Float(float64(receiver)))
}

func (Float[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '0')
	if err != nil {
		return err
	}
	*receiver = T(tok.Float())
	return nil
}

type FloatMapKey[T ~float64] struct{}

func (FloatMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	n := float64(receiver)
	switch {
	case math.Float64bits(n) == 0:
		out = append(out, "\"0\""...)
	case math.IsNaN(n):
		out = append(out, "\"NaN\""...)
	case math.IsInf(n, +1):
		out = append(out, "\"Infinity\""...)
	case math.IsInf(n, -1):
		out = append(out, "\"-Infinity\""...)
	default:
		out = append(out, '"')
		out = strconv.AppendFloat(out, n, 'f', -1, 64)
		out = append(out, '"')
	}
	return enc.WriteValue(out)
}

func (FloatMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	switch s := tok.String(); s {
	case "0":
		*receiver = T(0)
	case "NaN":
		*receiver = T(math.NaN())
	case "Infinity":
		*receiver = T(math.Inf(1))
	case "-Infinity":
		*receiver = T(math.Inf(-1))
	default:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return cj.WrapSyntaxError(dec, "invalid float", err)
		}
		*receiver = T(f)
	}
	return nil
}

type Int[T ~int | ~int8 | ~int16 | ~int32 | ~int64] struct{}

func (Int[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.Int(int64(receiver)))
}

func (Int[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '0')
	if err != nil {
		return err
	}
	*receiver = T(tok.Int())
	return nil
}

type IntMapKey[T ~int | ~int8 | ~int16 | ~int32 | ~int64] struct{}

func (IntMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(strconv.FormatInt(int64(receiver), 10)))
}

func (IntMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	i, err := strconv.ParseInt(tok.String(), 10, 64)
	if err != nil {
		return cj.WrapSyntaxError(dec, "invalid int", err)
	}
	*receiver = T(i)
	return nil
}

type RID[T rid.ResourceIdentifier] struct{}

func (RID[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(rid.ResourceIdentifier(receiver).String()))
}

func (RID[T]) Compare(a, b T) int {
	ra, rb := rid.ResourceIdentifier(a), rid.ResourceIdentifier(b)
	if cmp := strings.Compare(ra.Service, rb.Service); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(ra.Instance, rb.Instance); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(ra.Type, rb.Type); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(ra.Locator, rb.Locator); cmp != 0 {
		return cmp
	}
	return 0
}

func (RID[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	parsed, err := rid.ParseRID(tok.String())
	if err != nil {
		return cj.WrapSyntaxError(dec, "invalid resource identifier", err)
	}
	*receiver = T(parsed)
	return nil
}

type String[T ~string] struct{}

func (String[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(string(receiver)))
}

func (String[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	*receiver = T(tok.String())
	return nil
}

type StringerMarshaler[T fmt.Stringer] struct{}

func (StringerMarshaler[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(receiver.String()))
}

func (StringerMarshaler[T]) Compare(a, b T) int {
	return strings.Compare(a.String(), b.String())
}

type Uint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64] struct{}

func (Uint[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.Uint(uint64(receiver)))
}

func (Uint[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '0')
	if err != nil {
		return err
	}
	*receiver = T(tok.Uint())
	return nil
}

type UintMapKey[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64] struct{}

func (UintMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	out = append(out, '"')
	out = strconv.AppendUint(out, uint64(receiver), 10)
	out = append(out, '"')
	return enc.WriteValue(out)
}

func (UintMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	i, err := strconv.ParseUint(tok.String(), 10, 64)
	if err != nil {
		return cj.WrapSyntaxError(dec, "invalid uint", err)
	}
	*receiver = T(i)
	return nil
}

type UUID[T ~[16]byte] struct{}

func (UUID[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(uuid.UUID(receiver).String()))
}

func (t UUID[T]) Compare(a, b T) int {
	// UUIDs are 16 bytes, so we can compare them directly
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return int(a[i]) - int(b[i])
		}
	}
	return 0
}

func (UUID[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	parsed, err := uuid.ParseUUID(tok.String())
	if err != nil {
		return cj.WrapSyntaxError(dec, "invalid UUID", err)
	}
	*receiver = T(parsed)
	return nil
}

func readStringToken(dec *jsontext.Decoder) (jsontext.Token, error) {
	tok, err := dec.ReadToken()
	if err != nil {
		return tok, err
	}
	if kind := tok.Kind(); kind != '"' {
		return tok, cj.NewKindMismatchError(dec, kind, "json string")
	}
	return tok, nil
}

func readStringValue(dec *jsontext.Decoder) (jsontext.Value, error) {
	val, err := dec.ReadValue()
	if err != nil {
		return val, err
	}
	if kind := val.Kind(); kind != '"' {
		return val, cj.NewKindMismatchError(dec, kind, "json string")
	}
	return val, nil
}

func readTokenOfKind(dec *jsontext.Decoder, want ...jsontext.Kind) (jsontext.Token, error) {
	tok, err := dec.ReadToken()
	if err != nil {
		return jsontext.Token{}, err
	}
	got := tok.Kind()
	for _, k := range want {
		if got == k {
			return tok, nil
		}
	}
	return tok, cj.NewKindMismatchError(dec, got, fmt.Sprint(want))
}
