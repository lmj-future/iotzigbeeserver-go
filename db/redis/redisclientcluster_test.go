package redis

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-redis/redis"
)

//"github.com/gomodule/redigo/redis"
//"github.com/go-redis/redis"
//"github.com/gitstliu/go-redis-cluster"

// func TestClusterClient(t *testing.T) {
// 	cluster, _ := redis.NewCluster(
// 		&redis.Options{
// 			StartNodes:   []string{"127.0.0.1:6379", "127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005"},
// 			ConnTimeout:  50 * time.Millisecond,
// 			ReadTimeout:  50 * time.Millisecond,
// 			WriteTimeout: 50 * time.Millisecond,
// 			KeepAlive:    16,
// 			AliveTime:    60 * time.Second,
// 		})
// 	reply, err := redis.String(cluster.Do("set", "name", "itheima"))
// 	fmt.Println(reply, err)

// 	name, err := redis.String(cluster.Do("get", "name"))
// 	fmt.Println(name, err)
// 	//beego.Info(name)
// 	//this.Ctx.WriteString("集群创建成功")
// }

func TestAlone(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()
	reply, err := rdb.Set("mytest032501", "032501", 0).Result()
	fmt.Println(reply, err)
	reply, err = rdb.Get("mytest032501").Result()
	fmt.Println(reply, err)
}

func TestCluster(t *testing.T) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		//Addrs: []string{":6379"}, //":7001", ":7002", ":7003", ":7004", ":7005"},
		Addrs: []string{"127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002",
			"127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005"},
	})
	defer rdb.Close()
	reply, err := rdb.Set("mytest0325", "0325", 0).Result()
	fmt.Println(reply, err)
	reply, err = rdb.Get("mytest0325").Result()
	fmt.Println(reply, err)
}

func TestNewClusterClient(t *testing.T) {
	rdb, err := NewClusterClient("127.0.0.1", "6379", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rdb.CloseSession()
	reply, err := rdb.Set("mytest032501", "032501")
	fmt.Println(reply, err)
	reply, err = rdb.Get("mytest032501")
	fmt.Println(reply, err)
}
