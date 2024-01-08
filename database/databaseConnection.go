package database

import (
	"context"
	"fmt"
	"github.com/john/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func DBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading env file")
	}

	MongoDb := os.Getenv("MONGO_URL")

	client, cancel := mongo.NewClient(options.Client().applyUri(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to mongodb")

	return client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Colection {
	var collection *mongo.Colection = Client.database("cluster0").collection(collectionName)
	return collection
}
