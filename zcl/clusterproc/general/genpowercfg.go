package general

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/db/kafka"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/keepalive"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
	"github.com/h3c/iotzigbeeserver-go/zcl/zclmsgdown"
	uuid "github.com/satori/go.uuid"
)

func genPowerCfgProcAttribute(terminalInfo config.TerminalInfo, attributeName string, attribute *cluster.Attribute) {
	switch attributeName {
	case "BatteryVoltage": //电量值
	case "BatteryPercentageRemaining": //电量百分比
		// iotsmartspace publish msg to app
		if constant.Constant.Iotware {
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyBatteryPercentage{BatteryPercentage: attribute.Value.(uint64) / 2})
		} else if constant.Constant.Iotedge {
			power := attribute.Value.(uint64) / 2
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanCOSensorPropertyLeftElectric{DevEUI: terminalInfo.DevEUI, LeftElectric: power}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyLeftElectric{DevEUI: terminalInfo.DevEUI, LeftElectric: power}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHTSensorPropertyLeftElectric{DevEUI: terminalInfo.DevEUI, LeftElectric: power}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalPIRSensorHY0027:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanPIRSensorPropertyLeftElectric{DevEUI: terminalInfo.DevEUI, LeftElectric: power}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanPIRSensorPropertyLeftElectric{DevEUI: terminalInfo.DevEUI, LeftElectric: power}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSmokeSensorHY0024:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanSmokeSensorPropertyLeftElectric{DevEUI: terminalInfo.DevEUI, LeftElectric: power}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalWarningDevice005b0e12:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanWarningDevicePropertyLeftElectric{DevEUI: terminalInfo.DevEUI, LeftElectric: power}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanWaterSensorPropertyLeftElectric{DevEUI: terminalInfo.DevEUI, LeftElectric: power}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSceneSwitchEM30:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanSceneSwitchEM30PropertyLeftElectric{DevEUI: terminalInfo.DevEUI, LeftElectric: power}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30:
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM:
			default:
				globallogger.Log.Warnf("[devEUI: %v][genPowerCfgProcAttribute] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
			}
		} else if constant.Constant.Iotprivate {
			type appDataMsg struct {
				Electric string `json:"electric"`
			}
			kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
				Time:       time.Now(),
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData: appDataMsg{
					Electric: strconv.FormatUint(attribute.Value.(uint64)/2, 10) + "%",
				},
			})
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	default:
		globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "genPowerCfgProcAttribute unknow attributeName", attributeName)
	}
}

//genPowerCfgProcReadRsp 处理readRsp（0x01）消息
func genPowerCfgProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	for _, v := range command.(*cluster.ReadAttributesResponse).ReadAttributeStatuses {
		globallogger.Log.Infof("[devEUI: %v][genPowerCfgProcReadRsp]: readAttributesRsp: %+v", terminalInfo.DevEUI, v)
		if v.Status == cluster.ZclStatusSuccess {
			genPowerCfgProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
		}
	}
}

func procGenPowerCfgProcRead(devEUI string, dstEndpointIndex int, clusterID uint16) {
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.ReadAttribute,
		ClusterID:   clusterID,
		Command: common.Command{
			DstEndpointIndex: dstEndpointIndex,
		},
	})
}

func genPowerCfgProcKeepAlive(terminalInfo config.TerminalInfo, interval uint16) {
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.MAILEKE.ZigbeeTerminalPMT1004Detector,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalPMTSensor0001112b:
		keepalive.ProcKeepAlive(terminalInfo, interval)
		go procGenPowerCfgProcRead(terminalInfo.DevEUI, 0, 0x0402)
		go func() {
			timer := time.NewTimer(2 * time.Second)
			<-timer.C
			timer.Stop()
			procGenPowerCfgProcRead(terminalInfo.DevEUI, 0, 0x0405)
		}()
		go func() {
			timer := time.NewTimer(4 * time.Second)
			<-timer.C
			timer.Stop()
			procGenPowerCfgProcRead(terminalInfo.DevEUI, 0, 0x0415)
		}()
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
		keepalive.ProcKeepAlive(terminalInfo, interval)
		go procGenPowerCfgProcRead(terminalInfo.DevEUI, 0, 0x0402)
		go func() {
			timer := time.NewTimer(2 * time.Second)
			<-timer.C
			timer.Stop()
			procGenPowerCfgProcRead(terminalInfo.DevEUI, 0, 0x0405)
		}()
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalPIRSensorHY0027,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSceneSwitchEM30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSmokeSensorHY0024,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalWarningDevice005b0e12,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM:
		keepalive.ProcKeepAlive(terminalInfo, interval)
	default:
		globallogger.Log.Warnf("[genPowerCfgProcReport]: invalid tmnType: %s", terminalInfo.TmnType)
	}
}

// genPowerCfgProcReport 处理report（0x0a）消息
func genPowerCfgProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][genPowerCfgProcReport]: command: %+v", terminalInfo.DevEUI, Command)
	for _, v := range Command.AttributeReports {
		genPowerCfgProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
	}
	genPowerCfgProcKeepAlive(terminalInfo, uint16(terminalInfo.Interval))
}

// genPowerCfgProcConfigureReportingResponse 处理configureReportingResponse（0x07）消息
func genPowerCfgProcConfigureReportingResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ConfigureReportingResponse)
	globallogger.Log.Infof("[devEUI: %v][genPowerCfgProcConfigureReportingResponse]: command: %+v", terminalInfo.DevEUI, Command)
	for _, v := range Command.AttributeStatusRecords {
		if v.Status != cluster.ZclStatusSuccess {
			globallogger.Log.Infof("[devEUI: %v][genPowerCfgProcConfigureReportingResponse]: configReport failed: %x", terminalInfo.DevEUI, v.Status)
		} else {
			// proc keepalive
			if v.AttributeID == 0x0021 {
				genPowerCfgProcKeepAlive(terminalInfo, uint16(terminalInfo.Interval))
			}
		}
		// 空气盒子处理
		if terminalInfo.TmnType == constant.Constant.TMNTYPE.MAILEKE.ZigbeeTerminalPMT1004Detector {
			type confirmMessage struct {
				Response      string `json:"response"`
				OIDIndex      string `json:"OIDIndex"`
				PropertyTopic string `json:"propertyTopic"`
			}
			confirmMsg := confirmMessage{
				Response:      "success",
				OIDIndex:      terminalInfo.OIDIndex,
				PropertyTopic: "00000003",
			}
			if v.Status != cluster.ZclStatusSuccess {
				confirmMsg.Response = "failed"
			}
			confirmMsgByte, _ := json.Marshal(confirmMsg)
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicConfirmMessage, string(confirmMsgByte))
		}
	}
	if terminalInfo.TmnType == constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM {
		iotsmartspace.ProcHEIMANHTEM(terminalInfo, 2)
	}
}

// GenPowerCfgProc 处理clusterID 0x0001即genPowerCfg属性消息
func GenPowerCfgProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	switch zclFrame.CommandName {
	case "ReadAttributesResponse":
		genPowerCfgProcReadRsp(terminalInfo, zclFrame.Command)
	case "ConfigureReportingResponse":
		genPowerCfgProcConfigureReportingResponse(terminalInfo, zclFrame.Command)
	case "ReportAttributes":
		genPowerCfgProcReport(terminalInfo, zclFrame.Command)
	default:
		globallogger.Log.Warnf("[devEUI: %v][GenPowerCfgProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
