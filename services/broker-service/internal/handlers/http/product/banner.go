package http

import (
	"net/http"
)

// Banner Proxy Handlers

// GetAllBanners handles GET /api/v1/banners
func (c *Config) GetAllBanners(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/banners")
}

// CreateBanner handles POST /api/v1/banners
func (c *Config) CreateBanner(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/banners")
}

// GetBannerByID handles GET /api/v1/banners/{id}
func (c *Config) GetBannerByID(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/banners/{id}")
}

// UpdateBanner handles PUT /api/v1/banners/{id}
func (c *Config) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/banners/{id}")
}

// DeleteBanner handles DELETE /api/v1/banners/{id}
func (c *Config) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/banners/{id}")
}

// GetBannersByType handles GET /api/v1/banners/type/{type}
func (c *Config) GetBannersByType(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/banners/type/{type}")
}
