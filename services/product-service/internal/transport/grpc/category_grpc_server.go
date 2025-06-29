// services/product-service/internal/transport/grpc/category_server.go
package grpc

import (
	"context"
	"fmt"
	"log"
	"product-service/internal/domain"
	pb "product-service/proto/product"
	"time"
)

// CategoryGrpcServer represents the gRPC server for category service
type CategoryGrpcServer struct {
	pb.UnimplementedCategoryServiceServer
	categoryService domain.CategoryService
}

// NewCategoryGrpcServer creates a new gRPC server with the provided category service
func NewCategoryGrpcServer(categorySvc domain.CategoryService) *CategoryGrpcServer {
	log.Println("Creating new gRPC server for category service")
	return &CategoryGrpcServer{
		categoryService: categorySvc,
	}
}

// GetCategoryByID implements the GetCategoryByID RPC method
func (s *CategoryGrpcServer) GetCategoryByID(ctx context.Context, req *pb.GetCategoryByIDRequest) (*pb.Category, error) {
	category, err := s.categoryService.GetCategoryByID(int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return convertDomainCategoryToProto(category), nil
}

// GetCategoryBySlug implements the GetCategoryBySlug RPC method
func (s *CategoryGrpcServer) GetCategoryBySlug(ctx context.Context, req *pb.GetCategoryBySlugRequest) (*pb.Category, error) {
	category, err := s.categoryService.GetCategoryBySlug(req.Slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}

	return convertDomainCategoryToProto(category), nil
}

// ListCategories implements the ListCategories RPC method
func (s *CategoryGrpcServer) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	filters := make(map[string]string)
	for k, v := range req.Filters {
		filters[k] = v
	}

	categories, err := s.categoryService.GetAllCategories(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	protoCategories := make([]*pb.Category, 0, len(categories))
	for _, category := range categories {
		protoCategories = append(protoCategories, convertDomainCategoryToProto(category))
	}

	return &pb.ListCategoriesResponse{
		Categories: protoCategories,
	}, nil
}

// CreateCategory implements the CreateCategory RPC method
func (s *CategoryGrpcServer) CreateCategory(ctx context.Context, req *pb.Category) (*pb.CreateCategoryResponse, error) {
	domainCategory := convertProtoCategoryToDomain(req)

	id, err := s.categoryService.CreateCategory(domainCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	// Get the created category
	createdCategory, err := s.categoryService.GetCategoryByID(id)
	if err != nil {
		return nil, fmt.Errorf("category created but failed to retrieve: %w", err)
	}

	return &pb.CreateCategoryResponse{
		Id:       int32(id),
		Category: convertDomainCategoryToProto(createdCategory),
	}, nil
}

// UpdateCategory implements the UpdateCategory RPC method
func (s *CategoryGrpcServer) UpdateCategory(ctx context.Context, req *pb.Category) (*pb.StatusResponse, error) {
	domainCategory := convertProtoCategoryToDomain(req)

	err := s.categoryService.UpdateCategory(domainCategory)
	if err != nil {
		return &pb.StatusResponse{
			Success: false,
			Message: fmt.Sprintf("failed to update category: %v", err),
		}, nil
	}

	return &pb.StatusResponse{
		Success: true,
		Message: "category updated successfully",
	}, nil
}

// DeleteCategory implements the DeleteCategory RPC method
func (s *CategoryGrpcServer) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.StatusResponse, error) {
	err := s.categoryService.DeleteCategory(int(req.Id))
	if err != nil {
		return &pb.StatusResponse{
			Success: false,
			Message: fmt.Sprintf("failed to delete category: %v", err),
		}, nil
	}

	return &pb.StatusResponse{
		Success: true,
		Message: "category deleted successfully",
	}, nil
}

// SyncProductCounts implements the SyncProductCounts RPC method
func (s *CategoryGrpcServer) SyncProductCounts(ctx context.Context, req *pb.SyncProductCountsRequest) (*pb.StatusResponse, error) {
	err := s.categoryService.SyncProductCounts()
	if err != nil {
		return &pb.StatusResponse{
			Success: false,
			Message: fmt.Sprintf("failed to sync product counts: %v", err),
		}, nil
	}

	return &pb.StatusResponse{
		Success: true,
		Message: "product counts synced successfully",
	}, nil
}

// Helper functions for conversion between domain and proto models
func convertDomainCategoryToProto(category *domain.Category) *pb.Category {
	return &pb.Category{
		Id:           int32(category.ID),
		Name:         category.Name,
		Slug:         category.Slug,
		Description:  category.Description,
		ImageUrl:     category.ImageURL,
		ProductCount: int32(category.ProductCount),
		IsActive:     category.IsActive,
		IsVisible:    category.IsVisible,
		CreatedAt:    category.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    category.UpdatedAt.Format(time.RFC3339),
	}
}

func convertProtoCategoryToDomain(category *pb.Category) *domain.Category {
	return &domain.Category{
		ID:           int(category.Id),
		Name:         category.Name,
		Slug:         category.Slug,
		Description:  category.Description,
		ImageURL:     category.ImageUrl,
		ProductCount: int(category.ProductCount),
		IsActive:     category.IsActive,
		IsVisible:    category.IsVisible,
	}
}
