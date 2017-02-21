package main

import (
	"compress/gzip"
	"fmt"
	"net/http"

	"encoding/json"
	"github.com/auto-profile/server/driver"
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
		return
	}

	raw, err := ioutil.ReadAll(reader)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Bad application/gzip request"))
		return
	}

	fmt.Println(string(raw))

	var req driver.Entry
	err = json.Unmarshal(raw, &req)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Could not unmarshal request"))
		return
	}

	err = e.dataStore.Publish(req, req.AppName)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Could not write to data store"))
		return
	}

	return
}

type GetMetricsRequest struct {
	App        string `json:"app_name"`
	Resolution int    `json:"resolution"`
	Category   string `json:"category"`
}

func (e *Env) GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Bad get metrics request"))
		return
	}

	fmt.Println(string(raw))
	var req GetMetricsRequest
	err = json.Unmarshal(raw, &req)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Bad get metrics request"))
		return
	}

	results, err := e.dataStore.Get(req.App, req.Category, req.Resolution)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Could not find data"))
		return
	}

	raw, err = json.Marshal(results)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Could not marshal results"))
		return
	}

	w.WriteHeader(200)
	w.Write(raw)
}
