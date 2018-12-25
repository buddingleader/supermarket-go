package utils

import (
	"os"
	"path/filepath"
)

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// OpenFile open the file by path, create and open it if an error occours
func OpenFile(filePath string) *os.File {
	logFile, err := os.OpenFile(filePath, os.O_APPEND, os.ModePerm)
	if err != nil {
		if err = CreateFile(filePath); err != nil {
			panic(err)
		}

		logFile, err = os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return logFile
}

// CreateFile 调用os.MkdirAll递归创建文件夹
func CreateFile(filePath string) error {
	if !FileExists(filePath) {
		return os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	}
	return nil
}
