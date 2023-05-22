package mongohook

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoHook struct {
	*mongo.Client
}

func New(cl *mongo.Client) *mongoHook {
	return &mongoHook{cl}
}

func (mh *mongoHook) Fire(entry *logrus.Entry) error {

	data := make(logrus.Fields)
	data["Level"] = entry.Level.String()
	data["Time"] = entry.Time
	data["Message"] = entry.Message

	for k, v := range entry.Data {
		if errData, isError := v.(error); logrus.ErrorKey == k && v != nil && isError {
			data[k] = errData.Error()
		} else {
			data[k] = v
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := mh.Database("test_db").Collection("random_coll")
	_, err := coll.InsertOne(ctx, bson.M(data))

	if err != nil {
		return fmt.Errorf("failed to send log entry to mongodb: %v", err)
	}

	return nil
}

func (h *mongoHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
