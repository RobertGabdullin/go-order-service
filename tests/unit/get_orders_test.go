//go:build unit

package unit

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

func TestGetOrders_AssignArgs(t *testing.T) {
	mockService := new(mockStorageService)
	cmd := commands.NewGetOrds(mockService)

	t.Run("valid arguments", func(t *testing.T) {
		t.Parallel()
		args := map[string]string{"user": "1", "count": "5"}
		newCmd, err := cmd.AssignArgs(args)
		assert.NoError(t, err)
		assert.NotNil(t, newCmd)
	})

	t.Run("missing user argument", func(t *testing.T) {
		t.Parallel()
		args := map[string]string{"count": "5"}
		newCmd, err := cmd.AssignArgs(args)
		assert.Error(t, err)
		assert.Nil(t, newCmd)
	})
}

func TestGetOrders_Execute(t *testing.T) {
	t.Parallel()
	mockService := new(mockStorageService)
	cmd := commands.SetGetOrds(mockService, 1, 2)

	orders := []models.Order{
		{Id: 1, Recipient: 1, Status: "alive", Limit: time.Now().Add(24 * time.Hour)},
		{Id: 2, Recipient: 1, Status: "alive", Limit: time.Now().Add(24 * time.Hour)},
	}

	mockService.On("ListOrders", 1).Return(orders, nil)

	err := cmd.Execute(&sync.Mutex{})
	assert.NoError(t, err)

	mockService.AssertExpectations(t)
}
