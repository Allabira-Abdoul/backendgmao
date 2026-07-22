package service

import (
	"context"
	"fmt"
	"time"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"backend-gmao/pkg/middleware"
	"github.com/google/uuid"
)

// CreateMaintenanceSchedule creates a new maintenance schedule.
func (s *MaintenanceService) CreateMaintenanceSchedule(ctx context.Context, req domain.CreateMaintenanceScheduleRequest) (*domain.MaintenanceScheduleResponse, error) {
	assetID, err := uuid.Parse(req.AssetID)
	if err != nil {
		return nil, fmt.Errorf("invalid asset ID: %w", err)
	}

	schedule := &domain.MaintenanceSchedule{
		AssetID:               assetID,
		Title:                 req.Title,
		Description:           req.Description,
		Frequency:             req.Frequency,
		IntervalMonths:        req.IntervalMonths,
		IntervalHours:         req.IntervalHours,
		StartDate:             req.StartDate,
		EndDate:               req.EndDate,
		NextScheduledDate:     req.NextScheduledDate,
		NextScheduledUsage:    req.NextScheduledUsage,
		MaintenanceCategory:   req.MaintenanceCategory,
		MaintenanceType:       req.MaintenanceType,
		IsActive:              true, // Default
		RequireCounterReading: false,
	}

	if schedule.NextScheduledDate == nil && schedule.StartDate != nil {
		schedule.NextScheduledDate = schedule.StartDate
	}

	if req.IsActive != nil {
		schedule.IsActive = *req.IsActive
	}
	if req.RequireCounterReading != nil {
		schedule.RequireCounterReading = *req.RequireCounterReading
	}

	if err := s.maintenanceRepo.CreateMaintenanceSchedule(ctx, schedule); err != nil {
		return nil, fmt.Errorf("create schedule: %w", err)
	}

	s.auditLogger.LogAction(ctx, "MAINTENANCE_SCHEDULE_CREATED", fmt.Errorf("Schedule %s created for asset %s", schedule.ID, schedule.AssetID).Error())

	resp := schedule.ToResponse()
	return &resp, nil
}

// UpdateMaintenanceSchedule updates an existing maintenance schedule.
func (s *MaintenanceService) UpdateMaintenanceSchedule(ctx context.Context, id uuid.UUID, req domain.UpdateMaintenanceScheduleRequest) (*domain.MaintenanceScheduleResponse, error) {
	schedule, err := s.maintenanceRepo.FindMaintenanceScheduleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find schedule: %w", err)
	}

	if req.Title != nil {
		schedule.Title = *req.Title
	}
	if req.Description != nil {
		schedule.Description = *req.Description
	}
	if req.Frequency != nil {
		schedule.Frequency = *req.Frequency
	}
	if req.IntervalMonths != nil {
		schedule.IntervalMonths = req.IntervalMonths
	}
	if req.IntervalHours != nil {
		schedule.IntervalHours = req.IntervalHours
	}
	if req.StartDate != nil {
		schedule.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		schedule.EndDate = req.EndDate
	}
	if req.NextScheduledDate != nil {
		schedule.NextScheduledDate = req.NextScheduledDate
	}
	if req.NextScheduledUsage != nil {
		schedule.NextScheduledUsage = req.NextScheduledUsage
	}
	if req.MaintenanceCategory != nil {
		schedule.MaintenanceCategory = *req.MaintenanceCategory
	}
	if req.MaintenanceType != nil {
		schedule.MaintenanceType = *req.MaintenanceType
	}
	if req.IsActive != nil {
		schedule.IsActive = *req.IsActive
	}
	if req.RequireCounterReading != nil {
		schedule.RequireCounterReading = *req.RequireCounterReading
	}

	if err := s.maintenanceRepo.UpdateMaintenanceSchedule(ctx, schedule); err != nil {
		return nil, fmt.Errorf("update schedule: %w", err)
	}

	s.auditLogger.LogAction(ctx, "MAINTENANCE_SCHEDULE_UPDATED", fmt.Errorf("Schedule %s updated", schedule.ID).Error())

	resp := schedule.ToResponse()
	return &resp, nil
}

// GetMaintenanceSchedulesForAsset returns schedules for an asset.
func (s *MaintenanceService) GetMaintenanceSchedulesForAsset(ctx context.Context, assetID uuid.UUID) ([]domain.MaintenanceScheduleResponse, error) {
	schedules, err := s.maintenanceRepo.FindMaintenanceSchedulesByAssetID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("find schedules by asset ID: %w", err)
	}

	var res []domain.MaintenanceScheduleResponse
	for _, sch := range schedules {
		res = append(res, sch.ToResponse())
	}
	return res, nil
}

// GetAllMaintenanceSchedules returns all maintenance schedules across all assets.
func (s *MaintenanceService) GetAllMaintenanceSchedules(ctx context.Context) ([]domain.MaintenanceScheduleResponse, error) {
	schedules, err := s.maintenanceRepo.FindAllMaintenanceSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("find all schedules: %w", err)
	}

	var res []domain.MaintenanceScheduleResponse
	for _, sch := range schedules {
		res = append(res, sch.ToResponse())
	}
	return res, nil
}

// RecordCounterReading logs a new counter reading.
func (s *MaintenanceService) RecordCounterReading(ctx context.Context, req domain.CreateCounterReadingRequest) (*domain.CounterReadingResponse, error) {
	assetID, err := uuid.Parse(req.AssetID)
	if err != nil {
		return nil, fmt.Errorf("invalid asset ID: %w", err)
	}

	var inspectionIDPtr *uuid.UUID
	if req.InspectionID != nil && *req.InspectionID != "" {
		parsed, err := uuid.Parse(*req.InspectionID)
		if err != nil {
			return nil, fmt.Errorf("invalid inspection ID: %w", err)
		}
		inspectionIDPtr = &parsed
	}

	userIDStr, ok := ctx.Value(middleware.ContextKeyUserID).(string)
	var userID uuid.UUID
	if ok && userIDStr != "" {
		userID, _ = uuid.Parse(userIDStr)
	}

	reading := &domain.CounterReading{
		AssetID:      assetID,
		InspectionID: inspectionIDPtr,
		Value:        req.Value,
		Unit:         req.Unit,
		RecordedAt:   time.Now(),
		RecordedBy:   userID,
	}

	if err := s.maintenanceRepo.CreateCounterReading(ctx, reading); err != nil {
		return nil, fmt.Errorf("create counter reading: %w", err)
	}

	s.auditLogger.LogAction(ctx, "COUNTER_READING_RECORDED", fmt.Errorf("Reading %f %s recorded for asset %s", reading.Value, reading.Unit, reading.AssetID).Error())

	resp := reading.ToResponse()
	return &resp, nil
}

// GetCounterReadingsForAsset retrieves historical readings.
func (s *MaintenanceService) GetCounterReadingsForAsset(ctx context.Context, assetID uuid.UUID) ([]domain.CounterReadingResponse, error) {
	readings, err := s.maintenanceRepo.FindCounterReadingsByAssetID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("find readings by asset ID: %w", err)
	}

	var res []domain.CounterReadingResponse
	for _, r := range readings {
		res = append(res, r.ToResponse())
	}
	return res, nil
}
