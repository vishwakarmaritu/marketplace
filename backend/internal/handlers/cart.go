package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/models"
)

type CartRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

func AddToCart(w http.ResponseWriter, r *http.Request) {

	contextUserID := r.Context().Value("userID").(float64)
	buyerID := uint(contextUserID)

	var req CartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var cart models.Cart
	database.DB.FirstOrCreate(&cart, models.Cart{UserID: buyerID})

	var cartItem models.CartItem
	result := database.DB.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem)

	if result.Error == nil {

		cartItem.Quantity += req.Quantity
		database.DB.Save(&cartItem)
	} else {

		cartItem = models.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		database.DB.Create(&cartItem)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item successfully added to cart"})
}

func GetCart(w http.ResponseWriter, r *http.Request) {
	contextUserID := r.Context().Value("userID").(float64)
	buyerID := uint(contextUserID)

	var cart models.Cart

	result := database.DB.Preload("Items.Product").Where("user_id = ?", buyerID).First(&cart)

	if result.Error != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Cart is currently empty"})
		return
	}

	var totalPrice float64
	for _, item := range cart.Items {
		totalPrice += float64(item.Quantity) * item.Product.Price
	}

	response := map[string]interface{}{
		"cart_id":     cart.ID,
		"items":       cart.Items,
		"total_price": totalPrice,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
