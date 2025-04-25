package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routers() http.Handler {
	mux := chi.NewRouter()

	// config cors for router
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, // Allow cookies, session, token to be pass in request
		MaxAge:           300,  // Maximum value not ignored by any of major browsers
	}))

	mux.Use(middleware.Heartbeat("/ping")) // Check if server is alive

	// [POST] /broker
	// mux.Post("/broker", app.Broker)

	// // [POST] /handle => authentication-service
	// mux.Post("/handle", app.HandleSubmission)

	// [POST] /logs/gRPC
	// mux.Post("/log-grpc", app.LogViaGRPC)

	// [Product-Service]
	mux.Route("/api/v1", func(r chi.Router) {
		// Products
		r.Route("/health", func(r chi.Router) {
			r.Get("/", app.ProductHandler.GetProductHealth)
		})
		r.Route("/products", func(r chi.Router) {
			r.Get("/", app.ProductHandler.GetAllProducts)
			r.Post("/", app.ProductHandler.CreateProduct)
			r.Get("/{id}", app.ProductHandler.GetProductByID)
			r.Get("/slug/{slug}", app.ProductHandler.GetProductBySlug)
			r.Put("/{id}", app.ProductHandler.UpdateProduct)
			// r.Patch("/{id}", app.ProductHandler.PatchProduct)
			r.Patch("/{id}", app.httpProductHandler.PatchProduct) // using http handler to patch product for flexibilty
			r.Delete("/{id}", app.ProductHandler.DeleteProduct)

			// Product reviews
			r.Get("/{id}/reviews", app.ProductHandler.GetProductReviews)
			r.Post("/{id}/reviews", app.ProductHandler.AddProductReview)
			r.Put("/{id}/reviews/{reviewId}", app.ProductHandler.UpdateProductReview)

			// Product images
			r.Post("/{id}/images", app.ProductHandler.UploadProductImage)
			r.Delete("/{id}/images/{imageId}", app.ProductHandler.DeleteProductImage)
			r.Put("/{id}/images/{imageId}/primary", app.ProductHandler.SetPrimaryProductImage)
		})

		// Testimonials endpoint for HappyCustomers section
		r.Get("/testimonials", app.ProductHandler.GetTopRatedTestimonials)

		// Categories
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", app.CategoryHandler.GetAllCategories)
			r.Post("/", app.CategoryHandler.CreateCategory)
			r.Get("/{id}", app.CategoryHandler.GetCategoryByID)
			r.Get("/slug/{slug}", app.CategoryHandler.GetCategoryBySlug)
			r.Put("/{id}", app.CategoryHandler.UpdateCategory)
			r.Delete("/{id}", app.CategoryHandler.DeleteCategory)
		})

		// Banners
		r.Route("/banners", func(r chi.Router) {
			r.Get("/", app.BannerHandler.GetAllBanners)
			r.Post("/", app.BannerHandler.CreateBanner)
			r.Get("/{id}", app.BannerHandler.GetBannerByID)
			r.Put("/{id}", app.BannerHandler.UpdateBanner)
			r.Delete("/{id}", app.BannerHandler.DeleteBanner)
			r.Get("/type/{type}", app.BannerHandler.GetBannersByType)
		})

		// Ads
		r.Route("/ads", func(r chi.Router) {
			// Get ads by location for frontend display
			r.Get("/placement/{location}", app.AdsHandler.GetAdsByLocation)

			// CRUD operations on ad placements
			r.Get("/placements", app.AdsHandler.GetAllAdsPlacements)
			r.Post("/placements", app.AdsHandler.CreateAdsPlacement)
			r.Get("/placements/{id}", app.AdsHandler.GetAdsPlacementByID)
			r.Put("/placements/{id}", app.AdsHandler.UpdateAdsPlacement)
			r.Delete("/placements/{id}", app.AdsHandler.DeleteAdsPlacement)
		})
	})

	// Handle static files
	mux.Handle("/images/*", http.StripPrefix("/images/", http.FileServer(http.Dir("./static/images"))))
	// Route for product images - proxy directly to product service (trực tiếp)
	mux.HandleFunc("/api/products/images/{filename}", app.httpProductHandler.GetProductImage)

	return mux
}
