package service

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/daniyar23/subscribe-service/internal/model"
	"github.com/google/uuid"
)

// SubscriptionService provides methods for managing subscriptions.
type SubscriptionService struct {
	repo   SubscriptionRepository
	logger *zap.Logger
}

// NewSubscriptionService creates a new SubscriptionService with the given repository and logger.
func NewSubscriptionService(repo SubscriptionRepository, logger *zap.Logger) *SubscriptionService {
	return &SubscriptionService{
		repo:   repo,
		logger: logger,
	}
}

// Create creates a new subscription in the database.
func (s *SubscriptionService) Create(ctx context.Context, sub model.Subscription) (*model.Subscription, error) {

	s.logger.Info("service: create subscription",
		zap.String("user_id", sub.UserID.String()),
		zap.String("service_name", sub.ServiceName),
	)

	if sub.ServiceName == "" {
		s.logger.Error("service: create validation error: empty service name")
		return nil, fmt.Errorf("service name required")
	}

	if sub.Price <= 0 {
		s.logger.Error("service: create validation error: invalid price")
		return nil, fmt.Errorf("price must be positive")
	}

	result, err := s.repo.Create(ctx, sub)
	if err != nil {
		s.logger.Error("service: repo create error", zap.Error(err))
		return nil, err
	}

	return result, nil
}

// GetByID retrieves a subscription by its ID.
func (s *SubscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {

	s.logger.Info("service: get subscription by id",
		zap.String("id", id.String()),
	)

	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("service: repo getByID error", zap.Error(err))
		return nil, err
	}

	return sub, nil
}

// GetByUserID retrieves all subscriptions for a given user.
func (s *SubscriptionService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Subscription, error) {

	s.logger.Info("service: get subscriptions by user",
		zap.String("user_id", userID.String()),
	)

	subs, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("service: repo getByUserID error", zap.Error(err))
		return nil, err
	}

	return subs, nil
}

// GetAll retrieves all subscriptions from the database.
func (s *SubscriptionService) GetAll(ctx context.Context) ([]model.Subscription, error) {

	s.logger.Info("service: get all subscriptions")

	subs, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error("service: repo getAll error", zap.Error(err))
		return nil, err
	}

	return subs, nil
}

// Update updates an existing subscription in the database.
func (s *SubscriptionService) Update(ctx context.Context, sub model.Subscription) error {

	s.logger.Info("service: update subscription",
		zap.String("id", sub.ID.String()),
	)

	if sub.ID == uuid.Nil {
		s.logger.Error("service: update validation error: id required")
		return fmt.Errorf("id required")
	}

	err := s.repo.Update(ctx, sub)
	if err != nil {
		s.logger.Error("service: repo update error", zap.Error(err))
		return err
	}

	return nil
}

// Delete deletes a subscription by its ID.
func (s *SubscriptionService) Delete(ctx context.Context, id uuid.UUID) error {

	s.logger.Info("service: delete subscription",
		zap.String("id", id.String()),
	)

	if id == uuid.Nil {
		s.logger.Error("service: delete validation error: invalid id")
		return fmt.Errorf("invalid id")
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("service: repo delete error", zap.Error(err))
		return err
	}

	return nil
}

// SumByFilter sums the subscriptions for a given user and service within a time range.
func (s *SubscriptionService) SumByFilter(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	from time.Time,
	to time.Time,
) (int, error) {

	s.logger.Info("service: sum subscriptions",
		zap.String("user_id", userID.String()),
		zap.String("service_name", serviceName),
	)

	if userID == uuid.Nil {
		s.logger.Error("service: sum validation error: user id required")
		return 0, fmt.Errorf("user id required")
	}

	sum, err := s.repo.SumByFilter(ctx, userID, serviceName, from, to)
	if err != nil {
		s.logger.Error("service: repo sum error", zap.Error(err))
		return 0, err
	}

	return sum, nil
}
