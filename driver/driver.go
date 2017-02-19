package driver

type DatastoreCredentials struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Datastore interface {
	Connect(DatastoreCredentials) error
	Publish([]byte, string) error
	Close() error
}
