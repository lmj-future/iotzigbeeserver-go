package controllers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dyrkin/znp-go"
	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmemorycache"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalrediscache"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	zclmain "github.com/h3c/iotzigbeeserver-go/zcl"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
	"github.com/lib/pq"
)

var zigbeeServerKeepAliveTimerID = &sync.Map{}

/*
server 接收到ACK消息
1.如果消息队列中有相应消息，说明ack比rsp先到，置ackFlag标记
2.如果消息队列中没有相应消息，说明rsp比ack先到，不做处理
*/
func procACKMsg(jsonInfo publicstruct.JSONInfo) {
	var devEUI = jsonInfo.MessagePayload.Address
	var data = jsonInfo.MessagePayload.Data
	var APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
	var moduleID = jsonInfo.MessagePayload.ModuleID
	var SN = jsonInfo.TunnelHeader.FrameSN
	globallogger.Log.Infoln("devEUI : " + devEUI + " " + "procACKMsg receive msgType: " + data +
		globalmsgtype.MsgType.GetMsgTypeDOWNMsgMeaning(data) + " ACK success!")
	var key = constant.Constant.REDIS.ZigbeeRedisCtrlMsgKey + APMac + "_" + moduleID
	var isCtrlMsg = true
	if data == globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent {
		key = constant.Constant.REDIS.ZigbeeRedisDataMsgKey + APMac + "_" + moduleID
		isCtrlMsg = false
	}
	_, redisArray, _ := publicfunction.GetRedisLengthAndRangeRedis(key, devEUI)
	if len(redisArray) != 0 {
		var indexTemp = 0
		for itemIndex, item := range redisArray {
			if item != constant.Constant.REDIS.ZigbeeRedisRemoveRedis {
				var itemData publicstruct.RedisData
				json.Unmarshal([]byte(item), &itemData) //JSON.parse(item)
				if itemData.AckFlag == false {
					if SN == itemData.SN {
						// globallogger.Log.Infoln("devEUI : " + devEUI + " " + " procACKMsg ack is faster than rsp")
						itemData.AckFlag = true
						itemDataSerial, _ := json.Marshal(itemData)
						if constant.Constant.MultipleInstances {
							_, err := globalrediscache.RedisCache.SetRedis(key, itemIndex, itemDataSerial)
							if err != nil {
								globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" procACKMsg set redis error : ", err)
							} else {
								// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" procACKMsg set redis success : ", res)
							}
						} else {
							_, err := globalmemorycache.MemoryCache.SetMemory(key, itemIndex, itemDataSerial)
							if err != nil {
								globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" procACKMsg set momery error : ", err)
							} else {
								// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" procACKMsg set momery success : ", res)
							}
						}
						indexTemp = itemIndex + 1
					}
					if !isCtrlMsg {
						if indexTemp > 0 && indexTemp == itemIndex {
							publicfunction.CheckRedisAndSend(devEUI, key, itemData.SN, APMac)
						}
					}
				}
			} else {
				indexTemp++
			}
		}
	} else {
		// globallogger.Log.Infoln("devEUI : " + devEUI + " " + " procACKMsg ack is slowly than rsp, do nothing")
	}
}

func saveMsgDuplicationFlag(devEUI string, msgType string, SN string, isDataUpProc bool) (string, error) {
	var key = constant.Constant.REDIS.ZigbeeRedisDuplicationFlagKey + devEUI + "_" + msgType + "_" + SN //client上行的消息接收标记key
	if isDataUpProc {
		key = key + "_DataUpProc"
	}
	var errorCode string
	var err error
	if constant.Constant.MultipleInstances {
		res, err := globalrediscache.RedisCache.GetRedis(key)
		if err != nil {
			globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag get redis error : ", err)
			errorCode = "error"
		} else {
			// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag get redis success : ", res)
			if res != "" {
				// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "saveMsgDuplicationFlag server has already receive this msg, do nothing.")
				errorCode = "doNothing"
			} else {
				// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "saveMsgDuplicationFlag server is first receive this msg.")
				errorCode = "toDo"
				go func() {
					// jsonInfoSerial, _ := json.Marshal(jsonInfo)
					_, err = globalrediscache.RedisCache.UpdateRedis(key, []byte{1})
					if err != nil {
						globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag update redis error : ", err)
					} else {
						_, err := globalrediscache.RedisCache.SaddRedis(constant.Constant.REDIS.ZigbeeRedisDuplicationFlagSets, key)
						if err != nil {
							globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag sadd redis error : ", err)
						} else {
							// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag sadd redis success: ", resSadd)
						}
						// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag update redis success : ", resSadd)
						select {
						case <-time.After(time.Duration(constant.Constant.TIMER.ZigbeeTimerDuplicationFlag) * time.Second):
							_, err := globalrediscache.RedisCache.DeleteRedis(key)
							if err != nil {
								globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag delete redis error : ", err)
							} else {
								// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag delete redis success : ", res)
								_, err := globalrediscache.RedisCache.SremRedis(constant.Constant.REDIS.ZigbeeRedisDuplicationFlagSets, key)
								if err != nil {
									globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag srem redis error : ", err)
								} else {
									// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag srem redis success: ", redisData)
								}
							}
						}
					}
				}()
			}
		}
	} else {
		res, err := globalmemorycache.MemoryCache.GetMemory(key)
		if err != nil {
			globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag get memory error : ", err)
			errorCode = "error"
		} else {
			// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag get memory success : ", res)
			if res != "" {
				// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "saveMsgDuplicationFlag server has already receive this msg, do nothing.")
				errorCode = "doNothing"
			} else {
				// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "saveMsgDuplicationFlag server is first receive this msg.")
				errorCode = "toDo"
				go func() {
					// jsonInfoSerial, _ := json.Marshal(jsonInfo)
					_, err := globalmemorycache.MemoryCache.UpdateMemory(key, []byte{1})
					if err != nil {
						globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag update memory error : ", err)
					} else {
						// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag update memory success : ", resUpdate)
						select {
						case <-time.After(time.Duration(constant.Constant.TIMER.ZigbeeTimerDuplicationFlag) * time.Second):
							_, err := globalmemorycache.MemoryCache.DeleteMemory(key)
							if err != nil {
								globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag delete memory error : ", err)
							} else {
								// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" saveMsgDuplicationFlag delete memory success : ", res)
							}
						}
					}
				}()
			}
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
	var devEUI = jsonInfo.MessagePayload.Address
	var msgType = jsonInfo.MessageHeader.MsgType
	var SN = jsonInfo.TunnelHeader.FrameSN
	var firmTopic = jsonInfo.MessagePayload.Topic
	globallogger.Log.Infoln("devEUI : " + devEUI + " " + "procAnyMsg receive msgType: " + msgType +
		globalmsgtype.MsgType.GetMsgTypeUPMsgMeaning(msgType) + " success!")
	go publicfunction.SendACK(jsonInfo)
	res, _ := saveMsgDuplicationFlag(devEUI+jsonInfo.TunnelHeader.LinkInfo.APMac, msgType, SN, false)
	if res == "toDo" {
		if firmTopic == "80003002" {
			collihighProc(jsonInfo)
		} else {
			procMainMsg(jsonInfo, 1)
		}
	}
}

func procKeepAliveMsg(APMac string, t int) {
	if value, ok := zigbeeServerKeepAliveTimerID.Load(APMac); ok {
		value.(*time.Timer).Stop()
		// globallogger.Log.Infoln("APMac: "+APMac+" procKeepAliveMsg clearTimeout, time: ", t)
	}
	var timerID = time.NewTimer(time.Duration(t*4) * time.Second)
	zigbeeServerKeepAliveTimerID.Store(APMac, timerID)

	go func() {
		select {
		case <-timerID.C:
			if value, ok := zigbeeServerKeepAliveTimerID.Load(APMac); ok {
				if timerID == value.(*time.Timer) {
					zigbeeServerKeepAliveTimerID.Delete(APMac)
					keepAliveTimerInfo, _ := publicfunction.GetKeepAliveTimerByAPMac(APMac)
					if keepAliveTimerInfo != nil {
						var dateNow = time.Now()
						var num = dateNow.UnixNano() - keepAliveTimerInfo.UpdateTime.UnixNano()
						if num > int64(time.Duration(3*t)*time.Second) {
							globallogger.Log.Warnln("APMac: "+APMac+" procKeepAliveMsg server has not already recv keep alive msg for ",
								num/int64(time.Second), " seconds. Then offline all terminal")
							terminalList, err := publicfunction.ProcKeepAliveTerminalOffline(APMac)
							if err != nil {
								globallogger.Log.Errorln("procKeepAliveMsg ProcKeepAliveTerminalOffline err: ", err)
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
							go publicfunction.DeleteKeepAliveTimer(APMac)
							var zhaSNKey = constant.Constant.REDIS.ZigbeeRedisSNZHA + APMac
							var transferSNKey = constant.Constant.REDIS.ZigbeeRedisSNTransfer + APMac
							if constant.Constant.MultipleInstances {
								go globalrediscache.RedisCache.DeleteRedis(zhaSNKey)
								go globalrediscache.RedisCache.DeleteRedis(transferSNKey)
							} else {
								go globalmemorycache.MemoryCache.DeleteMemory(zhaSNKey)
								go globalmemorycache.MemoryCache.DeleteMemory(transferSNKey)
							}
						}
					}
				}
			}
		}
	}()
}

func procModuleRemoveEvent(jsonInfo publicstruct.JSONInfo) {
	var APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
	var moduleID = jsonInfo.MessagePayload.ModuleID
	var terminalList []config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalList, err = models.FindTerminalByAPMacAndModuleIDPG(APMac, moduleID)
	} else {
		terminalList, err = models.FindTerminalByAPMacAndModuleID(APMac, moduleID)
	}
	if err != nil {
		globallogger.Log.Errorln("procModuleRemoveEvent findTerminalByAPMacAndModuleID err: ", err)
		return
	}
	if terminalList != nil {
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
	var APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
	var moduleID = jsonInfo.MessagePayload.ModuleID
	var terminalList []config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalList, err = models.FindTerminalByAPMacAndModuleIDPG(APMac, moduleID)
	} else {
		terminalList, err = models.FindTerminalByAPMacAndModuleID(APMac, moduleID)
	}
	if err != nil {
		globallogger.Log.Errorln("procEnvironmentChangeEvent findTerminalByAPMacAndModuleID err: ", err)
	} else {
		if terminalList != nil {
			for _, item := range terminalList {
				publicfunction.TerminalOffline(item.DevEUI)
				if constant.Constant.UsePostgres {
					oMatch := make(map[string]interface{})
					oSet := make(map[string]interface{})
					oMatch["deveui"] = item.DevEUI
					oSet["leavestate"] = true
					_, err := models.FindTerminalAndUpdatePG(oMatch, oSet)
					if err != nil {
						globallogger.Log.Errorln("devEUI : "+item.DevEUI+" "+"procEnvironmentChangeEvent FindTerminalAndUpdatePG err : ", err)
					}
				} else {
					_, err := models.FindTerminalAndUpdate(bson.M{"devEUI": item.DevEUI}, bson.M{"leaveState": true})
					if err != nil {
						globallogger.Log.Errorln("devEUI : "+item.DevEUI+" "+"procEnvironmentChangeEvent FindTerminalAndUpdate err : ", err)
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

func findTerminalAndUpdate(jsonInfo publicstruct.JSONInfo) (*config.TerminalInfo, error) {
	var setData = config.TerminalInfo{}
	oSet := make(map[string]interface{})
	if jsonInfo.TunnelHeader.LinkInfo.ACMac != "" {
		setData.ACMac = jsonInfo.TunnelHeader.LinkInfo.ACMac
		oSet["acmac"] = jsonInfo.TunnelHeader.LinkInfo.ACMac
	}
	if jsonInfo.TunnelHeader.LinkInfo.T300ID != "" && jsonInfo.TunnelHeader.LinkInfo.T300ID != "000000000000" &&
		jsonInfo.TunnelHeader.LinkInfo.T300ID != "ffffffffffff" {
		setData.T300ID = jsonInfo.TunnelHeader.LinkInfo.T300ID
		oSet["t300id"] = jsonInfo.TunnelHeader.LinkInfo.T300ID
	}
	if jsonInfo.TunnelHeader.Version != "" {
		setData.UDPVersion = jsonInfo.TunnelHeader.Version
		oSet["udpversion"] = jsonInfo.TunnelHeader.Version
	}
	setData.APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
	oSet["apmac"] = jsonInfo.TunnelHeader.LinkInfo.APMac
	setData.ModuleID = jsonInfo.MessagePayload.ModuleID
	oSet["moduleid"] = jsonInfo.MessagePayload.ModuleID
	setData.UpdateTime = time.Now()
	oSet["updatetime"] = time.Now()
	var terminalInfo *config.TerminalInfo = nil
	var err error
	var oMatch = bson.M{"devEUI": jsonInfo.MessagePayload.Address}
	var oMatchPG = map[string]interface{}{"deveui": jsonInfo.MessagePayload.Address}
	if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 {
		oMatch = bson.M{"nwkAddr": jsonInfo.MessagePayload.Address, "T300ID": jsonInfo.TunnelHeader.LinkInfo.T300ID}
		oMatchPG = map[string]interface{}{"nwkaddr": jsonInfo.MessagePayload.Address, "t300id": jsonInfo.TunnelHeader.LinkInfo.T300ID}
		if jsonInfo.TunnelHeader.LinkInfo.T300ID == "000000000000" || jsonInfo.TunnelHeader.LinkInfo.T300ID == "ffffffffffff" {
			oMatch = bson.M{
				"nwkAddr":  jsonInfo.MessagePayload.Address,
				"ACMac":    jsonInfo.TunnelHeader.LinkInfo.ACMac,
				"moduleID": jsonInfo.MessagePayload.ModuleID,
			}
			oMatchPG = map[string]interface{}{
				"nwkaddr":  jsonInfo.MessagePayload.Address,
				"acmac":    jsonInfo.TunnelHeader.LinkInfo.ACMac,
				"moduleid": jsonInfo.MessagePayload.ModuleID,
			}
		}
	}
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.FindTerminalAndUpdatePG(oMatchPG, oSet)
	} else {
		terminalInfo, err = models.FindTerminalAndUpdate(oMatch, setData)
	}
	if err != nil {
		globallogger.Log.Errorln("devEUI : "+jsonInfo.MessagePayload.Address+" "+"findTerminalAndUpdate FindTerminalAndUpdate error : ", err)
	} else {
		if terminalInfo == nil {
			globallogger.Log.Warnln("devEUI : " + jsonInfo.MessagePayload.Address + " " + "findTerminalAndUpdate terminal is not exist, please first create !")
		} else {
			// globallogger.Log.Infoln("devEUI : " + jsonInfo.MessagePayload.Address + " " + "findTerminalAndUpdate FindTerminalAndUpdate success")
		}
	}
	return terminalInfo, err
}

func procMainMsg(jsonInfo publicstruct.JSONInfo, index int) {
	var msgType = jsonInfo.MessageHeader.MsgType
	if msgType == globalmsgtype.MsgType.UPMsg.ZigbeeTerminalJoinEvent {
		//设备入网
		procTerminalJoinEvent(jsonInfo)
	} else if msgType == globalmsgtype.MsgType.UPMsg.ZigbeeWholeNetworkIEEERspEvent {
		//整个网络拓扑获取
	} else if msgType == globalmsgtype.MsgType.UPMsg.ZigbeeNetworkEvent {
		//网络拓扑信息主动上报
	} else if msgType == globalmsgtype.MsgType.UPMsg.ZigbeeModuleRemoveEvent {
		//插卡下线事件
		procModuleRemoveEvent(jsonInfo)
	} else if msgType == globalmsgtype.MsgType.UPMsg.ZigbeeEnvironmentChangeEvent {
		//网络发生改变事件
		procEnvironmentChangeEvent(jsonInfo)
	} else if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 &&
		msgType == globalmsgtype.MsgType.UPMsg.ZigbeeTerminalLeaveEvent {
		//设备离网
		var setData = config.TerminalInfo{}
		oSet := make(map[string]interface{})
		if jsonInfo.TunnelHeader.LinkInfo.ACMac != "" {
			setData.ACMac = jsonInfo.TunnelHeader.LinkInfo.ACMac
			oSet["acmac"] = jsonInfo.TunnelHeader.LinkInfo.ACMac
		}
		if jsonInfo.TunnelHeader.LinkInfo.T300ID != "" {
			setData.T300ID = jsonInfo.TunnelHeader.LinkInfo.T300ID
			oSet["t300id"] = jsonInfo.TunnelHeader.LinkInfo.T300ID
		}
		if jsonInfo.TunnelHeader.Version != "" {
			setData.UDPVersion = jsonInfo.TunnelHeader.Version
			oSet["udpversion"] = jsonInfo.TunnelHeader.Version
		}
		setData.APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
		oSet["apmac"] = jsonInfo.TunnelHeader.LinkInfo.APMac
		setData.ModuleID = jsonInfo.MessagePayload.ModuleID
		oSet["moduleid"] = jsonInfo.MessagePayload.ModuleID
		setData.UpdateTime = time.Now()
		oSet["updatetime"] = time.Now()
		var oMatch = bson.M{"devEUI": strings.ToUpper(jsonInfo.MessagePayload.Data)}
		var oMatchPG = map[string]interface{}{"deveui": strings.ToUpper(jsonInfo.MessagePayload.Data)}
		var terminalInfo *config.TerminalInfo
		var err error
		if constant.Constant.UsePostgres {
			terminalInfo, err = models.FindTerminalAndUpdatePG(oMatchPG, oSet)
		} else {
			terminalInfo, err = models.FindTerminalAndUpdate(oMatch, setData)
		}
		if err == nil && terminalInfo != nil {
			procTerminalLeaveEvent(jsonInfo)
		}
	} else {
		terminalInfo, err := findTerminalAndUpdate(jsonInfo)
		if err == nil {
			if terminalInfo != nil {
				switch msgType {
				case globalmsgtype.MsgType.GENERALMsg.ZigbeeGeneralFailed,
					globalmsgtype.MsgType.GENERALMsgV0101.ZigbeeGeneralFailedV0101:
					//处理失败消息
					procFailedMsg(jsonInfo, terminalInfo.DevEUI, true)
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalCheckExistRspEvent:
					//检查存在否回复
					procTerminalCheckExistEvent(jsonInfo)
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalLeaveEvent:
					//设备离网
					procTerminalLeaveEvent(jsonInfo)
				case globalmsgtype.MsgType.UPMsg.ZigbeeDataUpEvent:
					//设备数据上报
					procDataUp(jsonInfo, *terminalInfo)
					if terminalInfo.Online == false {
						publicfunction.TerminalOnline(terminalInfo.DevEUI, false)
						if constant.Constant.Iotware {
							if terminalInfo.IsExist {
								iotsmartspace.StateTerminalOnlineIotware(*terminalInfo)
							}
						} else if constant.Constant.Iotedge {
							iotsmartspace.StateTerminalOnline(terminalInfo.DevEUI)
						}
					}
				case globalmsgtype.MsgType.UPMsg.ZigbeeWholeNetworkNWKRspEvent:
					//单个网络拓扑获取
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalEndpointRspEvent:
					//设备端口号上报
					procTerminalEndpointInfoUp(jsonInfo, *terminalInfo)
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalDiscoveryRspEvent:
					//设备服务发现回复
					procTerminalDiscoveryInfoUp(jsonInfo, *terminalInfo)
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalBindRspEvent,
					globalmsgtype.MsgType.UPMsg.ZigbeeTerminalUnbindRspEvent:
					//绑定/解绑回复
					procBindOrUnbindTerminalRsp(jsonInfo, *terminalInfo)
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalNetworkRspEvent:
					//设备邻居网络拓扑上报
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalLeaveRspEvent:
					//设备离网回复
					procTerminalLeaveRsp(jsonInfo, terminalInfo)
				case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalPermitJoinRspEvent:
					//允许设备入网回复
					procTerminalPermitJoinRsp(jsonInfo, *terminalInfo)
				default:
					globallogger.Log.Warnln("devEUI : " + terminalInfo.DevEUI + " " + "unknow msgType : " + msgType)
				}
			}
		}
	}
}

func procTerminalCheckExistEvent(jsonInfo publicstruct.JSONInfo) {
	globallogger.Log.Warnln("devEUI : " + jsonInfo.MessagePayload.Address + " " + "procTerminalCheckExistEvent, data: " + jsonInfo.MessagePayload.Data)
	var devEUI = jsonInfo.MessagePayload.Address
	var data = jsonInfo.MessagePayload.Data
	// var ProfileID = data[0:4]
	var MacAddr = data[4:20]
	// var NwkAddr = data[20:24]
	var Status = data[24:32]
	if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 {
		devEUI = strings.ToUpper(MacAddr)
	}
	if Status != "00000012" {
		publicfunction.TerminalOffline(devEUI)
		var terminalInfo *config.TerminalInfo
		var err error
		if constant.Constant.UsePostgres {
			_, err := models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, map[string]interface{}{"leavestate": true})
			if err != nil {
				globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procTerminalCheckExistEvent FindTerminalAndUpdatePG err : ", err)
			}
			terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
		} else {
			_, err := models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, bson.M{"leaveState": true})
			if err != nil {
				globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procTerminalCheckExistEvent FindTerminalAndUpdate err : ", err)
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
		var key = constant.Constant.REDIS.ZigbeeRedisCtrlMsgKey + jsonInfo.TunnelHeader.LinkInfo.APMac + "_" + jsonInfo.MessagePayload.ModuleID
		publicfunction.IfCtrlMsgExistThenDeleteAll(key, devEUI, jsonInfo.TunnelHeader.LinkInfo.APMac)
		key = constant.Constant.REDIS.ZigbeeRedisDataMsgKey + jsonInfo.TunnelHeader.LinkInfo.APMac + "_" + jsonInfo.MessagePayload.ModuleID
		publicfunction.IfDataMsgExistThenDeleteAll(key, devEUI, jsonInfo.TunnelHeader.LinkInfo.APMac)
	} else {
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procTerminalCheckExistEvent, terminal is exist")
	}
}

func procTerminalLeaveEvent(jsonInfo publicstruct.JSONInfo) {
	globallogger.Log.Warnln("devEUI : " + jsonInfo.MessagePayload.Address + " " + "procTerminalLeaveEvent")
	var devEUI = jsonInfo.MessagePayload.Address
	if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 {
		devEUI = strings.ToUpper(jsonInfo.MessagePayload.Data)
	}
	var key = constant.Constant.REDIS.ZigbeeRedisCtrlMsgKey + jsonInfo.TunnelHeader.LinkInfo.APMac + "_" + jsonInfo.MessagePayload.ModuleID
	publicfunction.TerminalOffline(devEUI)
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		_, err := models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, map[string]interface{}{"leavestate": true})
		if err != nil {
			globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procTerminalLeaveEvent FindTerminalAndUpdatePG err : ", err)
		}
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		_, err := models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, bson.M{"leaveState": true})
		if err != nil {
			globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procTerminalLeaveEvent FindTerminalAndUpdate err : ", err)
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
	publicfunction.IfCtrlMsgExistThenDeleteAll(key, devEUI, jsonInfo.TunnelHeader.LinkInfo.APMac)
	key = constant.Constant.REDIS.ZigbeeRedisDataMsgKey + jsonInfo.TunnelHeader.LinkInfo.APMac + "_" + jsonInfo.MessagePayload.ModuleID
	publicfunction.IfDataMsgExistThenDeleteAll(key, devEUI, jsonInfo.TunnelHeader.LinkInfo.APMac)
}

func procJoinTerminalIsExist(terminalInfo config.TerminalInfo, setData config.TerminalInfo,
	oSet map[string]interface{}, jsonInfo publicstruct.JSONInfo, devEUI string) {
	if terminalInfo.APMac != setData.APMac {
		jsonInfo := publicfunction.GetJSONInfo(terminalInfo)
		go publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
	}
	setData.Endpoint = terminalInfo.Endpoint
	setData.EndpointPG = terminalInfo.EndpointPG
	setData.EndpointCount = terminalInfo.EndpointCount
	setData.BindInfo = terminalInfo.BindInfo
	setData.BindInfoPG = terminalInfo.BindInfoPG
	setData.TmnType = "invalidType"
	globallogger.Log.Infoln("devEUI : " + devEUI + " " + "procJoinTerminalIsExist setData: " + fmt.Sprintf("%+v", setData))
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
		globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procJoinTerminalIsExist FindTerminalAndUpdate error : ", err)
	} else {
		go firstlyGetTerminalEndpoint(jsonInfo, devEUI)
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
		globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procJoinTerminalNotExist createTerminal error : ", err)
	} else {
		// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"procJoinTerminalNotExist createTerminal success ", setData)
		firstlyGetTerminalEndpoint(jsonInfo, devEUI)
	}
}

func joinProcBind(tmnInfo *config.TerminalInfo, NwkAddr string, devEUI string, oSet map[string]interface{}, setData config.TerminalInfo) {
	if constant.Constant.UsePostgres {
		models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, oSet)
	} else {
		models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, setData)
	}
	tmnInfo.NwkAddr = NwkAddr
	tmnInfo.ACMac = setData.ACMac
	tmnInfo.APMac = setData.APMac
	tmnInfo.T300ID = setData.T300ID
	tmnInfo.ModuleID = setData.ModuleID
	publicfunction.SendBindTerminalReq(*tmnInfo)
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
	globallogger.Log.Warnln("devEUI : " + jsonInfo.MessagePayload.Address + " " + "procTerminalJoinEvent, data: " + jsonInfo.MessagePayload.Data)
	var data = jsonInfo.MessagePayload.Data
	var devEUI = jsonInfo.MessagePayload.Address
	var APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
	var moduleID = jsonInfo.MessagePayload.ModuleID
	var NwkAddr = data[0:4]           //设备网络地址
	var MacAddr = data[4:20]          //设备物理地址
	var CapabilityFlags = data[20:22] //设备功能、性能标记
	var Type = data[22:24]            //设备入网方式
	devEUI = strings.ToUpper(MacAddr)

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
	setData.APMac = APMac
	oSet["apmac"] = APMac
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
	setData.ModuleID = moduleID
	oSet["moduleid"] = moduleID
	setData.UpdateTime = time.Now()
	oSet["updatetime"] = time.Now()
	if constant.Constant.UsePostgres {
		tmnInfo, err = models.FindTerminalByConditionPG(map[string]interface{}{"deveui": devEUI})
	} else {
		tmnInfo, err = models.FindTerminalByCondition(oMatch)
	}
	if err != nil {
		globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procTerminalJoinEvent FindTerminalByCondition err : ", err)
		publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
		return
	}
	if tmnInfo != nil && tmnInfo.NwkAddr != "" {
		if NwkAddr == tmnInfo.NwkAddr {
			globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procTerminalJoinEvent this join is rejoin and nwkAddr is not change, do nothing")
			return
		}
	}
	if constant.Constant.Iotprivate {
		httpTerminalInfo := publicfunction.HTTPRequestTerminalInfo(devEUI, "ZHA")
		if httpTerminalInfo != nil {
			// isPermitJoin := publicfunction.HTTPRequestCheckLicense(devEUI, terminalInfo.TenantID)
			// if !isPermitJoin {
			// 	globallogger.Log.Warnln("devEUI : "+devEUI+" "+"procTerminalJoinEvent publicfunction.HTTPRequestCheckLicense not permit join: ", isPermitJoin)
			// 	publicfunction.SendTerminalLeaveReq(jsonInfo)
			// 	return
			// }
			if !publicfunction.CheckTerminalSceneIsMatch(devEUI, httpTerminalInfo.SceneID, jsonInfo.TunnelHeader.VenderInfo) {
				globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procTerminalJoinEvent publicfunction.CheckTerminalSceneIsMatch not match")
				publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
				return
			}
			setData.TmnName = httpTerminalInfo.TmnName
			setData.OIDIndex = httpTerminalInfo.TmnOIDIndex
			setData.ScenarioID = httpTerminalInfo.SceneID
			setData.UserName = httpTerminalInfo.TenantID
			setData.FirmTopic = httpTerminalInfo.FirmTopic
			setData.ProfileID = httpTerminalInfo.LinkType
			setData.TmnType2 = httpTerminalInfo.TmnType
			setData.IsExist = true
			oSet["tmnname"] = httpTerminalInfo.TmnName
			oSet["oidindex"] = httpTerminalInfo.TmnOIDIndex
			oSet["scenarioid"] = httpTerminalInfo.SceneID
			oSet["userName"] = httpTerminalInfo.TenantID
			oSet["firmtopic"] = httpTerminalInfo.FirmTopic
			oSet["profileid"] = httpTerminalInfo.LinkType
			oSet["tmntype2"] = httpTerminalInfo.TmnType
			oSet["isexist"] = true
		} else {
			globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procTerminalJoinEvent HTTPRequestTerminalInfo nil")
			publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
			return
		}
	} else {
		if tmnInfo == nil || !tmnInfo.IsExist {
			if !constant.Constant.Iotware {
				key := constant.Constant.REDIS.ZigbeeRedisPermitJoinContentSets
				var permitJoinContentList []string
				if constant.Constant.MultipleInstances {
					permitJoinContentList, err = globalrediscache.RedisCache.FindAllRedisKeys(key)
					if err != nil {
						globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procTerminalJoinEvent FindAllRedisKeys err : ", err)
						publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
						return
					}
					globallogger.Log.Infoln("devEUI : "+devEUI+" "+"procTerminalJoinEvent FindAllRedisKeys success : ", permitJoinContentList)
					if len(permitJoinContentList) == 0 {
						globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procTerminalJoinEvent permit join is not enable")
						publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
						return
					}
				} else {
					permitJoinContentList, err = globalmemorycache.MemoryCache.FindAllMemoryKeys(key)
					if err != nil {
						globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procTerminalJoinEvent FindAllMemoryKeys err : ", err)
						publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
						return
					}
					globallogger.Log.Infoln("devEUI : "+devEUI+" "+"procTerminalJoinEvent FindAllMemoryKeys success : ", permitJoinContentList)
				}
				if len(permitJoinContentList) == 0 {
					globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procTerminalJoinEvent permit join is not enable")
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
					globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procTerminalJoinEvent permit join is not enable")
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
			globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" checkNextRedis get redis by index error : ", err)
		} else {
			// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" checkNextRedis get redis by index success : ", res)
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
	tempValue := "0000" + strconv.FormatInt(seqNum, 16)
	var SN = tempValue[len(tempValue)-4:]
	// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "procNeedCheckData get SN: " + SN)
	var key = constant.Constant.REDIS.ZigbeeRedisDataMsgKey + APMac + "_" + moduleID
	length, index, redisData, _ := publicfunction.GetRedisDataFromDataRedis(key, devEUI, SN)
	if redisData != nil {
		var err error
		// var res string
		if constant.Constant.MultipleInstances {
			_, err = globalrediscache.RedisCache.SetRedis(key, index, []byte(constant.Constant.REDIS.ZigbeeRedisRemoveRedis))
		} else {
			_, err = globalmemorycache.MemoryCache.SetMemory(key, index, []byte(constant.Constant.REDIS.ZigbeeRedisRemoveRedis))
		}
		if err != nil {
			globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" procNeedCheckData set redis error : ", err)
		} else {
			// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" procNeedCheckData set redis success : ", res)
			if redisData.AckFlag == false {
				// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "key: " + key + " procNeedCheckData rsp is faster than ack, send the next data msg ")
				index++
				go checkNextRedis(devEUI, APMac, key, index, length)
			} else {
				// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "key: " + key + " procNeedCheckData ack is faster than rsp, do nothing ")
			}
			// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "procNeedCheckData iotzigbeeserver send msg : " + fmt.Sprintf("%+v", afIncomingMessage))
			zclmain.ZclMain(globalmsgtype.MsgType.UPMsg.ZigbeeDataUpEvent, tmnInfo, afIncomingMessage, redisData.MsgID, redisData.Data, zclAfIncomingMessage)
		}
	} else {
		// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "procNeedCheckData iotzigbeeserver send msg : " + fmt.Sprintf("%+v", afIncomingMessage))
		zclmain.ZclMain(globalmsgtype.MsgType.UPMsg.ZigbeeDataUpEvent, tmnInfo, afIncomingMessage, "", "", zclAfIncomingMessage)
	}
}

func procBasicReadRsp(devEUI string, jsonInfo publicstruct.JSONInfo, terminalInfo config.TerminalInfo, command interface{}) {
	readAttributesRsp := command.(*cluster.ReadAttributesResponse)
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
	for _, value := range readAttributesRsp.ReadAttributeStatuses {
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
			globallogger.Log.Warnln("devEUI : "+devEUI+" "+"procBasicReadRsp unknow attributeName", value.AttributeName)
		}
	}
	globallogger.Log.Infof("[devEUI: %v][procBasicReadRsp] basicInfo: %+v", devEUI, basicInfo)
	modelIdentifier := basicInfo.ModelIdentifier
	if basicInfo.ModelIdentifier == "FTB56+PTH01MK2.2" {
		basicInfo.ModelIdentifier = constant.Constant.TMNTYPE.MAILEKE.ZigbeeTerminalPMT1004Detector
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
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procBasicReadRsp tmnType is not nil, this readBasic for check terminal state")
		return
	}
	if basicInfo.ModelIdentifier != "" {
		var info = publicstruct.TmnTypeAttr{}
		info = publicfunction.GetTmnTypeAndAttribute(basicInfo.ModelIdentifier)
		isNeedBind := publicfunction.CheckTerminalIsNeedBind(devEUI, info.TmnType)
		interval := publicfunction.GetTerminalInterval(basicInfo.ModelIdentifier)
		oSet := make(map[string]interface{})
		var setData = config.TerminalInfo{}
		if constant.Constant.UsePostgres {
			if !constant.Constant.Iotprivate {
				oSet["tmnname"] = info.TmnType + "-" + devEUI
			}
			oSet["interval"] = interval
			oSet["tmntype"] = info.TmnType
			oSet["manufacturername"] = strings.ToUpper(basicInfo.ManufacturerName)
			oSet["profileid"] = "ZHA"
			oSet["statuspg"] = pq.StringArray(info.Attribute.Status)
			oSet["attributedata"] = info.Attribute.Data
			oSet["latestcommand"] = info.Attribute.LatestCommand
			oSet["isneedbind"] = isNeedBind
			oSet["isreadbasic"] = true
			oSet["updatetime"] = time.Now()
			oSet["endpointpg"] = terminalInfo.EndpointPG
			if modelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c ||
				modelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55 ||
				modelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105 ||
				modelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106 ||
				modelIdentifier == constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM {
				oSet["endpointpg"] = pq.StringArray{"01"}
				oSet["endpointcount"] = 1
			}
		} else {
			if !constant.Constant.Iotprivate {
				setData.TmnName = info.TmnType + "-" + devEUI
			}
			setData.Interval = interval
			setData.TmnType = info.TmnType
			setData.ManufacturerName = strings.ToUpper(basicInfo.ManufacturerName)
			setData.ProfileID = "ZHA"
			setData.Attribute = info.Attribute
			setData.IsNeedBind = isNeedBind
			setData.IsReadBasic = true
			setData.UpdateTime = time.Now()
			setData.Endpoint = terminalInfo.Endpoint
			if modelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c ||
				modelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55 ||
				modelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105 ||
				modelIdentifier == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106 ||
				modelIdentifier == constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM {
				setData.Endpoint = []string{"01"}
				setData.EndpointCount = 1
			}
		}
		if !constant.Constant.Iotware {
			if constant.Constant.Iotprivate {
				thirdlyDiscoveryByEndpoint(devEUI, setData, terminalInfo)
			} else {
				if terminalInfo.IsExist {
					if constant.Constant.UsePostgres {
						thirdlyDiscoveryByEndpointPG(devEUI, oSet, terminalInfo)
					} else {
						thirdlyDiscoveryByEndpoint(devEUI, setData, terminalInfo)
					}
					return
				}
				key := constant.Constant.REDIS.ZigbeeRedisPermitJoinContentKey + basicInfo.ModelIdentifier
				var redisDataStr string
				var err error
				if constant.Constant.MultipleInstances {
					redisDataStr, err = globalrediscache.RedisCache.GetRedis(key)
				} else {
					redisDataStr, err = globalmemorycache.MemoryCache.GetMemory(key)
				}
				if err != nil {
					globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procBasicReadRsp GetRedis err : ", err)
					publicfunction.SendTerminalLeaveReq(jsonInfo, terminalInfo.DevEUI)
					return
				}
				if redisDataStr == "" {
					globallogger.Log.Warnln("devEUI : "+devEUI+" "+"procBasicReadRsp GetRedis redisDataStr null : ", redisDataStr)
					publicfunction.SendTerminalLeaveReq(jsonInfo, terminalInfo.DevEUI)
					return
				}
				redisData := make(map[string]interface{}, 4)
				err = json.Unmarshal([]byte(redisDataStr), &redisData)
				if err != nil {
					globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procBasicReadRsp Unmarshal err : ", err)
					publicfunction.SendTerminalLeaveReq(jsonInfo, terminalInfo.DevEUI)
					return
				}
				tmnType := basicInfo.ModelIdentifier
				if redisData["tmnType"] != nil && redisData["tmnType"] != "" {
					tmnType = redisData["tmnType"].(string)
				}
				if tmnType == basicInfo.ModelIdentifier {
					if constant.Constant.UsePostgres {
						thirdlyDiscoveryByEndpointPG(devEUI, oSet, terminalInfo)
					} else {
						thirdlyDiscoveryByEndpoint(devEUI, setData, terminalInfo)
					}
				} else {
					globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procBasicReadRsp tmnType is not match, real tmnType: " +
						basicInfo.ModelIdentifier + " permit join tmnType: " + tmnType)
					publicfunction.SendTerminalLeaveReq(jsonInfo, terminalInfo.DevEUI)
				}
			}
		} else {
			if constant.Constant.UsePostgres {
				thirdlyDiscoveryByEndpointPG(devEUI, oSet, terminalInfo)
			} else {
				thirdlyDiscoveryByEndpoint(devEUI, setData, terminalInfo)
			}
		}
	} else {
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procBasicReadRsp permit join is not enable")
		publicfunction.SendTerminalLeaveReq(jsonInfo, terminalInfo.DevEUI)
	}
}

func procDataUp(jsonInfo publicstruct.JSONInfo, terminalInfo config.TerminalInfo) {
	var data = jsonInfo.MessagePayload.Data
	var devEUI = terminalInfo.DevEUI
	var APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
	var moduleID = jsonInfo.MessagePayload.ModuleID
	var msgType = jsonInfo.MessageHeader.MsgType
	var profileID = data[0:4]
	GroupID, _ := strconv.ParseUint(data[4:8], 16, 16)
	ClusterID, _ := strconv.ParseUint(data[8:12], 16, 16)
	SrcAddr := data[12:16]
	SrcEndpoint, _ := strconv.ParseUint(data[16:18], 16, 8)
	DstEndpoint, _ := strconv.ParseUint(data[16:20], 16, 8)
	WasBroadcast, _ := strconv.ParseUint(data[20:22], 16, 8)
	LinkQuality, _ := strconv.ParseUint(data[22:24], 16, 8)
	SecurityUse, _ := strconv.ParseUint(data[24:26], 16, 8)
	Timestamp, _ := strconv.ParseUint(data[26:34], 16, 32)
	TransSeqNumber, _ := strconv.ParseUint(data[34:36], 16, 8)
	// Len, _ := strconv.ParseUint(data[36:38], 16, 8)
	Data, _ := hex.DecodeString(data[38:])
	globallogger.Log.Infoln("devEUI : " + terminalInfo.DevEUI + " " + "procDataUp, data: " + data[38:])
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
	switch profileID {
	case "0104":
		z := zcl.New()
		zclAfIncomingMessage, err := z.ToZclIncomingMessage(&afIncomingMessage)
		if err != nil {
			globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procDataUp ToZclIncomingMessage err : ", err)
			return
		}
		commandName := zclAfIncomingMessage.Data.CommandName
		seqNum := int64(zclAfIncomingMessage.Data.TransactionSequenceNumber)
		res, _ := saveMsgDuplicationFlag(devEUI+APMac, msgType, strconv.Itoa(int(seqNum)), true)
		if res == "toDo" {
			if commandName != "ReportAttributes" && commandName != "ZoneEnrollRequest" {
				go procNeedCheckData(seqNum, devEUI, APMac, moduleID, terminalInfo, afIncomingMessage, zclAfIncomingMessage)
				if commandName == "ReadAttributesResponse" && ClusterID == 0x0000 {
					procBasicReadRsp(devEUI, jsonInfo, terminalInfo, zclAfIncomingMessage.Data.Command)
				}
			} else {
				zclmain.ZclMain(globalmsgtype.MsgType.UPMsg.ZigbeeDataUpEvent, terminalInfo, afIncomingMessage, "", "", zclAfIncomingMessage)
			}
		} else {
			globallogger.Log.Infoln("devEUI : " + devEUI + " " + "procDataUp msg is duplication, drop this msg")
		}
	case "ZLL":
	default:
		//to do
	}
}

func procTerminalEndpointInfoUp(jsonInfo publicstruct.JSONInfo, terminalInfo config.TerminalInfo) {
	globallogger.Log.Infoln("devEUI : " + terminalInfo.DevEUI + " " + "procTerminalEndpointInfoUp, data: " + jsonInfo.MessagePayload.Data)
	var data = jsonInfo.MessagePayload.Data
	var devEUI = terminalInfo.DevEUI
	var APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
	var moduleID = jsonInfo.MessagePayload.ModuleID
	var SN = jsonInfo.TunnelHeader.FrameSN
	//var profileID = data[0:4]
	//var SrcAddr = data[4:8]         //设备网络地址
	var Status = data[8:10] //SUCCESS(0) or FAILURE(1)
	//var NWKAddr = data[10:14]       //设备短地址
	var ActiveEPCount = data[14:16] //设备端口号数量
	var ActiveEPList = data[16:]    //设备端口号列表，包含1个字节的端口号

	var key = constant.Constant.REDIS.ZigbeeRedisCtrlMsgKey + APMac + "_" + moduleID
	if Status == "00" {
		// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "procTerminalEndpointInfoUp success ")
		go publicfunction.IfMsgExistThenDelete(key, devEUI, SN, APMac)
		secondlyReadBasicByEndpoint(devEUI, ActiveEPCount, ActiveEPList)
	} else {
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procTerminalEndpointInfoUp failed ")
		redisData, _ := publicfunction.GetRedisDataFromRedis(key, devEUI, SN)
		if redisData != nil {
			publicfunction.CheckRedisAndReSend(devEUI, key, *redisData, APMac)
		}
	}
}

func procTerminalDiscoveryInfoUp(jsonInfo publicstruct.JSONInfo, terminalInfo config.TerminalInfo) {
	globallogger.Log.Infoln("devEUI : " + terminalInfo.DevEUI + " " + "procTerminalDiscoveryInfoUp, data: " + jsonInfo.MessagePayload.Data)
	var data = jsonInfo.MessagePayload.Data
	var devEUI = terminalInfo.DevEUI
	var APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
	var moduleID = jsonInfo.MessagePayload.ModuleID
	var SN = jsonInfo.TunnelHeader.FrameSN
	//var profileID = data[0:4]
	//var SrcAddr = data[4:8]         //发送消息设备的网络地址
	var Status = data[8:10] //SUCCESS(0) or FAILURE(1)
	//var NWKAddr = data[10:14]       //响应消息设备的网络地址
	//var Len = data[14:16]           //消息长度
	var Endpoint = data[16:18]  //设备端口号
	var ProfileID = data[18:22] //领域ID
	var DeviceID = data[22:26]  //设备类型ID
	//var DevVer = data[26:28]        //设备版本
	var NumInClusters = data[28:30] //InCluster数量
	tempValue, _ := strconv.ParseInt(NumInClusters, 16, 0)
	var numInClusters = int(tempValue)
	var InClusterList = data[30 : 30+4*numInClusters]                  //InCluster列表
	var NumOutClusters = data[30+4*numInClusters : 32+4*numInClusters] //OutCluster数量

	tempValue, _ = strconv.ParseInt(NumOutClusters, 16, 0)
	var numOutClusters = int(tempValue)
	var OutClusterList = data[32+4*numInClusters : 32+4*numInClusters+4*numOutClusters] //OutCluster列表

	var key = constant.Constant.REDIS.ZigbeeRedisCtrlMsgKey + APMac + "_" + moduleID
	if Status == "00" {
		// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "procTerminalDiscoveryInfoUp success ")
		_, err := lastlyUpdateBindInfo(devEUI, terminalInfo, ProfileID, DeviceID, numInClusters, InClusterList,
			numOutClusters, OutClusterList, Endpoint)
		if err == nil {
			go publicfunction.IfMsgExistThenDelete(key, devEUI, SN, APMac)
		}
	} else {
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procTerminalDiscoveryInfoUp failed ")
		redisData, _ := publicfunction.GetRedisDataFromRedis(key, devEUI, SN)
		if redisData != nil {
			publicfunction.CheckRedisAndReSend(devEUI, key, *redisData, APMac)
		}
	}
}

func procBindOrUnbindResult(devEUI string, msgType string, SrcEndpoint string, clusterID string, Status string, key string, APMac string,
	redisData publicstruct.RedisData, setData config.TerminalInfo, oSet map[string]interface{}, jsonInfo publicstruct.JSONInfo) {
	if Status == "00" {
		if msgType == globalmsgtype.MsgType.UPMsg.ZigbeeTerminalBindRspEvent {
			globallogger.Log.Infoln("devEUI : " + devEUI + " " + "procBindOrUnbindResult bind terminal success! endpoint: " +
				SrcEndpoint + " clusterID: " + clusterID)
			go publicfunction.IfMsgExistThenDelete(key, devEUI, redisData.SN, APMac)
		} else if msgType == globalmsgtype.MsgType.UPMsg.ZigbeeTerminalUnbindRspEvent {
			globallogger.Log.Infoln("devEUI : " + devEUI + " " + "procBindOrUnbindResult unbind terminal success! endpoint: " +
				SrcEndpoint + " clusterID: " + clusterID)
			var terminalInfoRes *config.TerminalInfo
			var err error
			if constant.Constant.UsePostgres {
				terminalInfoRes, err = models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, oSet)
			} else {
				terminalInfoRes, err = models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, setData)
			}
			if err != nil {
				globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procBindOrUnbindResult FindTerminalAndUpdate error : ", err)
			} else {
				globallogger.Log.Infoln("devEUI : "+devEUI+" "+"procBindOrUnbindResult FindTerminalAndUpdate success ", terminalInfoRes)
				go publicfunction.IfMsgExistThenDelete(key, devEUI, redisData.SN, APMac)
			}
		}
	} else if Status == "01" {
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procBindOrUnbindResult bind or unbind terminal failed! endpoint: " +
			SrcEndpoint + " clusterID: " + clusterID)
		go publicfunction.IfCtrlMsgExistThenDeleteAll(key, devEUI, APMac)
		go func() {
			select {
			case <-time.After(time.Duration(constant.Constant.TIMER.ZigbeeTimerAfterFive) * time.Second):
				publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
			}
		}()
	} else if Status == "8c" {
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procBindOrUnbindResult bind terminal failed! bind table is full! endpoint: " +
			SrcEndpoint + " clusterID: " + clusterID)
		go publicfunction.IfMsgExistThenDelete(key, devEUI, redisData.SN, APMac)
	}
}

func sensorSendWriteReq(terminalInfo config.TerminalInfo) {
	tmnInfo := config.TerminalInfo{
		DevEUI: terminalInfo.DevEUI,
	}
	cmd := common.Command{}
	cmd.DstEndpointIndex = 0
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      terminalInfo.DevEUI,
		T300ID:      terminalInfo.T300ID,
		CommandType: common.SensorWriteReq,
		ClusterID:   0x0500,
		Command:     cmd,
	}
	zclmain.ZclMain(globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent, tmnInfo, zclDownMsg, "", "", nil)
}

func procBindOrUnbindTerminalRsp(jsonInfo publicstruct.JSONInfo, terminalInfo config.TerminalInfo) {
	bindInfoTemp := []config.BindInfo{}
	if constant.Constant.UsePostgres {
		json.Unmarshal([]byte(terminalInfo.BindInfoPG), &bindInfoTemp)
	} else {
		bindInfoTemp = terminalInfo.BindInfo
	}
	globallogger.Log.Infoln("devEUI : " + terminalInfo.DevEUI + " " + "procBindOrUnbindTerminalRsp, data: " + jsonInfo.MessagePayload.Data)
	var msgType = jsonInfo.MessageHeader.MsgType
	var devEUI = terminalInfo.DevEUI
	var APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
	var moduleID = jsonInfo.MessagePayload.ModuleID
	var SN = jsonInfo.TunnelHeader.FrameSN
	var data = jsonInfo.MessagePayload.Data
	//var profileID = data[0:4]
	// var SrcAddr = data[4:8] //设备网络地址
	var Status = data[8:10] //SUCCESS(0) or FAILURE(1)

	var key = constant.Constant.REDIS.ZigbeeRedisCtrlMsgKey + APMac + "_" + moduleID
	redisData, _ := publicfunction.GetRedisDataFromRedis(key, devEUI, SN)
	if redisData != nil {
		var udpMsg, _ = hex.DecodeString(string(redisData.SendBuf))
		var JSONInfo = publicfunction.ParseUDPMsg(udpMsg, jsonInfo.Rinfo)
		var SrcEndpoint = JSONInfo.MessagePayload.Data[0:2]
		var clusterID = JSONInfo.MessagePayload.Data[2:6]
		if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 {
			SrcEndpoint = JSONInfo.MessagePayload.Data[16:18]
			clusterID = JSONInfo.MessagePayload.Data[18:22]
		}
		var setData = config.TerminalInfo{}
		oSet := make(map[string]interface{})
		if JSONInfo.MessagePayload.Data[len(JSONInfo.MessagePayload.Data)-2:] == "ff" { //绑定解绑T300回复处理
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
							(clusterID == item) && (Status == "00") &&
							(msgType == globalmsgtype.MsgType.UPMsg.ZigbeeTerminalBindRspEvent) {
							//只有最后一个clusterID绑定成功，才通知终端管理终端上线
							//传感器还需其他操作
							if terminalInfo.IsExist {
								if constant.Constant.Iotware {
									iotsmartspace.StateTerminalOnlineIotware(terminalInfo)
								} else if constant.Constant.Iotedge {
									iotsmartspace.StateTerminalOnline(devEUI)
								}
							} else {
								if constant.Constant.Iotedge {
									go iotsmartspace.ActionInsertReply(terminalInfo)
								}
							}
							if publicfunction.CheckTerminalIsSensor(devEUI, terminalInfo.TmnType) {
								go sensorSendWriteReq(terminalInfo)
							}
							publicfunction.TerminalOnline(devEUI, true)
							go func() {
								select {
								case <-time.After(10 * time.Second):
									iotsmartspace.ActionInsertSuccess(terminalInfo)
								}
							}()
						}
					}
					for i, v := range clusterIDArray {
						if v == "" {
							clusterIDArray = append(clusterIDArray[:i])
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
			procBindOrUnbindResult(devEUI, msgType, SrcEndpoint, clusterID, Status, key, APMac, *redisData, setData, oSet, jsonInfo)
		} else { //设备与设备绑定解绑回复
			publicfunction.IfMsgExistThenDelete(key, devEUI, redisData.SN, APMac)
		}
	}
}

func procTerminalLeaveRsp(jsonInfo publicstruct.JSONInfo, terminalInfo *config.TerminalInfo) {
	globallogger.Log.Infoln("devEUI : " + terminalInfo.DevEUI + " " + "procTerminalLeaveRsp, data: " + jsonInfo.MessagePayload.Data)
	var devEUI = terminalInfo.DevEUI
	var data = jsonInfo.MessagePayload.Data
	var APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
	var moduleID = jsonInfo.MessagePayload.ModuleID
	var SN = jsonInfo.TunnelHeader.FrameSN
	//var profileID = data[0:4]
	//var SrcAddr = data[4:8] //设备网络地址
	var Status = data[8:10] //SUCCESS(0) or FAILURE(1)
	var key = constant.Constant.REDIS.ZigbeeRedisCtrlMsgKey + APMac + "_" + moduleID

	if Status == "00" { //success
		// globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procTerminalLeaveRsp success ")
		if terminalInfo != nil {
			publicfunction.TerminalOffline(devEUI)
			if constant.Constant.UsePostgres {
				_, err := models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, map[string]interface{}{"leavestate": true})
				if err != nil {
					globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procTerminalLeaveRsp FindTerminalAndUpdate err : ", err)
				}
			} else {
				_, err := models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, bson.M{"leaveState": true})
				if err != nil {
					globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procTerminalLeaveRsp FindTerminalAndUpdate err : ", err)
				}
			}
			if constant.Constant.Iotware {
				iotsmartspace.StateTerminalLeaveIotware(*terminalInfo)
			} else if constant.Constant.Iotedge {
				iotsmartspace.StateTerminalLeave(devEUI)
			}
		}
		publicfunction.IfCtrlMsgExistThenDeleteAll(key, devEUI, APMac)
		key = constant.Constant.REDIS.ZigbeeRedisDataMsgKey + APMac + "_" + moduleID
		publicfunction.IfDataMsgExistThenDeleteAll(key, devEUI, APMac)
	} else {
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procTerminalLeaveRsp failed ")
		redisData, _ := publicfunction.GetRedisDataFromRedis(key, devEUI, SN)
		if redisData != nil {
			publicfunction.CheckRedisAndReSend(devEUI, key, *redisData, APMac)
		}
	}
}

func procTerminalPermitJoinRsp(jsonInfo publicstruct.JSONInfo, terminalInfo config.TerminalInfo) {
	globallogger.Log.Infoln("devEUI : " + terminalInfo.DevEUI + " " + "procTerminalPermitJoinRsp, data: " + jsonInfo.MessagePayload.Data)
	var devEUI = terminalInfo.DevEUI
	var data = jsonInfo.MessagePayload.Data
	var APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
	var moduleID = jsonInfo.MessagePayload.ModuleID
	var SN = jsonInfo.TunnelHeader.FrameSN
	//var profileID = data[0:4]
	//var SrcAddr = data[4:8] //设备网络地址
	var Status = data[8:10] //SUCCESS(0) or FAILURE(1)
	var key = constant.Constant.REDIS.ZigbeeRedisCtrlMsgKey + APMac + "_" + moduleID

	redisData, _ := publicfunction.GetRedisDataFromRedis(key, devEUI, SN)
	if redisData != nil {
		if Status == "00" { //success
			var setData = bson.M{}
			setData["PermitJoin"] = false //不允许加入成功
			if redisData.Data[0:2] == "ff" {
				setData["PermitJoin"] = true //允许加入成功
			}
			if constant.Constant.UsePostgres {
				_, err := models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI},
					map[string]interface{}{"permitjoin": setData["PermitJoin"]})
				if err != nil {
					globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procTerminalPermitJoinRsp FindTerminalAndUpdate error : ", err)
				}
			} else {
				_, err := models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, setData)
				if err != nil {
					globallogger.Log.Errorln("devEUI : "+devEUI+" "+"procTerminalPermitJoinRsp FindTerminalAndUpdate error : ", err)
				}
			}
			// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"procTerminalPermitJoinRsp FindTerminalAndUpdate success ", terminalInfoRes)
			go publicfunction.IfMsgExistThenDelete(key, devEUI, SN, APMac)
		} else {
			globallogger.Log.Warnln("devEUI : " + devEUI + " " + "procTerminalPermitJoinRsp failed ")
			publicfunction.CheckRedisAndReSend(devEUI, key, *redisData, APMac)
		}
	}
}
