package base_app

import (
	"fmt"
	"github.com/0xYeah/yeahBox/base_app/api"
	"github.com/0xYeah/yeahBox/base_app/app_cfg"
	"github.com/0xYeah/yeahBox/base_app/common"
	"github.com/0xYeah/yeahBox/base_app/custom_cmd"
	"github.com/0xYeah/yeahBox/config"
	"github.com/george012/gtbox"
	"github.com/george012/gtbox/gtbox_cmd"
	"github.com/george012/gtbox/gtbox_log"
	"time"
)

var (
	mRunMode       = ""
	mGitCommitHash = ""
	mGitCommitTime = ""
	mPackageOS     = ""
	mPackageTime   = ""
	mGoVersion     = ""
	mAppType       = ""
)

// setupApp 设置默认App
// projectName App名
// bundleID 包名
// description 描述
// acceptArgs 自定义 args []string 格式
// acceptArgsCallBack 自定义 args 放行的回调 自定义args 相关后续操作可以在这里执行
// Simole 1:
// common.SetupApp("testQueryCPU", "com.testing.testQueryCPU", "test Query CPU", nil, nil)
func setupApp(appType app_cfg.AppType) {
	runMode := gtbox.RunModeDebug
	switch mRunMode {
	case "debug":
		runMode = gtbox.RunModeDebug
	case "test":
		runMode = gtbox.RunModeTest
	case "release":
		runMode = gtbox.RunModeRelease
	default:
		runMode = gtbox.RunModeDebug
	}

	projectName := fmt.Sprintf("%s_%s", config.ProjectName, appType)
	bundleID := fmt.Sprintf("%s_%s", config.ProjectBundleID, appType)
	description := fmt.Sprintf("%s %s", config.ProjectName, appType)

	app_cfg.CurrentApp = app_cfg.NewApp(projectName, config.ProjectVersion, bundleID, description, runMode)

	app_cfg.CurrentApp.AppConfigFilePath = fmt.Sprintf("./conf/config_%s.yaml", appType)

	//	TODO 初始化gtbox及log分片
	if app_cfg.CurrentApp.CurrentRunMode == gtbox.RunModeDebug {
		cmdMap := map[string]string{
			"git_commit_hash": "git show -s --format=%H",
			"git_commit_time": "git show -s --format=\"%ci\" | cut -d ' ' -f 1,2 | sed 's/ /_/'",
			"build_os":        "go env GOOS",
			"go_version":      "go version | awk '{print $3}'",
		}
		cmdRes := gtbox_cmd.RunWith(cmdMap)

		if cmdRes != nil {
			mGitCommitHash = cmdRes["git_commit_hash"]
			mGitCommitTime = cmdRes["git_commit_time"]
			mPackageOS = cmdRes["build_os"]
			mGoVersion = cmdRes["go_version"]
			mPackageTime = time.Now().UTC().Format("2006-01-02_15:04:05")
		}
	}

	app_cfg.CurrentApp.GitCommitHash = mGitCommitHash
	app_cfg.CurrentApp.GitCommitTime = mGitCommitTime
	app_cfg.CurrentApp.GoVersion = mGoVersion
	app_cfg.CurrentApp.PackageOS = mPackageOS
	app_cfg.CurrentApp.PackageTime = mPackageTime

	custom_cmd.HandleCustomCmds(app_cfg.CurrentApp)

	gtbox.SetupGTBox(app_cfg.CurrentApp.AppName,
		app_cfg.CurrentApp.CurrentRunMode,
		app_cfg.CurrentApp.AppLogPath,
		30,
		gtbox_log.GTLogSaveHours,
		app_cfg.CurrentApp.HTTPRequestTimeOut,
	)

	gtbox_log.LogDebugf("this is debug log test")

}

func StartAppWithAppType(appType app_cfg.AppType) {
	setupApp(appType)

	app_cfg.SyncConfigFile(app_cfg.CurrentApp.AppConfigFilePath, func(err error) {
		if err != nil {
			gtbox_log.LogDebugf("%s", err.Error())
			common.ExitApp()
		}
	})

	api.StartAPIService(&app_cfg.GlobalConfig.Api)

	common.LoadSigHandle(nil, nil)
}
