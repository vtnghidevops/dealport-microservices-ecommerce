package domain

import "errors"

// Domain error definitions
var (
	// Common errors
	ErrNotFound       = errors.New("resource not found")
	ErrInternalServer = errors.New("internal server error")
	ErrInvalidRequest = errors.New("invalid request")

	// Product specific errors
	ErrInvalidProduct   = errors.New("invalid product data")
	ErrProductNotFound  = errors.New("product not found")
	ErrDuplicateProduct = errors.New("product already exists")

	// Category specific errors
	ErrInvalidCategory   = errors.New("invalid category data")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrDuplicateCategory = errors.New("category already exists")
	ErrCategoryInUse     = errors.New("category is in use by products and cannot be deleted")

	// Review specific errors
	ErrInvalidReview      = errors.New("invalid review data")
	ErrReviewNotFound     = errors.New("review not found")
	ErrUnauthorizedReview = errors.New("unauthorized to modify this review")

	// Banner errors
	ErrBannerNotFound = errors.New("banner not found")
	ErrInvalidBanner  = errors.New("invalid banner data")
)
