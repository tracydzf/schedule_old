package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"schedule/util/config"
	"schedule/util/data_schema"
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

// searchTask
func SearchTask(ctx context.Context, filter bson.D) (tasks []data_schema.TaskInfo, err error) {
	collection := mongoClient.Database(db).Collection("task")
	if collection == nil {
		log.ErrLogger.Printf("db.runoob.task is nil")
		return nil, fmt.Errorf("collection is nil")
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.ErrLogger.Printf("mongodb find task fail, filter:%+v, error:%+v", filter, err)
		return nil, err
	}
	tasks = make([]data_schema.TaskInfo, cursor.RemainingBatchLength())
	if err = cursor.All(ctx, &tasks); err != nil {
		log.ErrLogger.Printf("mongodb all function fail, filter:%+v, error:%+v", filter, err)
		return nil, err
	}
	return tasks, nil
}

// InsertOneSchedule 插入一条调度信息
func InsertOneSchedule(ctx context.Context, schedule data_schema.ScheduleHistory) (id string, err error) {
	collection := mongoClient.Database(db).Collection("schedule")
	if collection == nil {
		log.ErrLogger.Printf("db.runoob.schedule is nil")
		return "", fmt.Errorf("collection is nil")
	}
	var result *mongo.InsertOneResult
	var i int
	for i = 0; i < 3; i++ {
		if result, err = collection.InsertOne(ctx, schedule); err == nil {
			break
		}
	}

	if i >= 3 {
		log.ErrLogger.Printf("insert schedule into mongo fail, schedule:%+v, error:%+v", schedule, err)
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil

}
