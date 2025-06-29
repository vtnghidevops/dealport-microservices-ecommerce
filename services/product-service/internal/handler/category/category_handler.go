package category

import (
	"errors"
	"net/http"
	"product-service/internal/domain"
	"product-service/internal/handler"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetAll returns all categories
func GetAllCategories(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Build filters from query parameters
		filters := make(map[string]string)

		if isActive := r.URL.Query().Get("is_active"); isActive != "" {
			filters["is_active"] = isActive
		}

		if isVisible := r.URL.Query().Get("is_visible"); isVisible != "" {
			filters["is_visible"] = isVisible
		}

		// Get categories from service
		categories, err := app.CategoryService.GetAllCategories(filters)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status: http.StatusOK,
			Data:   categories,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// GetByID returns a category by ID
func GetCategoryByID(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			app.ErrorResponse(w, errors.New("missing category ID"), http.StatusBadRequest)
			return
		}

		// Convert string ID to integer
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid category ID format"), http.StatusBadRequest)
			return
		}

		category, err := app.CategoryService.GetCategoryByID(id)
		if err != nil {
			if errors.Is(err, domain.ErrCategoryNotFound) {
				app.ErrorResponse(w, err, http.StatusNotFound)
				return
			}
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status: http.StatusOK,
			Data:   category,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// GetBySlug returns a category by slug
func GetCategoryBySlug(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		if slug == "" {
			app.ErrorResponse(w, errors.New("missing category slug"), http.StatusBadRequest)
			return
		}

		category, err := app.CategoryService.GetCategoryBySlug(slug)
		if err != nil {
			if errors.Is(err, domain.ErrCategoryNotFound) {
				app.ErrorResponse(w, err, http.StatusNotFound)
				return
			}
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status: http.StatusOK,
			Data:   category,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// Create adds a new category
func CreateCategory(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var category domain.Category

		err := app.ParseJSON(w, r, &category)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		// Validate required fields
		if category.Name == "" || category.Slug == "" {
			app.ErrorResponse(w, errors.New("missing required fields"), http.StatusBadRequest)
			return
		}

		id, err := app.CategoryService.CreateCategory(&category)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		// Get the newly created category
		createdCategory, err := app.CategoryService.GetCategoryByID(id)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status:  http.StatusCreated,
			Message: "Category created successfully",
			Data:    createdCategory,
		}

		app.SuccessResponse(w, http.StatusCreated, payload)
	}
}

// Update updates an existing category
func UpdateCategory(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			app.ErrorResponse(w, errors.New("missing category ID"), http.StatusBadRequest)
			return
		}

		// Convert string ID to integer
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid category ID format"), http.StatusBadRequest)
			return
		}

		var category domain.Category

		err = app.ParseJSON(w, r, &category)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		// Ensure ID in URL matches category ID
		category.ID = id

		// Validate required fields
		if category.Name == "" || category.Slug == "" {
			app.ErrorResponse(w, errors.New("missing required fields"), http.StatusBadRequest)
			return
		}

		err = app.CategoryService.UpdateCategory(&category)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		// Get updated category
		updatedCategory, err := app.CategoryService.GetCategoryByID(id)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status:  http.StatusOK,
			Message: "Category updated successfully",
			Data:    updatedCategory,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// Delete removes a category
func DeleteCategory(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			app.ErrorResponse(w, errors.New("missing category ID"), http.StatusBadRequest)
			return
		}

		// Convert string ID to integer
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid category ID format"), http.StatusBadRequest)
			return
		}

		err = app.CategoryService.DeleteCategory(id)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status:  http.StatusOK,
			Message: "Category deleted successfully",
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// SyncProductCounts updates product counts for all categories
func SyncProductCounts(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := app.CategoryService.SyncProductCounts()
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status:  http.StatusOK,
			Message: "Product counts synchronized successfully",
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}
