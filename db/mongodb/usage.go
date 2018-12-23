package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	// client, err := mongo.NewClient("mongodb+srv://dbManager:fkQ6XFxnG9zWwzz@cexmarket-g8noj.azure.mongodb.net/test?retryWrites=true")
	// if err != nil {
	// 	panic(err)
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// if err = client.Connect(ctx); err != nil {
	// 	panic(err)
	// }
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	client, err := mongo.Connect(ctx, "mongodb://localhost:27017")
	// client, err := mongo.Connect(ctx, "mongodb://dbManager:fkQ6XFxnG9zWwzz@cexmarket-shard-00-00-g8noj.azure.mongodb.net:27017,cexmarket-shard-00-01-g8noj.azure.mongodb.net:27017,cexmarket-shard-00-02-g8noj.azure.mongodb.net:27017/test?ssl=true&replicaSet=cexmarket-shard-0&authSource=admin&retryWrites=true")
	if err != nil {
		panic(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 360*time.Second)
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	collection := client.Database("testing").Collection("numbers")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
	if err != nil {
		panic(err)
	}
	id := res.InsertedID
	fmt.Println("id:", id)

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, nil)
	if err != nil {
		panic(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			panic(err)
		}
		// do something with result....
		fmt.Println("result:", result)
	}
	if err := cur.Err(); err != nil {
		panic(err)
	}

	var result struct {
		Value float64
	}
	filter := bson.M{"name": "pi"}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		panic(err)
	}
	// Do something with result...
	fmt.Println("result:", result)
}
