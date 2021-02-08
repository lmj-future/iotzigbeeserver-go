package models

import (
	"reflect"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/mongo"
)

func getTmnNwkInfoModel() mongo.DBClient {
	cnames, _ := mongo.MongoClient.Database.CollectionNames()
	for _, name := range cnames {
		if name == "terminalNetworkInfo" {
			var c = mongo.MongoClient.Database.C("terminalNetworkInfo")
			var db mongo.Collection
			db.Collection = c
			return db
		}
	}

	APMacmoduleIDIndex := mgo.Index{
		Key:    []string{"APMac", "moduleID"},
		Unique: true,
	}
	mongo.MongoClient.Database.C("terminalNetworkInfo").EnsureIndex(APMacmoduleIDIndex)
	mongo.MongoClient.Database.C("terminalNetworkInfo").EnsureIndexKey("APMac")
	var c = mongo.MongoClient.Database.C("terminalNetworkInfo")
	var db mongo.Collection
	db.Collection = c
	return db
}

//GetTmnNwkInfo GetTmnNwkInfo
func GetTmnNwkInfo(query bson.M) (*config.TmnNwkInfo, error) {
	var db = getTmnNwkInfoModel()
	var tmnNwkInfo config.TmnNwkInfo
	err := db.FindOne(query, &tmnNwkInfo)
	if reflect.DeepEqual(tmnNwkInfo, config.TmnNwkInfo{}) {
		return nil, err
	}
	return &tmnNwkInfo, err
}

//FindTmnNwkAndUpdate FindTmnNwkAndUpdate
func FindTmnNwkAndUpdate(query bson.M, setData config.TmnNwkInfo) (*config.TmnNwkInfo, error) {
	var db = getTmnNwkInfoModel()
	var tmnNwkInfo config.TmnNwkInfo
	err := db.FindOneAndUpdate(query, setData, &tmnNwkInfo)
	if reflect.DeepEqual(tmnNwkInfo, config.TmnNwkInfo{}) {
		return nil, err
	}
	return &tmnNwkInfo, err
}

//CreateTmnNwk CreateTmnNwk
func CreateTmnNwk(setData config.TmnNwkInfo) error {
	var db = getTmnNwkInfoModel()
	setData.CreateTime = time.Now()
	err := db.Create(setData)
	return err
}

//DeleteTmnNwk DeleteTmnNwk
func DeleteTmnNwk(APMac string) error {
	var db = getTmnNwkInfoModel()
	err := db.Remove(bson.M{"APMac": APMac})
	return err
}

//DeleteTmnNwkByAPMacAndModuleID DeleteTmnNwkByAPMacAndModuleID
func DeleteTmnNwkByAPMacAndModuleID(APMac string, moduleID string) error {
	var db = getTmnNwkInfoModel()
	err := db.Remove(bson.M{"APMac": APMac, "moduleID": moduleID})
	return err
}
