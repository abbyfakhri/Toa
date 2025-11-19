package email

import (
	"context"

	"github.com/abbyfakhri/toa-api/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// single email client
// it accepts csv
// it store the email to database
// it iterates to send all the email
// for every 10 minutes or after it's all done it sends notif report to user

type EmailHandler interface {
	PostEmail(e echo.Context) error
}

type EmailUsecase interface {
	SendEmails(ctx context.Context, param models.PostEmailRequest) (batchId string, statusCode int, err error)
}

type EmailRepository interface {
	CreateBatch(ctx context.Context, tx *sqlx.Tx, param models.EmailBatch) error
	ReadBatch(ctx context.Context, db *sqlx.DB, batchId string) (models.EmailBatch, error)
	UpdateBatch(ctx context.Context, tx *sqlx.Tx, param models.EmailBatch) error
	DeleteBatch(ctx context.Context, tx *sqlx.Tx, batchId string) error

	CreateEmail(ctx context.Context, tx *sqlx.Tx, param models.Email) error
	ReadEmail(ctx context.Context, db *sqlx.DB, param models.GetEmailRequest) ([]models.Email, error)
	UpdateEmail(ctx context.Context, tx *sqlx.Tx, param models.Email) error
	DeleteEmail(ctx context.Context, tx *sqlx.Tx, emailId string) error
}
