package general

import (
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
)

//genScenesProcReadRsp 处理readRsp（0x01）消息
func genScenesProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	readAttributesRsp := command.(*cluster.ReadAttributesResponse)
	for _, v := range readAttributesRsp.ReadAttributeStatuses {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcReadRsp]: readAttributesRsp: %+v", terminalInfo.DevEUI, v)
		if v.Status == cluster.ZclStatusSuccess {
			globallogger.Log.Infof("[devEUI: %v][genScenesProcReadRsp]: read rsp success: status %v", terminalInfo.DevEUI, v.Status)
			switch v.AttributeName {
			case "SceneCount":
				globallogger.Log.Infof("[devEUI: %v][genScenesProcReadRsp]: SceneCount:  %v", terminalInfo.DevEUI, v.Attribute.Value.(uint64))
			case "CurrentScene":
				globallogger.Log.Infof("[devEUI: %v][genScenesProcReadRsp]: CurrentScene:  %v", terminalInfo.DevEUI, v.Attribute.Value.(uint64))
			case "CurrentGroup":
				globallogger.Log.Infof("[devEUI: %v][genScenesProcReadRsp]: CurrentGroup:  %v", terminalInfo.DevEUI, v.Attribute.Value.(uint64))
			case "SceneValid":
				globallogger.Log.Infof("[devEUI: %v][genScenesProcReadRsp]: SceneValid:  %v", terminalInfo.DevEUI, v.Attribute.Value.(bool))
			case "NameSupport":
				globallogger.Log.Infof("[devEUI: %v][genScenesProcReadRsp]: NameSupport:  %v", terminalInfo.DevEUI, v.Attribute.Value)
			case "LastConfiguredBy":
				globallogger.Log.Infof("[devEUI: %v][genScenesProcReadRsp]: LastConfiguredBy:  %v", terminalInfo.DevEUI, v.Attribute.Value)
			default:
				globallogger.Log.Warnln("devEUI : "+terminalInfo.DevEUI+" "+"genScenesProcReadRsp unknow attributeName", v.AttributeName)
			}
		} else {
			globallogger.Log.Infof("[devEUI: %v][genScenesProcReadRsp]: read rsp failed: status %v", terminalInfo.DevEUI, v.Status)
		}
	}
}

// genScenesProcReport 处理report（0x0a）消息
func genScenesProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][genScenesProcReport]: command: %+v", terminalInfo.DevEUI, Command.AttributeReports[0])
}

// genScenesProcDefaultResponse 处理defaultResponse（0x0b）消息
func genScenesProcDefaultResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.DefaultResponseCommand)
	globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse]: command: %+v", terminalInfo.DevEUI, Command)
	if Command.Status == cluster.ZclStatusSuccess {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse]: success", terminalInfo.DevEUI)
	} else {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse]: failed, status: %+v", terminalInfo.DevEUI, Command.Status)
	}
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal1SceneSwitch005f0cf1:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][1SceneSwitch005f0cf1]", terminalInfo.DevEUI)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][2SceneSwitch005f0cf3]", terminalInfo.DevEUI)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][3SceneSwitch005f0cf2]", terminalInfo.DevEUI)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][6SceneSwitch005f0c3b]", terminalInfo.DevEUI)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][SingleSwitch00500c32]", terminalInfo.DevEUI)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][DoubleSwitch00500c33]", terminalInfo.DevEUI)
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][TripleSwitch00500c35]", terminalInfo.DevEUI)
	default:
		globallogger.Log.Warnf("[devEUI: %v][genScenesProcDefaultResponse] invalid tmnType: %v", terminalInfo.DevEUI, terminalInfo.TmnType)
	}

	switch Command.CommandID {
	case 0x00:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][AddSceneResponse]", terminalInfo.DevEUI)
	case 0x01:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][ViewSceneResponse]", terminalInfo.DevEUI)
	case 0x02:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][RemoveSceneResponse]", terminalInfo.DevEUI)
	case 0x03:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][RemoveAllScenesResponse]", terminalInfo.DevEUI)
	case 0x04:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][StoreSceneResponse]", terminalInfo.DevEUI)
	case 0x06:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][GetSceneMembershipResponse]", terminalInfo.DevEUI)
	case 0x40:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][EnhancedAddSceneResponse]", terminalInfo.DevEUI)
	case 0x41:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][EnhancedViewSceneResponse]", terminalInfo.DevEUI)
	case 0x42:
		globallogger.Log.Infof("[devEUI: %v][genScenesProcDefaultResponse][CopySceneResponse]", terminalInfo.DevEUI)
	default:
		globallogger.Log.Warnf("[devEUI: %v][genScenesProcDefaultResponse] invalid command: %v", terminalInfo.DevEUI, Command.CommandID)
	}
}

// genScenesProcAddSceneResponse 处理AddSceneResponse command（0x00）消息
func genScenesProcAddSceneResponse(terminalInfo config.TerminalInfo, command interface{}, msgID interface{}, contentFrame *zcl.Frame, dataTemp string) {
	Command := command.(*cluster.AddSceneResponse)
	globallogger.Log.Infof("[devEUI: %v][genScenesProcAddSceneResponse]: command: %+v, dataTemp: %s", terminalInfo.DevEUI, Command, dataTemp)
	code := 0
	message := "success"
	if cluster.ZclStatus(Command.Status) != cluster.ZclStatusSuccess {
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
	if contentFrame != nil && contentFrame.CommandName == "AddScene" {
		contentCommand := contentFrame.Command.(*cluster.AddSceneCommand)
		// switch contentCommand.SceneName[1:] {
		// case "\x01\x01":
		// 	value = "01"
		// case "\x02\x01":
		// 	value = "02"
		// case "\x03\x01":
		// 	value = "03"
		// case "\x04\x01":
		// 	value = "04"
		// case "\x05\x01":
		// 	value = "05"
		// case "\x06\x01":
		// 	value = "06"
		// case "\x07\x01":
		// 	value = "07"
		// case "\x08\x01":
		// 	value = "08"
		// case "\x09\x01":
		// 	value = "09"
		// case "\x0a\x01":
		// 	value = "0a"
		// case "\x0b\x01":
		// 	value = "0b"
		// case "\x0c\x01":
		// 	value = "0c"
		// case "\x0d\x01":
		// 	value = "0d"
		// case "\x0e\x01":
		// 	value = "0e"
		// case "\x0f\x01":
		// 	value = "0f"
		// case "\x10\x01":
		// 	value = "10"
		// case "\x11\x01":
		// 	value = "11"
		// case "\x12\x01":
		// 	value = "12"
		// case "\x13\x01":
		// 	value = "13"
		// case "\x14\x01":
		// 	value = "14"
		// case "\x15\x01":
		// 	value = "15"
		// case "\x16\x01":
		// 	value = "16"
		// case "\x17\x01":
		// 	value = "17"
		// case "\x18\x01":
		// 	value = "18"
		// case "\x19\x01":
		// 	value = "19"
		// case "\x1a\x01":
		// 	value = "1a"
		// case "\x1b\x01":
		// 	value = "1b"
		// case "\x1c\x01":
		// 	value = "1c"
		// case "\x1d\x01":
		// 	value = "1d"
		// case "\x1e\x01":
		// 	value = "1e"
		// case "\x1f\x01":
		// 	value = "1f"
		// case " \x01":
		// 	value = "20"
		// case "!\x01":
		// 	value = "21"
		// }
		if dataTemp != "" && len(dataTemp) > 26 {
			switch dataTemp[18:26] {
			case "D3E9C0D6":
				value = "娱乐"
			case "D0DDCFA2":
				value = "休息"
			case "C9CFB0E0":
				value = "上班"
			case "CFC2B0E0":
				value = "下班"
			case "B0ECB9AB":
				value = "办公"
			}
		}
		switch contentCommand.SceneID {
		case 1:
			key = iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyLeftUpAddScene
			bPublish = true
		case 2:
			key = iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyLeftMiddleAddScene
			bPublish = true
		case 3:
			key = iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyLeftDownAddScene
			bPublish = true
		case 4:
			key = iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyRightUpAddScene
			bPublish = true
		case 5:
			key = iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyRightMiddleAddScene
			bPublish = true
		case 6:
			key = iotsmartspace.Honyar6SceneSwitch005f0c3bPropertyRightDownAddScene
			bPublish = true
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

// genScenesProcViewSceneResponse 处理ViewSceneResponse command（0x01）消息
func genScenesProcViewSceneResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ViewSceneResponse)
	globallogger.Log.Infof("[devEUI: %v][genScenesProcViewSceneResponse]: command: %+v", terminalInfo.DevEUI, Command)
	if cluster.ZclStatus(Command.Status) == cluster.ZclStatusSuccess {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcViewSceneResponse]: success, GroupID: %v, SceneID: %v, SceneName: %v, TransitionTime: %v",
			terminalInfo.DevEUI, Command.GroupID, Command.SceneID, Command.SceneName, Command.TransitionTime)
	} else {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcViewSceneResponse]: failed, status: %v", terminalInfo.DevEUI, Command.Status)
	}
}

// genScenesProcRemoveSceneResponse 处理RemoveSceneResponse command（0x02）消息
func genScenesProcRemoveSceneResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.RemoveSceneResponse)
	globallogger.Log.Infof("[devEUI: %v][genScenesProcRemoveSceneResponse]: command: %+v", terminalInfo.DevEUI, Command)
	if cluster.ZclStatus(Command.Status) == cluster.ZclStatusSuccess {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcRemoveSceneResponse]: success, GroupID: %v, SceneID: %v",
			terminalInfo.DevEUI, Command.GroupID, Command.SceneID)
	} else {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcRemoveSceneResponse]: failed, status: %v", terminalInfo.DevEUI, Command.Status)
	}
}

// genScenesProcRemoveAllScenesResponse 处理RemoveAllScenesResponse command（0x03）消息
func genScenesProcRemoveAllScenesResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.RemoveAllScenesResponse)
	globallogger.Log.Infof("[devEUI: %v][genScenesProcRemoveAllScenesResponse]: command: %+v", terminalInfo.DevEUI, Command)
	if cluster.ZclStatus(Command.Status) == cluster.ZclStatusSuccess {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcRemoveAllScenesResponse]: success, GroupID: %v",
			terminalInfo.DevEUI, Command.GroupID)
	} else {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcRemoveAllScenesResponse]: failed, status: %v", terminalInfo.DevEUI, Command.Status)
	}
}

// genScenesProcStoreSceneResponse 处理StoreSceneResponse command（0x04）消息
func genScenesProcStoreSceneResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.StoreSceneResponse)
	globallogger.Log.Infof("[devEUI: %v][genScenesProcStoreSceneResponse]: command: %+v", terminalInfo.DevEUI, Command)
	if cluster.ZclStatus(Command.Status) == cluster.ZclStatusSuccess {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcStoreSceneResponse]: success, GroupID: %v, SceneID: %v",
			terminalInfo.DevEUI, Command.GroupID, Command.SceneID)
	} else {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcStoreSceneResponse]: failed, status: %v", terminalInfo.DevEUI, Command.Status)
	}
}

// genScenesProcGetSceneMembershipResponse 处理GetSceneMembershipResponse command（0x06）消息
func genScenesProcGetSceneMembershipResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.GetSceneMembershipResponse)
	globallogger.Log.Infof("[devEUI: %v][genScenesProcGetSceneMembershipResponse]: command: %+v", terminalInfo.DevEUI, Command)
	if cluster.ZclStatus(Command.Status) == cluster.ZclStatusSuccess {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcGetSceneMembershipResponse]: success, Capacity: %v, GroupID: %v, SceneCount: %v, SceneList: %+v",
			terminalInfo.DevEUI, Command.Capacity, Command.GroupID, Command.SceneCount, Command.SceneList)
	} else {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcGetSceneMembershipResponse]: failed, status: %v", terminalInfo.DevEUI, Command.Status)
	}
}

// genScenesProcEnhancedAddSceneResponse 处理EnhancedAddSceneResponse command（0x40）消息
func genScenesProcEnhancedAddSceneResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.EnhancedAddSceneResponse)
	globallogger.Log.Infof("[devEUI: %v][genScenesProcEnhancedAddSceneResponse]: command: %+v", terminalInfo.DevEUI, Command)
	if cluster.ZclStatus(Command.Status) == cluster.ZclStatusSuccess {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcEnhancedAddSceneResponse]: success, GroupID: %v, SceneID: %v",
			terminalInfo.DevEUI, Command.GroupID, Command.SceneID)
	} else {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcEnhancedAddSceneResponse]: failed, status: %v", terminalInfo.DevEUI, Command.Status)
	}
}

// genScenesProcEnhancedViewSceneResponse 处理EnhancedViewSceneResponse command（0x41）消息
func genScenesProcEnhancedViewSceneResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.EnhancedViewSceneResponse)
	globallogger.Log.Infof("[devEUI: %v][genScenesProcEnhancedViewSceneResponse]: command: %+v", terminalInfo.DevEUI, Command)
	if cluster.ZclStatus(Command.Status) == cluster.ZclStatusSuccess {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcEnhancedViewSceneResponse]: success, GroupID: %v, SceneID: %v, SceneName: %v, TransitionTime: %v",
			terminalInfo.DevEUI, Command.GroupID, Command.SceneID, Command.SceneName, Command.TransitionTime)
	} else {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcEnhancedViewSceneResponse]: failed, status: %v", terminalInfo.DevEUI, Command.Status)
	}
}

// genScenesProcCopySceneResponse 处理CopySceneResponse command（0x42）消息
func genScenesProcCopySceneResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.CopySceneResponse)
	globallogger.Log.Infof("[devEUI: %v][genScenesProcCopySceneResponse]: command: %+v", terminalInfo.DevEUI, Command)
	if cluster.ZclStatus(Command.Status) == cluster.ZclStatusSuccess {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcCopySceneResponse]: success, FromGroupID: %v, FromSceneID: %v",
			terminalInfo.DevEUI, Command.FromGroupID, Command.FromSceneID)
	} else {
		globallogger.Log.Infof("[devEUI: %v][genScenesProcCopySceneResponse]: failed, status: %v", terminalInfo.DevEUI, Command.Status)
	}
}

//GenScenesProc 处理clusterID 0x0005即genScenes属性消息
func GenScenesProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame, msgID interface{}, contentFrame *zcl.Frame, dataTemp string) {
	// globallogger.Log.Infof("[devEUI: %v][GenScenesProc] Start......", terminalInfo.DevEUI)
	// globallogger.Log.Infof("[devEUI: %v][GenScenesProc] zclFrame: %+v", terminalInfo.DevEUI, zclFrame.Command)
	z := zcl.New()
	switch zclFrame.CommandName {
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandReadAttributesResponse)].Name:
		genScenesProcReadRsp(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandReportAttributes)].Name:
		genScenesProcReport(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandDefaultResponse)].Name:
		genScenesProcDefaultResponse(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandConfigureReportingResponse)].Name:
	case "AddSceneResponse":
		genScenesProcAddSceneResponse(terminalInfo, zclFrame.Command, msgID, contentFrame, dataTemp)
	case "ViewSceneResponse":
		genScenesProcViewSceneResponse(terminalInfo, zclFrame.Command)
	case "RemoveSceneResponse":
		genScenesProcRemoveSceneResponse(terminalInfo, zclFrame.Command)
	case "RemoveAllScenesResponse":
		genScenesProcRemoveAllScenesResponse(terminalInfo, zclFrame.Command)
	case "StoreSceneResponse":
		genScenesProcStoreSceneResponse(terminalInfo, zclFrame.Command)
	case "GetSceneMembershipResponse":
		genScenesProcGetSceneMembershipResponse(terminalInfo, zclFrame.Command)
	case "EnhancedAddSceneResponse":
		genScenesProcEnhancedAddSceneResponse(terminalInfo, zclFrame.Command)
	case "EnhancedViewSceneResponse":
		genScenesProcEnhancedViewSceneResponse(terminalInfo, zclFrame.Command)
	case "CopySceneResponse":
		genScenesProcCopySceneResponse(terminalInfo, zclFrame.Command)
	default:
		globallogger.Log.Warnf("[devEUI: %v][GenScenesProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
