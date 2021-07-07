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
func genOnOffProcKeepAlive(terminalInfo config.TerminalInfo, interval uint16) {
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
		keepalive.ProcKeepAlive(terminalInfo, interval)
		go procGenOnOffProcRead(terminalInfo.DevEUI, 0, 0x0702)
		go func() {
			timer := time.NewTimer(2 * time.Second)
			<-timer.C
			timer.Stop()
			procGenOnOffProcRead(terminalInfo.DevEUI, 0, 0x0b04)
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
		keepalive.ProcKeepAlive(terminalInfo, interval)
	default:
		globallogger.Log.Warnf("[genOnOffProcKeepAlive]: invalid tmnType: %s", terminalInfo.TmnType)
	}
}

func genOnOffProcReadRspOrReportIotware(terminalInfo config.TerminalInfo, srcEndpoint uint8, Value bool, attributeName string) {
	if (terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c ||
		terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55 ||
		terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105 ||
		terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106) &&
		(srcEndpoint == 2 || attributeName == "ChildLock" || attributeName == "BackGroundLight" || attributeName == "PowerOffMemory") {
		switch attributeName {
		case "OnOff":
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyUSBSwitch_1{USBSwitch_1: Value})
			}
		case "ChildLock":
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyChildLockSwitch{ChildLockSwitch: Value})
			}
		case "BackGroundLight":
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyBackLightSwitch{BackLightSwitch: Value})
			}
		case "PowerOffMemory":
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyMemorySwitch{MemorySwitch: !Value})
			}
		}
	} else {
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
			switch strconv.FormatUint(uint64(srcEndpoint-1), 10) {
			case "0":
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyPowerSwitch_1{PowerSwitch_1: Value})
			case "1":
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyPowerSwitch_2{PowerSwitch_2: Value})
			case "2":
				iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyPowerSwitch_3{PowerSwitch_3: Value})
			}
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket,
			constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, iotsmartspace.IotwarePropertyPowerSwitch{PowerSwitch: Value})
		}
	}
}
func genOnOffProcReadRspOrReportIotedge(terminalInfo config.TerminalInfo, srcEndpoint uint8, Value bool, attributeName string) {
	var value = "00"
	if Value {
		value = "01"
	}
	if (terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c ||
		terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55 ||
		terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105 ||
		terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106) &&
		(srcEndpoint == 2 || attributeName == "ChildLock" || attributeName == "BackGroundLight" || attributeName == "PowerOffMemory") {
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
			switch attributeName {
			case "OnOff":
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarSocket000a0c3cPropertyUSB{DevEUI: terminalInfo.DevEUI, USB: value}, uuid.NewV4().String())
			case "ChildLock":
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarSocket000a0c3cPropertyChildLock{DevEUI: terminalInfo.DevEUI, ChildLock: value}, uuid.NewV4().String())
			}
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
			switch attributeName {
			case "ChildLock":
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarSocketHY0106PropertyChildLock{DevEUI: terminalInfo.DevEUI, ChildLock: value}, uuid.NewV4().String())
			}
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
			switch attributeName {
			case "OnOff":
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarSocketHY0105PropertyUSB{DevEUI: terminalInfo.DevEUI, USB: value}, uuid.NewV4().String())
			case "ChildLock":
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarSocketHY0105PropertyChildLock{DevEUI: terminalInfo.DevEUI, ChildLock: value}, uuid.NewV4().String())
			case "BackGroundLight":
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarSocketHY0105PropertyBackGroundLight{DevEUI: terminalInfo.DevEUI, BackGroundLight: value}, uuid.NewV4().String())
			case "PowerOffMemory":
				if Value {
					value = "00"
				} else {
					value = "01"
				}
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarSocketHY0105PropertyPowerOffMemory{DevEUI: terminalInfo.DevEUI, PowerOffMemory: value}, uuid.NewV4().String())
			}
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
			switch attributeName {
			case "ChildLock":
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarSocketHY0106PropertyChildLock{DevEUI: terminalInfo.DevEUI, ChildLock: value}, uuid.NewV4().String())
			case "BackGroundLight":
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarSocketHY0106PropertyBackGroundLight{DevEUI: terminalInfo.DevEUI, BackGroundLight: value}, uuid.NewV4().String())
			case "PowerOffMemory":
				if Value {
					value = "00"
				} else {
					value = "01"
				}
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarSocketHY0106PropertyPowerOffMemory{DevEUI: terminalInfo.DevEUI, PowerOffMemory: value}, uuid.NewV4().String())
			}
		}
	} else {
		endpoint := strconv.FormatUint(uint64(srcEndpoint-1), 10)
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
				iotsmartspace.HeimanESocketPropertyOnOff{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
				iotsmartspace.HeimanSmartPlugPropertyOnOff{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30:
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
				iotsmartspace.HeimanHS2SW1LEFR30PropertyOnOff1{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30:
			if endpoint == "0" {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2SW2LEFR30PropertyOnOff1{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
			} else if endpoint == "1" {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2SW2LEFR30PropertyOnOff2{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
			}
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30:
			if endpoint == "0" {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2SW3LEFR30PropertyOnOff1{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
			} else if endpoint == "1" {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2SW3LEFR30PropertyOnOff2{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
			} else if endpoint == "2" {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HeimanHS2SW3LEFR30PropertyOnOff3{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
			}
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32:
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
				iotsmartspace.HonyarSingleSwitch00500c32PropertyOnOff1{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33:
			if endpoint == "0" {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarDoubleSwitch00500c33PropertyOnOff1{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
			} else if endpoint == "1" {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarDoubleSwitch00500c33PropertyOnOff2{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
			}
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35:
			if endpoint == "0" {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarTripleSwitch00500c35PropertyOnOff1{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
			} else if endpoint == "1" {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarTripleSwitch00500c35PropertyOnOff2{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
			} else if endpoint == "2" {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
					iotsmartspace.HonyarTripleSwitch00500c35PropertyOnOff3{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
			}
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141:
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142:
			// if endpoint == "0" {
			// } else if endpoint == "1" {
			// }
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
			// if endpoint == "0" {
			// } else if endpoint == "1" {
			// } else if endpoint == "2" {
			// }
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
				iotsmartspace.HonyarSocket000a0c3cPropertyOnOff{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
				iotsmartspace.HonyarSocketHY0106PropertyOnOff{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
				iotsmartspace.HonyarSocketHY0105PropertyOnOff{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp,
				iotsmartspace.HonyarSocketHY0106PropertyOnOff{DevEUI: terminalInfo.DevEUI, OnOff: value}, uuid.NewV4().String())
		default:
			globallogger.Log.Warnf("[devEUI: %v][genOnOffProcReadRspOrReportIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
		}
	}
}
func genOnOffProcReadRspOrReportIotprivate(terminalInfo config.TerminalInfo, srcEndpoint uint8, Value bool, attributeName string) {
	var appData string
	if (terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c ||
		terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55 ||
		terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105 ||
		terminalInfo.TmnType == constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106) &&
		(srcEndpoint == 2 || attributeName == "ChildLock" || attributeName == "BackGroundLight" || attributeName == "PowerOffMemory") {
		if Value {
			appData = "开"
			if attributeName == "PowerOffMemory" {
				appData = "关"
			}
		} else {
			appData = "关"
			if attributeName == "PowerOffMemory" {
				appData = "开"
			}
		}
		kafkaMsg := publicstruct.DataReportMsg{
			Time:       time.Now(),
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
	} else {
		// Value == true: ON; Value == false: OFF
		if Value {
			appData = "开"
		} else {
			appData = "关"
		}
		kafkaMsg := publicstruct.DataReportMsg{
			Time:       time.Now(),
			OIDIndex:   terminalInfo.OIDIndex,
			DevSN:      terminalInfo.DevEUI,
			LinkType:   terminalInfo.ProfileID,
			DeviceType: terminalInfo.TmnType2,
		}
		switch terminalInfo.TmnType {
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30, constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
			switch strconv.FormatUint(uint64(srcEndpoint-1), 10) {
			case "0":
				type appDataMsg struct {
					OnOffOneState string `json:"onOffOneState"`
				}
				kafkaMsg.AppData = appDataMsg{OnOffOneState: "一键" + appData}
			case "1":
				type appDataMsg struct {
					OnOffTwoState string `json:"onOffTwoState"`
				}
				kafkaMsg.AppData = appDataMsg{OnOffTwoState: "二键" + appData}
			case "2":
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
func genOnOffProcReadRspOrReport(terminalInfo config.TerminalInfo, srcEndpoint uint8, Value bool, attributeName string) {
	if constant.Constant.Iotware {
		genOnOffProcReadRspOrReportIotware(terminalInfo, srcEndpoint, Value, attributeName)
	} else if constant.Constant.Iotedge {
		genOnOffProcReadRspOrReportIotedge(terminalInfo, srcEndpoint, Value, attributeName)
	} else if constant.Constant.Iotprivate {
		genOnOffProcReadRspOrReportIotprivate(terminalInfo, srcEndpoint, Value, attributeName)
	}
}

//genOnOffProcReadRsp 处理readRsp（0x01）消息
func genOnOffProcReadRsp(terminalInfo config.TerminalInfo, srcEndpoint uint8, command interface{}) {
	for _, v := range command.(*cluster.ReadAttributesResponse).ReadAttributeStatuses {
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

func genOnOffProcWriteRspIotware(terminalInfo config.TerminalInfo, command interface{}, msgID interface{}, contentFrame *zcl.Frame) {
	Command := command.(*cluster.WriteAttributesResponse)
	globallogger.Log.Infof("[devEUI: %v][genOnOffProcWriteRspIotware]: command: %+v", terminalInfo.DevEUI, Command)
	message := "success"
	if Command.WriteAttributeStatuses[0].Status != cluster.ZclStatusSuccess {
		message = "failed"
	}
	if contentFrame != nil && contentFrame.CommandName == "WriteAttributes" {
		switch contentFrame.Command.(*cluster.WriteAttributesCommand).WriteAttributeRecords[0].AttributeName {
		case "ChildLock":
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
				iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, message, msgID)
			}
		case "BackGroundLight":
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
				iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, message, msgID)
			}
		case "PowerOffMemory":
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
				iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, message, msgID)
			}
		}
	}
}
func genOnOffProcWriteRspIotedge(terminalInfo config.TerminalInfo, command interface{}, msgID interface{}, contentFrame *zcl.Frame) {
	Command := command.(*cluster.WriteAttributesResponse)
	globallogger.Log.Infof("[devEUI: %v][genOnOffProcWriteRsp]: command: %+v", terminalInfo.DevEUI, Command)
	code := 0
	message := "success"
	if Command.WriteAttributeStatuses[0].Status != cluster.ZclStatusSuccess {
		code = -1
		message = "failed"
	}
	var value string
	if contentFrame != nil && contentFrame.CommandName == "WriteAttributes" {
		contentCommand := contentFrame.Command.(*cluster.WriteAttributesCommand)
		if contentCommand.WriteAttributeRecords[0].Attribute.Value.(bool) {
			value = "01"
		} else {
			value = "00"
		}
		switch contentCommand.WriteAttributeRecords[0].AttributeName {
		case "ChildLock":
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
					iotsmartspace.HonyarSocket000a0c3cPropertyChildLockWriteRsp{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, ChildLock: value}, msgID)
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
					iotsmartspace.HonyarSocketHY0106PropertyChildLockWriteRsp{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, ChildLock: value}, msgID)
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
					iotsmartspace.HonyarSocketHY0105PropertyChildLockWriteRsp{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, ChildLock: value}, msgID)
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
					iotsmartspace.HonyarSocketHY0106PropertyChildLockWriteRsp{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, ChildLock: value}, msgID)
			}
		case "BackGroundLight":
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
					iotsmartspace.HonyarSocketHY0105PropertyBackGroundLightWriteRsp{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, BackGroundLight: value}, msgID)
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
					iotsmartspace.HonyarSocketHY0106PropertyBackGroundLightWriteRsp{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, BackGroundLight: value}, msgID)
			}
		case "PowerOffMemory":
			switch terminalInfo.TmnType {
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
				if value == "00" {
					value = "01"
				} else if value == "01" {
					value = "00"
				}
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
					iotsmartspace.HonyarSocketHY0105PropertyPowerOffMemoryWriteRsp{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, PowerOffMemory: value}, msgID)
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
				if value == "00" {
					value = "01"
				} else if value == "01" {
					value = "00"
				}
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
					iotsmartspace.HonyarSocketHY0106PropertyPowerOffMemoryWriteRsp{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, PowerOffMemory: value}, msgID)
			}
		}
	}
}

// genOnOffProcWriteRsp 处理writeRsp（0x04）消息
func genOnOffProcWriteRsp(terminalInfo config.TerminalInfo, command interface{}, msgID interface{}, contentFrame *zcl.Frame) {
	if constant.Constant.Iotware {
		genOnOffProcWriteRspIotware(terminalInfo, command, msgID, contentFrame)
	} else if constant.Constant.Iotedge {
		genOnOffProcWriteRspIotedge(terminalInfo, command, msgID, contentFrame)
	}
}

// genOnOffProcConfigureReportingResponse 处理configureReportingResponse（0x07）消息
func genOnOffProcConfigureReportingResponse(terminalInfo config.TerminalInfo, srcEndpoint uint8, command interface{}) {
	for _, v := range command.(*cluster.ConfigureReportingResponse).AttributeStatusRecords {
		if v.Status != cluster.ZclStatusSuccess {
			globallogger.Log.Infof("[devEUI: %v][genOnOffProcConfigureReportingResponse]: configReport failed: %x", terminalInfo.DevEUI, v.Status)
		} else {
			// proc keepalive
			if v.AttributeID == 0 {
				genOnOffProcKeepAlive(terminalInfo, uint16(terminalInfo.Interval))
			}
		}
	}
}

// genOnOffProcReport 处理report（0x0a）消息
func genOnOffProcReport(terminalInfo config.TerminalInfo, srcEndpoint uint8, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][genOnOffProcReport]: command: %+v", terminalInfo.DevEUI, Command.AttributeReports[0])
	genOnOffProcReadRspOrReport(terminalInfo, srcEndpoint, Command.AttributeReports[0].Attribute.Value.(bool), Command.AttributeReports[0].AttributeName)
	genOnOffProcKeepAlive(terminalInfo, uint16(terminalInfo.Interval))
}

func genOnOffProcDefaultResponseIotware(terminalInfo config.TerminalInfo, command interface{}, msgID interface{}) {
	Command := command.(*cluster.DefaultResponseCommand)
	globallogger.Log.Infof("[devEUI: %v][genOnOffProcDefaultResponseIotware]: command: %+v", terminalInfo.DevEUI, Command)
	message := "success"
	if Command.Status != cluster.ZclStatusSuccess {
		message = "failed"
	}
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
		iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, message, msgID)
	default:
		globallogger.Log.Warnf("[devEUI: %v][genOnOffProcDefaultResponseIotware] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
	}
}
func genOnOffProcDefaultResponseIotedge(terminalInfo config.TerminalInfo, srcEndpoint uint8, command interface{}, msgID interface{}) {
	Command := command.(*cluster.DefaultResponseCommand)
	globallogger.Log.Infof("[devEUI: %v][genOnOffProcDefaultResponseIotedge]: command: %+v", terminalInfo.DevEUI, Command)
	code := 0
	message := "success"
	if Command.Status != cluster.ZclStatusSuccess {
		code = -1
		message = "failed"
	}
	var value string
	endpoint := strconv.FormatUint(uint64(srcEndpoint-1), 10)
	switch Command.CommandID {
	case 0x00:
		value = "00"
	case 0x01:
		value = "01"
	default:
		globallogger.Log.Warnf("[devEUI: %v][genOnOffProcDefaultResponseIotedge] invalid command: %v", terminalInfo.DevEUI, Command.CommandID)
	}
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
			iotsmartspace.HeimanESocketPropertyOnOffDefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
			iotsmartspace.HeimanSmartPlugPropertyOnOffDefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30:
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
			iotsmartspace.HeimanHS2SW1LEFR30PropertyOnOff1DefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30:
		if endpoint == "0" {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
				iotsmartspace.HeimanHS2SW2LEFR30PropertyOnOff1DefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
		} else if endpoint == "1" {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
				iotsmartspace.HeimanHS2SW2LEFR30PropertyOnOff2DefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
		}
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30:
		if endpoint == "0" {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
				iotsmartspace.HeimanHS2SW3LEFR30PropertyOnOff1DefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
		} else if endpoint == "1" {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
				iotsmartspace.HeimanHS2SW3LEFR30PropertyOnOff2DefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
		} else if endpoint == "2" {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
				iotsmartspace.HeimanHS2SW3LEFR30PropertyOnOff3DefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32:
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
			iotsmartspace.HonyarSingleSwitch00500c32PropertyOnOff1DefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33:
		if endpoint == "0" {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
				iotsmartspace.HonyarDoubleSwitch00500c33PropertyOnOff1DefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
		} else if endpoint == "1" {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
				iotsmartspace.HonyarDoubleSwitch00500c33PropertyOnOff2DefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35:
		if endpoint == "0" {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
				iotsmartspace.HonyarTripleSwitch00500c35PropertyOnOff1DefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
		} else if endpoint == "1" {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
				iotsmartspace.HonyarTripleSwitch00500c35PropertyOnOff2DefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
		} else if endpoint == "2" {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
				iotsmartspace.HonyarTripleSwitch00500c35PropertyOnOff3DefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141:
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142:
		// if endpoint == "0" {
		// } else if endpoint == "1" {
		// }
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
		// if endpoint == "0" {
		// } else if endpoint == "1" {
		// } else if endpoint == "2" {
		// }
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c:
		if srcEndpoint == 2 {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
				iotsmartspace.HonyarSocket000a0c3cPropertyUSBDefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, USB: value}, msgID)
		} else {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
				iotsmartspace.HonyarSocket000a0c3cPropertyOnOffDefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55:
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
			iotsmartspace.HonyarSocketHY0106PropertyOnOffDefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
		if srcEndpoint == 2 {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
				iotsmartspace.HonyarSocketHY0105PropertyUSBDefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, USB: value}, msgID)
		} else {
			iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
				iotsmartspace.HonyarSocketHY0105PropertyOnOffDefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
		}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
			iotsmartspace.HonyarSocketHY0106PropertyOnOffDefaultResponse{Code: code, Message: message, DevEUI: terminalInfo.DevEUI, OnOff: value}, msgID)
	default:
		globallogger.Log.Warnf("[devEUI: %v][genOnOffProcDefaultResponseIotedge] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
	}
}

// genOnOffProcDefaultResponse 处理defaultResponse（0x0b）消息
func genOnOffProcDefaultResponse(terminalInfo config.TerminalInfo, srcEndpoint uint8, command interface{}, msgID interface{}) {
	if constant.Constant.Iotware {
		genOnOffProcDefaultResponseIotware(terminalInfo, command, msgID)
	} else if constant.Constant.Iotedge {
		genOnOffProcDefaultResponseIotedge(terminalInfo, srcEndpoint, command, msgID)
	}
}

//GenOnOffProc 处理clusterID 0x0006即genOnOff属性消息
func GenOnOffProc(terminalInfo config.TerminalInfo, srcEndpoint uint8, zclFrame *zcl.Frame, msgID interface{}, contentFrame *zcl.Frame) {
	switch zclFrame.CommandName {
	case "ReadAttributesResponse":
		genOnOffProcReadRsp(terminalInfo, srcEndpoint, zclFrame.Command)
	case "WriteAttributesResponse":
		genOnOffProcWriteRsp(terminalInfo, zclFrame.Command, msgID, contentFrame)
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
