package consumer

import (
	"encoding/hex"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/db/rabbitmq"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/zclmsgdown"
	"github.com/lib/pq"
	"github.com/streadway/amqp"
)

type httpMsg struct {
	URL     string                 `json:"url"`
	Method  string                 `json:"method"`
	Headers map[string]interface{} `json:"headers"`
	Query   map[string]interface{} `json:"query"`
	Body    map[string]interface{} `json:"body"`
}
type terminalinfo struct {
	TmnName   string `json:"tmnName"`
	DevEUI    string `json:"devEUI"`
	FirmTopic string `json:"firmTopic"`
	NwkAddr   string `json:"nwkAddr"`
	Status    bool   `json:"status"`
	TmnType   string `json:"tmnType"`
}
type secondFloorData struct {
	TmnType  string         `json:"tmnType"`
	DataList []terminalinfo `json:"dataList"`
}
type firstFloorData struct {
	TmnType  string            `json:"tmnType"`
	DataList []secondFloorData `json:"dataList"`
}

// RabbitMQConnect RabbitMQConnect
type RabbitMQConnect struct {
}

// Consumer 实现接收者
func (r *RabbitMQConnect) Consumer(msg amqp.Delivery) error {
	globallogger.Log.Warnf("[RabbitMQ]Receive msg: %+v\n", msg)
	return nil
}

func pushTerminalInfoSecondFloor(terminalInfo terminalinfo, dataList []firstFloorData,
	firstFloorData firstFloorData, firstFloorIndex int, secondFloorTmnType string) []firstFloorData {
	tmpIndex := 0
	for secondFloorIndex, secondFloorData := range firstFloorData.DataList {
		if secondFloorData.TmnType == secondFloorTmnType {
			dataList[firstFloorIndex].DataList[secondFloorIndex].DataList = append(dataList[firstFloorIndex].DataList[secondFloorIndex].DataList, terminalInfo)
		} else {
			tmpIndex++
		}
	}
	if tmpIndex == len(firstFloorData.DataList) {
		dataList[firstFloorIndex].DataList = append(dataList[firstFloorIndex].DataList, secondFloorData{
			TmnType: secondFloorTmnType,
			DataList: []terminalinfo{
				terminalInfo,
			},
		})
	}
	return dataList
}
func pushTerminalInfoFirstFloor(terminalInfo terminalinfo, dataList []firstFloorData,
	firstFloorTmnType string, secondFloorTmnType string) []firstFloorData {
	firstFloorTmpData := firstFloorData{
		TmnType: firstFloorTmnType,
		DataList: []secondFloorData{
			{
				TmnType: secondFloorTmnType,
				DataList: []terminalinfo{
					terminalInfo,
				},
			},
		},
	}
	dataList = append(dataList, firstFloorTmpData)
	return dataList
}
func pushTerminalInfo(dataList []firstFloorData, terminalInfo terminalinfo,
	firstFloorTmnType string, secondFloorTmnType string) []firstFloorData {
	if len(dataList) > 0 {
		tmpIndex := 0
		for k, v := range dataList {
			if v.TmnType == firstFloorTmnType {
				dataList = pushTerminalInfoSecondFloor(terminalInfo, dataList, v, k, secondFloorTmnType)
			} else {
				tmpIndex++
			}
		}
		if tmpIndex == len(dataList) {
			dataList = pushTerminalInfoFirstFloor(terminalInfo, dataList, firstFloorTmnType, secondFloorTmnType)
		}
	} else {
		dataList = pushTerminalInfoFirstFloor(terminalInfo, dataList, firstFloorTmnType, secondFloorTmnType)
	}
	return dataList
}
func sortTerminalInfoList(terminalInfoList []config.TerminalInfo) []firstFloorData {
	dataList := []firstFloorData{}
	for _, v := range terminalInfoList {
		if !v.IsExist || v.TmnType == "invalidType" {
			continue
		}
		var terminalInfo terminalinfo = terminalinfo{}
		terminalInfo.DevEUI = v.DevEUI
		terminalInfo.FirmTopic = v.FirmTopic + "（" + v.ManufacturerName + "）"
		terminalInfo.NwkAddr = v.NwkAddr
		terminalInfo.TmnName = v.TmnName
		terminalInfo.Status = v.Online
		terminalInfo.TmnType = v.TmnType2
		switch v.TmnType {
		case constant.Constant.TMNTYPE.SWITCH.ZigbeeTerminalSingleSwitch,
			constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141:
			dataList = pushTerminalInfo(dataList, terminalInfo, "开关", "一键开关")
		case constant.Constant.TMNTYPE.SWITCH.ZigbeeTerminalDoubleSwitch,
			constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142:
			dataList = pushTerminalInfo(dataList, terminalInfo, "开关", "二键开关")
		case constant.Constant.TMNTYPE.SWITCH.ZigbeeTerminalTripleSwitch,
			constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
			dataList = pushTerminalInfo(dataList, terminalInfo, "开关", "三键开关")
		case constant.Constant.TMNTYPE.SWITCH.ZigbeeTerminalQuadrupleSwitch:
			dataList = pushTerminalInfo(dataList, terminalInfo, "开关", "四键开关")
		case constant.Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalSingleLightingSwitch:
			dataList = pushTerminalInfo(dataList, terminalInfo, "联动开关", "一键联动开关")
		case constant.Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalDoubleLightingSwitch:
			dataList = pushTerminalInfo(dataList, terminalInfo, "联动开关", "二键联动开关")
		case constant.Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalTripleLightingSwitch:
			dataList = pushTerminalInfo(dataList, terminalInfo, "联动开关", "三键联动开关")
		case constant.Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalQuadrupleLightingSwitch:
			dataList = pushTerminalInfo(dataList, terminalInfo, "联动开关", "四键联动开关")
		case constant.Constant.TMNTYPE.SCENE.ZigbeeTerminalScene:
			dataList = pushTerminalInfo(dataList, terminalInfo, "场景面板", "场景面板")
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSceneSwitchEM30:
			dataList = pushTerminalInfo(dataList, terminalInfo, "场景面板", "移动场景面板")
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
			dataList = pushTerminalInfo(dataList, terminalInfo, "场景面板", "墙插场景面板")
		case constant.Constant.TMNTYPE.CURTAIN.ZigbeeTerminalCurtain:
			dataList = pushTerminalInfo(dataList, terminalInfo, "窗帘", "智能窗帘")
		case constant.Constant.TMNTYPE.SOCKET.ZigbeeTerminalSocket,
			constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
			dataList = pushTerminalInfo(dataList, terminalInfo, "插座", "智能插座")
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
			dataList = pushTerminalInfo(dataList, terminalInfo, "插座", "可移动式计量插座")
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105:
			dataList = pushTerminalInfo(dataList, terminalInfo, "插座", "10A智能墙插")
		case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
			dataList = pushTerminalInfo(dataList, terminalInfo, "插座", "16A智能墙插")
		case constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalHumanDetector,
			constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalPIRSensorHY0027:
			dataList = pushTerminalInfo(dataList, terminalInfo, "传感器", "人体红外探测器")
		case constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalDoorSensor,
			constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30:
			dataList = pushTerminalInfo(dataList, terminalInfo, "传感器", "门磁")
		case constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalSmokeFireDetector,
			constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSmokeSensorHY0024:
			dataList = pushTerminalInfo(dataList, terminalInfo, "传感器", "烟火灾探测器")
		case constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalWaterDetector,
			constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM:
			dataList = pushTerminalInfo(dataList, terminalInfo, "传感器", "水浸探测器")
		case constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalGasDetector,
			constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalGASSensorHY0022:
			dataList = pushTerminalInfo(dataList, terminalInfo, "传感器", "可燃气体探测器")
		case constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalHumitureDetector,
			constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
			dataList = pushTerminalInfo(dataList, terminalInfo, "传感器", "温湿度探测器")
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM:
			dataList = pushTerminalInfo(dataList, terminalInfo, "传感器", "一氧化碳探测器")
		case constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalVoiceLightDetector,
			constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice,
			constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalWarningDevice005b0e12:
			dataList = pushTerminalInfo(dataList, terminalInfo, "报警器", "声光报警器")
		case constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalEmergencyButton:
			dataList = pushTerminalInfo(dataList, terminalInfo, "报警器", "无线紧急按钮")
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
			dataList = pushTerminalInfo(dataList, terminalInfo, "空气质量仪", "空气质量仪")
		case constant.Constant.TMNTYPE.MAILEKE.ZigbeeTerminalPMT1004Detector:
			dataList = pushTerminalInfo(dataList, terminalInfo, "空气盒子", "PM2.5检测仪")
		case constant.Constant.TMNTYPE.TRANSPONDER.ZigbeeTerminalInfraredTransponder:
			dataList = pushTerminalInfo(dataList, terminalInfo, "红外设备", "红外宝")
		case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM:
			dataList = pushTerminalInfo(dataList, terminalInfo, "红外设备", "红外转发器")
		default:
			dataList = pushTerminalInfo(dataList, terminalInfo, "其他", "未分类")
		}
	}
	return dataList
}
func procGetZigbeeNetwork(query map[string]interface{}) rabbitmq.Response {
	response := rabbitmq.Response{}
	if _, ok := query["devSN"]; !ok {
		response.Code = 1
		response.Message = "params is invalid"
		return response
	}
	if _, ok := query["moduleId"]; !ok {
		response.Code = 1
		response.Message = "params is invalid"
		return response
	}
	devSN := query["devSN"].(string)
	moduleIDNum, _ := strconv.Atoi(query["moduleId"].(string))
	moduleID := "00" + strconv.FormatInt(int64(moduleIDNum), 16)
	moduleID = moduleID[len(moduleID)-2:]
	var terminalInfos []config.TerminalInfo
	var err error
	if devSN[2:10] == "9801A26U" || devSN[2:10] == "9801A26N" {
		nodeIndexNum, _ := strconv.Atoi(query["nodeIndex"].(string))
		nodeIndex := "00" + strconv.FormatInt(int64(nodeIndexNum), 16)
		nodeIndex = nodeIndex[len(nodeIndex)-2:]
		moduleID = nodeIndex[1:] + moduleID[1:]
	}
	if constant.Constant.UsePostgres {
		terminalInfos, err = models.FindTerminalByAPMacAndModuleIDPG(hex.EncodeToString([]byte(devSN)), moduleID)
	} else {
		terminalInfos, err = models.FindTerminalByAPMacAndModuleID(hex.EncodeToString([]byte(devSN)), moduleID)
	}
	if err == nil && len(terminalInfos) > 0 {
		response.Data = sortTerminalInfoList(terminalInfos)
	}
	return response
}
func procSetLogLevel(query map[string]interface{}) rabbitmq.Response {
	response := rabbitmq.Response{}
	if _, ok := query["level"]; !ok {
		response.Code = 1
		response.Message = "params is invalid"
		return response
	}
	globallogger.SetLogLevel(query["level"].(string))
	return response
}
func heimanIRControlEMSendKeyCommand(devEUI string, ID string, KeyCode string, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanIRControlEMSendKeyCommand] ID: %s, KeyCode: %s", devEUI, ID, KeyCode)
	id, _ := strconv.Atoi(ID)
	keyCode, _ := strconv.Atoi(KeyCode)
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SendKeyCommand,
		ClusterID:   0xfc82,
		Command: common.Command{
			InfraredRemoteID:      uint8(id),
			InfraredRemoteKeyCode: uint8(keyCode),
		},
		MsgID: msgID,
	})
}
func heimanIRControlEMStudyKey(devEUI string, ID string, KeyCode string, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanIRControlEMStudyKey] ID: %s, KeyCode: %s", devEUI, ID, KeyCode)
	id, _ := strconv.Atoi(ID)
	keyCode, _ := strconv.Atoi(KeyCode)
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.StudyKey,
		ClusterID:   0xfc82,
		Command: common.Command{
			InfraredRemoteID:      uint8(id),
			InfraredRemoteKeyCode: uint8(keyCode),
		},
		MsgID: msgID,
	})
}
func heimanIRControlEMDeleteKey(devEUI string, ID string, KeyCode string, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanIRControlEMDeleteKey] ID: %s, KeyCode: %s", devEUI, ID, KeyCode)
	id, _ := strconv.Atoi(ID)
	keyCode, _ := strconv.Atoi(KeyCode)
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.DeleteKey,
		ClusterID:   0xfc82,
		Command: common.Command{
			InfraredRemoteID:      uint8(id),
			InfraredRemoteKeyCode: uint8(keyCode),
		},
		MsgID: msgID,
	})
}
func heimanIRControlEMCreateID(devEUI string, ModelType string, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanIRControlEMCreateID] ModelType: %s", devEUI, ModelType)
	modelType, _ := strconv.Atoi(ModelType)
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.CreateID,
		ClusterID:   0xfc82,
		Command: common.Command{
			InfraredRemoteModelType: uint8(modelType),
		},
		MsgID: msgID,
	})
}
func heimanIRControlEMGetIDAndKeyCodeList(devEUI string, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanIRControlEMGetIDAndKeyCodeList]", devEUI)
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.GetIDAndKeyCodeList,
		ClusterID:   0xfc82,
		MsgID:       msgID,
	})
}
func procTestTerminalJoin(body map[string]interface{}) rabbitmq.Response {
	response := rabbitmq.Response{}
	if body["devEUI"].(string) == "" || body["nwkAddr"].(string) == "" || body["tmnType"].(string) == "" || body["ACMac"].(string) == "" ||
		body["APMac"].(string) == "" || body["T300ID"].(string) == "" || body["moduleID"].(string) == "" || body["shopId"].(string) == "" ||
		body["userName"].(string) == "" {
		response.Code = 1
		response.Message = "params is invalid"
		return response
	}
	terminalInfo := config.TerminalInfo{
		DevEUI:       body["devEUI"].(string),
		NwkAddr:      body["nwkAddr"].(string),
		TmnType:      body["tmnType"].(string),
		ACMac:        body["ACMac"].(string),
		FirstAddr:    publicfunction.Transport16StringToString(body["ACMac"].(string)),
		APMac:        body["APMac"].(string),
		SecondAddr:   publicfunction.Transport16StringToString(body["APMac"].(string)),
		T300ID:       body["T300ID"].(string),
		ThirdAddr:    body["T300ID"].(string),
		ModuleID:     body["moduleID"].(string),
		ProfileID:    "ZHA",
		Online:       true,
		IsDiscovered: true,
		IsExist:      true,
		IsNeedBind:   true,
		IsReadBasic:  true,
	}
	terminalInfo.TmnName = terminalInfo.SecondAddr + "-" + terminalInfo.DevEUI[6:]
	switch body["tmnType"].(string) {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0001", "0009", "0500"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01"}
		terminalInfo.EndpointPG = pq.StringArray{"01"}
		terminalInfo.EndpointCount = 1
		terminalInfo.ManufacturerName = "HEIMAN"
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalGASSensorHY0022,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSmokeSensorHY0024,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalPIRSensorHY0027:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0001", "0009", "0500"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01"}
		terminalInfo.EndpointPG = pq.StringArray{"01"}
		terminalInfo.EndpointCount = 1
		terminalInfo.ManufacturerName = "HONYAR"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0006", "0702"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01"}
		terminalInfo.EndpointPG = pq.StringArray{"01"}
		terminalInfo.EndpointCount = 1
		terminalInfo.ManufacturerName = "HEIMAN"
		terminalInfo.Attribute.Status = []string{"OFF"}
		terminalInfo.Attribute.StatusPG = pq.StringArray{"OFF"}
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0001", "0009", "0402", "0405", "042a", "042b", "fc81"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
			{
				ClusterID:   []string{"0021"},
				SrcEndpoint: "f2",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01", "f2"}
		terminalInfo.EndpointPG = pq.StringArray{"01", "f2"}
		terminalInfo.EndpointCount = 2
		terminalInfo.ManufacturerName = "HEIMAN"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0006"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01"}
		terminalInfo.EndpointPG = pq.StringArray{"01"}
		terminalInfo.EndpointCount = 1
		terminalInfo.ManufacturerName = "HEIMAN"
		terminalInfo.Attribute.Status = []string{"OFF"}
		terminalInfo.Attribute.StatusPG = pq.StringArray{"OFF"}
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0006"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
			{
				ClusterID:   []string{"0006"},
				SrcEndpoint: "02",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01", "02"}
		terminalInfo.EndpointPG = pq.StringArray{"01", "02"}
		terminalInfo.EndpointCount = 2
		terminalInfo.ManufacturerName = "HEIMAN"
		terminalInfo.Attribute.Status = []string{"OFF", "OFF"}
		terminalInfo.Attribute.StatusPG = pq.StringArray{"OFF", "OFF"}
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0006"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
			{
				ClusterID:   []string{"0006"},
				SrcEndpoint: "02",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
			{
				ClusterID:   []string{"0006"},
				SrcEndpoint: "03",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01", "02", "03"}
		terminalInfo.EndpointPG = pq.StringArray{"01", "02", "03"}
		terminalInfo.EndpointCount = 3
		terminalInfo.ManufacturerName = "HEIMAN"
		terminalInfo.Attribute.Status = []string{"OFF", "OFF", "OFF"}
		terminalInfo.Attribute.StatusPG = pq.StringArray{"OFF", "OFF", "OFF"}
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0001", "0402"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
			{
				ClusterID:   []string{"0405"},
				SrcEndpoint: "02",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01", "02"}
		terminalInfo.EndpointPG = pq.StringArray{"01", "02"}
		terminalInfo.EndpointCount = 2
		terminalInfo.ManufacturerName = "HEIMAN"
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalPMTSensor0001112b:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0001", "0402", "0405"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01"}
		terminalInfo.EndpointPG = pq.StringArray{"01"}
		terminalInfo.EndpointCount = 1
		terminalInfo.ManufacturerName = "HONYAR"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSceneSwitchEM30:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0001"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01"}
		terminalInfo.EndpointPG = pq.StringArray{"01"}
		terminalInfo.EndpointCount = 1
		terminalInfo.ManufacturerName = "HEIMAN"
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0001", "0502"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01"}
		terminalInfo.EndpointPG = pq.StringArray{"01"}
		terminalInfo.EndpointCount = 1
		terminalInfo.ManufacturerName = "HEIMAN"
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalWarningDevice005b0e12:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0001", "0502"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01"}
		terminalInfo.EndpointPG = pq.StringArray{"01"}
		terminalInfo.EndpointCount = 1
		terminalInfo.ManufacturerName = "HONYAR"
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0006", "0702"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01"}
		terminalInfo.EndpointPG = pq.StringArray{"01"}
		terminalInfo.EndpointCount = 1
		terminalInfo.ManufacturerName = "HONYAR"
		terminalInfo.Attribute.Status = []string{"OFF"}
		terminalInfo.Attribute.StatusPG = pq.StringArray{"OFF"}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0006"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01"}
		terminalInfo.EndpointPG = pq.StringArray{"01"}
		terminalInfo.EndpointCount = 1
		terminalInfo.ManufacturerName = "HONYAR"
		terminalInfo.Attribute.Status = []string{"OFF"}
		terminalInfo.Attribute.StatusPG = pq.StringArray{"OFF"}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0006"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
			{
				ClusterID:   []string{"0006"},
				SrcEndpoint: "02",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01", "02"}
		terminalInfo.EndpointPG = pq.StringArray{"01", "02"}
		terminalInfo.EndpointCount = 2
		terminalInfo.ManufacturerName = "HONYAR"
		terminalInfo.Attribute.Status = []string{"OFF", "OFF"}
		terminalInfo.Attribute.StatusPG = pq.StringArray{"OFF", "OFF"}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"0006"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
			{
				ClusterID:   []string{"0006"},
				SrcEndpoint: "02",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
			{
				ClusterID:   []string{"0006"},
				SrcEndpoint: "03",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01", "02", "03"}
		terminalInfo.EndpointPG = pq.StringArray{"01", "02", "03"}
		terminalInfo.EndpointCount = 3
		terminalInfo.ManufacturerName = "HONYAR"
		terminalInfo.Attribute.Status = []string{"OFF", "OFF", "OFF"}
		terminalInfo.Attribute.StatusPG = pq.StringArray{"OFF", "OFF", "OFF"}
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal1SceneSwitch005f0cf1,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
		terminalInfo.BindInfo = []config.BindInfo{
			{
				ClusterID:   []string{"fe05"},
				SrcEndpoint: "01",
				DstEndpoint: "ff",
				DstAddrMode: "03",
				DstAddress:  "ffffffffffffffff",
			},
		}
		bindInfoByte, _ := json.Marshal(terminalInfo.BindInfo)
		terminalInfo.BindInfoPG = string(bindInfoByte)
		terminalInfo.Endpoint = []string{"01"}
		terminalInfo.EndpointPG = pq.StringArray{"01"}
		terminalInfo.EndpointCount = 1
		terminalInfo.ManufacturerName = "HONYAR"
	default:
	}
	if constant.Constant.Iotprivate {
		tmnTypeInfo := publicfunction.HTTPRequestTerminalTypeByAlias(terminalInfo.ManufacturerName + "_" + terminalInfo.TmnType)
		if tmnTypeInfo != nil {
			type tmnList struct {
				TmnName   string `json:"tmnName"`
				TmnDevSN  string `json:"tmnDevSN"`
				AddSource string `json:"addSource"`
			}
			type tmnInfo struct {
				FirmName     string    `json:"firmName"`
				TerminalType string    `json:"terminalType"`
				SceneID      string    `json:"sceneID"`
				TenantID     string    `json:"tenantID"`
				TmnList      []tmnList `json:"tmnList"`
			}
			tmnInfoTemp := tmnInfo{
				FirmName:     tmnTypeInfo.FirmName,
				TerminalType: tmnTypeInfo.TerminalType,
				SceneID:      body["shopId"].(string),
				TenantID:     body["userName"].(string),
				TmnList: []tmnList{
					{
						TmnName:   terminalInfo.TmnName,
						TmnDevSN:  terminalInfo.DevEUI,
						AddSource: "自动上线",
					},
				},
			}
			tmnInfoByte, _ := json.Marshal(tmnInfoTemp)
			if publicfunction.HTTPRequestAddTerminal(terminalInfo.DevEUI, tmnInfoByte) {
				httpTerminalInfo := publicfunction.HTTPRequestTerminalInfo(terminalInfo.DevEUI, "ZHA")
				if httpTerminalInfo != nil {
					if constant.Constant.UsePostgres {
					} else {
						terminalInfo.OIDIndex = httpTerminalInfo.TmnOIDIndex
						terminalInfo.ScenarioID = httpTerminalInfo.SceneID
						terminalInfo.UserName = httpTerminalInfo.TenantID
						terminalInfo.FirmTopic = httpTerminalInfo.FirmTopic
						terminalInfo.ProfileID = httpTerminalInfo.LinkType
						terminalInfo.TmnType2 = httpTerminalInfo.TmnType
						terminalInfo.IsExist = true
						models.CreateTerminal(terminalInfo)
					}
					publicfunction.TerminalOnline(terminalInfo.DevEUI, true)
				}
			}
		}
	}
	return response
}

// ConsumerByIotwebserver 处理来自iotwebserver的消息
func (r *RabbitMQConnect) ConsumerByIotwebserver(msg amqp.Delivery) (rabbitmq.Response, error) {
	defer func() {
		err := recover()
		if err != nil {
			globallogger.Log.Errorln("ConsumerByIotwebserver err :", err)
		}
	}()
	globallogger.Log.Warnf("[RabbitMQ]Receive msg: %+v", string(msg.Body))
	response := rabbitmq.Response{}
	httpMsg := httpMsg{}
	json.Unmarshal(msg.Body, &httpMsg)
	urlInfo, err := url.Parse(httpMsg.URL)
	if err == nil {
		switch urlInfo.Path {
		case "/iotzigbeeurl/getZigbeeNetwork":
			response = procGetZigbeeNetwork(httpMsg.Query)
		case "/iotzigbeeurl/logLevel":
			response = procSetLogLevel(httpMsg.Query)
		case "/iotzigbeeurl/sendKeyCommand":
			heimanIRControlEMSendKeyCommand(httpMsg.Query["devEUI"].(string), httpMsg.Query["ID"].(string), httpMsg.Query["KeyCode"].(string), 1)
		case "/iotzigbeeurl/studyKey":
			heimanIRControlEMStudyKey(httpMsg.Query["devEUI"].(string), httpMsg.Query["ID"].(string), httpMsg.Query["KeyCode"].(string), 1)
		case "/iotzigbeeurl/deleteKey":
			heimanIRControlEMDeleteKey(httpMsg.Query["devEUI"].(string), httpMsg.Query["ID"].(string), httpMsg.Query["KeyCode"].(string), 1)
		case "/iotzigbeeurl/createID":
			heimanIRControlEMCreateID(httpMsg.Query["devEUI"].(string), httpMsg.Query["ModelType"].(string), 1)
		case "/iotzigbeeurl/getIDAndKeyCodeList":
			heimanIRControlEMGetIDAndKeyCodeList(httpMsg.Query["devEUI"].(string), 1)
		case "/iotzigbeeurl/test/terminalJoin":
			response = procTestTerminalJoin(httpMsg.Body)
		}
	}
	return response, nil
}

// DeleteTerminalMsg DeleteTerminalMsg
type DeleteTerminalMsg struct {
	MsgType     string `json:"msgType"`
	TmnOIDIndex string `json:"tmnOIDIndex"`
}

// Consumer 实现接收者
func (d *DeleteTerminalMsg) Consumer(msg amqp.Delivery) error {
	globallogger.Log.Warnf("[RabbitMQ]Receive delete terminal msg: %+v\n", msg)
	deleteTerminalMsg := DeleteTerminalMsg{}
	err := json.Unmarshal(msg.Body, &deleteTerminalMsg)
	if err == nil {
		publicfunction.ProcTerminalDeleteByOIDIndex(deleteTerminalMsg.TmnOIDIndex)
	}
	return nil
}

// ConsumerByIotwebserver 处理来自iotwebserver的消息
func (d *DeleteTerminalMsg) ConsumerByIotwebserver(msg amqp.Delivery) (rabbitmq.Response, error) {
	globallogger.Log.Warnf("[RabbitMQ]Receive msg: %+v\n", msg)
	return rabbitmq.Response{}, nil
}

// RabbitMQConnection RabbitMQConnection
func RabbitMQConnection() {
	queueExchange := &rabbitmq.QueueExchange{
		QuName: "iotzigbeetmnmgr",
		RtKey:  "iotzigbeetmnmgr",
		ExName: "ex_direct",
		ExType: "direct",
	}
	mq := rabbitmq.New(queueExchange, "", "")
	d := &DeleteTerminalMsg{}
	mq.RegisterReceiver(d)
	mq.StartReceiver()

	queueExchange2 := &rabbitmq.QueueExchange{
		QuName: "iotzigbeeurl",
		RtKey:  "iotzigbeeurl",
		ExName: "ex_direct",
		ExType: "direct",
	}
	mq2 := rabbitmq.New(queueExchange2, "", "")
	r := &RabbitMQConnect{}
	mq2.RegisterReceiver(r)
	mq2.StartReceiver()
}
