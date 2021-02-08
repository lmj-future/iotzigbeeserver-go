package keepalive

import (
	"sync"
	"time"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/interactmodule/iotsmartspace"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
)

var terminalStatusTimerID = &sync.Map{}
var terminalStatusTimerIDExist = &sync.Map{}

func procKeepAliveCheckExist(devEUI string, interval uint16) {
	//1次未收到数据报文，检查“存在否”
	var timerIDExist = time.NewTimer(time.Duration(interval+30) * time.Second)
	terminalStatusTimerIDExist.Store(devEUI, timerIDExist)
	go func() {
		select {
		case <-timerIDExist.C:
			if value, ok := terminalStatusTimerIDExist.Load(devEUI); ok {
				if timerIDExist == value.(*time.Timer) {
					terminalStatusTimerIDExist.Delete(devEUI)
					terminalTimerInfo, _ := publicfunction.GetTerminalTimerByDevEUI(devEUI)
					if terminalTimerInfo != nil {
						var dateNow = time.Now()
						var num = dateNow.UnixNano() - terminalTimerInfo.UpdateTime.UnixNano()
						var terminalInfo *config.TerminalInfo
						if constant.Constant.UsePostgres {
							terminalInfo, _ = models.GetTerminalInfoByDevEUIPG(devEUI)
						} else {
							terminalInfo, _ = models.GetTerminalInfoByDevEUI(devEUI)
						}
						if terminalInfo != nil && terminalInfo.Interval != 0 {
							interval = uint16(terminalInfo.Interval)
						}
						if num > int64(time.Duration(interval)*time.Second) {
							globallogger.Log.Infoln("devEUI: "+devEUI+" procKeepAliveCheckExist server has not already recv data for ",
								num/int64(time.Second), " seconds. Then check terminal")
							jsonInfo := publicfunction.GetJSONInfo(*terminalInfo)
							publicfunction.SendTerminalCheckExistReq(jsonInfo, devEUI)
						}
					}
				}
			}
		}
	}()
}

// ProcKeepAlive 处理终端保活报文
func ProcKeepAlive(devEUI string, interval uint16) {
	globallogger.Log.Infoln("devEUI : " + devEUI + " " + "ProcKeepAlive")
	err := publicfunction.TerminalTimerUpdateOrCreate(devEUI)
	if err == nil {
		if value, ok := terminalStatusTimerID.Load(devEUI); ok {
			value.(*time.Timer).Stop()
			globallogger.Log.Infoln("devEUI: "+devEUI+" ProcKeepAlive clearTimeout, time: ", interval)
		}
		//3次未收到数据报文，通知离线
		var timerID = time.NewTimer(time.Duration(3*interval+10) * time.Second)
		terminalStatusTimerID.Store(devEUI, timerID)
		go func() {
			select {
			case <-timerID.C:
				if value, ok := terminalStatusTimerID.Load(devEUI); ok {
					if timerID == value.(*time.Timer) {
						terminalStatusTimerID.Delete(devEUI)
						terminalTimerInfo, _ := publicfunction.GetTerminalTimerByDevEUI(devEUI)
						if terminalTimerInfo != nil {
							var dateNow = time.Now()
							var num = dateNow.UnixNano() - terminalTimerInfo.UpdateTime.UnixNano()
							var terminalInfo *config.TerminalInfo
							var err error
							if constant.Constant.UsePostgres {
								terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
							} else {
								terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
							}
							if terminalInfo != nil && terminalInfo.Interval != 0 {
								interval = uint16(terminalInfo.Interval)
							}
							if num > int64(time.Duration(3*interval)*time.Second) {
								globallogger.Log.Infoln("devEUI: "+devEUI+" ProcKeepAlive server has not already recv data for ",
									num/int64(time.Second), " seconds. Then terminal offline")
								publicfunction.TerminalOffline(devEUI)
								if err == nil && terminalInfo != nil {
									if constant.Constant.Iotware {
										if !terminalInfo.LeaveState {
											iotsmartspace.StateTerminalOfflineIotware(*terminalInfo)
										} else {
											iotsmartspace.StateTerminalLeaveIotware(*terminalInfo)
										}
									} else if constant.Constant.Iotedge {
										if !terminalInfo.LeaveState {
											iotsmartspace.StateTerminalOffline(devEUI)
										} else {
											iotsmartspace.StateTerminalLeave(devEUI)
										}
									}
								}
								publicfunction.DeleteTerminalTimer(devEUI)
							}
						}
					}
				}
			}
		}()
		procKeepAliveCheckExist(devEUI, interval)
	}
}
