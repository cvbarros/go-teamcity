package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/assert"
)

// Ensure serialization/deserialization works as expected.
func TestSerialize(t *testing.T) {
	// Create a StepOctopusCreateRelease object.
	step, _ := teamcity.NewStepOctopusCreateRelease("Test step")
	step.Host = "web-14.smith.info"
	step.ApiKey = "DfkDxZSbSAIpblvdvcTv"
	step.OctopusServerVersion = "3.0+"
	step.Project = "Project"
	step.ReleaseNumber = "1.0.0"
	step.ChannelName = "Stage"
	step.Environments = "Stage"
	step.Tenants = "TenantA"
	step.TenantTags = ""
	step.WaitForDeployments = true
	step.AdditionalCommandLineArguments = ""

	// Serialize the step object.
	jsonStep, err := step.MarshalJSON()
	assert.Nil(t, err)
	assert.NotNil(t, jsonStep)

	// Then deserialize it.
	deserializeStep, _ := teamcity.NewStepOctopusCreateRelease("Deserialize test step")
	err = deserializeStep.UnmarshalJSON(jsonStep)
	assert.Nil(t, err)

	//Ensure it is equal to the original object.
	assert.Equal(t, step, deserializeStep)
}
