package main

import (
	"fmt"
	"net/http"

	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	"github.com/auto-profile/server/driver"
	stackimpact "github.com/auto-profile/stackimpact-go"
	"github.com/gorilla/mux"
)

type Env struct {
	dataStore driver.Datastore
}

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "configPath", "config.json", "Path to configuration file")
	flag.Parse()
}

func main() {
	config, err := os.Open(configPath)

	if err != nil {
		panic(fmt.Sprintf("Could not open configuration file: %s", err))
	}

	raw, err := ioutil.ReadAll(config)

	if err != nil {
		panic(fmt.Sprintf("Could not read configuration file data: %s", err))
	}

	var credentials driver.DatastoreCredentials
	err = json.Unmarshal(raw, &credentials)

	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal configuration data: %s", err))
	}

	env := Env{}

	switch credentials.Driver {
	case "mongo":
		env.dataStore = driver.NewMongoDriver()
	case "elasticsearch":
		env.dataStore = driver.NewElasticsearchDriver()
	default:
		panic(fmt.Sprintf("Use of unsupported driver: %s", credentials.Driver))
	}

	err = env.dataStore.Connect(credentials)

	if err != nil {
		panic(err)
	}

	agent := stackimpact.NewAgent()
	agent.Start(stackimpact.Options{
		DashboardAddress: "http://localhost:8080",
		AppName:          "test_server",
	})

	router := mux.NewRouter().StrictSlash(true)
	v1 := router.PathPrefix("/agent/v1").Subrouter()
	v1.HandleFunc("/metrics", env.GetMetricsHandler)
	v1.HandleFunc("/{type}", env.MetricHandler)
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}
