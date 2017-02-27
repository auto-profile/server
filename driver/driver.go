/*
Package driver implements a simple interface for storing captured metric data
in an external data store
*/
package driver

// DatastoreCredentials represents all of the details of the external data
// store that metrics are inserted into
type DatastoreCredentials struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Driver   string `json:"driver"`
}

// Datastore manages interaction with the external data store
type Datastore interface {
	Connect(DatastoreCredentials) error
	Publish(Entry, string) error
	Close() error
}

// Entry represents all of the data forwarded by the agent
type Entry struct {
	RuntimeVersion string  `json:"runtime_version"`
	AppEnvironment string  `json:"app_environment"`
	Payload        Payload `json:"payload"`
	RunTs          int     `json:"run_ts"`
	RuntimeType    string  `json:"runtime_type"`
	HostName       string  `json:"host_name"`
	RunID          string  `json:"run_id"`
	SentAt         int     `json:"sent_at"`
	AgentVersion   string  `json:"agent_version"`
	AppName        string  `json:"app_name"`
	AppVersion     string  `json:"app_version"`
	BuildID        string  `json:"build_id"`
}

// Payload represents a list of messages
type Payload struct {
	Messages []Message `json:"messages"`
}

// Message represents the content and topic
type Message struct {
	Content Content `json:"content"`
	Topic   string  `json:"topic"`
}

// Content represents data about a metric
type Content struct {
	Category    string      `json:"category"`
	ID          string      `json:"id"`
	Measurement Measurement `json:"measurement"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Unit        string      `json:"unit"`
}

// Measurement represents the actual value of the metric and other metadata
type Measurement struct {
	Value     int       `json:"value"`
	Breakdown Breakdown `json:"breakdown"`
	ID        string    `json:"id"`
	Timestamp int       `json:"timestamp"`
	Trigger   string    `json:"trigger"`
}

// Breakdown represents an optional further breakdown of the metric into by the
// call tree
type Breakdown struct {
	NumSamples  int         `json:"num_samples"`
	Measurement float64     `json:"measurement"`
	Name        string      `json:"name"`
	Children    []Breakdown `json:"children"`
}
