package yeah_box

import (
	"fmt"
	"github.com/wmyeah/yeah_box/config"
)

func GetVersion() string {
	vInfos := fmt.Sprintf("%s/%s/%s", config.ProjectName, config.ProjectVersion, config.ProjectBundleID)
	return vInfos
}
