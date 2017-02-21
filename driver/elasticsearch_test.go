package driver

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestElasticsearchNewDriver(t *testing.T) {
	d := NewElasticsearchDriver()

	if d == nil {
		t.Error("Elasticsearch driver failed to instantiate (got nil)")
	}
}

func TestElasticsearchDriver_HostSet(t *testing.T) {
	d := NewElasticsearchDriver()

	err := d.Connect(DatastoreCredentials{
		Hostname: "http://localhost",
		Port:     "9200",
	})

	if err != nil {
		t.Errorf("Didn't expect an error")
	}

	expected := "http://localhost:9200"
	if d.host != expected {
		t.Errorf("Expected %s for host, got %s", expected, d.host)
	}
}

func TestElasticsearchDriver_AlreadyClosedError(t *testing.T) {
	d := NewElasticsearchDriver()

	err := d.Close()
	if err != nil {
		t.Errorf("Did not expect an error: %v", err)
	}

	err = d.Connect(DatastoreCredentials{})
	if err == nil {
		t.Errorf("Expected error to be thrown for previously closed driver")
	}
}

func TestElasticsearchDriver_PublishAlreadyClosedError(t *testing.T) {
	d := NewElasticsearchDriver()

	err := d.Close()
	if err != nil {
		t.Errorf("Did not expect an error: %v", err)
	}

	err = d.Publish(Entry{}, "test")
	if err == nil {
		t.Errorf("Expected error to be thrown for previously closed driver")
	}
}

func TestElasticsearchSuccessfulInsert(t *testing.T) {
	var method string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		w.WriteHeader(201)
		w.Write([]byte("OK"))
	}))

	d := NewElasticsearchDriver()
	d.host = ts.URL
	err := d.Publish(Entry{}, "test")
	if err != nil {
		t.Error("Didn't expect an error")
	}

	if method != "POST" {
		t.Errorf("Expected POST request but got %s instead", method)
	}
}

func TestElasticsearchIndexSetCorrectly(t *testing.T) {
	var index string
	appName := "test"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		index = r.URL.Path
		w.WriteHeader(201)
		w.Write([]byte("OK"))
	}))

	d := NewElasticsearchDriver()
	d.host = ts.URL
	err := d.Publish(Entry{}, appName)
	if err != nil {
		t.Error("Didn't expect an error")
	}
	expected := fmt.Sprintf("/%s-%s/metrics", appName, time.Now().Format("2006-01-02"))
	if index != expected {
		t.Errorf("Expected index to be %s, got %s", expected, index)
	}
}

func TestElasticsearchUnsuccessfulInsert(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
	}))

	d := NewElasticsearchDriver()
	d.host = ts.URL
	err := d.Publish(Entry{}, "test")
	if err == nil {
		t.Error("Expected an error")
	}
}

func TestElasticsearchUnsuccessfulInsertTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))

	d := NewElasticsearchDriver()
	d.host = ts.URL
	err := d.Publish(Entry{}, "test")
	if err == nil {
		t.Error("Expected an error")
	}
}

func TestElasticsearchUnsuccessfulInsertErrorReturned(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte(`{"error": {"caused_by": {"reason": "bad request"}}}`))
	}))

	d := NewElasticsearchDriver()
	d.host = ts.URL
	err := d.Publish(Entry{}, "test")
	if err == nil {
		t.Error("Expected an error")
	}

	expected := "POST to Elasticsearch unsuccessful: bad request"
	if err.Error() != expected {
		t.Errorf("Expected error to be '%s', got '%s' instead", err, expected)
	}
}
