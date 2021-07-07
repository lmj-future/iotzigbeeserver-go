package smartenergy

import (
	"encoding/json"
	"strconv"
	"time"

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

func meteringProcReadRspIotware(terminalInfo config.TerminalInfo, command interface{}) {
	var value iotsmartspace.IotwarePropertyTotalConsumptionRealTimePower = iotsmartspace.IotwarePropertyTotalConsumptionRealTimePower{}
	var bPublish = false
	for _, v := range command.(*cluster.ReadAttributesResponse).ReadAttributeStatuses {
		if v.Status == cluster.ZclStatusSuccess {
			switch v.AttributeName {
			case "CurrentSummationDelivered":
				value.TotalConsumption = strconv.FormatFloat(float64(v.Attribute.Value.(uint64))*0.001, 'f', 2, 64)
				if terminalInfo.ManufacturerName == constant.Constant.MANUFACTURERNAME.HeiMan {
					value.TotalConsumption = strconv.FormatFloat(float64(v.Attribute.Value.(uint64))*0.1*0.001, 'f', 2, 64)
				}
				bPublish = true
			case "Status", "UnitofMeasure", "Multiplier", "Divisor", "SummationFormatting", "MeteringDeviceType":
			case "InstantaneousDemand":
				value.RealTimePower = strconv.FormatFloat(float64(v.Attribute.Value.(int64))*0.1, 'f', 2, 64)
				bPublish = true
			default:
				globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "meteringProcReadRspIotware unknow attributeName", v.AttributeName)
			}
		}
	}
	if bPublish {
		iotsmartspace.PublishTelemetryUpIotware(terminalInfo, value)
	}
}
func meteringProcReadRspIotedge(terminalInfo config.TerminalInfo, command interface{}) {
	for _, v := range command.(*cluster.ReadAttributesResponse).ReadAttributeStatuses {
		if v.Status == cluster.ZclStatusSuccess {
			switch v.AttributeName {
			case "CurrentSummationDelivered":
				currentSummationDelivered := strconv.FormatFloat(float64(v.Attribute.Value.(uint64))*0.001, 'f', 2, 64)
				// HEIMAN
				if terminalInfo.ManufacturerName == constant.Constant.MANUFACTURERNAME.HeiMan {
					currentSummationDelivered = strconv.FormatFloat(float64(v.Attribute.Value.(uint64))*0.1*0.001, 'f', 2, 64)
				}
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
					iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
						iotsmartspace.HeimanESocketPropertyElectric{DevEUI: terminalInfo.DevEUI, Electric: currentSummationDelivered}, uuid.NewV4().String())
				case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
					iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
						iotsmartspace.HeimanSmartPlugPropertyElectric{DevEUI: terminalInfo.DevEUI, Electric: currentSummationDelivered}, uuid.NewV4().String())
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
					iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
						iotsmartspace.HonyarSocket000a0c3cPropertyCurretSummationDelivered{
							DevEUI:                   terminalInfo.DevEUI,
							CurretSummationDelivered: currentSummationDelivered,
						}, uuid.NewV4().String())
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
					iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
						iotsmartspace.HonyarSocketHY0106PropertyCurretSummationDelivered{
							DevEUI:                   terminalInfo.DevEUI,
							CurretSummationDelivered: currentSummationDelivered,
						}, uuid.NewV4().String())
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
					iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
						iotsmartspace.HonyarSocketHY0105PropertyCurretSummationDelivered{
							DevEUI:                   terminalInfo.DevEUI,
							CurretSummationDelivered: currentSummationDelivered,
						}, uuid.NewV4().String())
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
					iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
						iotsmartspace.HonyarSocketHY0106PropertyCurretSummationDelivered{
							DevEUI:                   terminalInfo.DevEUI,
							CurretSummationDelivered: currentSummationDelivered,
						}, uuid.NewV4().String())
				default:
					globallogger.Log.Warnf("[devEUI: %v][meteringProcReadRspIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
				}
			case "Status", "UnitofMeasure", "Multiplier", "Divisor", "SummationFormatting", "MeteringDeviceType":
			case "InstantaneousDemand":
				instantaneousDemand := strconv.FormatFloat(float64(v.Attribute.Value.(int64))*0.1, 'f', 2, 64)
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
					iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
						iotsmartspace.HeimanESocketPropertyPower{DevEUI: terminalInfo.DevEUI, Power: instantaneousDemand}, uuid.NewV4().String())
				case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
					iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
						iotsmartspace.HeimanSmartPlugPropertyPower{DevEUI: terminalInfo.DevEUI, Power: instantaneousDemand}, uuid.NewV4().String())
				default:
					globallogger.Log.Warnf("[devEUI: %v][meteringProcReadRspIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
				}
			default:
				globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "meteringProcReadRspIotedge unknow attributeName", v.AttributeName)
			}
		}
	}
}
func meteringProcReadRspIotprivate(terminalInfo config.TerminalInfo, command interface{}) {
	for _, v := range command.(*cluster.ReadAttributesResponse).ReadAttributeStatuses {
		if v.Status == cluster.ZclStatusSuccess {
			switch v.AttributeName {
			case "CurrentSummationDelivered":
				type appDataMsg struct {
					Electric string `json:"electric"`
				}
				var electric = strconv.FormatFloat(float64(v.Attribute.Value.(uint64))*0.001, 'f', 2, 64) + "kW/h"
				// HEIMAN
				if terminalInfo.ManufacturerName == constant.Constant.MANUFACTURERNAME.HeiMan {
					electric = strconv.FormatFloat(float64(v.Attribute.Value.(uint64))*0.1*0.001, 'f', 2, 64) + "kW/h"
				}
				kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
					Time:       time.Now(),
					OIDIndex:   terminalInfo.OIDIndex,
					DevSN:      terminalInfo.DevEUI,
					LinkType:   terminalInfo.ProfileID,
					DeviceType: terminalInfo.TmnType2,
					AppData:    appDataMsg{Electric: electric},
				})
				kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
			case "Status", "UnitofMeasure", "Multiplier", "Divisor", "SummationFormatting", "MeteringDeviceType":
			case "InstantaneousDemand":
				type appDataMsg struct {
					Power string `json:"power"`
				}
				kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
					Time:       time.Now(),
					OIDIndex:   terminalInfo.OIDIndex,
					DevSN:      terminalInfo.DevEUI,
					LinkType:   terminalInfo.ProfileID,
					DeviceType: terminalInfo.TmnType2,
					AppData:    appDataMsg{Power: strconv.FormatFloat(float64(v.Attribute.Value.(int64))*0.1, 'f', 2, 64) + "W"},
				})
				kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
			default:
				globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "meteringProcReadRspIotprivate unknow attributeName", v.AttributeName)
			}
		}
	}
}

//meteringProcReadRsp 处理readRsp（0x01）消息
func meteringProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	if constant.Constant.Iotware {
		meteringProcReadRspIotware(terminalInfo, command)
	} else if constant.Constant.Iotedge {
		meteringProcReadRspIotedge(terminalInfo, command)
	} else if constant.Constant.Iotprivate {
		meteringProcReadRspIotprivate(terminalInfo, command)
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
	if constant.Constant.Iotware {
		iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, message, msgID)
	} else if constant.Constant.Iotedge {
		var value string
		value = "0"
		if contentFrame != nil && contentFrame.CommandName == "WriteAttributes" {
			contentCommand := contentFrame.Command.(*cluster.WriteAttributesCommand)
			if contentCommand.WriteAttributeRecords[0].Attribute.Value.(bool) {
				value = "1"
			}
			switch contentCommand.WriteAttributeRecords[0].AttributeName {
			case "HistoryElectricClear":
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
					iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
						iotsmartspace.HonyarSocketHY0105PropertyHistoryElectricClearWriteRsp{
							DevEUI:               terminalInfo.DevEUI,
							Code:                 code,
							Message:              message,
							HistoryElectricClear: value,
						}, msgID)
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
					iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
						iotsmartspace.HonyarSocketHY0106PropertyHistoryElectricClearWriteRsp{
							DevEUI:               terminalInfo.DevEUI,
							Code:                 code,
							Message:              message,
							HistoryElectricClear: value,
						}, msgID)
				}
			}
		}
	}
}

func meteringProcReportIotware(terminalInfo config.TerminalInfo, command interface{}) {
	var value iotsmartspace.IotwarePropertyTotalConsumptionRealTimePower = iotsmartspace.IotwarePropertyTotalConsumptionRealTimePower{}
	var bPublish = false
	for _, v := range command.(*cluster.ReportAttributesCommand).AttributeReports {
		switch v.AttributeName {
		case "CurrentSummationDelivered":
			value.TotalConsumption = strconv.FormatFloat(float64(v.Attribute.Value.(uint64))*0.001, 'f', 2, 64)
			if terminalInfo.ManufacturerName == constant.Constant.MANUFACTURERNAME.HeiMan {
				value.TotalConsumption = strconv.FormatFloat(float64(v.Attribute.Value.(uint64))*0.1*0.001, 'f', 2, 64)
			}
			bPublish = true
		case "Status", "UnitofMeasure", "Multiplier", "Divisor", "SummationFormatting", "MeteringDeviceType":
		case "InstantaneousDemand":
			value.RealTimePower = strconv.FormatFloat(float64(v.Attribute.Value.(int64))*0.1, 'f', 2, 64)
			bPublish = true
		default:
			globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "meteringProcReportIotware unknow attributeName", v.AttributeName)
		}
	}
	if bPublish {
		iotsmartspace.PublishTelemetryUpIotware(terminalInfo, value)
	}
}
func meteringProcReportIotedge(terminalInfo config.TerminalInfo, command interface{}) {
	for _, v := range command.(*cluster.ReportAttributesCommand).AttributeReports {
		switch v.AttributeName {
		case "CurrentSummationDelivered":
			currentSummationDelivered := strconv.FormatFloat(float64(v.Attribute.Value.(uint64))*0.001, 'f', 2, 64)
			// HEIMAN
			if terminalInfo.ManufacturerName == constant.Constant.MANUFACTURERNAME.HeiMan {
				currentSummationDelivered = strconv.FormatFloat(float64(v.Attribute.Value.(uint64))*0.1*0.001, 'f', 2, 64)
			}
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanESocketPropertyElectric{DevEUI: terminalInfo.DevEUI, Electric: currentSummationDelivered}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanSmartPlugPropertyElectric{DevEUI: terminalInfo.DevEUI, Electric: currentSummationDelivered}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarSocket000a0c3cPropertyCurretSummationDelivered{
						DevEUI:                   terminalInfo.DevEUI,
						CurretSummationDelivered: currentSummationDelivered,
					}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarSocketHY0106PropertyCurretSummationDelivered{
						DevEUI:                   terminalInfo.DevEUI,
						CurretSummationDelivered: currentSummationDelivered,
					}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarSocketHY0105PropertyCurretSummationDelivered{
						DevEUI:                   terminalInfo.DevEUI,
						CurretSummationDelivered: currentSummationDelivered,
					}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarSocketHY0106PropertyCurretSummationDelivered{
						DevEUI:                   terminalInfo.DevEUI,
						CurretSummationDelivered: currentSummationDelivered,
					}, uuid.NewV4().String())
			default:
				globallogger.Log.Warnf("[devEUI: %v][meteringProcReportIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
			}
		case "Status", "UnitofMeasure", "Multiplier", "Divisor", "SummationFormatting", "MeteringDeviceType":
		case "InstantaneousDemand":
			instantaneousDemand := strconv.FormatFloat(float64(v.Attribute.Value.(int64))*0.1, 'f', 2, 64)
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanESocketPropertyPower{DevEUI: terminalInfo.DevEUI, Power: instantaneousDemand}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanSmartPlugPropertyPower{DevEUI: terminalInfo.DevEUI, Power: instantaneousDemand}, uuid.NewV4().String())
			default:
				globallogger.Log.Warnf("[devEUI: %v][meteringProcReportIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
			}
		default:
			globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "meteringProcReportIotedge unknow attributeName", v.AttributeName)
		}
	}
}
func meteringProcReportIotprivate(terminalInfo config.TerminalInfo, command interface{}) {
	for _, v := range command.(*cluster.ReportAttributesCommand).AttributeReports {
		switch v.AttributeName {
		case "CurrentSummationDelivered":
			type appDataMsg struct {
				Electric string `json:"electric"`
			}
			var electric = strconv.FormatFloat(float64(v.Attribute.Value.(uint64))*0.001, 'f', 2, 64) + "kW/h"
			// HEIMAN
			if terminalInfo.ManufacturerName == constant.Constant.MANUFACTURERNAME.HeiMan {
				electric = strconv.FormatFloat(float64(v.Attribute.Value.(uint64))*0.1*0.001, 'f', 2, 64) + "kW/h"
			}
			kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
				Time:       time.Now(),
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData:    appDataMsg{Electric: electric},
			})
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		case "Status", "UnitofMeasure", "Multiplier", "Divisor", "SummationFormatting", "MeteringDeviceType":
		case "InstantaneousDemand":
			type appDataMsg struct {
				Power string `json:"power"`
			}
			kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
				Time:       time.Now(),
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData:    appDataMsg{Power: strconv.FormatFloat(float64(v.Attribute.Value.(int64))*0.1, 'f', 2, 64) + "W"},
			})
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		default:
			globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "meteringProcReportIotprivate unknow attributeName", v.AttributeName)
		}
	}
}

// meteringProcReport 处理report（0x0a）消息
func meteringProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	if constant.Constant.Iotware {
		meteringProcReportIotware(terminalInfo, command)
	} else if constant.Constant.Iotedge {
		meteringProcReportIotedge(terminalInfo, command)
	} else if constant.Constant.Iotprivate {
		meteringProcReportIotprivate(terminalInfo, command)
	}
}

// meteringProcConfigureReportingResponse 处理configureReportingResponse（0x07）消息
func meteringProcConfigureReportingResponse(terminalInfo config.TerminalInfo, command interface{}) {
	for _, v := range command.(*cluster.ConfigureReportingResponse).AttributeStatusRecords {
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
