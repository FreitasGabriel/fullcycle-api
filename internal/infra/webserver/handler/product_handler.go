package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FreitasGabriel/fullcycle-api/internal/dto"
	"github.com/FreitasGabriel/fullcycle-api/internal/entity"
	"github.com/FreitasGabriel/fullcycle-api/internal/infra/database"
	entityPKG "github.com/FreitasGabriel/fullcycle-api/pkg/entity"
	"github.com/go-chi/chi"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(productDB database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: productDB,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ph *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, err := entityPKG.ParseID(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := ph.ProductDB.FindById(id)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (ph *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}

	products, err := ph.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)

}

func (ph *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		fmt.Println("error to decode body request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product.ID, err = entityPKG.ParseID(id)
	if err != nil {
		fmt.Println("error to parse id", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = ph.ProductDB.FindById(id)
	if err != nil {
		fmt.Println("error to find product", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = ph.ProductDB.Update(&product)
	if err != nil {
		fmt.Println("error to update product", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ph *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := ph.ProductDB.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = ph.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
