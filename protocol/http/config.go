package http

import (
	"time"
)

// ServerInfo http server config
type ServerInfo struct {
	Address     string        `yaml:"address" json:"address"`
	Timeout     time.Duration `yaml:"timeout" json:"timeout" default:"5m"`
	Certificate `yaml:",inline" json:",inline"`
}

// ClientInfo http client config
type ClientInfo struct {
	Address     string        `yaml:"address" json:"address"`
	Timeout     time.Duration `yaml:"timeout" json:"timeout" default:"5m"`
	KeepAlive   time.Duration `yaml:"keepalive" json:"keepalive" default:"10m"`
	Username    string        `yaml:"username" json:"username"`
	Password    string        `yaml:"password" json:"password"`
	Certificate `yaml:",inline" json:",inline"`
}

// Certificate certificate config for mqtt server
type Certificate struct {
	CA       string `yaml:"ca" json:"ca"`
	Key      string `yaml:"key" json:"key"`
	Cert     string `yaml:"cert" json:"cert"`
	Insecure bool   `yaml:"insecure" json:"insecure"` // for client, for test purpose
}
