package honyar

import (
	"encoding/json"
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

func scenesProcReportIotware(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][scenesProcReportIotware]: command: %+v", terminalInfo.DevEUI, Command)
	for _, v := range Command.AttributeReports {
		switch v.AttributeName {
		case "ScenesCmd":
			switch v.Attribute.Value.(uint64) {
			case 1:
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal1SceneSwitch005f0cf1:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
					iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertySwitchLeftUp{SwitchLeftUp: "1"})
				default:
					globallogger.Log.Warnf("[devEUI: %v][scenesProcReportIotware] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
				}
			case 2:
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
					iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertySwitchLeftMid{SwitchLeftMid: "1"})
				default:
					globallogger.Log.Warnf("[devEUI: %v][scenesProcReportIotware] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
				}
			case 3:
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
					iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertySwitchLeftDown{SwitchLeftDown: "1"})
				default:
					globallogger.Log.Warnf("[devEUI: %v][scenesProcReportIotware] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
				}
			case 4:
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertySwitchRightUp{SwitchRightUp: "1"})
			case 5:
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertySwitchRightMid{SwitchRightMid: "1"})
			case 6:
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertySwitchRightDown{SwitchRightDown: "1"})
			default:
				globallogger.Log.Warnf("[devEUI: %v][scenesProcReportIotware] invalid Attribute: %v", terminalInfo.DevEUI, v.Attribute.Value.(uint64))
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][scenesProcReportIotware] invalid v.AttributeName: %v", terminalInfo.DevEUI, v.AttributeName)
		}
	}
}
func scenesProcReportIotedge(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][scenesProcReportIotedge]: command: %+v", terminalInfo.DevEUI, Command)
	for _, v := range Command.AttributeReports {
		switch v.AttributeName {
		case "ScenesCmd":
			switch v.Attribute.Value.(uint64) {
			case 1:
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal1SceneSwitch005f0cf1:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
					iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
						iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyLeftUp{DevEUI: terminalInfo.DevEUI, LeftUp: "1"}, uuid.NewV4().String())
				default:
					globallogger.Log.Warnf("[devEUI: %v][scenesProcReportIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
				}
			case 2:
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
					iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
						iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyLeftMiddle{DevEUI: terminalInfo.DevEUI, LeftMiddle: "1"}, uuid.NewV4().String())
				default:
					globallogger.Log.Warnf("[devEUI: %v][scenesProcReportIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
				}
			case 3:
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
					iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
						iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyLeftDown{DevEUI: terminalInfo.DevEUI, LeftDown: "1"}, uuid.NewV4().String())
				default:
					globallogger.Log.Warnf("[devEUI: %v][scenesProcReportIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
				}
			case 4:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyRightUp{DevEUI: terminalInfo.DevEUI, RightUp: "1"}, uuid.NewV4().String())
			case 5:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyRightMiddle{DevEUI: terminalInfo.DevEUI, RightMiddle: "1"}, uuid.NewV4().String())
			case 6:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyRightDown{DevEUI: terminalInfo.DevEUI, RightDown: "1"}, uuid.NewV4().String())
			default:
				globallogger.Log.Warnf("[devEUI: %v][scenesProcReportIotedge] invalid Attribute: %v", terminalInfo.DevEUI, v.Attribute.Value.(uint64))
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][scenesProcReportIotedge] invalid v.AttributeName: %v", terminalInfo.DevEUI, v.AttributeName)
		}
	}
}
func scenesProcReportIotprivate(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][scenesProcReportIotprivate]: command: %+v", terminalInfo.DevEUI, Command)
	kafkaMsg := publicstruct.DataReportMsg{
		Time:       time.Now(),
		OIDIndex:   terminalInfo.OIDIndex,
		DevSN:      terminalInfo.DevEUI,
		LinkType:   terminalInfo.ProfileID,
		DeviceType: terminalInfo.TmnType2,
	}
	for _, v := range Command.AttributeReports {
		switch v.AttributeName {
		case "ScenesCmd":
			switch v.Attribute.Value.(uint64) {
			case 1:
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal1SceneSwitch005f0cf1:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
					type appDataMsg struct {
						LeftUp string `json:"leftUp"`
					}
					kafkaMsg.AppData = appDataMsg{LeftUp: "左上键情景触发"}
				default:
					globallogger.Log.Warnf("[devEUI: %v][scenesProcReportIotprivate] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
				}
			case 2:
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
					type appDataMsg struct {
						LeftMiddle string `json:"leftMiddle"`
					}
					kafkaMsg.AppData = appDataMsg{LeftMiddle: "左中键情景触发"}
				default:
					globallogger.Log.Warnf("[devEUI: %v][scenesProcReportIotprivate] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
				}
			case 3:
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
					type appDataMsg struct {
						LeftDown string `json:"leftDown"`
					}
					kafkaMsg.AppData = appDataMsg{LeftDown: "左下键情景触发"}
				default:
					globallogger.Log.Warnf("[devEUI: %v][scenesProcReportIotprivate] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
				}
			case 4:
				type appDataMsg struct {
					RightUp string `json:"rightUp"`
				}
				kafkaMsg.AppData = appDataMsg{RightUp: "右上键情景触发"}
			case 5:
				type appDataMsg struct {
					RightMiddle string `json:"rightMiddle"`
				}
				kafkaMsg.AppData = appDataMsg{RightMiddle: "右中键情景触发"}
			case 6:
				type appDataMsg struct {
					RightDown string `json:"rightDown"`
				}
				kafkaMsg.AppData = appDataMsg{RightDown: "右下键情景触发"}
			default:
				globallogger.Log.Warnf("[devEUI: %v][scenesProcReportIotprivate] invalid Attribute: %v", terminalInfo.DevEUI, v.Attribute.Value.(uint64))
			}
		default:
			globallogger.Log.Warnf("[devEUI: %v][scenesProcReportIotprivate] invalid v.AttributeName: %v", terminalInfo.DevEUI, v.AttributeName)
		}
	}
	kafkaMsgByte, _ := json.Marshal(kafkaMsg)
	kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
}

// scenesProcReport 处理report（0x0a）消息
func scenesProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	if constant.Constant.Iotware {
		scenesProcReportIotware(terminalInfo, command)
	} else if constant.Constant.Iotedge {
		scenesProcReportIotedge(terminalInfo, command)
	} else if constant.Constant.Iotprivate {
		scenesProcReportIotprivate(terminalInfo, command)
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
