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

func TestReturnOrder_AssignArgs(t *testing.T) {
	t.Parallel()
	mockService := new(mockStorageService)
	cmd := commands.NewReturnOrd(mockService)

	t.Run("valid arguments", func(t *testing.T) {
		args := map[string]string{"order": "1"}
		newCmd, err := cmd.AssignArgs(args)
		assert.NoError(t, err)
		assert.NotNil(t, newCmd)
	})

	t.Run("missing order argument", func(t *testing.T) {
		args := map[string]string{}
		newCmd, err := cmd.AssignArgs(args)
		assert.Error(t, err)
		assert.Nil(t, newCmd)
	})
}

func TestReturnOrder_Execute(t *testing.T) {
	t.Parallel()
	mockService := new(mockStorageService)
	cmd := commands.SetReturnOrd(mockService, 1)

	order := models.Order{
		Id: 1, Status: "alive", Limit: time.Now().Add(-24 * time.Hour),
	}

	mockService.On("FindOrders", []int{1}).Return([]models.Order{order}, nil)
	mockService.On("DeleteOrder", 1).Return(nil)

	err := cmd.Execute(&sync.Mutex{})
	assert.NoError(t, err)

	mockService.AssertExpectations(t)
}
