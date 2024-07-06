//go:build unit

package commands

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/tests"
)

func TestAcceptReturn_AssignArgs(t *testing.T) {
	t.Parallel()
	cmd := NewAcceptReturn(nil)

	testTable := []struct {
		name      string
		args      map[string]string
		expectErr bool
		errString string
		expected  acceptReturn
	}{
		{
			name: "Valid arguments",
			args: map[string]string{
				"user":  "1",
				"order": "1",
			},
			expected: SetAcceptReturn(nil, 1, 1),
		},
		{
			name: "Missing required argument",
			args: map[string]string{
				"user": "1",
			},
			expectErr: true,
			errString: "invalid number of flags",
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
				assert.Equal(t, tt.expected, result)
				assert.NoError(t, err)
			}
		})
	}
}

func TestAcceptReturn_Execute(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)
	cmd := SetAcceptReturn(service, 1, 1)

	orders := []models.Order{
		{Id: 1, Recipient: 1, Status: "delivered", Expire: time.Now().Add(24 * time.Hour), DeliveredAt: time.Now().Add(-24 * time.Hour)},
	}

	service.On("FindOrders", []int{1}).Return(orders, nil)
	service.On("ChangeStatus", 1, "returned", mock.AnythingOfType("string")).Return(nil)

	err := cmd.Execute(&sync.Mutex{})
	assert.NoError(t, err)

	service.AssertCalled(t, "FindOrders", []int{1})
	service.AssertCalled(t, "ChangeStatus", 1, "returned", mock.AnythingOfType("string"))
	service.AssertExpectations(t)
}

func TestAcceptReturn_Execute_InvalidDeliverTime(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)
	cmd := SetAcceptReturn(service, 1, 1)

	orders := []models.Order{
		{Id: 1, Recipient: 1, Status: "delivered", Expire: time.Now().Add(24 * time.Hour), DeliveredAt: time.Now().Add(-50 * time.Hour)},
	}

	service.On("FindOrders", []int{1}).Return(orders, nil)

	err := cmd.Execute(&sync.Mutex{})
	assert.Error(t, err)
	tests.ErrorContains(t, err, "the order can only be returned within two days after issue")

	service.AssertCalled(t, "FindOrders", []int{1})
	service.AssertNotCalled(t, "ChangeStatus")
	service.AssertExpectations(t)
}

func TestAcceptReturn_Execute_InvalidTypeOrder(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)
	cmd := SetAcceptReturn(service, 1, 1)

	orders := []models.Order{
		{Id: 1, Recipient: 1, Status: "alive", Expire: time.Now().Add(24 * time.Hour), DeliveredAt: time.Now().Add(-24 * time.Hour)},
	}

	service.On("FindOrders", []int{1}).Return(orders, nil)

	err := cmd.Execute(&sync.Mutex{})
	tests.ErrorContains(t, err, "such an order has never been issued")

	service.AssertCalled(t, "FindOrders", []int{1})
	service.AssertNotCalled(t, "ChangeStatus")
	service.AssertExpectations(t)
}
