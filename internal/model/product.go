package model

type Product struct {
	ID    int
	Slug  string
	Name  string
	Quantity int
	Price float64
}

type CreateProductRequest struct {
	Slug     string  `json:"slug"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type UpdateProductRequest struct {
	Slug     *string  `json:"slug"`
	Name     *string  `json:"name"`
	Quantity *int     `json:"quantity"`
	Price    *float64 `json:"price"`
}

type ProductResponse struct {
	ID       string  `json:"id"`
	Slug     string  `json:"slug"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

