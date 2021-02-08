package redis

import "github.com/gomodule/redigo/redis"

//CloseSession CloseSession
func (rc *ClusterClient) CloseSession() error {
	crdb := rc.ClusterClient
	return crdb.Close()
}

//Rpush Rpush
func (rc *ClusterClient) Rpush(key string, data []byte) (int, error) {
	crdb := rc.ClusterClient
	//defer crdb.Close()
	reply, err := crdb.RPush(key, data).Result()
	if err == redis.ErrNil {
		err = nil
	}
	return int(reply), err
}

//Lpop Lpop
func (rc *ClusterClient) Lpop(key string) (string, error) {
	crdb := rc.ClusterClient
	//defer crdb.Close()
	reply, err := crdb.LPop(key).Result()
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Lpush Lpush
func (rc *ClusterClient) Lpush(key string, data []byte) (int, error) {
	crdb := rc.ClusterClient
	//defer crdb.Close()
	reply, err := crdb.LPush(key, data).Result()
	if err == redis.ErrNil {
		err = nil
	}
	return int(reply), err
}

//Del Del
func (rc *ClusterClient) Del(key string) (int, error) {
	crdb := rc.ClusterClient
	//defer crdb.Close()
	reply, err := crdb.Del(key).Result()
	if err == redis.ErrNil {
		err = nil
	}
	return int(reply), err
}

//Lrem Lrem
func (rc *ClusterClient) Lrem(key string, count int, value string) (int, error) {
	crdb := rc.ClusterClient
	//defer crdb.Close()
	reply, err := crdb.LRem(key, int64(count), value).Result()
	if err == redis.ErrNil {
		err = nil
	}
	return int(reply), err
}

//Lindex Lindex
func (rc *ClusterClient) Lindex(key string, start int) (string, error) {
	crdb := rc.ClusterClient
	defer crdb.Close()
	reply, err := crdb.LIndex(key, int64(start)).Result()
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Llen Llen
func (rc *ClusterClient) Llen(key string) (int, error) {
	crdb := rc.ClusterClient
	//defer crdb.Close()
	reply, err := crdb.LLen(key).Result()
	if err == redis.ErrNil {
		err = nil
	}
	return int(reply), err
}

//Lrange Lrange
func (rc *ClusterClient) Lrange(key string, start int, stop int) ([]string, error) {
	crdb := rc.ClusterClient
	//defer crdb.Close()
	reply, err := crdb.LRange(key, int64(start), int64(stop)).Result()
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Lset Lset
func (rc *ClusterClient) Lset(key string, index int, value []byte) (string, error) {
	crdb := rc.ClusterClient
	//defer crdb.Close()
	reply, err := crdb.LSet(key, int64(index), value).Result()
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Sadd Sadd
func (rc *ClusterClient) Sadd(key string, member string) (int, error) {
	crdb := rc.ClusterClient
	//defer crdb.Close()
	reply, err := crdb.SAdd(key, member).Result()
	if err == redis.ErrNil {
		err = nil
	}
	return int(reply), err
}

//Srem Srem
func (rc *ClusterClient) Srem(key string, member string) (int, error) {
	crdb := rc.ClusterClient
	//defer crdb.Close()
	reply, err := crdb.SRem(key, member).Result()
	if err == redis.ErrNil {
		err = nil
	}
	return int(reply), err
}

//Smembers Smembers
func (rc *ClusterClient) Smembers(key string) ([]string, error) {
	crdb := rc.ClusterClient
	//defer crdb.Close()
	reply, err := crdb.SMembers(key).Result()
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Set Set
func (rc *ClusterClient) Set(key string, value string) (string, error) {
	crdb := rc.ClusterClient
	//defer crdb.Close()
	reply, err := crdb.Set(key, value, 0).Result()
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Get Get
func (rc *ClusterClient) Get(key string) (string, error) {
	crdb := rc.ClusterClient
	//defer crdb.Close()
	reply, err := crdb.Get(key).Result()
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Incr Incr
func (rc *ClusterClient) Incr(transferKey string) (int, error) {
	crdb := rc.ClusterClient
	//defer crdb.Close()
	reply, err := crdb.Incr(transferKey).Result()
	if err == redis.ErrNil {
		err = nil
	}
	return int(reply), err
}
