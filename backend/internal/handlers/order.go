package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/models"
)

func Checkout(w http.ResponseWriter, r *http.Request) {

	contextUserID := r.Context().Value("userID").(float64)
	buyerID := uint(contextUserID)

	var cart models.Cart
	result := database.DB.Preload("Items.Product").Where("user_id = ?", buyerID).First(&cart)

	if result.Error != nil || len(cart.Items) == 0 {
		http.Error(w, "Cannot checkout: Cart is empty", http.StatusBadRequest)
		return
	}

	var totalAmount float64
	var orderItems []models.OrderItem

	for _, item := range cart.Items {
		totalAmount += float64(item.Quantity) * item.Product.Price

		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Product.Price,
		})
	}

	order := models.Order{
		BuyerID:     buyerID,
		TotalAmount: totalAmount,
		Items:       orderItems,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	database.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
	database.DB.Delete(&cart)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Order placed successfully",
		"order":   order,
	})
}
