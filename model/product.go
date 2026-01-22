package model

type Product struct {
	Id    int     `json:"product_id"`
	Name  string  `json:"product_name"`
	Price float64 `json:"product_price"`
}