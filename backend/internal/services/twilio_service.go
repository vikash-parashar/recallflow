package services

import (
	"fmt"

	"github.com/recallflow/backend/internal/config"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioService struct {
	client      *twilio.RestClient
	phoneNumber string
}

func NewTwilioService(cfg *config.Config) *TwilioService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.TwilioAccountSID,
		Password: cfg.TwilioAuthToken,
	})

	return &TwilioService{
		client:      client,
		phoneNumber: cfg.TwilioPhoneNumber,
	}
}

func (s *TwilioService) SendSMS(to, body string) (string, error) {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(s.phoneNumber)
	params.SetBody(body)

	resp, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return "", fmt.Errorf("failed to send SMS: %w", err)
	}

	return *resp.Sid, nil
}

func (s *TwilioService) SendMissedCallSMS(to, clinicName string) (string, error) {
	message := fmt.Sprintf(
		"Sorry we missed your call to %s. How can we help you today?",
		clinicName,
	)
	return s.SendSMS(to, message)
}
