package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routers() http.Handler {
	mux := chi.NewRouter()

	// config cors for router
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "https://test-payment.momo.vn", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "*"},
		ExposedHeaders:   []string{"Link", "Content-Disposition"},
		AllowCredentials: true, // Allow cookies, session, token to be pass in request
		MaxAge:           300,  // Maximum value not ignored by any of major browsers
	}))

	mux.Use(chimiddleware.Heartbeat("/ping")) // Check if server is alive

	// [POST] /broker
	// mux.Post("/broker", app.Broker)

	// // [POST] /handle => authentication-service
	// mux.Post("/handle", app.HandleSubmission)

	// [POST] /logs/gRPC
	// mux.Post("/log-grpc", app.LogViaGRPC)

	// [Product-Service]
	mux.Route("/api/v1", func(r chi.Router) {
		// Authentication routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.AuthHandler.Register)
			r.Post("/login", app.AuthHandler.Login)
			r.Post("/refresh", app.AuthHandler.RefreshToken)
			r.Get("/validate", app.AuthHandler.ValidateToken)
		})

		// User routes with authentication
		r.Route("/users", func(r chi.Router) {
			// Protected routes
			r.Group(func(r chi.Router) {
				r.Use(app.AuthMiddleware.RequireAuth)
				r.Get("/me", app.UserHandler.GetUserProfile)
				r.Put("/me", app.UserHandler.UpdateUserProfile)

				// Wishlist routes
				r.Route("/me/wishlist", func(r chi.Router) {
					r.Get("/", app.UserHandler.GetWishlist)
					r.Post("/", app.UserHandler.AddToWishlist)
					r.Delete("/{product_id}", app.UserHandler.RemoveFromWishlist)
				})
			})

			// Admin routes
			r.Group(func(r chi.Router) {
				r.Use(app.AuthMiddleware.RequireAuth)
				// TODO: Add admin role check middleware
				r.Get("/{id}", app.UserHandler.GetUser)
			})
		})

		// Cart routes with authentication
		r.Route("/cart", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(app.AuthMiddleware.RequireAuth)
				r.Get("/", app.CartHandler.GetCart)
				r.Post("/items", app.CartHandler.AddCartItem)
				r.Put("/items/{item_id}", app.CartHandler.UpdateCartItem)
				r.Delete("/items/{item_id}", app.CartHandler.RemoveCartItem)
				r.Delete("/", app.CartHandler.ClearCart)
				r.Post("/coupon", app.CartHandler.ApplyCoupon)
				r.Delete("/coupon", app.CartHandler.RemoveCoupon)
			})
		})

		// Checkout routes with authentication
		r.Route("/checkout", func(r chi.Router) {
			// Validate checkout without auth
			r.Post("/validate", app.CheckoutHandler.ValidateCheckout)

			// Protected checkout routes
			r.Group(func(r chi.Router) {
				r.Use(app.AuthMiddleware.RequireAuth)
				r.Post("/orders", app.CheckoutHandler.CreateOrder)
				r.Get("/orders/{id}", app.CheckoutHandler.GetOrder)
				r.Get("/orders", app.CheckoutHandler.ListOrders)
				r.Post("/orders/{id}/payment", app.CheckoutHandler.ProcessPayment)
			})
		})

		// Payment routes
		r.Route("/payments", func(r chi.Router) {
			r.Post("/momo/create", app.PaymentHandler.CreateMomoPayment)
			r.Post("/momo/verify", app.PaymentHandler.VerifyMomoPayment)
			// r.Post("/vnpay/create", app.PaymentHandler.CreateVnpayPayment)
			// r.Post("/vnpay/verify", app.PaymentHandler.VerifyVnpayPayment)
			// MoMo QuickPay endpoints
			r.Post("/momo/qr/create", app.PaymentHandler.CreateMomoQRPayment)
			r.Post("/momo/pos/create", app.PaymentHandler.CreateMomoPosPayment)
		})

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
