package api

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_log"
	"github.com/gorilla/mux"
	"github.com/wmyeah/yeah_box/api/api_config"
	"github.com/wmyeah/yeah_box/api/api_handler"
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

	apiCfg.UserAgentAllowed = append(apiCfg.UserAgentAllowed, fmt.Sprintf("%s/*", config.CurrentApp.AppName))

	apiCfg.APIMethodsAllowed = append(apiCfg.APIMethodsAllowed, apiMethods...)

	go func() {
		muxRouter := mux.NewRouter()
		muxRouter.Use(api_handler.Middleware) // 使用中间件
		muxRouter.HandleFunc("/", api_handler.HomeHandler).Methods("GET")
		muxRouter.HandleFunc("/api/v1", api_handler.ApiHandler).Methods("POST")

		addr := fmt.Sprintf("%s:%d", "0.0.0.0", apiCfg.Port)
		gtbox_log.LogInfof("API server Run On  [%s]", fmt.Sprintf("http://127.0.0.1:%d", apiCfg.Port))
		if err := http.ListenAndServe(addr, muxRouter); err != nil {
			gtbox_log.LogErrorf("Failed to start HTTP server: %v\n", err)
		}
	}()

}
