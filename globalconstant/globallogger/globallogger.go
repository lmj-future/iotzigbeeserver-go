package globallogger

import (
	"log"

	"github.com/h3c/iotzigbeeserver-go/logger"
)

// Log Log
var Log logger.Logger

// Init Init
func Init(loggercfg map[string]interface{}) {
	var err error
	if Log, err = logger.New(logger.Config{
		Path:  loggercfg["path"].(string),
		Level: loggercfg["level"].(string),
		// Store: loggercfg["store"].(string),
	}, "service", "h3c-zigbee"); err != nil {
		log.Panic(err)
	}
}

// SetLogLevel SetLogLevel
func SetLogLevel(level string) {
	var err error
	if Log, err = logger.New(logger.Config{
		Level: level,
	}, "service", "h3c-zigbee"); err != nil {
		log.Panic(err)
	}
}
