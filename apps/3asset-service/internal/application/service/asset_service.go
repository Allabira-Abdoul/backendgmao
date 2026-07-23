package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"backend-gmao/apps/asset-service/internal/core/domain"
	"backend-gmao/apps/asset-service/internal/core/ports/primary"
	"backend-gmao/apps/asset-service/internal/core/ports/secondary"
)

type assetService struct {
	repo secondary.AssetRepository
}

func NewAssetService(repo secondary.AssetRepository) primary.AssetService {
	return &assetService{repo: repo}
}

func (s *assetService) CreateSite(ctx context.Context, req domain.CreateSiteRequest) (*domain.SiteResponse, error) {
	site := &domain.Site{
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
	}
	if err := s.repo.CreateSite(ctx, site); err != nil {
		return nil, err
	}
	resp := site.ToResponse()
	return &resp, nil
}

func (s *assetService) GetSiteHierarchy(ctx context.Context, siteID uuid.UUID) (*domain.SiteResponse, error) {
	site, err := s.repo.GetSiteHierarchy(ctx, siteID)
	if err != nil {
		return nil, err
	}
	resp := site.ToResponse()
	return &resp, nil
}

func (s *assetService) GetAllSites(ctx context.Context) ([]domain.SiteResponse, error) {
	sites, err := s.repo.GetAllSites(ctx)
	if err != nil {
		return nil, err
	}
	var resps []domain.SiteResponse
	for _, site := range sites {
		resps = append(resps, site.ToResponse())
	}
	return resps, nil
}

func (s *assetService) CreateSystem(ctx context.Context, req domain.CreateSystemRequest) (*domain.SystemResponse, error) {
	sys := &domain.System{
		SiteID:      req.SiteID,
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.repo.CreateSystem(ctx, sys); err != nil {
		return nil, err
	}
	resp := sys.ToResponse()
	return &resp, nil
}

func (s *assetService) CreateAsset(ctx context.Context, req domain.CreateAssetRequest) (*domain.AssetResponse, error) {
	asset := &domain.Asset{
		SystemID:     req.SystemID,
		Name:         req.Name,
		Code:         req.Code,
		Model:        req.Model,
		Manufacturer: req.Manufacturer,
	}
	if err := s.repo.CreateAsset(ctx, asset); err != nil {
		return nil, err
	}
	resp := asset.ToResponse()
	return &resp, nil
}

func (s *assetService) UpdateAssetStatus(ctx context.Context, assetID uuid.UUID, status string) error {
	return s.repo.UpdateAssetStatus(ctx, assetID, status)
}

func (s *assetService) CreateSubsystem(ctx context.Context, req domain.CreateSubsystemRequest) (*domain.SubsystemResponse, error) {
	sub := &domain.Subsystem{
		AssetID:     req.AssetID,
		Name:        req.Name,
		Description: req.Description,
		Criticality: req.Criticality,
	}
	if err := s.repo.CreateSubsystem(ctx, sub); err != nil {
		return nil, err
	}
	resp := sub.ToResponse()
	return &resp, nil
}

func (s *assetService) CreateInventoryItem(ctx context.Context, req domain.CreateInventoryItemRequest) (*domain.InventoryItemResponse, error) {
	item := &domain.InventoryItem{
		ItemType:             req.ItemType,
		PartNumber:           req.PartNumber,
		Name:                 req.Name,
		Category:             req.Category,
		StockQuantity:        req.StockQuantity,
		ReorderPoint:         req.ReorderPoint,
		SupplierLeadTimeDays: req.SupplierLeadTimeDays,
		UnitOfMeasure:        req.UnitOfMeasure,
	}
	if err := s.repo.CreateInventoryItem(ctx, item); err != nil {
		return nil, err
	}
	resp := item.ToResponse()
	return &resp, nil
}

func (s *assetService) CreateComponent(ctx context.Context, req domain.CreateComponentRequest) (*domain.ComponentResponse, error) {
	// Business Rule: Validate inventory item is a SPARE_PART and has stock
	item, err := s.repo.GetInventoryItem(ctx, req.InventoryItemID)
	if err != nil {
		return nil, errors.New("inventory item not found")
	}

	if item.ItemType != "SPARE_PART" {
		return nil, errors.New("cannot create component from non-spare part")
	}

	if item.StockQuantity <= 0 {
		return nil, errors.New("insufficient stock to install this component")
	}

	comp := &domain.Component{
		SubsystemID:     req.SubsystemID,
		InventoryItemID: req.InventoryItemID,
		Name:            req.Name,
		SerialNumber:    req.SerialNumber,
	}

	if err := s.repo.CreateComponent(ctx, comp); err != nil {
		return nil, err
	}

	// Decrease stock
	if err := s.repo.UpdateInventoryStock(ctx, item.ID, -1); err != nil {
		return nil, err
	}

	// Check reorder point
	if item.StockQuantity-1 <= item.ReorderPoint {
		// Log an alert or trigger event (Mocking it for now)
		log.Printf("ALERT: Stock for %s (%s) has fallen below reorder point! Current: %d, Reorder Point: %d",
			item.Name, item.PartNumber, item.StockQuantity-1, item.ReorderPoint)
	}

	resp := comp.ToResponse()
	return &resp, nil
}
