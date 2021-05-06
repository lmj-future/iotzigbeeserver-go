package controllers

import (
	"encoding/hex"
	"strconv"

	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/crc"
	"github.com/h3c/iotzigbeeserver-go/dgram"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalsocket"
	"github.com/h3c/iotzigbeeserver-go/metrics"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
)

//BindUDPPort BindUDPPort
func BindUDPPort(rcvUDPMsgport int) {
	globalsocket.ServiceSocket = dgram.CreateUDPSocket(rcvUDPMsgport)
}

func checkUDPMsgIsLegal(msg []byte) bool {
	frameLen, _ := strconv.ParseInt(hex.EncodeToString(append(msg[:0:0], msg[2:4]...)), 16, 0) //string转int
	if int(frameLen) == len(msg) {
		return true
	}
	globallogger.Log.Warnln("[ WARNING ! WARNING ! WARNING ! ] Invalid UDP message, please check ! msg length:", len(msg))
	return false
}

//UDPMsgProc UDPMsgProc
func UDPMsgProc(msg []byte, rinfo dgram.RInfo) {
	defer func() {
		err := recover()
		if err != nil {
			globallogger.Log.Errorln("UDPMsgProc err :", err)
		}
	}()
	if crc.Check(msg) {
		if checkUDPMsgIsLegal(msg) {
			globallogger.Log.Infoln("Receive UDP msg:", hex.EncodeToString(msg), "from", rinfo.Address, ":", rinfo.Port)
			procUDPMsg(msg, rinfo)
		} else {
			globallogger.Log.Warnln("[ WARNING ! WARNING ! WARNING ! ] Receive illegal UDP msg:", hex.EncodeToString(msg),
				"from", rinfo.Address, ":", rinfo.Port)
		}
	}
}

// CreateUDPServer UDP服务端创建并监听
func CreateUDPServer(rcvUDPMsgport int) {
	globalsocket.ServiceSocket = dgram.CreateUDPSocket(rcvUDPMsgport)

	go func() {
		defer globalsocket.ServiceSocket.Close()
		var data = make([]byte, 1024)
		for {
			msg, rinfo, err := globalsocket.ServiceSocket.Receive(data)
			if err != nil {
				globallogger.Log.Errorln(err.Error())
				continue
			}
			go UDPMsgProc(msg, rinfo)
			metrics.CountUdpReceiveTotal()
			metrics.CountUdpReceiveByAddress(rinfo.Address)
		}
	}()
}

func checkMsgTypeIsLegal(msgType string) bool {
	var isLegal = false
	switch msgType {
	case globalmsgtype.MsgType.UPMsg.ZigbeeTerminalJoinEvent,
		globalmsgtype.MsgType.UPMsg.ZigbeeTerminalLeaveEvent,
		globalmsgtype.MsgType.UPMsg.ZigbeeNetworkEvent,
		globalmsgtype.MsgType.UPMsg.ZigbeeModuleRemoveEvent,
		globalmsgtype.MsgType.UPMsg.ZigbeeEnvironmentChangeEvent,
		globalmsgtype.MsgType.UPMsg.ZigbeeDataUpEvent,
		globalmsgtype.MsgType.UPMsg.ZigbeeWholeNetworkNWKRspEvent,
		globalmsgtype.MsgType.UPMsg.ZigbeeWholeNetworkIEEERspEvent,
		globalmsgtype.MsgType.UPMsg.ZigbeeTerminalEndpointRspEvent,
		globalmsgtype.MsgType.UPMsg.ZigbeeTerminalBindRspEvent,
		globalmsgtype.MsgType.UPMsg.ZigbeeTerminalUnbindRspEvent,
		globalmsgtype.MsgType.UPMsg.ZigbeeTerminalNetworkRspEvent,
		globalmsgtype.MsgType.UPMsg.ZigbeeTerminalLeaveRspEvent,
		globalmsgtype.MsgType.UPMsg.ZigbeeTerminalPermitJoinRspEvent,
		globalmsgtype.MsgType.UPMsg.ZigbeeTerminalDiscoveryRspEvent,
		globalmsgtype.MsgType.UPMsg.ZigbeeTerminalCheckExistRspEvent,
		globalmsgtype.MsgType.GENERALMsg.ZigbeeGeneralAck,
		globalmsgtype.MsgType.GENERALMsg.ZigbeeGeneralKeepalive,
		globalmsgtype.MsgType.GENERALMsg.ZigbeeGeneralFailed,
		globalmsgtype.MsgType.GENERALMsgV0101.ZigbeeGeneralAckV0101,
		globalmsgtype.MsgType.GENERALMsgV0101.ZigbeeGeneralKeepaliveV0101,
		globalmsgtype.MsgType.GENERALMsgV0101.ZigbeeGeneralFailedV0101:
		isLegal = true
		// globallogger.Log.Infoln("==========================msgType is legal, go on========================")
	default:
		globallogger.Log.Warnln("[ WARNING ! WARNING ! WARNING ! ] msgType is illegal, please check ! msgType:", msgType)
	}
	return isLegal
}

/**
 * Process UDP message
 */
func procUDPMsg(udpMsg []byte, rinfo dgram.RInfo) {
	var jsonInfo = publicfunction.ParseUDPMsg(udpMsg, rinfo)
	if jsonInfo.TunnelHeader.DevTypeInfo.DevType == "T320M" {
		jsonInfo.TunnelHeader.LinkInfo.APMac = jsonInfo.TunnelHeader.LinkInfo.ACMac
	}
	metrics.CountUdpReceiveByGwSN(jsonInfo.TunnelHeader.LinkInfo.APMac)
	metrics.CountUdpReceiveByGwSNAndModuleID(jsonInfo.TunnelHeader.LinkInfo.APMac, jsonInfo.MessagePayload.ModuleID)
	switch jsonInfo.MessageHeader.MsgType { // UDP新协议版本msgType判断，新协议版本号0101
	case globalmsgtype.MsgType.GENERALMsgV0101.ZigbeeGeneralAckV0101,
		globalmsgtype.MsgType.GENERALMsgV0101.ZigbeeGeneralKeepaliveV0101,
		globalmsgtype.MsgType.GENERALMsgV0101.ZigbeeGeneralFailedV0101:
		if jsonInfo.TunnelHeader.Version != constant.Constant.UDPVERSION.Version0101 &&
			jsonInfo.TunnelHeader.Version != constant.Constant.UDPVERSION.Version0102 {
			globallogger.Log.Warnln("[procUDPMsg] version is illegal, please check ! version:", jsonInfo.TunnelHeader.Version)
			return
		}
	}
	if checkMsgTypeIsLegal(jsonInfo.MessageHeader.MsgType) { //msgtyp合法性校验
		if jsonInfo.MessageHeader.MsgType == globalmsgtype.MsgType.GENERALMsg.ZigbeeGeneralAck ||
			jsonInfo.MessageHeader.MsgType == globalmsgtype.MsgType.GENERALMsgV0101.ZigbeeGeneralAckV0101 { //处理ack消息
			procACKMsg(jsonInfo)
			metrics.CountUdpReceiveByLabel("ack")
		} else if jsonInfo.MessageHeader.MsgType == globalmsgtype.MsgType.GENERALMsg.ZigbeeGeneralKeepalive ||
			jsonInfo.MessageHeader.MsgType == globalmsgtype.MsgType.GENERALMsgV0101.ZigbeeGeneralKeepaliveV0101 {
			keepAliveInterval, _ := strconv.ParseInt(hex.EncodeToString(append(udpMsg[:0:0], udpMsg[len(udpMsg)-4:len(udpMsg)-2]...)), 16, 0)
			publicfunction.UpdateSocketInfo(jsonInfo, keepAliveInterval)
			publicfunction.SendACK(jsonInfo)
			if constant.Constant.MultipleInstances {
				publicfunction.KeepAliveTimerRedisSet(jsonInfo.TunnelHeader.LinkInfo.APMac)
				procKeepAliveMsgRedis(jsonInfo.TunnelHeader.LinkInfo.APMac, int(keepAliveInterval))
			} else {
				publicfunction.KeepAliveTimerFreeCacheSet(jsonInfo.TunnelHeader.LinkInfo.APMac, int(keepAliveInterval))
				procKeepAliveMsgFreeCache(jsonInfo.TunnelHeader.LinkInfo.APMac, int(keepAliveInterval))
			}
			metrics.CountUdpReceiveByLabel("keepalive")
		} else {
			procAnyMsg(jsonInfo)
		}
	}
}
