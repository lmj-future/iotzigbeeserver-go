package publicfunction

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/coocood/freecache"
	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmemorycache"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalrediscache"
	"github.com/h3c/iotzigbeeserver-go/httprequest"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	"github.com/h3c/iotzigbeeserver-go/rabbitmqmsg/producer"
)

const (
	iotterminalmgrServiceHost     = "IOTTERMINALMGR_SERVICE_HOST"
	iotterminalmgrServicePortHTTP = "IOTTERMINALMGR_SERVICE_PORT_HTTP"
	iotfirmmgrServiceHost         = "IOTFIRMMGR_SERVICE_HOST"
	iotfirmmgrServicePortHTTP     = "IOTFIRMMGR_SERVICE_PORT_HTTP"
)

var zigbeeServerSocketCache = freecache.NewCache(1024 * 1024)

func GetCtrlMsgKey(APMac string, moduleID string) string {
	var ctrlKeyBuilder strings.Builder
	ctrlKeyBuilder.WriteString(constant.Constant.REDIS.ZigbeeRedisCtrlMsgKey)
	ctrlKeyBuilder.WriteString(APMac)
	ctrlKeyBuilder.WriteString("_")
	ctrlKeyBuilder.WriteString(moduleID)
	return ctrlKeyBuilder.String()
}

func GetDataMsgKey(APMac string, moduleID string) string {
	var dataKeyBuilder strings.Builder
	dataKeyBuilder.WriteString(constant.Constant.REDIS.ZigbeeRedisDataMsgKey)
	dataKeyBuilder.WriteString(APMac)
	dataKeyBuilder.WriteString("_")
	dataKeyBuilder.WriteString(moduleID)
	return dataKeyBuilder.String()
}

// Transport16StringToString Transport16StringToString
func Transport16StringToString(str string) string {
	var valueBuilder strings.Builder
	for len(str) != 0 {
		value, _ := strconv.ParseUint(strings.Repeat(str[:2], 1), 16, 8)
		valueBuilder.WriteString(string(rune(value)))
		str = strings.Repeat(str[2:], 1)
	}
	return valueBuilder.String()
}

// IfMsgExistThenDelete IfMsgExistThenDelete
func IfMsgExistThenDelete(key string, devEUI string, SN string, APMac string) {
	var isCtrlMsg = false
	if strings.Repeat(key[0:20], 1) == "zigbee_down_ctrl_msg" {
		isCtrlMsg = true
	}
	_, redisArray, _ := GetRedisLengthAndRangeRedis(key, devEUI)
	if len(redisArray) != 0 {
		var isNeedDown = false
		for itemIndex, item := range redisArray {
			if item != constant.Constant.REDIS.ZigbeeRedisRemoveRedis {
				var itemData publicstruct.RedisData
				json.Unmarshal([]byte(item), &itemData)
				if SN == itemData.SN {
					if constant.Constant.MultipleInstances {
						_, err := globalrediscache.RedisCache.SetRedis(key, itemIndex, []byte(constant.Constant.REDIS.ZigbeeRedisRemoveRedis))
						if err != nil {
							globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "IfMsgExistThenDelete set redis error :", err)
						}
					} else {
						_, err := globalmemorycache.MemoryCache.SetMemory(key, itemIndex, []byte(constant.Constant.REDIS.ZigbeeRedisRemoveRedis))
						if err != nil {
							globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "IfMsgExistThenDelete set memory error :", err)
						}
					}
					if isCtrlMsg {
						isNeedDown = true
					} else {
						if !itemData.AckFlag {
							isNeedDown = true
						}
					}
				} else {
					if isNeedDown {
						if isCtrlMsg {
							CheckRedisAndSend(itemData.DevEUI, key, itemData.SN, APMac)
						} else {
							if !itemData.AckFlag {
								CheckRedisAndSend(itemData.DevEUI, key, itemData.SN, APMac)
							}
						}
						isNeedDown = false
					}
				}
			}
		}
	}
}

// IfCtrlMsgExistThenDeleteAll IfCtrlMsgExistThenDeleteAll
func IfCtrlMsgExistThenDeleteAll(key string, devEUI string, APMac string) {
	_, redisArray, _ := GetRedisLengthAndRangeRedis(key, devEUI)
	if len(redisArray) != 0 {
		var isNeedDown = true
		var tempDevEUI = ""
		var tempSN = ""

		for index, item := range redisArray {
			if item != constant.Constant.REDIS.ZigbeeRedisRemoveRedis {
				var redisData struct {
					DevEUI string `json:"devEUI"`
					SN     string `json:"sn"`
				}
				json.Unmarshal([]byte(item), &redisData)
				if redisData.DevEUI == devEUI {
					if constant.Constant.MultipleInstances {
						_, err := globalrediscache.RedisCache.SetRedis(key, index, []byte(constant.Constant.REDIS.ZigbeeRedisRemoveRedis))
						if err != nil {
							globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "IfCtrlMsgExistThenDeleteAll set redis error :", err)
						}
					} else {
						_, err := globalmemorycache.MemoryCache.SetMemory(key, index, []byte(constant.Constant.REDIS.ZigbeeRedisRemoveRedis))
						if err != nil {
							globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "IfCtrlMsgExistThenDeleteAll set memory error :", err)
						}
					}
				} else {
					if isNeedDown {
						tempDevEUI = redisData.DevEUI
						tempSN = redisData.SN
					}
				}
				if isNeedDown && tempDevEUI != "" {
					CheckRedisAndSend(tempDevEUI, key, tempSN, APMac)
					isNeedDown = false
					tempDevEUI = ""
				}
			}
		}
	}
}

// IfDataMsgExistThenDeleteAll IfDataMsgExistThenDeleteAll
func IfDataMsgExistThenDeleteAll(key string, devEUI string, APMac string) {
	_, redisArray, _ := GetRedisLengthAndRangeRedis(key, devEUI)
	if len(redisArray) != 0 {
		var isSend = false
		for index, item := range redisArray {
			if item != constant.Constant.REDIS.ZigbeeRedisRemoveRedis {
				var redisData struct {
					DevEUI string `json:"devEUI"`
				}
				json.Unmarshal([]byte(item), &redisData)
				if redisData.DevEUI == devEUI {
					var err error
					if constant.Constant.MultipleInstances {
						_, err = globalrediscache.RedisCache.SetRedis(key, index, []byte(constant.Constant.REDIS.ZigbeeRedisRemoveRedis))
					} else {
						_, err = globalmemorycache.MemoryCache.SetMemory(key, index, []byte(constant.Constant.REDIS.ZigbeeRedisRemoveRedis))
					}
					if err != nil {
						globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "IfDataMsgExistThenDeleteAll set redis error :", err)
					} else {
						if !isSend {
							var tempIndex = index + 1
							for {
								if tempIndex >= len(redisArray) {
									break
								}
								if redisArray[tempIndex] != constant.Constant.REDIS.ZigbeeRedisRemoveRedis {
									var nextData struct {
										devEUI  string
										ackFlag bool
										SN      string
									}

									json.Unmarshal([]byte(redisArray[tempIndex]), &nextData)
									if (nextData.devEUI != devEUI) && (!nextData.ackFlag) {
										CheckRedisAndSend(nextData.devEUI, key, nextData.SN, APMac)
										isSend = true
									}
									tempIndex = len(redisArray)
								} else {
									tempIndex++
								}
							}
						}
					}
				}
			}
		}
	}
}

// TerminalOnline TerminalOnline
func TerminalOnline(devEUI string, isSendConfig bool) {
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI},
			map[string]interface{}{"leavestate": false, "online": true})
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "TerminalOnline FindTerminalAndUpdatePG err :", err)
		}
	} else {
		terminalInfo, err = models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, bson.M{"online": true, "leaveState": false})
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "TerminalOnline FindTerminalAndUpdate err :", err)
		}
	}
	if terminalInfo != nil {
		var key = terminalInfo.APMac + terminalInfo.ModuleID + terminalInfo.DevEUI
		if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
			key = terminalInfo.APMac + terminalInfo.ModuleID + terminalInfo.NwkAddr
		}
		DeleteTerminalInfoListCache(key)
	}
	if constant.Constant.UseRabbitMQ {
		var terminalStateInfo = producer.TerminalState{
			TmnDevSN:     devEUI,
			Time:         time.Now(),
			LinkType:     strings.ToLower(terminalInfo.ProfileID),
			IndexType:    "deveui",
			Index:        devEUI,
			TmnOIDIndex:  terminalInfo.OIDIndex,
			ExtType:      "terminalOnline",
			IsSendConfig: isSendConfig,
		}
		globallogger.Log.Warnf("devEUI : %+v TerminalOnline send msg to iotterminalmgr: %+v", devEUI, terminalStateInfo)
		producer.SendTerminalStateMsg(&terminalStateInfo, "iotterminalmgr")
	}
}

// TerminalOffline TerminalOffline
func TerminalOffline(devEUI string) {
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, map[string]interface{}{"online": false})
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "TerminalOffline FindTerminalAndUpdatePG err :", err)
		}
	} else {
		terminalInfo, err = models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, bson.M{"online": false})
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "TerminalOffline FindTerminalAndUpdate err :", err)
		}
	}
	if terminalInfo != nil {
		var key = terminalInfo.APMac + terminalInfo.ModuleID + terminalInfo.DevEUI
		if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
			key = terminalInfo.APMac + terminalInfo.ModuleID + terminalInfo.NwkAddr
		}
		DeleteTerminalInfoListCache(key)
	}
	if constant.Constant.UseRabbitMQ {
		var terminalStateInfo = producer.TerminalState{
			TmnDevSN:    devEUI,
			Time:        time.Now(),
			LinkType:    strings.ToLower(terminalInfo.ProfileID),
			IndexType:   "deveui",
			Index:       devEUI,
			TmnOIDIndex: terminalInfo.OIDIndex,
			ExtType:     "terminalOffline",
		}
		globallogger.Log.Warnf("devEUI : %+v TerminalOffline send msg to iotterminalmgr: %+v", devEUI, terminalStateInfo)
		producer.SendTerminalStateMsg(&terminalStateInfo, "iotterminalmgr")
	}
}

// ProcKeepAliveTerminalOffline ProcKeepAliveTerminalOffline
func ProcKeepAliveTerminalOffline(APMac string) ([]config.TerminalInfo, error) {
	if constant.Constant.UsePostgres {
		terminalList, err := models.FindTerminalByAPMacPG(APMac)
		if err != nil {
			globallogger.Log.Errorln("ProcKeepAliveTerminalOffline FindTerminalByAPMacPG err:", err)
		} else {
			if len(terminalList) > 0 {
				for _, item := range terminalList {
					TerminalOffline(item.DevEUI)
				}
			}
		}
		err = models.DeleteSocketPG(APMac)
		if err != nil {
			globallogger.Log.Errorln("ProcKeepAliveTerminalOffline DeleteSocketPG err:", err)
		}
		return terminalList, err
	}
	terminalList, err := models.FindTerminalByAPMac(APMac)
	if err != nil {
		globallogger.Log.Errorln("ProcKeepAliveTerminalOffline findTerminalByAPMac err:", err)
	} else {
		if len(terminalList) > 0 {
			for _, item := range terminalList {
				TerminalOffline(item.DevEUI)
			}
		}
	}
	err = models.DeleteSocket(APMac)
	if err != nil {
		globallogger.Log.Errorln("ProcKeepAliveTerminalOffline DeleteSocket err:", err)
	}
	return terminalList, err
}

// HTTPRequestTerminalTypeByAlias HTTPRequestTerminalTypeByAlias
func HTTPRequestTerminalTypeByAlias(alias string) *httprequest.Data4 {
	host := os.Getenv(iotfirmmgrServiceHost)
	port := os.Getenv(iotfirmmgrServicePortHTTP)
	var url = "http://" + host + ":" + port + "/iot/iotfirmmgr/findTerminalTypeByAlias?alias=" + alias
	data, err := httprequest.HTTPRequest(url, "GET", "terminalTypeByAlias", []byte{})
	httpData := data.(*httprequest.HTTPData4)
	if err != nil {
		globallogger.Log.Errorln("HTTPRequestTerminalTypeByAlias httprequest.HttpRequest error :", err)
	} else {
		if httpData != nil {
			globallogger.Log.Infof("HTTPRequestTerminalTypeByAlias httprequest.HttpRequest success , response: %+v", *httpData)
			if httpData.ErrCode == 0 {
				globallogger.Log.Infof("HTTPRequestTerminalTypeByAlias terminalType is exist, terminalType info: %+v", httpData.Data)
				return &httpData.Data
			}
		}
	}
	globallogger.Log.Warnln("HTTPRequestTerminalTypeByAlias terminalType is not exist")
	return nil
}

// HTTPRequestAddTerminal HTTPRequestAddTerminal
func HTTPRequestAddTerminal(devEUI string, terminalInfo []byte) bool {
	host := os.Getenv(iotterminalmgrServiceHost)
	port := os.Getenv(iotterminalmgrServicePortHTTP)
	var url = "http://" + host + ":" + port + "/iot/iotterminalmgr/addTerminalMulti"
	data, err := httprequest.HTTPRequest(url, "POST", "addTerminal", terminalInfo)
	httpData := data.(*httprequest.HTTPData3)
	if err != nil {
		globallogger.Log.Errorln("devEUI :", devEUI, "HTTPRequestAddTerminal httprequest.HttpRequest error :", err)
	} else {
		if httpData != nil {
			globallogger.Log.Infof("devEUI : %s HTTPRequestAddTerminal httprequest.HttpRequest success , response: %+v", devEUI, *httpData)
			if httpData.ErrCode == 0 {
				globallogger.Log.Infof("devEUI : %s HTTPRequestAddTerminal success, terminal info: %+v", devEUI, httpData.Data)
				return true
			}
		}
	}
	globallogger.Log.Warnln("devEUI :", devEUI, "HTTPRequestAddTerminal failed")
	return false
}

// HTTPRequestTerminalInfo HTTPRequestTerminalInfo
func HTTPRequestTerminalInfo(devEUI string, linkType string) *config.DevMgrInfo {
	host := os.Getenv(iotterminalmgrServiceHost)
	port := os.Getenv(iotterminalmgrServicePortHTTP)
	var url = "http://" + host + ":" + port + "/iot/iotterminalmgr/findOneTerminal?tmnDevSN=" + devEUI + "&linkType=" + linkType
	data, err := httprequest.HTTPRequest(url, "GET", "terminalInfoReq", []byte{})
	httpData := data.(*httprequest.HTTPData1)
	if err != nil {
		globallogger.Log.Errorln("devEUI :", devEUI, "HTTPRequestTerminalInfo httprequest.HttpRequest error :", err)
	} else {
		if httpData != nil {
			globallogger.Log.Infof("devEUI : %s HTTPRequestTerminalInfo httprequest.HttpRequest success , response: %+v", devEUI, *httpData)
			if httpData.ErrCode == 0 {
				globallogger.Log.Infof("devEUI : %s HTTPRequestTerminalInfo terminal is exist, terminal info: %+v", devEUI, *httpData.Data)
				return httpData.Data
			}
		}
	}
	globallogger.Log.Warnln("devEUI :", devEUI, "HTTPRequestTerminalInfo terminal is not exist")
	return nil
}

// HTTPRequestCheckLicense HTTPRequestCheckLicense
func HTTPRequestCheckLicense(devEUI string, userName string) bool {
	host := os.Getenv(iotterminalmgrServiceHost)
	port := os.Getenv(iotterminalmgrServicePortHTTP)
	var url = "http://" + host + ":" + port + "/iot/iotterminalmgr/checkLicense?tenantID=" + userName
	data, err := httprequest.HTTPRequest(url, "GET", "licenseCheck", []byte{})
	httpData := data.(*httprequest.HTTPData2)
	if err != nil {
		globallogger.Log.Errorln("devEUI :", devEUI, "HTTPRequestCheckLicense httprequest.HttpRequest error :", err)
	} else {
		globallogger.Log.Infoln("devEUI :", devEUI, "HTTPRequestCheckLicense httprequest.HttpRequest success", httpData)
		if httpData != nil {
			if httpData.ErrCode == 0 {
				globallogger.Log.Infoln("devEUI :", devEUI, "HTTPRequestCheckLicense is exist")
				if httpData.Data.Result != "" {
					return true
				}
			}
		}
	}
	globallogger.Log.Warnln("devEUI :", devEUI, "HTTPRequestCheckLicense is not exist")
	return false
}

//GetRedisLengthAndRangeRedis GetRedisLengthAndRangeRedis
func GetRedisLengthAndRangeRedis(key string, devEUI string) (int, []string, error) {
	var length int
	var res []string
	var err error
	if constant.Constant.MultipleInstances {
		length, err = globalrediscache.RedisCache.GetRedisLength(key)
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "getRedisLengthAndRangeRedis get redis length error :", err)
			return 0, nil, err
		}
		res, err = globalrediscache.RedisCache.RangeRedis(key, 0, length)
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "getRedisLengthAndRangeRedis range redis error :", err)
			return 0, nil, err
		}
	} else {
		length, err = globalmemorycache.MemoryCache.GetMemoryLength(key)
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "getRedisLengthAndRangeRedis get memory length error :", err)
			return 0, nil, err
		}
		res, err = globalmemorycache.MemoryCache.RangeMemory(key, 0, length)
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "getRedisLengthAndRangeRedis range memory error :", err)
			return 0, nil, err
		}
	}
	return length, res, err
}

//GetRedisDataFromRedis GetRedisDataFromRedis
func GetRedisDataFromRedis(key string, devEUI string, SN string) (*publicstruct.RedisData, error) {
	_, redisArray, _ := GetRedisLengthAndRangeRedis(key, devEUI)
	var redisData = publicstruct.RedisData{} //null;
	if len(redisArray) != 0 {
		for _, item := range redisArray {
			if item != constant.Constant.REDIS.ZigbeeRedisRemoveRedis {
				var itemData publicstruct.RedisData
				json.Unmarshal([]byte(item), &itemData) //JSON.parse(item)
				if itemData.SN == SN {
					redisData = itemData
				}
			}
		}
		return &redisData, nil
	}
	return nil, nil
}

//GetRedisDataFromDataRedis GetRedisDataFromDataRedis
func GetRedisDataFromDataRedis(key string, devEUI string, SN string) (int, int, *publicstruct.RedisData, error) {
	length, redisArray, _ := GetRedisLengthAndRangeRedis(key, devEUI)
	var redisData = publicstruct.RedisData{} //null
	var tempIndex = 0
	if len(redisArray) != 0 {
		for index, item := range redisArray {
			if item != constant.Constant.REDIS.ZigbeeRedisRemoveRedis {
				itemData := publicstruct.RedisData{}
				json.Unmarshal([]byte(item), &itemData)
				if itemData.DataSN == SN {
					redisData = itemData
					tempIndex = index
				}
			}
		}
		return length, tempIndex, &redisData, nil
	}
	return length, tempIndex, nil, nil
}

// CheckTerminalIsSensor CheckTerminalIsSensor
func CheckTerminalIsSensor(devEUI string, tmnType string) bool {
	var isSensor = false
	switch tmnType {
	case constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalEmergencyButton,
		constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalHumanDetector,
		constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalDoorSensor,
		constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalSmokeFireDetector,
		constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalWaterDetector,
		constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalGasDetector,
		constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalVoiceLightDetector,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM:
		// case constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalHumitureDetector:
		isSensor = true
		// globallogger.Log.Infoln("devEUI :", devEUI ,"", "CheckTerminalIsSensor is sensor:", tmnType)
	default:
		// globallogger.Log.Warnln("devEUI :", devEUI ,"", "CheckTerminalIsSensor is not sensor:", tmnType)
	}
	return isSensor
}

// CheckTerminalIsNeedBind CheckTerminalIsNeedBind
func CheckTerminalIsNeedBind(devEUI string, tmnType string) bool {
	var isNeed = true
	switch tmnType {
	case constant.Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalSingleLightingSwitch,
		constant.Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalDoubleLightingSwitch,
		constant.Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalTripleLightingSwitch,
		constant.Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalQuadrupleLightingSwitch,
		constant.Constant.TMNTYPE.SCENE.ZigbeeTerminalScene:
		isNeed = false
		// globallogger.Log.Infoln("devEUI :", devEUI ,"", "CheckTerminalIsNeedBind is need not to bind:", tmnType)
	default:
		// globallogger.Log.Warnln("devEUI :", devEUI ,"", "CheckTerminalIsNeedBind is need to bind:", tmnType)
	}
	return isNeed
}

// GetTerminalInfo GetTerminalInfo
func GetTerminalInfo(jsonInfo publicstruct.JSONInfo, devEUI string) config.TerminalInfo {
	var terminalInfo = config.TerminalInfo{}
	terminalInfo.DevEUI = devEUI
	if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 {
		terminalInfo.NwkAddr = jsonInfo.MessagePayload.Address
	}
	if jsonInfo.TunnelHeader.LinkInfo.ACMac != "" {
		terminalInfo.ACMac = jsonInfo.TunnelHeader.LinkInfo.ACMac
	}
	terminalInfo.APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
	if jsonInfo.TunnelHeader.LinkInfo.SecondAddr != "" {
		terminalInfo.SecondAddr = jsonInfo.TunnelHeader.LinkInfo.SecondAddr
	}
	terminalInfo.ModuleID = jsonInfo.MessagePayload.ModuleID
	if jsonInfo.TunnelHeader.LinkInfo.T300ID != "" {
		terminalInfo.T300ID = jsonInfo.TunnelHeader.LinkInfo.T300ID
	}
	if jsonInfo.TunnelHeader.Version != "" {
		terminalInfo.UDPVersion = jsonInfo.TunnelHeader.Version
	}
	return terminalInfo
}

// GetJSONInfo GetJSONInfo
func GetJSONInfo(terminalInfo config.TerminalInfo) publicstruct.JSONInfo {
	jsonInfo := publicstruct.JSONInfo{}
	jsonInfo.MessagePayload.Address = terminalInfo.DevEUI
	if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
		jsonInfo.TunnelHeader.Version = terminalInfo.UDPVersion
		jsonInfo.MessagePayload.Address = terminalInfo.NwkAddr
	}
	if terminalInfo.ACMac != "" {
		jsonInfo.TunnelHeader.LinkInfo.ACMac = terminalInfo.ACMac
	}
	jsonInfo.TunnelHeader.LinkInfo.APMac = terminalInfo.APMac
	if terminalInfo.SecondAddr != "" {
		jsonInfo.TunnelHeader.LinkInfo.SecondAddr = terminalInfo.SecondAddr
	}
	jsonInfo.MessagePayload.ModuleID = terminalInfo.ModuleID
	if terminalInfo.T300ID != "" {
		jsonInfo.TunnelHeader.LinkInfo.T300ID = terminalInfo.T300ID
	} else {
		jsonInfo.TunnelHeader.LinkInfo.T300ID = "ffffffffffffffff"
	}
	return jsonInfo
}

// GetTerminalInterval GetTerminalInterval
func GetTerminalInterval(tmnType string) int {
	var interval uint16 = 60 * 60
	switch tmnType {
	case constant.Constant.TMNTYPE.MAILEKE.ZigbeeTerminalPMT1004Detector:
		interval = constant.Constant.TIMER.MailekeKeepAliveTimer.ZigbeeTerminalPMT1004Detector
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalCOSensorEM
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHS2AQEM
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHTEM
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalPIRSensorEM
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalPIRILLSensorEF30
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSceneSwitchEM30:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalSceneSwitchEM30
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalSmokeSensorEM
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalWarningDevice
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalWaterSensorEM
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalDoorSensorEF30
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalIRControlEM
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalGASSensorEM
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalESocket
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHS2SW1LEFR30
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHS2SW2LEFR30
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHS2SW3LEFR30
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
		interval = constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalSmartPlug
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32:
		interval = constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSingleSwitch00500c32
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33:
		interval = constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalDoubleSwitch00500c33
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35:
		interval = constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalTripleSwitch00500c35
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141:
		interval = constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSingleSwitchHY0141
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142:
		interval = constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalDoubleSwitchHY0142
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
		interval = constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalTripleSwitchHY0143
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
		interval = constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSocket000a0c3c
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
		interval = constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSocket000a0c55
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
		interval = constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSocketHY0105
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
		interval = constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSocketHY0106
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal1SceneSwitch005f0cf1:
		interval = constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminal1SceneSwitch005f0cf1
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3:
		interval = constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminal2SceneSwitch005f0cf3
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
		interval = constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminal3SceneSwitch005f0cf2
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
		interval = constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminal6SceneSwitch005f0c3b
	default:
		globallogger.Log.Warnf("[GetTerminalInterval]: invalid tmnType: %s", tmnType)
	}
	return int(interval)
}

// GetTmnType GetTmnType
func GetTmnType(tmnType string) string {
	switch tmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
		tmnType = "空气质量仪"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
		tmnType = "10A智能插座U1"
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
		tmnType = "16A智能插座U1"
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
		tmnType = "10A智能插座U2"
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
		tmnType = "16A智能插座U2"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
		tmnType = "温湿度传感器"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM:
		tmnType = "一氧化碳探测器"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM:
		tmnType = "可燃气体探测器"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM:
		tmnType = "人体红外传感器"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30:
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM:
		tmnType = "烟火灾探测器"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice:
		tmnType = "声光报警器"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM:
		tmnType = "水浸探测器"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30:
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM:
		tmnType = "红外转发器"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30:
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32:
		tmnType = "一键零火开关"
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141:
		tmnType = "一键单火开关"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30:
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33:
		tmnType = "二键零火开关"
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142:
		tmnType = "二键单火开关"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30:
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35:
		tmnType = "三键零火开关"
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
		tmnType = "三键单火开关"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSceneSwitchEM30:
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal1SceneSwitch005f0cf1:
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3:
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
		tmnType = "六键情景开关"
	default:
		tmnType = "PMT1004_温湿度_PM"
	}

	return tmnType
}

// GetTmnTypeAndAttribute GetTmnTypeAndAttribute
func GetTmnTypeAndAttribute(tmnType string) publicstruct.TmnTypeAttr {
	var attribute = config.Attribute{}
	var status = []string{}
	var returnData = publicstruct.TmnTypeAttr{}
	switch tmnType {
	case "一联开关":
		tmnType = constant.Constant.TMNTYPE.SWITCH.ZigbeeTerminalSingleSwitch
		status = append(status, "OFF")
		attribute.Status = status
	case "二联开关":
		tmnType = constant.Constant.TMNTYPE.SWITCH.ZigbeeTerminalDoubleSwitch
		status = append(status, "OFF", "OFF")
		attribute.Status = status
	case "三联开关":
		tmnType = constant.Constant.TMNTYPE.SWITCH.ZigbeeTerminalTripleSwitch
		status = append(status, "OFF", "OFF", "OFF")
		attribute.Status = status
	case "四联开关":
		tmnType = constant.Constant.TMNTYPE.SWITCH.ZigbeeTerminalQuadrupleSwitch
		status = append(status, "OFF", "OFF", "OFF", "OFF")
		attribute.Status = status
	case "一联灯光联动":
		tmnType = constant.Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalSingleLightingSwitch
		status = append(status, "OFF", "OFF")
		attribute.Status = status
	case "二联灯光联动":
		tmnType = constant.Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalDoubleLightingSwitch
		status = append(status, "OFF", "OFF")
		attribute.Status = status
	case "三联灯光联动":
		tmnType = constant.Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalTripleLightingSwitch
		status = append(status, "OFF", "OFF", "OFF")
		attribute.Status = status
	case "四联灯光联动":
		tmnType = constant.Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalQuadrupleLightingSwitch
		status = append(status, "OFF", "OFF", "OFF", "OFF")
		attribute.Status = status
	case "场景面板":
		tmnType = constant.Constant.TMNTYPE.SCENE.ZigbeeTerminalScene
		status = append(status, "OFF")
		attribute.Status = status
	case "窗帘":
		tmnType = constant.Constant.TMNTYPE.CURTAIN.ZigbeeTerminalCurtain
		status = append(status, "OFF")
		attribute.Status = status
	case "插座":
		tmnType = constant.Constant.TMNTYPE.SOCKET.ZigbeeTerminalSocket
		status = append(status, "OFF")
		attribute.Status = status
	case "门磁":
		tmnType = constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalDoorSensor
	case "人体探测器":
		tmnType = constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalHumanDetector
	case "无线紧急按钮":
		tmnType = constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalEmergencyButton
	case "燃气泄漏探测器":
		tmnType = constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalGasDetector
	case "水浸探测器":
		tmnType = constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalWaterDetector
	case "烟火灾探测器":
		tmnType = constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalSmokeFireDetector
	case "声光探测器":
		tmnType = constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalVoiceLightDetector
	case "温湿度探测器":
		tmnType = constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalHumitureDetector
	case "红外宝":
		tmnType = constant.Constant.TMNTYPE.TRANSPONDER.ZigbeeTerminalInfraredTransponder
		var latestCommand = ""
		attribute.LatestCommand = latestCommand
	case "ZH-102-SN_CO2",
		"组合传感器_光照_PM2.5",
		"JZH-021-SN_温湿度_光照_CO2":
		var data = ""
		attribute.Data = data
	case "PMT1004_温湿度_PM":
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
		status = append(status, "OFF")
		attribute.Status = status
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM:
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141:
		status = append(status, "OFF")
		attribute.Status = status
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142:
		status = append(status, "OFF", "OFF")
		attribute.Status = status
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
		status = append(status, "OFF", "OFF", "OFF")
		attribute.Status = status
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSceneSwitchEM30:
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal1SceneSwitch005f0cf1:
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3:
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
	default:
		tmnType = ""
	}
	returnData.TmnType = tmnType
	returnData.Attribute = attribute

	return returnData
}

// UpdateSocketInfo 更新或创建UDP socket信息
func UpdateSocketInfo(jsonInfo publicstruct.JSONInfo, keepAliveInterval int64) {
	var APMac = jsonInfo.TunnelHeader.LinkInfo.APMac
	var setData = config.SocketInfo{}
	if jsonInfo.TunnelHeader.LinkInfo.ACMac != "" {
		setData.ACMac = jsonInfo.TunnelHeader.LinkInfo.ACMac
	}
	setData.APMac = APMac
	setData.Family = jsonInfo.Rinfo.Family
	setData.IPAddr = jsonInfo.Rinfo.Address
	setData.IPPort = jsonInfo.Rinfo.Port
	setData.UpdateTime = time.Now()
	setDataByte, _ := json.Marshal(setData)
	zigbeeServerSocketCache.Set([]byte(APMac), setDataByte, int(keepAliveInterval+5))
	if constant.Constant.UsePostgres {
		oSet := make(map[string]interface{})
		if jsonInfo.TunnelHeader.LinkInfo.ACMac != "" {
			oSet["acmac"] = jsonInfo.TunnelHeader.LinkInfo.ACMac
		}
		oSet["apmac"] = APMac
		oSet["family"] = jsonInfo.Rinfo.Family
		oSet["ipaddr"] = jsonInfo.Rinfo.Address
		oSet["ipport"] = jsonInfo.Rinfo.Port
		oSet["updatetime"] = time.Now()
		socketInfo, err := models.FindSocketInfoAndUpdatePG(APMac, oSet)
		if err != nil {
			globallogger.Log.Errorln("devEUI :", jsonInfo.MessagePayload.Address, "UpdateSocketInfo FindSocketInfoAndUpdatePG error :", err)
		} else {
			if socketInfo == nil {
				err = models.CreateSocketInfoPG(setData)
				if err != nil {
					globallogger.Log.Errorln("devEUI :", jsonInfo.MessagePayload.Address, "UpdateSocketInfo CreateSocketInfoPG error :", err)
					_, err = models.FindSocketInfoAndUpdatePG(APMac, oSet)
					if err != nil {
						globallogger.Log.Errorln("devEUI :", jsonInfo.MessagePayload.Address, "UpdateSocketInfo FindSocketInfoAndUpdatePG error :", err)
					}
				}
			}
		}
		return
	}
	socketInfo, err := models.FindSocketInfoAndUpdate(APMac, setData)
	if err != nil {
		globallogger.Log.Errorln("devEUI :", jsonInfo.MessagePayload.Address, "UpdateSocketInfo findSocketInfoAndUpdate error :", err)
	} else {
		if socketInfo == nil {
			err = models.CreateSocketInfo(setData)
			if err != nil {
				globallogger.Log.Errorln("devEUI :", jsonInfo.MessagePayload.Address, "UpdateSocketInfo createSocketInfo error :", err)
				_, err = models.FindSocketInfoAndUpdate(APMac, setData)
				if err != nil {
					globallogger.Log.Errorln("devEUI :", jsonInfo.MessagePayload.Address, "UpdateSocketInfo findSocketInfoAndUpdate error :", err)
				}
			}
		}
	}
}

// CheckTerminalSceneIsMatch CheckTerminalSceneIsMatch
func CheckTerminalSceneIsMatch(devEUI string, sceneID string, venderInfo publicstruct.VenderInfo) bool {
	var isMatch = false
	if venderInfo.VenderID != "" {
		if sceneID == venderInfo.VenderID {
			isMatch = true
		} else {
			globallogger.Log.Warnln("devEUI :", devEUI, "CheckTerminalSceneIsMatch this terminal is not match with scene. sceneID:",
				sceneID, "venderID:", venderInfo.VenderID)
		}
	} else { //若venderID未定义，说明不需要场景ID，不检测场景
		isMatch = true
	}
	return isMatch
}
