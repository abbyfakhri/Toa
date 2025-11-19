package email

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func Load(e *echo.Echo, db *sqlx.DB, emailClient EmailClient) {

	// init repository
	repository := NewRepository()

	// init usecase
	usecase := NewUsecase(db, repository, emailClient)

	// init handler
	handler := NewHandler(usecase)

	// init routes
	NewRoutes(e, handler)
}
