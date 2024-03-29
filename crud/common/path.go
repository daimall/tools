package common

import (
	"os"
	"path/filepath"
)

// getCurrPath 获取当前路径
func getCurrPath() (currpath string) {
	var err error
	if currpath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic("get current path from os.Args[0] failed")
	}
	return
}

// GetUploadPath ...
func GetUploadPath() (path string) {
	path = filepath.Join(getCurrPath(), "data", "temp", "upload")
	return
}

// GetTplPath
func GetTplPath() (path string) {
	path = filepath.Join(getCurrPath(), "data", "template")
	return
}

// GetLogsPath
func GetLogsPath() (path string) {
	path = filepath.Join(getCurrPath(), "logs")
	return
}

// GetPath 获取相对当前路径的通用路径
func GetPath(l []string) (path string) {
	nl := []string{}
	nl = append(nl, getCurrPath())
	nl = append(nl, l...)
	path = filepath.Join(nl...)
	return
}

func GetRootPath() (rootpath string) {
	return getCurrPath()
}
