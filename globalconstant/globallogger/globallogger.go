package globallogger

import (
	"log"

	"github.com/h3c/iotzigbeeserver-go/logger"
)

// Log Log
var Log logger.Logger

// Init Init
func Init(loggercfg map[string]interface{}) {
	cfg := logger.Config{
		Path:  loggercfg["path"].(string),
		Level: loggercfg["level"].(string),
		// Store: loggercfg["store"].(string),
	}
	var err error
	if Log, err = logger.New(cfg, "service", "h3c-zigbee"); err != nil {
		log.Panic(err)
	}
}

// SetLogLevel SetLogLevel
func SetLogLevel(level string) {
	cfg := logger.Config{
		Level: level,
	}
	var err error
	if Log, err = logger.New(cfg, "service", "h3c-zigbee"); err != nil {
		log.Panic(err)
	}
}
