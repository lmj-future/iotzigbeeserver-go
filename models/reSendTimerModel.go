package models

import (
	"reflect"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/mongo"
)

func getReSendTimerInfoModel() mongo.DBClient {
	cnames, _ := mongo.MongoClient.Database.CollectionNames()
	for _, name := range cnames {
		if name == "reSendTimerInfo" {
			var c = mongo.MongoClient.Database.C("reSendTimerInfo")
			var db mongo.Collection
			db.Collection = c
			return db
		}
	}

	msgKeyIndex := mgo.Index{
		Key:    []string{"msgKey"},
		Unique: true,
	}
	mongo.MongoClient.Database.C("reSendTimerInfo").EnsureIndex(msgKeyIndex)
	var c = mongo.MongoClient.Database.C("reSendTimerInfo")
	var db mongo.Collection
	db.Collection = c

	return db
}

//GetReSendTimerInfoByMsgKey GetReSendTimerInfoByMsgKey
func GetReSendTimerInfoByMsgKey(msgKey string) (*config.ReSendTimerInfo, error) {
	var db = getReSendTimerInfoModel()
	var filit = bson.M{"msgKey": msgKey}
	var reSendTimerInfo config.ReSendTimerInfo
	err := db.FindOne(filit, &reSendTimerInfo)
	if reflect.DeepEqual(reSendTimerInfo, config.ReSendTimerInfo{}) {
		return nil, err
	}
	return &reSendTimerInfo, err
}

//FindReSendTimerInfoAndUpdate FindReSendTimerInfoAndUpdate
func FindReSendTimerInfoAndUpdate(msgKey string, setData config.ReSendTimerInfo) (*config.ReSendTimerInfo, error) {
	var db = getReSendTimerInfoModel()
	var reSendTimerInfo config.ReSendTimerInfo
	err := db.FindOneAndUpdate(bson.M{"msgKey": msgKey}, setData, &reSendTimerInfo)
	if reflect.DeepEqual(reSendTimerInfo, config.ReSendTimerInfo{}) {
		return nil, err
	}
	return &reSendTimerInfo, err
}

//CreateReSendTimerInfo CreateReSendTimerInfo
func CreateReSendTimerInfo(setData config.ReSendTimerInfo) error {
	var db = getReSendTimerInfoModel()
	setData.CreateTime = time.Now()
	err := db.Create(setData)
	return err
}

//DeleteReSendTimer DeleteReSendTimer
func DeleteReSendTimer(msgKey string) error {
	var db = getReSendTimerInfoModel()
	err := db.Remove(bson.M{"msgKey": msgKey})
	return err
}
