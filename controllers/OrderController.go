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

func UpdateOrderStatus(c *gin.Context) {
	db := database.GetDatabase()
	ordersCollection := db.Collection("orders")

	orderIDStr := c.Param("order_id")
	orderID, err := primitive.ObjectIDFromHex(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var updateRequest struct {
		DeliveryStatus string `json:"delivery_status"`
	}

	if err := c.BindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	_, err = ordersCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": orderID},
		bson.M{"$set": bson.M{"delivery_status": updateRequest.DeliveryStatus}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Order status updated"})
}

func GetOrderStatus(c *gin.Context) {
	db := database.GetDatabase()
	ordersCollection := db.Collection("orders")

	orderIDStr := c.Param("order_id")
	orderID, err := primitive.ObjectIDFromHex(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	err = ordersCollection.FindOne(context.TODO(), bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find order"})
		return
	}

	c.JSON(http.StatusOK, order)
}
