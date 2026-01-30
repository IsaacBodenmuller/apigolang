package repository

import (
	"APIGolang/model"
	"database/sql"
	"fmt"
)

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts() ([]model.Product, error) {
	query := "SELECT product_id, product_name, product_price FROM product"
	rows, err := pr.connection.Query(query)
	if err != nil {
		fmt.Println(err)
		return []model.Product{}, err
	}

	var productList []model.Product
	var productObj model.Product

	for rows.Next(){
		err = rows.Scan(
			&productObj.Id,
			&productObj.Name,
			&productObj.Price)
		if err != nil{
			fmt.Println(err)
			return []model.Product{}, err
		}
		
		productList = append(productList, productObj)
	}
	rows.Close()
	return productList, nil

}

func (pr *ProductRepository) GetProductById(product_id int) (*model.Product, error) {
	
	query, err := pr.connection.Prepare("SELECT * FROM product WHERE product_id = $1")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var produto model.Product
	err = query.QueryRow(product_id).Scan(
		&produto.Id,
		&produto.Name,
		&produto.Price,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	query.Close()
	return &produto, nil
}

func (pr *ProductRepository) CreateProduct(product model.Product) (int, error) {

	var id int
	query, err := pr.connection.Prepare("INSERT INTO product" +
		"(product_name, product_price)" +
		" VALUES ($1, $2) RETURNING product_id")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	err = query.QueryRow(product.Name, product.Price).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	query.Close()
	return id, nil
}

func (pr *ProductRepository) UpdateProductById(product_id int, product model.Product) (*model.Product, error) {

	oldProduct, err := pr.GetProductById(product_id)
	if err != nil{
		fmt.Println(err)
		return nil, err
	}
	
	if product.Name == nil {
		product.Name = oldProduct.Name
	}
	if product.Price == nil {
		product.Price = oldProduct.Price
	}
	
	var updatedProduct model.Product

	query, err := pr.connection.Prepare("UPDATE product" +
		" SET product_name = $1, product_price = $2" + 
		" WHERE product_id = $3 RETURNING product_id, product_name, product_price")
	if err != nil{
		fmt.Println(err)
		return nil, err
	}

	err = query.QueryRow(product.Name, product.Price, product_id).Scan(&updatedProduct.Id, &updatedProduct.Name, &updatedProduct.Price)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	query.Close()
	return &updatedProduct, nil
}

func (pr *ProductRepository) DeleteProductById(product_id int) (bool, error) {

	query := "DELETE FROM product" + 
		" WHERE product_id = $1"

	result, err := pr.connection.Exec(query, product_id)
	if err != nil{
		fmt.Println(err)
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, nil
	}

	return true, nil

}