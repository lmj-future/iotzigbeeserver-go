package models

import (
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/postgres"
	"github.com/jinzhu/gorm"
)

var tmnNwkHistoryInfoModelCheckNum = 0

func getTmnNwkHistoryInfoModelPG() *gorm.DB {
	db := postgres.GetInstance().GetPostGresDB()
	if !db.HasTable("tmnnwkinfos") {
		db.Table("tmnnwkinfos").CreateTable(&config.TmnNwkInfo{})
	} else {
		if tmnNwkHistoryInfoModelCheckNum == 0 {
			tmnNwkHistoryInfoModelCheckNum++
			db.Table("tmnnwkinfos").AutoMigrate(&config.TmnNwkInfo{})
		} else {
			db.Table("tmnnwkinfos")
		}
	}
	return db
}

//GetTmnNwkHistoryInfoPG GetTmnNwkHistoryInfoPG
func GetTmnNwkHistoryInfoPG(query map[string]interface{}) (*config.TmnNwkInfo, error) {
	var db = getTmnNwkHistoryInfoModelPG()
	var tmnNwkInfo config.TmnNwkInfo
	err := db.Where(query).Take(&tmnNwkInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &tmnNwkInfo, err
}

//GetTmnNwkHistoryInfoCountPG GetTmnNwkHistoryInfoCountPG
func GetTmnNwkHistoryInfoCountPG(query map[string]interface{}) (int, error) {
	var db = getTmnNwkHistoryInfoModelPG()
	var count int = 0
	err := db.Model(config.TmnNwkInfo{}).Where(query).Count(&count).Error
	return count, err
}

//CreateTmnNwkHistoryPG CreateTmnNwkHistoryPG
func CreateTmnNwkHistoryPG(setData config.TmnNwkInfo) error {
	var db = getTmnNwkHistoryInfoModelPG()
	setData.CreateTime = time.Now()
	err := db.Create(&setData).Error
	return err
}

//DeleteTmnNwkHistoryPG DeleteTmnNwkHistoryPG
func DeleteTmnNwkHistoryPG(APMac string) error {
	var db = getTmnNwkHistoryInfoModelPG()
	err := db.Where("apmac = ?", APMac).Delete(&config.TmnNwkInfo{}).Error
	return err
}

//DeleteTmnNwkHistoryByOptionPG DeleteTmnNwkHistoryByOptionPG
func DeleteTmnNwkHistoryByOptionPG(option map[string]interface{}) error {
	var db = getTmnNwkHistoryInfoModelPG()
	err := db.Where(option).Delete(&config.TmnNwkInfo{}).Error
	return err
}

//DeleteTmnNwkHistoryByAPMacAndModuleIDPG DeleteTmnNwkHistoryByAPMacAndModuleIDPG
func DeleteTmnNwkHistoryByAPMacAndModuleIDPG(APMac string, moduleID string) error {
	var db = getTmnNwkHistoryInfoModelPG()
	err := db.Where(map[string]interface{}{"apmac": APMac, "moduleid": moduleID}).Delete(&config.TmnNwkInfo{}).Error
	return err
}
