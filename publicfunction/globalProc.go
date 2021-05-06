package publicfunction

import (
	"strconv"
	"sync"
	"time"

	"github.com/coocood/freecache"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmemorycache"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalrediscache"
)

var zigbeeServerMsgCheckTimerID = &sync.Map{}
var terminalInfoListCache = &sync.Map{}

var zigbeeServerMsgCheckTimerFreeCache = freecache.NewCache(10 * 1024 * 1024)
var zigbeeServerKeepAliveTimerFreeCache = freecache.NewCache(1024 * 1024)
var zigbeeServerTerminalTimerFreeCache = freecache.NewCache(1024 * 1024)

// LoadTerminalInfoListCache LoadTerminalInfoListCache
func LoadTerminalInfoListCache(key string) (interface{}, bool) {
	if v, ok := terminalInfoListCache.Load(key); ok {
		return v, ok
	}
	return nil, false
}

// StoreTerminalInfoListCache LoadTerminalInfoListCache
func StoreTerminalInfoListCache(key string, terminalInfo *config.TerminalInfo) {
	terminalInfoListCache.Store(key, terminalInfo)
}

// DeleteTerminalInfoListCache DeleteTerminalInfoListCache
func DeleteTerminalInfoListCache(key string) {
	if v, ok := terminalInfoListCache.Load(key); ok {
		terminalInfo := v.(*config.TerminalInfo)
		if terminalInfo != nil {
			terminalInfo = nil
		}
	}
	terminalInfoListCache.Delete(key)
}

/******************************************************************************************/
/* keepAliveTimerModel Func Start */
/******************************************************************************************/
func KeepAliveTimerRedisGet(APMac string) (int64, error) {
	var updateTime int64
	var err error
	value, err := globalrediscache.RedisCache.GetRedis(APMac + "ZigbeeServerKeepAliveTimer")
	if err == nil {
		updateTime, err = strconv.ParseInt(value, 10, 64)
	}
	return updateTime, err
}
func KeepAliveTimerRedisSet(APMac string) (int64, error) {
	updateTime := time.Now().UnixNano()
	_, err := globalrediscache.RedisCache.SetRedis(APMac+"ZigbeeServerKeepAliveTimer", 0, []byte(strconv.FormatInt(updateTime, 10)))
	globalrediscache.RedisCache.SaddRedis(constant.Constant.REDIS.ZigbeeRedisKeepAliveTimerSets, APMac+"ZigbeeServerKeepAliveTimer")
	return updateTime, err
}
func DeleteRedisKeepAliveTimer(APMac string) {
	globalrediscache.RedisCache.DeleteRedis(APMac + "ZigbeeServerKeepAliveTimer")
	globalrediscache.RedisCache.SremRedis(constant.Constant.REDIS.ZigbeeRedisKeepAliveTimerSets, APMac+"ZigbeeServerKeepAliveTimer")
}

func KeepAliveTimerFreeCacheGet(APMac string) (int64, error) {
	var updateTime int64
	var err error
	value, err := zigbeeServerKeepAliveTimerFreeCache.Get([]byte(APMac))
	if err == nil {
		updateTime, err = strconv.ParseInt(string(value), 10, 64)
	}
	return updateTime, err
}
func KeepAliveTimerFreeCacheSet(APMac string, keepAliveInterval int) (int64, error) {
	updateTime := time.Now().UnixNano()
	err := zigbeeServerKeepAliveTimerFreeCache.Set([]byte(APMac), []byte(strconv.FormatInt(updateTime, 10)), 4*keepAliveInterval)
	return updateTime, err
}
func DeleteFreeCacheKeepAliveTimer(APMac string) {
	zigbeeServerKeepAliveTimerFreeCache.Del([]byte(APMac))
}

/****************************************************************************************/
/* keepAliveTimerModel Func End */
/****************************************************************************************/

/*****************************************************************************************/
/* msgCheckTimerModel Func Start */
/*****************************************************************************************/
func msgCheckTimerRedisGet(MsgKey string) (int64, error) {
	var updateTime int64
	var err error
	value, err := globalrediscache.RedisCache.GetRedis(MsgKey + "Timer")
	if err == nil {
		updateTime, err = strconv.ParseInt(value, 10, 64)
	}
	return updateTime, err
}
func msgCheckTimerRedisSet(MsgKey string) (int64, error) {
	updateTime := time.Now().UnixNano()
	_, err := globalrediscache.RedisCache.SetRedis(MsgKey+"Timer", 0, []byte(strconv.FormatInt(updateTime, 10)))
	return updateTime, err
}
func deleteRedisMsgCheckTimer(MsgKey string) {
	globalrediscache.RedisCache.DeleteRedis(MsgKey + "Timer")
}

func msgCheckTimerFreeCacheGet(MsgKey string) (int64, error) {
	var updateTime int64
	var err error
	value, err := zigbeeServerMsgCheckTimerFreeCache.Get([]byte(MsgKey))
	if err == nil {
		updateTime, err = strconv.ParseInt(string(value), 10, 64)
	}
	return updateTime, err
}
func msgCheckTimerFreeCacheSet(MsgKey string) (int64, error) {
	updateTime := time.Now().UnixNano()
	err := zigbeeServerMsgCheckTimerFreeCache.Set([]byte(MsgKey), []byte(strconv.FormatInt(updateTime, 10)), constant.Constant.TIMER.ZigbeeTimerMsgDelete+5)
	return updateTime, err
}
func deleteFreeCacheMsgCheckTimer(MsgKey string) {
	zigbeeServerMsgCheckTimerFreeCache.Del([]byte(MsgKey))
}

/***************************************************************************************/
/* msgCheckTimerModel Func End */
/***************************************************************************************/

/******************************************************************************************/
/* terminalTimerModel Func Start */
/******************************************************************************************/
func TerminalTimerRedisGet(devEUI string) (int64, error) {
	var updateTime int64
	var err error
	value, err := globalrediscache.RedisCache.GetRedis(devEUI + "ZigbeeServerTerminalTimer")
	if err == nil {
		updateTime, err = strconv.ParseInt(value, 10, 64)
	}
	return updateTime, err
}
func TerminalTimerRedisSet(devEUI string) (int64, error) {
	updateTime := time.Now().UnixNano()
	_, err := globalrediscache.RedisCache.SetRedis(devEUI+"ZigbeeServerTerminalTimer", 0, []byte(strconv.FormatInt(updateTime, 10)))
	globalrediscache.RedisCache.SaddRedis(constant.Constant.REDIS.ZigbeeRedisTerminalTimerSets, devEUI+"ZigbeeServerTerminalTimer")
	return updateTime, err
}
func DeleteRedisTerminalTimer(devEUI string) {
	globalrediscache.RedisCache.DeleteRedis(devEUI + "ZigbeeServerTerminalTimer")
	globalrediscache.RedisCache.SremRedis(constant.Constant.REDIS.ZigbeeRedisTerminalTimerSets, devEUI+"ZigbeeServerTerminalTimer")
}

func TerminalTimerFreeCacheGet(devEUI string) (int64, error) {
	var updateTime int64
	var err error
	value, err := zigbeeServerTerminalTimerFreeCache.Get([]byte(devEUI))
	if err == nil {
		updateTime, err = strconv.ParseInt(string(value), 10, 64)
	}
	return updateTime, err
}
func TerminalTimerFreeCacheSet(devEUI string, interval int) (int64, error) {
	updateTime := time.Now().UnixNano()
	err := zigbeeServerTerminalTimerFreeCache.Set([]byte(devEUI), []byte(strconv.FormatInt(updateTime, 10)), 4*interval)
	return updateTime, err
}
func DeleteFreeCacheTerminalTimer(devEUI string) {
	zigbeeServerTerminalTimerFreeCache.Del([]byte(devEUI))
}

/****************************************************************************************/
/* terminalTimerModel Func End */
/****************************************************************************************/

func deleteRedisMsg(devEUI string, key string, t int) {
	zigbeeServerMsgCheckTimerID.Delete(key)
	updateTime, err := msgCheckTimerRedisGet(key)
	if err == nil {
		if time.Now().UnixNano()-updateTime > int64(time.Duration(t)*time.Second) {
			_, err := globalrediscache.RedisCache.DeleteRedis(key)
			if err != nil {
				globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "DeleteMsg delete redis error :", err)
			}
			_, err = globalrediscache.RedisCache.SremRedis(constant.Constant.REDIS.ZigbeeRedisDownMsgSets, key)
			if err != nil {
				globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "DeleteMsg srem redis error :", err)
			}
			deleteRedisMsgCheckTimer(key)
		}
	}
}

func deleteFreeCacheMsg(devEUI string, key string, t int) {
	zigbeeServerMsgCheckTimerID.Delete(key)
	updateTime, err := msgCheckTimerFreeCacheGet(key)
	if err == nil {
		if time.Now().UnixNano()-updateTime > int64(time.Duration(t)*time.Second) {
			_, err := globalmemorycache.MemoryCache.DeleteMemory(key)
			if err != nil {
				globallogger.Log.Errorln("devEUI :", devEUI, "key:", key, "DeleteMsg delete memory error :", err)
			}
			deleteFreeCacheMsgCheckTimer(key)
		}
	}
}

func timeoutCheckAndDeleteRedisMsg(key string, devEUI string, timerID *time.Timer, isRebuild bool) {
	length, redisArray, _ := GetRedisLengthAndRangeRedis(key, devEUI)
	if len(redisArray) != 0 {
		var tempLen = 0
		for _, item := range redisArray {
			if item == constant.Constant.REDIS.ZigbeeRedisRemoveRedis {
				tempLen++
			}
		}
		tempTimerID, ok := zigbeeServerMsgCheckTimerID.Load(key)
		if (tempLen == length) && (ok && timerID == tempTimerID.(*time.Timer)) {
			deleteRedisMsg(devEUI, key, constant.Constant.TIMER.ZigbeeTimerMsgDelete)
		} else if (tempLen < length) && (ok && timerID == tempTimerID.(*time.Timer)) {
			if isRebuild {
				deleteRedisMsg(devEUI, key, constant.Constant.TIMER.ZigbeeTimerServerRebuild)
			} else {
				zigbeeServerMsgCheckTimerID.Delete(key)
				updateTime, err := msgCheckTimerRedisGet(key)
				if err == nil {
					if time.Now().UnixNano()-updateTime > int64(time.Duration(constant.Constant.TIMER.ZigbeeTimerMsgDelete)*time.Second) {
						timerID = time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerMsgDelete) * time.Second)
						zigbeeServerMsgCheckTimerID.Store(key, timerID)
						go func() {
							<-timerID.C
							timerID.Stop()
							timeoutCheckAndDeleteRedisMsg(key, devEUI, timerID, false)
						}()
					}
				}
			}
		}
	}
}

func timeoutCheckAndDeleteFreeCacheMsg(key string, devEUI string, timerID *time.Timer) {
	length, redisArray, _ := GetRedisLengthAndRangeRedis(key, devEUI)
	if len(redisArray) != 0 {
		var tempLen = 0
		for _, item := range redisArray {
			if item == constant.Constant.REDIS.ZigbeeRedisRemoveRedis {
				tempLen++
			}
		}
		tempTimerID, ok := zigbeeServerMsgCheckTimerID.Load(key)
		if (tempLen == length) && (ok && timerID == tempTimerID.(*time.Timer)) {
			deleteFreeCacheMsg(devEUI, key, constant.Constant.TIMER.ZigbeeTimerMsgDelete)
		} else if (tempLen < length) && (ok && timerID == tempTimerID.(*time.Timer)) {
			zigbeeServerMsgCheckTimerID.Delete(key)
			updateTime, err := msgCheckTimerFreeCacheGet(key)
			if err == nil {
				if time.Now().UnixNano()-updateTime > int64(time.Duration(constant.Constant.TIMER.ZigbeeTimerMsgDelete)*time.Second) {
					timerID = time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerMsgDelete) * time.Second)
					zigbeeServerMsgCheckTimerID.Store(key, timerID)
					go func() {
						<-timerID.C
						timerID.Stop()
						timeoutCheckAndDeleteFreeCacheMsg(key, devEUI, timerID)
					}()
				}
			}
		}
	}
}

// CheckTimeoutAndDeleteRedisMsg CheckTimeoutAndDeleteRedisMsg
func CheckTimeoutAndDeleteRedisMsg(key string, devEUI string, isRebuild bool) {
	if isRebuild {
		// if value, ok := zigbeeServerMsgCheckTimerID.Load(key); ok {
		// 	value.(*time.Timer).Stop()
		// }
		var timerID = time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerAfterTwo) * time.Second)
		zigbeeServerMsgCheckTimerID.Store(key, timerID)
		go func() {
			<-timerID.C
			timerID.Stop()
			tempTimerID, ok := zigbeeServerMsgCheckTimerID.Load(key)
			if ok && timerID == tempTimerID.(*time.Timer) {
				timeoutCheckAndDeleteRedisMsg(key, devEUI, timerID, isRebuild)
			}
		}()
	} else {
		_, err := msgCheckTimerRedisSet(key)
		if err == nil {
			// if value, ok := zigbeeServerMsgCheckTimerID.Load(key); ok {
			// 	value.(*time.Timer).Stop()
			// }
			var timerID = time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerMsgDelete) * time.Second)
			zigbeeServerMsgCheckTimerID.Store(key, timerID)
			go func() {
				<-timerID.C
				timerID.Stop()
				tempTimerID, ok := zigbeeServerMsgCheckTimerID.Load(key)
				if ok && timerID == tempTimerID.(*time.Timer) {
					timeoutCheckAndDeleteRedisMsg(key, devEUI, timerID, isRebuild)
				}
			}()
		}
	}
}

// CheckTimeoutAndDeleteFreeCacheMsg CheckTimeoutAndDeleteFreeCacheMsg
func CheckTimeoutAndDeleteFreeCacheMsg(key string, devEUI string) {
	_, err := msgCheckTimerFreeCacheSet(key)
	if err == nil {
		// if value, ok := zigbeeServerMsgCheckTimerID.Load(key); ok {
		// 	value.(*time.Timer).Stop()
		// }
		var timerID = time.NewTimer(time.Duration(constant.Constant.TIMER.ZigbeeTimerMsgDelete+1) * time.Second)
		zigbeeServerMsgCheckTimerID.Store(key, timerID)
		go func() {
			<-timerID.C
			timerID.Stop()
			tempTimerID, ok := zigbeeServerMsgCheckTimerID.Load(key)
			if ok && timerID == tempTimerID.(*time.Timer) {
				timeoutCheckAndDeleteFreeCacheMsg(key, devEUI, timerID)
			}
		}()
	}
}
