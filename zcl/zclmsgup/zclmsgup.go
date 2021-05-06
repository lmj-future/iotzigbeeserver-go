package zclmsgup

import (
	"strconv"
	"strings"

	"github.com/dyrkin/znp-go"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/zcl/clusterproc"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
)

//ProcZclUpMsg ZCL上行数据报文处理
func ProcZclUpMsg(terminalInfo config.TerminalInfo, afIncomingMessage znp.AfIncomingMessage, msgID interface{},
	contentData string, zclAfIncomingMessage *zcl.IncomingMessage) {
	globallogger.Log.Infof("[devEUI: %v][ProcZclUpMsg] afIncomingMessage: %+v", terminalInfo.DevEUI, afIncomingMessage)
	// groupID := afIncomingMessage.GroupID
	// srcAddr := afIncomingMessage.SrcAddr
	// dstEndpoint := afIncomingMessage.DstEndpoint
	// wasBroadcast := afIncomingMessage.WasBroadcast
	// linkQuality := afIncomingMessage.LinkQuality
	// securityUse := afIncomingMessage.SecurityUse
	// timestamp := afIncomingMessage.Timestamp
	// transSeqNumber := afIncomingMessage.TransSeqNumber
	// zclData := afIncomingMessage.Data

	var err error
	if zclAfIncomingMessage == nil {
		zclAfIncomingMessage, err = zcl.ZCLObj().ToZclIncomingMessage(&afIncomingMessage)
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][ProcZclUpMsg] err: %v", terminalInfo.DevEUI, err)
	} else {
		var contentZclAfIncomingMessage *zcl.IncomingMessage = nil
		var dataTemp = ""
		if contentData != "" {
			contentData = strings.Repeat(contentData[10:], 1)
			dataTemp = contentData
			var data []uint8
			for len(contentData) != 0 {
				value, _ := strconv.ParseUint(strings.Repeat(contentData[:2], 1), 16, 8)
				data = append(data, uint8(value))
				contentData = strings.Repeat(contentData[2:], 1)
			}
			afIncomingMessage.Data = data
			contentZclAfIncomingMessage, err = zcl.ZCLObj().ToZclIncomingMessage(&afIncomingMessage)
			if err != nil {
				globallogger.Log.Errorf("[devEUI: %v][ProcZclUpMsg] contentZclAfIncomingMessage err: %v", terminalInfo.DevEUI, err)
			}
		}
		if zcl.ZCLObj().ClusterLibrary().Clusters()[cluster.ID(afIncomingMessage.ClusterID)] != nil {
			clusterIDByStr := zcl.ZCLObj().ClusterLibrary().Clusters()[cluster.ID(afIncomingMessage.ClusterID)].Name
			globallogger.Log.Infof("[devEUI: %v][ProcZclUpMsg] clusterID: %v", terminalInfo.DevEUI, clusterIDByStr)
			if contentZclAfIncomingMessage != nil && contentZclAfIncomingMessage.Data != nil {
				clusterproc.ClusterProc(terminalInfo, afIncomingMessage.SrcEndpoint, zclAfIncomingMessage.Data, clusterIDByStr, msgID, contentZclAfIncomingMessage.Data, dataTemp)
			} else {
				clusterproc.ClusterProc(terminalInfo, afIncomingMessage.SrcEndpoint, zclAfIncomingMessage.Data, clusterIDByStr, msgID, nil, "")
			}
		} else {
			globallogger.Log.Infof("[devEUI: %v][ProcZclUpMsg] invalid clusterID: %v", terminalInfo.DevEUI, afIncomingMessage.ClusterID)
		}
	}
}
