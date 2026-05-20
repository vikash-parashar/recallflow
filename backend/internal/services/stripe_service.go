package services

import (
	"fmt"
	"time"

	"github.com/recallflow/backend/internal/config"
	"github.com/recallflow/backend/internal/models"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/subscription"
)

type StripeService struct {
	apiKey string
}

func NewStripeService(cfg *config.Config) *StripeService {
	stripe.Key = cfg.StripeSecretKey
	return &StripeService{
		apiKey: cfg.StripeSecretKey,
	}
}

// CreateCustomer creates a new Stripe customer
func (s *StripeService) CreateCustomer(email, name, orgID string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
		Metadata: map[string]string{
			"organization_id": orgID,
		},
	}

	cust, err := customer.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return cust, nil
}

// CreateSubscription creates a subscription for a customer
func (s *StripeService) CreateSubscription(customerID, priceID string, orgID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
		Metadata: map[string]string{
			"organization_id": orgID,
		},
	}

	sub, err := subscription.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return sub, nil
}

// CancelSubscription cancels a subscription
func (s *StripeService) CancelSubscription(subscriptionID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionCancelParams{}
	sub, err := subscription.Cancel(subscriptionID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel subscription: %w", err)
	}

	return sub, nil
}

// GetSubscription retrieves a subscription
func (s *StripeService) GetSubscription(subscriptionID string) (*stripe.Subscription, error) {
	sub, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return sub, nil
}

// UpdateSubscription updates a subscription (e.g., change plan)
func (s *StripeService) UpdateSubscription(subscriptionID, newPriceID string) (*stripe.Subscription, error) {
	// First get the subscription to get the current item
	sub, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	if len(sub.Items.Data) == 0 {
		return nil, fmt.Errorf("subscription has no items")
	}

	params := &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{
			{
				ID:    stripe.String(sub.Items.Data[0].ID),
				Price: stripe.String(newPriceID),
			},
		},
	}

	updatedSub, err := subscription.Update(subscriptionID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return updatedSub, nil
}

// GetPriceIDForPlan returns the Stripe price ID for a given plan type
func (s *StripeService) GetPriceIDForPlan(planType string) string {
	// In production, these should be environment variables
	// For now, using placeholder IDs
	priceIDs := map[string]string{
		"solo":           "price_solo_clinic_99",      // Replace with actual Stripe price ID
		"multi_provider": "price_multi_provider_299",  // Replace with actual Stripe price ID
		"multi_location": "price_multi_location_999",  // Replace with actual Stripe price ID
	}

	if priceID, ok := priceIDs[planType]; ok {
		return priceID
	}

	return priceIDs["solo"] // Default to solo plan
}

// SyncSubscriptionStatus syncs the subscription status from Stripe to local database
func (s *StripeService) SyncSubscriptionStatus(sub *stripe.Subscription) *models.Subscription {
	return &models.Subscription{
		StripeSubscriptionID: sub.ID,
		StripeCustomerID:     sub.Customer.ID,
		Status:               string(sub.Status),
		CurrentPeriodStart:   time.Unix(sub.CurrentPeriodStart, 0),
		CurrentPeriodEnd:     time.Unix(sub.CurrentPeriodEnd, 0),
	}
}
