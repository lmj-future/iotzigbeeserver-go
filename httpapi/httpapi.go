package httpapi

import (
	"encoding/json"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/protocol/http"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/zclmsgdown"
	"github.com/lib/pq"
)

// TestTerminalJoinData TestTerminalJoinData
type TestTerminalJoinData struct {
	DevEUI   string `json:"devEUI"`
	NwkAddr  string `json:"nwkAddr"`
	TmnType  string `json:"tmnType"`
	ACMac    string `json:"ACMac"`
	APMac    string `json:"APMac"`
	T300ID   string `json:"T300ID"`
	ModuleID string `json:"moduleID"`
}

// TestTerminalDeleteData TestTerminalDeleteData
type TestTerminalDeleteData struct {
	DevEUI string `json:"devEUI"`
}

// TerminalListData TerminalListData
type TerminalListData struct {
	TerminalList []config.TerminalInfo `json:"terminalList"`
}

// ResultMessage ResultMessage
type ResultMessage struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// DataPage DataPage
type DataPage struct {
	List       interface{} `json:"list"`
	PageNum    int         `json:"pageNum"`
	PageSize   int         `json:"pageSize"`
	TotalCount int         `json:"totalCount"`
}

func procGetTerminalList(params http.Params, reqBody []byte) ([]byte, error) {
	defer func() {
		err := recover()
		if err != nil {
			globallogger.Log.Errorln("procGetTerminalList err :", err)
		}
	}()
	var resultMessage ResultMessage = ResultMessage{
		Code:    0,
		Message: "success",
	}
	var pageNum = 1
	var pageSize = 10
	var totalCount = 0
	if _, ok := params["pageNum"]; ok {
		pageNum, _ = strconv.Atoi(params["pageNum"])
	}
	if _, ok := params["pageSize"]; ok {
		pageSize, _ = strconv.Atoi(params["pageSize"])
	}
	dataPage := DataPage{
		PageNum:    pageNum,
		PageSize:   pageSize,
		TotalCount: totalCount,
	}
	var oMatch = map[string]interface{}{}
	var oMatchPG = map[string]interface{}{}
	var orMatchPG = map[string]interface{}{}
	oMatch["firstAddr"] = params["firstAddr"]
	oMatch["secondAddr"] = params["secondAddr"]
	oMatch["isExist"] = false
	oMatch["isDiscovered"] = true
	oMatchPG["firstaddr"] = params["firstAddr"]
	oMatchPG["secondaddr"] = params["secondAddr"]
	oMatchPG["isexist"] = false
	oMatchPG["isdiscovered"] = true
	orMatchPG["firstaddr"] = params["firstAddr"]
	orMatchPG["secondaddr"] = params["secondAddr"]
	orMatchPG["udpversion"] = constant.Constant.UDPVERSION.Version0101
	orMatchPG["isexist"] = false
	orMatchPG["isdiscovered"] = true
	if _, ok := params["thirdAddr"]; ok && len(params["thirdAddr"]) > 0 {
		oMatch["thirdAddr"] = strings.Replace(params["thirdAddr"], "-", "", -1)
		oMatchPG["thirdaddr"] = strings.Replace(params["thirdAddr"], "-", "", -1)
	}
	var field = map[string]interface{}{}
	field["devEUI"] = 1
	field["profileID"] = 1
	field["manufacturerName"] = 1
	field["tmnType"] = 1
	field["online"] = 1
	field["leaveState"] = 1
	field["firstAddr"] = 1
	field["secondAddr"] = 1
	field["thirdAddr"] = 1
	field["isExist"] = 1
	field["createTime"] = 1
	var terminalList []config.TerminalInfo2
	var terminalListPG []config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalListPG, totalCount, err = models.FindTerminalListByPagePG(oMatchPG, orMatchPG, []string{
			"deveui", "profileid", "manufacturername", "tmntype", "online", "leavestate", "firstaddr", "secondaddr", "thirdaddr", "isexist", "createtime",
		}, pageNum, pageSize)
		if len(terminalListPG) > 0 {
			terminalList = make([]config.TerminalInfo2, len(terminalListPG))
			for k, v := range terminalListPG {
				terminalList[k].DevEUI = v.DevEUI
				terminalList[k].ProfileID = v.ProfileID
				terminalList[k].TmnType = v.TmnType
				terminalList[k].Online = v.Online
				terminalList[k].LeaveState = v.LeaveState
				terminalList[k].FirstAddr = v.FirstAddr
				terminalList[k].SecondAddr = v.SecondAddr
				terminalList[k].ThirdAddr = v.ThirdAddr
				terminalList[k].IsExist = v.IsExist
				terminalList[k].CreateTime = v.CreateTime
				switch v.ManufacturerName {
				case constant.Constant.MANUFACTURERNAME.HeiMan:
					terminalList[k].ManufacturerName = "海曼"
				case constant.Constant.MANUFACTURERNAME.Honyar:
					terminalList[k].ManufacturerName = "鸿雁"
				}
			}
		}
		dataPage.TotalCount = totalCount
	} else {
		terminalList, err = models.FindTerminalListByCondition(oMatch, field)
	}
	if err != nil {
		globallogger.Log.Errorln("[HTTP][GET][/iot/iotzigbeeserver/terminalList] proc FindTerminalListByCondition err :", err)
		resultMessage.Code = 404
		resultMessage.Message = err.Error()
	} else {
		dataPage.List = terminalList
		resultMessage.Data = dataPage
	}
	jsonData, err := json.Marshal(resultMessage)
	if err != nil {
		return []byte(err.Error()), nil
	}
	return jsonData, nil
}
func procPostTerminalList(params http.Params, reqBody []byte) ([]byte, error) {
	defer func() {
		err := recover()
		if err != nil {
			globallogger.Log.Errorln("procPostTerminalList err :", err)
		}
	}()
	var resultMessage ResultMessage = ResultMessage{
		Code:    0,
		Message: "success",
	}
	var reqBodyJSONData []string
	err := json.Unmarshal(reqBody, &reqBodyJSONData)
	if err != nil {
		return []byte(err.Error()), nil
	}
	for _, devEUI := range reqBodyJSONData {
		var terminalInfo *config.TerminalInfo
		var err error
		if constant.Constant.UsePostgres {
			terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
		} else {
			terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
		}
		if err == nil && terminalInfo != nil {
			iotsmartspace.StateTerminalOnlineIotware(*terminalInfo)
		}
	}
	ticker := time.NewTicker(time.Millisecond * 500)
	for index := 0; index < 5; index++ {
		<-ticker.C
		for _, devEUI := range reqBodyJSONData {
			var terminalInfo *config.TerminalInfo
			var err error
			if constant.Constant.UsePostgres {
				terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
			} else {
				terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
			}
			if err != nil {
				break
			}
			if terminalInfo == nil {
				resultMessage.Code = 1
				resultMessage.Message = "设备[" + devEUI + "]加入失败"
				break
			}
			if !terminalInfo.IsExist {
				resultMessage.Code = 2
				resultMessage.Message = "设备[" + devEUI + "]加入失败"
				break
			} else {
				resultMessage.Code = 0
				resultMessage.Message = "success"
			}
		}
		if err != nil {
			ticker.Stop()
			return []byte(err.Error()), nil
		}
		if resultMessage.Code == 0 || resultMessage.Code == 1 {
			ticker.Stop()
			break
		}
	}
	if resultMessage.Code == 2 {
		resultMessage.Code = 1
	}
	jsonData, err := json.Marshal(resultMessage)
	if err != nil {
		return []byte(err.Error()), nil
	}
	return jsonData, nil
}
func procDeleteTerminalList(params http.Params, reqBody []byte) ([]byte, error) {
	defer func() {
		err := recover()
		if err != nil {
			globallogger.Log.Errorln("procDeleteTerminalList err :", err)
		}
	}()
	var resultMessage ResultMessage = ResultMessage{
		Code:    0,
		Message: "success",
	}
	var reqBodyJSONData []string
	err := json.Unmarshal(reqBody, &reqBodyJSONData)
	if err != nil {
		return []byte(err.Error()), nil
	}
	for _, devEUI := range reqBodyJSONData {
		publicfunction.ProcTerminalDelete(devEUI)
	}
	jsonData, err := json.Marshal(resultMessage)
	if err != nil {
		return []byte(err.Error()), nil
	}
	return jsonData, nil
}

func procTestTerminalDelete(params http.Params, reqBody []byte) ([]byte, error) {
	var jsonData TestTerminalDeleteData
	err := json.Unmarshal(reqBody, &jsonData)
	if err != nil {
		return []byte(err.Error()), nil
	}
	if jsonData.DevEUI == "" {
		return []byte("Body params is invalid, please check!"), nil
	}
	var terminalInfo *config.TerminalInfo
	if constant.Constant.UsePostgres {
		terminalInfo, _ = models.GetTerminalInfoByDevEUIPG(jsonData.DevEUI)
	} else {
		terminalInfo, _ = models.GetTerminalInfoByDevEUI(jsonData.DevEUI)
	}
	if terminalInfo != nil {
		var key = terminalInfo.APMac + terminalInfo.ModuleID + terminalInfo.NwkAddr
		publicfunction.DeleteTerminalInfoListCache(key)
	}
	if constant.Constant.UsePostgres {
		err = models.DeleteTerminalPG(jsonData.DevEUI)
	} else {
		err = models.DeleteTerminal(jsonData.DevEUI)
	}
	if err != nil {
		return []byte(err.Error()), nil
	}
	iotsmartspace.ActionTerminalLeave(jsonData.DevEUI)
	return []byte("success"), nil
}
func procTestTerminalJoin(params http.Params, reqBody []byte) ([]byte, error) {
	var jsonData TestTerminalJoinData
	err := json.Unmarshal(reqBody, &jsonData)
	if err != nil {
		return []byte(err.Error()), nil
	}
	if jsonData.DevEUI == "" || jsonData.NwkAddr == "" || jsonData.TmnType == "" || jsonData.ACMac == "" ||
		jsonData.APMac == "" || jsonData.T300ID == "" || jsonData.ModuleID == "" {
		return []byte("Body params is invalid, please check!"), nil
	}
	terminalInfo := config.TerminalInfo{
		DevEUI:       jsonData.DevEUI,
		NwkAddr:      jsonData.NwkAddr,
		TmnType:      jsonData.TmnType,
		ACMac:        jsonData.ACMac,
		FirstAddr:    publicfunction.Transport16StringToString(jsonData.ACMac),
		APMac:        jsonData.APMac,
		SecondAddr:   publicfunction.Transport16StringToString(jsonData.APMac),
		T300ID:       jsonData.T300ID,
		ThirdAddr:    jsonData.T300ID,
		ModuleID:     jsonData.ModuleID,
		ProfileID:    "ZHA",
		Online:       true,
		IsDiscovered: true,
		IsExist:      true,
		IsNeedBind:   true,
		IsReadBasic:  true,
	}
	switch jsonData.TmnType {
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
		return []byte("invalid tmnType: " + jsonData.TmnType), nil
	}
	if constant.Constant.UsePostgres {
		err = models.CreateTerminalPG(terminalInfo)
	} else {
		err = models.CreateTerminal(terminalInfo)
	}
	if err != nil {
		return []byte(err.Error()), nil
	}
	if constant.Constant.Iotware {
		iotsmartspace.StateTerminalOnlineIotware(terminalInfo)
	} else if constant.Constant.Iotedge {
		iotsmartspace.ActionInsertReplyTest(terminalInfo)
	}
	return []byte("success"), nil
}

// heimanIRControlEMSendKeyCommand 处理clusterID 0xfc82
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

// heimanIRControlEMStudyKey 处理clusterID 0xfc82
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

// heimanIRControlEMDeleteKey 处理clusterID 0xfc82
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

// heimanIRControlEMCreateID 处理clusterID 0xfc82
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

// heimanIRControlEMGetIDAndKeyCodeList 处理clusterID 0xfc82
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

// NewServer NewServer
func NewServer() {
	s, _ := http.NewServer(http.ServerInfo{Address: "tcp://" + "0.0.0.0" + ":" + constant.Constant.HTTPPort}, func(u, p string) bool {
		return u == "u" && p == "p"
	})
	defer s.Close()

	s.Handle(func(params http.Params, reqBody []byte) ([]byte, error) {
		return procGetTerminalList(params, reqBody)
	}, "GET", "/iot/iotzigbeeserver/terminalList", "firstAddr", "{firstAddr}", "secondAddr", "{secondAddr}", "thirdAddr", "{thirdAddr}")
	s.Handle(func(params http.Params, reqBody []byte) ([]byte, error) {
		return procGetTerminalList(params, reqBody)
	}, "GET", "/iot/iotzigbeeserver/terminalList", "firstAddr", "{firstAddr}", "secondAddr", "{secondAddr}")
	s.Handle(func(params http.Params, reqBody []byte) ([]byte, error) {
		return procPostTerminalList(params, reqBody)
	}, "POST", "/iot/iotzigbeeserver/terminalList")
	s.Handle(func(params http.Params, reqBody []byte) ([]byte, error) {
		return procDeleteTerminalList(params, reqBody)
	}, "POST", "/iot/iotzigbeeserver/terminalList/delete")
	s.Handle(func(params http.Params, reqBody []byte) ([]byte, error) {
		return procDeleteTerminalList(params, reqBody)
	}, "DELETE", "/iot/iotzigbeeserver/terminalList")
	s.Handle(func(params http.Params, reqBody []byte) ([]byte, error) {
		globallogger.SetLogLevel(params["level"])
		return []byte{0x4f, 0x4b}, nil
	}, "GET", "/iot/iotzigbeeserver/logLevel", "level", "{level}")
	s.Handle(func(params http.Params, reqBody []byte) ([]byte, error) {
		heimanIRControlEMSendKeyCommand(params["devEUI"], params["ID"], params["KeyCode"], 1)
		return []byte{0x4f, 0x4b}, nil
	}, "GET", "/iot/iotzigbeeserver/sendKeyCommand", "devEUI", "{devEUI}", "ID", "{ID}", "KeyCode", "{KeyCode}")
	s.Handle(func(params http.Params, reqBody []byte) ([]byte, error) {
		heimanIRControlEMStudyKey(params["devEUI"], params["ID"], params["KeyCode"], 1)
		return []byte{0x4f, 0x4b}, nil
	}, "GET", "/iot/iotzigbeeserver/studyKey", "devEUI", "{devEUI}", "ID", "{ID}", "KeyCode", "{KeyCode}")
	s.Handle(func(params http.Params, reqBody []byte) ([]byte, error) {
		heimanIRControlEMDeleteKey(params["devEUI"], params["ID"], params["KeyCode"], 1)
		return []byte{0x4f, 0x4b}, nil
	}, "GET", "/iot/iotzigbeeserver/deleteKey", "devEUI", "{devEUI}", "ID", "{ID}", "KeyCode", "{KeyCode}")
	s.Handle(func(params http.Params, reqBody []byte) ([]byte, error) {
		heimanIRControlEMCreateID(params["devEUI"], params["ModelType"], 1)
		return []byte{0x4f, 0x4b}, nil
	}, "GET", "/iot/iotzigbeeserver/createID", "devEUI", "{devEUI}", "ModelType", "{ModelType}")
	s.Handle(func(params http.Params, reqBody []byte) ([]byte, error) {
		heimanIRControlEMGetIDAndKeyCodeList(params["devEUI"], 1)
		return []byte{0x4f, 0x4b}, nil
	}, "GET", "/iot/iotzigbeeserver/getIDAndKeyCodeList", "devEUI", "{devEUI}")
	s.Handle(func(params http.Params, reqBody []byte) ([]byte, error) {
		return reqBody[:2], nil
	}, "PUT", "/h3c-zigbeeserver/test/put", "arg", "{arg}")
	s.Handle(func(params http.Params, reqBody []byte) ([]byte, error) {
		return procTestTerminalDelete(params, reqBody)
	}, "DELETE", "/h3c-zigbeeserver/test/terminalDelete")
	s.Handle(func(params http.Params, reqBody []byte) ([]byte, error) {
		return procTestTerminalJoin(params, reqBody)
	}, "POST", "/h3c-zigbeeserver/test/terminalJoin")
	s.Handle(func(params http.Params, reqBody []byte) ([]byte, error) {
		return procTestTerminalJoin(params, reqBody)
	}, "POST", "/iot/iotzigbeeserver/test/terminalJoin")
	s.HandlePrometheus("GET", "/iot/iotzigbeeserver/prometheus")

	s.Start()
	wait()
}

// wait wait
func wait() error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	signal.Ignore(syscall.SIGPIPE)
	<-sig
	return nil
}
