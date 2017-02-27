package driver

import (
	"gopkg.in/mgo.v2"
)

// MongoDriver implements the Datastore interface for inserting documents into
// MongoDB
type MongoDriver struct {
	session      *mgo.Session
	advancedDial func(*mgo.DialInfo) (*mgo.Session, error)
	basicDial    func(string) (*mgo.Session, error)
	insert       func(Entry, string) error
}

// NewMongoDriver returns a new MongoDriver
func NewMongoDriver() *MongoDriver {
	m := &MongoDriver{
		advancedDial: mgo.DialWithInfo,
		basicDial:    mgo.Dial,
	}
	m.insert = func(data Entry, appName string) error {
		return m.session.DB(appName).C("profile").Insert(data)
	}
	return m
}

// Connect creates a session to a MongoDB instance
func (m *MongoDriver) Connect(credentials DatastoreCredentials) (err error) {
	if credentials.Username != "" {
		m.session, err = m.advancedDial(&mgo.DialInfo{
			Username: credentials.Username,
			Password: credentials.Password,
			Addrs:    []string{credentials.Hostname + ":" + credentials.Port},
		})
		return err
	}

	m.session, err = m.basicDial(credentials.Hostname + ":" + credentials.Port)
	return err
}

// Publish inserts data into a MongoDB collection
func (m *MongoDriver) Publish(data Entry, appName string) error {
	return m.insert(data, appName)
}

// Close closes the MongoDB session
func (m *MongoDriver) Close() error {
	m.session.Close()
	return nil
}
