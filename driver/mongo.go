package driver

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDriver struct {
	session  *mgo.Session
	advancedDial func(*mgo.DialInfo) (*mgo.Session, error)
	basicDial func (string) (*mgo.Session, error)
	insert func([]byte, string) error
}

func NewMongoDriver() *MongoDriver {
	m := &MongoDriver{
		advancedDial: mgo.DialWithInfo,
		basicDial: mgo.Dial,
	}
	m.insert = func(data []byte, appName string) error {
		return m.session.DB(appName).C("profile").Insert(bson.M{"data": string(data)})
	}
	return m
}

func (m *MongoDriver) Connect(credentials DatastoreCredentials) (err error) {
	if credentials.Username != "" {
		m.session, err = m.advancedDial(&mgo.DialInfo{
			Username: credentials.Username,
			Password: credentials.Password,
			Addrs: []string{credentials.Hostname+":"+credentials.Port},
		})
		return err
	}

	m.session, err = m.basicDial(credentials.Hostname+":"+credentials.Port)
	return err
}

func (m *MongoDriver) Publish(data []byte, appName string) error {
	return m.insert(data, appName)
}

func (m *MongoDriver) Close() error {
	m.session.Close()
	return nil
}