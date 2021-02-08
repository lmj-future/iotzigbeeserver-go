package zclmain

import (
	"github.com/dyrkin/znp-go"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zclmsgdown"
	"github.com/h3c/iotzigbeeserver-go/zcl/zclmsgup"
)

//ZclMain 定义zcl模块主入口函数
func ZclMain(msgType string, terminalInfo config.TerminalInfo, data interface{}, msgID interface{},
	contentData string, zclAfIncomingMessage *zcl.IncomingMessage) {
	defer func() {
		err := recover()
		if err != nil {
			globallogger.Log.Errorln("ZclMain err : ", err)
		}
	}()
	// globallogger.Log.Infof("[devEUI: %v[ZclMain]] data: %v", terminalInfo.DevEUI, data)
	switch msgType {
	case globalmsgtype.MsgType.UPMsg.ZigbeeDataUpEvent:
		zclmsgup.ProcZclUpMsg(terminalInfo, data.(znp.AfIncomingMessage), msgID, contentData, zclAfIncomingMessage)
	case globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent:
		zclmsgdown.ProcZclDownMsg(data.(common.ZclDownMsg))
	default:
		globallogger.Log.Warnf("[devEUI: %v][ZclMain] invalid msgType: %v", terminalInfo.DevEUI, msgType)
	}
}
