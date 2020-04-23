package teamcity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentPools_List(t *testing.T) {
	client := setup()
	assert := assert.New(t)

	agentPools, err := client.AgentPools.List()
	assert.NoError(err)

	// whilst other pools may have been added by other tests - the Default pool
	// cannot be removed, so can be used as test data
	assert.True(len(agentPools.AgentPools) > 0, "At least one agent pool should exist")

	found := false
	for _, pool := range agentPools.AgentPools {
		if pool.Name == "Default" {
			found = true
		}
	}

	assert.True(found, "Default agent pool was not found")
}
