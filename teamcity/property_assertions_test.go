package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PropertyAssertions struct {
	a *assert.Assertions
	t *testing.T
}

func newPropertyAssertions(t *testing.T) *PropertyAssertions {
	return &PropertyAssertions{a: assert.New(t), t: t}
}

func (p *PropertyAssertions) assertPropertyValue(props *teamcity.Properties, name string, value string) {
	require.NotNil(p.t, props)

	propMap := props.Map()

	if v, ok := propMap[name]; ok {
		p.a.Equal(value, v)
	} else {
		p.a.Contains(propMap, name)
	}
}

func (p *PropertyAssertions) assertPropertyDoesNotExist(props *teamcity.Properties, name string) {
	require.NotNil(p.t, props)

	propMap := props.Map()

	p.a.NotContains(propMap, name)
}

func (p *PropertyAssertions) assertPropertyExists(props *teamcity.Properties, name string) {
	require.NotNil(p.t, props)

	propMap := props.Map()

	p.a.Contains(propMap, name)
}
