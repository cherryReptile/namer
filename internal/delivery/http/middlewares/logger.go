package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"namer/internal/customErrors"
	"namer/internal/domain"
	"net/http"
)

func ErrLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				customErr, ok := err.(customErrors.Error)
				if !ok {
					log.Errorln(errors.Wrap(errors.New("failed to cast error"), "ErrLogger #1"))

					return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				}

				log.Errorln(errors.Wrap(customErr.Err, "ErrLogger #2"))

				return c.JSON(customErr.StatusCode, domain.Response{
					Error: &customErr.Message,
				})
			}

			return nil
		}
	}
}
