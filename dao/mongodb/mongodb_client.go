package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"schedule/util/config"
	"schedule/util/log"
	"time"
)

var mongoClient *mongo.Client
var db string

func InitMongodb() error {
	clientSocket := config.Viper.GetString("mongodb.conn")
	clientOptions := options.Client().ApplyURI(clientSocket)
	clientOptions.SetConnectTimeout(config.Viper.GetDuration("mongodb.conn_timeout") * time.Second)
	clientOptions.SetSocketTimeout(config.Viper.GetDuration("mongodb.timeout") * time.Second)
	var err error

	// 连接到MongoDB
	mongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	mongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.ErrLogger.Printf("mongodb connect fail, err:%+v, clienSocket:%s", err.Error(), clientSocket)
		return err
	}

	// 检查连接
	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.ErrLogger.Printf("mongodb ping error, error:%+v", err.Error())
		return err
	}

	db = config.Viper.GetString("mongodb.database")
	return nil

}
