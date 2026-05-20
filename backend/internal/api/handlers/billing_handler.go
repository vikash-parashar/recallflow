package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"

	"github.com/recallflow/backend/internal/middleware"
	"github.com/recallflow/backend/internal/models"
	"github.com/recallflow/backend/internal/repositories"
	"github.com/recallflow/backend/internal/services"
)

type BillingHandler struct {
	stripeService *services.StripeService
	subRepo       *repositories.SubscriptionRepository
	orgRepo       *repositories.OrganizationRepository
	webhookSecret string
}

func NewBillingHandler(
	stripeService *services.StripeService,
	subRepo *repositories.SubscriptionRepository,
	orgRepo *repositories.OrganizationRepository,
	webhookSecret string,
) *BillingHandler {
	return &BillingHandler{
		stripeService: stripeService,
		subRepo:       subRepo,
		orgRepo:       orgRepo,
		webhookSecret: webhookSecret,
	}
}

type CreateSubscriptionRequest struct {
	PlanType string `json:"plan_type"` // solo, multi_provider, multi_location
}

type CreateSubscriptionResponse struct {
	SubscriptionID string `json:"subscription_id"`
	ClientSecret   string `json:"client_secret"`
	Status         string `json:"status"`
}

// CreateSubscription creates a new subscription for the organization
func (h *BillingHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate plan type
	if req.PlanType != "solo" && req.PlanType != "multi_provider" && req.PlanType != "multi_location" {
		RespondError(w, http.StatusBadRequest, "Invalid plan type")
		return
	}

	// Get organization
	org, err := h.orgRepo.GetByID(claims.OrganizationID)
	if err != nil {
		RespondError(w, http.StatusNotFound, "Organization not found")
		return
	}

	// Check if subscription already exists
	existingSub, _ := h.subRepo.GetByOrganizationID(claims.OrganizationID)
	if existingSub != nil {
		RespondError(w, http.StatusBadRequest, "Subscription already exists")
		return
	}

	// Create Stripe customer
	customer, err := h.stripeService.CreateCustomer(org.Email, org.Name, org.ID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to create customer")
		return
	}

	// Get price ID for plan
	priceID := h.stripeService.GetPriceIDForPlan(req.PlanType)

	// Create subscription
	stripeSub, err := h.stripeService.CreateSubscription(customer.ID, priceID, org.ID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to create subscription")
		return
	}

	// Save subscription to database
	subscription := &models.Subscription{
		ID:                   uuid.New().String(),
		OrganizationID:       org.ID,
		StripeCustomerID:     customer.ID,
		StripeSubscriptionID: stripeSub.ID,
		PlanType:             req.PlanType,
		Status:               string(stripeSub.Status),
		CurrentPeriodStart:   time.Unix(stripeSub.CurrentPeriodStart, 0),
		CurrentPeriodEnd:     time.Unix(stripeSub.CurrentPeriodEnd, 0),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := h.subRepo.Create(subscription); err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to save subscription")
		return
	}

	// Update organization plan type
	org.PlanType = req.PlanType

	response := CreateSubscriptionResponse{
		SubscriptionID: stripeSub.ID,
		ClientSecret:   stripeSub.LatestInvoice.PaymentIntent.ClientSecret,
		Status:         string(stripeSub.Status),
	}

	RespondJSON(w, http.StatusCreated, response)
}

// GetSubscription retrieves the current subscription
func (h *BillingHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	subscription, err := h.subRepo.GetByOrganizationID(claims.OrganizationID)
	if err != nil {
		RespondError(w, http.StatusNotFound, "No subscription found")
		return
	}

	RespondJSON(w, http.StatusOK, subscription)
}

// CancelSubscription cancels the current subscription
func (h *BillingHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	subscription, err := h.subRepo.GetByOrganizationID(claims.OrganizationID)
	if err != nil {
		RespondError(w, http.StatusNotFound, "No subscription found")
		return
	}

	_, err = h.stripeService.CancelSubscription(subscription.StripeSubscriptionID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to cancel subscription")
		return
	}

	// Update status in database
	if err := h.subRepo.UpdateStatus(subscription.StripeSubscriptionID, "canceled"); err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to update subscription status")
		return
	}

	RespondSuccess(w, "Subscription canceled successfully", nil)
}

// HandleStripeWebhook handles Stripe webhook events
func (h *BillingHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		RespondError(w, http.StatusServiceUnavailable, "Error reading request body")
		return
	}

	// Verify webhook signature
	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), h.webhookSecret)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid webhook signature")
		return
	}

	// Handle the event
	switch event.Type {
	case "customer.subscription.updated":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			RespondError(w, http.StatusBadRequest, "Error parsing webhook")
			return
		}
		h.handleSubscriptionUpdated(&subscription)

	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			RespondError(w, http.StatusBadRequest, "Error parsing webhook")
			return
		}
		h.handleSubscriptionDeleted(&subscription)

	case "invoice.payment_succeeded":
		// Handle successful payment
		// You might want to send a receipt email here

	case "invoice.payment_failed":
		// Handle failed payment
		// You might want to notify the customer here

	default:
		// Unhandled event type
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BillingHandler) handleSubscriptionUpdated(stripeSub *stripe.Subscription) {
	// Update subscription status in database
	h.subRepo.UpdateStatus(stripeSub.ID, string(stripeSub.Status))
}

func (h *BillingHandler) handleSubscriptionDeleted(stripeSub *stripe.Subscription) {
	// Mark subscription as canceled in database
	h.subRepo.UpdateStatus(stripeSub.ID, "canceled")
}
