package general

import (
	"encoding/json"

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

// genAlarmsProcAlarm 处理Alarm命令
func genAlarmsProcAlarm(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.AlarmCommand)
	globallogger.Log.Infof("[devEUI: %v][genAlarmsProcAlarm] zclFrame: %+v", terminalInfo.DevEUI, Command)
	params := make(map[string]interface{}, 2)
	var key string
	params["terminalId"] = terminalInfo.DevEUI
	var value string
	var warning string = "未知告警"
	var bPublish = false
	kafkaMsg := publicstruct.DataReportMsg{
		OIDIndex:   terminalInfo.OIDIndex,
		DevSN:      terminalInfo.DevEUI,
		LinkType:   terminalInfo.ProfileID,
		DeviceType: terminalInfo.TmnType2,
	}
	switch cluster.ID(Command.ClusterIdentifier) {
	case cluster.TemperatureMeasurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				value = "0"
				warning = "温度过高"
			} else if Command.AlarmCode == 0x02 {
				value = "1"
				warning = "温度过低"
			}
			key = iotsmartspace.HeimanHS2AQPropertyTemperatureAlarm
			bPublish = true
			type appDataMsg struct {
				TemperatureWarning string `json:"temperatureWarning"`
			}
			kafkaMsg.AppData = appDataMsg{TemperatureWarning: warning}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarm] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.RelativeHumidityMeasurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				value = "0"
				warning = "湿度过高"
			} else if Command.AlarmCode == 0x02 {
				value = "1"
				warning = "湿度过低"
			}
			key = iotsmartspace.HeimanHS2AQPropertyHumidityAlarm
			bPublish = true
			type appDataMsg struct {
				HumidityWarning string `json:"humidityWarning"`
			}
			kafkaMsg.AppData = appDataMsg{HumidityWarning: warning}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarm] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANPM25Measurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				value = "0"
				warning = "PM2.5超标"
			}
			key = iotsmartspace.HeimanHS2AQPropertyPM25Alarm
			bPublish = true
			type appDataMsg struct {
				PM25Warning string `json:"PM25Warning"`
			}
			kafkaMsg.AppData = appDataMsg{PM25Warning: warning}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarm] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANFormaldehydeMeasurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				value = "0"
				warning = "甲醛超标"
			}
			key = iotsmartspace.HeimanHS2AQPropertyFormaldehydeAlarm
			bPublish = true
			type appDataMsg struct {
				HCHOWarning string `json:"HCHOWarning"`
			}
			kafkaMsg.AppData = appDataMsg{HCHOWarning: warning}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarm] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANAirQualityPM10:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x02 {
				value = "0"
				warning = "PM10超标"
			}
			key = iotsmartspace.HeimanHS2AQPropertyPM10Alarm
			bPublish = true
			type appDataMsg struct {
				PM10Warning string `json:"PM10Warning"`
			}
			kafkaMsg.AppData = appDataMsg{PM10Warning: warning}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarm] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANAirQualityTVOC:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				value = "1"
			} else if Command.AlarmCode == 0x02 {
				value = "0"
			}
			key = iotsmartspace.HeimanHS2AQPropertyTVOCAlarm
			bPublish = true
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarm] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANAirQualityMeasurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				value = "0"
				warning = "AQI超标"
			}
			key = iotsmartspace.HeimanHS2AQPropertyAQIAlarm
			bPublish = true
			type appDataMsg struct {
				AQIWarning string `json:"AQIWarning"`
			}
			kafkaMsg.AppData = appDataMsg{AQIWarning: warning}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarm] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.PowerConfiguration:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				value = "0"
				warning = "电池电压过低"
			}
			key = iotsmartspace.HeimanHS2AQPropertyBatteryAlarm
			bPublish = true
			type appDataMsg struct {
				BatteryWarning string `json:"batteryWarning"`
			}
			kafkaMsg.AppData = appDataMsg{BatteryWarning: warning}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarm] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	default:
		globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarm] invalid clusterID: %v", terminalInfo.DevEUI, Command.ClusterIdentifier)
	}
	params[key] = value
	if bPublish {
		if constant.Constant.Iotware {
			values := make(map[string]interface{}, 1)
			values[iotsmartspace.IotwarePropertyWarning] = value
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
		} else if constant.Constant.Iotedge {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
		} else if constant.Constant.Iotprivate {
			kafkaMsgByte, _ := json.Marshal(kafkaMsg)
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	}
}

// GenAlarmsProc 处理clusterID 0x0009即genAlarms属性消息
func GenAlarmsProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	// globallogger.Log.Infof("[devEUI: %v][GenAlarmsProc] Start......", terminalInfo.DevEUI)
	// globallogger.Log.Infof("[devEUI: %v][GenAlarmsProc] zclFrame: %+v", terminalInfo.DevEUI, zclFrame.Command)
	switch zclFrame.CommandName {
	case "Alarm":
		genAlarmsProcAlarm(terminalInfo, zclFrame.Command)
	case "GetAlarmResponse":
	default:
		globallogger.Log.Warnf("[devEUI: %v][GenAlarmsProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
