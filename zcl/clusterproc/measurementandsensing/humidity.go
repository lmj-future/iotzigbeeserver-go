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

func relativeHumidityMeasurementProcAttribute(terminalInfo config.TerminalInfo, attributeName string, attribute *cluster.Attribute) {
	switch attributeName {
	case "MeasuredValue":
		humidity := float64(attribute.Value.(uint64)) * 0.01
		humidity, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", humidity), 64)

		// iotsmartspace publish msg to app
		if constant.Constant.Iotware {
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyCurrentHumidity{CurrentHumidity: humidity})
		} else if constant.Constant.Iotedge {
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHTSensorPropertyHumidity{DevEUI: terminalInfo.DevEUI, Humidity: humidity}, uuid.NewV4().String())
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyHumidity{DevEUI: terminalInfo.DevEUI, Humidity: humidity}, uuid.NewV4().String())
			default:
				globallogger.Log.Warnf("[devEUI: %v][relativeHumidityMeasurementProcAttribute] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
			}
		} else if constant.Constant.Iotprivate {
			type appDataMsg struct {
				Humidity string `json:"humidity"`
			}
			kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
				Time:       time.Now(),
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData: appDataMsg{
					Humidity: strconv.FormatFloat(humidity, 'f', 2, 64) + "%RH",
				},
			})
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	case "MinMeasuredValue":
	case "MaxMeasuredValue":
	case "Tolerance":
	default:
		globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "relativeHumidityMeasurementProcAttribute unknow attributeName", attributeName)
	}
}

//relativeHumidityMeasurementProcReadRsp 处理readRsp（0x01）消息
func relativeHumidityMeasurementProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	for _, v := range command.(*cluster.ReadAttributesResponse).ReadAttributeStatuses {
		globallogger.Log.Infof("[devEUI: %v][relativeHumidityMeasurementProcReadRsp]: readAttributesRsp: %+v", terminalInfo.DevEUI, v)
		if v.Status == cluster.ZclStatusSuccess {
			relativeHumidityMeasurementProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
		}
	}
}

// relativeHumidityMeasurementProcReport 处理report（0x0a）消息
func relativeHumidityMeasurementProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][relativeHumidityMeasurementProcReport]: command: %+v", terminalInfo.DevEUI, Command)
	attributeReports := Command.AttributeReports
	for _, v := range attributeReports {
		relativeHumidityMeasurementProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
	}
}

// RelativeHumidityMeasurementProc 处理clusterID 0x0405属性消息
func RelativeHumidityMeasurementProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	switch zclFrame.CommandName {
	case "ReadAttributesResponse":
		relativeHumidityMeasurementProcReadRsp(terminalInfo, zclFrame.Command)
	case "ConfigureReportingResponse":
		if terminalInfo.TmnType == constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM {
			iotsmartspace.ProcHEIMANHTEM(terminalInfo, 4)
		}
	case "ReportAttributes":
		relativeHumidityMeasurementProcReport(terminalInfo, zclFrame.Command)
	default:
		globallogger.Log.Warnf("[devEUI: %v][RelativeHumidityMeasurementProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
