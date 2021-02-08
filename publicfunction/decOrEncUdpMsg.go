package publicfunction

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/crc"
	"github.com/h3c/iotzigbeeserver-go/dgram"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
)

//ParseUDPMsg ParseUDPMsg
func ParseUDPMsg(udpMsg []byte, rinfo dgram.RInfo) publicstruct.JSONInfo {
	var tunnelHeader = parseTunnelHeader(udpMsg)
	var messageHeader = parseMessageHeader(udpMsg, tunnelHeader.TunnelHeaderLen)
	var messagePayload = parseMessagePayload(udpMsg, tunnelHeader.TunnelHeaderLen+messageHeader.MessageHeaderLen)
	var jsonInfo = publicstruct.JSONInfo{
		TunnelHeader:   tunnelHeader,   //传输头
		MessageHeader:  messageHeader,  //应用头
		MessagePayload: messagePayload, //应用数据
		Rinfo:          rinfo,
	}
	if messageHeader.MsgType == globalmsgtype.MsgType.GENERALMsg.ZigbeeGeneralKeepalive ||
		messageHeader.MsgType == globalmsgtype.MsgType.GENERALMsgV0101.ZigbeeGeneralKeepaliveV0101 {
		// globallogger.Log.Infoln("==========================KEEPALIVE===========================")
	} else if messageHeader.MsgType == globalmsgtype.MsgType.GENERALMsg.ZigbeeGeneralAck ||
		messageHeader.MsgType == globalmsgtype.MsgType.GENERALMsgV0101.ZigbeeGeneralAckV0101 {
		// globallogger.Log.Infoln("devEUI : " + messagePayload.Address + " " + "============================RECV ACK==========================")
	} else {
		// globallogger.Log.Infoln("devEUI : " + messagePayload.Address + " " + "========================parseUdpMsg start=====================")
		globallogger.Log.Infoln("devEUI : " + messagePayload.Address + fmt.Sprintf(" %+v", jsonInfo))
		// globallogger.Log.Infoln("devEUI : " + messagePayload.Address + fmt.Sprintf(" %+v", jsonInfo.TunnelHeader))
		// globallogger.Log.Infoln("devEUI : " + messagePayload.Address + fmt.Sprintf(" %+v", jsonInfo.MessageHeader))
		// globallogger.Log.Infoln("devEUI : " + messagePayload.Address + fmt.Sprintf(" %+v", jsonInfo.MessagePayload))
		// globallogger.Log.Infoln("devEUI : " + messagePayload.Address + fmt.Sprintf(" %+v", jsonInfo.Rinfo))
		// globallogger.Log.Infoln("devEUI : " + messagePayload.Address + " " + "========================parseUdpMsg end ======================")
	}

	return jsonInfo
}

func parseTunnelHeader(udpMsg []byte) publicstruct.TunnelHeader {
	var tunnelHeader = publicstruct.TunnelHeader{}
	tunnelHeader.Version = hex.EncodeToString(udpMsg[0:2])  //协议版本（2byte）
	tunnelHeader.FrameLen = hex.EncodeToString(udpMsg[2:4]) //帧长度（2byte）
	tunnelHeader.FrameSN = hex.EncodeToString(udpMsg[4:6])  //帧序列号（2byte）
	var linkInfo = publicstruct.LinkInfo{}                  //链路信息
	var addrInfo = publicstruct.AddrInfo{}                  //n级地址
	var linkInfoLen = 0
	linkInfo.AddrNum = hex.EncodeToString(udpMsg[6:7])        //地址个数
	addrNum, _ := strconv.ParseInt(linkInfo.AddrNum, 16, 0)   //字符转16进制数
	var address = make([]publicstruct.AddrInfo, int(addrNum)) //地址列表
	for index := 0; index < int(addrNum); index++ {
		addrInfo.AddrInfo = hex.EncodeToString(udpMsg[7+linkInfoLen : 8+linkInfoLen])
		addrLen, _ := strconv.ParseInt(addrInfo.AddrInfo, 16, 0) //n级地址长度
		var addrLength = int(addrLen) & 31                       //bit 4-0, length
		if (addrLen >> 5) == 0 {                                 //bit 7-5, 0b000,mac
			addrInfo.MACAddr = hex.EncodeToString(udpMsg[8+linkInfoLen : 8+linkInfoLen+addrLength])
		} else if (addrLen >> 5) == 1 { //bit 7-5, 0b001,sn
			addrInfo.SN = hex.EncodeToString(udpMsg[8+linkInfoLen : 8+linkInfoLen+addrLength])
		}
		address[index] = addrInfo
		addrInfo = publicstruct.AddrInfo{}
		linkInfoLen += (addrLength + 1)
	}
	if addrNum == 3 { //3级
		if address[0].SN != "" {
			linkInfo.ACMac = address[0].SN
		} else {
			linkInfo.ACMac = address[0].MACAddr
		}
		if address[1].SN != "" {
			linkInfo.APMac = address[1].SN
		} else {
			linkInfo.APMac = address[1].MACAddr
		}
		if address[2].SN != "" {
			linkInfo.T300ID = address[2].SN
		} else {
			linkInfo.T300ID = address[2].MACAddr
		}
		linkInfo.FirstAddr = linkInfo.ACMac
		linkInfo.SecondAddr = linkInfo.APMac
		linkInfo.ThirdAddr = linkInfo.T300ID
	} else if addrNum == 2 { //2级
		if address[0].SN != "" {
			linkInfo.ACMac = address[0].SN
		} else {
			linkInfo.ACMac = address[0].MACAddr
		}
		if address[1].SN != "" {
			linkInfo.APMac = address[1].SN
		} else {
			linkInfo.APMac = address[1].MACAddr
		}
		linkInfo.FirstAddr = linkInfo.ACMac
		linkInfo.SecondAddr = linkInfo.APMac
	} else if addrNum == 1 { //1级
		if address[0].SN != "" {
			linkInfo.APMac = address[0].SN
		} else {
			linkInfo.APMac = address[0].MACAddr
		}
		linkInfo.FirstAddr = linkInfo.APMac
	}
	linkInfo.Address = address
	tunnelHeader.LinkInfo = linkInfo           //��·��Ϣ   //链路信息
	var extendInfo = publicstruct.ExtendInfo{} //��չ�ֶ�,bit2:0b0������0b1���ܣ�bit1-0:0b00����ҪӦ��0b01��ҪӦ��0b10Ӧ��0b11����
	var extendData = hex.EncodeToString(udpMsg[7+linkInfoLen : 9+linkInfoLen])
	isNeedUserName, _ := strconv.ParseInt(extendData, 16, 0)
	isNeedUserName = isNeedUserName & 32
	isNeedDevType, _ := strconv.ParseInt(extendData, 16, 0)
	isNeedDevType = isNeedDevType & 16
	isNeedVender, _ := strconv.ParseInt(extendData, 16, 0)
	isNeedVender = isNeedVender & 8
	isNeedSec, _ := strconv.ParseInt(extendData, 16, 0)
	isNeedSec = isNeedVender & 4
	optionType, _ := strconv.ParseInt(extendData, 16, 0)
	optionType = optionType & 3
	//var isNeedVender = parseInt(extendData, 16) & 8
	//var isNeedSec = parseInt(extendData, 16) & 4
	//var optionType = parseInt(extendData, 16) & 3
	extendInfo.IsNeedUserName = (isNeedUserName == 32)
	extendInfo.IsNeedDevType = (isNeedDevType == 16)
	extendInfo.IsNeedVender = (isNeedVender == 8)
	extendInfo.IsNeedSec = (isNeedSec == 4)
	extendInfo.AckOptionType = int(optionType) //0,1,2
	extendInfo.ExtendData = extendData
	tunnelHeader.ExtendInfo = extendInfo //扩展字段
	var secInfo = publicstruct.SecInfo{}
	var secInfoLen = 0
	if extendInfo.IsNeedSec {
		secInfo.SecType = hex.EncodeToString(udpMsg[9+linkInfoLen : 10+linkInfoLen])
		secInfo.SecDataLen = hex.EncodeToString(udpMsg[10+linkInfoLen : 11+linkInfoLen])
		secDataLen, _ := strconv.ParseInt(secInfo.SecDataLen, 16, 0)
		secInfo.SecData = hex.EncodeToString(udpMsg[11+linkInfoLen : 11+linkInfoLen+int(secDataLen)])
		secInfo.SecID = hex.EncodeToString(udpMsg[11+linkInfoLen+int(secDataLen) : 12+linkInfoLen+int(secDataLen)])
		secInfoLen = 1 + 1 + int(secDataLen) + 1
	}
	var venderInfo = struct {
		VenderIDLen string
		VenderID    string
	}{}
	var venderInfoLen = 0
	if extendInfo.IsNeedVender {
		venderInfo.VenderIDLen = hex.EncodeToString(udpMsg[9+linkInfoLen+secInfoLen : 10+linkInfoLen+secInfoLen])
		venderIDLen, _ := strconv.ParseInt(venderInfo.VenderIDLen, 16, 0)
		//venderInfo.VenderID = hex.EncodeToString(udpMsg[10+linkInfoLen+secInfoLen : 10+linkInfoLen+secInfoLen+venderIDLen])
		venderInfo.VenderID = string(udpMsg[10+linkInfoLen+secInfoLen : 10+linkInfoLen+secInfoLen+int(venderIDLen)])
		venderInfoLen = 1 + int(venderIDLen)
	}
	var devTypeInfo = struct {
		DevTypeLen string
		DevType    string
	}{}
	var devTypeInfoLen = 0
	if extendInfo.IsNeedDevType {
		devTypeInfo.DevTypeLen = hex.EncodeToString(udpMsg[9+linkInfoLen+secInfoLen+venderInfoLen : 10+linkInfoLen+secInfoLen+venderInfoLen])
		devTypeLen, _ := strconv.ParseInt(devTypeInfo.DevTypeLen, 16, 0)
		devTypeInfo.DevType = string(udpMsg[10+linkInfoLen+secInfoLen+venderInfoLen : 10+linkInfoLen+secInfoLen+venderInfoLen+int(devTypeLen)])
		devTypeInfoLen = 1 + int(devTypeLen)
	}
	var userNameInfo = struct {
		UserNameLen string
		UserName    string
	}{}
	var userNameInfoLen = 0
	if extendInfo.IsNeedUserName {
		userNameInfo.UserNameLen = hex.EncodeToString(udpMsg[9+linkInfoLen+secInfoLen+venderInfoLen+devTypeInfoLen : 10+linkInfoLen+secInfoLen+venderInfoLen+devTypeInfoLen])
		userNameLen, _ := strconv.ParseInt(userNameInfo.UserNameLen, 16, 0)
		userNameInfo.UserName = string(udpMsg[10+linkInfoLen+secInfoLen+venderInfoLen+devTypeInfoLen : 10+linkInfoLen+secInfoLen+venderInfoLen+devTypeInfoLen+int(userNameLen)])
		userNameInfoLen = 1 + int(userNameLen)
	}
	tunnelHeader.SecInfo = secInfo       //������Ϣ    //加密字段
	tunnelHeader.VenderInfo = venderInfo //������Ϣ  //场景字段
	tunnelHeader.DevTypeInfo = devTypeInfo
	tunnelHeader.UserNameInfo = userNameInfo
	tunnelHeader.TunnelHeaderLen = 9 + linkInfoLen + secInfoLen + venderInfoLen + devTypeInfoLen + userNameInfoLen
	return tunnelHeader
}

func parseMessageHeader(udpMsg []byte, len int) publicstruct.MessageHeader {
	var messageHeader = publicstruct.MessageHeader{}
	messageHeader.OptionType = hex.EncodeToString(udpMsg[len : len+2]) //������,0b00����ҪӦ��0b01��ҪӦ��0b10Ӧ��0b11����
	messageHeader.SN = hex.EncodeToString(udpMsg[len+2 : len+4])       //����Ӧ�ò���ش����վܾ���ƥ��Ӧ��
	messageHeader.MsgType = hex.EncodeToString(udpMsg[len+4 : len+6])  //��Ϣ���ͣ���4Bit��0b0000����ϵͳ��0b0001����rfid0b0010����zigbee
	messageHeader.MessageHeaderLen = 6
	return messageHeader
}

func parseMessagePayload(udpMsg []byte, length int) publicstruct.MessagePayload {
	var messagePayload = publicstruct.MessagePayload{}
	messagePayload.ModuleID = hex.EncodeToString(udpMsg[length : length+1])
	//messagePayload.Ctrl = hex.EncodeToString(udpMsg[len+1 : len+2])
	ctrl, _ := strconv.ParseInt(hex.EncodeToString(udpMsg[length+1:length+2]), 16, 0)
	messagePayload.Ctrl = int(ctrl)
	var ctrlMsg = parseCtrl(udpMsg[length+2:], messagePayload.Ctrl)
	if len(ctrlMsg.Addr) == 4 {
		messagePayload.Address = ctrlMsg.Addr
	} else {
		messagePayload.Address = strings.ToUpper(ctrlMsg.Addr)
	}
	messagePayload.SubAddr = ctrlMsg.SubAddr
	messagePayload.PortID = ctrlMsg.PortID
	messagePayload.Topic = hex.EncodeToString(udpMsg[length+2+ctrlMsg.CtrlLen : length+2+ctrlMsg.CtrlLen+4])
	messagePayload.Data = hex.EncodeToString(udpMsg[length+2+ctrlMsg.CtrlLen+4 : len(udpMsg)-2])

	return messagePayload
}

func parseCtrl(udpMsg []byte, ctrl int) publicstruct.CtrlMsg {
	var reserve = (ctrl >> 6) & 3
	var portIDLen = (ctrl >> 4) & 3
	var subAddrLen = (ctrl >> 2) & 3
	var addrLen = ctrl & 3
	switch portIDLen {
	case 0:
		portIDLen = 0
	case 1:
		portIDLen = 2
	case 2:
		portIDLen = 4
	case 3:
		portIDLen = 8
	default:
		portIDLen = 0
		globallogger.Log.Warnln("invalid portIDLen: ", portIDLen)
	}
	switch subAddrLen {
	case 0:
		subAddrLen = 0
	case 1:
		subAddrLen = 2
	case 2:
		subAddrLen = 4
	case 3:
		subAddrLen = 8
	default:
		subAddrLen = 0
		globallogger.Log.Warnln("invalid subAddrLen: ", subAddrLen)
	}
	switch reserve {
	case 0: //����1
		switch addrLen {
		case 0:
			addrLen = 6
		case 1:
			addrLen = 8
		case 2:
			addrLen = 16
		case 3:
			addrLen = 20
		default:
			addrLen = 0
			globallogger.Log.Warnln("invalid addrLen: ", addrLen)
		}
	case 1: //����2
		switch addrLen {
		case 1:
			addrLen = 4
		case 2:
			addrLen = 2
		default:
			addrLen = 0
			globallogger.Log.Warnln("invalid addrLen: ", addrLen)
		}
	}
	var msg = publicstruct.CtrlMsg{}
	msg.Addr = hex.EncodeToString(udpMsg[0:addrLen])
	msg.SubAddr = hex.EncodeToString(udpMsg[addrLen : addrLen+subAddrLen])
	msg.PortID = hex.EncodeToString(udpMsg[addrLen+subAddrLen : addrLen+subAddrLen+portIDLen])
	msg.CtrlLen = addrLen + subAddrLen + portIDLen

	return msg
}

func encCtrl(addrLen int, subAddrLen int, portIDLen int) int {
	var reserve = 0
	switch addrLen {
	case 2, 4:
		reserve = 1
	case 6:
		addrLen = 0
	case 8:
		addrLen = 1
	case 16:
		addrLen = 2
	case 20:
		addrLen = 3
	default:
		globallogger.Log.Warnln("Invalid address length: ", addrLen)
	}

	switch subAddrLen {
	case 0:
		subAddrLen = 0
	case 2:
		subAddrLen = 1
	case 4:
		subAddrLen = 2
	case 8:
		subAddrLen = 3
	default:
		globallogger.Log.Warnln("Invalid sub address length: ", subAddrLen)
	}

	switch portIDLen {
	case 0:
		portIDLen = 0
	case 2:
		portIDLen = 1
	case 4:
		portIDLen = 2
	case 8:
		portIDLen = 3
	default:
		globallogger.Log.Warnln("Invalid port ID length: ", portIDLen)
	}

	return reserve<<6 | portIDLen<<4 | subAddrLen<<2 | addrLen
}

func getTunnelHeader(tunnelHeader TunnelHeaderSerial, devEUI string) []byte {
	var str = tunnelHeader.Version + tunnelHeader.FrameLen + tunnelHeader.FrameSN + tunnelHeader.LinkInfo +
		tunnelHeader.ExtendInfo + tunnelHeader.SecInfo + tunnelHeader.VenderInfo + tunnelHeader.DevTypeInfo + tunnelHeader.UserNameInfo
	if len(str)%2 != 0 {
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "getTunnelHeader WRONG STRING LENGTH!")
		str = "00"
	}
	tempBuffer, _ := hex.DecodeString(str)
	return tempBuffer //Buffer.from(string, "hex");
}

func getMessageHeader(messageHeader MessageHeaderSerial, devEUI string) []byte {
	var str = messageHeader.OptionType + messageHeader.SN + messageHeader.MsgType
	if len(str)%2 != 0 {
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "getMessageHeader WRONG STRING LENGTH!")
		str = "00"
	}
	tempBuffer, _ := hex.DecodeString(str)
	return tempBuffer //Buffer.from(string, "hex")
}

func getMessagePayload(messagePayload MessagePayloadSerial, devEUI string) []byte {
	var str = messagePayload.ModuleID + messagePayload.Ctrl + messagePayload.Address + messagePayload.Topic + messagePayload.Data
	if len(str)%2 != 0 {
		globallogger.Log.Warnln("devEUI : " + devEUI + " " + "getMessagePayload WRONG STRING LENGTH!")
		str = "00"
	}
	tempBuffer, _ := hex.DecodeString(str)
	return tempBuffer //Buffer.from(string, "hex");
}

//TunnelHeaderSerial TunnelHeaderSerial
type TunnelHeaderSerial struct {
	Version         string
	FrameLen        string
	FrameSN         string
	LinkInfo        string
	ExtendInfo      string
	SecInfo         string
	VenderInfo      string
	DevTypeInfo     string
	UserNameInfo    string
	TunnelHeaderLen string
}

//MessageHeaderSerial MessageHeaderSerial
type MessageHeaderSerial struct {
	OptionType       string
	SN               string
	MsgType          string
	MessageHeaderLen string
}

//MessagePayloadSerial MessagePayloadSerial
type MessagePayloadSerial struct {
	ModuleID string
	Ctrl     string
	Address  string
	SubAddr  string
	PortID   string
	Topic    string
	Data     string
}

//SendDownMsg SendDownMsg
type SendDownMsg struct {
	SN      string
	sendBuf []byte
}

func encDownBuf(encParameter EncParameter, UDPVersion string) SendDownMsg {
	var tunnelHeader = TunnelHeaderSerial{}
	var messageHeader = MessageHeaderSerial{}
	var messagePayload = MessagePayloadSerial{}
	if UDPVersion == "" {
		UDPVersion = "0001"
	}
	tunnelHeader.Version = UDPVersion
	tunnelHeader.FrameLen = "0000"
	tunnelHeader.FrameSN = encParameter.SN
	var linkInfo = "01"
	var APMacLen = len(encParameter.APMac) / 2
	if encParameter.T300ID != "" {
		linkInfo = "03"
		//var T300IDLen = len(encParameter.T300ID) / 2
	} else {
		if encParameter.ACMac != "" {
			linkInfo = "02"
		}
	}
	if encParameter.ACMac != "" {
		var ACMacLen = len(encParameter.ACMac) / 2
		if ACMacLen > 6 {
			//linkInfo = linkInfo + ("00" + (ACMacLen | (1 << 5)).toString(16)).slice(-2) + encParameter.ACMac
			tempValue := "00" + strconv.FormatInt(int64(ACMacLen|(1<<5)), 16)
			linkInfo = linkInfo + tempValue[len(tempValue)-2:] + encParameter.ACMac
		} else {
			//linkInfo = linkInfo + ("00" + ACMacLen.toString(16)).slice(-2) + encParameter.ACMac
			tempValue := "00" + strconv.FormatInt(int64(ACMacLen), 16)
			linkInfo = linkInfo + tempValue[len(tempValue)-2:] + encParameter.ACMac
		}
	}
	if APMacLen > 6 {
		//linkInfo = linkInfo + ("00" + (APMacLen | (1 << 5)).toString(16)).slice(-2) + encParameter.APMac
		tempValue := "00" + strconv.FormatInt(int64(APMacLen|(1<<5)), 16)
		linkInfo = linkInfo + tempValue[len(tempValue)-2:] + encParameter.APMac
	} else {
		//linkInfo = linkInfo + ("00" + APMacLen.toString(16)).slice(-2) + encParameter.APMac
		tempValue := "00" + strconv.FormatInt(int64(APMacLen), 16)
		linkInfo = linkInfo + tempValue[len(tempValue)-2:] + encParameter.APMac
	}
	if encParameter.T300ID != "" {
		var T300IDLen = len(encParameter.T300ID) / 2
		if T300IDLen > 8 {
			//linkInfo = linkInfo + ("00" + (T300IDLen | (1 << 5)).toString(16)).slice(-2) + encParameter.T300ID
			tempValue := "00" + strconv.FormatInt(int64(T300IDLen|(1<<5)), 16)
			linkInfo = linkInfo + tempValue[len(tempValue)-2:] + encParameter.T300ID
		} else {
			//linkInfo = linkInfo + ("00" + T300IDLen.toString(16)).slice(-2) + encParameter.T300ID
			tempValue := "00" + strconv.FormatInt(int64(T300IDLen), 16)
			linkInfo = linkInfo + tempValue[len(tempValue)-2:] + encParameter.T300ID
		}
	}

	tunnelHeader.LinkInfo = linkInfo
	tunnelHeader.ExtendInfo = encParameter.ackOptionType
	tunnelHeader.SecInfo = ""
	tunnelHeader.VenderInfo = ""
	tunnelHeader.DevTypeInfo = ""
	tunnelHeader.UserNameInfo = ""
	messageHeader.OptionType = encParameter.optionType
	messageHeader.SN = "0000"
	messageHeader.MsgType = encParameter.msgType
	messagePayload.ModuleID = encParameter.moduleID
	var ctrl int
	if UDPVersion == constant.Constant.UDPVERSION.Version0102 {
		ctrl = encCtrl(2, 0, 0)
	} else {
		ctrl = encCtrl(8, 0, 0)
	}
	tempCtrl := "00" + strconv.FormatInt(int64(ctrl), 16)
	messagePayload.Ctrl = tempCtrl[len(tempCtrl)-2:]
	messagePayload.Address = encParameter.devEUI
	messagePayload.Topic = "ffffffff"
	messagePayload.Data = encParameter.data

	if encParameter.msgType == globalmsgtype.MsgType.GENERALMsg.ZigbeeGeneralAck ||
		encParameter.msgType == globalmsgtype.MsgType.GENERALMsgV0101.ZigbeeGeneralAckV0101 {
		// globallogger.Log.Infoln("devEUI : " + messagePayload.Address + " " + "========================SEND ACK==============================")
	} else {
		// globallogger.Log.Infoln("devEUI : " + messagePayload.Address + " " + "=====================encDownBuf start=========================")
		globallogger.Log.Infoln("devEUI : " + messagePayload.Address + " down msg: " + fmt.Sprintf("%+v", tunnelHeader) +
			fmt.Sprintf("%+v", messageHeader) + fmt.Sprintf("%+v", messagePayload))
		// globallogger.Log.Infoln("devEUI : " + messagePayload.Address + " " + "=====================encDownBuf end  =========================")
	}

	var tunnelHeaderBuf = getTunnelHeader(tunnelHeader, encParameter.devEUI)
	var messageHeaderBuf = getMessageHeader(messageHeader, encParameter.devEUI)
	var messagePayloadBuf = getMessagePayload(messagePayload, encParameter.devEUI)
	var length = len(tunnelHeaderBuf) + len(messageHeaderBuf) + len(messagePayloadBuf) + 2
	//tunnelHeader.FrameLen = ("0000" + length.toString(16)).slice(-4)
	tempFrameLen := "0000" + strconv.FormatInt(int64(length), 16)
	tunnelHeader.FrameLen = tempFrameLen[len(tempFrameLen)-4:]
	tunnelHeaderBuf = getTunnelHeader(tunnelHeader, encParameter.devEUI)

	var returnData = SendDownMsg{}
	returnData.SN = tunnelHeader.FrameSN
	//var buf = Buffer.concat([tunnelHeaderBuf, messageHeaderBuf, messagePayloadBuf], length-2)
	var buf = []byte(string(tunnelHeaderBuf) + string(messageHeaderBuf) + string(messagePayloadBuf))[:length-2]
	var crc = crc.CRC(buf)
	//returnData.sendBuf = Buffer.concat([tunnelHeaderBuf, messageHeaderBuf, messagePayloadBuf, Buffer.from(crc,"hex")], length)
	tempBuffer, _ := hex.DecodeString(crc)
	returnData.sendBuf = []byte(string(tunnelHeaderBuf) + string(messageHeaderBuf) + string(messagePayloadBuf) + string(tempBuffer))[:length]

	return returnData
}
