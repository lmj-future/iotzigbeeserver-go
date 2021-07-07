package iotsmartspace

import (
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
	"github.com/h3c/iotzigbeeserver-go/zcl/zclmsgdown"
	"github.com/lib/pq"
)

// genOnOffSwitchCommandDown 处理clusterID 0x0006
func genOnOffSwitchCommandDown(devEUI string, value string, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][genOnOffSwitchCommandDown] value: %s", devEUI, value)
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][genOnOffSwitchCommandDown] GetTerminalInfoByDevEUI err: %+v", devEUI, err)
		return
	}
	if terminalInfo != nil {
		var endpointTemp pq.StringArray
		if constant.Constant.UsePostgres {
			endpointTemp = terminalInfo.EndpointPG
		} else {
			endpointTemp = terminalInfo.Endpoint
		}
		var dstEndpointIndex int
		for i, v := range endpointTemp {
			if v == "0"+string([]byte(value)[:1]) {
				dstEndpointIndex = i
			}
		}
		var commandID uint8
		if string([]byte(value)[1:2]) == "0" {
			commandID = 0x00
		} else {
			commandID = 0x01
		}
		zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
			MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
			DevEUI:      devEUI,
			CommandType: common.SwitchCommand,
			ClusterID:   uint16(cluster.OnOff),
			Command: common.Command{
				Cmd:              commandID,
				DstEndpointIndex: dstEndpointIndex,
			},
			MsgID: msgID,
		})
	}
}

// honyarSocketCommandDown 处理鸿雁sokcet command down
func honyarSocketCommandDown(devEUI string, value string, msgID interface{}, commandType string) {
	globallogger.Log.Infof("[devEUI: %v][honyarSocketCommandDown] value: %s", devEUI, value)
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][honyarSocketCommandDown] GetTerminalInfoByDevEUI err: %+v", devEUI, err)
		return
	}
	if terminalInfo != nil {
		var endpointTemp pq.StringArray
		if constant.Constant.UsePostgres {
			endpointTemp = terminalInfo.EndpointPG
		} else {
			endpointTemp = terminalInfo.Endpoint
		}
		var dstEndpointIndex int
		for i, v := range endpointTemp {
			if v == "0"+string([]byte(value)[:1]) {
				dstEndpointIndex = i
			}
		}
		cmd := common.Command{
			DstEndpointIndex: dstEndpointIndex,
		}
		var clusterID uint16
		switch commandType {
		case common.ChildLock:
			clusterID = uint16(cluster.OnOff)
			cmd.Value = string([]byte(value)[1:2]) == "1"
			cmd.AttributeID = 0x8000
			cmd.AttributeName = "ChildLock"
		case common.OnOffUSB:
			clusterID = uint16(cluster.OnOff)
			if string([]byte(value)[1:2]) == "0" {
				cmd.Cmd = 0x00
			} else {
				cmd.Cmd = 0x01
			}
			cmd.DstEndpointIndex = 1
		case common.BackGroundLight:
			clusterID = uint16(cluster.OnOff)
			cmd.Value = string([]byte(value)[1:2]) == "1"
			cmd.AttributeID = 0x8005
			cmd.AttributeName = "BackGroundLight"
		case common.PowerOffMemory:
			clusterID = uint16(cluster.OnOff)
			cmd.Value = string([]byte(value)[1:2]) == "0"
			cmd.AttributeID = 0x8006
			cmd.AttributeName = "PowerOffMemory"
		case common.HistoryElectricClear:
			clusterID = uint16(cluster.SmartEnergyMetering)
			cmd.Value = string([]byte(value)[1:2]) == "1"
			cmd.AttributeID = 0x8000
			cmd.AttributeName = "HistoryElectricClear"
		}
		if commandType == common.OnOffUSB {
			commandType = common.SocketCommand
		}
		zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
			MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
			DevEUI:      devEUI,
			CommandType: commandType,
			ClusterID:   clusterID,
			Command:     cmd,
			MsgID:       msgID,
		})
	}
}

// heimanWarningDeviceProcAlarm 处理clusterID 0x0502
func heimanWarningDeviceProcAlarm(devEUI string, alarmTime string, alarmStart string, warningMode float64, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanWarningDeviceProcAlarm] alarmTime: %s, alarmStart: %s, warningMode: %v", devEUI, alarmTime, alarmStart, warningMode)
	time, err := strconv.ParseUint(alarmTime, 10, 16)
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanWarningDeviceProcAlarm] ParseUint time err: %+v", devEUI, err)
	}
	var warningControl uint8 = 0x03
	if alarmStart == "0" {
		warningControl = 0x17
		switch warningMode {
		case 1: // Burglar
			warningControl = 0x17
		case 2: // Fire
			warningControl = 0x27
		case 3: // Emergency
			warningControl = 0x37
		case 4: // Police panic
			warningControl = 0x47
		case 5: // Fire panic
			warningControl = 0x57
		case 6: // Emergency panic
			warningControl = 0x67
		}
	}
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.StartWarning,
		ClusterID:   0x0502,
		Command: common.Command{
			DstEndpointIndex: 0,
			WarningControl:   warningControl,
			WarningTime:      uint16(time),
		},
		MsgID: msgID,
	})
}

// heimanWarningDeviceProcArmed 处理clusterID 0x0502
func heimanWarningDeviceProcArmed(devEUI string, warningMode float64, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanWarningDeviceProcArmed] warningMode: %v", devEUI, warningMode)
	var squawkMode uint8
	switch warningMode {
	case 0: // armed
		squawkMode = 0x03
	case 1: // disarmed
		squawkMode = 0x13
	}
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.Squawk,
		ClusterID:   0x0502,
		Command: common.Command{
			DstEndpointIndex: 0,
			SquawkMode:       squawkMode,
		},
		MsgID: msgID,
	})
}

// heimanTemperatureAlarmThresholdSet 处理clusterID 0xfc81
func heimanTemperatureAlarmThresholdSet(devEUI string, value string, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanTemperatureAlarmThresholdSet] value: %s", devEUI, value)
	valueArr := strings.Split(value, ",")
	if len(valueArr) == 1 {
		valueArr = strings.Split(value, "，")
	}
	minValue, err := strconv.ParseInt(valueArr[0], 10, 16)
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanTemperatureAlarmThresholdSet] ParseInt minValue err: %+v", devEUI, err)
	}
	maxValue, err := strconv.ParseInt(valueArr[1], 10, 16)
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanTemperatureAlarmThresholdSet] ParseInt maxValue err: %+v", devEUI, err)
	}
	minValue = minValue * 100
	maxValue = maxValue * 100
	if minValue < 0 {
		minValue = minValue + 65536
	}
	if maxValue < 0 {
		maxValue = maxValue + 65536
	}
	var terminalInfo *config.TerminalInfo
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanTemperatureAlarmThresholdSet] GetTerminalInfoByDevEUI err: %+v", devEUI, err)
	}
	var disturb uint16 = 0x00f0
	if terminalInfo.Disturb == "80c0" {
		disturb = 0x80c0
	}
	var endpointTemp pq.StringArray
	if constant.Constant.UsePostgres {
		endpointTemp = terminalInfo.EndpointPG
	} else {
		endpointTemp = terminalInfo.Endpoint
	}
	endpointIndex := 0
	if len(endpointTemp) > 1 && endpointTemp[1] == "01" {
		endpointIndex = 1
	}
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SetThreshold,
		ClusterID:   0xfc81,
		Command: common.Command{
			DstEndpointIndex: endpointIndex,
			MaxTemperature:   int16(maxValue),
			MinTemperature:   int16(minValue),
			Disturb:          disturb,
		},
		MsgID: msgID,
	})
}

// heimanHumidityAlarmThresholdSet 处理clusterID 0xfc81
func heimanHumidityAlarmThresholdSet(devEUI string, value string, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanHumidityAlarmThresholdSet] value: %s", devEUI, value)
	valueArr := strings.Split(value, ",")
	if len(valueArr) == 1 {
		valueArr = strings.Split(value, "，")
	}
	minValue, err := strconv.ParseUint(valueArr[0], 10, 16)
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanHumidityAlarmThresholdSet] ParseUint minValue err: %+v", devEUI, err)
	}
	maxValue, err := strconv.ParseUint(valueArr[1], 10, 16)
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanHumidityAlarmThresholdSet] ParseUint maxValue err: %+v", devEUI, err)
	}
	minValue = minValue * 100
	maxValue = maxValue * 100
	var terminalInfo *config.TerminalInfo
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanHumidityAlarmThresholdSet] GetTerminalInfoByDevEUI err: %+v", devEUI, err)
	}
	var disturb uint16 = 0x00f0
	if terminalInfo.Disturb == "80c0" {
		disturb = 0x80c0
	}
	var endpointTemp pq.StringArray
	if constant.Constant.UsePostgres {
		endpointTemp = terminalInfo.EndpointPG
	} else {
		endpointTemp = terminalInfo.Endpoint
	}
	endpointIndex := 0
	if len(endpointTemp) > 1 && endpointTemp[1] == "01" {
		endpointIndex = 1
	}
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SetThreshold,
		ClusterID:   0xfc81,
		Command: common.Command{
			DstEndpointIndex: endpointIndex,
			MaxHumidity:      uint16(maxValue),
			MinHumidity:      uint16(minValue),
			Disturb:          disturb,
		},
		MsgID: msgID,
	})
}

// heimanDisturbAlarmThresholdSet 处理clusterID 0xfc81
func heimanDisturbAlarmThresholdSet(devEUI string, value string, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanDisturbAlarmThresholdSet] value: %s", devEUI, value)
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanDisturbAlarmThresholdSet] GetTerminalInfoByDevEUI err: %+v", devEUI, err)
	}
	var endpointTemp pq.StringArray
	if constant.Constant.UsePostgres {
		endpointTemp = terminalInfo.EndpointPG
	} else {
		endpointTemp = terminalInfo.Endpoint
	}
	endpointIndex := 0
	if len(endpointTemp) > 1 && endpointTemp[1] == "01" {
		endpointIndex = 1
	}
	var disturb uint16 = 0x00f0
	if value == "1" { // 开启勿扰模式
		disturb = 0x80c0
	}
	var maxTemperature int64 = 6000
	var minTemperature int64 = -1000
	var maxHumidity uint64 = 10000
	var minHumidity uint64 = 0
	if terminalInfo.MaxTemperature != "" {
		maxTemperature, _ = strconv.ParseInt(terminalInfo.MaxTemperature, 10, 16)
	}
	if terminalInfo.MinTemperature != "" {
		minTemperature, _ = strconv.ParseInt(terminalInfo.MinTemperature, 10, 16)
	}
	if terminalInfo.MaxHumidity != "" {
		maxHumidity, _ = strconv.ParseUint(terminalInfo.MaxHumidity, 10, 16)
	}
	if terminalInfo.MinHumidity != "" {
		minHumidity, _ = strconv.ParseUint(terminalInfo.MinHumidity, 10, 16)
	}
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SetThreshold,
		ClusterID:   0xfc81,
		Command: common.Command{
			DstEndpointIndex: endpointIndex,
			MaxTemperature:   int16(maxTemperature),
			MinTemperature:   int16(minTemperature),
			MaxHumidity:      uint16(maxHumidity),
			MinHumidity:      uint16(minHumidity),
			Disturb:          disturb,
		},
		MsgID: msgID,
	})
}

// heimanIRControlEMSendKeyCommand 处理clusterID 0xfc82
func heimanIRControlEMSendKeyCommand(devEUI string, value cluster.HEIMANInfraredRemoteSendKeyCommand, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanIRControlEMSendKeyCommand] value: %+v", devEUI, value)
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SendKeyCommand,
		ClusterID:   0xfc82,
		Command: common.Command{
			InfraredRemoteID:      value.ID,
			InfraredRemoteKeyCode: value.KeyCode,
		},
		MsgID: msgID,
	})
}

// heimanIRControlEMStudyKey 处理clusterID 0xfc82
func heimanIRControlEMStudyKey(devEUI string, value cluster.HEIMANInfraredRemoteStudyKey, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanIRControlEMStudyKey] value: %+v", devEUI, value)
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.StudyKey,
		ClusterID:   0xfc82,
		Command: common.Command{
			InfraredRemoteID:      value.ID,
			InfraredRemoteKeyCode: value.KeyCode,
		},
		MsgID: msgID,
	})
}

// heimanIRControlEMDeleteKey 处理clusterID 0xfc82
func heimanIRControlEMDeleteKey(devEUI string, value cluster.HEIMANInfraredRemoteDeleteKey, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanIRControlEMDeleteKey] value: %+v", devEUI, value)
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.DeleteKey,
		ClusterID:   0xfc82,
		Command: common.Command{
			InfraredRemoteID:      value.ID,
			InfraredRemoteKeyCode: value.KeyCode,
		},
		MsgID: msgID,
	})
}

// heimanIRControlEMCreateID 处理clusterID 0xfc82
func heimanIRControlEMCreateID(devEUI string, value cluster.HEIMANInfraredRemoteCreateID, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanIRControlEMCreateID] value: %+v", devEUI, value)
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.CreateID,
		ClusterID:   0xfc82,
		Command: common.Command{
			InfraredRemoteModelType: value.ModelType,
		},
		MsgID: msgID,
	})
}

// heimanIRControlEMGetIDAndKeyCodeList 处理clusterID 0xfc82
func heimanIRControlEMGetIDAndKeyCodeList(devEUI string, value cluster.HEIMANInfraredRemoteGetIDAndKeyCodeList, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanIRControlEMGetIDAndKeyCodeList] value: %+v", devEUI, value)
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.GetIDAndKeyCodeList,
		ClusterID:   0xfc82,
		MsgID:       msgID,
	})
}

// heimanHS2AQSetLanguage 处理clusterID 0xfc81
func heimanHS2AQSetLanguage(devEUI string, language uint8, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanHS2AQSetLanguage] language: %+v", devEUI, language)
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SetLanguage,
		ClusterID:   0xfc81,
		Command: common.Command{
			Cmd: language,
		},
		MsgID: msgID,
	})
}

// heimanHS2AQSetUnitOfTemperature 处理clusterID 0xfc81
func heimanHS2AQSetUnitOfTemperature(devEUI string, unitOfTemperature uint8, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanHS2AQSetUnitOfTemperature] unitOfTemperature: %+v", devEUI, unitOfTemperature)
	unitoftemperature := "C"
	if unitOfTemperature == 0 {
		unitoftemperature = "F"
	}
	var terminalInfo *config.TerminalInfo = nil
	if constant.Constant.UsePostgres {
		terminalInfo, _ = models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, map[string]interface{}{"unitoftemperature": unitoftemperature})
	} else {
		terminalInfo, _ = models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, bson.M{"unitOfTemperature": unitoftemperature})
	}
	if terminalInfo != nil {
		var keyBuilder strings.Builder
		if terminalInfo.UDPVersion == constant.Constant.UDPVERSION.Version0102 {
			keyBuilder.WriteString(terminalInfo.APMac)
			keyBuilder.WriteString(terminalInfo.ModuleID)
			keyBuilder.WriteString(terminalInfo.NwkAddr)
		} else {
			keyBuilder.WriteString(terminalInfo.APMac)
			keyBuilder.WriteString(terminalInfo.ModuleID)
			keyBuilder.WriteString(terminalInfo.DevEUI)
		}
		publicfunction.DeleteTerminalInfoListCache(keyBuilder.String())
	}
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SetUnitOfTemperature,
		ClusterID:   0xfc81,
		Command: common.Command{
			Cmd: unitOfTemperature,
		},
		MsgID: msgID,
	})
}

// procPropertyIotware 处理属性消息
func procPropertyIotware(mqttMsg publicstruct.RPCIotware) {
	globallogger.Log.Infof("[procPropertyIotware]: mqttMsg.Data.Params: %+v", mqttMsg.Data.Params)
	switch mqttMsg.Data.Method {
	case IotwarePropertyMethodPowerSwitch:
		if mqttMsg.Data.Params.PowerSwitch {
			genOnOffSwitchCommandDown(mqttMsg.Device, "11", mqttMsg.Data.ID)
		} else {
			genOnOffSwitchCommandDown(mqttMsg.Device, "10", mqttMsg.Data.ID)
		}
	case IotwarePropertyMethodPowerSwitch_1:
		if mqttMsg.Data.Params.PowerSwitch_1 {
			genOnOffSwitchCommandDown(mqttMsg.Device, "11", mqttMsg.Data.ID)
		} else {
			genOnOffSwitchCommandDown(mqttMsg.Device, "10", mqttMsg.Data.ID)
		}
	case IotwarePropertyMethodPowerSwitch_2:
		if mqttMsg.Data.Params.PowerSwitch_2 {
			genOnOffSwitchCommandDown(mqttMsg.Device, "21", mqttMsg.Data.ID)
		} else {
			genOnOffSwitchCommandDown(mqttMsg.Device, "20", mqttMsg.Data.ID)
		}
	case IotwarePropertyMethodPowerSwitch_3:
		if mqttMsg.Data.Params.PowerSwitch_3 {
			genOnOffSwitchCommandDown(mqttMsg.Device, "31", mqttMsg.Data.ID)
		} else {
			genOnOffSwitchCommandDown(mqttMsg.Device, "30", mqttMsg.Data.ID)
		}
	case IotwarePropertyMethodActiveAlarm:
		heimanWarningDeviceProcAlarm(mqttMsg.Device, strconv.FormatUint(uint64(mqttMsg.Data.Params.Time), 10), "0", mqttMsg.Data.Params.WarningMode, mqttMsg.Data.ID)
	case IotwarePropertyMethodClearAlarm:
		heimanWarningDeviceProcAlarm(mqttMsg.Device, "1", "1", 0, mqttMsg.Data.ID)
	case IotwarePropertyMethodArmed:
		heimanWarningDeviceProcArmed(mqttMsg.Device, 0, mqttMsg.Data.ID)
	case IotwarePropertyMethodDisarmed:
		heimanWarningDeviceProcArmed(mqttMsg.Device, 1, mqttMsg.Data.ID)
	case IotwarePropertyMethodTemperatureThreshold:
		heimanTemperatureAlarmThresholdSet(mqttMsg.Device,
			strconv.FormatUint(uint64(mqttMsg.Data.Params.Lower), 10)+","+strconv.FormatUint(uint64(mqttMsg.Data.Params.Upper), 10), mqttMsg.Data.ID)
	case IotwarePropertyMethodHumidityThreshold:
		heimanHumidityAlarmThresholdSet(mqttMsg.Device,
			strconv.FormatUint(uint64(mqttMsg.Data.Params.Lower), 10)+","+strconv.FormatUint(uint64(mqttMsg.Data.Params.Upper), 10), mqttMsg.Data.ID)
	case IotwarePropertyMethodNoDisturb:
		if mqttMsg.Data.Params.NoDisturb {
			heimanDisturbAlarmThresholdSet(mqttMsg.Device, "1", mqttMsg.Data.ID)
		} else {
			heimanDisturbAlarmThresholdSet(mqttMsg.Device, "0", mqttMsg.Data.ID)
		}
	case IotwarePropertyMethodDownControl:
		heimanIRControlEMSendKeyCommand(mqttMsg.Device,
			cluster.HEIMANInfraredRemoteSendKeyCommand{ID: uint8(mqttMsg.Data.Params.IRId), KeyCode: uint8(mqttMsg.Data.Params.KeyCode)}, mqttMsg.Data.ID)
	case IotwarePropertyMethodLearnControl:
		heimanIRControlEMStudyKey(mqttMsg.Device,
			cluster.HEIMANInfraredRemoteStudyKey{ID: uint8(mqttMsg.Data.Params.IRId), KeyCode: uint8(mqttMsg.Data.Params.KeyCode)}, mqttMsg.Data.ID)
	case IotwarePropertyMethodDeleteControl:
		heimanIRControlEMDeleteKey(mqttMsg.Device,
			cluster.HEIMANInfraredRemoteDeleteKey{ID: uint8(mqttMsg.Data.Params.IRId), KeyCode: uint8(mqttMsg.Data.Params.KeyCode)}, mqttMsg.Data.ID)
	case IotwarePropertyMethodCreateModelType:
		heimanIRControlEMCreateID(mqttMsg.Device,
			cluster.HEIMANInfraredRemoteCreateID{ModelType: uint8(mqttMsg.Data.Params.ModelType)}, mqttMsg.Data.ID)
	case IotwarePropertyMethodGetDevAndCmdList:
		heimanIRControlEMGetIDAndKeyCodeList(mqttMsg.Device, cluster.HEIMANInfraredRemoteGetIDAndKeyCodeList{}, mqttMsg.Data.ID)
	case IotwarePropertyMethodChildLockSwitch:
		if mqttMsg.Data.Params.ChildLockSwitch {
			honyarSocketCommandDown(mqttMsg.Device, "11", mqttMsg.Data.ID, common.ChildLock)
		} else {
			honyarSocketCommandDown(mqttMsg.Device, "10", mqttMsg.Data.ID, common.ChildLock)
		}
	case IotwarePropertyMethodUSBSwitch_1:
		if mqttMsg.Data.Params.USBSwitch_1 {
			honyarSocketCommandDown(mqttMsg.Device, "11", mqttMsg.Data.ID, common.OnOffUSB)
		} else {
			honyarSocketCommandDown(mqttMsg.Device, "10", mqttMsg.Data.ID, common.OnOffUSB)
		}
	case IotwarePropertyMethodBackLightSwitch:
		if mqttMsg.Data.Params.BackLightSwitch {
			honyarSocketCommandDown(mqttMsg.Device, "11", mqttMsg.Data.ID, common.BackGroundLight)
		} else {
			honyarSocketCommandDown(mqttMsg.Device, "10", mqttMsg.Data.ID, common.BackGroundLight)
		}
	case IotwarePropertyMethodMemorySwitch:
		if mqttMsg.Data.Params.MemorySwitch {
			honyarSocketCommandDown(mqttMsg.Device, "11", mqttMsg.Data.ID, common.PowerOffMemory)
		} else {
			honyarSocketCommandDown(mqttMsg.Device, "10", mqttMsg.Data.ID, common.PowerOffMemory)
		}
	case IotwarePropertyMethodClearHistory:
		if mqttMsg.Data.Params.ClearHistory == "1" {
			honyarSocketCommandDown(mqttMsg.Device, "11", mqttMsg.Data.ID, common.HistoryElectricClear)
		} else {
			honyarSocketCommandDown(mqttMsg.Device, "10", mqttMsg.Data.ID, common.HistoryElectricClear)
		}
	case IotwarePropertyMethodSwitchLeftUpLogo:
		honyarAddSceneCommandDown(mqttMsg.Device, mqttMsg.Data.Params.SwitchLeftUpLogo, mqttMsg.Data.ID, IotwarePropertyMethodSwitchLeftUpLogo, common.AddScene)
	case IotwarePropertyMethodSwitchLeftMidLogo:
		honyarAddSceneCommandDown(mqttMsg.Device, mqttMsg.Data.Params.SwitchLeftMidLogo, mqttMsg.Data.ID, IotwarePropertyMethodSwitchLeftMidLogo, common.AddScene)
	case IotwarePropertyMethodSwitchLeftDownLogo:
		honyarAddSceneCommandDown(mqttMsg.Device, mqttMsg.Data.Params.SwitchLeftDownLogo, mqttMsg.Data.ID, IotwarePropertyMethodSwitchLeftDownLogo, common.AddScene)
	case IotwarePropertyMethodSwitchRightUpLogo:
		honyarAddSceneCommandDown(mqttMsg.Device, mqttMsg.Data.Params.SwitchRightUpLogo, mqttMsg.Data.ID, IotwarePropertyMethodSwitchRightUpLogo, common.AddScene)
	case IotwarePropertyMethodSwitchRightMidLogo:
		honyarAddSceneCommandDown(mqttMsg.Device, mqttMsg.Data.Params.SwitchRightMidLogo, mqttMsg.Data.ID, IotwarePropertyMethodSwitchRightMidLogo, common.AddScene)
	case IotwarePropertyMethodSwitchRightDownLogo:
		honyarAddSceneCommandDown(mqttMsg.Device, mqttMsg.Data.Params.SwitchRightDownLogo, mqttMsg.Data.ID, IotwarePropertyMethodSwitchRightDownLogo, common.AddScene)
	case IotwarePropertyMethodLanguage:
		if mqttMsg.Data.Params.Language == "0" {
			heimanHS2AQSetLanguage(mqttMsg.Device, 0x00, mqttMsg.Data.ID)
		} else if mqttMsg.Data.Params.Language == "1" {
			heimanHS2AQSetLanguage(mqttMsg.Device, 0x01, mqttMsg.Data.ID)
		}
	case IotwarePropertyMethodUnitOfTemperature:
		if mqttMsg.Data.Params.UnitOfTemperature == "0" {
			heimanHS2AQSetUnitOfTemperature(mqttMsg.Device, 0x01, mqttMsg.Data.ID)
		} else if mqttMsg.Data.Params.UnitOfTemperature == "1" {
			heimanHS2AQSetUnitOfTemperature(mqttMsg.Device, 0x00, mqttMsg.Data.ID)
		}
	default:
		globallogger.Log.Warnf("[procPropertyIotware]: invalid method: %s", mqttMsg.Data.Method)
	}
}

// procProperty 处理属性消息
func procProperty(mqttMsg mqttMsgSt) {
	globallogger.Log.Infof("[procProperty]: mqttMsg.Params: %+v", mqttMsg.Params)
	mapParams := mqttMsg.Params.(map[string]interface{})
	var alarmTime string
	var alarmStart string
	for key, value := range mapParams {
		if key == "terminalId" {
			continue
		}
		switch key {
		case HeimanSmartPlugPropertyKeyOnOff:
			fallthrough
		case HeimanESocketPropertyKeyOnOff, HonyarSocket000a0c3cPropertyKeyOnOff,
			HonyarSocketHY0105PropertyKeyOnOff, HonyarSocketHY0106PropertyKeyOnOff:
			globallogger.Log.Errorf("[devEUI: %v][procProperty] subscribe time: %+v", mapParams["terminalId"], time.Now())
			if value.(string) == "00" {
				genOnOffSwitchCommandDown(mapParams["terminalId"].(string), "10", mqttMsg.ID)
			} else if value.(string) == "01" {
				genOnOffSwitchCommandDown(mapParams["terminalId"].(string), "11", mqttMsg.ID)
			}
		case HeimanHS2SW1LEFR30PropertyKeyOnOff1, HonyarSingleSwitch00500c32PropertyKeyOnOff1:
			if value.(string) == "00" {
				genOnOffSwitchCommandDown(mapParams["terminalId"].(string), "10", mqttMsg.ID)
			} else if value.(string) == "01" {
				genOnOffSwitchCommandDown(mapParams["terminalId"].(string), "11", mqttMsg.ID)
			}
		case HeimanHS2SW2LEFR30PropertyKeyOnOff1, HonyarDoubleSwitch00500c33PropertyKeyOnOff1:
			if value.(string) == "00" {
				genOnOffSwitchCommandDown(mapParams["terminalId"].(string), "10", mqttMsg.ID)
			} else if value.(string) == "01" {
				genOnOffSwitchCommandDown(mapParams["terminalId"].(string), "11", mqttMsg.ID)
			}
		case HeimanHS2SW2LEFR30PropertyKeyOnOff2, HonyarDoubleSwitch00500c33PropertyKeyOnOff2:
			if value.(string) == "00" {
				genOnOffSwitchCommandDown(mapParams["terminalId"].(string), "20", mqttMsg.ID)
			} else if value.(string) == "01" {
				genOnOffSwitchCommandDown(mapParams["terminalId"].(string), "21", mqttMsg.ID)
			}
		case HeimanHS2SW3LEFR30PropertyKeyOnOff1, HonyarTripleSwitch00500c35PropertyKeyOnOff1:
			if value.(string) == "00" {
				genOnOffSwitchCommandDown(mapParams["terminalId"].(string), "10", mqttMsg.ID)
			} else if value.(string) == "01" {
				genOnOffSwitchCommandDown(mapParams["terminalId"].(string), "11", mqttMsg.ID)
			}
		case HeimanHS2SW3LEFR30PropertyKeyOnOff2, HonyarTripleSwitch00500c35PropertyKeyOnOff2:
			if value.(string) == "00" {
				genOnOffSwitchCommandDown(mapParams["terminalId"].(string), "20", mqttMsg.ID)
			} else if value.(string) == "01" {
				genOnOffSwitchCommandDown(mapParams["terminalId"].(string), "21", mqttMsg.ID)
			}
		case HeimanHS2SW3LEFR30PropertyKeyOnOff3, HonyarTripleSwitch00500c35PropertyKeyOnOff3:
			if value.(string) == "00" {
				genOnOffSwitchCommandDown(mapParams["terminalId"].(string), "30", mqttMsg.ID)
			} else if value.(string) == "01" {
				genOnOffSwitchCommandDown(mapParams["terminalId"].(string), "31", mqttMsg.ID)
			}
		case HeimanWarningDevicePropertyAlarmTime:
			alarmTime = value.(string)
			if alarmStart != "" {
				heimanWarningDeviceProcAlarm(mapParams["terminalId"].(string), alarmTime, alarmStart, 1, mqttMsg.ID)
			}
		case HeimanWarningDevicePropertyAlarmStart:
			alarmStart = value.(string)
			if alarmTime != "" {
				heimanWarningDeviceProcAlarm(mapParams["terminalId"].(string), alarmTime, alarmStart, 1, mqttMsg.ID)
			}
			if alarmStart == "1" {
				alarmTime = "1"
				heimanWarningDeviceProcAlarm(mapParams["terminalId"].(string), alarmTime, alarmStart, 0, mqttMsg.ID)
			}
		case HeimanHS2AQPropertyKeyTemperatureAlarmThreshold:
			heimanTemperatureAlarmThresholdSet(mapParams["terminalId"].(string), value.(string), mqttMsg.ID)
		case HeimanHS2AQPropertyKeyHumidityAlarmThreshold:
			heimanHumidityAlarmThresholdSet(mapParams["terminalId"].(string), value.(string), mqttMsg.ID)
		case HeimanHS2AQPropertyKeyUnDisturb:
			heimanDisturbAlarmThresholdSet(mapParams["terminalId"].(string), value.(string), mqttMsg.ID)
		case HeimanIRControlEMPropertyKeySendKeyCommand:
			heimanIRControlEMSendKeyCommand(mapParams["terminalId"].(string), value.(cluster.HEIMANInfraredRemoteSendKeyCommand), mqttMsg.ID)
		case HeimanIRControlEMPropertyKeyStudyKey:
			heimanIRControlEMStudyKey(mapParams["terminalId"].(string), value.(cluster.HEIMANInfraredRemoteStudyKey), mqttMsg.ID)
		case HeimanIRControlEMPropertyKeyDeleteKey:
			heimanIRControlEMDeleteKey(mapParams["terminalId"].(string), value.(cluster.HEIMANInfraredRemoteDeleteKey), mqttMsg.ID)
		case HeimanIRControlEMPropertyKeyCreateID:
			heimanIRControlEMCreateID(mapParams["terminalId"].(string), value.(cluster.HEIMANInfraredRemoteCreateID), mqttMsg.ID)
		case HeimanIRControlEMPropertyKeyGetIDAndKeyCodeList:
			heimanIRControlEMGetIDAndKeyCodeList(mapParams["terminalId"].(string), value.(cluster.HEIMANInfraredRemoteGetIDAndKeyCodeList), mqttMsg.ID)
		case HonyarSocket000a0c3cPropertyKeyChildLock, HonyarSocketHY0105PropertyKeyChildLock, HonyarSocketHY0106PropertyKeyChildLock:
			if value.(string) == "00" {
				honyarSocketCommandDown(mapParams["terminalId"].(string), "10", mqttMsg.ID, common.ChildLock)
			} else if value.(string) == "01" {
				honyarSocketCommandDown(mapParams["terminalId"].(string), "11", mqttMsg.ID, common.ChildLock)
			}
		case HonyarSocket000a0c3cPropertyKeyUSB, HonyarSocketHY0105PropertyKeyUSB:
			if value.(string) == "00" {
				honyarSocketCommandDown(mapParams["terminalId"].(string), "10", mqttMsg.ID, common.OnOffUSB)
			} else if value.(string) == "01" {
				honyarSocketCommandDown(mapParams["terminalId"].(string), "11", mqttMsg.ID, common.OnOffUSB)
			}
		case HonyarSocketHY0105PropertyKeyBackGroundLight, HonyarSocketHY0106PropertyKeyBackGroundLight:
			if value.(string) == "00" {
				honyarSocketCommandDown(mapParams["terminalId"].(string), "10", mqttMsg.ID, common.BackGroundLight)
			} else if value.(string) == "01" {
				honyarSocketCommandDown(mapParams["terminalId"].(string), "11", mqttMsg.ID, common.BackGroundLight)
			}
		case HonyarSocketHY0105PropertyKeyPowerOffMemory, HonyarSocketHY0106PropertyKeyPowerOffMemory:
			if value.(string) == "00" {
				honyarSocketCommandDown(mapParams["terminalId"].(string), "10", mqttMsg.ID, common.PowerOffMemory)
			} else if value.(string) == "01" {
				honyarSocketCommandDown(mapParams["terminalId"].(string), "11", mqttMsg.ID, common.PowerOffMemory)
			}
		case HonyarSocketHY0105PropertyKeyHistoryElectricClear, HonyarSocketHY0106PropertyKeyHistoryElectricClear:
			if value.(string) == "0" {
				honyarSocketCommandDown(mapParams["terminalId"].(string), "10", mqttMsg.ID, common.HistoryElectricClear)
			} else if value.(string) == "1" {
				honyarSocketCommandDown(mapParams["terminalId"].(string), "11", mqttMsg.ID, common.HistoryElectricClear)
			}
		case Honyar6SceneSwitch005f0c3bPropertyKeyLeftUpAddScene,
			Honyar6SceneSwitch005f0c3bPropertyKeyLeftMiddleAddScene,
			Honyar6SceneSwitch005f0c3bPropertyKeyLeftDownAddScene,
			Honyar6SceneSwitch005f0c3bPropertyKeyRightUpAddScene,
			Honyar6SceneSwitch005f0c3bPropertyKeyRightMiddleAddScene,
			Honyar6SceneSwitch005f0c3bPropertyKeyRightDownAddScene:
			honyarAddSceneCommandDown(mapParams["terminalId"].(string), value.(string), mqttMsg.ID, key, common.AddScene)
		default:
			globallogger.Log.Warnf("[procProperty]: invalid key: %s", key)
		}
	}
}

// procScenesCommand 处理scenes command
func procScenesCommand(devEUI string, commandType string, command common.Command, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][procScenesCommand] command: %+v", devEUI, command)
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][procScenesCommand] GetTerminalInfoByDevEUI err: %+v", devEUI, err)
	}
	var endpointTemp pq.StringArray
	if constant.Constant.UsePostgres {
		endpointTemp = terminalInfo.EndpointPG
	} else {
		endpointTemp = terminalInfo.Endpoint
	}
	endpointIndex := 0
	if len(endpointTemp) > 1 && endpointTemp[1] == "01" {
		endpointIndex = 1
	}
	command.DstEndpointIndex = endpointIndex
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: commandType,
		ClusterID:   0x0005,
		Command:     command,
		MsgID:       msgID,
	})
}

// honyarAddSceneCommandDown 处理场景消息
func honyarAddSceneCommandDown(devEUI string, value string, msgID interface{}, key string, commandType string) {
	value = strings.ToUpper(hex.EncodeToString([]byte(mahonia.NewEncoder("gbk").ConvertString(value))))
	command := common.Command{}
	command.GroupID = 1
	command.TransitionTime = 0
	command.SceneName = "0" + strconv.FormatInt(int64(len(value)/2), 16) + value
	command.KeyID = 1
	switch key {
	case IotwarePropertyMethodSwitchLeftUpLogo, Honyar6SceneSwitch005f0c3bPropertyKeyLeftUpAddScene:
		command.SceneID = 1
	case IotwarePropertyMethodSwitchLeftMidLogo, Honyar6SceneSwitch005f0c3bPropertyKeyLeftMiddleAddScene:
		command.SceneID = 2
	case IotwarePropertyMethodSwitchLeftDownLogo, Honyar6SceneSwitch005f0c3bPropertyKeyLeftDownAddScene:
		command.SceneID = 3
	case IotwarePropertyMethodSwitchRightUpLogo, Honyar6SceneSwitch005f0c3bPropertyKeyRightUpAddScene:
		command.SceneID = 4
	case IotwarePropertyMethodSwitchRightMidLogo, Honyar6SceneSwitch005f0c3bPropertyKeyRightMiddleAddScene:
		command.SceneID = 5
	case IotwarePropertyMethodSwitchRightDownLogo, Honyar6SceneSwitch005f0c3bPropertyKeyRightDownAddScene:
		command.SceneID = 6
	}
	procScenesCommand(devEUI, commandType, command, msgID)
}

// procScenes 处理场景消息
func procScenes(mqttMsg mqttMsgSt) {
	globallogger.Log.Infof("[procScenes]: mqttMsg.Params: %+v", mqttMsg.Params)
	mapParams := mqttMsg.Params.(map[string]interface{})
	command := common.Command{}
	for k, v := range mapParams {
		switch k {
		case "GroupID":
			tempGroupID, _ := strconv.ParseUint(v.(string), 10, 16)
			command.GroupID = uint16(tempGroupID)
		case "SceneID":
			tempSceneID, _ := strconv.ParseUint(v.(string), 10, 8)
			command.SceneID = uint8(tempSceneID)
		case "TransitionTime":
			tempTransitionTime, _ := strconv.ParseUint(v.(string), 10, 16)
			command.TransitionTime = uint16(tempTransitionTime)
		case "SceneName":
			command.SceneName = v.(string)
		case "KeyID":
			tempKeyID, _ := strconv.ParseUint(v.(string), 10, 8)
			command.KeyID = uint8(tempKeyID)
		case "Mode":
			tempMode, _ := strconv.ParseUint(v.(string), 10, 8)
			command.Mode = uint8(tempMode)
		case "FromGroupID":
			tempFromGroupID, _ := strconv.ParseUint(v.(string), 10, 16)
			command.FromGroupID = uint16(tempFromGroupID)
		case "FromSceneID":
			tempFromSceneID, _ := strconv.ParseUint(v.(string), 10, 16)
			command.FromSceneID = uint16(tempFromSceneID)
		case "ToGroupID":
			tempToGroupID, _ := strconv.ParseUint(v.(string), 10, 16)
			command.ToGroupID = uint16(tempToGroupID)
		case "ToSceneID":
			tempToSceneID, _ := strconv.ParseUint(v.(string), 10, 16)
			command.ToSceneID = uint16(tempToSceneID)
		}
	}
	for key := range mapParams {
		if key == "terminalId" || key == "GroupID" || key == "SceneID" || key == "TransitionTime" || key == "SceneName" ||
			key == "KeyID" || key == "Mode" || key == "FromGroupID" || key == "FromSceneID" || key == "ToGroupID" || key == "ToSceneID" {
			continue
		}
		switch key {
		case common.AddScene, common.ViewScene, common.RemoveScene, common.RemoveAllScenes, common.StoreScene,
			common.RecallScene, common.GetSceneMembership, common.EnhancedAddScene, common.EnhancedViewScene, common.CopyScene:
			procScenesCommand(mapParams["terminalId"].(string), key, command, mqttMsg.ID)
		default:
			globallogger.Log.Warnf("[procScenes]: invalid key: %s", key)
		}
	}
}

// procReadAttribute 处理read req
func procReadAttribute(mqttMsg mqttMsgSt) {
	globallogger.Log.Infof("[procReadAttribute]: mqttMsg.Params: %+v", mqttMsg.Params)
	mapParams := mqttMsg.Params.(map[string]interface{})
	var clusterID uint16
	var endpoint string
	for k, v := range mapParams {
		switch k {
		case "clusterID":
			tempClusterID, _ := strconv.ParseUint(v.(string), 10, 16)
			clusterID = uint16(tempClusterID)
		case "endpoint":
			endpoint = v.(string)
		}
	}
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(mapParams["terminalId"].(string))
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(mapParams["terminalId"].(string))
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][procReadAttribute] GetTerminalInfoByDevEUI err: %+v", mapParams["terminalId"], err)
	}
	var endpointTemp pq.StringArray
	if constant.Constant.UsePostgres {
		endpointTemp = terminalInfo.EndpointPG
	} else {
		endpointTemp = terminalInfo.Endpoint
	}
	endpointIndex := 0
	if len(endpointTemp) > 1 {
		for index, EP := range endpointTemp {
			if endpoint == EP {
				endpointIndex = index
			}
		}
	}
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      mapParams["terminalId"].(string),
		CommandType: common.ReadAttribute,
		ClusterID:   clusterID,
		Command: common.Command{
			DstEndpointIndex: endpointIndex,
		},
	})
}
