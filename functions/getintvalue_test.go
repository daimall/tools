package functions

import (
	"testing"
)

func TestGetIntV(t *testing.T) {
	testCases := []struct {
		Input       interface{}
		ExpectedInt int
		ExpectedErr bool
	}{
		{123, 123, false},         // 测试整数类型
		{3.14, 3, false},          // 测试浮点数类型
		{"456", 456, false},       // 测试字符串数字
		{"abc", 0, true},          // 测试非数字字符串
		{true, 0, true},           // 测试布尔值
		{uint32(789), 789, false}, // 测试无符号整数类型
		{int64(999), 999, false},  // 测试有符号整数类型
	}

	for _, testCase := range testCases {
		t.Run("GetIntV", func(t *testing.T) {
			result, err := GetIntV(testCase.Input)
			if err != nil && !testCase.ExpectedErr {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != testCase.ExpectedInt {
				t.Errorf("Expected %d, but got %d", testCase.ExpectedInt, result)
			}
		})
	}
}
