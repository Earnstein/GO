package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"example.com/earnstein-api/controller"
	"example.com/earnstein-api/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	server         *gin.Engine
	userservice    services.UserService
	usercontroller controller.UserController
	ctx            context.Context
	usercollection *mongo.Collection
	mongoclient    *mongo.Client
	err            error
	MONGO_URL      string
)

func init() {
	godotenv.Load()
	ctx := context.TODO()
	MONGO_URL = os.Getenv("MONGO_URL")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	mongoconn := options.Client().ApplyURI(MONGO_URL).SetServerAPIOptions(serverAPI)
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
	}

	if err := mongoclient.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}

	fmt.Println("connected to mongodb")
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	usercollection = mongoclient.Database("Userdb").Collection("users")
	userservice = services.NewUserServiceImpl(usercollection, ctx)
	usercontroller = controller.New(userservice)
	server = gin.Default()
}
func main() {
	defer mongoclient.Disconnect(ctx)

	basepath := server.Group("/api")
	usercontroller.RegisterUserRoutes(basepath)
	log.Fatal(server.Run(":8080"))
}
