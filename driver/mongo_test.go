package driver

import (
	"gopkg.in/mgo.v2"
	"testing"
)

type MockSession struct {}
type MockDatabase struct {}
type MockCollection struct {}
type MockQuery struct {}
type MockIterator struct {
	Continue bool
}

func (m *MockSession) DB(database string) Database {
	return &MockDatabase{}
}

func (m *MockSession) Close() {}

func (m *MockDatabase) C(collection string) Collection {
	return &MockCollection{}
}

func (m *MockCollection) Find(doc interface{}) Query {
	return &MockQuery{}
}

func (m *MockCollection) Insert(docs ...interface{}) error {
	return nil
}

func (m *MockQuery) Iter() Iterator {
	return &MockIterator{true}
}

func (m *MockIterator) Next(doc interface{}) bool {
	if m.Continue {
		m.Continue = false
		e := doc.(*Entry)
		message := Message{}
		message.Content = Content{}
		message.Content.Category = "memory"
		e.Payload.Messages = append(e.Payload.Messages, message)
		return true
	}
	return false
}

func (m *MockIterator) Close() error {
	return nil
}

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
		session: &MockSession{},
	}
	err := m.Publish(Entry{}, "test_app")

	if err != nil {
		t.Errorf("Shouldn't have received error with nil insert func, got %s", err)
	}
}

func TestMongoDriver_Get(t *testing.T) {
	m := &MongoDriver{
		session: &MockSession{},
	}
	_, err := m.Get("test_app", "memory", 60)

	if err != nil {
		t.Errorf("Shouldn't have received error with nil insert func, got %s", err)
	}
}

func TestMongoDriver_Close(t *testing.T) {
	m := &MongoDriver{
		session: &MockSession{},
	}
	m.Close()
}

func TestMongoProxyMethods(t *testing.T) {
	m := &MongoSession{&mgo.Session{}}
	db := m.DB("test")

	if db == nil {
		t.Error("Should have received DB when calling DB on MongoSession, got nil")
	}

	collection := db.C("test")

	if collection == nil {
		t.Error("Should have received collection when calling C on MongoDatabase, got nil")
	}

	query := collection.Find(nil)

	if query == nil {
		t.Error("Should have received Query when calling Find on MongoCollection, got nil")
	}

	m.Close()
}