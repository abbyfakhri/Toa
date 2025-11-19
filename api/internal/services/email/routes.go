package email

import (
	"log"

	"github.com/labstack/echo/v4"
)

func NewRoutes(e *echo.Echo, handler EmailHandler) {
	e.POST("/email", handler.PostEmail)
	log.Print("loaded")
}
