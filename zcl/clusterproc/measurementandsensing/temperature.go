package measurementandsensing

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

func temperatureMeasurementProcAttribute(terminalInfo config.TerminalInfo, attributeName string, attribute *cluster.Attribute) {
	switch attributeName {
	case "MeasuredValue":
		temperature := float64(attribute.Value.(int64)) * 0.01
		if attribute.Value.(int64) > 38220 {
			temperature = temperature - 655.36
		}
		// if terminalInfo.UnitOfTemperature == "F" {
		// 	temperature = 32 + temperature*1.8
		// }
		temperature, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", temperature), 64)

		// iotsmartspace publish msg to app
		if constant.Constant.Iotware {
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyCurrentTemperature{CurrentTemperature: temperature})
		} else if constant.Constant.Iotedge {
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHTSensorPropertyTemperature{DevEUI: terminalInfo.DevEUI, Temperature: temperature}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyTemperature{DevEUI: terminalInfo.DevEUI, Temperature: temperature}, uuid.NewV4().String())
			default:
				globallogger.Log.Warnf("[devEUI: %v][temperatureMeasurementProcAttribute] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
			}
		} else if constant.Constant.Iotprivate {
			type appDataMsg struct {
				Temperature string `json:"temperature"`
			}
			unitOfTemperature := "℃"
			if terminalInfo.UnitOfTemperature == "F" {
				unitOfTemperature = "℉"
			}
			kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
				Time:       time.Now(),
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData: appDataMsg{
					Temperature: strconv.FormatFloat(temperature, 'f', 2, 64) + unitOfTemperature,
				},
			})
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	case "MinMeasuredValue":
	case "MaxMeasuredValue":
	case "Tolerance":
	default:
		globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "temperatureMeasurementProcAttribute unknow attributeName", attributeName)
	}
}

//temperatureMeasurementProcReadRsp 处理readRsp（0x01）消息
func temperatureMeasurementProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	for _, v := range command.(*cluster.ReadAttributesResponse).ReadAttributeStatuses {
		globallogger.Log.Infof("[devEUI: %v][temperatureMeasurementProcReadRsp]: readAttributesRsp: %+v", terminalInfo.DevEUI, v)
		if v.Status == cluster.ZclStatusSuccess {
			temperatureMeasurementProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
		}
	}
}

// temperatureMeasurementProcReport 处理report（0x0a）消息
func temperatureMeasurementProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][temperatureMeasurementProcReport]: command: %+v", terminalInfo.DevEUI, Command)
	for _, v := range Command.AttributeReports {
		temperatureMeasurementProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
	}
}

// temperatureMeasurementProcConfigureReportingResponse 处理configureReportingResponse（0x07）消息
func temperatureMeasurementProcConfigureReportingResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ConfigureReportingResponse)
	globallogger.Log.Infof("[devEUI: %v][temperatureMeasurementProcConfigureReportingResponse]: command: %+v", terminalInfo.DevEUI, Command)
	for _, v := range Command.AttributeStatusRecords {
		globallogger.Log.Infof("[devEUI: %v][temperatureMeasurementProcConfigureReportingResponse]: AttributeStatusRecord: %+v", terminalInfo.DevEUI, v)
		if v.Status != cluster.ZclStatusSuccess {
			globallogger.Log.Infof("[devEUI: %v][temperatureMeasurementProcConfigureReportingResponse]: configReport failed: %x", terminalInfo.DevEUI, v.Status)
		}
	}
	if terminalInfo.TmnType == constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM {
		iotsmartspace.ProcHEIMANHTEM(terminalInfo, 3)
	}
}

// TemperatureMeasurementProc 处理clusterID 0x0402属性消息
func TemperatureMeasurementProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	switch zclFrame.CommandName {
	case "ReadAttributesResponse":
		temperatureMeasurementProcReadRsp(terminalInfo, zclFrame.Command)
	case "ConfigureReportingResponse":
		temperatureMeasurementProcConfigureReportingResponse(terminalInfo, zclFrame.Command)
	case "ReportAttributes":
		temperatureMeasurementProcReport(terminalInfo, zclFrame.Command)
	default:
		globallogger.Log.Warnf("[devEUI: %v][TemperatureMeasurementProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
