package teamcity_test

import (
	"testing"

	u "github.com/cvbarros/go-teamcity-sdk/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func init() {
	u.InitTest()
}

func TestBasicAuth(t *testing.T) {
	t.Run("Basic auth works against local instance", func(t *testing.T) {
		success, err := u.Client.Validate()
		if err != nil {
			t.Fatalf("Error when validating client: %s", err)
		}

		assert.Equal(t, true, success)
	})
}
