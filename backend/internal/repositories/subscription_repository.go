package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/recallflow/backend/internal/models"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(sub *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, organization_id, stripe_customer_id, stripe_subscription_id, plan_type, status, current_period_start, current_period_end, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(query,
		sub.ID,
		sub.OrganizationID,
		sub.StripeCustomerID,
		sub.StripeSubscriptionID,
		sub.PlanType,
		sub.Status,
		sub.CurrentPeriodStart,
		sub.CurrentPeriodEnd,
		sub.CreatedAt,
		sub.UpdatedAt,
	)
	return err
}

func (r *SubscriptionRepository) GetByOrganizationID(orgID string) (*models.Subscription, error) {
	query := `
		SELECT id, organization_id, stripe_customer_id, stripe_subscription_id, plan_type, status, current_period_start, current_period_end, created_at, updated_at
		FROM subscriptions
		WHERE organization_id = $1
	`

	sub := &models.Subscription{}
	
	err := r.db.QueryRow(query, orgID).Scan(
		&sub.ID,
		&sub.OrganizationID,
		&sub.StripeCustomerID,
		&sub.StripeSubscriptionID,
		&sub.PlanType,
		&sub.Status,
		&sub.CurrentPeriodStart,
		&sub.CurrentPeriodEnd,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subscription not found")
	}

	return sub, err
}

func (r *SubscriptionRepository) GetByStripeSubscriptionID(stripeSubID string) (*models.Subscription, error) {
	query := `
		SELECT id, organization_id, stripe_customer_id, stripe_subscription_id, plan_type, status, current_period_start, current_period_end, created_at, updated_at
		FROM subscriptions
		WHERE stripe_subscription_id = $1
	`

	sub := &models.Subscription{}
	
	err := r.db.QueryRow(query, stripeSubID).Scan(
		&sub.ID,
		&sub.OrganizationID,
		&sub.StripeCustomerID,
		&sub.StripeSubscriptionID,
		&sub.PlanType,
		&sub.Status,
		&sub.CurrentPeriodStart,
		&sub.CurrentPeriodEnd,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subscription not found")
	}

	return sub, err
}

func (r *SubscriptionRepository) Update(sub *models.Subscription) error {
	query := `
		UPDATE subscriptions
		SET stripe_customer_id = $1, stripe_subscription_id = $2, plan_type = $3, status = $4, 
		    current_period_start = $5, current_period_end = $6, updated_at = $7
		WHERE id = $8
	`
	_, err := r.db.Exec(query,
		sub.StripeCustomerID,
		sub.StripeSubscriptionID,
		sub.PlanType,
		sub.Status,
		sub.CurrentPeriodStart,
		sub.CurrentPeriodEnd,
		time.Now(),
		sub.ID,
	)
	return err
}

func (r *SubscriptionRepository) UpdateStatus(stripeSubID, status string) error {
	query := `
		UPDATE subscriptions
		SET status = $1, updated_at = $2
		WHERE stripe_subscription_id = $3
	`
	_, err := r.db.Exec(query, status, time.Now(), stripeSubID)
	return err
}

