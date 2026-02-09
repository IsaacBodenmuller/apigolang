package controller

import (
	"APIGolang/internal/model"
	"APIGolang/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type productController struct {
	productUsecase usecase.ProductUsecase
}

func NewProductController(usecase usecase.ProductUsecase) productController {
	return productController{
		productUsecase: usecase,
	}
}

// GetProducts godoc
// @Summary Listar produtos
// @Description Retorna todos os produtos
// @Tags Products
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Product
// @Failure 401 {object} map[string]string
// @Router /products [get]
func (p *productController) GetProducts(ctx *gin.Context) {
	
	products, err := p.productUsecase.GetProducts()
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, products)
}

// GetProductById godoc
// @Summary Buscar produto por ID
// @Tags Products
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do produto"
// @Success 200 {object} model.Product
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Router /products/{id} [get]
func (p *productController) GetProductById(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		response := model.Response{
			Message: "Id do produto não pode ser nulo",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	productId, err := strconv.Atoi(id)
	if err != nil {
		response := model.Response{
			Message: "Id do produto precisa ser um número",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	product, err := p.productUsecase.GetProductById(productId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	if product == nil {
		response := model.Response{
			Message: "Produto não foi encontrado na base de dados",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary Criar produto
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body model.Product true "Produto"
// @Success 201 {object} model.Product
// @Failure 400 {object} map[string]string
// @Router /products [post]
func (p *productController) CreateProduct(ctx *gin.Context) {

	var product model.Product
	err := ctx.BindJSON(&product)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, err)
		return
	}

	insertedProduct, err := p.productUsecase.CreateProduct(product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, insertedProduct)
}

// UpdateProductById godoc
// @Summary Atualizar produto
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do produto"
// @Param product body model.Product true "Campos para atualizar"
// @Success 200 {object} model.Product
// @Failure 400 {object} model.Response
// @Router /products/{id} [put]
func (p *productController) UpdateProductById(ctx *gin.Context) {

	var product model.Product
	err := ctx.BindJSON(&product)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, err)
		return
	}
	if product.Name == nil && product.Price == nil {
		response := model.Response{
			Message: "É necessário preencher ao menos um campo para ser atualizado",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	id := ctx.Param("id")
	if id == "" {
		response := model.Response{
			Message: "Id do produto não pode ser nulo",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	productId, err := strconv.Atoi(id)
	if err != nil {
		response := model.Response{
			Message: "Id do produto precisa ser um número",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	updatedProduct, err := p.productUsecase.UpdateProductById(productId, product)
	if updatedProduct == nil {
		response := model.Response{
			Message: "Produto não foi encontrado na base de dados",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, updatedProduct)
}

// DeleteProductById godoc
// @Summary Deletar produto
// @Tags Products
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do produto"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /products/{id} [delete]
func (p *productController) DeleteProductById(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		response := model.Response{
			Message: "Id do produto não pode ser nulo",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	productId, err := strconv.Atoi(id)
	if err != nil {
		response := model.Response{
			Message: "Id do produto precisa ser um número",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	isSucess, err := p.productUsecase.DeleteProductById(productId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	if isSucess {
		response := model.Response{
			Message: "O produto foi deletado com sucesso",
		}
		ctx.JSON(http.StatusOK, response)
	}
}