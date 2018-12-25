package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var noExistedFile = "/tmp/not_existed_file"

func Test_FileExists(t *testing.T) {
	if !FileExists("./fileUtils.go") {
		t.Errorf("./fileUtils.go should exists, but it didn't")
	}

	if FileExists(noExistedFile) {
		t.Errorf("Weird, how could this file exists: %s", noExistedFile)
	}
}

func Test_OpenFile(t *testing.T) {
	filePath := "2018.log"
	file := OpenFile(filePath)
	assert.NotNil(t, file)
	file.Close()
	err := os.RemoveAll(filePath)
	assert.NoError(t, err)

	filePath = "./logs/2018.log"
	file = OpenFile(filePath)
	assert.NotNil(t, file)
	err = os.RemoveAll(filePath)
	assert.NoError(t, err)

	filePath = "c://temp/logs/2018.log"
	file = OpenFile(filePath)
	assert.NotNil(t, file)
	err = os.RemoveAll(filePath)
	assert.NoError(t, err)
}

func Test_CreateFile(t *testing.T) {
	filePath := "2018.log"
	err := CreateFile(filePath)
	assert.NoError(t, err)

	filePath = "./logs/2018.log"
	err = CreateFile(filePath)
	assert.NoError(t, err)
	err = os.RemoveAll(filepath.Dir(filePath))
	assert.NoError(t, err)

	filePath = "./logs1/logs2/2018.log"
	err = CreateFile(filePath)
	assert.NoError(t, err)
	err = os.RemoveAll(filepath.Dir(filePath))
	assert.NoError(t, err)
}

func Test_filepath(t *testing.T) {

}
