package general

import (
	"encoding/json"
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

func genAlarmsProcAlarmIotware(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.AlarmCommand)
	globallogger.Log.Infof("[devEUI: %v][genAlarmsProcAlarmIotware] zclFrame: %+v", terminalInfo.DevEUI, Command)
	switch cluster.ID(Command.ClusterIdentifier) {
	case cluster.TemperatureMeasurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyTemperatureAlarm{TemperatureAlarm: "0"})
			} else if Command.AlarmCode == 0x02 {
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyTemperatureAlarm{TemperatureAlarm: "1"})
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotware] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.RelativeHumidityMeasurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyHumidityAlarm{HumidityAlarm: "0"})
			} else if Command.AlarmCode == 0x02 {
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyHumidityAlarm{HumidityAlarm: "1"})
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotware] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANPM25Measurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyPM25Alarm{PM25Alarm: "0"})
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotware] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANFormaldehydeMeasurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyHCHOAlarm{HCHOAlarm: "0"})
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotware] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANAirQualityPM10:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x02 {
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyPM10Alarm{PM10Alarm: "0"})
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotware] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANAirQualityTVOC:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			// if Command.AlarmCode == 0x01 {
			// 	// "1"
			// } else if Command.AlarmCode == 0x02 {
			// 	// "0"
			// }
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotware] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANAirQualityMeasurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyAQIAlarm{AQIAlarm: "0"})
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotware] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.PowerConfiguration:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyPowerAlarm{PowerAlarm: "0"})
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotware] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	default:
		globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotware] invalid clusterID: %v", terminalInfo.DevEUI, Command.ClusterIdentifier)
	}
}

func genAlarmsProcAlarmIotedge(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.AlarmCommand)
	globallogger.Log.Infof("[devEUI: %v][genAlarmsProcAlarmIotedge] zclFrame: %+v", terminalInfo.DevEUI, Command)
	switch cluster.ID(Command.ClusterIdentifier) {
	case cluster.TemperatureMeasurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyTemperatureAlarm{DevEUI: terminalInfo.DevEUI, TemperatureAlarm: "0"}, uuid.NewV4().String())
			} else if Command.AlarmCode == 0x02 {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyTemperatureAlarm{DevEUI: terminalInfo.DevEUI, TemperatureAlarm: "1"}, uuid.NewV4().String())
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.RelativeHumidityMeasurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyHumidityAlarm{DevEUI: terminalInfo.DevEUI, HumidityAlarm: "0"}, uuid.NewV4().String())
			} else if Command.AlarmCode == 0x02 {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyHumidityAlarm{DevEUI: terminalInfo.DevEUI, HumidityAlarm: "1"}, uuid.NewV4().String())
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANPM25Measurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyPM25Alarm{DevEUI: terminalInfo.DevEUI, PM25Alarm: "0"}, uuid.NewV4().String())
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANFormaldehydeMeasurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyFormaldehydeAlarm{DevEUI: terminalInfo.DevEUI, FormaldehydeAlarm: "0"}, uuid.NewV4().String())
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANAirQualityPM10:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x02 {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyPM10Alarm{DevEUI: terminalInfo.DevEUI, PM10Alarm: "0"}, uuid.NewV4().String())
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANAirQualityTVOC:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyTVOCAlarm{DevEUI: terminalInfo.DevEUI, TVOCAlarm: "1"}, uuid.NewV4().String())
			} else if Command.AlarmCode == 0x02 {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyTVOCAlarm{DevEUI: terminalInfo.DevEUI, TVOCAlarm: "0"}, uuid.NewV4().String())
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANAirQualityMeasurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyAQIAlarm{DevEUI: terminalInfo.DevEUI, AQIAlarm: "0"}, uuid.NewV4().String())
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.PowerConfiguration:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyBatteryAlarm{DevEUI: terminalInfo.DevEUI, BatteryAlarm: "0"}, uuid.NewV4().String())
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	default:
		globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotedge] invalid clusterID: %v", terminalInfo.DevEUI, Command.ClusterIdentifier)
	}
}

func genAlarmsProcAlarmIotprivate(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.AlarmCommand)
	globallogger.Log.Infof("[devEUI: %v][genAlarmsProcAlarmIotprivate] zclFrame: %+v", terminalInfo.DevEUI, Command)
	var warning string = "未知告警"
	kafkaMsg := publicstruct.DataReportMsg{
		Time:       time.Now(),
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
				warning = "温度过高"
			} else if Command.AlarmCode == 0x02 {
				warning = "温度过低"
			}
			type appDataMsg struct {
				TemperatureWarning string `json:"temperatureWarning"`
			}
			kafkaMsg.AppData = appDataMsg{TemperatureWarning: warning}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotprivate] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.RelativeHumidityMeasurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				warning = "湿度过高"
			} else if Command.AlarmCode == 0x02 {
				warning = "湿度过低"
			}
			type appDataMsg struct {
				HumidityWarning string `json:"humidityWarning"`
			}
			kafkaMsg.AppData = appDataMsg{HumidityWarning: warning}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotprivate] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANPM25Measurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				warning = "PM2.5超标"
			}
			type appDataMsg struct {
				PM25Warning string `json:"PM25Warning"`
			}
			kafkaMsg.AppData = appDataMsg{PM25Warning: warning}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotprivate] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANFormaldehydeMeasurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				warning = "甲醛超标"
			}
			type appDataMsg struct {
				HCHOWarning string `json:"HCHOWarning"`
			}
			kafkaMsg.AppData = appDataMsg{HCHOWarning: warning}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotprivate] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANAirQualityPM10:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x02 {
				warning = "PM10超标"
			}
			type appDataMsg struct {
				PM10Warning string `json:"PM10Warning"`
			}
			kafkaMsg.AppData = appDataMsg{PM10Warning: warning}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotprivate] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANAirQualityTVOC:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			// if Command.AlarmCode == 0x01 {
			// } else if Command.AlarmCode == 0x02 {
			// }
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotprivate] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.HEIMANAirQualityMeasurement:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				warning = "AQI超标"
			}
			type appDataMsg struct {
				AQIWarning string `json:"AQIWarning"`
			}
			kafkaMsg.AppData = appDataMsg{AQIWarning: warning}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotprivate] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	case cluster.PowerConfiguration:
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			if Command.AlarmCode == 0x01 {
				warning = "电池电压过低"
			}
			type appDataMsg struct {
				BatteryWarning string `json:"batteryWarning"`
			}
			kafkaMsg.AppData = appDataMsg{BatteryWarning: warning}
		default:
			globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotprivate] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
			return
		}
	default:
		globallogger.Log.Warnf("[devEUI: %v][genAlarmsProcAlarmIotprivate] invalid clusterID: %v", terminalInfo.DevEUI, Command.ClusterIdentifier)
		return
	}
	kafkaMsgByte, err := json.Marshal(kafkaMsg)
	if err == nil {
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
}

// genAlarmsProcAlarm 处理Alarm命令
func genAlarmsProcAlarm(terminalInfo config.TerminalInfo, command interface{}) {
	if constant.Constant.Iotware {
		genAlarmsProcAlarmIotware(terminalInfo, command)
	} else if constant.Constant.Iotedge {
		genAlarmsProcAlarmIotedge(terminalInfo, command)
	} else if constant.Constant.Iotprivate {
		genAlarmsProcAlarmIotprivate(terminalInfo, command)
	}
}

// GenAlarmsProc 处理clusterID 0x0009即genAlarms属性消息
func GenAlarmsProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	switch zclFrame.CommandName {
	case "Alarm":
		genAlarmsProcAlarm(terminalInfo, zclFrame.Command)
	case "GetAlarmResponse":
	default:
		globallogger.Log.Warnf("[devEUI: %v][GenAlarmsProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
