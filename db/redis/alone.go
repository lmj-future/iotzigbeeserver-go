package redis

import "github.com/gomodule/redigo/redis"

//CloseSession CloseSession
func (rc *Client) CloseSession() error {
	return rc.Pool.Close()
}

//Rpush Rpush
func (rc *Client) Rpush(key string, data []byte) (int, error) {
	conn := rc.Pool.Get()
	defer conn.Close()
	reply, err := redis.Int(conn.Do("rpush", key, data))
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Lpop Lpop
func (rc *Client) Lpop(key string) (string, error) {
	conn := rc.Pool.Get()
	defer conn.Close()
	reply, err := redis.String(conn.Do("lpop", key))
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Lpush Lpush
func (rc *Client) Lpush(key string, data []byte) (int, error) {
	conn := rc.Pool.Get()
	defer conn.Close()
	reply, err := redis.Int(conn.Do("lpush", key, data))
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Del Del
func (rc *Client) Del(key string) (int, error) {
	conn := rc.Pool.Get()
	defer conn.Close()
	reply, err := redis.Int(conn.Do("del", key))
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Lrem Lrem
func (rc *Client) Lrem(key string, count int, value string) (int, error) {
	conn := rc.Pool.Get()
	defer conn.Close()
	reply, err := redis.Int(conn.Do("lrem", key, count, value))
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Lindex Lindex
func (rc *Client) Lindex(key string, start int) (string, error) {
	conn := rc.Pool.Get()
	defer conn.Close()
	reply, err := redis.String(conn.Do("lindex", key, start))
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Llen Llen
func (rc *Client) Llen(key string) (int, error) {
	conn := rc.Pool.Get()
	defer conn.Close()
	reply, err := redis.Int(conn.Do("llen", key))
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Lrange Lrange
func (rc *Client) Lrange(key string, start int, stop int) ([]string, error) {
	conn := rc.Pool.Get()
	defer conn.Close()
	reply, err := redis.Strings(conn.Do("lrange", key, start, stop))
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Lset Lset
func (rc *Client) Lset(key string, index int, value []byte) (string, error) {
	conn := rc.Pool.Get()
	defer conn.Close()
	reply, err := redis.String(conn.Do("lset", key, index, value))
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Sadd Sadd
func (rc *Client) Sadd(key string, member string) (int, error) {
	conn := rc.Pool.Get()
	defer conn.Close()
	reply, err := redis.Int(conn.Do("sadd", key, member))
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Srem Srem
func (rc *Client) Srem(key string, member string) (int, error) {
	conn := rc.Pool.Get()
	defer conn.Close()
	reply, err := redis.Int(conn.Do("srem", key, member))
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Smembers Smembers
func (rc *Client) Smembers(key string) ([]string, error) {
	conn := rc.Pool.Get()
	defer conn.Close()
	reply, err := redis.Strings(conn.Do("smembers", key))
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Set Set
func (rc *Client) Set(key string, value string) (string, error) {
	conn := rc.Pool.Get()
	defer conn.Close()
	reply, err := redis.String(conn.Do("set", key, value))
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Get Get
func (rc *Client) Get(key string) (string, error) {
	conn := rc.Pool.Get()
	defer conn.Close()
	reply, err := redis.String(conn.Do("get", key))
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}

//Incr Incr
func (rc *Client) Incr(transferKey string) (int, error) {
	conn := rc.Pool.Get()
	defer conn.Close()
	reply, err := redis.Int(conn.Do("incr", transferKey))
	if err == redis.ErrNil {
		err = nil
	}
	return reply, err
}
