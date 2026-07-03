package service

import (
	"context"
	"time"

	"backend-gmao/apps/audit-service/internal/core/domain"
	"backend-gmao/apps/audit-service/internal/core/ports/secondary"
	"github.com/google/uuid"
)

// AuditService implements primary.AuditUseCase.
type AuditService struct {
	auditRepo secondary.AuditRepository
}

// NewAuditService initializes a new AuditService instance.
func NewAuditService(auditRepo secondary.AuditRepository) *AuditService {
	return &AuditService{
		auditRepo: auditRepo,
	}
}

func (s *AuditService) RecordAction(ctx context.Context, log domain.AuditLog) error {
	// Ensure ID and timestamps are set correctly before saving
	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}
	if log.PerformedAt.IsZero() {
		log.PerformedAt = time.Now()
	}

	return s.auditRepo.Save(ctx, log)
}

func (s *AuditService) GetLogs(ctx context.Context, filter domain.AuditFilter) ([]domain.AuditLogResponse, int64, error) {
	logs, total, err := s.auditRepo.Find(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]domain.AuditLogResponse, len(logs))
	for i, l := range logs {
		responses[i] = l.ToResponse()
	}
	return responses, total, nil
}
