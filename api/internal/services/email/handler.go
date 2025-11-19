package email

import (
	"net/http"

	"github.com/abbyfakhri/toa-api/internal/models"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase EmailUsecase
}

// PostEmail implements EmailHandler.
func (h Handler) PostEmail(e echo.Context) error {
	var request models.PostEmailRequest

	// bind request
	if err := e.Bind(&request); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// validate request
	if err := e.Validate(&request); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// call usecase
	batchId, statusCode, err := h.usecase.SendEmails(e.Request().Context(), request)
	if err != nil {
		return e.JSON(statusCode, map[string]string{
			"error": err.Error(),
		})
	}

	// return response
	return e.JSON(http.StatusAccepted, map[string]any{
		"data": map[string]string{
			"batchId": batchId,
		},
	})

}

func NewHandler(usecase EmailUsecase) EmailHandler {
	return Handler{
		usecase: usecase,
	}
}
