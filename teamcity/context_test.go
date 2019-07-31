package teamcity_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/cvbarros/go-teamcity/teamcity"
)

type TestContext struct {
	Prefix string
	Client *teamcity.Client
	T      *testing.T
}

func NewTc(prefix string, t *testing.T) *TestContext {
	return &TestContext{
		Prefix: prefix,
		T:      t,
		Client: setup(),
	}
}

func (tc *TestContext) RandomName() string {
	reseed()
	return fmt.Sprintf("%s-%d", tc.Prefix, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}

func reseed() {
	rand.Seed(time.Now().UTC().UnixNano())
}

var BuildTypeContextOptionsDefault = BuildTypeContextOptions{
	AttachVcsRoot: false,
}

type BuildTypeContextOptions struct {
	AttachVcsRoot bool
}
type BuildTypeContext struct {
	TC        *TestContext
	BuildType *teamcity.BuildType
	Project   *teamcity.Project
	VcsRoot   *teamcity.VcsRootReference

	ready bool
}

func (b *BuildTypeContext) Setup(t *TestContext) {
	b.SetupWithOpt(t, BuildTypeContextOptionsDefault)
}

func (b *BuildTypeContext) SetupWithOpt(t *TestContext, opt BuildTypeContextOptions) {
	b.TC = t
	b.Project = createTestProject(t.T, t.Client, t.RandomName())
	b.BuildType = b.NewBuildType()
	b.ready = true

	if opt.AttachVcsRoot {
		gitVcs := getTestVcsRootData(b.Project.ID).(*teamcity.GitVcsRoot)
		created, _ := t.Client.VcsRoots.Create(b.Project.ID, gitVcs)
		b.VcsRoot = created

		t.Client.BuildTypes.AttachVcsRoot(b.BuildType.ID, created)
	}
}

func (b *BuildTypeContext) Teardown() {
	if b.ready {
		cleanUpProject(b.TC.T, b.TC.Client, b.Project.ID)
	}
}

func (b *BuildTypeContext) NewBuildType() *teamcity.BuildType {
	return createTestBuildTypeWithName(b.TC.T, b.TC.Client, b.Project.ID, b.TC.RandomName(), false)
}
