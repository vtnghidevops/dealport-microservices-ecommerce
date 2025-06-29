// services/product-service/internal/transport/grpc/ads_server.go
package grpc

import (
	"context"
	"fmt"
	"log"
	"product-service/internal/domain"
	pb "product-service/proto/product"
	"time"
)

// AdsGrpcServer represents the gRPC server for ads service
type AdsGrpcServer struct {
	pb.UnimplementedAdsServiceServer
	adsService domain.AdsService
}

// NewAdsGrpcServer creates a new gRPC server with the provided ads service
func NewAdsGrpcServer(adsSvc domain.AdsService) *AdsGrpcServer {
	log.Println("Creating new gRPC server for ads service")
	return &AdsGrpcServer{
		adsService: adsSvc,
	}
}

// GetAdsByLocation implements the GetAdsByLocation RPC method
func (s *AdsGrpcServer) GetAdsByLocation(ctx context.Context, req *pb.GetAdsByLocationRequest) (*pb.GetAdsByLocationResponse, error) {
	ads, err := s.adsService.GetAdsByLocation(req.Location)
	if err != nil {
		return nil, fmt.Errorf("failed to get ads by location: %w", err)
	}

	// Convert domain ads to proto ads
	protoAds := make([]*pb.AdsData, 0, len(ads))
	for _, ad := range ads {
		protoAds = append(protoAds, convertDomainAdsDataToProto(ad))
	}

	return &pb.GetAdsByLocationResponse{
		Ads: protoAds,
	}, nil
}

// GetAdsPlacementByID implements the GetAdsPlacementByID RPC method
func (s *AdsGrpcServer) GetAdsPlacementByID(ctx context.Context, req *pb.GetAdsPlacementByIDRequest) (*pb.AdsPlacement, error) {
	ad, err := s.adsService.GetAdsPlacementByID(int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to get ads placement: %w", err)
	}

	return convertDomainAdsPlacementToProto(ad), nil
}

// ListAdsPlacements implements the ListAdsPlacements RPC method
func (s *AdsGrpcServer) ListAdsPlacements(ctx context.Context, req *pb.ListAdsPlacementsRequest) (*pb.ListAdsPlacementsResponse, error) {
	ads, total, err := s.adsService.GetAllAdsPlacements(int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, fmt.Errorf("failed to list ads placements: %w", err)
	}

	// Convert domain ads to proto ads
	protoAds := make([]*pb.AdsPlacement, 0, len(ads))
	for _, ad := range ads {
		protoAds = append(protoAds, convertDomainAdsPlacementToProto(ad))
	}

	return &pb.ListAdsPlacementsResponse{
		Placements: protoAds,
		Total:      int32(total),
	}, nil
}

// CreateAdsPlacement implements the CreateAdsPlacement RPC method
func (s *AdsGrpcServer) CreateAdsPlacement(ctx context.Context, req *pb.AdsPlacement) (*pb.CreateAdsPlacementResponse, error) {
	domainAd := convertProtoAdsPlacementToDomain(req)

	id, err := s.adsService.CreateAdsPlacement(domainAd)
	if err != nil {
		return nil, fmt.Errorf("failed to create ads placement: %w", err)
	}

	// Get the created ad
	createdAd, err := s.adsService.GetAdsPlacementByID(id)
	if err != nil {
		return nil, fmt.Errorf("ads placement created but failed to retrieve: %w", err)
	}

	return &pb.CreateAdsPlacementResponse{
		Id:        int32(id),
		Placement: convertDomainAdsPlacementToProto(createdAd),
	}, nil
}

// UpdateAdsPlacement implements the UpdateAdsPlacement RPC method
func (s *AdsGrpcServer) UpdateAdsPlacement(ctx context.Context, req *pb.AdsPlacement) (*pb.StatusResponse, error) {
	domainAd := convertProtoAdsPlacementToDomain(req)

	err := s.adsService.UpdateAdsPlacement(domainAd)
	if err != nil {
		return &pb.StatusResponse{
			Success: false,
			Message: fmt.Sprintf("failed to update ads placement: %v", err),
		}, nil
	}

	return &pb.StatusResponse{
		Success: true,
		Message: "ads placement updated successfully",
	}, nil
}

// DeleteAdsPlacement implements the DeleteAdsPlacement RPC method
func (s *AdsGrpcServer) DeleteAdsPlacement(ctx context.Context, req *pb.DeleteAdsPlacementRequest) (*pb.StatusResponse, error) {
	err := s.adsService.DeleteAdsPlacement(int(req.Id))
	if err != nil {
		return &pb.StatusResponse{
			Success: false,
			Message: fmt.Sprintf("failed to delete ads placement: %v", err),
		}, nil
	}

	return &pb.StatusResponse{
		Success: true,
		Message: "ads placement deleted successfully",
	}, nil
}

// Helper functions for conversion between domain and proto models
func convertDomainAdsPlacementToProto(ad *domain.AdsPlacement) *pb.AdsPlacement {
	uiSettings := ""
	if ad.UISettings != nil {
		uiSettings = string(ad.UISettings)
	}

	return &pb.AdsPlacement{
		Id:             int32(ad.ID),
		Location:       ad.Location,
		ReferenceType:  ad.ReferenceType,
		ReferenceId:    int32(ad.ReferenceID),
		DisplayOrder:   int32(ad.DisplayOrder),
		CustomTitle:    ad.CustomTitle,
		CustomImageUrl: ad.CustomImageURL,
		UiSettings:     uiSettings,
		IsActive:       ad.IsActive,
		CreatedAt:      ad.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      ad.UpdatedAt.Format(time.RFC3339),
	}
}

func convertProtoAdsPlacementToDomain(ad *pb.AdsPlacement) *domain.AdsPlacement {
	var uiSettings []byte
	if ad.UiSettings != "" {
		uiSettings = []byte(ad.UiSettings)
	}

	return &domain.AdsPlacement{
		ID:             int(ad.Id),
		Location:       ad.Location,
		ReferenceType:  ad.ReferenceType,
		ReferenceID:    int(ad.ReferenceId),
		DisplayOrder:   int(ad.DisplayOrder),
		CustomTitle:    ad.CustomTitle,
		CustomImageURL: ad.CustomImageUrl,
		UISettings:     uiSettings,
		IsActive:       ad.IsActive,
	}
}

func convertDomainAdsDataToProto(ad *domain.AdsData) *pb.AdsData {
	uiSettings := ""
	if ad.UISettings != nil {
		uiSettings = string(ad.UISettings)
	}

	return &pb.AdsData{
		Id:           int32(ad.ID),
		Type:         ad.Type,
		Name:         ad.Name,
		Slug:         ad.Slug,
		ImageUrl:     ad.ImageURL,
		Price:        ad.Price,
		CategoryId:   int32(ad.CategoryID),
		CategorySlug: ad.CategorySlug,
		CustomTitle:  ad.CustomTitle,
		UiSettings:   uiSettings,
	}
}
