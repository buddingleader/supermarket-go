package log

import (
	"log"
	"path"
	"strings"

	"github.com/wangff15386/supermarket-go/common/utils"
	"github.com/wangff15386/supermarket-go/conf"
)

// GetLogger get a logger to record log
func GetLogger(fileName string, level string) *log.Logger {
	if ext := path.Ext(fileName); !strings.Contains(ext, "log") {
		fileName += ".log"
	}

	logFile := utils.OpenFile(conf.Config.LogPath + fileName)
	return log.New(logFile, level, log.LstdFlags|log.Llongfile)
}
