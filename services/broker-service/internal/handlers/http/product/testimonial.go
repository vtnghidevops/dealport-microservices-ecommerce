package http

import (
	"net/http"
)

// Testimonial Proxy Handlers

// GetTopRatedTestimonials handles GET /api/v1/testimonials
func (c *Config) GetTopRatedTestimonials(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/testimonials")
}
