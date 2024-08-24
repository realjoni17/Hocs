package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/realjoni17/Hdocs/database"
	"github.com/realjoni17/Hdocs/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateService(c *gin.Context) {
	db := database.GetDatabase()
	servicesCollection := db.Collection("services")

	var newService models.Service
	if err := c.ShouldBindJSON(&newService); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the new service into the collection
	_, err := servicesCollection.InsertOne(context.TODO(), newService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "while creating service"})
		return
	}

	c.JSON(http.StatusCreated, newService)
}

func GetService(c *gin.Context) {
	db := database.GetDatabase()
	servicesCollection := db.Collection("services")

	cursor, err := servicesCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot find services"})
		return
	}
	defer cursor.Close(context.TODO())

	var services []models.Service
	if err := cursor.All(context.TODO(), &services); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode services"})
		return
	}

	c.JSON(http.StatusOK, services)
}

func AddServiceToCart(c *gin.Context) {
	db := database.GetDatabase()
	cartsCollection := db.Collection("usercarts")

	userIDStr := c.Param("user_id")
	serviceIDStr := c.Param("service_id")

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	serviceID, err := primitive.ObjectIDFromHex(serviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	userCart := models.UserCart{
		UserID:    userID,
		ServiceID: serviceID,
	}

	_, err = cartsCollection.InsertOne(context.TODO(), userCart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add service to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service added to cart"})
}

func GetUserTotalServices(c *gin.Context) {
	db := database.GetDatabase()
	cartsCollection := db.Collection("usercarts")
	servicesCollection := db.Collection("services")

	userIDStr := c.Param("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find all services in the user's cart
	cursor, err := cartsCollection.Find(context.TODO(), bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find services in cart"})
		return
	}
	defer cursor.Close(context.TODO())

	var userCarts []models.UserCart
	if err := cursor.All(context.TODO(), &userCarts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user cart"})
		return
	}

	// Find all services by their IDs
	var serviceIDs []primitive.ObjectID
	for _, cart := range userCarts {
		serviceIDs = append(serviceIDs, cart.ServiceID)
	}

	servicesCursor, err := servicesCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": serviceIDs}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find services"})
		return
	}
	defer servicesCursor.Close(context.TODO())

	var services []models.Service
	if err := servicesCursor.All(context.TODO(), &services); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode services"})
		return
	}

	c.JSON(http.StatusOK, services)
}
