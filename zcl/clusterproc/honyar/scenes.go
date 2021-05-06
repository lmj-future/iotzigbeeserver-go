package honyar

import (
	"encoding/json"

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

func scenesProcMsg2Kafka(terminalInfo config.TerminalInfo, values map[string]interface{}) {
	kafkaMsg := publicstruct.DataReportMsg{
		OIDIndex:   terminalInfo.OIDIndex,
		DevSN:      terminalInfo.DevEUI,
		LinkType:   terminalInfo.ProfileID,
		DeviceType: terminalInfo.TmnType2,
	}
	if _, ok := values[iotsmartspace.IotwarePropertySwitchLeftUp]; ok {
		type appDataMsg struct {
			LeftUp string `json:"leftUp"`
		}
		kafkaMsg.AppData = appDataMsg{LeftUp: "左上键情景触发"}
		kafkaMsgByte, _ := json.Marshal(kafkaMsg)
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
	if _, ok := values[iotsmartspace.IotwarePropertySwitchLeftMid]; ok {
		type appDataMsg struct {
			LeftMiddle string `json:"leftMiddle"`
		}
		kafkaMsg.AppData = appDataMsg{LeftMiddle: "左中键情景触发"}
		kafkaMsgByte, _ := json.Marshal(kafkaMsg)
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
	if _, ok := values[iotsmartspace.IotwarePropertySwitchLeftDown]; ok {
		type appDataMsg struct {
			LeftDown string `json:"leftDown"`
		}
		kafkaMsg.AppData = appDataMsg{LeftDown: "左下键情景触发"}
		kafkaMsgByte, _ := json.Marshal(kafkaMsg)
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
	if _, ok := values[iotsmartspace.IotwarePropertySwitchRightUp]; ok {
		type appDataMsg struct {
			RightUp string `json:"rightUp"`
		}
		kafkaMsg.AppData = appDataMsg{RightUp: "右上键情景触发"}
		kafkaMsgByte, _ := json.Marshal(kafkaMsg)
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
	if _, ok := values[iotsmartspace.IotwarePropertySwitchRightMid]; ok {
		type appDataMsg struct {
			RightMiddle string `json:"rightMiddle"`
		}
		kafkaMsg.AppData = appDataMsg{RightMiddle: "右中键情景触发"}
		kafkaMsgByte, _ := json.Marshal(kafkaMsg)
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
	if _, ok := values[iotsmartspace.IotwarePropertySwitchRightDown]; ok {
		type appDataMsg struct {
			RightDown string `json:"rightDown"`
		}
		kafkaMsg.AppData = appDataMsg{RightDown: "右下键情景触发"}
		kafkaMsgByte, _ := json.Marshal(kafkaMsg)
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
}

// scenesProcReport 处理report（0x0a）消息
func scenesProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][scenesProcReport]: command: %+v", terminalInfo.DevEUI, Command)
	attributeReports := Command.AttributeReports
	params := make(map[string]interface{}, 2)
	values := make(map[string]interface{}, 1)
	params["terminalId"] = terminalInfo.DevEUI
	var key string
	var bPublish = false
	for _, v := range attributeReports {
		switch v.AttributeName {
		case "ScenesCmd":
			Value := v.Attribute.Value.(uint64)
			switch Value {
			case 1:
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal1SceneSwitch005f0cf1:
					key = iotsmartspace.Honyar1SceneSwitch005f0cf1Property
					values[iotsmartspace.IotwarePropertySwitchOne] = "1"
					bPublish = true
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3:
					key = iotsmartspace.Honyar2SceneSwitch005f0cf3PropertyLeft
					values[iotsmartspace.IotwarePropertySwitchLeft] = "1"
					bPublish = true
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
					key = iotsmartspace.Honyar3SceneSwitch005f0cf2PropertyLeft
					values[iotsmartspace.IotwarePropertySwitchLeft] = "1"
					bPublish = true
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
					key = iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyLeftUp
					values[iotsmartspace.IotwarePropertySwitchLeftUp] = "1"
					bPublish = true
				default:
					globallogger.Log.Warnf("[devEUI: %v][scenesProcReport] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
				}
			case 2:
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3:
					key = iotsmartspace.Honyar2SceneSwitch005f0cf3PropertyRight
					values[iotsmartspace.IotwarePropertySwitchRight] = "1"
					bPublish = true
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
					key = iotsmartspace.Honyar3SceneSwitch005f0cf2PropertyMiddle
					values[iotsmartspace.IotwarePropertySwitchMid] = "1"
					bPublish = true
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
					key = iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyLeftMiddle
					values[iotsmartspace.IotwarePropertySwitchLeftMid] = "1"
					bPublish = true
				default:
					globallogger.Log.Warnf("[devEUI: %v][scenesProcReport] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
				}
			case 3:
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
					key = iotsmartspace.Honyar3SceneSwitch005f0cf2PropertyRight
					values[iotsmartspace.IotwarePropertySwitchRight] = "1"
					bPublish = true
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
					key = iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyLeftDown
					values[iotsmartspace.IotwarePropertySwitchLeftDown] = "1"
					bPublish = true
				default:
					globallogger.Log.Warnf("[devEUI: %v][scenesProcReport] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
				}
			case 4:
				key = iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyRightUp
				values[iotsmartspace.IotwarePropertySwitchRightUp] = "1"
				bPublish = true
			case 5:
				key = iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyRightMiddle
				values[iotsmartspace.IotwarePropertySwitchRightMid] = "1"
				bPublish = true
			case 6:
				key = iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyRightDown
				values[iotsmartspace.IotwarePropertySwitchRightDown] = "1"
				bPublish = true
			default:
				globallogger.Log.Warnf("[devEUI: %v][scenesProcReport] invalid Attribute: %v", terminalInfo.DevEUI, Value)
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][scenesProcReport] invalid v.AttributeName: %v", terminalInfo.DevEUI, v.AttributeName)
		}
	}
	params[key] = "1"
	if bPublish {
		// iotsmartspace publish msg to app
		if constant.Constant.Iotware {
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
		} else if constant.Constant.Iotedge {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
		} else if constant.Constant.Iotprivate {
			scenesProcMsg2Kafka(terminalInfo, values)
		}
	}
}

// ScenesProc 处理clusterID 0xfe05属性消息
func ScenesProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	switch zclFrame.CommandName {
	case "ReportAttributes":
		scenesProcReport(terminalInfo, zclFrame.Command)
	default:
		globallogger.Log.Warnf("[devEUI: %v][ScenesProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
