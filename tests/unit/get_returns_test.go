//go:build unit

package unit

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

func TestGetReturns_AssignArgs(t *testing.T) {
	mockService := new(mockStorageService)
	cmd := commands.NewGetReturns(mockService)

	t.Run("valid arguments", func(t *testing.T) {
		t.Parallel()
		args := map[string]string{"offset": "10", "limit": "20"}
		newCmd, err := cmd.AssignArgs(args)
		assert.NoError(t, err)
		assert.NotNil(t, newCmd)
	})

	t.Run("invalid offset argument", func(t *testing.T) {
		t.Parallel()
		args := map[string]string{"offset": "-1", "limit": "20"}
		newCmd, err := cmd.AssignArgs(args)
		assert.Error(t, err)
		assert.Nil(t, newCmd)
	})
}

func TestGetReturns_Execute(t *testing.T) {
	t.Parallel()
	mockService := new(mockStorageService)
	cmd := commands.SetGetReturns(mockService, 0, 2)

	returns := []models.Order{
		{Id: 1, Recipient: 1, Status: "returned"},
		{Id: 2, Recipient: 1, Status: "returned"},
	}

	mockService.On("GetReturns", 0, 2).Return(returns, nil)

	err := cmd.Execute(&sync.Mutex{})
	assert.NoError(t, err)

	mockService.AssertExpectations(t)
}
