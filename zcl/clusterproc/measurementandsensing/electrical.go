package measurementandsensing

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

func electricalMeasurementProcAttribute(terminalInfo config.TerminalInfo, RMSVoltage uint64, ACVoltageDivisor float64, RMSCurrent uint64,
	ACCurrentDivisor float64, ActivePower int64, ACPowerDivisor float64, bPublish bool,
	params map[string]interface{}, values map[string]interface{}) (bool, map[string]interface{}, map[string]interface{}) {
	RMSVoltageValue := strconv.FormatFloat(float64(RMSVoltage)/ACVoltageDivisor, 'f', 2, 64)
	RMSCurrentValue := strconv.FormatFloat(float64(RMSCurrent)/ACCurrentDivisor, 'f', 2, 64)
	ActivePowerValue := strconv.FormatFloat(float64(ActivePower)/ACPowerDivisor, 'f', 2, 64)
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			params[iotsmartspace.HeimanESocketPropertyPower] = ActivePowerValue
			values[iotsmartspace.IotwarePropertyPower] = ActivePowerValue
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			params[iotsmartspace.HeimanSmartPlugPropertyPower] = ActivePowerValue
			values[iotsmartspace.IotwarePropertyPower] = ActivePowerValue
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
		if RMSVoltageValue != "NaN" && RMSVoltageValue != "+Inf" {
			params[iotsmartspace.HonyarSocket000a0c3cPropertyRMSVoltage] = RMSVoltageValue
			values[iotsmartspace.IotwarePropertyVoltage] = RMSVoltageValue
			bPublish = true
		}
		if RMSCurrentValue != "NaN" && RMSCurrentValue != "+Inf" {
			params[iotsmartspace.HonyarSocket000a0c3cPropertyRMSCurrent] = RMSCurrentValue
			values[iotsmartspace.IotwarePropertyCurrent] = RMSCurrentValue
			bPublish = true
		}
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			params[iotsmartspace.HonyarSocket000a0c3cPropertyActivePower] = ActivePowerValue
			values[iotsmartspace.IotwarePropertyPower] = ActivePowerValue
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
		if RMSVoltageValue != "NaN" && RMSVoltageValue != "+Inf" {
			params[iotsmartspace.HonyarSocketHY0106PropertyRMSVoltage] = RMSVoltageValue
			values[iotsmartspace.IotwarePropertyVoltage] = RMSVoltageValue
			bPublish = true
		}
		if RMSCurrentValue != "NaN" && RMSCurrentValue != "+Inf" {
			params[iotsmartspace.HonyarSocketHY0106PropertyRMSCurrent] = RMSCurrentValue
			values[iotsmartspace.IotwarePropertyCurrent] = RMSCurrentValue
			bPublish = true
		}
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			params[iotsmartspace.HonyarSocketHY0106PropertyActivePower] = ActivePowerValue
			values[iotsmartspace.IotwarePropertyPower] = ActivePowerValue
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
		if RMSVoltageValue != "NaN" && RMSVoltageValue != "+Inf" {
			params[iotsmartspace.HonyarSocketHY0105PropertyRMSVoltage] = RMSVoltageValue
			values[iotsmartspace.IotwarePropertyVoltage] = RMSVoltageValue
			bPublish = true
		}
		if RMSCurrentValue != "NaN" && RMSCurrentValue != "+Inf" {
			params[iotsmartspace.HonyarSocketHY0105PropertyRMSCurrent] = RMSCurrentValue
			values[iotsmartspace.IotwarePropertyCurrent] = RMSCurrentValue
			bPublish = true
		}
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			params[iotsmartspace.HonyarSocketHY0105PropertyActivePower] = ActivePowerValue
			values[iotsmartspace.IotwarePropertyPower] = ActivePowerValue
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
		if RMSVoltageValue != "NaN" && RMSVoltageValue != "+Inf" {
			params[iotsmartspace.HonyarSocketHY0106PropertyRMSVoltage] = RMSVoltageValue
			values[iotsmartspace.IotwarePropertyVoltage] = RMSVoltageValue
			bPublish = true
		}
		if RMSCurrentValue != "NaN" && RMSCurrentValue != "+Inf" {
			params[iotsmartspace.HonyarSocketHY0106PropertyRMSCurrent] = RMSCurrentValue
			values[iotsmartspace.IotwarePropertyCurrent] = RMSCurrentValue
			bPublish = true
		}
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			params[iotsmartspace.HonyarSocketHY0106PropertyActivePower] = ActivePowerValue
			values[iotsmartspace.IotwarePropertyPower] = ActivePowerValue
			bPublish = true
		}
	default:
		globallogger.Log.Warnf("[devEUI: %v][electricalMeasurementProcAttribute] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
	}

	return bPublish, params, values
}

func electricalMeasurementProcMsg2Kafka(terminalInfo config.TerminalInfo, values map[string]interface{}) {
	kafkaMsg := publicstruct.DataReportMsg{
		OIDIndex:   terminalInfo.OIDIndex,
		DevSN:      terminalInfo.DevEUI,
		LinkType:   terminalInfo.ProfileID,
		DeviceType: terminalInfo.TmnType2,
	}
	if voltage, ok := values[iotsmartspace.IotwarePropertyVoltage]; ok {
		type appDataMsg struct {
			Voltage string `json:"voltage"`
		}
		kafkaMsg.AppData = appDataMsg{Voltage: voltage.(string) + "V"}
		kafkaMsgByte, _ := json.Marshal(kafkaMsg)
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
	if current, ok := values[iotsmartspace.IotwarePropertyCurrent]; ok {
		type appDataMsg struct {
			Current string `json:"current"`
		}
		kafkaMsg.AppData = appDataMsg{Current: current.(string) + "A"}
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

//electricalMeasurementProcReadAttributesResponse 处理readRsp（0x01）消息
func electricalMeasurementProcReadAttributesResponse(terminalInfo config.TerminalInfo, command interface{}) {
	readAttributesRsp := command.(*cluster.ReadAttributesResponse)
	params := make(map[string]interface{}, 3)
	values := make(map[string]interface{}, 3)
	var bPublish = false
	params["terminalId"] = terminalInfo.DevEUI
	var RMSVoltage uint64
	var ACVoltageDivisor uint64
	var RMSCurrent uint64
	var ACCurrentDivisor uint64
	var ActivePower int64
	var ACPowerDivisor uint64
	for _, v := range readAttributesRsp.ReadAttributeStatuses {
		if v.Status == cluster.ZclStatusSuccess {
			switch v.AttributeName {
			case "MeasurementType", "PowerFactor", "ACVoltageMultiplier", "ACCurrentMultiplier", "ACPowerMultiplier",
				"ACAlarmsMask", "ACVoltageOverload", "ACCurrentOverload", "ACActivePowerOverload":
				// globallogger.Log.Infof("[devEUI: %v][electricalMeasurementProcReadAttributesResponse]: %s: %+v", terminalInfo.DevEUI, attributeName, attribute.Value)
			case "RMSVoltage":
				RMSVoltage = v.Attribute.Value.(uint64)
			case "ACVoltageDivisor":
				ACVoltageDivisor = v.Attribute.Value.(uint64)
			case "RMSCurrent":
				RMSCurrent = v.Attribute.Value.(uint64)
			case "ACCurrentDivisor":
				ACCurrentDivisor = v.Attribute.Value.(uint64)
			case "ActivePower":
				ActivePower = v.Attribute.Value.(int64)
			case "ACPowerDivisor":
				ACPowerDivisor = v.Attribute.Value.(uint64)
			default:
				globallogger.Log.Warnf("[devEUI: %v][electricalMeasurementProcReadAttributesResponse] invalid attributeName: %v",
					terminalInfo.DevEUI, v.AttributeName)
			}
		}
	}
	bPublish, params, values = electricalMeasurementProcAttribute(terminalInfo, RMSVoltage, float64(ACVoltageDivisor),
		RMSCurrent, float64(ACCurrentDivisor), ActivePower, float64(ACPowerDivisor), bPublish, params, values)
	if bPublish {
		if constant.Constant.Iotware {
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
		} else if constant.Constant.Iotedge {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
		} else if constant.Constant.Iotprivate {
			electricalMeasurementProcMsg2Kafka(terminalInfo, values)
		}
	}
}

// electricalMeasurementProcReport 处理report（0x0a）消息
func electricalMeasurementProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][electricalMeasurementProcReport]: command: %+v", terminalInfo.DevEUI, Command)
	attributeReports := Command.AttributeReports
	params := make(map[string]interface{}, 3)
	values := make(map[string]interface{}, 3)
	var bPublish = false
	params["terminalId"] = terminalInfo.DevEUI
	var RMSVoltage uint64
	var ACVoltageDivisor uint64
	var RMSCurrent uint64
	var ACCurrentDivisor uint64
	var ActivePower int64
	var ACPowerDivisor uint64
	for _, v := range attributeReports {
		switch v.AttributeName {
		case "MeasurementType", "PowerFactor", "ACVoltageMultiplier", "ACCurrentMultiplier", "ACPowerMultiplier",
			"ACAlarmsMask", "ACVoltageOverload", "ACCurrentOverload", "ACActivePowerOverload":
			// globallogger.Log.Infof("[devEUI: %v][electricalMeasurementProcReadAttributesResponse]: %s: %+v", terminalInfo.DevEUI, attributeName, attribute.Value)
		case "RMSVoltage":
			RMSVoltage = v.Attribute.Value.(uint64)
		case "ACVoltageDivisor":
			ACVoltageDivisor = v.Attribute.Value.(uint64)
		case "RMSCurrent":
			RMSCurrent = v.Attribute.Value.(uint64)
		case "ACCurrentDivisor":
			ACCurrentDivisor = v.Attribute.Value.(uint64)
		case "ActivePower":
			ActivePower = v.Attribute.Value.(int64)
		case "ACPowerDivisor":
			ACPowerDivisor = v.Attribute.Value.(uint64)
		default:
			globallogger.Log.Warnf("[devEUI: %v][electricalMeasurementProcReadAttributesResponse] invalid attributeName: %v",
				terminalInfo.DevEUI, v.AttributeName)
		}
	}
	bPublish, params, values = electricalMeasurementProcAttribute(terminalInfo, RMSVoltage, float64(ACVoltageDivisor),
		RMSCurrent, float64(ACCurrentDivisor), ActivePower, float64(ACPowerDivisor), bPublish, params, values)
	if bPublish {
		if constant.Constant.Iotware {
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
		} else if constant.Constant.Iotedge {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
		} else if constant.Constant.Iotprivate {
			electricalMeasurementProcMsg2Kafka(terminalInfo, values)
		}
	}
}

// electricalMeasurementProcConfigureReportingResponse 处理configureReportingResponse（0x07）消息
func electricalMeasurementProcConfigureReportingResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ConfigureReportingResponse)
	for _, v := range Command.AttributeStatusRecords {
		if v.Status != cluster.ZclStatusSuccess {
			globallogger.Log.Infof("[devEUI: %v][electricalMeasurementProcConfigureReportingResponse]: configReport failed: %x", terminalInfo.DevEUI, v.Status)
		}
	}
}

// ElectricalMeasurementProc 处理clusterID 0x0b04属性消息
func ElectricalMeasurementProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	switch zclFrame.CommandName {
	case "ReadAttributesResponse":
		electricalMeasurementProcReadAttributesResponse(terminalInfo, zclFrame.Command)
	case "ConfigureReportingResponse":
		electricalMeasurementProcConfigureReportingResponse(terminalInfo, zclFrame.Command)
	case "ReportAttributes":
		electricalMeasurementProcReport(terminalInfo, zclFrame.Command)
	default:
		globallogger.Log.Warnf("[devEUI: %v][ElectricalMeasurementProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
