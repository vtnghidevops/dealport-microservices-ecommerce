package http

import (
	"net/http"
	"product-service/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Server represents the HTTP server
type Server struct {
	Handler *handler.Config
}

// NewServer creates a new HTTP server
func NewServer(handler *handler.Config) *Server {
	return &Server{
		Handler: handler,
	}
}

// Routes returns the HTTP routes
func (s *Server) Routes() http.Handler {
	router := chi.NewRouter()

	// Middleware (phương tiện trung gian)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	// CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check route
	router.Get("/health", s.Handler.HealthCheck)

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		// Products
		r.Get("/products", s.Handler.GetAllProducts)
		r.Get("/products/{id}", s.Handler.GetProductByID)
		r.Get("/products/slug/{slug}", s.Handler.GetProductBySlug)
		r.Post("/products", s.Handler.CreateProduct)
		r.Put("/products/{id}", s.Handler.UpdateProduct)
		r.Patch("/products/{id}", s.Handler.PatchProduct)
		r.Delete("/products/{id}", s.Handler.DeleteProduct)

		// Product reviews
		r.Get("/products/{id}/reviews", s.Handler.GetProductReviews)
		r.Post("/products/{id}/reviews", s.Handler.AddProductReview)

		// Categories
		r.Get("/categories", s.Handler.GetAllCategories)
		r.Get("/categories/{id}", s.Handler.GetCategoryByID)
		r.Get("/categories/slug/{slug}", s.Handler.GetCategoryBySlug)
		r.Get("/categories/tree", s.Handler.GetCategoryTree)
	})

	return router
}
