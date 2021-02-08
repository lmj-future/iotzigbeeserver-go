package iotsmartspace

import (
	"encoding/json"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmemorycache"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalrediscache"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/zclmsgdown"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

const (
	paramKeyGroupID        string = "groupId"
	paramKeyTerminalType   string = "terminalType"
	paramKeyTerminalID     string = "terminalId"
	paramKeyTerminalIDList string = "terminalIdList"
	paramKeyCode           string = "code"
	paramKeyMessage        string = "message"
)

// ActionTerminalLeave ActionTerminalLeave
func ActionTerminalLeave(devEUI string) {
	params := make(map[string]interface{}, 1)
	params[paramKeyTerminalID] = devEUI
	Publish(TopicZigbeeserverIotsmartspaceAction, MethodActionOffline, params, uuid.NewV4().String())
}

// ActionInsertReply ActionInsertReply
func ActionInsertReply(terminalInfo config.TerminalInfo) {
	devEUI := terminalInfo.DevEUI
	jsonInfo := publicfunction.GetJSONInfo(terminalInfo)
	key := constant.Constant.REDIS.ZigbeeRedisPermitJoinContentKey + terminalInfo.TmnType
	var redisDataStr string
	var err error
	if constant.Constant.MultipleInstances {
		redisDataStr, err = globalrediscache.RedisCache.GetRedis(key)
	} else {
		redisDataStr, err = globalmemorycache.MemoryCache.GetMemory(key)
	}
	if err != nil {
		globallogger.Log.Errorln("devEUI : "+devEUI+" "+"ActionInsertReply GetRedis err : ", err)
		publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
		return
	}
	if redisDataStr == "" {
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "ActionInsertReply permit join is not enable")
		publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
		return
	}
	redisDataMap := make(map[string]interface{}, 4)
	err = json.Unmarshal([]byte(redisDataStr), &redisDataMap)
	if err != nil {
		globallogger.Log.Errorln("devEUI : "+devEUI+" "+"ActionInsertReply Unmarshal err: ", err)
		publicfunction.SendTerminalLeaveReq(jsonInfo, devEUI)
		return
	}
	params := make(map[string]interface{}, 3)
	params[paramKeyGroupID] = redisDataMap["groupID"]
	params[paramKeyTerminalType] = redisDataMap["terminalType"]
	params[paramKeyTerminalID] = devEUI
	Publish(TopicZigbeeserverIotsmartspaceAction, MethodActionInsertReply, params, redisDataMap["msgID"])

	if constant.Constant.MultipleInstances {
		_, err := globalrediscache.RedisCache.DeleteRedis(key)
		if err != nil {
			globallogger.Log.Errorf("[ActionInsertReply]: DeleteRedis err: %+v", err)
		} else {
			// globallogger.Log.Infof("[ActionInsertReply]: DeleteRedis success, key: %s", key)
			globalrediscache.RedisCache.SremRedis(constant.Constant.REDIS.ZigbeeRedisPermitJoinContentSets, key)
		}
	} else {
		_, err := globalmemorycache.MemoryCache.DeleteMemory(key)
		if err != nil {
			globallogger.Log.Errorf("[ActionInsertReply]: DeleteMemory err: %+v", err)
		} else {
			// globallogger.Log.Infof("[ActionInsertReply]: DeleteMemory success, key: %s", key)
			globalmemorycache.MemoryCache.SremMemory(constant.Constant.REDIS.ZigbeeRedisPermitJoinContentSets, key)
		}
	}
}

// ActionInsertReplyTest ActionInsertReplyTest
func ActionInsertReplyTest(terminalInfo config.TerminalInfo) {
	params := make(map[string]interface{}, 3)
	params[paramKeyGroupID] = "root"
	params[paramKeyTerminalType] = terminalInfo.TmnType
	params[paramKeyTerminalID] = terminalInfo.DevEUI
	Publish(TopicZigbeeserverIotsmartspaceAction, MethodActionInsertReplyTest, params, uuid.NewV4().String())
}

func procActionCheckRedis(timeStamp int64, key string) {
	globallogger.Log.Infof("[procActionCheckRedis]: timeStamp: %v, key: %v", timeStamp, key)
	go func() {
		select {
		case <-time.After(time.Duration(constant.Constant.TIMER.ZigbeeTimerPermitJoinContent) * time.Second):
			redisData := make(map[string]interface{}, 4)
			if constant.Constant.MultipleInstances {
				data, _ := globalrediscache.RedisCache.GetRedis(key)
				if data != "" {
					err := json.Unmarshal([]byte(data), &redisData)
					if err != nil {
						globallogger.Log.Errorf("[procActionCheckRedis]: Unmarshal err: %+v", err)
					} else {
						if int64(redisData["timeStamp"].(float64)) == timeStamp {
							_, err := globalrediscache.RedisCache.DeleteRedis(key)
							if err != nil {
								globallogger.Log.Errorf("[procActionCheckRedis]: DeleteRedis err: %+v", err)
							} else {
								// globallogger.Log.Infof("[procActionCheckRedis]: DeleteRedis success, key: %s", key)
								globalrediscache.RedisCache.SremRedis(constant.Constant.REDIS.ZigbeeRedisPermitJoinContentSets, key)
							}
						}
					}
				}
			} else {
				data, _ := globalmemorycache.MemoryCache.GetMemory(key)
				if data != "" {
					err := json.Unmarshal([]byte(data), &redisData)
					if err != nil {
						globallogger.Log.Errorf("[procActionCheckRedis]: Unmarshal err: %+v", err)
					} else {
						if int64(redisData["timeStamp"].(float64)) == timeStamp {
							_, err := globalmemorycache.MemoryCache.DeleteMemory(key)
							if err != nil {
								globallogger.Log.Errorf("[procActionCheckRedis]: DeleteMemory err: %+v", err)
							} else {
								// globallogger.Log.Infof("[procActionCheckRedis]: DeleteMemory success, key: %s", key)
								globalmemorycache.MemoryCache.SremMemory(constant.Constant.REDIS.ZigbeeRedisPermitJoinContentSets, key)
							}
						}
					}
				}
			}
		}
	}()
}

func procActionInsert(mapParams map[string]interface{}, msgID interface{}) {
	if mapParams[paramKeyTerminalType] != nil {
		terminalType := mapParams[paramKeyTerminalType].(string)
		timeStamp := time.Now().UnixNano() / 1e6
		redisData := make(map[string]interface{}, 4)
		redisData["msgID"] = msgID
		redisData["terminalType"] = terminalType
		redisData["groupID"] = mapParams[paramKeyGroupID].(string)
		redisData["timeStamp"] = timeStamp
		key := constant.Constant.REDIS.ZigbeeRedisPermitJoinContentKey + terminalType
		data, _ := json.Marshal(redisData)
		if constant.Constant.MultipleInstances {
			_, err := globalrediscache.RedisCache.UpdateRedis(key, data)
			if err != nil {
				globallogger.Log.Errorf("[procAction]: UpdateRedis err: %+v", err)
			} else {
				globalrediscache.RedisCache.SaddRedis(constant.Constant.REDIS.ZigbeeRedisPermitJoinContentSets, key)
				procActionCheckRedis(timeStamp, key)
			}
		} else {
			_, err := globalmemorycache.MemoryCache.UpdateMemory(key, data)
			if err != nil {
				globallogger.Log.Errorf("[procAction]: UpdateMemory err: %+v", err)
			} else {
				globalmemorycache.MemoryCache.SaddMemory(constant.Constant.REDIS.ZigbeeRedisPermitJoinContentSets, key)
				procActionCheckRedis(timeStamp, key)
			}
		}
	}
}

func procActionInsertSuccessRead(devEUI string, dstEndpointIndex int, clusterID uint16) {
	cmd := common.Command{
		DstEndpointIndex: dstEndpointIndex,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.ReadAttribute,
		ClusterID:   clusterID,
		Command:     cmd,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

func procActionInsertSuccess(mapParams map[string]interface{}) {
	devEUI := mapParams[paramKeyTerminalID].(string)
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, map[string]interface{}{"isexist": true})
	} else {
		terminalInfo, err = models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, bson.M{"isExist": true})
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI]: %s [procActionInsertSuccess]: GetTerminalInfoByDevEUI err: %+v", devEUI, err)
		return
	}
	if terminalInfo != nil {
		endpointTemp := pq.StringArray{}
		if constant.Constant.UsePostgres {
			endpointTemp = terminalInfo.EndpointPG
		} else {
			endpointTemp = terminalInfo.Endpoint
		}
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM:
			go procActionInsertSuccessRead(devEUI, 0, 0x0001)
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
			go procActionInsertSuccessRead(devEUI, 0, 0x0006)
			go procActionInsertSuccessRead(devEUI, 0, 0x0702)
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			dstEndpoint := 0
			if len(endpointTemp) > 1 && endpointTemp[1] == "01" {
				dstEndpoint = 1
			}
			go procActionInsertSuccessRead(devEUI, dstEndpoint, 0x0001)
			go procActionInsertSuccessRead(devEUI, dstEndpoint, 0x0402)
			go procActionInsertSuccessRead(devEUI, dstEndpoint, 0x0405)
			go procActionInsertSuccessRead(devEUI, dstEndpoint, 0x042a)
			go procActionInsertSuccessRead(devEUI, dstEndpoint, 0x042b)
			go procActionInsertSuccessRead(devEUI, dstEndpoint, 0xfc81)
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
			if len(endpointTemp) > 0 && endpointTemp[0] == "01" {
				go procActionInsertSuccessRead(devEUI, 0, 0x0001)
				go procActionInsertSuccessRead(devEUI, 0, 0x0402)
				go procActionInsertSuccessRead(devEUI, 1, 0x0405)
			} else {
				go procActionInsertSuccessRead(devEUI, 1, 0x0001)
				go procActionInsertSuccessRead(devEUI, 1, 0x0402)
				go procActionInsertSuccessRead(devEUI, 0, 0x0405)
			}
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM:
			go procActionInsertSuccessRead(devEUI, 0, 0x0001)
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30:
			go procActionInsertSuccessRead(devEUI, 0, 0x0001)
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
			go procActionInsertSuccessRead(devEUI, 0, 0x0006)
			go procActionInsertSuccessRead(devEUI, 0, 0x0702)
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM:
			go procActionInsertSuccessRead(devEUI, 0, 0x0001)
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice:
			go procActionInsertSuccessRead(devEUI, 0, 0x0001)
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM:
			go procActionInsertSuccessRead(devEUI, 0, 0x0001)
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30:
			go procActionInsertSuccessRead(devEUI, 0, 0x0006)
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30:
			go procActionInsertSuccessRead(devEUI, 0, 0x0006)
			go procActionInsertSuccessRead(devEUI, 1, 0x0006)
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30:
			go procActionInsertSuccessRead(devEUI, 0, 0x0006)
			go procActionInsertSuccessRead(devEUI, 1, 0x0006)
			go procActionInsertSuccessRead(devEUI, 2, 0x0006)
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSceneSwitchEM30:
			go procActionInsertSuccessRead(devEUI, 0, 0x0001)
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30:
			go procActionInsertSuccessRead(devEUI, 0, 0x0001)
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM:
			go procActionInsertSuccessRead(devEUI, 0, 0x0001)
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141:
			go procActionInsertSuccessRead(devEUI, 0, 0x0006)
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142:
			go procActionInsertSuccessRead(devEUI, 0, 0x0006)
			go procActionInsertSuccessRead(devEUI, 1, 0x0006)
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
			go procActionInsertSuccessRead(devEUI, 0, 0x0006)
			go procActionInsertSuccessRead(devEUI, 1, 0x0006)
			go procActionInsertSuccessRead(devEUI, 2, 0x0006)
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
			go procActionInsertSuccessRead(devEUI, 0, 0x0006)
			go procActionInsertSuccessRead(devEUI, 0, 0x0702)
			go procActionInsertSuccessRead(devEUI, 0, 0x0b04)
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal1SceneSwitch005f0cf1:
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3:
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
		default:
			globallogger.Log.Warnf("[devEUI]: %s [procActionInsertSuccess]: invalid tmnType: %v", devEUI, terminalInfo.TmnType)
		}
	} else {
		var setData = config.TerminalInfo{
			DevEUI:  devEUI,
			IsExist: true,
		}
		if constant.Constant.UsePostgres {
			models.CreateTerminalPG(setData)
		} else {
			models.CreateTerminal(setData)
		}
	}
}

func terminalIntervalReq(devEUI string, clusterID uint16, endpointIndex int, interval uint16) {
	cmd := common.Command{
		Interval: interval,
	}
	cmd.DstEndpointIndex = endpointIndex
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.IntervalCommand,
		ClusterID:   clusterID,
		Command:     cmd,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

func syncTime(devEUI string, clusterID uint16, endpointIndex int) {
	cmd := common.Command{}
	cmd.DstEndpointIndex = endpointIndex
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SyncTime,
		ClusterID:   clusterID,
		Command:     cmd,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

// ActionInsertSuccess ActionInsertSuccess
func ActionInsertSuccess(terminalInfo config.TerminalInfo) {
	endpointTemp := pq.StringArray{}
	if constant.Constant.UsePostgres {
		endpointTemp = terminalInfo.EndpointPG
	} else {
		endpointTemp = terminalInfo.Endpoint
	}
	devEUI := terminalInfo.DevEUI
	endpointIndex := 0
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM:
		terminalIntervalReq(devEUI, 0x0001, endpointIndex, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalCOSensorEM)
		go procActionInsertSuccessRead(devEUI, 0, 0x0001)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
		terminalIntervalReq(devEUI, 0x0006, endpointIndex, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalESocket)
		go procActionInsertSuccessRead(devEUI, 0, 0x0006)
		go procActionInsertSuccessRead(devEUI, 0, 0x0702)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM:
		terminalIntervalReq(devEUI, 0x0500, endpointIndex, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalGASSensorEM)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
		if len(endpointTemp) > 1 && endpointTemp[1] == "01" {
			endpointIndex = 1
		}
		terminalIntervalReq(devEUI, 0x0001, endpointIndex, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHS2AQEM)
		go procActionInsertSuccessRead(devEUI, endpointIndex, 0x0001)
		go procActionInsertSuccessRead(devEUI, endpointIndex, 0x0402)
		go procActionInsertSuccessRead(devEUI, endpointIndex, 0x0405)
		go procActionInsertSuccessRead(devEUI, endpointIndex, 0x042a)
		go procActionInsertSuccessRead(devEUI, endpointIndex, 0x042b)
		go procActionInsertSuccessRead(devEUI, endpointIndex, 0xfc81)
		go func() {
			select {
			case <-time.After(5 * time.Second):
				terminalIntervalReq(devEUI, 0x0402, endpointIndex, 30*60)
			}
		}()
		go func() {
			select {
			case <-time.After(10 * time.Second):
				terminalIntervalReq(devEUI, 0x0405, endpointIndex, 30*60)
			}
		}()
		go func() {
			select {
			case <-time.After(15 * time.Second):
				syncTime(devEUI, 0x000a, endpointIndex)
			}
		}()
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
		if len(endpointTemp) > 1 && endpointTemp[1] == "01" {
			go procActionInsertSuccessRead(devEUI, 0, 0x0001)
			go procActionInsertSuccessRead(devEUI, 0, 0x0402)
			go procActionInsertSuccessRead(devEUI, 1, 0x0405)
			go func() {
				select {
				case <-time.After(5 * time.Second):
					terminalIntervalReq(devEUI, 0x0402, 1, 30*60)
				}
			}()
			go func() {
				select {
				case <-time.After(10 * time.Second):
					terminalIntervalReq(devEUI, 0x0405, 0, 30*60)
				}
			}()
		} else {
			go procActionInsertSuccessRead(devEUI, 1, 0x0001)
			go procActionInsertSuccessRead(devEUI, 1, 0x0402)
			go procActionInsertSuccessRead(devEUI, 0, 0x0405)
			go func() {
				select {
				case <-time.After(5 * time.Second):
					terminalIntervalReq(devEUI, 0x0402, 0, 30*60)
				}
			}()
			go func() {
				select {
				case <-time.After(10 * time.Second):
					terminalIntervalReq(devEUI, 0x0405, 1, 30*60)
				}
			}()
		}
		terminalIntervalReq(devEUI, 0x0001, endpointIndex, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHTEM)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM:
		terminalIntervalReq(devEUI, 0x0001, endpointIndex, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalPIRSensorEM)
		go procActionInsertSuccessRead(devEUI, 0, 0x0001)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30:
		terminalIntervalReq(devEUI, 0x0001, endpointIndex, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalPIRILLSensorEF30)
		go procActionInsertSuccessRead(devEUI, 0, 0x0001)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
		terminalIntervalReq(devEUI, 0x0006, endpointIndex, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalSmartPlug)
		go procActionInsertSuccessRead(devEUI, 0, 0x0006)
		go procActionInsertSuccessRead(devEUI, 0, 0x0702)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM:
		terminalIntervalReq(devEUI, 0x0001, endpointIndex, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalSmokeSensorEM)
		go procActionInsertSuccessRead(devEUI, 0, 0x0001)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice:
		terminalIntervalReq(devEUI, 0x0001, endpointIndex, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalWarningDevice)
		go procActionInsertSuccessRead(devEUI, 0, 0x0001)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM:
		terminalIntervalReq(devEUI, 0x0001, endpointIndex, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalWaterSensorEM)
		go procActionInsertSuccessRead(devEUI, 0, 0x0001)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30:
		terminalIntervalReq(devEUI, 0x0006, endpointIndex, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHS2SW1LEFR30)
		go procActionInsertSuccessRead(devEUI, 0, 0x0006)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30:
		terminalIntervalReq(devEUI, 0x0006, 0, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHS2SW2LEFR30)
		go procActionInsertSuccessRead(devEUI, 0, 0x0006)
		go procActionInsertSuccessRead(devEUI, 1, 0x0006)
		go func() {
			select {
			case <-time.After(5 * time.Second):
				terminalIntervalReq(devEUI, 0x0006, 1, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHS2SW2LEFR30)
			}
		}()
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30:
		terminalIntervalReq(devEUI, 0x0006, 0, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHS2SW3LEFR30)
		go procActionInsertSuccessRead(devEUI, 0, 0x0006)
		go procActionInsertSuccessRead(devEUI, 1, 0x0006)
		go procActionInsertSuccessRead(devEUI, 2, 0x0006)
		go func() {
			select {
			case <-time.After(5 * time.Second):
				terminalIntervalReq(devEUI, 0x0006, 1, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHS2SW3LEFR30)
			}
		}()
		go func() {
			select {
			case <-time.After(10 * time.Second):
				terminalIntervalReq(devEUI, 0x0006, 2, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHS2SW3LEFR30)
			}
		}()
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSceneSwitchEM30:
		terminalIntervalReq(devEUI, 0x0001, endpointIndex, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalSceneSwitchEM30)
		go procActionInsertSuccessRead(devEUI, 0, 0x0001)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30:
		terminalIntervalReq(devEUI, 0x0001, endpointIndex, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalDoorSensorEF30)
		go procActionInsertSuccessRead(devEUI, 0, 0x0001)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM:
		terminalIntervalReq(devEUI, 0x0001, endpointIndex, constant.Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalIRControlEM)
		go procActionInsertSuccessRead(devEUI, 0, 0x0001)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
		terminalIntervalReq(devEUI, 0x0000, endpointIndex, constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSocket000a0c3c)
		go procActionInsertSuccessRead(devEUI, 0, 0x0006)
		go procActionInsertSuccessRead(devEUI, 0, 0x0702)
		go procActionInsertSuccessRead(devEUI, 0, 0x0b04)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
		terminalIntervalReq(devEUI, 0x0000, endpointIndex, constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSocket000a0c55)
		go procActionInsertSuccessRead(devEUI, 0, 0x0006)
		go procActionInsertSuccessRead(devEUI, 0, 0x0702)
		go procActionInsertSuccessRead(devEUI, 0, 0x0b04)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
		terminalIntervalReq(devEUI, 0x0000, endpointIndex, constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSocketHY0105)
		go procActionInsertSuccessRead(devEUI, 0, 0x0006)
		go procActionInsertSuccessRead(devEUI, 0, 0x0702)
		go procActionInsertSuccessRead(devEUI, 0, 0x0b04)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
		terminalIntervalReq(devEUI, 0x0000, endpointIndex, constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSocketHY0106)
		go procActionInsertSuccessRead(devEUI, 0, 0x0006)
		go procActionInsertSuccessRead(devEUI, 0, 0x0702)
		go procActionInsertSuccessRead(devEUI, 0, 0x0b04)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141:
		terminalIntervalReq(devEUI, 0x0006, 0, constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSingleSwitch00500c32)
		go procActionInsertSuccessRead(devEUI, 0, 0x0006)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142:
		terminalIntervalReq(devEUI, 0x0006, 0, constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalDoubleSwitch00500c33)
		go procActionInsertSuccessRead(devEUI, 0, 0x0006)
		go procActionInsertSuccessRead(devEUI, 1, 0x0006)
		go func() {
			select {
			case <-time.After(5 * time.Second):
				terminalIntervalReq(devEUI, 0x0006, 1, constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalDoubleSwitch00500c33)
			}
		}()
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
		terminalIntervalReq(devEUI, 0x0006, 0, constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalTripleSwitch00500c35)
		go procActionInsertSuccessRead(devEUI, 0, 0x0006)
		go procActionInsertSuccessRead(devEUI, 1, 0x0006)
		go procActionInsertSuccessRead(devEUI, 2, 0x0006)
		go func() {
			select {
			case <-time.After(5 * time.Second):
				terminalIntervalReq(devEUI, 0x0006, 1, constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalTripleSwitch00500c35)
			}
		}()
		go func() {
			select {
			case <-time.After(10 * time.Second):
				terminalIntervalReq(devEUI, 0x0006, 2, constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalTripleSwitch00500c35)
			}
		}()
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal1SceneSwitch005f0cf1:
		terminalIntervalReq(devEUI, 0x0000, endpointIndex, constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminal1SceneSwitch005f0cf1)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3:
		terminalIntervalReq(devEUI, 0x0000, endpointIndex, constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminal2SceneSwitch005f0cf3)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
		terminalIntervalReq(devEUI, 0x0000, endpointIndex, constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminal3SceneSwitch005f0cf2)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
		terminalIntervalReq(devEUI, 0x0000, endpointIndex, constant.Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminal6SceneSwitch005f0c3b)
	default:
		globallogger.Log.Warnf("[procActionInsertSuccess]: invalid tmnType: %s", terminalInfo.TmnType)
	}
}

func procActionInsertFail(mapParams map[string]interface{}) {
	devEUI := mapParams[paramKeyTerminalID].(string)
	code := mapParams[paramKeyCode].(string)
	message := mapParams[paramKeyMessage].(string)
	if code == "1" {
		globallogger.Log.Infof("[devEUI]: %s [procActionInsertFail]: code: %s, message: %s", devEUI, code, message)
		publicfunction.ProcTerminalDelete(devEUI)
	} else {
		globallogger.Log.Infof("[devEUI]: %s [procActionInsertFail]: code: %s, message: %s", devEUI, code, message)
	}
}

func procActionDelete(mapParams map[string]interface{}) {
	for k, v := range mapParams {
		switch k {
		case paramKeyTerminalID:
			publicfunction.ProcTerminalDelete(v.(string))
		case paramKeyTerminalIDList:
			for _, devEUI := range v.([]interface{}) {
				publicfunction.ProcTerminalDelete(devEUI.(string))
			}
		}
	}
}

// procAction 处理终端增删改查事件
func procAction(mqttMsg mqttMsgSt) {
	globallogger.Log.Infof("[procAction]: mqttMsg.Params: %+v", mqttMsg.Params)
	mapParams := mqttMsg.Params.(map[string]interface{})
	if mqttMsg.Method != "" {
		switch mqttMsg.Method {
		case MethodActionInsert:
			procActionInsert(mapParams, mqttMsg.ID)
		case MethodActionInsertSuccess:
			procActionInsertSuccess(mapParams)
		case MethodActionInsertFail:
			procActionInsertFail(mapParams)
		case MethodActionDelete:
			procActionDelete(mapParams)
		default:
			globallogger.Log.Warnf("[procAction]: invalid method : %+v", mqttMsg.Method)
		}
	}
}

// procActionIotware 处理终端增删改查事件
func procActionIotware(networkInAckMsg publicstruct.NetworkInAckIotware) {
	globallogger.Log.Infof("[procActionIotware]: networkInAckMsg: %+v", networkInAckMsg)
	mapParams := make(map[string]interface{}, 1)
	mapParams[paramKeyTerminalID] = networkInAckMsg.Device
	if networkInAckMsg.IsOK {
		procActionInsertSuccess(mapParams)
	} else {
		publicfunction.ProcTerminalDelete(networkInAckMsg.Device)
	}
}
