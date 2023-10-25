package http

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"namer/internal/domain"
	"namer/internal/storage/usecase/person"
	"namer/pkg/utils"
	"net/http"
	"strconv"
)

type Usecase interface {
	NewPerson(req *domain.Person) (*domain.Response, error)
	GetByID(id int) (*domain.Response, error)
	GetWithFilterAndPagination(req *domain.FilterWithPagination) (*domain.Response, error)
	Update(req *domain.Person) (*domain.Response, error)
	Delete(id int) (*domain.Response, error)
}

var (
	invalidParameter   = "invalid uri parameter"
	invalidRequestBody = "invalid request body"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		usecase: person.NewUsecase(db),
	}
}

func (h *Handler) NewPerson(c *gin.Context) {
	var req domain.Person

	if err := c.BindJSON(&req); err != nil || req.Name == "" || req.Surname == "" {
		log.Errorln(errors.Wrap(err, "NewPerson #1"))

		c.AbortWithStatusJSON(http.StatusBadRequest, domain.Response{
			Error: &invalidRequestBody,
		})

		return
	}

	res, err := h.usecase.NewPerson(&req)
	if err != nil {
		log.Errorln(errors.Wrap(err, "NewPerson #2"))

		c.AbortWithStatusJSON(http.StatusInternalServerError, domain.Response{
			Error: utils.StringToPtr(http.StatusText(http.StatusInternalServerError)),
		})

		return
	}

	c.JSON(res.StatusCode, res)
}

func (h *Handler) GetPerson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorln(errors.Wrap(err, "GetPerson #1"))

		c.AbortWithStatusJSON(http.StatusBadRequest, domain.Response{
			Error: &invalidParameter,
		})

		return
	}

	res, err := h.usecase.GetByID(id)
	if err != nil {
		log.Errorln(errors.Wrap(err, "GetPerson #2"))

		c.AbortWithStatusJSON(http.StatusInternalServerError, domain.Response{
			Error: utils.StringToPtr(http.StatusText(http.StatusInternalServerError)),
		})

		return
	}

	c.JSON(res.StatusCode, res)
}

func (h *Handler) GetPersons(c *gin.Context) {
	var req domain.FilterWithPagination

	if err := c.BindJSON(&req); err != nil {
		log.Errorln(errors.Wrap(err, "GetPersons #1"))

		c.AbortWithStatusJSON(http.StatusBadRequest, domain.Response{
			Error: &invalidRequestBody,
		})

		return
	}

	res, err := h.usecase.GetWithFilterAndPagination(&req)
	if err != nil {
		log.Errorln(errors.Wrap(err, "GetPersons #2"))

		c.AbortWithStatusJSON(http.StatusInternalServerError, domain.Response{
			Error: utils.StringToPtr(http.StatusText(http.StatusInternalServerError)),
		})

		return
	}

	c.JSON(res.StatusCode, res)
}

func (h *Handler) UpdatePerson(c *gin.Context) {
	var (
		req domain.Person
		err error
	)

	if req.ID, err = strconv.Atoi(c.Param("id")); err != nil {
		log.Errorln(errors.Wrap(err, "UpdatePerson #1"))

		c.AbortWithStatusJSON(http.StatusBadRequest, domain.Response{
			Error: &invalidParameter,
		})

		return
	}

	if err = c.BindJSON(&req); err != nil {
		log.Errorln(errors.Wrap(err, "UpdatePerson #2"))

		c.AbortWithStatusJSON(http.StatusBadRequest, domain.Response{
			Error: &invalidRequestBody,
		})

		return
	}

	res, err := h.usecase.Update(&req)
	if err != nil {
		log.Errorln(errors.Wrap(err, "UpdatePerson #3"))

		c.AbortWithStatusJSON(http.StatusInternalServerError, domain.Response{
			Error: utils.StringToPtr(http.StatusText(http.StatusInternalServerError)),
		})

		return
	}

	c.JSON(res.StatusCode, res)
}

func (h *Handler) DeletePerson(c *gin.Context) {
	var err error

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorln(errors.Wrap(err, "DeletePerson #1"))

		c.AbortWithStatusJSON(http.StatusBadRequest, domain.Response{
			Error: &invalidParameter,
		})

		return
	}

	res, err := h.usecase.Delete(id)
	if err != nil {
		log.Errorln(errors.Wrap(err, "DeletePerson #2"))

		c.AbortWithStatusJSON(http.StatusInternalServerError, domain.Response{
			Error: utils.StringToPtr(http.StatusText(http.StatusInternalServerError)),
		})

		return
	}

	c.JSON(res.StatusCode, res)
}
