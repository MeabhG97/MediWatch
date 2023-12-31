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
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "Users")
var validate = validator.New()

const defaultNumberCompartments = 7

var LoggedInUser primitive.ObjectID

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

	c.JSON(http.StatusOK, views.UserView{Status: http.StatusOK, Message: "Success", Data: user})
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

	c.JSON(http.StatusOK, views.UserView{Status: http.StatusOK, Message: "Success", Data: users})
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

	newPillbox := models.Pillbox{
		Id:                 primitive.NewObjectID(),
		NumberCompartments: defaultNumberCompartments,
	}

	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Email:    user.Email,
		Password: user.Password,
		Pillbox:  newPillbox,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusInternalServerError, Message: "Error", Data: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, views.UserView{Status: http.StatusCreated, Message: "success", Data: result})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusBadRequest, Message: "Error", Data: err.Error()})
		return
	}

	if err = c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, views.UserView{Status: http.StatusBadRequest, Message: "Error", Data: err.Error()})
		return
	}

	if err := validate.Struct(&user); err != nil {
		c.JSON(http.StatusBadRequest, views.UserView{Status: http.StatusBadRequest, Message: "Error", Data: err.Error()})
		return
	}

	update := bson.M{"email": user.Email, "password": user.Password}

	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusInternalServerError, Message: "Error", Data: err.Error()})
		return
	}

	var updatedUser models.User

	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusInternalServerError, Message: "Error", Data: err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, views.UserView{Status: http.StatusOK, Message: "Updated", Data: updatedUser})
}

func DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusBadRequest, Message: "Error", Data: err.Error()})
		return
	}

	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusInternalServerError, Message: "Error", Data: err.Error()})
		return
	}

	if result.DeletedCount < 1 {
		c.JSON(http.StatusNotFound, views.UserView{Status: http.StatusNotFound, Message: "Not Found", Data: "Matching id not found"})
		return
	}

	c.JSON(http.StatusOK, views.UserView{Status: http.StatusOK, Message: "Deleted", Data: "User deleted"})
}

func RegisterUser(c *gin.Context) {
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

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)

	if err != nil {
		c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusInternalServerError, Message: "Error", Data: err.Error()})
		return
	}

	newPillbox := models.Pillbox{
		Id:                 primitive.NewObjectID(),
		NumberCompartments: defaultNumberCompartments,
	}

	newMedications := []models.Medication{}
	newSchedule := []models.Schedule{}
	newHistory := []models.History{}
	newEvents := []models.Event{}

	newUser := models.User{
		Id:          primitive.NewObjectID(),
		Email:       user.Email,
		Password:    string(passwordHash),
		Pillbox:     newPillbox,
		Medications: newMedications,
		Schedule:    newSchedule,
		History:     newHistory,
		Events:      newEvents,
	}

	err = userCollection.FindOne(ctx, bson.M{"email": newUser.Email}).Decode(&user)

	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusInternalServerError, Message: "Error", Data: err.Error()})
		return
	}

	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusBadRequest, views.UserView{Status: http.StatusNotFound, Message: "User already exists", Data: "User already exists"})
		return
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusInternalServerError, Message: "Error", Data: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, views.UserView{Status: http.StatusCreated, Message: "success", Data: result})
}

func LoginUser(c *gin.Context) {
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

	loginUser := models.User{
		Email:    user.Email,
		Password: user.Password,
	}

	err := userCollection.FindOne(ctx, bson.M{"email": loginUser.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, views.UserView{Status: http.StatusNotFound, Message: "Not Found", Data: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusInternalServerError, Message: "Error", Data: err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, views.UserView{Status: http.StatusBadRequest, Message: "Error", Data: err.Error()})
		return
	}

	loginUser = models.User{
		Id:          user.Id,
		Email:       user.Email,
		Pillbox:     user.Pillbox,
		Medications: user.Medications,
		Schedule:    user.Schedule,
		History:     user.History,
	}

	LoggedInUser = loginUser.Id

	c.JSON(http.StatusOK, views.UserView{Status: http.StatusOK, Message: "Success", Data: loginUser})
}

func GoogleLoginUser(c *gin.Context) {
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

	loginUser := models.User{
		Email: user.Email,
	}

	err := userCollection.FindOne(ctx, bson.M{"email": loginUser.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			newPillbox := models.Pillbox{
				Id:                 primitive.NewObjectID(),
				NumberCompartments: defaultNumberCompartments,
			}

			newMedications := []models.Medication{}
			newSchedule := []models.Schedule{}
			newHistory := []models.History{}
			newEvents := []models.Event{}

			newUser := models.User{
				Id:          primitive.NewObjectID(),
				Email:       user.Email,
				Pillbox:     newPillbox,
				Medications: newMedications,
				Schedule:    newSchedule,
				History:     newHistory,
				Events:      newEvents,
			}

			result, err := userCollection.InsertOne(ctx, newUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusInternalServerError, Message: "Error", Data: err.Error()})
				return
			}

			err = userCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusInternalServerError, Message: "Error", Data: err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, views.UserView{Status: http.StatusInternalServerError, Message: "Error", Data: err.Error()})
			return
		}
	}

	loginUser = models.User{
		Id:          user.Id,
		Email:       user.Email,
		Pillbox:     user.Pillbox,
		Medications: user.Medications,
		Schedule:    user.Schedule,
		History:     user.History,
	}

	LoggedInUser = loginUser.Id

	c.JSON(http.StatusOK, views.UserView{Status: http.StatusOK, Message: "Success", Data: loginUser})
}
