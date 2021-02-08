package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/db/kafka"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmemorycache"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalrediscache"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
)

var zigbeeServerColliHighTimerID = &sync.Map{}

func collihigh4To20mAParse(min int, max int, value float64) float64 {
	return float64(max-min) / (20 - 4) * (value - 4)
}

func collihighSensorParsePM25(devEUI string, c0 string, c0Type string, c0Data string) float64 {
	var value float64 = 0
	if c0 != "c0" {
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "collihighSensorParsePM25 c0 is not match : " + c0)
	} else {
		switch c0Type {
		case "03":
			tempVaule, _ := strconv.ParseInt(c0Data, 16, 0)
			value = float64(tempVaule)
			value /= 1000
			if value < 4 || value > 20 {
				value = 0
			} else {
				value = collihigh4To20mAParse(0, 999, value)
			}
			globallogger.Log.Infoln("devEUI : "+devEUI+" "+"collihighSensorParsePM25 value: ", value)
		default:
			globallogger.Log.Warnln("devEUI : " + devEUI + " " + "collihighSensorParsePM25 c0Type is unknow : " + c0Type)
		}
	}
	return value
}

func collihighSensorParseLighting(devEUI string, c1 string, c1Type string, c1Data string) float64 {
	var value float64 = 0
	if c1 != "c1" {
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "collihighSensorParseLighting c1 is not match : " + c1)
	} else {
		switch c1Type {
		case "03":
			tempVaule, _ := strconv.ParseInt(c1Data, 16, 0)
			value = float64(tempVaule)
			value /= 1000
			if value < 4 || value > 20 {
				value = 0
			} else {
				value = collihigh4To20mAParse(0, 20000, value)
			}
			globallogger.Log.Infoln("devEUI : "+devEUI+" "+"collihighSensorParseLighting value: ", value)
		default:
			globallogger.Log.Warnln("devEUI : " + devEUI + " " + "collihighSensorParseLighting c1Type is unknow : " + c1Type)
		}
	}
	return value
}

//SendorAttr SendorAttr
type SendorAttr struct {
	Attribute string  `json:"attrbute"`
	Data      float64 `json:"data"`
}

func collihighSensorParse(devEUI string, Data string) []SendorAttr {
	var c0 = Data[26:28]
	var c0Type = Data[28:30]
	var c0Data = Data[30:34]
	var c1 = Data[34:36]
	var c1Type = Data[36:38]
	var c1Data = Data[38:42]
	//var c2 = Data[42:44]
	//var c2Type = Data[44:46]
	//var c2Data = Data[46:50]
	//var c3 = Data[50:52]
	//var c3Type = Data[52:54]
	//var c3Data = Data[54:58]
	var c0Value float64 = 0
	var c1Value float64 = 0
	var c2Value float64 = 0
	var c3Value float64 = 0

	c0Value = collihighSensorParsePM25(devEUI, c0, c0Type, c0Data)
	c1Value = collihighSensorParseLighting(devEUI, c1, c1Type, c1Data)
	return []SendorAttr{
		{
			Attribute: "51",
			Data:      c0Value,
		},
		{
			Attribute: "03",
			Data:      c1Value,
		},
		{
			Attribute: "00",
			Data:      c2Value,
		},
		{
			Attribute: "00",
			Data:      c3Value,
		},
	}
}

func collihighSensorDataProc(devEUI string, data string, terminalInfo config.TerminalInfo) {
	globallogger.Log.Infoln("devEUI : " + devEUI + " " + "collihighSensorDataProc")
	// var profileID = data[0:4]
	// var GroupID = data[4:8]
	// var ClusterID = data[8:12]
	// var SrcAddr = data[12:16]
	// var SrcEndpoint = data[16:18]
	// var DestEndpoint = data[18:20]
	// var WasBroadcast = data[20:22]
	// var LinkQuality = data[22:24]
	// var SecurityUse = data[24:26]
	// var Timestamp = data[26:34]
	// var TransSeqNumber = data[34:36]
	// var Len = data[36:38]
	var Data = data[38:]

	var sendMsg = struct {
		TmnName   string      `json:"tmnName"`
		DevEUI    string      `json:"devEUI"`
		ShopID    string      `json:"shopId"`
		OIDIndex  string      `json:"OIDIndex"`
		FirmTopic string      `json:"firmTopic"`
		TmnType   string      `json:"tmnType"`
		Time      time.Time   `json:"time"`
		Data      interface{} `json:"data"`
	}{}
	sendMsg.TmnName = terminalInfo.TmnName
	sendMsg.DevEUI = devEUI
	sendMsg.ShopID = terminalInfo.ScenarioID
	sendMsg.OIDIndex = terminalInfo.OIDIndex
	sendMsg.FirmTopic = terminalInfo.FirmTopic
	sendMsg.TmnType = terminalInfo.TmnType
	sendMsg.Time = time.Now() //new Date().getTime();
	if terminalInfo.TmnType == constant.Constant.TMNTYPE.COLLIHIGH.ZigbeeTerminalLightingPM25 {
		sendMsg.Data = collihighSensorParse(devEUI, Data)
	} else {
		sendMsg.Data = Data
	}

	globallogger.Log.Infoln("devEUI : " + devEUI + " " + "collihighSensorDataProc sendMsg: " + fmt.Sprintf("%+v", sendMsg))
	sendMsgSerial, _ := json.Marshal(sendMsg)
	//mqhd.SendMsg("iotenvmonitorns", sendMsgSerial)
	kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicCollihighSensorDataUp, string(sendMsgSerial))
	if constant.Constant.UsePostgres {
		models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, map[string]interface{}{"attributedata": Data})
	} else {
		models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, config.TerminalInfo{Attribute: config.Attribute{Data: Data}})
	}
}

func collihighTerminalStatusProc(devEUI string, t int) {
	if value, ok := zigbeeServerColliHighTimerID.Load(devEUI); ok {
		value.(*time.Timer).Stop()
		globallogger.Log.Infoln("devEUI: "+devEUI+" collihighTerminalStatusProc clearTimeout, time: ", t)
	}
	//3次未收到数据报文，通知离线
	var timerID = time.NewTimer(time.Duration(3*t) * time.Second)
	zigbeeServerColliHighTimerID.Store(devEUI, timerID)
	go func() {
		select {
		case <-timerID.C:
			if value, ok := zigbeeServerColliHighTimerID.Load(devEUI); ok {
				if timerID == value.(*time.Timer) {
					zigbeeServerColliHighTimerID.Delete(devEUI)
					terminalTimerInfo, _ := publicfunction.GetTerminalTimerByDevEUI(devEUI)
					if terminalTimerInfo != nil {
						var dateNow = time.Now()
						var num = dateNow.UnixNano() - terminalTimerInfo.UpdateTime.UnixNano()
						var terminalInfo *config.TerminalInfo
						var err error
						if constant.Constant.UsePostgres {
							terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
						} else {
							terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
						}
						if terminalInfo != nil && terminalInfo.Interval != 0 {
							t = terminalInfo.Interval
						}
						if num > int64(time.Duration(3*t)*time.Second) {
							globallogger.Log.Infoln("devEUI: "+devEUI+" collihighTerminalStatusProc server has not already recv collihigh data for ",
								num/int64(time.Second), " seconds. Then terminal offline")
							publicfunction.TerminalOffline(devEUI)
							if err == nil && terminalInfo != nil {
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
							publicfunction.DeleteTerminalTimer(devEUI)
						}
					}
				}
			}
		}
	}()
}

func findCollihighTerminalAndUpdate(jsonInfo publicstruct.JSONInfo) (*config.TerminalInfo, error) {
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		oMatch := make(map[string]interface{})
		oSet := make(map[string]interface{})
		oMatch["deveui"] = jsonInfo.MessagePayload.Address
		if jsonInfo.TunnelHeader.LinkInfo.ACMac != "" {
			oSet["acmac"] = jsonInfo.TunnelHeader.LinkInfo.ACMac
		}
		if jsonInfo.TunnelHeader.LinkInfo.T300ID != "" {
			oSet["t300id"] = jsonInfo.TunnelHeader.LinkInfo.T300ID
		}
		oSet["apmac"] = jsonInfo.TunnelHeader.LinkInfo.APMac
		oSet["moduleid"] = jsonInfo.MessagePayload.ModuleID
		oSet["updatetime"] = time.Now()
		terminalInfo, err = models.FindTerminalAndUpdatePG(oMatch, oSet)
	} else {
		var setData = config.TerminalInfo{}
		if jsonInfo.TunnelHeader.LinkInfo.ACMac != "" {
			setData.ACMac = jsonInfo.TunnelHeader.LinkInfo.ACMac
		}
		if jsonInfo.TunnelHeader.LinkInfo.T300ID != "" {
			setData.T300ID = jsonInfo.TunnelHeader.LinkInfo.T300ID
		}
		setData.APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
		setData.ModuleID = jsonInfo.MessagePayload.ModuleID
		setData.UpdateTime = time.Now()
		terminalInfo, err = models.FindTerminalAndUpdate(bson.M{"devEUI": jsonInfo.MessagePayload.Address}, setData)
	}
	if err != nil {
		globallogger.Log.Errorln("devEUI : "+jsonInfo.MessagePayload.Address+" "+"findCollihighTerminalAndUpdate error : ", err)
	} else {
		if terminalInfo == nil {
			globallogger.Log.Warnln("devEUI : " + jsonInfo.MessagePayload.Address + " " +
				"findCollihighTerminalAndUpdate terminal is not exist, please first create !")
		} else {
			// globallogger.Log.Infoln("devEUI : " + jsonInfo.MessagePayload.Address + " " + "findCollihighTerminalAndUpdate success")
		}
	}
	return terminalInfo, err
}

func collihighTerminalCreate(devEUI string, jsonInfo publicstruct.JSONInfo, terminalInfo config.DevMgrInfo) (*config.TerminalInfo, error) {
	var setData = config.TerminalInfo{}
	var info = publicstruct.TmnTypeAttr{}
	info = publicfunction.GetTmnTypeAndAttribute(terminalInfo.TmnType)
	setData.TmnName = terminalInfo.TmnName
	setData.DevEUI = terminalInfo.TmnDevSN
	setData.OIDIndex = terminalInfo.TmnOIDIndex
	setData.ScenarioID = terminalInfo.SceneID
	setData.UserName = terminalInfo.TenantID
	setData.FirmTopic = terminalInfo.FirmTopic
	setData.ProfileID = terminalInfo.LinkType
	setData.TmnType = info.TmnType
	setData.Attribute = info.Attribute
	setData.Online = true
	if jsonInfo.TunnelHeader.LinkInfo.ACMac != "" {
		setData.ACMac = jsonInfo.TunnelHeader.LinkInfo.ACMac
		if len(jsonInfo.TunnelHeader.LinkInfo.ACMac) == 40 {
			setData.FirstAddr = publicfunction.Transport16StringToString(jsonInfo.TunnelHeader.LinkInfo.ACMac)
		} else {
			setData.FirstAddr = jsonInfo.TunnelHeader.LinkInfo.ACMac
		}
	}
	setData.APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
	if len(jsonInfo.TunnelHeader.LinkInfo.APMac) == 40 {
		setData.SecondAddr = publicfunction.Transport16StringToString(jsonInfo.TunnelHeader.LinkInfo.APMac)
	} else {
		setData.SecondAddr = jsonInfo.TunnelHeader.LinkInfo.APMac
	}
	if jsonInfo.TunnelHeader.LinkInfo.T300ID != "" {
		setData.T300ID = jsonInfo.TunnelHeader.LinkInfo.T300ID
		setData.ThirdAddr = jsonInfo.TunnelHeader.LinkInfo.T300ID
	}
	setData.ModuleID = jsonInfo.MessagePayload.ModuleID
	setData.UpdateTime = time.Now()
	var err error
	if constant.Constant.UsePostgres {
		err = models.CreateTerminalPG(setData)
	} else {
		err = models.CreateTerminal(setData)
	}
	if err != nil {
		globallogger.Log.Errorln("devEUI : "+devEUI+" "+"collihighTerminalCreate createTerminal error : ", err)
	} else {
		// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"collihighTerminalCreate createTerminal success ", setData)
	}
	//var interval = 60
	if terminalInfo.Property != nil {
		if len(terminalInfo.Property) > 0 {
			//tempVaule, _ := strconv.ParseInt(terminalInfo.property[0].interval, 10, 0) //string转10进制数值
			//interval = int(tempVaule)
			if constant.Constant.MultipleInstances {
				globalrediscache.RedisCache.SetRedisSet(constant.Constant.REDIS.ZigbeeRedisCollihighInterval+devEUI, terminalInfo.Property[0].Interval)
			} else {
				globalmemorycache.MemoryCache.SetMemorySet(constant.Constant.REDIS.ZigbeeRedisCollihighInterval+devEUI, terminalInfo.Property[0].Interval)
			}
		}
	}
	return &setData, err
}

func collihighTerminalUpdate(jsonInfo publicstruct.JSONInfo, devEUI string) (*config.TerminalInfo, string) {
	var errorCode string = ""
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		oMatch := make(map[string]interface{})
		oSet := make(map[string]interface{})
		oMatch["deveui"] = jsonInfo.MessagePayload.Address
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
		terminalInfo, err = models.FindTerminalAndUpdatePG(oMatch, oSet)
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
		terminalInfo, err = models.FindTerminalAndUpdate(bson.M{"devEUI": jsonInfo.MessagePayload.Address}, setData)
	}
	if err != nil {
		globallogger.Log.Errorln("devEUI : "+jsonInfo.MessagePayload.Address+" "+"collihighTerminalUpdate FindTerminalAndUpdate error : ", err)
	} else {
		if terminalInfo == nil {
			globallogger.Log.Warnln("devEUI : " + jsonInfo.MessagePayload.Address + " " +
				"collihighTerminalUpdate terminal is not exist, please first create !")
		} else {
			// globallogger.Log.Infoln("devEUI : " + jsonInfo.MessagePayload.Address + " " + "collihighTerminalUpdate FindTerminalAndUpdate success")
		}
	}
	if err == nil {
		if terminalInfo != nil {
			if publicfunction.CheckTerminalSceneIsMatch(devEUI, terminalInfo.ScenarioID, jsonInfo.TunnelHeader.VenderInfo) {
				publicfunction.TerminalOnline(devEUI, false)
				if constant.Constant.Iotware {
					if terminalInfo.IsExist {
						iotsmartspace.StateTerminalOnlineIotware(*terminalInfo)
					}
				} else if constant.Constant.Iotedge {
					iotsmartspace.StateTerminalOnline(devEUI)
				}
			} else {
				errorCode = "scene is not match"
			}
		} else {
			httpTerminalInfo := publicfunction.HTTPRequestTerminalInfo(devEUI, "klha")
			if httpTerminalInfo != nil {
				if publicfunction.CheckTerminalSceneIsMatch(devEUI, httpTerminalInfo.SceneID, jsonInfo.TunnelHeader.VenderInfo) {
					terminalInfo, err = findCollihighTerminalAndUpdate(jsonInfo)
					if err == nil {
						if terminalInfo != nil {
							publicfunction.TerminalOnline(devEUI, false)
							if constant.Constant.Iotware {
								if terminalInfo.IsExist {
									iotsmartspace.StateTerminalOnlineIotware(*terminalInfo)
								}
							} else if constant.Constant.Iotedge {
								iotsmartspace.StateTerminalOnline(devEUI)
							}
						} else {
							terminal, _ := collihighTerminalCreate(devEUI, jsonInfo, *httpTerminalInfo)
							publicfunction.TerminalOnline(devEUI, false)
							if constant.Constant.Iotware {
								if terminalInfo.IsExist {
									iotsmartspace.StateTerminalOnlineIotware(*terminal)
								}
							} else if constant.Constant.Iotedge {
								iotsmartspace.StateTerminalOnline(devEUI)
							}
							terminalInfo = terminal
						}
					} else {
						errorCode = err.Error()
					}
				} else {
					errorCode = "scene is not match"
				}
			}
		}
	} else {
		errorCode = err.Error()
	}
	return terminalInfo, errorCode
}

func collihighProc(jsonInfo publicstruct.JSONInfo) {
	globallogger.Log.Infoln("devEUI : " + jsonInfo.MessagePayload.Address + " " + "collihighProc")
	var devEUI = jsonInfo.MessagePayload.Address
	var msgType = jsonInfo.MessageHeader.MsgType
	var data = jsonInfo.MessagePayload.Data
	terminalInfo, errorCode := collihighTerminalUpdate(jsonInfo, devEUI)
	if errorCode == "" && terminalInfo != nil {
		if msgType == globalmsgtype.MsgType.UPMsg.ZigbeeDataUpEvent {
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.COLLIHIGH.ZigbeeTerminalCO2,
				constant.Constant.TMNTYPE.COLLIHIGH.ZigbeeTerminalLightingPM25,
				constant.Constant.TMNTYPE.COLLIHIGH.ZigbeeTerminalHumitureLightingCO2:
				var interval = 60
				var value string
				var err error
				if constant.Constant.MultipleInstances {
					value, err = globalrediscache.RedisCache.GetRedisGet(constant.Constant.REDIS.ZigbeeRedisCollihighInterval + devEUI)
				} else {
					value, err = globalmemorycache.MemoryCache.GetMemoryGet(constant.Constant.REDIS.ZigbeeRedisCollihighInterval + devEUI)
				}
				if err == nil && value != "" {
					tempVaule, _ := strconv.ParseInt(value, 10, 0) //string转10进制数值
					interval = int(tempVaule)
				}
				err = publicfunction.TerminalTimerUpdateOrCreate(devEUI)
				if err == nil {
					collihighTerminalStatusProc(devEUI, interval)
				}
				collihighSensorDataProc(devEUI, data, *terminalInfo)
			default:
				globallogger.Log.Warnln("devEUI : " + devEUI + " " + "unknow tmnType : " + terminalInfo.TmnType)
			}
		}
	}
}
