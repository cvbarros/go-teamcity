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

func Test_LocatorTypeProvider(t *testing.T) {
	sut := LocatorTypeProvider("OAuthProvider", "teamcity-vault")
	actual := sut.String()

	expected := "type%3AOAuthProvider%2Cproperty%28name%3AproviderType%2Cvalue%3Ateamcity-vault%29"

	assert.Equal(t, expected, actual)
}
