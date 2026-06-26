package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/models"
)

type DisputeRequest struct {
	Reason string `json:"reason"`
}

type ResolveDisputeRequest struct {
	Status        string `json:"status"`
	AdminResponse string `json:"admin_response"`
}

func RaiseDispute(c *fiber.Ctx) error {
	orderID, err := c.ParamsInt("id")
	if err != nil || orderID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	buyerID := uint(c.Locals("userID").(float64))

	var req DisputeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}
	if order.BuyerID != buyerID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized to dispute this order"})
	}

	dispute := models.Dispute{
		OrderID: uint(orderID),
		BuyerID: buyerID,
		Reason:  req.Reason,
	}

	if err := database.DB.Create(&dispute).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create dispute"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Dispute raised successfully",
		"dispute": dispute,
	})
}

func ResolveDispute(c *fiber.Ctx) error {
	disputeID, err := c.ParamsInt("id")
	if err != nil || disputeID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid dispute ID"})
	}

	var req ResolveDisputeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var dispute models.Dispute
	if err := database.DB.First(&dispute, disputeID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Dispute not found"})
	}

	dispute.Status = req.Status
	dispute.AdminResponse = req.AdminResponse
	database.DB.Save(&dispute)

	return c.JSON(fiber.Map{
		"message": "Dispute updated successfully",
		"dispute": dispute,
	})
}
