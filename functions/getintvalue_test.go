package functions

import (
	"errors"
	"testing"
)

func TestGetIntV(t *testing.T) {
	tests := []struct {
		input interface{}
		want  int
		err   error
	}{
		// Test cases covering different types
		{int(42), 42, nil},                         // Test case for int input
		{int8(42), 42, nil},                        // Test case for int8 input
		{int16(42), 42, nil},                       // Test case for int16 input
		{int32(42), 42, nil},                       // Test case for int32 input
		{int64(42), 42, nil},                       // Test case for int64 input
		{float32(42.5), 42, nil},                   // Test case for float32 input
		{float64(42.5), 42, nil},                   // Test case for float64 input
		{uint(42), 42, nil},                        // Test case for uint input
		{uint8(42), 42, nil},                       // Test case for uint8 input
		{uint16(42), 42, nil},                      // Test case for uint16 input
		{uint32(42), 42, nil},                      // Test case for uint32 input
		{uint64(42), 42, nil},                      // Test case for uint64 input
		{"42", 42, nil},                            // Test case for string input
		{"invalid", 0, errors.New("unknown type")}, // Test case for unknown string input

		// Add more test cases here if needed
	}

	for _, test := range tests {
		got, err := GetIntV(test.input)
		if got != test.want || (err == nil && test.err != nil) || (err != nil && test.err == nil) {
			t.Errorf("GetIntV(%v) = %v, %v, want %v, %v", test.input, got, err, test.want, test.err)
		}
	}
}
