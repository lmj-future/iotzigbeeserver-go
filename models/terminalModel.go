package models

import (
	"reflect"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/mongo"
)

func getTerminalInfoModel() mongo.DBClient {
	cnames, _ := mongo.MongoClient.Database.CollectionNames()
	for _, name := range cnames {
		if name == "terminalInfo" {
			var c = mongo.MongoClient.Database.C("terminalInfo")
			var db mongo.Collection
			db.Collection = c
			return db
		}
	}

	devEUIIndex := mgo.Index{
		Key:    []string{"devEUI"},
		Unique: true,
	}
	mongo.MongoClient.Database.C("terminalInfo").EnsureIndex(devEUIIndex)
	mongo.MongoClient.Database.C("terminalInfo").EnsureIndexKey("OIDIndex")
	mongo.MongoClient.Database.C("terminalInfo").EnsureIndexKey("APMac", "devEUI")
	mongo.MongoClient.Database.C("terminalInfo").EnsureIndexKey("APMac", "moduleID")
	mongo.MongoClient.Database.C("terminalInfo").EnsureIndexKey("APMac")
	mongo.MongoClient.Database.C("terminalInfo").EnsureIndexKey("scenarioID")
	var c = mongo.MongoClient.Database.C("terminalInfo")
	var db mongo.Collection
	db.Collection = c
	return db
}

//GetTerminalInfoByDevEUI GetTerminalInfoByDevEUI
func GetTerminalInfoByDevEUI(devEUI string) (*config.TerminalInfo, error) {
	var db = getTerminalInfoModel()
	var filit = bson.M{
		"devEUI": devEUI,
	}
	var terminalInfo config.TerminalInfo
	err := db.FindOne(filit, &terminalInfo)
	if reflect.DeepEqual(terminalInfo, config.TerminalInfo{}) {
		return nil, err
	}
	return &terminalInfo, err
}

//GetTerminalInfoByOIDIndex GetTerminalInfoByOIDIndex
func GetTerminalInfoByOIDIndex(OIDIndex string) (*config.TerminalInfo, error) {
	var db = getTerminalInfoModel()
	var filit = bson.M{
		"OIDIndex": OIDIndex,
	}
	var terminalInfo config.TerminalInfo
	err := db.FindOne(filit, &terminalInfo)
	if reflect.DeepEqual(terminalInfo, config.TerminalInfo{}) {
		return nil, err
	}
	return &terminalInfo, err
}

//FindTerminalAndUpdateByOIDIndex FindTerminalAndUpdateByOIDIndex
func FindTerminalAndUpdateByOIDIndex(OIDIndex string, setData interface{}) (*config.TerminalInfo, error) {
	var db = getTerminalInfoModel()
	var terminalInfo config.TerminalInfo
	err := db.FindOneAndUpdate(bson.M{"OIDIndex": OIDIndex}, setData, &terminalInfo)
	if reflect.DeepEqual(terminalInfo, config.TerminalInfo{}) {
		return nil, err
	}
	return &terminalInfo, err
}

//FindTerminalAndUpdate FindTerminalAndUpdate
func FindTerminalAndUpdate(oMatch bson.M, setData interface{}) (*config.TerminalInfo, error) {
	var db = getTerminalInfoModel()
	var terminalInfo config.TerminalInfo
	err := db.FindOneAndUpdate(oMatch, setData, &terminalInfo)
	if reflect.DeepEqual(terminalInfo, config.TerminalInfo{}) {
		return nil, err
	}
	return &terminalInfo, err
}

//FindTerminalAndUpdateByAPMac FindTerminalAndUpdateByAPMac
// func FindTerminalAndUpdateByAPMac(devEUI string, APMac string, setData config.TerminalInfo) (*config.TerminalInfo, error) {
// 	var db = getTerminalInfoModel()
// 	var terminalInfo config.TerminalInfo
// 	err := db.FindOneAndUpdate(bson.M{"APMac": APMac, "devEUI": devEUI}, setData, &terminalInfo)
// 	if reflect.DeepEqual(terminalInfo, config.TerminalInfo{}) {
// 		return nil, err
// 	}
// 	return &terminalInfo, err
// }

//CreateTerminal CreateTerminal
func CreateTerminal(setData config.TerminalInfo) error {
	var db = getTerminalInfoModel()
	setData.CreateTime = time.Now()
	err := db.Create(setData)
	return err
}

//DeleteTerminal DeleteTerminal
func DeleteTerminal(devEUI string) error {
	var db = getTerminalInfoModel()
	err := db.Remove(bson.M{"devEUI": devEUI})
	return err
}

//DeleteTerminalList DeleteTerminalList
// func DeleteTerminalList(devEUIList []string) error {
// 	var db = getTerminalInfoModel()
// 	err := db.Remove(bson.M{"$in": devEUIList})
// 	return err
// }

//FindTerminalBySecondAddrAndModuleID FindTerminalBySecondAddrAndModuleID
func FindTerminalBySecondAddrAndModuleID(secondAddr string, moduleID string) ([]config.TerminalInfo, error) {
	var db = getTerminalInfoModel()
	var terminalInfos []config.TerminalInfo
	err := db.Find(bson.M{"secondAddr": secondAddr, "moduleID": moduleID}, &terminalInfos)
	if reflect.DeepEqual(terminalInfos, []config.TerminalInfo{}) {
		return nil, err
	}
	return terminalInfos, err
}

//FindTerminalByAPMacAndModuleID FindTerminalByAPMacAndModuleID
func FindTerminalByAPMacAndModuleID(APMac string, moduleID string) ([]config.TerminalInfo, error) {
	var db = getTerminalInfoModel()
	var terminalInfos []config.TerminalInfo
	err := db.Find(bson.M{"APMac": APMac, "moduleID": moduleID}, &terminalInfos)
	if reflect.DeepEqual(terminalInfos, []config.TerminalInfo{}) {
		return nil, err
	}
	return terminalInfos, err
}

//FindTerminalByAPMac FindTerminalByAPMac
func FindTerminalByAPMac(APMac string) ([]config.TerminalInfo, error) {
	var db = getTerminalInfoModel()
	var terminalInfos []config.TerminalInfo
	err := db.Find(bson.M{"APMac": APMac}, &terminalInfos)
	if reflect.DeepEqual(terminalInfos, []config.TerminalInfo{}) {
		return nil, err
	}
	return terminalInfos, err
}

//DeleteTerminalByScenarioID DeleteTerminalByScenarioID
func DeleteTerminalByScenarioID(scenarioID string) error {
	var db = getTerminalInfoModel()
	err := db.Remove(bson.M{"scenarioID": scenarioID})
	return err
}

//FindTerminalByCondition FindTerminalByCondition
func FindTerminalByCondition(condition bson.M) (*config.TerminalInfo, error) {
	var db = getTerminalInfoModel()
	var terminalInfo config.TerminalInfo
	err := db.FindOne(condition, &terminalInfo)
	if reflect.DeepEqual(terminalInfo, config.TerminalInfo{}) {
		return nil, err
	}
	return &terminalInfo, err
}

//FindTerminalListByCondition FindTerminalListByCondition
func FindTerminalListByCondition(condition bson.M, field bson.M) ([]config.TerminalInfo2, error) {
	var db = getTerminalInfoModel()
	var terminalInfos []config.TerminalInfo2
	err := db.FindField(condition, field, &terminalInfos)
	if reflect.DeepEqual(terminalInfos, []config.TerminalInfo2{}) {
		return nil, err
	}
	return terminalInfos, err
}

//FindAllTerminalByCodition FindAllTerminalByCodition
func FindAllTerminalByCodition(condition bson.M) ([]config.TerminalInfo, error) {
	var db = getTerminalInfoModel()
	var terminalInfos []config.TerminalInfo
	err := db.Find(condition, &terminalInfos)
	if reflect.DeepEqual(terminalInfos, []config.TerminalInfo{}) {
		return nil, err
	}
	return terminalInfos, err
}
