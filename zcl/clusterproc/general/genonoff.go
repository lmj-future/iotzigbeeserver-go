package general

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/db/kafka"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/keepalive"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
	"github.com/h3c/iotzigbeeserver-go/zcl/zclmsgdown"
	uuid "github.com/satori/go.uuid"
)

func getGenOnOffParams(terminalID string, tmnType string, value string, endpoint string) (interface{}, bool) {
	var bPublish = false
	params := make(map[string]interface{}, 2)
	var key string
	params["terminalId"] = terminalID
	switch tmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
		key = iotsmartspace.HeimanESocketPropertyOnOff
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
		key = iotsmartspace.HeimanSmartPlugPropertyOnOff
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30:
		key = iotsmartspace.HeimanHS2SW1LEFR30PropertyOnOff1
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30:
		if endpoint == "0" {
			key = iotsmartspace.HeimanHS2SW2LEFR30PropertyOnOff1
			bPublish = true
		} else if endpoint == "1" {
			key = iotsmartspace.HeimanHS2SW2LEFR30PropertyOnOff2
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30:
		if endpoint == "0" {
			key = iotsmartspace.HeimanHS2SW3LEFR30PropertyOnOff1
			bPublish = true
		} else if endpoint == "1" {
			key = iotsmartspace.HeimanHS2SW3LEFR30PropertyOnOff2
			bPublish = true
		} else if endpoint == "2" {
			key = iotsmartspace.HeimanHS2SW3LEFR30PropertyOnOff3
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32:
		key = iotsmartspace.HonyarSingleSwitch00500c32PropertyOnOff1
		bPublish = true
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33:
		if endpoint == "0" {
			key = iotsmartspace.HonyarDoubleSwitch00500c33PropertyOnOff1
			bPublish = true
		} else if endpoint == "1" {
			key = iotsmartspace.HonyarDoubleSwitch00500c33PropertyOnOff2
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35:
		if endpoint == "0" {
			key = iotsmartspace.HonyarTripleSwitch00500c35PropertyOnOff1
			bPublish = true
		} else if endpoint == "1" {
			key = iotsmartspace.HonyarTripleSwitch00500c35PropertyOnOff2
			bPublish = true
		} else if endpoint == "2" {
			key = iotsmartspace.HonyarTripleSwitch00500c35PropertyOnOff3
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141:
		bPublish = true
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142:
		if endpoint == "0" {
			bPublish = true
		} else if endpoint == "1" {
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
		if endpoint == "0" {
			bPublish = true
		} else if endpoint == "1" {
			bPublish = true
		} else if endpoint == "2" {
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
		key = iotsmartspace.HonyarSocket000a0c3cPropertyOnOff
		bPublish = true
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
		key = iotsmartspace.HonyarSocketHY0106PropertyOnOff
		bPublish = true
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
		key = iotsmartspace.HonyarSocketHY0105PropertyOnOff
		bPublish = true
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
		key = iotsmartspace.HonyarSocketHY0106PropertyOnOff
		bPublish = true
	default:
		globallogger.Log.Warnf("[devEUI: %v][getGenOnOffParams] invalid tmnType: %v", terminalID, tmnType)
	}
	params[key] = value
	return params, bPublish
}

func genOnOffProcReadRspOrReport(terminalInfo config.TerminalInfo, srcEndpoint uint8, Value bool, attributeName string) {
	var bPublish = false
	var appData string
	if (terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c ||
		terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55 ||
		terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105 ||
		terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106) &&
		(srcEndpoint == 2 || attributeName == "ChildLock" || attributeName == "BackGroundLight" || attributeName == "PowerOffMemory") {
		params := make(map[string]interface{}, 2)
		var key string
		var iotwareKey string
		params["terminalId"] = terminalInfo.DevEUI
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
			key = iotsmartspace.HonyarSocket000a0c3cPropertyUSB
			switch attributeName {
			case "OnOff":
				key = iotsmartspace.HonyarSocket000a0c3cPropertyUSB
				iotwareKey = iotsmartspace.IotwarePropertyUSBSwitch
				bPublish = true
			case "ChildLock":
				key = iotsmartspace.HonyarSocket000a0c3cPropertyChildLock
				iotwareKey = iotsmartspace.IotwarePropertyChildLock
				bPublish = true
			}
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
			switch attributeName {
			case "OnOff":
				bPublish = false
			case "ChildLock":
				key = iotsmartspace.HonyarSocketHY0106PropertyChildLock
				iotwareKey = iotsmartspace.IotwarePropertyChildLock
				bPublish = true
			}
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
			switch attributeName {
			case "OnOff":
				key = iotsmartspace.HonyarSocketHY0105PropertyUSB
				iotwareKey = iotsmartspace.IotwarePropertyUSBSwitch
				bPublish = true
			case "ChildLock":
				key = iotsmartspace.HonyarSocketHY0105PropertyChildLock
				iotwareKey = iotsmartspace.IotwarePropertyChildLock
				bPublish = true
			case "BackGroundLight":
				key = iotsmartspace.HonyarSocketHY0105PropertyBackGroundLight
				iotwareKey = iotsmartspace.IotwarePropertyBackLightSwitch
				bPublish = true
			case "PowerOffMemory":
				key = iotsmartspace.HonyarSocketHY0105PropertyPowerOffMemory
				iotwareKey = iotsmartspace.IotwarePropertyMemorySwitch
				bPublish = true
			}
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
			switch attributeName {
			case "OnOff":
				bPublish = false
			case "ChildLock":
				key = iotsmartspace.HonyarSocketHY0106PropertyChildLock
				iotwareKey = iotsmartspace.IotwarePropertyChildLock
				bPublish = true
			case "BackGroundLight":
				key = iotsmartspace.HonyarSocketHY0106PropertyBackGroundLight
				iotwareKey = iotsmartspace.IotwarePropertyBackLightSwitch
				bPublish = true
			case "PowerOffMemory":
				key = iotsmartspace.HonyarSocketHY0106PropertyPowerOffMemory
				iotwareKey = iotsmartspace.IotwarePropertyMemorySwitch
				bPublish = true
			}
		}
		if Value {
			params[key] = "01"
			appData = "开"
			if attributeName == "PowerOffMemory" {
				params[key] = "00"
				appData = "关"
			}
		} else {
			params[key] = "00"
			appData = "关"
			if attributeName == "PowerOffMemory" {
				params[key] = "01"
				appData = "开"
			}
		}
		if bPublish {
			if constant.Constant.Iotware {
				values := make(map[string]interface{}, 1)
				values[iotwareKey] = params[key]
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
			} else if constant.Constant.Iotedge {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
			} else if constant.Constant.Iotprivate {
				kafkaMsg := publicstruct.DataReportMsg{
					OIDIndex:   terminalInfo.OIDIndex,
					DevSN:      terminalInfo.DevEUI,
					LinkType:   terminalInfo.ProfileID,
					DeviceType: terminalInfo.TmnType2,
				}
				switch attributeName {
				case "OnOff":
					type appDataMsg struct {
						USBState string `json:"USBState"`
					}
					kafkaMsg.AppData = appDataMsg{USBState: appData}
				case "ChildLock":
					type appDataMsg struct {
						ChildLockState string `json:"childLockState"`
					}
					kafkaMsg.AppData = appDataMsg{ChildLockState: appData}
				case "BackGroundLight":
					type appDataMsg struct {
						BackLightState string `json:"backLightState"`
					}
					kafkaMsg.AppData = appDataMsg{BackLightState: appData}
				case "PowerOffMemory":
					type appDataMsg struct {
						MemoryState string `json:"memoryState"`
					}
					kafkaMsg.AppData = appDataMsg{MemoryState: appData}
				}
				kafkaMsgByte, _ := json.Marshal(kafkaMsg)
				kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
			}
		}
	} else {
		value := ""
		// Value == true: ON; Value == false: OFF
		if Value {
			value = "01"
			appData = "开"
		} else {
			value = "00"
			appData = "关"
		}
		endpoint := strconv.FormatUint(uint64(srcEndpoint-1), 10)

		if constant.Constant.Iotware {
			values := make(map[string]interface{}, 1)
			var key string
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30,
				constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30,
				constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
				if endpoint == "0" {
					key = iotsmartspace.IotwarePropertyOnOffOne
					bPublish = true
				} else if endpoint == "1" {
					key = iotsmartspace.IotwarePropertyOnOffTwo
					bPublish = true
				} else if endpoint == "2" {
					key = iotsmartspace.IotwarePropertyOnOffThree
					bPublish = true
				}
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket,
				constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
				key = iotsmartspace.IotwarePropertyOnOff
				bPublish = true
			}
			if bPublish {
				values[key] = value
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
			}
		} else if constant.Constant.Iotedge {
			params, bPublish := getGenOnOffParams(terminalInfo.DevEUI, terminalInfo.TmnType, value, endpoint)
			if bPublish {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
			}
		} else if constant.Constant.Iotprivate {
			kafkaMsg := publicstruct.DataReportMsg{
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
			}
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30, constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
				if endpoint == "0" {
					type appDataMsg struct {
						OnOffOneState string `json:"onOffOneState"`
					}
					kafkaMsg.AppData = appDataMsg{OnOffOneState: "一键" + appData}
				} else if endpoint == "1" {
					type appDataMsg struct {
						OnOffTwoState string `json:"onOffTwoState"`
					}
					kafkaMsg.AppData = appDataMsg{OnOffTwoState: "二键" + appData}
				} else if endpoint == "2" {
					type appDataMsg struct {
						OnOffThreeState string `json:"onOffThreeState"`
					}
					kafkaMsg.AppData = appDataMsg{OnOffThreeState: "三键" + appData}
				}
			default:
				type appDataMsg struct {
					OnOffState string `json:"onOffState"`
				}
				kafkaMsg.AppData = appDataMsg{OnOffState: appData}
			}
			kafkaMsgByte, _ := json.Marshal(kafkaMsg)
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	}
}

//genOnOffProcReadRsp 处理readRsp（0x01）消息
func genOnOffProcReadRsp(terminalInfo config.TerminalInfo, srcEndpoint uint8, command interface{}) {
	readAttributesRsp := command.(*cluster.ReadAttributesResponse)
	for _, v := range readAttributesRsp.ReadAttributeStatuses {
		globallogger.Log.Infof("[devEUI: %v][genOnOffProcReadRsp]: readAttributesRsp: %+v", terminalInfo.DevEUI, v)
		switch v.AttributeName {
		case "OnOff", "ChildLock", "BackGroundLight", "PowerOffMemory":
			if v.Status == cluster.ZclStatusSuccess {
				genOnOffProcReadRspOrReport(terminalInfo, srcEndpoint, v.Attribute.Value.(bool), v.AttributeName)
			}
		default:
			globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "genOnOffProcReadRsp unknow attributeName", v.AttributeName)
		}
	}
}

// genOnOffProcWriteRsp 处理writeRsp（0x04）消息
func genOnOffProcWriteRsp(terminalInfo config.TerminalInfo, srcEndpoint uint8, command interface{}, msgID interface{}, contentFrame *zcl.Frame) {
	Command := command.(*cluster.WriteAttributesResponse)
	globallogger.Log.Infof("[devEUI: %v][genOnOffProcWriteRsp]: command: %+v", terminalInfo.DevEUI, Command)
	code := 0
	message := "success"
	if Command.WriteAttributeStatuses[0].Status != cluster.ZclStatusSuccess {
		code = -1
		message = "failed"
	}
	params := make(map[string]interface{}, 4)
	params["code"] = code
	params["message"] = message
	params["terminalId"] = terminalInfo.DevEUI
	var key string
	var value string
	var attributeName string
	var bPublish = false
	if contentFrame != nil && contentFrame.CommandName == "WriteAttributes" {
		contentCommand := contentFrame.Command.(*cluster.WriteAttributesCommand)
		if contentCommand.WriteAttributeRecords[0].Attribute.Value.(bool) {
			value = "01"
		} else {
			value = "00"
		}
		attributeName = contentCommand.WriteAttributeRecords[0].AttributeName
	}
	switch attributeName {
	case "ChildLock":
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
			key = iotsmartspace.HonyarSocket000a0c3cPropertyChildLock
			bPublish = true
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
			key = iotsmartspace.HonyarSocketHY0106PropertyChildLock
			bPublish = true
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
			key = iotsmartspace.HonyarSocketHY0105PropertyChildLock
			bPublish = true
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
			key = iotsmartspace.HonyarSocketHY0106PropertyChildLock
			bPublish = true
		}
	case "BackGroundLight":
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
			key = iotsmartspace.HonyarSocketHY0105PropertyBackGroundLight
			bPublish = true
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
			key = iotsmartspace.HonyarSocketHY0106PropertyBackGroundLight
			bPublish = true
		}
	case "PowerOffMemory":
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
			key = iotsmartspace.HonyarSocketHY0105PropertyPowerOffMemory
			bPublish = true
			if value == "00" {
				value = "01"
			} else if value == "01" {
				value = "00"
			}
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
			key = iotsmartspace.HonyarSocketHY0106PropertyPowerOffMemory
			bPublish = true
			if value == "00" {
				value = "01"
			} else if value == "01" {
				value = "00"
			}
		}
	}
	params[key] = value
	if bPublish {
		if constant.Constant.Iotware {
			iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, message, msgID)
		} else if constant.Constant.Iotedge {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply, params, msgID)
		}
	}
}

func procGenOnOffProcRead(devEUI string, dstEndpointIndex int, clusterID uint16) {
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.ReadAttribute,
		ClusterID:   clusterID,
		Command: common.Command{
			DstEndpointIndex: dstEndpointIndex,
		},
	})
}

func genOnOffProcKeepAlive(devEUI string, tmnType string, interval uint16) {
	switch tmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
		keepalive.ProcKeepAlive(devEUI, interval)
		go procGenOnOffProcRead(devEUI, 0, 0x0702)
		go func() {
			timer := time.NewTimer(2 * time.Second)
			<-timer.C
			timer.Stop()
			procGenOnOffProcRead(devEUI, 0, 0x0b04)
		}()
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
		keepalive.ProcKeepAlive(devEUI, interval)
	default:
		globallogger.Log.Warnf("[genOnOffProcKeepAlive]: invalid tmnType: %s", tmnType)
	}
}

// genOnOffProcReport 处理report（0x0a）消息
func genOnOffProcReport(terminalInfo config.TerminalInfo, srcEndpoint uint8, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][genOnOffProcReport]: command: %+v", terminalInfo.DevEUI, Command.AttributeReports[0])
	genOnOffProcReadRspOrReport(terminalInfo, srcEndpoint, Command.AttributeReports[0].Attribute.Value.(bool), Command.AttributeReports[0].AttributeName)
	genOnOffProcKeepAlive(terminalInfo.DevEUI, terminalInfo.TmnType, uint16(terminalInfo.Interval))
}

// genOnOffProcDefaultResponse 处理defaultResponse（0x0b）消息
func genOnOffProcDefaultResponse(terminalInfo config.TerminalInfo, srcEndpoint uint8, command interface{}, msgID interface{}) {
	Command := command.(*cluster.DefaultResponseCommand)
	globallogger.Log.Infof("[devEUI: %v][genOnOffProcDefaultResponse]: command: %+v", terminalInfo.DevEUI, Command)
	code := 0
	message := "success"
	if Command.Status != cluster.ZclStatusSuccess {
		code = -1
		message = "failed"
	}
	params := make(map[string]interface{}, 4)
	params["code"] = code
	params["message"] = message
	params["terminalId"] = terminalInfo.DevEUI
	var key string
	var value string
	var bPublish = false
	endpoint := strconv.FormatUint(uint64(srcEndpoint-1), 10)
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
		key = iotsmartspace.HeimanESocketPropertyOnOff
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
		key = iotsmartspace.HeimanSmartPlugPropertyOnOff
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30:
		key = iotsmartspace.HeimanHS2SW1LEFR30PropertyOnOff1
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30:
		if endpoint == "0" {
			key = iotsmartspace.HeimanHS2SW2LEFR30PropertyOnOff1
			bPublish = true
		} else if endpoint == "1" {
			key = iotsmartspace.HeimanHS2SW2LEFR30PropertyOnOff2
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30:
		if endpoint == "0" {
			key = iotsmartspace.HeimanHS2SW3LEFR30PropertyOnOff1
			bPublish = true
		} else if endpoint == "1" {
			key = iotsmartspace.HeimanHS2SW3LEFR30PropertyOnOff2
			bPublish = true
		} else if endpoint == "2" {
			key = iotsmartspace.HeimanHS2SW3LEFR30PropertyOnOff3
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32:
		key = iotsmartspace.HonyarSingleSwitch00500c32PropertyOnOff1
		bPublish = true
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33:
		if endpoint == "0" {
			key = iotsmartspace.HonyarDoubleSwitch00500c33PropertyOnOff1
			bPublish = true
		} else if endpoint == "1" {
			key = iotsmartspace.HonyarDoubleSwitch00500c33PropertyOnOff2
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35:
		if endpoint == "0" {
			key = iotsmartspace.HonyarTripleSwitch00500c35PropertyOnOff1
			bPublish = true
		} else if endpoint == "1" {
			key = iotsmartspace.HonyarTripleSwitch00500c35PropertyOnOff2
			bPublish = true
		} else if endpoint == "2" {
			key = iotsmartspace.HonyarTripleSwitch00500c35PropertyOnOff3
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141:
		bPublish = true
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142:
		if endpoint == "0" {
			bPublish = true
		} else if endpoint == "1" {
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
		if endpoint == "0" {
			bPublish = true
		} else if endpoint == "1" {
			bPublish = true
		} else if endpoint == "2" {
			bPublish = true
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
		key = iotsmartspace.HonyarSocket000a0c3cPropertyOnOff
		if srcEndpoint == 2 {
			key = iotsmartspace.HonyarSocket000a0c3cPropertyUSB
		}
		bPublish = true
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
		key = iotsmartspace.HonyarSocketHY0106PropertyOnOff
		bPublish = true
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
		key = iotsmartspace.HonyarSocketHY0105PropertyOnOff
		if srcEndpoint == 2 {
			key = iotsmartspace.HonyarSocketHY0105PropertyUSB
		}
		bPublish = true
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
		key = iotsmartspace.HonyarSocketHY0106PropertyOnOff
		bPublish = true
	default:
		globallogger.Log.Warnf("[devEUI: %v][genOnOffProcDefaultResponse] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
	}
	switch Command.CommandID {
	case 0x00:
		value = "00"
	case 0x01:
		value = "01"
	default:
		globallogger.Log.Warnf("[devEUI: %v][genOnOffProcDefaultResponse] invalid command: %v", terminalInfo.DevEUI, Command.CommandID)
	}
	params[key] = value
	if bPublish {
		if constant.Constant.Iotware {
			iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, message, msgID)
		} else if constant.Constant.Iotedge {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply, params, msgID)
		}
	}
}

// genOnOffProcConfigureReportingResponse 处理configureReportingResponse（0x07）消息
func genOnOffProcConfigureReportingResponse(terminalInfo config.TerminalInfo, srcEndpoint uint8, command interface{}) {
	Command := command.(*cluster.ConfigureReportingResponse)
	for _, v := range Command.AttributeStatusRecords {
		if v.Status != cluster.ZclStatusSuccess {
			globallogger.Log.Infof("[devEUI: %v][genOnOffProcConfigureReportingResponse]: configReport failed: %x", terminalInfo.DevEUI, v.Status)
		} else {
			// proc keepalive
			if v.AttributeID == 0 {
				genOnOffProcKeepAlive(terminalInfo.DevEUI, terminalInfo.TmnType, uint16(terminalInfo.Interval))
			}
		}
	}
}

//GenOnOffProc 处理clusterID 0x0006即genOnOff属性消息
func GenOnOffProc(terminalInfo config.TerminalInfo, srcEndpoint uint8, zclFrame *zcl.Frame, msgID interface{}, contentFrame *zcl.Frame) {
	switch zclFrame.CommandName {
	case "ReadAttributesResponse":
		genOnOffProcReadRsp(terminalInfo, srcEndpoint, zclFrame.Command)
	case "WriteAttributesResponse":
		genOnOffProcWriteRsp(terminalInfo, srcEndpoint, zclFrame.Command, msgID, contentFrame)
	case "ConfigureReportingResponse":
		genOnOffProcConfigureReportingResponse(terminalInfo, srcEndpoint, zclFrame.Command)
	case "ReportAttributes":
		genOnOffProcReport(terminalInfo, srcEndpoint, zclFrame.Command)
	case "DefaultResponse":
		genOnOffProcDefaultResponse(terminalInfo, srcEndpoint, zclFrame.Command, msgID)
	default:
		globallogger.Log.Warnf("[devEUI: %v][GenOnOffProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
