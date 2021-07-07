package general

import (
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/keepalive"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
	"github.com/h3c/iotzigbeeserver-go/zcl/zclmsgdown"
)

//genBasicProcProcReadRsp 处理readRsp（0x01）消息
func genBasicProcProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	readAttributesRsp := command.(*cluster.ReadAttributesResponse)
	globallogger.Log.Infof("[devEUI: %v][genBasicProcProcReadRsp]: readAttributesRsp: %+v", terminalInfo.DevEUI, readAttributesRsp.ReadAttributeStatuses[0])
	type BasicInfo struct {
		ZLibraryVersion     uint8
		ApplicationVersion  uint8
		StackVersion        uint8
		HWVersion           uint8
		ManufacturerName    string
		ModelIdentifier     string
		DateCode            string
		PowerSource         uint8
		LocationDescription string
		PhysicalEnvironment uint8
		DeviceEnabled       bool
		AlarmMask           interface{}
		DisableLocalConfig  interface{}
		SWBuildID           string
	}
	basicInfo := BasicInfo{}
	for _, v := range readAttributesRsp.ReadAttributeStatuses {
		switch v.AttributeName {
		case "ZLibraryVersion":
			if v.Status == cluster.ZclStatusSuccess {
				basicInfo.ZLibraryVersion = uint8(v.Attribute.Value.(uint64))
			}
		case "ApplicationVersion":
			if v.Status == cluster.ZclStatusSuccess {
				basicInfo.ApplicationVersion = uint8(v.Attribute.Value.(uint64))
			}
		case "StackVersion":
			if v.Status == cluster.ZclStatusSuccess {
				basicInfo.StackVersion = uint8(v.Attribute.Value.(uint64))
			}
		case "HWVersion":
			if v.Status == cluster.ZclStatusSuccess {
				basicInfo.HWVersion = uint8(v.Attribute.Value.(uint64))
			}
		case "ManufacturerName":
			if v.Status == cluster.ZclStatusSuccess {
				basicInfo.ManufacturerName = v.Attribute.Value.(string)
			}
		case "ModelIdentifier":
			if v.Status == cluster.ZclStatusSuccess {
				basicInfo.ModelIdentifier = v.Attribute.Value.(string)
			}
		case "DateCode":
			if v.Status == cluster.ZclStatusSuccess {
				basicInfo.DateCode = v.Attribute.Value.(string)
			}
		case "PowerSource":
			if v.Status == cluster.ZclStatusSuccess {
				basicInfo.PowerSource = uint8(v.Attribute.Value.(uint64))
			}
		case "LocationDescription":
			if v.Status == cluster.ZclStatusSuccess {
				basicInfo.LocationDescription = v.Attribute.Value.(string)
			}
		case "PhysicalEnvironment":
			if v.Status == cluster.ZclStatusSuccess {
				basicInfo.PhysicalEnvironment = uint8(v.Attribute.Value.(uint64))
			}
		case "DeviceEnabled":
			if v.Status == cluster.ZclStatusSuccess {
				basicInfo.DeviceEnabled = v.Attribute.Value.(bool)
			}
		case "AlarmMask":
			if v.Status == cluster.ZclStatusSuccess {
				basicInfo.AlarmMask = v.Attribute.Value
			}
		case "DisableLocalConfig":
			if v.Status == cluster.ZclStatusSuccess {
				basicInfo.DisableLocalConfig = v.Attribute.Value
			}
		case "SWBuildID":
			if v.Status == cluster.ZclStatusSuccess {
				basicInfo.SWBuildID = v.Attribute.Value.(string)
			}
		default:
			globallogger.Log.Warnln("devEUI :", terminalInfo.DevEUI, "procBasicReadRsp unknow attributeName", v.AttributeName)
		}
	}
	globallogger.Log.Infof("[devEUI: %v][genBasicProcProcReadRsp]: read rsp: basicInfo %+v", terminalInfo.DevEUI, basicInfo)
}

func procGenBasicProcRead(devEUI string, dstEndpointIndex int, clusterID uint16) {
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

func genBasicProcKeepAlive(terminalInfo config.TerminalInfo, interval uint16) {
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
		keepalive.ProcKeepAlive(terminalInfo, interval)
		go procGenBasicProcRead(terminalInfo.DevEUI, 0, 0x0006)
		go func() {
			timer := time.NewTimer(2 * time.Second)
			<-timer.C
			timer.Stop()
			procGenBasicProcRead(terminalInfo.DevEUI, 0, 0x0702)
		}()
		go func() {
			timer := time.NewTimer(4 * time.Second)
			<-timer.C
			timer.Stop()
			procGenBasicProcRead(terminalInfo.DevEUI, 0, 0x0b04)
		}()
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal1SceneSwitch005f0cf1,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
		keepalive.ProcKeepAlive(terminalInfo, interval)
	default:
		globallogger.Log.Warnf("[genBasicProcKeepAlive]: invalid tmnType: %s", terminalInfo.TmnType)
	}
}

// genBasicProcReport 处理report（0x0a）消息
func genBasicProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	genBasicProcKeepAlive(terminalInfo, uint16(terminalInfo.Interval))
}

// genBasicProcConfigureReportingResponse 处理configureReportingResponse（0x07）消息
func genBasicProcConfigureReportingResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ConfigureReportingResponse)
	globallogger.Log.Infof("[devEUI: %v][genBasicProcConfigureReportingResponse]: command: %+v", terminalInfo.DevEUI, Command)
	for _, v := range Command.AttributeStatusRecords {
		if v.Status != cluster.ZclStatusSuccess {
			globallogger.Log.Infof("[devEUI: %v][genBasicProcConfigureReportingResponse]: configReport failed: %x", terminalInfo.DevEUI, v.Status)
		} else {
			// proc keepalive
			if v.AttributeID == 0 {
				genBasicProcKeepAlive(terminalInfo, uint16(terminalInfo.Interval))
			}
		}
	}
}

//GenBasicProc 处理clusterID 0x0000即genBasic属性消息
func GenBasicProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	switch zclFrame.CommandName {
	case "ReadAttributesResponse":
		genBasicProcProcReadRsp(terminalInfo, zclFrame.Command)
	case "ReportAttributes":
		genBasicProcReport(terminalInfo, zclFrame.Command)
	case "ConfigureReportingResponse":
		genBasicProcConfigureReportingResponse(terminalInfo, zclFrame.Command)
	default:
		globallogger.Log.Warnf("[devEUI: %v][GenBasicProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
