package http

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"namer/internal/storage/usecase /person"
	"net/http"
)

type Handler struct {
	db      *sql.DB
	usecase person.Usecase
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db:      db,
		usecase: person.NewUsecase(db),
	}
}

func (h *Handler) NewPerson(c *gin.Context) {
	c.String(http.StatusOK, "hello")
}

func (h *Handler) GetPerson(c *gin.Context) {

}

func (h *Handler) GetPersons(c *gin.Context) {

}

func (h *Handler) UpdatePerson(c *gin.Context) {

}

func (h *Handler) DeletePerson(c *gin.Context) {

}
