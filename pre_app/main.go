package main

import (
	"fmt"
	"github.com/george012/gtbox"
	"github.com/george012/gtbox/gtbox_cmd"
	"github.com/george012/gtbox/gtbox_log"
	"github.com/wmyeah/yeah_box/config"
	"pre_app/api"
	"pre_app/common"
	"pre_app/custom_cmd"
	"pre_app/pre_app_cfg"
	"time"
)

var (
	mRunMode       = ""
	mGitCommitHash = ""
	mGitCommitTime = ""
	mPackageOS     = ""
	mPackageTime   = ""
	mGoVersion     = ""
	mRunWith       = ""
)

// setupApp 设置默认App
// projectName App名
// bundleID 包名
// description 描述
// acceptArgs 自定义 args []string 格式
// acceptArgsCallBack 自定义 args 放行的回调 自定义args 相关后续操作可以在这里执行
// Simole 1:
// common.SetupApp("testQueryCPU", "com.testing.testQueryCPU", "test Query CPU", nil, nil)
func setupApp() {
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

	if mRunWith != "agent" && mRunWith != "server" {
		fmt.Print("!!!---error---!!!: banary is not server or agent\n")
		fmt.Print("!!!---error---!!!: banary is not server or agent\n")
		fmt.Print("!!!---error---!!!: banary is not server or agent\n")
		common.ExitApp()
	}

	projectName := fmt.Sprintf("%s_%s", config.ProjectName, mRunWith)
	bundleID := fmt.Sprintf("%s_%s", config.ProjectBundleID, mRunWith)
	description := fmt.Sprintf("%s %s", config.ProjectName, mRunWith)

	pre_app_cfg.CurrentApp = pre_app_cfg.NewApp(projectName, bundleID, description, runMode)

	pre_app_cfg.CurrentApp.CurrentRunWith = pre_app_cfg.YeahBoxRunWith(mRunWith)
	pre_app_cfg.CurrentApp.AppConfigFilePath = fmt.Sprintf("./conf/config_%s.json", pre_app_cfg.CurrentApp.CurrentRunWith)

	//	TODO 初始化gtbox及log分片
	if pre_app_cfg.CurrentApp.CurrentRunMode == gtbox.RunModeDebug {
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

	pre_app_cfg.CurrentApp.GitCommitHash = mGitCommitHash
	pre_app_cfg.CurrentApp.GitCommitTime = mGitCommitTime
	pre_app_cfg.CurrentApp.GoVersion = mGoVersion
	pre_app_cfg.CurrentApp.PackageOS = mPackageOS
	pre_app_cfg.CurrentApp.PackageTime = mPackageTime

	custom_cmd.HandleCustomCmds(pre_app_cfg.CurrentApp)

	gtbox.SetupGTBox(pre_app_cfg.CurrentApp.AppName,
		pre_app_cfg.CurrentApp.CurrentRunMode,
		pre_app_cfg.CurrentApp.AppLogPath,
		30,
		gtbox_log.GTLogSaveHours,
		int(pre_app_cfg.CurrentApp.HTTPRequestTimeOut.Seconds()),
	)

	gtbox_log.LogDebugf("this is debug log test")

}

func main() {
	setupApp()

	pre_app_cfg.SyncConfigFile(func(err error) {
		if err != nil {
			gtbox_log.LogDebugf("%s", err.Error())
			common.ExitApp()
		}
	})

	api.StartAPIService(&pre_app_cfg.GlobalConfig.Api)

	common.LoadSigHandle(nil, nil)

}
