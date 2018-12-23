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
	MongoURL string
	AppName  string
	AppPath  string
	HTTPURL  string
	WorkPath string
	LogPath  string
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
func Initial() {
	var err error
	if Config.AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}

	if Config.WorkPath, err = os.Getwd(); err != nil {
		panic(err)
	}
	appConfigPath := filepath.Join(Config.WorkPath, filename)
	fmt.Println("WorkPath:", appConfigPath)
	if !utils.FileExists(appConfigPath) {
		appConfigPath = filepath.Join(Config.AppPath, filename)
		fmt.Println("AppPath:", appConfigPath)
		if !utils.FileExists(appConfigPath) {
			panic("Cannot found the app.conf in WorkPath and AppPath")
		}
	}

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
}

func newAppConf() *AppConfig {
	return &AppConfig{
		AppName: "supermarket-go",
	}
}

func saveAppConf(conf *AppConfig) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println("conf:", string(b))
	return ioutil.WriteFile(filename, b, 0666)
}
