package email

import (
	"github.com/labstack/echo/v4"
)

func NewRoutes(e *echo.Echo, handler EmailHandler) {
	e.POST("/email", handler.PostEmail)
	e.POST("/email/csv", handler.PostEmailWithCsv)
}
