package pre_build_cfg

import (
	"github.com/george012/gtbox"
	"github.com/george012/gtbox/gtbox_app"
	"github.com/wmyeah/yeah_box/config"
)

var (
	CurrentApp *ExtendApp
)

type ExtendApp struct {
	*gtbox_app.App
	NetListenPortStratumDefault int
	ApiPort                     int
	CurrentRunWith              YeahBoxRunWith
}

func NewApp(appName, bundleID, description string, runMode gtbox.RunMode) *ExtendApp {
	app := &ExtendApp{
		App:                         gtbox_app.NewApp(appName, config.ProjectVersion, bundleID, description, runMode),
		NetListenPortStratumDefault: 0,
		ApiPort:                     apiPortDefault,
	}

	return app
}
