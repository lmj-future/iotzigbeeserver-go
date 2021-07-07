package securityandsafety

import (
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
)

// iasWarningDeviceProcDefaultResponse 处理defaultResponse（0x0b）消息
func iasWarningDeviceProcDefaultResponse(terminalInfo config.TerminalInfo, command interface{}, msgID interface{}) {
	Command := command.(*cluster.DefaultResponseCommand)
	globallogger.Log.Infof("[devEUI: %v][iasWarningDeviceProcDefaultResponse]: command: %+v", terminalInfo.DevEUI, Command)
	code := 0
	message := "success"
	if Command.Status != cluster.ZclStatusSuccess {
		code = -1
		message = "failed"
	}
	if constant.Constant.Iotware {
		iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, message, msgID)
	} else if constant.Constant.Iotedge {
		type warningDeviceDefaultResponse struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
			warningDeviceDefaultResponse{Code: code, Message: message}, msgID)
	}

}

// IASWarningDeviceProc 处理clusterID 0x0502即IASWarningDevice属性消息
func IASWarningDeviceProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame, msgID interface{}) {
	switch zclFrame.CommandName {
	case "DefaultResponse":
		iasWarningDeviceProcDefaultResponse(terminalInfo, zclFrame.Command, msgID)
	default:
		globallogger.Log.Warnf("[devEUI: %v][IASWarningDeviceProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
