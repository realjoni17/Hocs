package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id"`
	Items          []OrderItem        `bson:"items" json:"items"`
	TotalPrice     float64            `bson:"total_price" json:"total_price"`
	DeliveryStatus string             `bson:"delivery_status" json:"delivery_status"` // e.g., "Pending", "Shipped", "Delivered"
	CreatedAt      int64              `bson:"created_at" json:"created_at"`           // Timestamp for order creation
}

type OrderItem struct {
	ServiceID primitive.ObjectID `bson:"service_id" json:"service_id"`
	Quantity  int                `bson:"quantity" json:"quantity"`
}
