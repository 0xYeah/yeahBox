package app_cfg

import (
	"encoding/json"
	"errors"
	"github.com/0xYeah/yeahBox/base_app/api/api_config"
	"github.com/george012/gtbox"
	"github.com/george012/gtbox/gtbox_app"
	"github.com/george012/gtbox/gtbox_log"
	"os"
	"path/filepath"
)

type AppType string

const (
	AppTypeAgent  AppType = "agent"
	AppTypeServer AppType = "server"
)

var (
	CurrentApp *ExtendApp
)

type ExtendApp struct {
	*gtbox_app.App
	NetListenPortStratumDefault int
	ApiPort                     int
	AppType                     AppType
}

func NewApp(appName string, version string, bundleID string, description string, runMode gtbox.RunMode) *ExtendApp {
	app := &ExtendApp{
		App:                         gtbox_app.NewApp(appName, version, bundleID, description, runMode),
		NetListenPortStratumDefault: 0,
		ApiPort:                     apiPortDefault,
	}

	return app
}

const (
	apiPortDefault = 17173
)

var GlobalConfig *FileConfig

type FileConfig struct {
	AppType       AppType              `json:"app_type"`
	Api           api_config.ApiConfig `yaml:"api" json:"api"`
	BaseUploadDir string               `yaml:"upload_dir" json:"upload_dir"`
}

func LoadConfig(file string) error {
	fInfo, err := os.Stat(file)
	if err != nil {
		return err
	}
	if fInfo.IsDir() {
		return errors.New("config file can not be a dir")
	}

	buf, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &GlobalConfig)
	//err = yaml.Unmarshal(buf, &GlobalConfig)
	if err != nil {
		return err
	}

	return nil
}

func SaveConfig(file string) error {
	if file == "" {
		file = CurrentApp.AppConfigFilePath
	}
	//config, err := yaml.Marshal(GlobalConfig)
	config, err := json.MarshalIndent(GlobalConfig, "", "    ")

	if err != nil {
		return err
	}

	err = os.WriteFile(file, config, 0644)
	if err != nil {
		return err
	}

	return nil
}

func generateDefaultConfig() *FileConfig {
	aApiPort := apiPortDefault
	switch CurrentApp.AppType {
	case AppTypeAgent:
		aApiPort = apiPortDefault
	case AppTypeServer:
		aApiPort = aApiPort + 1
	}

	fileCfg := &FileConfig{
		AppType: CurrentApp.AppType,
		Api: api_config.ApiConfig{
			Enabled: true,
			Port:    aApiPort,
		},
		BaseUploadDir: "./uploads",
	}
	return fileCfg
}

func SyncConfigFile(configFile string, endFunc func(error)) {
	if CurrentApp == nil {
		endFunc(errors.New("App Not Setup "))
		return
	}

	gtbox_log.LogInfof("加载配置文件 [%s]", CurrentApp.AppConfigFilePath)
	_, err := os.Stat(CurrentApp.AppConfigFilePath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		// 获取配置文件的父目录路径
		dir := filepath.Dir(CurrentApp.AppConfigFilePath)

		// 检查父目录是否存在
		if _, err = os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			// 创建父目录
			if err = os.MkdirAll(dir, 0755); err != nil {
				gtbox_log.LogErrorf("无法创建目录 [%s]: %s", dir, err.Error())
				return
			}
		}

		// 写入默认配置文件内容
		jd, _ := json.MarshalIndent(generateDefaultConfig(), "", "  ")
		err = os.WriteFile(CurrentApp.AppConfigFilePath, jd, 0755)
		if err != nil {
			gtbox_log.LogErrorf("无法写入配置文件 [%s]: %s", CurrentApp.AppConfigFilePath, err.Error())
			return
		}
	} else {
		buf, err := os.ReadFile(CurrentApp.AppConfigFilePath)
		if err != nil {
			gtbox_log.LogErrorf("读取配置文件 [%s] 错误: %s", CurrentApp.AppConfigFilePath, err.Error())

			return
		}
		if len(buf) == 0 {
			gtbox_log.LogErrorf("配置文件重置")
			jd, _ := json.MarshalIndent(generateDefaultConfig(), "", "  ")
			// 写入默认配置文件内容
			err = os.WriteFile(CurrentApp.AppConfigFilePath, jd, 0755)
			if err != nil {
				gtbox_log.LogErrorf("无法写入配置文件 [%s]: %s", CurrentApp.AppConfigFilePath, err.Error())
				return
			}
		}
	}

	err = LoadConfig(CurrentApp.AppConfigFilePath)

	if err != nil {
		gtbox_log.LogErrorf("无法加载配置文件 [%s]: %s", CurrentApp.AppConfigFilePath, err.Error())
		return
	}

}
