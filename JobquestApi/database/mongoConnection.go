package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var Client *mongo.Client = DBConnection()

func DBConnection() *mongo.Client {
	err := godotenv.Load(".env")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 100*time.Second)
	defer cancel()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	MONGO_URL := os.Getenv("MONGO_URL")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
 	mongoconn := options.Client().ApplyURI(MONGO_URL).SetServerAPIOptions(serverAPI)
	mongoclient, err := mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal("Error connecting to mongodb", err.Error())
	}
	if err:= mongoclient.Database("admin").RunCommand(ctx,bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to mongodb")
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return mongoclient
}


func OpenCollection(client *mongo.Client, dbName, collectionName string) *mongo.Collection {
	collection := client.Database(dbName).Collection(collectionName)
	return collection
}