package heiman

import (
	"encoding/json"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
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
	params := make(map[string]interface{}, 3)
	params["code"] = code
	params["message"] = message
	params[iotsmartspace.HeimanIRControlEMPropertyStudyKey] = cluster.HEIMANInfraredRemoteStudyKey{
		ID:      Command.ID,
		KeyCode: Command.KeyCode,
	}
	if constant.Constant.Iotware {
		iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, message, msgID)
	} else if constant.Constant.Iotedge {
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply, params, msgID)
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
	params := make(map[string]interface{}, 3)
	params["code"] = code
	params["message"] = message
	params[iotsmartspace.HeimanIRControlEMPropertyCreateID] = cluster.HEIMANInfraredRemoteCreateIDResponse{
		ID:        Command.ID,
		ModelType: Command.ModelType,
	}
	if constant.Constant.Iotware {
		iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, message, msgID)
	} else if constant.Constant.Iotedge {
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply, params, msgID)
	}
}

// infraredRemoteProcGetIDAndKeyCodeListResponse 处理StudyKeyResponse（0xf2）消息
func infraredRemoteProcGetIDAndKeyCodeListResponse(terminalInfo config.TerminalInfo, command interface{}, msgID interface{}) {
	type idKeyCode struct {
		ID        uint8   `json:"ID"`
		ModelType uint8   `json:"ModelType"`
		KeyNum    uint8   `json:"KeyNum"`
		KeyCode   []uint8 `json:"KeyCode"`
	}
	type getIDAndKeyCodeListResponse struct {
		PacketNumSum        uint8       `json:"ID"`
		CurrentPacketNum    uint8       `json:"CurrentPacketNum"`
		CurrentPacketLength uint8       `json:"CurrentPacketLength"`
		IDKeyCodeList       []idKeyCode `json:"IDKeyCodeList"`
	}
	Command := command.(*cluster.HEIMANInfraredRemoteGetIDAndKeyCodeListResponse)
	globallogger.Log.Infof("[devEUI: %v][infraredRemoteProcGetIDAndKeyCodeListResponse]: command: %+v", terminalInfo.DevEUI, Command)
	value := getIDAndKeyCodeListResponse{
		PacketNumSum:        Command.PacketNumSum,
		CurrentPacketNum:    Command.CurrentPacketNum,
		CurrentPacketLength: Command.CurrentPacketLength,
	}
	for len(Command.IDKeyCodeList) > 0 {
		IDKeyCode := idKeyCode{
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
	code := 0
	message := "success"
	params := make(map[string]interface{}, 3)
	params["code"] = code
	params["message"] = message
	params[iotsmartspace.HeimanIRControlEMPropertyGetIDAndKeyCodeList] = value
	if constant.Constant.Iotware {
		result, _ := json.Marshal(value)
		iotsmartspace.PublishRPCRspIotware(terminalInfo.DevEUI, string(result), msgID)
	} else if constant.Constant.Iotedge {
		iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyDownReply, params, msgID)
	}
}

// InfraredRemoteProc 处理clusterID 0xfc82属性消息
func InfraredRemoteProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame, msgID interface{}, contentFrame *zcl.Frame) {
	// globallogger.Log.Infof("[devEUI: %v][InfraredRemoteProc] Start......", terminalInfo.DevEUI)
	// globallogger.Log.Infof("[devEUI: %v][InfraredRemoteProc] zclFrame: %+v", terminalInfo.DevEUI, zclFrame)
	z := zcl.New()
	switch zclFrame.CommandName {
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandDefaultResponse)].Name:
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
