package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveEmptyInterfaces(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []interface{}
	}{
		{
			name:     "empty slice",
			input:    []interface{}{},
			expected: []interface{}{},
		},
		{
			name:     "no nil values",
			input:    []interface{}{1, "test", true},
			expected: []interface{}{1, "test", true},
		},
		{
			name:     "all nil values",
			input:    []interface{}{nil, nil, nil},
			expected: []interface{}{},
		},
		{
			name:     "mixed values",
			input:    []interface{}{1, nil, "test", nil, true, nil},
			expected: []interface{}{1, "test", true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveEmptyInterfaces(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
