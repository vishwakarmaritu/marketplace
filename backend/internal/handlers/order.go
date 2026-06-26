package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/models"
)

type CheckoutRequest struct {
	CouponCode  string `json:"coupon_code"`
	PaymentCard string `json:"payment_card"`
}

func processMockPayment(cardNumber string) bool {
	return cardNumber == "1234-mock"
}

func Checkout(c *fiber.Ctx) error {
	buyerID := uint(c.Locals("userID").(float64))

	var req CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var cart models.Cart
	result := database.DB.Preload("Items.Product").Where("user_id = ?", buyerID).First(&cart)
	if result.Error != nil || len(cart.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot checkout: Cart is empty"})
	}

	var subtotal float64
	var orderItems []models.OrderItem

	for _, item := range cart.Items {
		subtotal += float64(item.Quantity) * item.Product.Price
		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Product.Price,
		})
	}

	finalTotal := subtotal

	if req.CouponCode != "" {
		var coupon models.Coupon
		if err := database.DB.Where("code = ? AND is_active = ?", strings.ToUpper(req.CouponCode), true).First(&coupon).Error; err == nil {

			if coupon.ExpiryDate.After(time.Now()) {

				discountAmount := (subtotal * coupon.Discount) / 100
				finalTotal = subtotal - discountAmount
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Coupon has expired"})
			}
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid coupon code"})
		}
	}

	paymentSuccessful := processMockPayment(req.PaymentCard)
	orderStatus := "Placed"
	if !paymentSuccessful {
		orderStatus = "Failed"
	}

	order := models.Order{
		BuyerID:     buyerID,
		Status:      orderStatus,
		TotalAmount: finalTotal,
		Items:       orderItems,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order"})
	}

	if !paymentSuccessful {
		return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
			"error": "Payment failed. Invalid card number.",
			"order": order,
		})
	}

	database.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
	database.DB.Delete(&cart)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Order placed successfully and payment approved!",
		"order":   order,
	})
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

func GetOrderTracking(c *fiber.Ctx) error {
	orderID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	userID := uint(c.Locals("userID").(float64))
	userRole := c.Locals("userRole").(string)

	var order models.Order
	if err := database.DB.Preload("Items.Product").First(&order, orderID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	if userRole == "buyer" && order.BuyerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized to view this order"})
	}

	return c.JSON(order)
}

func UpdateOrderStatus(c *fiber.Ctx) error {
	orderID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var req UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	validStatuses := map[string]bool{
		"Placed":    true,
		"Packed":    true,
		"Shipped":   true,
		"Delivered": true,
		"Failed":    true,
	}

	if !validStatuses[req.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid status string. Use Packed, Shipped, or Delivered."})
	}

	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	order.Status = req.Status
	database.DB.Save(&order)

	return c.JSON(fiber.Map{
		"message": "Order status updated successfully",
		"order":   order,
	})
}
