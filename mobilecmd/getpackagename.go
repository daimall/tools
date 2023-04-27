package mobilecmd

import (
	"errors"
	"strings"
)

// 获取包名称
func GetPackagePath(packageName string) (string, error) {
	pmPathOutput, err := Command{
		Args:  []string{"pm", "path", packageName},
		Shell: true,
	}.CombinedOutputString()
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(pmPathOutput, "package:") {
		return "", errors.New("invalid pm path output: " + pmPathOutput)
	}
	packagePath := strings.TrimSpace(pmPathOutput[len("package:"):])
	return packagePath, nil
}
