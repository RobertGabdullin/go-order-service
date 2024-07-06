//go:build unit

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/r_gabdullin/homework-1/tests"
)

type mockParser struct{}

func (m *mockParser) Parse(line string) (string, map[string]string, error) {
	return "acceptOrd", map[string]string{
		"user":      "1",
		"order":     "1",
		"weight":    "5",
		"basePrice": "100",
		"expire":    "2025-06-05T10",
		"wrapper":   "pack",
	}, nil
}

func TestCLI_Find(t *testing.T) {
	t.Parallel()
	service := &tests.MockStorageService{}
	parser := &mockParser{}
	cli := NewCLI(service, parser)

	cmd, err := cli.Find("acceptOrd")
	assert.NoError(t, err)
	assert.Equal(t, "acceptOrd", cmd.GetName())

	cmd, err = cli.Find("unknownCmd")
	assert.Error(t, err)
	assert.Nil(t, cmd)
}
