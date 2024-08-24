package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserCart struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	ServiceID primitive.ObjectID `bson:"service_id" json:"service_id"`
	Quantity  int                `bson:"quantity" json:"quantity"` // Quantity of the service in the cart
}
