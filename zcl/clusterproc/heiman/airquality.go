package heiman

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/db/kafka"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
	uuid "github.com/satori/go.uuid"
)

func airQualityMeasurementProcReadRspIotware(terminalInfo config.TerminalInfo, command interface{}) {
	oSet := bson.M{}
	oSetPG := make(map[string]interface{})
	var bPublish = false
	var value iotsmartspace.IotwarePropertyPM10PowerStateAQI = iotsmartspace.IotwarePropertyPM10PowerStateAQI{}
	for _, v := range command.(*cluster.ReadAttributesResponse).ReadAttributeStatuses {
		globallogger.Log.Infof("[devEUI: %v][airQualityMeasurementProcReadRspIotware]: readAttributesRsp: %+v", terminalInfo.DevEUI, v)
		if v.Status != cluster.ZclStatusSuccess {
			continue
		}
		switch v.AttributeName {
		case "Language":
			Value := v.Attribute.Value.(uint64)
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
			Value := v.Attribute.Value.(uint64)
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
			Value := v.Attribute.Value.(uint64)
			if Value == 0 {
				oSet["batteryState"] = "unCharge"
				oSetPG["batterystate"] = "unCharge"
				value.PowerState = "0"
				bPublish = true
			} else if Value == 1 {
				oSet["batteryState"] = "charging"
				oSetPG["batterystate"] = "charging"
				value.PowerState = "1"
				bPublish = true
			} else if Value == 2 {
				oSet["batteryState"] = "fullCharged"
				oSetPG["batterystate"] = "fullCharged"
				value.PowerState = "2"
				bPublish = true
			} else {
				oSet["batteryState"] = "Invalid"
				oSetPG["batterystate"] = "Invalid"
			}
		case "PM10MeasuredValue":
			value.PM10 = strconv.FormatUint(v.Attribute.Value.(uint64), 10)
			oSet["PM10"] = value.PM10
			oSetPG["pm10"] = value.PM10
			bPublish = true
		case "TVOCMeasuredValue":
			Value := v.Attribute.Value.(uint64)
			oSet["TVOC"] = strconv.FormatUint(Value, 10)
			oSetPG["tvoc"] = strconv.FormatUint(Value, 10)
		case "AQIMeasuredValue":
			value.AQI = strconv.FormatUint(v.Attribute.Value.(uint64), 10)
			oSet["AQI"] = value.AQI
			oSetPG["aqi"] = value.AQI
			bPublish = true
		case "MaxTemperature":
			Value := v.Attribute.Value.(int64)
			oSet["maxTemperature"] = strconv.FormatInt(Value, 10)
			oSetPG["maxtemperature"] = strconv.FormatInt(Value, 10)
		case "MinTemperature":
			Value := v.Attribute.Value.(int64)
			oSet["minTemperature"] = strconv.FormatInt(Value, 10)
			oSetPG["mintemperature"] = strconv.FormatInt(Value, 10)
		case "MaxHumidity":
			Value := v.Attribute.Value.(uint64)
			oSet["maxHumidity"] = strconv.FormatUint(Value, 10)
			oSetPG["maxhumidity"] = strconv.FormatUint(Value, 10)
		case "MinHumidity":
			Value := v.Attribute.Value.(uint64)
			oSet["minHumidity"] = strconv.FormatUint(Value, 10)
			oSetPG["minhumidity"] = strconv.FormatUint(Value, 10)
		case "Disturb":
			Value := v.Attribute.Value.(uint64)
			disturb := strconv.FormatUint(Value, 16)
			disturb = "0000" + disturb
			disturb = disturb[len(disturb)-4:]
			oSet["disturb"] = disturb
			oSetPG["disturb"] = disturb
		default:
			globallogger.Log.Warnf("[devEUI: %v][airQualityMeasurementProcReadRspIotware] invalid attributeName: %v", terminalInfo.DevEUI, v.AttributeName)
		}
	}
	if oSet["disturb"] != nil || oSet["minHumidity"] != nil || oSet["maxHumidity"] != nil || oSet["minTemperature"] != nil || oSet["maxTemperature"] != nil {
		if constant.Constant.UsePostgres {
			models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": terminalInfo.DevEUI}, oSetPG)
		} else {
			models.FindTerminalAndUpdate(bson.M{"devEUI": terminalInfo.DevEUI}, oSet)
		}
		var key = terminalInfo.APMac + terminalInfo.ModuleID + terminalInfo.DevEUI
		if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
			key = terminalInfo.APMac + terminalInfo.ModuleID + terminalInfo.NwkAddr
		}
		publicfunction.DeleteTerminalInfoListCache(key)
	}

	if bPublish {
		iotsmartspace.PublishTelemetryUpIotware(terminalInfo, value)
	}
}
func airQualityMeasurementProcReadRspIotedge(terminalInfo config.TerminalInfo, command interface{}) {
	oSet := bson.M{}
	oSetPG := make(map[string]interface{})
	var bPublish = false
	var value iotsmartspace.HeimanHS2AQProperty = iotsmartspace.HeimanHS2AQProperty{
		DevEUI: terminalInfo.DevEUI,
	}
	for _, v := range command.(*cluster.ReadAttributesResponse).ReadAttributeStatuses {
		globallogger.Log.Infof("[devEUI: %v][airQualityMeasurementProcReadRspIotedge]: readAttributesRsp: %+v", terminalInfo.DevEUI, v)
		if v.Status != cluster.ZclStatusSuccess {
			continue
		}
		switch v.AttributeName {
		case "Language":
			Value := v.Attribute.Value.(uint64)
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
			Value := v.Attribute.Value.(uint64)
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
			Value := v.Attribute.Value.(uint64)
			if Value == 0 {
				oSet["batteryState"] = "unCharge"
				oSetPG["batterystate"] = "unCharge"
				value.BatteryState = "0"
				bPublish = true
			} else if Value == 1 {
				oSet["batteryState"] = "charging"
				oSetPG["batterystate"] = "charging"
				value.BatteryState = "1"
				bPublish = true
			} else if Value == 2 {
				oSet["batteryState"] = "fullCharged"
				oSetPG["batterystate"] = "fullCharged"
				value.BatteryState = "2"
				bPublish = true
			} else {
				oSet["batteryState"] = "Invalid"
				oSetPG["batterystate"] = "Invalid"
			}
		case "PM10MeasuredValue":
			value.PM10 = strconv.FormatUint(v.Attribute.Value.(uint64), 10)
			oSet["PM10"] = value.PM10
			oSetPG["pm10"] = value.PM10
			bPublish = true
		case "TVOCMeasuredValue":
			value.TVOC = strconv.FormatUint(v.Attribute.Value.(uint64), 10)
			oSet["TVOC"] = value.TVOC
			oSetPG["tvoc"] = value.TVOC
			bPublish = true
		case "AQIMeasuredValue":
			value.AQI = strconv.FormatUint(v.Attribute.Value.(uint64), 10)
			oSet["AQI"] = value.AQI
			oSetPG["aqi"] = value.AQI
			bPublish = true
		case "MaxTemperature":
			Value := v.Attribute.Value.(int64)
			oSet["maxTemperature"] = strconv.FormatInt(Value, 10)
			oSetPG["maxtemperature"] = strconv.FormatInt(Value, 10)
		case "MinTemperature":
			Value := v.Attribute.Value.(int64)
			oSet["minTemperature"] = strconv.FormatInt(Value, 10)
			oSetPG["mintemperature"] = strconv.FormatInt(Value, 10)
		case "MaxHumidity":
			Value := v.Attribute.Value.(uint64)
			oSet["maxHumidity"] = strconv.FormatUint(Value, 10)
			oSetPG["maxhumidity"] = strconv.FormatUint(Value, 10)
		case "MinHumidity":
			Value := v.Attribute.Value.(uint64)
			oSet["minHumidity"] = strconv.FormatUint(Value, 10)
			oSetPG["minhumidity"] = strconv.FormatUint(Value, 10)
		case "Disturb":
			Value := v.Attribute.Value.(uint64)
			disturb := strconv.FormatUint(Value, 16)
			disturb = "0000" + disturb
			disturb = disturb[len(disturb)-4:]
			oSet["disturb"] = disturb
			oSetPG["disturb"] = disturb
		default:
			globallogger.Log.Warnf("[devEUI: %v][airQualityMeasurementProcReadRspIotedge] invalid attributeName: %v", terminalInfo.DevEUI, v.AttributeName)
		}
	}
	if oSet["disturb"] != nil || oSet["minHumidity"] != nil || oSet["maxHumidity"] != nil || oSet["minTemperature"] != nil || oSet["maxTemperature"] != nil {
		if constant.Constant.UsePostgres {
			models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": terminalInfo.DevEUI}, oSetPG)
		} else {
			models.FindTerminalAndUpdate(bson.M{"devEUI": terminalInfo.DevEUI}, oSet)
		}
		var key = terminalInfo.APMac + terminalInfo.ModuleID + terminalInfo.DevEUI
		if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
			key = terminalInfo.APMac + terminalInfo.ModuleID + terminalInfo.NwkAddr
		}
		publicfunction.DeleteTerminalInfoListCache(key)
	}

	if bPublish {
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, value, uuid.NewV4().String())
	}
}
func airQualityMeasurementProcReadRspIotprivate(terminalInfo config.TerminalInfo, command interface{}) {
	oSet := bson.M{}
	oSetPG := make(map[string]interface{})
	for _, v := range command.(*cluster.ReadAttributesResponse).ReadAttributeStatuses {
		globallogger.Log.Infof("[devEUI: %v][airQualityMeasurementProcReadRspIotprivate]: readAttributesRsp: %+v", terminalInfo.DevEUI, v)
		if v.Status != cluster.ZclStatusSuccess {
			continue
		}
		switch v.AttributeName {
		case "Language":
			Value := v.Attribute.Value.(uint64)
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
			Value := v.Attribute.Value.(uint64)
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
			type appDataMsg struct {
				BatteryState string `json:"batteryState"`
			}
			Value := v.Attribute.Value.(uint64)
			var batteryState string
			if Value == 0 {
				oSet["batteryState"] = "unCharge"
				oSetPG["batterystate"] = "unCharge"
				batteryState = "未充电"
			} else if Value == 1 {
				oSet["batteryState"] = "charging"
				oSetPG["batterystate"] = "charging"
				batteryState = "充电中"
			} else if Value == 2 {
				oSet["batteryState"] = "fullCharged"
				oSetPG["batterystate"] = "fullCharged"
				batteryState = "已充满"
			} else {
				oSet["batteryState"] = "Invalid"
				oSetPG["batterystate"] = "Invalid"
			}
			kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
				Time:       time.Now(),
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData:    appDataMsg{BatteryState: batteryState},
			})
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		case "PM10MeasuredValue":
			Value := v.Attribute.Value.(uint64)
			oSet["PM10"] = strconv.FormatUint(Value, 10)
			oSetPG["pm10"] = strconv.FormatUint(Value, 10)
			type appDataMsg struct {
				PM10 string `json:"PM10"`
			}
			kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
				Time:       time.Now(),
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData:    appDataMsg{PM10: strconv.FormatUint(Value, 10) + "ug/m^3"},
			})
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		case "TVOCMeasuredValue":
			Value := v.Attribute.Value.(uint64)
			oSet["TVOC"] = strconv.FormatUint(Value, 10)
			oSetPG["tvoc"] = strconv.FormatUint(Value, 10)
		case "AQIMeasuredValue":
			Value := v.Attribute.Value.(uint64)
			oSet["AQI"] = strconv.FormatUint(Value, 10)
			oSetPG["aqi"] = strconv.FormatUint(Value, 10)
			type appDataMsg struct {
				AQI string `json:"AQI"`
			}
			kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
				Time:       time.Now(),
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData:    appDataMsg{AQI: strconv.FormatUint(Value, 10)},
			})
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		case "MaxTemperature":
			Value := v.Attribute.Value.(int64)
			oSet["maxTemperature"] = strconv.FormatInt(Value, 10)
			oSetPG["maxtemperature"] = strconv.FormatInt(Value, 10)
		case "MinTemperature":
			Value := v.Attribute.Value.(int64)
			oSet["minTemperature"] = strconv.FormatInt(Value, 10)
			oSetPG["mintemperature"] = strconv.FormatInt(Value, 10)
		case "MaxHumidity":
			Value := v.Attribute.Value.(uint64)
			oSet["maxHumidity"] = strconv.FormatUint(Value, 10)
			oSetPG["maxhumidity"] = strconv.FormatUint(Value, 10)
		case "MinHumidity":
			Value := v.Attribute.Value.(uint64)
			oSet["minHumidity"] = strconv.FormatUint(Value, 10)
			oSetPG["minhumidity"] = strconv.FormatUint(Value, 10)
		case "Disturb":
			Value := v.Attribute.Value.(uint64)
			disturb := strconv.FormatUint(Value, 16)
			disturb = "0000" + disturb
			disturb = disturb[len(disturb)-4:]
			oSet["disturb"] = disturb
			oSetPG["disturb"] = disturb
		default:
			globallogger.Log.Warnf("[devEUI: %v][airQualityMeasurementProcReadRspIotprivate] invalid attributeName: %v", terminalInfo.DevEUI, v.AttributeName)
		}
	}
	if oSet["disturb"] != nil || oSet["minHumidity"] != nil || oSet["maxHumidity"] != nil || oSet["minTemperature"] != nil || oSet["maxTemperature"] != nil {
		if constant.Constant.UsePostgres {
			models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": terminalInfo.DevEUI}, oSetPG)
		} else {
			models.FindTerminalAndUpdate(bson.M{"devEUI": terminalInfo.DevEUI}, oSet)
		}
		var key = terminalInfo.APMac + terminalInfo.ModuleID + terminalInfo.DevEUI
		if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
			key = terminalInfo.APMac + terminalInfo.ModuleID + terminalInfo.NwkAddr
		}
		publicfunction.DeleteTerminalInfoListCache(key)
	}
}

//airQualityMeasurementProcReadRsp 处理readRsp（0x01）消息
func airQualityMeasurementProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	if constant.Constant.Iotware {
		airQualityMeasurementProcReadRspIotware(terminalInfo, command)
	} else if constant.Constant.Iotedge {
		airQualityMeasurementProcReadRspIotedge(terminalInfo, command)
	} else if constant.Constant.Iotprivate {
		airQualityMeasurementProcReadRspIotprivate(terminalInfo, command)
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
			for _, v := range contentFrame.Command.(*cluster.WriteAttributesCommand).WriteAttributeRecords {
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
			var key = terminalInfo.APMac + terminalInfo.ModuleID + terminalInfo.DevEUI
			if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
				key = terminalInfo.APMac + terminalInfo.ModuleID + terminalInfo.NwkAddr
			}
			publicfunction.DeleteTerminalInfoListCache(key)
			if constant.Constant.UsePostgres {
				models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": terminalInfo.DevEUI}, oSetPG)
			} else {
				models.FindTerminalAndUpdate(bson.M{"devEUI": terminalInfo.DevEUI}, oSet)
			}
		}
	}
	if constant.Constant.Iotware {
		iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, message, msgID)
	} else if constant.Constant.Iotedge {
		type response struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply, response{Code: code, Message: message}, msgID)
	}
}

// airQualityMeasurementProcReport 处理report（0x0a）消息
func airQualityMeasurementProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	globallogger.Log.Infof("[devEUI: %v][airQualityMeasurementProcReport]: command: %+v", terminalInfo.DevEUI, command.(*cluster.ReportAttributesCommand))
	for _, v := range command.(*cluster.ReportAttributesCommand).AttributeReports {
		switch v.AttributeName {
		case "BatteryState":
			var airState = "3"
			if v.Attribute.Value.(uint64) == 0 {
				airState = "0"
			} else if v.Attribute.Value.(uint64) == 1 {
				airState = "1"
			} else if v.Attribute.Value.(uint64) == 2 {
				airState = "2"
			}
			if constant.Constant.Iotware {
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyPowerState{PowerState: airState})
			} else if constant.Constant.Iotedge {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyBatteryState{DevEUI: terminalInfo.DevEUI, BatteryState: airState}, uuid.NewV4().String())
			} else if constant.Constant.Iotprivate {
				type appDataMsg struct {
					BatteryState string `json:"batteryState"`
				}
				if airState == "0" {
					airState = "未充电"
				} else if airState == "1" {
					airState = "充电中"
				} else if airState == "2" {
					airState = "已充满"
				}
				kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
					Time:       time.Now(),
					OIDIndex:   terminalInfo.OIDIndex,
					DevSN:      terminalInfo.DevEUI,
					LinkType:   terminalInfo.ProfileID,
					DeviceType: terminalInfo.TmnType2,
					AppData:    appDataMsg{BatteryState: airState},
				})
				kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
			}
		case "PM10MeasuredValue":
			if constant.Constant.Iotware {
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyPM10{PM10: strconv.FormatUint(v.Attribute.Value.(uint64), 10)})
			} else if constant.Constant.Iotedge {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyPM10{DevEUI: terminalInfo.DevEUI, PM10: strconv.FormatUint(v.Attribute.Value.(uint64), 10)}, uuid.NewV4().String())
			} else if constant.Constant.Iotprivate {
				type appDataMsg struct {
					PM10 string `json:"PM10"`
				}
				kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
					Time:       time.Now(),
					OIDIndex:   terminalInfo.OIDIndex,
					DevSN:      terminalInfo.DevEUI,
					LinkType:   terminalInfo.ProfileID,
					DeviceType: terminalInfo.TmnType2,
					AppData:    appDataMsg{PM10: strconv.FormatUint(v.Attribute.Value.(uint64), 10) + "ug/m^3"},
				})
				kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
			}
		case "AQIMeasuredValue":
			if constant.Constant.Iotware {
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyAQI{AQI: strconv.FormatUint(v.Attribute.Value.(uint64), 10)})
			} else if constant.Constant.Iotedge {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2AQPropertyAQI{DevEUI: terminalInfo.DevEUI, AQI: strconv.FormatUint(v.Attribute.Value.(uint64), 10)}, uuid.NewV4().String())
			} else if constant.Constant.Iotprivate {
				type appDataMsg struct {
					AQI string `json:"AQI"`
				}
				kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
					Time:       time.Now(),
					OIDIndex:   terminalInfo.OIDIndex,
					DevSN:      terminalInfo.DevEUI,
					LinkType:   terminalInfo.ProfileID,
					DeviceType: terminalInfo.TmnType2,
					AppData:    appDataMsg{AQI: strconv.FormatUint(v.Attribute.Value.(uint64), 10)},
				})
				kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][airQualityMeasurementProcReport] invalid attributeName: %v", terminalInfo.DevEUI, v.AttributeName)
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
	switch zclFrame.CommandName {
	case "ReadAttributesResponse":
		airQualityMeasurementProcReadRsp(terminalInfo, zclFrame.Command)
	case "WriteAttributesResponse":
		airQualityMeasurementProcWriteResponse(terminalInfo, zclFrame.Command, msgID, contentFrame)
	case "ReportAttributes":
		airQualityMeasurementProcReport(terminalInfo, zclFrame.Command)
	case "DefaultResponse":
		airQualityMeasurementProcDefaultResponse(terminalInfo, zclFrame.Command, msgID, contentFrame)
	default:
		globallogger.Log.Warnf("[devEUI: %v][AirQualityMeasurementProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
