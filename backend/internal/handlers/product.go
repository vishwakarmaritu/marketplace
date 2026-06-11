package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/models"
)

type ProductRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {

	contextUserID := r.Context().Value("userID").(float64)
	sellerID := uint(contextUserID)

	var req ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product := models.Product{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SellerID:    sellerID,
	}

	if result := database.DB.Create(&product); result.Error != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product

	searchQuery := r.URL.Query().Get("search")

	if searchQuery != "" {

		database.DB.Where("title ILIKE ?", "%"+searchQuery+"%").Find(&products)
	} else {

		database.DB.Find(&products)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func ManageProduct(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	productID, err := strconv.Atoi(idStr)
	if err != nil || productID == 0 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	if result := database.DB.First(&product, productID); result.Error != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodDelete {
		database.DB.Delete(&product)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
		return
	}

	if r.Method == http.MethodPut {
		var req ProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		product.Title = req.Title
		product.Description = req.Description
		product.Price = req.Price
		product.Stock = req.Stock

		database.DB.Save(&product)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(product)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
