package functions

import (
	"testing"
)

// 结构体用于记录函数的所有输入和输出
type PathTestInputOutput struct {
	Input  string
	Output []string
}

func TestPathSplit(t *testing.T) {
	testCases := []PathTestInputOutput{
		{"path/to/file", []string{"path", "to", "file"}}, // 测试常见路径
		{"another/path", []string{"another", "path"}},    // 测试另一个路径
		{"", []string{}},               // 测试空路径
		{"single", []string{"single"}}, // 测试单个名称
	}

	for _, testCase := range testCases {
		t.Run(testCase.Input, func(t *testing.T) {
			result := PathSplit(testCase.Input)
			if len(result) != len(testCase.Output) {
				t.Errorf("Expected %d segments, but got %d", len(testCase.Output), len(result))
			}
			for i, segment := range testCase.Output {
				if result[i] != segment {
					t.Errorf("Expected segment '%s', but got '%s'", segment, result[i])
				}
			}
		})
	}
}
