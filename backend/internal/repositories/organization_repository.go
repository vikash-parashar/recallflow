package repositories

import (
	"database/sql"

	"github.com/recallflow/backend/internal/models"
)

type OrganizationRepository struct {
	db *sql.DB
}

func NewOrganizationRepository(db *sql.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Create(org *models.Organization) error {
	query := `
		INSERT INTO organizations (id, name, email, phone, address, plan_type, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(query,
		org.ID,
		org.Name,
		org.Email,
		org.Phone,
		org.Address,
		org.PlanType,
		org.IsActive,
		org.CreatedAt,
		org.UpdatedAt,
	)
	return err
}

func (r *OrganizationRepository) GetByID(id string) (*models.Organization, error) {
	query := `
		SELECT id, name, email, phone, address, plan_type, is_active, created_at, updated_at
		FROM organizations
		WHERE id = $1 AND deleted_at IS NULL
	`
	
	org := &models.Organization{}
	err := r.db.QueryRow(query, id).Scan(
		&org.ID,
		&org.Name,
		&org.Email,
		&org.Phone,
		&org.Address,
		&org.PlanType,
		&org.IsActive,
		&org.CreatedAt,
		&org.UpdatedAt,
	)
	
	return org, err
}
