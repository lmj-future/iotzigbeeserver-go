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
	globallogger.Log.Warnf("Pub Client Topic : %s ", msg.Topic())
	globallogger.Log.Warnf("Pub Client msg : %s ", msg.Payload())
}

//创建全局mqtt sub消息处理 handler
// var messageSubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
// 	fmt.Printf("Sub Client Topic : %s \n", msg.Topic())
// 	fmt.Printf("Sub Client msg : %s \n", msg.Payload())
// }

var client mqtt.Client
var taskID = "h3c-zigbeeserver"

func keepAlive(mqttHost string, mqttPort string, mqttUserName string, mqttPassword string) {
	timer := time.NewTimer(time.Minute)
	select {
	case <-timer.C:
		if client.IsConnected() {
			globallogger.Log.Warnf("[MQTT][KeepAlive] Ping success")
			keepAlive(mqttHost, mqttPort, mqttUserName, mqttPassword)
		} else {
			ConnectMQTT(mqttHost, mqttPort, mqttUserName, mqttPassword)
			//客户端连接判断
			if token := client.Connect(); !token.WaitTimeout(time.Duration(5)*time.Second) && token.Wait() && token.Error() != nil {
				globallogger.Log.Errorf("[MQTT][KeepAlive] mqtt connect error, taskID: %s, error: %s ", taskID, token.Error())
				//客户端重连
				ConnectMQTT(mqttHost, mqttPort, mqttUserName, mqttPassword)
			} else {
				globallogger.Log.Infof("[MQTT][KeepAlive] connect success %+v, %+v, %+v, %+v, %+v", token, mqttHost, mqttPort, mqttUserName, mqttPassword)
				keepAlive(mqttHost, mqttPort, mqttUserName, mqttPassword)
			}
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
	//设置连接超时
	clientOptions.SetConnectTimeout(time.Duration(5) * time.Second)
	//设置保活时长
	// clientOptions.SetKeepAlive(60 * time.Second)
	//创建客户端连接
	client = mqtt.NewClient(clientOptions)
	//客户端连接判断
	if token := client.Connect(); !token.WaitTimeout(time.Duration(5)*time.Second) && token.Wait() && token.Error() != nil {
		globallogger.Log.Errorf("[MQTT] mqtt connect error, taskID: %s, error: %s ", taskID, token.Error())
		//客户端重连
		ConnectMQTT(mqttHost, mqttPort, mqttUserName, mqttPassword)
	} else {
		globallogger.Log.Infof("[MQTT] connect success %+v, %+v, %+v, %+v, %+v", token, mqttHost, mqttPort, mqttUserName, mqttPassword)
		// keepAlive(mqttHost, mqttPort, mqttUserName, mqttPassword)
	}
}

// Publish Publish
func Publish(topic string, msg string) {
	//发布消息
	client.Publish(topic, 1, false, msg)
	// globallogger.Log.Warnf("[Pub] end publish msg to mqtt broker, taskID: %s, token : %s ", taskID, token)
	// token.Wait()
}

// Subscribe Subscribe
func Subscribe(topic string, cb func(topic string, msg []byte)) {
	client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		globallogger.Log.Warnf("Sub Client Topic : %s, msg : %s ", msg.Topic(), msg.Payload())
		go cb(msg.Topic(), msg.Payload())
	})
	// fmt.Printf("[Pub] end Subscribe msg from mqtt broker, taskID: %s, token : %s ", taskID, token)
	// token.Wait()
}
