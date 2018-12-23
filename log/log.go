package log

import (
	"path/filepath"

	"github.com/wangff15386/supermarket-go/common/utils"

	"github.com/sirupsen/logrus"
	"github.com/wangff15386/supermarket-go/conf"
)

// GetLogger get a logger to record log
func GetLogger(module string) *logrus.Entry {
	log := logrus.New()

	// log.Out = os.Stdout
	logpath := filepath.Join(conf.Config.Log.Path, module+".log")
	log.Out = utils.OpenFile(logpath)

	lvl, err := logrus.ParseLevel(conf.Config.Log.Level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	log.Level = lvl

	return log.WithField("module", module)
}
