package email

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/abbyfakhri/toa-api/internal/models"
	"github.com/abbyfakhri/toa-api/internal/utils"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type Usecase struct {
	db          *sqlx.DB
	repository  EmailRepository
	emailClient EmailClient
}

// SendEmailsWithCsv implements EmailUsecase.
func (u Usecase) SendEmailsWithCsv(ctx context.Context, param models.PostEmailRequestCsv, r io.Reader) (batchId string, statusCode int, err error) {
	batchId = uuid.NewV4().String()

	// read file
	emails, err := utils.ReadCsv(r, param.TargetColumn)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	if len(emails) == 0 {
		return "", http.StatusBadRequest, fmt.Errorf("email destination cannot be empty")
	}

	go func() {
		for _, email := range emails {
			err := u.emailClient.SendMail(Email{
				To:       email,
				Subject:  param.Subject,
				Body:     param.Body,
				Template: param.Template,
			})

			if err != nil {
				log.Printf("fail to send email: %s", err.Error())
				// TODO: LOG THIS TO DATABASE
			}
		}
	}()
	return batchId, http.StatusAccepted, nil
}

// SendEmails implements EmailUsecase.
func (u Usecase) SendEmails(ctx context.Context, param models.PostEmailRequest) (batchId string, statusCode int, err error) {
	batchId = uuid.NewV4().String()

	if len(param.Destinations) == 0 {
		return "", http.StatusBadRequest, fmt.Errorf("email destination cannot be empty")
	}
	go func() {
		for _, address := range param.Destinations {
			err := u.emailClient.SendMail(Email{
				To:       address,
				Subject:  param.Subject,
				Body:     param.Body,
				Template: param.Template,
			})

			if err != nil {
				log.Printf("fail to send email: %s", err.Error())
				// TODO: LOG THIS TO DATABASE
			}
		}
	}()

	return batchId, http.StatusAccepted, nil
}

func NewUsecase(db *sqlx.DB, repository EmailRepository, emailClient EmailClient) EmailUsecase {
	return Usecase{
		db:          db,
		repository:  repository,
		emailClient: emailClient,
	}
}
