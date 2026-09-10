package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Product struct {

	// Primitive não é usado na v2 do mongo-driver, então substituímos por bson.ObjectID
	// Motivo: Consolidaram alguns tipos no pacote bson, como ObjectID, Decimal128, DateTime, etc. E também consolidaram os pacotes bson e bson/primitive em um único pacote bson.
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Slug      string        `bson:"slug" json:"slug"`
	Name      string        `bson:"name" json:"name"`
	Quantity  int           `bson:"quantity" json:"quantity"`
	Price     float64       `bson:"price" json:"price"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

