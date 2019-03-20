package visitors

import (
	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/conjure/types"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestExternalTypeFallback(t *testing.T) {

	t.Run("External Fallback", func(t *testing.T) {
		def := spec.ExternalReference{
			ExternalReference: spec.TypeName{
				Name:    "Foo",
				Package: "com.example.foo",
			},
			Fallback: spec.NewTypeFromPrimitive(spec.PrimitiveTypeString),
		}

		provider := newExternalVisitor(def)

		info := types.NewPkgInfo("", nil)
		typ, err := provider.ParseType(info)

		assert.Equal(t, err, nil)
		assert.Equal(t, types.String, typ)
	})

}
