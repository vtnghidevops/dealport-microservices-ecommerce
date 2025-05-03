package http

import (
	"net/http"
)

// Category Proxy Handlers

// GetAllCategories handles GET /api/v1/categories
func (c *Config) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/categories")
}

// CreateCategory handles POST /api/v1/categories
func (c *Config) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/categories")
}

// GetCategoryByID handles GET /api/v1/categories/{id}
func (c *Config) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/categories/{id}")
}

// GetCategoryBySlug handles GET /api/v1/categories/slug/{slug}
func (c *Config) GetCategoryBySlug(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/categories/slug/{slug}")
}

// UpdateCategory handles PUT /api/v1/categories/{id}
func (c *Config) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/categories/{id}")
}

// DeleteCategory handles DELETE /api/v1/categories/{id}
func (c *Config) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/categories/{id}")
}
