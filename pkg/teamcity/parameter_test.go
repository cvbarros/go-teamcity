package teamcity_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/require"
)

func Test_ParameterConfiguration_Serialization(t *testing.T) {
	assert := assert.New(t)
	sut, _ := teamcity.NewParameter(teamcity.ParameterTypes.Configuration, "param1", "value1")
	jsonBytes, err := sut.MarshalJSON()

	require.NoError(t, err)
	require.Equal(t, string(jsonBytes), `{"name":"param1","value":"value1"}`)
	actual := &teamcity.Parameter{}
	if err := json.Unmarshal([]byte(jsonBytes), &actual); err != nil {
		t.Error(err)
	}

	assert.Equal("param1", actual.Name)
	assert.Equal("value1", actual.Value)
	assert.Equal(string(teamcity.ParameterTypes.Configuration), actual.Type)
}

func Test_ParameterSystem_Serialization(t *testing.T) {
	assert := assert.New(t)
	sut, _ := teamcity.NewParameter(teamcity.ParameterTypes.System, "param1", "value1")
	jsonBytes, err := sut.MarshalJSON()

	require.NoError(t, err)
	require.Equal(t, string(jsonBytes), `{"name":"system.param1","value":"value1"}`)

	actual := &teamcity.Parameter{}
	if err := json.Unmarshal([]byte(jsonBytes), &actual); err != nil {
		t.Error(err)
	}

	assert.Equal("param1", actual.Name)
	assert.Equal("value1", actual.Value)
	assert.Equal(string(teamcity.ParameterTypes.System), actual.Type)
}

func Test_ParameterEnvironmentVariable_Serialization(t *testing.T) {
	assert := assert.New(t)
	sut, _ := teamcity.NewParameter(teamcity.ParameterTypes.EnvironmentVariable, "param1", "value1")
	jsonBytes, err := sut.MarshalJSON()

	require.NoError(t, err)
	require.Equal(t, string(jsonBytes), `{"name":"env.param1","value":"value1"}`)

	actual := &teamcity.Parameter{}
	if err := json.Unmarshal([]byte(jsonBytes), &actual); err != nil {
		t.Error(err)
	}

	assert.Equal("param1", actual.Name)
	assert.Equal("value1", actual.Value)
	assert.Equal(string(teamcity.ParameterTypes.EnvironmentVariable), actual.Type)
}
