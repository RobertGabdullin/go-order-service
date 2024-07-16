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

func TestGetOrders_AssignArgs(t *testing.T) {
	t.Parallel()
	cmd := NewGetOrds(nil)

	testTable := []struct {
		name      string
		args      map[string]string
		expectErr bool
		errString string
		expected  getOrders
	}{
		{
			name: "Valid arguments with count",
			args: map[string]string{
				"user":  "3",
				"count": "5",
			},
			expected: SetGetOrds(nil, 3, 5),
		},
		{
			name: "Valid arguments without count",
			args: map[string]string{
				"user": "17",
			},
			expected: SetGetOrds(nil, 17, -1),
		},
		{
			name: "Missing user argument",
			args: map[string]string{
				"count": "123",
			},
			expectErr: true,
			errString: "missing user flag",
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

func TestGetOrders_Execute(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)
	cmd := SetGetOrds(service, 1, 2)

	orders := []models.Order{
		{Id: 1, Recipient: 1, Status: "alive", Expire: time.Now().Add(24 * time.Hour)},
		{Id: 2, Recipient: 1, Status: "alive", Expire: time.Now().Add(24 * time.Hour)},
	}

	service.On("ListOrders", 1).Return(orders, nil)

	_, err := cmd.Execute(&sync.Mutex{})
	assert.NoError(t, err)

	service.AssertCalled(t, "ListOrders", 1)
	service.AssertExpectations(t)
}
