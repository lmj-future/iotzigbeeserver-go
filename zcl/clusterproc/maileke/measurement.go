package maileke

import (
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
)

func measurementProcAttribute(terminalInfo config.TerminalInfo, attributeName string, attribute *cluster.Attribute) {
	switch attributeName {
	case "PM1MeasuredValue":
	case "PM25MeasuredValue":
		// PM25 := attribute.Value.(uint64)
	case "PM10MeasuredValue":
	default:
		globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "measurementProcAttribute unknow attributeName", attributeName)
	}
}

//measurementProcReadRsp 处理readRsp（0x01）消息
func measurementProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	readAttributesRsp := command.(*cluster.ReadAttributesResponse)
	for _, v := range readAttributesRsp.ReadAttributeStatuses {
		globallogger.Log.Infof("[devEUI: %v][measurementProcReadRsp]: readAttributesRsp: %+v", terminalInfo.DevEUI, v)
		if v.Status == cluster.ZclStatusSuccess {
			measurementProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
		}
	}
}

// measurementProcReport 处理report（0x0a）消息
func measurementProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][measurementProcReport]: command: %+v", terminalInfo.DevEUI, Command)
	attributeReports := Command.AttributeReports
	for _, v := range attributeReports {
		measurementProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
	}
}

// MeasurementProc 处理clusterID 0x0415属性消息
func MeasurementProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	switch zclFrame.CommandName {
	case "ReadAttributesResponse":
		measurementProcReadRsp(terminalInfo, zclFrame.Command)
	case "ReportAttributes":
		measurementProcReport(terminalInfo, zclFrame.Command)
	default:
		globallogger.Log.Warnf("[devEUI: %v][MeasurementProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
