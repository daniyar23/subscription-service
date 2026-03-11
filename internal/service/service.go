package service

import (
	"context"
	"fmt"
	"log"
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

	log.Println("service: create subscription", sub.UserID, sub.ServiceName)

	if sub.ServiceName == "" {
		log.Println("service: create validation error: empty service name")
		return nil, fmt.Errorf("service name required")
	}

	if sub.Price <= 0 {
		log.Println("service: create validation error: invalid price")
		return nil, fmt.Errorf("price must be positive")
	}

	result, err := s.repo.Create(ctx, sub)
	if err != nil {
		log.Println("service: repo create error:", err)
		return nil, err
	}

	return result, nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {

	log.Println("service: get subscription by id", id)

	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Println("service: repo getByID error:", err)
		return nil, err
	}

	return sub, nil
}

func (s *SubscriptionService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Subscription, error) {

	log.Println("service: get subscriptions by user", userID)

	subs, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		log.Println("service: repo getByUserID error:", err)
		return nil, err
	}

	return subs, nil
}

func (s *SubscriptionService) GetAll(ctx context.Context) ([]model.Subscription, error) {

	log.Println("service: get all subscriptions")

	subs, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Println("service: repo getAll error:", err)
		return nil, err
	}

	return subs, nil
}

func (s *SubscriptionService) Update(ctx context.Context, sub model.Subscription) error {

	log.Println("service: update subscription", sub.ID)

	if sub.ID == uuid.Nil {
		log.Println("service: update validation error: id required")
		return fmt.Errorf("id required")
	}

	err := s.repo.Update(ctx, sub)
	if err != nil {
		log.Println("service: repo update error:", err)
		return err
	}

	return nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id uuid.UUID) error {

	log.Println("service: delete subscription", id)

	if id == uuid.Nil {
		log.Println("service: delete validation error: invalid id")
		return fmt.Errorf("invalid id")
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Println("service: repo delete error:", err)
		return err
	}

	return nil
}

func (s *SubscriptionService) SumByFilter(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	from time.Time,
	to time.Time,
) (int, error) {

	log.Println("service: sum subscriptions", userID, serviceName)

	if userID == uuid.Nil {
		log.Println("service: sum validation error: user id required")
		return 0, fmt.Errorf("user id required")
	}

	sum, err := s.repo.SumByFilter(ctx, userID, serviceName, from, to)
	if err != nil {
		log.Println("service: repo sum error:", err)
		return 0, err
	}

	return sum, nil
}
