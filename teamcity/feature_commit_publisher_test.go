package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/assert"
)

func TestFeatureCommitPublisher_UnmarshallProperties_Github(t *testing.T) {
	assert := assert.New(t)
	var actual teamcity.FeatureCommitStatusPublisher
	const json = `
	{
		"id": "BUILD_EXT_1",
		"type": "commit-status-publisher",
		"properties": {
			"count": 6,
			"property": [
				{
					"name": "github_authentication_type",
					"value": "password"
				},
				{
					"name": "github_host",
					"value": "https://api.github.com"
				},
				{
					"name": "github_username",
					"value": "me@me.com"
				},
				{
					"name": "publisherId",
					"value": "githubStatusPublisher"
				},
				{
					"name": "secure:github_access_token"
				},
				{
					"name": "secure:github_password"
				},
				{
					"name": "vcsRootId",
					"value": "Project_VcsRootId"
				}
			]
		}
	}
	`
	actual.UnmarshalJSON([]byte(json))

	assert.Equal("BUILD_EXT_1", actual.ID())
	assert.Equal("Project_VcsRootId", actual.VcsRootID())
	assert.IsType(new(teamcity.StatusPublisherGithubOptions), actual.Options)
	actualOpt := actual.Options.(*teamcity.StatusPublisherGithubOptions)

	assert.Equal("password", actualOpt.AuthenticationType)
	assert.Equal("me@me.com", actualOpt.Username)
}
