package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vishwakarmaritu/marketplace/internal/handlers"
	"github.com/vishwakarmaritu/marketplace/internal/middleware"
)

func RegisterRoutes(app *fiber.App) {

	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	api.Post("/auth/signup", handlers.Signup)
	api.Post("/auth/login", handlers.Login)
	api.Get("/products", handlers.GetProducts)

	sellerRoutes := api.Group("/", middleware.RequireRole("seller", "admin"))
	sellerRoutes.Post("/products", handlers.CreateProduct)
	sellerRoutes.Put("/products/:id", handlers.ManageProduct)
	sellerRoutes.Delete("/products/:id", handlers.ManageProduct)

	api.Post("/coupons", middleware.RequireRole("seller"), handlers.CreateCoupon)

	buyerRoutes := api.Group("/", middleware.RequireRole("buyer"))
	buyerRoutes.Get("/cart", handlers.GetCart)
	buyerRoutes.Post("/cart", handlers.AddToCart)
	buyerRoutes.Post("/checkout", handlers.Checkout)
}
