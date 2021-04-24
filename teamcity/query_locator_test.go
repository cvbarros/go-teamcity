package teamcity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildQueryLocator_Empty(t *testing.T) {
	q := buildQueryLocator()
	require.Equal(t, "", q.value)
}

func TestBuildQueryLocator_OneElement(t *testing.T) {
	q := buildQueryLocator(LocatorStart(0))
	require.Equal(t, "locator=start%3A0", q.String())
}

func TestBuildQueryLocator_TwoElements(t *testing.T) {
	q := buildQueryLocator(LocatorStart(0), LocatorCount(1))
	require.Equal(t, "locator=start%3A0,count%3A1", q.String())
}
