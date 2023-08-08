package collection

import (
	"testing"
)

func TestMap(t *testing.T) {
	// Test cases for Map function
	tests := []struct {
		name     string
		input    []interface{}
		mapping  func(interface{}) interface{}
		expected []interface{}
	}{
		{
			name:     "Test Map with square mapping",
			input:    []interface{}{1, 2, 3, 4},
			mapping:  func(x interface{}) interface{} { return x.(int) * x.(int) },
			expected: []interface{}{1, 4, 9, 16},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.input, tt.mapping)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, but got %d", len(tt.expected), len(result))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("At index %d, expected %v, but got %v", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

func TestReduce(t *testing.T) {
	// Test cases for Reduce function
	tests := []struct {
		name     string
		input    []interface{}
		reducer  func(interface{}, interface{}) interface{}
		initial  interface{}
		expected interface{}
	}{
		{
			name:     "Test Reduce with sum reducer",
			input:    []interface{}{1, 2, 3, 4},
			reducer:  func(acc interface{}, x interface{}) interface{} { return acc.(int) + x.(int) },
			initial:  0,
			expected: 10,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reduce(tt.input, tt.reducer, tt.initial)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
