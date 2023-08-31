package functions

import (
	"testing"
)

func TestStr2Uint(t *testing.T) {
	negative := -123
	testCases := []struct {
		Input       string
		ExpectedID  uint
		ExpectedErr bool
	}{
		{"123", 123, false},               // 测试正常字符串转换
		{"0", 0, false},                   // 测试0字符串转换
		{"4294967295", 4294967295, false}, // 测试最大值字符串转换
		{"-123", uint(negative), true},    // 测试负数字符串转换
		{"abc", 0, true},                  // 测试非数字字符串转换
	}

	for _, testCase := range testCases {
		t.Run(testCase.Input, func(t *testing.T) {
			result, err := Str2Uint(testCase.Input)
			if err != nil && !testCase.ExpectedErr {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != testCase.ExpectedID {
				t.Errorf("Expected %d, but got %d", testCase.ExpectedID, result)
			}
		})
	}
}
