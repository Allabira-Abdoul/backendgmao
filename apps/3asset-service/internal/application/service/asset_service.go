package service

import (
	"context"
	"errors"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"backend-gmao/apps/asset-service/internal/core/ports/secondary"
	"backend-gmao/pkg/middleware"
	"fmt"
	"github.com/google/uuid"
)

type assetService struct {
	repo           secondary.AssetRepository
	eventPublisher secondary.EventPublisher
}

func NewAssetService(repo secondary.AssetRepository, eventPublisher secondary.EventPublisher) *assetService {
	return &assetService{repo: repo, eventPublisher: eventPublisher}
}

func getUserIDFromContext(ctx context.Context) *uuid.UUID {
	userIDStr, ok := ctx.Value(middleware.ContextKeyUserID).(string)
	if ok && userIDStr != "" {
		uid, err := uuid.Parse(userIDStr)
		if err == nil {
			return &uid
		}
	}
	return nil
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

func (s *assetService) GetEquipmentModelByID(ctx context.Context, id uuid.UUID) (domain.EquipmentModelResponse, error) {
	model, err := s.repo.GetEquipmentModelByID(ctx, id)
	if err != nil || model == nil {
		return domain.EquipmentModelResponse{}, errors.New("equipment model not found")
	}

	return model.ToResponse(), nil
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
				s.eventPublisher.PublishAuditLog(ctx, "CREATE", "PART_INSTANCE", partInst.ID.String(), getUserIDFromContext(ctx), map[string]interface{}{
					"auto_generated":        true,
					"equipment_instance_id": instance.ID.String(),
				})
			}
		}
	}

	// Emit Audit Log
	s.eventPublisher.PublishAuditLog(ctx, "CREATE", "EQUIPMENT_INSTANCE", instance.ID.String(), getUserIDFromContext(ctx), map[string]interface{}{
		"code":     instance.Code,
		"model_id": instance.EquipmentModelID.String(),
	})

	// Emit Domain Event for Analytics
	s.eventPublisher.PublishAssetCreated(ctx, instance.ID, instance.EquipmentModelID, model.Category, []string{}, getUserIDFromContext(ctx))

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

	s.eventPublisher.PublishAuditLog(ctx, "CREATE", "PART_INSTANCE", instance.ID.String(), getUserIDFromContext(ctx), map[string]interface{}{
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

	s.eventPublisher.PublishAuditLog(ctx, "STATUS_CHANGE", "EQUIPMENT_INSTANCE", instance.ID.String(), getUserIDFromContext(ctx), map[string]interface{}{
		"old_status": oldStatus,
		"new_status": newStatus,
	})

	// Emit Domain Event for Analytics
	if oldStatus != newStatus {
		s.eventPublisher.PublishAssetStateChanged(ctx, instance.ID, oldStatus, newStatus, getUserIDFromContext(ctx))
	}

	return nil
}

func (s *assetService) UpdateEquipmentLocation(ctx context.Context, id uuid.UUID, newLocation string) error {
	if !domain.IsValidLocation(newLocation) {
		return fmt.Errorf("invalid location: %s", newLocation)
	}

	instance, err := s.repo.GetEquipmentInstanceByID(ctx, id)
	if err != nil || instance == nil {
		return errors.New("equipment instance not found")
	}

	oldLocation := instance.Location
	instance.Location = newLocation

	if err := s.repo.UpdateEquipmentInstance(ctx, instance); err != nil {
		return err
	}

	for i := range instance.Parts {
		instance.Parts[i].CurrentLocation = newLocation
		if err := s.repo.UpdatePartInstance(ctx, &instance.Parts[i]); err != nil {
			return err
		}
	}

	s.eventPublisher.PublishAuditLog(ctx, "LOCATION_CHANGE", "EQUIPMENT_INSTANCE", instance.ID.String(), getUserIDFromContext(ctx), map[string]interface{}{
		"old_location": oldLocation,
		"new_location": newLocation,
	})

	return nil
}

func (s *assetService) MovePartInstance(ctx context.Context, partInstanceID uuid.UUID, req domain.MovePartInstanceRequest) (domain.PartInstanceResponse, error) {
	instance, err := s.repo.GetPartInstanceByID(ctx, partInstanceID)
	if err != nil || instance == nil {
		return domain.PartInstanceResponse{}, errors.New("part instance not found")
	}

	oldEq := instance.EquipmentInstanceID
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

	s.eventPublisher.PublishAuditLog(ctx, "MOVE", "PART_INSTANCE", instance.ID.String(), getUserIDFromContext(ctx), map[string]interface{}{
		"from_equipment": oldEq,
		"to_equipment":   eqID,
		"new_location":   req.CurrentLocation,
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

func (s *assetService) RecordUsage(ctx context.Context, id uuid.UUID, req domain.RecordUsageRequest) error {
	instance, err := s.repo.GetEquipmentInstanceByID(ctx, id)
	if err != nil || instance == nil {
		return errors.New("equipment instance not found")
	}

	instance.UsageHours = req.UsageHours

	if err := s.repo.UpdateEquipmentInstance(ctx, instance); err != nil {
		return err
	}

	return nil
}

// --- Consumables ---

func (s *assetService) CreateConsumable(ctx context.Context, req domain.CreateConsumableRequest) (domain.ConsumableResponse, error) {
	consumable := &domain.Consumable{
		ID:            uuid.New(),
		Name:          req.Name,
		Category:      req.Category,
		UnitOfMeasure: req.UnitOfMeasure,
		TotalStock:    0,
	}

	if err := s.repo.CreateConsumable(ctx, consumable); err != nil {
		return domain.ConsumableResponse{}, err
	}

	return consumable.ToResponse(), nil
}

func (s *assetService) GetConsumables(ctx context.Context) ([]domain.ConsumableResponse, error) {
	consumables, err := s.repo.GetConsumables(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]domain.ConsumableResponse, len(consumables))
	for i, c := range consumables {
		res[i] = c.ToResponse()
	}
	return res, nil
}

func (s *assetService) GetConsumableByID(ctx context.Context, id uuid.UUID) (domain.ConsumableResponse, error) {
	consumable, err := s.repo.GetConsumableByID(ctx, id)
	if err != nil || consumable == nil {
		return domain.ConsumableResponse{}, errors.New("consumable not found")
	}

	return consumable.ToResponse(), nil
}

func (s *assetService) AddConsumableStock(ctx context.Context, id uuid.UUID, req domain.AddConsumableStockRequest) (domain.ConsumableResponse, error) {
	consumable, err := s.repo.GetConsumableByID(ctx, id)
	if err != nil || consumable == nil {
		return domain.ConsumableResponse{}, errors.New("consumable not found")
	}

	stock, err := s.repo.GetConsumableStock(ctx, id, req.Location)
	if err != nil || stock == nil {
		stock = &domain.ConsumableLocationStock{
			ID:           uuid.New(),
			ConsumableID: id,
			Location:     req.Location,
			Quantity:     req.Quantity,
		}
		// Since we didn't add a CreateConsumableStock method explicitly, UpdateConsumableStock should use GORM Save/Upsert
		if err := s.repo.UpdateConsumableStock(ctx, stock); err != nil {
			return domain.ConsumableResponse{}, err
		}
	} else {
		stock.Quantity += req.Quantity
		if err := s.repo.UpdateConsumableStock(ctx, stock); err != nil {
			return domain.ConsumableResponse{}, err
		}
	}

	consumable.TotalStock += req.Quantity
	if err := s.repo.UpdateConsumable(ctx, consumable); err != nil {
		return domain.ConsumableResponse{}, err
	}

	// Refetch to get updated stock relationships
	consumable, _ = s.repo.GetConsumableByID(ctx, id)
	return consumable.ToResponse(), nil
}

func (s *assetService) ConsumeConsumable(ctx context.Context, req domain.ConsumeConsumableRequest, userID uuid.UUID) error {
	consumable, err := s.repo.GetConsumableByID(ctx, req.ConsumableID)
	if err != nil || consumable == nil {
		return errors.New("consumable not found")
	}

	stock, err := s.repo.GetConsumableStock(ctx, req.ConsumableID, req.Location)
	if err != nil || stock == nil {
		return errors.New("no stock found at this location")
	}

	if stock.Quantity < req.Quantity {
		return errors.New("insufficient consumable quantity at location")
	}

	stock.Quantity -= req.Quantity
	if err := s.repo.UpdateConsumableStock(ctx, stock); err != nil {
		return err
	}

	consumable.TotalStock -= req.Quantity
	if err := s.repo.UpdateConsumable(ctx, consumable); err != nil {
		return err
	}

	log := &domain.ConsumableConsumptionLog{
		ID:           uuid.New(),
		ConsumableID: req.ConsumableID,
		QuantityUsed: req.Quantity,
		WorkOrderID:  req.WorkOrderID,
		ConsumedBy:   userID,
		Notes:        req.Notes,
	}

	if err := s.repo.CreateConsumableConsumptionLog(ctx, log); err != nil {
		return err
	}

	s.eventPublisher.PublishAuditLog(ctx, "CONSUME", "CONSUMABLE", req.ConsumableID.String(), &userID, map[string]interface{}{
		"quantity_used": req.Quantity,
		"location":      req.Location,
		"work_order_id": req.WorkOrderID,
		"notes":         req.Notes,
	})

	return nil
}

