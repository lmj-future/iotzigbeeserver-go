package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// Client represents a Redis client
type Client struct {
	Pool *redis.Pool // A thread-safe pool of connections to Redis
}

//NewClientPool NewClientPool
func NewClientPool(host string, port string, password string) (*Client, error) {
	var r *Client
	connectionString := fmt.Sprintf("%s:%s", host, port)
	opts := []redis.DialOption{
		redis.DialPassword(password),
		//redis.DialConnectTimeout(time.Duration(0) * time.Second),
	}
	dialFunc := func() (redis.Conn, error) {
		conn, err := redis.Dial(
			"tcp", connectionString, opts...,
		)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
	r = &Client{
		Pool: &redis.Pool{
			//IdleTimeout: 0,
			/* The current implementation processes nested structs using concurrent connections.
			 * With the deepest nesting level being 3, three shall be the number of maximum open
			 * idle connections in the pool, to allow reuse.
			 * TODO: Once we have a concurrent benchmark, this should be revisited.
			 * TODO: Longer term, once the objects are clean of external dependencies, the use
			 * of another serializer should make this moot.
			 */
			//MaxIdle: 10,
			Dial: dialFunc,
		},
	}
	conn := r.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("ping")
	if err != nil {
		return nil, err
	}
	return r, nil
}
