package teamcity_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testProjectId          = "ProjectTest"
	testBuildTypeProjectId = "BuildTypeProjectTest"
	testVcsRootProjectId   = "VcsRootProjectTest"
	testBuildTypeId        = "BuildTypeTest"
	testParameterProjectId = "ParameterProjectId"
)

func setup() *teamcity.Client {
	client, _ := teamcity.NewClient(teamcity.BasicAuth("admin", "admin"), http.DefaultClient)
	return client
}

func TestClient_BasicAuth(t *testing.T) {
	t.Run("Basic auth works against local instance", func(t *testing.T) {
		client := setup()
		success, err := client.Validate()
		if err != nil {
			t.Fatalf("Error when validating client: %s", err)
		}

		assert.Equal(t, true, success)
	})
}

func TestClient_Address(t *testing.T) {
	t.Run("Specify address from alternate constructor", func(t *testing.T) {
		address := os.Getenv("TEAMCITY_ADDR")
		client, err := teamcity.NewClientWithAddress(teamcity.BasicAuth("admin", "admin"), address, http.DefaultClient)
		require.NoError(t, err)
		success, err := client.Validate()
		if err != nil {
			t.Fatalf("Error when validating client: %s", err)
		}

		assert.Equal(t, true, success)
	})
}
