package models

import (
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/db/postgres"
	"github.com/jinzhu/gorm"
)

var terminalInfoModelCheckNum = 0

func getTerminalInfoModelPG() *gorm.DB {
	db := postgres.GetInstance().GetPostGresDB()
	if !db.HasTable("zigbeeterminalinfos") {
		db.Table("zigbeeterminalinfos").CreateTable(&config.TerminalInfo{})
	} else {
		if terminalInfoModelCheckNum == 0 {
			terminalInfoModelCheckNum++
			db.Table("zigbeeterminalinfos").AutoMigrate(&config.TerminalInfo{})
		} else {
			db.Table("zigbeeterminalinfos")
		}
	}
	return db
}

//GetTerminalInfoByDevEUIPG GetTerminalInfoByDevEUIPG
func GetTerminalInfoByDevEUIPG(devEUI string) (*config.TerminalInfo, error) {
	var db = getTerminalInfoModelPG()
	var terminalInfo config.TerminalInfo
	err := db.Where("deveui = ?", devEUI).Take(&terminalInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &terminalInfo, err
}

//GetTerminalInfoByOIDIndexPG GetTerminalInfoByOIDIndexPG
func GetTerminalInfoByOIDIndexPG(OIDIndex string) (*config.TerminalInfo, error) {
	var db = getTerminalInfoModelPG()
	var terminalInfo config.TerminalInfo
	err := db.Where("oidindex = ?", OIDIndex).Take(&terminalInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &terminalInfo, err
}

//FindTerminalAndUpdateByOIDIndexPG FindTerminalAndUpdateByOIDIndexPG
func FindTerminalAndUpdateByOIDIndexPG(OIDIndex string, oSet map[string]interface{}) (*config.TerminalInfo, error) {
	var db = getTerminalInfoModelPG()
	var terminalInfo config.TerminalInfo
	err := db.Where("oidindex = ?", OIDIndex).Take(&terminalInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	err = db.Model(config.TerminalInfo{}).Where("oidindex = ?", OIDIndex).Update(oSet).Error
	return &terminalInfo, err
}

//FindTerminalAndUpdatePG FindTerminalAndUpdatePG
func FindTerminalAndUpdatePG(oMatch map[string]interface{}, oSet map[string]interface{}) (*config.TerminalInfo, error) {
	var db = getTerminalInfoModelPG()
	var terminalInfo config.TerminalInfo
	err := db.Where(oMatch).Take(&terminalInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	err = db.Model(config.TerminalInfo{}).Where(oMatch).Update(oSet).Error
	return &terminalInfo, err
}

// FindTerminalAndUpdateByAPMacPG FindTerminalAndUpdateByAPMacPG
func FindTerminalAndUpdateByAPMacPG(devEUI string, APMac string, oSet map[string]interface{}) (*config.TerminalInfo, error) {
	var db = getTerminalInfoModelPG()
	var terminalInfo config.TerminalInfo
	err := db.Where(map[string]interface{}{"apmac": APMac, "deveui": devEUI}).Take(&terminalInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	err = db.Model(config.TerminalInfo{}).Where(map[string]interface{}{"apmac": APMac, "deveui": devEUI}).Update(oSet).Error
	return &terminalInfo, err
}

//CreateTerminalPG CreateTerminalPG
func CreateTerminalPG(setData config.TerminalInfo) error {
	var db = getTerminalInfoModelPG()
	setData.CreateTime = time.Now()
	err := db.Create(&setData).Error
	if err != nil {
		db.AutoMigrate(&config.TerminalInfo{})
		err = db.Create(&setData).Error
	}
	return err
}

// DeleteTerminalPG DeleteTerminalPG
func DeleteTerminalPG(devEUI string) error {
	var db = getTerminalInfoModelPG()
	err := db.Where("deveui = ?", devEUI).Delete(&config.TerminalInfo{}).Error
	return err
}

// DeleteTerminalListPG DeleteTerminalListPG
// func DeleteTerminalListPG(devEUIList []string) error {
// 	var db = getTerminalInfoModelPG()
// 	err := db.Where("deveui IN ?", devEUIList).Delete(&config.TerminalInfo{}).Error
// 	return err
// }

// FindTerminalBySecondAddrAndModuleIDPG FindTerminalBySecondAddrAndModuleIDPG
func FindTerminalBySecondAddrAndModuleIDPG(secondAddr string, moduleID string) ([]config.TerminalInfo, error) {
	var db = getTerminalInfoModelPG()
	var terminalInfos []config.TerminalInfo
	err := db.Where(map[string]interface{}{"secondaddr": secondAddr, "moduleid": moduleID}).Find(&terminalInfos).Error
	return terminalInfos, err
}

// FindTerminalByAPMacAndModuleIDPG FindTerminalByAPMacAndModuleIDPG
func FindTerminalByAPMacAndModuleIDPG(APMac string, moduleID string) ([]config.TerminalInfo, error) {
	var db = getTerminalInfoModelPG()
	var terminalInfos []config.TerminalInfo
	err := db.Where(map[string]interface{}{"apmac": APMac, "moduleid": moduleID}).Find(&terminalInfos).Error
	return terminalInfos, err
}

//FindTerminalByAPMacPG FindTerminalByAPMacPG
func FindTerminalByAPMacPG(APMac string) ([]config.TerminalInfo, error) {
	var db = getTerminalInfoModelPG()
	var terminalInfos []config.TerminalInfo
	err := db.Where("apmac = ?", APMac).Find(&terminalInfos).Error
	return terminalInfos, err
}

//DeleteTerminalByScenarioIDPG DeleteTerminalByScenarioIDPG
func DeleteTerminalByScenarioIDPG(scenarioID string) error {
	var db = getTerminalInfoModelPG()
	err := db.Where("scenarioid = ?", scenarioID).Delete(&config.TerminalInfo{}).Error
	return err
}

//FindTerminalByConditionPG FindTerminalByConditionPG
func FindTerminalByConditionPG(condition map[string]interface{}) (*config.TerminalInfo, error) {
	var db = getTerminalInfoModelPG()
	var terminalInfo config.TerminalInfo
	err := db.Where(condition).Take(&terminalInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &terminalInfo, err
}

//FindTerminalListByConditionPG FindTerminalListByConditionPG
func FindTerminalListByConditionPG(condition map[string]interface{}, field []string) ([]config.TerminalInfo, error) {
	var db = getTerminalInfoModelPG()
	var terminalInfos []config.TerminalInfo
	result := db.Select(field).Where(condition).Find(&terminalInfos)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil, result.Error
	}
	return terminalInfos, result.Error
}

//FindAllTerminalByConditionPG FindAllTerminalByConditionPG
func FindAllTerminalByConditionPG(condition map[string]interface{}) ([]config.TerminalInfo, error) {
	var db = getTerminalInfoModelPG()
	var terminalInfos []config.TerminalInfo
	result := db.Where(condition).Find(&terminalInfos)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil, result.Error
	}
	return terminalInfos, result.Error
}

//FindTerminalListByPagePG FindTerminalListByPagePG
func FindTerminalListByPagePG(condition map[string]interface{}, or map[string]interface{},
	field []string, pageNum int, pageSize int) ([]config.TerminalInfo, int, error) {
	var db = getTerminalInfoModelPG()
	var terminalInfos []config.TerminalInfo
	var totalCount int
	db.Model(config.TerminalInfo{}).Where(condition).Or(or).Count(&totalCount)
	result := db.Select(field).Where(condition).Or(or).Order("id").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&terminalInfos)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil, 0, result.Error
	}
	return terminalInfos, totalCount, result.Error
}
