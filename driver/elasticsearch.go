package driver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// ErrorResponse is the JSON response for non-201 insertion of new documents in
// Elasticsearch
type ErrorResponse struct {
	Error struct {
		CausedBy struct {
			Reason string `json:"reason"`
		} `json:"caused_by"`
	} `json:"error"`
}

// ElasticsearchDriver implements the Datastore interface for inserting
// documents into Elasticsearch
type ElasticsearchDriver struct {
	client *http.Client
	host   string
}

// NewElasticsearchDriver returns a new ElasticsearchDriver
func NewElasticsearchDriver() *ElasticsearchDriver {
	return &ElasticsearchDriver{
		client: &http.Client{
			Timeout: time.Second,
		},
	}
}

// Connect constructs the host that requests will be sent to
func (e *ElasticsearchDriver) Connect(credentials DatastoreCredentials) (err error) {
	if e.client == nil {
		return errors.New("Driver already closed")
	}

	e.host = fmt.Sprintf("%s:%s", credentials.Hostname, credentials.Port)
	return
}

// Publish sends a POST request to Elasticsearch
//
// The name of the index that data will be inserted to will be
// <appName>-<YYYY-MM-DD>
func (e *ElasticsearchDriver) Publish(data Entry, appName string) (err error) {
	if e.client == nil {
		return errors.New("Driver already closed")
	}

	index := fmt.Sprintf("%s-%s", appName, time.Now().Format("2006-01-02"))
	raw, err := json.Marshal(data)
	if err != nil {
		return
	}

	resp, err := e.client.Post(fmt.Sprintf("%s/%s/%s", e.host, index, "metrics"), "application/json", bytes.NewBuffer(raw))
	if err != nil {
		return
	}

	if resp.StatusCode != 201 {
		defer resp.Body.Close()
		raw, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New("POST to Elasticsearch unsuccessful")
		}

		var errResp ErrorResponse
		err = json.Unmarshal(raw, &errResp)
		if err != nil {
			return errors.New("POST to Elasticsearch unsuccessful")
		}
		return fmt.Errorf("POST to Elasticsearch unsuccessful: %s", errResp.Error.CausedBy.Reason)
	}

	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()

	return
}

// Close will close the driver by deleting its http.Client
func (e *ElasticsearchDriver) Close() (err error) {
	e.client = nil
	e.host = ""
	return
}
