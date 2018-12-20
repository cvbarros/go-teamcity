package teamcity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_Get(t *testing.T) {
	client := setup()
	server, err := client.Server.Get()
	if err != nil {
		t.Fatalf("Failed to GetServer: %s", err)
	}

	if server == nil {
		t.Fatalf("GetServer did not return a server or serialization failure.")
	}

	assert.Equal(t, int32(2018), server.VersionMajor)
	assert.NotEmpty(t, server.WebURL)
}
