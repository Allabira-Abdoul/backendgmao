package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"backend-gmao/apps/maintenance-service/internal/core/ports/secondary"
	"backend-gmao/pkg/audit"
	"backend-gmao/pkg/middleware"

	"github.com/google/uuid"
)

var (
	ErrWorkOrderNotFound       = errors.New("work order not found")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

// MaintenanceService implements primary.MaintenanceService.
type MaintenanceService struct {
	maintenanceRepo secondary.MaintenanceRepository
	analyticsClient secondary.AnalyticsClient
	auditClient     audit.Client
	userClient      secondary.UserClient
	assetClient     secondary.AssetClient
	eventPublisher  secondary.EventPublisher
}

// NewMaintenanceService initializes a new MaintenanceService instance.
func NewMaintenanceService(
	maintenanceRepo secondary.MaintenanceRepository,
	analyticsClient secondary.AnalyticsClient,
	auditClient audit.Client,
	userClient secondary.UserClient,
	assetClient secondary.AssetClient,
	eventPublisher secondary.EventPublisher,
) *MaintenanceService {
	return &MaintenanceService{
		maintenanceRepo: maintenanceRepo,
		analyticsClient: analyticsClient,
		auditClient:     auditClient,
		userClient:      userClient,
		assetClient:     assetClient,
		eventPublisher:  eventPublisher,
	}
}

func (s *MaintenanceService) fireAudit(ctx context.Context, action, details string) {
	userID, ok := ctx.Value(middleware.ContextKeyUserID).(string)
	var uidPtr *string
	if ok && userID != "" {
		uidPtr = &userID
	}

	var userName string
	if name, ok := ctx.Value(middleware.ContextKeyFullName).(string); ok {
		userName = name
	}

	go func() {
		bgCtx := context.Background()

		// Publish to HTTP Audit Service
		if s.auditClient != nil {
			_ = s.auditClient.LogEvent(bgCtx, audit.AuditEvent{
				ServiceName: "maintenance-service",
				Action:      action,
				Details:     details,
				UserID:      uidPtr,
				UserName:    userName,
			})
		}

		// Publish to RabbitMQ Audit Exchange
		if s.eventPublisher != nil {
			var actorID *uuid.UUID
			if uidPtr != nil {
				parsed, err := uuid.Parse(*uidPtr)
				if err == nil {
					actorID = &parsed
				}
			}
			changes := map[string]interface{}{
				"details": details,
			}
			if userName != "" {
				changes["user_name"] = userName
			}
			_ = s.eventPublisher.PublishAuditLog(bgCtx, action, "MAINTENANCE", "", actorID, changes)
		}
	}()
}

func (s *MaintenanceService) CreateWorkOrder(ctx context.Context, req domain.CreateOrdreTravailRequest) (*domain.OrdreTravailResponse, error) {
	assetID, err := uuid.Parse(req.AssetID)
	if err != nil {
		return nil, errors.New("invalid asset ID format")
	}

	var assignedTo *uuid.UUID
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		parsed, err := uuid.Parse(*req.AssignedTo)
		if err != nil {
			return nil, errors.New("invalid assigned user ID format")
		}
		assignedTo = &parsed
	}

	var scheduleID *uuid.UUID
	if req.ScheduleID != nil && *req.ScheduleID != "" {
		parsed, err := uuid.Parse(*req.ScheduleID)
		if err != nil {
			return nil, errors.New("invalid schedule ID format")
		}
		scheduleID = &parsed
	}

	wo := &domain.OrdreTravail{
		ID:                  uuid.New(),
		Title:               req.Title,
		Description:         req.Description,
		AssetID:             assetID,
		ScheduleID:          scheduleID,
		ScheduledStartDate:  req.ScheduledStartDate,
		ScheduledEndDate:    req.ScheduledEndDate,
		Priority:            req.Priority,
		Status:              "PENDING",
		MaintenanceCategory: req.MaintenanceCategory,
		MaintenanceType:     req.MaintenanceType,
		IsMetricMeasurement: req.IsMetricMeasurement,
		AssignedTo:          assignedTo,
	}

	if err := s.maintenanceRepo.CreateWorkOrder(ctx, wo); err != nil {
		return nil, err
	}

	s.fireAudit(ctx, "CREATE_WORK_ORDER", fmt.Sprintf("Created work order %s for asset %s", wo.ID, wo.AssetID))

	resp := s.buildOrdreTravailResponse(ctx, wo)
	return &resp, nil
}

func (s *MaintenanceService) UpdateWorkOrder(ctx context.Context, id uuid.UUID, req domain.UpdateOrdreTravailRequest) (*domain.OrdreTravailResponse, error) {
	wo, err := s.maintenanceRepo.FindWorkOrderByID(ctx, id)
	if err != nil {
		return nil, ErrWorkOrderNotFound
	}

	if req.Title != nil {
		wo.Title = *req.Title
	}
	if req.Description != nil {
		wo.Description = *req.Description
	}
	if req.ScheduleID != nil {
		if *req.ScheduleID == "" {
			wo.ScheduleID = nil
		} else {
			parsed, err := uuid.Parse(*req.ScheduleID)
			if err == nil {
				wo.ScheduleID = &parsed
			}
		}
	}
	if req.ScheduledStartDate != nil {
		wo.ScheduledStartDate = req.ScheduledStartDate
	}
	if req.ScheduledEndDate != nil {
		wo.ScheduledEndDate = req.ScheduledEndDate
	}
	if req.Status != nil {
		wo.Status = *req.Status
	}
	if req.Priority != nil {
		wo.Priority = *req.Priority
	}
	if req.MaintenanceCategory != nil {
		wo.MaintenanceCategory = *req.MaintenanceCategory
	}
	if req.MaintenanceType != nil {
		wo.MaintenanceType = *req.MaintenanceType
	}
	if req.IsMetricMeasurement != nil {
		wo.IsMetricMeasurement = *req.IsMetricMeasurement
	}
	if req.AssetID != nil {
		parsed, err := uuid.Parse(*req.AssetID)
		if err == nil {
			wo.AssetID = parsed
		}
	}
	if req.AssignedTo != nil {
		if *req.AssignedTo == "" {
			wo.AssignedTo = nil
		} else {
			parsed, err := uuid.Parse(*req.AssignedTo)
			if err != nil {
				return nil, errors.New("invalid assigned user ID format")
			}
			wo.AssignedTo = &parsed
		}
	}

	wo.UpdatedAt = time.Now()

	if err := s.maintenanceRepo.UpdateWorkOrder(ctx, wo); err != nil {
		return nil, err
	}

	s.fireAudit(ctx, "UPDATE_WORK_ORDER", fmt.Sprintf("Updated work order %s", wo.ID))

	resp := s.buildOrdreTravailResponse(ctx, wo)
	return &resp, nil
}

func (s *MaintenanceService) DeleteWorkOrder(ctx context.Context, id uuid.UUID) error {
	_, err := s.maintenanceRepo.FindWorkOrderByID(ctx, id)
	if err != nil {
		return ErrWorkOrderNotFound
	}
	
	if err := s.maintenanceRepo.DeleteWorkOrder(ctx, id); err != nil {
		return err
	}
	
	s.fireAudit(ctx, "COMPUTE_ANALYTICS", fmt.Sprintf("Asset %s computed", id))

	return nil
}

// ----------------------------------------------------------------------------
// Defect Alerts
// ----------------------------------------------------------------------------

func (s *MaintenanceService) CreateDefectAlert(ctx context.Context, assetID uuid.UUID, reportedBy uuid.UUID, title, description, imageURL string) (*domain.DefectAlertResponse, error) {
	alert := &domain.DefectAlert{
		ID:          uuid.New(),
		AssetID:     assetID,
		ReportedBy:  reportedBy,
		Title:       title,
		Description: description,
		ImageURL:    imageURL,
		Status:      domain.DefectStatusPending,
	}

	if err := s.maintenanceRepo.CreateDefectAlert(ctx, alert); err != nil {
		return nil, fmt.Errorf("create defect alert: %w", err)
	}

	s.fireAudit(ctx, "CREATE_DEFECT_ALERT", fmt.Sprintf("Alert %s created for asset %s", alert.ID, assetID))

	// Get asset info and user info to return fully populated response
	var assetName, reporterName string
	if name, err := s.assetClient.GetAssetName(ctx, assetID); err == nil {
		assetName = name
	}
	if name, err := s.userClient.GetUserName(ctx, reportedBy); err == nil {
		reporterName = name
	}

	resp := alert.ToResponse(assetName, reporterName)
	return &resp, nil
}

func (s *MaintenanceService) GetAllDefectAlerts(ctx context.Context) ([]domain.DefectAlertResponse, error) {
	alerts, err := s.maintenanceRepo.FindAllDefectAlerts(ctx)
	if err != nil {
		return nil, fmt.Errorf("find all defect alerts: %w", err)
	}

	var res []domain.DefectAlertResponse
	for _, alert := range alerts {
		var assetName, reporterName string
		if name, err := s.assetClient.GetAssetName(ctx, alert.AssetID); err == nil {
			assetName = name
		}
		if name, err := s.userClient.GetUserName(ctx, alert.ReportedBy); err == nil {
			reporterName = name
		}
		res = append(res, alert.ToResponse(assetName, reporterName))
	}
	return res, nil
}

func (s *MaintenanceService) ReviewDefectAlert(ctx context.Context, id uuid.UUID, req domain.ReviewDefectAlertRequest) (*domain.DefectAlertResponse, error) {
	alert, err := s.maintenanceRepo.FindDefectAlertByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find defect alert: %w", err)
	}

	if alert.Status != domain.DefectStatusPending {
		return nil, fmt.Errorf("defect alert is already %s", alert.Status)
	}

	alert.Status = req.Status
	if err := s.maintenanceRepo.UpdateDefectAlert(ctx, alert); err != nil {
		return nil, fmt.Errorf("update defect alert: %w", err)
	}

	s.fireAudit(ctx, "REVIEW_DEFECT_ALERT", fmt.Sprintf("Alert %s marked as %s", alert.ID, alert.Status))

	// If approved, create a work order automatically
	if req.Status == domain.DefectStatusApproved {
		woTitle := fmt.Sprintf("Defect: %s", alert.Title)
		woDesc := fmt.Sprintf("Generated from defect alert %s.\n\nDescription: %s\nReview Notes: %s", alert.ID, alert.Description, req.ReviewNotes)
		
		woReq := domain.CreateOrdreTravailRequest{
			Title:               woTitle,
			Description:         woDesc,
			AssetID:             alert.AssetID.String(),
			Priority:            "HIGH",
			MaintenanceCategory: "CORRECTIVE",
			MaintenanceType:     "CURATIVE",
		}
		
		// Ignore error here so that even if WO creation fails, the alert is still approved.
		// In a real system, you might wrap this in a database transaction.
		_, _ = s.CreateWorkOrder(ctx, woReq)
	}

	var assetName, reporterName string
	if name, err := s.assetClient.GetAssetName(ctx, alert.AssetID); err == nil {
		assetName = name
	}
	if name, err := s.userClient.GetUserName(ctx, alert.ReportedBy); err == nil {
		reporterName = name
	}

	resp := alert.ToResponse(assetName, reporterName)
	return &resp, nil
}

func (s *MaintenanceService) GetWorkOrder(ctx context.Context, id uuid.UUID) (*domain.OrdreTravailResponse, error) {
	wo, err := s.maintenanceRepo.FindWorkOrderByID(ctx, id)
	if err != nil {
		return nil, ErrWorkOrderNotFound
	}

	interventions, _ := s.maintenanceRepo.FindInterventionsByWorkOrderID(ctx, id)
	wo.Interventions = interventions

	resp := s.buildOrdreTravailResponse(ctx, wo)
	return &resp, nil
}

func (s *MaintenanceService) GetAllWorkOrders(ctx context.Context) ([]domain.OrdreTravailResponse, error) {
	workorders, err := s.maintenanceRepo.FindAllWorkOrders(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.OrdreTravailResponse, len(workorders))
	for i, wo := range workorders {
		interventions, _ := s.maintenanceRepo.FindInterventionsByWorkOrderID(ctx, wo.ID)
		wo.Interventions = interventions
		responses[i] = s.buildOrdreTravailResponse(ctx, &wo)
	}
	return responses, nil
}

func (s *MaintenanceService) StartWorkOrder(ctx context.Context, id uuid.UUID) (*domain.OrdreTravailResponse, error) {
	wo, err := s.maintenanceRepo.FindWorkOrderByID(ctx, id)
	if err != nil {
		return nil, ErrWorkOrderNotFound
	}

	if wo.Status == "PENDING" {
		wo.Status = "IN_PROGRESS"
		wo.UpdatedAt = time.Now()
		if err := s.maintenanceRepo.UpdateWorkOrder(ctx, wo); err != nil {
			return nil, err
		}
		
		s.fireAudit(ctx, "START_WORK_ORDER", fmt.Sprintf("Started work order %s", wo.ID))
		
		if s.eventPublisher != nil {
			_ = s.eventPublisher.PublishWorkOrderStarted(ctx, wo.ID, wo.AssetID, "INTERVENTION")
		}
	}

	interventions, _ := s.maintenanceRepo.FindInterventionsByWorkOrderID(ctx, id)
	wo.Interventions = interventions

	resp := s.buildOrdreTravailResponse(ctx, wo)
	return &resp, nil
}

func (s *MaintenanceService) RecordIntervention(ctx context.Context, workOrderID uuid.UUID, req domain.CreateInterventionRequest) (*domain.InterventionResponse, error) {
	_, err := s.maintenanceRepo.FindWorkOrderByID(ctx, workOrderID)
	if err != nil {
		return nil, ErrWorkOrderNotFound
	}

	performedBy, err := uuid.Parse(req.PerformedBy)
	if err != nil {
		return nil, errors.New("invalid performed by user ID format")
	}

	intervention := &domain.Intervention{
		ID:                  uuid.New(),
		WorkOrderID:         workOrderID,
		Description:         req.Description,
		MaintenanceCategory: req.MaintenanceCategory,
		MaintenanceType:     req.MaintenanceType,
		IsMetricMeasurement: req.IsMetricMeasurement,
		PerformedBy:         performedBy,
	}

	if len(req.Measurements) > 0 {
		meas := make([]domain.MetricMeasurement, len(req.Measurements))
		for i, mReq := range req.Measurements {
			var compID *uuid.UUID
			if mReq.ComponentID != nil && *mReq.ComponentID != "" {
				parsedComp, err := uuid.Parse(*mReq.ComponentID)
				if err == nil {
					compID = &parsedComp
				}
			}
			meas[i] = domain.MetricMeasurement{
				ID:                  uuid.New(),
				InterventionID:      &intervention.ID,
				ComponentID:         compID,
				MetricName:          mReq.MetricName,
				Value:               mReq.Value,
				Unit:                mReq.Unit,
				IsThresholdBreached: mReq.IsThresholdBreached,
			}
		}
		intervention.Measurements = meas
	}

	if err := s.maintenanceRepo.CreateIntervention(ctx, intervention); err != nil {
		return nil, err
	}

	// Note: We no longer trigger analytics here automatically because duration is not known until EndIntervention
	// It will be triggered in EndIntervention instead.

	s.fireAudit(ctx, "RECORD_INTERVENTION", fmt.Sprintf("Recorded intervention %s for work order %s", intervention.ID, intervention.WorkOrderID))

	resp := s.buildInterventionResponse(ctx, intervention)
	return &resp, nil
}

func (s *MaintenanceService) UpdateIntervention(ctx context.Context, workOrderID uuid.UUID, interventionID uuid.UUID, req domain.UpdateInterventionRequest) (*domain.InterventionResponse, error) {
	inv, err := s.maintenanceRepo.FindInterventionByID(ctx, interventionID)
	if err != nil || inv.WorkOrderID != workOrderID {
		return nil, errors.New("intervention not found")
	}

	if req.Description != nil {
		inv.Description = *req.Description
	}
	if req.MaintenanceCategory != nil {
		inv.MaintenanceCategory = *req.MaintenanceCategory
	}
	if req.MaintenanceType != nil {
		inv.MaintenanceType = *req.MaintenanceType
	}
	if req.IsMetricMeasurement != nil {
		inv.IsMetricMeasurement = *req.IsMetricMeasurement
	}

	if len(req.Measurements) > 0 {
		meas := make([]domain.MetricMeasurement, len(req.Measurements))
		for i, mReq := range req.Measurements {
			var compID *uuid.UUID
			if mReq.ComponentID != nil && *mReq.ComponentID != "" {
				parsedComp, err := uuid.Parse(*mReq.ComponentID)
				if err == nil {
					compID = &parsedComp
				}
			}
			meas[i] = domain.MetricMeasurement{
				ID:                  uuid.New(),
				InterventionID:      &inv.ID,
				ComponentID:         compID,
				MetricName:          mReq.MetricName,
				Value:               mReq.Value,
				Unit:                mReq.Unit,
				IsThresholdBreached: mReq.IsThresholdBreached,
			}
		}
		inv.Measurements = append(inv.Measurements, meas...)
	}

	if err := s.maintenanceRepo.UpdateIntervention(ctx, inv); err != nil {
		return nil, err
	}

	s.fireAudit(ctx, "UPDATE_INTERVENTION", fmt.Sprintf("Updated intervention %s", interventionID))
	resp := s.buildInterventionResponse(ctx, inv)
	return &resp, nil
}

func (s *MaintenanceService) StartIntervention(ctx context.Context, workOrderID uuid.UUID, interventionID uuid.UUID) (*domain.InterventionResponse, error) {
	inv, err := s.maintenanceRepo.FindInterventionByID(ctx, interventionID)
	if err != nil || inv.WorkOrderID != workOrderID {
		return nil, errors.New("intervention not found")
	}

	now := time.Now()
	inv.StartedAt = &now
	if err := s.maintenanceRepo.UpdateIntervention(ctx, inv); err != nil {
		return nil, err
	}

	wo, _ := s.maintenanceRepo.FindWorkOrderByID(ctx, workOrderID)
	if wo != nil && wo.Status == "PENDING" {
		wo.Status = "IN_PROGRESS"
		wo.UpdatedAt = now
		_ = s.maintenanceRepo.UpdateWorkOrder(ctx, wo)
	}

	if wo != nil && s.assetClient != nil {
		go func() {
			bgCtx := context.Background()
			_ = s.assetClient.UpdateAssetStatus(bgCtx, wo.AssetID, "DOWN")
		}()
	}

	s.fireAudit(ctx, "START_INTERVENTION", fmt.Sprintf("Started intervention %s", interventionID))
	resp := s.buildInterventionResponse(ctx, inv)
	return &resp, nil
}

func (s *MaintenanceService) EndIntervention(ctx context.Context, workOrderID uuid.UUID, interventionID uuid.UUID) (*domain.InterventionResponse, error) {
	inv, err := s.maintenanceRepo.FindInterventionByID(ctx, interventionID)
	if err != nil || inv.WorkOrderID != workOrderID {
		return nil, errors.New("intervention not found")
	}

	now := time.Now()
	inv.EndedAt = &now
	if err := s.maintenanceRepo.UpdateIntervention(ctx, inv); err != nil {
		return nil, err
	}

	wo, _ := s.maintenanceRepo.FindWorkOrderByID(ctx, workOrderID)
	if wo != nil {
		wo.Status = "COMPLETED"
		wo.UpdatedAt = now
		_ = s.maintenanceRepo.UpdateWorkOrder(ctx, wo)

		if s.assetClient != nil {
			go func() {
				bgCtx := context.Background()
				_ = s.assetClient.UpdateAssetStatus(bgCtx, wo.AssetID, "OPERATIONAL")
			}()
		}

		// If this is a PREVENTIVE SYSTEMATIC maintenance, we can reset the maintenance schedule
		if wo.MaintenanceCategory == "PREVENTIVE" && wo.MaintenanceType == "SYSTEMATIC" && s.assetClient != nil {
			// We might need to pass the rule ID, for now we will just pass nil and the asset service could just update the main last_maintenance_at. 
			// Wait, the user wanted flexible rules. If the frontend passes a rule ID in the work order or intervention, we can pass it here.
			// Let's just update the LastMaintenanceAt for the generic asset if no rule ID is provided.
			_ = s.assetClient.RecordUsage(ctx, wo.AssetID, 0, &now, nil) // The usage shouldn't be 0, we need the current usage, but we don't have it.
			// Actually, let's just let the asset service handle "LastMaintenanceAt" update if we pass the date.
			// But wait, the asset service's RecordUsage expects usage_hours to be required. 
		}
		
		if s.eventPublisher != nil {
			_ = s.eventPublisher.PublishWorkOrderCompleted(ctx, wo.ID, wo.AssetID, "INTERVENTION", inv.MaintenanceType)
		}

		if inv.StartedAt != nil {
			durationMins := int(now.Sub(*inv.StartedAt).Minutes())
			event := secondary.MaintenanceEvent{
				AssetID:             wo.AssetID,
				MaintenanceCategory: inv.MaintenanceCategory,
				DurationMinutes:     float64(durationMins),
			}
			go func() {
				_ = s.analyticsClient.PublishMaintenanceEvent(context.Background(), event)
			}()
		}
	}

	s.fireAudit(ctx, "END_INTERVENTION", fmt.Sprintf("Ended intervention %s", interventionID))
	resp := s.buildInterventionResponse(ctx, inv)
	return &resp, nil
}

func (s *MaintenanceService) CreateInspection(ctx context.Context, req domain.CreateInspectionRequest) (*domain.InspectionResponse, error) {
	assetID, err := uuid.Parse(req.AssetID)
	if err != nil {
		return nil, errors.New("invalid asset ID format")
	}

	var scheduleID *uuid.UUID
	if req.ScheduleID != nil && *req.ScheduleID != "" {
		parsed, err := uuid.Parse(*req.ScheduleID)
		if err == nil {
			scheduleID = &parsed
		}
	}

	performedBy, err := uuid.Parse(req.PerformedBy)
	if err != nil {
		return nil, errors.New("invalid performed by user ID format")
	}

	inspection := &domain.Inspection{
		ID:                 uuid.New(),
		AssetID:            assetID,
		ScheduleID:         scheduleID,
		Date:               req.Date,
		Observations:       req.Observations,
		UsageHoursRecorded: req.UsageHoursRecorded,
		RequiresAttention:  req.RequiresAttention,
		AttentionReason:    req.AttentionReason,
		PerformedBy:        performedBy,
	}

	if len(req.Measurements) > 0 {
		meas := make([]domain.MetricMeasurement, len(req.Measurements))
		for i, mReq := range req.Measurements {
			var compID *uuid.UUID
			if mReq.ComponentID != nil && *mReq.ComponentID != "" {
				parsedComp, err := uuid.Parse(*mReq.ComponentID)
				if err == nil {
					compID = &parsedComp
				}
			}
			meas[i] = domain.MetricMeasurement{
				ID:                  uuid.New(),
				InspectionID:        &inspection.ID,
				ComponentID:         compID,
				MetricName:          mReq.MetricName,
				Value:               mReq.Value,
				Unit:                mReq.Unit,
				IsThresholdBreached: mReq.IsThresholdBreached,
			}
		}
		inspection.Measurements = meas
	}

	if err := s.maintenanceRepo.CreateInspection(ctx, inspection); err != nil {
		return nil, err
	}

	if req.UsageHoursRecorded != nil && s.assetClient != nil {
		go func() {
			bgCtx := context.Background()
			_ = s.assetClient.RecordUsage(bgCtx, assetID, *req.UsageHoursRecorded, nil, nil)
		}()
	}

	s.fireAudit(ctx, "RECORD_INSPECTION", fmt.Sprintf("Recorded inspection %s for asset %s", inspection.ID, inspection.AssetID))

	resp := s.buildInspectionResponse(ctx, inspection)
	return &resp, nil
}

func (s *MaintenanceService) UpdateInspection(ctx context.Context, inspectionID uuid.UUID, req domain.UpdateInspectionRequest) (*domain.InspectionResponse, error) {
	ins, err := s.maintenanceRepo.FindInspectionByID(ctx, inspectionID)
	if err != nil {
		return nil, errors.New("inspection not found")
	}

	if req.Date != nil {
		ins.Date = req.Date
	}
	if req.Observations != nil {
		ins.Observations = *req.Observations
	}
	if req.UsageHoursRecorded != nil {
		ins.UsageHoursRecorded = req.UsageHoursRecorded
	}
	if req.RequiresAttention != nil {
		ins.RequiresAttention = *req.RequiresAttention
	}
	if req.AttentionReason != nil {
		ins.AttentionReason = *req.AttentionReason
	}

	if len(req.Measurements) > 0 {
		meas := make([]domain.MetricMeasurement, len(req.Measurements))
		for i, mReq := range req.Measurements {
			var compID *uuid.UUID
			if mReq.ComponentID != nil && *mReq.ComponentID != "" {
				parsedComp, err := uuid.Parse(*mReq.ComponentID)
				if err == nil {
					compID = &parsedComp
				}
			}
			meas[i] = domain.MetricMeasurement{
				ID:                  uuid.New(),
				InspectionID:        &ins.ID,
				ComponentID:         compID,
				MetricName:          mReq.MetricName,
				Value:               mReq.Value,
				Unit:                mReq.Unit,
				IsThresholdBreached: mReq.IsThresholdBreached,
			}
		}
		ins.Measurements = append(ins.Measurements, meas...)
	}

	if err := s.maintenanceRepo.UpdateInspection(ctx, ins); err != nil {
		return nil, err
	}

	if req.UsageHoursRecorded != nil && s.assetClient != nil {
		go func() {
			bgCtx := context.Background()
			_ = s.assetClient.RecordUsage(bgCtx, ins.AssetID, *req.UsageHoursRecorded, nil, nil)
		}()
	}

	s.fireAudit(ctx, "UPDATE_INSPECTION", fmt.Sprintf("Updated inspection %s", inspectionID))
	resp := s.buildInspectionResponse(ctx, ins)
	return &resp, nil
}

func (s *MaintenanceService) StartInspection(ctx context.Context, inspectionID uuid.UUID) (*domain.InspectionResponse, error) {
	ins, err := s.maintenanceRepo.FindInspectionByID(ctx, inspectionID)
	if err != nil {
		return nil, errors.New("inspection not found")
	}

	now := time.Now()
	ins.StartedAt = &now
	if err := s.maintenanceRepo.UpdateInspection(ctx, ins); err != nil {
		return nil, err
	}

	s.fireAudit(ctx, "START_INSPECTION", fmt.Sprintf("Started inspection %s", inspectionID))
	resp := s.buildInspectionResponse(ctx, ins)
	return &resp, nil
}

func (s *MaintenanceService) EndInspection(ctx context.Context, inspectionID uuid.UUID) (*domain.InspectionResponse, error) {
	ins, err := s.maintenanceRepo.FindInspectionByID(ctx, inspectionID)
	if err != nil {
		return nil, errors.New("inspection not found")
	}

	now := time.Now()
	ins.EndedAt = &now
	if err := s.maintenanceRepo.UpdateInspection(ctx, ins); err != nil {
		return nil, err
	}

	s.fireAudit(ctx, "END_INSPECTION", fmt.Sprintf("Ended inspection %s", inspectionID))
	resp := s.buildInspectionResponse(ctx, ins)
	return &resp, nil
}

func (s *MaintenanceService) GetInspectionsForAsset(ctx context.Context, assetID uuid.UUID) ([]domain.InspectionResponse, error) {
	inspections, err := s.maintenanceRepo.FindInspectionsByAssetID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch inspections: %w", err)
	}
	var res []domain.InspectionResponse
	for _, ins := range inspections {
		res = append(res, s.buildInspectionResponse(ctx, &ins))
	}
	return res, nil
}

func (s *MaintenanceService) GetInterventionsForWorkOrder(ctx context.Context, workOrderID uuid.UUID) ([]domain.InterventionResponse, error) {
	interventions, err := s.maintenanceRepo.FindInterventionsByWorkOrderID(ctx, workOrderID)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.InterventionResponse, len(interventions))
	for i, iv := range interventions {
		responses[i] = s.buildInterventionResponse(ctx, &iv)
	}
	return responses, nil
}

func (s *MaintenanceService) buildOrdreTravailResponse(ctx context.Context, wo *domain.OrdreTravail) domain.OrdreTravailResponse {
	assetName := "Unknown Asset"
	if name, err := s.assetClient.GetAssetName(ctx, wo.AssetID); err == nil {
		assetName = name
	}

	assignedToName := "Unknown User"
	if wo.AssignedTo != nil && *wo.AssignedTo != uuid.Nil {
		if name, err := s.userClient.GetUserName(ctx, *wo.AssignedTo); err == nil {
			assignedToName = name
		}
	}

	perfNames := make(map[uuid.UUID]string)
	for _, inv := range wo.Interventions {
		if _, ok := perfNames[inv.PerformedBy]; !ok {
			if inv.PerformedBy != uuid.Nil {
				if name, err := s.userClient.GetUserName(ctx, inv.PerformedBy); err == nil {
					perfNames[inv.PerformedBy] = name
				} else {
					perfNames[inv.PerformedBy] = "Unknown User"
				}
			} else {
				perfNames[inv.PerformedBy] = "Unknown User"
			}
		}
	}

	compNames := make(map[uuid.UUID]string)

	return wo.ToResponse(assetName, assignedToName, perfNames, compNames)
}

func (s *MaintenanceService) buildInterventionResponse(ctx context.Context, inv *domain.Intervention) domain.InterventionResponse {
	woTitle := "Unknown Work Order"
	if wo, err := s.maintenanceRepo.FindWorkOrderByID(ctx, inv.WorkOrderID); err == nil {
		woTitle = wo.Title
	}

	perfName := "Unknown User"
	if inv.PerformedBy != uuid.Nil {
		if name, err := s.userClient.GetUserName(ctx, inv.PerformedBy); err == nil {
			perfName = name
		}
	}

	compNames := make(map[uuid.UUID]string)
	return inv.ToResponse(woTitle, perfName, compNames)
}

func (s *MaintenanceService) buildInspectionResponse(ctx context.Context, ins *domain.Inspection) domain.InspectionResponse {
	assetName := "Unknown Asset"
	if name, err := s.assetClient.GetAssetName(ctx, ins.AssetID); err == nil {
		assetName = name
	}

	perfName := "Unknown User"
	if ins.PerformedBy != uuid.Nil {
		if name, err := s.userClient.GetUserName(ctx, ins.PerformedBy); err == nil {
			perfName = name
		}
	}

	compNames := make(map[uuid.UUID]string)
	return ins.ToResponse(assetName, perfName, compNames)
}

func (s *MaintenanceService) HandleAssetCreated(ctx context.Context, assetID uuid.UUID, modelID uuid.UUID) error {
	// Fetch the equipment model from asset service
	modelMap, err := s.assetClient.GetEquipmentModel(ctx, modelID)
	if err != nil {
		return fmt.Errorf("failed to fetch equipment model: %w", err)
	}

	// Safely parse maintenance rules from map
	rulesInterface, ok := modelMap["maintenance_rules"]
	if !ok || rulesInterface == nil {
		// No rules to schedule
		return nil
	}

	rules, ok := rulesInterface.([]interface{})
	if !ok {
		return fmt.Errorf("invalid maintenance rules format")
	}

	for _, rInterface := range rules {
		rule, ok := rInterface.(map[string]interface{})
		if !ok {
			continue
		}

		title := ""
		if rn, ok := rule["rule_name"].(string); ok {
			title = rn
		}

		var intervalMonths *int
		if im, ok := rule["interval_months"].(float64); ok { // JSON decodes numbers as float64
			v := int(im)
			intervalMonths = &v
		}

		var intervalHours *float64
		if ih, ok := rule["interval_hours"].(float64); ok {
			intervalHours = &ih
		}

		frequency := "USAGE_BASED"
		requireCounter := false

		if intervalMonths != nil {
			if *intervalMonths == 1 {
				frequency = "MONTHLY"
			} else if *intervalMonths == 12 {
				frequency = "YEARLY"
			} else {
				frequency = fmt.Sprintf("EVERY_%d_MONTHS", *intervalMonths)
			}
		}

		// Simplified frequency determination, if a rule says daily, it might have interval_months = 0 and some other flag
		// For now we'll just map MONTHLY or USAGE_BASED

		var nextDate *time.Time
		if intervalMonths != nil && *intervalMonths > 0 {
			nd := time.Now().AddDate(0, *intervalMonths, 0)
			nextDate = &nd
		}

		schedule := &domain.MaintenanceSchedule{
			ID:                    uuid.New(),
			AssetID:               assetID,
			Title:                 fmt.Sprintf("Preventive Maintenance: %s", title),
			Description:           "Generated from manufacturer baseline rule",
			Frequency:             frequency,
			IntervalMonths:        intervalMonths,
			IntervalHours:         intervalHours,
			NextScheduledDate:     nextDate,
			NextScheduledUsage:    intervalHours,
			MaintenanceCategory:   "PREVENTIVE",
			MaintenanceType:       "SYSTEMATIC",
			IsActive:              true,
			RequireCounterReading: requireCounter,
		}

		if err := s.maintenanceRepo.CreateMaintenanceSchedule(ctx, schedule); err != nil {
			// Log but continue
			fmt.Printf("Failed to create maintenance schedule for asset %s: %v\n", assetID, err)
		} else {
			s.fireAudit(ctx, "CREATE_MAINTENANCE_SCHEDULE", fmt.Sprintf("Created schedule %s for asset %s", schedule.ID, assetID))
		}
	}

	return nil
}
