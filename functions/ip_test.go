package functions

import (
	"testing"
)

// 测试IsLocalIP函数
func TestIsLocalIP(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{"127.0.0.1", true}, // 测试本地IP
		{"8.8.8.8", false},  // 测试非本地IP
	}

	for _, testCase := range testCases {
		t.Run(testCase.Input, func(t *testing.T) {
			result, err := IsLocalIP(testCase.Input)
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			if result != testCase.Expected {
				t.Errorf("Expected %v, but got %v", testCase.Expected, result)
			}
		})
	}
}

// 测试IsIntranet函数
func TestIsIntranet(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{"192.168.1.1", true}, // 测试局域网IP
		{"8.8.8.8", false},    // 测试非内网IP
	}

	for _, testCase := range testCases {
		t.Run(testCase.Input, func(t *testing.T) {
			result := IsIntranet(testCase.Input)
			if result != testCase.Expected {
				t.Errorf("Expected %v, but got %v", testCase.Expected, result)
			}
		})
	}
}

// 测试checkIp函数
func TestCheckIp(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{"127.0.0.1", true},  // 测试有效IP
		{"256.0.0.1", false}, // 测试无效IP
		{"notanip", false},   // 测试非IP字符串
	}

	for _, testCase := range testCases {
		t.Run(testCase.Input, func(t *testing.T) {
			result := checkIp(testCase.Input)
			if result != testCase.Expected {
				t.Errorf("Expected %v, but got %v", testCase.Expected, result)
			}
		})
	}
}

// 测试inetAton函数
func TestInetAton(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected int64
	}{
		{"192.168.1.1", 3232235777}, // 测试IPv4
		{"8.8.8.8", 134744072},      // 测试IPv4
	}

	for _, testCase := range testCases {
		t.Run(testCase.Input, func(t *testing.T) {
			result := inetAton(testCase.Input)
			if result != testCase.Expected {
				t.Errorf("Expected %v, but got %v", testCase.Expected, result)
			}
		})
	}
}
