package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/json-iterator/go"
	"github.com/wangff15386/supermarket-go/common/utils"
)

// AppConfig the config of app
type AppConfig struct {
	Mongo     *MongoConfig
	LevelPath string
	AppName   string
	AppPath   string
	HTTPURL   string
	WorkPath  string
	Log       *LogConfig
}

// MongoConfig for mongodb
type MongoConfig struct {
	URL     string
	TimeOut int64 // Second
}

// LogConfig for log
type LogConfig struct {
	Path  string
	Level string
}

// Config the gloable config
var (
	Config   *AppConfig
	filename = "app.conf"
)

func init() {
	Config = newAppConf()
}

// Initial config
func Initial(appConfigPath string) {
	if !utils.FileExists(appConfigPath) {
		var err error
		if Config.AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
			panic(err)
		}

		if Config.WorkPath, err = os.Getwd(); err != nil {
			panic(err)
		}
		appConfigPath = filepath.Join(Config.WorkPath, filename)
		fmt.Println("WorkPath:", appConfigPath)
		if !utils.FileExists(appConfigPath) {
			appConfigPath = filepath.Join(Config.AppPath, filename)
			fmt.Println("AppPath:", appConfigPath)
			if !utils.FileExists(appConfigPath) {
				panic("Cannot found the app.conf in WorkPath and AppPath")
			}
		}
	}

	fmt.Println("appConfigPath:", appConfigPath)
	b, err := ioutil.ReadFile(appConfigPath)
	if err != nil {
		panic(err)
	}

	if b != nil && len(b) > 0 {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		if err = json.Unmarshal(b, &Config); err != nil {
			panic(err)
		}
	}

	if err := formatAppConf(Config, appConfigPath); err != nil {
		fmt.Println("Failed to format config:", err)
	}
}

func newAppConf() *AppConfig {
	return &AppConfig{
		AppName:   "supermarket-go",
		LevelPath: "/leveldb",
		Mongo:     &MongoConfig{},
		Log:       &LogConfig{},
	}
}

func formatAppConf(conf *AppConfig, appConfigPath string) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(appConfigPath, b, 0666)
}
