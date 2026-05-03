package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/daniyar23/subscribe-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// SubscriptionRepository handles database operations for subscriptions.
type SubscriptionRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

// NewSubscriptionRepo creates a new SubscriptionRepository.
func NewSubscriptionRepo(db *pgxpool.Pool, logger *zap.Logger) *SubscriptionRepository {
	return &SubscriptionRepository{
		db:     db,
		logger: logger,
	}
}

// Create inserts a new subscription into the database.
func (r *SubscriptionRepository) Create(ctx context.Context, sub model.Subscription) (*model.Subscription, error) {

	r.logger.Info("repo: create subscription",
		zap.String("user_id", sub.UserID.String()),
		zap.String("service_name", sub.ServiceName),
	)

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
		r.logger.Error("repo: create subscription error",
			zap.Error(err),
		)
		return nil, fmt.Errorf("Create repo error: %w", err)
	}

	return &sub, nil
}

// GetByID retrieves a subscription by its ID.
func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {

	r.logger.Info("repo: get subscription by id",
		zap.String("id", id.String()),
	)

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
		r.logger.Error("repo: get by id error",
			zap.Error(err),
		)
		return nil, fmt.Errorf("GetByID repo error: %w", err)
	}

	return &sub, nil
}

// GetByUserID retrieves all subscriptions for a given user.
func (r *SubscriptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Subscription, error) {

	r.logger.Info("repo: get subscriptions by user",
		zap.String("user_id", userID.String()),
	)

	query := `
	SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions
	WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		r.logger.Error("repo: get by user id error",
			zap.Error(err),
		)
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
			r.logger.Error("repo: scan error",
				zap.Error(err),
			)
			return nil, err
		}

		subs = append(subs, sub)
	}

	return subs, nil
}

// GetAll retrieves all subscriptions from the database.
func (r *SubscriptionRepository) GetAll(ctx context.Context) ([]model.Subscription, error) {

	r.logger.Info("repo: get all subscriptions")

	query := `
	SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error("repo: get all error",
			zap.Error(err),
		)
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
			r.logger.Error("repo: scan error",
				zap.Error(err),
			)
			return nil, err
		}

		subs = append(subs, sub)
	}

	return subs, nil
}

// Update updates an existing subscription in the database.
func (r *SubscriptionRepository) Update(ctx context.Context, sub model.Subscription) error {

	r.logger.Info("repo: update subscription",
		zap.String("id", sub.ID.String()),
	)

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
		r.logger.Error("repo: update error",
			zap.Error(err),
		)
		return fmt.Errorf("Update repo error: %w", err)
	}

	return nil
}

// Delete removes a subscription from the database by its ID.
func (r *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {

	r.logger.Info("repo: delete subscription",
		zap.String("id", id.String()),
	)

	query := `
	DELETE FROM subscriptions
	WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("repo: delete error",
			zap.Error(err),
		)
		return fmt.Errorf("Delete repo error: %w", err)
	}

	return nil
}

// SumByFilter calculates the sum of subscription prices for a given filter.
func (r *SubscriptionRepository) SumByFilter(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	from time.Time,
	to time.Time,
) (int, error) {

	r.logger.Info("repo: sum by filter",
		zap.String("user_id", userID.String()),
		zap.String("service_name", serviceName),
	)

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
		r.logger.Error("repo: sum query error",
			zap.Error(err),
		)
		return 0, err
	}

	return sum, nil
}
