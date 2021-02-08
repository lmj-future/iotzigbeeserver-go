package securityandsafety

import (
	"encoding/json"

	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/db/kafka"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/models"
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
		// globallogger.Log.Infof("[devEUI: %v][iasZoneProcWriteRsp]: WriteAttributeStatuses: %+v", terminalInfo.DevEUI, Command.WriteAttributeStatuses[0])
		cmd := common.Command{
			DstEndpointIndex: 0,
		}
		zclDownMsg := common.ZclDownMsg{
			MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
			DevEUI:      terminalInfo.DevEUI,
			CommandType: common.SensorEnrollRsp,
			ClusterID:   0x0500,
			Command:     cmd,
		}
		zclmsgdown.ProcZclDownMsg(zclDownMsg)
	}
}

// iasZoneProcZoneStatusChangeNotificationAlarm1 处理alarm1告警
func iasZoneProcZoneStatusChangeNotificationAlarm1(terminalInfo config.TerminalInfo, alarm1 uint16) {
	globallogger.Log.Infof("[devEUI: %v][iasZoneProcZoneStatusChangeNotificationAlarm1]: alarm1: %v", terminalInfo.DevEUI, alarm1)
	params := make(map[string]interface{}, 2)
	var key string
	var bPublish = false
	params["terminalId"] = terminalInfo.DevEUI
	var value string
	var warning string = "未知告警"
	if alarm1 == 0x01 {
		value = "0" //报警
		warning = "告警"
	} else {
		value = "1" //报警解除
		warning = "告警解除"
	}
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM:
		key = iotsmartspace.HeimanCOSensorPropertyAlarm
		bPublish = true
		warning = "一氧化碳" + warning
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM:
		key = iotsmartspace.HeimanGASSensorPropertyAlarm
		bPublish = true
		warning = "可燃气体" + warning
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM:
		key = iotsmartspace.HeimanPIRSensorPropertyAlarm
		if value == "0" {
			value = "1"
			warning = "有人移动"
		} else {
			value = "0"
			warning = "无人移动"
		}
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30:
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM:
		key = iotsmartspace.HeimanSmokeSensorPropertyAlarm
		bPublish = true
		warning = "烟火灾" + warning
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice:
		key = iotsmartspace.HeimanWarningDevicePropertyAlarmTime
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM:
		key = iotsmartspace.HeimanWaterSensorPropertyAlarm
		bPublish = true
		warning = "水浸" + warning
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30:
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM:
		bPublish = true
	default:
		globallogger.Log.Warnf("[devEUI: %v][iasZoneProcZoneStatusChangeNotificationAlarm1] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
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
			type appDataMsg struct {
				Warning string `json:"warning"`
			}
			kafkaMsg := publicstruct.DataReportMsg{
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData: appDataMsg{
					Warning: warning,
				},
			}
			kafkaMsgByte, _ := json.Marshal(kafkaMsg)
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	}
}

// iasZoneProcZoneStatusChangeNotificationTamper 处理tamper告警
func iasZoneProcZoneStatusChangeNotificationTamper(terminalInfo config.TerminalInfo, tamper uint16) {
	globallogger.Log.Infof("[devEUI: %v][iasZoneProcZoneStatusChangeNotificationTamper]: tamper: %v", terminalInfo.DevEUI, tamper)
	params := make(map[string]interface{}, 2)
	var key string
	var bPublish = false
	params["terminalId"] = terminalInfo.DevEUI
	var value string
	var tamperEvident string = "未知告警"
	if tamper == 0x01 {
		value = "0" //拆除报警
		tamperEvident = "设备已拆除"
	} else {
		value = "1" //报警解除
		tamperEvident = "设备已安装"
	}
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM:
		key = iotsmartspace.HeimanPIRSensorPropertyTamper
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30:
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM:
		key = iotsmartspace.HeimanWaterSensorPropertyTamper
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30:
		bPublish = true
	default:
		globallogger.Log.Warnf("[devEUI: %v][iasZoneProcZoneStatusChangeNotificationTamper] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
	}
	params[key] = value
	if bPublish {
		if constant.Constant.Iotware {
			values := make(map[string]interface{}, 1)
			values[iotsmartspace.IotwarePropertyAntiDismantle] = value
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
		} else if constant.Constant.Iotedge {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
		} else if constant.Constant.Iotprivate {
			type appDataMsg struct {
				TamperEvident string `json:"tamper-evident"`
			}
			kafkaMsg := publicstruct.DataReportMsg{
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData: appDataMsg{
					TamperEvident: tamperEvident,
				},
			}
			kafkaMsgByte, _ := json.Marshal(kafkaMsg)
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	}
}

// iasZoneProcZoneStatusChangeNotification 处理command generated ZoneStatusChangeNotification消息
func iasZoneProcZoneStatusChangeNotification(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ZoneStatusChangeNotificationCommand)
	globallogger.Log.Infof("[devEUI: %v][iasZoneProcZoneStatusChangeNotification]: command: %+v", terminalInfo.DevEUI, Command)
	type ZoneStatus struct {
		data           uint16
		alarm1         uint16
		alarm2         uint16
		tamper         uint16
		battery        uint16
		supervisionRpt uint16
		restoreRpt     uint16
		trouble        uint16
		ACMains        uint16
		test           uint16
		batteryDefect  uint16
	}
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
	// globallogger.Log.Infof("[devEUI: %v][iasZoneProcZoneStatusChangeNotification]: update terminal info: attribute %+v", terminalInfo.DevEUI, attribute)
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
		// globallogger.Log.Infof("[devEUI: %v][iasZoneProcZoneStatusChangeNotification]: get terminal info success tmnInfo %+v", terminalInfo.DevEUI, tmnInfo)
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
}

// iasZoneProcZoneEnrollReq 处理command generated ZoneEnrollRequest消息
func iasZoneProcZoneEnrollReq(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ZoneEnrollCommand)
	globallogger.Log.Infof("[devEUI: %v][iasZoneProcZoneEnrollReq]: command: %+v", terminalInfo.DevEUI, Command)
	cmd := common.Command{
		DstEndpointIndex: 0,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      terminalInfo.DevEUI,
		CommandType: common.SensorEnrollRsp,
		ClusterID:   0x0500,
		Command:     cmd,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

func procIASZoneProcRead(devEUI string, dstEndpointIndex int, clusterID uint16) {
	cmd := common.Command{
		DstEndpointIndex: dstEndpointIndex,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.ReadAttribute,
		ClusterID:   clusterID,
		Command:     cmd,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

func iasZoneProcKeepAlive(devEUI string, tmnType string, interval uint16) {
	switch tmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM:
		keepalive.ProcKeepAlive(devEUI, interval)
	default:
		globallogger.Log.Warnf("[iasZoneProcKeepAlive]: invalid tmnType: %s", tmnType)
	}
}

// iasZoneProcReport 处理report（0x0a）消息
func iasZoneProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	iasZoneProcKeepAlive(terminalInfo.DevEUI, terminalInfo.TmnType, uint16(terminalInfo.Interval))
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
				iasZoneProcKeepAlive(terminalInfo.DevEUI, terminalInfo.TmnType, uint16(terminalInfo.Interval))
			}
		}
	}
}

// IASZoneProc 处理clusterID 0x0500即IASZone属性消息
func IASZoneProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	// globallogger.Log.Infof("[devEUI: %v][IASZoneProc] Start......", terminalInfo.DevEUI)
	// globallogger.Log.Infof("[devEUI: %v][IASZoneProc] zclFrame: %+v", terminalInfo.DevEUI, zclFrame.Command)
	switch zclFrame.CommandName {
	case "ReadAttributesResponse":
	case "WriteAttributesResponse":
		iasZoneProcWriteRsp(terminalInfo, zclFrame.Command)
	case "DefaultResponse":
	case "ZoneStatusChangeNotification":
		iasZoneProcZoneStatusChangeNotification(terminalInfo, zclFrame.Command)
	case "ZoneEnrollRequest":
		iasZoneProcZoneEnrollReq(terminalInfo, zclFrame.Command)
	case "ConfigureReportingResponse":
		iasZoneProcConfigureReportingResponse(terminalInfo, zclFrame.Command)
	case "ReportAttributes":
		iasZoneProcReport(terminalInfo, zclFrame.Command)
	default:
		globallogger.Log.Warnf("[devEUI: %v][IASZoneProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
