package rabbitmq

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/streadway/amqp"
)

// 定义全局变量,指针类型
// var mqConn *amqp.Connection
// var mqChan *amqp.Channel

// Producer 定义生产者接口
type Producer interface {
	MsgContent() string
}

// Receiver 定义接收者接口
type Receiver interface {
	Consumer(amqp.Delivery) error
	ConsumerByIotwebserver(amqp.Delivery) (Response, error)
}

// RabbitMQ 定义RabbitMQ对象
type RabbitMQ struct {
	connection   *amqp.Connection
	channel      *amqp.Channel
	appID        string // application id(server name)
	messageID    string // message id
	queueName    string // 队列名称
	routingKey   string // key名称
	exchangeName string // 交换机名称
	exchangeType string // 交换机类型
	producerList []Producer
	receiverList []Receiver
	mu           sync.RWMutex
}

// QueueExchange 定义队列交换机对象
type QueueExchange struct {
	QuName string // 队列名称
	RtKey  string // key值
	ExName string // 交换机名称
	ExType string // 交换机类型
}

// Response Response
type Response struct {
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
}

// MsgContent 实现发送者
func (res *Response) MsgContent() string {
	tJSONByte, _ := json.Marshal(res)
	globallogger.Log.Warnf("[RabbitMQ] reply to iotwebserver msg: %+v\n", string(tJSONByte))
	return string(tJSONByte)
}

// 链接rabbitMQ
func (r *RabbitMQ) mqConnect() {
	var err error
	var defaultConfigs = make(map[string]interface{})
	err = config.LoadJSON("./config/production.json", &defaultConfigs)
	if err != nil {
		panic(err)
	}
	rabbitMQConnParas := defaultConfigs["rabbitMQConnParas"].(map[string]interface{})
	userName := rabbitMQConnParas["userName"].(string)
	password := rabbitMQConnParas["password"].(string)
	host := rabbitMQConnParas["host"].(string)
	port := rabbitMQConnParas["port"].(string)
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
	RabbitURL := "amqp://" + userName + ":" + password + "@" + host + ":" + port + "/"
	mqConn, err := amqp.Dial(RabbitURL)
	r.connection = mqConn // 赋值给RabbitMQ对象
	if err != nil {
		globallogger.Log.Errorf("MQ打开链接失败:%s", err)
	}
	mqChan, err := mqConn.Channel()
	r.channel = mqChan // 赋值给RabbitMQ对象
	if err != nil {
		globallogger.Log.Errorf("MQ打开管道失败:%s", err)
	}
	globallogger.Log.Warnf("[RabbitMQ] connect success: %+v", r)
}

// 关闭RabbitMQ连接
func (r *RabbitMQ) mqClose() {
	// 先关闭管道,再关闭链接
	err := r.channel.Close()
	if err != nil {
		globallogger.Log.Errorf("MQ管道关闭失败:%s", err)
	}
	err = r.connection.Close()
	if err != nil {
		globallogger.Log.Errorf("MQ链接关闭失败:%s", err)
	}
}

// New 创建一个新的操作对象
func New(q *QueueExchange, appID string, messageID string) *RabbitMQ {
	return &RabbitMQ{
		appID:        appID,
		messageID:    messageID,
		queueName:    q.QuName,
		routingKey:   q.RtKey,
		exchangeName: q.ExName,
		exchangeType: q.ExType,
	}
}

// StartProducer 启动RabbitMQ客户端,并初始化
func (r *RabbitMQ) StartProducer() {
	// 开启监听生产者发送任务
	for _, producer := range r.producerList {
		go r.listenProducer(producer)
	}
}

// StartReceiver 启动RabbitMQ客户端,并初始化
func (r *RabbitMQ) StartReceiver() {
	// 开启监听接收者接收任务
	for _, receiver := range r.receiverList {
		go r.listenReceiver(receiver)
	}
}

// RegisterProducer 注册发送指定队列指定路由的生产者
func (r *RabbitMQ) RegisterProducer(producer Producer) {
	r.producerList = append(r.producerList, producer)
}

// 发送任务
func (r *RabbitMQ) listenProducer(producer Producer) {
	// 验证链接是否正常,否则重新链接
	if r.channel == nil {
		r.mqConnect()
	}
	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err := r.channel.QueueDeclarePassive(r.queueName, true, false, false, true, nil)
	if err != nil {
		// 队列不存在,声明队列
		// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
		// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
		_, err = r.channel.QueueDeclare(r.queueName, true, false, false, true, nil)
		if err != nil {
			globallogger.Log.Errorf("MQ注册队列失败:%s", err)
			return
		}
	}
	// 队列绑定
	err = r.channel.QueueBind(r.queueName, r.routingKey, r.exchangeName, true, nil)
	if err != nil {
		globallogger.Log.Errorf("MQ绑定队列失败:%s", err)
		return
	}
	// 用于检查交换机是否存在,已经存在不需要重复声明
	err = r.channel.ExchangeDeclarePassive(r.exchangeName, r.exchangeType, true, false, false, true, nil)
	if err != nil {
		// 注册交换机
		// name:交换机名称,kind:交换机类型,durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;
		// noWait:是否非阻塞, true为是,不等待RMQ返回信息;args:参数,传nil即可; internal:是否为内部
		err = r.channel.ExchangeDeclare(r.exchangeName, r.exchangeType, true, false, false, true, nil)
		if err != nil {
			globallogger.Log.Errorf("MQ注册交换机失败:%s", err)
			return
		}
	}
	// 发送任务消息
	err = r.channel.Publish(r.exchangeName, r.routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(producer.MsgContent()),
		AppId:       r.appID,
		MessageId:   r.messageID,
	})
	if err != nil {
		globallogger.Log.Errorf("MQ任务发送失败:%s", err)
		return
	}
}

// RegisterReceiver 注册接收指定队列指定路由的数据接收者
func (r *RabbitMQ) RegisterReceiver(receiver Receiver) {
	r.mu.Lock()
	r.receiverList = append(r.receiverList, receiver)
	r.mu.Unlock()
}

// 监听接收者接收任务
func (r *RabbitMQ) listenReceiver(receiver Receiver) {
	// 处理结束关闭链接
	defer r.mqClose()
	// 验证链接是否正常
	if r.channel == nil {
		r.mqConnect()
	}
	// 用于检查队列是否存在,已经存在不需要重复声明
	// _, err := r.channel.QueueDeclarePassive(r.queueName, true, false, false, true, nil)
	// if err != nil {
	// 队列不存在,声明队列
	// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
	// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
	_, err := r.channel.QueueDeclare(r.queueName, true, false, false, true, nil)
	if err != nil {
		globallogger.Log.Errorf("MQ注册队列失败:%s", err)
		return
	}
	// }
	// 绑定任务
	err = r.channel.QueueBind(r.queueName, r.routingKey, r.exchangeName, true, nil)
	if err != nil {
		globallogger.Log.Errorf("绑定队列失败:%s", err)
		return
	}
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	r.channel.Qos(50, 0, true)
	msgList, err := r.channel.Consume(r.queueName, "", false, false, false, false, nil)
	if err != nil {
		globallogger.Log.Errorf("获取消费通道异常:%s", err)
	} else {
		for msg := range msgList {
			// 处理数据
			if msg.AppId == "iotwebserver" {
				response, _ := receiver.ConsumerByIotwebserver(msg)
				queueExchange := &QueueExchange{
					QuName: msg.ReplyTo,
					RtKey:  msg.ReplyTo,
					ExName: "ex_direct",
					ExType: "direct",
				}
				mq := New(queueExchange, "iotzigbeeurl", msg.MessageId)
				mq.RegisterProducer(&response)
				mq.StartProducer()
			} else {
				err = receiver.Consumer(msg)
			}
			if err != nil {
				err = msg.Ack(true)
				if err != nil {
					globallogger.Log.Errorf("确认消息未完成异常:%s", err)
				}
			} else {
				// 确认消息,必须为false
				err = msg.Ack(false)
				if err != nil {
					globallogger.Log.Errorf("确认消息完成异常:%s", err)
				}
			}
		}
	}
}
