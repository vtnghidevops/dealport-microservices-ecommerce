package http

import (
	"net/http"
	"product-service/internal/handler"
	"product-service/internal/handler/ads"
	"product-service/internal/handler/banner"
	"product-service/internal/handler/category"
	"product-service/internal/handler/products"

	"time"

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

	// Middleware
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not cluttering logs
	}))

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		// Products
		r.Route("/products", func(r chi.Router) {
			r.Get("/", products.GetAll(s.Handler))
			r.Post("/", products.Create(s.Handler))
			r.Get("/{id}", products.GetByID(s.Handler))
			r.Get("/slug/{slug}", products.GetBySlug(s.Handler))
			r.Put("/{id}", products.Update(s.Handler))
			r.Delete("/{id}", products.Delete(s.Handler))

			// Product reviews
			r.Get("/{id}/reviews", products.GetReviews(s.Handler))
			r.Post("/{id}/reviews", products.AddReview(s.Handler))
			r.Put("/{id}/reviews/{reviewId}", products.UpdateReview(s.Handler))
		})

		// Testimonials endpoint for HappyCustomers section
		r.Get("/testimonials", products.GetTopRatedTestimonials(s.Handler))

		// Categories
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", category.GetAllCategories(s.Handler))
			r.Post("/", category.CreateCategory(s.Handler))
			r.Get("/{id}", category.GetCategoryByID(s.Handler))
			r.Get("/slug/{slug}", category.GetCategoryBySlug(s.Handler))
			r.Put("/{id}", category.UpdateCategory(s.Handler))
			r.Delete("/{id}", category.DeleteCategory(s.Handler))
		})

		// Banners
		r.Route("/banners", func(r chi.Router) {
			r.Get("/", banner.GetAllBanners(s.Handler))
			r.Post("/", banner.CreateBanner(s.Handler))
			r.Get("/{id}", banner.GetBannerByID(s.Handler))
			r.Put("/{id}", banner.UpdateBanner(s.Handler))
			r.Delete("/{id}", banner.DeleteBanner(s.Handler))
			r.Get("/type/{type}", banner.GetBannersByType(s.Handler))
		})

		// Ads
		r.Route("/ads", func(r chi.Router) {
			// Get ads by location for frontend display
			r.Get("/placement/{location}", ads.GetAdsByLocation(s.Handler))

			// CRUD operations on ad placements
			r.Get("/placements", ads.GetAllAdsPlacements(s.Handler))
			r.Post("/placements", ads.CreateAdsPlacement(s.Handler))
			r.Get("/placements/{id}", ads.GetAdsPlacementByID(s.Handler))
			r.Put("/placements/{id}", ads.UpdateAdsPlacement(s.Handler))
			r.Delete("/placements/{id}", ads.DeleteAdsPlacement(s.Handler))
		})
	})

	return router
}
