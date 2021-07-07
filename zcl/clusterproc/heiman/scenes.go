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
	if constant.Constant.Iotware {
		switch zclFrame.CommandName {
		case "AtHome":
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertySwitchRight{SwitchRight: "1"})
		case "GoOut":
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertySwitchLeft{SwitchLeft: "1"})
		case "Cinema":
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertySwitchUp{SwitchUp: "1"})
		case "Repast":
		case "Sleep":
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertySwitchDown{SwitchDown: "1"})
		default:
			globallogger.Log.Warnf("[devEUI: %v][ScenesProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
		}
	} else if constant.Constant.Iotedge {
		switch zclFrame.CommandName {
		case "AtHome":
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
				iotsmartspace.HeimanSceneSwitchEM30PropertyAtHome{DevEUI: terminalInfo.DevEUI, AtHome: "1"}, uuid.NewV4().String())
		case "GoOut":
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
				iotsmartspace.HeimanSceneSwitchEM30PropertyGoOut{DevEUI: terminalInfo.DevEUI, GoOut: "1"}, uuid.NewV4().String())
		case "Cinema":
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
				iotsmartspace.HeimanSceneSwitchEM30PropertyCinema{DevEUI: terminalInfo.DevEUI, Cinema: "1"}, uuid.NewV4().String())
		case "Repast":
		case "Sleep":
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
				iotsmartspace.HeimanSceneSwitchEM30PropertySleep{DevEUI: terminalInfo.DevEUI, Sleep: "1"}, uuid.NewV4().String())
		default:
			globallogger.Log.Warnf("[devEUI: %v][ScenesProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
		}
	}
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      terminalInfo.DevEUI,
		CommandType: common.HeimanScenesDefaultResponse,
		ClusterID:   0xfc80,
		Command: common.Command{
			TransactionID:    zclFrame.TransactionSequenceNumber,
			Cmd:              zclFrame.CommandIdentifier,
			DstEndpointIndex: 0,
		},
	})
}
