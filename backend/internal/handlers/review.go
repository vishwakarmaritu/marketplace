package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/models"
)

type ReviewRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

func CreateReview(c *fiber.Ctx) error {

	productID, err := c.ParamsInt("id")
	if err != nil || productID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	buyerID := uint(c.Locals("userID").(float64))

	var req ReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Rating < 1 || req.Rating > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Rating must be between 1 and 5 stars"})
	}

	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	review := models.Review{
		ProductID: uint(productID),
		BuyerID:   buyerID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	if err := database.DB.Create(&review).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save review"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Review added successfully",
		"review":  review,
	})
}
