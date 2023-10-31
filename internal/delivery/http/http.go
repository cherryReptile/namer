package http

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"namer/internal/customErrors"
	"namer/internal/domain"
	"namer/internal/storage/usecase/person"
	"net/http"
	"strconv"
)

//go:generate mockery --name Usecase
type Usecase interface {
	NewPerson(req *domain.Person) (*domain.Response, error)
	GetByID(id int) (*domain.Response, error)
	GetWithFilterAndPagination(req *domain.FilterWithPagination) (*domain.Response, error)
	Update(req *domain.Person) (*domain.Response, error)
	Delete(id int) (*domain.Response, error)
}

var (
	invalidParameterErr   = "invalid uri parameter"
	invalidRequestBodyErr = "invalid request body"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		usecase: person.NewUsecase(db),
	}
}

func (h *Handler) NewPerson(c echo.Context) error {
	var req domain.Person

	if err := c.Bind(&req); err != nil {
		return customErrors.New(
			invalidRequestBodyErr,
			errors.Wrap(err, "NewPerson #1"),
			http.StatusBadRequest,
		)
	}

	res, err := h.usecase.NewPerson(&req)
	if err != nil {
		return customErrors.Wrap(err, "NewPerson #2")
	}

	return c.JSON(res.StatusCode, res)
}

func (h *Handler) GetPerson(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return customErrors.New(
			invalidParameterErr,
			errors.Wrap(err, "GetPerson #1"),
			http.StatusBadRequest,
		)
	}

	res, err := h.usecase.GetByID(id)
	if err != nil {
		return customErrors.Wrap(err, "GetPerson #2")
	}

	return c.JSON(res.StatusCode, res)
}

func (h *Handler) GetPersons(c echo.Context) error {
	var req domain.FilterWithPagination

	if err := c.Bind(&req); err != nil {
		return customErrors.New(
			invalidRequestBodyErr,
			errors.Wrap(err, "GetPersons #1"),
			http.StatusBadRequest,
		)
	}

	res, err := h.usecase.GetWithFilterAndPagination(&req)
	if err != nil {
		return customErrors.Wrap(err, "GetPersons #2")
	}

	resB, ok := res.Data.([]byte)
	if !ok {
		return c.JSON(res.StatusCode, res)
	}

	return c.JSONBlob(res.StatusCode, resB)
}

func (h *Handler) UpdatePerson(c echo.Context) error {
	var (
		req domain.Person
		err error
	)

	if req.ID, err = strconv.Atoi(c.Param("id")); err != nil {
		return customErrors.New(
			invalidParameterErr,
			errors.Wrap(err, "UpdatePerson #1"),
			http.StatusBadRequest,
		)
	}

	if err = c.Bind(&req); err != nil {
		return customErrors.New(
			invalidRequestBodyErr,
			errors.Wrap(err, "UpdatePerson #2"),
			http.StatusBadRequest,
		)
	}

	res, err := h.usecase.Update(&req)
	if err != nil {
		return customErrors.Wrap(err, "UpdatePerson #3")
	}

	return c.JSON(res.StatusCode, res)
}

func (h *Handler) DeletePerson(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return customErrors.New(
			invalidParameterErr,
			errors.Wrap(err, "DeletePerson #1"),
			http.StatusBadRequest,
		)
	}

	res, err := h.usecase.Delete(id)
	if err != nil {
		return customErrors.Wrap(err, "DeletePerson #2")
	}

	return c.JSON(res.StatusCode, res)
}
