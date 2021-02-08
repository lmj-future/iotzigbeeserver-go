package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type mongoCfg struct {
	Credentials  options.Credential
	uri          string
	databaseName string
	collection   string
	timeout      time.Duration
}

type mongoHook struct {
	Client     *mongo.Client
	collection *mongo.Collection
	config     mongoCfg
	level      logrus.Level
}

// NewMongoLogger create a new mongodb logger useing default config
func NewMongoLogger(c Config, fields ...string) (Logger, error) {
	// 使用默认配置初始化日志
	mgo := &mongoHook{}
	mgo, err := mgo.newMongoLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to init mongo logger: %v", err)
	}
	logLevel, err := logrus.ParseLevel(c.Level)
	if err != nil {
		log.Printf("failed to parse log level (%s), use default level (info)\n", c.Level)
		logLevel = logrus.InfoLevel
	}
	mgo.level = logLevel
	logrusFields := logrus.Fields{}
	for index := 0; index < len(fields)-1; index = index + 2 {
		logrusFields[fields[index]] = fields[index+1]
	}

	// 根据时间逆序对系统日志创建索引
	indexOpts := options.CreateIndexes().SetMaxTime(60 * time.Second)
	indexView := mgo.collection.Indexes()
	result, err := indexView.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.M{"time": -1},
			Options: options.Index().SetBackground(true),
		},
		indexOpts,
	)
	if result == "" || err != nil {
		log.Printf("Create time -1 index failed, error: %v\n", err)
	}
	// 根据时间顺序对系统日志创建索引
	indexOpts = options.CreateIndexes().SetMaxTime(60 * time.Second)
	indexView = mgo.collection.Indexes()
	result, err = indexView.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.M{"time": 1},
			Options: options.Index().SetBackground(true),
		},
		indexOpts,
	)
	if result == "" || err != nil {
		log.Printf("Create time 1 index failed, error: %v\n", err)
	}

	entry := logrus.NewEntry(logrus.New())
	entry.Level = logLevel
	entry.Logger.Level = logLevel
	entry.Logger.Formatter = newFormatter(c.Format)
	if mgo != nil {
		entry.Logger.Hooks.Add(mgo)
	}
	return &logger{entry.WithFields(logrusFields)}, nil
}

func (mgo *mongoHook) newMongoLogger() (*mongoHook, error) {
	cfg := mgo.defaultCfg()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout*time.Second)
	defer cancel()

	opt := options.Client().ApplyURI(cfg.uri)
	opt.SetAuth(cfg.Credentials)

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("connect mongdb failed: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("ping mongdb failed: %v", err)
	}

	return &mongoHook{
		Client:     client,
		config:     cfg,
		collection: client.Database("h3c-log").Collection("system"),
	}, nil
}

func (mgo *mongoHook) defaultCfg() mongoCfg {
	return mongoCfg{
		uri: "mongodb://h3c-db:27017",
		Credentials: options.Credential{
			Username:   "h3c-log",
			Password:   "h3c-iot-edge-log",
			AuthSource: "h3c-log",
		},
		databaseName: "h3c-log",
		collection:   "system",
		timeout:      5,
	}
}

func (mgo *mongoHook) Disconnect() error {
	if mgo.Client == nil {
		return nil
	}

	if err := mgo.Client.Disconnect(context.TODO()); err != nil {
		return err
	}

	return nil
}

func (mgo *mongoHook) Fire(entry *logrus.Entry) error {
	if mgo.level < entry.Level {
		return nil
	}

	fields := make(logrus.Fields)
	fields["level"] = entry.Level.String()
	fields["time"] = entry.Time
	fields["message"] = entry.Message
	fields["raw_time"] = entry.Time.Format("2006-01-02 15:04:05")

	for k, v := range entry.Data {
		if errInfo, isErr := v.(error); k == logrus.ErrorKey && v != nil && isErr {
			fields[k] = errInfo.Error()
		} else {
			fields[k] = v
		}
	}

	_, err := mgo.collection.InsertOne(context.TODO(), bson.M(fields))
	if err != nil {
		log.Printf("failed to send log entry to mongdb: %v\n", err)
		return err
	}

	return nil
}

func (mgo *mongoHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// NewMongoLoggerWithDatabase create a new mongodb logger using given mongo database
func NewMongoLoggerWithDatabase(c Config, client *mongo.Database, fields ...string) (Logger, error) {
	if client == nil {
		return NewMongoLogger(c, fields...)
	}

	// 解析日志级别
	logLevel, err := logrus.ParseLevel(c.Level)
	if err != nil {
		log.Printf("failed to parse log level (%s), use default level (info)", c.Level)
		logLevel = logrus.InfoLevel
	}

	mgo := &mongoHook{}
	collection := "log"
	logrusFields := logrus.Fields{}
	for index := 0; index < len(fields)-1; index = index + 2 {
		if fields[index] == "collection" {
			collection = fields[index+1]
			continue
		}
		logrusFields[fields[index]] = fields[index+1]
	}

	mgo.level = logLevel
	mgo.collection = client.Collection(collection)

	entry := logrus.NewEntry(logrus.New())
	entry.Level = logLevel
	entry.Logger.Level = logLevel
	entry.Logger.Out = os.Stdout
	entry.Logger.Formatter = newFormatter(c.Format)
	entry.Logger.Hooks.Add(mgo)

	return &logger{entry.WithFields(logrusFields)}, nil
}
