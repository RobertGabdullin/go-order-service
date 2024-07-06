//go:build unit

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/r_gabdullin/homework-1/tests"
)

func TestGetArgs(t *testing.T) {
	t.Parallel()
	testTable := []struct {
		name      string
		args      []string
		expected  map[string]string
		expectErr bool
		errString string
	}{
		{
			name:     "Valid short flags with =",
			args:     []string{"-key=value"},
			expected: map[string]string{"key": "value"},
		},
		{
			name:     "Valid long flags with =",
			args:     []string{"--key=value"},
			expected: map[string]string{"key": "value"},
		},
		{
			name:     "Valid short flags without =",
			args:     []string{"-key", "value"},
			expected: map[string]string{"key": "value"},
		},
		{
			name:     "Valid long flags without =",
			args:     []string{"--key", "value"},
			expected: map[string]string{"key": "value"},
		},
		{
			name:      "Invalid flag format",
			args:      []string{"key", "value"},
			expectErr: true,
			errString: "invalid flag",
		},
		{
			name:      "Missing flag value",
			args:      []string{"--key="},
			expectErr: true,
			errString: "invalid flag",
		},
	}

	parser := ArgsParser{}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := parser.GetArgs(tt.args)
			if tt.expectErr {
				tests.ErrorContains(t, err, tt.errString)
			} else {
				assert.NoError(t, err)
				tests.MapEqual(t, tt.expected, result)
			}
		})
	}
}

func TestParse(t *testing.T) {
	t.Parallel()
	testTable := []struct {
		name         string
		input        string
		expectedCmd  string
		expectedArgs map[string]string
		expectErr    bool
		errString    string
	}{
		{
			name:         "Valid input with short flags",
			input:        "command -key=value",
			expectedCmd:  "command",
			expectedArgs: map[string]string{"key": "value"},
		},
		{
			name:         "Valid input with long flags",
			input:        "command --key=value",
			expectedCmd:  "command",
			expectedArgs: map[string]string{"key": "value"},
		},
		{
			name:      "Empty input",
			input:     "",
			expectErr: true,
			errString: "empty line",
		},
		{
			name:         "Input with only command",
			input:        "command",
			expectedCmd:  "command",
			expectedArgs: map[string]string{},
		},
		{
			name:      "Input with invalid flags",
			input:     "command key value",
			expectErr: true,
			errString: "invalid flag",
		},
	}

	parser := ArgsParser{}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd, args, err := parser.Parse(tt.input)
			if tt.expectErr {
				tests.ErrorContains(t, err, tt.errString)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCmd, cmd)
				tests.MapEqual(t, tt.expectedArgs, args)
			}
		})
	}
}
