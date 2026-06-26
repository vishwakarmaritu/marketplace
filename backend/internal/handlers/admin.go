package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/models"
)

func GetAllDisputes(c *fiber.Ctx) error {
	var disputes []models.Dispute

	if err := database.DB.Where("status = ?", "Open").Find(&disputes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch disputes"})
	}

	return c.JSON(fiber.Map{
		"active_disputes": disputes,
	})
}

func AdminResolveDispute(c *fiber.Ctx) error {
	disputeID, err := c.ParamsInt("id")
	if err != nil || disputeID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid dispute ID"})
	}

	var dispute models.Dispute
	if err := database.DB.First(&dispute, disputeID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Dispute not found"})
	}

	dispute.Status = "Resolved"
	dispute.AdminResponse = "Resolved by System Administrator."

	database.DB.Save(&dispute)

	return c.JSON(fiber.Map{
		"message": "Dispute successfully marked as resolved",
		"dispute": dispute,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	targetUserID, err := c.ParamsInt("id")
	if err != nil || targetUserID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var user models.User
	if err := database.DB.First(&user, targetUserID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.JSON(fiber.Map{
		"message": "User successfully banned and deleted from the system",
	})
}
