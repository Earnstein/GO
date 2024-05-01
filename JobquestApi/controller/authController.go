package controller

import (
	"JobquestApi/database"
	"JobquestApi/models"
	"JobquestApi/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var validate = validator.New()


func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			utils.HandleBadRequest(c, err.Error())
			return
		}

		if err := validate.Struct(user); err != nil {
			utils.HandleBadRequest(c, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			utils.HandleInternalServerError(c, err.Error())
			return
		}
		if count > 0 {		
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
			return
		}
		_, err = userCollection.InsertOne(ctx, user)
		if err != nil {
			utils.HandleInternalServerError(c, err.Error())
			return
		}
		c.JSON(http.StatusCreated, gin.H{"data": user})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}