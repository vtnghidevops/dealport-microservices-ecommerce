// internal/handler/banner/banner_handler.go
package banner

import (
	"errors"
	"net/http"
	"product-service/internal/domain"
	"product-service/internal/handler"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetByType returns banners by type
func GetBannersByType(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bannerType := chi.URLParam(r, "type")
		
		limitStr := r.URL.Query().Get("limit")
		limit := 0
		if limitStr != "" {
			var err error
			limit, err = strconv.Atoi(limitStr)
			if err != nil || limit < 0 {
				limit = 0
			}
		}

		// Get banners from service
		banners, err := app.BannerService.GetBannersByType(bannerType, limit)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status: http.StatusOK,
			Data:   banners,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// GetAll returns all banners with pagination and filtering
func GetAllBanners(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get page and page_size from query parameters
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("page_size")

		page := 1
		pageSize := 10

		if pageStr != "" {
			pageVal, err := strconv.Atoi(pageStr)
			if err == nil && pageVal > 0 {
				page = pageVal
			}
		}

		if pageSizeStr != "" {
			pageSizeVal, err := strconv.Atoi(pageSizeStr)
			if err == nil && pageSizeVal > 0 {
				pageSize = pageSizeVal
			}
		}

		// Build filters from query parameters
		filters := make(map[string]string)

		if bannerType := r.URL.Query().Get("type"); bannerType != "" {
			filters["type"] = bannerType
		}

		if isActive := r.URL.Query().Get("is_active"); isActive != "" {
			filters["is_active"] = isActive
		}

		if productID := r.URL.Query().Get("product_id"); productID != "" {
			filters["product_id"] = productID
		}

		if categoryID := r.URL.Query().Get("category_id"); categoryID != "" {
			filters["category_id"] = categoryID
		}

		// Get banners from service
		banners, totalCount, err := app.BannerService.GetAllBanners(page, pageSize, filters)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		// Prepare pagination metadata
		meta := handler.PaginationMeta{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalItems:  totalCount,
			TotalPages:  (totalCount + pageSize - 1) / pageSize,
		}

		payload := handler.Response{
			Status: http.StatusOK,
			Data:   banners,
			Meta:   meta,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// GetByID returns a banner by ID
func GetBannerByID(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			app.ErrorResponse(w, errors.New("missing banner ID"), http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid banner ID format"), http.StatusBadRequest)
			return
		}

		banner, err := app.BannerService.GetBannerByID(id)
		if err != nil {
			if errors.Is(err, domain.ErrBannerNotFound) {
				app.ErrorResponse(w, err, http.StatusNotFound)
			} else {
				app.ErrorResponse(w, err, http.StatusInternalServerError)
			}
			return
		}

		payload := handler.Response{
			Status: http.StatusOK,
			Data:   banner,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// Create creates a new banner
func CreateBanner(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var banner domain.Banner

		err := app.ParseJSON(w, r, &banner)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		// Create banner
		id, err := app.BannerService.CreateBanner(&banner)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		// Get the newly created banner
		createdBanner, err := app.BannerService.GetBannerByID(id)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status:  http.StatusCreated,
			Message: "Banner created successfully",
			Data:    createdBanner,
		}

		app.SuccessResponse(w, http.StatusCreated, payload)
	}
}

// Update updates an existing banner
func UpdateBanner(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			app.ErrorResponse(w, errors.New("missing banner ID"), http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid banner ID format"), http.StatusBadRequest)
			return
		}

		// Check if banner exists
		existing, err := app.BannerService.GetBannerByID(id)
		if err != nil {
			if errors.Is(err, domain.ErrBannerNotFound) {
				app.ErrorResponse(w, err, http.StatusNotFound)
			} else {
				app.ErrorResponse(w, err, http.StatusInternalServerError)
			}
			return
		}

		// Parse updated banner data
		var banner domain.Banner
		err = app.ParseJSON(w, r, &banner)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		// Ensure ID matches
		banner.ID = existing.ID

		// Update banner
		err = app.BannerService.UpdateBanner(&banner)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		// Get updated banner
		updatedBanner, err := app.BannerService.GetBannerByID(id)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status:  http.StatusOK,
			Message: "Banner updated successfully",
			Data:    updatedBanner,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// Delete deletes a banner by ID
func DeleteBanner(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			app.ErrorResponse(w, errors.New("missing banner ID"), http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid banner ID format"), http.StatusBadRequest)
			return
		}

		// Check if banner exists
		_, err = app.BannerService.GetBannerByID(id)
		if err != nil {
			if errors.Is(err, domain.ErrBannerNotFound) {
				app.ErrorResponse(w, err, http.StatusNotFound)
			} else {
				app.ErrorResponse(w, err, http.StatusInternalServerError)
			}
			return
		}

		// Delete banner
		err = app.BannerService.DeleteBanner(id)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status:  http.StatusOK,
			Message: "Banner deleted successfully",
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}