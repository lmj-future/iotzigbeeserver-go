package heiman

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
)

// infraredRemoteProcStudyKeyResponse 处理StudyKeyResponse（0xf2）消息
func infraredRemoteProcStudyKeyResponse(terminalInfo config.TerminalInfo, command interface{}, msgID interface{}) {
	Command := command.(*cluster.HEIMANInfraredRemoteStudyKeyResponse)
	globallogger.Log.Infof("[devEUI: %v][infraredRemoteProcStudyKeyResponse]: command: %+v", terminalInfo.DevEUI, Command)
	code := 0
	message := "success"
	if Command.ResultStatus == 0 {
		code = -1
		message = "failed"
	}
	if constant.Constant.Iotware {
		iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, message, msgID)
	} else if constant.Constant.Iotedge {
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
			iotsmartspace.HeimanIRControlEMPropertyStudyKeyRsp{
				DevEUI:  terminalInfo.DevEUI,
				Code:    code,
				Message: message,
				StudyKey: cluster.HEIMANInfraredRemoteStudyKey{
					ID:      Command.ID,
					KeyCode: Command.KeyCode,
				},
			}, msgID)
	} else if constant.Constant.Iotprivate {
		type appDataMsg struct {
			StudyKeyResponse cluster.HEIMANInfraredRemoteStudyKeyResponse `json:"studyKeyResponse"`
		}
		kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
			Time:       time.Now(),
			OIDIndex:   terminalInfo.OIDIndex,
			DevSN:      terminalInfo.DevEUI,
			LinkType:   terminalInfo.ProfileID,
			DeviceType: terminalInfo.TmnType2,
			AppData: appDataMsg{
				StudyKeyResponse: cluster.HEIMANInfraredRemoteStudyKeyResponse{
					ID:           Command.ID,
					KeyCode:      Command.KeyCode,
					ResultStatus: Command.ResultStatus,
				},
			},
		})
		globallogger.Log.Infof("[devEUI: %v][infraredRemoteProcStudyKeyResponse]: kafkaMsg: %s", terminalInfo.DevEUI, string(kafkaMsgByte))
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
}

// infraredRemoteProcCreateIDResponse 处理CreateIDResponse（0xf5）消息
func infraredRemoteProcCreateIDResponse(terminalInfo config.TerminalInfo, command interface{}, msgID interface{}) {
	Command := command.(*cluster.HEIMANInfraredRemoteCreateIDResponse)
	globallogger.Log.Infof("[devEUI: %v][infraredRemoteProcCreateIDResponse]: command: %+v", terminalInfo.DevEUI, Command)
	code := 0
	message := "success"
	if Command.ID == 0xff {
		code = -1
		message = "failed"
	}
	if constant.Constant.Iotware {
		iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, message, msgID)
	} else if constant.Constant.Iotedge {
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
			iotsmartspace.HeimanIRControlEMPropertyCreateIDRsp{
				DevEUI:  terminalInfo.DevEUI,
				Code:    code,
				Message: message,
				CreateID: cluster.HEIMANInfraredRemoteCreateIDResponse{
					ID:        Command.ID,
					ModelType: Command.ModelType,
				},
			}, msgID)
	} else if constant.Constant.Iotprivate {
		type appDataMsg struct {
			CreateIDResponse cluster.HEIMANInfraredRemoteCreateIDResponse `json:"createIDResponse"`
		}
		kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
			Time:       time.Now(),
			OIDIndex:   terminalInfo.OIDIndex,
			DevSN:      terminalInfo.DevEUI,
			LinkType:   terminalInfo.ProfileID,
			DeviceType: terminalInfo.TmnType2,
			AppData: appDataMsg{
				CreateIDResponse: cluster.HEIMANInfraredRemoteCreateIDResponse{
					ID:        Command.ID,
					ModelType: Command.ModelType,
				},
			},
		})
		globallogger.Log.Infof("[devEUI: %v][infraredRemoteProcCreateIDResponse]: kafkaMsg: %s", terminalInfo.DevEUI, string(kafkaMsgByte))
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
}

// infraredRemoteProcGetIDAndKeyCodeListResponse 处理GetIDAndKeyCodeListResponse（0xf6）消息
func infraredRemoteProcGetIDAndKeyCodeListResponse(terminalInfo config.TerminalInfo, command interface{}, msgID interface{}) {
	Command := command.(*cluster.HEIMANInfraredRemoteGetIDAndKeyCodeListResponse)
	globallogger.Log.Infof("[devEUI: %v][infraredRemoteProcGetIDAndKeyCodeListResponse]: command: %+v", terminalInfo.DevEUI, Command)
	value := iotsmartspace.HeimanIRControlEMPropertyGetIDAndKeyCodeListResponse{
		PacketNumSum:        Command.PacketNumSum,
		CurrentPacketNum:    Command.CurrentPacketNum,
		CurrentPacketLength: Command.CurrentPacketLength,
		IDKeyCodeList:       []iotsmartspace.HeimanIRControlEMPropertyIdKeyCode{},
	}
	for len(Command.IDKeyCodeList) > 0 {
		IDKeyCode := iotsmartspace.HeimanIRControlEMPropertyIdKeyCode{
			ID:        Command.IDKeyCodeList[0],
			ModelType: Command.IDKeyCodeList[1],
			KeyNum:    Command.IDKeyCodeList[2],
			KeyCode:   Command.IDKeyCodeList[3 : 3+int(Command.IDKeyCodeList[2])],
		}
		value.IDKeyCodeList = append(value.IDKeyCodeList, IDKeyCode)
		if len(Command.IDKeyCodeList) > 3+int(Command.IDKeyCodeList[2]) {
			Command.IDKeyCodeList = Command.IDKeyCodeList[3+int(Command.IDKeyCodeList[2]):]
		} else {
			Command.IDKeyCodeList = []uint8{}
		}
	}
	globallogger.Log.Infof("[devEUI: %v][infraredRemoteProcGetIDAndKeyCodeListResponse]: value: %+v", terminalInfo.DevEUI, value)
	if constant.Constant.Iotware {
		result, _ := json.Marshal(value)
		iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, string(result), msgID)
	} else if constant.Constant.Iotedge {
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply,
			iotsmartspace.HeimanIRControlEMPropertyGetIDAndKeyCodeListRsp{
				DevEUI:              terminalInfo.DevEUI,
				Code:                0,
				Message:             "success",
				GetIDAndKeyCodeList: value,
			}, msgID)
	} else if constant.Constant.Iotprivate {
		kafkaMsgByte, _ := json.Marshal(publicstruct.DataReportMsg{
			Time:       time.Now(),
			OIDIndex:   terminalInfo.OIDIndex,
			DevSN:      terminalInfo.DevEUI,
			LinkType:   terminalInfo.ProfileID,
			DeviceType: terminalInfo.TmnType2,
			AppData:    value,
		})
		globallogger.Log.Infof("[devEUI: %v][infraredRemoteProcGetIDAndKeyCodeListResponse]: kafkaMsg: %s", terminalInfo.DevEUI, string(kafkaMsgByte))
		kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
	}
}

// InfraredRemoteProc 处理clusterID 0xfc82属性消息
func InfraredRemoteProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame, msgID interface{}, contentFrame *zcl.Frame) {
	switch zclFrame.CommandName {
	case "StudyKeyResponse":
		infraredRemoteProcStudyKeyResponse(terminalInfo, zclFrame.Command, msgID)
	case "CreateIDResponse":
		infraredRemoteProcCreateIDResponse(terminalInfo, zclFrame.Command, msgID)
	case "GetIDAndKeyCodeListResponse":
		infraredRemoteProcGetIDAndKeyCodeListResponse(terminalInfo, zclFrame.Command, msgID)
	default:
		globallogger.Log.Warnf("[devEUI: %v][InfraredRemoteProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
