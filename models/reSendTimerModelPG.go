package models

import (
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/postgres"
	"github.com/jinzhu/gorm"
)

var reSendTimerInfoModelCheckNum = 0

func getReSendTimerInfoModelPG() *gorm.DB {
	db := postgres.GetInstance().GetPostGresDB()
	if !db.HasTable("resendtimerinfos") {
		db.Table("resendtimerinfos").CreateTable(&config.ReSendTimerInfo{})
	} else {
		if reSendTimerInfoModelCheckNum == 0 {
			reSendTimerInfoModelCheckNum++
			db.Table("resendtimerinfos").AutoMigrate(&config.ReSendTimerInfo{})
		} else {
			db.Table("resendtimerinfos")
		}
	}
	return db
}

//GetReSendTimerInfoByMsgKeyPG GetReSendTimerInfoByMsgKeyPG
func GetReSendTimerInfoByMsgKeyPG(msgKey string) (*config.ReSendTimerInfo, error) {
	var db = getReSendTimerInfoModelPG()
	var reSendTimerInfo config.ReSendTimerInfo
	err := db.Where("msgkey = ?", msgKey).Take(&reSendTimerInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &reSendTimerInfo, err
}

//FindReSendTimerInfoAndUpdatePG FindReSendTimerInfoAndUpdatePG
func FindReSendTimerInfoAndUpdatePG(msgKey string, oSet map[string]interface{}) (*config.ReSendTimerInfo, error) {
	var db = getReSendTimerInfoModelPG()
	var reSendTimerInfo config.ReSendTimerInfo
	err := db.Where("msgkey = ?", msgKey).Take(&reSendTimerInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	err = db.Model(config.ReSendTimerInfo{}).Where("msgkey = ?", msgKey).Update(oSet).Error
	return &reSendTimerInfo, err
}

//CreateReSendTimerInfoPG CreateReSendTimerInfoPG
func CreateReSendTimerInfoPG(setData config.ReSendTimerInfo) error {
	var db = getReSendTimerInfoModelPG()
	setData.CreateTime = time.Now()
	err := db.Create(&setData).Error
	return err
}

//DeleteReSendTimerPG DeleteReSendTimerPG
func DeleteReSendTimerPG(msgKey string) error {
	var db = getReSendTimerInfoModelPG()
	err := db.Where("msgkey = ?", msgKey).Delete(&config.ReSendTimerInfo{}).Error
	return err
}
