package heiman

import (
	"encoding/json"
	"strconv"

	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/db/kafka"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
	uuid "github.com/satori/go.uuid"
)

func airQualityMeasurementProcAttribute(devEUI string, attributeName string, attribute *cluster.Attribute, bPublish bool, oSet bson.M,
	oSetPG map[string]interface{}, params map[string]interface{}, values map[string]interface{}) (bool, bson.M, map[string]interface{},
	map[string]interface{}, map[string]interface{}) {
	switch attributeName {
	case "Language":
		Value := attribute.Value.(uint64)
		if Value == 0 {
			oSet["language"] = "Chinese"
			oSetPG["language"] = "Chinese"
		} else if Value == 1 {
			oSet["language"] = "English"
			oSetPG["language"] = "English"
		} else {
			oSet["language"] = "Invalid"
			oSetPG["language"] = "Invalid"
		}
	case "UnitOfTemperature":
		Value := attribute.Value.(uint64)
		if Value == 0 {
			oSet["unitOfTemperature"] = "C"
			oSetPG["unitoftemperature"] = "C"
		} else if Value == 1 {
			oSet["unitOfTemperature"] = "F"
			oSetPG["unitoftemperature"] = "F"
		} else {
			oSet["unitOfTemperature"] = "Invalid"
			oSetPG["unitoftemperature"] = "Invalid"
		}
	case "BatteryState":
		Value := attribute.Value.(uint64)
		if Value == 0 {
			oSet["batteryState"] = "unCharge"
			oSetPG["batterystate"] = "unCharge"
			params[iotsmartspace.HeimanHS2AQPropertyBatteryState] = "0"
			values[iotsmartspace.IotwarePropertyAirState] = "0"
			bPublish = true
		} else if Value == 1 {
			oSet["batteryState"] = "charging"
			oSetPG["batterystate"] = "charging"
			params[iotsmartspace.HeimanHS2AQPropertyBatteryState] = "1"
			values[iotsmartspace.IotwarePropertyAirState] = "1"
			bPublish = true
		} else if Value == 2 {
			oSet["batteryState"] = "fullCharged"
			oSetPG["batterystate"] = "fullCharged"
			params[iotsmartspace.HeimanHS2AQPropertyBatteryState] = "2"
			values[iotsmartspace.IotwarePropertyAirState] = "2"
			bPublish = true
		} else {
			oSet["batteryState"] = "Invalid"
			oSetPG["batterystate"] = "Invalid"
		}
	case "PM10MeasuredValue":
		Value := attribute.Value.(uint64)
		oSet["PM10"] = strconv.FormatUint(Value, 10)
		oSetPG["pm10"] = strconv.FormatUint(Value, 10)
		params[iotsmartspace.HeimanHS2AQPropertyPM10] = strconv.FormatUint(Value, 10)
		values[iotsmartspace.IotwarePropertyPM10] = strconv.FormatUint(Value, 10)
		bPublish = true
	case "TVOCMeasuredValue":
		Value := attribute.Value.(uint64)
		oSet["TVOC"] = strconv.FormatUint(Value, 10)
		oSetPG["tvoc"] = strconv.FormatUint(Value, 10)
		params[iotsmartspace.HeimanHS2AQPropertyTVOC] = strconv.FormatUint(Value, 10)
		values[iotsmartspace.IotwarePropertyTVOC] = strconv.FormatUint(Value, 10)
		bPublish = true
	case "AQIMeasuredValue":
		Value := attribute.Value.(uint64)
		oSet["AQI"] = strconv.FormatUint(Value, 10)
		oSetPG["aqi"] = strconv.FormatUint(Value, 10)
		params[iotsmartspace.HeimanHS2AQPropertyAQI] = strconv.FormatUint(Value, 10)
		values[iotsmartspace.IotwarePropertyAQI] = strconv.FormatUint(Value, 10)
		bPublish = true
	case "MaxTemperature":
		Value := attribute.Value.(int64)
		oSet["maxTemperature"] = strconv.FormatInt(Value, 10)
		oSetPG["maxtemperature"] = strconv.FormatInt(Value, 10)
	case "MinTemperature":
		Value := attribute.Value.(int64)
		oSet["minTemperature"] = strconv.FormatInt(Value, 10)
		oSetPG["mintemperature"] = strconv.FormatInt(Value, 10)
	case "MaxHumidity":
		Value := attribute.Value.(uint64)
		oSet["maxHumidity"] = strconv.FormatUint(Value, 10)
		oSetPG["maxhumidity"] = strconv.FormatUint(Value, 10)
	case "MinHumidity":
		Value := attribute.Value.(uint64)
		oSet["minHumidity"] = strconv.FormatUint(Value, 10)
		oSetPG["minhumidity"] = strconv.FormatUint(Value, 10)
	case "Disturb":
		Value := attribute.Value.(uint64)
		disturb := strconv.FormatUint(Value, 16)
		disturb = "0000" + disturb
		disturb = disturb[len(disturb)-4:]
		oSet["disturb"] = disturb
		oSetPG["disturb"] = disturb
	default:
		globallogger.Log.Warnf("[devEUI: %v][airQualityMeasurementProcAttribute] invalid attributeName: %v", devEUI, attributeName)
	}
	return bPublish, oSet, oSetPG, params, values
}

func airQualityMeasurementProcMsg2Kafka(terminalInfo config.TerminalInfo, values map[string]interface{}) {
	kafkaMsg := publicstruct.DataReportMsg{
		OIDIndex:   terminalInfo.OIDIndex,
		DevSN:      terminalInfo.DevEUI,
		LinkType:   terminalInfo.ProfileID,
		DeviceType: terminalInfo.TmnType2,
	}
	if PM10, ok := values[iotsmartspace.IotwarePropertyPM10]; ok {
		type appDataMsg struct {
			PM10 string `json:"PM10"`
		}
		kafkaMsg.AppData = appDataMsg{PM10: PM10.(string) + "ug/m^3"}
		kafkaMsgByte, _ := json.Marshal(kafkaMsg)
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
	if AQI, ok := values[iotsmartspace.IotwarePropertyAQI]; ok {
		type appDataMsg struct {
			AQI string `json:"AQI"`
		}
		kafkaMsg.AppData = appDataMsg{AQI: AQI.(string)}
		kafkaMsgByte, _ := json.Marshal(kafkaMsg)
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
	if batteryState, ok := values[iotsmartspace.IotwarePropertyAirState]; ok {
		type appDataMsg struct {
			BatteryState string `json:"batteryState"`
		}
		if batteryState == "0" {
			batteryState = "未充电"
		} else if batteryState == "1" {
			batteryState = "充电中"
		} else if batteryState == "2" {
			batteryState = "已充满"
		}
		kafkaMsg.AppData = appDataMsg{BatteryState: batteryState.(string)}
		kafkaMsgByte, _ := json.Marshal(kafkaMsg)
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
}

//airQualityMeasurementProcReadRsp 处理readRsp（0x01）消息
func airQualityMeasurementProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	readAttributesRsp := command.(*cluster.ReadAttributesResponse)
	oSet := bson.M{}
	oSetPG := make(map[string]interface{})
	params := make(map[string]interface{}, 11)
	values := make(map[string]interface{}, 11)
	params["terminalId"] = terminalInfo.DevEUI
	var bPublish = false
	for _, v := range readAttributesRsp.ReadAttributeStatuses {
		globallogger.Log.Infof("[devEUI: %v][airQualityMeasurementProcReadRsp]: readAttributesRsp: %+v", terminalInfo.DevEUI, v)
		if v.Status != cluster.ZclStatusSuccess {
			continue
		}
		bPublish, oSet, oSetPG, params, values = airQualityMeasurementProcAttribute(terminalInfo.DevEUI, v.AttributeName, v.Attribute,
			bPublish, oSet, oSetPG, params, values)
	}
	// globallogger.Log.Infof("[devEUI: %v][airQualityMeasurementProcReadRsp]: update terminal info: oSet %+v", terminalInfo.DevEUI, oSet)
	if oSet["disturb"] != nil || oSet["minHumidity"] != nil || oSet["maxHumidity"] != nil || oSet["minTemperature"] != nil || oSet["maxTemperature"] != nil {
		if constant.Constant.UsePostgres {
			models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": terminalInfo.DevEUI}, oSetPG)
		} else {
			models.FindTerminalAndUpdate(bson.M{"devEUI": terminalInfo.DevEUI}, oSet)
		}
	}

	if bPublish {
		// iotsmartspace publish msg to app
		if constant.Constant.Iotware {
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
		} else if constant.Constant.Iotedge {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
		} else if constant.Constant.Iotprivate {
			airQualityMeasurementProcMsg2Kafka(terminalInfo, values)
		}
	}
}

// airQualityMeasurementProcWriteResponse 处理writeRsp（0x04）消息
func airQualityMeasurementProcWriteResponse(terminalInfo config.TerminalInfo, command interface{}, msgID interface{}, contentFrame *zcl.Frame) {
	Command := command.(*cluster.WriteAttributesResponse)
	globallogger.Log.Infof("[devEUI: %v][airQualityMeasurementProcWriteResponse]: command: %+v", terminalInfo.DevEUI, Command)
	code := 0
	message := "success"
	if Command.WriteAttributeStatuses[0].Status != cluster.ZclStatusSuccess {
		code = -1
		message = "failed"
	} else {
		if contentFrame != nil && contentFrame.CommandName == "WriteAttributes" {
			oSet := bson.M{}
			oSetPG := make(map[string]interface{})
			contentCommand := contentFrame.Command.(*cluster.WriteAttributesCommand)
			for _, v := range contentCommand.WriteAttributeRecords {
				switch v.AttributeName {
				case "MaxTemperature":
					oSet["maxTemperature"] = strconv.FormatInt(v.Attribute.Value.(int64), 10)
					oSetPG["maxtemperature"] = strconv.FormatInt(v.Attribute.Value.(int64), 10)
				case "MinTemperature":
					oSet["minTemperature"] = strconv.FormatInt(v.Attribute.Value.(int64), 10)
					oSetPG["mintemperature"] = strconv.FormatInt(v.Attribute.Value.(int64), 10)
				case "MaxHumidity":
					oSet["maxHumidity"] = strconv.FormatUint(v.Attribute.Value.(uint64), 10)
					oSetPG["maxhumidity"] = strconv.FormatUint(v.Attribute.Value.(uint64), 10)
				case "MinHumidity":
					oSet["minHumidity"] = strconv.FormatUint(v.Attribute.Value.(uint64), 10)
					oSetPG["minhumidity"] = strconv.FormatUint(v.Attribute.Value.(uint64), 10)
				case "Disturb":
					disturb := strconv.FormatUint(v.Attribute.Value.(uint64), 16)
					disturb = "0000" + disturb
					disturb = disturb[len(disturb)-4:]
					oSet["disturb"] = disturb
					oSetPG["disturb"] = disturb
				default:
					globallogger.Log.Infof("[devEUI: %v][airQualityMeasurementProcWriteResponse]: invalid attributeName: %+v",
						terminalInfo.DevEUI, v.AttributeName)
				}
			}
			// globallogger.Log.Infof("[devEUI: %v][airQualityMeasurementProcWriteResponse]: oSet: %+v", terminalInfo.DevEUI, oSet)
			if constant.Constant.UsePostgres {
				models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": terminalInfo.DevEUI}, oSetPG)
			} else {
				models.FindTerminalAndUpdate(bson.M{"devEUI": terminalInfo.DevEUI}, oSet)
			}
		}
	}
	params := make(map[string]interface{}, 2)
	params["code"] = code
	params["message"] = message
	if constant.Constant.Iotware {
		iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, message, msgID)
	} else if constant.Constant.Iotedge {
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply, params, msgID)
	}
}

// airQualityMeasurementProcReport 处理report（0x0a）消息
func airQualityMeasurementProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][airQualityMeasurementProcReport]: command: %+v", terminalInfo.DevEUI, Command)
	attributeReports := Command.AttributeReports
	oSet := bson.M{}
	oSetPG := make(map[string]interface{})
	params := make(map[string]interface{}, 11)
	values := make(map[string]interface{}, 11)
	params["terminalId"] = terminalInfo.DevEUI
	var bPublish = false
	for _, v := range attributeReports {
		bPublish, oSet, oSetPG, params, values = airQualityMeasurementProcAttribute(terminalInfo.DevEUI, v.AttributeName, v.Attribute,
			bPublish, oSet, oSetPG, params, values)
	}

	// globallogger.Log.Infof("[devEUI: %v][airQualityMeasurementProcReport]: update terminal info: oSet %+v", terminalInfo.DevEUI, oSet)
	if oSet["disturb"] != nil || oSet["minHumidity"] != nil || oSet["maxHumidity"] != nil || oSet["minTemperature"] != nil || oSet["maxTemperature"] != nil {
		if constant.Constant.UsePostgres {
			models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": terminalInfo.DevEUI}, oSetPG)
		} else {
			models.FindTerminalAndUpdate(bson.M{"devEUI": terminalInfo.DevEUI}, oSet)
		}
	}

	if bPublish {
		// iotsmartspace publish msg to app
		if constant.Constant.Iotware {
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
		} else if constant.Constant.Iotedge {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
		} else if constant.Constant.Iotprivate {
			airQualityMeasurementProcMsg2Kafka(terminalInfo, values)
		}
	}
}

// airQualityMeasurementProcDefaultResponse 处理defaultResponse（0x0b）消息
func airQualityMeasurementProcDefaultResponse(terminalInfo config.TerminalInfo, command interface{}, msgID interface{}, contentFrame *zcl.Frame) {
	Command := command.(*cluster.DefaultResponseCommand)
	globallogger.Log.Infof("[devEUI: %v][airQualityMeasurementProcDefaultResponse]: command: %+v", terminalInfo.DevEUI, Command)
	result := "success"
	if Command.Status != cluster.ZclStatusSuccess {
		result = "failed"
	}
	var bPublish = false
	if contentFrame != nil {
		switch contentFrame.CommandName {
		case "SetLanguage":
			bPublish = true
		case "SetUnitOfTemperature":
			bPublish = true
		default:
			globallogger.Log.Warnf("[devEUI: %v][airQualityMeasurementProcDefaultResponse] invalid CommandName: %v",
				terminalInfo.DevEUI, contentFrame.CommandName)
		}
	}
	if bPublish {
		if constant.Constant.Iotware {
			iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, result, msgID)
		}
	}
}

// AirQualityMeasurementProc 处理clusterID 0xfc81属性消息
func AirQualityMeasurementProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame, msgID interface{}, contentFrame *zcl.Frame) {
	// globallogger.Log.Infof("[devEUI: %v][AirQualityMeasurementProc] Start......", terminalInfo.DevEUI)
	// globallogger.Log.Infof("[devEUI: %v][AirQualityMeasurementProc] zclFrame: %+v", terminalInfo.DevEUI, zclFrame)
	z := zcl.New()
	switch zclFrame.CommandName {
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandReadAttributesResponse)].Name:
		airQualityMeasurementProcReadRsp(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandWriteAttributesResponse)].Name:
		airQualityMeasurementProcWriteResponse(terminalInfo, zclFrame.Command, msgID, contentFrame)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandReportAttributes)].Name:
		airQualityMeasurementProcReport(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandConfigureReportingResponse)].Name:
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandDefaultResponse)].Name:
		airQualityMeasurementProcDefaultResponse(terminalInfo, zclFrame.Command, msgID, contentFrame)
	default:
		globallogger.Log.Warnf("[devEUI: %v][AirQualityMeasurementProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
