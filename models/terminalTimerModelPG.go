package models

import (
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/postgres"
	"github.com/jinzhu/gorm"
)

var terminalTimerInfoModelCheckNum = 0

func getTerminalTimerInfoModelPG() *gorm.DB {
	db := postgres.GetInstance().GetPostGresDB()
	if !db.HasTable("terminaltimerinfos") {
		db.Table("terminaltimerinfos").CreateTable(&config.TerminalTimerInfo{})
	} else {
		if terminalTimerInfoModelCheckNum == 0 {
			terminalTimerInfoModelCheckNum++
			db.Table("terminaltimerinfos").AutoMigrate(&config.TerminalTimerInfo{})
		} else {
			db.Table("terminaltimerinfos")
		}
	}
	return db
}

//GetTerminalTimerInfoByDevEUIPG GetTerminalTimerInfoByDevEUIPG
func GetTerminalTimerInfoByDevEUIPG(devEUI string) (*config.TerminalTimerInfo, error) {
	var db = getTerminalTimerInfoModelPG()
	var terminalTimer config.TerminalTimerInfo
	err := db.Where("deveui = ?", devEUI).Take(&terminalTimer).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &terminalTimer, err
}

//FindTerminalTimerInfoAndUpdatePG FindTerminalTimerInfoAndUpdatePG
func FindTerminalTimerInfoAndUpdatePG(devEUI string, oSet map[string]interface{}) (*config.TerminalTimerInfo, error) {
	var db = getTerminalTimerInfoModelPG()
	var terminalTimer config.TerminalTimerInfo
	err := db.Where("deveui = ?", devEUI).Take(&terminalTimer).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	err = db.Model(config.TerminalTimerInfo{}).Where("deveui = ?", devEUI).Update(oSet).Error
	return &terminalTimer, err
}

//CreateTerminalTimerInfoPG CreateTerminalTimerInfoPG
func CreateTerminalTimerInfoPG(setData config.TerminalTimerInfo) error {
	var db = getTerminalTimerInfoModelPG()
	setData.CreateTime = time.Now()
	err := db.Create(&setData).Error
	return err
}

//DeleteTerminalTimerPG DeleteTerminalTimerPG
func DeleteTerminalTimerPG(devEUI string) error {
	var db = getTerminalTimerInfoModelPG()
	err := db.Where("deveui = ?", devEUI).Delete(&config.TerminalTimerInfo{}).Error
	return err
}
