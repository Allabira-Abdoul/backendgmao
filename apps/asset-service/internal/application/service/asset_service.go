package service

import (
	"context"
	"errors"
	"time"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"backend-gmao/apps/asset-service/internal/core/ports/secondary"
	"fmt"
	"github.com/google/uuid"
)

type assetService struct {
	repo           secondary.AssetRepository
	measurementRepo secondary.MeasurementRepository
	eventPublisher secondary.EventPublisher
}

// NewAssetService creates a new Asset Service instance.
func NewAssetService(repo secondary.AssetRepository, measurementRepo secondary.MeasurementRepository, eventPublisher secondary.EventPublisher) *assetService {
	return &assetService{repo: repo, measurementRepo: measurementRepo, eventPublisher: eventPublisher}
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

	for _, pr := range req.PartRequirements {
		reqModel := &domain.EquipmentModelPartRequirement{
			ID:               uuid.New(),
			EquipmentModelID: model.ID,
			PartModelID:      pr.PartModelID,
			Quantity:         pr.Quantity,
		}
		_ = s.repo.CreateEquipmentModelPartRequirement(ctx, reqModel)
		model.PartRequirements = append(model.PartRequirements, *reqModel)
	}

	return model.ToResponse(), nil
}

func (s *assetService) UpdateEquipmentModel(ctx context.Context, id uuid.UUID, req domain.UpdateEquipmentModelRequest) (domain.EquipmentModelResponse, error) {
	model, err := s.repo.GetEquipmentModelByID(ctx, id)
	if err != nil || model == nil {
		return domain.EquipmentModelResponse{}, errors.New("equipment model not found")
	}

	if req.Name != nil { model.Name = *req.Name }
	if req.Category != nil { model.Category = *req.Category }
	if req.Description != nil { model.Description = *req.Description }

	if err := s.repo.UpdateEquipmentModel(ctx, model); err != nil {
		return domain.EquipmentModelResponse{}, err
	}

	if req.PartRequirements != nil {
		_ = s.repo.DeleteEquipmentModelPartRequirements(ctx, id)
		var newReqs []domain.EquipmentModelPartRequirement
		for _, pr := range req.PartRequirements {
			reqModel := &domain.EquipmentModelPartRequirement{
				ID:               uuid.New(),
				EquipmentModelID: id,
				PartModelID:      pr.PartModelID,
				Quantity:         pr.Quantity,
			}
			_ = s.repo.CreateEquipmentModelPartRequirement(ctx, reqModel)
			newReqs = append(newReqs, *reqModel)
		}
		model.PartRequirements = newReqs
	}

	return model.ToResponse(), nil
}

func (s *assetService) CreatePartModel(ctx context.Context, req domain.CreatePartModelRequest) (domain.PartModelResponse, error) {
	model := &domain.PartModel{
		ID:            uuid.New(),
		Name:          req.Name,
		Category:      req.Category,
		SpareQuantity: req.SpareQuantity,
		IsSerialized:  req.IsSerialized,
	}

	if err := s.repo.CreatePartModel(ctx, model); err != nil {
		return domain.PartModelResponse{}, err
	}

	return model.ToResponse(), nil
}

func (s *assetService) UpdatePartModel(ctx context.Context, id uuid.UUID, req domain.UpdatePartModelRequest) (domain.PartModelResponse, error) {
	model, err := s.repo.GetPartModelByID(ctx, id)
	if err != nil || model == nil {
		return domain.PartModelResponse{}, errors.New("part model not found")
	}

	if req.Name != nil { model.Name = *req.Name }
	if req.Category != nil { model.Category = *req.Category }
	if req.SpareQuantity != nil { model.SpareQuantity = *req.SpareQuantity }
	if req.IsSerialized != nil { model.IsSerialized = *req.IsSerialized }

	if err := s.repo.UpdatePartModel(ctx, model); err != nil {
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

// --- Suppliers ---

func (s *assetService) CreateSupplier(ctx context.Context, req domain.CreateSupplierRequest) (domain.SupplierResponse, error) {
	supplier := &domain.Supplier{
		ID:          uuid.New(),
		Name:        req.Name,
		ContactInfo: req.ContactInfo,
	}

	if err := s.repo.CreateSupplier(ctx, supplier); err != nil {
		return domain.SupplierResponse{}, err
	}

	return supplier.ToResponse(), nil
}

func (s *assetService) GetSuppliers(ctx context.Context) ([]domain.SupplierResponse, error) {
	suppliers, err := s.repo.GetSuppliers(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]domain.SupplierResponse, len(suppliers))
	for i, sup := range suppliers {
		res[i] = sup.ToResponse()
	}
	return res, nil
}

func (s *assetService) UpdateSupplier(ctx context.Context, id uuid.UUID, req domain.UpdateSupplierRequest) (domain.SupplierResponse, error) {
	supplier, err := s.repo.GetSupplierByID(ctx, id)
	if err != nil || supplier == nil {
		return domain.SupplierResponse{}, errors.New("supplier not found")
	}

	if req.Name != nil { supplier.Name = *req.Name }
	if req.ContactInfo != nil { supplier.ContactInfo = *req.ContactInfo }

	if err := s.repo.UpdateSupplier(ctx, supplier); err != nil {
		return domain.SupplierResponse{}, err
	}
	return supplier.ToResponse(), nil
}

func (s *assetService) DeleteSupplier(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteSupplier(ctx, id)
}

func (s *assetService) AddSupplierToEquipmentModel(ctx context.Context, modelID uuid.UUID, req domain.AddModelSupplierRequest) (domain.ModelSupplierResponse, error) {
	supplier, err := s.repo.GetSupplierByID(ctx, req.SupplierID)
	if err != nil || supplier == nil {
		return domain.ModelSupplierResponse{}, errors.New("supplier not found")
	}

	model, err := s.repo.GetEquipmentModelByID(ctx, modelID)
	if err != nil || model == nil {
		return domain.ModelSupplierResponse{}, errors.New("equipment model not found")
	}

	modelSupplier := &domain.ModelSupplier{
		ID:                    uuid.New(),
		SupplierID:            req.SupplierID,
		EquipmentModelID:      &modelID,
		SupplierReferenceCode: req.SupplierReferenceCode,
		TechnicalDocReference: req.TechnicalDocReference,
	}

	if err := s.repo.AddModelSupplier(ctx, modelSupplier); err != nil {
		return domain.ModelSupplierResponse{}, err
	}

	modelSupplier.Supplier = supplier // for response serialization
	return modelSupplier.ToResponse(), nil
}

func (s *assetService) AddSupplierToPartModel(ctx context.Context, modelID uuid.UUID, req domain.AddModelSupplierRequest) (domain.ModelSupplierResponse, error) {
	supplier, err := s.repo.GetSupplierByID(ctx, req.SupplierID)
	if err != nil || supplier == nil {
		return domain.ModelSupplierResponse{}, errors.New("supplier not found")
	}

	model, err := s.repo.GetPartModelByID(ctx, modelID)
	if err != nil || model == nil {
		return domain.ModelSupplierResponse{}, errors.New("part model not found")
	}

	modelSupplier := &domain.ModelSupplier{
		ID:                    uuid.New(),
		SupplierID:            req.SupplierID,
		PartModelID:           &modelID,
		SupplierReferenceCode: req.SupplierReferenceCode,
		TechnicalDocReference: req.TechnicalDocReference,
	}

	if err := s.repo.AddModelSupplier(ctx, modelSupplier); err != nil {
		return domain.ModelSupplierResponse{}, err
	}

	modelSupplier.Supplier = supplier // for response serialization
	return modelSupplier.ToResponse(), nil
}

func (s *assetService) CreateEquipmentInstance(ctx context.Context, req domain.CreateEquipmentInstanceRequest) (domain.EquipmentInstanceResponse, error) {
	// Verify model exists
	model, err := s.repo.GetEquipmentModelByID(ctx, req.EquipmentModelID)
	if err != nil || model == nil {
		return domain.EquipmentInstanceResponse{}, errors.New("equipment model not found")
	}

	if !domain.IsValidLocation(req.Location) {
		return domain.EquipmentInstanceResponse{}, fmt.Errorf("invalid location: %s", req.Location)
	}

	instance := &domain.EquipmentInstance{
		ID:               uuid.New(),
		Code:             req.Code,
		EquipmentModelID: req.EquipmentModelID,
		SupplierID:       req.SupplierID,
		Status:           "OPERATIONAL",
		Location:         req.Location,
	}

	if err := s.repo.CreateEquipmentInstance(ctx, instance); err != nil {
		return domain.EquipmentInstanceResponse{}, err
	}

	// Auto-generate required parts based on EquipmentModel blueprint
	for _, reqPart := range model.PartRequirements {
		for i := 0; i < reqPart.Quantity; i++ {
			partInst := &domain.PartInstance{
				ID:                  uuid.New(),
				EquipmentInstanceID: &instance.ID,
				PartModelID:         reqPart.PartModelID,
				SerialNumber:        fmt.Sprintf("SN-AUTO-%s-%d", uuid.New().String()[:8], i+1),
				Status:              "OPERATIONAL",
				CurrentLocation:     instance.Location,
			}
			if err := s.repo.CreatePartInstance(ctx, partInst); err == nil {
				s.eventPublisher.PublishAuditLog(ctx, "CREATE", "PART_INSTANCE", partInst.ID.String(), nil, map[string]interface{}{
					"auto_generated":        true,
					"equipment_instance_id": instance.ID.String(),
				})
			}
		}
	}

	// Emit Audit Log (ActorID is typically passed via context, but we can set it to nil for system actions if not present)
	s.eventPublisher.PublishAuditLog(ctx, "CREATE", "EQUIPMENT_INSTANCE", instance.ID.String(), nil, map[string]interface{}{
		"code":     instance.Code,
		"model_id": instance.EquipmentModelID.String(),
	})

	// Emit Domain Event for Analytics
	s.eventPublisher.PublishAssetCreated(ctx, instance.ID, instance.EquipmentModelID, model.Category, []string{})

	return instance.ToResponse(), nil
}

func (s *assetService) CreatePartInstance(ctx context.Context, req domain.CreatePartInstanceRequest) (domain.PartInstanceResponse, error) {
	// Verify equipment instance exists if provided
	var eqID *uuid.UUID
	if req.EquipmentInstanceID != nil && *req.EquipmentInstanceID != "" {
		parsed, err := uuid.Parse(*req.EquipmentInstanceID)
		if err == nil {
			_, err = s.repo.GetEquipmentInstanceByID(ctx, parsed)
			if err != nil {
				return domain.PartInstanceResponse{}, errors.New("equipment instance not found")
			}
			eqID = &parsed
		}
	}

	// Verify part model exists and is serialized
	partModel, err := s.repo.GetPartModelByID(ctx, req.PartModelID)
	if err != nil || partModel == nil {
		return domain.PartInstanceResponse{}, errors.New("part model not found")
	}
	if !partModel.IsSerialized {
		return domain.PartInstanceResponse{}, errors.New("part model is not serialized, cannot create instance")
	}

	instance := &domain.PartInstance{
		ID:                  uuid.New(),
		EquipmentInstanceID: eqID,
		PartModelID:         req.PartModelID,
		SupplierID:          req.SupplierID,
		SerialNumber:        req.SerialNumber,
		Status:              "OPERATIONAL",
		CurrentLocation:     req.CurrentLocation,
	}

	if err := s.repo.CreatePartInstance(ctx, instance); err != nil {
		return domain.PartInstanceResponse{}, err
	}

	s.eventPublisher.PublishAuditLog(ctx, "CREATE", "PART_INSTANCE", instance.ID.String(), nil, map[string]interface{}{
		"part_model_id":         instance.PartModelID.String(),
		"equipment_instance_id": eqID,
		"serial_number":         instance.SerialNumber,
	})

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

func (s *assetService) UpdateEquipmentStatus(ctx context.Context, id uuid.UUID, newStatus string) error {
	instance, err := s.repo.GetEquipmentInstanceByID(ctx, id)
	if err != nil || instance == nil {
		return errors.New("equipment instance not found")
	}

	oldStatus := instance.Status
	instance.Status = newStatus

	if err := s.repo.UpdateEquipmentInstance(ctx, instance); err != nil {
		return err
	}

	s.eventPublisher.PublishAuditLog(ctx, "STATUS_CHANGE", "EQUIPMENT_INSTANCE", instance.ID.String(), nil, map[string]interface{}{
		"old_status": oldStatus,
		"new_status": newStatus,
	})

	// Emit Domain Event for Analytics
	if oldStatus != newStatus {
		s.eventPublisher.PublishAssetStateChanged(ctx, instance.ID, oldStatus, newStatus)
	}

	return nil
}

func (s *assetService) MovePartInstance(ctx context.Context, partInstanceID uuid.UUID, req domain.MovePartInstanceRequest) (domain.PartInstanceResponse, error) {
	instance, err := s.repo.GetPartInstanceByID(ctx, partInstanceID)
	if err != nil || instance == nil {
		return domain.PartInstanceResponse{}, errors.New("part instance not found")
	}

	var eqID *uuid.UUID
	if req.EquipmentInstanceID != nil && *req.EquipmentInstanceID != "" {
		parsed, err := uuid.Parse(*req.EquipmentInstanceID)
		if err == nil {
			_, err = s.repo.GetEquipmentInstanceByID(ctx, parsed)
			if err != nil {
				return domain.PartInstanceResponse{}, errors.New("target equipment instance not found")
			}
			eqID = &parsed
		}
	}

	instance.EquipmentInstanceID = eqID
	instance.CurrentLocation = req.CurrentLocation

	if err := s.repo.UpdatePartInstance(ctx, instance); err != nil {
		return domain.PartInstanceResponse{}, err
	}

	s.eventPublisher.PublishAuditLog(ctx, "MOVE", "PART_INSTANCE", instance.ID.String(), nil, map[string]interface{}{
		"new_equipment_instance_id": eqID,
		"new_location":              req.CurrentLocation,
	})

	return instance.ToResponse(), nil
}

func (s *assetService) ConsumePart(ctx context.Context, req domain.ConsumePartRequest, userID uuid.UUID) error {
	partModel, err := s.repo.GetPartModelByID(ctx, req.PartModelID)
	if err != nil || partModel == nil {
		return errors.New("part model not found")
	}

	if partModel.IsSerialized {
		return errors.New("part model is serialized, must be managed via instances instead of consumed")
	}

	if partModel.SpareQuantity < req.Quantity {
		return errors.New("insufficient spare quantity")
	}

	// Decrement
	partModel.SpareQuantity -= req.Quantity
	if err := s.repo.UpdatePartModel(ctx, partModel); err != nil {
		return err
	}

	// Log it
	log := &domain.PartConsumptionLog{
		ID:           uuid.New(),
		PartModelID:  req.PartModelID,
		QuantityUsed: req.Quantity,
		WorkOrderID:  req.WorkOrderID,
		ConsumedBy:   userID,
		Notes:        req.Notes,
	}

	if err := s.repo.CreatePartConsumptionLog(ctx, log); err != nil {
		return err
	}

	s.eventPublisher.PublishAuditLog(ctx, "CONSUME", "PART_MODEL", req.PartModelID.String(), &userID, map[string]interface{}{
		"quantity_used": req.Quantity,
		"work_order_id": req.WorkOrderID,
		"notes":         req.Notes,
	})

	return nil
}

func (s *assetService) IngestMeasurement(ctx context.Context, req domain.IngestMeasurementRequest, userID *uuid.UUID) (domain.MeasurementResponse, error) {
	if req.EquipmentInstanceID == nil && req.PartInstanceID == nil {
		return domain.MeasurementResponse{}, errors.New("must specify either equipment_instance_id or part_instance_id")
	}

	recordedAt := time.Now()
	if req.RecordedAt != nil {
		recordedAt = *req.RecordedAt
	}

	measurement := &domain.Measurement{
		ID:                  uuid.New(),
		EquipmentInstanceID: req.EquipmentInstanceID,
		PartInstanceID:      req.PartInstanceID,
		MetricName:          req.MetricName,
		Value:               req.Value,
		Unit:                req.Unit,
		RecordedAt:          recordedAt,
		RecordedBy:          userID,
	}

	if err := s.measurementRepo.CreateMeasurement(ctx, measurement); err != nil {
		return domain.MeasurementResponse{}, err
	}

	// Fetch Thresholds
	thresholds, err := s.repo.GetMetricThresholds(ctx, req.MetricName, req.EquipmentInstanceID, req.PartInstanceID)
	if err != nil || len(thresholds) == 0 {
		return domain.MeasurementResponse{}, fmt.Errorf("metric not permitted: no threshold defined for metric name '%s'", req.MetricName)
	}

	// Evaluate the first matching threshold
	t := thresholds[0]
	breached := false
		reason := ""
		if t.MinValue != nil && req.Value < *t.MinValue {
			breached = true
			reason = "Below minimum threshold"
		} else if t.MaxValue != nil && req.Value > *t.MaxValue {
			breached = true
			reason = "Above maximum threshold"
		}

		if breached {
			s.eventPublisher.PublishAuditLog(ctx, "THRESHOLD_ALERT", "MEASUREMENT", measurement.ID.String(), nil, map[string]interface{}{
				"metric_name": req.MetricName,
				"value":       req.Value,
				"reason":      reason,
			})
		}

	return measurement.ToResponse(), nil
}

func (s *assetService) GetMeasurements(ctx context.Context, targetType string, targetID uuid.UUID, since string) ([]domain.MeasurementResponse, error) {
	var sinceTime time.Time
	if since != "" {
		parsed, err := time.Parse(time.RFC3339, since)
		if err == nil {
			sinceTime = parsed
		}
	} else {
		sinceTime = time.Now().AddDate(0, -1, 0) // Default to last 1 month
	}

	var measurements []domain.Measurement
	var err error

	if targetType == "equipment" {
		measurements, err = s.measurementRepo.GetMeasurementsByEquipment(ctx, targetID, sinceTime)
	} else if targetType == "part" {
		measurements, err = s.measurementRepo.GetMeasurementsByPart(ctx, targetID, sinceTime)
	} else {
		return nil, errors.New("invalid targetType: must be 'equipment' or 'part'")
	}

	if err != nil {
		return nil, err
	}

	res := make([]domain.MeasurementResponse, len(measurements))
	for i, m := range measurements {
		res[i] = m.ToResponse()
	}
	return res, nil
}

// --- Thresholds ---

func (s *assetService) CreateMetricThreshold(ctx context.Context, req domain.CreateMetricThresholdRequest) (domain.MetricThresholdResponse, error) {
	threshold := &domain.MetricThreshold{
		ID:                  uuid.New(),
		EquipmentModelID:    req.EquipmentModelID,
		PartModelID:         req.PartModelID,
		EquipmentInstanceID: req.EquipmentInstanceID,
		PartInstanceID:      req.PartInstanceID,
		MetricName:          req.MetricName,
		MinValue:            req.MinValue,
		MaxValue:            req.MaxValue,
		Unit:                req.Unit,
	}

	if err := s.repo.CreateMetricThreshold(ctx, threshold); err != nil {
		return domain.MetricThresholdResponse{}, err
	}
	return threshold.ToResponse(), nil
}

func (s *assetService) UpdateMetricThreshold(ctx context.Context, id uuid.UUID, req domain.UpdateMetricThresholdRequest) (domain.MetricThresholdResponse, error) {
	t, err := s.repo.GetMetricThresholdByID(ctx, id)
	if err != nil || t == nil {
		return domain.MetricThresholdResponse{}, errors.New("threshold not found")
	}

	if req.MinValue != nil { t.MinValue = req.MinValue }
	if req.MaxValue != nil { t.MaxValue = req.MaxValue }

	if err := s.repo.UpdateMetricThreshold(ctx, t); err != nil {
		return domain.MetricThresholdResponse{}, err
	}
	return t.ToResponse(), nil
}

func (s *assetService) DeleteMetricThreshold(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteMetricThreshold(ctx, id)
}

