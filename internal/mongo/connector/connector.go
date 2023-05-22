package connector

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(uri string ) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(uri))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %v", err)
	}

	return client, nil
}