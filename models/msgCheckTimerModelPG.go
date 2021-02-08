package models

import (
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/postgres"
	"github.com/jinzhu/gorm"
)

var msgCheckTimerInfoModelCheckNum = 0

func getMsgCheckTimerInfoModelPG() *gorm.DB {
	db := postgres.GetInstance().GetPostGresDB()
	if !db.HasTable("msgchecktimerinfos") {
		db.Table("msgchecktimerinfos").CreateTable(&config.MsgCheckTimerInfo{})
	} else {
		if msgCheckTimerInfoModelCheckNum == 0 {
			msgCheckTimerInfoModelCheckNum++
			db.Table("msgchecktimerinfos").AutoMigrate(&config.MsgCheckTimerInfo{})
		} else {
			db.Table("msgchecktimerinfos")
		}
		db.Table("msgchecktimerinfos")
	}
	return db
}

//GetMsgCheckTimerInfoByMsgKeyPG GetMsgCheckTimerInfoByMsgKeyPG
func GetMsgCheckTimerInfoByMsgKeyPG(msgKey string) (*config.MsgCheckTimerInfo, error) {
	var db = getMsgCheckTimerInfoModelPG()
	var msgCheckTimerInfo config.MsgCheckTimerInfo
	err := db.Where("msgkey = ?", msgKey).Take(&msgCheckTimerInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &msgCheckTimerInfo, err
}

//FindMsgCheckTimerInfoAndUpdatePG FindMsgCheckTimerInfoAndUpdatePG
func FindMsgCheckTimerInfoAndUpdatePG(msgKey string, oSet map[string]interface{}) (*config.MsgCheckTimerInfo, error) {
	var db = getMsgCheckTimerInfoModelPG()
	var msgCheckTimerInfo config.MsgCheckTimerInfo
	err := db.Where("msgkey = ?", msgKey).Take(&msgCheckTimerInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	err = db.Model(config.MsgCheckTimerInfo{}).Where("msgkey = ?", msgKey).Update(oSet).Error
	return &msgCheckTimerInfo, err
}

//CreateMsgCheckTimerInfoPG CreateMsgCheckTimerInfoPG
func CreateMsgCheckTimerInfoPG(setData config.MsgCheckTimerInfo) error {
	var db = getMsgCheckTimerInfoModelPG()
	setData.CreateTime = time.Now()
	err := db.Create(&setData).Error
	return err
}

//DeleteMsgCheckTimerPG DeleteMsgCheckTimerPG
func DeleteMsgCheckTimerPG(msgKey string) error {
	var db = getMsgCheckTimerInfoModelPG()
	err := db.Where("msgkey = ?", msgKey).Delete(&config.MsgCheckTimerInfo{}).Error
	return err
}
