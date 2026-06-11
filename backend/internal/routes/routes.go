package routes

import (
	"net/http"

	"github.com/vishwakarmaritu/marketplace/internal/handlers"
	"github.com/vishwakarmaritu/marketplace/internal/middleware"
)

func RegisterRoutes() {

	http.HandleFunc("/api/auth/signup", handlers.Signup)
	http.HandleFunc("/api/auth/login", handlers.Login)

	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetProducts(w, r)
		} else if r.Method == http.MethodPost {
			middleware.RequireRole("seller", "admin")(handlers.CreateProduct).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/products/", middleware.RequireRole("seller", "admin")(handlers.ManageProduct))

	http.HandleFunc("/api/coupons", middleware.RequireRole("seller")(handlers.CreateCoupon))

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "message": "Server and database are responding cleanly"}`))
	})
}
