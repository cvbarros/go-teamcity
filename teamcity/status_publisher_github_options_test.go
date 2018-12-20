package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const githubHost = "https://api.github.com"

func TestFeatureCommitPublisher_Invariants(t *testing.T) {
	t.Run("AuthenticationType Required", func(t *testing.T) {
		opt := teamcity.StatusPublisherGithubOptions{}
		_, err := teamcity.NewFeatureCommitStatusPublisherGithub(opt)
		assert.Error(t, err)
	})
	t.Run("AuthenticationType Valid", func(t *testing.T) {
		opt := teamcity.StatusPublisherGithubOptions{AuthenticationType: "anything"}
		_, err := teamcity.NewFeatureCommitStatusPublisherGithub(opt)
		assert.Error(t, err)
	})
	t.Run("Host Required", func(t *testing.T) {
		opt := teamcity.StatusPublisherGithubOptions{AuthenticationType: "password"}
		_, err := teamcity.NewFeatureCommitStatusPublisherGithub(opt)
		require.EqualError(t, err, "Host is required")
	})
}

func TestFeatureCommitPublisher_GithubAuthenticationPassword(t *testing.T) {
	t.Run("Username Required", func(t *testing.T) {
		opt := teamcity.StatusPublisherGithubOptions{AuthenticationType: "password", Host: githubHost, Password: "1234"}
		_, err := teamcity.NewFeatureCommitStatusPublisherGithub(opt)
		require.EqualError(t, err, "username/password required for auth type 'password'")
	})
	t.Run("Password Required", func(t *testing.T) {
		opt := teamcity.StatusPublisherGithubOptions{AuthenticationType: "password", Host: githubHost, Username: "bob"}
		_, err := teamcity.NewFeatureCommitStatusPublisherGithub(opt)
		require.EqualError(t, err, "username/password required for auth type 'password'")
	})

	t.Run("Correct Properties", func(t *testing.T) {
		assert := assert.New(t)
		opt := teamcity.NewCommitStatusPublisherGithubOptionsPassword(githubHost, "bob", "1234")

		actual := opt.Properties().Map()

		assert.Equal("githubStatusPublisher", actual["publisherId"])
		assert.Equal(githubHost, actual["github_host"])
		assert.Equal("password", actual["github_authentication_type"])
		assert.Equal("bob", actual["github_username"])
		assert.Equal("1234", actual["secure:github_password"])
	})
}

func TestFeatureCommitPublisher_GithubAuthenticationToken(t *testing.T) {
	t.Run("AccessToken Required", func(t *testing.T) {
		opt := teamcity.StatusPublisherGithubOptions{AuthenticationType: "token", Host: githubHost}
		_, err := teamcity.NewFeatureCommitStatusPublisherGithub(opt)
		require.EqualError(t, err, "accesstoken required for auth type 'token'")
	})

	t.Run("Correct Properties", func(t *testing.T) {
		assert := assert.New(t)
		opt := teamcity.NewCommitStatusPublisherGithubOptionsToken(githubHost, "1234")

		actual := opt.Properties().Map()

		assert.Equal("githubStatusPublisher", actual["publisherId"])
		assert.Equal(githubHost, actual["github_host"])
		assert.Equal("token", actual["github_authentication_type"])
		assert.Equal("1234", actual["secure:github_access_token"])
	})
}

func TestFeatureCommitPublisher_GithubFromProperties(t *testing.T) {
	props := teamcity.NewProperties([]*teamcity.Property{
		teamcity.NewProperty("github_host", githubHost),
		teamcity.NewProperty("github_authentication_type", "password"),
		teamcity.NewProperty("github_username", "bob"),
		teamcity.NewProperty("secure:github_password", "1234"),
		teamcity.NewProperty("secure:github_access_token", "1234"),
	}...)

	t.Run("From Password AuthType", func(t *testing.T) {
		assert := assert.New(t)

		actual, err := teamcity.CommitStatusPublisherGithubOptionsFromProperties(props)

		require.NoError(t, err)
		assert.Equal("password", actual.AuthenticationType)
		assert.Equal(githubHost, actual.Host)
		assert.Equal("bob", actual.Username)
		assert.Empty(actual.Password)
		assert.Empty(actual.AccessToken)
	})

	t.Run("From Token AuthType", func(t *testing.T) {
		assert := assert.New(t)
		props.AddOrReplaceValue("github_authentication_type", "token")
		actual, err := teamcity.CommitStatusPublisherGithubOptionsFromProperties(props)

		require.NoError(t, err)
		assert.Equal("token", actual.AuthenticationType)
		assert.Equal(githubHost, actual.Host)
		assert.Empty(actual.Username)
		assert.Empty(actual.Password)
		assert.Empty(actual.AccessToken)
	})
}
