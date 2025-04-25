// services/product-service/internal/transport/grpc/banner_server.go
package grpc

import (
	"context"
	"fmt"
	"log"
	"product-service/internal/domain"
	pb "product-service/proto/product"
	"time"
)

// BannerGrpcServer represents the gRPC server for banner service
type BannerGrpcServer struct {
	pb.UnimplementedBannerServiceServer
	bannerService domain.BannerService
}

// NewBannerGrpcServer creates a new gRPC server with the provided banner service
func NewBannerGrpcServer(bannerSvc domain.BannerService) *BannerGrpcServer {
	log.Println("Creating new gRPC server for banner service")
	return &BannerGrpcServer{
		bannerService: bannerSvc,
	}
}

// GetBannerByID implements the GetBannerByID RPC method
func (s *BannerGrpcServer) GetBannerByID(ctx context.Context, req *pb.GetBannerByIDRequest) (*pb.Banner, error) {
	banner, err := s.bannerService.GetBannerByID(int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to get banner: %w", err)
	}

	return convertDomainBannerToProto(banner), nil
}

// GetBannersByType implements the GetBannersByType RPC method
func (s *BannerGrpcServer) GetBannersByType(ctx context.Context, req *pb.GetBannersByTypeRequest) (*pb.GetBannersByTypeResponse, error) {
	banners, err := s.bannerService.GetBannersByType(req.Type, int(req.Limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get banners by type: %w", err)
	}

	protoBanners := make([]*pb.Banner, 0, len(banners))
	for _, banner := range banners {
		protoBanners = append(protoBanners, convertDomainBannerToProto(banner))
	}

	return &pb.GetBannersByTypeResponse{
		Banners: protoBanners,
	}, nil
}

// ListBanners implements the ListBanners RPC method
func (s *BannerGrpcServer) ListBanners(ctx context.Context, req *pb.ListBannersRequest) (*pb.ListBannersResponse, error) {
	filters := make(map[string]string)
	for k, v := range req.Filters {
		filters[k] = v
	}

	banners, total, err := s.bannerService.GetAllBanners(int(req.Page), int(req.PageSize), filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list banners: %w", err)
	}

	protoBanners := make([]*pb.Banner, 0, len(banners))
	for _, banner := range banners {
		protoBanners = append(protoBanners, convertDomainBannerToProto(banner))
	}

	return &pb.ListBannersResponse{
		Banners: protoBanners,
		Total:   int32(total),
	}, nil
}

// CreateBanner implements the CreateBanner RPC method
func (s *BannerGrpcServer) CreateBanner(ctx context.Context, req *pb.Banner) (*pb.CreateBannerResponse, error) {
	domainBanner := convertProtoBannerToDomain(req)

	id, err := s.bannerService.CreateBanner(domainBanner)
	if err != nil {
		return nil, fmt.Errorf("failed to create banner: %w", err)
	}

	// Get the created banner
	createdBanner, err := s.bannerService.GetBannerByID(id)
	if err != nil {
		return nil, fmt.Errorf("banner created but failed to retrieve: %w", err)
	}

	return &pb.CreateBannerResponse{
		Id:     int32(id),
		Banner: convertDomainBannerToProto(createdBanner),
	}, nil
}

// UpdateBanner implements the UpdateBanner RPC method
func (s *BannerGrpcServer) UpdateBanner(ctx context.Context, req *pb.Banner) (*pb.StatusResponse, error) {
	domainBanner := convertProtoBannerToDomain(req)

	err := s.bannerService.UpdateBanner(domainBanner)
	if err != nil {
		return &pb.StatusResponse{
			Success: false,
			Message: fmt.Sprintf("failed to update banner: %v", err),
		}, nil
	}

	return &pb.StatusResponse{
		Success: true,
		Message: "banner updated successfully",
	}, nil
}

// DeleteBanner implements the DeleteBanner RPC method
func (s *BannerGrpcServer) DeleteBanner(ctx context.Context, req *pb.DeleteBannerRequest) (*pb.StatusResponse, error) {
	err := s.bannerService.DeleteBanner(int(req.Id))
	if err != nil {
		return &pb.StatusResponse{
			Success: false,
			Message: fmt.Sprintf("failed to delete banner: %v", err),
		}, nil
	}

	return &pb.StatusResponse{
		Success: true,
		Message: "banner deleted successfully",
	}, nil
}

// Helper functions for conversion between domain and proto models
func convertDomainBannerToProto(banner *domain.Banner) *pb.Banner {
	var productID, categoryID int32
	if banner.ProductID != nil {
		productID = int32(*banner.ProductID)
	}
	if banner.CategoryID != nil {
		categoryID = int32(*banner.CategoryID)
	}

	return &pb.Banner{
		Id:              int32(banner.ID),
		Title:           banner.Title,
		Subtitle:        banner.Subtitle,
		Description:     banner.Description,
		Discount:        banner.Discount,
		HighlightText:   banner.HighlightText,
		ImageUrl:        banner.ImageURL,
		LinkUrl:         banner.LinkURL,
		ActionText:      banner.ActionText,
		BackgroundColor: banner.BackgroundColor,
		TextColor:       banner.TextColor,
		AnimationType:   banner.AnimationType,
		IsActive:        banner.IsActive,
		Priority:        int32(banner.Priority),
		Type:            banner.Type,
		ProductId:       productID,
		CategoryId:      categoryID,
		CreatedAt:       banner.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       banner.UpdatedAt.Format(time.RFC3339),
	}
}

func convertProtoBannerToDomain(banner *pb.Banner) *domain.Banner {
	var productID, categoryID *int
	if banner.ProductId > 0 {
		productIDInt := int(banner.ProductId)
		productID = &productIDInt
	}
	if banner.CategoryId > 0 {
		categoryIDInt := int(banner.CategoryId)
		categoryID = &categoryIDInt
	}

	return &domain.Banner{
		ID:              int(banner.Id),
		Title:           banner.Title,
		Subtitle:        banner.Subtitle,
		Description:     banner.Description,
		Discount:        banner.Discount,
		HighlightText:   banner.HighlightText,
		ImageURL:        banner.ImageUrl,
		LinkURL:         banner.LinkUrl,
		ActionText:      banner.ActionText,
		BackgroundColor: banner.BackgroundColor,
		TextColor:       banner.TextColor,
		AnimationType:   banner.AnimationType,
		IsActive:        banner.IsActive,
		Priority:        int(banner.Priority),
		Type:            banner.Type,
		ProductID:       productID,
		CategoryID:      categoryID,
	}
}
