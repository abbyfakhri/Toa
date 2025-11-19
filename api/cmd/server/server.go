package server

import (
	"net/http"

	"github.com/abbyfakhri/toa-api/internal/services"
	"github.com/abbyfakhri/toa-api/internal/services/email"
	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerConfig struct {
	Port        string
	Db          *sqlx.DB
	EmailClient email.EmailClient
}

type Server struct {
	cfg ServerConfig
}

func NewServer(cfg ServerConfig) Server {
	return Server{
		cfg: cfg,
	}
}

func (s *Server) Start() (*echo.Echo, error) {
	e := echo.New()
	e.Validator = &customValidator{
		validator: validator.New(),
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	services.LoadServices(e, s.cfg.Db, s.cfg.EmailClient)

	if err := e.Start(":" + s.cfg.Port); err != nil {
		return nil, err
	}

	return e, nil
}

// setup custom validator

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}
