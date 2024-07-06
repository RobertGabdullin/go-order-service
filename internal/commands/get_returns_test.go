//go:build unit

package commands

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/tests"
)

func TestGetReturns_AssignArgs(t *testing.T) {
	t.Parallel()
	cmd := NewGetReturns(nil)

	testTable := []struct {
		name      string
		args      map[string]string
		expectErr bool
		errString string
		expected  getReturns
	}{
		{
			name: "Valid arguments",
			args: map[string]string{
				"offset": "1",
				"limit":  "1",
			},
			expected: SetGetReturns(nil, 1, 1),
		},
		{
			name: "Valid arguments without offset",
			args: map[string]string{
				"offset": "1",
			},
			expected: SetGetReturns(nil, 1, -1),
		},
		{
			name: "Invalid value for offset",
			args: map[string]string{
				"offset": "-123",
			},
			expectErr: true,
			errString: "invalid flag value",
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

func TestGetReturns_Execute(t *testing.T) {
	t.Parallel()
	service := new(tests.MockStorageService)
	cmd := SetGetReturns(service, 0, 2)

	returns := []models.Order{
		{Id: 1, Recipient: 1, Status: "returned"},
		{Id: 2, Recipient: 1, Status: "returned"},
	}

	service.On("GetReturns", 0, 2).Return(returns, nil)

	err := cmd.Execute(&sync.Mutex{})
	assert.NoError(t, err)

	service.AssertCalled(t, "GetReturns", 0, 2)
	service.AssertExpectations(t)
}
