package models

import (
	"time"
)

// Organization represents a clinic organization
type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	PlanType  string    `json:"plan_type"` // solo, multi_provider, multi_location
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// Location represents a clinic location
type Location struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	PhoneNumber    string    `json:"phone_number"`
	Address        string    `json:"address"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// User represents a user in the system
type User struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"` // Never expose in JSON
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Role           string    `json:"role"` // owner, admin, staff
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// Call represents an inbound call
type Call struct {
	ID             string    `json:"id"`
	LocationID     string    `json:"location_id"`
	OrganizationID string    `json:"organization_id"`
	TwilioCallSID  string    `json:"twilio_call_sid"`
	FromNumber     string    `json:"from_number"`
	ToNumber       string    `json:"to_number"`
	Status         string    `json:"status"` // missed, answered, voicemail
	Duration       int       `json:"duration"`
	CallTime       time.Time `json:"call_time"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Conversation represents an SMS conversation
type Conversation struct {
	ID             string    `json:"id"`
	CallID         string    `json:"call_id"`
	OrganizationID string    `json:"organization_id"`
	LocationID     string    `json:"location_id"`
	PatientPhone   string    `json:"patient_phone"`
	Status         string    `json:"status"` // active, resolved, escalated
	Intent         string    `json:"intent"` // appointment, billing, emergency, insurance, general
	Summary        string    `json:"summary"`
	IsResolved     bool      `json:"is_resolved"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// SMSMessage represents an individual SMS message
type SMSMessage struct {
	ID              string    `json:"id"`
	ConversationID  string    `json:"conversation_id"`
	TwilioMessageSID string   `json:"twilio_message_sid"`
	Direction       string    `json:"direction"` // inbound, outbound
	FromNumber      string    `json:"from_number"`
	ToNumber        string    `json:"to_number"`
	Body            string    `json:"body"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

// AIClassification represents AI analysis of a message
type AIClassification struct {
	ID              string    `json:"id"`
	ConversationID  string    `json:"conversation_id"`
	MessageID       string    `json:"message_id"`
	Intent          string    `json:"intent"`
	Confidence      float64   `json:"confidence"`
	IsEmergency     bool      `json:"is_emergency"`
	RequiresStaff   bool      `json:"requires_staff"`
	SuggestedAction string    `json:"suggested_action"`
	CreatedAt       time.Time `json:"created_at"`
}

// Notification represents a notification to staff
type Notification struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	UserID         string    `json:"user_id"`
	ConversationID string    `json:"conversation_id"`
	Type           string    `json:"type"` // missed_call, urgent, new_message
	Title          string    `json:"title"`
	Message        string    `json:"message"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
}

// Subscription represents billing subscription
type Subscription struct {
	ID                   string    `json:"id"`
	OrganizationID       string    `json:"organization_id"`
	StripeCustomerID     string    `json:"stripe_customer_id"`
	StripeSubscriptionID string    `json:"stripe_subscription_id"`
	PlanType             string    `json:"plan_type"`
	Status               string    `json:"status"` // active, canceled, past_due
	CurrentPeriodStart   time.Time `json:"current_period_start"`
	CurrentPeriodEnd     time.Time `json:"current_period_end"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// Analytics represents aggregated metrics
type Analytics struct {
	OrganizationID      string    `json:"organization_id"`
	Date                time.Time `json:"date"`
	TotalCalls          int       `json:"total_calls"`
	MissedCalls         int       `json:"missed_calls"`
	RecoveredLeads      int       `json:"recovered_leads"`
	SMSSent             int       `json:"sms_sent"`
	SMSResponses        int       `json:"sms_responses"`
	AppointmentsBooked  int       `json:"appointments_booked"`
	EstimatedRevenue    float64   `json:"estimated_revenue"`
	ResponseRate        float64   `json:"response_rate"`
}
