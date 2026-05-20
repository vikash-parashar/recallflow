package repositories

import (
	"database/sql"
	"fmt"

	"github.com/recallflow/backend/internal/models"
)

type LocationRepository struct {
	db *sql.DB
}

func NewLocationRepository(db *sql.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) Create(location *models.Location) error {
	query := `
		INSERT INTO locations (id, organization_id, name, phone_number, address, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query,
		location.ID,
		location.OrganizationID,
		location.Name,
		location.PhoneNumber,
		location.Address,
		location.IsActive,
		location.CreatedAt,
		location.UpdatedAt,
	)
	return err
}

func (r *LocationRepository) GetByID(id string) (*models.Location, error) {
	query := `
		SELECT id, organization_id, name, phone_number, address, is_active, created_at, updated_at
		FROM locations
		WHERE id = $1 AND deleted_at IS NULL
	`
	
	location := &models.Location{}
	err := r.db.QueryRow(query, id).Scan(
		&location.ID,
		&location.OrganizationID,
		&location.Name,
		&location.PhoneNumber,
		&location.Address,
		&location.IsActive,
		&location.CreatedAt,
		&location.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("location not found")
	}
	
	return location, err
}

func (r *LocationRepository) GetByPhoneNumber(phoneNumber string) (*models.Location, error) {
	query := `
		SELECT id, organization_id, name, phone_number, address, is_active, created_at, updated_at
		FROM locations
		WHERE phone_number = $1 AND deleted_at IS NULL AND is_active = true
	`
	
	location := &models.Location{}
	err := r.db.QueryRow(query, phoneNumber).Scan(
		&location.ID,
		&location.OrganizationID,
		&location.Name,
		&location.PhoneNumber,
		&location.Address,
		&location.IsActive,
		&location.CreatedAt,
		&location.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("location not found for phone number")
	}
	
	return location, err
}

func (r *LocationRepository) GetByOrganization(orgID string) ([]*models.Location, error) {
	query := `
		SELECT id, organization_id, name, phone_number, address, is_active, created_at, updated_at
		FROM locations
		WHERE organization_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var locations []*models.Location
	for rows.Next() {
		location := &models.Location{}
		err := rows.Scan(
			&location.ID,
			&location.OrganizationID,
			&location.Name,
			&location.PhoneNumber,
			&location.Address,
			&location.IsActive,
			&location.CreatedAt,
			&location.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	
	return locations, nil
}
