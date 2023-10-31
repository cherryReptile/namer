package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"time"
)

func newServer(e *echo.Echo) *http.Server {
	s := &http.Server{
		Addr: fmt.Sprintf(
			"%s:%s",
			os.Getenv("SERVER_HOST"),
			os.Getenv("SERVER_PORT"),
		),
		Handler:      e,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	return s
}
