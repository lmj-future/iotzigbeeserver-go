package publicfunction

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/dgram"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmemorycache"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalrediscache"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalredisclient"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalsocket"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globaltransfersn"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/lib/pq"
)

var zigbeeServerReSendMsgCheckTimerID = &sync.Map{}

// Body Body
type Body struct {
	rinfo dgram.RInfo
	msg   string
}

//SendData SendData
type SendData struct {
	extType      string
	devEUI       string
	body         Body
	isSendConfig bool
}

// sendDownMsg UDP 下行发送接口
func sendDownMsg(devEUI string, sendBuf []byte, APMac string) {
	var socketInfo *config.SocketInfo
	var err error
	if constant.Constant.UsePostgres {
		socketInfo, err = models.GetSocketInfoByAPMacPG(APMac)
	} else {
		socketInfo, err = models.GetSocketInfoByAPMac(APMac)
	}
	if err != nil {
		globallogger.Log.Infof("devEUI : "+devEUI+" "+"sendDownMsg GetSocketInfoByAPMac err: %+v", err)
		return
	}
	if socketInfo == nil {
		globallogger.Log.Infof("devEUI : " + devEUI + " " + "sendDownMsg GetSocketInfoByAPMac socketInfo is nil")
		return
	}
	// globallogger.Log.Infof("devEUI : "+devEUI+" "+"sendDownMsg GetSocketInfoByAPMac: %+v", socketInfo)
	if socketInfo.Family == "IPv6" {
		var sendData = SendData{
			extType: "iotzigbeeserver",
			body: Body{
				rinfo: dgram.RInfo{
					Address: socketInfo.IPAddr,
					Port:    socketInfo.IPPort,
				},
				msg: hex.EncodeToString(sendBuf), //sendBuf.toString("hex")
			},
		}
		globallogger.Log.Infoln("devEUI : "+devEUI+" "+"mqhd send ipv6 msg to IPv6Server:", fmt.Sprintf("%+v", sendData))
		//sendDataSerial, _ := json.Marshal(sendData)
		//mqhd.SendMsg("IPv6Server", sendDataSerial)
	} else {
		err = globalsocket.ServiceSocket.Send(sendBuf, socketInfo.IPPort, socketInfo.IPAddr)
		if err != nil {
			globallogger.Log.Infof("devEUI : "+devEUI+" "+"sendDownMsg UDP send err: %+v", err)
			return
		}
		globallogger.Log.Infoln("devEUI : "+devEUI+" "+"send to "+socketInfo.IPAddr+":", socketInfo.IPPort, " send buf: "+hex.EncodeToString(sendBuf))
		// fmt.Println("devEUI : "+devEUI+" "+"send to "+socketInfo.IPAddr+":", socketInfo.IPPort, " send buf: "+hex.EncodeToString(sendBuf))
	}
}

// createAckTimer createAckTimer
func createAckTimer(index int, devEUI string, key string, SN string, APMac string, isReSend bool) {
	go func() {
		select {
		case <-time.After(time.Duration(constant.Constant.TIMER.ZigbeeTimerRetran) * time.Second):
			if index == 3 {
				globallogger.Log.Warnln("devEUI : " + devEUI + " " + "createAckTimer msg has already resend three times. Msg lose")
				if constant.Constant.MultipleInstances {
					_, err := globalrediscache.RedisCache.DeleteRedis(key)
					if err != nil {
						globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" createAckTimer delete redis error : ", err)
					} else {
						// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" createAckTimer delete redis success : ", res)
					}
				} else {
					_, err := globalmemorycache.MemoryCache.DeleteMemory(key)
					if err != nil {
						globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" createAckTimer delete memory error : ", err)
					} else {
						// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" createAckTimer delete memory success : ", res)
					}
				}
			} else if index < 3 {
				index++
				checkAndSendDownMsg(index, devEUI, key, SN, APMac, false)
			}
		}
	}()

	if index == 0 {
		if isReSend {
			reSendTimerUpdateTime := ReSendTimerUpdateOrCreate(key + SN)
			if value, ok := zigbeeServerReSendMsgCheckTimerID.Load(key + SN); ok {
				value.(*time.Timer).Stop()
				// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "key+SN: " + key + SN + " createAckTimer clearTimeout")
			}
			var timerID = time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerSetMsgRemove) * time.Second)
			zigbeeServerReSendMsgCheckTimerID.Store(key+SN, timerID)
			go func() {
				select {
				case <-timerID.C:
					zigbeeServerReSendMsgCheckTimerID.Delete(key + SN)
					reSendTimerInfo, err := GetReSendTimerByMsgKey(key + SN)
					if err != nil {
						globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key+SN: "+key+SN+" createAckTimer GetReSendTimerByMsgKey err: ", err)
						return
					}
					if reSendTimerInfo != nil {
						var dateNow = time.Now()
						var num = dateNow.UnixNano() - reSendTimerUpdateTime.UnixNano()
						if num > int64(time.Duration(constant.Constant.TIMER.ZigbeeTimerSetMsgRemove)*time.Second) {
							IfMsgExistThenDelete(key, devEUI, SN, APMac)
							go DeleteReSendTimer(key + SN)
						}
					}
				}
			}()
		} else {
			go func() {
				select {
				case <-time.After(time.Duration(constant.Constant.TIMER.ZigbeeTimerSetMsgRemove) * time.Second):
					IfMsgExistThenDelete(key, devEUI, SN, APMac)
				}
			}()
		}
	}
}

// checkAndSendDownMsg checkAndSendDownMsg
func checkAndSendDownMsg(index int, devEUI string, key string, SN string, APMac string, isReSend bool) {
	var isCtrlMsg = false
	if key[0:20] == "zigbee_down_ctrl_msg" {
		isCtrlMsg = true
	}
	_, redisArray, _ := GetRedisLengthAndRangeRedis(key, devEUI)
	if len(redisArray) != 0 {
		var tempBuffer = []byte{}
		var isNeedDown = true
		for itemIndex, item := range redisArray {
			if item != constant.Constant.REDIS.ZigbeeRedisRemoveRedis {
				var itemData publicstruct.RedisData
				_ = json.Unmarshal([]byte(item), &itemData) //JSON.parse(item)
				if SN == itemData.SN {
					tempBuffer, _ = hex.DecodeString(string(itemData.SendBuf)) //Buffer.from(itemData.sendBuf, "hex")
					if index > 0 {                                             //����û���յ�ack����ʱ�ش�����Ϣ
						if itemData.AckFlag == false {
							globallogger.Log.Warnln("devEUI : "+devEUI+" "+"createAckTimer server has not receive ACK, ", index, " times resend this msg.")
							sendDownMsg(devEUI, tempBuffer, APMac)
							createAckTimer(index, devEUI, key, SN, APMac, isReSend)
						} else {
							// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "createAckTimer server has already receive ACK, do nothing.")
						}
					} else {
						if itemIndex == 0 { //�����׸���Ϣ���϶��·�
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
					if isCtrlMsg { //���Ʊ��ģ�ֻҪǰһ����Ϣ������Ϣ�����У��Ͳ��·���һ����Ϣ
						isNeedDown = false
					} else {
						if itemData.AckFlag == false { //���ݱ��ģ����ǰһ����Ϣû���յ�ack�����·���һ����Ϣ
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
	var redisData = publicstruct.RedisData{}
	redisData.DevEUI = devEUI                       //�豸��ʶ
	redisData.MsgType = msgType                     //��Ϣ����
	redisData.OptionType = optionType               //Ӧ������
	redisData.DataSN = dataSN                       //���ݱ������к�
	redisData.SN = SN                               //����㱨�����к�
	redisData.AckFlag = false                       //ack���
	redisData.ReSendTime = 0                        //�ش�����
	redisData.Data = data                           //�����ֶ�
	redisData.SendBuf = hex.EncodeToString(sendBuf) //��Ҫ���͵�buffer
	redisData.MsgID = msgID
	redisDataSerial, _ := json.Marshal(redisData)
	if constant.Constant.MultipleInstances {
		_, err := globalrediscache.RedisCache.InsertRedis(key, redisDataSerial)
		if err != nil {
			globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" msgInsertRedis insert redis error : ", err)
		} else {
			// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" msgInsertRedis insert redis success : ", res)
			checkAndSendDownMsg(index, devEUI, key, SN, APMac, false)
		}
		_, err = globalrediscache.RedisCache.SaddRedis("zigbee_down_msg", key)
		if err != nil {
			globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" msgInsertRedis sadd redis error : ", err)
		} else {
			// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" msgInsertRedis sadd redis success: ", res)
		}
	} else {
		_, err := globalmemorycache.MemoryCache.InsertMemory(key, redisDataSerial)
		if err != nil {
			globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" msgInsertRedis insert memory error : ", err)
		} else {
			// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" msgInsertRedis insert memory success : ", res)
			checkAndSendDownMsg(index, devEUI, key, SN, APMac, false)
		}
	}
	go CheckTimeoutAndDeleteRedisMsg(key, devEUI, false)
}

// encAndSendDownMsg encAndSendDownMsg
func encAndSendDownMsg(terminalInfo config.TerminalInfo, ackOptionType string, msgType string, optionType string,
	dataSN int, msgID interface{}, data string) {
	var transferKey = constant.Constant.REDIS.ZigbeeRedisSNTransfer + terminalInfo.APMac
	var num1 int
	var err error
	if constant.Constant.MultipleInstances {
		num1, err = globalredisclient.MyZigbeeServerRedisClient.Incr(transferKey)
	} else {
		globaltransfersn.TransferSN.Lock()
		if globaltransfersn.TransferSN.SN[transferKey] < 1 {
			globaltransfersn.TransferSN.SN[transferKey] = 1
		} else {
			globaltransfersn.TransferSN.SN[transferKey] = globaltransfersn.TransferSN.SN[transferKey] + 1
		}
		num1 = globaltransfersn.TransferSN.SN[transferKey]
		globaltransfersn.TransferSN.Unlock()
	}
	if err == nil && num1 != 0 {
		var transferSN = num1 % 0xffff
		if transferSN == 0 {
			transferSN = 1
		}
		tempTransferSNstr := "0000" + strconv.FormatInt(int64(transferSN), 16)
		transferSNstr := tempTransferSNstr[len(tempTransferSNstr)-4:]

		var encParameter = EncParameter{
			devEUI:        terminalInfo.DevEUI,
			APMac:         terminalInfo.APMac,
			moduleID:      terminalInfo.ModuleID,
			SN:            transferSNstr,
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
		var index = 0 //������ʾ�ڼ����ش�
		var key = constant.Constant.REDIS.ZigbeeRedisCtrlMsgKey + terminalInfo.APMac + "_" + terminalInfo.ModuleID
		if msgType == globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent {
			key = constant.Constant.REDIS.ZigbeeRedisDataMsgKey + terminalInfo.APMac + "_" + terminalInfo.ModuleID
		}
		var dataSNstr string
		if dataSN != 0 {
			tempVaule := "0000" + strconv.FormatInt(int64(dataSN), 16)
			dataSNstr = tempVaule[len(tempVaule)-4:]
		}
		msgInsertRedis(index, terminalInfo.DevEUI, key, msgType, optionType, returnData.SN, dataSNstr, returnData.sendBuf, data, terminalInfo.APMac, msgID)
		if num1 > 512345678 {
			go func() {
				select {
				case <-time.After(time.Duration(constant.Constant.TIMER.ZigbeeTimerOneDay) * time.Second):
					if constant.Constant.MultipleInstances {
						num2, err := globalredisclient.MyZigbeeServerRedisClient.Get(transferKey)
						if err == nil && num2 != "" {
							tempValue, _ := strconv.ParseInt(num2, 10, 0)
							if int(tempValue) == num1 {
								globalredisclient.MyZigbeeServerRedisClient.Del(transferKey)
							}
						}
					} else {
						globaltransfersn.TransferSN.RLock()
						num2 := globaltransfersn.TransferSN.SN[transferKey]
						globaltransfersn.TransferSN.RUnlock()
						if num2 == num1 {
							globaltransfersn.TransferSN.Lock()
							delete(globaltransfersn.TransferSN.SN, transferKey)
							globaltransfersn.TransferSN.Unlock()
						}
					}
				}
			}()
		} else {
			go func() {
				select {
				case <-time.After(time.Duration(60*60) * time.Second):
					if constant.Constant.MultipleInstances {
						num2, err := globalredisclient.MyZigbeeServerRedisClient.Get(transferKey)
						if err == nil && num2 != "" {
							tempValue, _ := strconv.ParseInt(num2, 10, 0)
							if int(tempValue) == num1 {
								globalredisclient.MyZigbeeServerRedisClient.Del(transferKey)
							}
						}
					} else {
						globaltransfersn.TransferSN.RLock()
						num2 := globaltransfersn.TransferSN.SN[transferKey]
						globaltransfersn.TransferSN.RUnlock()
						if num2 == num1 {
							globaltransfersn.TransferSN.Lock()
							delete(globaltransfersn.TransferSN.SN, transferKey)
							globaltransfersn.TransferSN.Unlock()
						}
					}
				}
			}()
		}
	}
}

// SendTerminalCheckExistReq SendTerminalCheckExistReq
func SendTerminalCheckExistReq(jsonInfo publicstruct.JSONInfo, devEUI string) {
	var terminalInfo = GetTerminalInfo(jsonInfo, devEUI)
	var ackOptionType = globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp
	var msgType = globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalCheckExistReqEvent
	var optionType = globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp
	var data = terminalInfo.DevEUI + terminalInfo.NwkAddr
	globallogger.Log.Warnln("devEUI : " + jsonInfo.MessagePayload.Address + " " + "SendTerminalCheckExistReq")
	encAndSendDownMsg(terminalInfo, ackOptionType, msgType, optionType, 0, "", data)
}

// SendTerminalLeaveReq SendTerminalLeaveReq
func SendTerminalLeaveReq(jsonInfo publicstruct.JSONInfo, devEUI string) {
	var terminalInfo = GetTerminalInfo(jsonInfo, devEUI)
	var ackOptionType = globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp
	var msgType = globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalLeaveReqEvent
	var optionType = globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp
	var data = "00" //RemoveChildren/Rejoin
	if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 {
		data = devEUI + data
	}
	globallogger.Log.Warnln("devEUI : " + jsonInfo.MessagePayload.Address + " " + "SendTerminalLeaveReq")
	encAndSendDownMsg(terminalInfo, ackOptionType, msgType, optionType, 0, "", data)
}

// SendPermitJoinReq SendPermitJoinReq
func SendPermitJoinReq(terminalInfo config.TerminalInfo, duration string, TCSignificance string) {
	var ackOptionType = globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp
	var msgType = globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalPermitJoinReqEvent
	var optionType = globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp
	var data = duration + TCSignificance //duration + TCSignificance
	globallogger.Log.Infoln("devEUI : " + terminalInfo.DevEUI + " " + "SendPermitJoinReq")
	encAndSendDownMsg(terminalInfo, ackOptionType, msgType, optionType, 0, "", data)
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
		globallogger.Log.Errorln("devEUI : "+devEUI+" "+"ProcTerminalDelete GetTerminalInfoByDevEUI err: ", err)
	}
	if terminalInfo != nil {
		if constant.Constant.UsePostgres {
			err = models.DeleteTerminalPG(devEUI)
		} else {
			err = models.DeleteTerminal(devEUI)
		}
		if err != nil {
			globallogger.Log.Errorln("devEUI : "+devEUI+" "+"ProcTerminalDelete DeleteTerminal err: ", err)
		} else {
			var ackOptionType = globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp
			var msgType = globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalLeaveReqEvent
			var optionType = globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp
			var data = "00" //RemoveChildren/Rejoin
			if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
				data = devEUI + data
			}
			globallogger.Log.Warnln("devEUI : " + devEUI + " " + "ProcTerminalDelete terminal leave req")
			encAndSendDownMsg(*terminalInfo, ackOptionType, msgType, optionType, 0, "", data)
			DeleteTerminalTimer(devEUI)

			var key = constant.Constant.REDIS.ZigbeeRedisCtrlMsgKey + terminalInfo.APMac + "_" + terminalInfo.ModuleID
			IfCtrlMsgExistThenDeleteAll(key, devEUI, terminalInfo.APMac)
			key = constant.Constant.REDIS.ZigbeeRedisDataMsgKey + terminalInfo.APMac + "_" + terminalInfo.ModuleID
			IfDataMsgExistThenDeleteAll(key, devEUI, terminalInfo.APMac)
		}
	} else {
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "ProcTerminalDelete terminal is not exist")
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
		globallogger.Log.Errorln("OIDIndex : "+OIDIndex+" "+"ProcTerminalDeleteByOIDIndex GetTerminalInfoByOIDIndex err: ", err)
	}
	if terminalInfo != nil {
		if constant.Constant.UsePostgres {
			err = models.DeleteTerminalPG(terminalInfo.DevEUI)
		} else {
			err = models.DeleteTerminal(terminalInfo.DevEUI)
		}
		if err != nil {
			globallogger.Log.Errorln("devEUI : "+terminalInfo.DevEUI+" "+"ProcTerminalDeleteByOIDIndex DeleteTerminal err: ", err)
		} else {
			var ackOptionType = globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp
			var msgType = globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalLeaveReqEvent
			var optionType = globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp
			var data = "00" //RemoveChildren/Rejoin
			if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
				data = terminalInfo.DevEUI + data
			}
			globallogger.Log.Warnln("devEUI : " + terminalInfo.DevEUI + " " + "ProcTerminalDeleteByOIDIndex terminal leave req")
			encAndSendDownMsg(*terminalInfo, ackOptionType, msgType, optionType, 0, "", data)
			DeleteTerminalTimer(terminalInfo.DevEUI)

			var key = constant.Constant.REDIS.ZigbeeRedisCtrlMsgKey + terminalInfo.APMac + "_" + terminalInfo.ModuleID
			IfCtrlMsgExistThenDeleteAll(key, terminalInfo.DevEUI, terminalInfo.APMac)
			key = constant.Constant.REDIS.ZigbeeRedisDataMsgKey + terminalInfo.APMac + "_" + terminalInfo.ModuleID
			IfDataMsgExistThenDeleteAll(key, terminalInfo.DevEUI, terminalInfo.APMac)
		}
	} else {
		globallogger.Log.Warnln("OIDIndex : " + OIDIndex + " " + "ProcTerminalDeleteByOIDIndex terminal is not exist")
	}
}

// function sendBindOrUnbindReq(recvData, terminalInfo) {
//     var ackOptionType = globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp;
//     var optionType = globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp;
//     console.log("devEUI : " + terminalInfo.devEUI + " " + "sendBindOrUnbindReq");
//     var data = recvData.srcEndpoint + recvData.clusterID + recvData.dstAddrMode + recvData.dstAddress + recvData.dstEndpoint;
//     encAndSendDownMsg(terminalInfo, ackOptionType, recvData.msgType, optionType, null, "", data);
// }

// function sendBindOrUnbindTerminalReq(srcEndpoint, msgType, terminalInfo, bindOrUnbindInfo) {
//     var ackOptionType = globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp;
//     var optionType = globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp;
//     console.log("devEUI : " + terminalInfo.devEUI + " " + "sendBindOrUnbindTerminalReq");
//     if (srcEndpoint) {
//         console.log("devEUI : " + terminalInfo.devEUI + " " + "bind or unbind one, srcEndpoint: " + srcEndpoint);
//         bindOrUnbindInfo.forEach(function(item, bindInfoIndex) {
//             if (item.srcEndpoint && item.srcEndpoint === srcEndpoint) {
//                 item.clusterID.forEach(function(clusterID, clusterIDIndex) {
//                     var data = item.srcEndpoint + clusterID + item.dstAddrMode + item.dstAddress + item.dstEndpoint;
//                     encAndSendDownMsg(terminalInfo, ackOptionType, msgType, optionType, null, "", data);
//                 });
//             }
//         });
//     } else {
//         console.log("devEUI : " + terminalInfo.devEUI + " " + "bind or unbind all");
//         bindOrUnbindInfo.forEach(function(item, bindInfoIndex) {
//             item.clusterID.forEach(function(clusterID, clusterIDIndex) {
//                 var data = item.srcEndpoint + clusterID + item.dstAddrMode + item.dstAddress + item.dstEndpoint;
//                 encAndSendDownMsg(terminalInfo, ackOptionType, msgType, optionType, null, "", data);
//             });
//         });
//     }
// }

// SendBindTerminalReq SendBindTerminalReq
func SendBindTerminalReq(terminalInfo config.TerminalInfo) {
	var ackOptionType = globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp
	var optionType = globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp
	var msgType = globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalBindReqEvent
	// globallogger.Log.Infoln("devEUI : " + terminalInfo.DevEUI + " " + "SendBindTerminalReq")
	bindInfo := []config.BindInfo{}
	if constant.Constant.UsePostgres {
		json.Unmarshal([]byte(terminalInfo.BindInfoPG), &bindInfo)
	} else {
		bindInfo = terminalInfo.BindInfo
	}
	for bindInfoIndex, item := range bindInfo {
		for clusterIDIndex, clusterID := range item.ClusterID {
			globallogger.Log.Infoln("devEUI : "+terminalInfo.DevEUI+" "+"SendBindTerminalReq, bindInfoIndex: ", bindInfoIndex,
				" clusterIDIndex: ", clusterIDIndex, " clusterID: ", clusterID)
			var data = item.SrcEndpoint + clusterID + item.DstAddrMode + item.DstAddress + item.DstEndpoint
			if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
				data = terminalInfo.DevEUI + data
			}
			encAndSendDownMsg(terminalInfo, ackOptionType, msgType, optionType, 0, "", data)
		}
	}
	SendPermitJoinReq(terminalInfo, "00", "00")
}

// GetTerminalEndpoint GetTerminalEndpoint
func GetTerminalEndpoint(terminalInfo config.TerminalInfo) {
	globallogger.Log.Infoln("devEUI : " + terminalInfo.DevEUI + " " + "GetTerminalEndpoint")
	var ackOptionType = globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp
	var msgType = globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalEndpointReqEvent
	var optionType = globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp
	var data = ""
	encAndSendDownMsg(terminalInfo, ackOptionType, msgType, optionType, 0, "", data)
}

// SendTerminalDiscovery SendTerminalDiscovery
func SendTerminalDiscovery(terminalInfo config.TerminalInfo, endpoint string) {
	globallogger.Log.Infoln("devEUI : " + terminalInfo.DevEUI + " " + "SendTerminalDiscovery")
	var ackOptionType = globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp
	var msgType = globalmsgtype.MsgType.DOWNMsg.ZigbeeTerminalDiscoveryReqEvent
	var optionType = globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp
	encAndSendDownMsg(terminalInfo, ackOptionType, msgType, optionType, 0, "", endpoint)
}

// SendZHADownMsg SendZHADownMsg
func SendZHADownMsg(downMsg common.DownMsg) {
	devEUI := downMsg.DevEUI
	// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "SendZHADownMsg")
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
	}
	if err != nil {
		globallogger.Log.Errorln("devEUI : "+devEUI+" "+"SendZHADownMsg GetTerminalInfoByDevEUI err: ", err)
		return
	}
	if terminalInfo != nil {
		endpointTemp := pq.StringArray{}
		if constant.Constant.UsePostgres {
			endpointTemp = terminalInfo.EndpointPG
		} else {
			endpointTemp = terminalInfo.Endpoint
		}
		if len(endpointTemp) > 0 {
			var ackOptionType = globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp
			var msgType = downMsg.MsgType
			var optionType = globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNeedRsp
			tempProfileID := "0000" + strconv.FormatUint(uint64(downMsg.ProfileID), 16)
			ProfileID := tempProfileID[len(tempProfileID)-4:]
			tempClusterID := "0000" + strconv.FormatUint(uint64(downMsg.ClusterID), 16)
			ClusterID := tempClusterID[len(tempClusterID)-4:]
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
			data := ProfileID + ClusterID + DstEndpoint + downMsg.ZclData
			encAndSendDownMsg(*terminalInfo, ackOptionType, msgType, optionType, int(downMsg.SN), downMsg.MsgID, data)
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
	var ackOptionType = globalmsgtype.MsgType.ACKOPTIONType.ZigbeeAckOptionTypeRsp
	var msgType = globalmsgtype.MsgType.GENERALMsg.ZigbeeGeneralAck
	if jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0101 ||
		jsonInfo.TunnelHeader.Version == constant.Constant.UDPVERSION.Version0102 {
		msgType = globalmsgtype.MsgType.GENERALMsgV0101.ZigbeeGeneralAckV0101
	}
	var optionType = globalmsgtype.MsgType.OPTIONType.ZigbeeOptionTypeNoRsp
	var data = jsonInfo.MessageHeader.MsgType
	globallogger.Log.Infoln("devEUI : " + jsonInfo.MessagePayload.Address + " APMAC: " + jsonInfo.TunnelHeader.LinkInfo.APMac +
		" SendACK, msgType: " + jsonInfo.MessageHeader.MsgType + globalmsgtype.MsgType.GetMsgTypeUPMsgMeaning(jsonInfo.MessageHeader.MsgType))
	// fmt.Println("devEUI : " + jsonInfo.MessagePayload.Address + " " + "SendACK, msgType: " + jsonInfo.MessageHeader.MsgType + globalmsgtype.MsgType.GetMsgTypeUPMsgMeaning(jsonInfo.MessageHeader.MsgType))
	var encParameter = EncParameter{
		devEUI:        terminalInfo.DevEUI,
		APMac:         terminalInfo.APMac,
		moduleID:      terminalInfo.ModuleID,
		SN:            jsonInfo.TunnelHeader.FrameSN,
		ackOptionType: ackOptionType,
		msgType:       msgType,
		optionType:    optionType,
		data:          data,
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
			select {
			case <-time.After(time.Duration(constant.Constant.TIMER.ZigbeeTimerRetran) * time.Second):
				redisDataSerial, _ := json.Marshal(redisData)
				if constant.Constant.MultipleInstances {
					_, err := globalrediscache.RedisCache.SetRedis(key, 0, redisDataSerial)
					if err != nil {
						globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" CheckRedisAndReSend set redis error : ", err)
					} else {
						// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" CheckRedisAndReSend set redis success : ", res)
						checkAndSendDownMsg(0, devEUI, key, redisData.SN, APMac, true)
					}
				} else {
					_, err := globalmemorycache.MemoryCache.SetMemory(key, 0, redisDataSerial)
					if err != nil {
						globallogger.Log.Errorln("devEUI : "+devEUI+" "+"key: "+key+" CheckRedisAndReSend set memory error : ", err)
					} else {
						// globallogger.Log.Infoln("devEUI : "+devEUI+" "+"key: "+key+" CheckRedisAndReSend set memory success : ", res)
						checkAndSendDownMsg(0, devEUI, key, redisData.SN, APMac, true)
					}
				}
			}
		}()
	} else {
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "key: " + key + " CheckRedisAndReSend has already resend 3 times, do nothing ")
	}
}

// CheckRedisAndSend CheckRedisAndSend
func CheckRedisAndSend(devEUI string, key string, SN string, APMac string) {
	checkAndSendDownMsg(0, devEUI, key, SN, APMac, false)
}
