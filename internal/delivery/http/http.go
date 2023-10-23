package http

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"namer/internal/domain"
	"namer/internal/storage/usecase/person"
	"net/http"
	"strconv"
)

type Handler struct {
	usecase person.Usecase
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		usecase: person.NewUsecase(db),
	}
}

func (h *Handler) NewPerson(c *gin.Context) {
	var (
		req domain.Person
		err error
	)

	defer func() {
		if err != nil {
			log.Errorln(err)

			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	}()

	if err = c.Bind(&req); err != nil {
		err = errors.Wrap(err, "NewPerson #1")

		return
	}

	res, err := h.usecase.NewPerson(&req)
	if err != nil {
		err = errors.Wrap(err, "NewPerson #2")

		return
	}

	c.JSON(res.StatusCode, res)
}

func (h *Handler) GetPerson(c *gin.Context) {
	var err error

	defer func() {
		if err != nil {
			log.Error(err)

			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	}()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err = errors.Wrap(err, "GetPerson #1")

		return
	}

	res, err := h.usecase.GetByID(id)
	if err != nil {
		err = errors.Wrap(err, "GetPerson #2")

		return
	}

	c.JSON(res.StatusCode, res)
}

func (h *Handler) GetPersons(c *gin.Context) {

}

func (h *Handler) UpdatePerson(c *gin.Context) {

}

func (h *Handler) DeletePerson(c *gin.Context) {
	var err error

	defer func() {
		if err != nil {
			log.Error(err)

			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	}()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err = errors.Wrap(err, "DeletePerson #1")

		return
	}

	res, err := h.usecase.Delete(id)
	if err != nil {
		err = errors.Wrap(err, "DeletePerson #2")

		return
	}

	c.JSON(res.StatusCode, res)
}
