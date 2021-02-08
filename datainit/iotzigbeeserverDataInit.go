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
			select {
			case <-time.After(time.Duration(5) * time.Second):
				go func() {
					redisDataList, err := globalrediscache.RedisCache.FindAllRedisKeys(constant.Constant.REDIS.ZigbeeRedisDownMsgSets)
					if err != nil {
						globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit findAllRedisKeys error: ", err)
					} else {
						globallogger.Log.Infoln("|||||| iotzigbeeserverDataInit findAllRedisKeys sucess:", redisDataList)
						if redisDataList != nil {
							for _, item := range redisDataList {
								publicfunction.CheckTimeoutAndDeleteRedisMsg(item, "POD IS REBUILD", true)
							}
						}
					}
				}()
				go func() {
					redisDataList, err := globalrediscache.RedisCache.FindAllRedisKeys(constant.Constant.REDIS.ZigbeeRedisDuplicationFlagSets)
					if err != nil {
						globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit findAllRedisKeys error: ", err)
					} else {
						globallogger.Log.Infoln("|||||| iotzigbeeserverDataInit findAllRedisKeys sucess:", redisDataList)
						if redisDataList != nil {
							for _, item := range redisDataList {
								redisData, err := globalrediscache.RedisCache.DeleteRedis(item)
								if err != nil {
									globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit key: "+item+" delete redis error : ", err)
								} else {
									globallogger.Log.Infoln("|||||| iotzigbeeserverDataInit key: "+item+" delete redis success: ", redisData)
								}
								redisData, err = globalrediscache.RedisCache.SremRedis(constant.Constant.REDIS.ZigbeeRedisDuplicationFlagSets, item)
								if err != nil {
									globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit key: "+item+" srem redis error : ", err)
								} else {
									globallogger.Log.Infoln("|||||| iotzigbeeserverDataInit key: "+item+" srem redis success: ", redisData)
								}
							}
						}
					}
				}()
				go func() {
					redisDataList, err := globalrediscache.RedisCache.FindAllRedisKeys(constant.Constant.REDIS.ZigbeeRedisPermitJoinContentSets)
					if err != nil {
						globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit findAllRedisKeys error: ", err)
					} else {
						globallogger.Log.Infoln("|||||| iotzigbeeserverDataInit findAllRedisKeys sucess:", redisDataList)
						if redisDataList != nil {
							for _, item := range redisDataList {
								redisData, err := globalrediscache.RedisCache.DeleteRedis(item)
								if err != nil {
									globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit key: "+item+" delete redis error : ", err)
								} else {
									globallogger.Log.Infoln("|||||| iotzigbeeserverDataInit key: "+item+" delete redis success: ", redisData)
								}
								redisData, err = globalrediscache.RedisCache.SremRedis(constant.Constant.REDIS.ZigbeeRedisDuplicationFlagSets, item)
								if err != nil {
									globallogger.Log.Errorln("|||||| iotzigbeeserverDataInit key: "+item+" srem redis error : ", err)
								} else {
									globallogger.Log.Infoln("|||||| iotzigbeeserverDataInit key: "+item+" srem redis success: ", redisData)
								}
							}
						}
					}
				}()
			}
		}()
	} else {
		globalmemorycache.MemoryCache = memorycache.MemoryCache{}
	}
}
