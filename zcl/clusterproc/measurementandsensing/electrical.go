package measurementandsensing

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

func electricalMeasurementProcReadOrReportIotware(terminalInfo config.TerminalInfo, RMSVoltageValue string, RMSCurrentValue string, ActivePowerValue string) {
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket, constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyVoltageCurrentPower{RealTimePower: ActivePowerValue})
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
		var value iotsmartspace.IotwarePropertyVoltageCurrentPower = iotsmartspace.IotwarePropertyVoltageCurrentPower{}
		if RMSVoltageValue != "NaN" && RMSVoltageValue != "+Inf" {
			value.CurrentVoltage = RMSVoltageValue
		}
		if RMSCurrentValue != "NaN" && RMSCurrentValue != "+Inf" {
			value.Current = RMSCurrentValue
		}
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			value.RealTimePower = ActivePowerValue
		}
		iotsmartspace.PublishTelemetryUpIotware(terminalInfo, value)
	default:
		globallogger.Log.Warnf("[devEUI: %v][electricalMeasurementProcReadOrReportIotware] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
	}
}

func electricalMeasurementProcReadOrReportIotedge(terminalInfo config.TerminalInfo, RMSVoltageValue string, RMSCurrentValue string, ActivePowerValue string) {
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
				iotsmartspace.HeimanESocketPropertyPower{DevEUI: terminalInfo.DevEUI, Power: ActivePowerValue}, uuid.NewV4().String())
		}
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
				iotsmartspace.HeimanSmartPlugPropertyPower{DevEUI: terminalInfo.DevEUI, Power: ActivePowerValue}, uuid.NewV4().String())
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
		var value iotsmartspace.HonyarSocket000a0c3cPropertyVoltageCurrentPower = iotsmartspace.HonyarSocket000a0c3cPropertyVoltageCurrentPower{
			DevEUI: terminalInfo.DevEUI,
		}
		if RMSVoltageValue != "NaN" && RMSVoltageValue != "+Inf" {
			value.RMSVoltage = RMSVoltageValue
		}
		if RMSCurrentValue != "NaN" && RMSCurrentValue != "+Inf" {
			value.RMSCurrent = RMSCurrentValue
		}
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			value.ActivePower = ActivePowerValue
		}
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, value, uuid.NewV4().String())
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
		var value iotsmartspace.HonyarSocketHY0106PropertyVoltageCurrentPower = iotsmartspace.HonyarSocketHY0106PropertyVoltageCurrentPower{
			DevEUI: terminalInfo.DevEUI,
		}
		if RMSVoltageValue != "NaN" && RMSVoltageValue != "+Inf" {
			value.RMSVoltage = RMSVoltageValue
		}
		if RMSCurrentValue != "NaN" && RMSCurrentValue != "+Inf" {
			value.RMSCurrent = RMSCurrentValue
		}
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			value.ActivePower = ActivePowerValue
		}
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, value, uuid.NewV4().String())
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
		var value iotsmartspace.HonyarSocketHY0105PropertyVoltageCurrentPower = iotsmartspace.HonyarSocketHY0105PropertyVoltageCurrentPower{
			DevEUI: terminalInfo.DevEUI,
		}
		if RMSVoltageValue != "NaN" && RMSVoltageValue != "+Inf" {
			value.RMSVoltage = RMSVoltageValue
		}
		if RMSCurrentValue != "NaN" && RMSCurrentValue != "+Inf" {
			value.RMSCurrent = RMSCurrentValue
		}
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			value.ActivePower = ActivePowerValue
		}
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, value, uuid.NewV4().String())
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
		var value iotsmartspace.HonyarSocketHY0106PropertyVoltageCurrentPower = iotsmartspace.HonyarSocketHY0106PropertyVoltageCurrentPower{
			DevEUI: terminalInfo.DevEUI,
		}
		if RMSVoltageValue != "NaN" && RMSVoltageValue != "+Inf" {
			value.RMSVoltage = RMSVoltageValue
		}
		if RMSCurrentValue != "NaN" && RMSCurrentValue != "+Inf" {
			value.RMSCurrent = RMSCurrentValue
		}
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			value.ActivePower = ActivePowerValue
		}
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, value, uuid.NewV4().String())
	default:
		globallogger.Log.Warnf("[devEUI: %v][electricalMeasurementProcReadOrReportIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
	}
}

func electricalMeasurementProcReadOrReportIotprivate(terminalInfo config.TerminalInfo, RMSVoltageValue string, RMSCurrentValue string, ActivePowerValue string) {
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket, constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			type appDataMsg struct {
				Power string `json:"power"`
			}
			kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
				Time:       time.Now(),
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData:    appDataMsg{Power: ActivePowerValue + "W"},
			})
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
		if RMSVoltageValue != "NaN" && RMSVoltageValue != "+Inf" {
			type appDataMsg struct {
				Voltage string `json:"voltage"`
			}
			kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
				Time:       time.Now(),
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData:    appDataMsg{Voltage: RMSVoltageValue + "V"},
			})
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
		if RMSCurrentValue != "NaN" && RMSCurrentValue != "+Inf" {
			type appDataMsg struct {
				Current string `json:"current"`
			}
			kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
				Time:       time.Now(),
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData:    appDataMsg{Current: RMSCurrentValue + "A"},
			})
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
		if ActivePowerValue != "NaN" && ActivePowerValue != "+Inf" {
			type appDataMsg struct {
				Power string `json:"power"`
			}
			kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
				Time:       time.Now(),
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData:    appDataMsg{Power: ActivePowerValue + "W"},
			})
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	default:
		globallogger.Log.Warnf("[devEUI: %v][electricalMeasurementProcReadOrReportIotprivate] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
	}
}

func electricalMeasurementProcReadOrReport(terminalInfo config.TerminalInfo, RMSVoltage float64, ACVoltageDivisor float64,
	RMSCurrent float64, ACCurrentDivisor float64, ActivePower float64, ACPowerDivisor float64) {
	if constant.Constant.Iotware {
		electricalMeasurementProcReadOrReportIotware(terminalInfo, strconv.FormatFloat(RMSVoltage/ACVoltageDivisor, 'f', 2, 64),
			strconv.FormatFloat(RMSCurrent/ACCurrentDivisor, 'f', 2, 64), strconv.FormatFloat(ActivePower/ACPowerDivisor, 'f', 2, 64))
	} else if constant.Constant.Iotedge {
		electricalMeasurementProcReadOrReportIotedge(terminalInfo, strconv.FormatFloat(RMSVoltage/ACVoltageDivisor, 'f', 2, 64),
			strconv.FormatFloat(RMSCurrent/ACCurrentDivisor, 'f', 2, 64), strconv.FormatFloat(ActivePower/ACPowerDivisor, 'f', 2, 64))
	} else if constant.Constant.Iotprivate {
		electricalMeasurementProcReadOrReportIotprivate(terminalInfo, strconv.FormatFloat(RMSVoltage/ACVoltageDivisor, 'f', 2, 64),
			strconv.FormatFloat(RMSCurrent/ACCurrentDivisor, 'f', 2, 64), strconv.FormatFloat(ActivePower/ACPowerDivisor, 'f', 2, 64))
	}
}

//electricalMeasurementProcReadAttributesResponse 处理readRsp（0x01）消息
func electricalMeasurementProcReadAttributesResponse(terminalInfo config.TerminalInfo, command interface{}) {
	var RMSVoltage uint64
	var ACVoltageDivisor uint64
	var RMSCurrent uint64
	var ACCurrentDivisor uint64
	var ActivePower int64
	var ACPowerDivisor uint64
	for _, v := range command.(*cluster.ReadAttributesResponse).ReadAttributeStatuses {
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
	electricalMeasurementProcReadOrReport(terminalInfo, float64(RMSVoltage), float64(ACVoltageDivisor), float64(RMSCurrent),
		float64(ACCurrentDivisor), float64(ActivePower), float64(ACPowerDivisor))
}

// electricalMeasurementProcReport 处理report（0x0a）消息
func electricalMeasurementProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][electricalMeasurementProcReport]: command: %+v", terminalInfo.DevEUI, Command)
	var RMSVoltage uint64
	var ACVoltageDivisor uint64
	var RMSCurrent uint64
	var ACCurrentDivisor uint64
	var ActivePower int64
	var ACPowerDivisor uint64
	for _, v := range Command.AttributeReports {
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
	electricalMeasurementProcReadOrReport(terminalInfo, float64(RMSVoltage), float64(ACVoltageDivisor), float64(RMSCurrent),
		float64(ACCurrentDivisor), float64(ActivePower), float64(ACPowerDivisor))
}

// electricalMeasurementProcConfigureReportingResponse 处理configureReportingResponse（0x07）消息
func electricalMeasurementProcConfigureReportingResponse(terminalInfo config.TerminalInfo, command interface{}) {
	for _, v := range command.(*cluster.ConfigureReportingResponse).AttributeStatusRecords {
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
