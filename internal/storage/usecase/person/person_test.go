package person

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"namer/internal/domain"
	"namer/internal/domain/external"
	"namer/internal/storage/usecase/person/mocks"
	"namer/pkg/utils"
	"net/http"
	"testing"
)

func newUsecase(apiRepo APIRepository, personRepo PersonRepository) *Usecase {
	return &Usecase{
		apiRepository:    apiRepo,
		personRepository: personRepo,
	}
}

func TestNewPerson(t *testing.T) {
	mockAPIRepository := new(mocks.APIRepository)
	mockPersonRepository := new(mocks.PersonRepository)

	usecase := newUsecase(mockAPIRepository, mockPersonRepository)

	person := domain.Person{
		Name:    "Helen",
		Surname: "Johnson",
	}

	extRes := external.ExternalResponse{
		Agify: &external.AgifyResponse{
			Count: utils.IntToPtr(30),
			Name:  utils.StringToPtr("Helen"),
			Age:   utils.IntToPtr(30),
		},
		Genderize: &external.GenderizeResponse{
			Count:       utils.IntToPtr(30),
			Name:        utils.StringToPtr("Helen"),
			Gender:      utils.StringToPtr("female"),
			Probability: utils.Float64ToPtr(0.1),
		},
		Nationalize: &external.NationalizeResponse{
			Count: utils.IntToPtr(30),
			Name:  utils.StringToPtr("Helen"),
			Country: []struct {
				CountryId   string  `json:"country_id"`
				Probability float64 `json:"probability"`
			}{
				{CountryId: "GB", Probability: 0.1},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		mockAPIRepository.On("GetNameInfo", person.Name).Return(&extRes, nil).Once()

		mockPersonRepository.On("Create", &person).Return(nil).Once()

		res, err := usecase.NewPerson(&person)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("error_empty_name_or_surname", func(t *testing.T) {
		res, err := usecase.NewPerson(&domain.Person{})
		if assert.Error(t, err) {
			assert.Equal(t, emptyNameOrSurnameErr, errors.Cause(err).Error())
		}
		assert.Nil(t, res)
	})

	t.Run("error_external_api", func(t *testing.T) {
		mockAPIRepository.On("GetNameInfo", person.Name).Return(
			nil,
			errors.New(http.StatusText(http.StatusInternalServerError)),
		).Once()

		res, err := usecase.NewPerson(&person)
		if assert.Error(t, err) {
			assert.Equal(t, http.StatusText(http.StatusInternalServerError), errors.Cause(err).Error())
		}

		assert.Nil(t, res)
	})

	t.Run("error_external_bad_request", func(t *testing.T) {
		mockAPIRepository.On("GetNameInfo", person.Name).Return(
			&external.ExternalResponse{
				Error: utils.StringToPtr("Missing 'name' parameter"),
			},
			nil,
		).Once()

		res, err := usecase.NewPerson(&person)
		if assert.Error(t, err) {
			assert.Equal(t, "Missing 'name' parameter", errors.Cause(err).Error())
		}

		assert.Nil(t, res)
	})

	t.Run("error_postgres_create", func(t *testing.T) {
		mockAPIRepository.On("GetNameInfo", person.Name).Return(&extRes, nil).Once()

		mockPersonRepository.On("Create", &person).Return(
			errors.New("pg_error"),
		).Once()

		res, err := usecase.NewPerson(&person)
		if assert.Error(t, err) {
			assert.Equal(t, "pg_error", errors.Cause(err).Error())
		}

		assert.Nil(t, res)
	})
}

func TestGetByID(t *testing.T) {
	mockAPIRepository := new(mocks.APIRepository)
	mockPersonRepository := new(mocks.PersonRepository)

	usecase := newUsecase(mockAPIRepository, mockPersonRepository)

	t.Run("success", func(t *testing.T) {
		mockPersonRepository.On("GetByID", 1).Return(
			&domain.Person{},
			nil,
		).Once()

		res, err := usecase.GetByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("error_not_found", func(t *testing.T) {
		mockPersonRepository.On("GetByID", 1).Return(
			nil,
			sql.ErrNoRows,
		).Once()

		res, err := usecase.GetByID(1)
		if assert.Error(t, err) {
			assert.Equal(t, sql.ErrNoRows.Error(), errors.Cause(err).Error())
		}

		assert.Nil(t, res)
	})

	t.Run("error_postgres_create", func(t *testing.T) {
		mockPersonRepository.On("GetByID", 1).Return(
			nil,
			errors.New("pg_error"),
		).Once()

		res, err := usecase.GetByID(1)
		if assert.Error(t, err) {
			assert.Equal(t, "pg_error", errors.Cause(err).Error())
		}

		assert.Nil(t, res)
	})
}

func TestGetWithFilterAndPagination(t *testing.T) {
	mockAPIRepository := new(mocks.APIRepository)
	mockPersonRepository := new(mocks.PersonRepository)

	usecase := newUsecase(mockAPIRepository, mockPersonRepository)

	req := domain.FilterWithPagination{
		Filter: []domain.Filter{
			{
				Field: "name",
				Value: "en",
			},
		},
		Pagination: &domain.Pagination{
			Page:  1,
			Limit: 5,
		},
	}

	t.Run("success", func(t *testing.T) {
		mockPersonRepository.On("GetWithFilterAndPagination", mock.Anything, mock.Anything).Return(
			[]byte("test"),
			nil,
		).Once()

		res, err := usecase.GetWithFilterAndPagination(&req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("error_bad_request", func(t *testing.T) {
		mockReq := domain.FilterWithPagination{
			Filter: []domain.Filter{
				{
					Field: "test test",
					Value: "test",
				},
			},
		}

		res, err := usecase.GetWithFilterAndPagination(&mockReq)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("error_postgres_get", func(t *testing.T) {
		mockPersonRepository.On("GetWithFilterAndPagination", mock.Anything, mock.Anything).Return(
			nil,
			errors.New("pg_error"),
		).Once()

		res, err := usecase.GetWithFilterAndPagination(&req)
		if assert.Error(t, err) {
			assert.Equal(t, "pg_error", errors.Cause(err).Error())
		}

		assert.Nil(t, res)
	})
}

func TestUpdate(t *testing.T) {
	mockAPIRepository := new(mocks.APIRepository)
	mockPersonRepository := new(mocks.PersonRepository)

	usecase := newUsecase(mockAPIRepository, mockPersonRepository)

	person := domain.Person{
		ID:      1,
		Name:    "Helen",
		Surname: "Johnson",
	}

	t.Run("success", func(t *testing.T) {
		mockPersonRepository.On("Update", &person).Return(
			nil,
		).Once()

		res, err := usecase.Update(&person)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("error_not_found", func(t *testing.T) {
		mockPersonRepository.On("Update", &person).Return(
			sql.ErrNoRows,
		).Once()

		res, err := usecase.Update(&person)
		if assert.Error(t, err) {
			assert.Equal(t, sql.ErrNoRows, errors.Cause(err))
		}

		assert.Nil(t, res)
	})

	t.Run("error_postgres_update", func(t *testing.T) {
		mockPersonRepository.On("Update", &person).Return(
			errors.New("pg_error"),
		).Once()

		res, err := usecase.Update(&person)
		if assert.Error(t, err) {
			assert.Equal(t, "pg_error", errors.Cause(err).Error())
		}

		assert.Nil(t, res)
	})
}

func TestDelete(t *testing.T) {
	mockAPIRepository := new(mocks.APIRepository)
	mockPersonRepository := new(mocks.PersonRepository)

	usecase := newUsecase(mockAPIRepository, mockPersonRepository)

	t.Run("success", func(t *testing.T) {
		mockPersonRepository.On("Delete", 1).Return(
			utils.Int64ToPtr(1),
			nil,
		).Once()

		res, err := usecase.Delete(1)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("error_postgres_delete", func(t *testing.T) {
		mockPersonRepository.On("Delete", 1).Return(
			nil,
			errors.New("pg_error"),
		).Once()

		res, err := usecase.Delete(1)
		if assert.Error(t, err) {
			assert.Equal(t, "pg_error", errors.Cause(err).Error())
		}

		assert.Nil(t, res)
	})

	t.Run("error_not_found", func(t *testing.T) {
		mockPersonRepository.On("Delete", 1).Return(
			utils.Int64ToPtr(0),
			nil,
		).Once()

		res, err := usecase.Delete(1)
		if assert.Error(t, err) {
			assert.Equal(t, personNotFoundErr, errors.Cause(err).Error())
		}

		assert.Nil(t, res)
	})
}
