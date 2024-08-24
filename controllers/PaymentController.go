package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
	"github.com/realjoni17/Hdocs/database"
	"github.com/realjoni17/Hdocs/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Payment(c *gin.Context) {
	db := database.GetDatabase()
	cartsCollection := db.Collection("usercarts")
	servicesCollection := db.Collection("services")
	ordersCollection := db.Collection("orders")

	userIDStr := c.Param("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Fetch the user's cart items
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

	// Calculate the total price
	serviceMap := make(map[primitive.ObjectID]models.Service)
	for _, service := range services {
		serviceMap[service.ID] = service
	}

	totalPrice := 0
	var orderItems []models.OrderItem
	for _, cart := range userCarts {
		if service, exists := serviceMap[cart.ServiceID]; exists {
			totalPrice += int(service.Price) * cart.Quantity
			orderItems = append(orderItems, models.OrderItem{
				ServiceID: cart.ServiceID,
				Quantity:  cart.Quantity,
			})
		}
	}

	// Initialize Razorpay client
	razorpayClient := razorpay.NewClient("rzp_test_69hZiriEnqQIjk", "HP4Hd5QC5cV8smh4q8Eqpob6")

	// Create payment options
	paymentOptions := map[string]interface{}{
		"amount":   totalPrice,
		"currency": "INR",
	}

	// Create a payment order
	payment, err := razorpayClient.Order.Create(paymentOptions, nil)
	if err != nil {
		fmt.Println("Razorpay Payment Creation Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment creation failed"})
		return
	}

	// Move cart items to orders collection
	order := models.Order{
		ID:             primitive.NewObjectID(),
		UserID:         userID,
		Items:          orderItems,
		TotalPrice:     float64(totalPrice),
		DeliveryStatus: "Pending",
		CreatedAt:      time.Now().Unix(),
	}

	_, err = ordersCollection.InsertOne(context.TODO(), order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Clear the user's cart
	_, err = cartsCollection.DeleteMany(context.TODO(), bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment": payment, "order": order})
}
