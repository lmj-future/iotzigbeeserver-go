package mqtt

import (
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
)

//创建全局mqtt publish消息处理 handler
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	globallogger.Log.Warnf("Pub Client Topic : %s", msg.Topic())
	globallogger.Log.Warnf("Pub Client msg : %s", msg.Payload())
}

//创建全局mqtt sub消息处理 handler
// var messageSubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
// 	fmt.Printf("Sub Client Topic : %s \n", msg.Topic())
// 	fmt.Printf("Sub Client msg : %s \n", msg.Payload())
// }

var client mqtt.Client
var taskID = "h3c-zigbeeserver"
var connectting = false
var connectFlag = false

func connect(clientOptions *mqtt.ClientOptions) error {
	//创建客户端连接
	client = mqtt.NewClient(clientOptions)
	//客户端连接判断
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		globallogger.Log.Errorf("[MQTT] mqtt connect error, taskID: %s, error: %s", taskID, token.Error())
		return token.Error()
	} else {
		globallogger.Log.Errorf("[MQTT] connect success taskID:%s, token:%+v, clientOptions:%+v", taskID, token, clientOptions)
		return nil
	}
}
func connectLoop(clientOptions *mqtt.ClientOptions) {
	for {
		if err := connect(clientOptions); err != nil {
			time.Sleep(time.Second * 2)
		} else {
			connectFlag = true
			connectting = false
			break
		}
	}
}

// ConnectMQTT ConnectMQTT
func ConnectMQTT(mqttHost string, mqttPort string, mqttUserName string, mqttPassword string) {
	//设置连接参数
	clientOptions := mqtt.NewClientOptions().AddBroker("tcp://" + mqttHost + ":" + mqttPort).SetUsername(mqttUserName).SetPassword(mqttPassword)
	//设置客户端ID
	clientOptions.SetClientID(fmt.Sprintf("%s-%d", taskID, rand.Int()))
	//设置handler
	clientOptions.SetDefaultPublishHandler(messagePubHandler)
	//设置保活时长
	clientOptions.SetKeepAlive(30 * time.Second)
	clientOptions.SetAutoReconnect(true)
	clientOptions.SetMaxReconnectInterval(time.Minute)
	connectting = true
	connectLoop(clientOptions)
	go func() {
		for {
			if !connectting && !client.IsConnectionOpen() {
				connectting = true
				client.Disconnect(250)
				connectLoop(clientOptions)
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

// Publish Publish
func Publish(topic string, msg string) {
	//发布消息
	client.Publish(topic, 1, false, msg)
}

// Subscribe Subscribe
func Subscribe(topic string, cb func(topic string, msg []byte)) {
	client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		globallogger.Log.Warnf("Sub Client Topic : %s, msg : %s", msg.Topic(), msg.Payload())
		go cb(msg.Topic(), msg.Payload())
	})
}

func GetConnectFlag() bool {
	return connectFlag
}

func SetConnectFlag(flag bool) {
	connectFlag = flag
}
