package config

import (
	"context"
	"khanhanhtr/sample/mongo"
	"time"
)

func NewMongoDatabase(env *Env) (mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoPath := env.Mongo.MongoPath

	mongo_client, err := mongo.NewClient(ctx, mongoPath)
	if err != nil {
		return nil, err

	}

	if err := mongo_client.Ping(ctx); err != nil {
		return nil, err
	}

	return mongo_client, nil

}
