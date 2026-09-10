package model

import (
	"time"
	"github.com/notoliveira/estoque/domain"
)

type CreateProductRequest struct {
	Slug     string  `json:"slug"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type UpdateProductRequest struct {
	Slug     string  `json:"slug"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type ProductResponse struct {
	ID        string  `json:"id"`
	Slug      string  `json:"slug"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (r CreateProductRequest) ToDomain() *domain.Product {
	return &domain.Product{
		Slug:      r.Slug,
		Name:      r.Name,
		Quantity:  r.Quantity,
		Price:     r.Price,
	}
}

func (r UpdateProductRequest) ToDomain() *domain.Product {
	return &domain.Product{
		Slug:      r.Slug,
		Name:      r.Name,
		Quantity:  r.Quantity,
		Price:     r.Price,
	}
}

func ToResponse(p domain.Product) ProductResponse {
	return ProductResponse{
		ID:        p.ID.Hex(),
		Slug:      p.Slug,
		Name:      p.Name,
		Quantity:  p.Quantity,
		Price:     p.Price,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}