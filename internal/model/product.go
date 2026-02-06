package model

type Product struct {
	Id    int      `json:"productId"`
	Name  *string  `json:"productName"`
	Price *float64 `json:"productPrice"`
}