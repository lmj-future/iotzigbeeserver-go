package heiman

import (
	"encoding/json"
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

func getPM25MeasurementParams(terminalID string, tmnType string, PM25 uint64) interface{} {
	params := make(map[string]interface{}, 2)
	params["terminalId"] = terminalID
	var key string
	switch tmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
		key = iotsmartspace.HeimanHS2AQPropertyPM25
	default:
		globallogger.Log.Warnf("[devEUI: %v][getTemperatureMeasurementParams] invalid tmnType: %v", terminalID, tmnType)
	}
	params[key] = PM25
	return params
}

func pm25MeasurementProcAttribute(terminalInfo config.TerminalInfo, attributeName string, attribute *cluster.Attribute) {
	switch attributeName {
	case "MeasuredValue":
		PM25 := attribute.Value.(uint64)

		// iotsmartspace publish msg to app
		if constant.Constant.Iotware {
			values := make(map[string]interface{}, 1)
			values[iotsmartspace.IotwarePropertyPM25] = PM25
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
		} else if constant.Constant.Iotedge {
			params := getPM25MeasurementParams(terminalInfo.DevEUI, terminalInfo.TmnType, PM25)
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
		} else if constant.Constant.Iotprivate {
			type appDataMsg struct {
				PM25 string `json:"PM25"`
			}
			kafkaMsg := publicstruct.DataReportMsg{
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData: appDataMsg{
					PM25: strconv.FormatUint(PM25, 10) + "ug/m^3",
				},
			}
			kafkaMsgByte, _ := json.Marshal(kafkaMsg)
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	case "MinMeasuredValue":
	case "MaxMeasuredValue":
	case "Tolerance":
	default:
		globallogger.Log.Warnln("devEUI : "+terminalInfo.DevEUI+" "+"pm25MeasurementProcAttribute unknow attributeName", attributeName)
	}
}

//pm25MeasurementProcReadRsp 处理readRsp（0x01）消息
func pm25MeasurementProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	readAttributesRsp := command.(*cluster.ReadAttributesResponse)
	for _, v := range readAttributesRsp.ReadAttributeStatuses {
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
	attributeReports := Command.AttributeReports
	for _, v := range attributeReports {
		pm25MeasurementProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
	}
}

// PM25MeasurementProc 处理clusterID 0x042a属性消息
func PM25MeasurementProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	// globallogger.Log.Infof("[devEUI: %v][PM25MeasurementProc] Start......", terminalInfo.DevEUI)
	// globallogger.Log.Infof("[devEUI: %v][PM25MeasurementProc] zclFrame: %+v", terminalInfo.DevEUI, zclFrame)
	z := zcl.New()
	switch zclFrame.CommandName {
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandReadAttributesResponse)].Name:
		pm25MeasurementProcReadRsp(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandReportAttributes)].Name:
		pm25MeasurementProcReport(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandConfigureReportingResponse)].Name:
	default:
		globallogger.Log.Warnf("[devEUI: %v][PM25MeasurementProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
