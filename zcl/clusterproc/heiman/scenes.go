package heiman

import (
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zclmsgdown"
	uuid "github.com/satori/go.uuid"
)

// ScenesProc 处理clusterID 0xfc80属性消息
func ScenesProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	// globallogger.Log.Infof("[devEUI: %v][ScenesProc] Start......", terminalInfo.DevEUI)
	// globallogger.Log.Infof("[devEUI: %v][ScenesProc] zclFrame: %+v", terminalInfo.DevEUI, zclFrame)
	params := make(map[string]interface{}, 2)
	values := make(map[string]interface{}, 1)
	params["terminalId"] = terminalInfo.DevEUI
	var key string
	var bPublish = false
	switch zclFrame.CommandName {
	case "AtHome":
		key = iotsmartspace.HeimanSceneSwitchEM30PropertyAtHome
		values[iotsmartspace.IotwarePropertySwitchRight] = "1"
		bPublish = true
	case "GoOut":
		key = iotsmartspace.HeimanSceneSwitchEM30PropertyGoOut
		values[iotsmartspace.IotwarePropertySwitchLeft] = "1"
		bPublish = true
	case "Cinema":
		key = iotsmartspace.HeimanSceneSwitchEM30PropertyCinema
		values[iotsmartspace.IotwarePropertySwitchUp] = "1"
		bPublish = true
	case "Repast":
		key = iotsmartspace.HeimanSceneSwitchEM30PropertyRepast
		bPublish = true
	case "Sleep":
		key = iotsmartspace.HeimanSceneSwitchEM30PropertySleep
		values[iotsmartspace.IotwarePropertySwitchDown] = "1"
		bPublish = true
	default:
		globallogger.Log.Warnf("[devEUI: %v][ScenesProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
	params[key] = "1"
	// iotsmartspace publish msg to app
	if bPublish {
		if constant.Constant.Iotware {
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
		} else if constant.Constant.Iotedge {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
		}
	}
	cmd := common.Command{
		TransactionID:    zclFrame.TransactionSequenceNumber,
		Cmd:              zclFrame.CommandIdentifier,
		DstEndpointIndex: 0,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      terminalInfo.DevEUI,
		CommandType: common.HeimanScenesDefaultResponse,
		ClusterID:   0xfc80,
		Command:     cmd,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}
