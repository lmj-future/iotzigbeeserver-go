package models

import (
	"reflect"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/mongo"
)

func getTmnNwkHistoryInfoModel() mongo.DBClient {
	cnames, _ := mongo.MongoClient.Database.CollectionNames()
	for _, name := range cnames {
		if name == "terminalNetworkHistoryInfo" {
			var c = mongo.MongoClient.Database.C("terminalNetworkHistoryInfo")
			var db mongo.Collection
			db.Collection = c
			return db
		}
	}

	APMacmoduleIDIndex := mgo.Index{
		Key: []string{"APMac", "moduleID"},
	}
	mongo.MongoClient.Database.C("terminalNetworkHistoryInfo").EnsureIndex(APMacmoduleIDIndex)
	mongo.MongoClient.Database.C("terminalNetworkHistoryInfo").EnsureIndexKey("APMac")
	var c = mongo.MongoClient.Database.C("terminalNetworkHistoryInfo")
	var db mongo.Collection
	db.Collection = c
	return db
}

//GetTmnNwkHistoryInfo GetTmnNwkHistoryInfo
func GetTmnNwkHistoryInfo(query bson.M) (*config.TmnNwkInfo, error) {
	var db = getTmnNwkHistoryInfoModel()
	var tmnNwkInfo config.TmnNwkInfo
	err := db.FindOne(query, &tmnNwkInfo)
	if reflect.DeepEqual(tmnNwkInfo, config.TmnNwkInfo{}) {
		return nil, err
	}
	return &tmnNwkInfo, err
}

//GetTmnNwkHistoryInfoCount GetTmnNwkHistoryInfoCount
func GetTmnNwkHistoryInfoCount(query bson.M) (int, error) {
	var db = getTmnNwkHistoryInfoModel()
	count, err := db.Count(query)
	return count, err
}

//CreateTmnNwkHistory CreateTmnNwkHistory
func CreateTmnNwkHistory(setData config.TmnNwkInfo) error {
	var db = getTmnNwkHistoryInfoModel()
	setData.CreateTime = time.Now()
	err := db.Create(setData)
	return err
}

//DeleteTmnNwkHistory DeleteTmnNwkHistory
func DeleteTmnNwkHistory(APMac string) error {
	var db = getTmnNwkHistoryInfoModel()
	err := db.Remove(bson.M{"APMac": APMac})
	return err
}

//DeleteTmnNwkHistoryByOption DeleteTmnNwkHistoryByOption
func DeleteTmnNwkHistoryByOption(option bson.M) error {
	var db = getTmnNwkHistoryInfoModel()
	err := db.Remove(option)
	return err
}

//DeleteTmnNwkHistoryByAPMacAndModuleID DeleteTmnNwkHistoryByAPMacAndModuleID
func DeleteTmnNwkHistoryByAPMacAndModuleID(APMac string, moduleID string) error {
	var db = getTmnNwkHistoryInfoModel()
	err := db.Remove(bson.M{"APMac": APMac, "moduleID": moduleID})
	return err
}
