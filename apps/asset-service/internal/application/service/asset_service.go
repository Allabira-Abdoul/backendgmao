package service

import (
	"context"
	"errors"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"backend-gmao/apps/asset-service/internal/core/ports/secondary"
	"backend-gmao/pkg/audit"
	"github.com/google/uuid"
)

type assetService struct {
	repo        secondary.AssetRepository
	auditClient audit.Client
}

// NewAssetService creates a new Asset Service instance.
func NewAssetService(repo secondary.AssetRepository, auditClient audit.Client) *assetService {
	return &assetService{repo: repo, auditClient: auditClient}
}

func (s *assetService) CreateEquipmentModel(ctx context.Context, req domain.CreateEquipmentModelRequest) (domain.EquipmentModelResponse, error) {
	model := &domain.EquipmentModel{
		ID:          uuid.New(),
		Name:        req.Name,
		Category:    req.Category,
		Description: req.Description,
	}

	if err := s.repo.CreateEquipmentModel(ctx, model); err != nil {
		return domain.EquipmentModelResponse{}, err
	}

	return model.ToResponse(), nil
}

func (s *assetService) CreatePartModel(ctx context.Context, req domain.CreatePartModelRequest) (domain.PartModelResponse, error) {
	model := &domain.PartModel{
		ID:            uuid.New(),
		Name:          req.Name,
		Category:      req.Category,
		SpareQuantity: req.SpareQuantity,
	}

	if err := s.repo.CreatePartModel(ctx, model); err != nil {
		return domain.PartModelResponse{}, err
	}

	return model.ToResponse(), nil
}

func (s *assetService) GetEquipmentModels(ctx context.Context) ([]domain.EquipmentModelResponse, error) {
	models, err := s.repo.GetEquipmentModels(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]domain.EquipmentModelResponse, len(models))
	for i, m := range models {
		res[i] = m.ToResponse()
	}
	return res, nil
}

func (s *assetService) GetPartModels(ctx context.Context) ([]domain.PartModelResponse, error) {
	models, err := s.repo.GetPartModels(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]domain.PartModelResponse, len(models))
	for i, m := range models {
		res[i] = m.ToResponse()
	}
	return res, nil
}

func (s *assetService) CreateEquipmentInstance(ctx context.Context, req domain.CreateEquipmentInstanceRequest) (domain.EquipmentInstanceResponse, error) {
	// Verify model exists
	model, err := s.repo.GetEquipmentModelByID(ctx, req.EquipmentModelID)
	if err != nil || model == nil {
		return domain.EquipmentInstanceResponse{}, errors.New("equipment model not found")
	}

	instance := &domain.EquipmentInstance{
		ID:               uuid.New(),
		Code:             req.Code,
		EquipmentModelID: req.EquipmentModelID,
		Status:           "OPERATIONAL",
		Location:         req.Location,
		PurchaseDate:     req.PurchaseDate,
		PurchaseValue:    req.PurchaseValue,
	}

	if err := s.repo.CreateEquipmentInstance(ctx, instance); err != nil {
		return domain.EquipmentInstanceResponse{}, err
	}

	return instance.ToResponse(), nil
}

func (s *assetService) CreatePartInstance(ctx context.Context, req domain.CreatePartInstanceRequest) (domain.PartInstanceResponse, error) {
	// Verify equipment instance exists
	_, err := s.repo.GetEquipmentInstanceByID(ctx, req.EquipmentInstanceID)
	if err != nil {
		return domain.PartInstanceResponse{}, errors.New("equipment instance not found")
	}

	// Verify part model exists
	_, err = s.repo.GetPartModelByID(ctx, req.PartModelID)
	if err != nil {
		return domain.PartInstanceResponse{}, errors.New("part model not found")
	}

	instance := &domain.PartInstance{
		ID:                  uuid.New(),
		EquipmentInstanceID: req.EquipmentInstanceID,
		PartModelID:         req.PartModelID,
		Status:              "OPERATIONAL",
	}

	if err := s.repo.CreatePartInstance(ctx, instance); err != nil {
		return domain.PartInstanceResponse{}, err
	}

	return instance.ToResponse(), nil
}

func (s *assetService) GetEquipmentInstances(ctx context.Context) ([]domain.EquipmentInstanceResponse, error) {
	instances, err := s.repo.GetEquipmentInstances(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]domain.EquipmentInstanceResponse, len(instances))
	for i, inst := range instances {
		res[i] = inst.ToResponse()
	}
	return res, nil
}

func (s *assetService) GetEquipmentInstanceByCode(ctx context.Context, code string) (domain.EquipmentInstanceResponse, error) {
	instance, err := s.repo.GetEquipmentInstanceByCode(ctx, code)
	if err != nil || instance == nil {
		return domain.EquipmentInstanceResponse{}, errors.New("equipment instance not found")
	}

	return instance.ToResponse(), nil
}

func (s *assetService) GetEquipmentInstanceByID(ctx context.Context, id uuid.UUID) (domain.EquipmentInstanceResponse, error) {
	instance, err := s.repo.GetEquipmentInstanceByID(ctx, id)
	if err != nil || instance == nil {
		return domain.EquipmentInstanceResponse{}, errors.New("equipment instance not found")
	}

	return instance.ToResponse(), nil
}
