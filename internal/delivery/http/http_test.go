package http

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"namer/internal/customErrors"
	"namer/internal/delivery/http/middlewares"
	"namer/internal/delivery/http/mocks"
	"namer/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewPerson(t *testing.T) {
	person := domain.Person{
		ID:      1,
		Name:    "Helen",
		Surname: "Johnson",
	}

	res := domain.Response{
		Data:       &person,
		StatusCode: http.StatusCreated,
	}

	t.Run("success", func(t *testing.T) {
		mockUsecase := new(mocks.Usecase)

		h := &Handler{
			usecase: mockUsecase,
		}

		e := echo.New()

		e.POST("/api/person", h.NewPerson)

		rec := httptest.NewRecorder()

		mockUsecase.On("NewPerson", &person).Return(&res, nil).Once()

		b, err := json.Marshal(person)
		assert.NoError(t, err)
		assert.NotNil(t, b)

		req := httptest.NewRequest(http.MethodPost, "/api/person", bytes.NewBuffer(b))
		req.Header.Add("content-type", "application/json")

		e.ServeHTTP(rec, req)

		httpResponse := make(map[string]*domain.Person)

		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &httpResponse))

		responsePerson, ok := httpResponse["data"]
		assert.True(t, ok)

		if assert.NotNil(t, responsePerson) {
			assert.Equal(t, person, *responsePerson)
		}

		mockUsecase.AssertExpectations(t)
	})

	t.Run("error_bind", func(t *testing.T) {
		mockUsecase := new(mocks.Usecase)

		h := &Handler{
			usecase: mockUsecase,
		}

		e := echo.New()

		e.POST("/api/person", h.NewPerson, middlewares.ErrLogger())

		rec := httptest.NewRecorder()

		badRequest := map[string]any{"name": 1}

		b, err := json.Marshal(badRequest)
		assert.NoError(t, err)
		assert.NotNil(t, b)

		req := httptest.NewRequest(http.MethodPost, "/api/person", bytes.NewBuffer(b))
		req.Header.Add("content-type", "application/json")

		e.ServeHTTP(rec, req)

		var httpResponse domain.Response

		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &httpResponse))
		assert.Equal(t, rec.Code, http.StatusBadRequest)
		assert.Equal(t, invalidRequestBodyErr, *httpResponse.Error)
		assert.Nil(t, httpResponse.Data)
	})

	t.Run("error_bad_request", func(t *testing.T) {
		mockUsecase := new(mocks.Usecase)

		h := &Handler{
			usecase: mockUsecase,
		}

		e := echo.New()

		e.POST("/api/person", h.NewPerson, middlewares.ErrLogger())

		rec := httptest.NewRecorder()

		badPerson := domain.Person{}

		mockUsecase.On("NewPerson", &badPerson).Return(
			nil,
			customErrors.New(
				"empty name or surname",
				errors.New("empty name or surname"),
				http.StatusBadRequest,
			),
		).Once()

		b, err := json.Marshal(badPerson)
		assert.NoError(t, err)
		assert.NotNil(t, b)

		req := httptest.NewRequest(http.MethodPost, "/api/person", bytes.NewBuffer(b))
		req.Header.Add("content-type", "application/json")

		e.ServeHTTP(rec, req)

		var httpResponse domain.Response

		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &httpResponse))
		assert.Equal(t, rec.Code, http.StatusBadRequest)
		assert.Equal(t, "empty name or surname", *httpResponse.Error)
		assert.Nil(t, httpResponse.Data)

		mockUsecase.AssertExpectations(t)
	})
}

func TestGetPerson(t *testing.T) {
	//
}

func TestGetPersons(t *testing.T) {
	//
}

func TestUpdate(t *testing.T) {
	//
}

func TestDelete(t *testing.T) {
	//
}
