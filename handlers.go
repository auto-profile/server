package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func (e *Env) MetricHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("Type of request: %s\n", vars["type"])
	reader, err := gzip.NewReader(r.Body)
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(os.Stdout, reader); err != nil {
		panic(err)
	}
	return
}
