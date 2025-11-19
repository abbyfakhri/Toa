package email

import (
	"context"

	"github.com/abbyfakhri/toa-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
}

// CreateBatch implements EmailRepository.
func (r Repository) CreateBatch(ctx context.Context, tx *sqlx.Tx, param models.EmailBatch) error {
	panic("unimplemented")
}

// CreateEmail implements EmailRepository.
func (r Repository) CreateEmail(ctx context.Context, tx *sqlx.Tx, param models.Email) error {
	panic("unimplemented")
}

// DeleteBatch implements EmailRepository.
func (r Repository) DeleteBatch(ctx context.Context, tx *sqlx.Tx, batchId string) error {
	panic("unimplemented")
}

// DeleteEmail implements EmailRepository.
func (r Repository) DeleteEmail(ctx context.Context, tx *sqlx.Tx, emailId string) error {
	panic("unimplemented")
}

// ReadBatch implements EmailRepository.
func (r Repository) ReadBatch(ctx context.Context, db *sqlx.DB, batchId string) (models.EmailBatch, error) {
	panic("unimplemented")
}

// ReadEmail implements EmailRepository.
func (r Repository) ReadEmail(ctx context.Context, db *sqlx.DB, param models.GetEmailRequest) ([]models.Email, error) {
	panic("unimplemented")
}

// UpdateBatch implements EmailRepository.
func (r Repository) UpdateBatch(ctx context.Context, tx *sqlx.Tx, param models.EmailBatch) error {
	panic("unimplemented")
}

// UpdateEmail implements EmailRepository.
func (r Repository) UpdateEmail(ctx context.Context, tx *sqlx.Tx, param models.Email) error {
	panic("unimplemented")
}

func NewRepository() EmailRepository {
	return Repository{}
}
