package http

import (
	"net/http"
)

// Ads Proxy Handlers

// GetAdsByLocation handles GET /api/v1/ads/placement/{location}
func (c *Config) GetAdsByLocation(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/ads/placement/{location}")
}

// GetAllAdsPlacements handles GET /api/v1/ads/placements
func (c *Config) GetAllAdsPlacements(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/ads/placements")
}

// CreateAdsPlacement handles POST /api/v1/ads/placements
func (c *Config) CreateAdsPlacement(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/ads/placements")
}

// GetAdsPlacementByID handles GET /api/v1/ads/placements/{id}
func (c *Config) GetAdsPlacementByID(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/ads/placements/{id}")
}

// UpdateAdsPlacement handles PUT /api/v1/ads/placements/{id}
func (c *Config) UpdateAdsPlacement(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/ads/placements/{id}")
}

// DeleteAdsPlacement handles DELETE /api/v1/ads/placements/{id}
func (c *Config) DeleteAdsPlacement(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/ads/placements/{id}")
}
