package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRouter(r *mux.Router) {
	//添加日志记录器
	// r.Use(logMidware)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

}
