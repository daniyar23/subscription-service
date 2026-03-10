package service

import (
	"context"
	"fmt"
	"time"

	"github.com/daniyar23/subscribe-service/internal/model"
	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo SubscriptionRepository
}

func NewSubscriptionService(repo SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
	}
}

func (s *SubscriptionService) Create(ctx context.Context, sub model.Subscription) (*model.Subscription, error) {

	if sub.ServiceName == "" {
		return nil, fmt.Errorf("service name required")
	}

	if sub.Price <= 0 {
		return nil, fmt.Errorf("price must be positive")
	}

	return s.repo.Create(ctx, sub)
}

func (s *SubscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SubscriptionService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Subscription, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *SubscriptionService) GetAll(ctx context.Context) ([]model.Subscription, error) {
	return s.repo.GetAll(ctx)
}

func (s *SubscriptionService) Update(ctx context.Context, sub model.Subscription) error {

	if sub.ID == uuid.Nil {
		return fmt.Errorf("id required")
	}

	return s.repo.Update(ctx, sub)
}

func (s *SubscriptionService) Delete(ctx context.Context, id uuid.UUID) error {

	if id == uuid.Nil {
		return fmt.Errorf("invalid id")
	}

	return s.repo.Delete(ctx, id)
}

func (s *SubscriptionService) SumByFilter(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	from time.Time,
	to time.Time,
) (int, error) {

	if userID == uuid.Nil {
		return 0, fmt.Errorf("user id required")
	}

	return s.repo.SumByFilter(ctx, userID, serviceName, from, to)
}
