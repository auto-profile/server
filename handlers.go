package main

import (
	"compress/gzip"
	"fmt"
	"net/http"

	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"github.com/auto-profile/server/driver"
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

	fmt.Println(string(raw))

	var req driver.Entry
	err = json.Unmarshal(raw, &req)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Could not unmarshal request"))
	}

	err = e.dataStore.Publish(req, req.AppName)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Could not write to data store"))
	}

	return
}
