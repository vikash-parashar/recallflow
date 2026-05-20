package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/recallflow/backend/internal/models"
	"github.com/recallflow/backend/internal/repositories"
)

type ConversationService struct {
	conversationRepo *repositories.ConversationRepository
	smsRepo          *repositories.SMSRepository
	openaiService    *OpenAIService
	twilioService    *TwilioService
}

func NewConversationService(
	conversationRepo *repositories.ConversationRepository,
	smsRepo *repositories.SMSRepository,
	openaiService *OpenAIService,
	twilioService *TwilioService,
) *ConversationService {
	return &ConversationService{
		conversationRepo: conversationRepo,
		smsRepo:          smsRepo,
		openaiService:    openaiService,
		twilioService:    twilioService,
	}
}

func (s *ConversationService) CreateConversation(callID, orgID, locationID, patientPhone string) (*models.Conversation, error) {
	conversation := &models.Conversation{
		ID:             uuid.New().String(),
		CallID:         callID,
		OrganizationID: orgID,
		LocationID:     locationID,
		PatientPhone:   patientPhone,
		Status:         "active",
		IsResolved:     false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.conversationRepo.Create(conversation); err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return conversation, nil
}

func (s *ConversationService) ProcessInboundSMS(conversationID, messageSID, fromNumber, toNumber, body string) error {
	// Save the inbound SMS
	message := &models.SMSMessage{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		TwilioMessageSID: messageSID,
		Direction:      "inbound",
		FromNumber:     fromNumber,
		ToNumber:       toNumber,
		Body:           body,
		Status:         "received",
		CreatedAt:      time.Now(),
	}

	if err := s.smsRepo.Create(message); err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	// Classify intent with AI
	ctx := context.Background()
	classification, err := s.openaiService.ClassifyIntent(ctx, body)
	if err != nil {
		// Continue even if AI fails
		fmt.Printf("AI classification failed: %v\n", err)
		classification = &IntentClassification{
			Intent:          "general",
			RequiresStaff:   true,
			SuggestedAction: "Review manually",
		}
	}

	// Update conversation with intent
	if err := s.conversationRepo.UpdateIntent(conversationID, classification.Intent); err != nil {
		return fmt.Errorf("failed to update intent: %w", err)
	}

	// If emergency, notify staff immediately
	if classification.IsEmergency {
		// TODO: Send urgent notification
		fmt.Printf("EMERGENCY detected in conversation %s\n", conversationID)
	}

	// Generate and send AI response
	response, err := s.openaiService.GenerateResponse(ctx, classification.Intent, body)
	if err != nil {
		response = "Thank you for your message. Our team will contact you shortly."
	}

	// Send response via Twilio
	responseSID, err := s.twilioService.SendSMS(fromNumber, response)
	if err != nil {
		return fmt.Errorf("failed to send response SMS: %w", err)
	}

	// Save outbound SMS
	outboundMessage := &models.SMSMessage{
		ID:               uuid.New().String(),
		ConversationID:   conversationID,
		TwilioMessageSID: responseSID,
		Direction:        "outbound",
		FromNumber:       toNumber,
		ToNumber:         fromNumber,
		Body:             response,
		Status:           "sent",
		CreatedAt:        time.Now(),
	}

	if err := s.smsRepo.Create(outboundMessage); err != nil {
		return fmt.Errorf("failed to save outbound message: %w", err)
	}

	return nil
}

func (s *ConversationService) GetConversationsByOrg(orgID string) ([]*models.Conversation, error) {
	return s.conversationRepo.GetByOrganization(orgID)
}

func (s *ConversationService) GetConversation(conversationID, orgID string) (*models.Conversation, error) {
	return s.conversationRepo.GetByID(conversationID, orgID)
}

func (s *ConversationService) GetMessages(conversationID, orgID string) ([]*models.SMSMessage, error) {
	return s.smsRepo.GetByConversation(conversationID)
}

func (s *ConversationService) ResolveConversation(conversationID, orgID string) error {
	return s.conversationRepo.MarkResolved(conversationID, orgID)
}
