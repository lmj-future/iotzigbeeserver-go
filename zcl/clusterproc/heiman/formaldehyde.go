package heiman

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

func getFormaldehydeMeasurementParams(terminalID string, tmnType string, formaldehyde float64) interface{} {
	params := make(map[string]interface{}, 2)
	params["terminalId"] = terminalID
	var key string
	switch tmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
		key = iotsmartspace.HeimanHS2AQPropertyFormaldehyde
	default:
		globallogger.Log.Warnf("[devEUI: %v][getTemperatureMeasurementParams] invalid tmnType: %v", terminalID, tmnType)
	}
	params[key] = formaldehyde
	return params
}

func formaldehydeMeasurementProcAttribute(terminalInfo config.TerminalInfo, attributeName string, attribute *cluster.Attribute) {
	switch attributeName {
	case "MeasuredValue":
		formaldehyde := float64(attribute.Value.(uint64)) * 0.001
		formaldehyde, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", formaldehyde), 64)

		// iotsmartspace publish msg to app
		if constant.Constant.Iotware {
			values := make(map[string]interface{}, 1)
			values[iotsmartspace.IotwarePropertyHCHO] = formaldehyde
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
		} else if constant.Constant.Iotedge {
			params := getFormaldehydeMeasurementParams(terminalInfo.DevEUI, terminalInfo.TmnType, formaldehyde)
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
		} else if constant.Constant.Iotprivate {
			type appDataMsg struct {
				HCHO string `json:"HCHO"`
			}
			kafkaMsg := publicstruct.DataReportMsg{
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData: appDataMsg{
					HCHO: strconv.FormatFloat(formaldehyde, 'f', 2, 64) + "mg/m^3",
				},
			}
			kafkaMsgByte, _ := json.Marshal(kafkaMsg)
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	case "MinMeasuredValue":
	case "MaxMeasuredValue":
	case "Tolerance":
	default:
		globallogger.Log.Warnln("devEUI : "+terminalInfo.DevEUI+" "+"formaldehydeMeasurementProcAttribute unknow attributeName", attributeName)
	}
}

//formaldehydeMeasurementProcReadRsp 处理readRsp（0x01）消息
func formaldehydeMeasurementProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	readAttributesRsp := command.(*cluster.ReadAttributesResponse)
	for _, v := range readAttributesRsp.ReadAttributeStatuses {
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
	attributeReports := Command.AttributeReports
	for _, v := range attributeReports {
		formaldehydeMeasurementProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
	}
}

// FormaldehydeMeasurementProc 处理clusterID 0x042b属性消息
func FormaldehydeMeasurementProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	// globallogger.Log.Infof("[devEUI: %v][FormaldehydeMeasurementProc] Start......", terminalInfo.DevEUI)
	// globallogger.Log.Infof("[devEUI: %v][FormaldehydeMeasurementProc] zclFrame: %+v", terminalInfo.DevEUI, zclFrame)
	z := zcl.New()
	switch zclFrame.CommandName {
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandReadAttributesResponse)].Name:
		formaldehydeMeasurementProcReadRsp(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandReportAttributes)].Name:
		formaldehydeMeasurementProcReport(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandConfigureReportingResponse)].Name:
	default:
		globallogger.Log.Warnf("[devEUI: %v][FormaldehydeMeasurementProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
