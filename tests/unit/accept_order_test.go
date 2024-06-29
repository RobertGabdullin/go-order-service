//go:build unit

package unit

import (
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

func TestAcceptOrder_AssignArgs(t *testing.T) {
	service := new(mockStorageService)
	cmd := commands.NewAcceptOrd(service)

	tests := []struct {
		name      string
		args      map[string]string
		expectErr bool
	}{
		{
			name: "Valid arguments",
			args: map[string]string{
				"user":      "1",
				"order":     "1",
				"weight":    "5",
				"basePrice": "100",
				"expire":    "2024-06-05T10",
				"wrapper":   "pack",
			},
		},
		{
			name: "Missing required argument",
			args: map[string]string{
				"user":      "1",
				"order":     "1",
				"weight":    "5",
				"basePrice": "100",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := cmd.AssignArgs(tt.args)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAcceptOrder_Execute(t *testing.T) {
	t.Parallel()
	service := new(mockStorageService)

	service.On("GetWrapper", "none").Return(models.Wrapper{MaxWeight: sql.NullInt64{5, true}, Id: 1, Markup: 5, Type: "pack"}, nil)
	service.On("AddOrder", mock.AnythingOfType("models.Order")).Return(nil)

	cmd := commands.SetAcceptOrd(service, 1, 1, 5, 100, time.Now().Add(24*time.Hour), "none")

	err := cmd.Execute(&sync.Mutex{})
	assert.NoError(t, err)

	service.AssertExpectations(t)
}
