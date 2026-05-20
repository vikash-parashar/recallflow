package repositories

import (
	"database/sql"

	"github.com/recallflow/backend/internal/models"
)

type SMSRepository struct {
	db *sql.DB
}

func NewSMSRepository(db *sql.DB) *SMSRepository {
	return &SMSRepository{db: db}
}

func (r *SMSRepository) Create(message *models.SMSMessage) error {
	query := `
		INSERT INTO sms_messages (id, conversation_id, twilio_message_sid, direction, from_number, to_number, body, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(query,
		message.ID,
		message.ConversationID,
		message.TwilioMessageSID,
		message.Direction,
		message.FromNumber,
		message.ToNumber,
		message.Body,
		message.Status,
		message.CreatedAt,
	)
	return err
}

func (r *SMSRepository) GetByConversation(conversationID string) ([]*models.SMSMessage, error) {
	query := `
		SELECT id, conversation_id, twilio_message_sid, direction, from_number, to_number, body, status, created_at
		FROM sms_messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`
	
	rows, err := r.db.Query(query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var messages []*models.SMSMessage
	for rows.Next() {
		message := &models.SMSMessage{}
		err := rows.Scan(
			&message.ID,
			&message.ConversationID,
			&message.TwilioMessageSID,
			&message.Direction,
			&message.FromNumber,
			&message.ToNumber,
			&message.Body,
			&message.Status,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	
	return messages, nil
}

func (r *SMSRepository) GetByID(id string) (*models.SMSMessage, error) {
	query := `
		SELECT id, conversation_id, twilio_message_sid, direction, from_number, to_number, body, status, created_at
		FROM sms_messages
		WHERE id = $1
	`
	
	message := &models.SMSMessage{}
	err := r.db.QueryRow(query, id).Scan(
		&message.ID,
		&message.ConversationID,
		&message.TwilioMessageSID,
		&message.Direction,
		&message.FromNumber,
		&message.ToNumber,
		&message.Body,
		&message.Status,
		&message.CreatedAt,
	)
	
	return message, err
}
