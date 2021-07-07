package clusterproc

import (
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/zcl/clusterproc/general"
	"github.com/h3c/iotzigbeeserver-go/zcl/clusterproc/heiman"
	"github.com/h3c/iotzigbeeserver-go/zcl/clusterproc/honyar"
	"github.com/h3c/iotzigbeeserver-go/zcl/clusterproc/maileke"
	"github.com/h3c/iotzigbeeserver-go/zcl/clusterproc/measurementandsensing"
	"github.com/h3c/iotzigbeeserver-go/zcl/clusterproc/securityandsafety"
	"github.com/h3c/iotzigbeeserver-go/zcl/clusterproc/smartenergy"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
)

// ClusterProc 处理general clusterID属性消息
func ClusterProc(terminalInfo config.TerminalInfo, srcEndpoint uint8, zclFrame *zcl.Frame, clusterIDByStr string, msgID interface{},
	contentFrame *zcl.Frame, dataTemp string) {
	globallogger.Log.Infof("[devEUI: %v][ClusterProc] zclFrame: FrameControl[%+v] ManufacturerCode[%+v] TransactionSequenceNumber[%+v] CommandIdentifier[%+v] CommandName[%+v] Command[%+v]",
		terminalInfo.DevEUI, zclFrame.FrameControl, zclFrame.ManufacturerCode, zclFrame.TransactionSequenceNumber, zclFrame.CommandIdentifier, zclFrame.CommandName, zclFrame.Command)
	switch clusterIDByStr {
	case "Basic":
		// 处理genBasic消息,0x0000
		general.GenBasicProc(terminalInfo, zclFrame)
	case "PowerConfiguration":
		// 处理genPowerCfg消息,0x0001
		general.GenPowerCfgProc(terminalInfo, zclFrame)
	case "Scenes":
		// 处理genScenes消息，0x0005
		general.GenScenesProc(terminalInfo, zclFrame, msgID, contentFrame, dataTemp)
	case "OnOff":
		// 处理genOnOff消息,0x0006
		general.GenOnOffProc(terminalInfo, srcEndpoint, zclFrame, msgID, contentFrame)
	case "Alarms":
		// 处理genAlarms消息,0x0009
		general.GenAlarmsProc(terminalInfo, zclFrame)
	case "TemperatureMeasurement":
		// 处理TemperatureMeasurement消息，0x0402
		measurementandsensing.TemperatureMeasurementProc(terminalInfo, zclFrame)
	case "RelativeHumidityMeasurement":
		// 处理RelativeHumidityMeasurement消息，0x0405
		measurementandsensing.RelativeHumidityMeasurementProc(terminalInfo, zclFrame)
	case "ElectricalMeasurement":
		// 处理ElectricalMeasurement消息，0x0b04
		measurementandsensing.ElectricalMeasurementProc(terminalInfo, zclFrame)
	case "IASZone":
		// 处理IASZone消息，0x0500
		securityandsafety.IASZoneProc(terminalInfo, zclFrame)
	case "IASWarningDevice":
		// 处理IASWarningDevice消息，0x0502
		securityandsafety.IASWarningDeviceProc(terminalInfo, zclFrame, msgID)
	case "SmartEnergyMetering":
		// 处理SmartEnergyMetering消息，0x0702
		smartenergy.MeteringProc(terminalInfo, zclFrame, msgID, contentFrame)
	case "HEIMANScenes":
		// 处理HEIMAN私有属性Scenes消息，0xfc80
		heiman.ScenesProc(terminalInfo, zclFrame)
	case "HEIMANPM25Measurement":
		// 处理HEIMAN私有属性PM25Measurement消息，0x042a
		heiman.PM25MeasurementProc(terminalInfo, zclFrame)
	case "HEIMANFormaldehydeMeasurement":
		// 处理HEIMAN私有属性FormaldehydeMeasurement消息，0x042b
		heiman.FormaldehydeMeasurementProc(terminalInfo, zclFrame)
	case "HEIMANAirQualityMeasurement":
		// 处理HEIMAN私有属性AirQualityMeasurement消息，0xfc81
		heiman.AirQualityMeasurementProc(terminalInfo, zclFrame, msgID, contentFrame)
	case "HEIMANInfraredRemote":
		// 处理HEIMAN私有属性InfraredRemote消息，0xfc82
		heiman.InfraredRemoteProc(terminalInfo, zclFrame, msgID, contentFrame)
	case "HONYARScenes":
		// 处理HONYAR私有属性Scenes消息，0xfe05
		honyar.ScenesProc(terminalInfo, zclFrame)
	case "HONYARPM25Measurement":
		// 处理HONYAR私有属性PM25Measurement消息，0xfe02
		honyar.PM25MeasurementProc(terminalInfo, zclFrame)
	case "MAILEKEMeasurement":
		// 处理MAILEKE私有属性Measurement消息，0x0415
		maileke.MeasurementProc(terminalInfo, zclFrame)
	default:
		globallogger.Log.Errorf("[devEUI: %v][ClusterProc] invalid clusterID: %v", terminalInfo.DevEUI, clusterIDByStr)
	}
}
