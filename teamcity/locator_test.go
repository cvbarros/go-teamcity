package teamcity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LocatorNameWithSpaces(t *testing.T) {
	sut := LocatorName("<Root Project>")
	actual := sut.String()

	assert.Equal(t, "name%3A%3CRoot%20Project%3E", actual)
}

func Test_LocatorId(t *testing.T) {
	sut := LocatorID("_Root")
	actual := sut.String()

	assert.Equal(t, "id%3A_Root", actual)
}
