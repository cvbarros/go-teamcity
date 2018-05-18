package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
)

const (
	testProjectId          = "ProjectTest"
	testBuildTypeProjectId = "BuildTypeProjectTest"
	testVcsRootProjectId   = "VcsRootProjectTest"
	testBuildTypeId        = "BuildTypeTest"
	testParameterProjectId = "ParameterProjectId"
)

func setup() (client *teamcity.Client) {
	return teamcity.New("admin", "admin")
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
