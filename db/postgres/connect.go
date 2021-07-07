package postgres

import (
	"os"
	"sync"
	"time"

	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/jinzhu/gorm"
)

// PgConnectPool PgConnectPool
type PgConnectPool struct {
}

var instance *PgConnectPool
var once sync.Once

var db *gorm.DB
var err error

// GetInstance GetInstance
func GetInstance() *PgConnectPool {
	once.Do(func() {
		instance = &PgConnectPool{}
	})
	return instance
}

func reConnectDB(postgresURL string) {
	timer := time.NewTimer(30 * time.Second)
	<-timer.C
	timer.Stop()
	db, err = gorm.Open("postgres", postgresURL)
	if err != nil {
		globallogger.Log.Errorln("[Postgres][Reconnect] couldn't open postgres:", err.Error())
		reConnectDB(postgresURL)
	} else {
		globallogger.Log.Errorf("[Postgres][Reconnect] connect success: db: %+v, db.DB(): %+v", db, db.DB())
		keepAliveDB(postgresURL)
	}
}

func keepAliveDB(postgresURL string) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		<-ticker.C
		err = db.DB().Ping()
		if err != nil {
			db, err = gorm.Open("postgres", postgresURL)
			if err != nil {
				globallogger.Log.Errorln("[Postgres][KeepAlive] couldn't open postgres:", err.Error())
			} else {
				globallogger.Log.Errorf("[Postgres][KeepAlive] connect success: db: %+v, db.DB(): %+v", db, db.DB())
			}
		} else {
			globallogger.Log.Errorf("[Postgres][KeepAlive] Ping success： db: %+v, db.DB(): %+v, status: %+v", db, db.DB(), db.DB().Stats())
		}
	}
}

// InitDataPool 初始化数据库连接(可在mail()适当位置调用)
func (pg *PgConnectPool) InitDataPool(postgresConnParas map[string]interface{}) (isSuccess bool) {
	userName := postgresConnParas["userName"].(string)
	password := postgresConnParas["password"].(string)
	host := postgresConnParas["host"].(string)
	port := postgresConnParas["port"].(string)
	dbName := postgresConnParas["dbName"].(string)
	sslMode := postgresConnParas["sslMode"].(string)
	if os.Getenv(userName) != "" {
		userName = os.Getenv(userName)
	}
	if os.Getenv(password) != "" {
		password = os.Getenv(password)
	}
	if os.Getenv(host) != "" {
		host = os.Getenv(host)
	}
	if os.Getenv(port) != "" {
		port = os.Getenv(port)
	}
	if os.Getenv(dbName) != "" {
		dbName = os.Getenv(dbName)
	}
	if os.Getenv(sslMode) != "" {
		sslMode = os.Getenv(sslMode)
	}
	postgresURL := "host=" + host + " port=" + port + " user=" + userName + " dbname=" + dbName + " password=" + password + " sslmode=" + sslMode
	db, err = gorm.Open("postgres", postgresURL)
	if err != nil {
		globallogger.Log.Errorln("couldn't open postgres:", err.Error())
		go reConnectDB(postgresURL)
		return false
	}
	db.DB().SetMaxOpenConns(5)
	globallogger.Log.Errorf("[Postgres] connect success: db: %+v, db.DB(): %+v", db, db.DB())
	//关闭数据库，db会被多个goroutine共享，可以不调用
	// defer db.Close()
	go keepAliveDB(postgresURL)
	return true
}

// GetPostGresDB 对外获取数据库连接对象db
func (pg *PgConnectPool) GetPostGresDB() (dbCon *gorm.DB) {
	return db
}
