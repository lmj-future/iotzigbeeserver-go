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
	var timerIDExist = time.NewTimer(time.Duration(interval+3) * time.Second)
	terminalStatusTimerIDExist.Store(devEUI, timerIDExist)
	go func() {
		<-timerIDExist.C
		timerIDExist.Stop()
		if value, ok := terminalStatusTimerIDExist.Load(devEUI); ok {
			if timerIDExist == value.(*time.Timer) {
				terminalStatusTimerIDExist.Delete(devEUI)
				var terminalTimerUpdateTime int64
				var err error
				if constant.Constant.MultipleInstances {
					terminalTimerUpdateTime, err = publicfunction.TerminalTimerRedisGet(devEUI)
				} else {
					terminalTimerUpdateTime, err = publicfunction.TerminalTimerFreeCacheGet(devEUI)
				}
				if err == nil {
					var terminalInfo *config.TerminalInfo
					if constant.Constant.UsePostgres {
						terminalInfo, _ = models.GetTerminalInfoByDevEUIPG(devEUI)
					} else {
						terminalInfo, _ = models.GetTerminalInfoByDevEUI(devEUI)
					}
					if terminalInfo != nil && terminalInfo.Interval != 0 {
						interval = uint16(terminalInfo.Interval)
					}
					if time.Now().UnixNano()-terminalTimerUpdateTime > int64(time.Duration(interval)*time.Second) {
						globallogger.Log.Infoln("devEUI:", devEUI, "procKeepAliveCheckExist server has not already recv data for",
							(time.Now().UnixNano()-terminalTimerUpdateTime)/int64(time.Second), "seconds. Then check terminal")
						publicfunction.SendTerminalCheckExistReq(publicfunction.GetJSONInfo(*terminalInfo), devEUI)
					}
				}
			}
		}
	}()
}

// ProcKeepAlive 处理终端保活报文
func ProcKeepAlive(tmnInfo config.TerminalInfo, interval uint16) {
	globallogger.Log.Infoln("devEUI :", tmnInfo.DevEUI, "ProcKeepAlive")
	var err error
	if constant.Constant.Iotware {
		iotsmartspace.PublishTelemetryUpIotware(tmnInfo, iotsmartspace.IotwarePropertyKeepAliveTime{KeepAliveTime: interval})
	}
	if constant.Constant.MultipleInstances {
		_, err = publicfunction.TerminalTimerRedisSet(tmnInfo.DevEUI)
	} else {
		_, err = publicfunction.TerminalTimerFreeCacheSet(tmnInfo.DevEUI, int(interval))
	}
	if err == nil {
		//3次未收到数据报文，通知离线
		var timerID = time.NewTimer(time.Duration(3*interval+3) * time.Second)
		terminalStatusTimerID.Store(tmnInfo.DevEUI, timerID)
		go func() {
			<-timerID.C
			timerID.Stop()
			if value, ok := terminalStatusTimerID.Load(tmnInfo.DevEUI); ok {
				if timerID == value.(*time.Timer) {
					terminalStatusTimerID.Delete(tmnInfo.DevEUI)
					var terminalTimerUpdateTime int64
					if constant.Constant.MultipleInstances {
						terminalTimerUpdateTime, err = publicfunction.TerminalTimerRedisGet(tmnInfo.DevEUI)
					} else {
						terminalTimerUpdateTime, err = publicfunction.TerminalTimerFreeCacheGet(tmnInfo.DevEUI)
					}
					if err == nil {
						var terminalInfo *config.TerminalInfo
						var err error
						if constant.Constant.UsePostgres {
							terminalInfo, err = models.GetTerminalInfoByDevEUIPG(tmnInfo.DevEUI)
						} else {
							terminalInfo, err = models.GetTerminalInfoByDevEUI(tmnInfo.DevEUI)
						}
						if terminalInfo != nil && terminalInfo.Interval != 0 {
							interval = uint16(terminalInfo.Interval)
						}
						if time.Now().UnixNano()-terminalTimerUpdateTime > int64(time.Duration(3*interval)*time.Second) {
							globallogger.Log.Infoln("devEUI:", tmnInfo.DevEUI, "ProcKeepAlive server has not already recv data for",
								(time.Now().UnixNano()-terminalTimerUpdateTime)/int64(time.Second), "seconds. Then terminal offline")
							publicfunction.TerminalOffline(tmnInfo.DevEUI)
							if err == nil && terminalInfo != nil {
								if constant.Constant.Iotware {
									if !terminalInfo.LeaveState {
										iotsmartspace.StateTerminalOfflineIotware(*terminalInfo)
									} else {
										iotsmartspace.StateTerminalLeaveIotware(*terminalInfo)
									}
								} else if constant.Constant.Iotedge {
									if !terminalInfo.LeaveState {
										iotsmartspace.StateTerminalOffline(tmnInfo.DevEUI)
									} else {
										iotsmartspace.StateTerminalLeave(tmnInfo.DevEUI)
									}
								}
							}
							if constant.Constant.MultipleInstances {
								publicfunction.DeleteRedisTerminalTimer(tmnInfo.DevEUI)
							} else {
								publicfunction.DeleteFreeCacheTerminalTimer(tmnInfo.DevEUI)
							}
						}
					}
				}
			}
		}()
		procKeepAliveCheckExist(tmnInfo.DevEUI, interval)
	}
}
