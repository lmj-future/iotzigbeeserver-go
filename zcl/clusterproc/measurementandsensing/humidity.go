package measurementandsensing

import (
	"encoding/json"
	"fmt"
	"strconv"

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

func getRelativeHumidityParams(terminalID string, tmnType string, humidity float64) (interface{}, bool) {
	bPublish := false
	params := make(map[string]interface{}, 2)
	var key string
	params["terminalId"] = terminalID
	switch tmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
		key = iotsmartspace.HeimanHTSensorPropertyHumidity
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
		key = iotsmartspace.HeimanHS2AQPropertyHumidity
		bPublish = true
	default:
		globallogger.Log.Warnf("[devEUI: %v][getRelativeHumidityParams] invalid tmnType: %v", terminalID, tmnType)
	}
	params[key] = humidity
	return params, bPublish
}

func relativeHumidityMeasurementProcAttribute(terminalInfo config.TerminalInfo, attributeName string, attribute *cluster.Attribute) {
	switch attributeName {
	case "MeasuredValue":
		humidity := float64(attribute.Value.(uint64)) * 0.01
		humidity, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", humidity), 64)

		// iotsmartspace publish msg to app
		if constant.Constant.Iotware {
			values := make(map[string]interface{}, 1)
			values[iotsmartspace.IotwarePropertyHumidity] = humidity
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
		} else if constant.Constant.Iotedge {
			params, bPublish := getRelativeHumidityParams(terminalInfo.DevEUI, terminalInfo.TmnType, humidity)
			if bPublish {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
			}
		} else if constant.Constant.Iotprivate {
			type appDataMsg struct {
				Humidity string `json:"humidity"`
			}
			kafkaMsg := publicstruct.DataReportMsg{
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData: appDataMsg{
					Humidity: strconv.FormatFloat(humidity, 'f', 2, 64) + "%RH",
				},
			}
			kafkaMsgByte, _ := json.Marshal(kafkaMsg)
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	case "MinMeasuredValue":
	case "MaxMeasuredValue":
	case "Tolerance":
	default:
		globallogger.Log.Warnln("devEUI : "+terminalInfo.DevEUI+" "+"relativeHumidityMeasurementProcAttribute unknow attributeName", attributeName)
	}
}

//relativeHumidityMeasurementProcReadRsp 处理readRsp（0x01）消息
func relativeHumidityMeasurementProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	readAttributesRsp := command.(*cluster.ReadAttributesResponse)
	for _, v := range readAttributesRsp.ReadAttributeStatuses {
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
	// globallogger.Log.Infof("[devEUI: %v][RelativeHumidityMeasurementProc] Start......", terminalInfo.DevEUI)
	// globallogger.Log.Infof("[devEUI: %v][RelativeHumidityMeasurementProc] zclFrame: %+v", terminalInfo.DevEUI, zclFrame)
	z := zcl.New()
	switch zclFrame.CommandName {
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandReadAttributesResponse)].Name:
		relativeHumidityMeasurementProcReadRsp(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandReportAttributes)].Name:
		relativeHumidityMeasurementProcReport(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandConfigureReporting)].Name:
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandConfigureReportingResponse)].Name:
	default:
		globallogger.Log.Warnf("[devEUI: %v][RelativeHumidityMeasurementProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
