package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wangff15386/supermarket-go/conf"
)

func Test_GetLogger(t *testing.T) {
	conf.Initial("../conf/app.conf")
	path := ""
	conf.Config.Log.Path = path

	logger := GetLogger("test")
	assert.NotNil(t, logger)

	logger.Println("hello world!")
	logger.Traceln("hello world!")
	logger.Debugln("hello world!")
	logger.Infoln("hello world!")
	logger.Warnln("hello world!")
	logger.Warningln("hello world!")
	logger.Errorln("hello world!")
	// logger.Fatalln("hello world!")
	// logger.Panicln("hello world!")

	path = "/logs"
	conf.Config.Log.Path = path
	logger = GetLogger("test1")
	assert.NotNil(t, logger)
}

func Test_Open(t *testing.T) {
	path := "./log/0181224.log"
	_, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	assert.NoError(t, err)
}
