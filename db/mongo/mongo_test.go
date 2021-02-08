package mongo

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
)

func TestClient(t *testing.T) {
	//[mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	//mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb
	//url := "mongodb://iotzigbeeserver:iotzigbeeserver@localhost:27017/iotzigbeeserver"
	url := "localhost:27017"
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Println("connect failed")
		os.Exit(1)
	}
	defer session.Close()
	dbnames, err := session.DatabaseNames()
	fmt.Println(dbnames)
	// err = session.DB("fortest").AddUser("mydb", "mydb", false)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	err = session.DB("mydb").C("cc").Create(&mgo.CollectionInfo{})
	if err != nil {
		fmt.Println("create collection cc err", err)
	}
	cnames, err := session.DB("mydb").CollectionNames()
	fmt.Println(cnames)
	// err = session.DB("mydb").C("bb").Insert(map[string]interface{}{"testbb": "bb"})
	// if err != nil {
	// 	fmt.Println("create collection bb err", err)
	// }
	err = session.DB("mydb").C("bb").Insert(map[string]interface{}{"testbb": "bhhh"})
	if err != nil {
		fmt.Println("insert many error", err)
	}
	//database := session.DB("fortest")
	//database.C("fortest").Insert("fortestvalue")
}

func TestNewClient(t *testing.T) {
	url := "localhost:27017"
	var dbclient DBClient
	var collection Collection
	client, err := NewClient(url)
	if err != nil {
		fmt.Println("connect mongodb err", err)
		os.Exit(1)
	}
	defer client.CloseSession()
	collection.Collection = client.Database.C("mytest")
	dbclient = collection

	err = dbclient.Create(map[string]interface{}{"test": "ccc"})
	if err != nil {
		fmt.Println("insert err", err)
	}
	err = dbclient.Create(map[string]interface{}{"test": "ddd"})
	if err != nil {
		fmt.Println("insert err", err)
	}
}

func TestFunction(t *testing.T) {
	url := "localhost:27017/iotzigbeeserver"
	var dbclient DBClient
	var collection Collection
	client, err := NewClient(url)
	if err != nil {
		fmt.Println("connect mongodb err", err)
		os.Exit(1)
	}
	defer client.CloseSession()

	// devEUIIndex := mgo.Index{
	// 	Key:    []string{"devEUI"},
	// 	Unique: true,
	// }
	// client.Database.C("mytest0317").EnsureIndex(devEUIIndex)
	// client.Database.C("mytest0317").EnsureIndexKey("OIDIndex")
	// client.Database.C("mytest0317").EnsureIndexKey("APMac", "devEUI")
	// client.Database.C("mytest0317").EnsureIndexKey("APMac", "moduleID")
	// client.Database.C("mytest0317").EnsureIndexKey("APMac")
	// client.Database.C("mytest0317").EnsureIndexKey("scenarioID")
	collection.Collection = client.Database.C("mytest0402")
	dbclient = collection
	setdata := config.TerminalInfo{
		//APMac:  "1111",
		//DevEUI: "2222",
		Online: false,
	}

	aa, _ := json.Marshal(setdata)
	bb, _ := bson.Marshal(setdata)
	fmt.Printf("%s\n", aa)
	fmt.Printf("%s\n", bb)
	var result config.TerminalInfo

	err = dbclient.FindOneAndUpdate(bson.M{"APMac": ""}, setdata, &result)
	if err != nil {
		fmt.Println("FindOneAndUpdate err", err)
	}

	// err = dbclient.Create(setdata)
	// if err != nil {
	// 	fmt.Println("insert err", err)
	// }
	// err = dbclient.Create(map[string]interface{}{"test": "ddd"})
	// if err != nil {
	// 	fmt.Println("insert err", err)
	// }
}
