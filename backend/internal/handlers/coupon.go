package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/models"
)

type CouponRequest struct {
	Code       string    `json:"code"`
	Discount   float64   `json:"discount"`
	ExpiryDate time.Time `json:"expiry_date"`
}

func CreateCoupon(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CouponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body. Ensure expiry_date is in ISO-8601 format.", http.StatusBadRequest)
		return
	}

	coupon := models.Coupon{
		Code:       strings.ToUpper(req.Code),
		Discount:   req.Discount,
		ExpiryDate: req.ExpiryDate,
		IsActive:   true,
	}

	if result := database.DB.Create(&coupon); result.Error != nil {
		http.Error(w, "Failed to create coupon (Code might already exist)", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(coupon)
}
