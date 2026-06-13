package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/models"
)

type ProductRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func CreateProduct(c *fiber.Ctx) error {
	sellerID := uint(c.Locals("userID").(float64))

	var req ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	product := models.Product{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SellerID:    sellerID,
	}

	if result := database.DB.Create(&product); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func GetProducts(c *fiber.Ctx) error {
	var products []models.Product
	searchQuery := c.Query("search")

	if searchQuery != "" {
		database.DB.Where("title ILIKE ?", "%"+searchQuery+"%").Find(&products)
	} else {
		database.DB.Find(&products)
	}

	return c.JSON(products)
}

func ManageProduct(c *fiber.Ctx) error {
	productID, err := c.ParamsInt("id")
	if err != nil || productID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	var product models.Product
	if result := database.DB.First(&product, productID); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	if c.Method() == fiber.MethodDelete {
		database.DB.Delete(&product)
		return c.JSON(fiber.Map{"message": "Product deleted successfully"})
	}

	if c.Method() == fiber.MethodPut {
		var req ProductRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		product.Title = req.Title
		product.Description = req.Description
		product.Price = req.Price
		product.Stock = req.Stock

		database.DB.Save(&product)
		return c.JSON(product)
	}

	return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{"error": "Method not allowed"})
}
