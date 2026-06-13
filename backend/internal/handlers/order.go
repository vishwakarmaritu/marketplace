package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/models"
)

func Checkout(c *fiber.Ctx) error {
	buyerID := uint(c.Locals("userID").(float64))

	var cart models.Cart
	result := database.DB.Preload("Items.Product").Where("user_id = ?", buyerID).First(&cart)

	if result.Error != nil || len(cart.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot checkout: Cart is empty"})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order"})
	}

	database.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
	database.DB.Delete(&cart)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Order placed successfully",
		"order":   order,
	})
}
