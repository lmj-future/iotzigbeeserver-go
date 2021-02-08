package zclmsgdown

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/frame"
)

// procZclMsgDown 下行公共封装接口
func procZclMsgDown(zclDownMsg common.ZclDownMsg, frameHexString string, transactionID uint8) {
	// globallogger.Log.Infof("[devEUI: %v][procZclMsgDown] frameHexString: %v transactionID: %v", zclDownMsg.DevEUI, frameHexString, transactionID)
	downMsg := common.DownMsg{
		ProfileID:        0x0104,
		ClusterID:        zclDownMsg.ClusterID,
		DstEndpointIndex: zclDownMsg.Command.DstEndpointIndex,
		ZclData:          frameHexString,
		SN:               transactionID,
		DevEUI:           zclDownMsg.DevEUI,
		MsgType:          globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		MsgID:            zclDownMsg.MsgID,
	}
	globallogger.Log.Infof("[devEUI: %v][procZclMsgDown] downMsg: %+v", zclDownMsg.DevEUI, downMsg)
	publicfunction.SendZHADownMsg(downMsg)
}

// procReadAttribute 处理command read 下行命令
func procReadAttribute(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procReadAttribute] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	var attributes []uint16
	switch cluster.ID(zclDownMsg.ClusterID) {
	case cluster.PowerConfiguration:
		attributes = append(attributes, 0x0021)
	case cluster.Scenes:
		attributes = append(attributes, 0x0000, 0x0001, 0x0002, 0x0003, 0x0004, 0x0005)
	case cluster.OnOff:
		attributes = append(attributes, 0x0000)
	case cluster.TemperatureMeasurement:
		attributes = append(attributes, 0x0000, 0x0001, 0x0002, 0x0003)
	case cluster.RelativeHumidityMeasurement:
		attributes = append(attributes, 0x0000, 0x0001, 0x0002, 0x0003)
	case cluster.ElectricalMeasurement:
		attributes = append(attributes, 0x0000, 0x0505, 0x0508, 0x050b, 0x0510, 0x0600, 0x0601, 0x0602, 0x0603, 0x0604, 0x0605, 0x0800, 0x0801, 0x0802, 0x0803)
	case cluster.IASZone:
		attributes = append(attributes, 0x0000, 0x0001, 0x0002)
	case cluster.SmartEnergyMetering:
		attributes = append(attributes, 0x0000, 0x0200, 0x0300, 0x0301, 0x0302, 0x0303, 0x0306, 0x0400)
	case cluster.HEIMANAirQualityMeasurement:
		attributes = append(attributes, 0xf000, 0xf001, 0xf002, 0xf003, 0xf004, 0xf005, 0xf006, 0xf007, 0xf008, 0xf009, 0xf00a)
	case cluster.HEIMANPM25Measurement:
		attributes = append(attributes, 0x0000, 0x0001, 0x0002, 0x0003)
	case cluster.HEIMANFormaldehydeMeasurement:
		attributes = append(attributes, 0x0000, 0x0001, 0x0002, 0x0003)
	default:
		globallogger.Log.Warnf("[devEUI: %v][procReadAttribute] invalid clusterID : %+v", zclDownMsg.DevEUI, cluster.ID(zclDownMsg.ClusterID))
	}
	command := cluster.ReadAttributesCommand{
		AttributeIDs: attributes,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeGlobal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x00,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procHonyarSocketRead 处理command read 下行命令
func procHonyarSocketRead(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procHonyarSocketRead] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	var attributes []uint16
	switch cluster.ID(zclDownMsg.ClusterID) {
	case cluster.OnOff:
		attributes = append(attributes, 0x0000, 0x8000, 0x8005, 0x8006)
	default:
		globallogger.Log.Warnf("[devEUI: %v][procHonyarSocketRead] invalid clusterID : %+v", zclDownMsg.DevEUI, cluster.ID(zclDownMsg.ClusterID))
	}
	command := cluster.ReadAttributesCommand{
		AttributeIDs: attributes,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeGlobal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x00,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procHonyarSocketReadUSB 处理command read 下行命令
func procHonyarSocketReadUSB(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procHonyarSocketReadUSB] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	var attributes []uint16
	switch cluster.ID(zclDownMsg.ClusterID) {
	case cluster.OnOff:
		attributes = append(attributes, 0x0000)
	default:
		globallogger.Log.Warnf("[devEUI: %v][procHonyarSocketReadUSB] invalid clusterID : %+v", zclDownMsg.DevEUI, cluster.ID(zclDownMsg.ClusterID))
	}
	command := cluster.ReadAttributesCommand{
		AttributeIDs: attributes,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeGlobal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x00,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	zclDownMsg.Command.DstEndpointIndex = 1
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procReadBasic 处理command read basic下行命令
func procReadBasic(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procReadBasic] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.ReadAttributesCommand{
		AttributeIDs: []uint16{
			0x0000, //ZLibraryVersion
			0x0001, //ApplicationVersion
			0x0002, //StackVersion
			0x0003, //HWVersion
			0x0004, //ManufacturerName
			0x0005, //ModelIdentifier
			0x0006, //DateCode
			0x0007, //PowerSource
			// 0x0010, //LocationDescription
			// 0x0011, //PhysicalEnvironment
			// 0x0012, //DeviceEnabled
			// 0x0013, //AlarmMask
			// 0x0014, //DisableLocalConfig
			// 0x4000, //SWBuildID
		},
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeGlobal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x00,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procSwitchCommand 处理command on off下行命令
func procSwitchCommand(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procSwitchCommand] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	var command interface{}
	if zclDownMsg.Command.Cmd == 0x00 {
		command = cluster.OffCommand{}
	} else {
		command = cluster.OnCommand{}
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           zclDownMsg.Command.Cmd,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procGetSwitchState 处理command read下行命令
func procGetSwitchState(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procGetSwitchState] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.ReadAttributesCommand{
		AttributeIDs: []uint16{0x0000},
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeGlobal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x00,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procSocketCommand 处理command on off下行命令
func procSocketCommand(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procSocketCommand] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	var command interface{}
	if zclDownMsg.Command.Cmd == 0x00 {
		command = cluster.OffCommand{}
	} else {
		command = cluster.OnCommand{}
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           zclDownMsg.Command.Cmd,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procGetSocketState 处理command read下行命令
func procGetSocketState(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procGetSocketState] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.ReadAttributesCommand{
		AttributeIDs: []uint16{0x0000},
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeGlobal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x00,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procHonyarSocketWriteReq 处理Honyar command write下行命令
func procHonyarSocketWriteReq(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procHonyarSocketWriteReq] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.WriteAttributesCommand{
		WriteAttributeRecords: []*cluster.WriteAttributeRecord{
			{
				AttributeName: zclDownMsg.Command.AttributeName,
				AttributeID:   zclDownMsg.Command.AttributeID,
				Attribute: &cluster.Attribute{
					DataType: cluster.ZclDataTypeBoolean,
					Value:    zclDownMsg.Command.Value,
				},
			},
		},
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeGlobal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x02,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procSensorWriteReq 处理command received ZoneEnrollResponse下行命令
func procSensorWriteReq(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procSensorWriteReq] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	T300ID := zclDownMsg.T300ID
	command := cluster.WriteAttributesCommand{
		WriteAttributeRecords: []*cluster.WriteAttributeRecord{
			{
				AttributeName: "IAS_CIE_Address",
				AttributeID:   0x0010,
				Attribute: &cluster.Attribute{
					DataType: cluster.ZclDataTypeIeeeAddr,
					Value:    "0x" + T300ID,
				},
			},
		},
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeGlobal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x02,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procSensorEnrollRsp 处理command received ZoneEnrollResponse下行命令
func procSensorEnrollRsp(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procSensorEnrollRsp] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.ZoneEnrollResponse{
		ResponseCode: 0x00,
		ZoneID:       0x01,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x00,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procSensorReadReq 处理command read下行命令
func procSensorReadReq(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procSensorReadReq] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
}

// procIntervalCommand 处理command received ZoneEnrollResponse下行命令
func procIntervalCommand(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procIntervalCommand] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	type tempSt struct {
		attributeName     string
		attributeID       uint16
		attributeDataType cluster.ZclDataType
	}
	data := tempSt{}
	var temp []tempSt = make([]tempSt, 1, 20)
	switch cluster.ID(zclDownMsg.ClusterID) {
	case cluster.Basic:
		data.attributeName = "ZLibraryVersion"
		data.attributeID = 0x0000
		data.attributeDataType = cluster.ZclDataTypeUint8
		temp[0] = data
	case cluster.PowerConfiguration:
		data.attributeName = "BatteryPercentageRemaining"
		data.attributeID = 0x0021
		data.attributeDataType = cluster.ZclDataTypeUint8
		temp[0] = data
	case cluster.OnOff:
		data.attributeName = "OnOff"
		data.attributeID = 0x0000
		data.attributeDataType = cluster.ZclDataTypeBoolean
		temp[0] = data
	case cluster.TemperatureMeasurement:
		data.attributeName = "MeasuredValue"
		data.attributeID = 0x0000
		data.attributeDataType = cluster.ZclDataTypeInt16
		temp[0] = data
	case cluster.RelativeHumidityMeasurement:
		data.attributeName = "MeasuredValue"
		data.attributeID = 0x0000
		data.attributeDataType = cluster.ZclDataTypeUint16
		temp[0] = data
	case cluster.IASZone:
		data.attributeName = "ZoneState"
		data.attributeID = 0x0000
		data.attributeDataType = cluster.ZclDataTypeEnum8
		temp[0] = data
	case cluster.SmartEnergyMetering:
		data.attributeName = "CurrentSummationDelivered" // 耗电量 kW/H
		data.attributeID = 0x0000
		data.attributeDataType = cluster.ZclDataTypeUint48
		temp[0] = data
	case cluster.ElectricalMeasurement:
		data.attributeName = "RMSVoltage" // 电压V
		data.attributeID = 0x0505
		data.attributeDataType = cluster.ZclDataTypeUint16
		temp[0] = data
		data.attributeName = "RMSCurrent" // 电流A
		data.attributeID = 0x0508
		data.attributeDataType = cluster.ZclDataTypeUint16
		temp = append(temp, data)
		data.attributeName = "ActivePower" // 负载功率W
		data.attributeID = 0x050b
		data.attributeDataType = cluster.ZclDataTypeInt16
		temp = append(temp, data)
	default:
		globallogger.Log.Warnf("[devEUI: %v][procIntervalCommand] invalid clusterID: %v", zclDownMsg.DevEUI, cluster.ID(zclDownMsg.ClusterID))
	}
	command := cluster.ConfigureReportingCommand{}
	var attributeReportingConfigurationRecords []*cluster.AttributeReportingConfigurationRecord = make([]*cluster.AttributeReportingConfigurationRecord, len(temp))
	for i, v := range temp {
		attributeReportingConfigurationRecords[i] = &cluster.AttributeReportingConfigurationRecord{
			Direction:                cluster.ReportDirectionAttributeReported,
			AttributeName:            v.attributeName,
			AttributeID:              v.attributeID,
			AttributeDataType:        v.attributeDataType,
			MinimumReportingInterval: zclDownMsg.Command.Interval,
			MaximumReportingInterval: zclDownMsg.Command.Interval,
			ReportableChange:         &cluster.Attribute{},
			TimeoutPeriod:            0,
		}
	}
	command.AttributeReportingConfigurationRecords = attributeReportingConfigurationRecords
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeGlobal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x06,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	// if zclDownMsg.ClusterID == 0x0702 {
	// 	frameHexString = frameHexString[:6] + "00000025ac0d1b0e0100000000000000042a0200481c0f0000"
	// }
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procSetThreshold 处理command write下行命令
func procSetThreshold(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procSetThreshold] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	var command interface{}
	switch zclDownMsg.ClusterID {
	case 0xfc81:
		if (zclDownMsg.Command.MaxTemperature != 0 || zclDownMsg.Command.MinTemperature != 0) &&
			(zclDownMsg.Command.MaxHumidity != 0 || zclDownMsg.Command.MinHumidity != 0) {
			command = cluster.WriteAttributesCommand{
				WriteAttributeRecords: []*cluster.WriteAttributeRecord{
					{
						AttributeName: "MaxTemperature",
						AttributeID:   0xf006,
						Attribute: &cluster.Attribute{
							DataType: cluster.ZclDataTypeInt16,
							Value:    int64(zclDownMsg.Command.MaxTemperature),
						},
					},
					{
						AttributeName: "MinTemperature",
						AttributeID:   0xf007,
						Attribute: &cluster.Attribute{
							DataType: cluster.ZclDataTypeInt16,
							Value:    int64(zclDownMsg.Command.MinTemperature),
						},
					},
					{
						AttributeName: "MaxHumidity",
						AttributeID:   0xf008,
						Attribute: &cluster.Attribute{
							DataType: cluster.ZclDataTypeUint16,
							Value:    uint64(zclDownMsg.Command.MaxHumidity),
						},
					},
					{
						AttributeName: "MinHumidity",
						AttributeID:   0xf009,
						Attribute: &cluster.Attribute{
							DataType: cluster.ZclDataTypeUint16,
							Value:    uint64(zclDownMsg.Command.MinHumidity),
						},
					},
					{
						AttributeName: "Disturb",
						AttributeID:   0xf00a,
						Attribute: &cluster.Attribute{
							DataType: cluster.ZclDataTypeUint16,
							Value:    uint64(zclDownMsg.Command.Disturb),
						},
					},
				},
			}
		} else if zclDownMsg.Command.MaxTemperature != 0 || zclDownMsg.Command.MinTemperature != 0 {
			command = cluster.WriteAttributesCommand{
				WriteAttributeRecords: []*cluster.WriteAttributeRecord{
					{
						AttributeName: "MaxTemperature",
						AttributeID:   0xf006,
						Attribute: &cluster.Attribute{
							DataType: cluster.ZclDataTypeInt16,
							Value:    int64(zclDownMsg.Command.MaxTemperature),
						},
					},
					{
						AttributeName: "MinTemperature",
						AttributeID:   0xf007,
						Attribute: &cluster.Attribute{
							DataType: cluster.ZclDataTypeInt16,
							Value:    int64(zclDownMsg.Command.MinTemperature),
						},
					},
					{
						AttributeName: "Disturb",
						AttributeID:   0xf00a,
						Attribute: &cluster.Attribute{
							DataType: cluster.ZclDataTypeUint16,
							Value:    uint64(zclDownMsg.Command.Disturb),
						},
					},
				},
			}
		} else if zclDownMsg.Command.MaxHumidity != 0 || zclDownMsg.Command.MinHumidity != 0 {
			command = cluster.WriteAttributesCommand{
				WriteAttributeRecords: []*cluster.WriteAttributeRecord{
					{
						AttributeName: "MaxHumidity",
						AttributeID:   0xf008,
						Attribute: &cluster.Attribute{
							DataType: cluster.ZclDataTypeUint16,
							Value:    uint64(zclDownMsg.Command.MaxHumidity),
						},
					},
					{
						AttributeName: "MinHumidity",
						AttributeID:   0xf009,
						Attribute: &cluster.Attribute{
							DataType: cluster.ZclDataTypeUint16,
							Value:    uint64(zclDownMsg.Command.MinHumidity),
						},
					},
					{
						AttributeName: "Disturb",
						AttributeID:   0xf00a,
						Attribute: &cluster.Attribute{
							DataType: cluster.ZclDataTypeUint16,
							Value:    uint64(zclDownMsg.Command.Disturb),
						},
					},
				},
			}
		} else if zclDownMsg.Command.Disturb != 0 {
			command = cluster.WriteAttributesCommand{
				WriteAttributeRecords: []*cluster.WriteAttributeRecord{
					{
						AttributeName: "Disturb",
						AttributeID:   0xf00a,
						Attribute: &cluster.Attribute{
							DataType: cluster.ZclDataTypeUint16,
							Value:    uint64(zclDownMsg.Command.Disturb),
						},
					},
				},
			}
		}
	default:
		globallogger.Log.Warnf("[devEUI: %v][procSetThreshold] invalid clusterID: %+v", zclDownMsg.DevEUI, zclDownMsg.ClusterID)
		return
	}

	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeGlobal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x02,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procStartWarning 处理command StartWarning下行命令
func procStartWarning(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procStartWarning] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.StartWarning{
		WarningControl:  zclDownMsg.Command.WarningControl,
		WarningDuration: zclDownMsg.Command.WarningTime,
		StrobeDutyCycle: 0x28,
		StrobeLevel:     0x01,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x00,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procHeimanScenesDefaultResponse 处理heiman command 0xfc80下行命令
func procHeimanScenesDefaultResponse(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procHeimanScenesDefaultResponse] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.DefaultResponseCommand{
		CommandID: zclDownMsg.Command.Cmd,
		Status:    cluster.ZclStatusSuccess,
	}
	frameConfig := frame.Configuration{
		FrameType:                  frame.FrameTypeGlobal,
		FrameTypeConfigured:        true,
		ManufacturerCode:           0x120b,
		ManufacturerCodeConfigured: true,
		Direction:                  frame.DirectionServerClient,
		DirectionConfigured:        true,
		CommandID:                  0x0b,
		CommandIDConfigured:        true,
		Command:                    command,
		CommandConfigured:          true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, zclDownMsg.Command.TransactionID)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procAddScene 处理command AddScene下行命令
func procAddScene(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procAddScene] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.AddSceneCommand{
		GroupID:        zclDownMsg.Command.GroupID,
		SceneID:        zclDownMsg.Command.SceneID,
		TransitionTime: zclDownMsg.Command.TransitionTime,
		SceneName:      zclDownMsg.Command.SceneName,
		KeyID:          zclDownMsg.Command.KeyID,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x00,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	tempBuf, _ := hex.DecodeString(frameHexString[16 : len(frameHexString)-2])
	tempString := ""
	for _, v := range tempBuf {
		tempString += string(v)
	}
	frameHexString = string(frameHexString[:16]) + tempString + string(frameHexString[len(frameHexString)-2:])

	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procViewScene 处理command ViewScene下行命令
func procViewScene(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procViewScene] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.ViewSceneCommand{
		GroupID: zclDownMsg.Command.GroupID,
		SceneID: zclDownMsg.Command.SceneID,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x01,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procRemoveScene 处理command RemoveScene下行命令
func procRemoveScene(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procRemoveScene] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.RemoveSceneCommand{
		GroupID: zclDownMsg.Command.GroupID,
		SceneID: zclDownMsg.Command.SceneID,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x02,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procRemoveAllScenes 处理command RemoveAllScenes下行命令
func procRemoveAllScenes(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procRemoveAllScenes] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.RemoveAllScenesCommand{
		GroupID: zclDownMsg.Command.GroupID,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x03,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procStoreScene 处理command StoreScene下行命令
func procStoreScene(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procStoreScene] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.StoreSceneCommand{
		GroupID: zclDownMsg.Command.GroupID,
		SceneID: zclDownMsg.Command.SceneID,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x04,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procRecallScene 处理command RecallScene下行命令
func procRecallScene(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procRecallScene] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.RecallSceneCommand{
		GroupID: zclDownMsg.Command.GroupID,
		SceneID: zclDownMsg.Command.SceneID,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x05,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procGetSceneMembership 处理command GetSceneMembership下行命令
func procGetSceneMembership(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procGetSceneMembership] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.GetSceneMembership{
		GroupID: zclDownMsg.Command.GroupID,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x06,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procEnhancedAddScene 处理command EnhancedAddScene下行命令
func procEnhancedAddScene(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procEnhancedAddScene] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.EnhancedAddSceneCommand{
		GroupID:        zclDownMsg.Command.GroupID,
		SceneID:        zclDownMsg.Command.SceneID,
		TransitionTime: zclDownMsg.Command.TransitionTime,
		SceneName:      zclDownMsg.Command.SceneName,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x40,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procEnhancedViewScene 处理command EnhancedViewScene下行命令
func procEnhancedViewScene(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procEnhancedViewScene] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.EnhancedViewSceneCommand{
		GroupID: zclDownMsg.Command.GroupID,
		SceneID: zclDownMsg.Command.SceneID,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x41,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procCopyScene 处理command CopyScene下行命令
func procCopyScene(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procCopyScene] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.CopySceneCommand{
		Mode:        zclDownMsg.Command.Mode,
		FromGroupID: zclDownMsg.Command.FromGroupID,
		FromSceneID: zclDownMsg.Command.FromSceneID,
		ToGroupID:   zclDownMsg.Command.ToGroupID,
		ToSceneID:   zclDownMsg.Command.ToSceneID,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x42,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procSendKeyCommand 处理command SendKeyCommand下行命令
func procSendKeyCommand(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procSendKeyCommand] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.HEIMANInfraredRemoteSendKeyCommand{
		ID:      zclDownMsg.Command.InfraredRemoteID,
		KeyCode: zclDownMsg.Command.InfraredRemoteKeyCode,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0xf0,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procStudyKey 处理command StudyKey下行命令
func procStudyKey(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procStudyKey] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.HEIMANInfraredRemoteStudyKey{
		ID:      zclDownMsg.Command.InfraredRemoteID,
		KeyCode: zclDownMsg.Command.InfraredRemoteKeyCode,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0xf1,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procDeleteKey 处理command DeleteKey下行命令
func procDeleteKey(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procDeleteKey] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.HEIMANInfraredRemoteDeleteKey{
		ID:      zclDownMsg.Command.InfraredRemoteID,
		KeyCode: zclDownMsg.Command.InfraredRemoteKeyCode,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0xf3,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procCreateID 处理command CreateID下行命令
func procCreateID(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procCreateID] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.HEIMANInfraredRemoteCreateID{
		ModelType: zclDownMsg.Command.InfraredRemoteModelType,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0xf4,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procGetIDAndKeyCodeList 处理command GetIDAndKeyCodeList下行命令
func procGetIDAndKeyCodeList(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procGetIDAndKeyCodeList] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.HEIMANInfraredRemoteGetIDAndKeyCodeList{}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0xf6,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procSyncTime 处理时间同步下行命令
func procSyncTime(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procSyncTime] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	// 空气盒子UTC时间从2000年开始计算，终端厂商的问题
	var date = "2000-01-01 00:00:00"
	loc1, _ := time.LoadLocation("Local")
	// 将时间字符串转成本地时间
	result, err := time.ParseInLocation("2006-01-02 15:05:06", date, loc1)
	if err != nil {
		fmt.Println("err", err)
	}
	loc2, _ := time.LoadLocation("") // 相当于LoadLocation("UTC")
	day, _ := time.ParseDuration("8h")
	endTime := result.In(loc2).Add(day) // 转成UTC时间并加8h
	command := cluster.WriteAttributesCommand{
		WriteAttributeRecords: []*cluster.WriteAttributeRecord{
			{
				AttributeName: "Time",
				AttributeID:   0x0000,
				Attribute: &cluster.Attribute{
					DataType: cluster.ZclDataTypeUtc,
					Value:    uint32(time.Now().UTC().Unix() - endTime.UTC().Unix()),
				},
			},
			{
				AttributeName: "TimeZone",
				AttributeID:   0x0002,
				Attribute: &cluster.Attribute{
					DataType: cluster.ZclDataTypeInt32,
					Value:    int64(8 * 60 * 60),
				},
			},
		},
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeGlobal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x02,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procSetLanguage 处理command SetLanguage下行命令
func procSetLanguage(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procSetLanguage] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.HEIMANSetLanguage{
		Language: zclDownMsg.Command.Cmd,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x00,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

// procSetUnitOfTemperature 处理command SetUnitOfTemperature下行命令
func procSetUnitOfTemperature(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][procSetUnitOfTemperature] start: %+v", zclDownMsg.DevEUI, zclDownMsg)
	command := cluster.HEIMANSetUnitOfTemperature{
		UnitOfTemperature: zclDownMsg.Command.Cmd,
	}
	frameConfig := frame.Configuration{
		FrameType:           frame.FrameTypeLocal,
		FrameTypeConfigured: true,
		Direction:           frame.DirectionClientServer,
		DirectionConfigured: true,
		CommandID:           0x01,
		CommandIDConfigured: true,
		Command:             command,
		CommandConfigured:   true,
	}
	z := zcl.New()
	frameHexString, transactionID := z.EncFrameConfigurationToHexString(frameConfig, 0)
	procZclMsgDown(zclDownMsg, frameHexString, transactionID)
}

//ProcZclDownMsg ZCL下行命令处理
func ProcZclDownMsg(zclDownMsg common.ZclDownMsg) {
	globallogger.Log.Infof("[devEUI: %v][ProcZclDownMsg] zclDownMsg: %v", zclDownMsg.DevEUI, zclDownMsg)
	switch zclDownMsg.CommandType {
	case common.ReadAttribute:
		procReadAttribute(zclDownMsg)
	case common.HonyarSocketRead:
		procHonyarSocketRead(zclDownMsg)
		go procHonyarSocketReadUSB(zclDownMsg)
	case common.ReadBasic:
		procReadBasic(zclDownMsg)
	case common.SwitchCommand:
		procSwitchCommand(zclDownMsg)
	case common.GetSwitchState:
		procGetSwitchState(zclDownMsg)
	case common.SocketCommand:
		procSocketCommand(zclDownMsg)
	case common.GetSocketState:
		procGetSocketState(zclDownMsg)
	case common.ChildLock, common.BackGroundLight, common.PowerOffMemory, common.HistoryElectricClear:
		procHonyarSocketWriteReq(zclDownMsg)
	case common.SensorWriteReq:
		procSensorWriteReq(zclDownMsg)
	case common.SensorEnrollRsp:
		procSensorEnrollRsp(zclDownMsg)
	case common.SensorReadReq:
		procSensorReadReq(zclDownMsg)
	case common.IntervalCommand:
		procIntervalCommand(zclDownMsg)
	case common.SetThreshold:
		procSetThreshold(zclDownMsg)
	case common.StartWarning:
		procStartWarning(zclDownMsg)
	case common.HeimanScenesDefaultResponse:
		procHeimanScenesDefaultResponse(zclDownMsg)
	case common.AddScene:
		procAddScene(zclDownMsg)
	case common.ViewScene:
		procViewScene(zclDownMsg)
	case common.RemoveScene:
		procRemoveScene(zclDownMsg)
	case common.RemoveAllScenes:
		procRemoveAllScenes(zclDownMsg)
	case common.StoreScene:
		procStoreScene(zclDownMsg)
	case common.RecallScene:
		procRecallScene(zclDownMsg)
	case common.GetSceneMembership:
		procGetSceneMembership(zclDownMsg)
	case common.EnhancedAddScene:
		procEnhancedAddScene(zclDownMsg)
	case common.EnhancedViewScene:
		procEnhancedViewScene(zclDownMsg)
	case common.CopyScene:
		procCopyScene(zclDownMsg)
	case common.SendKeyCommand:
		procSendKeyCommand(zclDownMsg)
	case common.StudyKey:
		procStudyKey(zclDownMsg)
	case common.DeleteKey:
		procDeleteKey(zclDownMsg)
	case common.CreateID:
		procCreateID(zclDownMsg)
	case common.GetIDAndKeyCodeList:
		procGetIDAndKeyCodeList(zclDownMsg)
	case common.SyncTime:
		procSyncTime(zclDownMsg)
	case common.SetLanguage:
		procSetLanguage(zclDownMsg)
	case common.SetUnitOfTemperature:
		procSetUnitOfTemperature(zclDownMsg)
	default:
		globallogger.Log.Warnf("[devEUI: %v][ProcZclDownMsg] invalid commandType: %v", zclDownMsg.DevEUI, zclDownMsg.CommandType)
	}
}
