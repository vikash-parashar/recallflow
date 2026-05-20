package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/recallflow/backend/internal/middleware"
	"github.com/recallflow/backend/internal/services"
)

type ConversationHandler struct {
	conversationService *services.ConversationService
}

func NewConversationHandler(conversationService *services.ConversationService) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationService,
	}
}

func (h *ConversationHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	conversations, err := h.conversationService.GetConversationsByOrg(claims.OrganizationID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to fetch conversations")
		return
	}

	RespondJSON(w, http.StatusOK, conversations)
}

func (h *ConversationHandler) GetConversation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationID := vars["id"]

	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	conversation, err := h.conversationService.GetConversation(conversationID, claims.OrganizationID)
	if err != nil {
		RespondError(w, http.StatusNotFound, "Conversation not found")
		return
	}

	RespondJSON(w, http.StatusOK, conversation)
}

func (h *ConversationHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationID := vars["id"]

	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	messages, err := h.conversationService.GetMessages(conversationID, claims.OrganizationID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to fetch messages")
		return
	}

	RespondJSON(w, http.StatusOK, messages)
}

func (h *ConversationHandler) ResolveConversation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationID := vars["id"]

	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := h.conversationService.ResolveConversation(conversationID, claims.OrganizationID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to resolve conversation")
		return
	}

	RespondSuccess(w, "Conversation resolved successfully", nil)
}
