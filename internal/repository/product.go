package repository

import (
	"context"

	"github.com/notoliveira/estoque/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductRepository interface {
	// retorno do produto não é necessário, pois o produto já é passado como referência, então qualquer alteração feita no produto dentro do método será refletida fora dele
	Create(ctx context.Context, product *domain.Product) error
	FindAll(ctx context.Context) ([]*domain.Product, error)
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	Update(ctx context.Context, id string, product *domain.Product) error
	Delete(ctx context.Context, id string) error
}

type mongoProductRepository struct {
	collection *mongo.Collection
}

func NewMongoProductRepository(collection *mongo.Collection) ProductRepository {
	return &mongoProductRepository{
		collection: collection,
	}
}

func (r *mongoProductRepository) Create(ctx context.Context, product *domain.Product) error {

	res, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return err
	}

	id := res.InsertedID

	// Type assertion -> Dizer ao Go qual o tipo do objeto
	// Nesse caso, estamos "transformando" o res.InsertedId (do tipo any) em primitive.ObjectId ao fazer id.(primitive.ObjectID)
	// Esse .() serve justamente para "tipar" o objeto
    mongoId := id.(primitive.ObjectID)

	product.ID = mongoId

	return nil
}

func (r *mongoProductRepository) FindAll(ctx context.Context) ([]*domain.Product, error) {
	return nil, nil
}

func (r *mongoProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {

	// Transformando o id em primitive.ObjectId
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
        return nil, err
    }
	
	// Criando o filtro (_id = <id>)
	f := bson.M{"_id": objID}

	// Criando product
	var product domain.Product

	err = r.collection.FindOne(ctx, f).Decode(&product)
	if err != nil {
        // Se o erro for "não encontrou nada", você pode retornar nil para o produto
        if err == mongo.ErrNoDocuments {
            return nil, nil
        }
        // Se for outro erro (ex: banco caiu), retorna o erro real
        return nil, err
    }

	return &product, nil
}

func (r *mongoProductRepository) Update(ctx context.Context, id string, product *domain.Product) error {
	return nil
}

func (r *mongoProductRepository) Delete(ctx context.Context, id string) error {
	return nil
}
