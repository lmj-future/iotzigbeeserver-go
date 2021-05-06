package smartenergy

import (
	"encoding/json"
	"strconv"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/db/kafka"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
	uuid "github.com/satori/go.uuid"
)

func meteringProcAttribute(terminalInfo config.TerminalInfo, attributeName string, attribute *cluster.Attribute, bPublish bool,
	params map[string]interface{}, values map[string]interface{}) (bool, map[string]interface{}, map[string]interface{}) {
	var key string
	switch attributeName {
	case "CurrentSummationDelivered":
		currentSummationDelivered := attribute.Value.(uint64)
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
			key = iotsmartspace.HeimanESocketPropertyElectric
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
			key = iotsmartspace.HeimanSmartPlugPropertyElectric
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
			key = iotsmartspace.HonyarSocket000a0c3cPropertyCurretSummationDelivered
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
			key = iotsmartspace.HonyarSocketHY0106PropertyCurretSummationDelivered
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
			key = iotsmartspace.HonyarSocketHY0105PropertyCurretSummationDelivered
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
			key = iotsmartspace.HonyarSocketHY0106PropertyCurretSummationDelivered
		default:
			globallogger.Log.Warnf("[devEUI: %v][meteringProcAttribute] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
		// HEIMAN
		// params[key] = strconv.FormatFloat(float64(currentSummationDelivered)*0.1*0.001, 'f', 2, 64)
		params[key] = strconv.FormatFloat(float64(currentSummationDelivered)*0.001, 'f', 2, 64)
		values[iotsmartspace.IotwarePropertyElectric] = params[key]
		bPublish = true
	case "Status", "UnitofMeasure", "Multiplier", "Divisor", "SummationFormatting", "MeteringDeviceType":
	case "InstantaneousDemand":
		instantaneousDemand := attribute.Value.(int64)
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
			key = iotsmartspace.HeimanESocketPropertyPower
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
			key = iotsmartspace.HeimanSmartPlugPropertyPower
		default:
			globallogger.Log.Warnf("[devEUI: %v][meteringProcAttribute] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
		params[key] = strconv.FormatFloat(float64(instantaneousDemand)*0.1, 'f', 2, 64)
		values[iotsmartspace.IotwarePropertyPower] = params[key]
		bPublish = true
	default:
		globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "meteringProcAttribute unknow attributeName", attributeName)
	}
	return bPublish, params, values
}

func meteringProcMsg2Kafka(terminalInfo config.TerminalInfo, values map[string]interface{}) {
	kafkaMsg := publicstruct.DataReportMsg{
		OIDIndex:   terminalInfo.OIDIndex,
		DevSN:      terminalInfo.DevEUI,
		LinkType:   terminalInfo.ProfileID,
		DeviceType: terminalInfo.TmnType2,
	}
	if electric, ok := values[iotsmartspace.IotwarePropertyElectric]; ok {
		type appDataMsg struct {
			Electric string `json:"electric"`
		}
		kafkaMsg.AppData = appDataMsg{Electric: electric.(string) + "kW/h"}
		kafkaMsgByte, _ := json.Marshal(kafkaMsg)
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
	if power, ok := values[iotsmartspace.IotwarePropertyPower]; ok {
		type appDataMsg struct {
			Power string `json:"power"`
		}
		kafkaMsg.AppData = appDataMsg{Power: power.(string) + "W"}
		kafkaMsgByte, _ := json.Marshal(kafkaMsg)
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
}

//meteringProcReadRsp 处理readRsp（0x01）消息
func meteringProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	readAttributesRsp := command.(*cluster.ReadAttributesResponse)
	params := make(map[string]interface{}, 2)
	values := make(map[string]interface{}, 2)
	params["terminalId"] = terminalInfo.DevEUI
	var bPublish = false
	for _, v := range readAttributesRsp.ReadAttributeStatuses {
		if v.Status == cluster.ZclStatusSuccess {
			bPublish, params, values = meteringProcAttribute(terminalInfo, v.AttributeName, v.Attribute, bPublish, params, values)
		}
	}

	if bPublish {
		// iotsmartspace publish msg to app
		if constant.Constant.Iotware {
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
		} else if constant.Constant.Iotedge {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
		} else if constant.Constant.Iotprivate {
			meteringProcMsg2Kafka(terminalInfo, values)
		}
	}
}

// meteringProcWriteRsp 处理writeRsp（0x04）消息
func meteringProcWriteRsp(terminalInfo config.TerminalInfo, command interface{}, msgID interface{}, contentFrame *zcl.Frame) {
	Command := command.(*cluster.WriteAttributesResponse)
	globallogger.Log.Infof("[devEUI: %v][meteringProcWriteRsp]: command: %+v", terminalInfo.DevEUI, Command)
	code := 0
	message := "success"
	if Command.WriteAttributeStatuses[0].Status != cluster.ZclStatusSuccess {
		code = -1
		message = "failed"
	}
	params := make(map[string]interface{}, 4)
	params["code"] = code
	params["message"] = message
	params["terminalId"] = terminalInfo.DevEUI
	var key string
	var value string
	var attributeName string
	var bPublish = false
	value = "0"
	if contentFrame != nil && contentFrame.CommandName == "WriteAttributes" {
		contentCommand := contentFrame.Command.(*cluster.WriteAttributesCommand)
		if contentCommand.WriteAttributeRecords[0].Attribute.Value.(bool) {
			value = "1"
		}
		attributeName = contentCommand.WriteAttributeRecords[0].AttributeName
	}
	switch attributeName {
	case "HistoryElectricClear":
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
			key = iotsmartspace.HonyarSocketHY0105PropertyHistoryElectricClear
			bPublish = true
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
			key = iotsmartspace.HonyarSocketHY0106PropertyHistoryElectricClear
			bPublish = true
		}
	}
	params[key] = value
	if bPublish {
		if constant.Constant.Iotware {
			iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, message, msgID)
		} else if constant.Constant.Iotedge {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply, params, msgID)
		}
	}
}

// meteringProcReport 处理report（0x0a）消息
func meteringProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][meteringProcReport]: command: %+v", terminalInfo.DevEUI, Command)
	params := make(map[string]interface{}, 2)
	values := make(map[string]interface{}, 2)
	params["terminalId"] = terminalInfo.DevEUI
	var bPublish = false
	attributeReports := Command.AttributeReports
	for _, v := range attributeReports {
		bPublish, params, values = meteringProcAttribute(terminalInfo, v.AttributeName, v.Attribute, bPublish, params, values)
	}

	if bPublish {
		// iotsmartspace publish msg to app
		if constant.Constant.Iotware {
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
		} else if constant.Constant.Iotedge {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
		} else if constant.Constant.Iotprivate {
			meteringProcMsg2Kafka(terminalInfo, values)
		}
	}
}

// meteringProcConfigureReportingResponse 处理configureReportingResponse（0x07）消息
func meteringProcConfigureReportingResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ConfigureReportingResponse)
	for _, v := range Command.AttributeStatusRecords {
		if v.Status != cluster.ZclStatusSuccess {
			globallogger.Log.Infof("[devEUI: %v][meteringProcConfigureReportingResponse]: configReport failed: %x", terminalInfo.DevEUI, v.Status)
		}
	}
}

// MeteringProc 处理clusterID 0x0702属性消息
func MeteringProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame, msgID interface{}, contentFrame *zcl.Frame) {
	switch zclFrame.CommandName {
	case "ReadAttributesResponse":
		meteringProcReadRsp(terminalInfo, zclFrame.Command)
	case "WriteAttributesResponse":
		meteringProcWriteRsp(terminalInfo, zclFrame.Command, msgID, contentFrame)
	case "ConfigureReportingResponse":
		meteringProcConfigureReportingResponse(terminalInfo, zclFrame.Command)
	case "ReportAttributes":
		meteringProcReport(terminalInfo, zclFrame.Command)
	default:
		globallogger.Log.Warnf("[devEUI: %v][MeteringProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
