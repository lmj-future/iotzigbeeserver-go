package securityandsafety

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/db/kafka"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/keepalive"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
	"github.com/h3c/iotzigbeeserver-go/zcl/zclmsgdown"
	uuid "github.com/satori/go.uuid"
)

// iasZoneProcWriteRsp 处理writeRsp（0x04）消息
func iasZoneProcWriteRsp(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.WriteAttributesResponse)
	globallogger.Log.Infof("[devEUI: %v][iasZoneProcWriteRsp]: command: %+v", terminalInfo.DevEUI, Command)
	if Command.WriteAttributeStatuses[0].Status == cluster.ZclStatusSuccess {
		zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
			MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
			DevEUI:      terminalInfo.DevEUI,
			CommandType: common.SensorEnrollRsp,
			ClusterID:   0x0500,
			Command: common.Command{
				DstEndpointIndex: 0,
			},
		})
	}
}

func iasZoneProcZoneStatusChangeNotificationAlarm1Iotware(terminalInfo config.TerminalInfo, alarm1 uint16) {
	globallogger.Log.Infof("[devEUI: %v][iasZoneProcZoneStatusChangeNotificationAlarm1Iotware]: alarm1: %v", terminalInfo.DevEUI, alarm1)
	var value string
	if alarm1 == 0x01 {
		value = "0" //报警
	} else {
		value = "1" //报警解除
	}
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM:
		iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyCOAlarm{COAlarm: value})
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalGASSensorHY0022:
		iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyGasAlarm{GasAlarm: value})
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM:
		if value == "0" {
			value = "1"
		} else {
			value = "0"
		}
		iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyMoveAlarm{MoveAlarm: value})
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalPIRSensorHY0027:
		iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyMoveAlarm{MoveAlarm: value})
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30:
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSmokeSensorHY0024:
		iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyFireAlarm{FireAlarm: value})
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalWarningDevice005b0e12:
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM:
		iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyWaterSensorAlarm{WaterSensorAlarm: value})
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30:
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM:
	default:
		globallogger.Log.Warnf("[devEUI: %v][iasZoneProcZoneStatusChangeNotificationAlarm1Iotware] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
	}
}

func iasZoneProcZoneStatusChangeNotificationAlarm1Iotedge(terminalInfo config.TerminalInfo, alarm1 uint16) {
	globallogger.Log.Infof("[devEUI: %v][iasZoneProcZoneStatusChangeNotificationAlarm1Iotedge]: alarm1: %v", terminalInfo.DevEUI, alarm1)
	var value string
	if alarm1 == 0x01 {
		value = "0" //报警
	} else {
		value = "1" //报警解除
	}
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM:
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
			iotsmartspace.HeimanCOSensorPropertyAlarm{DevEUI: terminalInfo.DevEUI, Alarm: value}, uuid.NewV4().String())
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalGASSensorHY0022:
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
			iotsmartspace.HeimanGASSensorPropertyAlarm{DevEUI: terminalInfo.DevEUI, Alarm: value}, uuid.NewV4().String())
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM:
		if value == "0" {
			value = "1"
		} else {
			value = "0"
		}
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
			iotsmartspace.HeimanPIRSensorPropertyAlarm{DevEUI: terminalInfo.DevEUI, Alarm: value}, uuid.NewV4().String())
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalPIRSensorHY0027:
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
			iotsmartspace.HeimanPIRSensorPropertyAlarm{DevEUI: terminalInfo.DevEUI, Alarm: value}, uuid.NewV4().String())
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30:
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSmokeSensorHY0024:
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
			iotsmartspace.HeimanSmokeSensorPropertyAlarm{DevEUI: terminalInfo.DevEUI, Alarm: value}, uuid.NewV4().String())
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalWarningDevice005b0e12:
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM:
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
			iotsmartspace.HeimanWaterSensorPropertyAlarm{DevEUI: terminalInfo.DevEUI, Alarm: value}, uuid.NewV4().String())
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30:
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM:
	default:
		globallogger.Log.Warnf("[devEUI: %v][iasZoneProcZoneStatusChangeNotificationAlarm1Iotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
	}
}

func iasZoneProcZoneStatusChangeNotificationAlarm1Iotprivate(terminalInfo config.TerminalInfo, alarm1 uint16) {
	globallogger.Log.Infof("[devEUI: %v][iasZoneProcZoneStatusChangeNotificationAlarm1Iotprivate]: alarm1: %v", terminalInfo.DevEUI, alarm1)
	var warning string = "未知告警"
	if alarm1 == 0x01 {
		warning = "告警"
	} else {
		warning = "告警解除"
	}
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM:
		warning = "一氧化碳" + warning
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalGASSensorHY0022:
		warning = "可燃气体" + warning
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM:
		if alarm1 == 0x00 {
			warning = "有人移动"
		} else {
			warning = "无人移动"
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalPIRSensorHY0027:
		if alarm1 == 0x01 {
			warning = "有人移动"
		} else {
			warning = "无人移动"
		}
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30:
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSmokeSensorHY0024:
		warning = "烟火灾" + warning
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalWarningDevice005b0e12:
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM:
		warning = "水浸" + warning
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30:
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM:
	default:
		globallogger.Log.Warnf("[devEUI: %v][iasZoneProcZoneStatusChangeNotificationAlarm1Iotprivate] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
	}
	type appDataMsg struct {
		Warning string `json:"warning"`
	}
	kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
		Time:       time.Now(),
		OIDIndex:   terminalInfo.OIDIndex,
		DevSN:      terminalInfo.DevEUI,
		LinkType:   terminalInfo.ProfileID,
		DeviceType: terminalInfo.TmnType2,
		AppData: appDataMsg{
			Warning: warning,
		},
	})
	kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
}

// iasZoneProcZoneStatusChangeNotificationAlarm1 处理alarm1告警
func iasZoneProcZoneStatusChangeNotificationAlarm1(terminalInfo config.TerminalInfo, alarm1 uint16) {
	if constant.Constant.Iotware {
		iasZoneProcZoneStatusChangeNotificationAlarm1Iotware(terminalInfo, alarm1)
	} else if constant.Constant.Iotedge {
		iasZoneProcZoneStatusChangeNotificationAlarm1Iotedge(terminalInfo, alarm1)
	} else if constant.Constant.Iotprivate {
		iasZoneProcZoneStatusChangeNotificationAlarm1Iotprivate(terminalInfo, alarm1)
	}
}

func iasZoneProcZoneStatusChangeNotificationTamperIotware(terminalInfo config.TerminalInfo, tamper uint16) {
	globallogger.Log.Infof("[devEUI: %v][iasZoneProcZoneStatusChangeNotificationTamperIotware]: tamper: %v", terminalInfo.DevEUI, tamper)
	var value string
	if tamper == 0x01 {
		value = "0" //拆除报警
	} else {
		value = "1" //报警解除
	}
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM, constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalPIRSensorHY0027, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalWarningDevice005b0e12:
		iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyDismantleAlarm{DismantleAlarm: value})
	default:
		globallogger.Log.Warnf("[devEUI: %v][iasZoneProcZoneStatusChangeNotificationTamperIotware] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
	}
}
func iasZoneProcZoneStatusChangeNotificationTamperIotedge(terminalInfo config.TerminalInfo, tamper uint16) {
	globallogger.Log.Infof("[devEUI: %v][iasZoneProcZoneStatusChangeNotificationTamperIotedge]: tamper: %v", terminalInfo.DevEUI, tamper)
	var value string
	if tamper == 0x01 {
		value = "0" //拆除报警
	} else {
		value = "1" //报警解除
	}
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalPIRSensorHY0027:
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
			iotsmartspace.HeimanPIRSensorPropertyTamper{DevEUI: terminalInfo.DevEUI, Tamper: value}, uuid.NewV4().String())
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM:
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
			iotsmartspace.HeimanWaterSensorPropertyTamper{DevEUI: terminalInfo.DevEUI, Tamper: value}, uuid.NewV4().String())
	default:
		globallogger.Log.Warnf("[devEUI: %v][iasZoneProcZoneStatusChangeNotificationTamperIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
	}
}
func iasZoneProcZoneStatusChangeNotificationTamperIotprivate(terminalInfo config.TerminalInfo, tamper uint16) {
	globallogger.Log.Infof("[devEUI: %v][iasZoneProcZoneStatusChangeNotificationTamperIotprivate]: tamper: %v", terminalInfo.DevEUI, tamper)
	var tamperEvident string = "未知告警"
	if tamper == 0x01 {
		tamperEvident = "设备已拆除"
	} else {
		tamperEvident = "设备已安装"
	}
	type appDataMsg struct {
		TamperEvident string `json:"tamper-evident"`
	}
	kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
		Time:       time.Now(),
		OIDIndex:   terminalInfo.OIDIndex,
		DevSN:      terminalInfo.DevEUI,
		LinkType:   terminalInfo.ProfileID,
		DeviceType: terminalInfo.TmnType2,
		AppData: appDataMsg{
			TamperEvident: tamperEvident,
		},
	})
	kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
}

// iasZoneProcZoneStatusChangeNotificationTamper 处理tamper告警
func iasZoneProcZoneStatusChangeNotificationTamper(terminalInfo config.TerminalInfo, tamper uint16) {
	if constant.Constant.Iotware {
		iasZoneProcZoneStatusChangeNotificationTamperIotware(terminalInfo, tamper)
	} else if constant.Constant.Iotedge {
		iasZoneProcZoneStatusChangeNotificationTamperIotedge(terminalInfo, tamper)
	} else if constant.Constant.Iotprivate {
		iasZoneProcZoneStatusChangeNotificationTamperIotprivate(terminalInfo, tamper)
	}
}

// iasZoneProcZoneStatusChangeNotification 处理command generated ZoneStatusChangeNotification消息
func iasZoneProcZoneStatusChangeNotification(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ZoneStatusChangeNotificationCommand)
	globallogger.Log.Infof("[devEUI: %v][iasZoneProcZoneStatusChangeNotification]: command: %+v", terminalInfo.DevEUI, Command)
	attribute := config.Attribute{
		ZoneStatusChangeNotification: config.ZoneStatusChangeNotification{
			ZoneStatus: config.ZoneStatus{
				Data:           Command.ZoneStatus,
				Alarm1:         Command.ZoneStatus & 1,
				Alarm2:         (Command.ZoneStatus >> 1) & 1,
				Tamper:         (Command.ZoneStatus >> 2) & 1,
				Battery:        (Command.ZoneStatus >> 3) & 1,
				SupervisionRpt: (Command.ZoneStatus >> 4) & 1,
				RestoreRpt:     (Command.ZoneStatus >> 5) & 1,
				Trouble:        (Command.ZoneStatus >> 6) & 1,
				ACMains:        (Command.ZoneStatus >> 7) & 1,
				Test:           (Command.ZoneStatus >> 8) & 1,
				BatteryDefect:  (Command.ZoneStatus >> 9) & 1,
			},
			ExtendedStatus: Command.ExtendedStatus,
			ZoneID:         Command.ZoneID,
			Delay:          Command.Delay,
		},
	}
	var tmnInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		tmnInfo, err = models.GetTerminalInfoByDevEUIPG(terminalInfo.DevEUI)
	} else {
		tmnInfo, err = models.GetTerminalInfoByDevEUI(terminalInfo.DevEUI)
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][iasZoneProcZoneStatusChangeNotification]: get terminal info err %+v", terminalInfo.DevEUI, err)
		return
	}
	if tmnInfo != nil {
		if tmnInfo.Attribute.ZoneStatusChangeNotification.ZoneStatus.Alarm1 != attribute.ZoneStatusChangeNotification.ZoneStatus.Alarm1 {
			iasZoneProcZoneStatusChangeNotificationAlarm1(terminalInfo, attribute.ZoneStatusChangeNotification.ZoneStatus.Alarm1)
		}
		if tmnInfo.Attribute.ZoneStatusChangeNotification.ZoneStatus.Tamper != attribute.ZoneStatusChangeNotification.ZoneStatus.Tamper {
			iasZoneProcZoneStatusChangeNotificationTamper(terminalInfo, attribute.ZoneStatusChangeNotification.ZoneStatus.Tamper)
		}
	} else {
		globallogger.Log.Infof("[devEUI: %v][iasZoneProcZoneStatusChangeNotification]: get terminal info nil", terminalInfo.DevEUI)
		iasZoneProcZoneStatusChangeNotificationAlarm1(terminalInfo, attribute.ZoneStatusChangeNotification.ZoneStatus.Alarm1)
		iasZoneProcZoneStatusChangeNotificationTamper(terminalInfo, attribute.ZoneStatusChangeNotification.ZoneStatus.Tamper)
	}
	if constant.Constant.UsePostgres {
		models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": terminalInfo.DevEUI}, map[string]interface{}{
			"alarm1": attribute.ZoneStatusChangeNotification.ZoneStatus.Alarm1,
			"tamper": attribute.ZoneStatusChangeNotification.ZoneStatus.Tamper,
		})
	} else {
		models.FindTerminalAndUpdate(bson.M{"devEUI": terminalInfo.DevEUI}, bson.M{"attribute": attribute})
	}
	if tmnInfo != nil {
		var keyBuilder strings.Builder
		if tmnInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
			keyBuilder.WriteString(tmnInfo.APMac)
			keyBuilder.WriteString(tmnInfo.ModuleID)
			keyBuilder.WriteString(tmnInfo.NwkAddr)
		} else {
			keyBuilder.WriteString(tmnInfo.APMac)
			keyBuilder.WriteString(tmnInfo.ModuleID)
			keyBuilder.WriteString(tmnInfo.DevEUI)
		}
		publicfunction.DeleteTerminalInfoListCache(keyBuilder.String())
	}
}

// iasZoneProcZoneEnrollReq 处理command generated ZoneEnrollRequest消息
func iasZoneProcZoneEnrollReq(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ZoneEnrollCommand)
	globallogger.Log.Infof("[devEUI: %v][iasZoneProcZoneEnrollReq]: command: %+v", terminalInfo.DevEUI, Command)
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      terminalInfo.DevEUI,
		CommandType: common.SensorEnrollRsp,
		ClusterID:   0x0500,
		Command: common.Command{
			DstEndpointIndex: 0,
		},
	})
}

func iasZoneProcKeepAlive(terminalInfo config.TerminalInfo, interval uint16) {
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalGASSensorHY0022:
		keepalive.ProcKeepAlive(terminalInfo, interval)
	default:
		globallogger.Log.Warnf("[iasZoneProcKeepAlive]: invalid tmnType: %s", terminalInfo.TmnType)
	}
}

// iasZoneProcReport 处理report（0x0a）消息
func iasZoneProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	iasZoneProcKeepAlive(terminalInfo, uint16(terminalInfo.Interval))
}

// iasZoneProcConfigureReportingResponse 处理configureReportingResponse（0x07）消息
func iasZoneProcConfigureReportingResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ConfigureReportingResponse)
	globallogger.Log.Infof("[devEUI: %v][iasZoneProcConfigureReportingResponse]: command: %+v", terminalInfo.DevEUI, Command)
	for _, v := range Command.AttributeStatusRecords {
		if v.Status != cluster.ZclStatusSuccess {
			globallogger.Log.Infof("[devEUI: %v][iasZoneProcConfigureReportingResponse]: configReport failed: %x", terminalInfo.DevEUI, v.Status)
		} else {
			// proc keepalive
			if v.AttributeID == 0 {
				iasZoneProcKeepAlive(terminalInfo, uint16(terminalInfo.Interval))
			}
		}
	}
}

// IASZoneProc 处理clusterID 0x0500即IASZone属性消息
func IASZoneProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	switch zclFrame.CommandName {
	case "WriteAttributesResponse":
		iasZoneProcWriteRsp(terminalInfo, zclFrame.Command)
	case "ConfigureReportingResponse":
		iasZoneProcConfigureReportingResponse(terminalInfo, zclFrame.Command)
	case "ReportAttributes":
		iasZoneProcReport(terminalInfo, zclFrame.Command)
	case "ZoneStatusChangeNotification":
		iasZoneProcZoneStatusChangeNotification(terminalInfo, zclFrame.Command)
	case "ZoneEnrollRequest":
		iasZoneProcZoneEnrollReq(terminalInfo, zclFrame.Command)
	default:
		globallogger.Log.Warnf("[devEUI: %v][IASZoneProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
