package ads

import (
	"errors"
	"net/http"
	"product-service/internal/domain"
	"product-service/internal/handler"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetAdsByLocation returns ads for a specific location
func GetAdsByLocation(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get location from URL parameter
		location := chi.URLParam(r, "location")
		if location == "" {
			app.ErrorResponse(w, errors.New("missing location parameter"), http.StatusBadRequest)
			return
		}

		// Validate location value
		validLocations := map[string]bool{
			"banner":      true,
			"display":     true,
			"gaming":      true,
			"new_fashion": true,
		}

		if !validLocations[location] {
			app.ErrorResponse(w, errors.New("invalid location parameter"), http.StatusBadRequest)
			return
		}

		ads, err := app.AdsService.GetAdsByLocation(location)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		// Post-process ads for specific locations if needed
		if location == "new_fashion" && len(ads) > 0 {
			// For new_fashion, we only return the first item in the frontend
			payload := handler.Response{
				Status:  http.StatusOK,
				Message: "Ads retrieved successfully",
				Data:    ads[0:1], // Return only the first item
			}
			app.SuccessResponse(w, http.StatusOK, payload)
			return
		}

		payload := handler.Response{
			Status:  http.StatusOK,
			Message: "Ads retrieved successfully",
			Data:    ads,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// CreateAdsPlacement creates a new ads placement
func CreateAdsPlacement(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ad domain.AdsPlacement

		err := app.ParseJSON(w, r, &ad)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		// Validate required fields
		if ad.Location == "" || ad.ReferenceType == "" || ad.ReferenceID <= 0 {
			app.ErrorResponse(w, errors.New("missing required fields"), http.StatusBadRequest)
			return
		}

		// Validate reference type
		if ad.ReferenceType != "product" && ad.ReferenceType != "category" {
			app.ErrorResponse(w, errors.New("reference_type must be 'product' or 'category'"), http.StatusBadRequest)
			return
		}

		// Set default values
		if ad.DisplayOrder <= 0 {
			ad.DisplayOrder = 999 // Set a high value to place it at the end by default
		}
		ad.IsActive = true // Default to active

		id, err := app.AdsService.CreateAdsPlacement(&ad)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		// Get the created ad
		createdAd, err := app.AdsService.GetAdsPlacementByID(id)
		if err != nil {
			// Return just the ID if we can't fetch the complete ad
			payload := handler.Response{
				Status:  http.StatusCreated,
				Message: "Ad placement created successfully",
				Data: map[string]interface{}{
					"id": id,
				},
			}
			app.SuccessResponse(w, http.StatusCreated, payload)
			return
		}

		payload := handler.Response{
			Status:  http.StatusCreated,
			Message: "Ad placement created successfully",
			Data:    createdAd,
		}

		app.SuccessResponse(w, http.StatusCreated, payload)
	}
}

// GetAdsPlacementByID returns an ads placement by ID
func GetAdsPlacementByID(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from URL parameter
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			app.ErrorResponse(w, errors.New("missing ID parameter"), http.StatusBadRequest)
			return
		}

		// Convert string ID to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid ID format"), http.StatusBadRequest)
			return
		}

		ad, err := app.AdsService.GetAdsPlacementByID(id)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusNotFound)
			return
		}

		payload := handler.Response{
			Status:  http.StatusOK,
			Message: "Ad placement retrieved successfully",
			Data:    ad,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// UpdateAdsPlacement updates an existing ads placement
func UpdateAdsPlacement(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from URL parameter
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			app.ErrorResponse(w, errors.New("missing ID parameter"), http.StatusBadRequest)
			return
		}

		// Convert string ID to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid ID format"), http.StatusBadRequest)
			return
		}

		// Check if ad exists
		_, err = app.AdsService.GetAdsPlacementByID(id)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusNotFound)
			return
		}

		// Parse update data
		var ad domain.AdsPlacement
		err = app.ParseJSON(w, r, &ad)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		// Set ID from URL
		ad.ID = id

		// Update the ad
		err = app.AdsService.UpdateAdsPlacement(&ad)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		// Get the updated ad
		updatedAd, err := app.AdsService.GetAdsPlacementByID(id)
		if err != nil {
			// Return success message if we can't fetch the complete ad
			payload := handler.Response{
				Status:  http.StatusOK,
				Message: "Ad placement updated successfully",
			}
			app.SuccessResponse(w, http.StatusOK, payload)
			return
		}

		payload := handler.Response{
			Status:  http.StatusOK,
			Message: "Ad placement updated successfully",
			Data:    updatedAd,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// DeleteAdsPlacement deletes an ads placement by ID
func DeleteAdsPlacement(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from URL parameter
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			app.ErrorResponse(w, errors.New("missing ID parameter"), http.StatusBadRequest)
			return
		}

		// Convert string ID to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid ID format"), http.StatusBadRequest)
			return
		}

		err = app.AdsService.DeleteAdsPlacement(id)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status:  http.StatusOK,
			Message: "Ad placement deleted successfully",
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// GetAllAdsPlacements returns all ads placements with pagination
func GetAllAdsPlacements(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get pagination parameters from query params
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("page_size")

		page := 1      // Default page
		pageSize := 10 // Default page size

		// Parse page parameter if provided
		if pageStr != "" {
			p, err := strconv.Atoi(pageStr)
			if err == nil && p > 0 {
				page = p
			}
		}

		// Parse page_size parameter if provided
		if pageSizeStr != "" {
			ps, err := strconv.Atoi(pageSizeStr)
			if err == nil && ps > 0 {
				pageSize = ps
			}
		}

		ads, total, err := app.AdsService.GetAllAdsPlacements(page, pageSize)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status:  http.StatusOK,
			Message: "Ad placements retrieved successfully",
			Data: map[string]interface{}{
				"items":       ads,
				"total":       total,
				"page":        page,
				"page_size":   pageSize,
				"total_pages": (total + pageSize - 1) / pageSize,
			},
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}
