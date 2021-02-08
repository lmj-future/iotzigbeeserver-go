package models

import (
	"reflect"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/mongo"
)

func getTerminalTimerInfoModel() mongo.DBClient {
	cnames, _ := mongo.MongoClient.Database.CollectionNames()
	for _, name := range cnames {
		if name == "terminalTimerInfo" {
			var c = mongo.MongoClient.Database.C("terminalTimerInfo")
			var db mongo.Collection
			db.Collection = c
			return db
		}
	}

	devEUIIndex := mgo.Index{
		Key:    []string{"devEUI"},
		Unique: true,
	}
	mongo.MongoClient.Database.C("terminalTimerInfo").EnsureIndex(devEUIIndex)
	var c = mongo.MongoClient.Database.C("terminalTimerInfo")
	var db mongo.Collection
	db.Collection = c
	return db
}

//GetTerminalTimerInfoByDevEUI GetTerminalTimerInfoByDevEUI
func GetTerminalTimerInfoByDevEUI(devEUI string) (*config.TerminalTimerInfo, error) {
	var db = getTerminalTimerInfoModel()
	var filit = bson.M{"devEUI": devEUI}
	var terminalTimer config.TerminalTimerInfo
	err := db.FindOne(filit, &terminalTimer)
	if reflect.DeepEqual(terminalTimer, config.TerminalTimerInfo{}) {
		return nil, err
	}
	return &terminalTimer, err
}

//FindTerminalTimerInfoAndUpdate FindTerminalTimerInfoAndUpdate
func FindTerminalTimerInfoAndUpdate(devEUI string, setData config.TerminalTimerInfo) (*config.TerminalTimerInfo, error) {
	var db = getTerminalTimerInfoModel()
	var terminalTimer config.TerminalTimerInfo
	err := db.FindOneAndUpdate(bson.M{"devEUI": devEUI}, setData, &terminalTimer)
	if reflect.DeepEqual(terminalTimer, config.TerminalTimerInfo{}) {
		return nil, err
	}
	return &terminalTimer, err
}

//CreateTerminalTimerInfo CreateTerminalTimerInfo
func CreateTerminalTimerInfo(setData config.TerminalTimerInfo) error {
	var db = getTerminalTimerInfoModel()
	setData.CreateTime = time.Now()
	err := db.Create(setData)
	return err
}

//DeleteTerminalTimer DeleteTerminalTimer
func DeleteTerminalTimer(devEUI string) error {
	var db = getTerminalTimerInfoModel()
	err := db.Remove(bson.M{"devEUI": devEUI})
	return err
}
