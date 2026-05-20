package handlers

import (
	"net/http"

	"github.com/recallflow/backend/internal/repositories"
	"github.com/recallflow/backend/internal/services"
)

type TwilioHandler struct {
	callService    *services.CallService
	twilioService  *services.TwilioService
	locationRepo   *repositories.LocationRepository
}

func NewTwilioHandler(
	callService *services.CallService,
	twilioService *services.TwilioService,
	locationRepo *repositories.LocationRepository,
) *TwilioHandler {
	return &TwilioHandler{
		callService:   callService,
		twilioService: twilioService,
		locationRepo:  locationRepo,
	}
}

// HandleVoiceWebhook processes incoming call webhooks from Twilio
func (h *TwilioHandler) HandleVoiceWebhook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		RespondError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	callSID := r.FormValue("CallSid")
	fromNumber := r.FormValue("From")
	toNumber := r.FormValue("To")
	callStatus := r.FormValue("CallStatus")

	// Find location by phone number
	location, err := h.locationRepo.GetByPhoneNumber(toNumber)
	if err != nil {
		RespondError(w, http.StatusNotFound, "Location not found for this number")
		return
	}

	// Process the call
	err = h.callService.ProcessIncomingCall(callSID, fromNumber, toNumber, callStatus, location.ID, location.OrganizationID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to process call")
		return
	}

	// Return TwiML response (tell Twilio to hang up)
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><Response><Hangup/></Response>`))
}

// HandleSMSWebhook processes incoming SMS messages from Twilio
func (h *TwilioHandler) HandleSMSWebhook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		RespondError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	messageSID := r.FormValue("MessageSid")
	fromNumber := r.FormValue("From")
	toNumber := r.FormValue("To")
	body := r.FormValue("Body")

	// Process the SMS
	err := h.callService.ProcessIncomingSMS(messageSID, fromNumber, toNumber, body)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to process SMS")
		return
	}

	// Return empty TwiML response
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><Response></Response>`))
}

// HandleStatusCallback processes status updates from Twilio
func (h *TwilioHandler) HandleStatusCallback(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		RespondError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	callSID := r.FormValue("CallSid")
	callStatus := r.FormValue("CallStatus")
	callDuration := r.FormValue("CallDuration")

	// Update call status in database
	// This is where we detect if a call was truly missed
	err := h.callService.UpdateCallStatus(callSID, callStatus, callDuration)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to update call status")
		return
	}

	w.WriteHeader(http.StatusOK)
}
