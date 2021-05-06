package iotsmartspace

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/db/mqtt"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
)

// 定义method常量
const (
	MethodPropertyUp               = "property.up"
	MethodPropertyUpReply          = "property.up.reply"
	MethodPropertyDown             = "property.down"
	MethodPropertyDownReply        = "property.down.reply"
	MethodStateUp                  = "state.up"
	MethodStateUpReply             = "state.up.reply"
	MethodStateDown                = "state.down"
	MethodStateDownReply           = "state.down.reply"
	MethodStateLeave               = "state.leave"
	MethodStateLeaveReply          = "state.leave.reply"
	MethodActionInsert             = "action.insert"
	MethodActionInsertReply        = "action.insert.reply"
	MethodActionInsertReplyTest    = "action.insert.reply.test"
	MethodActionInsertSuccess      = "action.insert.success"
	MethodActionInsertSuccessReply = "action.insert.success.reply"
	MethodActionInsertFail         = "action.insert.fail"
	MethodActionInsertFailReply    = "action.insert.fail.reply"
	MethodActionDelete             = "action.delete"
	MethodActionDeleteReply        = "action.delete.reply"
	MethodActionOffline            = "action.offline"
	MethodActionOfflineReply       = "action.offline.reply"
)

// 定义topic常量
const (
	TopicZigbeeserverIotsmartspaceProperty = "/zigbeeserver/iotsmartspace/property"
	TopicIotsmartspaceZigbeeserverProperty = "/iotsmartspace/zigbeeserver/property"
	TopicZigbeeserverIotsmartspaceState    = "/zigbeeserver/iotsmartspace/state"
	TopicIotsmartspaceZigbeeserverState    = "/iotsmartspace/zigbeeserver/state"
	TopicZigbeeserverIotsmartspaceAction   = "/zigbeeserver/iotsmartspace/action"
	TopicIotsmartspaceZigbeeserverAction   = "/iotsmartspace/zigbeeserver/action"
)

// 定义topic常量(iotware)
const (
	TopicV1GatewayTelemetry         = "v1/gateway/telemetry"
	TopicV1GatewayRPC               = "v1/gateway/rpc"
	TopicV1GatewayConnect           = "v1/gateway/connect"
	TopicV1GatewayDisconnect        = "v1/gateway/disconnect"
	TopicV1GatewayNetworkIn         = "v1/gateway/network/in"
	TopicV1GatewayNetworkInAck      = "v1/gateway/network/in/ack"
	TopicV1GatewayNetworkOut        = "v1/gateway/network/out"
	TopicV1GatewayEventDeviceDelete = "v1/gateway/event/device/delete"
)

// 定义终端属性
const (
	// 海曼智能计量插座（0x1000001x）
	//      0x10000011（开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	//      0x10000012（功率）：string
	//      0x10000013（电量）：string
	HeimanSmartPlugPropertyOnOff    = "10000011"
	HeimanSmartPlugPropertyPower    = "10000012"
	HeimanSmartPlugPropertyElectric = "10000013"
	// 海曼智能墙面插座（0x1000002x）
	// 		0x10000021（开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	// 		0x10000022（功率）：string
	// 		0x10000023（电量）：string
	HeimanESocketPropertyOnOff    = "10000021"
	HeimanESocketPropertyPower    = "10000022"
	HeimanESocketPropertyElectric = "10000023"
	// 海曼火灾探测器（0x1000003x）
	// 		0x10000031（告警）："0"告警上报，"1"告警解除
	// 		0x10000032（剩余电量）：string，百分比
	HeimanSmokeSensorPropertyAlarm        = "10000031"
	HeimanSmokeSensorPropertyLeftElectric = "10000032"
	// 海曼水浸传感器（0x1000004x）
	// 		0x10000041（告警）："0"告警上报，"1"告警解除
	// 		0x10000042（剩余电量）：string，百分比
	// 		0x10000043（防拆）："0"设备拆除，"1"设备拆除解除
	HeimanWaterSensorPropertyAlarm        = "10000041"
	HeimanWaterSensorPropertyLeftElectric = "10000042"
	HeimanWaterSensorPropertyTamper       = "10000043"
	// 海曼人体传感器（0x1000005x）
	// 		0x10000051（人体检测）："0"无人移动，"1"有人移动
	// 		0x10000052（剩余电量）：string，百分比
	// 		0x10000053（防拆）："0"设备拆除，"1"设备拆除解除
	HeimanPIRSensorPropertyAlarm        = "10000051"
	HeimanPIRSensorPropertyLeftElectric = "10000052"
	HeimanPIRSensorPropertyTamper       = "10000053"
	// 海曼温湿度传感器（0x1000006x）
	// 		0x10000061（温度）：[-10, 60]，string
	// 		0x10000062（湿度）：[0, 100]，string
	// 		0x10000063（剩余电量）：string，百分比
	HeimanHTSensorPropertyTemperature  = "10000061"
	HeimanHTSensorPropertyHumidity     = "10000062"
	HeimanHTSensorPropertyLeftElectric = "10000063"
	// 海曼一氧化碳探测器（0x1000007x）
	// 		0x10000071（告警）："0"告警上报，"1"告警解除
	// 		0x10000072（剩余电量）：string，百分比
	HeimanCOSensorPropertyAlarm        = "10000071"
	HeimanCOSensorPropertyLeftElectric = "10000072"
	// 海曼智能声光报警器（0x1000008x）
	// 		0x10000081（报警时长设置）：string
	// 		0x10000082（剩余电量）：string，百分比
	// 		0x10000083（下发警报）："0"触发警报，"1"解除警报
	HeimanWarningDevicePropertyAlarmTime    = "10000081"
	HeimanWarningDevicePropertyLeftElectric = "10000082"
	HeimanWarningDevicePropertyAlarmStart   = "10000083"
	// 海曼可燃气体探测器（0x1000009x）
	// 		0x10000091（告警）："0"告警上报，"1"告警解除
	HeimanGASSensorPropertyAlarm = "10000091"
	// 海曼空气质量仪（0x100000Ax）
	// 		0x100000A1（温度）：[-10, 60]，string
	// 		0x100000A2（湿度）：[0, 100]，string
	// 		0x100000A3（甲醛）：string，百分比
	// 		0x100000A4（PM2.5）：string
	// 		0x100000A5（PM10）：string
	// 		0x100000A6（剩余电量）：string，百分比
	// 		0x100000A7（温度报警阈值设置）：string
	// 		0x100000A8（湿度报警阈值设置）：string
	// 		0x100000A9（温度报警）："0"过高，"1"过高解除，"2"过低，"3"过低解除
	// 		0x100000AA（湿度报警）："0"过高，"1"过高解除，"2"过低，"3"过低解除
	// 		0x100000AB（TVOC）：string
	// 		0x100000AC（AQI）：string
	// 		0x100000AD（电池状态）："0"未充电，"1"充电中，"2"已充满
	//      0x100000AE（勿扰模式）："0"关闭，"1"开启 默认关闭
	//      0x100000AF（甲醛告警）："0"超标警报上报
	//      0x100000B1（PM2.5告警）："0"超标警报上报
	//      0x100000B2（PM10告警）："0"超标警报上报
	//      0x100000B3（TVOC告警）："0"轻度污染上报，"1"中度污染上报
	//      0x100000B4（AQI告警）："0"超标警报上报
	//      0x100000B5（电池告警）："0"电池低压警报上报
	HeimanHS2AQPropertyTemperature               = "100000A1"
	HeimanHS2AQPropertyHumidity                  = "100000A2"
	HeimanHS2AQPropertyFormaldehyde              = "100000A3"
	HeimanHS2AQPropertyPM25                      = "100000A4"
	HeimanHS2AQPropertyPM10                      = "100000A5"
	HeimanHS2AQPropertyLeftElectric              = "100000A6"
	HeimanHS2AQPropertyTemperatureAlarmThreshold = "100000A7"
	HeimanHS2AQPropertyHumidityAlarmThreshold    = "100000A8"
	HeimanHS2AQPropertyTemperatureAlarm          = "100000A9"
	HeimanHS2AQPropertyHumidityAlarm             = "100000AA"
	HeimanHS2AQPropertyTVOC                      = "100000AB"
	HeimanHS2AQPropertyAQI                       = "100000AC"
	HeimanHS2AQPropertyBatteryState              = "100000AD"
	HeimanHS2AQPropertyUnDisturb                 = "100000AE"
	HeimanHS2AQPropertyFormaldehydeAlarm         = "100000AF"
	HeimanHS2AQPropertyPM25Alarm                 = "100000B1"
	HeimanHS2AQPropertyPM10Alarm                 = "100000B2"
	HeimanHS2AQPropertyTVOCAlarm                 = "100000B3"
	HeimanHS2AQPropertyAQIAlarm                  = "100000B4"
	HeimanHS2AQPropertyBatteryAlarm              = "100000B5"
	// 海曼智能开关（0x100000Cx）
	// 		0x100000Cx（一路开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	HeimanHS2SW1LEFR30PropertyOnOff1 = "100000C1"
	// 		0x100000Dx（二路开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	HeimanHS2SW2LEFR30PropertyOnOff1 = "100000D1"
	HeimanHS2SW2LEFR30PropertyOnOff2 = "100000D2"
	// 		0x100000Ex（三路开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	HeimanHS2SW3LEFR30PropertyOnOff1 = "100000E1"
	HeimanHS2SW3LEFR30PropertyOnOff2 = "100000E2"
	HeimanHS2SW3LEFR30PropertyOnOff3 = "100000E3"
	// 海曼智能情景开关（0x100000Fx）
	HeimanSceneSwitchEM30PropertyAtHome       = "100000F2" //0xf0, 在家模式
	HeimanSceneSwitchEM30PropertyGoOut        = "100000F1" //0xf1, 外出模式
	HeimanSceneSwitchEM30PropertyCinema       = "100000F3" //0xf2, 影院模式
	HeimanSceneSwitchEM30PropertyRepast       = "100000F6" //0xf3, 就餐模式
	HeimanSceneSwitchEM30PropertySleep        = "100000F4" //0xf4, 就寝模式
	HeimanSceneSwitchEM30PropertyLeftElectric = "100000F5" //剩余电量,string，百分比
	// 海曼智能红外转发器（0x1000011x）
	HeimanIRControlEMPropertySendKeyCommand      = "10000111" //0xf0, 发送遥控指令
	HeimanIRControlEMPropertyStudyKey            = "10000112" //0xf1, 学习遥控指令
	HeimanIRControlEMPropertyDeleteKey           = "10000113" //0xf3, 删除遥控指令
	HeimanIRControlEMPropertyCreateID            = "10000114" //0xf4, 创建遥控设备
	HeimanIRControlEMPropertyGetIDAndKeyCodeList = "10000115" //0xf6, 获取遥控设备及其指令列表

	// 鸿雁智能开关（0x2000001x）
	// 		0x2000001x（一路开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	HonyarSingleSwitch00500c32PropertyOnOff1 = "20000011"
	// 		0x2000002x（二路开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	HonyarDoubleSwitch00500c33PropertyOnOff1 = "20000021"
	HonyarDoubleSwitch00500c33PropertyOnOff2 = "20000022"
	// 		0x2000003x（三路开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	HonyarTripleSwitch00500c35PropertyOnOff1 = "20000031"
	HonyarTripleSwitch00500c35PropertyOnOff2 = "20000032"
	HonyarTripleSwitch00500c35PropertyOnOff3 = "20000033"
	// 鸿雁智能墙面插座（0x2000004x）
	// 		0x20000041（开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	// 		0x20000042（总耗电量）：string
	// 		0x20000043（电压）：string
	// 		0x20000044（电流）：string
	// 		0x20000045（负载功率）：string
	// 		0x20000046（童锁开关）："00"关，"01"开
	// 		0x20000047（USB开关）："00"关，"01"开
	HonyarSocket000a0c3cPropertyOnOff                    = "20000041"
	HonyarSocket000a0c3cPropertyCurretSummationDelivered = "20000042"
	HonyarSocket000a0c3cPropertyRMSVoltage               = "20000043"
	HonyarSocket000a0c3cPropertyRMSCurrent               = "20000044"
	HonyarSocket000a0c3cPropertyActivePower              = "20000045"
	HonyarSocket000a0c3cPropertyChildLock                = "20000046"
	HonyarSocket000a0c3cPropertyUSB                      = "20000047"
	// 鸿雁二键情景开关（0x2000005x）
	Honyar2SceneSwitch005f0cf3PropertyLeft  = "20000051" //0x01, 左键
	Honyar2SceneSwitch005f0cf3PropertyRight = "20000052" //0x02, 右键
	// 鸿雁六键情景开关（0x2000006x）
	Honyar6SceneSwitch005f0c3bPropertyLeftUp              = "20000061" //0x01, 左上键
	Honyar6SceneSwitch005f0c3bPropertyLeftMiddle          = "20000062" //0x02, 左中键
	Honyar6SceneSwitch005f0c3bPropertyLeftDown            = "20000063" //0x03, 左下键
	Honyar6SceneSwitch005f0c3bPropertyRightUp             = "20000064" //0x04, 右上键
	Honyar6SceneSwitch005f0c3bPropertyRightMiddle         = "20000065" //0x05, 右中键
	Honyar6SceneSwitch005f0c3bPropertyRightDown           = "20000066" //0x06, 右下键
	Honyar6SceneSwitch005f0c3bPropertyLeftUpAddScene      = "20000067" //0x01~0x21，对应1~33个场景
	Honyar6SceneSwitch005f0c3bPropertyLeftMiddleAddScene  = "20000068" //0x01~0x21，对应1~33个场景
	Honyar6SceneSwitch005f0c3bPropertyLeftDownAddScene    = "20000069" //0x01~0x21，对应1~33个场景
	Honyar6SceneSwitch005f0c3bPropertyRightUpAddScene     = "2000006A" //0x01~0x21，对应1~33个场景
	Honyar6SceneSwitch005f0c3bPropertyRightMiddleAddScene = "2000006B" //0x01~0x21，对应1~33个场景
	Honyar6SceneSwitch005f0c3bPropertyRightDownAddScene   = "2000006C" //0x01~0x21，对应1~33个场景
	// 鸿雁一键情景开关（0x2000007x）
	Honyar1SceneSwitch005f0cf1Property = "20000071" //0x01
	// 鸿雁三键情景开关（0x2000008x）
	Honyar3SceneSwitch005f0cf2PropertyLeft   = "20000081" //0x01, 左键
	Honyar3SceneSwitch005f0cf2PropertyMiddle = "20000082" //0x02, 中键
	Honyar3SceneSwitch005f0cf2PropertyRight  = "20000083" //0x03, 右键
	// 鸿雁智能墙面插座（0x2000009x）
	// 		0x20000091（开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	// 		0x20000092（总耗电量）：string
	// 		0x20000093（电压）：string
	// 		0x20000094（电流）：string
	// 		0x20000095（负载功率）：string
	// 		0x20000096（童锁开关）："00"关，"01"开
	// 		0x20000097（USB开关）："00"关，"01"开
	// 		0x20000098（背光开关）："00"关，"01"开
	// 		0x20000099（断电记忆开关）："00"关，"01"开
	// 		0x2000009A（清除历史电量）：1-清除
	HonyarSocketHY0105PropertyOnOff                    = "20000091"
	HonyarSocketHY0105PropertyCurretSummationDelivered = "20000092"
	HonyarSocketHY0105PropertyRMSVoltage               = "20000093"
	HonyarSocketHY0105PropertyRMSCurrent               = "20000094"
	HonyarSocketHY0105PropertyActivePower              = "20000095"
	HonyarSocketHY0105PropertyChildLock                = "20000096"
	HonyarSocketHY0105PropertyUSB                      = "20000097"
	HonyarSocketHY0105PropertyBackGroundLight          = "20000098"
	HonyarSocketHY0105PropertyPowerOffMemory           = "20000099"
	HonyarSocketHY0105PropertyHistoryElectricClear     = "2000009A"
	// 鸿雁16A智能墙面插座（0x200000Ax）
	// 		0x200000A1（开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	// 		0x200000A2（总耗电量）：string
	// 		0x200000A3（电压）：string
	// 		0x200000A4（电流）：string
	// 		0x200000A5（负载功率）：string
	// 		0x200000A6（童锁开关）："00"关，"01"开
	// 		0x200000A7（背光开关）："00"关，"01"开
	// 		0x200000A8（断电记忆开关）："00"关，"01"开
	// 		0x200000A9（清除历史电量）：1-清除
	HonyarSocketHY0106PropertyOnOff                    = "200000A1"
	HonyarSocketHY0106PropertyCurretSummationDelivered = "200000A2"
	HonyarSocketHY0106PropertyRMSVoltage               = "200000A3"
	HonyarSocketHY0106PropertyRMSCurrent               = "200000A4"
	HonyarSocketHY0106PropertyActivePower              = "200000A5"
	HonyarSocketHY0106PropertyChildLock                = "200000A6"
	HonyarSocketHY0106PropertyBackGroundLight          = "200000A7"
	HonyarSocketHY0106PropertyPowerOffMemory           = "200000A8"
	HonyarSocketHY0106PropertyHistoryElectricClear     = "200000A9"

	// Iotware属性字段名称
	IotwarePropertyOnOff               = "onoff"
	IotwarePropertyElectric            = "electric"
	IotwarePropertyVoltage             = "voltage"
	IotwarePropertyCurrent             = "current"
	IotwarePropertyPower               = "power"
	IotwarePropertyChildLock           = "childLock"
	IotwarePropertyUSBSwitch           = "usbSwitch"
	IotwarePropertyBackLightSwitch     = "backLightSwitch"
	IotwarePropertyMemorySwitch        = "memorySwitch"
	IotwarePropertyHistoryClear        = "historyClear"
	IotwarePropertyOnOffOne            = "onoffOne"
	IotwarePropertyOnOffTwo            = "onoffTwo"
	IotwarePropertyOnOffThree          = "onoffThree"
	IotwarePropertySwitchOne           = "switchOne"
	IotwarePropertySwitchLeft          = "switchLeft"
	IotwarePropertySwitchLeftUp        = "switchLeftUp"
	IotwarePropertySwitchLeftMid       = "switchLeftMid"
	IotwarePropertySwitchLeftDown      = "switchLeftDown"
	IotwarePropertySwitchRight         = "switchRight"
	IotwarePropertySwitchRightUp       = "switchRightUp"
	IotwarePropertySwitchRightMid      = "switchRightMid"
	IotwarePropertySwitchRightDown     = "switchRightDown"
	IotwarePropertySwitchLeftUpLogo    = "switchLeftUpLogo"
	IotwarePropertySwitchLeftMidLogo   = "switchLeftMidLogo"
	IotwarePropertySwitchLeftDownLogo  = "switchLeftDownLogo"
	IotwarePropertySwitchRightUpLogo   = "switchRightUpLogo"
	IotwarePropertySwitchRightMidLogo  = "switchRightMidLogo"
	IotwarePropertySwitchRightDownLogo = "switchRightDownLogo"
	IotwarePropertySwitchMid           = "switchMid"
	IotwarePropertySwitchUp            = "switchUp"
	IotwarePropertySwitchDown          = "switchDown"
	IotwarePropertyAntiDismantle       = "antiDismantle"
	IotwarePropertyTemperature         = "temperature"
	IotwarePropertyHumidity            = "humidity"
	IotwarePropertyHCHO                = "hcho"
	IotwarePropertyPM25                = "pm25"
	IotwarePropertyPM10                = "pm10"
	IotwarePropertyWarning             = "warning"
	IotwarePropertyTWarning            = "tWarning"
	IotwarePropertyHWarning            = "hWarning"
	IotwarePropertyTVOC                = "tvoc"
	IotwarePropertyAQI                 = "aqi"
	IotwarePropertyAirState            = "airState"
	IotwarePropertyHCHOWarning         = "hchoWarning"
	IotwarePropertyPM25Warning         = "pm25Warning"
	IotwarePropertyPM10Warning         = "pm10Warning"
	IotwarePropertyTVOCWarning         = "tvocWarning"
	IotwarePropertyCellWarning         = "cellWarning"
	IotwarePropertyAQIWarning          = "aqiWarning"
	IotwarePropertyDownWarning         = "downWarning"
	IotwarePropertyClearWarning        = "clearWarning"
	IotwarePropertyTThreshold          = "tThreshold"
	IotwarePropertyHThreshold          = "hThreshold"
	IotwarePropertyNotDisturb          = "notDisturb"
	IotwarePropertySendKeyCommand      = "downControl"
	IotwarePropertyStudyKey            = "learnControl"
	IotwarePropertyDeleteKey           = "deleteControl"
	IotwarePropertyCreateID            = "createModelType"
	IotwarePropertyGetIDAndKeyCodeList = "getDevAndCmdList"
	IotwarePropertyLanguage            = "language"
	IotwarePropertyUnitOfTemperature   = "unitOfTemperature"
)

// func getPropertyMeaning(srcString string) string {
// 	srcString = strings.Replace(srcString, HeimanSmartPlugPropertyOnOff, HeimanSmartPlugPropertyOnOff+"(HeimanSmartPlugPropertyOnOff)", 1)
// 	srcString = strings.Replace(srcString, HeimanSmartPlugPropertyPower, HeimanSmartPlugPropertyPower+"(HeimanSmartPlugPropertyPower)", 1)
// 	srcString = strings.Replace(srcString, HeimanSmartPlugPropertyElectric, HeimanSmartPlugPropertyElectric+"(HeimanSmartPlugPropertyElectric)", 1)
// 	srcString = strings.Replace(srcString, HeimanESocketPropertyOnOff, HeimanESocketPropertyOnOff+"(HeimanESocketPropertyOnOff)", 1)
// 	srcString = strings.Replace(srcString, HeimanESocketPropertyPower, HeimanESocketPropertyPower+"(HeimanESocketPropertyPower)", 1)
// 	srcString = strings.Replace(srcString, HeimanESocketPropertyElectric, HeimanESocketPropertyElectric+"(HeimanESocketPropertyElectric)", 1)
// 	srcString = strings.Replace(srcString, HeimanSmokeSensorPropertyAlarm, HeimanSmokeSensorPropertyAlarm+"(HeimanSmokeSensorPropertyAlarm)", 1)
// 	srcString = strings.Replace(srcString, HeimanSmokeSensorPropertyLeftElectric, HeimanSmokeSensorPropertyLeftElectric+"(HeimanSmokeSensorPropertyLeftElectric)", 1)
// 	srcString = strings.Replace(srcString, HeimanWaterSensorPropertyAlarm, HeimanWaterSensorPropertyAlarm+"(HeimanWaterSensorPropertyAlarm)", 1)
// 	srcString = strings.Replace(srcString, HeimanWaterSensorPropertyLeftElectric, HeimanWaterSensorPropertyLeftElectric+"(HeimanWaterSensorPropertyLeftElectric)", 1)
// 	srcString = strings.Replace(srcString, HeimanWaterSensorPropertyTamper, HeimanWaterSensorPropertyTamper+"(HeimanWaterSensorPropertyTamper)", 1)
// 	srcString = strings.Replace(srcString, HeimanPIRSensorPropertyAlarm, HeimanPIRSensorPropertyAlarm+"(HeimanPIRSensorPropertyAlarm)", 1)
// 	srcString = strings.Replace(srcString, HeimanPIRSensorPropertyLeftElectric, HeimanPIRSensorPropertyLeftElectric+"(HeimanPIRSensorPropertyLeftElectric)", 1)
// 	srcString = strings.Replace(srcString, HeimanPIRSensorPropertyTamper, HeimanPIRSensorPropertyTamper+"(HeimanPIRSensorPropertyTamper)", 1)
// 	srcString = strings.Replace(srcString, HeimanHTSensorPropertyTemperature, HeimanHTSensorPropertyTemperature+"(HeimanHTSensorPropertyTemperature)", 1)
// 	srcString = strings.Replace(srcString, HeimanHTSensorPropertyHumidity, HeimanHTSensorPropertyHumidity+"(HeimanHTSensorPropertyHumidity)", 1)
// 	srcString = strings.Replace(srcString, HeimanHTSensorPropertyLeftElectric, HeimanHTSensorPropertyLeftElectric+"(HeimanHTSensorPropertyLeftElectric)", 1)
// 	srcString = strings.Replace(srcString, HeimanCOSensorPropertyAlarm, HeimanCOSensorPropertyAlarm+"(HeimanCOSensorPropertyAlarm)", 1)
// 	srcString = strings.Replace(srcString, HeimanCOSensorPropertyLeftElectric, HeimanCOSensorPropertyLeftElectric+"(HeimanCOSensorPropertyLeftElectric)", 1)
// 	srcString = strings.Replace(srcString, HeimanWarningDevicePropertyAlarmTime, HeimanWarningDevicePropertyAlarmTime+"(HeimanWarningDevicePropertyAlarmTime)", 1)
// 	srcString = strings.Replace(srcString, HeimanWarningDevicePropertyLeftElectric, HeimanWarningDevicePropertyLeftElectric+"(HeimanWarningDevicePropertyLeftElectric)", 1)
// 	srcString = strings.Replace(srcString, HeimanWarningDevicePropertyAlarmStart, HeimanWarningDevicePropertyAlarmStart+"(HeimanWarningDevicePropertyAlarmStart)", 1)
// 	srcString = strings.Replace(srcString, HeimanGASSensorPropertyAlarm, HeimanGASSensorPropertyAlarm+"(HeimanGASSensorPropertyAlarm)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyTemperature, HeimanHS2AQPropertyTemperature+"(HeimanHS2AQPropertyTemperature)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyHumidity, HeimanHS2AQPropertyHumidity+"(HeimanHS2AQPropertyHumidity)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyFormaldehyde, HeimanHS2AQPropertyFormaldehyde+"(HeimanHS2AQPropertyFormaldehyde)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyPM25, HeimanHS2AQPropertyPM25+"(HeimanHS2AQPropertyPM25)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyPM10, HeimanHS2AQPropertyPM10+"(HeimanHS2AQPropertyPM10)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyLeftElectric, HeimanHS2AQPropertyLeftElectric+"(HeimanHS2AQPropertyLeftElectric)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyTemperatureAlarmThreshold, HeimanHS2AQPropertyTemperatureAlarmThreshold+"(HeimanHS2AQPropertyTemperatureAlarmThreshold)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyHumidityAlarmThreshold, HeimanHS2AQPropertyHumidityAlarmThreshold+"(HeimanHS2AQPropertyHumidityAlarmThreshold)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyTemperatureAlarm, HeimanHS2AQPropertyTemperatureAlarm+"(HeimanHS2AQPropertyTemperatureAlarm)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyHumidityAlarm, HeimanHS2AQPropertyHumidityAlarm+"(HeimanHS2AQPropertyHumidityAlarm)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyTVOC, HeimanHS2AQPropertyTVOC+"(HeimanHS2AQPropertyTVOC)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyAQI, HeimanHS2AQPropertyAQI+"(HeimanHS2AQPropertyAQI)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyBatteryState, HeimanHS2AQPropertyBatteryState+"(HeimanHS2AQPropertyBatteryState)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyUnDisturb, HeimanHS2AQPropertyUnDisturb+"(HeimanHS2AQPropertyUnDisturb)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyFormaldehydeAlarm, HeimanHS2AQPropertyFormaldehydeAlarm+"(HeimanHS2AQPropertyFormaldehydeAlarm)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyPM25Alarm, HeimanHS2AQPropertyPM25Alarm+"(HeimanHS2AQPropertyPM25Alarm)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyPM10Alarm, HeimanHS2AQPropertyPM10Alarm+"(HeimanHS2AQPropertyPM10Alarm)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyTVOCAlarm, HeimanHS2AQPropertyTVOCAlarm+"(HeimanHS2AQPropertyTVOCAlarm)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyAQIAlarm, HeimanHS2AQPropertyAQIAlarm+"(HeimanHS2AQPropertyAQIAlarm)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2AQPropertyBatteryAlarm, HeimanHS2AQPropertyBatteryAlarm+"(HeimanHS2AQPropertyBatteryAlarm)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2SW1LEFR30PropertyOnOff1, HeimanHS2SW1LEFR30PropertyOnOff1+"(HeimanHS2SW1LEFR30PropertyOnOff1)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2SW2LEFR30PropertyOnOff1, HeimanHS2SW2LEFR30PropertyOnOff1+"(HeimanHS2SW2LEFR30PropertyOnOff1)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2SW2LEFR30PropertyOnOff2, HeimanHS2SW2LEFR30PropertyOnOff2+"(HeimanHS2SW2LEFR30PropertyOnOff2)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2SW3LEFR30PropertyOnOff1, HeimanHS2SW3LEFR30PropertyOnOff1+"(HeimanHS2SW3LEFR30PropertyOnOff1)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2SW3LEFR30PropertyOnOff2, HeimanHS2SW3LEFR30PropertyOnOff2+"(HeimanHS2SW3LEFR30PropertyOnOff2)", 1)
// 	srcString = strings.Replace(srcString, HeimanHS2SW3LEFR30PropertyOnOff3, HeimanHS2SW3LEFR30PropertyOnOff3+"(HeimanHS2SW3LEFR30PropertyOnOff3)", 1)
// 	srcString = strings.Replace(srcString, HeimanSceneSwitchEM30PropertyAtHome, HeimanSceneSwitchEM30PropertyAtHome+"(HeimanSceneSwitchEM30PropertyAtHome)", 1)
// 	srcString = strings.Replace(srcString, HeimanSceneSwitchEM30PropertyGoOut, HeimanSceneSwitchEM30PropertyGoOut+"(HeimanSceneSwitchEM30PropertyGoOut)", 1)
// 	srcString = strings.Replace(srcString, HeimanSceneSwitchEM30PropertyCinema, HeimanSceneSwitchEM30PropertyCinema+"(HeimanSceneSwitchEM30PropertyCinema)", 1)
// 	srcString = strings.Replace(srcString, HeimanSceneSwitchEM30PropertyRepast, HeimanSceneSwitchEM30PropertyRepast+"(HeimanSceneSwitchEM30PropertyRepast)", 1)
// 	srcString = strings.Replace(srcString, HeimanSceneSwitchEM30PropertySleep, HeimanSceneSwitchEM30PropertySleep+"(HeimanSceneSwitchEM30PropertySleep)", 1)
// 	srcString = strings.Replace(srcString, HeimanSceneSwitchEM30PropertyLeftElectric, HeimanSceneSwitchEM30PropertyLeftElectric+"(HeimanSceneSwitchEM30PropertyLeftElectric)", 1)
// 	srcString = strings.Replace(srcString, HeimanIRControlEMPropertySendKeyCommand, HeimanIRControlEMPropertySendKeyCommand+"(HeimanIRControlEMPropertySendKeyCommand)", 1)
// 	srcString = strings.Replace(srcString, HeimanIRControlEMPropertyStudyKey, HeimanIRControlEMPropertyStudyKey+"(HeimanIRControlEMPropertyStudyKey)", 1)
// 	srcString = strings.Replace(srcString, HeimanIRControlEMPropertyDeleteKey, HeimanIRControlEMPropertyDeleteKey+"(HeimanIRControlEMPropertyDeleteKey)", 1)
// 	srcString = strings.Replace(srcString, HeimanIRControlEMPropertyCreateID, HeimanIRControlEMPropertyCreateID+"(HeimanIRControlEMPropertyCreateID)", 1)
// 	srcString = strings.Replace(srcString, HeimanIRControlEMPropertyGetIDAndKeyCodeList, HeimanIRControlEMPropertyGetIDAndKeyCodeList+"(HeimanIRControlEMPropertyGetIDAndKeyCodeList)", 1)
// 	srcString = strings.Replace(srcString, HonyarSingleSwitch00500c32PropertyOnOff1, HonyarSingleSwitch00500c32PropertyOnOff1+"(HonyarSingleSwitch00500c32PropertyOnOff1)", 1)
// 	srcString = strings.Replace(srcString, HonyarDoubleSwitch00500c33PropertyOnOff1, HonyarDoubleSwitch00500c33PropertyOnOff1+"(HonyarDoubleSwitch00500c33PropertyOnOff1)", 1)
// 	srcString = strings.Replace(srcString, HonyarDoubleSwitch00500c33PropertyOnOff2, HonyarDoubleSwitch00500c33PropertyOnOff2+"(HonyarDoubleSwitch00500c33PropertyOnOff2)", 1)
// 	srcString = strings.Replace(srcString, HonyarTripleSwitch00500c35PropertyOnOff1, HonyarTripleSwitch00500c35PropertyOnOff1+"(HonyarTripleSwitch00500c35PropertyOnOff1)", 1)
// 	srcString = strings.Replace(srcString, HonyarTripleSwitch00500c35PropertyOnOff2, HonyarTripleSwitch00500c35PropertyOnOff2+"(HonyarTripleSwitch00500c35PropertyOnOff2)", 1)
// 	srcString = strings.Replace(srcString, HonyarTripleSwitch00500c35PropertyOnOff3, HonyarTripleSwitch00500c35PropertyOnOff3+"(HonyarTripleSwitch00500c35PropertyOnOff3)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocket000a0c3cPropertyOnOff, HonyarSocket000a0c3cPropertyOnOff+"(HonyarSocket000a0c3cPropertyOnOff)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocket000a0c3cPropertyCurretSummationDelivered, HonyarSocket000a0c3cPropertyCurretSummationDelivered+"(HonyarSocket000a0c3cPropertyCurretSummationDelivered)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocket000a0c3cPropertyRMSVoltage, HonyarSocket000a0c3cPropertyRMSVoltage+"(HonyarSocket000a0c3cPropertyRMSVoltage)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocket000a0c3cPropertyRMSCurrent, HonyarSocket000a0c3cPropertyRMSCurrent+"(HonyarSocket000a0c3cPropertyRMSCurrent)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocket000a0c3cPropertyActivePower, HonyarSocket000a0c3cPropertyActivePower+"(HonyarSocket000a0c3cPropertyActivePower)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocket000a0c3cPropertyChildLock, HonyarSocket000a0c3cPropertyChildLock+"(HonyarSocket000a0c3cPropertyChildLock)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocket000a0c3cPropertyUSB, HonyarSocket000a0c3cPropertyUSB+"(HonyarSocket000a0c3cPropertyUSB)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0105PropertyOnOff, HonyarSocketHY0105PropertyOnOff+"(HonyarSocketHY0105PropertyOnOff)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0105PropertyCurretSummationDelivered, HonyarSocketHY0105PropertyCurretSummationDelivered+"(HonyarSocketHY0105PropertyCurretSummationDelivered)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0105PropertyRMSVoltage, HonyarSocketHY0105PropertyRMSVoltage+"(HonyarSocketHY0105PropertyRMSVoltage)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0105PropertyRMSCurrent, HonyarSocketHY0105PropertyRMSCurrent+"(HonyarSocketHY0105PropertyRMSCurrent)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0105PropertyActivePower, HonyarSocketHY0105PropertyActivePower+"(HonyarSocketHY0105PropertyActivePower)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0105PropertyChildLock, HonyarSocketHY0105PropertyChildLock+"(HonyarSocketHY0105PropertyChildLock)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0105PropertyUSB, HonyarSocketHY0105PropertyUSB+"(HonyarSocketHY0105PropertyUSB)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0105PropertyBackGroundLight, HonyarSocketHY0105PropertyBackGroundLight+"(HonyarSocketHY0105PropertyBackGroundLight)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0105PropertyPowerOffMemory, HonyarSocketHY0105PropertyPowerOffMemory+"(HonyarSocketHY0105PropertyPowerOffMemory)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0105PropertyHistoryElectricClear, HonyarSocketHY0105PropertyHistoryElectricClear+"(HonyarSocketHY0105PropertyHistoryElectricClear)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0106PropertyOnOff, HonyarSocketHY0106PropertyOnOff+"(HonyarSocketHY0106PropertyOnOff)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0106PropertyCurretSummationDelivered, HonyarSocketHY0106PropertyCurretSummationDelivered+"(HonyarSocketHY0106PropertyCurretSummationDelivered)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0106PropertyRMSVoltage, HonyarSocketHY0106PropertyRMSVoltage+"(HonyarSocketHY0106PropertyRMSVoltage)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0106PropertyRMSCurrent, HonyarSocketHY0106PropertyRMSCurrent+"(HonyarSocketHY0106PropertyRMSCurrent)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0106PropertyActivePower, HonyarSocketHY0106PropertyActivePower+"(HonyarSocketHY0106PropertyActivePower)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0106PropertyChildLock, HonyarSocketHY0106PropertyChildLock+"(HonyarSocketHY0106PropertyChildLock)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0106PropertyBackGroundLight, HonyarSocketHY0106PropertyBackGroundLight+"(HonyarSocketHY0106PropertyBackGroundLight)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0106PropertyPowerOffMemory, HonyarSocketHY0106PropertyPowerOffMemory+"(HonyarSocketHY0106PropertyPowerOffMemory)", 1)
// 	srcString = strings.Replace(srcString, HonyarSocketHY0106PropertyHistoryElectricClear, HonyarSocketHY0106PropertyHistoryElectricClear+"(HonyarSocketHY0106PropertyHistoryElectricClear)", 1)
// 	srcString = strings.Replace(srcString, Honyar1SceneSwitch005f0cf1Property, Honyar1SceneSwitch005f0cf1Property+"(Honyar1SceneSwitch005f0cf1Property)", 1)
// 	srcString = strings.Replace(srcString, Honyar2SceneSwitch005f0cf3PropertyLeft, Honyar2SceneSwitch005f0cf3PropertyLeft+"(Honyar2SceneSwitch005f0cf3PropertyLeft)", 1)
// 	srcString = strings.Replace(srcString, Honyar2SceneSwitch005f0cf3PropertyRight, Honyar2SceneSwitch005f0cf3PropertyRight+"(Honyar2SceneSwitch005f0cf3PropertyRight)", 1)
// 	srcString = strings.Replace(srcString, Honyar3SceneSwitch005f0cf2PropertyLeft, Honyar3SceneSwitch005f0cf2PropertyLeft+"(Honyar3SceneSwitch005f0cf2PropertyLeft)", 1)
// 	srcString = strings.Replace(srcString, Honyar3SceneSwitch005f0cf2PropertyMiddle, Honyar3SceneSwitch005f0cf2PropertyMiddle+"(Honyar3SceneSwitch005f0cf2PropertyMiddle)", 1)
// 	srcString = strings.Replace(srcString, Honyar3SceneSwitch005f0cf2PropertyRight, Honyar3SceneSwitch005f0cf2PropertyRight+"(Honyar3SceneSwitch005f0cf2PropertyRight)", 1)
// 	srcString = strings.Replace(srcString, Honyar6SceneSwitch005f0c3bPropertyLeftUp, Honyar6SceneSwitch005f0c3bPropertyLeftUp+"(Honyar6SceneSwitch005f0c3bPropertyLeftUp)", 1)
// 	srcString = strings.Replace(srcString, Honyar6SceneSwitch005f0c3bPropertyLeftMiddle, Honyar6SceneSwitch005f0c3bPropertyLeftMiddle+"(Honyar6SceneSwitch005f0c3bPropertyLeftMiddle)", 1)
// 	srcString = strings.Replace(srcString, Honyar6SceneSwitch005f0c3bPropertyLeftDown, Honyar6SceneSwitch005f0c3bPropertyLeftDown+"(Honyar6SceneSwitch005f0c3bPropertyLeftDown)", 1)
// 	srcString = strings.Replace(srcString, Honyar6SceneSwitch005f0c3bPropertyRightUp, Honyar6SceneSwitch005f0c3bPropertyRightUp+"(Honyar6SceneSwitch005f0c3bPropertyRightUp)", 1)
// 	srcString = strings.Replace(srcString, Honyar6SceneSwitch005f0c3bPropertyRightMiddle, Honyar6SceneSwitch005f0c3bPropertyRightMiddle+"(Honyar6SceneSwitch005f0c3bPropertyRightMiddle)", 1)
// 	srcString = strings.Replace(srcString, Honyar6SceneSwitch005f0c3bPropertyRightDown, Honyar6SceneSwitch005f0c3bPropertyRightDown+"(Honyar6SceneSwitch005f0c3bPropertyRightDown)", 1)
// 	srcString = strings.Replace(srcString, Honyar6SceneSwitch005f0c3bPropertyLeftUpAddScene, Honyar6SceneSwitch005f0c3bPropertyLeftUpAddScene+"(Honyar6SceneSwitch005f0c3bPropertyLeftUpAddScene)", 1)
// 	srcString = strings.Replace(srcString, Honyar6SceneSwitch005f0c3bPropertyLeftMiddleAddScene, Honyar6SceneSwitch005f0c3bPropertyLeftMiddleAddScene+"(Honyar6SceneSwitch005f0c3bPropertyLeftMiddleAddScene)", 1)
// 	srcString = strings.Replace(srcString, Honyar6SceneSwitch005f0c3bPropertyLeftDownAddScene, Honyar6SceneSwitch005f0c3bPropertyLeftDownAddScene+"(Honyar6SceneSwitch005f0c3bPropertyLeftDownAddScene)", 1)
// 	srcString = strings.Replace(srcString, Honyar6SceneSwitch005f0c3bPropertyRightUpAddScene, Honyar6SceneSwitch005f0c3bPropertyRightUpAddScene+"(Honyar6SceneSwitch005f0c3bPropertyRightUpAddScene)", 1)
// 	srcString = strings.Replace(srcString, Honyar6SceneSwitch005f0c3bPropertyRightMiddleAddScene, Honyar6SceneSwitch005f0c3bPropertyRightMiddleAddScene+"(Honyar6SceneSwitch005f0c3bPropertyRightMiddleAddScene)", 1)
// 	srcString = strings.Replace(srcString, Honyar6SceneSwitch005f0c3bPropertyRightDownAddScene, Honyar6SceneSwitch005f0c3bPropertyRightDownAddScene+"(Honyar6SceneSwitch005f0c3bPropertyRightDownAddScene)", 1)
// 	return srcString
// }

// 定义mqtt消息结构体
type mqttMsgSt struct {
	ID      interface{} `json:"id"`
	Version string      `json:"version"`
	Params  interface{} `json:"params"`
	Method  string      `json:"method"`
}

// Publish publish mqtt msg
func Publish(topic string, method string, params interface{}, msgID interface{}) {
	jsonMsg, err := json.Marshal(mqttMsgSt{
		ID:      msgID,
		Version: "1.0",
		Params:  params,
		Method:  method,
	})
	if err != nil {
		globallogger.Log.Errorf("[Publish]: JSON marshaling failed: %s", err)
	}
	mqtt.Publish(topic, string(jsonMsg))
	// globallogger.Log.Warnf("[Publish]: mqtt msg publish: topic: %s, mqttMsg %v", topic, getPropertyMeaning(string(jsonMsg)))
	globallogger.Log.Warnf("[Publish]: mqtt msg publish: topic: %s, mqttMsg %s", topic, string(jsonMsg))
}

// PublishRPCRspIotware publish rpc rsp
func PublishRPCRspIotware(devEUI string, result string, msgID interface{}) {
	jsonMsg, err := json.Marshal(publicstruct.RPCRspIotware{
		Device: devEUI,
		ID:     msgID,
		Data: publicstruct.DataRPCRsp{
			Result: result,
		},
	})
	if err != nil {
		globallogger.Log.Errorf("[Publish]: JSON marshaling failed: %s", err)
	}
	PublishIotware(TopicV1GatewayRPC, jsonMsg)
}

// PublishTelemetryUpIotware publish telemetry msg
func PublishTelemetryUpIotware(terminalInfo config.TerminalInfo, values interface{}) {
	if terminalInfo.IsExist {
		params := make(map[string]interface{}, 1)
		var typeBuilder strings.Builder
		typeBuilder.WriteString(terminalInfo.ManufacturerName)
		typeBuilder.WriteString("-")
		typeBuilder.WriteString(terminalInfo.TmnType)
		params[terminalInfo.DevEUI] = publicstruct.DeviceData{
			Type: typeBuilder.String(),
			Data: []publicstruct.DataItem{
				{
					TimeStamp: time.Now().UnixNano() / 1e6,
					Values:    values,
				},
			},
		}
		jsonMsg, err := json.Marshal(publicstruct.TelemetryIotware{
			Gateway: publicstruct.GatewayIotware{
				FirstAddr:  terminalInfo.FirstAddr,
				SecondAddr: terminalInfo.SecondAddr,
				ThirdAddr:  terminalInfo.ThirdAddr,
				PortID:     terminalInfo.ModuleID,
			},
			DeviceData: params,
		})
		if err != nil {
			globallogger.Log.Errorf("[Publish]: JSON marshaling failed: %s", err)
		}
		PublishIotware(TopicV1GatewayTelemetry, jsonMsg)
		params = nil
	}
}

// PublishIotware publish mqtt msg
func PublishIotware(topic string, jsonMsg []byte) {
	mqtt.Publish(topic, string(jsonMsg))
	globallogger.Log.Warnf("[PublishIotware]: mqtt msg publish: topic: %s, mqttMsg %v", topic, string(jsonMsg))
}

// ProcSubMsg process subscribe mqtt msg
func ProcSubMsg(topic string, jsonMsg []byte) {
	defer func() {
		err := recover()
		if err != nil {
			globallogger.Log.Errorln("ProcSubMsg err :", err)
		}
	}()
	if constant.Constant.Iotware {
		switch topic {
		case TopicV1GatewayRPC:
			var mqttMsg publicstruct.RPCIotware
			err := json.Unmarshal(jsonMsg, &mqttMsg)
			if err != nil {
				globallogger.Log.Errorln("[Subscribe]: JSON marshaling failed:", err)
				return
			}
			globallogger.Log.Warnf("[Subscribe]: mqtt msg subscribe: topic: %s, mqttMsg %+v", topic, mqttMsg)
			procPropertyIotware(mqttMsg)
		case TopicV1GatewayNetworkInAck:
			var networkInAckMsg publicstruct.NetworkInAckIotware
			err := json.Unmarshal(jsonMsg, &networkInAckMsg)
			if err != nil {
				globallogger.Log.Errorln("[Subscribe]: JSON marshaling failed:", err)
				return
			}
			globallogger.Log.Warnf("[Subscribe]: mqtt msg subscribe: topic: %s, networkInAckMsg %+v", topic, networkInAckMsg)
			procActionIotware(networkInAckMsg)
		case TopicV1GatewayEventDeviceDelete:
			var deviceDeleteMsg publicstruct.EventDeviceDeleteIotware
			err := json.Unmarshal(jsonMsg, &deviceDeleteMsg)
			if err != nil {
				globallogger.Log.Errorln("[Subscribe]: JSON marshaling failed:", err)
				return
			}
			globallogger.Log.Warnf("[Subscribe]: mqtt msg subscribe: topic: %s, deviceDeleteMsg %+v", topic, deviceDeleteMsg)
			for _, devEUI := range deviceDeleteMsg.DeviceList {
				publicfunction.ProcTerminalDelete(devEUI)
			}
		default:
			globallogger.Log.Warnf("[Subscribe]: invalid topic: %s", topic)
		}
	} else if constant.Constant.Iotedge {
		var mqttMsg mqttMsgSt
		err := json.Unmarshal(jsonMsg, &mqttMsg)
		if err != nil {
			globallogger.Log.Errorln("[Subscribe]: JSON marshaling failed:", err)
			return
		}
		globallogger.Log.Warnf("[Subscribe]: mqtt msg subscribe: topic: %s, mqttMsg %+v", topic, mqttMsg)

		switch topic {
		case TopicIotsmartspaceZigbeeserverProperty:
			procProperty(mqttMsg)
		case TopicIotsmartspaceZigbeeserverState:
			procState(mqttMsg)
		case TopicIotsmartspaceZigbeeserverAction:
			procAction(mqttMsg)
		case "scenes":
			procScenes(mqttMsg)
		case common.ReadAttribute:
			procReadAttribute(mqttMsg)
		default:
			globallogger.Log.Warnf("[Subscribe]: invalid topic: %s", topic)
		}
	}
}
