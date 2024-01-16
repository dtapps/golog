package golog

import (
	"context"
	"errors"
	"go.dtapp.net/gotime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoDelete 删除N小时前数据
func (am *ApiMongo) MongoDelete(ctx context.Context, hour int64) error {
	return am.MongoDeleteDataCustom(ctx, am.mongoConfig.databaseName, am.mongoConfig.collectionName, hour)
}

// MongoDeleteDataCustom 删除N小时前数据
func (am *ApiMongo) MongoDeleteDataCustom(ctx context.Context, databaseName string, collectionName string, hour int64) error {
	if am.mongoConfig.stats == false {
		return nil
	}

	if databaseName == "" {
		return errors.New("没有设置库名")
	}
	if collectionName == "" {
		return errors.New("没有设置集合名")
	}
	filter := bson.D{{"log_time", bson.D{{"$lt", primitive.NewDateTimeFromTime(gotime.Current().BeforeHour(hour).Time)}}}}
	_, err := am.mongoClient.Database(databaseName).
		Collection(collectionName).
		DeleteMany(ctx, filter)
	return err
}
