package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/recallflow/backend/internal/models"
	"github.com/recallflow/backend/internal/repositories"
)

type CallService struct {
	callRepo            *repositories.CallRepository
	conversationService *ConversationService
}

func NewCallService(
	callRepo *repositories.CallRepository,
	conversationService *ConversationService,
) *CallService {
	return &CallService{
		callRepo:            callRepo,
		conversationService: conversationService,
	}
}

func (s *CallService) ProcessIncomingCall(callSID, fromNumber, toNumber, callStatus, locationID, orgID string) error {
	// Create call record
	call := &models.Call{
		ID:             uuid.New().String(),
		LocationID:     locationID,
		OrganizationID: orgID,
		TwilioCallSID:  callSID,
		FromNumber:     fromNumber,
		ToNumber:       toNumber,
		Status:         "pending", // Will be updated by status callback
		CallTime:       time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.callRepo.Create(call); err != nil {
		return fmt.Errorf("failed to create call: %w", err)
	}

	return nil
}

func (s *CallService) UpdateCallStatus(callSID, callStatus, callDuration string) error {
	duration, _ := strconv.Atoi(callDuration)

	// Determine if call was missed
	status := "answered"
	if callStatus == "no-answer" || callStatus == "busy" || duration == 0 {
		status = "missed"
	}

	call, err := s.callRepo.GetByTwilioSID(callSID)
	if err != nil {
		return fmt.Errorf("failed to find call: %w", err)
	}

	// Update call status
	if err := s.callRepo.UpdateStatus(call.ID, status, duration); err != nil {
		return fmt.Errorf("failed to update call status: %w", err)
	}

	// If call was missed, trigger SMS workflow
	if status == "missed" {
		if err := s.triggerMissedCallWorkflow(call); err != nil {
			fmt.Printf("Failed to trigger missed call workflow: %v\n", err)
		}
	}

	return nil
}

func (s *CallService) triggerMissedCallWorkflow(call *models.Call) error {
	// Create conversation
	conversation, err := s.conversationService.CreateConversation(
		call.ID,
		call.OrganizationID,
		call.LocationID,
		call.FromNumber,
	)
	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	// TODO: Send initial SMS via worker queue
	// For now, we'll mark it as pending for async processing
	fmt.Printf("Missed call workflow triggered for conversation: %s\n", conversation.ID)

	return nil
}

func (s *CallService) ProcessIncomingSMS(messageSID, fromNumber, toNumber, body string) error {
	// Find existing conversation by patient phone
	conversation, err := s.callRepo.FindActiveConversationByPhone(fromNumber)
	if err != nil || conversation == nil {
		return fmt.Errorf("no active conversation found for this number")
	}

	// Process the SMS through conversation service
	return s.conversationService.ProcessInboundSMS(
		conversation.ID,
		messageSID,
		fromNumber,
		toNumber,
		body,
	)
}
