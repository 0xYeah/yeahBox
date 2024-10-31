package main

import (
	"fmt"
	"github.com/wmyeah/yeah_box/common"
	"github.com/wmyeah/yeah_box/config"
	"github.com/wmyeah/yeah_box/pre_builds"
)

func main() {

	pre_builds.SetupApp(
		fmt.Sprintf("%s_agent", config.ProjectName),
		fmt.Sprintf("%s_agent", config.ProjectBundleID),
		fmt.Sprintf("%s agent binary", config.ProjectName),
		nil,
	)

	common.LoadSigHandle(nil, nil)
}
