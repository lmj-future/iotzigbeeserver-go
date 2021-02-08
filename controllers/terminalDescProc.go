package controllers

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	zclmain "github.com/h3c/iotzigbeeserver-go/zcl"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
	"github.com/lib/pq"
)

func firstlyGetTerminalEndpoint(jsonInfo publicstruct.JSONInfo, devEUI string) {
	var terminalInfo = publicfunction.GetTerminalInfo(jsonInfo, devEUI)
	publicfunction.GetTerminalEndpoint(terminalInfo)
}

func secondlyReadBasicByEndpoint(devEUI string, ActiveEPCount string, ActiveEPList string) {
	tempActiveEPCount, _ := strconv.ParseInt(ActiveEPCount, 16, 0)
	var endpointArray = make([]string, int(tempActiveEPCount))
	for i := 0; i < int(tempActiveEPCount); i++ {
		tempValue, _ := strconv.ParseInt(ActiveEPCount, 16, 0)
		endpointArray[int(tempValue)-i-1] = ActiveEPList[0:2]
		ActiveEPList = ActiveEPList[2:]
	}
	var err error
	if constant.Constant.UsePostgres {
		oMatch := make(map[string]interface{})
		oSet := make(map[string]interface{})
		oMatch["deveui"] = devEUI
		oSet["endpointpg"] = pq.StringArray(endpointArray)
		oSet["endpointcount"] = len(endpointArray)
		oSet["updatetime"] = time.Now()
		_, err = models.FindTerminalAndUpdatePG(oMatch, oSet)
	} else {
		var setData = config.TerminalInfo{}
		setData.Endpoint = endpointArray
		setData.EndpointCount = len(endpointArray)
		// if terminalInfo.FirmTopic == "8001d001" && terminalInfo.TmnType == constant.Constant.TMNTYPE.SWITCH.ZigbeeTerminalQuadrupleSwitch {
		// 	//针对狄耐克厂商的四联开关做特殊处理
		// 	var temp = setData.Endpoint[1]
		// 	setData.Endpoint[1] = setData.Endpoint[2]
		// 	setData.Endpoint[2] = temp
		// }
		setData.UpdateTime = time.Now()
		_, err = models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, setData)
	}
	if err != nil {
		globallogger.Log.Errorln("devEUI : "+devEUI+" "+"secondlyReadBasicByEndpoint FindTerminalAndUpdate err : ", err)
		return
	}
	endpointIndex := 0
	for index, value := range endpointArray {
		if value == "01" {
			endpointIndex = index
		}
	}
	terminalInfo := config.TerminalInfo{
		DevEUI: devEUI,
	}
	cmd := common.Command{
		DstEndpointIndex: endpointIndex,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.ReadBasic,
		ClusterID:   0x0000,
		Command:     cmd,
	}
	zclmain.ZclMain(globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent, terminalInfo, zclDownMsg, "", "", nil)
}

func thirdlyDiscoveryByEndpointFor(terminalInfo config.TerminalInfo, endpoint string) {
	go publicfunction.SendTerminalDiscovery(terminalInfo, endpoint)
}

func thirdlyDiscoveryByEndpoint(devEUI string, setData config.TerminalInfo, terminalInfo config.TerminalInfo) {
	_, err := models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, setData)
	if err != nil {
		globallogger.Log.Errorln("devEUI : "+devEUI+" "+"thirdlyDiscoveryByEndpoint FindTerminalAndUpdate err : ", err)
		return
	}
	for _, item := range setData.Endpoint {
		thirdlyDiscoveryByEndpointFor(terminalInfo, item)
	}
}

func thirdlyDiscoveryByEndpointPG(devEUI string, oSet map[string]interface{}, terminalInfo config.TerminalInfo) {
	_, err := models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, oSet)
	if err != nil {
		globallogger.Log.Errorln("devEUI : "+devEUI+" "+"thirdlyDiscoveryByEndpointPG FindTerminalAndUpdatePG err : ", err)
		return
	}
	for _, item := range oSet["endpointpg"].(pq.StringArray) {
		thirdlyDiscoveryByEndpointFor(terminalInfo, item)
	}
}

func getClusterIDByDeviceID(tmnType string, ProfileID string, DeviceID string, manufacturerName string, endpoint string) []string {
	var clusterIDArray = []string{}
	deviceID, _ := strconv.ParseUint(DeviceID, 16, 16)
	clusterID := ""
	switch uint16(deviceID) {
	case cluster.SmartPlugDevice.DeviceID,
		cluster.OnOffOutputDevice.DeviceID:
		switch manufacturerName {
		case constant.Constant.MANUFACTURERNAME.HeiMan:
			switch tmnType {
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket,
				constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug:
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.OnOff), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.SmartEnergyMetering), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.ElectricalMeasurement), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30,
				constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30,
				constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30:
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.OnOff), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
			}
		case constant.Constant.MANUFACTURERNAME.Honyar:
			switch tmnType {
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.OnOff), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.OnOff), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.SmartEnergyMetering), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.ElectricalMeasurement), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
			}
		default:
			globallogger.Log.Warnln("getClusterIDByDeviceID unknow manufacturerName: ", manufacturerName)
			clusterID = "0000" + strconv.FormatUint(uint64(cluster.OnOff), 16)
			clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
		}
	case cluster.SceneSelectorDevice.DeviceID: //情景面板
		switch manufacturerName {
		case constant.Constant.MANUFACTURERNAME.HeiMan:
			clusterID = "0000" + strconv.FormatUint(uint64(cluster.PowerConfiguration), 16)
			clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
		case constant.Constant.MANUFACTURERNAME.Honyar:
		default:
			globallogger.Log.Warnln("getClusterIDByDeviceID unknow manufacturerName: ", manufacturerName)
		}
	case cluster.RemoteControlDevice.DeviceID: //远程控制器
		switch manufacturerName {
		case constant.Constant.MANUFACTURERNAME.HeiMan:
			clusterID = "0000" + strconv.FormatUint(uint64(cluster.PowerConfiguration), 16)
			clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
			// clusterID = "0000" + strconv.FormatUint(uint64(cluster.HEIMANInfraredRemote), 16)
			// clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
		}
	case cluster.OnOffLightDevice.DeviceID: //灯光面板
		clusterID = "0000" + strconv.FormatUint(uint64(cluster.OnOff), 16)
		clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
	case cluster.MainsPowerOutletDevice.DeviceID: //智能插座
		switch manufacturerName {
		case constant.Constant.MANUFACTURERNAME.HeiMan:
			switch tmnType {
			case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket:
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.OnOff), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.SmartEnergyMetering), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.ElectricalMeasurement), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
			}
		case constant.Constant.MANUFACTURERNAME.Honyar:
			switch tmnType {
			case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105,
				constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106:
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.OnOff), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.SmartEnergyMetering), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.ElectricalMeasurement), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
			}
		default:
			globallogger.Log.Warnln("getClusterIDByDeviceID unknow manufacturerName: ", manufacturerName)
			clusterID = "0000" + strconv.FormatUint(uint64(cluster.OnOff), 16)
			clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
		}
	case cluster.SceneSelectorDevice.DeviceID: //场景面板
	case cluster.OnOffLightSwitchDevice.DeviceID: //灯光联动面板
	case cluster.OccupancySensorDevice.DeviceID: //占位传感器
		clusterID = "0000" + strconv.FormatUint(uint64(cluster.IASZone), 16)
		clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
		clusterID = "0000" + strconv.FormatUint(uint64(cluster.PowerConfiguration), 16)
		clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
	case cluster.WindowCoveringDeviceDevice.DeviceID: //窗帘
		// clusterID = "0000" + strconv.FormatUint(uint64(cluster.WindowCovering), 16)
		// clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
	case cluster.WindowCoveringControllerDevice.DeviceID: //辅助开关
	case cluster.IasZoneDevice.DeviceID: //紧急按钮、门磁、人体探测器、燃气泄漏、水浸探测、烟火探测
		clusterID = "0000" + strconv.FormatUint(uint64(cluster.IASZone), 16)
		clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
		switch manufacturerName {
		case constant.Constant.MANUFACTURERNAME.HeiMan:
			if tmnType != constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30 &&
				tmnType != constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM {
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.PowerConfiguration), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
			}
		}
	case cluster.TemperatureSensorDevice.DeviceID: //温湿度探测器
		switch manufacturerName {
		case constant.Constant.MANUFACTURERNAME.HeiMan:
			if tmnType == constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM {
				if endpoint == "01" {
					clusterID = "0000" + strconv.FormatUint(uint64(cluster.PowerConfiguration), 16)
					clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
					clusterID = "0000" + strconv.FormatUint(uint64(cluster.TemperatureMeasurement), 16)
					clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				} else if endpoint == "02" {
					clusterID = "0000" + strconv.FormatUint(uint64(cluster.RelativeHumidityMeasurement), 16)
					clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				}
			} else if tmnType == constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM {
				if endpoint == "01" {
					clusterID = "0000" + strconv.FormatUint(uint64(cluster.PowerConfiguration), 16)
					clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
					clusterID = "0000" + strconv.FormatUint(uint64(cluster.TemperatureMeasurement), 16)
					clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
					clusterID = "0000" + strconv.FormatUint(uint64(cluster.RelativeHumidityMeasurement), 16)
					clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
					clusterID = "0000" + strconv.FormatUint(uint64(cluster.HEIMANPM25Measurement), 16)
					clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
					clusterID = "0000" + strconv.FormatUint(uint64(cluster.HEIMANFormaldehydeMeasurement), 16)
					clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
					clusterID = "0000" + strconv.FormatUint(uint64(cluster.HEIMANAirQualityMeasurement), 16)
					clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				} else if endpoint == "f2" {
					clusterIDArray = append(clusterIDArray, "0021")
				}
			}
		default:
			globallogger.Log.Warnln("getClusterIDByDeviceID unknow manufacturerName: ", manufacturerName)
			if tmnType == constant.Constant.TMNTYPE.SENSOR.ZigbeeTerminalHumitureDetector {
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.IASZone), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.PowerConfiguration), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.TemperatureMeasurement), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.RelativeHumidityMeasurement), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
			} else if tmnType == constant.Constant.TMNTYPE.MAILEKE.ZigbeeTerminalPMT1004Detector {
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.PowerConfiguration), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.TemperatureMeasurement), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
				clusterID = "0000" + strconv.FormatUint(uint64(cluster.RelativeHumidityMeasurement), 16)
				clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
			}
		}
	case cluster.IasWarningDevice.DeviceID: //声光探测
		switch manufacturerName {
		case constant.Constant.MANUFACTURERNAME.HeiMan:
			clusterID = "0000" + strconv.FormatUint(uint64(cluster.IASZone), 16)
			clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
			clusterID = "0000" + strconv.FormatUint(uint64(cluster.IASWarningDevice), 16)
			clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
		default:
			globallogger.Log.Warnln("getClusterIDByDeviceID unknow manufacturerName: ", manufacturerName)
			clusterID = "0000" + strconv.FormatUint(uint64(cluster.IASZone), 16)
			clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
			clusterID = "0000" + strconv.FormatUint(uint64(cluster.IASWarningDevice), 16)
			clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:], "ffff")
		}
	case 0xc000: //狄耐克红外宝
		clusterIDArray = append(clusterIDArray, "fc01")
	case 0x0061: //麦乐克PM2.5检测仪
		clusterIDArray = append(clusterIDArray, "fe02")
	case 0x0309: //狄耐克PM2.5检测仪
		clusterID = "0000" + strconv.FormatUint(uint64(cluster.PowerConfiguration), 16)
		clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
		clusterID = "0000" + strconv.FormatUint(uint64(cluster.TemperatureMeasurement), 16)
		clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:])
		clusterID = "0000" + strconv.FormatUint(uint64(cluster.RelativeHumidityMeasurement), 16)
		clusterIDArray = append(clusterIDArray, clusterID[len(clusterID)-4:], "0415")
	default:
		globallogger.Log.Warnln("getClusterIDByDeviceID unknow deviceID: ", deviceID)
	}
	return clusterIDArray
}

func procNeedntBindTerminal(terminalInfo config.TerminalInfo) {
	if terminalInfo.IsExist {
		publicfunction.TerminalOnline(terminalInfo.DevEUI, false)
		if constant.Constant.Iotware {
			iotsmartspace.StateTerminalOnlineIotware(terminalInfo)
		} else if constant.Constant.Iotedge {
			iotsmartspace.StateTerminalOnline(terminalInfo.DevEUI)
		}
	} else {
		if constant.Constant.Iotedge {
			go iotsmartspace.ActionInsertReply(terminalInfo)
		}
	}
	if publicfunction.CheckTerminalIsSensor(terminalInfo.DevEUI, terminalInfo.TmnType) {
		go sensorSendWriteReq(terminalInfo)
	}
	publicfunction.TerminalOnline(terminalInfo.DevEUI, true)
	go func() {
		select {
		case <-time.After(time.Duration(2) * time.Second):
			iotsmartspace.ActionInsertSuccess(terminalInfo)
		}
	}()
}

func updateTerminalBindInfo(devEUI string, terminalInfo config.TerminalInfo, ProfileID string, DeviceID string,
	clusterIDInArray []string, clusterIDOutArray []string, Endpoint string) (*config.TerminalInfo, error) {
	bindInfoTemp := []config.BindInfo{}
	endpointTemp := pq.StringArray{}
	if constant.Constant.UsePostgres {
		endpointTemp = terminalInfo.EndpointPG
		bindInfoTemp = make([]config.BindInfo, len(endpointTemp))
		json.Unmarshal([]byte(terminalInfo.BindInfoPG), &bindInfoTemp)
	} else {
		bindInfoTemp = terminalInfo.BindInfo
		endpointTemp = terminalInfo.Endpoint
	}
	globallogger.Log.Infof("[devEUI: %v][updateTerminalBindInfo] DeviceID: %+v clusterIDInArray: %+v clusterIDOutArray: %+v Endpoint: %+v",
		terminalInfo.DevEUI, DeviceID, clusterIDInArray, clusterIDOutArray, Endpoint)
	var bindInfoArray = make([]config.BindInfo, len(endpointTemp))
	var bindInfo = config.BindInfo{}
	var bindIndex = 0
	var setData = bson.M{}
	oSet := make(map[string]interface{})
	var isEnd = false
	bindInfo.ClusterIDIn = clusterIDInArray
	bindInfo.ClusterIDOut = clusterIDOutArray
	bindInfo.ClusterID = getClusterIDByDeviceID(terminalInfo.TmnType, ProfileID, DeviceID, terminalInfo.ManufacturerName, Endpoint)
	bindInfo.SrcEndpoint = Endpoint
	bindInfo.DstAddrMode = "03"
	bindInfo.DstAddress = "ffffffffffffffff"
	bindInfo.DstEndpoint = "ff"
	if len(bindInfoTemp) > 0 {
		bindInfoArray = bindInfoTemp
	}
	if len(bindInfoArray) > 0 {
		for index, item := range endpointTemp {
			if Endpoint == item {
				bindIndex = index
			}
		}
	}
	bindInfoArray[bindIndex] = bindInfo
	setData["bindInfo"] = bindInfoArray
	setData["updateTime"] = time.Now()
	setData["isDiscovered"] = true
	bindInfoByte, _ := json.Marshal(bindInfoArray)
	oSet["bindinfopg"] = string(bindInfoByte)
	oSet["updatetime"] = time.Now()
	oSet["isdiscovered"] = true
	isEnd = true
	for _, value := range bindInfoArray {
		if reflect.DeepEqual(value, config.BindInfo{}) {
			setData["isDiscovered"] = false
			oSet["isdiscovered"] = false
			isEnd = false
			break
		}
	}
	// globallogger.Log.Infoln("devEUI : " + devEUI + " " + "updateTerminalBindInfo setData: " + fmt.Sprintf("%+v", setData))
	var terminalInfoRes *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfoRes, err = models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, oSet)
	} else {
		terminalInfoRes, err = models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, setData)
	}
	if err == nil {
		if isEnd {
			terminalInfo.BindInfo = bindInfoArray
			bindInfoByte, _ = json.Marshal(bindInfoArray)
			terminalInfo.BindInfoPG = string(bindInfoByte)
			if terminalInfo.IsNeedBind {
				switch terminalInfo.TmnType {
				case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal1SceneSwitch005f0cf1,
					constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3,
					constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2,
					constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b,
					constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141,
					constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142,
					constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
					procNeedntBindTerminal(terminalInfo)
				default:
					go func() {
						select {
						case <-time.After(time.Second):
							publicfunction.SendBindTerminalReq(terminalInfo)
						}
					}()
				}
			} else {
				if terminalInfo.Online {
					publicfunction.TerminalOnline(devEUI, false)
					if constant.Constant.Iotware {
						if terminalInfo.IsExist {
							iotsmartspace.StateTerminalOnlineIotware(terminalInfo)
						}
					} else if constant.Constant.Iotedge {
						iotsmartspace.StateTerminalOnline(devEUI)
					}
				}
				go func() {
					select {
					case <-time.After(time.Duration(2) * time.Second):
						iotsmartspace.ActionInsertSuccess(terminalInfo)
					}
				}()
			}
		}
	}
	return terminalInfoRes, err
}

func lastlyUpdateBindInfo(devEUI string, terminalInfo config.TerminalInfo, ProfileID string, DeviceID string, numInClusters int,
	InClusterList string, numOutClusters int, OutClusterList string, Endpoint string) (*config.TerminalInfo, error) {
	var clusterIDInArray = []string{}
	var clusterIDOutArray = []string{}
	for i := 0; i < numInClusters; i++ {
		clusterIDInArray = append(clusterIDInArray, InClusterList[0:4])
		InClusterList = InClusterList[4:]
	}
	for j := 0; j < numOutClusters; j++ {
		clusterIDOutArray = append(clusterIDOutArray, OutClusterList[0:4])
		OutClusterList = OutClusterList[4:]
	}
	terminalInfoRes, err := updateTerminalBindInfo(devEUI, terminalInfo, ProfileID, DeviceID, clusterIDInArray, clusterIDOutArray, Endpoint)
	return terminalInfoRes, err
}
