package main

import (
	"github.com/0xYeah/yeahBox/base_app"
	"github.com/0xYeah/yeahBox/base_app/app_cfg"
)

func main() {
	base_app.StartAppWithAppType(app_cfg.AppTypeAgent)
}
