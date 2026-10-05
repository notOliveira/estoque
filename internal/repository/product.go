package repository

import (
	"context"
	"errors"

	"github.com/notoliveira/estoque/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrNotFound = errors.New("product not found")

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
	mongoId := id.(bson.ObjectID)

	product.ID = mongoId

	return nil
}

func (r *mongoProductRepository) FindAll(ctx context.Context) ([]*domain.Product, error) {

	// Find retorna um cursor porque a consulta pode trazer vários documentos any
	// bson.M{} = filtro vazio, portanto busca todos os produtos
	cursor, err := r.collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	// O cursor é um recurso que precisa ser fechado depois da consulta
	// Com defer, Close será executado quando FindAll terminar
	defer cursor.Close(ctx)

	// Slice que receberá todos os produtos encontrados
	// Cada elemento é um ponteiro para domain.Product
	var products []*domain.Product

	// All percorre todos os documentos do cursor e preenche o slice
	// Passamos &products porque All precisa modificar a variável products
	// & = endereço da variável (ponteiro), permitindo que a função a altere
	// Já faz o decode de cada documento para domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *mongoProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {

	// Transformando o id em primitive.ObjectId
	objID, err := bson.ObjectIDFromHex(id)
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
			return nil, ErrNotFound
		}
		// Se for outro erro (ex: banco caiu), retorna o erro real
		return nil, err
	}

	return &product, nil
}

func (r *mongoProductRepository) Update(ctx context.Context, id string, product *domain.Product) error {

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	f := bson.M{"_id": objID}
	result, err := r.collection.ReplaceOne(ctx, f, product)
	if err != nil {
		return err
	}

	// Verifica se foi encontrado algum documento para atualizar. Se não, retorna ErrNotFound
	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *mongoProductRepository) Delete(ctx context.Context, id string) error {

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	f := bson.M{"_id": objID}
	result, err := r.collection.DeleteOne(ctx, f)
	if err != nil {
		return err
	}

	// Retorna a quantidade de documentos deletados. Se for 0, significa que o documento não foi encontrado
	if result.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}
