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

func TestAcceptReturn_AssignArgs(t *testing.T) {
	t.Parallel()
	service := new(mockStorageService)
	cmd := commands.NewAcceptReturn(service)

	tests := []struct {
		name      string
		args      map[string]string
		expectErr bool
	}{
		{
			name: "Valid arguments",
			args: map[string]string{
				"user":  "1",
				"order": "1",
			},
		},
		{
			name: "Missing required argument",
			args: map[string]string{
				"user": "1",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cmd.AssignArgs(tt.args)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAcceptReturn_Execute(t *testing.T) {
	t.Parallel()
	service := new(mockStorageService)
	cmd := commands.SetAcceptReturn(service, 1, 1)

	orders := []models.Order{
		{Id: 1, Recipient: 1, Status: "delivered", Limit: time.Now().Add(24 * time.Hour)},
	}

	service.On("FindOrders", []int{1}).Return(orders, nil)

	err := cmd.Execute(&sync.Mutex{})
	assert.Error(t, err)
}
