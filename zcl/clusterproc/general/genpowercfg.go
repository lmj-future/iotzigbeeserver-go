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

func getGenPowerCfgParams(terminalID string, tmnType string, power uint64) (interface{}, bool) {
	bPublish := false
	params := make(map[string]interface{}, 2)
	var key string
	params["terminalId"] = terminalID
	switch tmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM:
		key = iotsmartspace.HeimanCOSensorPropertyLeftElectric
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
		key = iotsmartspace.HeimanHS2AQPropertyLeftElectric
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
		key = iotsmartspace.HeimanHTSensorPropertyLeftElectric
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM:
		key = iotsmartspace.HeimanPIRSensorPropertyLeftElectric
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30:
		key = iotsmartspace.HeimanPIRSensorPropertyLeftElectric
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM:
		key = iotsmartspace.HeimanSmokeSensorPropertyLeftElectric
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice:
		key = iotsmartspace.HeimanWarningDevicePropertyLeftElectric
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM:
		key = iotsmartspace.HeimanWaterSensorPropertyLeftElectric
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSceneSwitchEM30:
		key = iotsmartspace.HeimanSceneSwitchEM30PropertyLeftElectric
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30:
		bPublish = true
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM:
		bPublish = true
	default:
		globallogger.Log.Warnf("[devEUI: %v][getParams] invalid tmnType: %v", terminalID, tmnType)
	}
	params[key] = power
	return params, bPublish
}

func genPowerCfgProcAttribute(terminalInfo config.TerminalInfo, attributeName string, attribute *cluster.Attribute) {
	switch attributeName {
	case "BatteryVoltage": //电量值
	case "BatteryPercentageRemaining": //电量百分比
		// iotsmartspace publish msg to app
		if constant.Constant.Iotware {
			values := make(map[string]interface{}, 1)
			values[iotsmartspace.IotwarePropertyElectric] = attribute.Value.(uint64) / 2
			iotsmartspace.PublishTelemetryUpIotware(terminalInfo, values)
		} else if constant.Constant.Iotedge {
			params, bPublish := getGenPowerCfgParams(terminalInfo.DevEUI, terminalInfo.TmnType, attribute.Value.(uint64)/2)
			if bPublish {
				iotsmartspace.Publish(iotsmartspace.TopicZigbeeserverIotsmartspaceProperty, iotsmartspace.MethodPropertyUp, params, uuid.NewV4().String())
			}
		} else if constant.Constant.Iotprivate {
			type appDataMsg struct {
				Electric string `json:"electric"`
			}
			kafkaMsg := publicstruct.DataReportMsg{
				OIDIndex:   terminalInfo.OIDIndex,
				DevSN:      terminalInfo.DevEUI,
				LinkType:   terminalInfo.ProfileID,
				DeviceType: terminalInfo.TmnType2,
				AppData: appDataMsg{
					Electric: strconv.FormatUint(attribute.Value.(uint64)/2, 10) + "%",
				},
			}
			kafkaMsgByte, _ := json.Marshal(kafkaMsg)
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg, string(kafkaMsgByte))
		}
	default:
		globallogger.Log.Warnln("devEUI : "+terminalInfo.DevEUI+" "+"genPowerCfgProcAttribute unknow attributeName", attributeName)
	}
}

//genPowerCfgProcReadRsp 处理readRsp（0x01）消息
func genPowerCfgProcReadRsp(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReadAttributesResponse)
	for _, v := range Command.ReadAttributeStatuses {
		globallogger.Log.Infof("[devEUI: %v][genPowerCfgProcReadRsp]: readAttributesRsp: %+v", terminalInfo.DevEUI, v)
		if v.Status == cluster.ZclStatusSuccess {
			genPowerCfgProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
		}
	}
}

func procGenPowerCfgProcRead(devEUI string, dstEndpointIndex int, clusterID uint16) {
	cmd := common.Command{
		DstEndpointIndex: dstEndpointIndex,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.ReadAttribute,
		ClusterID:   clusterID,
		Command:     cmd,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

func genPowerCfgProcKeepAlive(devEUI string, tmnType string, interval uint16) {
	switch tmnType {
	case constant.Constant.TMNTYPE.MAILEKE.ZigbeeTerminalPMT1004Detector:
		keepalive.ProcKeepAlive(devEUI, interval)
		go procGenPowerCfgProcRead(devEUI, 0, 0x0402)
		go func() {
			select {
			case <-time.After(time.Duration(2) * time.Second):
				procGenPowerCfgProcRead(devEUI, 0, 0x0405)
			}
		}()
		go func() {
			select {
			case <-time.After(time.Duration(4) * time.Second):
				procGenPowerCfgProcRead(devEUI, 0, 0x0415)
			}
		}()
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
		keepalive.ProcKeepAlive(devEUI, interval)
		go procGenPowerCfgProcRead(devEUI, 0, 0x0402)
		go func() {
			select {
			case <-time.After(time.Duration(2) * time.Second):
				procGenPowerCfgProcRead(devEUI, 0, 0x0405)
			}
		}()
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSceneSwitchEM30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM:
		keepalive.ProcKeepAlive(devEUI, interval)
	default:
		globallogger.Log.Warnf("[genPowerCfgProcReport]: invalid tmnType: %s", tmnType)
	}
}

// genPowerCfgProcReport 处理report（0x0a）消息
func genPowerCfgProcReport(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ReportAttributesCommand)
	globallogger.Log.Infof("[devEUI: %v][genPowerCfgProcReport]: command: %+v", terminalInfo.DevEUI, Command)
	for _, v := range Command.AttributeReports {
		genPowerCfgProcAttribute(terminalInfo, v.AttributeName, v.Attribute)
	}
	genPowerCfgProcKeepAlive(terminalInfo.DevEUI, terminalInfo.TmnType, uint16(terminalInfo.Interval))
}

// genPowerCfgProcConfigureReportingResponse 处理configureReportingResponse（0x07）消息
func genPowerCfgProcConfigureReportingResponse(terminalInfo config.TerminalInfo, command interface{}) {
	Command := command.(*cluster.ConfigureReportingResponse)
	globallogger.Log.Infof("[devEUI: %v][genPowerCfgProcConfigureReportingResponse]: command: %+v", terminalInfo.DevEUI, Command)
	for _, v := range Command.AttributeStatusRecords {
		if v.Status != cluster.ZclStatusSuccess {
			globallogger.Log.Infof("[devEUI: %v][genPowerCfgProcConfigureReportingResponse]: configReport failed: %x", terminalInfo.DevEUI, v.Status)
		} else {
			// proc keepalive
			if v.AttributeID == 0x0021 {
				genPowerCfgProcKeepAlive(terminalInfo.DevEUI, terminalInfo.TmnType, uint16(terminalInfo.Interval))
			}
		}
		// 空气盒子处理
		if terminalInfo.TmnType == constant.Constant.TMNTYPE.MAILEKE.ZigbeeTerminalPMT1004Detector {
			type confirmMessage struct {
				Response      string `json:"response"`
				OIDIndex      string `json:"OIDIndex"`
				PropertyTopic string `json:"propertyTopic"`
			}
			confirmMsg := confirmMessage{
				Response:      "success",
				OIDIndex:      terminalInfo.OIDIndex,
				PropertyTopic: "00000003",
			}
			if v.Status != cluster.ZclStatusSuccess {
				confirmMsg.Response = "failed"
			}
			confirmMsgByte, _ := json.Marshal(confirmMsg)
			kafka.Producer(constant.Constant.KAFKA.ZigbeeKafkaProduceTopicConfirmMessage, string(confirmMsgByte))
		}
	}
}

// GenPowerCfgProc 处理clusterID 0x0001即genPowerCfg属性消息
func GenPowerCfgProc(terminalInfo config.TerminalInfo, zclFrame *zcl.Frame) {
	// globallogger.Log.Infof("[devEUI: %v][GenPowerCfgProc] Start......", terminalInfo.DevEUI)
	// globallogger.Log.Infof("[devEUI: %v][GenPowerCfgProc] zclFrame: %+v", terminalInfo.DevEUI, zclFrame)
	z := zcl.New()
	switch zclFrame.CommandName {
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandReadAttributesResponse)].Name:
		genPowerCfgProcReadRsp(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandReportAttributes)].Name:
		genPowerCfgProcReport(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandConfigureReportingResponse)].Name:
		genPowerCfgProcConfigureReportingResponse(terminalInfo, zclFrame.Command)
	case z.ClusterLibrary().Global()[uint8(cluster.ZclCommandConfigureReporting)].Name:
	default:
		globallogger.Log.Warnf("[devEUI: %v][GenPowerCfgProc] invalid commandName: %v", terminalInfo.DevEUI, zclFrame.CommandName)
	}
}
