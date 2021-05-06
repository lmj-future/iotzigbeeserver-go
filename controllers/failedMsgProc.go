package controllers

import (
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalerrorcode"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmemorycache"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalrediscache"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
)

func deleteCtrlAndDataRedis(devEUI string, ctrlKey string, dataKey string) {
	if constant.Constant.MultipleInstances {
		_, err := globalrediscache.RedisCache.DeleteRedis(ctrlKey)
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "ctrlKey:", ctrlKey, "deleteCtrlAndDataRedis delete ctrl redis error :", err)
		}
		_, err = globalrediscache.RedisCache.DeleteRedis(dataKey)
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "dataKey:", dataKey, "deleteCtrlAndDataRedis delete data redis error :", err)
		}
	} else {
		_, err := globalmemorycache.MemoryCache.DeleteMemory(ctrlKey)
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "ctrlKey:", ctrlKey, "deleteCtrlAndDataRedis delete ctrl memory error :", err)
		}
		_, err = globalmemorycache.MemoryCache.DeleteMemory(dataKey)
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "dataKey:", dataKey, "deleteCtrlAndDataRedis delete data memory error :", err)
		}
	}
}

func procErrorFailed(devEUI string, lastMsgType string, key string, APMac string, SN string) {
	globallogger.Log.Warnln("devEUI :", devEUI, "operation failed, resend msg:", lastMsgType,
		globalmsgtype.MsgType.GetMsgTypeDOWNMsgMeaning(lastMsgType))
	redisData, _ := publicfunction.GetRedisDataFromRedis(key, devEUI, SN)
	if redisData != nil {
		publicfunction.CheckRedisAndReSend(devEUI, key, *redisData, APMac)
	}
}

func procTerminalOffline(devEUI string) {
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
	}
	if err == nil && terminalInfo != nil {
		if terminalInfo.Online {
			publicfunction.TerminalOffline(devEUI)
			if constant.Constant.Iotware {
				if !terminalInfo.LeaveState {
					iotsmartspace.StateTerminalOfflineIotware(*terminalInfo)
				} else {
					iotsmartspace.StateTerminalLeaveIotware(*terminalInfo)
				}
			} else if constant.Constant.Iotedge {
				if !terminalInfo.LeaveState {
					iotsmartspace.StateTerminalOffline(devEUI)
				} else {
					iotsmartspace.StateTerminalLeave(devEUI)
				}
			}
		}
	}
}

func procAPNotExist(devEUI string, ctrlKey string, dataKey string, APMac string) {
	globallogger.Log.Warnln("devEUI :", devEUI, "AP:", APMac, "is not exist. Please rebuild the environment !")
	procTerminalOffline(devEUI)
	deleteCtrlAndDataRedis(devEUI, ctrlKey, dataKey)
}

func procAPOffline(devEUI string, ctrlKey string, dataKey string, APMac string) {
	globallogger.Log.Warnln("devEUI :", devEUI, "AP:", APMac, "is offline. Please rebuild the environment !")
	procTerminalOffline(devEUI)
	deleteCtrlAndDataRedis(devEUI, ctrlKey, dataKey)
}

func procModuleNotExist(devEUI string, ctrlKey string, dataKey string, moduleID string) {
	globallogger.Log.Warnln("devEUI :", devEUI, "Module:", moduleID, "is not exist. Please check !")
	procTerminalOffline(devEUI)
	deleteCtrlAndDataRedis(devEUI, ctrlKey, dataKey)
}

func procModuleOffline(devEUI string, ctrlKey string, dataKey string, moduleID string) {
	globallogger.Log.Warnln("devEUI :", devEUI, "Module:", moduleID, "is offline. Please rebuild the environment !")
	procTerminalOffline(devEUI)
	deleteCtrlAndDataRedis(devEUI, ctrlKey, dataKey)
}

func procTerminalNotExist(devEUI string, ctrlKey string, dataKey string, APMac string) {
	globallogger.Log.Warnln("devEUI :", devEUI, "this terminal is not exist. Please let the terminal reJoin !")
	procTerminalOffline(devEUI)
	go publicfunction.IfCtrlMsgExistThenDeleteAll(ctrlKey, devEUI, APMac)
	go publicfunction.IfDataMsgExistThenDeleteAll(dataKey, devEUI, APMac)
}

func procChipNotMatch(devEUI string, ctrlKey string, dataKey string, moduleID string) {
	globallogger.Log.Warnln("devEUI :", devEUI, "Module:", moduleID, "chip is not match. Please check the chip type !")
	procTerminalOffline(devEUI)
	deleteCtrlAndDataRedis(devEUI, ctrlKey, dataKey)
}

func procTimeout(devEUI string, lastMsgType string, ctrlKey string, dataKey string, SN string, APMac string,
	jsonInfo publicstruct.JSONInfo, terminalIsExist bool) {
	if lastMsgType == globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalEndpointReqEvent {
		globallogger.Log.Warnln("devEUI :", devEUI, "timeout ! ZigbeeTerminalEndpointReqEvent !")
		procTerminalOffline(devEUI)
		go publicfunction.IfCtrlMsgExistThenDeleteAll(ctrlKey, devEUI, APMac)
		go publicfunction.IfDataMsgExistThenDeleteAll(dataKey, devEUI, APMac)
		if !terminalIsExist {
			go func() {
				timer := time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerAfterFive) * time.Second)
				<-timer.C
				publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
				timer.Stop()
			}()
		}
	} else if lastMsgType == globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalBindReqEvent {
		globallogger.Log.Warnln("devEUI :", devEUI, "timeout ! ZigbeeTerminalBindReqEvent !")
		go publicfunction.IfCtrlMsgExistThenDeleteAll(ctrlKey, devEUI, APMac)
		go publicfunction.IfDataMsgExistThenDeleteAll(dataKey, devEUI, APMac)
		go func() {
			timer := time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerAfterFive) * time.Second)
			<-timer.C
			publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
			timer.Stop()
		}()
	} else {
		globallogger.Log.Warnln("devEUI :", devEUI, "timeout ! drop this msg !")
		if terminalIsExist {
			if lastMsgType == globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent {
				publicfunction.IfMsgExistThenDelete(dataKey, devEUI, SN, APMac)
			} else {
				publicfunction.IfMsgExistThenDelete(ctrlKey, devEUI, SN, APMac)
			}
		} else {
			go publicfunction.IfCtrlMsgExistThenDeleteAll(ctrlKey, devEUI, APMac)
			go publicfunction.IfDataMsgExistThenDeleteAll(dataKey, devEUI, APMac)
		}
	}
}

func procFailedMsg(jsonInfo publicstruct.JSONInfo, devEUI string, terminalIsExist bool) {
	globallogger.Log.Warnln("devEUI :", jsonInfo.MessagePayload.Address, "procFailedMsg")
	var ctrlKey = publicfunction.GetCtrlMsgKey(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
	var dataKey = publicfunction.GetDataMsgKey(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)

	if strings.Repeat(jsonInfo.MessagePayload.Data[0:4], 1) == globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalCheckExistReqEvent {
		publicfunction.TerminalOffline(devEUI)
		var terminalInfo *config.TerminalInfo
		var err error
		if constant.Constant.UsePostgres {
			_, err = models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, map[string]interface{}{"leavestate": true})
			if err != nil {
				globallogger.Log.Errorln("devEUI :", devEUI, "procFailedMsg FindTerminalAndUpdatePG err :", err)
			}
			terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
		} else {
			_, err = models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, bson.M{"leaveState": true})
			if err != nil {
				globallogger.Log.Errorln("devEUI :", devEUI, "procFailedMsg FindTerminalAndUpdate err :", err)
			}
			terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
		}
		if err == nil && terminalInfo != nil {
			if constant.Constant.Iotware {
				iotsmartspace.StateTerminalLeaveIotware(*terminalInfo)
			} else if constant.Constant.Iotedge {
				iotsmartspace.StateTerminalLeave(devEUI)
			}
		}
		publicfunction.IfCtrlMsgExistThenDeleteAll(ctrlKey, devEUI, jsonInfo.TunnelHeader.LinkInfo.APMac)
		publicfunction.IfDataMsgExistThenDeleteAll(dataKey, devEUI, jsonInfo.TunnelHeader.LinkInfo.APMac)
	} else {
		switch strings.Repeat(jsonInfo.MessagePayload.Data[4:12], 1) {
		case globalerrorcode.ErrorCode.ZigbeeErrorFailed:
			if strings.Repeat(jsonInfo.MessagePayload.Data[0:4], 1) != globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent { //数据报文失败，不管
				procErrorFailed(devEUI, strings.Repeat(jsonInfo.MessagePayload.Data[0:4], 1), ctrlKey,
					jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.TunnelHeader.FrameSN)
			}
		case globalerrorcode.ErrorCode.ZigbeeErrorAPNotExist:
			procAPNotExist(devEUI, ctrlKey, dataKey, jsonInfo.TunnelHeader.LinkInfo.APMac)
		case globalerrorcode.ErrorCode.ZigbeeErrorAPOffline:
			procAPOffline(devEUI, ctrlKey, dataKey, jsonInfo.TunnelHeader.LinkInfo.APMac)
		case globalerrorcode.ErrorCode.ZigbeeErrorModuleNotExist:
			procModuleNotExist(devEUI, ctrlKey, dataKey, jsonInfo.MessagePayload.ModuleID)
		case globalerrorcode.ErrorCode.ZigbeeErrorModuleOffline:
			procModuleOffline(devEUI, ctrlKey, dataKey, jsonInfo.MessagePayload.ModuleID)
		case globalerrorcode.ErrorCode.ZigbeeErrorTerminalNotExist:
			procTerminalNotExist(devEUI, ctrlKey, dataKey, jsonInfo.TunnelHeader.LinkInfo.APMac)
		case globalerrorcode.ErrorCode.ZigbeeErrorChipNotMatch:
			procChipNotMatch(devEUI, ctrlKey, dataKey, jsonInfo.MessagePayload.ModuleID)
		case globalerrorcode.ErrorCode.ZigbeeErrorAPBusy:
			if strings.Repeat(jsonInfo.MessagePayload.Data[0:4], 1) != globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent {
				procTimeout(devEUI, strings.Repeat(jsonInfo.MessagePayload.Data[0:4], 1), ctrlKey, dataKey, jsonInfo.TunnelHeader.FrameSN,
					jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo, terminalIsExist)
			}
		case globalerrorcode.ErrorCode.ZigbeeErrorTimeout:
			procTimeout(devEUI, strings.Repeat(jsonInfo.MessagePayload.Data[0:4], 1), ctrlKey, dataKey, jsonInfo.TunnelHeader.FrameSN,
				jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo, terminalIsExist)
		default:
			globallogger.Log.Warnln("devEUI :", devEUI, "unknow errorCode:", strings.Repeat(jsonInfo.MessagePayload.Data[4:12], 1))
		}
	}
}
