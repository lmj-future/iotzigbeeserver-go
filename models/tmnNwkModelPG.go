package models

import (
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/postgres"
	"github.com/jinzhu/gorm"
)

var tmnNwkInfoModelCheckNum = 0

func getTmnNwkInfoModelPG() *gorm.DB {
	db := postgres.GetInstance().GetPostGresDB()
	if !db.HasTable("tmnnwkinfos") {
		db.Table("tmnnwkinfos").CreateTable(&config.TmnNwkInfo{})
	} else {
		if tmnNwkInfoModelCheckNum == 0 {
			tmnNwkInfoModelCheckNum++
			db.Table("tmnnwkinfos").AutoMigrate(&config.TmnNwkInfo{})
		} else {
			db.Table("tmnnwkinfos")
		}
	}
	return db
}

//GetTmnNwkInfoPG GetTmnNwkInfoPG
func GetTmnNwkInfoPG(query map[string]interface{}) (*config.TmnNwkInfo, error) {
	var db = getTmnNwkInfoModelPG()
	var tmnNwkInfo config.TmnNwkInfo
	err := db.Where(query).Take(&tmnNwkInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &tmnNwkInfo, err
}

//FindTmnNwkAndUpdatePG FindTmnNwkAndUpdatePG
func FindTmnNwkAndUpdatePG(query map[string]interface{}, oSet map[string]interface{}) (*config.TmnNwkInfo, error) {
	var db = getTmnNwkInfoModelPG()
	var tmnNwkInfo config.TmnNwkInfo
	err := db.Where(query).Take(&tmnNwkInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	err = db.Model(config.TmnNwkInfo{}).Where(query).Update(oSet).Error
	return &tmnNwkInfo, err
}

//CreateTmnNwkPG CreateTmnNwkPG
func CreateTmnNwkPG(setData config.TmnNwkInfo) error {
	var db = getTmnNwkInfoModelPG()
	setData.CreateTime = time.Now()
	err := db.Create(&setData).Error
	return err
}

//DeleteTmnNwkPG DeleteTmnNwkPG
func DeleteTmnNwkPG(APMac string) error {
	var db = getTmnNwkInfoModelPG()
	err := db.Where("apmac = ?", APMac).Delete(&config.TmnNwkInfo{}).Error
	return err
}

//DeleteTmnNwkByAPMacAndModuleIDPG DeleteTmnNwkByAPMacAndModuleIDPG
func DeleteTmnNwkByAPMacAndModuleIDPG(APMac string, moduleID string) error {
	var db = getTmnNwkInfoModelPG()
	err := db.Where(map[string]interface{}{"apmac": APMac, "moduleid": moduleID}).Delete(&config.TmnNwkInfo{}).Error
	return err
}
