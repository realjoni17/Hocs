package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/realjoni17/Hdocs/database"
	"github.com/realjoni17/Hdocs/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(c *gin.Context) {
	db := database.GetDatabase()
	usersCollection := db.Collection("users")

	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser.ID = primitive.NewObjectID() // Generate a new ObjectID
	_, err := usersCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "while creating user"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}
