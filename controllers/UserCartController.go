package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/realjoni17/Hdocs/database"
	"github.com/realjoni17/Hdocs/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserCart(c *gin.Context) {
	db := database.GetDatabase()
	cartsCollection := db.Collection("usercarts")
	servicesCollection := db.Collection("services")

	userIDStr := c.Param("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find all items in the user's cart
	cursor, err := cartsCollection.Find(context.TODO(), bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find items in cart"})
		return
	}
	defer cursor.Close(context.TODO())

	var userCarts []models.UserCart
	if err := cursor.All(context.TODO(), &userCarts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user cart"})
		return
	}

	// Retrieve service details for each item in the cart
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

	// Create a map to associate services with user cart items
	serviceMap := make(map[primitive.ObjectID]models.Service)
	for _, service := range services {
		serviceMap[service.ID] = service
	}

	// Create a response struct that includes both cart items and services
	type CartItemWithService struct {
		UserCart models.UserCart `json:"cart_item"`
		Service  models.Service  `json:"service"`
	}

	var response []CartItemWithService
	for _, cart := range userCarts {
		if service, exists := serviceMap[cart.ServiceID]; exists {
			response = append(response, CartItemWithService{
				UserCart: cart,
				Service:  service,
			})
		}
	}

	c.JSON(http.StatusOK, response)
}

// Retrieve service details for each

func AddItemToCart(c *gin.Context) {
	db := database.GetDatabase()
	cartsCollection := db.Collection("usercarts")

	// Extract service_id and user_id from the URL
	serviceIDStr := c.Param("service_id")
	userIDStr := c.Param("user_id")

	// Convert serviceIDStr and userIDStr to ObjectID
	serviceID, err := primitive.ObjectIDFromHex(serviceIDStr)
	if err != nil {
		log.Printf("Error converting service ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		log.Printf("Error converting user ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Create a new UserCart item
	newItem := models.UserCart{
		ServiceID: serviceID,
		UserID:    userID,
	}

	// Check if the item already exists in the cart
	existingItem := new(models.UserCart)
	err = cartsCollection.FindOne(context.TODO(), bson.M{"user_id": userID, "service_id": serviceID}).Decode(existingItem)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Item already exists in the cart"})
		return
	}

	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing item"})
		return
	}

	// Insert the new item into the cart
	insertResult, err := cartsCollection.InsertOne(context.TODO(), newItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"inserted_id": insertResult.InsertedID,
		"item":        newItem,
	})
}

func DeleteItemFromCart(c *gin.Context) {
	db := database.GetDatabase()
	cartsCollection := db.Collection("usercarts")

	serviceIDStr := c.Param("service_id")
	userIDStr := c.Param("user_id")

	serviceID, err := primitive.ObjectIDFromHex(serviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	result, err := cartsCollection.DeleteOne(context.TODO(), bson.M{"user_id": userID, "service_id": serviceID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found in cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}
