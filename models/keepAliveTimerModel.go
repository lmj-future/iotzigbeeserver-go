package models

import (
	"reflect"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/mongo"
)

func getKeepAliveTimerInfoModel() mongo.DBClient {
	cnames, _ := mongo.MongoClient.Database.CollectionNames()
	for _, name := range cnames {
		if name == "keepAliveTimerInfo" {
			var c = mongo.MongoClient.Database.C("keepAliveTimerInfo")
			var db mongo.Collection
			db.Collection = c
			return db
		}
	}

	APMacIndex := mgo.Index{
		Key:    []string{"APMac"},
		Unique: true,
	}
	mongo.MongoClient.Database.C("keepAliveTimerInfo").EnsureIndex(APMacIndex)
	var c = mongo.MongoClient.Database.C("keepAliveTimerInfo")
	var db mongo.Collection
	db.Collection = c
	return db
}

//GetKeepAliveTimerInfoByAPMac GetKeepAliveTimerInfoByAPMac
func GetKeepAliveTimerInfoByAPMac(APMac string) (*config.KeepAliveTimerInfo, error) {
	var db = getKeepAliveTimerInfoModel()
	var filit = bson.M{"APMac": APMac}
	var keepAliveTimerInfo config.KeepAliveTimerInfo
	err := db.FindOne(filit, &keepAliveTimerInfo)
	if reflect.DeepEqual(keepAliveTimerInfo, config.KeepAliveTimerInfo{}) {
		return nil, err
	}
	return &keepAliveTimerInfo, err
}

//FindKeepAliveTimerInfoAndUpdate FindKeepAliveTimerInfoAndUpdate
func FindKeepAliveTimerInfoAndUpdate(APMac string, setData config.KeepAliveTimerInfo) (*config.KeepAliveTimerInfo, error) {
	var db = getKeepAliveTimerInfoModel()
	var keepAliveTimerInfo config.KeepAliveTimerInfo
	err := db.FindOneAndUpdate(bson.M{"APMac": APMac}, setData, &keepAliveTimerInfo)
	if reflect.DeepEqual(keepAliveTimerInfo, config.KeepAliveTimerInfo{}) {
		return nil, err
	}
	return &keepAliveTimerInfo, err
}

//CreateKeepAliveTimerInfo CreateKeepAliveTimerInfo
func CreateKeepAliveTimerInfo(setData config.KeepAliveTimerInfo) error {
	var db = getKeepAliveTimerInfoModel()
	setData.CreateTime = time.Now()
	err := db.Create(setData)
	return err
}

//DeleteKeepAliveTimer DeleteKeepAliveTimer
func DeleteKeepAliveTimer(APMac string) error {
	var db = getKeepAliveTimerInfoModel()
	err := db.Remove(bson.M{"APMac": APMac})
	return err
}
