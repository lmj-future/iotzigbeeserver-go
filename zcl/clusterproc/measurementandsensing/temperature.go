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

func getTemperatureMeasurementParams(terminalID string, tmnType string, temperature float64) (interface{}, bool) {
	bPublish := false
	params := make(map[string]interface{}, 2)
	var key string
	params["terminalId"] = terminalID
	switch tmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
		key = iotsmartspace.HeimanHTSensorPropertyTemperature
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
		key = iotsmartspace.HeimanHS2AQPropertyTemperature
		bPublish = true
	default:
		globallogger.Log.Warnf("[devEUI: %v][getTemperatureMeasurementParams] invalid tmnType: %v", terminalID, tmnType)
	}
	params[key] = temperature
	return params, bPublish
}

func temperatureMeasurementProcAttribute(terminalInfo config.TerminalInfo, attributeName string, attribute *cluster.Attribute) {
	switch attributeName {
	case "MeasuredValue":
		temperature := float64(attribute.Value.(int64)) * 0.01
		if attribute.Value.(int64) > 38220 {
			temperature = temperature - 655.36
		}
		if terminalInfo.UnitOfTemperature == "F" {
			temperature = 32 + temperature*1.8
		}
		temperature, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", temperature), 64)

		// iotsmartspace publish msg to app
		if constant.Constant.Iotware {
			values := make(map[string]interface{}, 1)
			values[iotsmartspace.IotwarePropertyTemperature] = temperature
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
		} else if constant.Constant.Iotedge {
			params, bPublish := getTemperatureMeasurementParams(terminalInfo.DevEUI, terminalInfo.TmnType, temperature)
			if bPublish {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
			}
		} else if constant.Constant.Iotprivate {
			type appDataMsg struct {
				Temperature string `json:"temperature"`
			}
			unitOfTemperature := "℃"
			if terminalInfo.UnitOfTemperature == "F" {
				unitOfTemperature = "℉"
			}
			kafkaMsg := publicstruct.DataReportMsg{
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData: appDataMsg{
					Temperature: strconv.FormatFloat(temperature, 'f', 2, 64) + unitOfTemperature,
				},
			}
			kafkaMsgByte, _ := json.Marshal(kafkaMsg)
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	case "MinMeasuredValue":
	case "MaxMeasuredValue":
	case "Tolerance":
	default:
		globallogger.Log.Warnln("devEUI : "+terminalInfo.DevEUI+" "+"temperatureMeasurementProcAttribute unknow attributeName", attributeName)
	}
}

//temperatureMeasurementProcReadRsp 处理readRsp（0x01）消息
func temperatureMeasurementProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	readAttributesRsp := command.(*cluster.ReadAttributesResponse)
	for _, v := range readAttributesRsp.ReadAttributeStatuses {
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
	attributeReports := Command.AttributeReports
	for _, v := range attributeReports {
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
}

// TemperatureMeasurementProc 处理clusterID 0x0402属性消息
func TemperatureMeasurementProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	// globallogger.Log.Infof("[devEUI: %v][TemperatureMeasurementProc] Start......", terminalInfo.DevEUI)
	// globallogger.Log.Infof("[devEUI: %v][TemperatureMeasurementProc] zclFrame: %+v", terminalInfo.DevEUI, zclFrame)
	z := zcl.New()
	switch zclFrame.CommandName {
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandReadAttributesResponse)].Name:
		temperatureMeasurementProcReadRsp(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandReportAttributes)].Name:
		temperatureMeasurementProcReport(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandConfigureReporting)].Name:
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandConfigureReportingResponse)].Name:
		temperatureMeasurementProcConfigureReportingResponse(terminalInfo, zclFrame.Command)
	default:
		globallogger.Log.Warnf("[devEUI: %v][TemperatureMeasurementProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
