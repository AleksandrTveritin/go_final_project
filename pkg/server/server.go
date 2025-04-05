package server

import (
	"fmt"
	"net/http"

	"github.com/AleksandrTveritin/go_final_project/pkg/api"
	"github.com/AleksandrTveritin/go_final_project/pkg/config"
)

func Run() error {
	api.Init()

	http.Handle("/", http.FileServer(http.Dir("web")))

	return http.ListenAndServe(fmt.Sprintf(":%d", config.AppConfig.Port), nil)
}
