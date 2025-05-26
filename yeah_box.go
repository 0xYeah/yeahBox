package yeahBox

import (
	"fmt"
	"github.com/0xYeah/yeahBox/config"
)

func GetVersion() string {
	vInfos := fmt.Sprintf("%s/%s/%s", config.ProjectName, config.ProjectVersion, config.ProjectBundleID)
	return vInfos
}
