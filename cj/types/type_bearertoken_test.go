package types_test

import (
	"testing"

	"github.com/palantir/conjure-go/v6/cj/types"
	"github.com/palantir/pkg/bearertoken"
)

func TestBearerToken(t *testing.T) {
	for name, test := range map[string]typeTest{
		"bearertoken": typeTestCase[bearertoken.Token, types.BearerToken[bearertoken.Token], types.BearerToken[bearertoken.Token]]{
			Value: "foo", JSON: "\"foo\"",
		},
		"null": typeTestCase[bearertoken.Token, types.BearerToken[bearertoken.Token], types.BearerToken[bearertoken.Token]]{
			JSON: "null", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "KindMismatchError at 4: want json string, got null",
		},
		"invalid": typeTestCase[bearertoken.Token, types.BearerToken[bearertoken.Token], types.BearerToken[bearertoken.Token]]{
			JSON: "\" \"", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "InvalidValueError at 3: invalid character for bearer token",
		},
		"empty": typeTestCase[bearertoken.Token, types.BearerToken[bearertoken.Token], types.BearerToken[bearertoken.Token]]{
			JSON: "\"\"", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "InvalidValueError at 2: empty bearer token",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}
