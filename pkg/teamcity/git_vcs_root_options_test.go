package teamcity

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GitVcsRootOptionsConstructor(t *testing.T) {
	assert := assert.New(t)
	propAssert := newPropertyAssertions(t)

	t.Run("Properties initialized correctly", func(t *testing.T) {
		actual, _ := NewGitVcsRootOptions("refs/heads/master", "fetch", "push", GitAuthMethodPassword, "admin", "admin")
		require.NotNil(t, actual)

		assert.Equal("refs/heads/master", actual.DefaultBranch)
		assert.Equal("fetch", actual.FetchURL)
		assert.Equal("push", actual.PushURL)
		assert.Equal("PASSWORD", string(actual.AuthMethod))
		assert.Equal("admin", actual.Username)
		assert.Equal("admin", actual.Password)

		props := actual.gitVcsRootProperties()
		propAssert.assertPropertyValue(props, "authMethod", string(actual.AuthMethod))
		propAssert.assertPropertyValue(props, "branch", actual.DefaultBranch)
		propAssert.assertPropertyValue(props, "push_url", actual.PushURL)
		propAssert.assertPropertyValue(props, "url", actual.FetchURL)
		propAssert.assertPropertyValue(props, "username", actual.Username)
		propAssert.assertPropertyValue(props, "secure:password", actual.Password)

		propAssert.assertPropertyValue(props, "agentCleanFilesPolicy", "ALL_UNTRACKED")
		propAssert.assertPropertyValue(props, "agentCleanPolicy", "ON_BRANCH_CHANGE")
		propAssert.assertPropertyValue(props, "ignoreKnownHosts", "true")
		propAssert.assertPropertyValue(props, "submoduleCheckout", "CHECKOUT")
		propAssert.assertPropertyValue(props, "useAlternates", "true")
		propAssert.assertPropertyValue(props, "usernameStyle", "USERID")
	})
	t.Run("PushURL should use FetchURL if empty", func(t *testing.T) {
		actual, _ := NewGitVcsRootOptions("refs/heads/master", "fetch", "", GitAuthMethodAnonymous, "", "")
		require.NotNil(t, actual)
		assert.Equal(actual.FetchURL, actual.PushURL)
	})
	t.Run("authMethod is required", func(t *testing.T) {
		_, err := NewGitVcsRootOptions("refs/heads/master", "fetch", "", "", "", "")
		require.EqualError(t, err, "auth is required")
	})
	t.Run("defaultBranch is required", func(t *testing.T) {
		_, err := NewGitVcsRootOptions("", "fetch", "push", GitAuthMethodAnonymous, "", "")
		require.EqualError(t, err, "defaultBranch is required")
	})
	t.Run("fetchURL is required", func(t *testing.T) {
		_, err := NewGitVcsRootOptions("refs/heads/master", "", "push", GitAuthMethodAnonymous, "", "")
		require.EqualError(t, err, "fetchURL is required")
	})
}

func Test_GitVcsRootOptionsVcsRootProperties_AnonymousAuth(t *testing.T) {
	propAssert := newPropertyAssertions(t)

	actual, _ := NewGitVcsRootOptions("refs/heads/master", "fetch", "", GitAuthMethodAnonymous, "", "")
	props := actual.gitVcsRootProperties()

	// If using anonymous, don't consider username/password properties
	propAssert.assertPropertyValue(props, "authMethod", string(GitAuthMethodAnonymous))
	propAssert.assertPropertyDoesNotExist(props, "username")
	propAssert.assertPropertyDoesNotExist(props, "secure:password")
}

func Test_GitVcsRootOptionsVcsRootProperties_UsernamePasswordAuth(t *testing.T) {
	propAssert := newPropertyAssertions(t)

	actual, _ := NewGitVcsRootOptions("refs/heads/master", "fetch", "", GitAuthMethodPassword, "admin", "admin")
	props := actual.gitVcsRootProperties()

	propAssert.assertPropertyValue(props, "authMethod", string(GitAuthMethodPassword))
	propAssert.assertPropertyValue(props, "username", actual.Username)
	propAssert.assertPropertyValue(props, "secure:password", actual.Password)
}

func Test_GitVcsRootOptionsVcsRootProperties_UsernamePasswordAuth_UsernameRequired(t *testing.T) {
	_, err := NewGitVcsRootOptions("refs/heads/master", "fetch", "", GitAuthMethodPassword, "", "admin")

	assert.Errorf(t, err, "username is required if using auth method '%s'", GitAuthMethodPassword)
}

func Test_GitVcsRootOptionsVcsRootProperties_UploadedKeyAuth(t *testing.T) {
	propAssert := newPropertyAssertions(t)

	actual, _ := NewGitVcsRootOptions("refs/heads/master", "fetch", "", GitAuthSSHUploadedKey, "admin", "admin")
	actual.PrivateKeySource = "MyUploadedKey"
	props := actual.gitVcsRootProperties()

	propAssert.assertPropertyValue(props, "authMethod", string(GitAuthSSHUploadedKey))
	propAssert.assertPropertyValue(props, "username", actual.Username)
	propAssert.assertPropertyValue(props, "teamcitySshKey", actual.PrivateKeySource)
	propAssert.assertPropertyValue(props, "secure:passphrase", actual.Password)
}

func Test_GitVcsRootOptionsVcsRootProperties_CustomKeyAuth(t *testing.T) {
	propAssert := newPropertyAssertions(t)

	actual, _ := NewGitVcsRootOptions("refs/heads/master", "fetch", "", GitAuthSSHCustomKey, "admin", "admin")
	actual.PrivateKeySource = "~/.ssh/id_rsa"
	props := actual.gitVcsRootProperties()

	propAssert.assertPropertyValue(props, "authMethod", string(GitAuthSSHCustomKey))
	propAssert.assertPropertyValue(props, "username", actual.Username)
	propAssert.assertPropertyValue(props, "privateKeyPath", actual.PrivateKeySource)
	propAssert.assertPropertyValue(props, "secure:passphrase", actual.Password)
}

func Test_GitVcsRootOptionsVcsRootProperties_BranchSpec(t *testing.T) {
	propAssert := newPropertyAssertions(t)

	actual, _ := NewGitVcsRootOptions("refs/heads/master", "fetch", "", GitAuthMethodAnonymous, "", "")
	actual.BranchSpec = []string{"+:refs/(pull/*)/head", "+:refs/heads/develop"}
	props := actual.gitVcsRootProperties()

	propAssert.assertPropertyValue(props, "teamcity:branchSpec", "+:refs/(pull/*)/head\\n+:refs/heads/develop")
}

func Test_GitVcsRootOptionsVcsRootProperties_EnableTagsInBranchSpec(t *testing.T) {
	propAssert := newPropertyAssertions(t)

	actual, _ := NewGitVcsRootOptions("refs/heads/master", "fetch", "", GitAuthMethodAnonymous, "", "")
	actual.EnableTagsInBranchSpec = true
	props := actual.gitVcsRootProperties()

	propAssert.assertPropertyValue(props, "reportTagRevisions", "true")
}

func Test_GitVcsRootOptionsVcsRootProperties_DefaultAgentSettings(t *testing.T) {
	propAssert := newPropertyAssertions(t)
	assert := assert.New(t)

	actual, _ := NewGitVcsRootOptions("refs/heads/master", "fetch", "", GitAuthMethodAnonymous, "", "")
	require.NotNil(t, actual.AgentSettings)

	s := actual.AgentSettings
	assert.Equal(s.CleanFilesPolicy, CleanFilesPolicyAllUntracked)
	assert.Equal(s.CleanPolicy, CleanPolicyBranchChange)
	assert.Equal("", s.GitPath)
	assert.Equal(true, s.UseMirrors)

	props := actual.gitVcsRootProperties()

	propAssert.assertPropertyValue(props, "agentCleanPolicy", string(s.CleanPolicy))
	propAssert.assertPropertyValue(props, "agentCleanFilesPolicy", string(s.CleanFilesPolicy))
	propAssert.assertPropertyValue(props, "useAlternates", strconv.FormatBool(s.UseMirrors))

	propAssert.assertPropertyDoesNotExist(props, "agentGitPath")
}

func Test_PropertiesToGitVcsRootOptions(t *testing.T) {
	assert := assert.New(t)

	sut := NewProperties([]*Property{
		NewProperty("branch", "refs/head/master"),
		NewProperty("reportTagRevisions", "true"),
		NewProperty("teamcity:branchSpec", "+:refs/(pull/*)/head\\n+:refs/heads/develop"),
		NewProperty("authMethod", string(GitAuthMethodPassword)),
		NewProperty("username", "admin"),
		NewProperty("submoduleCheckout", "CHECKOUT"),
		NewProperty("usernameStyle", string(GitVcsUsernameStyleUserID)),
		NewProperty("agentGitPath", "gitPath"),
		NewProperty("agentCleanPolicy", string(CleanPolicyBranchChange)),
		NewProperty("agentCleanFilesPolicy", string(CleanFilesPolicyAllUntracked)),
		NewProperty("useAlternates", "true"),
	}...)

	actual := sut.gitVcsOptions()
	require.NotNil(t, actual)

	assert.Equal("refs/head/master", actual.DefaultBranch)
	assert.ElementsMatch([]string{"+:refs/(pull/*)/head", "+:refs/heads/develop"}, actual.BranchSpec)
	assert.Equal(true, actual.EnableTagsInBranchSpec)
	assert.Equal(GitAuthMethodPassword, actual.AuthMethod)
	assert.Equal("admin", actual.Username)
	assert.Equal("CHECKOUT", actual.SubModuleCheckout)
	assert.Equal(GitVcsUsernameStyleUserID, actual.UsernameStyle)

	agentSettings := sut.gitAgentSettings()

	require.NotNil(t, agentSettings)
	assert.Equal("gitPath", agentSettings.GitPath)
	assert.Equal(CleanPolicyBranchChange, agentSettings.CleanPolicy)
	assert.Equal(CleanFilesPolicyAllUntracked, agentSettings.CleanFilesPolicy)
	assert.Equal(true, agentSettings.UseMirrors)
}
