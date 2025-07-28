package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrderID(t *testing.T) {
	cases := []struct {
		name        string
		value       int64
		expected    OrderID
		expectedErr error
	}{
		{
			name:        "valid",
			value:       9278923470,
			expected:    OrderID{9278923470},
			expectedErr: nil,
		},
		{
			name:        "invalid",
			value:       12345,
			expected:    OrderID{0},
			expectedErr: ErrOrderID,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v, err := NewOrderID(c.value)
			assert.Equal(t, c.expected, v)
			assert.Equal(t, c.expectedErr, err)
		})
	}
}
