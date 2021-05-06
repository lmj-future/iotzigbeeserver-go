package mongo

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
)

// MongoClient MongoClient
var MongoClient *Client // Singleton used so that mongoEvent can use it to de-reference readings

// Client Client
type Client struct {
	Session  *mgo.Session  // Mongo database session
	Database *mgo.Database // Mongo database
}

//Collection Collection
type Collection struct {
	*mgo.Collection
}

// NewClient Return a pointer to the MongoClient
func NewClient(url string) (*Client, error) {
	var m Client
	mongoInfo, _ := mgo.ParseURL(url)
	mongoInfo.Timeout = 5 * time.Second
	globallogger.Log.Infof("%+v", mongoInfo)
	session, err := mgo.DialWithInfo(mongoInfo)
	if err != nil {
		return nil, err
	}
	err = session.Ping()
	if err != nil {
		globallogger.Log.Errorln(err)
		return nil, err
	}
	m.Session = session
	m.Database = session.DB(mongoInfo.Database)
	return &m, nil
}

//CloseSession CloseSession
func (mc Client) CloseSession() {
	mc.Session.Close()
}
