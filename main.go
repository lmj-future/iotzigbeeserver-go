package main

import (
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/controllers"
	"github.com/h3c/iotzigbeeserver-go/datainit"
	"github.com/h3c/iotzigbeeserver-go/db/kafka"
	"github.com/h3c/iotzigbeeserver-go/db/mongo"
	"github.com/h3c/iotzigbeeserver-go/db/mqtt"
	"github.com/h3c/iotzigbeeserver-go/db/postgres"
	"github.com/h3c/iotzigbeeserver-go/db/redis"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmemorycache"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalredisclient"
	"github.com/h3c/iotzigbeeserver-go/httpapi"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/rabbitmqmsg/consumer"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
)

var defaultConfigs = make(map[string]interface{})

func loggerInit() {
	globallogger.Init(defaultConfigs["logger"].(map[string]interface{}))
}

func connectToRedis() error {
	globallogger.Log.WithField("db", "redis").Infoln("create redis client")
	redisConnParas := defaultConfigs["redisConnParas"].(map[string]interface{})
	host := redisConnParas["host"].(string)
	port := redisConnParas["port"].(string)
	password := redisConnParas["password"].(string)
	if os.Getenv(host) != "" {
		host = os.Getenv(host)
	}
	if os.Getenv(port) != "" {
		port = os.Getenv(port)
	}
	if os.Getenv(password) != "" {
		password = os.Getenv(password)
	}
	var err error
	if redisConnParas["useAlone"].(bool) {
		globalredisclient.MyZigbeeServerRedisClient, err = redis.NewClientPool(host, port, password)
	} else {
		servers := []string{host + ":" + port}
		globalredisclient.MyZigbeeServerRedisClient, err = redis.NewClusterClient(servers, password)
	}
	if err != nil {
		globallogger.Log.Errorln("couldn't create redis client:", err.Error())
		return err
	}
	globallogger.Log.Errorf("create redis client sucess: %+v", globalredisclient.MyZigbeeServerRedisClient)
	return nil
}

func reConnectMongoDB(mongoURL string) {
	var err error
	timer := time.NewTimer(30 * time.Second)
	<-timer.C
	mongo.MongoClient, err = mongo.NewClient(mongoURL)
	if err != nil {
		globallogger.Log.Errorln("[MongoDB][Reconnect] couldn't open mongo:", err.Error())
		reConnectMongoDB(mongoURL)
	} else {
		globallogger.Log.Warnf("[MongoDB][Reconnect] connect success: %+v", mongo.MongoClient)
		keepAliveMongoDB(mongoURL, mongo.MongoClient)
	}
}
func keepAliveMongoDB(mongoURL string, client *mongo.Client) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		<-ticker.C
		err := client.Session.Ping()
		if err != nil {
			mongo.MongoClient, err = mongo.NewClient(mongoURL)
			if err != nil {
				globallogger.Log.Errorln("[MongoDB][KeepAlive] couldn't open mongo:", err.Error())
			} else {
				globallogger.Log.Errorf("[MongoDB][KeepAlive] connect success: %+v", mongo.MongoClient)
			}
		} else {
			globallogger.Log.Warnf("[MongoDB][KeepAlive] Ping success")
		}
	}
}
func connectToMongoDB() error {
	mongoConnParas := defaultConfigs["mongoConnParas"].(map[string]interface{})
	userName := mongoConnParas["userName"].(string)
	password := mongoConnParas["password"].(string)
	host := mongoConnParas["host"].(string)
	port := mongoConnParas["port"].(string)
	dbName := mongoConnParas["dbName"].(string)
	if os.Getenv(userName) != "" {
		userName = os.Getenv(userName)
	}
	if os.Getenv(password) != "" {
		password = os.Getenv(password)
	}
	if os.Getenv(host) != "" {
		host = os.Getenv(host)
	}
	if os.Getenv(port) != "" {
		port = os.Getenv(port)
	}
	if os.Getenv(dbName) != "" {
		dbName = os.Getenv(dbName)
	}
	mongoURL := "mongodb://" + userName + ":" + password + "@" + host + ":" + port + "/" + dbName
	// Create a database client
	client, err := mongo.NewClient(mongoURL)
	if err != nil {
		globallogger.Log.Errorln("couldn't create mongodb client:", err.Error())
		go reConnectMongoDB(mongoURL)
		return err
	}
	globallogger.Log.Errorf("[Mongodb] connect success: %+v", client)
	mongo.MongoClient = client
	go keepAliveMongoDB(mongoURL, client)
	return nil
}

func connectToMQTT() {
	mqttConnParas := defaultConfigs["mqttConnParas"].(map[string]interface{})
	userName := mqttConnParas["userName"].(string)
	password := mqttConnParas["password"].(string)
	host := mqttConnParas["host"].(string)
	port := mqttConnParas["port"].(string)
	if os.Getenv(userName) != "" {
		userName = os.Getenv(userName)
	}
	if os.Getenv(password) != "" {
		password = os.Getenv(password)
	}
	if os.Getenv(host) != "" {
		host = os.Getenv(host)
	}
	if os.Getenv(port) != "" {
		port = os.Getenv(port)
	}
	mqtt.ConnectMQTT(host, port, userName, password)
}

func subscribeFromMQTT() {
	subTopics := make([]string, 0)
	if constant.Constant.Iotware {
		subTopics = append(subTopics, iotsmartspace.TopicV1GatewayRPC)
		subTopics = append(subTopics, iotsmartspace.TopicV1GatewayNetworkInAck)
		subTopics = append(subTopics, iotsmartspace.TopicV1GatewayEventDeviceDelete)
	} else if constant.Constant.Iotedge {
		subTopics = append(subTopics, iotsmartspace.TopicIotsmartspaceZigbeeserverAction)
		subTopics = append(subTopics, iotsmartspace.TopicIotsmartspaceZigbeeserverProperty)
		subTopics = append(subTopics, iotsmartspace.TopicIotsmartspaceZigbeeserverState)
		subTopics = append(subTopics, common.ReadAttribute)
		subTopics = append(subTopics, "scenes")
	}
	for _, topic := range subTopics {
		mqtt.Subscribe(topic, func(topic string, msg []byte) {
			iotsmartspace.ProcSubMsg(topic, msg)
		})
	}
}

func connectToRabbitMQ() {
	consumer.RabbitMQConnection()
}

func connectToKafka(kafkaConnParas map[string]interface{}) {
	addrs := []string{}
	host1 := kafkaConnParas["host1"].(string)
	port1 := kafkaConnParas["port1"].(string)
	host2 := kafkaConnParas["host2"].(string)
	port2 := kafkaConnParas["port2"].(string)
	host3 := kafkaConnParas["host3"].(string)
	port3 := kafkaConnParas["port3"].(string)
	if os.Getenv(host1) != "" {
		host1 = os.Getenv(host1)
		if os.Getenv(port1) != "" {
			port1 = os.Getenv(port1)
		}
	}
	if os.Getenv(host2) != "" {
		host2 = os.Getenv(host2)
		if os.Getenv(port2) != "" {
			port2 = os.Getenv(port2)
		}
	}
	if os.Getenv(host3) != "" {
		host3 = os.Getenv(host3)
		if os.Getenv(port3) != "" {
			port3 = os.Getenv(port3)
		}
	}
	addrs = append(addrs, host1+":"+port1)
	addrs = append(addrs, host2+":"+port2)
	addrs = append(addrs, host3+":"+port3)
	kafka.NewClient(addrs)
}

// func terminalStateSmooth() {
// 	var oMatch = map[string]interface{}{}
// 	var oMatchPG = map[string]interface{}{}
// 	oMatch["isExist"] = true
// 	oMatchPG["isexist"] = true
// 	var terminalList []config.TerminalInfo
// 	var err error
// 	if constant.Constant.UsePostgres {
// 		terminalList, err = models.FindAllTerminalByConditionPG(oMatchPG)
// 	} else {
// 		terminalList, err = models.FindAllTerminalByCodition(oMatch)
// 	}
// 	if err == nil {
// 		for _, terminalInfo := range terminalList {
// 			if terminalInfo.Online {
// 				publicfunction.TerminalOnline(terminalInfo.DevEUI, false)
// 				if constant.Constant.Iotware {
// 					iotsmartspace.StateTerminalOnlineIotware(terminalInfo)
// 				} else if constant.Constant.Iotedge {
// 					iotsmartspace.StateTerminalOnline(terminalInfo.DevEUI)
// 				}
// 			} else {
// 				publicfunction.TerminalOffline(terminalInfo.DevEUI)
// 				if constant.Constant.Iotware {
// 					iotsmartspace.StateTerminalOfflineIotware(terminalInfo)
// 				} else if constant.Constant.Iotedge {
// 					iotsmartspace.StateTerminalOffline(terminalInfo.DevEUI)
// 				}
// 			}
// 		}
// 	}
// }
func getLocalHost() {
	if constant.Constant.Wcg {
		constant.Constant.LocalHost = "172.25.252.9"
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			os.Exit(1)
		}
		for _, address := range addrs {
			globallogger.Log.Errorln(address)
			// 检查ip地址判断是否回环地址
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					if ipnet.IP.String() == "172.25.252.9" {
						constant.Constant.LocalHost = ipnet.IP.String()
					}
				}
			}
		}
	} else {
		constant.Constant.LocalHost = "0.0.0.0"
	}
}
func main() {
	// err := config.LoadJSON("./config/default.json", &defaultConfigs)
	err := config.LoadJSON("./config/production.json", &defaultConfigs)
	if err != nil {
		panic(err)
	}
	loggerInit()
	globallogger.Log.Errorln("**************** h3c-zigbeeserver start ****************")

	/* 连接数据库 */
	if defaultConfigs["usePostgres"].(bool) {
		// Postgres
		postgres.GetInstance().InitDataPool(defaultConfigs["postgresConnParas"].(map[string]interface{}))
	} else {
		// MongoDB
		connectToMongoDB()
	}
	if defaultConfigs["multipleInstances"].(bool) {
		/* 连接redis 多实例使用redis 单实例使用内存 */
		connectToRedis()
	}
	if defaultConfigs["useMQTT"].(bool) {
		go func() {
			/* 连接MQTT */
			connectToMQTT()
		}()
	}
	if defaultConfigs["useRabbitMQ"].(bool) {
		/* 连接rabbitMQ */
		connectToRabbitMQ()
	}
	if defaultConfigs["useKafka"].(bool) {
		/* 连接kafka */
		connectToKafka(defaultConfigs["kafkaConnParas"].(map[string]interface{}))
		go func() {
			for {
				if kafka.GetConnectFlag() {
					kafka.Consumer()
					kafka.SetConnectFlag(false)
				}
				time.Sleep(time.Minute)
			}
		}()
	}
	/* 初始化表项 */
	datainit.IotzigbeeserverDataInit()
	/* 获取本地地址 */
	getLocalHost()
	/* UDP监听 */
	controllers.CreateUDPServer(int(defaultConfigs["udpPort"].(float64)))
	/* HTTP API */
	go httpapi.NewServer()
	/* Terminal State Smooth */
	go func() {
		for {
			if mqtt.GetConnectFlag() {
				/* MQTT订阅消息 */
				subscribeFromMQTT()
				mqtt.SetConnectFlag(false)
				// go func() {
				// 	var timer = time.NewTimer(5 * time.Minute)
				// 	defer timer.Stop()
				// 	<-timer.C
				// 	if !mqtt.GetConnectFlag() {
				// 		terminalStateSmooth()
				// 	}
				// }()
			}
			time.Sleep(time.Minute)
		}
	}()

	// memoryCache内存监测
	if defaultConfigs["memoryCacheMonitor"].(bool) {
		var tickerID = time.NewTicker(time.Duration(10*60) * time.Second)
		defer tickerID.Stop()
		go func() {
			for {
				<-tickerID.C
				globalmemorycache.MemoryCache.GetMemorySize()
			}
		}()
	}

	// prometheus监测
	// if defaultConfigs["prometheus"].(bool) {
	// 	go func() {
	// 		metrics.PrometheusStart()
	// 	}()
	// }

	go func() {
		http.ListenAndServe("0.0.0.0:80", nil)
	}()

	wait(defaultConfigs["multipleInstances"].(bool), defaultConfigs["useKafka"].(bool))
	globallogger.Log.Errorln("**************** h3c-zigbeeserver exit ****************")
}

func wait(multipleInstances bool, useKafka bool) { //<-chan os.Signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	signal.Ignore(syscall.SIGPIPE)
	<-sig
	if multipleInstances {
		globalredisclient.MyZigbeeServerRedisClient.CloseSession()
	}
	if useKafka {
		kafka.Group.Close()
		kafka.Produce.AsyncClose()
	}
}
