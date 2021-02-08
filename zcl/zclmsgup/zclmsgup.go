package zclmsgup

import (
	"strconv"

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
	clusterID := afIncomingMessage.ClusterID
	// srcAddr := afIncomingMessage.SrcAddr
	srcEndpoint := afIncomingMessage.SrcEndpoint
	// dstEndpoint := afIncomingMessage.DstEndpoint
	// wasBroadcast := afIncomingMessage.WasBroadcast
	// linkQuality := afIncomingMessage.LinkQuality
	// securityUse := afIncomingMessage.SecurityUse
	// timestamp := afIncomingMessage.Timestamp
	// transSeqNumber := afIncomingMessage.TransSeqNumber
	// zclData := afIncomingMessage.Data

	z := zcl.New()
	var err error
	if zclAfIncomingMessage == nil {
		zclAfIncomingMessage, err = z.ToZclIncomingMessage(&afIncomingMessage)
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][ProcZclUpMsg] err: %v", terminalInfo.DevEUI, err)
	} else {
		// globallogger.Log.Infof("[devEUI: %v][ProcZclUpMsg] zclAfIncomingMessage: %+v", terminalInfo.DevEUI, zclAfIncomingMessage)
		var contentZclAfIncomingMessage *zcl.IncomingMessage = nil
		var dataTemp = ""
		if contentData != "" {
			contentData = contentData[10:]
			dataTemp = contentData
			var data []uint8
			for len(contentData) != 0 {
				value, _ := strconv.ParseUint(contentData[:2], 16, 8)
				data = append(data, uint8(value))
				contentData = contentData[2:]
			}
			afIncomingMessage.Data = data
			contentZclAfIncomingMessage, err = z.ToZclIncomingMessage(&afIncomingMessage)
			if err != nil {
				globallogger.Log.Errorf("[devEUI: %v][ProcZclUpMsg] contentZclAfIncomingMessage err: %v", terminalInfo.DevEUI, err)
			}
		}
		if z.ClusterLibrary().Clusters()[cluster.ID(clusterID)] != nil {
			clusterIDByStr := z.ClusterLibrary().Clusters()[cluster.ID(clusterID)].Name
			globallogger.Log.Infof("[devEUI: %v][ProcZclUpMsg] clusterID: %v", terminalInfo.DevEUI, clusterIDByStr)
			if contentZclAfIncomingMessage != nil && contentZclAfIncomingMessage.Data != nil {
				clusterproc.ClusterProc(terminalInfo, srcEndpoint, zclAfIncomingMessage.Data, clusterIDByStr, msgID, contentZclAfIncomingMessage.Data, dataTemp)
			} else {
				clusterproc.ClusterProc(terminalInfo, srcEndpoint, zclAfIncomingMessage.Data, clusterIDByStr, msgID, nil, "")
			}
		} else {
			globallogger.Log.Infof("[devEUI: %v][ProcZclUpMsg] invalid clusterID: %v", terminalInfo.DevEUI, clusterID)
		}
	}
}
