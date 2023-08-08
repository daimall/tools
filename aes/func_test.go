package aes

import (
	"testing"
)

func TestGetPriAesKey(t *testing.T) {
	tests := []struct {
		name           string
		appid          string
		model          string
		timestamp      string
		expectedResult string
	}{
		// Test cases go here
		// For each test case, provide appid, model, timestamp, and the expected result.
		// Example:
		{"Test case 1", "1234567890", "model123", "9876543210", "af59789dea9bfaa780286d3be2cdc7fe"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPriAesKey(tt.appid, tt.model, tt.timestamp)

			if result != tt.expectedResult {
				t.Errorf("Result mismatch, got %v, expected %v", result, tt.expectedResult)
			}
		})
	}
}

func TestGetSha1(t *testing.T) {
	tests := []struct {
		name           string
		src            string
		expectedResult string
	}{
		// Test cases go here
		// For each test case, provide src and the expected result.
		// Example:
		{"Test case 1", "hello", "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := GetSha1(tt.src)

			if result != tt.expectedResult {
				t.Errorf("Result mismatch, got %v, expected %v", result, tt.expectedResult)
			}
		})
	}
}
