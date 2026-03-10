package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/daniyar23/subscribe-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepo(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub model.Subscription) (*model.Subscription, error) {
	query := `
	INSERT INTO subscriptions
	(service_name, price, user_id, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(&sub.ID)

	if err != nil {
		return nil, fmt.Errorf("Create repo error: %w", err)
	}

	return &sub, nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	query := `
	SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions
	WHERE id = $1
	`

	var sub model.Subscription

	err := r.db.QueryRow(ctx, query, id).
		Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		)

	if err != nil {
		return nil, fmt.Errorf("GetByID repo error: %w", err)
	}

	return &sub, nil
}

func (r *SubscriptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Subscription, error) {
	query := `
	SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions
	WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetByUserID repo error: %w", err)
	}
	defer rows.Close()

	var subs []model.Subscription

	for rows.Next() {
		var sub model.Subscription

		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		)

		if err != nil {
			return nil, err
		}

		subs = append(subs, sub)
	}

	return subs, nil
}

func (r *SubscriptionRepository) GetAll(ctx context.Context) ([]model.Subscription, error) {
	query := `
	SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetAll repo error: %w", err)
	}
	defer rows.Close()

	var subs []model.Subscription

	for rows.Next() {
		var sub model.Subscription

		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		)

		if err != nil {
			return nil, err
		}

		subs = append(subs, sub)
	}

	return subs, nil
}
func (r *SubscriptionRepository) Update(ctx context.Context, sub model.Subscription) error {
	query := `
	UPDATE subscriptions
	SET service_name = $1,
	    price = $2,
	    user_id = $3,
	    start_date = $4,
	    end_date = $5,
	    updated_at = NOW()
	WHERE id = $6
	`

	_, err := r.db.Exec(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)

	if err != nil {
		return fmt.Errorf("Update repo error: %w", err)
	}

	return nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM subscriptions
	WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Delete repo error: %w", err)
	}

	return nil
}

func (r *SubscriptionRepository) SumByFilter(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	from time.Time,
	to time.Time,
) (int, error) {

	query := `
	SELECT COALESCE(SUM(price),0)
	FROM subscriptions
	WHERE user_id = $1
	AND service_name = $2
	AND start_date <= $4
	AND (end_date IS NULL OR end_date >= $3)
	`

	var sum int

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		serviceName,
		from,
		to,
	).Scan(&sum)

	if err != nil {
		return 0, err
	}

	return sum, nil
}
