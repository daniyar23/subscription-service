package service

import (
	"context"
	"time"

	"github.com/daniyar23/subscribe-service/internal/model"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub model.Subscription) (*model.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Subscription, error)
	GetAll(ctx context.Context) ([]model.Subscription, error)
	Update(ctx context.Context, sub model.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	SumByFilter(ctx context.Context, userID uuid.UUID, serviceName string, from time.Time, to time.Time) (int, error)
}
