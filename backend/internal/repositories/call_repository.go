package repositories

import (
	"database/sql"
	"fmt"

	"github.com/recallflow/backend/internal/models"
)

type CallRepository struct {
	db *sql.DB
}

func NewCallRepository(db *sql.DB) *CallRepository {
	return &CallRepository{db: db}
}

func (r *CallRepository) Create(call *models.Call) error {
	query := `
		INSERT INTO calls (id, location_id, organization_id, twilio_call_sid, from_number, to_number, status, duration, call_time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(query,
		call.ID,
		call.LocationID,
		call.OrganizationID,
		call.TwilioCallSID,
		call.FromNumber,
		call.ToNumber,
		call.Status,
		call.Duration,
		call.CallTime,
		call.CreatedAt,
		call.UpdatedAt,
	)
	return err
}

func (r *CallRepository) GetByID(id string) (*models.Call, error) {
	query := `
		SELECT id, location_id, organization_id, twilio_call_sid, from_number, to_number, status, duration, call_time, created_at, updated_at
		FROM calls
		WHERE id = $1
	`
	
	call := &models.Call{}
	err := r.db.QueryRow(query, id).Scan(
		&call.ID,
		&call.LocationID,
		&call.OrganizationID,
		&call.TwilioCallSID,
		&call.FromNumber,
		&call.ToNumber,
		&call.Status,
		&call.Duration,
		&call.CallTime,
		&call.CreatedAt,
		&call.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("call not found")
	}
	
	return call, err
}

func (r *CallRepository) GetByTwilioSID(twilioSID string) (*models.Call, error) {
	query := `
		SELECT id, location_id, organization_id, twilio_call_sid, from_number, to_number, status, duration, call_time, created_at, updated_at
		FROM calls
		WHERE twilio_call_sid = $1
	`
	
	call := &models.Call{}
	err := r.db.QueryRow(query, twilioSID).Scan(
		&call.ID,
		&call.LocationID,
		&call.OrganizationID,
		&call.TwilioCallSID,
		&call.FromNumber,
		&call.ToNumber,
		&call.Status,
		&call.Duration,
		&call.CallTime,
		&call.CreatedAt,
		&call.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("call not found")
	}
	
	return call, err
}

func (r *CallRepository) UpdateStatus(callID, status string, duration int) error {
	query := `
		UPDATE calls
		SET status = $1, duration = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.Exec(query, status, duration, callID)
	return err
}

func (r *CallRepository) GetByOrganization(orgID string, limit int) ([]*models.Call, error) {
	query := `
		SELECT id, location_id, organization_id, twilio_call_sid, from_number, to_number, status, duration, call_time, created_at, updated_at
		FROM calls
		WHERE organization_id = $1
		ORDER BY call_time DESC
		LIMIT $2
	`
	
	rows, err := r.db.Query(query, orgID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var calls []*models.Call
	for rows.Next() {
		call := &models.Call{}
		err := rows.Scan(
			&call.ID,
			&call.LocationID,
			&call.OrganizationID,
			&call.TwilioCallSID,
			&call.FromNumber,
			&call.ToNumber,
			&call.Status,
			&call.Duration,
			&call.CallTime,
			&call.CreatedAt,
			&call.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		calls = append(calls, call)
	}
	
	return calls, nil
}

func (r *CallRepository) FindActiveConversationByPhone(phone string) (*models.Conversation, error) {
	query := `
		SELECT c.id, c.call_id, c.organization_id, c.location_id, c.patient_phone, c.status, c.intent, c.summary, c.is_resolved, c.created_at, c.updated_at
		FROM conversations c
		WHERE c.patient_phone = $1 AND c.status = 'active'
		ORDER BY c.created_at DESC
		LIMIT 1
	`
	
	conversation := &models.Conversation{}
	err := r.db.QueryRow(query, phone).Scan(
		&conversation.ID,
		&conversation.CallID,
		&conversation.OrganizationID,
		&conversation.LocationID,
		&conversation.PatientPhone,
		&conversation.Status,
		&conversation.Intent,
		&conversation.Summary,
		&conversation.IsResolved,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no active conversation found")
	}
	
	return conversation, err
}
