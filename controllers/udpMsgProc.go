package controllers

import (
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coocood/freecache"
	"github.com/dyrkin/znp-go"
	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmemorycache"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalrediscache"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"

	// "github.com/h3c/iotzigbeeserver-go/metrics"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	zclmain "github.com/h3c/iotzigbeeserver-go/zcl"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/keepalive"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/frame"
	"github.com/lib/pq"
)

var zigbeeServerKeepAliveTimerID = &sync.Map{}
var zigbeeServerDuplicationCache = freecache.NewCache(10 * 1024 * 1024)

/*
server 接收到ACK消息
1.如果消息队列中有相应消息，说明ack比rsp先到，置ackFlag标记
2.如果消息队列中没有相应消息，说明rsp比ack先到，不做处理
*/
func procACKMsg(jsonInfo publicstruct.JSONInfo) {
	globallogger.Log.Infoln("devEUI :", jsonInfo.MessagePayload.Address, "procACKMsg receive msgType:", jsonInfo.MessagePayload.Data,
		globalmsgtype.MsgType.GetMsgTypeDOWNMsgMeaning(jsonInfo.MessagePayload.Data), "ACK success!")
	var key string
	var isCtrlMsg = true
	if jsonInfo.MessagePayload.Data == globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent {
		key = publicfunction.GetDataMsgKey(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
		isCtrlMsg = false
	} else {
		key = publicfunction.GetCtrlMsgKey(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
	}
	_, redisArray, _ := publicfunction.GetRedisLengthAndRangeRedis(key, jsonInfo.MessagePayload.Address)
	if len(redisArray) != 0 {
		var indexTemp = 0
		for itemIndex, item := range redisArray {
			if item != constant.Constant.REDIS.ZigbeeRedisRemoveRedis {
				var itemData publicstruct.RedisData
				json.Unmarshal([]byte(item), &itemData)
				if !itemData.AckFlag {
					if jsonInfo.TunnelHeader.FrameSN == itemData.SN {
						itemData.AckFlag = true
						itemDataSerial, _ := json.Marshal(itemData)
						if constant.Constant.MultipleInstances {
							_, err := globalrediscache.RedisCache.SetRedis(key, itemIndex, itemDataSerial)
							if err != nil {
								globallogger.Log.Errorln("devEUI :", jsonInfo.MessagePayload.Address, "key:", key,
									" procACKMsg set redis error :", err)
							}
						} else {
							_, err := globalmemorycache.MemoryCache.SetMemory(key, itemIndex, itemDataSerial)
							if err != nil {
								globallogger.Log.Errorln("devEUI :", jsonInfo.MessagePayload.Address, "key:", key,
									" procACKMsg set momery error :", err)
							}
						}
						indexTemp = itemIndex + 1
					}
					if !isCtrlMsg {
						if indexTemp > 0 && indexTemp == itemIndex {
							publicfunction.CheckRedisAndSend(jsonInfo.MessagePayload.Address, key,
								itemData.SN, jsonInfo.TunnelHeader.LinkInfo.APMac)
						}
					}
				}
			} else {
				indexTemp++
			}
		}
	}
}

func saveMsgDuplicationFlag(devEUI string, msgType string, SN string, isDataUpProc bool) (string, error) {
	if isDataUpProc || msgType == globalmsgtype.MsgType.UPMsg.ZigbeeDataUpEvent {
		// key = key + "_DataUpProc"
		return "toDo", nil
	}
	var keyBuilder strings.Builder
	keyBuilder.WriteString(constant.Constant.REDIS.ZigbeeRedisDuplicationFlagKey)
	keyBuilder.WriteString(devEUI)
	keyBuilder.WriteString("_")
	keyBuilder.WriteString(msgType)
	keyBuilder.WriteString("_")
	keyBuilder.WriteString(SN)
	var key = keyBuilder.String()
	var errorCode string
	var err error
	if constant.Constant.MultipleInstances {
		res, err := globalrediscache.RedisCache.GetRedis(key)
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "saveMsgDuplicationFlag get redis error :", err)
			errorCode = "error"
		} else {
			if res != "" {
				errorCode = "doNothing"
			} else {
				errorCode = "toDo"
				go func() {
					_, err = globalrediscache.RedisCache.UpdateRedis(key, []byte{1})
					if err != nil {
						globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "saveMsgDuplicationFlag update redis error :", err)
					} else {
						_, err := globalrediscache.RedisCache.SaddRedis(constant.Constant.REDIS.ZigbeeRedisDuplicationFlagSets, key)
						if err != nil {
							globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "saveMsgDuplicationFlag sadd redis error :", err)
						}
						timer := time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerDuplicationFlag) * time.Second)
						<-timer.C
						_, err = globalrediscache.RedisCache.DeleteRedis(key)
						if err != nil {
							globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "saveMsgDuplicationFlag delete redis error :", err)
						} else {
							_, err := globalrediscache.RedisCache.SremRedis(constant.Constant.REDIS.ZigbeeRedisDuplicationFlagSets, key)
							if err != nil {
								globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "saveMsgDuplicationFlag srem redis error :", err)
							}
						}
						timer.Stop()
					}
				}()
			}
		}
	} else {
		_, err := zigbeeServerDuplicationCache.Get([]byte(key))
		if err != nil {
			if err.Error() == "Entry not found" {
				errorCode = "toDo"
				if msgType == globalmsgtype.MsgType.UPMsg.ZigbeeDataUpEvent {
					zigbeeServerDuplicationCache.Set([]byte(key), []byte{1}, 1)
				} else {
					zigbeeServerDuplicationCache.Set([]byte(key), []byte{1}, constant.Constant.TIMER.ZigbeeTimerDuplicationFlag)
				}
			} else {
				globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "saveMsgDuplicationFlag get memory error :", err)
				errorCode = "error"
			}
		} else {
			errorCode = "doNothing"
		}
	}
	return errorCode, err
}

/*
server 接收到任何数据消息，
1.先回ACK；
2.再检查redis中有没有标记：
  如果有标记，表示server之前已经收到此消息，这次是client没有收到ACK的重传消息，无需再处理；
  如果没有标记，表示server第一次收到此消息，redis中存入标记，并起一个10秒定时器，超时后删除标记，继续处理相应消息；
3.消息处理过程中，如果需要向client发送消息请求，先入消息队列
4.再发送消息请求，存redis上下文
5.同时起一个1秒定时器，超时后检查消息队列中有没有该消息：
  如果有该消息，表示server没有收到ACK，消息重传（最多三次重传，三次重传失败后删除消息）；
  如果没有该消息，表示server收到了ACK，无需再处理；
*/
func procAnyMsg(jsonInfo publicstruct.JSONInfo) {
	globallogger.Log.Infoln("devEUI :", jsonInfo.MessagePayload.Address, "procAnyMsg receive msgType:", jsonInfo.MessageHeader.MsgType,
		globalmsgtype.MsgType.GetMsgTypeUPMsgMeaning(jsonInfo.MessageHeader.MsgType), "success!")
	go publicfunction.SendACK(jsonInfo)
	res, _ := saveMsgDuplicationFlag(jsonInfo.MessagePayload.Address+jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessageHeader.MsgType,
		jsonInfo.TunnelHeader.FrameSN, false)
	if res == "toDo" {
		if jsonInfo.MessagePayload.Topic == "80003002" {
			collihighProc(jsonInfo)
		} else {
			procMainMsg(jsonInfo)
		}
	}
}

func procKeepAliveMsgTimeout(APMac string) {
	terminalList, err := publicfunction.ProcKeepAliveTerminalOffline(APMac)
	if err != nil {
		globallogger.Log.Errorln("procKeepAliveMsg ProcKeepAliveTerminalOffline err:", err)
	} else {
		if terminalList != nil {
			if constant.Constant.Iotware {
				for _, item := range terminalList {
					if !item.LeaveState {
						iotsmartspace.StateTerminalOfflineIotware(item)
					} else {
						iotsmartspace.StateTerminalLeaveIotware(item)
					}
				}
			} else if constant.Constant.Iotedge {
				for _, item := range terminalList {
					if !item.LeaveState {
						iotsmartspace.StateTerminalOffline(item.DevEUI)
					} else {
						iotsmartspace.StateTerminalLeave(item.DevEUI)
					}
				}
			}
		}
	}
}
func procKeepAliveMsgRedis(APMac string, t int) {
	// if value, ok := zigbeeServerKeepAliveTimerID.Load(APMac); ok {
	// 	value.(*time.Timer).Stop()
	// }
	var timerID = time.NewTimer(time.Duration(t*3+3) * time.Second)
	zigbeeServerKeepAliveTimerID.Store(APMac, timerID)

	<-timerID.C
	timerID.Stop()
	if value, ok := zigbeeServerKeepAliveTimerID.Load(APMac); ok {
		if timerID == value.(*time.Timer) {
			zigbeeServerKeepAliveTimerID.Delete(APMac)
			updateTime, err := publicfunction.KeepAliveTimerRedisGet(APMac)
			if err == nil {
				if time.Now().UnixNano()-updateTime > int64(time.Duration(3*t)*time.Second) {
					globallogger.Log.Warnln("APMac:", APMac, "procKeepAliveMsgRedis server has not already recv keep alive msg for",
						(time.Now().UnixNano()-updateTime)/int64(time.Second), "seconds. Then offline all terminal")
					procKeepAliveMsgTimeout(APMac)
					go publicfunction.DeleteRedisKeepAliveTimer(APMac)
					go globalrediscache.RedisCache.DeleteRedis(constant.Constant.REDIS.ZigbeeRedisSNZHA + APMac)
					go globalrediscache.RedisCache.DeleteRedis(constant.Constant.REDIS.ZigbeeRedisSNTransfer + APMac)
				}
			}
		}
	}
}
func procKeepAliveMsgFreeCache(APMac string, t int) {
	// if value, ok := zigbeeServerKeepAliveTimerID.Load(APMac); ok {
	// 	value.(*time.Timer).Stop()
	// }
	var timerID = time.NewTimer(time.Duration(t*3+3) * time.Second)
	zigbeeServerKeepAliveTimerID.Store(APMac, timerID)

	<-timerID.C
	timerID.Stop()
	if value, ok := zigbeeServerKeepAliveTimerID.Load(APMac); ok {
		if timerID == value.(*time.Timer) {
			zigbeeServerKeepAliveTimerID.Delete(APMac)
			updateTime, err := publicfunction.KeepAliveTimerFreeCacheGet(APMac)
			if err == nil {
				if time.Now().UnixNano()-updateTime > int64(time.Duration(3*t)*time.Second) {
					globallogger.Log.Warnln("APMac:", APMac, "procKeepAliveMsgFreeCache server has not already recv keep alive msg for",
						(time.Now().UnixNano()-updateTime)/int64(time.Second), "seconds. Then offline all terminal")
					procKeepAliveMsgTimeout(APMac)
					go publicfunction.DeleteFreeCacheKeepAliveTimer(APMac)
					go globalmemorycache.MemoryCache.DeleteMemory(constant.Constant.REDIS.ZigbeeRedisSNZHA + APMac)
					go globalmemorycache.MemoryCache.DeleteMemory(constant.Constant.REDIS.ZigbeeRedisSNTransfer + APMac)
				}
			}
		}
	}
}

func procModuleRemoveEvent(jsonInfo publicstruct.JSONInfo) {
	var terminalList []config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalList, err = models.FindTerminalByAPMacAndModuleIDPG(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
	} else {
		terminalList, err = models.FindTerminalByAPMacAndModuleID(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
	}
	if err != nil {
		globallogger.Log.Errorln("procModuleRemoveEvent findTerminalByAPMacAndModuleID err:", err)
		return
	}
	if len(terminalList) > 0 {
		for _, item := range terminalList {
			publicfunction.TerminalOffline(item.DevEUI)
			if !item.LeaveState {
				if constant.Constant.Iotware {
					iotsmartspace.StateTerminalOfflineIotware(item)
				} else if constant.Constant.Iotedge {
					iotsmartspace.StateTerminalOffline(item.DevEUI)
				}
			} else {
				if constant.Constant.Iotware {
					iotsmartspace.StateTerminalLeaveIotware(item)
				} else if constant.Constant.Iotedge {
					iotsmartspace.StateTerminalLeave(item.DevEUI)
				}
			}
		}
	}
}

func procEnvironmentChangeEvent(jsonInfo publicstruct.JSONInfo) {
	var terminalList []config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalList, err = models.FindTerminalByAPMacAndModuleIDPG(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
	} else {
		terminalList, err = models.FindTerminalByAPMacAndModuleID(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
	}
	if err != nil {
		globallogger.Log.Errorln("procEnvironmentChangeEvent findTerminalByAPMacAndModuleID err:", err)
	} else {
		if len(terminalList) > 0 {
			for _, item := range terminalList {
				publicfunction.TerminalOffline(item.DevEUI)
				if constant.Constant.UsePostgres {
					oMatch := make(map[string]interface{}, 1)
					oSet := make(map[string]interface{}, 1)
					oMatch["deveui"] = item.DevEUI
					oSet["leavestate"] = true
					_, err := models.FindTerminalAndUpdatePG(oMatch, oSet)
					oMatch, oSet = nil, nil
					if err != nil {
						globallogger.Log.Errorln("devEUI :", item.DevEUI, "procEnvironmentChangeEvent FindTerminalAndUpdatePG err :", err)
					}
				} else {
					_, err := models.FindTerminalAndUpdate(bson.M{"devEUI": item.DevEUI}, bson.M{"leaveState": true})
					if err != nil {
						globallogger.Log.Errorln("devEUI :", item.DevEUI, "procEnvironmentChangeEvent FindTerminalAndUpdate err :", err)
					}
				}
				if constant.Constant.Iotware {
					iotsmartspace.StateTerminalLeaveIotware(item)
				} else if constant.Constant.Iotedge {
					iotsmartspace.StateTerminalLeave(item.DevEUI)
				}
			}
		}
	}
}

func findTerminalAndUpdateSame(jsonInfo publicstruct.JSONInfo) (*config.TerminalInfo, error) {
	var terminalInfo *config.TerminalInfo = nil
	var err error
	if constant.Constant.UsePostgres {
		var oMatchPG = map[string]interface{}{"deveui": jsonInfo.MessagePayload.Address}
		if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 {
			oMatchPG = map[string]interface{}{
				"nwkaddr":  jsonInfo.MessagePayload.Address,
				"apmac":    jsonInfo.TunnelHeader.LinkInfo.APMac,
				"moduleid": jsonInfo.MessagePayload.ModuleID,
			}
		}
		oSet := make(map[string]interface{}, 6)
		if jsonInfo.TunnelHeader.LinkInfo.ACMac != "" {
			oSet["acmac"] = jsonInfo.TunnelHeader.LinkInfo.ACMac
		}
		if jsonInfo.TunnelHeader.LinkInfo.T300ID != "" && jsonInfo.TunnelHeader.LinkInfo.T300ID != "000000000000" &&
			jsonInfo.TunnelHeader.LinkInfo.T300ID != "ffffffffffff" {
			oSet["t300id"] = jsonInfo.TunnelHeader.LinkInfo.T300ID
		}
		if jsonInfo.TunnelHeader.Version != "" {
			oSet["udpversion"] = jsonInfo.TunnelHeader.Version
		}
		oSet["apmac"] = jsonInfo.TunnelHeader.LinkInfo.APMac
		oSet["moduleid"] = jsonInfo.MessagePayload.ModuleID
		oSet["updatetime"] = time.Now()
		terminalInfo, err = models.FindTerminalAndUpdatePG(oMatchPG, oSet)
		oSet = nil
	} else {
		var oMatch = bson.M{"devEUI": jsonInfo.MessagePayload.Address}
		if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 {
			oMatch = bson.M{
				"nwkAddr":  jsonInfo.MessagePayload.Address,
				"APMac":    jsonInfo.TunnelHeader.LinkInfo.APMac,
				"moduleID": jsonInfo.MessagePayload.ModuleID,
			}
		}
		var setData = config.TerminalInfo{}
		if jsonInfo.TunnelHeader.LinkInfo.ACMac != "" {
			setData.ACMac = jsonInfo.TunnelHeader.LinkInfo.ACMac
		}
		if jsonInfo.TunnelHeader.LinkInfo.T300ID != "" && jsonInfo.TunnelHeader.LinkInfo.T300ID != "000000000000" &&
			jsonInfo.TunnelHeader.LinkInfo.T300ID != "ffffffffffff" {
			setData.T300ID = jsonInfo.TunnelHeader.LinkInfo.T300ID
		}
		if jsonInfo.TunnelHeader.Version != "" {
			setData.UDPVersion = jsonInfo.TunnelHeader.Version
		}
		setData.APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
		setData.ModuleID = jsonInfo.MessagePayload.ModuleID
		setData.UpdateTime = time.Now()
		terminalInfo, err = models.FindTerminalAndUpdate(oMatch, setData)
	}
	if err != nil {
		globallogger.Log.Errorln("devEUI :", jsonInfo.MessagePayload.Address, "findTerminalAndUpdate FindTerminalAndUpdate error :", err)
	} else {
		if terminalInfo == nil {
			globallogger.Log.Warnln("devEUI :", jsonInfo.MessagePayload.Address, "findTerminalAndUpdate terminal is not exist, please first create !")
		}
	}
	return terminalInfo, err
}

func findTerminalAndUpdate(jsonInfo publicstruct.JSONInfo) (*config.TerminalInfo, error) {
	var terminalInfo *config.TerminalInfo = nil
	var err error
	if jsonInfo.MessageHeader.MsgType == globalmsgtype.MsgType.UPMsg.ZigbeeDataUpEvent {
		var keyBuilder strings.Builder
		keyBuilder.WriteString(jsonInfo.TunnelHeader.LinkInfo.APMac)
		keyBuilder.WriteString(jsonInfo.MessagePayload.ModuleID)
		keyBuilder.WriteString(jsonInfo.MessagePayload.Address)
		var key = keyBuilder.String()
		var isNeedUpdate = false
		if v, ok := publicfunction.LoadTerminalInfoListCache(key); ok {
			terminalInfo = v.(*config.TerminalInfo)
			if terminalInfo.ACMac != jsonInfo.TunnelHeader.LinkInfo.ACMac || terminalInfo.APMac != jsonInfo.TunnelHeader.LinkInfo.APMac ||
				terminalInfo.T300ID != jsonInfo.TunnelHeader.LinkInfo.T300ID || terminalInfo.ModuleID != jsonInfo.MessagePayload.ModuleID ||
				terminalInfo.UDPVersion != jsonInfo.TunnelHeader.Version {
				isNeedUpdate = true
				globallogger.Log.Errorf("findTerminalAndUpdate: key: %s, terminalInfo: %+v, jsonInfo: %+v", key, terminalInfo, jsonInfo)
			}
		} else {
			isNeedUpdate = true
			globallogger.Log.Errorf("findTerminalAndUpdate: key: %s, LoadTerminalInfoListCache not ok", key)
		}
		if isNeedUpdate {
			terminalInfo, err = findTerminalAndUpdateSame(jsonInfo)
			if terminalInfo != nil && terminalInfo.Online && terminalInfo.IsExist && terminalInfo.TmnType != "invalidType" {
				publicfunction.StoreTerminalInfoListCache(key, terminalInfo)
			}
		}
	} else {
		terminalInfo, err = findTerminalAndUpdateSame(jsonInfo)
	}
	return terminalInfo, err
}

func procMainMsg(jsonInfo publicstruct.JSONInfo) {
	if jsonInfo.MessageHeader.MsgType == globalmsgtype.MsgType.UPMsg.ZigbeeTerminalJoinEvent {
		//设备入网
		procTerminalJoinEvent(jsonInfo)
		// metrics.CountUdpReceiveByLabel("ZigbeeTerminalJoinEvent")
	} else if jsonInfo.MessageHeader.MsgType == globalmsgtype.MsgType.UPMsg.ZigbeeWholeNetworkIEEERspEvent {
		//整个网络拓扑获取
	} else if jsonInfo.MessageHeader.MsgType == globalmsgtype.MsgType.UPMsg.ZigbeeNetworkEvent {
		//网络拓扑信息主动上报
	} else if jsonInfo.MessageHeader.MsgType == globalmsgtype.MsgType.UPMsg.ZigbeeModuleRemoveEvent {
		//插卡下线事件
		procModuleRemoveEvent(jsonInfo)
		// metrics.CountUdpReceiveByLabel("ZigbeeModuleRemoveEvent")
	} else if jsonInfo.MessageHeader.MsgType == globalmsgtype.MsgType.UPMsg.ZigbeeEnvironmentChangeEvent {
		//网络发生改变事件
		procEnvironmentChangeEvent(jsonInfo)
		// metrics.CountUdpReceiveByLabel("ZigbeeEnvironmentChangeEvent")
	} else if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 &&
		jsonInfo.MessageHeader.MsgType == globalmsgtype.MsgType.UPMsg.ZigbeeTerminalLeaveEvent {
		//设备离网
		var terminalInfo *config.TerminalInfo
		var err error
		if constant.Constant.UsePostgres {
			oSet := make(map[string]interface{}, 6)
			if jsonInfo.TunnelHeader.LinkInfo.ACMac != "" {
				oSet["acmac"] = jsonInfo.TunnelHeader.LinkInfo.ACMac
			}
			if jsonInfo.TunnelHeader.LinkInfo.T300ID != "" {
				oSet["t300id"] = jsonInfo.TunnelHeader.LinkInfo.T300ID
			}
			if jsonInfo.TunnelHeader.Version != "" {
				oSet["udpversion"] = jsonInfo.TunnelHeader.Version
			}
			oSet["apmac"] = jsonInfo.TunnelHeader.LinkInfo.APMac
			oSet["moduleid"] = jsonInfo.MessagePayload.ModuleID
			oSet["updatetime"] = time.Now()
			var oMatchPG = map[string]interface{}{"deveui": strings.ToUpper(jsonInfo.MessagePayload.Data)}
			terminalInfo, err = models.FindTerminalAndUpdatePG(oMatchPG, oSet)
			oSet = nil
		} else {
			var setData = config.TerminalInfo{}
			if jsonInfo.TunnelHeader.LinkInfo.ACMac != "" {
				setData.ACMac = jsonInfo.TunnelHeader.LinkInfo.ACMac
			}
			if jsonInfo.TunnelHeader.LinkInfo.T300ID != "" {
				setData.T300ID = jsonInfo.TunnelHeader.LinkInfo.T300ID
			}
			if jsonInfo.TunnelHeader.Version != "" {
				setData.UDPVersion = jsonInfo.TunnelHeader.Version
			}
			setData.APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
			setData.ModuleID = jsonInfo.MessagePayload.ModuleID
			setData.UpdateTime = time.Now()
			var oMatch = bson.M{"devEUI": strings.ToUpper(jsonInfo.MessagePayload.Data)}
			terminalInfo, err = models.FindTerminalAndUpdate(oMatch, setData)
		}
		if err == nil && terminalInfo != nil {
			procTerminalLeaveEvent(jsonInfo)
		}
		// metrics.CountUdpReceiveByLabel("ZigbeeTerminalLeaveEvent")
	} else {
		terminalInfo, err := findTerminalAndUpdate(jsonInfo)
		if err == nil {
			if terminalInfo != nil {
				switch jsonInfo.MessageHeader.MsgType {
				case globalmsgtype.MsgType.GENERALMsg.ZigbeeGeneralFailed,
					globalmsgtype.MsgType.GENERALMsgV0101.ZigbeeGeneralFailedV0101:
					//处理失败消息
					procFailedMsg(jsonInfo, terminalInfo.DevEUI, true)
					// metrics.CountUdpReceiveByLabel("ZigbeeGeneralFailed")
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalCheckExistRspEvent:
					//检查存在否回复
					procTerminalCheckExistEvent(jsonInfo)
					// metrics.CountUdpReceiveByLabel("ZigbeeTerminalCheckExistRspEvent")
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalLeaveEvent:
					//设备离网
					procTerminalLeaveEvent(jsonInfo)
					// metrics.CountUdpReceiveByLabel("ZigbeeTerminalLeaveEvent")
				case globalmsgtype.MsgType.UPMsg.ZigbeeDataUpEvent:
					//设备数据上报
					procDataUp(jsonInfo, *terminalInfo)
					if !terminalInfo.Online {
						publicfunction.TerminalOnline(terminalInfo.DevEUI, false)
						if constant.Constant.Iotware {
							if terminalInfo.IsExist {
								iotsmartspace.StateTerminalOnlineIotware(*terminalInfo)
							}
						} else if constant.Constant.Iotedge {
							iotsmartspace.StateTerminalOnline(terminalInfo.DevEUI)
						}
						if terminalInfo.Interval != 0 {
							keepalive.ProcKeepAlive(*terminalInfo, uint16(terminalInfo.Interval))
						}
					}
					// metrics.CountUdpReceiveByLabel("ZigbeeDataUpEvent")
				case globalmsgtype.MsgType.UPMsg.ZigbeeWholeNetworkNWKRspEvent:
					//单个网络拓扑获取
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalEndpointRspEvent:
					//设备端口号上报
					procTerminalEndpointInfoUp(jsonInfo, *terminalInfo)
					// metrics.CountUdpReceiveByLabel("ZigbeeTerminalEndpointRspEvent")
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalDiscoveryRspEvent:
					//设备服务发现回复
					procTerminalDiscoveryInfoUp(jsonInfo, *terminalInfo)
					// metrics.CountUdpReceiveByLabel("ZigbeeTerminalDiscoveryRspEvent")
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalBindRspEvent,
					globalmsgtype.MsgType.UPMsg.ZigbeeTerminalUnbindRspEvent:
					//绑定/解绑回复
					procBindOrUnbindTerminalRsp(jsonInfo, *terminalInfo)
					// metrics.CountUdpReceiveByLabel("ZigbeeTerminalBindRspEvent")
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalNetworkRspEvent:
					//设备邻居网络拓扑上报
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalLeaveRspEvent:
					//设备离网回复
					procTerminalLeaveRsp(jsonInfo, terminalInfo)
					// metrics.CountUdpReceiveByLabel("ZigbeeTerminalLeaveRspEvent")
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalPermitJoinRspEvent:
					//允许设备入网回复
					procTerminalPermitJoinRsp(jsonInfo, *terminalInfo)
					// metrics.CountUdpReceiveByLabel("ZigbeeTerminalPermitJoinRspEvent")
				default:
					globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "unknow msgType :", jsonInfo.MessageHeader.MsgType)
				}
				// metrics.CountUdpReceiveByDevSN(terminalInfo.DevEUI)
			}
		}
	}
}

func procTerminalCheckExistEvent(jsonInfo publicstruct.JSONInfo) {
	globallogger.Log.Warnln("devEUI :", jsonInfo.MessagePayload.Address, "procTerminalCheckExistEvent, data:", jsonInfo.MessagePayload.Data)
	var devEUI = jsonInfo.MessagePayload.Address
	// var ProfileID = strings.Repeat(jsonInfo.MessagePayload.Data[0:4], 1)
	// var NwkAddr = strings.Repeat(jsonInfo.MessagePayload.Data[20:24], 1)
	if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 {
		devEUI = strings.ToUpper(strings.Repeat(jsonInfo.MessagePayload.Data[4:20], 1))
	}
	if strings.Repeat(jsonInfo.MessagePayload.Data[24:32], 1) != "00000012" {
		publicfunction.TerminalOffline(devEUI)
		var terminalInfo *config.TerminalInfo = nil
		var err error
		if constant.Constant.UsePostgres {
			_, err = models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, map[string]interface{}{"leavestate": true})
			if err != nil {
				globallogger.Log.Errorln("devEUI :", devEUI, "procTerminalCheckExistEvent FindTerminalAndUpdatePG err :", err)
			}
			terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
		} else {
			_, err = models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, bson.M{"leaveState": true})
			if err != nil {
				globallogger.Log.Errorln("devEUI :", devEUI, "procTerminalCheckExistEvent FindTerminalAndUpdate err :", err)
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
		var key = publicfunction.GetCtrlMsgKey(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
		publicfunction.IfCtrlMsgExistThenDeleteAll(key, devEUI, jsonInfo.TunnelHeader.LinkInfo.APMac)
		key = publicfunction.GetDataMsgKey(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
		publicfunction.IfDataMsgExistThenDeleteAll(key, devEUI, jsonInfo.TunnelHeader.LinkInfo.APMac)
	} else {
		globallogger.Log.Warnln("devEUI :", devEUI, "procTerminalCheckExistEvent, terminal is exist")
	}
}

func procTerminalLeaveEvent(jsonInfo publicstruct.JSONInfo) {
	globallogger.Log.Warnln("devEUI :", jsonInfo.MessagePayload.Address, "procTerminalLeaveEvent")
	var devEUI = jsonInfo.MessagePayload.Address
	if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 {
		devEUI = strings.ToUpper(jsonInfo.MessagePayload.Data)
	}
	publicfunction.TerminalOffline(devEUI)
	var terminalInfo *config.TerminalInfo = nil
	var err error
	if constant.Constant.UsePostgres {
		_, err = models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, map[string]interface{}{"leavestate": true})
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "procTerminalLeaveEvent FindTerminalAndUpdatePG err :", err)
		}
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		_, err = models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, bson.M{"leaveState": true})
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "procTerminalLeaveEvent FindTerminalAndUpdate err :", err)
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
	var key = publicfunction.GetCtrlMsgKey(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
	publicfunction.IfCtrlMsgExistThenDeleteAll(key, devEUI, jsonInfo.TunnelHeader.LinkInfo.APMac)
	key = publicfunction.GetDataMsgKey(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
	publicfunction.IfDataMsgExistThenDeleteAll(key, devEUI, jsonInfo.TunnelHeader.LinkInfo.APMac)
}

func procJoinTerminalIsExist(terminalInfo config.TerminalInfo, setData config.TerminalInfo,
	oSet map[string]interface{}, jsonInfo publicstruct.JSONInfo, devEUI string) {
	if terminalInfo.APMac != setData.APMac {
		go publicfunction.SendTerminalLeaveReq(publicfunction.GetJSONInfo(terminalInfo), devEUI)
	}
	setData.Endpoint = terminalInfo.Endpoint
	setData.EndpointPG = terminalInfo.EndpointPG
	setData.EndpointCount = terminalInfo.EndpointCount
	setData.BindInfo = terminalInfo.BindInfo
	setData.BindInfoPG = terminalInfo.BindInfoPG
	setData.TmnType = "invalidType"
	globallogger.Log.Infof("devEUI : %+v procJoinTerminalIsExist setData: %+v", devEUI, setData)
	var err error
	if constant.Constant.UsePostgres {
		oSet["endpointpg"] = terminalInfo.EndpointPG
		oSet["endpointcount"] = terminalInfo.EndpointCount
		oSet["bindinfopg"] = terminalInfo.BindInfoPG
		oSet["tmntype"] = "invalidType"
		_, err = models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, oSet)
	} else {
		_, err = models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, setData)
	}
	if err != nil {
		globallogger.Log.Errorln("devEUI :", devEUI, "procJoinTerminalIsExist FindTerminalAndUpdate error :", err)
	} else {
		firstlyGetTerminalEndpoint(jsonInfo, devEUI)
	}
}

func procJoinTerminalNotExist(setData config.TerminalInfo, devEUI string, jsonInfo publicstruct.JSONInfo) {
	var err error
	if constant.Constant.UsePostgres {
		err = models.CreateTerminalPG(setData)
	} else {
		err = models.CreateTerminal(setData)
	}
	if err != nil {
		globallogger.Log.Errorln("devEUI :", devEUI, "procJoinTerminalNotExist createTerminal error :", err)
	} else {
		firstlyGetTerminalEndpoint(jsonInfo, devEUI)
	}
}

/*
CapabilityFlags:
CAPINFO_ALTPANCOORD       0x01    协调器
CAPINFO_DEVICETYPE_FFD    0x02
CAPINFO_DEVICETYPE_RFD    0x00
CAPINFO_POWER_AC          0x04    长供电
CAPINFO_RCVR_ON_IDLE      0x08
CAPINFO_SECURITY_CAPABLE  0x40
CAPINFO_ALLOC_ADDR        0x80

Type:
NWK_ASSOC_JOIN               0  第一次入网
NWK_ASSOC_REJOIN_UNSECURE    1  不加密重新入网
NWK_ASSOC_REJOIN_SECURE      2  加密重新入网
*/
func procTerminalJoinEvent(jsonInfo publicstruct.JSONInfo) {
	globallogger.Log.Warnln("devEUI :", jsonInfo.MessagePayload.Address, "procTerminalJoinEvent, data:", jsonInfo.MessagePayload.Data)
	var devEUI = jsonInfo.MessagePayload.Address
	var NwkAddr = strings.Repeat(jsonInfo.MessagePayload.Data[0:4], 1)           //设备网络地址
	var CapabilityFlags = strings.Repeat(jsonInfo.MessagePayload.Data[20:22], 1) //设备功能、性能标记
	var Type = strings.Repeat(jsonInfo.MessagePayload.Data[22:24], 1)            //设备入网方式
	devEUI = strings.ToUpper(strings.Repeat(jsonInfo.MessagePayload.Data[4:20], 1))
	res, _ := saveMsgDuplicationFlag(devEUI+jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessageHeader.MsgType, NwkAddr, false)
	if res != "toDo" {
		globallogger.Log.Warnln("devEUI :", devEUI, "procTerminalJoinEvent this join is rejoin and nwkAddr is not change, do nothing")
		return
	}
	// time.Sleep(2 * time.Second)
	// publicfunction.SendWholeNetworkIEEE(setData)
	time.Sleep(3 * time.Second)

	var oMatch = map[string]interface{}{}
	oMatch["devEUI"] = devEUI
	var tmnInfo *config.TerminalInfo
	var err error
	var setData = config.TerminalInfo{}
	oSet := make(map[string]interface{})
	setData.DevEUI = devEUI
	oSet["deveui"] = devEUI
	setData.Online = true
	oSet["online"] = true
	setData.NwkAddr = NwkAddr
	oSet["nwkaddr"] = NwkAddr
	setData.CapabilityFlags = CapabilityFlags
	oSet["capabilityflags"] = CapabilityFlags
	setData.JoinType = Type
	oSet["jointype"] = Type
	setData.IsNeedBind = true
	oSet["isneedbind"] = true
	setData.LeaveState = false
	oSet["leavestate"] = false
	if jsonInfo.TunnelHeader.LinkInfo.ACMac != "" {
		setData.ACMac = jsonInfo.TunnelHeader.LinkInfo.ACMac
		oSet["acmac"] = jsonInfo.TunnelHeader.LinkInfo.ACMac
	}
	setData.APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
	oSet["apmac"] = jsonInfo.TunnelHeader.LinkInfo.APMac
	if jsonInfo.TunnelHeader.LinkInfo.T300ID != "" {
		setData.T300ID = jsonInfo.TunnelHeader.LinkInfo.T300ID
		oSet["t300id"] = jsonInfo.TunnelHeader.LinkInfo.T300ID
	}
	if jsonInfo.TunnelHeader.Version != "" {
		setData.UDPVersion = jsonInfo.TunnelHeader.Version
		oSet["udpversion"] = jsonInfo.TunnelHeader.Version
	}
	if jsonInfo.TunnelHeader.LinkInfo.FirstAddr != "" {
		if len(jsonInfo.TunnelHeader.LinkInfo.FirstAddr) == 40 {
			setData.FirstAddr = publicfunction.Transport16StringToString(jsonInfo.TunnelHeader.LinkInfo.FirstAddr)
			oSet["firstaddr"] = publicfunction.Transport16StringToString(jsonInfo.TunnelHeader.LinkInfo.FirstAddr)
		} else {
			setData.FirstAddr = jsonInfo.TunnelHeader.LinkInfo.FirstAddr
			oSet["firstaddr"] = jsonInfo.TunnelHeader.LinkInfo.FirstAddr
		}
	}
	if jsonInfo.TunnelHeader.LinkInfo.SecondAddr != "" {
		if len(jsonInfo.TunnelHeader.LinkInfo.SecondAddr) == 40 {
			setData.SecondAddr = publicfunction.Transport16StringToString(jsonInfo.TunnelHeader.LinkInfo.SecondAddr)
			oSet["secondaddr"] = publicfunction.Transport16StringToString(jsonInfo.TunnelHeader.LinkInfo.SecondAddr)
		} else {
			setData.SecondAddr = jsonInfo.TunnelHeader.LinkInfo.SecondAddr
			oSet["secondaddr"] = jsonInfo.TunnelHeader.LinkInfo.SecondAddr
		}
	}
	if jsonInfo.TunnelHeader.LinkInfo.ThirdAddr != "" {
		setData.ThirdAddr = jsonInfo.TunnelHeader.LinkInfo.ThirdAddr
		oSet["thirdaddr"] = jsonInfo.TunnelHeader.LinkInfo.ThirdAddr
	}
	setData.ModuleID = jsonInfo.MessagePayload.ModuleID
	oSet["moduleid"] = jsonInfo.MessagePayload.ModuleID
	setData.UpdateTime = time.Now()
	oSet["updatetime"] = time.Now()
	if constant.Constant.UsePostgres {
		tmnInfo, err = models.FindTerminalByConditionPG(map[string]interface{}{"deveui": devEUI})
	} else {
		tmnInfo, err = models.FindTerminalByCondition(oMatch)
	}
	if err != nil {
		globallogger.Log.Errorln("devEUI :", devEUI, "procTerminalJoinEvent FindTerminalByCondition err :", err)
		publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
		return
	}
	if tmnInfo != nil && tmnInfo.NwkAddr != "" {
		if NwkAddr == tmnInfo.NwkAddr {
			globallogger.Log.Warnln("devEUI :", devEUI, "procTerminalJoinEvent this join is rejoin and nwkAddr is not change, do nothing")
			return
		}
		var key = tmnInfo.APMac + tmnInfo.ModuleID + devEUI
		if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 {
			key = tmnInfo.APMac + tmnInfo.ModuleID + tmnInfo.NwkAddr
		}
		publicfunction.DeleteTerminalInfoListCache(key)
	}
	if constant.Constant.Iotprivate {
	} else {
		if tmnInfo == nil || !tmnInfo.IsExist {
			if !constant.Constant.Iotware {
				var permitJoinContentList []string
				if constant.Constant.MultipleInstances {
					permitJoinContentList, err = globalrediscache.RedisCache.FindAllRedisKeys(constant.Constant.REDIS.ZigbeeRedisPermitJoinContentSets)
					if err != nil {
						globallogger.Log.Errorln("devEUI :", devEUI, "procTerminalJoinEvent FindAllRedisKeys err :", err)
						publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
						return
					}
					globallogger.Log.Infoln("devEUI :", devEUI, "procTerminalJoinEvent FindAllRedisKeys success :", permitJoinContentList)
					if len(permitJoinContentList) == 0 {
						globallogger.Log.Warnln("devEUI :", devEUI, "procTerminalJoinEvent permit join is not enable")
						publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
						return
					}
				} else {
					permitJoinContentList, err = globalmemorycache.MemoryCache.FindAllMemoryKeys(constant.Constant.REDIS.ZigbeeRedisPermitJoinContentSets)
					if err != nil {
						globallogger.Log.Errorln("devEUI :", devEUI, "procTerminalJoinEvent FindAllMemoryKeys err :", err)
						publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
						return
					}
					globallogger.Log.Infoln("devEUI :", devEUI, "procTerminalJoinEvent FindAllMemoryKeys success :", permitJoinContentList)
				}
				if len(permitJoinContentList) == 0 {
					globallogger.Log.Warnln("devEUI :", devEUI, "procTerminalJoinEvent permit join is not enable")
					publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
					return
				}
				isNeedLeave := true
				for _, v := range permitJoinContentList {
					if v != "" {
						isNeedLeave = false
					}
				}
				if isNeedLeave {
					globallogger.Log.Warnln("devEUI :", devEUI, "procTerminalJoinEvent permit join is not enable")
					publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
					return
				}
			}
		}
	}
	if tmnInfo != nil {
		isNeedBind := publicfunction.CheckTerminalIsNeedBind(devEUI, tmnInfo.TmnType)
		setData.JoinType = Type
		oSet["jointype"] = Type
		setData.IsNeedBind = isNeedBind
		oSet["isneedbind"] = isNeedBind
		procJoinTerminalIsExist(*tmnInfo, setData, oSet, jsonInfo, devEUI)
	} else {
		procJoinTerminalNotExist(setData, devEUI, jsonInfo)
	}
}

func checkNextRedis(devEUI string, APMac string, key string, index int, length int) {
	if index < length {
		var res string
		var err error
		if constant.Constant.MultipleInstances {
			res, err = globalrediscache.RedisCache.GetRedisByIndex(key, index)
		} else {
			res, err = globalmemorycache.MemoryCache.GetMemoryByIndex(key, index)
		}
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "checkNextRedis get redis by index error :", err)
		} else {
			if res != "" {
				if res == constant.Constant.REDIS.ZigbeeRedisRemoveRedis {
					index++
					checkNextRedis(devEUI, APMac, key, index, length)
				} else {
					var redisData struct {
						SN string
					}
					json.Unmarshal([]byte(res), &redisData)
					publicfunction.CheckRedisAndSend(devEUI, key, redisData.SN, APMac)
				}
			} else {
				index++
				checkNextRedis(devEUI, APMac, key, index, length)
			}
		}
	}
}

func procNeedCheckData(seqNum int64, devEUI string, APMac string, moduleID string,
	tmnInfo config.TerminalInfo, afIncomingMessage znp.AfIncomingMessage, zclAfIncomingMessage *zcl.IncomingMessage) {
	var SNBuilder strings.Builder
	SNBuilder.WriteString("0000")
	SNBuilder.WriteString(strconv.FormatInt(seqNum, 16))
	var key = publicfunction.GetDataMsgKey(APMac, moduleID)
	length, index, redisData, _ := publicfunction.GetRedisDataFromDataRedis(key, devEUI, strings.Repeat(SNBuilder.String()[SNBuilder.Len()-4:], 1))
	if redisData != nil {
		var err error
		if constant.Constant.MultipleInstances {
			_, err = globalrediscache.RedisCache.SetRedis(key, index, []byte(constant.Constant.REDIS.ZigbeeRedisRemoveRedis))
		} else {
			_, err = globalmemorycache.MemoryCache.SetMemory(key, index, []byte(constant.Constant.REDIS.ZigbeeRedisRemoveRedis))
		}
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "procNeedCheckData set redis error :", err)
		} else {
			if !redisData.AckFlag {
				index++
				go checkNextRedis(devEUI, APMac, key, index, length)
			}
			zclmain.ZclMain(globalmsgtype.MsgType.UPMsg.ZigbeeDataUpEvent, tmnInfo, afIncomingMessage, redisData.MsgID, redisData.Data, zclAfIncomingMessage)
		}
	} else {
		zclmain.ZclMain(globalmsgtype.MsgType.UPMsg.ZigbeeDataUpEvent, tmnInfo, afIncomingMessage, "", "", zclAfIncomingMessage)
	}
}

func procBasicReadRsp(devEUI string, jsonInfo publicstruct.JSONInfo, terminalInfo config.TerminalInfo, command interface{}) {
	type BasicInfo struct {
		ZLibraryVersion     uint8
		ApplicationVersion  uint8
		StackVersion        uint8
		HWVersion           uint8
		ManufacturerName    string
		ModelIdentifier     string
		DateCode            string
		PowerSource         uint8
		LocationDescription string
		PhysicalEnvironment uint8
		DeviceEnabled       bool
		AlarmMask           interface{}
		DisableLocalConfig  interface{}
		SWBuildID           string
	}
	basicInfo := BasicInfo{}
	for _, value := range command.(*cluster.ReadAttributesResponse).ReadAttributeStatuses {
		switch value.AttributeName {
		case "ZLibraryVersion":
			if value.Status == cluster.ZclStatusSuccess {
				basicInfo.ZLibraryVersion = uint8(value.Attribute.Value.(uint64))
			}
		case "ApplicationVersion":
			if value.Status == cluster.ZclStatusSuccess {
				basicInfo.ApplicationVersion = uint8(value.Attribute.Value.(uint64))
			}
		case "StackVersion":
			if value.Status == cluster.ZclStatusSuccess {
				basicInfo.StackVersion = uint8(value.Attribute.Value.(uint64))
			}
		case "HWVersion":
			if value.Status == cluster.ZclStatusSuccess {
				basicInfo.HWVersion = uint8(value.Attribute.Value.(uint64))
			}
		case "ManufacturerName":
			if value.Status == cluster.ZclStatusSuccess {
				basicInfo.ManufacturerName = value.Attribute.Value.(string)
			}
		case "ModelIdentifier":
			if value.Status == cluster.ZclStatusSuccess {
				basicInfo.ModelIdentifier = value.Attribute.Value.(string)
			}
		case "DateCode":
			if value.Status == cluster.ZclStatusSuccess {
				basicInfo.DateCode = value.Attribute.Value.(string)
			}
		case "PowerSource":
			if value.Status == cluster.ZclStatusSuccess {
				basicInfo.PowerSource = uint8(value.Attribute.Value.(uint64))
			}
		case "LocationDescription":
			if value.Status == cluster.ZclStatusSuccess {
				basicInfo.LocationDescription = value.Attribute.Value.(string)
			}
		case "PhysicalEnvironment":
			if value.Status == cluster.ZclStatusSuccess {
				basicInfo.PhysicalEnvironment = uint8(value.Attribute.Value.(uint64))
			}
		case "DeviceEnabled":
			if value.Status == cluster.ZclStatusSuccess {
				basicInfo.DeviceEnabled = value.Attribute.Value.(bool)
			}
		case "AlarmMask":
			if value.Status == cluster.ZclStatusSuccess {
				basicInfo.AlarmMask = value.Attribute.Value
			}
		case "DisableLocalConfig":
			if value.Status == cluster.ZclStatusSuccess {
				basicInfo.DisableLocalConfig = value.Attribute.Value
			}
		case "SWBuildID":
			if value.Status == cluster.ZclStatusSuccess {
				basicInfo.SWBuildID = value.Attribute.Value.(string)
			}
		default:
			globallogger.Log.Warnln("devEUI :", devEUI, "procBasicReadRsp unknow attributeName", value.AttributeName)
		}
	}
	globallogger.Log.Infof("[devEUI: %v][procBasicReadRsp] basicInfo: %+v", devEUI, basicInfo)
	if basicInfo.ModelIdentifier == "FTB56+PTH01MK2.2" {
		basicInfo.ModelIdentifier = constant.Constant.TMNTYPE.MAILEKE.ZigbeeTerminalPMT1004Detector
	}
	if basicInfo.ModelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106 {
		basicInfo.ModelIdentifier = constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55
	}
	if basicInfo.ManufacturerName == "Feibit,co.ltd   " {
		basicInfo.ManufacturerName = "DINAIKE"
	}
	if basicInfo.ModelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c {
		basicInfo.ManufacturerName = constant.Constant.MANUFACTURERNAME.Honyar
	}
	if basicInfo.ManufacturerName == "REXENSE" {
		basicInfo.ManufacturerName = constant.Constant.MANUFACTURERNAME.Honyar
	}
	if terminalInfo.TmnType != "" && terminalInfo.TmnType != "invalidType" {
		globallogger.Log.Warnln("devEUI :", devEUI, "procBasicReadRsp tmnType is not nil, this readBasic for check terminal state")
		return
	}
	if basicInfo.ModelIdentifier != "" {
		var info = publicstruct.TmnTypeAttr{}
		info = publicfunction.GetTmnTypeAndAttribute(basicInfo.ModelIdentifier)
		oSet := make(map[string]interface{})
		var setData = config.TerminalInfo{}
		if constant.Constant.UsePostgres {
			oSet["tmnname"] = terminalInfo.SecondAddr + "-" + devEUI[6:]
			oSet["interval"] = publicfunction.GetTerminalInterval(basicInfo.ModelIdentifier)
			oSet["tmntype"] = info.TmnType
			oSet["manufacturername"] = strings.ToUpper(basicInfo.ManufacturerName)
			oSet["profileid"] = "ZHA"
			oSet["statuspg"] = pq.StringArray(info.Attribute.Status)
			oSet["attributedata"] = info.Attribute.Data
			oSet["latestcommand"] = info.Attribute.LatestCommand
			oSet["isneedbind"] = publicfunction.CheckTerminalIsNeedBind(devEUI, info.TmnType)
			oSet["isreadbasic"] = true
			oSet["updatetime"] = time.Now()
			oSet["endpointpg"] = terminalInfo.EndpointPG
			if basicInfo.ModelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c ||
				basicInfo.ModelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55 ||
				basicInfo.ModelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105 ||
				basicInfo.ModelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106 ||
				basicInfo.ModelIdentifier == constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM {
				oSet["endpointpg"] = pq.StringArray{"01"}
				oSet["endpointcount"] = 1
			}
		} else {
			setData.TmnName = terminalInfo.SecondAddr + "-" + devEUI[6:]
			setData.Interval = publicfunction.GetTerminalInterval(basicInfo.ModelIdentifier)
			setData.TmnType = info.TmnType
			setData.ManufacturerName = strings.ToUpper(basicInfo.ManufacturerName)
			setData.ProfileID = "ZHA"
			setData.Attribute = info.Attribute
			setData.IsNeedBind = publicfunction.CheckTerminalIsNeedBind(devEUI, info.TmnType)
			setData.IsReadBasic = true
			setData.UpdateTime = time.Now()
			setData.Endpoint = terminalInfo.Endpoint
			if basicInfo.ModelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c ||
				basicInfo.ModelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55 ||
				basicInfo.ModelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105 ||
				basicInfo.ModelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106 ||
				basicInfo.ModelIdentifier == constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM {
				setData.Endpoint = []string{"01"}
				setData.EndpointCount = 1
			}
		}
		if !constant.Constant.Iotware {
			if constant.Constant.Iotprivate {
				thirdlyDiscoveryByEndpoint(devEUI, setData, terminalInfo, jsonInfo)
			} else {
				if terminalInfo.IsExist {
					if constant.Constant.UsePostgres {
						thirdlyDiscoveryByEndpointPG(devEUI, oSet, terminalInfo, jsonInfo)
					} else {
						thirdlyDiscoveryByEndpoint(devEUI, setData, terminalInfo, jsonInfo)
					}
					return
				}
				var redisDataStr string
				var err error
				if constant.Constant.MultipleInstances {
					redisDataStr, err = globalrediscache.RedisCache.GetRedis(constant.Constant.REDIS.ZigbeeRedisPermitJoinContentKey + basicInfo.ModelIdentifier)
				} else {
					redisDataStr, err = globalmemorycache.MemoryCache.GetMemory(constant.Constant.REDIS.ZigbeeRedisPermitJoinContentKey + basicInfo.ModelIdentifier)
				}
				if err != nil {
					globallogger.Log.Errorln("devEUI :", devEUI, "procBasicReadRsp GetRedis err :", err)
					publicfunction.SendTerminalLeaveReq(jsonInfo, terminalInfo.DevEUI)
					return
				}
				if redisDataStr == "" {
					globallogger.Log.Warnln("devEUI :", devEUI, "procBasicReadRsp GetRedis redisDataStr null :", redisDataStr)
					publicfunction.SendTerminalLeaveReq(jsonInfo, terminalInfo.DevEUI)
					return
				}
				redisData := make(map[string]interface{}, 4)
				err = json.Unmarshal([]byte(redisDataStr), &redisData)
				if err != nil {
					globallogger.Log.Errorln("devEUI :", devEUI, "procBasicReadRsp Unmarshal err :", err)
					publicfunction.SendTerminalLeaveReq(jsonInfo, terminalInfo.DevEUI)
					return
				}
				tmnType := basicInfo.ModelIdentifier
				if redisData["tmnType"] != nil && redisData["tmnType"] != "" {
					tmnType = redisData["tmnType"].(string)
				}
				if tmnType == basicInfo.ModelIdentifier {
					if constant.Constant.UsePostgres {
						thirdlyDiscoveryByEndpointPG(devEUI, oSet, terminalInfo, jsonInfo)
					} else {
						thirdlyDiscoveryByEndpoint(devEUI, setData, terminalInfo, jsonInfo)
					}
				} else {
					globallogger.Log.Warnln("devEUI :", devEUI, "procBasicReadRsp tmnType is not match, real tmnType:",
						basicInfo.ModelIdentifier, "permit join tmnType:", tmnType)
					publicfunction.SendTerminalLeaveReq(jsonInfo, terminalInfo.DevEUI)
				}
			}
		} else {
			if constant.Constant.UsePostgres {
				thirdlyDiscoveryByEndpointPG(devEUI, oSet, terminalInfo, jsonInfo)
			} else {
				thirdlyDiscoveryByEndpoint(devEUI, setData, terminalInfo, jsonInfo)
			}
		}
	} else {
		globallogger.Log.Warnln("devEUI :", devEUI, "procBasicReadRsp permit join is not enable")
		publicfunction.SendTerminalLeaveReq(jsonInfo, terminalInfo.DevEUI)
	}
}

func procDataUp(jsonInfo publicstruct.JSONInfo, terminalInfo config.TerminalInfo) {
	// var msgType = jsonInfo.MessageHeader.MsgType
	ProfileID, _ := strconv.ParseUint(strings.Repeat(jsonInfo.MessagePayload.Data[0:4], 1), 16, 16)
	GroupID, _ := strconv.ParseUint(strings.Repeat(jsonInfo.MessagePayload.Data[4:8], 1), 16, 16)
	ClusterID, _ := strconv.ParseUint(strings.Repeat(jsonInfo.MessagePayload.Data[8:12], 1), 16, 16)
	SrcAddr := strings.Repeat(jsonInfo.MessagePayload.Data[12:16], 1)
	SrcEndpoint, _ := strconv.ParseUint(strings.Repeat(jsonInfo.MessagePayload.Data[16:18], 1), 16, 8)
	DstEndpoint, _ := strconv.ParseUint(strings.Repeat(jsonInfo.MessagePayload.Data[18:20], 1), 16, 8)
	WasBroadcast, _ := strconv.ParseUint(strings.Repeat(jsonInfo.MessagePayload.Data[20:22], 1), 16, 8)
	LinkQuality, _ := strconv.ParseUint(strings.Repeat(jsonInfo.MessagePayload.Data[22:24], 1), 16, 8)
	SecurityUse, _ := strconv.ParseUint(strings.Repeat(jsonInfo.MessagePayload.Data[24:26], 1), 16, 8)
	Timestamp, _ := strconv.ParseUint(strings.Repeat(jsonInfo.MessagePayload.Data[26:34], 1), 16, 32)
	TransSeqNumber, _ := strconv.ParseUint(strings.Repeat(jsonInfo.MessagePayload.Data[34:36], 1), 16, 8)
	// Len, _ := strconv.ParseUint(strings.Repeat(jsonInfo.MessagePayload.Data[36:38], 1), 16, 8)
	Data, _ := hex.DecodeString(strings.Repeat(jsonInfo.MessagePayload.Data[38:], 1))
	globallogger.Log.Infoln("devEUI :", terminalInfo.DevEUI, "procDataUp, data:", strings.Repeat(jsonInfo.MessagePayload.Data[38:], 1))
	if !terminalInfo.IsDiscovered && ClusterID != 0x0000 {
		return
	}

	afIncomingMessage := znp.AfIncomingMessage{
		GroupID:        uint16(GroupID),
		ClusterID:      uint16(ClusterID),
		SrcAddr:        SrcAddr,
		SrcEndpoint:    uint8(SrcEndpoint),
		DstEndpoint:    uint8(DstEndpoint),
		WasBroadcast:   uint8(WasBroadcast),
		LinkQuality:    uint8(LinkQuality),
		SecurityUse:    uint8(SecurityUse),
		Timestamp:      uint32(Timestamp),
		TransSeqNumber: uint8(TransSeqNumber),
		Data:           []uint8(Data),
	}
	switch ProfileID {
	case 0x0104:
		zclAfIncomingMessage, err := zcl.ZCLObj().ToZclIncomingMessage(&afIncomingMessage)
		if err != nil {
			globallogger.Log.Errorln("devEUI :", terminalInfo.DevEUI, "procDataUp ToZclIncomingMessage err :", err)
			return
		}
		if !zclAfIncomingMessage.Data.FrameControl.DisableDefaultResponse {
			go func() {
				frameHexString, transactionID := zcl.ZCLObj().EncFrameConfigurationToHexString(frame.Configuration{
					FrameType:                        frame.FrameTypeGlobal,
					FrameTypeConfigured:              true,
					Direction:                        frame.DirectionClientServer,
					DirectionConfigured:              true,
					DisableDefaultResponse:           true,
					DisableDefaultResponseConfigured: true,
					CommandID:                        0x0b,
					CommandIDConfigured:              true,
					Command: cluster.DefaultResponseCommand{
						CommandID: zclAfIncomingMessage.Data.CommandIdentifier,
						Status:    cluster.ZclStatusSuccess,
					},
					CommandConfigured: true,
				}, zclAfIncomingMessage.Data.TransactionSequenceNumber)
				publicfunction.SendZHADefaultResponseDownMsg(
					common.DownMsg{
						ProfileID: uint16(ProfileID),
						ClusterID: uint16(ClusterID),
						ZclData:   frameHexString,
						SN:        transactionID,
						MsgType:   globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
						DevEUI:    terminalInfo.DevEUI,
					}, &terminalInfo, strings.Repeat(jsonInfo.MessagePayload.Data[16:18], 1))
			}()
		}
		// res, _ := saveMsgDuplicationFlag(terminalInfo.DevEUI+jsonInfo.TunnelHeader.LinkInfo.APMac, msgType+commandName, strconv.Itoa(int(seqNum)), true)
		// if res == "toDo" {
		if zclAfIncomingMessage.Data.CommandName != "ReportAttributes" && zclAfIncomingMessage.Data.CommandName != "ZoneEnrollRequest" {
			go procNeedCheckData(int64(zclAfIncomingMessage.Data.TransactionSequenceNumber), terminalInfo.DevEUI, jsonInfo.TunnelHeader.LinkInfo.APMac,
				jsonInfo.MessagePayload.ModuleID, terminalInfo, afIncomingMessage, zclAfIncomingMessage)
			if zclAfIncomingMessage.Data.CommandName == "ReadAttributesResponse" && ClusterID == 0x0000 {
				procBasicReadRsp(terminalInfo.DevEUI, jsonInfo, terminalInfo, zclAfIncomingMessage.Data.Command)
			}
		} else {
			zclmain.ZclMain(globalmsgtype.MsgType.UPMsg.ZigbeeDataUpEvent, terminalInfo, afIncomingMessage, "", "", zclAfIncomingMessage)
		}
		// } else {
		// 	globallogger.Log.Infoln("devEUI :", terminalInfo.DevEUI ,"", "procDataUp msg is duplication, drop this msg")
		// }
	case 0xc05e:
	default:
		//to do
	}
}

func procTerminalEndpointInfoUp(jsonInfo publicstruct.JSONInfo, terminalInfo config.TerminalInfo) {
	globallogger.Log.Infoln("devEUI :", terminalInfo.DevEUI, "procTerminalEndpointInfoUp, data:", jsonInfo.MessagePayload.Data)
	//var profileID = strings.Repeat(data[0:4], 1)
	//var SrcAddr = strings.Repeat(data[4:8], 1)         //设备网络地址
	//var NWKAddr = strings.Repeat(data[10:14], 1)       //设备短地址

	var key = publicfunction.GetCtrlMsgKey(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
	if strings.Repeat(jsonInfo.MessagePayload.Data[8:10], 1) == "00" {
		go publicfunction.IfMsgExistThenDelete(key, terminalInfo.DevEUI, jsonInfo.TunnelHeader.FrameSN, jsonInfo.TunnelHeader.LinkInfo.APMac)
		secondlyReadBasicByEndpoint(terminalInfo.DevEUI, strings.Repeat(jsonInfo.MessagePayload.Data[14:16], 1),
			strings.Repeat(jsonInfo.MessagePayload.Data[16:], 1))
	} else {
		globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "procTerminalEndpointInfoUp failed ")
		redisData, _ := publicfunction.GetRedisDataFromRedis(key, terminalInfo.DevEUI, jsonInfo.TunnelHeader.FrameSN)
		if redisData != nil {
			publicfunction.CheckRedisAndReSend(terminalInfo.DevEUI, key, *redisData, jsonInfo.TunnelHeader.LinkInfo.APMac)
		}
	}
}

func procTerminalDiscoveryInfoUp(jsonInfo publicstruct.JSONInfo, terminalInfo config.TerminalInfo) {
	globallogger.Log.Infoln("devEUI :", terminalInfo.DevEUI, "procTerminalDiscoveryInfoUp, data:", jsonInfo.MessagePayload.Data)
	//var profileID = strings.Repeat(jsonInfo.MessagePayload.Data[0:4], 1)
	//var SrcAddr = strings.Repeat(jsonInfo.MessagePayload.Data[4:8], 1)         //发送消息设备的网络地址
	//var NWKAddr = strings.Repeat(jsonInfo.MessagePayload.Data[10:14], 1)       //响应消息设备的网络地址
	//var Len = strings.Repeat(jsonInfo.MessagePayload.Data[14:16], 1)           //消息长度
	//var DevVer = strings.Repeat(jsonInfo.MessagePayload.Data[26:28], 1)        //设备版本
	var NumInClusters = strings.Repeat(jsonInfo.MessagePayload.Data[28:30], 1) //InCluster数量
	numInClusters, _ := strconv.ParseInt(NumInClusters, 16, 0)
	var NumOutClusters = strings.Repeat(jsonInfo.MessagePayload.Data[30+4*numInClusters:32+4*numInClusters], 1) //OutCluster数量
	numOutClusters, _ := strconv.ParseInt(NumOutClusters, 16, 0)

	var key = publicfunction.GetCtrlMsgKey(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
	if strings.Repeat(jsonInfo.MessagePayload.Data[8:10], 1) == "00" {
		_, err := lastlyUpdateBindInfo(terminalInfo.DevEUI, terminalInfo, strings.Repeat(jsonInfo.MessagePayload.Data[18:22], 1),
			strings.Repeat(jsonInfo.MessagePayload.Data[22:26], 1), int(numInClusters),
			strings.Repeat(jsonInfo.MessagePayload.Data[30:30+4*numInClusters], 1), int(numOutClusters),
			strings.Repeat(jsonInfo.MessagePayload.Data[32+4*numInClusters:32+4*numInClusters+4*numOutClusters], 1),
			strings.Repeat(jsonInfo.MessagePayload.Data[16:18], 1), jsonInfo)
		if err == nil {
			publicfunction.IfMsgExistThenDelete(key, terminalInfo.DevEUI, jsonInfo.TunnelHeader.FrameSN, jsonInfo.TunnelHeader.LinkInfo.APMac)
		}
	} else {
		globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "procTerminalDiscoveryInfoUp failed ")
		redisData, _ := publicfunction.GetRedisDataFromRedis(key, terminalInfo.DevEUI, jsonInfo.TunnelHeader.FrameSN)
		if redisData != nil {
			publicfunction.CheckRedisAndReSend(terminalInfo.DevEUI, key, *redisData, jsonInfo.TunnelHeader.LinkInfo.APMac)
		}
	}
}

func procBindOrUnbindResult(devEUI string, msgType string, SrcEndpoint string, clusterID string, Status string, key string, APMac string,
	redisData publicstruct.RedisData, setData config.TerminalInfo, oSet map[string]interface{}, jsonInfo publicstruct.JSONInfo) {
	if Status == "00" {
		if msgType == globalmsgtype.MsgType.UPMsg.ZigbeeTerminalBindRspEvent {
			globallogger.Log.Infoln("devEUI :", devEUI, "procBindOrUnbindResult bind terminal success! endpoint:",
				SrcEndpoint, "clusterID:", clusterID)
			go publicfunction.IfMsgExistThenDelete(key, devEUI, redisData.SN, APMac)
		} else if msgType == globalmsgtype.MsgType.UPMsg.ZigbeeTerminalUnbindRspEvent {
			globallogger.Log.Infoln("devEUI :", devEUI, "procBindOrUnbindResult unbind terminal success! endpoint:",
				SrcEndpoint, "clusterID:", clusterID)
			var terminalInfoRes *config.TerminalInfo
			var err error
			if constant.Constant.UsePostgres {
				terminalInfoRes, err = models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, oSet)
			} else {
				terminalInfoRes, err = models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, setData)
			}
			if err != nil {
				globallogger.Log.Errorln("devEUI :", devEUI, "procBindOrUnbindResult FindTerminalAndUpdate error :", err)
			} else {
				globallogger.Log.Infoln("devEUI :", devEUI, "procBindOrUnbindResult FindTerminalAndUpdate success", terminalInfoRes)
				go publicfunction.IfMsgExistThenDelete(key, devEUI, redisData.SN, APMac)
			}
		}
	} else if Status == "01" {
		globallogger.Log.Warnln("devEUI :", devEUI, "procBindOrUnbindResult bind or unbind terminal failed! endpoint:",
			SrcEndpoint, "clusterID:", clusterID)
		go publicfunction.IfCtrlMsgExistThenDeleteAll(key, devEUI, APMac)
		go func() {
			timer := time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerAfterFive) * time.Second)
			<-timer.C
			publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
			timer.Stop()
		}()
	} else if Status == "8c" {
		globallogger.Log.Warnln("devEUI :", devEUI, "procBindOrUnbindResult bind terminal failed! bind table is full! endpoint:",
			SrcEndpoint, "clusterID:", clusterID)
		go publicfunction.IfMsgExistThenDelete(key, devEUI, redisData.SN, APMac)
	}
}

func sensorSendWriteReq(terminalInfo config.TerminalInfo) {
	zclmain.ZclMain(globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		config.TerminalInfo{
			DevEUI: terminalInfo.DevEUI,
		}, common.ZclDownMsg{
			MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
			DevEUI:      terminalInfo.DevEUI,
			T300ID:      terminalInfo.T300ID,
			CommandType: common.SensorWriteReq,
			ClusterID:   0x0500,
			Command: common.Command{
				DstEndpointIndex: 0,
			},
		}, "", "", nil)
}

func procBindOrUnbindTerminalRsp(jsonInfo publicstruct.JSONInfo, terminalInfo config.TerminalInfo) {
	bindInfoTemp := []config.BindInfo{}
	if constant.Constant.UsePostgres {
		json.Unmarshal([]byte(terminalInfo.BindInfoPG), &bindInfoTemp)
	} else {
		bindInfoTemp = terminalInfo.BindInfo
	}
	globallogger.Log.Infoln("devEUI :", terminalInfo.DevEUI, "procBindOrUnbindTerminalRsp, data:", jsonInfo.MessagePayload.Data)
	//var profileID = strings.Repeat(jsonInfo.MessagePayload.Data[0:4], 1)
	// var SrcAddr = strings.Repeat(jsonInfo.MessagePayload.Data[4:8], 1) //设备网络地址

	var key = publicfunction.GetCtrlMsgKey(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
	redisData, _ := publicfunction.GetRedisDataFromRedis(key, terminalInfo.DevEUI, jsonInfo.TunnelHeader.FrameSN)
	if redisData != nil {
		var udpMsg, _ = hex.DecodeString(string(redisData.SendBuf))
		var JSONInfo = publicfunction.ParseUDPMsg(udpMsg, jsonInfo.Rinfo)
		var SrcEndpoint = strings.Repeat(JSONInfo.MessagePayload.Data[0:2], 1)
		var clusterID = strings.Repeat(JSONInfo.MessagePayload.Data[2:6], 1)
		if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 {
			SrcEndpoint = strings.Repeat(JSONInfo.MessagePayload.Data[16:18], 1)
			clusterID = strings.Repeat(JSONInfo.MessagePayload.Data[18:22], 1)
		}
		var setData = config.TerminalInfo{}
		oSet := make(map[string]interface{})
		if strings.Repeat(JSONInfo.MessagePayload.Data[len(JSONInfo.MessagePayload.Data)-2:], 1) == "ff" { //绑定解绑T300回复处理
			var bindOrUnbindInfoArray = make([]config.BindInfo, len(bindInfoTemp))
			var clusterIDArray = make([]string, 20)
			var index = 0
			for bindInfoIndex, bindInfo := range bindInfoTemp {
				bindOrUnbindInfoArray[bindInfoIndex] = bindInfo
				if SrcEndpoint == bindInfo.SrcEndpoint {
					for clusterIDIndex, item := range bindInfo.ClusterID {
						if clusterID != item {
							clusterIDArray[index] = item
							index++
						}
						if (len(bindInfoTemp) == bindInfoIndex+1) &&
							(len(bindInfo.ClusterID) == clusterIDIndex+1) &&
							(clusterID == item) && (strings.Repeat(jsonInfo.MessagePayload.Data[8:10], 1) == "00") &&
							(jsonInfo.MessageHeader.MsgType == globalmsgtype.MsgType.UPMsg.ZigbeeTerminalBindRspEvent) {
							//只有最后一个clusterID绑定成功，才通知终端管理终端上线
							//传感器还需其他操作
							procTerminalOnlineToAPP(terminalInfo, jsonInfo)
						}
					}
					for i, v := range clusterIDArray {
						if v == "" {
							clusterIDArray = clusterIDArray[:i]
							break
						}
					}
					bindOrUnbindInfoArray[bindInfoIndex].ClusterID = clusterIDArray
					if clusterIDArray == nil {
						bindOrUnbindInfoArray[bindInfoIndex] = config.BindInfo{}
					}
				}
			}
			setData.BindInfo = bindOrUnbindInfoArray
			bindInfoByte, _ := json.Marshal(bindOrUnbindInfoArray)
			oSet["bindinfopg"] = string(bindInfoByte)
			setData.BindInfoPG = string(bindInfoByte)
			procBindOrUnbindResult(terminalInfo.DevEUI, jsonInfo.MessageHeader.MsgType, SrcEndpoint, clusterID,
				strings.Repeat(jsonInfo.MessagePayload.Data[8:10], 1),
				key, jsonInfo.TunnelHeader.LinkInfo.APMac, *redisData, setData, oSet, jsonInfo)
		} else { //设备与设备绑定解绑回复
			publicfunction.IfMsgExistThenDelete(key, terminalInfo.DevEUI, redisData.SN, jsonInfo.TunnelHeader.LinkInfo.APMac)
		}
	}
}

func procTerminalLeaveRsp(jsonInfo publicstruct.JSONInfo, terminalInfo *config.TerminalInfo) {
	globallogger.Log.Infoln("devEUI :", terminalInfo.DevEUI, "procTerminalLeaveRsp, data:", jsonInfo.MessagePayload.Data)
	//var profileID = strings.Repeat(jsonInfo.MessagePayload.Data[0:4], 1)
	//var SrcAddr = strings.Repeat(jsonInfo.MessagePayload.Data[4:8], 1) //设备网络地址
	var key = publicfunction.GetCtrlMsgKey(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)

	if strings.Repeat(jsonInfo.MessagePayload.Data[8:10], 1) == "00" { //success
		if terminalInfo != nil {
			publicfunction.TerminalOffline(terminalInfo.DevEUI)
			if constant.Constant.UsePostgres {
				_, err := models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": terminalInfo.DevEUI}, map[string]interface{}{"leavestate": true})
				if err != nil {
					globallogger.Log.Errorln("devEUI :", terminalInfo.DevEUI, "procTerminalLeaveRsp FindTerminalAndUpdate err :", err)
				}
			} else {
				_, err := models.FindTerminalAndUpdate(bson.M{"devEUI": terminalInfo.DevEUI}, bson.M{"leaveState": true})
				if err != nil {
					globallogger.Log.Errorln("devEUI :", terminalInfo.DevEUI, "procTerminalLeaveRsp FindTerminalAndUpdate err :", err)
				}
			}
			if constant.Constant.Iotware {
				iotsmartspace.StateTerminalLeaveIotware(*terminalInfo)
			} else if constant.Constant.Iotedge {
				iotsmartspace.StateTerminalLeave(terminalInfo.DevEUI)
			}
		}
		publicfunction.IfCtrlMsgExistThenDeleteAll(key, terminalInfo.DevEUI, jsonInfo.TunnelHeader.LinkInfo.APMac)
		key = publicfunction.GetDataMsgKey(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
		publicfunction.IfDataMsgExistThenDeleteAll(key, terminalInfo.DevEUI, jsonInfo.TunnelHeader.LinkInfo.APMac)
	} else {
		globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "procTerminalLeaveRsp failed ")
		redisData, _ := publicfunction.GetRedisDataFromRedis(key, terminalInfo.DevEUI, jsonInfo.TunnelHeader.FrameSN)
		if redisData != nil {
			publicfunction.CheckRedisAndReSend(terminalInfo.DevEUI, key, *redisData, jsonInfo.TunnelHeader.LinkInfo.APMac)
		}
	}
}

func procTerminalPermitJoinRsp(jsonInfo publicstruct.JSONInfo, terminalInfo config.TerminalInfo) {
	globallogger.Log.Infoln("devEUI :", terminalInfo.DevEUI, "procTerminalPermitJoinRsp, data:", jsonInfo.MessagePayload.Data)
	//var profileID = strings.Repeat(jsonInfo.MessagePayload.Data[0:4], 1)
	//var SrcAddr = strings.Repeat(jsonInfo.MessagePayload.Data[4:8], 1) //设备网络地址
	var key = publicfunction.GetCtrlMsgKey(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)

	redisData, _ := publicfunction.GetRedisDataFromRedis(key, terminalInfo.DevEUI, jsonInfo.TunnelHeader.FrameSN)
	if redisData != nil {
		if strings.Repeat(jsonInfo.MessagePayload.Data[8:10], 1) == "00" { //success
			var setData = bson.M{}
			setData["PermitJoin"] = false //不允许加入成功
			if redisData.Data[0:2] == "ff" {
				setData["PermitJoin"] = true //允许加入成功
			}
			if constant.Constant.UsePostgres {
				_, err := models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": terminalInfo.DevEUI},
					map[string]interface{}{"permitjoin": setData["PermitJoin"]})
				if err != nil {
					globallogger.Log.Errorln("devEUI :", terminalInfo.DevEUI, "procTerminalPermitJoinRsp FindTerminalAndUpdate error :", err)
				}
			} else {
				_, err := models.FindTerminalAndUpdate(bson.M{"devEUI": terminalInfo.DevEUI}, setData)
				if err != nil {
					globallogger.Log.Errorln("devEUI :", terminalInfo.DevEUI, "procTerminalPermitJoinRsp FindTerminalAndUpdate error :", err)
				}
			}
			publicfunction.IfMsgExistThenDelete(key, terminalInfo.DevEUI, jsonInfo.TunnelHeader.FrameSN, jsonInfo.TunnelHeader.LinkInfo.APMac)
		} else {
			globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "procTerminalPermitJoinRsp failed ")
			publicfunction.CheckRedisAndReSend(terminalInfo.DevEUI, key, *redisData, jsonInfo.TunnelHeader.LinkInfo.APMac)
		}
	}
}
