package driver

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type MongoDriver struct {
	session      Session
	advancedDial func(*mgo.DialInfo) (*mgo.Session, error)
	basicDial    func(string) (*mgo.Session, error)
}

type MongoSession struct {
	session *mgo.Session
}

func (m *MongoSession) DB(database string) Database {
	return &MongoDatabase{m.session.DB(database)}
}

func (m *MongoSession) Close() {
	m.session.Close()
}

type Session interface {
	DB(string) Database
	Close()
}

type MongoDatabase struct {
	Database *mgo.Database
}

func (m *MongoDatabase) C(collection string) Collection {
	return &MongoCollection{m.Database.C(collection)}
}

type Database interface {
	C(string) Collection
}

type MongoCollection struct {
	Collection *mgo.Collection
}

func (m *MongoCollection) Find(doc interface{}) Query {
	return &MongoQuery{m.Collection.Find(doc)}
}

func (m *MongoCollection) Insert(docs ...interface{}) error {
	return m.Collection.Insert(docs...)
}

type Collection interface {
	Find(interface{}) Query
	Insert(docs ...interface{}) error
}

type Query interface {
	Iter() Iterator
}

type MongoQuery struct {
	Query *mgo.Query
}

func (m *MongoQuery) Iter() Iterator {
	return &MongoIterator{m.Query.Iter()}
}

type MongoIterator struct {
	Iterator *mgo.Iter
}

func (m *MongoIterator) Next(doc interface{}) bool {
	return m.Iterator.Next(doc)
}

func (m *MongoIterator) Close() error {
	return m.Iterator.Close()
}

type Iterator interface {
	Next(interface{}) bool
	Close() error
}

func NewMongoDriver() *MongoDriver {
	return &MongoDriver{
		advancedDial: mgo.DialWithInfo,
		basicDial:    mgo.Dial,
	}
}

func (m *MongoDriver) Connect(credentials DatastoreCredentials) (err error) {
	if credentials.Username != "" {
		session, err := m.advancedDial(&mgo.DialInfo{
			Username: credentials.Username,
			Password: credentials.Password,
			Addrs:    []string{credentials.Hostname + ":" + credentials.Port},
		})
		m.session = &MongoSession{session}
		return err
	}

	session, err := m.basicDial(credentials.Hostname + ":" + credentials.Port)
	m.session = &MongoSession{session}
	return err
}

func (m *MongoDriver) Publish(data Entry, appName string) error {
	return m.session.DB(appName).C("profile").Insert(data)
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
