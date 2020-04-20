package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/suite"
)

type SuiteBuildFeature struct {
	suite.Suite
	TC               *TestContext
	BuildTypeContext *BuildTypeContext
	BuildTypeID      string
	VcsRootID        string

	Github *teamcity.FeatureCommitStatusPublisher
}

func NewSuiteBuildFeature(t *testing.T) *SuiteBuildFeature {
	return &SuiteBuildFeature{TC: NewTc("SuiteBuildFeature", t), BuildTypeContext: new(BuildTypeContext)}
}

func (suite *SuiteBuildFeature) SetupTest() {
	suite.BuildTypeContext.SetupWithOpt(suite.TC, BuildTypeContextOptions{AttachVcsRoot: true})
	suite.BuildTypeID = suite.BuildTypeContext.BuildType.ID
	suite.Require().NotNil(suite.BuildTypeContext.VcsRoot)
	suite.VcsRootID = suite.BuildTypeContext.VcsRoot.ID

	opt := teamcity.NewCommitStatusPublisherGithubOptionsToken("https://api.github.com", "1234")
	suite.Github, _ = teamcity.NewFeatureCommitStatusPublisherGithub(opt, suite.BuildTypeContext.VcsRoot.ID)
}

func (suite *SuiteBuildFeature) TearDownTest() {
	suite.BuildTypeContext.Teardown()
}

func (suite *SuiteBuildFeature) Service() *teamcity.BuildFeatureService {
	return suite.TC.Client.BuildFeatureService(suite.BuildTypeContext.BuildType.ID)
}

func (suite *SuiteBuildFeature) TestCommitPublisher_Create() {
	sut := suite.Service()
	actual, err := sut.Create(suite.Github)
	suite.Require().NoError(err)

	suite.Require().IsType(new(teamcity.FeatureCommitStatusPublisher), actual)

	csp := actual.(*teamcity.FeatureCommitStatusPublisher)

	suite.NotEqual("", csp.ID())
	suite.Equal(suite.BuildTypeID, csp.BuildTypeID())
	suite.Equal("commit-status-publisher", csp.Type())
	suite.Equal(false, csp.Disabled())
	suite.Equal(suite.VcsRootID, csp.VcsRootID())
}

func (suite *SuiteBuildFeature) TestCommitPublisher_Get() {
	sut := suite.Service()
	actual, err := sut.Create(suite.Github)
	suite.Require().NoError(err)

	actual, err = sut.GetByID(actual.ID())
	suite.Require().NoError(err)
	suite.Require().IsType(new(teamcity.FeatureCommitStatusPublisher), actual)

	csp := actual.(*teamcity.FeatureCommitStatusPublisher)

	suite.NotEqual("", csp.ID())
	suite.Equal(suite.BuildTypeID, csp.BuildTypeID())
	suite.Equal("commit-status-publisher", csp.Type())
	suite.Equal(false, csp.Disabled())
	suite.Equal(suite.VcsRootID, csp.VcsRootID())
}

func TestBuildFeatureSuite(t *testing.T) {
	suite.Run(t, NewSuiteBuildFeature(t))
}
