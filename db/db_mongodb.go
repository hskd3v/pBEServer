package db

import (
	"context"
	"time"

	"github.com/harriklein/pBE/pBEServer/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongodbInit = func() {

	log.Log.Println("Connecting into DB (MongoDB)...")

	var err error
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://twittor:twittor@cluster0.zy6ew.mongodb.net/twittor"))
	if err != nil {
		log.Log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Log.Fatal(err)
	}
	defer client.Disconnect(ctx)
}
