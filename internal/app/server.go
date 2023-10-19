package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func NewServer(g *gin.Engine) *http.Server {
	s := &http.Server{
		Addr: fmt.Sprintf(
			"%s:%s",
			os.Getenv("SERVER_HOST"),
			os.Getenv("SERVER_PORT"),
		),
		Handler:      g,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	return s
}
