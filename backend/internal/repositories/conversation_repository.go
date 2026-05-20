package repositories

import (
	"database/sql"
	"fmt"

	"github.com/recallflow/backend/internal/models"
)

type ConversationRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) Create(conversation *models.Conversation) error {
	query := `
		INSERT INTO conversations (id, call_id, organization_id, location_id, patient_phone, status, intent, summary, is_resolved, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(query,
		conversation.ID,
		conversation.CallID,
		conversation.OrganizationID,
		conversation.LocationID,
		conversation.PatientPhone,
		conversation.Status,
		conversation.Intent,
		conversation.Summary,
		conversation.IsResolved,
		conversation.CreatedAt,
		conversation.UpdatedAt,
	)
	return err
}

func (r *ConversationRepository) GetByID(id, orgID string) (*models.Conversation, error) {
	query := `
		SELECT id, call_id, organization_id, location_id, patient_phone, status, intent, summary, is_resolved, created_at, updated_at
		FROM conversations
		WHERE id = $1 AND organization_id = $2
	`
	
	conversation := &models.Conversation{}
	err := r.db.QueryRow(query, id, orgID).Scan(
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
		return nil, fmt.Errorf("conversation not found")
	}
	
	return conversation, err
}

func (r *ConversationRepository) GetByOrganization(orgID string) ([]*models.Conversation, error) {
	query := `
		SELECT id, call_id, organization_id, location_id, patient_phone, status, intent, summary, is_resolved, created_at, updated_at
		FROM conversations
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`
	
	rows, err := r.db.Query(query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var conversations []*models.Conversation
	for rows.Next() {
		conversation := &models.Conversation{}
		err := rows.Scan(
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
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conversation)
	}
	
	return conversations, nil
}

func (r *ConversationRepository) UpdateIntent(id, intent string) error {
	query := `
		UPDATE conversations
		SET intent = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.Exec(query, intent, id)
	return err
}

func (r *ConversationRepository) MarkResolved(id, orgID string) error {
	query := `
		UPDATE conversations
		SET is_resolved = true, status = 'resolved', updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
	`
	_, err := r.db.Exec(query, id, orgID)
	return err
}
