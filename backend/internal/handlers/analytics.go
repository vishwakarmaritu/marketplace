package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/models"
)

func GetMarketplaceAnalytics(c *fiber.Ctx) error {
	var totalSales float64
	var successfulOrders int64
	var failedOrders int64
	var totalStock int64
	var outOfStockProducts []models.Product

	database.DB.Model(&models.Order{}).
		Select("COALESCE(SUM(total_amount), 0)").
		Where("status = ?", "Delivered").
		Scan(&totalSales)

	database.DB.Model(&models.Order{}).
		Where("status != ?", "Failed").
		Count(&successfulOrders)

	database.DB.Model(&models.Order{}).
		Where("status = ?", "Failed").
		Count(&failedOrders)

	database.DB.Model(&models.Product{}).
		Select("COALESCE(SUM(stock), 0)").
		Scan(&totalStock)

	database.DB.Where("stock = ?", 0).Find(&outOfStockProducts)

	return c.JSON(fiber.Map{
		"sales": fiber.Map{
			"total_delivered_revenue": totalSales,
		},
		"orders": fiber.Map{
			"successful_count": successfulOrders,
			"cancelled_count":  failedOrders,
		},
		"inventory": fiber.Map{
			"total_units_in_stock": totalStock,
			"out_of_stock_count":   len(outOfStockProducts),
			"out_of_stock_items":   outOfStockProducts,
		},
	})
}
