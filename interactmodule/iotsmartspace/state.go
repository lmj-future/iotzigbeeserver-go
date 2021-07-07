package iotsmartspace

import (
	"encoding/json"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	uuid "github.com/satori/go.uuid"
)

// StateTerminalOnline StateTerminalOnline
func StateTerminalOnline(devEUI string) {
	params := make(map[string]interface{}, 1)
	params["terminalId"] = devEUI
	Publish(TopicZigbeeserverIotsmartspaceState, MethodStateUp, params, uuid.NewV4().String())
}

// StateTerminalOffline StateTerminalOffline
func StateTerminalOffline(devEUI string) {
	params := make(map[string]interface{}, 1)
	params["terminalId"] = devEUI
	Publish(TopicZigbeeserverIotsmartspaceState, MethodStateDown, params, uuid.NewV4().String())
}

// StateTerminalLeave StateTerminalLeave
func StateTerminalLeave(devEUI string) {
	params := make(map[string]interface{}, 1)
	params["terminalId"] = devEUI
	Publish(TopicZigbeeserverIotsmartspaceState, MethodStateLeave, params, uuid.NewV4().String())
}

type gateway struct {
	FirstAddr  string `json:"firstAddr"`
	SecondAddr string `json:"secondAddr"`
	ThirdAddr  string `json:"thirdAddr"`
	PortID     string `json:"portID"`
}
type iotwareConnect struct {
	Device  string  `json:"device"`
	Type    string  `json:"type"`
	Gateway gateway `json:"gateway"`
}

// StateTerminalOnlineIotware StateTerminalOnlineIotware
func StateTerminalOnlineIotware(terminalInfo config.TerminalInfo) {
	jsonMsg, err := json.Marshal(iotwareConnect{
		Device: terminalInfo.DevEUI,
		Type:   terminalInfo.ManufacturerName + "-" + terminalInfo.TmnType,
		Gateway: gateway{
			FirstAddr:  terminalInfo.FirstAddr,
			SecondAddr: terminalInfo.SecondAddr,
			ThirdAddr:  terminalInfo.ThirdAddr,
			PortID:     terminalInfo.ModuleID,
		},
	})
	if err != nil {
		globallogger.Log.Errorf("[StateTerminalOnlineIotware]: JSON marshaling failed: %s", err)
	}
	PublishIotware(TopicV1GatewayConnect, jsonMsg)
}

// StateTerminalOfflineIotware StateTerminalOfflineIotware
func StateTerminalOfflineIotware(terminalInfo config.TerminalInfo) {
	jsonMsg, err := json.Marshal(iotwareConnect{
		Device: terminalInfo.DevEUI,
		Type:   terminalInfo.ManufacturerName + "-" + terminalInfo.TmnType,
		Gateway: gateway{
			FirstAddr:  terminalInfo.FirstAddr,
			SecondAddr: terminalInfo.SecondAddr,
			ThirdAddr:  terminalInfo.ThirdAddr,
			PortID:     terminalInfo.ModuleID,
		},
	})
	if err != nil {
		globallogger.Log.Errorf("[StateTerminalOfflineIotware]: JSON marshaling failed: %s", err)
	}
	PublishIotware(TopicV1GatewayDisconnect, jsonMsg)
}

// StateTerminalJoinIotware StateTerminalJoinIotware
func StateTerminalJoinIotware(terminalInfo config.TerminalInfo) {
	jsonMsg, err := json.Marshal(iotwareConnect{
		Device: terminalInfo.DevEUI,
		Type:   terminalInfo.ManufacturerName + "-" + terminalInfo.TmnType,
		Gateway: gateway{
			FirstAddr:  terminalInfo.FirstAddr,
			SecondAddr: terminalInfo.SecondAddr,
			ThirdAddr:  terminalInfo.ThirdAddr,
			PortID:     terminalInfo.ModuleID,
		},
	})
	if err != nil {
		globallogger.Log.Errorf("[StateTerminalJoinIotware]: JSON marshaling failed: %s", err)
	}
	PublishIotware(TopicV1GatewayNetworkIn, jsonMsg)
}

// StateTerminalLeaveIotware StateTerminalLeaveIotware
func StateTerminalLeaveIotware(terminalInfo config.TerminalInfo) {
	jsonMsg, err := json.Marshal(iotwareConnect{
		Device: terminalInfo.DevEUI,
		Type:   terminalInfo.ManufacturerName + "-" + terminalInfo.TmnType,
		Gateway: gateway{
			FirstAddr:  terminalInfo.FirstAddr,
			SecondAddr: terminalInfo.SecondAddr,
			ThirdAddr:  terminalInfo.ThirdAddr,
			PortID:     terminalInfo.ModuleID,
		},
	})
	if err != nil {
		globallogger.Log.Errorf("[StateTerminalLeaveIotware]: JSON marshaling failed: %s", err)
	}
	PublishIotware(TopicV1GatewayNetworkOut, jsonMsg)
}

// procState 处理终端上下线消息
func procState(mqttMsg mqttMsgSt) {
	globallogger.Log.Infof("[procState]: mqttMsg.Params: %+v", mqttMsg.Params)
}
