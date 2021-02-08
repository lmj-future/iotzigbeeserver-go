package globalmsgtype

import "github.com/h3c/iotzigbeeserver-go/publicstruct"

var dataMsgType publicstruct.DataMsgType

func init() {
	dataMsgType.UPMsg.ZigbeeTerminalJoinEvent = "2001"
	dataMsgType.UPMsg.ZigbeeTerminalLeaveEvent = "2002"
	dataMsgType.UPMsg.ZigbeeNetworkEvent = "2003"
	dataMsgType.UPMsg.ZigbeeModuleRemoveEvent = "2004"
	dataMsgType.UPMsg.ZigbeeEnvironmentChangeEvent = "2005"
	dataMsgType.UPMsg.ZigbeeDataUpEvent = "2102"                //zigbee data up, WLAN_ZB_AF_INCOMING_MSG
	dataMsgType.UPMsg.ZigbeeWholeNetworkNWKRspEvent = "2202"    //ZDO_NWK_ADDR_RSP
	dataMsgType.UPMsg.ZigbeeWholeNetworkIEEERspEvent = "2204"   //ZDO_IEEE_ADDR_RSP
	dataMsgType.UPMsg.ZigbeeTerminalEndpointRspEvent = "2206"   //ZDO_ACTIVE_EP_RSP
	dataMsgType.UPMsg.ZigbeeTerminalBindRspEvent = "2208"       //ZDO_BIND_RSP
	dataMsgType.UPMsg.ZigbeeTerminalUnbindRspEvent = "220a"     //ZDO_UNBIND_RSP
	dataMsgType.UPMsg.ZigbeeTerminalNetworkRspEvent = "220c"    //ZDO_MGMT_LQI_RSP
	dataMsgType.UPMsg.ZigbeeTerminalLeaveRspEvent = "220e"      //ZDO_MGMT_LEAVE_RSP
	dataMsgType.UPMsg.ZigbeeTerminalPermitJoinRspEvent = "2210" //ZDO_MGMT_PERMIT_JOIN_RSP
	dataMsgType.UPMsg.ZigbeeTerminalDiscoveryRspEvent = "2212"  //ZDO_SIMPLE_DESC_RSP
	dataMsgType.UPMsg.ZigbeeTerminalCheckExistRspEvent = "2214"

	dataMsgType.DOWNMsg.ZigbeeCmdRequestEvent = "2101"            //zigbee cmd down, WLAN_ZB_AF_DATA_REQUEST
	dataMsgType.DOWNMsg.ZigbeeWholeNetworkNWKReqEvent = "2201"    //ZDO_NWK_ADDR_REQ
	dataMsgType.DOWNMsg.ZigbeeWholeNetworkIEEEReqEvent = "2203"   //ZDO_IEEE_ADDR_REQ
	dataMsgType.DOWNMsg.ZigbeeTerminalEndpointReqEvent = "2205"   //ZDO_ACTIVE_EP_REQ
	dataMsgType.DOWNMsg.ZigbeeTerminalBindReqEvent = "2207"       //ZDO_BIND_REQ
	dataMsgType.DOWNMsg.ZigbeeTerminalUnbindReqEvent = "2209"     //ZDO_UNBIND_REQ
	dataMsgType.DOWNMsg.ZigbeeTerminalNetworkReqEvent = "220b"    //ZDO_MGMT_LQI_REQ
	dataMsgType.DOWNMsg.ZigbeeTerminalLeaveReqEvent = "220d"      //ZDO_MGMT_LEAVE_REQ
	dataMsgType.DOWNMsg.ZigbeeTerminalPermitJoinReqEvent = "220f" //ZDO_MGMT_PERMIT_JOIN_REQ
	dataMsgType.DOWNMsg.ZigbeeTerminalDiscoveryReqEvent = "2211"  //ZDO_SIMPLE_DESC_REQ
	dataMsgType.DOWNMsg.ZigbeeTerminalCheckExistReqEvent = "2213"

	dataMsgType.GENERALMsg.ZigbeeGeneralAck = "0001"       //通用ack回复
	dataMsgType.GENERALMsg.ZigbeeGeneralKeepalive = "0002" //保活报文
	dataMsgType.GENERALMsg.ZigbeeGeneralFailed = "0003"    //通用处理失败回复

	dataMsgType.GENERALMsgV0101.ZigbeeGeneralAckV0101 = "2006"       //通用ack回复
	dataMsgType.GENERALMsgV0101.ZigbeeGeneralKeepaliveV0101 = "2007" //保活报文
	dataMsgType.GENERALMsgV0101.ZigbeeGeneralFailedV0101 = "2008"    //通用处理失败回复

	dataMsgType.ACKOPTIONType.ZigbeeAckOptionTypeNoRsp = "0000"   //传输层报文，不需要应答
	dataMsgType.ACKOPTIONType.ZigbeeAckOptionTypeNeedRsp = "0001" //传输层报文，需要应答
	dataMsgType.ACKOPTIONType.ZigbeeAckOptionTypeRsp = "0002"     //传输层报文，应答报文

	dataMsgType.OPTIONType.ZigbeeOptionTypeNoRsp = "0000"   //应用层报文，不需要应答
	dataMsgType.OPTIONType.ZigbeeOptionTypeNeedRsp = "0001" //应用层报文，需要应答
	dataMsgType.OPTIONType.ZigbeeOptionTypeRsp = "0002"     //应用层报文，应答报文
}

// MsgType MsgType
var MsgType = &dataMsgType
