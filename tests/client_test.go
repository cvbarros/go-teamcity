package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	client = initTest()
}

func TestBasicAuth(t *testing.T) {
	t.Run("Basic auth works against local instance", func(t *testing.T) {
		success, err := client.Validate()
		if err != nil {
			t.Fatalf("Error when validating client: %s", err)
		}

		assert.Equal(t, true, success)
	})
}
