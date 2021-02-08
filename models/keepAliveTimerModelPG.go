package models

import (
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/postgres"
	"github.com/jinzhu/gorm"
)

var keepAliveTimerInfoModelCheckNum = 0

func getKeepAliveTimerInfoModelPG() *gorm.DB {
	db := postgres.GetInstance().GetPostGresDB()
	if !db.HasTable("keepalivetimerinfos") {
		db.Table("keepalivetimerinfos").CreateTable(&config.KeepAliveTimerInfo{})
	} else {
		if keepAliveTimerInfoModelCheckNum == 0 {
			keepAliveTimerInfoModelCheckNum++
			db.Table("keepalivetimerinfos").AutoMigrate(&config.KeepAliveTimerInfo{})
		} else {
			db.Table("keepalivetimerinfos")
		}
	}
	return db
}

//GetKeepAliveTimerInfoByAPMacPG GetKeepAliveTimerInfoByAPMacPG
func GetKeepAliveTimerInfoByAPMacPG(APMac string) (*config.KeepAliveTimerInfo, error) {
	var db = getKeepAliveTimerInfoModelPG()
	var keepAliveTimerInfo config.KeepAliveTimerInfo
	err := db.Where("apmac = ?", APMac).Take(&keepAliveTimerInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &keepAliveTimerInfo, err
}

//FindKeepAliveTimerInfoAndUpdatePG FindKeepAliveTimerInfoAndUpdatePG
func FindKeepAliveTimerInfoAndUpdatePG(APMac string, oSet map[string]interface{}) (*config.KeepAliveTimerInfo, error) {
	var db = getKeepAliveTimerInfoModelPG()
	var keepAliveTimerInfo config.KeepAliveTimerInfo
	err := db.Where("apmac = ?", APMac).Take(&keepAliveTimerInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	err = db.Model(config.KeepAliveTimerInfo{}).Where("apmac = ?", APMac).Update(oSet).Error
	return &keepAliveTimerInfo, err
}

//CreateKeepAliveTimerInfoPG CreateKeepAliveTimerInfoPG
func CreateKeepAliveTimerInfoPG(setData config.KeepAliveTimerInfo) error {
	var db = getKeepAliveTimerInfoModelPG()
	setData.CreateTime = time.Now()
	err := db.Create(&setData).Error
	return err
}

//DeleteKeepAliveTimerPG DeleteKeepAliveTimerPG
func DeleteKeepAliveTimerPG(APMac string) error {
	var db = getKeepAliveTimerInfoModelPG()
	err := db.Where("apmac = ?", APMac).Delete(&config.KeepAliveTimerInfo{}).Error
	return err
}
