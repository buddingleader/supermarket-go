package utils

import (
	"testing"
)

var noExistedFile = "/tmp/not_existed_file"

func TestFileExists(t *testing.T) {
	if !FileExists("./fileUtils.go") {
		t.Errorf("./fileUtils.go should exists, but it didn't")
	}

	if FileExists(noExistedFile) {
		t.Errorf("Weird, how could this file exists: %s", noExistedFile)
	}
}
