package main

import (
	"fmt"
	"net/http"

	stackimpact "github.com/auto-profile/stackimpact-go"
	"github.com/gorilla/mux"
)

func main() {
	agent := stackimpact.NewAgent()
	agent.Start(stackimpact.Options{
		DashboardAddress: "http://localhost:8080",
		AppName:          "test_server",
	})
	router := mux.NewRouter().StrictSlash(true)
	v1 := router.PathPrefix("/agent/v1").Subrouter()
	v1.HandleFunc("/{type}", MetricHandler)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}
