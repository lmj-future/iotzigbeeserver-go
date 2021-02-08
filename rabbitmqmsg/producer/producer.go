package producer

import (
	"encoding/json"
	"time"

	"github.com/h3c/iotzigbeeserver-go/db/rabbitmq"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
)

// TerminalState TerminalState
type TerminalState struct {
	TmnDevSN     string    `json:"tmnDevSN"`
	Time         time.Time `json:"time"`
	LinkType     string    `json:"linkType"`
	IndexType    string    `json:"indexType"`
	Index        string    `json:"index"`
	TmnOIDIndex  string    `json:"tmnOIDIndex"`
	ExtType      string    `json:"extType"`
	IsSendConfig bool      `json:"isSendConfig"`
}

// MsgContent 实现发送者
func (t *TerminalState) MsgContent() string {
	tJSONByte, _ := json.Marshal(t)
	globallogger.Log.Warnf("[RabbitMQ] Send terminal state msg: %+v\n", string(tJSONByte))
	return string(tJSONByte)
}

// SendTerminalStateMsg SendTerminalStateMsg
func SendTerminalStateMsg(t *TerminalState, queueName string) {
	queueExchange := &rabbitmq.QueueExchange{
		QuName: queueName,
		RtKey:  queueName,
		ExName: "ex_direct",
		ExType: "direct",
	}
	mq := rabbitmq.New(queueExchange, "iotzigbeetmnmgr", "")
	mq.RegisterProducer(t)
	mq.StartProducer()
}
