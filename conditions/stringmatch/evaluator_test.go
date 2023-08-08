package stringmatch

import (
	"strings"
	"testing"
)

func TestCalculate(t *testing.T) {
	// Mock boolean function for testing
	boolFunc := func(ident string) bool {
		// You can implement your own logic here to return boolean values
		// based on the identifier.
		// For testing purposes, you can use a map to provide predefined values.
		// Example:
		// predefinedValues := map[string]bool{"identifier1": true, "identifier2": false}
		// return predefinedValues[ident]
		return !strings.Contains(ident, "NOK")
	}

	tests := []struct {
		name           string
		expression     string
		stackSize      int
		expectedResult bool
		expectedError  bool
	}{
		{
			name:           "Test Calculate with AND and OR operations",
			expression:     "OK AND X OR Y",
			stackSize:      10,
			expectedResult: true, // Replace with the expected result
			expectedError:  false,
		},
		{
			name:           "Test Calculate with AND and OR operations",
			expression:     "OK AND (NOK OR Y)",
			stackSize:      10,
			expectedResult: true, // Replace with the expected result
			expectedError:  false,
		},
		{
			name:           "Test Calculate with AND and OR operations",
			expression:     "NOK AND (X OR Y)",
			stackSize:      10,
			expectedResult: false, // Replace with the expected result
			expectedError:  false,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Calculate(tt.expression, tt.stackSize, boolFunc)

			if (err != nil) != tt.expectedError {
				t.Errorf("Error mismatch, got error %v, expected error %v", err, tt.expectedError)
			}

			if result != tt.expectedResult {
				t.Errorf("Result mismatch, got %v, expected %v", result, tt.expectedResult)
			}
		})
	}
}
