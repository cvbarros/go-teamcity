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

func TestClient_TokenAuth(t *testing.T) {
	t.Run("Token auth works against local instance", func(t *testing.T) {
		//This token was created for user 'admin' on the pre-warmed integration testing server. It is named 'integration_tests'
		client, err := teamcity.NewClient(
			teamcity.TokenAuth("eyJ0eXAiOiAiVENWMiJ9.bWZ3QWswa1ViWk5CUFlrRC1GQUVYQkM1cmZz.ODViNDA2MDctZmFkNS00YTc0LTlmYTgtM2MwMzkxMmY2ZGY5"),
			http.DefaultClient)
		if err != nil {
			t.Fatalf("Error when connecting to server: %s", err)
		}
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
