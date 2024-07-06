package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func MapEqual(t *testing.T, expected, actual map[string]string) {
	assert.Len(t, expected, len(actual))
	for key, elem := range expected {
		require.Contains(t, actual, key)
		assert.Equal(t, actual[key], elem)
	}
}
