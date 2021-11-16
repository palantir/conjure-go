package imports

import (
	"testing"

	"github.com/palantir/conjure-go/v6/integration_test/testgenerated/imports/pkg3/api"
)

// TestImportNameConflict ensures that we can import a type which includes types of the same package name.
// This is only testing that the relevant code compiles, and does not do anything in the test body.
func TestImportNameConflict(t *testing.T) {
	_ = api.Union{}
}
