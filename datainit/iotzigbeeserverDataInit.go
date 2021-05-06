package datainit

import (
	"time"

	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmemorycache"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalrediscache"
	"github.com/h3c/iotzigbeeserver-go/memorycache"
	"github.com/h3c/iotzigbeeserver-go/publicfunction"
	"github.com/h3c/iotzigbeeserver-go/rediscache"
)

//IotzigbeeserverDataInit IotzigbeeserverDataInit
func IotzigbeeserverDataInit() {
	if constant.Constant.MultipleInstances {
		globalrediscache.RedisCache = rediscache.RedisCache{}
		go func() {
			timer := time.NewTimer(5 * time.Second)
			<-timer.C
			timer.Stop()
			go func() {
				redisDataList, err := globalrediscache.RedisCache.FindAllRedisKeys(constant.Constant.REDIS.ZigbeeRedisDownMsgSets)
				if err != nil {
					globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit findAllRedisKeys error:", err)
				} else {
					if len(redisDataList) > 0 {
						for _, item := range redisDataList {
							publicfunction.CheckTimeoutAndDeleteRedisMsg(item, "POD IS REBUILD", true)
						}
					}
				}
			}()
			go func() {
				redisDataList, err := globalrediscache.RedisCache.FindAllRedisKeys(constant.Constant.REDIS.ZigbeeRedisDuplicationFlagSets)
				if err != nil {
					globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit findAllRedisKeys error:", err)
				} else {
					if len(redisDataList) > 0 {
						for _, item := range redisDataList {
							_, err = globalrediscache.RedisCache.DeleteRedis(item)
							if err != nil {
								globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit key:", item, "delete redis error :", err)
							}
							_, err = globalrediscache.RedisCache.SremRedis(constant.Constant.REDIS.ZigbeeRedisDuplicationFlagSets, item)
							if err != nil {
								globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit key:", item, "srem redis error :", err)
							}
						}
					}
				}
			}()
			go func() {
				redisDataList, err := globalrediscache.RedisCache.FindAllRedisKeys(constant.Constant.REDIS.ZigbeeRedisPermitJoinContentSets)
				if err != nil {
					globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit findAllRedisKeys error:", err)
				} else {
					if len(redisDataList) > 0 {
						for _, item := range redisDataList {
							_, err = globalrediscache.RedisCache.DeleteRedis(item)
							if err != nil {
								globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit key:", item, "delete redis error :", err)
							}
							_, err = globalrediscache.RedisCache.SremRedis(constant.Constant.REDIS.ZigbeeRedisPermitJoinContentSets, item)
							if err != nil {
								globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit key:", item, "srem redis error :", err)
							}
						}
					}
				}
			}()
			go func() {
				redisDataList, err := globalrediscache.RedisCache.FindAllRedisKeys(constant.Constant.REDIS.ZigbeeRedisKeepAliveTimerSets)
				if err != nil {
					globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit findAllRedisKeys error:", err)
				} else {
					if len(redisDataList) > 0 {
						for _, item := range redisDataList {
							_, err = globalrediscache.RedisCache.DeleteRedis(item)
							if err != nil {
								globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit key:", item, "delete redis error :", err)
							}
							_, err = globalrediscache.RedisCache.SremRedis(constant.Constant.REDIS.ZigbeeRedisKeepAliveTimerSets, item)
							if err != nil {
								globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit key:", item, "srem redis error :", err)
							}
						}
					}
				}
			}()
			go func() {
				redisDataList, err := globalrediscache.RedisCache.FindAllRedisKeys(constant.Constant.REDIS.ZigbeeRedisTerminalTimerSets)
				if err != nil {
					globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit findAllRedisKeys error:", err)
				} else {
					if len(redisDataList) > 0 {
						for _, item := range redisDataList {
							_, err = globalrediscache.RedisCache.DeleteRedis(item)
							if err != nil {
								globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit key:", item, "delete redis error :", err)
							}
							_, err = globalrediscache.RedisCache.SremRedis(constant.Constant.REDIS.ZigbeeRedisTerminalTimerSets, item)
							if err != nil {
								globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit key:", item, "srem redis error :", err)
							}
						}
					}
				}
			}()
		}()
	} else {
		globalmemorycache.MemoryCache = memorycache.MemoryCache{}
	}
}
