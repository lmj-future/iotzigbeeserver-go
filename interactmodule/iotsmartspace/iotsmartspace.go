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
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
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
	HeimanSmartPlugPropertyKeyOnOff = "10000011"
	// HeimanSmartPlugPropertyPower    = "10000012"
	// HeimanSmartPlugPropertyElectric = "10000013"
	// 海曼智能墙面插座（0x1000002x）
	// 		0x10000021（开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	// 		0x10000022（功率）：string
	// 		0x10000023（电量）：string
	HeimanESocketPropertyKeyOnOff = "10000021"
	// HeimanESocketPropertyPower    = "10000022"
	// HeimanESocketPropertyElectric = "10000023"
	// 海曼火灾探测器（0x1000003x）
	// 		0x10000031（告警）："0"告警上报，"1"告警解除
	// 		0x10000032（剩余电量）：string，百分比
	// HeimanSmokeSensorPropertyAlarm = "10000031"
	// HeimanSmokeSensorPropertyLeftElectric = "10000032"
	// 海曼水浸传感器（0x1000004x）
	// 		0x10000041（告警）："0"告警上报，"1"告警解除
	// 		0x10000042（剩余电量）：string，百分比
	// 		0x10000043（防拆）："0"设备拆除，"1"设备拆除解除
	// HeimanWaterSensorPropertyAlarm = "10000041"
	// HeimanWaterSensorPropertyLeftElectric = "10000042"
	// HeimanWaterSensorPropertyTamper = "10000043"
	// 海曼人体传感器（0x1000005x）
	// 		0x10000051（人体检测）："0"无人移动，"1"有人移动
	// 		0x10000052（剩余电量）：string，百分比
	// 		0x10000053（防拆）："0"设备拆除，"1"设备拆除解除
	// HeimanPIRSensorPropertyAlarm = "10000051"
	// HeimanPIRSensorPropertyLeftElectric = "10000052"
	// HeimanPIRSensorPropertyTamper = "10000053"
	// 海曼温湿度传感器（0x1000006x）
	// 		0x10000061（温度）：[-10, 60]，string
	// 		0x10000062（湿度）：[0, 100]，string
	// 		0x10000063（剩余电量）：string，百分比
	// HeimanHTSensorPropertyTemperature  = "10000061"
	// HeimanHTSensorPropertyHumidity     = "10000062"
	// HeimanHTSensorPropertyLeftElectric = "10000063"
	// 海曼一氧化碳探测器（0x1000007x）
	// 		0x10000071（告警）："0"告警上报，"1"告警解除
	// 		0x10000072（剩余电量）：string，百分比
	// HeimanCOSensorPropertyAlarm = "10000071"
	// HeimanCOSensorPropertyLeftElectric = "10000072"
	// 海曼智能声光报警器（0x1000008x）
	// 		0x10000081（报警时长设置）：string
	// 		0x10000082（剩余电量）：string，百分比
	// 		0x10000083（下发警报）："0"触发警报，"1"解除警报
	HeimanWarningDevicePropertyAlarmTime = "10000081"
	// HeimanWarningDevicePropertyLeftElectric = "10000082"
	HeimanWarningDevicePropertyAlarmStart = "10000083"
	// 海曼可燃气体探测器（0x1000009x）
	// 		0x10000091（告警）："0"告警上报，"1"告警解除
	// HeimanGASSensorPropertyAlarm = "10000091"
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
	// HeimanHS2AQPropertyTemperature               = "100000A1"
	// HeimanHS2AQPropertyHumidity                  = "100000A2"
	// HeimanHS2AQPropertyFormaldehyde              = "100000A3"
	// HeimanHS2AQPropertyPM25                      = "100000A4"
	// HeimanHS2AQPropertyPM10                      = "100000A5"
	// HeimanHS2AQPropertyLeftElectric              = "100000A6"
	HeimanHS2AQPropertyKeyTemperatureAlarmThreshold = "100000A7"
	HeimanHS2AQPropertyKeyHumidityAlarmThreshold    = "100000A8"
	// HeimanHS2AQPropertyTemperatureAlarm          = "100000A9"
	// HeimanHS2AQPropertyHumidityAlarm             = "100000AA"
	// HeimanHS2AQPropertyTVOC         = "100000AB"
	// HeimanHS2AQPropertyAQI          = "100000AC"
	// HeimanHS2AQPropertyBatteryState = "100000AD"
	HeimanHS2AQPropertyKeyUnDisturb = "100000AE"
	// HeimanHS2AQPropertyFormaldehydeAlarm         = "100000AF"
	// HeimanHS2AQPropertyPM25Alarm                 = "100000B1"
	// HeimanHS2AQPropertyPM10Alarm                 = "100000B2"
	// HeimanHS2AQPropertyTVOCAlarm                 = "100000B3"
	// HeimanHS2AQPropertyAQIAlarm                  = "100000B4"
	// HeimanHS2AQPropertyBatteryAlarm              = "100000B5"
	// 海曼智能开关（0x100000Cx）
	// 		0x100000Cx（一路开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	HeimanHS2SW1LEFR30PropertyKeyOnOff1 = "100000C1"
	// 		0x100000Dx（二路开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	HeimanHS2SW2LEFR30PropertyKeyOnOff1 = "100000D1"
	HeimanHS2SW2LEFR30PropertyKeyOnOff2 = "100000D2"
	// 		0x100000Ex（三路开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	HeimanHS2SW3LEFR30PropertyKeyOnOff1 = "100000E1"
	HeimanHS2SW3LEFR30PropertyKeyOnOff2 = "100000E2"
	HeimanHS2SW3LEFR30PropertyKeyOnOff3 = "100000E3"
	// 海曼智能情景开关（0x100000Fx）
	// HeimanSceneSwitchEM30PropertyAtHome = "100000F2" //0xf0, 在家模式
	// HeimanSceneSwitchEM30PropertyGoOut  = "100000F1" //0xf1, 外出模式
	// HeimanSceneSwitchEM30PropertyCinema = "100000F3" //0xf2, 影院模式
	// HeimanSceneSwitchEM30PropertyRepast = "100000F6" //0xf3, 就餐模式
	// HeimanSceneSwitchEM30PropertySleep  = "100000F4" //0xf4, 就寝模式
	// HeimanSceneSwitchEM30PropertyLeftElectric = "100000F5" //剩余电量,string，百分比
	// 海曼智能红外转发器（0x1000011x）
	HeimanIRControlEMPropertyKeySendKeyCommand      = "10000111" //0xf0, 发送遥控指令
	HeimanIRControlEMPropertyKeyStudyKey            = "10000112" //0xf1, 学习遥控指令
	HeimanIRControlEMPropertyKeyDeleteKey           = "10000113" //0xf3, 删除遥控指令
	HeimanIRControlEMPropertyKeyCreateID            = "10000114" //0xf4, 创建遥控设备
	HeimanIRControlEMPropertyKeyGetIDAndKeyCodeList = "10000115" //0xf6, 获取遥控设备及其指令列表

	// 鸿雁智能开关（0x2000001x）
	// 		0x2000001x（一路开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	HonyarSingleSwitch00500c32PropertyKeyOnOff1 = "20000011"
	// 		0x2000002x（二路开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	HonyarDoubleSwitch00500c33PropertyKeyOnOff1 = "20000021"
	HonyarDoubleSwitch00500c33PropertyKeyOnOff2 = "20000022"
	// 		0x2000003x（三路开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	HonyarTripleSwitch00500c35PropertyKeyOnOff1 = "20000031"
	HonyarTripleSwitch00500c35PropertyKeyOnOff2 = "20000032"
	HonyarTripleSwitch00500c35PropertyKeyOnOff3 = "20000033"
	// 鸿雁智能墙面插座（0x2000004x）
	// 		0x20000041（开关）："00"1号开关关，"01"1号开关开，"10"2号开关关，"11"2号开关开......
	// 		0x20000042（总耗电量）：string
	// 		0x20000043（电压）：string
	// 		0x20000044（电流）：string
	// 		0x20000045（负载功率）：string
	// 		0x20000046（童锁开关）："00"关，"01"开
	// 		0x20000047（USB开关）："00"关，"01"开
	HonyarSocket000a0c3cPropertyKeyOnOff = "20000041"
	// HonyarSocket000a0c3cPropertyCurretSummationDelivered = "20000042"
	// HonyarSocket000a0c3cPropertyRMSVoltage               = "20000043"
	// HonyarSocket000a0c3cPropertyRMSCurrent               = "20000044"
	// HonyarSocket000a0c3cPropertyActivePower              = "20000045"
	HonyarSocket000a0c3cPropertyKeyChildLock = "20000046"
	HonyarSocket000a0c3cPropertyKeyUSB       = "20000047"
	// 鸿雁二键情景开关（0x2000005x）
	// Honyar2SceneSwitch005f0cf3PropertyLeft  = "20000051" //0x01, 左键
	// Honyar2SceneSwitch005f0cf3PropertyRight = "20000052" //0x02, 右键
	// 鸿雁六键情景开关（0x2000006x）
	// Honyar6SceneSwitch005f0c3bPropertyLeftUp              = "20000061" //0x01, 左上键
	// Honyar6SceneSwitch005f0c3bPropertyLeftMiddle          = "20000062" //0x02, 左中键
	// Honyar6SceneSwitch005f0c3bPropertyLeftDown            = "20000063" //0x03, 左下键
	// Honyar6SceneSwitch005f0c3bPropertyRightUp             = "20000064" //0x04, 右上键
	// Honyar6SceneSwitch005f0c3bPropertyRightMiddle         = "20000065" //0x05, 右中键
	// Honyar6SceneSwitch005f0c3bPropertyRightDown           = "20000066" //0x06, 右下键
	Honyar6SceneSwitch005f0c3bPropertyKeyLeftUpAddScene      = "20000067" //0x01~0x21，对应1~33个场景
	Honyar6SceneSwitch005f0c3bPropertyKeyLeftMiddleAddScene  = "20000068" //0x01~0x21，对应1~33个场景
	Honyar6SceneSwitch005f0c3bPropertyKeyLeftDownAddScene    = "20000069" //0x01~0x21，对应1~33个场景
	Honyar6SceneSwitch005f0c3bPropertyKeyRightUpAddScene     = "2000006A" //0x01~0x21，对应1~33个场景
	Honyar6SceneSwitch005f0c3bPropertyKeyRightMiddleAddScene = "2000006B" //0x01~0x21，对应1~33个场景
	Honyar6SceneSwitch005f0c3bPropertyKeyRightDownAddScene   = "2000006C" //0x01~0x21，对应1~33个场景
	// 鸿雁一键情景开关（0x2000007x）
	// Honyar1SceneSwitch005f0cf1Property = "20000071" //0x01
	// 鸿雁三键情景开关（0x2000008x）
	// Honyar3SceneSwitch005f0cf2PropertyLeft   = "20000081" //0x01, 左键
	// Honyar3SceneSwitch005f0cf2PropertyMiddle = "20000082" //0x02, 中键
	// Honyar3SceneSwitch005f0cf2PropertyRight  = "20000083" //0x03, 右键
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
	HonyarSocketHY0105PropertyKeyOnOff = "20000091"
	// HonyarSocketHY0105PropertyCurretSummationDelivered = "20000092"
	// HonyarSocketHY0105PropertyRMSVoltage               = "20000093"
	// HonyarSocketHY0105PropertyRMSCurrent               = "20000094"
	// HonyarSocketHY0105PropertyActivePower              = "20000095"
	HonyarSocketHY0105PropertyKeyChildLock            = "20000096"
	HonyarSocketHY0105PropertyKeyUSB                  = "20000097"
	HonyarSocketHY0105PropertyKeyBackGroundLight      = "20000098"
	HonyarSocketHY0105PropertyKeyPowerOffMemory       = "20000099"
	HonyarSocketHY0105PropertyKeyHistoryElectricClear = "2000009A"
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
	HonyarSocketHY0106PropertyKeyOnOff = "200000A1"
	// HonyarSocketHY0106PropertyCurretSummationDelivered = "200000A2"
	// HonyarSocketHY0106PropertyRMSVoltage               = "200000A3"
	// HonyarSocketHY0106PropertyRMSCurrent               = "200000A4"
	// HonyarSocketHY0106PropertyActivePower              = "200000A5"
	HonyarSocketHY0106PropertyKeyChildLock            = "200000A6"
	HonyarSocketHY0106PropertyKeyBackGroundLight      = "200000A7"
	HonyarSocketHY0106PropertyKeyPowerOffMemory       = "200000A8"
	HonyarSocketHY0106PropertyKeyHistoryElectricClear = "200000A9"

	// Iotware属性字段名称
	// IotwarePropertyOnOff               = "onoff"
	// IotwarePropertyElectric            = "electric"
	// IotwarePropertyVoltage             = "voltage"
	// IotwarePropertyCurrent             = "current"
	// IotwarePropertyPower               = "power"
	// IotwarePropertyChildLock           = "childLock"
	// IotwarePropertyUSBSwitch           = "usbSwitch"
	// IotwarePropertyBackLightSwitch     = "backLightSwitch"
	// IotwarePropertyMemorySwitch        = "memorySwitch"
	// IotwarePropertyHistoryClear        = "historyClear"
	// IotwarePropertyOnOffOne            = "onoffOne"
	// IotwarePropertyOnOffTwo            = "onoffTwo"
	// IotwarePropertyOnOffThree          = "onoffThree"
	// IotwarePropertySwitchOne           = "switchOne"
	// IotwarePropertySwitchLeft          = "switchLeft"
	// IotwarePropertySwitchLeftUp        = "switchLeftUp"
	// IotwarePropertySwitchLeftMid       = "switchLeftMid"
	// IotwarePropertySwitchLeftDown      = "switchLeftDown"
	// IotwarePropertySwitchRight         = "switchRight"
	// IotwarePropertySwitchRightUp       = "switchRightUp"
	// IotwarePropertySwitchRightMid      = "switchRightMid"
	// IotwarePropertySwitchRightDown     = "switchRightDown"
	// IotwarePropertySwitchLeftUpLogo    = "switchLeftUpLogo"
	// IotwarePropertySwitchLeftMidLogo   = "switchLeftMidLogo"
	// IotwarePropertySwitchLeftDownLogo  = "switchLeftDownLogo"
	// IotwarePropertySwitchRightUpLogo   = "switchRightUpLogo"
	// IotwarePropertySwitchRightMidLogo  = "switchRightMidLogo"
	// IotwarePropertySwitchRightDownLogo = "switchRightDownLogo"
	// IotwarePropertySwitchMid           = "switchMid"
	// IotwarePropertySwitchUp            = "switchUp"
	// IotwarePropertySwitchDown          = "switchDown"
	// IotwarePropertyAntiDismantle       = "antiDismantle"
	// IotwarePropertyTemperature         = "temperature"
	// IotwarePropertyHumidity            = "humidity"
	// IotwarePropertyHCHO                = "hcho"
	// IotwarePropertyPM25                = "pm25"
	// IotwarePropertyPM10                = "pm10"
	// IotwarePropertyWarning             = "warning"
	// IotwarePropertyTWarning            = "tWarning"
	// IotwarePropertyHWarning            = "hWarning"
	// IotwarePropertyTVOC                = "tvoc"
	// IotwarePropertyAQI                 = "aqi"
	// IotwarePropertyAirState            = "airState"
	// IotwarePropertyHCHOWarning         = "hchoWarning"
	// IotwarePropertyPM25Warning         = "pm25Warning"
	// IotwarePropertyPM10Warning         = "pm10Warning"
	// IotwarePropertyTVOCWarning         = "tvocWarning"
	// IotwarePropertyCellWarning         = "cellWarning"
	// IotwarePropertyAQIWarning          = "aqiWarning"
	// IotwarePropertyDownWarning         = "downWarning"
	// IotwarePropertyClearWarning        = "clearWarning"
	// IotwarePropertyTThreshold          = "tThreshold"
	// IotwarePropertyHThreshold          = "hThreshold"
	// IotwarePropertyNotDisturb          = "notDisturb"
	// IotwarePropertySendKeyCommand      = "downControl"
	// IotwarePropertyStudyKey            = "learnControl"
	// IotwarePropertyDeleteKey           = "deleteControl"
	// IotwarePropertyCreateID            = "createModelType"
	// IotwarePropertyGetIDAndKeyCodeList = "getDevAndCmdList"
	// IotwarePropertyLanguage            = "language"
	// IotwarePropertyUnitOfTemperature   = "unitOfTemperature"
)

// 定义Iotedge终端属性结构体
type HeimanSmartPlugPropertyOnOff struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"10000011"`
}
type HeimanSmartPlugPropertyOnOffDefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"10000011"`
}
type HeimanSmartPlugPropertyPower struct {
	DevEUI string `json:"terminalId"`
	Power  string `json:"10000012"`
}
type HeimanSmartPlugPropertyElectric struct {
	DevEUI   string `json:"terminalId"`
	Electric string `json:"10000013"`
}
type HeimanESocketPropertyOnOff struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"10000021"`
}
type HeimanESocketPropertyOnOffDefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"10000021"`
}
type HeimanESocketPropertyPower struct {
	DevEUI string `json:"terminalId"`
	Power  string `json:"10000022"`
}
type HeimanESocketPropertyElectric struct {
	DevEUI   string `json:"terminalId"`
	Electric string `json:"10000023"`
}
type HeimanSmokeSensorPropertyAlarm struct {
	DevEUI string `json:"terminalId"`
	Alarm  string `json:"10000031"`
}
type HeimanSmokeSensorPropertyLeftElectric struct {
	DevEUI       string `json:"terminalId"`
	LeftElectric uint64 `json:"10000032"`
}
type HeimanWaterSensorPropertyAlarm struct {
	DevEUI string `json:"terminalId"`
	Alarm  string `json:"10000041"`
}
type HeimanWaterSensorPropertyLeftElectric struct {
	DevEUI       string `json:"terminalId"`
	LeftElectric uint64 `json:"10000042"`
}
type HeimanWaterSensorPropertyTamper struct {
	DevEUI string `json:"terminalId"`
	Tamper string `json:"10000043"`
}
type HeimanPIRSensorPropertyAlarm struct {
	DevEUI string `json:"terminalId"`
	Alarm  string `json:"10000051"`
}
type HeimanPIRSensorPropertyLeftElectric struct {
	DevEUI       string `json:"terminalId"`
	LeftElectric uint64 `json:"10000052"`
}
type HeimanPIRSensorPropertyTamper struct {
	DevEUI string `json:"terminalId"`
	Tamper string `json:"10000053"`
}
type HeimanHTSensorPropertyTemperature struct {
	DevEUI      string  `json:"terminalId"`
	Temperature float64 `json:"10000061"`
}
type HeimanHTSensorPropertyHumidity struct {
	DevEUI   string  `json:"terminalId"`
	Humidity float64 `json:"10000062"`
}
type HeimanHTSensorPropertyLeftElectric struct {
	DevEUI       string `json:"terminalId"`
	LeftElectric uint64 `json:"10000063"`
}
type HeimanCOSensorPropertyAlarm struct {
	DevEUI string `json:"terminalId"`
	Alarm  string `json:"10000071"`
}
type HeimanCOSensorPropertyLeftElectric struct {
	DevEUI       string `json:"terminalId"`
	LeftElectric uint64 `json:"10000072"`
}
type HeimanWarningDevicePropertyLeftElectric struct {
	DevEUI       string `json:"terminalId"`
	LeftElectric uint64 `json:"10000082"`
}
type HeimanGASSensorPropertyAlarm struct {
	DevEUI string `json:"terminalId"`
	Alarm  string `json:"10000091"`
}
type HeimanHS2AQPropertyTemperature struct {
	DevEUI      string  `json:"terminalId"`
	Temperature float64 `json:"100000A1"`
}
type HeimanHS2AQPropertyHumidity struct {
	DevEUI   string  `json:"terminalId"`
	Humidity float64 `json:"100000A2"`
}
type HeimanHS2AQPropertyFormaldehyde struct {
	DevEUI       string  `json:"terminalId"`
	Formaldehyde float64 `json:"100000A3"`
}
type HeimanHS2AQPropertyPM25 struct {
	DevEUI string `json:"terminalId"`
	PM25   uint64 `json:"100000A4"`
}
type HeimanHS2AQPropertyPM10 struct {
	DevEUI string `json:"terminalId"`
	PM10   string `json:"100000A5"`
}
type HeimanHS2AQPropertyLeftElectric struct {
	DevEUI       string `json:"terminalId"`
	LeftElectric uint64 `json:"100000A6"`
}

// type HeimanHS2AQPropertyTemperatureAlarmThreshold struct {
// 	DevEUI                    string `json:"terminalId"`
// 	TemperatureAlarmThreshold string `json:"100000A7"`
// }
// type HeimanHS2AQPropertyHumidityAlarmThreshold struct {
// 	DevEUI                 string `json:"terminalId"`
// 	HumidityAlarmThreshold string `json:"100000A8"`
// }
type HeimanHS2AQPropertyTemperatureAlarm struct {
	DevEUI           string `json:"terminalId"`
	TemperatureAlarm string `json:"100000A9"`
}
type HeimanHS2AQPropertyHumidityAlarm struct {
	DevEUI        string `json:"terminalId"`
	HumidityAlarm string `json:"100000AA"`
}

// type HeimanHS2AQPropertyTVOC struct {
// 	DevEUI string `json:"terminalId"`
// 	TVOC   string `json:"100000AB"`
// }
type HeimanHS2AQPropertyAQI struct {
	DevEUI string `json:"terminalId"`
	AQI    string `json:"100000AC"`
}
type HeimanHS2AQPropertyBatteryState struct {
	DevEUI       string `json:"terminalId"`
	BatteryState string `json:"100000AD"`
}
type HeimanHS2AQProperty struct {
	DevEUI       string `json:"terminalId"`
	PM10         string `json:"100000A5,omitempty"`
	TVOC         string `json:"100000AB,omitempty"`
	AQI          string `json:"100000AC,omitempty"`
	BatteryState string `json:"100000AD,omitempty"`
}

// type HeimanHS2AQPropertyUnDisturb struct {
// 	DevEUI    string `json:"terminalId"`
// 	UnDisturb string `json:"100000AE"`
// }
type HeimanHS2AQPropertyFormaldehydeAlarm struct {
	DevEUI            string `json:"terminalId"`
	FormaldehydeAlarm string `json:"100000AF"`
}
type HeimanHS2AQPropertyPM25Alarm struct {
	DevEUI    string `json:"terminalId"`
	PM25Alarm string `json:"100000B1"`
}
type HeimanHS2AQPropertyPM10Alarm struct {
	DevEUI    string `json:"terminalId"`
	PM10Alarm string `json:"100000B2"`
}
type HeimanHS2AQPropertyTVOCAlarm struct {
	DevEUI    string `json:"terminalId"`
	TVOCAlarm string `json:"100000B3"`
}
type HeimanHS2AQPropertyAQIAlarm struct {
	DevEUI   string `json:"terminalId"`
	AQIAlarm string `json:"100000B4"`
}
type HeimanHS2AQPropertyBatteryAlarm struct {
	DevEUI       string `json:"terminalId"`
	BatteryAlarm string `json:"100000B5"`
}
type HeimanHS2SW1LEFR30PropertyOnOff1 struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"100000C1"`
}
type HeimanHS2SW1LEFR30PropertyOnOff1DefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"100000C1"`
}
type HeimanHS2SW2LEFR30PropertyOnOff1 struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"100000D1"`
}
type HeimanHS2SW2LEFR30PropertyOnOff1DefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"100000D1"`
}
type HeimanHS2SW2LEFR30PropertyOnOff2 struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"100000D2"`
}
type HeimanHS2SW2LEFR30PropertyOnOff2DefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"100000D2"`
}
type HeimanHS2SW3LEFR30PropertyOnOff1 struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"100000E1"`
}
type HeimanHS2SW3LEFR30PropertyOnOff1DefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"100000E1"`
}
type HeimanHS2SW3LEFR30PropertyOnOff2 struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"100000E2"`
}
type HeimanHS2SW3LEFR30PropertyOnOff2DefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"100000E2"`
}
type HeimanHS2SW3LEFR30PropertyOnOff3 struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"100000E3"`
}
type HeimanHS2SW3LEFR30PropertyOnOff3DefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"100000E3"`
}
type HeimanSceneSwitchEM30PropertyAtHome struct {
	DevEUI string `json:"terminalId"`
	AtHome string `json:"100000F2"` //0xf0, 在家模式
}
type HeimanSceneSwitchEM30PropertyGoOut struct {
	DevEUI string `json:"terminalId"`
	GoOut  string `json:"100000F1"` //0xf1, 外出模式
}
type HeimanSceneSwitchEM30PropertyCinema struct {
	DevEUI string `json:"terminalId"`
	Cinema string `json:"100000F3"` //0xf2, 影院模式
}
type HeimanSceneSwitchEM30PropertyRepast struct {
	DevEUI string `json:"terminalId"`
	Repast string `json:"100000F6"` //0xf3, 就餐模式
}
type HeimanSceneSwitchEM30PropertySleep struct {
	DevEUI string `json:"terminalId"`
	Sleep  string `json:"100000F4"` //0xf4, 就寝模式
}
type HeimanSceneSwitchEM30PropertyLeftElectric struct {
	DevEUI       string `json:"terminalId"`
	LeftElectric uint64 `json:"100000F5"`
}

// type HeimanIRControlEMPropertySendKeyCommand struct {
// 	DevEUI         string `json:"terminalId"`
// 	SendKeyCommand uint64 `json:"10000111"` //0xf0, 发送遥控指令
// }
// type HeimanIRControlEMPropertyStudyKey struct {
// 	DevEUI   string `json:"terminalId"`
// 	StudyKey uint64 `json:"10000112"` //0xf1, 学习遥控指令
// }
type HeimanIRControlEMPropertyStudyKeyRsp struct {
	Code     int                                  `json:"code"`
	Message  string                               `json:"message"`
	DevEUI   string                               `json:"terminalId"`
	StudyKey cluster.HEIMANInfraredRemoteStudyKey `json:"10000112"`
}

// type HeimanIRControlEMPropertyDeleteKey struct {
// 	DevEUI    string `json:"terminalId"`
// 	DeleteKey uint64 `json:"10000113"` //0xf3, 删除遥控指令
// }
// type HeimanIRControlEMPropertyCreateID struct {
// 	DevEUI   string `json:"terminalId"`
// 	CreateID uint64 `json:"10000114"` //0xf4, 创建遥控设备
// }
type HeimanIRControlEMPropertyCreateIDRsp struct {
	Code     int                                          `json:"code"`
	Message  string                                       `json:"message"`
	DevEUI   string                                       `json:"terminalId"`
	CreateID cluster.HEIMANInfraredRemoteCreateIDResponse `json:"10000114"`
}

// type HeimanIRControlEMPropertyGetIDAndKeyCodeList struct {
// 	DevEUI              string `json:"terminalId"`
// 	GetIDAndKeyCodeList uint64 `json:"10000115"` //0xf6, 获取遥控设备及其指令列表
// }
type HeimanIRControlEMPropertyIdKeyCode struct {
	ID        uint8   `json:"ID"`
	ModelType uint8   `json:"modelType"`
	KeyNum    uint8   `json:"keyNum"`
	KeyCode   []uint8 `json:"keyCode"`
}
type HeimanIRControlEMPropertyGetIDAndKeyCodeListResponse struct {
	PacketNumSum        uint8                                `json:"packetNumSum"`
	CurrentPacketNum    uint8                                `json:"currentPacketNum"`
	CurrentPacketLength uint8                                `json:"currentPacketLength"`
	IDKeyCodeList       []HeimanIRControlEMPropertyIdKeyCode `json:"IDKeyCodeList"`
}
type HeimanIRControlEMPropertyGetIDAndKeyCodeListRsp struct {
	Code                int                                                  `json:"code"`
	Message             string                                               `json:"message"`
	DevEUI              string                                               `json:"terminalId"`
	GetIDAndKeyCodeList HeimanIRControlEMPropertyGetIDAndKeyCodeListResponse `json:"10000115"`
}
type HonyarSingleSwitch00500c32PropertyOnOff1 struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"20000011"`
}
type HonyarSingleSwitch00500c32PropertyOnOff1DefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"20000011"`
}
type HonyarDoubleSwitch00500c33PropertyOnOff1 struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"20000021"`
}
type HonyarDoubleSwitch00500c33PropertyOnOff1DefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"20000021"`
}
type HonyarDoubleSwitch00500c33PropertyOnOff2 struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"20000022"`
}
type HonyarDoubleSwitch00500c33PropertyOnOff2DefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"20000022"`
}
type HonyarTripleSwitch00500c35PropertyOnOff1 struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"20000031"`
}
type HonyarTripleSwitch00500c35PropertyOnOff1DefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"20000031"`
}
type HonyarTripleSwitch00500c35PropertyOnOff2 struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"20000032"`
}
type HonyarTripleSwitch00500c35PropertyOnOff2DefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"20000032"`
}
type HonyarTripleSwitch00500c35PropertyOnOff3 struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"20000033"`
}
type HonyarTripleSwitch00500c35PropertyOnOff3DefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"20000033"`
}
type HonyarSocket000a0c3cPropertyOnOff struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"20000041"`
}
type HonyarSocket000a0c3cPropertyOnOffDefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"20000041"`
}
type HonyarSocket000a0c3cPropertyCurretSummationDelivered struct {
	DevEUI                   string `json:"terminalId"`
	CurretSummationDelivered string `json:"20000042"`
}
type HonyarSocket000a0c3cPropertyVoltageCurrentPower struct {
	DevEUI      string `json:"terminalId"`
	RMSVoltage  string `json:"20000043,omitempty"`
	RMSCurrent  string `json:"20000044,omitempty"`
	ActivePower string `json:"20000045,omitempty"`
}
type HonyarSocket000a0c3cPropertyChildLock struct {
	DevEUI    string `json:"terminalId"`
	ChildLock string `json:"20000046"`
}
type HonyarSocket000a0c3cPropertyChildLockWriteRsp struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	DevEUI    string `json:"terminalId"`
	ChildLock string `json:"20000046"`
}
type HonyarSocket000a0c3cPropertyUSB struct {
	DevEUI string `json:"terminalId"`
	USB    string `json:"20000047"`
}
type HonyarSocket000a0c3cPropertyUSBDefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	USB     string `json:"20000047"`
}
type Honyar6SceneSwitch005f0c3bPropertyLeftUp struct {
	DevEUI string `json:"terminalId"`
	LeftUp string `json:"20000061"` //0x01, 左上键
}
type Honyar6SceneSwitch005f0c3bPropertyLeftMiddle struct {
	DevEUI     string `json:"terminalId"`
	LeftMiddle string `json:"20000062"` //0x02, 左中键
}
type Honyar6SceneSwitch005f0c3bPropertyLeftDown struct {
	DevEUI   string `json:"terminalId"`
	LeftDown string `json:"20000063"` //0x03, 左下键
}
type Honyar6SceneSwitch005f0c3bPropertyRightUp struct {
	DevEUI  string `json:"terminalId"`
	RightUp string `json:"20000064"` //0x04, 右上键
}
type Honyar6SceneSwitch005f0c3bPropertyRightMiddle struct {
	DevEUI      string `json:"terminalId"`
	RightMiddle string `json:"20000065"` //0x05, 右中键
}
type Honyar6SceneSwitch005f0c3bPropertyRightDown struct {
	DevEUI    string `json:"terminalId"`
	RightDown string `json:"20000066"` //0x06, 右下键
}

// type Honyar6SceneSwitch005f0c3bPropertyLeftUpAddScene struct {
// 	DevEUI   string `json:"terminalId"`
// 	AddScene string `json:"20000067"` //0x01~0x21，对应1~33个场景
// }
type Honyar6SceneSwitch005f0c3bPropertyLeftUpAddSceneRsp struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	DevEUI   string `json:"terminalId"`
	AddScene string `json:"20000067"` //0x01~0x21，对应1~33个场景
}

// type Honyar6SceneSwitch005f0c3bPropertyLeftMiddleAddScene struct {
// 	DevEUI   string `json:"terminalId"`
// 	AddScene string `json:"20000068"` //0x01~0x21，对应1~33个场景
// }
type Honyar6SceneSwitch005f0c3bPropertyLeftMiddleAddSceneRsp struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	DevEUI   string `json:"terminalId"`
	AddScene string `json:"20000068"` //0x01~0x21，对应1~33个场景
}

// type Honyar6SceneSwitch005f0c3bPropertyLeftDownAddScene struct {
// 	DevEUI   string `json:"terminalId"`
// 	AddScene string `json:"20000069"` //0x01~0x21，对应1~33个场景
// }
type Honyar6SceneSwitch005f0c3bPropertyLeftDownAddSceneRsp struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	DevEUI   string `json:"terminalId"`
	AddScene string `json:"20000069"` //0x01~0x21，对应1~33个场景
}

// type Honyar6SceneSwitch005f0c3bPropertyRightUpAddScene struct {
// 	DevEUI   string `json:"terminalId"`
// 	AddScene string `json:"2000006A"` //0x01~0x21，对应1~33个场景
// }
type Honyar6SceneSwitch005f0c3bPropertyRightUpAddSceneRsp struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	DevEUI   string `json:"terminalId"`
	AddScene string `json:"2000006A"` //0x01~0x21，对应1~33个场景
}

// type Honyar6SceneSwitch005f0c3bPropertyRightMiddleAddScene struct {
// 	DevEUI   string `json:"terminalId"`
// 	AddScene string `json:"2000006B"` //0x01~0x21，对应1~33个场景
// }
type Honyar6SceneSwitch005f0c3bPropertyRightMiddleAddSceneRsp struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	DevEUI   string `json:"terminalId"`
	AddScene string `json:"2000006B"` //0x01~0x21，对应1~33个场景
}

// type Honyar6SceneSwitch005f0c3bPropertyRightDownAddScene struct {
// 	DevEUI   string `json:"terminalId"`
// 	AddScene string `json:"2000006C"` //0x01~0x21，对应1~33个场景
// }
type Honyar6SceneSwitch005f0c3bPropertyRightDownAddSceneRsp struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	DevEUI   string `json:"terminalId"`
	AddScene string `json:"2000006C"` //0x01~0x21，对应1~33个场景
}
type HonyarSocketHY0105PropertyOnOff struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"20000091"`
}
type HonyarSocketHY0105PropertyOnOffDefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"20000091"`
}
type HonyarSocketHY0105PropertyCurretSummationDelivered struct {
	DevEUI                   string `json:"terminalId"`
	CurretSummationDelivered string `json:"20000092"`
}
type HonyarSocketHY0105PropertyVoltageCurrentPower struct {
	DevEUI      string `json:"terminalId"`
	RMSVoltage  string `json:"20000093,omitempty"`
	RMSCurrent  string `json:"20000094,omitempty"`
	ActivePower string `json:"20000095,omitempty"`
}
type HonyarSocketHY0105PropertyChildLock struct {
	DevEUI    string `json:"terminalId"`
	ChildLock string `json:"20000096"`
}
type HonyarSocketHY0105PropertyChildLockWriteRsp struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	DevEUI    string `json:"terminalId"`
	ChildLock string `json:"20000096"`
}
type HonyarSocketHY0105PropertyUSB struct {
	DevEUI string `json:"terminalId"`
	USB    string `json:"20000097"`
}
type HonyarSocketHY0105PropertyUSBDefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	USB     string `json:"20000097"`
}
type HonyarSocketHY0105PropertyBackGroundLight struct {
	DevEUI          string `json:"terminalId"`
	BackGroundLight string `json:"20000098"`
}
type HonyarSocketHY0105PropertyBackGroundLightWriteRsp struct {
	Code            int    `json:"code"`
	Message         string `json:"message"`
	DevEUI          string `json:"terminalId"`
	BackGroundLight string `json:"20000098"`
}
type HonyarSocketHY0105PropertyPowerOffMemory struct {
	DevEUI         string `json:"terminalId"`
	PowerOffMemory string `json:"20000099"`
}
type HonyarSocketHY0105PropertyPowerOffMemoryWriteRsp struct {
	Code           int    `json:"code"`
	Message        string `json:"message"`
	DevEUI         string `json:"terminalId"`
	PowerOffMemory string `json:"20000099"`
}
type HonyarSocketHY0105PropertyHistoryElectricClearWriteRsp struct {
	Code                 int    `json:"code"`
	Message              string `json:"message"`
	DevEUI               string `json:"terminalId"`
	HistoryElectricClear string `json:"2000009A"`
}
type HonyarSocketHY0106PropertyOnOff struct {
	DevEUI string `json:"terminalId"`
	OnOff  string `json:"200000A1"`
}
type HonyarSocketHY0106PropertyOnOffDefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	DevEUI  string `json:"terminalId"`
	OnOff   string `json:"200000A1"`
}
type HonyarSocketHY0106PropertyCurretSummationDelivered struct {
	DevEUI                   string `json:"terminalId"`
	CurretSummationDelivered string `json:"200000A2"`
}
type HonyarSocketHY0106PropertyVoltageCurrentPower struct {
	DevEUI      string `json:"terminalId"`
	RMSVoltage  string `json:"200000A3,omitempty"`
	RMSCurrent  string `json:"200000A4,omitempty"`
	ActivePower string `json:"200000A5,omitempty"`
}
type HonyarSocketHY0106PropertyChildLock struct {
	DevEUI    string `json:"terminalId"`
	ChildLock string `json:"200000A6"`
}
type HonyarSocketHY0106PropertyChildLockWriteRsp struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	DevEUI    string `json:"terminalId"`
	ChildLock string `json:"200000A6"`
}
type HonyarSocketHY0106PropertyBackGroundLight struct {
	DevEUI          string `json:"terminalId"`
	BackGroundLight string `json:"200000A7"`
}
type HonyarSocketHY0106PropertyBackGroundLightWriteRsp struct {
	Code            int    `json:"code"`
	Message         string `json:"message"`
	DevEUI          string `json:"terminalId"`
	BackGroundLight string `json:"200000A7"`
}
type HonyarSocketHY0106PropertyPowerOffMemory struct {
	DevEUI         string `json:"terminalId"`
	PowerOffMemory string `json:"200000A8"`
}
type HonyarSocketHY0106PropertyPowerOffMemoryWriteRsp struct {
	Code           int    `json:"code"`
	Message        string `json:"message"`
	DevEUI         string `json:"terminalId"`
	PowerOffMemory string `json:"200000A8"`
}
type HonyarSocketHY0106PropertyHistoryElectricClearWriteRsp struct {
	Code                 int    `json:"code"`
	Message              string `json:"message"`
	DevEUI               string `json:"terminalId"`
	HistoryElectricClear string `json:"200000A9"`
}

// 定义Iotware属性method
const (
	IotwarePropertyMethodPowerSwitch          = "PowerSwitch"
	IotwarePropertyMethodChildLockSwitch      = "ChildLockSwitch"
	IotwarePropertyMethodUSBSwitch_1          = "USBSwitch_1"
	IotwarePropertyMethodBackLightSwitch      = "BackLightSwitch"
	IotwarePropertyMethodMemorySwitch         = "MemorySwitch"
	IotwarePropertyMethodClearHistory         = "ClearHistory"
	IotwarePropertyMethodPowerSwitch_1        = "PowerSwitch_1"
	IotwarePropertyMethodPowerSwitch_2        = "PowerSwitch_2"
	IotwarePropertyMethodPowerSwitch_3        = "PowerSwitch_3"
	IotwarePropertyMethodSwitchLeftUpLogo     = "SwitchLeftUpLogo"
	IotwarePropertyMethodSwitchLeftMidLogo    = "SwitchLeftMidLogo"
	IotwarePropertyMethodSwitchLeftDownLogo   = "SwitchLeftDownLogo"
	IotwarePropertyMethodSwitchRightUpLogo    = "SwitchRightUpLogo"
	IotwarePropertyMethodSwitchRightMidLogo   = "SwitchRightMidLogo"
	IotwarePropertyMethodSwitchRightDownLogo  = "SwitchRightDownLogo"
	IotwarePropertyMethodActiveAlarm          = "ActiveAlarm"
	IotwarePropertyMethodClearAlarm           = "ClearAlarm"
	IotwarePropertyMethodArmed                = "Armed"
	IotwarePropertyMethodDisarmed             = "Disarmed"
	IotwarePropertyMethodTemperatureThreshold = "TemperatureThreshold"
	IotwarePropertyMethodHumidityThreshold    = "HumidityThreshold"
	IotwarePropertyMethodNoDisturb            = "NoDisturb"
	IotwarePropertyMethodDownControl          = "DownControl"
	IotwarePropertyMethodLearnControl         = "LearnControl"
	IotwarePropertyMethodDeleteControl        = "DeleteControl"
	IotwarePropertyMethodCreateModelType      = "CreateModelType"
	IotwarePropertyMethodGetDevAndCmdList     = "GetDevAndCmdList"
	IotwarePropertyMethodLanguage             = "Language"
	IotwarePropertyMethodUnitOfTemperature    = "UnitOfTemperature"
)

// 定义Iotware终端属性结构体
type IotwarePropertyKeepAliveTime struct {
	KeepAliveTime uint16 `json:"KeepAliveTime"`
}

// type IotwarePropertyTime struct {
// 	Time float64 `json:"time"`
// }
// type IotwarePropertyLowerUpper struct {
// 	Lower float64 `json:"lower"`
// 	Upper float64 `json:"upper"`
// }
// type IotwarePropertyNoDisturb struct {
// 	NoDisturb bool `json:"NoDisturb"`
// }
// type IotwarePropertyLanguage struct {
// 	Language string `json:"Language"`
// }
// type IotwarePropertyUnitOfTemperature struct {
// 	UnitOfTemperature string `json:"UnitOfTemperature"`
// }
type IotwarePropertyPowerSwitch struct {
	PowerSwitch bool `json:"PowerSwitch"`
}

// type IotwarePropertyTotalConsumption struct {
// 	TotalConsumption string `json:"TotalConsumption"`
// }
type IotwarePropertyTotalConsumptionRealTimePower struct {
	TotalConsumption string `json:"TotalConsumption,omitempty"`
	RealTimePower    string `json:"RealTimePower,omitempty"`
}

// type IotwarePropertyCurrentVoltage struct {
// 	CurrentVoltage string `json:"CurrentVoltage"`
// }
// type IotwarePropertyCurrent struct {
// 	Current string `json:"Current"`
// }
// type IotwarePropertyRealTimePower struct {
// 	RealTimePower string `json:"RealTimePower"`
// }
type IotwarePropertyVoltageCurrentPower struct {
	CurrentVoltage string `json:"CurrentVoltage,omitempty"`
	Current        string `json:"Current,omitempty"`
	RealTimePower  string `json:"RealTimePower,omitempty"`
}
type IotwarePropertyChildLockSwitch struct {
	ChildLockSwitch bool `json:"ChildLockSwitch"`
}
type IotwarePropertyUSBSwitch_1 struct {
	USBSwitch_1 bool `json:"USBSwitch_1"`
}
type IotwarePropertyBackLightSwitch struct {
	BackLightSwitch bool `json:"BackLightSwitch"`
}
type IotwarePropertyMemorySwitch struct {
	MemorySwitch bool `json:"MemorySwitch"`
}

// type IotwarePropertyClearHistory struct {
// 	ClearHistory string `json:"ClearHistory"`
// }
type IotwarePropertyPowerSwitch_1 struct {
	PowerSwitch_1 bool `json:"PowerSwitch_1"`
}
type IotwarePropertyPowerSwitch_2 struct {
	PowerSwitch_2 bool `json:"PowerSwitch_2"`
}
type IotwarePropertyPowerSwitch_3 struct {
	PowerSwitch_3 bool `json:"PowerSwitch_3"`
}
type IotwarePropertySwitchLeftUp struct {
	SwitchLeftUp string `json:"SwitchLeftUp"`
}
type IotwarePropertySwitchLeftMid struct {
	SwitchLeftMid string `json:"SwitchLeftMid"`
}
type IotwarePropertySwitchLeftDown struct {
	SwitchLeftDown string `json:"SwitchLeftDown"`
}
type IotwarePropertySwitchRightUp struct {
	SwitchRightUp string `json:"SwitchRightUp"`
}
type IotwarePropertySwitchRightMid struct {
	SwitchRightMid string `json:"SwitchRightMid"`
}
type IotwarePropertySwitchRightDown struct {
	SwitchRightDown string `json:"SwitchRightDown"`
}
type IotwarePropertyFireAlarm struct {
	FireAlarm string `json:"FireAlarm"`
}
type IotwarePropertyBatteryPercentage struct {
	BatteryPercentage uint64 `json:"BatteryPercentage"`
}
type IotwarePropertyWaterSensorAlarm struct {
	WaterSensorAlarm string `json:"WaterSensorAlarm"`
}
type IotwarePropertyDismantleAlarm struct {
	DismantleAlarm string `json:"DismantleAlarm"`
}
type IotwarePropertyMoveAlarm struct {
	MoveAlarm string `json:"MoveAlarm"`
}
type IotwarePropertyCurrentTemperature struct {
	CurrentTemperature float64 `json:"CurrentTemperature"`
}
type IotwarePropertyCurrentHumidity struct {
	CurrentHumidity float64 `json:"CurrentHumidity"`
}
type IotwarePropertyCOAlarm struct {
	COAlarm string `json:"COAlarm"`
}
type IotwarePropertyGasAlarm struct {
	GasAlarm string `json:"GasAlarm"`
}
type IotwarePropertyHCHO struct {
	HCHO float64 `json:"HCHO"`
}
type IotwarePropertyPM25 struct {
	PM25 uint64 `json:"PM25"`
}
type IotwarePropertyPM10 struct {
	PM10 string `json:"PM10"`
}
type IotwarePropertyTemperatureAlarm struct {
	TemperatureAlarm string `json:"TemperatureAlarm"`
}
type IotwarePropertyHumidityAlarm struct {
	HumidityAlarm string `json:"HumidityAlarm"`
}
type IotwarePropertyAQI struct {
	AQI string `json:"AQI"`
}
type IotwarePropertyPowerState struct {
	PowerState string `json:"PowerState"`
}
type IotwarePropertyPM10PowerStateAQI struct {
	PM10       string `json:"PM10,omitempty"`
	PowerState string `json:"PowerState,omitempty"`
	AQI        string `json:"AQI,omitempty"`
}
type IotwarePropertyHCHOAlarm struct {
	HCHOAlarm string `json:"HCHOAlarm"`
}
type IotwarePropertyPM25Alarm struct {
	PM25Alarm string `json:"PM25Alarm"`
}
type IotwarePropertyPM10Alarm struct {
	PM10Alarm string `json:"PM10Alarm"`
}
type IotwarePropertyPowerAlarm struct {
	PowerAlarm string `json:"PowerAlarm"`
}
type IotwarePropertyAQIAlarm struct {
	AQIAlarm string `json:"AQIAlarm"`
}
type IotwarePropertySwitchLeft struct {
	SwitchLeft string `json:"SwitchLeft"`
}
type IotwarePropertySwitchRight struct {
	SwitchRight string `json:"SwitchRight"`
}
type IotwarePropertySwitchUp struct {
	SwitchUp string `json:"SwitchUp"`
}
type IotwarePropertySwitchDown struct {
	SwitchDown string `json:"SwitchDown"`
}

// type IotwarePropertySwitchLeftUpLogo struct {
// 	SwitchLeftUpLogo string `json:"SwitchLeftUpLogo"`
// }
// type IotwarePropertySwitchLeftMidLogo struct {
// 	SwitchLeftMidLogo string `json:"SwitchLeftMidLogo"`
// }
// type IotwarePropertySwitchLeftDownLogo struct {
// 	SwitchLeftDownLogo string `json:"SwitchLeftDownLogo"`
// }
// type IotwarePropertySwitchRightUpLogo struct {
// 	SwitchRightUpLogo string `json:"SwitchRightUpLogo"`
// }
// type IotwarePropertySwitchRightMidLogo struct {
// 	SwitchRightMidLogo string `json:"SwitchRightMidLogo"`
// }
// type IotwarePropertySwitchRightDownLogo struct {
// 	SwitchRightDownLogo string `json:"SwitchRightDownLogo"`
// }
// type IotwarePropertyIRId struct {
// 	IRId uint16 `json:"IRId"`
// }
// type IotwarePropertyKeyCode struct {
// 	KeyCode uint16 `json:"KeyCode"`
// }
// type IotwarePropertyModelType struct {
// 	ModelType uint16 `json:"ModelType"`
// }
// type IotwarePropertyDownControl struct {
// 	IRId    float64 `json:"IRId"`
// 	KeyCode float64 `json:"KeyCode"`
// }
// type IotwarePropertyLearnControl struct {
// 	IRId    float64 `json:"IRId"`
// 	KeyCode float64 `json:"KeyCode"`
// }
// type IotwarePropertyDeleteControl struct {
// 	IRId    float64 `json:"IRId"`
// 	KeyCode float64 `json:"KeyCode"`
// }
// type IotwarePropertyCreateModelType struct {
// 	ModelType float64 `json:"ModelType"`
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
			globallogger.Log.Errorf("[Subscribe]: mqtt msg subscribe: topic: %s, networkInAckMsg %+v", topic, networkInAckMsg)
			procActionIotware(networkInAckMsg)
		case TopicV1GatewayEventDeviceDelete:
			var deviceDeleteMsg publicstruct.EventDeviceDeleteIotware
			err := json.Unmarshal(jsonMsg, &deviceDeleteMsg)
			if err != nil {
				globallogger.Log.Errorln("[Subscribe]: JSON marshaling failed:", err)
				return
			}
			globallogger.Log.Errorf("[Subscribe]: mqtt msg subscribe: topic: %s, deviceDeleteMsg %+v", topic, deviceDeleteMsg)
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
