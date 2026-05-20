package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/recallflow/backend/internal/middleware"
)

type DashboardHandler struct {
	db *sql.DB
}

func NewDashboardHandler(db *sql.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

type DashboardStats struct {
	TotalCalls         int     `json:"total_calls"`
	MissedCalls        int     `json:"missed_calls"`
	RecoveredLeads     int     `json:"recovered_leads"`
	ActiveConversations int    `json:"active_conversations"`
	ResponseRate       float64 `json:"response_rate"`
	EstimatedRevenue   float64 `json:"estimated_revenue"`
}

func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get stats for the last 30 days
	startDate := time.Now().AddDate(0, 0, -30)

	var stats DashboardStats
	
	// Total calls
	h.db.QueryRow(`
		SELECT COUNT(*) FROM calls 
		WHERE organization_id = $1 AND call_time >= $2
	`, claims.OrganizationID, startDate).Scan(&stats.TotalCalls)

	// Missed calls
	h.db.QueryRow(`
		SELECT COUNT(*) FROM calls 
		WHERE organization_id = $1 AND status = 'missed' AND call_time >= $2
	`, claims.OrganizationID, startDate).Scan(&stats.MissedCalls)

	// Active conversations
	h.db.QueryRow(`
		SELECT COUNT(*) FROM conversations 
		WHERE organization_id = $1 AND status = 'active'
	`, claims.OrganizationID).Scan(&stats.ActiveConversations)

	// Recovered leads (conversations with responses)
	h.db.QueryRow(`
		SELECT COUNT(DISTINCT conversation_id) FROM sms_messages 
		WHERE conversation_id IN (
			SELECT id FROM conversations WHERE organization_id = $1
		) AND direction = 'inbound' AND created_at >= $2
	`, claims.OrganizationID, startDate).Scan(&stats.RecoveredLeads)

	// Calculate response rate
	if stats.MissedCalls > 0 {
		stats.ResponseRate = float64(stats.RecoveredLeads) / float64(stats.MissedCalls) * 100
	}

	// Estimated revenue (assuming $200 per recovered lead)
	stats.EstimatedRevenue = float64(stats.RecoveredLeads) * 200

	RespondJSON(w, http.StatusOK, stats)
}

func (h *DashboardHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get daily analytics for the last 30 days
	rows, err := h.db.Query(`
		SELECT 
			DATE(call_time) as date,
			COUNT(*) as total_calls,
			COUNT(CASE WHEN status = 'missed' THEN 1 END) as missed_calls
		FROM calls
		WHERE organization_id = $1 AND call_time >= NOW() - INTERVAL '30 days'
		GROUP BY DATE(call_time)
		ORDER BY date DESC
	`, claims.OrganizationID)
	
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to fetch analytics")
		return
	}
	defer rows.Close()

	type DailyAnalytics struct {
		Date        string `json:"date"`
		TotalCalls  int    `json:"total_calls"`
		MissedCalls int    `json:"missed_calls"`
	}

	var analytics []DailyAnalytics
	for rows.Next() {
		var a DailyAnalytics
		rows.Scan(&a.Date, &a.TotalCalls, &a.MissedCalls)
		analytics = append(analytics, a)
	}

	RespondJSON(w, http.StatusOK, analytics)
}
