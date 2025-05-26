package api

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_log"
	"github.com/george012/gtbox/gtbox_net"
	"github.com/gorilla/mux"
	"github.com/wmyeah/yeah_box/base_app/api/api_config"
	"github.com/wmyeah/yeah_box/base_app/api/api_handler"
	"github.com/wmyeah/yeah_box/base_app/app_cfg"
	"github.com/wmyeah/yeah_box/config"
	"net/http"
)

var apiMethods = []string{"auth", "logout"}

func StartAPIService(apiCfg *api_config.ApiConfig) {

	if apiCfg.Port < 1 || apiCfg.Port > 65535 {
		gtbox_log.LogErrorf("api port must be between 1 and 65535")
		return
	}

	api_config.CurrentApiConfig = apiCfg

	apiCfg.UserAgentAllowed = append(apiCfg.UserAgentAllowed, fmt.Sprintf("%s/*", app_cfg.CurrentApp.AppName))

	apiCfg.APIMethodsAllowed = append(apiCfg.APIMethodsAllowed, apiMethods...)

	go func() {
		apiGroup := "/api/v1"

		muxRouter := mux.NewRouter()
		muxRouter.Use(api_handler.Middleware) // 使用中间件
		muxRouter.HandleFunc("/", api_handler.HomeHandler).Methods("GET")
		muxRouter.HandleFunc(apiGroup, api_handler.ApiHandler).Methods("POST")

		runWith := pre_app_cfg.CurrentApp.CurrentRunWith

		switch runWith {
		case config.YeahBoxRunWithServer:
			loadApiMethodForServer(apiGroup, api_handler.ApiHandler)
		case config.YeahBoxRunWithAgent:
			loadApiMethodForAgent(apiGroup, api_handler.ApiHandler)
		}

		addr := fmt.Sprintf("%s:%d", "0.0.0.0", apiCfg.Port)
		localAddr := gtbox_net.GTGetLocalIPV4WithCurrentActive()
		pubAddr := gtbox_net.GTGetPublicIPV4()
		gtbox_log.LogInfof("API server Run On  [%s]", fmt.Sprintf("http://127.0.0.1:%d", apiCfg.Port))
		gtbox_log.LogInfof("API server Run as local internet [%s]", fmt.Sprintf("http://%s:%d", localAddr, apiCfg.Port))
		gtbox_log.LogInfof("API server Run as public internet [%s]", fmt.Sprintf("http://%s:%d", pubAddr, apiCfg.Port))

		if err := http.ListenAndServe(addr, muxRouter); err != nil {
			gtbox_log.LogErrorf("Failed to start HTTP server: %v\n", err)
		}
	}()

}
