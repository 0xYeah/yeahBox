package main

import (
	"github.com/wmyeah/yeah_box/base_app"
	"github.com/wmyeah/yeah_box/base_app/app_cfg"
)

func main() {
	base_app.StartAppWithAppType(app_cfg.AppTypeAgent)
}
