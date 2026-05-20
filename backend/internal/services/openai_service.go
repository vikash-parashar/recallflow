package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OpenAIService struct {
	apiKey string
	client *http.Client
}

func NewOpenAIService(apiKey string) *OpenAIService {
	return &OpenAIService{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

type IntentClassification struct {
	Intent          string  `json:"intent"`
	Confidence      float64 `json:"confidence"`
	IsEmergency     bool    `json:"is_emergency"`
	RequiresStaff   bool    `json:"requires_staff"`
	SuggestedAction string  `json:"suggested_action"`
}

func (s *OpenAIService) ClassifyIntent(ctx context.Context, message string) (*IntentClassification, error) {
	prompt := fmt.Sprintf(`You are analyzing a patient message to a healthcare clinic. 
Classify the intent and provide details.

Patient message: "%s"

Respond with a JSON object containing:
- intent: one of [appointment, billing, emergency, insurance, prescription, hours, general]
- confidence: float between 0 and 1
- is_emergency: boolean
- requires_staff: boolean
- suggested_action: brief suggestion for clinic staff

Respond ONLY with valid JSON, no other text.`, message)

	requestBody := map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a healthcare communication assistant."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.3,
		"max_tokens":  300,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error: %s", string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	var classification IntentClassification
	content := result.Choices[0].Message.Content
	
	// Try to extract JSON from the response
	if err := json.Unmarshal([]byte(content), &classification); err != nil {
		// Fallback classification
		return &IntentClassification{
			Intent:          "general",
			Confidence:      0.5,
			IsEmergency:     false,
			RequiresStaff:   true,
			SuggestedAction: "Review patient message",
		}, nil
	}

	return &classification, nil
}

func (s *OpenAIService) GenerateResponse(ctx context.Context, intent, patientMessage string) (string, error) {
	prompt := fmt.Sprintf(`Generate a brief, professional response to a patient who contacted a clinic.
Patient's intent: %s
Patient's message: "%s"

The response should be:
- Professional and friendly
- Healthcare-appropriate
- Brief (1-2 sentences)
- Never give medical advice
- Inform them staff will follow up if needed

Response:`, intent, patientMessage)

	requestBody := map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a professional healthcare clinic assistant."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
		"max_tokens":  150,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", strings.NewReader(string(jsonData)))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "Thank you for contacting us. Our team will get back to you shortly.", nil
	}

	return result.Choices[0].Message.Content, nil
}
