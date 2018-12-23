package mongodb

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/wangff15386/supermarket-go/conf"
)

// DBNAME for mongo
const (
	DBNAME   = "supermarkt-mongo"
	FOUNTNOT = "mongo: no documents in result"
)

// GetConn get mongodb client
func GetConn() (*mongo.Database, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Config.Mongo.TimeOut)*time.Second)
	client, err := mongo.Connect(ctx, conf.Config.Mongo.URL)
	if err != nil {
		panic(err)
	}

	return client.Database(DBNAME), ctx, cancel
}
