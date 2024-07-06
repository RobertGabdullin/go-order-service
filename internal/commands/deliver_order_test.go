//go:build unit

package commands

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/tests"
)

func TestDeliverOrder_AssignArgs(t *testing.T) {
	t.Parallel()
	cmd := NewDeliverOrd(nil)

	testTable := []struct {
		name      string
		args      map[string]string
		expectErr bool
		errString string
		expected  deliverOrder
	}{
		{
			name: "Valid arguments",
			args: map[string]string{
				"orders": "[1,2,3]",
			},
			expected: SetDeliverOrd(nil, []int{1, 2, 3}),
		},
		{
			name: "Missing required argument",
			args: map[string]string{
				"user": "1",
			},
			expectErr: true,
			errString: "missing orders flag",
		},
		{
			name: "Invalid number of flags",
			args: map[string]string{
				"orders": "[1,2,3]",
				"user":   "1",
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
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDeliverOrder_Execute(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)

	service.On("FindOrders", mock.AnythingOfType("[]int")).Return([]models.Order{
		{Id: 1, Recipient: 1, Status: "alive", Expire: time.Now().Add(24 * time.Hour)},
		{Id: 2, Recipient: 1, Status: "alive", Expire: time.Now().Add(24 * time.Hour)},
	}, nil)
	service.On("ChangeStatus", mock.AnythingOfType("int"), "delivered", mock.AnythingOfType("string")).Return(nil)

	cmd := SetDeliverOrd(service, []int{1, 2})

	err := cmd.Execute(&sync.Mutex{})
	assert.NoError(t, err)

	service.AssertExpectations(t)
}

func TestDeliverOrder_Execute_OrderNotFound(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)

	service.On("FindOrders", mock.AnythingOfType("[]int")).Return(nil, errors.New("order not found"))

	cmd := SetDeliverOrd(service, []int{1, 2})

	err := cmd.Execute(&sync.Mutex{})
	tests.ErrorContains(t, err, "order not found")

	service.AssertCalled(t, "FindOrders", []int{1, 2})
	service.AssertExpectations(t)
}

func TestDeliverOrder_Execute_InvalidUser(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)

	service.On("FindOrders", mock.AnythingOfType("[]int")).Return([]models.Order{
		{Id: 1, Recipient: 1, Status: "alive", Expire: time.Now().Add(24 * time.Hour)},
		{Id: 2, Recipient: 2, Status: "alive", Expire: time.Now().Add(24 * time.Hour)},
	}, nil)

	cmd := SetDeliverOrd(service, []int{1, 2})

	err := cmd.Execute(&sync.Mutex{})
	assert.Error(t, err)
	tests.ErrorContains(t, err, "list of orders should belong only to one person")

	service.AssertCalled(t, "FindOrders", []int{1, 2})
	service.AssertExpectations(t)
}

func TestDeliverOrder_Execute_OrderNotAvailable(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)

	service.On("FindOrders", mock.AnythingOfType("[]int")).Return([]models.Order{
		{Id: 1, Recipient: 1, Status: "delivered", Expire: time.Now().Add(24 * time.Hour)},
		{Id: 2, Recipient: 1, Status: "delivered", Expire: time.Now().Add(24 * time.Hour)},
	}, nil)

	cmd := SetDeliverOrd(service, []int{1, 2})

	err := cmd.Execute(&sync.Mutex{})
	assert.Error(t, err)
	tests.ErrorContains(t, err, "some orders are not available")

	service.AssertCalled(t, "FindOrders", []int{1, 2})
	service.AssertExpectations(t)
}

func TestDeliverOrder_Execute_OrderOutOfStorageLimit(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)

	service.On("FindOrders", mock.AnythingOfType("[]int")).Return([]models.Order{
		{Id: 1, Recipient: 1, Status: "alive", Expire: time.Now().Add(-24 * time.Hour)},
	}, nil)

	cmd := SetDeliverOrd(service, []int{1})

	err := cmd.Execute(&sync.Mutex{})
	assert.Error(t, err)
	tests.ErrorContains(t, err, "some orders is out of storage limit date")

	service.AssertCalled(t, "FindOrders", []int{1})
	service.AssertExpectations(t)
}
