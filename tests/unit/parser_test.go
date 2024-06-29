//go:build unit

package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/parser"
)

func TestGetArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		expected  map[string]string
		expectErr bool
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
		},
		{
			name:      "Missing flag value",
			args:      []string{"--key="},
			expectErr: true,
		},
	}

	parser := parser.ArgsParser{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := parser.GetArgs(tt.args)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedCmd  string
		expectedArgs map[string]string
		expectErr    bool
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
		},
	}

	parser := parser.ArgsParser{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd, args, err := parser.Parse(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCmd, cmd)
				assert.Equal(t, tt.expectedArgs, args)
			}
		})
	}
}
