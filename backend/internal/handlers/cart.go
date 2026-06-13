package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/models"
)

type CartRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

func AddToCart(c *fiber.Ctx) error {
	buyerID := uint(c.Locals("userID").(float64))

	var req CartRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
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

	return c.JSON(fiber.Map{"message": "Item successfully added to cart"})
}

func GetCart(c *fiber.Ctx) error {
	buyerID := uint(c.Locals("userID").(float64))

	var cart models.Cart
	result := database.DB.Preload("Items.Product").Where("user_id = ?", buyerID).First(&cart)

	if result.Error != nil {
		return c.JSON(fiber.Map{"message": "Cart is currently empty"})
	}

	var totalPrice float64
	for _, item := range cart.Items {
		totalPrice += float64(item.Quantity) * item.Product.Price
	}

	return c.JSON(fiber.Map{
		"cart_id":     cart.ID,
		"items":       cart.Items,
		"total_price": totalPrice,
	})
}
