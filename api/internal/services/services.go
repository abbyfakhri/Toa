package services

import (
	"github.com/abbyfakhri/toa-api/internal/services/email"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func LoadServices(e *echo.Echo, db *sqlx.DB, emailClient email.EmailClient) {
	email.Load(e, db, emailClient)
}

type ServiceLoader func(e *echo.Echo, db *sqlx.DB, emailClient email.EmailClient)
