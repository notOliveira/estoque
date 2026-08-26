package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Slug      string             `bson:"slug" json:"slug"`
	Name      string             `bson:"name" json:"name"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	Price     float64            `bson:"price" json:"price"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}