package constant

import "github.com/h3c/iotzigbeeserver-go/config"

// Constant Constant
var Constant Constants

//TIMER TIMER
var TIMER = 1

//Constants Constants
type Constants struct {
	KAFKA             Kafka
	REDIS             Redis
	TIMER             Timer
	TMNTYPE           TmnType
	MANUFACTURERNAME  ManufacturerName
	UDPVERSION        UDPVersion
	MultipleInstances bool
	UsePostgres       bool
	UseRabbitMQ       bool
	Iotware           bool
	Iotprivate        bool
	Iotedge           bool
	HTTPPort          string
}

// Kafka Kafka
type Kafka struct {
	ZigbeeKafkaGroupName string

	ZigbeeKafkaConsumeTopicBuildInPropertyByCustom string
	ZigbeeKafkaConsumeTopicBuildInProperty         string
	ZigbeeKafkaConsumeTopicDelTerminal             string
	ZigbeeKafkaConsumeTopicUpdateTerminal          string

	ZigbeeKafkaProduceTopicConfirmMessage        string
	ZigbeeKafkaProduceTopicCollihighSensorDataUp string
	ZigbeeKafkaProduceTopicDataReportMsg         string
}

//Redis Redis
type Redis struct {
	ZigbeeRedisRemoveRedis           string
	ZigbeeRedisDownMsgSets           string
	ZigbeeRedisDuplicationFlagSets   string
	ZigbeeRedisPermitJoinContentSets string
	ZigbeeRedisCtrlMsgKey            string
	ZigbeeRedisDataMsgKey            string
	ZigbeeRedisDuplicationFlagKey    string
	ZigbeeRedisSNZHA                 string
	ZigbeeRedisSNTransfer            string
	ZigbeeRedisCollihighInterval     string
	ZigbeeRedisPermitJoinContentKey  string
}

//Timer Timer
type Timer struct {
	ZigbeeTimerRetran            int
	ZigbeeTimerDuplicationFlag   int
	ZigbeeTimerAPBusy            int
	ZigbeeTimerAfterFive         int
	ZigbeeTimerAfterTwo          int
	ZigbeeTimerServerRebuild     int
	ZigbeeTimerSetMsgRemove      int
	ZigbeeTimerMsgDelete         int
	ZigbeeTimerOneDay            int
	ZigbeeTimerPermitJoinContent int

	HeimanKeepAliveTimer  HeimanKeepAliveTimer
	HonyarKeepAliveTimer  HonyarKeepAliveTimer
	MailekeKeepAliveTimer MailekeKeepAliveTimer
}

// HeimanKeepAliveTimer HeimanKeepAliveTimer
type HeimanKeepAliveTimer struct {
	ZigbeeTerminalHS2AQEM          uint16
	ZigbeeTerminalSmartPlug        uint16
	ZigbeeTerminalESocket          uint16
	ZigbeeTerminalHTEM             uint16
	ZigbeeTerminalSmokeSensorEM    uint16
	ZigbeeTerminalWaterSensorEM    uint16
	ZigbeeTerminalPIRSensorEM      uint16
	ZigbeeTerminalPIRILLSensorEF30 uint16
	ZigbeeTerminalCOSensorEM       uint16
	ZigbeeTerminalWarningDevice    uint16
	ZigbeeTerminalGASSensorEM      uint16
	ZigbeeTerminalHS2SW1LEFR30     uint16
	ZigbeeTerminalHS2SW2LEFR30     uint16
	ZigbeeTerminalHS2SW3LEFR30     uint16
	ZigbeeTerminalSceneSwitchEM30  uint16
	ZigbeeTerminalDoorSensorEF30   uint16
	ZigbeeTerminalIRControlEM      uint16
}

// HonyarKeepAliveTimer HonyarKeepAliveTimer
type HonyarKeepAliveTimer struct {
	ZigbeeTerminalSingleSwitch00500c32 uint16
	ZigbeeTerminalDoubleSwitch00500c33 uint16
	ZigbeeTerminalTripleSwitch00500c35 uint16
	ZigbeeTerminalSingleSwitchHY0141   uint16
	ZigbeeTerminalDoubleSwitchHY0142   uint16
	ZigbeeTerminalTripleSwitchHY0143   uint16
	ZigbeeTerminalSocket000a0c3c       uint16
	ZigbeeTerminalSocket000a0c55       uint16
	ZigbeeTerminalSocketHY0105         uint16
	ZigbeeTerminalSocketHY0106         uint16
	ZigbeeTerminal1SceneSwitch005f0cf1 uint16
	ZigbeeTerminal2SceneSwitch005f0cf3 uint16
	ZigbeeTerminal3SceneSwitch005f0cf2 uint16
	ZigbeeTerminal6SceneSwitch005f0c3b uint16
}

// MailekeKeepAliveTimer MailekeKeepAliveTimer
type MailekeKeepAliveTimer struct {
	ZigbeeTerminalPMT1004Detector uint16
}

//TmnType TmnType
type TmnType struct {
	COLLIHIGH      Collihigh
	SWITCH         Switch
	LIGHTINGSWITCH LightingSwitch
	SCENE          Scene
	CURTAIN        Curtain
	SOCKET         Socket
	SENSOR         Sensor
	TRANSPONDER    Transponder
	MAILEKE        Maileke
	HEIMAN         Heiman
	HONYAR         Honyar
}

//Collihigh Collihigh
type Collihigh struct {
	ZigbeeTerminalCO2                 string
	ZigbeeTerminalLightingPM25        string
	ZigbeeTerminalHumitureLightingCO2 string
}

//Switch Switch
type Switch struct {
	ZigbeeTerminalSingleSwitch    string
	ZigbeeTerminalDoubleSwitch    string
	ZigbeeTerminalTripleSwitch    string
	ZigbeeTerminalQuadrupleSwitch string
}

//LightingSwitch LightingSwitch
type LightingSwitch struct {
	ZigbeeTerminalSingleLightingSwitch    string
	ZigbeeTerminalDoubleLightingSwitch    string
	ZigbeeTerminalTripleLightingSwitch    string
	ZigbeeTerminalQuadrupleLightingSwitch string
}

//Scene Scene
type Scene struct {
	ZigbeeTerminalScene string
}

//Curtain Curtain
type Curtain struct {
	ZigbeeTerminalCurtain string
}

//Socket Socket
type Socket struct {
	ZigbeeTerminalSocket string
}

//Sensor Sensor
type Sensor struct {
	ZigbeeTerminalDoorSensor         string
	ZigbeeTerminalHumanDetector      string
	ZigbeeTerminalEmergencyButton    string
	ZigbeeTerminalGasDetector        string
	ZigbeeTerminalWaterDetector      string
	ZigbeeTerminalSmokeFireDetector  string
	ZigbeeTerminalVoiceLightDetector string
	ZigbeeTerminalHumitureDetector   string
}

//Transponder Transponder
type Transponder struct {
	ZigbeeTerminalInfraredTransponder string
}

//Maileke Maileke
type Maileke struct {
	ZigbeeTerminalPMT1004Detector string
}

// Heiman Heiman
type Heiman struct {
	ZigbeeTerminalHS2AQEM          string
	ZigbeeTerminalSmartPlug        string
	ZigbeeTerminalESocket          string
	ZigbeeTerminalHTEM             string
	ZigbeeTerminalSmokeSensorEM    string
	ZigbeeTerminalWaterSensorEM    string
	ZigbeeTerminalPIRSensorEM      string
	ZigbeeTerminalPIRILLSensorEF30 string
	ZigbeeTerminalCOSensorEM       string
	ZigbeeTerminalWarningDevice    string
	ZigbeeTerminalGASSensorEM      string
	ZigbeeTerminalHS2SW1LEFR30     string
	ZigbeeTerminalHS2SW2LEFR30     string
	ZigbeeTerminalHS2SW3LEFR30     string
	ZigbeeTerminalSceneSwitchEM30  string
	ZigbeeTerminalDoorSensorEF30   string
	ZigbeeTerminalIRControlEM      string
}

// Honyar Honyar
type Honyar struct {
	ZigbeeTerminalSingleSwitch00500c32 string //零火
	ZigbeeTerminalDoubleSwitch00500c33 string
	ZigbeeTerminalTripleSwitch00500c35 string
	ZigbeeTerminalSingleSwitchHY0141   string //单火
	ZigbeeTerminalDoubleSwitchHY0142   string
	ZigbeeTerminalTripleSwitchHY0143   string
	ZigbeeTerminalSocket000a0c3c       string //10A U1 双USB
	ZigbeeTerminalSocket000a0c55       string //16A U1
	ZigbeeTerminalSocketHY0105         string //10A U2 单USB
	ZigbeeTerminalSocketHY0106         string //16A U2
	ZigbeeTerminal1SceneSwitch005f0cf1 string
	ZigbeeTerminal2SceneSwitch005f0cf3 string
	ZigbeeTerminal3SceneSwitch005f0cf2 string
	ZigbeeTerminal6SceneSwitch005f0c3b string
}

// ManufacturerName ManufacturerName
type ManufacturerName struct {
	HeiMan string
	Honyar string
}

// UDPVersion UDPVersion
type UDPVersion struct {
	Version0101 string
	Version0102 string
}

func init() {
	var defaultConfigs = make(map[string]interface{})
	err := config.LoadJSON("./config/production.json", &defaultConfigs)
	if err != nil {
		panic(err)
	}
	Constant.MultipleInstances = defaultConfigs["multipleInstances"].(bool)
	Constant.UsePostgres = defaultConfigs["usePostgres"].(bool)
	Constant.UseRabbitMQ = defaultConfigs["useRabbitMQ"].(bool)
	Constant.Iotware = defaultConfigs["iotware"].(bool)
	Constant.Iotprivate = defaultConfigs["iotprivate"].(bool)
	Constant.Iotedge = defaultConfigs["iotedge"].(bool)
	Constant.HTTPPort = defaultConfigs["httpPort"].(string)

	Constant.KAFKA.ZigbeeKafkaGroupName = "iotzigbeeserver_kafka"
	Constant.KAFKA.ZigbeeKafkaConsumeTopicBuildInPropertyByCustom = "custom2iotterminalmgr2buildInProperty2ZigBee"
	Constant.KAFKA.ZigbeeKafkaConsumeTopicBuildInProperty = "iotterminalmgr2buildInProperty2ZigBee"
	Constant.KAFKA.ZigbeeKafkaConsumeTopicDelTerminal = "iotterminalmgr2delTerminal"
	Constant.KAFKA.ZigbeeKafkaConsumeTopicUpdateTerminal = "iotterminalmgr2updateTerminal"
	Constant.KAFKA.ZigbeeKafkaProduceTopicConfirmMessage = "confirmMessage"
	Constant.KAFKA.ZigbeeKafkaProduceTopicCollihighSensorDataUp = "collihighSensorDataUp"
	Constant.KAFKA.ZigbeeKafkaProduceTopicDataReportMsg = "dataReportMsg"

	Constant.REDIS.ZigbeeRedisRemoveRedis = "lmj"
	Constant.REDIS.ZigbeeRedisDownMsgSets = "zigbee_down_msg"
	Constant.REDIS.ZigbeeRedisDuplicationFlagSets = "zigbee_duplication_flag"
	Constant.REDIS.ZigbeeRedisPermitJoinContentSets = "zigbee_permit_join_content"
	Constant.REDIS.ZigbeeRedisCtrlMsgKey = "zigbee_down_ctrl_msg_"
	Constant.REDIS.ZigbeeRedisDataMsgKey = "zigbee_down_data_msg_"
	Constant.REDIS.ZigbeeRedisDuplicationFlagKey = "zigbee_duplication_flag_"
	Constant.REDIS.ZigbeeRedisSNZHA = "zigbee_sn_zha_"
	Constant.REDIS.ZigbeeRedisSNTransfer = "zigbee_sn_transfer_"
	Constant.REDIS.ZigbeeRedisCollihighInterval = "zigbee_collihigh_interval_"
	Constant.REDIS.ZigbeeRedisPermitJoinContentKey = "zigbee_permit_join_content_"

	Constant.TIMER.ZigbeeTimerRetran = TIMER
	Constant.TIMER.ZigbeeTimerDuplicationFlag = (3*TIMER + 10)
	Constant.TIMER.ZigbeeTimerAPBusy = (3*TIMER + 10)
	Constant.TIMER.ZigbeeTimerAfterFive = 5
	Constant.TIMER.ZigbeeTimerAfterTwo = 2
	Constant.TIMER.ZigbeeTimerServerRebuild = 7
	Constant.TIMER.ZigbeeTimerSetMsgRemove = (3*3*TIMER + 5)
	Constant.TIMER.ZigbeeTimerMsgDelete = (3*3*TIMER + 10)
	Constant.TIMER.ZigbeeTimerOneDay = (24 * 60 * 60)
	Constant.TIMER.ZigbeeTimerPermitJoinContent = 250

	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalCOSensorEM = (60 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalESocket = (10 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalGASSensorEM = (10 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHS2AQEM = (60 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHTEM = (60 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalPIRSensorEM = (60 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalPIRILLSensorEF30 = (60 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalSmartPlug = (10 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalSmokeSensorEM = (60 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalWarningDevice = (60 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalWaterSensorEM = (60 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHS2SW1LEFR30 = (10 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHS2SW2LEFR30 = (10 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalHS2SW3LEFR30 = (10 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalSceneSwitchEM30 = (60 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalDoorSensorEF30 = (60 * 60)
	Constant.TIMER.HeimanKeepAliveTimer.ZigbeeTerminalIRControlEM = (60 * 60)

	Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSingleSwitch00500c32 = (10 * 60)
	Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalDoubleSwitch00500c33 = (10 * 60)
	Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalTripleSwitch00500c35 = (10 * 60)
	Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSingleSwitchHY0141 = (10 * 60)
	Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalDoubleSwitchHY0142 = (10 * 60)
	Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalTripleSwitchHY0143 = (10 * 60)
	Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSocket000a0c3c = (10 * 60)
	Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSocket000a0c55 = (10 * 60)
	Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSocketHY0105 = (10 * 60)
	Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminalSocketHY0106 = (10 * 60)
	Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminal1SceneSwitch005f0cf1 = (10 * 60)
	Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminal2SceneSwitch005f0cf3 = (10 * 60)
	Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminal3SceneSwitch005f0cf2 = (10 * 60)
	Constant.TIMER.HonyarKeepAliveTimer.ZigbeeTerminal6SceneSwitch005f0c3b = (10 * 60)

	Constant.TIMER.MailekeKeepAliveTimer.ZigbeeTerminalPMT1004Detector = (5 * 60)

	Constant.TMNTYPE.COLLIHIGH.ZigbeeTerminalCO2 = "ZH-102-SN_CO2"
	Constant.TMNTYPE.COLLIHIGH.ZigbeeTerminalLightingPM25 = "组合传感器_光照_PM2.5"
	Constant.TMNTYPE.COLLIHIGH.ZigbeeTerminalHumitureLightingCO2 = "JZH-021-SN_温湿度_光照_CO2"

	Constant.TMNTYPE.SWITCH.ZigbeeTerminalSingleSwitch = "single_switch"
	Constant.TMNTYPE.SWITCH.ZigbeeTerminalDoubleSwitch = "double_switch"
	Constant.TMNTYPE.SWITCH.ZigbeeTerminalTripleSwitch = "triple_switch"
	Constant.TMNTYPE.SWITCH.ZigbeeTerminalQuadrupleSwitch = "quadruple_switch"

	Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalSingleLightingSwitch = "single_LIGHTINGSWITCH"
	Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalDoubleLightingSwitch = "double_LIGHTINGSWITCH"
	Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalTripleLightingSwitch = "triple_LIGHTINGSWITCH"
	Constant.TMNTYPE.LIGHTINGSWITCH.ZigbeeTerminalQuadrupleLightingSwitch = "quadruple_LIGHTINGSWITCH"

	Constant.TMNTYPE.SCENE.ZigbeeTerminalScene = "scene"

	Constant.TMNTYPE.CURTAIN.ZigbeeTerminalCurtain = "curtain"

	Constant.TMNTYPE.SOCKET.ZigbeeTerminalSocket = "socket"

	Constant.TMNTYPE.SENSOR.ZigbeeTerminalDoorSensor = "door_sensor"
	Constant.TMNTYPE.SENSOR.ZigbeeTerminalHumanDetector = "human_detector"
	Constant.TMNTYPE.SENSOR.ZigbeeTerminalEmergencyButton = "emergency_button"
	Constant.TMNTYPE.SENSOR.ZigbeeTerminalGasDetector = "gas_detector"
	Constant.TMNTYPE.SENSOR.ZigbeeTerminalWaterDetector = "water_detector"
	Constant.TMNTYPE.SENSOR.ZigbeeTerminalSmokeFireDetector = "smoke_fire_detector"
	Constant.TMNTYPE.SENSOR.ZigbeeTerminalVoiceLightDetector = "voice_light_detector"
	Constant.TMNTYPE.SENSOR.ZigbeeTerminalHumitureDetector = "humiture_detector"

	Constant.TMNTYPE.TRANSPONDER.ZigbeeTerminalInfraredTransponder = "infrared_transponder"

	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM = "COSensor-EM"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket = "E_Socket"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM = "GASSensor-EM"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM = "HS2AQ-EM"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM = "HT-EM"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM = "PIRSensor-EM"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30 = "PIRILLSensor-EF-3.0"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug = "SmartPlug"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM = "SmokeSensor-EM"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice = "WarningDevice"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM = "WaterSensor-EM"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30 = "HS2SW1L-EFR-3.0"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30 = "HS2SW2L-EFR-3.0"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30 = "HS2SW3L-EFR-3.0"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSceneSwitchEM30 = "SceneSwitch-EM-3.0"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30 = "DoorSensor-EF-3.0"
	Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM = "IRControl-EM"

	Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32 = "00500c32"
	Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33 = "00500c33"
	Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35 = "00500c35"
	Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141 = "HY0141"
	Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142 = "HY0142"
	Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143 = "HY0143"
	Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c = "000a0c3c"
	Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55 = "000a0c55"
	Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105 = "HY0105"
	Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106 = "HY0106"
	Constant.TMNTYPE.HONYAR.ZigbeeTerminal1SceneSwitch005f0cf1 = "005f0cf1"
	Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3 = "005f0cf3"
	Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2 = "005f0cf2"
	Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b = "005f0c3b"

	Constant.TMNTYPE.MAILEKE.ZigbeeTerminalPMT1004Detector = "PMT1004_温湿度_PM"

	Constant.MANUFACTURERNAME.HeiMan = "HEIMAN"
	Constant.MANUFACTURERNAME.Honyar = "HONYAR"

	Constant.UDPVERSION.Version0101 = "0101"
	Constant.UDPVERSION.Version0102 = "0102"
}
