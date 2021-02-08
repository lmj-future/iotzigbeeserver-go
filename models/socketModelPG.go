package models

import (
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/postgres"
	"github.com/jinzhu/gorm"
)

var socketInfoModelCheckNum = 0

func getSocketInfoModelPG() *gorm.DB {
	db := postgres.GetInstance().GetPostGresDB()
	if !db.HasTable("socketinfos") {
		db.Table("socketinfos").CreateTable(&config.SocketInfo{})
	} else {
		if socketInfoModelCheckNum == 0 {
			socketInfoModelCheckNum++
			db.Table("socketinfos").AutoMigrate(&config.SocketInfo{})
		} else {
			db.Table("socketinfos")
		}
	}
	return db
}

//GetSocketInfoByAPMacPG GetSocketInfoByAPMacPG
func GetSocketInfoByAPMacPG(APMac string) (*config.SocketInfo, error) {
	var db = getSocketInfoModelPG()
	var socketInfo config.SocketInfo
	err := db.Where("apmac = ?", APMac).Take(&socketInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &socketInfo, err
}

//FindSocketInfoAndUpdatePG FindSocketInfoAndUpdatePG
func FindSocketInfoAndUpdatePG(APMac string, oSet map[string]interface{}) (*config.SocketInfo, error) {
	var db = getSocketInfoModelPG()
	var socketInfo config.SocketInfo
	err := db.Where("apmac = ?", APMac).Take(&socketInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	err = db.Model(config.SocketInfo{}).Where("apmac = ?", APMac).Update(oSet).Error
	return &socketInfo, err
}

//CreateSocketInfoPG CreateSocketInfoPG
func CreateSocketInfoPG(setData config.SocketInfo) error {
	var db = getSocketInfoModelPG()
	setData.CreateTime = time.Now()
	err := db.Create(&setData).Error
	return err
}

//DeleteSocketPG DeleteSocketPG
func DeleteSocketPG(APMac string) error {
	var db = getSocketInfoModelPG()
	err := db.Where("apmac = ?", APMac).Delete(&config.SocketInfo{}).Error
	return err
}
