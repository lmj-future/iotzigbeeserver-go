package models

import (
	"reflect"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/mongo"
)

func getSocketInfoModel() mongo.DBClient {
	cnames, _ := mongo.MongoClient.Database.CollectionNames()
	for _, name := range cnames {
		if name == "socketInfo" {
			var c = mongo.MongoClient.Database.C("socketInfo")
			var db mongo.Collection
			db.Collection = c
			return db
		}
	}

	APMacIndex := mgo.Index{
		Key:    []string{"APMac"},
		Unique: true,
	}
	mongo.MongoClient.Database.C("socketInfo").EnsureIndex(APMacIndex)
	var c = mongo.MongoClient.Database.C("socketInfo")
	var db mongo.Collection
	db.Collection = c

	return db
}

//GetSocketInfoByAPMac GetSocketInfoByAPMac
func GetSocketInfoByAPMac(APMac string) (*config.SocketInfo, error) {
	var db = getSocketInfoModel()
	var filit = bson.M{"APMac": APMac}
	var socketInfo config.SocketInfo
	err := db.FindOne(filit, &socketInfo)
	if reflect.DeepEqual(socketInfo, config.SocketInfo{}) {
		return nil, err
	}
	return &socketInfo, err
}

//FindSocketInfoAndUpdate FindSocketInfoAndUpdate
func FindSocketInfoAndUpdate(APMac string, setData config.SocketInfo) (*config.SocketInfo, error) {
	var db = getSocketInfoModel() //从mongodb查找拓扑信息
	var socketInfo config.SocketInfo
	err := db.FindOneAndUpdate(bson.M{"APMac": APMac}, setData, &socketInfo) //更新数据库model信息
	if reflect.DeepEqual(socketInfo, config.SocketInfo{}) {
		return nil, err
	}
	return &socketInfo, err
}

//CreateSocketInfo CreateSocketInfo
func CreateSocketInfo(setData config.SocketInfo) error {
	var db = getSocketInfoModel()
	setData.CreateTime = time.Now()
	err := db.Create(setData)
	return err
}

//DeleteSocket DeleteSocket
func DeleteSocket(APMac string) error {
	var db = getSocketInfoModel()
	err := db.Remove(bson.M{"APMac": APMac})
	return err
}
