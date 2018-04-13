package tests

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk"
	"github.com/stretchr/testify/assert"
)

var (
	client   *teamcity.Client
	username string
	password string
)

func init() {
	client = initTest()
}

func initTest() *teamcity.Client {
	username = "admin"
	password = "admin"

	return teamcity.NewClient(username, password)
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
