package driver

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type MongoDriver struct {
	session      *mgo.Session
	advancedDial func(*mgo.DialInfo) (*mgo.Session, error)
	basicDial    func(string) (*mgo.Session, error)
	insert       func(Entry, string) error
}

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

func (m *MongoDriver) Publish(data Entry, appName string) error {
	return m.insert(data, appName)
}

func (m *MongoDriver) Get(app string, category string, resolution int) (entries []Entry, err error) {
	var entry Entry
	it := m.session.DB(app).C("profile").Find(bson.M{
		"payload.messages.content.category": category,
		"appname":                           app,
		"runts":                             bson.M{"$gte": time.Now().Unix() - int64(resolution)}}).Iter()

	for it.Next(&entry) {
		var messages []Message
		for _, msg := range entry.Payload.Messages {
			if msg.Content.Category == category {
				messages = append(messages, msg)
			}
		}
		entry.Payload.Messages = messages
		entries = append(entries, entry)
	}

	err = it.Close()
	return entries, err
}

func (m *MongoDriver) Close() error {
	m.session.Close()
	return nil
}
