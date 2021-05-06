package redis

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/letsfire/redigo"
	"github.com/letsfire/redigo/mode"
	"github.com/letsfire/redigo/mode/cluster"
)

func TestClient(t *testing.T) {
	host := "127.0.0.1"
	port := 6379
	connectionString := fmt.Sprintf("%s:%d", host, port)
	opts := []redis.DialOption{
		//redis.DialPassword(""),
		//redis.DialConnectTimeout(time.Duration(0) * time.Second),
	}
	conn, err := redis.Dial(
		"tcp", connectionString, opts...,
	)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	conn.Do("set", "aa", "11")
	reply, err := redis.String(conn.Do("get", "aa"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reply)
}

func TestNewClientPool(t *testing.T) {
	client, err := NewClientPool("127.0.0.1", "6379", "")
	if err != nil {
		fmt.Println("connect err", err)
		os.Exit(1)
	}
	defer client.CloseSession()
	reply, err := client.Set("mytest", "sss")
	if err != nil {
		fmt.Println("set err", err)
		os.Exit(1)
	}
	fmt.Println("set", reply)

	reply, err = client.Get("mytest")
	if err != nil {
		fmt.Println("set err", err)
		os.Exit(1)
	}
	fmt.Println("get", reply)
}

func TestClusterClient(t *testing.T) {
	var echoStr = "hello world"
	var clusterMode = cluster.New(
		cluster.Nodes([]string{
			//"192.168.0.110:30001", "192.168.0.110:30002", "192.168.0.110:30003",
			//"192.168.0.110:30004", "192.168.0.110:30005", "192.168.0.110:30006",
			"127.0.0.1:6379", //"127.0.0.1:30002", "127.0.0.1:30003",
			//"127.0.0.1:30004", "127.0.0.1:30005", "127.0.0.1:30006",
		}),
		cluster.PoolOpts(
			mode.MaxActive(0),       // 最大连接数，默认0无限制
			mode.MaxIdle(0),         // 最多保持空闲连接数，默认2*runtime.GOMAXPROCS(0)
			mode.Wait(false),        // 连接耗尽时是否等待，默认false
			mode.IdleTimeout(0),     // 空闲连接超时时间，默认0不超时
			mode.MaxConnLifetime(0), // 连接的生命周期，默认0不失效
			mode.TestOnBorrow(nil),  // 空间连接取出后检测是否健康，默认nil
		),
		cluster.DialOpts(
			redis.DialReadTimeout(time.Second),    // 读取超时，默认time.Second
			redis.DialWriteTimeout(time.Second),   // 写入超时，默认time.Second
			redis.DialConnectTimeout(time.Second), // 连接超时，默认500*time.Millisecond
			redis.DialPassword(""),                // 鉴权密码，默认空
			redis.DialDatabase(0),                 // 数据库号，默认0
			redis.DialKeepAlive(time.Minute*5),    // 默认5*time.Minute
			redis.DialNetDial(nil),                // 自定义dial，默认nil
			redis.DialUseTLS(false),               // 是否用TLS，默认false
			redis.DialTLSSkipVerify(false),        // 服务器证书校验，默认false
			redis.DialTLSConfig(nil),              // 默认nil，详见tls.Config
		),
	)
	var instance = redigo.New(clusterMode)

	res, _ := instance.String(func(c redis.Conn) (res interface{}, err error) {
		return c.Do("ECHO", echoStr)
	})
	fmt.Println(res)

	res, _ = instance.String(func(c redis.Conn) (res interface{}, err error) {
		return c.Do("set", "test0324", "0324")
	})
	fmt.Println(res)

	res, err := instance.String(func(c redis.Conn) (res interface{}, err error) {
		return c.Do("get", "test0324")
	})
	fmt.Println(res)

	if err != nil {
		log.Fatal(err)
	} else if res != echoStr {
		log.Fatalf("unexpected result, expect = %s, but = %s", echoStr, res)

	}
}
