///go:build unit

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLI_Find(t *testing.T) {
	t.Parallel()
	cli := NewCLI(nil, nil, nil)

	cmd, err := cli.Find("acceptOrd")
	assert.NoError(t, err)
	assert.Equal(t, "acceptOrd", cmd.GetName())

	cmd, err = cli.Find("unknownCmd")
	assert.Error(t, err)
	assert.Nil(t, cmd)
}
