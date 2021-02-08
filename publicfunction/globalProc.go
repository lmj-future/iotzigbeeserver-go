package publicfunction

import (
	"math/rand"
	"sync"
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmemorycache"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalrediscache"
	"github.com/h3c/iotzigbeeserver-go/models"
)

var zigbeeServerMsgCheckTimerID = &sync.Map{}

/******************************************************************************************/
/* keepAliveTimerModel Func Start */
/******************************************************************************************/

// GetKeepAliveTimerByAPMac GetKeepAliveTimerByAPMac
func GetKeepAliveTimerByAPMac(APMac string) (*config.KeepAliveTimerInfo, error) {
	if constant.Constant.UsePostgres {
		keepAliveTimerInfo, err := models.GetKeepAliveTimerInfoByAPMacPG(APMac)
		if err != nil {
			globallogger.Log.Errorln("GetKeepAliveTimerByAPMac GetKeepAliveTimerInfoByAPMacPG error: ", err)
		} else {
			// globallogger.Log.Infoln("GetKeepAliveTimerByAPMac GetKeepAliveTimerInfoByAPMacPG success")
		}
		return keepAliveTimerInfo, err
	}
	keepAliveTimerInfo, err := models.GetKeepAliveTimerInfoByAPMac(APMac)
	if err != nil {
		globallogger.Log.Errorln("GetKeepAliveTimerByAPMac getKeepAliveTimerInfoByAPMac error: ", err)
	} else {
		// globallogger.Log.Infoln("GetKeepAliveTimerByAPMac getKeepAliveTimerInfoByAPMac success")
	}
	return keepAliveTimerInfo, err
}

// FindKeepAliveTimerAndUpdate FindKeepAliveTimerAndUpdate
func FindKeepAliveTimerAndUpdate(APMac string, setData config.KeepAliveTimerInfo) (*config.KeepAliveTimerInfo, error) {
	keepAliveTimerInfo, err := models.FindKeepAliveTimerInfoAndUpdate(APMac, setData)
	if err != nil {
		globallogger.Log.Errorln("FindKeepAliveTimerAndUpdate findKeepAliveTimerInfoAndUpdate error: ", err)
	} else {
		// globallogger.Log.Infoln("FindKeepAliveTimerAndUpdate findKeepAliveTimerInfoAndUpdate success")
	}
	return keepAliveTimerInfo, err
}

// FindKeepAliveTimerAndUpdatePG FindKeepAliveTimerAndUpdatePG
func FindKeepAliveTimerAndUpdatePG(APMac string, oSet map[string]interface{}) (*config.KeepAliveTimerInfo, error) {
	keepAliveTimerInfo, err := models.FindKeepAliveTimerInfoAndUpdatePG(APMac, oSet)
	if err != nil {
		globallogger.Log.Errorln("FindKeepAliveTimerAndUpdatePG findKeepAliveTimerInfoAndUpdatePG error: ", err)
	} else {
		// globallogger.Log.Infoln("FindKeepAliveTimerAndUpdatePG findKeepAliveTimerInfoAndUpdatePG success")
	}
	return keepAliveTimerInfo, err
}

// CreateKeepAliveTimer CreateKeepAliveTimer
func CreateKeepAliveTimer(APMac string) error {
	var setData = config.KeepAliveTimerInfo{}
	setData.APMac = APMac
	setData.UpdateTime = time.Now()
	if constant.Constant.UsePostgres {
		err := models.CreateKeepAliveTimerInfoPG(setData)
		if err != nil {
			globallogger.Log.Errorln("CreateKeepAliveTimer CreateKeepAliveTimerInfoPG error: ", err)
			oSet := make(map[string]interface{})
			oSet["apmac"] = APMac
			oSet["updatetime"] = time.Now()
			_, err = FindKeepAliveTimerAndUpdatePG(APMac, oSet)
		} else {
			// globallogger.Log.Infoln("CreateKeepAliveTimer CreateKeepAliveTimerInfoPG success")
		}
		return err
	}
	err := models.CreateKeepAliveTimerInfo(setData)
	if err != nil {
		globallogger.Log.Errorln("CreateKeepAliveTimer CreateKeepAliveTimerInfo error: ", err)
		_, err = FindKeepAliveTimerAndUpdate(APMac, setData)
	} else {
		// globallogger.Log.Infoln("CreateKeepAliveTimer CreateKeepAliveTimerInfo success")
	}
	return err
}

// DeleteKeepAliveTimer DeleteKeepAliveTimer
func DeleteKeepAliveTimer(APMac string) {
	var err error
	if constant.Constant.UsePostgres {
		err = models.DeleteKeepAliveTimerPG(APMac)
	} else {
		err = models.DeleteKeepAliveTimer(APMac)
	}
	if err != nil {
		globallogger.Log.Errorln("DeleteKeepAliveTimer DeleteKeepAliveTimer error: ", err)
	} else {
		// globallogger.Log.Infoln("DeleteKeepAliveTimer DeleteKeepAliveTimer success")
	}
}

// KeepAliveTimerUpdateOrCreate KeepAliveTimerUpdateOrCreate
func KeepAliveTimerUpdateOrCreate(APMac string) {
	var setData = config.KeepAliveTimerInfo{}
	setData.UpdateTime = time.Now()
	var keepAliveTimerInfo *config.KeepAliveTimerInfo
	if constant.Constant.UsePostgres {
		oSet := make(map[string]interface{})
		oSet["updatetime"] = time.Now()
		keepAliveTimerInfo, _ = FindKeepAliveTimerAndUpdatePG(APMac, oSet)
	} else {
		keepAliveTimerInfo, _ = FindKeepAliveTimerAndUpdate(APMac, setData)
	}
	if keepAliveTimerInfo == nil {
		CreateKeepAliveTimer(APMac)
	}
}

/****************************************************************************************/
/* keepAliveTimerModel Func End */
/****************************************************************************************/

/*****************************************************************************************/
/* msgCheckTimerModel Func Start */
/*****************************************************************************************/

// GetMsgCheckTimerByMsgKey GetMsgCheckTimerByMsgKey
func GetMsgCheckTimerByMsgKey(MsgKey string) (*config.MsgCheckTimerInfo, error) {
	if constant.Constant.UsePostgres {
		msgCheckTimerInfo, err := models.GetMsgCheckTimerInfoByMsgKeyPG(MsgKey)
		if err != nil {
			globallogger.Log.Errorln("GetMsgCheckTimerByMsgKey GetMsgCheckTimerInfoByMsgKeyPG error: ", err)
		} else {
			// globallogger.Log.Infoln("GetMsgCheckTimerByMsgKey GetMsgCheckTimerInfoByMsgKeyPG success")
		}
		return msgCheckTimerInfo, err
	}
	msgCheckTimerInfo, err := models.GetMsgCheckTimerInfoByMsgKey(MsgKey)
	if err != nil {
		globallogger.Log.Errorln("GetMsgCheckTimerByMsgKey GetMsgCheckTimerInfoByMsgKey error: ", err)
	} else {
		// globallogger.Log.Infoln("GetMsgCheckTimerByMsgKey GetMsgCheckTimerInfoByMsgKey success")
	}
	return msgCheckTimerInfo, err
}

// FindMsgCheckTimerAndUpdate FindMsgCheckTimerAndUpdate
func FindMsgCheckTimerAndUpdate(MsgKey string, setData config.MsgCheckTimerInfo) (*config.MsgCheckTimerInfo, error) {
	msgCheckTimerInfo, err := models.FindMsgCheckTimerInfoAndUpdate(MsgKey, setData)
	if err != nil {
		globallogger.Log.Errorln("FindMsgCheckTimerAndUpdate FindMsgCheckTimerInfoAndUpdate error: ", err)
	} else {
		// globallogger.Log.Infoln("FindMsgCheckTimerAndUpdate FindMsgCheckTimerInfoAndUpdate success")
	}
	return msgCheckTimerInfo, err
}

// FindMsgCheckTimerAndUpdatePG FindMsgCheckTimerAndUpdatePG
func FindMsgCheckTimerAndUpdatePG(MsgKey string, oSet map[string]interface{}) (*config.MsgCheckTimerInfo, error) {
	msgCheckTimerInfo, err := models.FindMsgCheckTimerInfoAndUpdatePG(MsgKey, oSet)
	if err != nil {
		globallogger.Log.Errorln("FindMsgCheckTimerAndUpdatePG FindMsgCheckTimerInfoAndUpdatePG error: ", err)
	} else {
		// globallogger.Log.Infoln("FindMsgCheckTimerAndUpdatePG FindMsgCheckTimerInfoAndUpdatePG success")
	}
	return msgCheckTimerInfo, err
}

// CreateMsgCheckTimer CreateMsgCheckTimer
func CreateMsgCheckTimer(MsgKey string) (*config.MsgCheckTimerInfo, error) {
	var setData = config.MsgCheckTimerInfo{}
	setData.MsgKey = MsgKey
	setData.UpdateTime = time.Now()
	if constant.Constant.UsePostgres {
		err := models.CreateMsgCheckTimerInfoPG(setData)
		if err != nil {
			globallogger.Log.Errorln("CreateMsgCheckTimer CreateMsgCheckTimerInfoPG error: ", err)
			oSet := make(map[string]interface{})
			oSet["msgkey"] = MsgKey
			oSet["updatetime"] = time.Now()
			_, err = FindMsgCheckTimerAndUpdatePG(MsgKey, oSet)
		} else {
			// globallogger.Log.Infoln("CreateMsgCheckTimer CreateMsgCheckTimerInfoPG success")
		}
		return nil, err
	}
	err := models.CreateMsgCheckTimerInfo(setData)
	if err != nil {
		// globallogger.Log.Warnln("CreateMsgCheckTimer CreateMsgCheckTimerInfo error: ", err)
		msgCheckTimerInfo, err := FindMsgCheckTimerAndUpdate(MsgKey, setData)
		return msgCheckTimerInfo, err
	}
	// globallogger.Log.Infoln("CreateMsgCheckTimer CreateMsgCheckTimerInfo success")
	return nil, err
}

// DeleteMsgCheckTimer DeleteMsgCheckTimer
func DeleteMsgCheckTimer(MsgKey string) {
	var err error
	if constant.Constant.UsePostgres {
		err = models.DeleteMsgCheckTimerPG(MsgKey)
	} else {
		err = models.DeleteMsgCheckTimer(MsgKey)
	}
	if err != nil {
		globallogger.Log.Errorln("DeleteMsgCheckTimer DeleteMsgCheckTimer error: ", err)
	} else {
		// globallogger.Log.Infoln("DeleteMsgCheckTimer DeleteMsgCheckTimer success")
	}
}

// MsgCheckTimerUpdateOrCreate MsgCheckTimerUpdateOrCreate
func MsgCheckTimerUpdateOrCreate(MsgKey string) (*config.MsgCheckTimerInfo, error) {
	var setData = config.MsgCheckTimerInfo{}
	setData.UpdateTime = time.Now()
	var err error
	var msgCheckTimerInfo *config.MsgCheckTimerInfo
	if constant.Constant.UsePostgres {
		oSet := make(map[string]interface{})
		oSet["updatetime"] = time.Now()
		msgCheckTimerInfo, err = FindMsgCheckTimerAndUpdatePG(MsgKey, oSet)
	} else {
		msgCheckTimerInfo, err = FindMsgCheckTimerAndUpdate(MsgKey, setData)
	}
	if msgCheckTimerInfo == nil {
		msgCheckTimerInfo, err = CreateMsgCheckTimer(MsgKey)
	}
	return msgCheckTimerInfo, err
}

/***************************************************************************************/
/* msgCheckTimerModel Func End */
/***************************************************************************************/

/***************************************************************************************/
/* reSendTimerModel Func Start */
/***************************************************************************************/

// GetReSendTimerByMsgKey GetReSendTimerByMsgKey
func GetReSendTimerByMsgKey(MsgKey string) (*config.ReSendTimerInfo, error) {
	if constant.Constant.UsePostgres {
		reSendTimerInfo, err := models.GetReSendTimerInfoByMsgKeyPG(MsgKey)
		if err != nil {
			globallogger.Log.Errorln("GetReSendTimerByMsgKey GetReSendTimerInfoByMsgKeyPG error: ", err)
		} else {
			// globallogger.Log.Infoln("GetReSendTimerByMsgKey GetReSendTimerInfoByMsgKeyPG success")
		}
		return reSendTimerInfo, err
	}
	reSendTimerInfo, err := models.GetReSendTimerInfoByMsgKey(MsgKey)
	if err != nil {
		globallogger.Log.Errorln("GetReSendTimerByMsgKey GetReSendTimerInfoByMsgKey error: ", err)
	} else {
		// globallogger.Log.Infoln("GetReSendTimerByMsgKey GetReSendTimerInfoByMsgKey success")
	}
	return reSendTimerInfo, err
}

// FindReSendTimerAndUpdate FindReSendTimerAndUpdate
func FindReSendTimerAndUpdate(MsgKey string, setData config.ReSendTimerInfo) (*config.ReSendTimerInfo, error) {
	reSendTimerInfo, err := models.FindReSendTimerInfoAndUpdate(MsgKey, setData)
	if err != nil {
		globallogger.Log.Errorln("FindReSendTimerAndUpdate FindReSendTimerInfoAndUpdate error: ", err)
	} else {
		// globallogger.Log.Infoln("FindReSendTimerAndUpdate FindReSendTimerInfoAndUpdate success")
	}
	return reSendTimerInfo, err
}

// FindReSendTimerAndUpdatePG FindReSendTimerAndUpdatePG
func FindReSendTimerAndUpdatePG(MsgKey string, oSet map[string]interface{}) (*config.ReSendTimerInfo, error) {
	reSendTimerInfo, err := models.FindReSendTimerInfoAndUpdatePG(MsgKey, oSet)
	if err != nil {
		globallogger.Log.Errorln("FindReSendTimerAndUpdatePG FindReSendTimerInfoAndUpdatePG error: ", err)
	} else {
		// globallogger.Log.Infoln("FindReSendTimerAndUpdatePG FindReSendTimerInfoAndUpdatePG success")
	}
	return reSendTimerInfo, err
}

// CreateReSendTimer CreateReSendTimer
func CreateReSendTimer(MsgKey string) error {
	var setData = config.ReSendTimerInfo{}
	setData.MsgKey = MsgKey
	setData.UpdateTime = time.Now()
	if constant.Constant.UsePostgres {
		err := models.CreateReSendTimerInfoPG(setData)
		if err != nil {
			globallogger.Log.Errorln("CreateReSendTimer CreateReSendTimerInfoPG error: ", err)
			oSet := make(map[string]interface{})
			oSet["msgkey"] = MsgKey
			oSet["updatetime"] = time.Now()
			_, err = FindReSendTimerAndUpdatePG(MsgKey, oSet)
		} else {
			// globallogger.Log.Infoln("CreateReSendTimer CreateReSendTimerInfoPG success")
		}
		return err
	}
	err := models.CreateReSendTimerInfo(setData)
	if err != nil {
		globallogger.Log.Errorln("CreateReSendTimer CreateReSendTimerInfo error: ", err)
		_, err = FindReSendTimerAndUpdate(MsgKey, setData)
	} else {
		// globallogger.Log.Infoln("CreateReSendTimer CreateReSendTimerInfo success")
	}
	return err
}

// DeleteReSendTimer DeleteReSendTimer
func DeleteReSendTimer(MsgKey string) {
	var err error
	if constant.Constant.UsePostgres {
		err = models.DeleteReSendTimerPG(MsgKey)
	} else {
		err = models.DeleteReSendTimer(MsgKey)
	}
	if err != nil {
		globallogger.Log.Errorln("DeleteReSendTimer DeleteReSendTimer error: ", err)
	} else {
		// globallogger.Log.Infoln("DeleteReSendTimer DeleteReSendTimer success")
	}
}

// ReSendTimerUpdateOrCreate ReSendTimerUpdateOrCreate
func ReSendTimerUpdateOrCreate(MsgKey string) time.Time {
	var setData = config.ReSendTimerInfo{}
	setData.UpdateTime = time.Now()

	go func() {
		select {
		case <-time.After(time.Duration(rand.Int()) * time.Millisecond):
			var reSendTimerInfo *config.ReSendTimerInfo
			if constant.Constant.UsePostgres {
				oSet := make(map[string]interface{})
				oSet["updatetime"] = time.Now()
				reSendTimerInfo, _ = FindReSendTimerAndUpdatePG(MsgKey, oSet)
			} else {
				reSendTimerInfo, _ = FindReSendTimerAndUpdate(MsgKey, setData)
			}
			if reSendTimerInfo != nil {
				go func() {
					select {
					case <-time.After(time.Duration(rand.Int()) * time.Millisecond):
						CreateReSendTimer(MsgKey)
					}
				}()
			}
		}
	}()

	return setData.UpdateTime
}

/*************************************************************************************/
/* reSendTimerModel Func End */
/*************************************************************************************/

/******************************************************************************************/
/* terminalTimerModel Func Start */
/******************************************************************************************/

// GetTerminalTimerByDevEUI GetTerminalTimerByDevEUI
func GetTerminalTimerByDevEUI(devEUI string) (*config.TerminalTimerInfo, error) {
	if constant.Constant.UsePostgres {
		terminalTimerInfo, err := models.GetTerminalTimerInfoByDevEUIPG(devEUI)
		if err != nil {
			globallogger.Log.Errorln("GetTerminalTimerByDevEUI GetTerminalTimerInfoByDevEUIPG error: ", err)
		} else {
			// globallogger.Log.Infoln("GetTerminalTimerByDevEUI GetTerminalTimerInfoByDevEUIPG success")
		}
		return terminalTimerInfo, err
	}
	terminalTimerInfo, err := models.GetTerminalTimerInfoByDevEUI(devEUI)
	if err != nil {
		globallogger.Log.Errorln("GetTerminalTimerByDevEUI GetTerminalTimerInfoByDevEUI error: ", err)
	} else {
		// globallogger.Log.Infoln("GetTerminalTimerByDevEUI GetTerminalTimerInfoByDevEUI success")
	}
	return terminalTimerInfo, err
}

// FindTerminalTimerAndUpdate FindTerminalTimerAndUpdate
func FindTerminalTimerAndUpdate(devEUI string, setData config.TerminalTimerInfo) (*config.TerminalTimerInfo, error) {
	terminalTimerInfo, err := models.FindTerminalTimerInfoAndUpdate(devEUI, setData)
	if err != nil {
		globallogger.Log.Errorln("FindTerminalTimerAndUpdate findTerminalTimerInfoAndUpdate error: ", err)
	} else {
		// globallogger.Log.Infoln("FindTerminalTimerAndUpdate findTerminalTimerInfoAndUpdate success")
	}
	return terminalTimerInfo, err
}

// FindTerminalTimerAndUpdatePG FindTerminalTimerAndUpdatePG
func FindTerminalTimerAndUpdatePG(devEUI string, oSet map[string]interface{}) (*config.TerminalTimerInfo, error) {
	terminalTimerInfo, err := models.FindTerminalTimerInfoAndUpdatePG(devEUI, oSet)
	if err != nil {
		globallogger.Log.Errorln("FindTerminalTimerAndUpdatePG FindTerminalTimerInfoAndUpdatePG error: ", err)
	} else {
		// globallogger.Log.Infoln("FindTerminalTimerAndUpdatePG FindTerminalTimerInfoAndUpdatePG success")
	}
	return terminalTimerInfo, err
}

// CreateTerminalTimer CreateTerminalTimer
func CreateTerminalTimer(devEUI string) error {
	var setData = config.TerminalTimerInfo{}
	setData.DevEUI = devEUI
	setData.UpdateTime = time.Now()
	if constant.Constant.UsePostgres {
		err := models.CreateTerminalTimerInfoPG(setData)
		if err != nil {
			globallogger.Log.Errorln("CreateTerminalTimer CreateTerminalTimerInfoPG error: ", err)
			oSet := make(map[string]interface{})
			oSet["deveui"] = devEUI
			oSet["updatetime"] = time.Now()
			_, err = FindTerminalTimerAndUpdatePG(devEUI, oSet)
		} else {
			// globallogger.Log.Infoln("CreateTerminalTimer CreateTerminalTimerInfoPG success")
		}
		return err
	}
	err := models.CreateTerminalTimerInfo(setData)
	if err != nil {
		globallogger.Log.Errorln("CreateTerminalTimer CreateTerminalTimerInfo error: ", err)
		_, err = FindTerminalTimerAndUpdate(devEUI, setData)
	} else {
		// globallogger.Log.Infoln("CreateTerminalTimer CreateTerminalTimerInfo success")
	}
	return err
}

// DeleteTerminalTimer DeleteTerminalTimer
func DeleteTerminalTimer(devEUI string) {
	var err error
	if constant.Constant.UsePostgres {
		err = models.DeleteTerminalTimerPG(devEUI)
	} else {
		err = models.DeleteTerminalTimer(devEUI)
	}
	if err != nil {
		globallogger.Log.Errorln("DeleteTerminalTimer DeleteTerminalTimer error: ", err)
	} else {
		// globallogger.Log.Infoln("DeleteTerminalTimer DeleteTerminalTimer success")
	}
}

// TerminalTimerUpdateOrCreate TerminalTimerUpdateOrCreate
func TerminalTimerUpdateOrCreate(devEUI string) error {
	var setData = config.TerminalTimerInfo{}
	setData.UpdateTime = time.Now()
	var err error
	var terminalTimerInfo *config.TerminalTimerInfo
	if constant.Constant.UsePostgres {
		oSet := make(map[string]interface{})
		oSet["updatetime"] = time.Now()
		terminalTimerInfo, err = FindTerminalTimerAndUpdatePG(devEUI, oSet)
	} else {
		terminalTimerInfo, err = FindTerminalTimerAndUpdate(devEUI, setData)
	}
	if terminalTimerInfo == nil {
		CreateTerminalTimer(devEUI)
	}
	return err
}

/****************************************************************************************/
/* terminalTimerModel Func End */
/****************************************************************************************/

// DeleteMsg DeleteMsg
func DeleteMsg(devEUI string, key string, t int) {
	zigbeeServerMsgCheckTimerID.Delete(key)
	msgCheckTimerInfo, _ := GetMsgCheckTimerByMsgKey(key)
	if msgCheckTimerInfo != nil {
		var dateNow = time.Now()
		var num = dateNow.UnixNano() - msgCheckTimerInfo.UpdateTime.UnixNano()
		if num > int64(time.Duration(t)*time.Second) {
			if constant.Constant.MultipleInstances {
				_, err := globalrediscache.RedisCache.DeleteRedis(key)
				if err != nil {
					globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" DeleteMsg delete redis error : ", err)
				} else {
					// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" DeleteMsg delete redis success: ", redisData)
				}
				_, err = globalrediscache.RedisCache.SremRedis("zigbee_down_msg", key)
				if err != nil {
					globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" DeleteMsg srem redis error : ", err)
				} else {
					// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" DeleteMsg srem redis success: ", redisData)
				}
			} else {
				_, err := globalmemorycache.MemoryCache.DeleteMemory(key)
				if err != nil {
					globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" DeleteMsg delete memory error : ", err)
				} else {
					// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" DeleteMsg delete memory success: ", memoryData)
				}
			}
			DeleteMsgCheckTimer(key)
		}
	}
}

// TimeoutCheckAndDeleteMsg TimeoutCheckAndDeleteMsg
func TimeoutCheckAndDeleteMsg(key string, devEUI string, timerID *time.Timer, isRebuild bool) {
	length, redisArray, _ := GetRedisLengthAndRangeRedis(key, devEUI)
	if len(redisArray) != 0 {
		var tempLen = 0
		for _, item := range redisArray {
			if item == constant.Constant.REDIS.ZigbeeRedisRemoveRedis {
				tempLen++
			}
		}
		tempTimerID, ok := zigbeeServerMsgCheckTimerID.Load(key)
		if (tempLen == length) && (ok && timerID == tempTimerID.(*time.Timer)) {
			DeleteMsg(devEUI, key, constant.Constant.TIMER.ZigbeeTimerMsgDelete)
		} else if (tempLen < length) && (ok && timerID == tempTimerID.(*time.Timer)) {
			if isRebuild {
				DeleteMsg(devEUI, key, constant.Constant.TIMER.ZigbeeTimerServerRebuild)
			} else {
				zigbeeServerMsgCheckTimerID.Delete(key)
				msgCheckTimerInfo, _ := GetMsgCheckTimerByMsgKey(key)
				if msgCheckTimerInfo != nil {
					var dateNow = time.Now()
					var num = dateNow.UnixNano() - msgCheckTimerInfo.UpdateTime.UnixNano()
					if num > int64(time.Duration(constant.Constant.TIMER.ZigbeeTimerMsgDelete)*time.Second) {
						timerID = time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerMsgDelete) * time.Second)
						zigbeeServerMsgCheckTimerID.Store(key, timerID)
						go func() {
							select {
							case <-timerID.C:
								TimeoutCheckAndDeleteMsg(key, devEUI, timerID, false)
							}
						}()
					}
				}
			}
		}
	}
}

// CheckTimeoutAndDeleteRedisMsg CheckTimeoutAndDeleteRedisMsg
func CheckTimeoutAndDeleteRedisMsg(key string, devEUI string, isRebuild bool) {
	if isRebuild {
		if value, ok := zigbeeServerMsgCheckTimerID.Load(key); ok {
			value.(*time.Timer).Stop()
			// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "key: " + key + " CheckTimeoutAndDeleteRedisMsg clearTimeout")
		}
		var timerID = time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerAfterTwo) * time.Second)
		zigbeeServerMsgCheckTimerID.Store(key, timerID)
		go func() {
			select {
			case <-timerID.C:
				TimeoutCheckAndDeleteMsg(key, devEUI, timerID, isRebuild)
			}
		}()
	} else {
		_, err := MsgCheckTimerUpdateOrCreate(key)
		if err == nil {
			if value, ok := zigbeeServerMsgCheckTimerID.Load(key); ok {
				value.(*time.Timer).Stop()
				// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "key: " + key + " CheckTimeoutAndDeleteRedisMsg clearTimeout")
			}
			var timerID = time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerMsgDelete) * time.Second)
			zigbeeServerMsgCheckTimerID.Store(key, timerID)
			go func() {
				select {
				case <-timerID.C:
					TimeoutCheckAndDeleteMsg(key, devEUI, timerID, isRebuild)
				}
			}()
		}
	}
}
