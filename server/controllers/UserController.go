package controllers

import (
	"context"
	"mediwatch/server/configs"
	"mediwatch/server/models"
	"mediwatch/server/views"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "Users")
var validate = validator.New()

func GetUser(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusBadRequest, Message: "Error", Data: err.Error()})
		return
	}

	err = userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, views.UserView{Status: http.StatusNotFound, Message: "Not Found", Data: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusInternalServerError, Message: "Error", Data: err.Error()})
		return
	}

	c.JSON(http.StatusFound, views.UserView{Status: http.StatusFound, Message: "Success", Data: user})
}

func GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User

	results, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusInternalServerError, Message: "Error", Data: err.Error()})
		return
	}

	defer results.Close(ctx)

	for results.Next(ctx) {
		var user models.User

		if err = results.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusInternalServerError, Message: "Error", Data: err.Error()})
			return
		}

		users = append(users, user)
	}

	c.JSON(http.StatusFound, views.UserView{Status: http.StatusFound, Message: "Success", Data: users})
}

func CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, views.UserView{Status: http.StatusBadRequest, Message: "Error", Data: err.Error()})
		return
	}

	if validationErr := validate.Struct(&user); validationErr != nil {
		c.JSON(http.StatusBadRequest, views.UserView{Status: http.StatusBadRequest, Message: "Error", Data: validationErr.Error()})
		return
	}

	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Email:    user.Email,
		Password: user.Password,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusInternalServerError, Message: "Error", Data: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, views.UserView{Status: http.StatusCreated, Message: "success", Data: result})
}
