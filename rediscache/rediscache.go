package rediscache

import (
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalredisclient"
)

//RedisCache RedisCache
type RedisCache struct{}

//Interface Interface
type Interface interface {
	InsertRedis(key string, data []byte) (int, error)
	UpdateRedis(key string, data []byte) (string, error)
	DeleteRedis(key string) (int, error)
	DeleteRedisOne(key string, value string) (int, error)
	GetRedis(key string) (string, error)
	GetRedisByIndex(key string, index int) (string, error)
	GetRedisEnd(key string) (interface{}, error)
	GetRedisLength(key string) (int, error)
	PopRedis(key string) (string, error)
	RangeRedis(key string, start int, stop int) ([]string, error)
	SetRedis(key string, index int, value []byte) (string, error)
	RemoveRedis(key string, count int, value string) (int, error)
	SaddRedis(key string, member string) (int, error)
	SremRedis(key string, member string) (int, error)
	FindAllRedisKeys(key string) ([]string, error)
	SetRedisSet(key string, value string) (string, error)
	GetRedisGet(key string) (string, error)
}

//InsertRedis InsertRedis
func (RedisCache) InsertRedis(key string, data []byte) (int, error) {
	//往队列尾部添加数据
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Rpush(key, data)
	return reply, err
}

//UpdateRedis UpdateRedis
func (RedisCache) UpdateRedis(key string, data []byte) (string, error) {
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Lpop(key)
	if err == nil {
		globalredisclient.MyZigbeeServerRedisClient.Lpush(key, data)
	}
	return reply, err
}

//DeleteRedis DeleteRedis
func (RedisCache) DeleteRedis(key string) (int, error) {
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Del(key)
	return reply, err
}

//DeleteRedisOne DeleteRedisOne
func (RedisCache) DeleteRedisOne(key string, value string) (int, error) {
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Lrem(key, 0, value)
	return reply, err
}

//GetRedis GetRedis
func (RedisCache) GetRedis(key string) (string, error) {
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Lindex(key, 0)
	return reply, err
}

//GetRedisByIndex GetRedisByIndex
func (RedisCache) GetRedisByIndex(key string, index int) (string, error) {
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Lindex(key, index)
	return reply, err
}

//GetRedisEnd GetRedisEnd
func (RedisCache) GetRedisEnd(key string) (interface{}, error) {
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Llen(key)
	if err == nil {
		var index = reply
		if reply > 0 {
			index = reply - 1
		}
		reply, err := globalredisclient.MyZigbeeServerRedisClient.Lindex(key, index)
		return reply, err
	}
	return reply, err
}

//GetRedisLength GetRedisLength
func (RedisCache) GetRedisLength(key string) (int, error) {
	len, err := globalredisclient.MyZigbeeServerRedisClient.Llen(key)
	return len, err
}

//PopRedis PopRedis
func (RedisCache) PopRedis(key string) (string, error) {
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Lpop(key)
	return reply, err
}

//RangeRedis RangeRedis
func (RedisCache) RangeRedis(key string, start int, stop int) ([]string, error) {
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Lrange(key, start, stop)
	return reply, err
}

//SetRedis SetRedis
func (RedisCache) SetRedis(key string, index int, value []byte) (string, error) {
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Lset(key, index, value)
	return reply, err
}

//RemoveRedis RemoveRedis
func (RedisCache) RemoveRedis(key string, count int, value string) (int, error) {
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Lrem(key, count, value)
	return reply, err
}

//SaddRedis SaddRedis
func (RedisCache) SaddRedis(key string, member string) (int, error) {
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Sadd(key, member)
	return reply, err
}

//SremRedis SremRedis
func (RedisCache) SremRedis(key string, member string) (int, error) {
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Srem(key, member)
	return reply, err
}

//FindAllRedisKeys FindAllRedisKeys
func (RedisCache) FindAllRedisKeys(key string) ([]string, error) {
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Smembers(key)
	return reply, err
}

//SetRedisSet SetRedisSet
func (RedisCache) SetRedisSet(key string, value string) (string, error) {
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Set(key, value)
	return reply, err
}

//GetRedisGet GetRedisGet
func (RedisCache) GetRedisGet(key string) (string, error) {
	reply, err := globalredisclient.MyZigbeeServerRedisClient.Get(key)
	return reply, err
}
