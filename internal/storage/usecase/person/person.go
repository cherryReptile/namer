package person

import (
	"database/sql"
	"github.com/pkg/errors"
	"namer/internal/customErrors"
	"namer/internal/domain"
	"namer/internal/domain/external"
	personAPI "namer/internal/storage/repository/api/person"
	personPostgres "namer/internal/storage/repository/postgres/person"
	"namer/pkg/utils"
	"net/http"
)

type APIRepository interface {
	GetNameInfo(name string) (*external.ExternalResponse, error)
	GetAge(name string) (*external.AgifyResponse, error)
	GetGender(name string) (*external.GenderizeResponse, error)
	GetNation(name string) (*external.NationalizeResponse, error)
}

type PersonRepository interface {
	Create(req *domain.Person) error
	GetByID(id int) (*domain.Person, error)
	GetWithFilterAndPagination(filter, pagination string) ([]byte, error)
	Update(req *domain.Person) error
	Delete(id int) (*int64, error)
}

var (
	personNotFoundErr     = "person not found"
	emptyNameOrSurnameErr = "empty name or surname"
)

type Usecase struct {
	apiRepository    APIRepository
	personRepository PersonRepository
}

func NewUsecase(db *sql.DB) *Usecase {
	return &Usecase{
		apiRepository:    personAPI.NewRepository(),
		personRepository: personPostgres.NewRepository(db),
	}
}

func (u *Usecase) NewPerson(req *domain.Person) (*domain.Response, error) {
	utils.PrepareRequest(req)

	if req.Name == "" || req.Surname == "" {
		return nil, customErrors.New(
			emptyNameOrSurnameErr,
			errors.Wrap(errors.New(emptyNameOrSurnameErr), "NewPerson #1"),
			http.StatusBadRequest,
		)
	}

	info, err := u.apiRepository.GetNameInfo(req.Name)
	if err != nil {
		return nil, customErrors.New(
			http.StatusText(http.StatusInternalServerError),
			errors.Wrap(err, "NewPerson #2"),
			http.StatusInternalServerError,
		)
	}

	if info.Error != nil {
		return nil, customErrors.New(
			*info.Error,
			errors.Wrap(errors.New(*info.Error), "NewPerson #3"),
			info.StatusCode,
		)
	}

	req.Age, req.Gender = info.Agify.Age, info.Genderize.Gender

	if len(info.Nationalize.Country) != 0 {
		req.Nation = &info.Nationalize.Country[0].CountryId
	}

	if err = u.personRepository.Create(req); err != nil {
		return nil, customErrors.New(
			http.StatusText(http.StatusInternalServerError),
			errors.Wrap(err, "NewPerson #4"),
			http.StatusInternalServerError,
		)
	}

	return &domain.Response{
		Data:       req,
		StatusCode: http.StatusCreated,
	}, nil
}

func (u *Usecase) GetByID(id int) (*domain.Response, error) {
	person, err := u.personRepository.GetByID(id)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, customErrors.New(
				personNotFoundErr,
				errors.Wrap(err, "GetByID #1"),
				http.StatusBadRequest,
			)
		}

		return nil, customErrors.New(
			http.StatusText(http.StatusInternalServerError),
			errors.Wrap(err, "GetByID #2"),
			http.StatusInternalServerError,
		)
	}

	return &domain.Response{
		Data:       person,
		StatusCode: http.StatusOK,
	}, nil
}

func (u *Usecase) GetWithFilterAndPagination(req *domain.FilterWithPagination) (*domain.Response, error) {
	result, err := utils.GetFilterAndPagination(req, "pt")
	if err != nil {
		return nil, customErrors.New(
			err.Error(),
			errors.Wrap(err, "GetWithFilterAndPagination #1"),
			http.StatusBadRequest,
		)
	}

	b, err := u.personRepository.GetWithFilterAndPagination(result[0], result[1])
	if err != nil {
		return nil, customErrors.New(
			http.StatusText(http.StatusInternalServerError),
			errors.Wrap(err, "GetWithFilterAndPagination #2"),
			http.StatusInternalServerError,
		)
	}

	return &domain.Response{
		StatusCode: http.StatusOK,
		Data:       b,
	}, nil
}

func (u *Usecase) Update(req *domain.Person) (*domain.Response, error) {
	utils.PrepareRequest(req)

	if err := u.personRepository.Update(req); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, customErrors.New(
				personNotFoundErr,
				errors.Wrap(err, "Update #1"),
				http.StatusNotFound,
			)
		}

		return nil, customErrors.New(
			http.StatusText(http.StatusInternalServerError),
			errors.Wrap(err, "Update #2"),
			http.StatusInternalServerError,
		)
	}

	return &domain.Response{
		Data:       req,
		StatusCode: http.StatusOK,
	}, nil
}

func (u *Usecase) Delete(id int) (*domain.Response, error) {
	aff, err := u.personRepository.Delete(id)
	if err != nil {
		return nil, customErrors.New(
			http.StatusText(http.StatusInternalServerError),
			errors.Wrap(err, "Delete #1"),
			http.StatusInternalServerError,
		)
	}

	if *aff == 0 {
		return nil, customErrors.New(
			personNotFoundErr,
			errors.Wrap(errors.New(personNotFoundErr), "Delete #2"),
			http.StatusNotFound,
		)
	}

	return &domain.Response{
		Data:       "record deleted successfully",
		StatusCode: http.StatusOK,
	}, nil
}
