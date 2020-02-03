package teamcity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericBuildFeature_UnmarshalJSON(t *testing.T) {
	assert := assert.New(t)
	var actual GenericBuildFeature
	const json = `
	{
		"id": "DockerSupport",
		"type": "DockerSupport",
		"properties": {
			"count": 2,
			"property": [
				{
					"name": "cleanupPushed",
					"value": "true"
				},
				{
					"name": "testProperty",
					"value": "testValue"
				}
			]
		}
	}
	`
	actual.UnmarshalJSON([]byte(json))

	assert.Equal("DockerSupport", actual.ID())
	assert.Equal("DockerSupport", actual.Type())
	assert.Equal(&Property{Name: "cleanupPushed", Value: "true"}, actual.properties.Items[0])
	assert.Equal(&Property{Name: "testProperty", Value: "testValue"}, actual.properties.Items[1])
}
