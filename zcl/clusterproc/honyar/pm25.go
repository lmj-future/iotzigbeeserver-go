package honyar

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

func pm25MeasurementProcAttribute(terminalInfo config.TerminalInfo, attributeName string, attribute *cluster.Attribute) {
	switch attributeName {
	case "MeasuredValue":
		PM25 := attribute.Value.(uint64)
		// iotsmartspace publish msg to app
		if constant.Constant.Iotware {
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyPM25{PM25: PM25})
		} else if constant.Constant.Iotedge {
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyPM25{DevEUI: terminalInfo.DevEUI, PM25: PM25}, uuid.NewV4().String())
			default:
				globallogger.Log.Warnf("[devEUI: %v][pm25MeasurementProcAttribute] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
			}
		} else if constant.Constant.Iotprivate {
			type appDataMsg struct {
				PM25 string `json:"PM25"`
			}
			kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
				Time:       time.Now(),
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData: appDataMsg{
					PM25: strconv.FormatUint(PM25, 10) + "ug/m^3",
				},
			})
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	default:
		globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "pm25MeasurementProcAttribute unknow attributeName", attributeName)
	}
}

//pm25MeasurementProcReadRsp 处理readRsp（0x01）消息
func pm25MeasurementProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	for _, v := range command.(*cluster.ReadAttributesResponse).ReadAttributeStatuses {
		globallogger.Log.Infof("[devEUI: %v][pm25MeasurementProcReadRsp]: readAttributesRsp: %+v", terminalInfo.DevEUI, v)
		if v.Status == cluster.ZclStatusSuccess {
			pm25MeasurementProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
		}
	}
}

// pm25MeasurementProcReport 处理report（0x0a）消息
func pm25MeasurementProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][pm25MeasurementProcReport]: command: %+v", terminalInfo.DevEUI, Command)
	for _, v := range Command.AttributeReports {
		pm25MeasurementProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
	}
}

// PM25MeasurementProc 处理clusterID 0x042a属性消息
func PM25MeasurementProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	switch zclFrame.CommandName {
	case "ReadAttributesResponse":
		pm25MeasurementProcReadRsp(terminalInfo, zclFrame.Command)
	case "ReportAttributes":
		pm25MeasurementProcReport(terminalInfo, zclFrame.Command)
	default:
		globallogger.Log.Warnf("[devEUI: %v][PM25MeasurementProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
