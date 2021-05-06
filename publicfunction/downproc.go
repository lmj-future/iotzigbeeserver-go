package publicfunction

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmemorycache"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalrediscache"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalredisclient"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalsocket"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globaltransfersn"
	"github.com/h3c/iotzigbeeserver-go/metrics"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/lib/pq"
)

// sendDownMsg UDP 下行发送接口
func sendDownMsg(devEUI string, sendBuf []byte, APMac string) {
	var socketInfo *config.SocketInfo
	var err error
	socketInfoByte, err := zigbeeServerSocketCache.Get([]byte(APMac))
	if err == nil {
		json.Unmarshal(socketInfoByte, &socketInfo)
	} else {
		if constant.Constant.UsePostgres {
			socketInfo, err = models.GetSocketInfoByAPMacPG(APMac)
		} else {
			socketInfo, err = models.GetSocketInfoByAPMac(APMac)
		}
	}
	if err != nil {
		globallogger.Log.Errorf("devEUI : %s sendDownMsg GetSocketInfoByAPMac err: %+v", devEUI, err)
		return
	}
	if socketInfo == nil {
		globallogger.Log.Warnln("devEUI :", devEUI, "sendDownMsg GetSocketInfoByAPMac socketInfo is nil")
		return
	}
	if socketInfo.Family == "IPv6" {
		// to do
	} else {
		err = globalsocket.ServiceSocket.Send(sendBuf, socketInfo.IPPort, socketInfo.IPAddr)
		if err != nil {
			globallogger.Log.Errorf("devEUI : %s sendDownMsg UDP send err: %+v", devEUI, err)
			return
		}
		globallogger.Log.Infoln("devEUI :", devEUI, "send to", socketInfo.IPAddr, ":", socketInfo.IPPort, "send buf:", hex.EncodeToString(sendBuf))
		metrics.CountUdpSendTotal()
	}
}

// createAckTimer createAckTimer
func createAckTimer(index int, devEUI string, key string, SN string, APMac string, isReSend bool) {
	go func() {
		timer := time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerRetran) * time.Second)
		<-timer.C
		timer.Stop()
		if index == 3 {
			globallogger.Log.Warnln("devEUI :", devEUI, "createAckTimer msg has already resend three times. Msg lose")
			if constant.Constant.MultipleInstances {
				_, err := globalrediscache.RedisCache.DeleteRedis(key)
				if err != nil {
					globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "createAckTimer delete redis error :", err)
				}
			} else {
				_, err := globalmemorycache.MemoryCache.DeleteMemory(key)
				if err != nil {
					globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "createAckTimer delete memory error :", err)
				}
			}
		} else if index < 3 {
			index++
			checkAndSendDownMsg(index, devEUI, key, SN, APMac, false)
		}
	}()

	if index == 0 {
		go func() {
			timer := time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerSetMsgRemove) * time.Second)
			<-timer.C
			timer.Stop()
			IfMsgExistThenDelete(key, devEUI, SN, APMac)
		}()
	}
}

// checkAndSendDownMsg checkAndSendDownMsg
func checkAndSendDownMsg(index int, devEUI string, key string, SN string, APMac string, isReSend bool) {
	var isCtrlMsg = false
	if strings.Repeat(key[0:20], 1) == "zigbee_down_ctrl_msg" {
		isCtrlMsg = true
	}
	_, redisArray, _ := GetRedisLengthAndRangeRedis(key, devEUI)
	if len(redisArray) != 0 {
		var tempBuffer = []byte{}
		var isNeedDown = true
		for itemIndex, item := range redisArray {
			if item != constant.Constant.REDIS.ZigbeeRedisRemoveRedis {
				var itemData publicstruct.RedisData
				_ = json.Unmarshal([]byte(item), &itemData)
				if SN == itemData.SN {
					tempBuffer, _ = hex.DecodeString(string(itemData.SendBuf))
					if index > 0 {
						if !itemData.AckFlag {
							globallogger.Log.Warnln("devEUI :", devEUI, "createAckTimer server has not receive ACK,", index, "times resend this msg.")
							sendDownMsg(devEUI, tempBuffer, APMac)
							createAckTimer(index, devEUI, key, SN, APMac, isReSend)
						}
					} else {
						if itemIndex == 0 {
							sendDownMsg(devEUI, tempBuffer, APMac)
							createAckTimer(index, devEUI, key, SN, APMac, isReSend)
						} else {
							if isNeedDown {
								sendDownMsg(devEUI, tempBuffer, APMac)
								createAckTimer(index, devEUI, key, SN, APMac, isReSend)
								isNeedDown = false
							}
						}
					}
				} else {
					if isCtrlMsg {
						isNeedDown = false
					} else {
						if !itemData.AckFlag {
							isNeedDown = false
						}
					}
				}
			}
		}
	}
}

// msgInsertRedis msgInsertRedis
func msgInsertRedis(index int, devEUI string, key string, msgType string, optionType string,
	SN string, dataSN string, sendBuf []byte, data string, APMac string, msgID interface{}) {
	redisDataSerial, _ := json.Marshal(publicstruct.RedisData{
		DevEUI:     devEUI,
		MsgType:    msgType,
		OptionType: optionType,
		DataSN:     dataSN,
		SN:         SN,
		AckFlag:    false,
		ReSendTime: 0,
		Data:       data,
		SendBuf:    hex.EncodeToString(sendBuf),
		MsgID:      msgID,
	})
	if constant.Constant.MultipleInstances {
		_, err := globalrediscache.RedisCache.InsertRedis(key, redisDataSerial)
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "msgInsertRedis insert redis error :", err)
		} else {
			checkAndSendDownMsg(index, devEUI, key, SN, APMac, false)
		}
		_, err = globalrediscache.RedisCache.SaddRedis(constant.Constant.REDIS.ZigbeeRedisDownMsgSets, key)
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "msgInsertRedis sadd redis error :", err)
		}
		go CheckTimeoutAndDeleteRedisMsg(key, devEUI, false)
	} else {
		_, err := globalmemorycache.MemoryCache.InsertMemory(key, redisDataSerial)
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "msgInsertRedis insert memory error :", err)
		} else {
			checkAndSendDownMsg(index, devEUI, key, SN, APMac, false)
		}
		go CheckTimeoutAndDeleteFreeCacheMsg(key, devEUI)
	}
}

// encAndSendDownMsg encAndSendDownMsg
func encAndSendDownMsg(terminalInfo config.TerminalInfo, ackOptionType string, msgType string, optionType string,
	dataSN int, msgID interface{}, data string) {
	var transferKeyBuilder strings.Builder
	transferKeyBuilder.WriteString(constant.Constant.REDIS.ZigbeeRedisSNTransfer)
	transferKeyBuilder.WriteString(terminalInfo.APMac)
	var transferKey = transferKeyBuilder.String()
	var num int
	var err error
	if constant.Constant.MultipleInstances {
		num, err = globalredisclient.MyZigbeeServerRedisClient.Incr(transferKey)
		if err == nil && num > 512345678 {
			globalredisclient.MyZigbeeServerRedisClient.Del(transferKey)
			num, err = globalredisclient.MyZigbeeServerRedisClient.Incr(transferKey)
		}
	} else {
		globaltransfersn.TransferSN.Lock()
		if globaltransfersn.TransferSN.SN[transferKey] < 1 {
			globaltransfersn.TransferSN.SN[transferKey] = 1
		} else {
			globaltransfersn.TransferSN.SN[transferKey] = globaltransfersn.TransferSN.SN[transferKey] + 1
		}
		num = globaltransfersn.TransferSN.SN[transferKey]
		if err == nil && num > 512345678 {
			delete(globaltransfersn.TransferSN.SN, transferKey)
			globaltransfersn.TransferSN.SN[transferKey] = 1
			num = globaltransfersn.TransferSN.SN[transferKey]
		}
		globaltransfersn.TransferSN.Unlock()
	}
	if err == nil && num != 0 {
		var transferSN = num % 0xffff
		if transferSN == 0 {
			transferSN = 1
		}
		var tempTransferSNstrBuilder strings.Builder
		tempTransferSNstrBuilder.WriteString("0000")
		tempTransferSNstrBuilder.WriteString(strconv.FormatInt(int64(transferSN), 16))

		var encParameter = EncParameter{
			devEUI:        terminalInfo.DevEUI,
			APMac:         terminalInfo.APMac,
			moduleID:      terminalInfo.ModuleID,
			SN:            tempTransferSNstrBuilder.String()[tempTransferSNstrBuilder.Len()-4:],
			ackOptionType: ackOptionType,
			msgType:       msgType,
			optionType:    optionType,
			data:          data,
		}
		if terminalInfo.SecondAddr != "" {
			if len(terminalInfo.APMac) != len(terminalInfo.SecondAddr) {
				secondAddr := ""
				for _, v := range terminalInfo.SecondAddr {
					secondAddr += fmt.Sprintf("%x", v)
				}
				encParameter.APMac = secondAddr
			} else {
				encParameter.APMac = terminalInfo.SecondAddr
			}
		}
		if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
			encParameter.devEUI = terminalInfo.NwkAddr
		}
		if terminalInfo.ACMac != "" {
			encParameter.ACMac = terminalInfo.ACMac
		}
		if terminalInfo.T300ID != "" {
			encParameter.T300ID = terminalInfo.T300ID
		}
		var returnData = encDownBuf(encParameter, terminalInfo.UDPVersion)
		var key = GetCtrlMsgKey(terminalInfo.APMac, terminalInfo.ModuleID)
		if msgType == globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent {
			key = GetDataMsgKey(terminalInfo.APMac, terminalInfo.ModuleID)
		}
		var dataSNstr string
		if dataSN != 0 {
			var dataSNstrBuilder strings.Builder
			dataSNstrBuilder.WriteString("0000")
			dataSNstrBuilder.WriteString(strconv.FormatInt(int64(dataSN), 16))
			dataSNstr = dataSNstrBuilder.String()[dataSNstrBuilder.Len()-4:]
		}
		msgInsertRedis(0, terminalInfo.DevEUI, key, msgType, optionType, returnData.SN, dataSNstr,
			returnData.sendBuf, data, terminalInfo.APMac, msgID)
	}
}

// SendTerminalCheckExistReq SendTerminalCheckExistReq
func SendTerminalCheckExistReq(jsonInfo publicstruct.JSONInfo, devEUI string) {
	var terminalInfo = GetTerminalInfo(jsonInfo, devEUI)
	globallogger.Log.Warnln("devEUI :", jsonInfo.MessagePayload.Address, "SendTerminalCheckExistReq")
	encAndSendDownMsg(terminalInfo, globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp,
		globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalCheckExistReqEvent,
		globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp, 0, "", terminalInfo.DevEUI+terminalInfo.NwkAddr)
}

// SendTerminalLeaveReq SendTerminalLeaveReq
func SendTerminalLeaveReq(jsonInfo publicstruct.JSONInfo, devEUI string) {
	var terminalInfo = GetTerminalInfo(jsonInfo, devEUI)
	var data = "00" //RemoveChildren/Rejoin
	if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 {
		data = devEUI + data
	}
	globallogger.Log.Warnln("devEUI :", jsonInfo.MessagePayload.Address, "SendTerminalLeaveReq")
	encAndSendDownMsg(terminalInfo, globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp,
		globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalLeaveReqEvent,
		globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp, 0, "", data)
}

// SendPermitJoinReq SendPermitJoinReq
func SendPermitJoinReq(terminalInfo config.TerminalInfo, duration string, TCSignificance string) {
	globallogger.Log.Infoln("devEUI :", terminalInfo.DevEUI, "SendPermitJoinReq")
	encAndSendDownMsg(terminalInfo, globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp,
		globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalPermitJoinReqEvent,
		globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp, 0, "", duration+TCSignificance)
}

// ProcTerminalDelete ProcTerminalDelete
func ProcTerminalDelete(devEUI string) {
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
	}
	if err != nil {
		globallogger.Log.Errorln("devEUI :", devEUI, "ProcTerminalDelete GetTerminalInfoByDevEUI err:", err)
	}
	if terminalInfo != nil {
		var keyBuilder strings.Builder
		if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
			keyBuilder.WriteString(terminalInfo.APMac)
			keyBuilder.WriteString(terminalInfo.ModuleID)
			keyBuilder.WriteString(terminalInfo.NwkAddr)
		} else {
			keyBuilder.WriteString(terminalInfo.APMac)
			keyBuilder.WriteString(terminalInfo.ModuleID)
			keyBuilder.WriteString(terminalInfo.DevEUI)
		}
		DeleteTerminalInfoListCache(keyBuilder.String())
		if constant.Constant.UsePostgres {
			err = models.DeleteTerminalPG(devEUI)
		} else {
			err = models.DeleteTerminal(devEUI)
		}
		if err != nil {
			globallogger.Log.Errorln("devEUI :", devEUI, "ProcTerminalDelete DeleteTerminal err:", err)
		} else {
			var data = "00" //RemoveChildren/Rejoin
			if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
				data = devEUI + data
			}
			globallogger.Log.Warnln("devEUI :", devEUI, "ProcTerminalDelete terminal leave req")
			encAndSendDownMsg(*terminalInfo, globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp,
				globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalLeaveReqEvent,
				globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp, 0, "", data)
			if constant.Constant.MultipleInstances {
				DeleteRedisTerminalTimer(devEUI)
			} else {
				DeleteFreeCacheTerminalTimer(devEUI)
			}

			var key = GetCtrlMsgKey(terminalInfo.APMac, terminalInfo.ModuleID)
			IfCtrlMsgExistThenDeleteAll(key, devEUI, terminalInfo.APMac)
			key = GetDataMsgKey(terminalInfo.APMac, terminalInfo.ModuleID)
			IfDataMsgExistThenDeleteAll(key, devEUI, terminalInfo.APMac)
		}
	} else {
		globallogger.Log.Warnln("devEUI :", devEUI, "ProcTerminalDelete terminal is not exist")
	}
}

// ProcTerminalDeleteByOIDIndex ProcTerminalDeleteByOIDIndex
func ProcTerminalDeleteByOIDIndex(OIDIndex string) {
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByOIDIndexPG(OIDIndex)
	} else {
		terminalInfo, err = models.GetTerminalInfoByOIDIndex(OIDIndex)
	}
	if err != nil {
		globallogger.Log.Errorln("OIDIndex :", OIDIndex, "ProcTerminalDeleteByOIDIndex GetTerminalInfoByOIDIndex err:", err)
	}
	if terminalInfo != nil {
		var keyBuilder strings.Builder
		if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
			keyBuilder.WriteString(terminalInfo.APMac)
			keyBuilder.WriteString(terminalInfo.ModuleID)
			keyBuilder.WriteString(terminalInfo.NwkAddr)
		} else {
			keyBuilder.WriteString(terminalInfo.APMac)
			keyBuilder.WriteString(terminalInfo.ModuleID)
			keyBuilder.WriteString(terminalInfo.DevEUI)
		}
		DeleteTerminalInfoListCache(keyBuilder.String())
		if constant.Constant.UsePostgres {
			err = models.DeleteTerminalPG(terminalInfo.DevEUI)
		} else {
			err = models.DeleteTerminal(terminalInfo.DevEUI)
		}
		if err != nil {
			globallogger.Log.Errorln("devEUI :", terminalInfo.DevEUI, "ProcTerminalDeleteByOIDIndex DeleteTerminal err:", err)
		} else {
			var data = "00" //RemoveChildren/Rejoin
			if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
				data = terminalInfo.DevEUI + data
			}
			globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "ProcTerminalDeleteByOIDIndex terminal leave req")
			encAndSendDownMsg(*terminalInfo, globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp,
				globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalLeaveReqEvent,
				globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp, 0, "", data)
			if constant.Constant.MultipleInstances {
				DeleteRedisTerminalTimer(terminalInfo.DevEUI)
			} else {
				DeleteFreeCacheTerminalTimer(terminalInfo.DevEUI)
			}

			var key = GetCtrlMsgKey(terminalInfo.APMac, terminalInfo.ModuleID)
			IfCtrlMsgExistThenDeleteAll(key, terminalInfo.DevEUI, terminalInfo.APMac)
			key = GetDataMsgKey(terminalInfo.APMac, terminalInfo.ModuleID)
			IfDataMsgExistThenDeleteAll(key, terminalInfo.DevEUI, terminalInfo.APMac)
		}
	} else {
		globallogger.Log.Warnln("OIDIndex :", OIDIndex, "ProcTerminalDeleteByOIDIndex terminal is not exist")
	}
}

// SendBindTerminalReq SendBindTerminalReq
func SendBindTerminalReq(terminalInfo config.TerminalInfo) {
	bindInfo := []config.BindInfo{}
	if constant.Constant.UsePostgres {
		json.Unmarshal([]byte(terminalInfo.BindInfoPG), &bindInfo)
	} else {
		bindInfo = terminalInfo.BindInfo
	}
	var dataBuilder strings.Builder
	for bindInfoIndex, item := range bindInfo {
		for clusterIDIndex, clusterID := range item.ClusterID {
			globallogger.Log.Infoln("devEUI :", terminalInfo.DevEUI, "SendBindTerminalReq, bindInfoIndex:", bindInfoIndex,
				" clusterIDIndex:", clusterIDIndex, "clusterID:", clusterID)
			dataBuilder.Reset()
			if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
				dataBuilder.WriteString(terminalInfo.DevEUI)
				dataBuilder.WriteString(item.SrcEndpoint)
				dataBuilder.WriteString(clusterID)
				dataBuilder.WriteString(item.DstAddrMode)
				dataBuilder.WriteString(item.DstAddress)
				dataBuilder.WriteString(item.DstEndpoint)
			} else {
				dataBuilder.WriteString(item.SrcEndpoint)
				dataBuilder.WriteString(clusterID)
				dataBuilder.WriteString(item.DstAddrMode)
				dataBuilder.WriteString(item.DstAddress)
				dataBuilder.WriteString(item.DstEndpoint)
			}
			encAndSendDownMsg(terminalInfo, globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp,
				globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalBindReqEvent,
				globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp, 0, "", dataBuilder.String())
		}
	}
	SendPermitJoinReq(terminalInfo, "00", "00")
}

// GetTerminalEndpoint GetTerminalEndpoint
func GetTerminalEndpoint(terminalInfo config.TerminalInfo) {
	globallogger.Log.Infoln("devEUI :", terminalInfo.DevEUI, "GetTerminalEndpoint")
	encAndSendDownMsg(terminalInfo, globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp,
		globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalEndpointReqEvent,
		globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp, 0, "", "")
}

// SendTerminalDiscovery SendTerminalDiscovery
func SendTerminalDiscovery(terminalInfo config.TerminalInfo, endpoint string) {
	globallogger.Log.Infoln("devEUI :", terminalInfo.DevEUI, "SendTerminalDiscovery")
	encAndSendDownMsg(terminalInfo, globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp,
		globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalDiscoveryReqEvent,
		globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp, 0, "", endpoint)
}

// SendZHADefaultResponseDownMsg SendZHADefaultResponseDownMsg
func SendZHADefaultResponseDownMsg(downMsg common.DownMsg, terminalInfo *config.TerminalInfo, dstEndpoint string) {
	var profileIDBuilder strings.Builder
	profileIDBuilder.WriteString("0000")
	profileIDBuilder.WriteString(strconv.FormatUint(uint64(downMsg.ProfileID), 16))
	var clusterIDBuilder strings.Builder
	clusterIDBuilder.WriteString("0000")
	clusterIDBuilder.WriteString(strconv.FormatUint(uint64(downMsg.ClusterID), 16))
	var dataBuilder strings.Builder
	dataBuilder.WriteString(profileIDBuilder.String()[profileIDBuilder.Len()-4:])
	dataBuilder.WriteString(clusterIDBuilder.String()[clusterIDBuilder.Len()-4:])
	dataBuilder.WriteString(dstEndpoint)
	dataBuilder.WriteString(downMsg.ZclData)
	var transferKeyBuilder strings.Builder
	transferKeyBuilder.WriteString(constant.Constant.REDIS.ZigbeeRedisSNTransfer)
	transferKeyBuilder.WriteString(terminalInfo.APMac)
	var transferKey = transferKeyBuilder.String()
	var num int
	var err error
	if constant.Constant.MultipleInstances {
		num, err = globalredisclient.MyZigbeeServerRedisClient.Incr(transferKey)
		if err == nil && num > 512345678 {
			globalredisclient.MyZigbeeServerRedisClient.Del(transferKey)
			num, err = globalredisclient.MyZigbeeServerRedisClient.Incr(transferKey)
		}
	} else {
		globaltransfersn.TransferSN.Lock()
		if globaltransfersn.TransferSN.SN[transferKey] < 1 {
			globaltransfersn.TransferSN.SN[transferKey] = 1
		} else {
			globaltransfersn.TransferSN.SN[transferKey] = globaltransfersn.TransferSN.SN[transferKey] + 1
		}
		num = globaltransfersn.TransferSN.SN[transferKey]
		if err == nil && num > 512345678 {
			delete(globaltransfersn.TransferSN.SN, transferKey)
			globaltransfersn.TransferSN.SN[transferKey] = 1
			num = globaltransfersn.TransferSN.SN[transferKey]
		}
		globaltransfersn.TransferSN.Unlock()
	}
	if err == nil && num != 0 {
		var transferSN = num % 0xffff
		if transferSN == 0 {
			transferSN = 1
		}
		var tempTransferSNstrBuilder strings.Builder
		tempTransferSNstrBuilder.WriteString("0000")
		tempTransferSNstrBuilder.WriteString(strconv.FormatInt(int64(transferSN), 16))

		var encParameter = EncParameter{
			devEUI:        terminalInfo.DevEUI,
			APMac:         terminalInfo.APMac,
			moduleID:      terminalInfo.ModuleID,
			SN:            tempTransferSNstrBuilder.String()[tempTransferSNstrBuilder.Len()-4:],
			ackOptionType: globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp,
			msgType:       downMsg.MsgType,
			optionType:    globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeRsp,
			data:          dataBuilder.String(),
		}
		if terminalInfo.SecondAddr != "" {
			if len(terminalInfo.APMac) != len(terminalInfo.SecondAddr) {
				secondAddr := ""
				for _, v := range terminalInfo.SecondAddr {
					secondAddr += fmt.Sprintf("%x", v)
				}
				encParameter.APMac = secondAddr
			} else {
				encParameter.APMac = terminalInfo.SecondAddr
			}
		}
		if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
			encParameter.devEUI = terminalInfo.NwkAddr
		}
		if terminalInfo.ACMac != "" {
			encParameter.ACMac = terminalInfo.ACMac
		}
		if terminalInfo.T300ID != "" {
			encParameter.T300ID = terminalInfo.T300ID
		}
		var returnData = encDownBuf(encParameter, terminalInfo.UDPVersion)
		sendDownMsg(terminalInfo.DevEUI, returnData.sendBuf, terminalInfo.APMac)
	}
}

// SendZHADownMsg SendZHADownMsg
func SendZHADownMsg(downMsg common.DownMsg) {
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(downMsg.DevEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(downMsg.DevEUI)
	}
	if err != nil {
		globallogger.Log.Errorln("devEUI :", downMsg.DevEUI, "SendZHADownMsg GetTerminalInfoByDevEUI err:", err)
		return
	}
	if terminalInfo != nil {
		var endpointTemp pq.StringArray
		if constant.Constant.UsePostgres {
			endpointTemp = terminalInfo.EndpointPG
		} else {
			endpointTemp = terminalInfo.Endpoint
		}
		if len(endpointTemp) > 0 {
			var profileIDBuilder strings.Builder
			profileIDBuilder.WriteString("0000")
			profileIDBuilder.WriteString(strconv.FormatUint(uint64(downMsg.ProfileID), 16))
			var clusterIDBuilder strings.Builder
			clusterIDBuilder.WriteString("0000")
			clusterIDBuilder.WriteString(strconv.FormatUint(uint64(downMsg.ClusterID), 16))
			var DstEndpoint string
			if terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105 ||
				terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106 ||
				terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c ||
				terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55 {
				if downMsg.DstEndpointIndex == 1 {
					DstEndpoint = "02"
				} else {
					DstEndpoint = endpointTemp[downMsg.DstEndpointIndex]
				}
			} else {
				DstEndpoint = endpointTemp[downMsg.DstEndpointIndex]
			}
			var dataBuilder strings.Builder
			dataBuilder.WriteString(profileIDBuilder.String()[profileIDBuilder.Len()-4:])
			dataBuilder.WriteString(clusterIDBuilder.String()[clusterIDBuilder.Len()-4:])
			dataBuilder.WriteString(DstEndpoint)
			dataBuilder.WriteString(downMsg.ZclData)
			encAndSendDownMsg(*terminalInfo, globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp, downMsg.MsgType,
				globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp, int(downMsg.SN), downMsg.MsgID, dataBuilder.String())
		}
	}
}

//EncParameter EncParameter
type EncParameter struct {
	devEUI        string
	APMac         string
	moduleID      string
	SN            string
	ackOptionType string
	msgType       string
	optionType    string
	data          string
	ACMac         string
	T300ID        string
}

// SendACK SendACK
func SendACK(jsonInfo publicstruct.JSONInfo) {
	var terminalInfo = GetTerminalInfo(jsonInfo, jsonInfo.MessagePayload.Address)
	var msgType = globalmsgtype.MsgType.GENERALMsg.ZigbeeGeneralAck
	if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0101 ||
		jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 {
		msgType = globalmsgtype.MsgType.GENERALMsgV0101.ZigbeeGeneralAckV0101
	}
	globallogger.Log.Infoln("devEUI :", jsonInfo.MessagePayload.Address, "APMAC:", jsonInfo.TunnelHeader.LinkInfo.APMac,
		" SendACK, msgType:", jsonInfo.MessageHeader.MsgType, globalmsgtype.MsgType.GetMsgTypeUPMsgMeaning(jsonInfo.MessageHeader.MsgType))
	var encParameter = EncParameter{
		devEUI:        terminalInfo.DevEUI,
		APMac:         terminalInfo.APMac,
		moduleID:      terminalInfo.ModuleID,
		SN:            jsonInfo.TunnelHeader.FrameSN,
		ackOptionType: globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeRsp,
		msgType:       msgType,
		optionType:    globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNoRsp,
		data:          jsonInfo.MessageHeader.MsgType,
	}
	if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
		encParameter.devEUI = terminalInfo.NwkAddr
	}
	if terminalInfo.ACMac != "" {
		encParameter.ACMac = terminalInfo.ACMac
	}
	if terminalInfo.T300ID != "" {
		encParameter.T300ID = terminalInfo.T300ID
	}
	var returnData = encDownBuf(encParameter, jsonInfo.TunnelHeader.Version)
	sendDownMsg(terminalInfo.DevEUI, returnData.sendBuf, terminalInfo.APMac)
}

// CheckRedisAndReSend CheckRedisAndReSend
func CheckRedisAndReSend(devEUI string, key string, redisData publicstruct.RedisData, APMac string) {
	if redisData.AckFlag {
		redisData.AckFlag = false
	}
	if redisData.ReSendTime < 3 {
		redisData.ReSendTime++
		go func() {
			timer := time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerRetran) * time.Second)
			<-timer.C
			timer.Stop()
			redisDataSerial, _ := json.Marshal(redisData)
			if constant.Constant.MultipleInstances {
				_, err := globalrediscache.RedisCache.SetRedis(key, 0, redisDataSerial)
				if err != nil {
					globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "CheckRedisAndReSend set redis error :", err)
				} else {
					checkAndSendDownMsg(0, devEUI, key, redisData.SN, APMac, true)
				}
			} else {
				_, err := globalmemorycache.MemoryCache.SetMemory(key, 0, redisDataSerial)
				if err != nil {
					globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "CheckRedisAndReSend set memory error :", err)
				} else {
					checkAndSendDownMsg(0, devEUI, key, redisData.SN, APMac, true)
				}
			}
		}()
	} else {
		globallogger.Log.Warnln("devEUI :", devEUI, "key:", key, "CheckRedisAndReSend has already resend 3 times, do nothing ")
	}
}

// CheckRedisAndSend CheckRedisAndSend
func CheckRedisAndSend(devEUI string, key string, SN string, APMac string) {
	checkAndSendDownMsg(0, devEUI, key, SN, APMac, false)
}
