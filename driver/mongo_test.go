package driver

import (
	"testing"
	"gopkg.in/mgo.v2"
)

func TestNewMongoDriver(t *testing.T) {
	m := NewMongoDriver()

	if m == nil {
		t.Error("Mongo driver failed to instantiate (got nil)")
	}
}

func TestMongoDriver_ConnectWithCredentials(t *testing.T) {
	m := &MongoDriver{
		advancedDial: func(*mgo.DialInfo) (*mgo.Session, error) {
			return nil, nil
		},
	}
	err := m.Connect(DatastoreCredentials{
		Username: "Test",
		Password: "secret",
	})

	if err != nil {
		t.Errorf("Shouldn't have received error during mock initialization with credentials, got %s", err)
	}
}

func TestMongoDriver_ConnectWithoutCredentials(t *testing.T) {
	m := &MongoDriver{
		basicDial: func(string) (*mgo.Session, error) {
			return nil, nil
		},
	}
	err := m.Connect(DatastoreCredentials{})

	if err != nil {
		t.Errorf("Shouldn't have received error during mock initialization with credentials, got %s", err)
	}
}

func TestMongoDriver_Publish(t *testing.T) {
	m := &MongoDriver{
		insert: func(data []byte, app string) error {
			return nil
		},
	}
	err := m.Publish([]byte("test"), "test_app")

	if err != nil {
		t.Errorf("Shouldn't have received error with nil insert func, got %s", err)
	}
}