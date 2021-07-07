package kafka

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
	"github.com/h3c/iotzigbeeserver-go/zcl/zclmsgdown"
	"github.com/lib/pq"
)

// Client 全局kafka客户端变量
var Client sarama.Client
var Produce sarama.AsyncProducer
var Group sarama.ConsumerGroup
var connectting = false
var connectFlag = false

func GetConnectFlag() bool {
	return connectFlag
}

func SetConnectFlag(flag bool) {
	connectFlag = flag
}

func connectLoop(addrs []string, config *sarama.Config) {
	var err error
	for {
		if Client, err = sarama.NewClient(addrs, config); err != nil {
			globallogger.Log.Errorf("[Kafka] [connectLoop] failed err:%s", err.Error())
			time.Sleep(time.Second * 2)
		} else {
			globallogger.Log.Errorf("[Kafka] new client success :%+v", addrs)
			connectFlag = true
			connectting = false
			Produce, err = sarama.NewAsyncProducerFromClient(Client)
			if err != nil {
				globallogger.Log.Errorf("[Kafka] [connectLoop] create producer error :%s", err.Error())
			}
			break
		}
	}
}

// NewClient 连接kafka
func NewClient(addrs []string) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = false
	config.Consumer.Return.Errors = true
	config.Version = sarama.V0_11_0_2
	connectting = true
	connectLoop(addrs, config)
	go func() {
		for {
			if !connectting && Client.Closed() {
				connectting = true
				connectLoop(addrs, config)
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

// Producer 生产者
func Producer(topic string, value string) {
	// send to chain
	Produce.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}

	select {
	case suc := <-Produce.Successes():
		globallogger.Log.Warnf("[Kafka] [Producer] partition: %d, offset: %d, timestamp: %s, topic: %s, value: %s",
			suc.Partition, suc.Offset, suc.Timestamp.String(), topic, value)
	case fail := <-Produce.Errors():
		globallogger.Log.Errorf("[Kafka] [Producer] failed, err: %s", fail.Err.Error())
	default:
		globallogger.Log.Warnf("[Kafka] [Producer] default, topic: %s, value: %s", topic, value)
	}
}

type consumerGroupHandler struct {
	name string
}

// Consumer 消费者
func Consumer() {
	var err error
	Group, err = sarama.NewConsumerGroupFromClient(constant.Constant.KAFKA.ZigbeeKafkaGroupName, Client)
	if err != nil {
		globallogger.Log.Errorf("[Kafka] [Consumer] create consumer error %s", err.Error())
		return
	}
	go consume(&Group, constant.Constant.KAFKA.ZigbeeKafkaGroupName)
}

func consume(group *sarama.ConsumerGroup, name string) {
	globallogger.Log.Infoln("[Kafka] [Consume]", name, "start")
	ctx := context.Background()
	topics := []string{
		constant.Constant.KAFKA.ZigbeeKafkaConsumeTopicBuildInPropertyByCustom,
		constant.Constant.KAFKA.ZigbeeKafkaConsumeTopicSetTransConfigurationByIottransparent,
		constant.Constant.KAFKA.ZigbeeKafkaConsumeTopicBuildInProperty,
		constant.Constant.KAFKA.ZigbeeKafkaConsumeTopicDelTerminal,
		constant.Constant.KAFKA.ZigbeeKafkaConsumeTopicUpdateTerminal,
	}
	handler := consumerGroupHandler{name: name}
	err := (*group).Consume(ctx, topics, handler)
	if err != nil {
		globallogger.Log.Errorln("[Kafka] [Consume]", err.Error())
	}
}
func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	defer func() {
		err := recover()
		if err != nil {
			globallogger.Log.Errorln("[Kafka] [Consume] err :", err)
		}
	}()
	for msg := range claim.Messages() {
		globallogger.Log.Warnf("[Kafka] [Consume] %s Message topic:%q partition:%d offset:%d  value:%s",
			h.name, msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		// 手动确认消息
		sess.MarkMessage(msg, "")
		switch msg.Topic {
		case constant.Constant.KAFKA.ZigbeeKafkaConsumeTopicBuildInPropertyByCustom,
			constant.Constant.KAFKA.ZigbeeKafkaConsumeTopicSetTransConfigurationByIottransparent:
			procBuildInPropertyByCustom(msg.Value)
		case constant.Constant.KAFKA.ZigbeeKafkaConsumeTopicBuildInProperty:
			procBuildInProperty(msg.Value)
		case constant.Constant.KAFKA.ZigbeeKafkaConsumeTopicDelTerminal:
			procDelTerminal(msg.Value)
		case constant.Constant.KAFKA.ZigbeeKafkaConsumeTopicUpdateTerminal:
			procUpdateTerminal(msg.Value)
		}
	}
	return nil
}

// func handleErrors(group *sarama.ConsumerGroup, wg *sync.WaitGroup) {
// 	wg.Done()
// 	for err := range (*group).Errors() {
// 		globallogger.Log.Errorln("[Kafka] [Consume] ERROR", err)
// 	}
// }

func procIntervalCommand(terminalInfo *config.TerminalInfo, interval int) {
	if constant.Constant.UsePostgres {
		models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": terminalInfo.DevEUI}, map[string]interface{}{"interval": interval})
	} else {
		models.FindTerminalAndUpdate(bson.M{"devEUI": terminalInfo.DevEUI}, bson.M{"interval": interval})
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
	var endpointTemp pq.StringArray
	if constant.Constant.UsePostgres {
		endpointTemp = terminalInfo.EndpointPG
	} else {
		endpointTemp = terminalInfo.Endpoint
	}
	var endpointIndex int = 0
	var clusterID uint16 = 0x0000
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM, constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
		if len(endpointTemp) > 1 && endpointTemp[1] == "01" {
			endpointIndex = 1
		}
		clusterID = 0x0001
	case constant.Constant.TMNTYPE.MAILEKE.ZigbeeTerminalPMT1004Detector, constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalCOSensorEM,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRSensorEM, constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalPIRILLSensorEF30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmokeSensorEM, constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWarningDevice,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalWaterSensorEM, constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSceneSwitchEM30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalDoorSensorEF30, constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalIRControlEM,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalPMTSensor0001112b, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalWarningDevice005b0e12,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSmokeSensorHY0024, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalPIRSensorHY0027:
		clusterID = 0x0001
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalGASSensorEM, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalGASSensorHY0022:
		clusterID = 0x0500
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalESocket, constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalSmartPlug,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW1LEFR30, constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW2LEFR30,
		constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2SW3LEFR30, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitch00500c32,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSingleSwitchHY0141, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitch00500c33,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalDoubleSwitchHY0142, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitch00500c35,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalTripleSwitchHY0143:
		clusterID = 0x0006
	case constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c3c, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocket000a0c55,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0105, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminalSocketHY0106,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal1SceneSwitch005f0cf1, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal2SceneSwitch005f0cf3,
		constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal3SceneSwitch005f0cf2, constant.Constant.TMNTYPE.HONYAR.ZigbeeTerminal6SceneSwitch005f0c3b:
		clusterID = 0x0000
	default:
		globallogger.Log.Warnf("[procIntervalCommand]: invalid tmnType: %s", terminalInfo.TmnType)
	}
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      terminalInfo.DevEUI,
		CommandType: common.IntervalCommand,
		ClusterID:   clusterID,
		Command: common.Command{
			Interval:         uint16(interval),
			DstEndpointIndex: endpointIndex,
		},
	})
}

func procAddScene(terminalInfo *config.TerminalInfo, sceneID uint8, sceneName string) {
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
		DevEUI:      terminalInfo.DevEUI,
		CommandType: common.AddScene,
		ClusterID:   uint16(cluster.Scenes),
		Command: common.Command{
			GroupID:          1,
			TransitionTime:   0,
			SceneName:        "0" + strconv.FormatInt(int64(len(sceneName)/2), 16) + sceneName,
			KeyID:            1,
			SceneID:          sceneID,
			DstEndpointIndex: endpointIndex,
		},
	})
}

func procOnOff(terminalInfo *config.TerminalInfo, value string) {
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
		DevEUI:      terminalInfo.DevEUI,
		CommandType: common.SwitchCommand,
		ClusterID:   uint16(cluster.OnOff),
		Command: common.Command{
			Cmd:              commandID,
			DstEndpointIndex: dstEndpointIndex,
		},
	})
}

func procSocketCommand(terminalInfo *config.TerminalInfo, value string, commandType string) {
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
		DevEUI:      terminalInfo.DevEUI,
		CommandType: commandType,
		ClusterID:   clusterID,
		Command:     cmd,
	})
}

func procIRControlEMCommand(devEUI string, cmd common.Command, commandType string) {
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: commandType,
		ClusterID:   uint16(cluster.HEIMANInfraredRemote),
		Command:     cmd,
	})
}

func procHS2AQEMTemperatureAlarmThresholdSet(terminalInfo *config.TerminalInfo, minValue int, maxValue int) {
	minValue = minValue * 100
	maxValue = maxValue * 100
	if minValue < 0 {
		minValue = minValue + 65536
	}
	if maxValue < 0 {
		maxValue = maxValue + 65536
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
		DevEUI:      terminalInfo.DevEUI,
		CommandType: common.SetThreshold,
		ClusterID:   uint16(cluster.HEIMANAirQualityMeasurement),
		Command: common.Command{
			DstEndpointIndex: endpointIndex,
			MaxTemperature:   int16(maxValue),
			MinTemperature:   int16(minValue),
			Disturb:          disturb,
		},
	})
}

func procHS2AQEMHumidityAlarmThresholdSet(terminalInfo *config.TerminalInfo, minValue int, maxValue int) {
	minValue = minValue * 100
	maxValue = maxValue * 100
	if minValue < 0 {
		minValue = minValue + 65536
	}
	if maxValue < 0 {
		maxValue = maxValue + 65536
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
		DevEUI:      terminalInfo.DevEUI,
		CommandType: common.SetThreshold,
		ClusterID:   uint16(cluster.HEIMANAirQualityMeasurement),
		Command: common.Command{
			DstEndpointIndex: endpointIndex,
			MaxHumidity:      uint16(maxValue),
			MinHumidity:      uint16(minValue),
			Disturb:          disturb,
		},
	})
}

func procHS2AQEMDisturbMode(terminalInfo *config.TerminalInfo, value string) {
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
		DevEUI:      terminalInfo.DevEUI,
		CommandType: common.SetThreshold,
		ClusterID:   uint16(cluster.HEIMANAirQualityMeasurement),
		Command: common.Command{
			DstEndpointIndex: endpointIndex,
			MaxTemperature:   int16(maxTemperature),
			MinTemperature:   int16(minTemperature),
			MaxHumidity:      uint16(maxHumidity),
			MinHumidity:      uint16(minHumidity),
			Disturb:          disturb,
		},
	})
}

func procHS2AQSetLanguage(devEUI string, language uint8) {
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SetLanguage,
		ClusterID:   uint16(cluster.HEIMANAirQualityMeasurement),
		Command: common.Command{
			Cmd: language,
		},
	})
}

func procHS2AQSetUnitOfTemperature(devEUI string, unitOfTemperature uint8) {
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
		ClusterID:   uint16(cluster.HEIMANAirQualityMeasurement),
		Command: common.Command{
			Cmd: unitOfTemperature,
		},
	})
}

func procWarningDeviceAlarm(devEUI string, alarmTime string, warningMode string) {
	time, _ := strconv.ParseUint(alarmTime, 10, 16)
	var warningControl uint8 = 0x17
	if alarmTime == "0" {
		warningControl = 0x03
	} else {
		switch warningMode {
		case "1": // Burglar
			warningControl = 0x17
		case "2": // Fire
			warningControl = 0x27
		case "3": // Emergency
			warningControl = 0x37
		case "4": // Police panic
			warningControl = 0x47
		case "5": // Fire panic
			warningControl = 0x57
		case "6": // Emergency panic
			warningControl = 0x67
		}
	}
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.StartWarning,
		ClusterID:   uint16(cluster.IASWarningDevice),
		Command: common.Command{
			DstEndpointIndex: 0,
			WarningControl:   warningControl,
			WarningTime:      uint16(time),
		},
	})
}

func procWarningDeviceSquawk(devEUI string, warningMode string) {
	var squawkMode uint8
	switch warningMode {
	case "0": // armed
		squawkMode = 0x03
	case "1": // disarmed
		squawkMode = 0x13
	}
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.Squawk,
		ClusterID:   uint16(cluster.IASWarningDevice),
		Command: common.Command{
			DstEndpointIndex: 0,
			SquawkMode:       squawkMode,
		},
	})
}

func procInterval(terminalInfo *config.TerminalInfo, clusterID uint16, interval uint16) {
	var endpointTemp pq.StringArray
	if constant.Constant.UsePostgres {
		endpointTemp = terminalInfo.EndpointPG
	} else {
		endpointTemp = terminalInfo.Endpoint
	}
	endpointIndex := 0
	switch terminalInfo.TmnType {
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHS2AQEM:
		if len(endpointTemp) > 1 && endpointTemp[1] == "01" {
			endpointIndex = 1
		}
	case constant.Constant.TMNTYPE.HEIMAN.ZigbeeTerminalHTEM:
		if len(endpointTemp) > 1 && endpointTemp[1] == "01" {
			if clusterID == uint16(cluster.TemperatureMeasurement) {
				endpointIndex = 1
			} else if clusterID == uint16(cluster.RelativeHumidityMeasurement) {
				endpointIndex = 0
			}
		} else {
			if clusterID == uint16(cluster.TemperatureMeasurement) {
				endpointIndex = 0
			} else if clusterID == uint16(cluster.RelativeHumidityMeasurement) {
				endpointIndex = 1
			}
		}
	}
	zclmsgdown.ProcZclDownMsg(common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      terminalInfo.DevEUI,
		CommandType: common.IntervalCommand,
		ClusterID:   clusterID,
		Command: common.Command{
			Interval:         interval,
			DstEndpointIndex: endpointIndex,
		},
	})
}

func procBuildInPropertyByCustom(value []byte) {
	type data struct {
		CapacityIdentifier string `json:"capacityIdentifier"`
		NeedAck            string `json:"needAck"`
		PropertyTopic      string `json:"propertyTopic"`
		ProductKey         string `json:"productKey"`
		Addr               string `json:"addr"`
		Command            string `json:"command"`
	}
	d := data{}
	err := json.Unmarshal(value, &d)
	if err == nil {
		var terminalInfo *config.TerminalInfo = nil
		if constant.Constant.UsePostgres {
			terminalInfo, err = models.GetTerminalInfoByDevEUIPG(d.Addr)
		} else {
			terminalInfo, err = models.GetTerminalInfoByDevEUI(d.Addr)
		}
		if err == nil && terminalInfo != nil {
			switch d.CapacityIdentifier {
			case "keepAlive":
				interval, _ := strconv.Atoi(d.Command)
				procIntervalCommand(terminalInfo, interval)
			case "addScene":
				commandArr := strings.Split(d.Command, ",")
				if len(commandArr) == 1 {
					commandArr = strings.Split(d.Command, "，")
				}
				sceneID, _ := strconv.Atoi(commandArr[0])
				procAddScene(terminalInfo, uint8(sceneID), commandArr[1])
			case "onOff":
				procOnOff(terminalInfo, "1"+d.Command)
			case "onOffOne":
				procOnOff(terminalInfo, "1"+d.Command)
			case "onOffTwo":
				procOnOff(terminalInfo, "2"+d.Command)
			case "onOffThree":
				procOnOff(terminalInfo, "3"+d.Command)
			case "childLock":
				procSocketCommand(terminalInfo, "1"+d.Command, common.ChildLock)
			case "USB":
				procSocketCommand(terminalInfo, "1"+d.Command, common.OnOffUSB)
			case "backLight":
				procSocketCommand(terminalInfo, "1"+d.Command, common.BackGroundLight)
			case "memory":
				procSocketCommand(terminalInfo, "1"+d.Command, common.PowerOffMemory)
			case "clearHistoryElectric":
				procSocketCommand(terminalInfo, "11", common.HistoryElectricClear)
			case "sendKeyCommand":
				commandArr := strings.Split(d.Command, ",")
				if len(commandArr) == 1 {
					commandArr = strings.Split(d.Command, "，")
				}
				infraredRemoteID, _ := strconv.Atoi(commandArr[0])
				infraredRemoteKeyCode, _ := strconv.Atoi(commandArr[1])
				cmd := common.Command{
					InfraredRemoteID:      uint8(infraredRemoteID),
					InfraredRemoteKeyCode: uint8(infraredRemoteKeyCode),
				}
				procIRControlEMCommand(terminalInfo.DevEUI, cmd, common.SendKeyCommand)
			case "studyKey":
				commandArr := strings.Split(d.Command, ",")
				if len(commandArr) == 1 {
					commandArr = strings.Split(d.Command, "，")
				}
				infraredRemoteID, _ := strconv.Atoi(commandArr[0])
				infraredRemoteKeyCode, _ := strconv.Atoi(commandArr[1])
				cmd := common.Command{
					InfraredRemoteID:      uint8(infraredRemoteID),
					InfraredRemoteKeyCode: uint8(infraredRemoteKeyCode),
				}
				procIRControlEMCommand(terminalInfo.DevEUI, cmd, common.StudyKey)
			case "deleteKey":
				commandArr := strings.Split(d.Command, ",")
				if len(commandArr) == 1 {
					commandArr = strings.Split(d.Command, "，")
				}
				infraredRemoteID, _ := strconv.Atoi(commandArr[0])
				infraredRemoteKeyCode, _ := strconv.Atoi(commandArr[1])
				cmd := common.Command{
					InfraredRemoteID:      uint8(infraredRemoteID),
					InfraredRemoteKeyCode: uint8(infraredRemoteKeyCode),
				}
				procIRControlEMCommand(terminalInfo.DevEUI, cmd, common.DeleteKey)
			case "createID":
				infraredRemoteModelType, _ := strconv.Atoi(d.Command)
				cmd := common.Command{
					InfraredRemoteModelType: uint8(infraredRemoteModelType),
				}
				procIRControlEMCommand(terminalInfo.DevEUI, cmd, common.CreateID)
			case "getIDAndKeyCodeList":
				procIRControlEMCommand(terminalInfo.DevEUI, common.Command{}, common.GetIDAndKeyCodeList)
			case "temperatureMinMax":
				commandArr := strings.Split(d.Command, ",")
				if len(commandArr) == 1 {
					commandArr = strings.Split(d.Command, "，")
				}
				minValue, _ := strconv.Atoi(commandArr[0])
				maxValue, _ := strconv.Atoi(commandArr[1])
				procHS2AQEMTemperatureAlarmThresholdSet(terminalInfo, minValue, maxValue)
			case "humidityMinMax":
				commandArr := strings.Split(d.Command, ",")
				if len(commandArr) == 1 {
					commandArr = strings.Split(d.Command, "，")
				}
				minValue, _ := strconv.Atoi(commandArr[0])
				maxValue, _ := strconv.Atoi(commandArr[1])
				procHS2AQEMHumidityAlarmThresholdSet(terminalInfo, minValue, maxValue)
			case "disturbMode":
				procHS2AQEMDisturbMode(terminalInfo, d.Command)
			case "language":
				language, _ := strconv.Atoi(d.Command)
				procHS2AQSetLanguage(terminalInfo.DevEUI, uint8(language))
			case "unitOfTemperature":
				unitOfTemperature, _ := strconv.Atoi(d.Command)
				if unitOfTemperature == 1 {
					unitOfTemperature = 0
				} else if unitOfTemperature == 0 {
					unitOfTemperature = 1
				}
				procHS2AQSetUnitOfTemperature(terminalInfo.DevEUI, uint8(unitOfTemperature))
			case "warning":
				commandArr := strings.Split(d.Command, ",")
				if len(commandArr) == 1 {
					commandArr = strings.Split(d.Command, "，")
				}
				if len(commandArr) == 1 {
					procWarningDeviceAlarm(terminalInfo.DevEUI, d.Command, "1")
				} else {
					procWarningDeviceAlarm(terminalInfo.DevEUI, commandArr[0], commandArr[1])
				}
			case "squawk":
				procWarningDeviceSquawk(terminalInfo.DevEUI, d.Command)
			case "temperatureInterval":
				interval, _ := strconv.Atoi(d.Command)
				procInterval(terminalInfo, uint16(cluster.TemperatureMeasurement), uint16(interval))
			case "humidityInterval":
				interval, _ := strconv.Atoi(d.Command)
				procInterval(terminalInfo, uint16(cluster.RelativeHumidityMeasurement), uint16(interval))
			}
		}
	}
}

func procBuildInProperty(value []byte) {
	type data struct {
		DevType       string `json:"devType"`
		DownTopic     string `json:"downTopic"`
		MsgType       string `json:"msgType"`
		SubAddrLen    int    `json:"subAddrLen"`
		Ctrl          string `json:"ctrl"`
		PortIDLen     int    `json:"portIDLen"`
		DataLen       int    `json:"dataLen"`
		AddrLen       int    `json:"addrLen"`
		PropertyName  string `json:"propertyName"`
		NeedAck       string `json:"needAck"`
		PropertyTopic string `json:"propertyTopic"`
		Interval      string `json:"interval"`
		Addr          string `json:"addr"`
	}
	d := data{}
	err := json.Unmarshal(value, &d)
	if err == nil {
		interval, _ := strconv.Atoi(d.Interval)
		var terminalInfo *config.TerminalInfo = nil
		if constant.Constant.UsePostgres {
			terminalInfo, err = models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": d.Addr}, map[string]interface{}{"interval": interval})
		} else {
			terminalInfo, err = models.FindTerminalAndUpdate(bson.M{"devEUI": d.Addr}, bson.M{"interval": interval})
		}
		if err == nil && terminalInfo != nil {
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
			switch d.PropertyTopic {
			case "0x00000003":
				procIntervalCommand(terminalInfo, interval)
			}
		}
	}
}

func procDelTerminal(value []byte) {
	type data struct {
		AccessType  string `json:"accessType"`
		SceneID     string `json:"sceneID"`
		TenantID    string `json:"tenantID"`
		TmnDevSN    string `json:"tmnDevSN"`
		ProductKey  string `json:"productKey"`
		TmnOIDIndex string `json:"tmnOIDIndex"`
	}
	d := data{}
	err := json.Unmarshal(value, &d)
	if err == nil {
		publicfunction.ProcTerminalDeleteByOIDIndex(d.TmnOIDIndex)
	}
}

func procUpdateTerminal(value []byte) {
	type data struct {
		TmnName     string `json:"tmnName"`
		TmnOIDIndex string `json:"tmnOIDIndex"`
	}
	d := data{}
	err := json.Unmarshal(value, &d)
	if err == nil {
		var terminalInfo *config.TerminalInfo = nil
		if constant.Constant.UsePostgres {
			oSet := make(map[string]interface{}, 1)
			oSet["tmnname"] = d.TmnName
			terminalInfo, err = models.FindTerminalAndUpdateByOIDIndexPG(d.TmnOIDIndex, oSet)
			oSet = nil
			if err != nil {
				globallogger.Log.Errorln("[Kafka] [Consume] OIDIndex :", d.TmnOIDIndex, "procUpdateTerminal FindTerminalAndUpdateByOIDIndexPG err :", err)
			}
		} else {
			terminalInfo, err = models.FindTerminalAndUpdateByOIDIndex(d.TmnOIDIndex, bson.M{"tmnName": d.TmnName})
			if err != nil {
				globallogger.Log.Errorln("[Kafka] [Consume] OIDIndex :", d.TmnOIDIndex, "procUpdateTerminal FindTerminalAndUpdateByOIDIndex err :", err)
			}
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
	}
}
