package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/models"
)

type CouponRequest struct {
	Code       string    `json:"code"`
	Discount   float64   `json:"discount"`
	ExpiryDate time.Time `json:"expiry_date"`
}

func CreateCoupon(c *fiber.Ctx) error {
	var req CouponRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body. Ensure expiry_date is in ISO-8601 format."})
	}

	coupon := models.Coupon{
		Code:       strings.ToUpper(req.Code),
		Discount:   req.Discount,
		ExpiryDate: req.ExpiryDate,
		IsActive:   true,
	}

	if result := database.DB.Create(&coupon); result.Error != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Failed to create coupon (Code might already exist)"})
	}

	return c.Status(fiber.StatusCreated).JSON(coupon)
}
