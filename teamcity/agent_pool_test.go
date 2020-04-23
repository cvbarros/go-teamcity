package teamcity_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/assert"
)

func TestAgentPools_GetDefaultProject(t *testing.T) {
	client := setup()
	assert := assert.New(t)

	// this is hard-coded in TeamCity so we may as well do the same
	defaultAgentPoolId := 0

	retrievedPool, err := client.AgentPools.Get(defaultAgentPoolId)
	assert.NoError(err)
	assert.Equal("Default", retrievedPool.Name)
	assert.Nil(retrievedPool.MaxAgents)
	assert.True(len(retrievedPool.Projects.Project) == 1)
}

func TestAgentPools_Lifecycle(t *testing.T) {
	client := setup()
	assert := assert.New(t)

	agentPool := teamcity.CreateAgentPool{
		Name: fmt.Sprintf("test-%d", time.Now().Unix()),
	}
	createdPool, err := client.AgentPools.Create(agentPool)
	assert.NoError(err)
	assert.NotEmpty(createdPool.Id)
	assert.Equal(agentPool.Name, createdPool.Name)

	retrievedPool, err := client.AgentPools.Get(createdPool.Id)
	assert.NoError(err)
	assert.Equal(agentPool.Name, retrievedPool.Name)
	assert.Nil(retrievedPool.MaxAgents)

	assert.NoError(client.AgentPools.Delete(createdPool.Id))

	// confirm it's gone
	agentPools, err := client.AgentPools.List()
	assert.NoError(err)
	for _, pool := range agentPools.AgentPools {
		if pool.Name == agentPool.Name {
			t.Fatalf("Created agent pool still exists!")
		}
	}
}

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

	assert.True(found, "Default Agent Pool was not found")
}
