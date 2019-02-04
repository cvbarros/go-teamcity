package teamcity_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	teamcity "github.com/cvbarros/go-teamcity-sdk/teamcity"
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

func Test_ParameterConfiguration_FullSpec_Serialization(t *testing.T) {
	assert := assert.New(t)
	sut, _ := teamcity.NewParameter(teamcity.ParameterTypes.Configuration, "param1", "value1")
	sut.ControlType = "password"
	sut.Description = "some description"
	sut.Display = "prompt"
	sut.Label = "some label"
	sut.ReadOnly = "true"
	jsonBytes, err := sut.MarshalJSON()

	require.NoError(t, err)
	require.Equal(t, string(jsonBytes), `{"name":"param1","type":{"rawValue":"password display='prompt' description='some description' readOnly='true' label='some label'"},"value":"value1"}`)
	actual := &teamcity.Parameter{}
	if err := json.Unmarshal([]byte(jsonBytes), &actual); err != nil {
		t.Error(err)
	}

	assert.Equal("param1", actual.Name)
	assert.Equal("value1", actual.Value)
	assert.Equal("password", actual.ControlType)
	assert.Equal("some description", actual.Description)
	assert.Equal("prompt", actual.Display)
	assert.Equal("true", actual.ReadOnly)
	assert.Equal("some label", actual.Label)
	assert.Equal(string(teamcity.ParameterTypes.Configuration), actual.Type)
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

func Test_ParameterCollection_Serialization(t *testing.T) {
	sut := teamcity.NewParametersEmpty()

	sut.AddOrReplaceValue(teamcity.ParameterTypes.Configuration, "config", "value_config")
	sut.AddOrReplaceValue(teamcity.ParameterTypes.System, "system", "value_system")
	sut.AddOrReplaceValue(teamcity.ParameterTypes.EnvironmentVariable, "env", "value_env")

	jsonBytes, err := json.Marshal(sut)
	require.NoError(t, err)
	require.Equal(t, `{"count":3,"property":[{"name":"config","value":"value_config"},{"name":"system.system","value":"value_system"},{"name":"env.env","value":"value_env"}]}`, string(jsonBytes))
}

func Test_ParameterConvertToProperty(t *testing.T) {
	assert := assert.New(t)
	sut, _ := teamcity.NewParameter(teamcity.ParameterTypes.Configuration, "name", "value_config")
	actual := sut.Property()

	assert.Equal("name", actual.Name)
	assert.Equal("value_config", actual.Value)

	sut, _ = teamcity.NewParameter(teamcity.ParameterTypes.System, "name", "value_system")
	actual = sut.Property()

	assert.Equal("system.name", actual.Name)
	assert.Equal("value_system", actual.Value)

	sut, _ = teamcity.NewParameter(teamcity.ParameterTypes.EnvironmentVariable, "name", "value_env")
	actual = sut.Property()

	assert.Equal("env.name", actual.Name)
	assert.Equal("value_env", actual.Value)
}
