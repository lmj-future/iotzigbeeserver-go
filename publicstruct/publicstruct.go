package publicstruct

import (
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/dgram"
)

// DataReportMsg DataReportMsg
type DataReportMsg struct {
	OIDIndex   string      `json:"OIDIndex"`
	DevSN      string      `json:"devSN"`
	LinkType   string      `json:"linkType"`
	DeviceType string      `json:"deviceType"`
	AppData    interface{} `json:"appData"`
}

// NetworkInAckIotware NetworkInAckIotware
type NetworkInAckIotware struct {
	Device string `json:"device"`
	IsOK   bool   `json:"isOK"`
}

// EventDeviceDeleteIotware EventDeviceDeleteIotware
type EventDeviceDeleteIotware struct {
	DeviceList []string `json:"deviceList"`
}

// DataTemp DataTemp
type DataTemp struct {
	ID     interface{} `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// RPCIotware RPCIotware
type RPCIotware struct {
	Device string   `json:"device"`
	Data   DataTemp `json:"data"`
}

// DataRPCRsp DataRPCRsp
type DataRPCRsp struct {
	Result string `json:"result"`
}

// RPCRspIotware RPCRspIotware
type RPCRspIotware struct {
	Device string      `json:"device"`
	ID     interface{} `json:"id"`
	Data   DataRPCRsp  `json:"data"`
}

// GatewayIotware GatewayIotware
type GatewayIotware struct {
	FirstAddr  string `json:"firstAddr"`
	SecondAddr string `json:"secondAddr"`
	ThirdAddr  string `json:"thirdAddr"`
	PortID     string `json:"portID"`
}

// DataItem DataItem
type DataItem struct {
	TimeStamp int64       `json:"ts"`
	Values    interface{} `json:"values"`
}

// DeviceData DeviceData
type DeviceData struct {
	Type string     `json:"type"`
	Data []DataItem `json:"data"`
}

// TelemetryIotware TelemetryIotware
type TelemetryIotware struct {
	Gateway    GatewayIotware `json:"gateway"`
	DeviceData interface{}    `json:"deviceData"`
}

//RedisData RedisData
type RedisData struct {
	AckFlag    bool        `json:"ackFlag"`
	SN         string      `json:"sn"`
	SendBuf    string      `json:"sendBuf"`
	DevEUI     string      `json:"devEUI"`
	MsgType    string      `json:"msgType"`
	OptionType string      `json:"optionType"`
	DataSN     string      `json:"dataSN"`
	ReSendTime int         `json:"reSendTime"`
	Data       string      `json:"data"`
	MsgID      interface{} `json:"msgID"`
}

//TmnTypeAttr TmnTypeAttr
type TmnTypeAttr struct {
	TmnType   string
	Attribute config.Attribute
}

//JSONInfo JSONInfo
type JSONInfo struct {
	TunnelHeader   TunnelHeader   `json:"tunnelHeader"`   //传输头
	MessageHeader  MessageHeader  `json:"messageHeader"`  //应用头
	MessagePayload MessagePayload `json:"messagePayload"` //应用数据
	Rinfo          dgram.RInfo    `json:"rinfo"`
}

//TunnelHeader TunnelHeader
type TunnelHeader struct {
	Version         string
	FrameLen        string
	FrameSN         string
	LinkInfo        LinkInfo
	ExtendInfo      ExtendInfo
	SecInfo         SecInfo
	VenderInfo      VenderInfo
	DevTypeInfo     DevTypeInfo
	UserNameInfo    UserNameInfo
	TunnelHeaderLen int
}

//LinkInfo LinkInfo
type LinkInfo struct {
	AddrNum    string
	ACMac      string
	APMac      string
	T300ID     string
	FirstAddr  string
	SecondAddr string
	ThirdAddr  string
	Address    []AddrInfo
}

//AddrInfo AddrInfo
type AddrInfo struct {
	AddrInfo string
	MACAddr  string
	SN       string
}

//ExtendInfo ExtendInfo
type ExtendInfo struct {
	IsNeedUserName bool
	IsNeedDevType  bool
	IsNeedVender   bool
	IsNeedSec      bool
	AckOptionType  int
	ExtendData     string
}

//SecInfo SecInfo
type SecInfo struct {
	SecType    string
	SecDataLen string
	SecData    string
	SecID      string
}

//VenderInfo VenderInfo
type VenderInfo struct {
	VenderIDLen string
	VenderID    string
}

//DevTypeInfo DevTypeInfo
type DevTypeInfo struct {
	DevTypeLen string
	DevType    string
}

//UserNameInfo UserNameInfo
type UserNameInfo struct {
	UserNameLen string
	UserName    string
}

//MessageHeader MessageHeader
type MessageHeader struct {
	OptionType       string
	SN               string
	MsgType          string
	MessageHeaderLen int
}

//CtrlMsg CtrlMsg
type CtrlMsg struct {
	Addr    string
	SubAddr string
	PortID  string
	CtrlLen int
}

//MessagePayload MessagePayload
type MessagePayload struct {
	ModuleID string
	Ctrl     int
	Address  string
	SubAddr  string
	PortID   string
	Topic    string
	Data     string
}

//DataMsgType DataMsgType
type DataMsgType struct {
	UPMsg           UPMsg
	DOWNMsg         DOWNMsg
	GENERALMsg      GENERALMsg
	GENERALMsgV0101 GENERALMsgV0101
	ACKOPTIONType   ACKOPTIONType
	OPTIONType      OPTIONType
}

//UPMsg UPMsg
type UPMsg struct {
	ZigbeeTerminalJoinEvent          string
	ZigbeeTerminalLeaveEvent         string
	ZigbeeNetworkEvent               string
	ZigbeeModuleRemoveEvent          string
	ZigbeeEnvironmentChangeEvent     string
	ZigbeeDataUpEvent                string //zigbee data up, WLAN_ZB_AF_INCOMING_MSG
	ZigbeeWholeNetworkNWKRspEvent    string //ZDO_NWK_ADDR_RSP
	ZigbeeWholeNetworkIEEERspEvent   string //ZDO_IEEE_ADDR_RSP
	ZigbeeTerminalEndpointRspEvent   string //ZDO_ACTIVE_EP_RSP
	ZigbeeTerminalBindRspEvent       string //ZDO_BIND_RSP
	ZigbeeTerminalUnbindRspEvent     string //ZDO_UNBIND_RSP
	ZigbeeTerminalNetworkRspEvent    string //ZDO_MGMT_LQI_RSP
	ZigbeeTerminalLeaveRspEvent      string //ZDO_MGMT_LEAVE_RSP
	ZigbeeTerminalPermitJoinRspEvent string //ZDO_MGMT_PERMIT_JOIN_RSP
	ZigbeeTerminalDiscoveryRspEvent  string //ZDO_SIMPLE_DESC_RSP
	ZigbeeTerminalCheckExistRspEvent string
}

//DOWNMsg DOWNMsg
type DOWNMsg struct {
	ZigbeeCmdRequestEvent            string //zigbee cmd down, WLAN_ZB_AF_DATA_REQUEST
	ZigbeeWholeNetworkNWKReqEvent    string //ZDO_NWK_ADDR_REQ
	ZigbeeWholeNetworkIEEEReqEvent   string //ZDO_IEEE_ADDR_REQ
	ZigbeeTerminalEndpointReqEvent   string //ZDO_ACTIVE_EP_REQ
	ZigbeeTerminalBindReqEvent       string //ZDO_BIND_REQ
	ZigbeeTerminalUnbindReqEvent     string //ZDO_UNBIND_REQ
	ZigbeeTerminalNetworkReqEvent    string //ZDO_MGMT_LQI_REQ
	ZigbeeTerminalLeaveReqEvent      string //ZDO_MGMT_LEAVE_REQ
	ZigbeeTerminalPermitJoinReqEvent string //ZDO_MGMT_PERMIT_JOIN_REQ
	ZigbeeTerminalDiscoveryReqEvent  string //ZDO_SIMPLE_DESC_REQ
	ZigbeeTerminalCheckExistReqEvent string
}

//GENERALMsg GENERALMsg
type GENERALMsg struct {
	ZigbeeGeneralAck       string //通用ack回复
	ZigbeeGeneralKeepalive string //保活报文
	ZigbeeGeneralFailed    string //通用处理失败回复
}

//GENERALMsgV0101 GENERALMsgV0101
type GENERALMsgV0101 struct {
	ZigbeeGeneralAckV0101       string //通用ack回复
	ZigbeeGeneralKeepaliveV0101 string //保活报文
	ZigbeeGeneralFailedV0101    string //通用处理失败回复
}

//ACKOPTIONType ACKOPTIONType
type ACKOPTIONType struct {
	ZigbeeAckOptionTypeNoRsp   string //传输层报文，不需要应答
	ZigbeeAckOptionTypeNeedRsp string //传输层报文，需要应答
	ZigbeeAckOptionTypeRsp     string //传输层报文，应答报文
}

//OPTIONType OPTIONType
type OPTIONType struct {
	ZigbeeOptionTypeNoRsp   string //应用层报文，不需要应答
	ZigbeeOptionTypeNeedRsp string //应用层报文，需要应答
	ZigbeeOptionTypeRsp     string //应用层报文，应答报文
}

//ErrCode ErrCode
type ErrCode struct {
	ZigbeeErrorSuccess          string
	ZigbeeErrorFailed           string
	ZigbeeErrorModuleNotExist   string
	ZigbeeErrorModuleOffline    string
	ZigbeeErrorAPNotExist       string
	ZigbeeErrorAPOffline        string
	ZigbeeErrorTerminalNotExist string
	ZigbeeErrorChipNotMatch     string
	ZigbeeErrorAPBusy           string
	ZigbeeErrorTimeout          string
}

// GetMsgTypeUPMsgMeaning GetMsgTypeUPMsgMeaning
func (msgtp DataMsgType) GetMsgTypeUPMsgMeaning(msgType string) string {
	var meaning = ""
	switch msgType {
	case "2001":
		meaning = "( ZigbeeTerminalJoinEvent )"
	case "2002":
		meaning = "( ZigbeeTerminalLeaveEvent )"
	case "2003":
		meaning = "( ZigbeeNetworkEvent )"
	case "2004":
		meaning = "( ZigbeeModuleRemoveEvent )"
	case "2005":
		meaning = "( ZigbeeEnvironmentChangeEvent )"
	case "2102":
		meaning = "( ZigbeeDataUpEvent )"
	case "2202":
		meaning = "( ZigbeeWholeNetworkNWKRspEvent )"
	case "2204":
		meaning = "( ZigbeeWholeNetworkIEEERspEvent )"
	case "2206":
		meaning = "( ZigbeeTerminalEndpointRspEvent )"
	case "2208":
		meaning = "( ZigbeeTerminalBindRspEvent )"
	case "220a":
		meaning = "( ZigbeeTerminalUnbindRspEvent )"
	case "220c":
		meaning = "( ZigbeeTerminalNetworkRspEvent )"
	case "220e":
		meaning = "( ZigbeeTerminalLeaveRspEvent )"
	case "2210":
		meaning = "( ZigbeeTerminalPermitJoinRspEvent )"
	case "2212":
		meaning = "( ZigbeeTerminalDiscoveryRspEvent )"
	case "2214":
		meaning = "( ZigbeeTerminalCheckExistRspEvent )"
	case "0002", "2007":
		meaning = "( ZigbeeGeneralKeepalive )"
	case "0003", "2008":
		meaning = "( ZigbeeGeneralFailed )"
	default:
		meaning = "WRONG UPMsg MSGTYPE"
	}
	return meaning
}

// GetMsgTypeDOWNMsgMeaning GetMsgTypeDOWNMsgMeaning
func (msgtp DataMsgType) GetMsgTypeDOWNMsgMeaning(msgType string) string {
	var meaning = ""
	switch msgType {
	case "2101":
		meaning = "( ZigbeeCmdRequestEvent )"
	case "2201":
		meaning = "( ZigbeeWholeNetworkNWKReqEvent )"
	case "2203":
		meaning = "( ZigbeeWholeNetworkIEEEReqEvent )"
	case "2205":
		meaning = "( ZigbeeTerminalEndpointReqEvent )"
	case "2207":
		meaning = "( ZigbeeTerminalBindReqEvent )"
	case "2209":
		meaning = "( ZigbeeTerminalUnbindReqEvent )"
	case "220b":
		meaning = "( ZigbeeTerminalNetworkReqEvent )"
	case "220d":
		meaning = "( ZigbeeTerminalLeaveReqEvent )"
	case "220f":
		meaning = "( ZigbeeTerminalPermitJoinReqEvent )"
	case "2211":
		meaning = "( ZigbeeTerminalDiscoveryReqEvent )"
	case "2213":
		meaning = "( ZigbeeTerminalCheckExistReqEvent )"
	default:
		meaning = "WRONG DOWNMsg MSGTYPE"
	}
	return meaning
}

// GetMsgTypeGENERALMsgMeaning GetMsgTypeGENERALMsgMeaning
func (msgtp DataMsgType) GetMsgTypeGENERALMsgMeaning(msgType string) string {
	var meaning = ""
	switch msgType {
	case "0001", "2006":
		meaning = "( ZigbeeGeneralAck )"
		break
	case "0002", "2007":
		meaning = "( ZigbeeGeneralKeepalive )"
		break
	case "0003", "2008":
		meaning = "( ZigbeeGeneralFailed )"
		break
	default:
		meaning = "WRONG GENERALMsg MSGTYPE"
		break
	}
	return meaning
}
