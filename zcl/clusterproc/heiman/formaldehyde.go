package heiman

import (
	"encoding/json"
	"fmt"
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

func formaldehydeMeasurementProcAttribute(terminalInfo config.TerminalInfo, attributeName string, attribute *cluster.Attribute) {
	switch attributeName {
	case "MeasuredValue":
		formaldehyde := float64(attribute.Value.(uint64)) * 0.001
		formaldehyde, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", formaldehyde), 64)
		// iotsmartspace publish msg to app
		if constant.Constant.Iotware {
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyHCHO{HCHO: formaldehyde})
		} else if constant.Constant.Iotedge {
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyFormaldehyde{DevEUI: terminalInfo.DevEUI, Formaldehyde: formaldehyde}, uuid.NewV4().String())
			default:
				globallogger.Log.Warnf("[devEUI: %v][formaldehydeMeasurementProcAttribute] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
			}
		} else if constant.Constant.Iotprivate {
			type appDataMsg struct {
				HCHO string `json:"HCHO"`
			}
			kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
				Time:       time.Now(),
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData: appDataMsg{
					HCHO: strconv.FormatFloat(formaldehyde, 'f', 2, 64) + "mg/m^3",
				},
			})
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	case "MinMeasuredValue":
	case "MaxMeasuredValue":
	case "Tolerance":
	default:
		globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "formaldehydeMeasurementProcAttribute unknow attributeName", attributeName)
	}
}

//formaldehydeMeasurementProcReadRsp 处理readRsp（0x01）消息
func formaldehydeMeasurementProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	for _, v := range command.(*cluster.ReadAttributesResponse).ReadAttributeStatuses {
		globallogger.Log.Infof("[devEUI: %v][formaldehydeMeasurementProcReadRsp]: readAttributesRsp: %+v", terminalInfo.DevEUI, v)
		if v.Status == cluster.ZclStatusSuccess {
			formaldehydeMeasurementProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
		}
	}
}

// formaldehydeMeasurementProcReport 处理report（0x0a）消息
func formaldehydeMeasurementProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][formaldehydeMeasurementProcReport]: command: %+v", terminalInfo.DevEUI, Command)
	for _, v := range Command.AttributeReports {
		formaldehydeMeasurementProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
	}
}

// FormaldehydeMeasurementProc 处理clusterID 0x042b属性消息
func FormaldehydeMeasurementProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	switch zclFrame.CommandName {
	case "ReadAttributesResponse":
		formaldehydeMeasurementProcReadRsp(terminalInfo, zclFrame.Command)
	case "ReportAttributes":
		formaldehydeMeasurementProcReport(terminalInfo, zclFrame.Command)
	default:
		globallogger.Log.Warnf("[devEUI: %v][FormaldehydeMeasurementProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
