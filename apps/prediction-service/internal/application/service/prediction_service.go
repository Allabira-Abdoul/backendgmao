package service

import (
	"context"
	"errors"

	"backend-gmao/apps/prediction-service/internal/core/domain"
	"backend-gmao/apps/prediction-service/internal/core/ports/secondary"
	"github.com/google/uuid"
)

var (
	ErrPredictionNotFound = errors.New("prediction not found")
)

// PredictionService implements primary.PredictionService.
type PredictionService struct {
	predictionRepo secondary.PredictionRepository
	assetClient    secondary.AssetClient
}

// NewPredictionService initializes a new PredictionService instance.
func NewPredictionService(predictionRepo secondary.PredictionRepository, assetClient secondary.AssetClient) *PredictionService {
	return &PredictionService{
		predictionRepo: predictionRepo,
		assetClient:    assetClient,
	}
}

func (s *PredictionService) RecordPrediction(ctx context.Context, req domain.CreatePredictionRequest) (*domain.PredictionResponse, error) {
	assetUUID, err := uuid.Parse(req.AssetID)
	if err != nil {
		return nil, errors.New("invalid asset ID format")
	}

	prediction := &domain.Prediction{
		ID:                   uuid.New(),
		AssetID:              assetUUID,
		FailureProbability:   req.FailureProbability,
		PredictedFailureDate: req.PredictedFailureDate,
	}

	if err := s.predictionRepo.Save(ctx, prediction); err != nil {
		return nil, err
	}

	resp := s.buildPredictionResponse(ctx, prediction)
	return &resp, nil
}

func (s *PredictionService) GetPrediction(ctx context.Context, id uuid.UUID) (*domain.PredictionResponse, error) {
	prediction, err := s.predictionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrPredictionNotFound
	}
	resp := s.buildPredictionResponse(ctx, prediction)
	return &resp, nil
}

func (s *PredictionService) GetPredictionsForAsset(ctx context.Context, assetID uuid.UUID) ([]domain.PredictionResponse, error) {
	predictions, err := s.predictionRepo.FindByAssetID(ctx, assetID)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.PredictionResponse, len(predictions))
	for i, p := range predictions {
		responses[i] = s.buildPredictionResponse(ctx, &p)
	}
	return responses, nil
}

func (s *PredictionService) GetAllPredictions(ctx context.Context) ([]domain.PredictionResponse, error) {
	predictions, err := s.predictionRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.PredictionResponse, len(predictions))
	for i, p := range predictions {
		responses[i] = s.buildPredictionResponse(ctx, &p)
	}
	return responses, nil
}

func (s *PredictionService) buildPredictionResponse(ctx context.Context, p *domain.Prediction) domain.PredictionResponse {
	assetName := "Unknown Asset"
	if name, err := s.assetClient.GetAssetName(ctx, p.AssetID); err == nil {
		assetName = name
	}
	return p.ToResponse(assetName)
}
