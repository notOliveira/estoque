package repository

import (
	"context"
	"estoque/domain"
	"go.mongodb.org/mongo-driver/mongo"
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
	return nil
}

func (r *mongoProductRepository) FindAll(ctx context.Context) ([]*domain.Product, error) {
	return nil, nil
}

func (r *mongoProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	return nil, nil
}

func (r *mongoProductRepository) Update(ctx context.Context, id string, product *domain.Product) error {
	return nil
}

func (r *mongoProductRepository) Delete(ctx context.Context, id string) error {
	return nil
}