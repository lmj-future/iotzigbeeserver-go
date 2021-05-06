package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
)

//ClusterClient ClusterClient
type ClusterClient struct {
	ClusterClient *redis.ClusterClient
}

//NewClusterClient NewClusterClient
func NewClusterClient(servers []string, password string) (*ClusterClient, error) {
	var r *ClusterClient = new(ClusterClient)
	r.ClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:       servers,
		DialTimeout: 5 * time.Second,
		Password:    password,
	})
	if r.ClusterClient == nil {
		return nil, fmt.Errorf("create cluster client nil")
	}
	if r.ClusterClient.ClientGetName().Name() == "" {
		return nil, fmt.Errorf("create cluster client failed")
	}
	_, err := r.ClusterClient.Ping().Result()
	if err != nil {
		globallogger.Log.Errorln(err)
		return nil, err
	}
	return r, nil
}
