package main

import (
	"fmt"
	"net/http"

	"github.com/auto-profile/server/driver"
	stackimpact "github.com/auto-profile/stackimpact-go"
	"github.com/gorilla/mux"
)

type Env struct {
	dataStore driver.Datastore
}

func main() {
	agent := stackimpact.NewAgent()
	agent.Start(stackimpact.Options{
		DashboardAddress: "http://localhost:8080",
		AppName:          "test_server",
	})

	// TODO: Load config file here and create a Datastore

	env := Env{}

	router := mux.NewRouter().StrictSlash(true)
	v1 := router.PathPrefix("/agent/v1").Subrouter()
	v1.HandleFunc("/{type}", env.MetricHandler)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}
