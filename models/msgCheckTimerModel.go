package models

import (
	"reflect"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/mongo"
)

func getMsgCheckTimerInfoModel() mongo.DBClient {
	cnames, _ := mongo.MongoClient.Database.CollectionNames()
	for _, name := range cnames {
		if name == "msgCheckTimerInfo" {
			var c = mongo.MongoClient.Database.C("msgCheckTimerInfo")
			var db mongo.Collection
			db.Collection = c
			return db
		}
	}

	msgKeyIndex := mgo.Index{
		Key:    []string{"msgKey"},
		Unique: true,
	}
	mongo.MongoClient.Database.C("msgCheckTimerInfo").EnsureIndex(msgKeyIndex)
	var c = mongo.MongoClient.Database.C("msgCheckTimerInfo")
	var db mongo.Collection
	db.Collection = c
	return db
}

//GetMsgCheckTimerInfoByMsgKey GetMsgCheckTimerInfoByMsgKey
func GetMsgCheckTimerInfoByMsgKey(msgKey string) (*config.MsgCheckTimerInfo, error) {
	var db = getMsgCheckTimerInfoModel()
	var filit = bson.M{"msgKey": msgKey}
	var msgCheckTimerInfo config.MsgCheckTimerInfo
	err := db.FindOne(filit, &msgCheckTimerInfo)
	if reflect.DeepEqual(msgCheckTimerInfo, config.MsgCheckTimerInfo{}) {
		return nil, err
	}
	return &msgCheckTimerInfo, err
}

//FindMsgCheckTimerInfoAndUpdate FindMsgCheckTimerInfoAndUpdate
func FindMsgCheckTimerInfoAndUpdate(msgKey string, setData config.MsgCheckTimerInfo) (*config.MsgCheckTimerInfo, error) {
	var db = getMsgCheckTimerInfoModel()
	var msgCheckTimerInfo config.MsgCheckTimerInfo
	err := db.FindOneAndUpdate(bson.M{"msgKey": msgKey}, setData, &msgCheckTimerInfo)
	if reflect.DeepEqual(msgCheckTimerInfo, config.MsgCheckTimerInfo{}) {
		return nil, err
	}
	return &msgCheckTimerInfo, err
}

//CreateMsgCheckTimerInfo CreateMsgCheckTimerInfo
func CreateMsgCheckTimerInfo(setData config.MsgCheckTimerInfo) error {
	var db = getMsgCheckTimerInfoModel()
	setData.CreateTime = time.Now()
	err := db.Create(setData)
	return err
}

//DeleteMsgCheckTimer DeleteMsgCheckTimer
func DeleteMsgCheckTimer(msgKey string) error {
	var db = getMsgCheckTimerInfoModel()
	err := db.Remove(bson.M{"msgKey": msgKey})
	return err
}
