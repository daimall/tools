package functions

import (
	"fmt"
	"testing"
)

// 结构体用于记录函数的所有输入和输出
type TestInputOutput struct {
	Input  int
	Output string
}

func TestCreateRandomNumber(t *testing.T) {
	testCases := []TestInputOutput{
		{1, CreateRandomNumber(1)},   // 测试生成1位随机数字
		{5, CreateRandomNumber(5)},   // 测试生成5位随机数字
		{10, CreateRandomNumber(10)}, // 测试生成10位随机数字
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Input: %d", testCase.Input), func(t *testing.T) {
			result := CreateRandomNumber(testCase.Input)
			if len(result) != testCase.Input {
				t.Errorf("Expected length of %d, but got length of %d", testCase.Input, len(result))
				return
			}
			if !IsNumeric(result) {
				t.Errorf("result is not numeric, %s", result)
				return
			}
			// 可以根据需要添加其他断言，确保生成的数字符合要求
			if result == testCase.Output {
				t.Errorf("Not Expected %s, but got %s", testCase.Output, result)
				return
			}
		})
	}
}

func TestCreateRandomString(t *testing.T) {
	testCases := []TestInputOutput{
		{1, CreateRandomString(1)},   // 测试生成1位随机字符串
		{5, CreateRandomString(5)},   // 测试生成5位随机字符串
		{10, CreateRandomString(10)}, // 测试生成10位随机字符串
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Input: %d", testCase.Input), func(t *testing.T) {
			result := CreateRandomString(testCase.Input)
			if len(result) != testCase.Input {
				t.Errorf("Expected length of %d, but got length of %d", testCase.Input, len(result))
			}
			// 可以根据需要添加其他断言，确保生成的字符串符合要求
			if result == testCase.Output {
				t.Errorf("Not Expected %s, but got %s", testCase.Output, result)
				return
			}
		})
	}
}
