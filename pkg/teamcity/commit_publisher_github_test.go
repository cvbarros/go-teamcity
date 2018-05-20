package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeatureCommitPublisher_GithubRequiredProperties(t *testing.T) {
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
		opt := teamcity.StatusPublisherGithubOptions{AuthenticationType: "password", Host: "https://api.github.com", Password: "1234"}
		_, err := teamcity.NewFeatureCommitStatusPublisherGithub(opt)
		require.EqualError(t, err, "username/password required for auth type 'password'")
	})
	t.Run("Password Required", func(t *testing.T) {
		opt := teamcity.StatusPublisherGithubOptions{AuthenticationType: "password", Host: "https://api.github.com", Username: "bob"}
		_, err := teamcity.NewFeatureCommitStatusPublisherGithub(opt)
		require.EqualError(t, err, "username/password required for auth type 'password'")
	})

	t.Run("Correct Properties", func(t *testing.T) {
		assert := assert.New(t)
		opt := teamcity.NewCommitStatusPublisherGithubOptionsPassword("https://api.github.com", "bob", "1234")

		actual := opt.Properties().Map()

		assert.Equal("githubStatusPublisher", actual["publisherId"])
		assert.Equal("https://api.github.com", actual["github_host"])
		assert.Equal("password", actual["github_authentication_type"])
		assert.Equal("bob", actual["github_username"])
		assert.Equal("1234", actual["github_password"])
	})
}

func TestFeatureCommitPublisher_GithubAuthenticationToken(t *testing.T) {
	t.Run("AccessToken Required", func(t *testing.T) {
		opt := teamcity.StatusPublisherGithubOptions{AuthenticationType: "token", Host: "https://api.github.com"}
		_, err := teamcity.NewFeatureCommitStatusPublisherGithub(opt)
		require.EqualError(t, err, "accesstoken required for auth type 'token'")
	})

	t.Run("Correct Properties", func(t *testing.T) {
		assert := assert.New(t)
		opt := teamcity.NewCommitStatusPublisherGithubOptionsToken("https://api.github.com", "1234")

		actual := opt.Properties().Map()

		assert.Equal("githubStatusPublisher", actual["publisherId"])
		assert.Equal("https://api.github.com", actual["github_host"])
		assert.Equal("token", actual["github_authentication_type"])
		assert.Equal("1234", actual["github_access_token"])
	})
}
