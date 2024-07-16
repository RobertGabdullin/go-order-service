//go:build unit

package commands

import (
	"database/sql"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/tests"
)

func TestAcceptOrder_AssignArgs(t *testing.T) {
	t.Parallel()
	cmd := NewAcceptOrd(nil)

	expireTime, _ := time.Parse("2006-01-02T15", "2024-06-05T10")

	testTable := []struct {
		name      string
		args      map[string]string
		expectErr bool
		errString string
		expected  AcceptOrder
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
			expected: SetAcceptOrd(nil, 1, 1, 5, 100, expireTime, "pack"),
		},
		{
			name: "Missing required argument",
			args: map[string]string{
				"user":      "1",
				"order":     "1",
				"weight":    "5",
				"basePrice": "100",
				"wrapper":   "pack",
			},
			expectErr: true,
			errString: "missing expire flag",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := cmd.AssignArgs(tt.args)
			if tt.expectErr {
				assert.Error(t, err)
				tests.ErrorContains(t, err, tt.errString)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestAcceptOrder_Execute(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)

	service.On("GetWrapper", "none").Return(models.Wrapper{MaxWeight: sql.NullInt64{Int64: 5, Valid: true}, Id: 1, Markup: 5, Type: "pack"}, nil)
	service.On("AddOrder", mock.AnythingOfType("models.Order")).Return(nil)

	cmd := SetAcceptOrd(service, 1, 1, 5, 100, time.Now().Add(24*time.Hour), "none")

	_, err := cmd.Execute(&sync.Mutex{})
	assert.NoError(t, err)

	service.AssertCalled(t, "GetWrapper", "none")
	service.AssertCalled(t, "AddOrder", mock.AnythingOfType("models.Order"))
	service.AssertExpectations(t)
}

func TestAcceptOrder_Execute_InvalidWrapper(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)

	service.On("GetWrapper", "none").Return(models.Wrapper{}, errors.New("wrapper not found"))

	cmd := SetAcceptOrd(service, 1, 1, 5, 100, time.Now().Add(24*time.Hour), "none")

	_, err := cmd.Execute(&sync.Mutex{})
	tests.ErrorContains(t, err, "wrapper not found")

	service.AssertCalled(t, "GetWrapper", "none")
	service.AssertNotCalled(t, "AddOrder", mock.AnythingOfType("models.Order"))
	service.AssertExpectations(t)
}

func TestAcceptOrder_Execute_InvalidOrder(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)

	service.On("GetWrapper", "pack").Return(models.Wrapper{MaxWeight: sql.NullInt64{Int64: 5, Valid: true}, Id: 1, Markup: 5, Type: "pack"}, nil)
	service.On("AddOrder", mock.AnythingOfType("models.Order")).Return(errors.New("order not valid"))

	cmd := SetAcceptOrd(service, 1, 1, 5, 100, time.Now().Add(24*time.Hour), "pack")

	_, err := cmd.Execute(&sync.Mutex{})
	tests.ErrorContains(t, err, "order not valid")

	service.AssertCalled(t, "GetWrapper", "pack")
	service.AssertCalled(t, "AddOrder", mock.AnythingOfType("models.Order"))
	service.AssertExpectations(t)
}
