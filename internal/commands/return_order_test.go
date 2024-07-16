//go:build unit

package commands

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/tests"
)

func TestReturnOrder_AssignArgs(t *testing.T) {
	t.Parallel()
	cmd := NewReturnOrd(nil)

	testTable := []struct {
		name      string
		args      map[string]string
		expectErr bool
		errString string
		expected  returnOrders
	}{
		{
			name: "Valid arguments",
			args: map[string]string{
				"order": "1",
			},
			expected: SetReturnOrd(nil, 1),
		},
		{
			name:      "Missing order argument",
			args:      map[string]string{},
			expectErr: true,
			errString: "invalid number of flags",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := cmd.AssignArgs(tt.args)
			if tt.expectErr {
				tests.ErrorContains(t, err, tt.errString)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestReturnOrder_Execute(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)
	cmd := SetReturnOrd(service, 1)

	order := models.Order{
		Id: 1, Status: "alive", Expire: time.Now().Add(-24 * time.Hour),
	}

	service.On("FindOrders", []int{1}).Return([]models.Order{order}, nil)
	service.On("DeleteOrder", 1).Return(nil)

	_, err := cmd.Execute(&sync.Mutex{})
	assert.NoError(t, err)

	service.AssertCalled(t, "FindOrders", []int{1})
	service.AssertCalled(t, "DeleteOrder", 1)
	service.AssertExpectations(t)
}

func TestReturnOrder_Execute_OrderNotFound(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)
	cmd := SetReturnOrd(service, 1)

	service.On("FindOrders", []int{1}).Return([]models.Order{}, nil)

	_, err := cmd.Execute(&sync.Mutex{})
	tests.ErrorContains(t, err, "such order does not exist")

	service.AssertCalled(t, "FindOrders", []int{1})
	service.AssertNotCalled(t, "DeleteOrder")
	service.AssertExpectations(t)
}

func TestReturnOrder_Execute_OrderNotInStorage(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)
	cmd := SetReturnOrd(service, 1)

	order := models.Order{
		Id: 1, Status: "delivered", Expire: time.Now().Add(-24 * time.Hour),
	}

	service.On("FindOrders", []int{1}).Return([]models.Order{order}, nil)

	_, err := cmd.Execute(&sync.Mutex{})
	tests.ErrorContains(t, err, "order is not at storage")

	service.AssertCalled(t, "FindOrders", []int{1})
	service.AssertNotCalled(t, "DeleteOrder")
	service.AssertExpectations(t)
}

func TestReturnOrder_Execute_OrderNotExpired(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)
	cmd := SetReturnOrd(service, 1)

	order := models.Order{
		Id: 1, Status: "alive", Expire: time.Now().Add(24 * time.Hour),
	}

	service.On("FindOrders", []int{1}).Return([]models.Order{order}, nil)

	_, err := cmd.Execute(&sync.Mutex{})
	tests.ErrorContains(t, err, "order should be out of storage limit")

	service.AssertCalled(t, "FindOrders", []int{1})
	service.AssertNotCalled(t, "DeleteOrder")
	service.AssertExpectations(t)
}
