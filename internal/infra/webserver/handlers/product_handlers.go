package handlers

import (
	"encoding/json"
	"github.com/andre2ar/go-products/internal/dto"
	"github.com/andre2ar/go-products/internal/entity"
	"github.com/andre2ar/go-products/internal/infra/database"
	"net/http"
)

type ProductHandler struct {
	ProductRepository database.ProductRepositoryInterface
}

func NewProductHandler(db database.ProductRepositoryInterface) *ProductHandler {
	return &ProductHandler{ProductRepository: db}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newProduct, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ProductRepository.Create(newProduct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
