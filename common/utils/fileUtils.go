package utils

import (
	"io"
	"os"
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

// OpenFile open the file by filepath, create it if an error occours
func OpenFile(filepath string) io.Writer {
	logFile, err := os.OpenFile(filepath, os.O_APPEND, 0666)
	if err != nil {
		logFile, err = os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666) // CREATE FILE
		if err != nil {
			panic(err)
		}
	}
	return logFile
}
