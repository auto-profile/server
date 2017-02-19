package main

import (
	"compress/gzip"
	"fmt"
	"net/http"

	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
)

func (e *Env) MetricHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("Type of request: %s\n", vars["type"])
	reader, err := gzip.NewReader(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Bad application/gzip request"))
	}

	raw, err := ioutil.ReadAll(reader)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Bad application/gzip request"))
	}

	var req map[string]interface{}
	err = json.Unmarshal(raw, &req)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Could not unmarshal request"))
	}

	appName, ok := req["app_name"]

	if !ok {
		w.WriteHeader(400)
		w.Write([]byte("Request missing required key: app_name"))
	}

	app, ok := appName.(string)

	if !ok {
		w.WriteHeader(400)
		w.Write([]byte("app_name must be string"))
	}

	err = e.dataStore.Publish(raw, app)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Could not write to data store"))
	}

	return
}
