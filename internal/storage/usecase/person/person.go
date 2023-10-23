package person

import (
	"database/sql"
	"github.com/pkg/errors"
	"namer/internal/domain"
	personAPI "namer/internal/storage/repository/api/person"
	personPostgres "namer/internal/storage/repository/postgres/person"
	"net/http"
)

var personNotFoundErr = "person not found"

type Usecase interface {
	NewPerson(req *domain.Person) (*domain.Response, error)
	GetByID(id int) (*domain.Response, error)
	GetWithFilterAndPagination(req *domain.FilterWithPagination) (*domain.Response, error)
	Update(req *domain.Person) (*domain.Response, error)
	Delete(id int) (*domain.Response, error)
}

type usecase struct {
	apiRepository    personAPI.Repository
	personRepository personPostgres.Repository
}

func NewUsecase(db *sql.DB) Usecase {
	return &usecase{
		apiRepository:    personAPI.NewRepository(),
		personRepository: personPostgres.NewRepository(db),
	}
}

func (u *usecase) NewPerson(req *domain.Person) (*domain.Response, error) {
	prepareRequest(req)

	info, err := u.apiRepository.GetNameInfo(req.Name)
	if err != nil {
		return nil, errors.Wrap(err, "NewPerson #1")
	}

	if info.Error != nil {
		return &domain.Response{
			Error:      info.Error,
			StatusCode: info.StatusCode,
		}, nil
	}

	req.Age, req.Gender = info.Agify.Age, info.Genderize.Gender

	if len(info.Nationalize.Country) != 0 {
		req.Nation = &info.Nationalize.Country[0].CountryId
	}

	if err = u.personRepository.Create(req); err != nil {
		return nil, errors.Wrap(err, "NewPerson #2")
	}

	return &domain.Response{
		Data:       req,
		StatusCode: http.StatusCreated,
	}, nil
}

func (u *usecase) GetByID(id int) (*domain.Response, error) {
	person, err := u.personRepository.GetByID(id)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return &domain.Response{
				Error:      &personNotFoundErr,
				StatusCode: http.StatusNotFound,
			}, nil
		}

		return nil, errors.Wrap(err, "GetByID #1")
	}

	return &domain.Response{
		Data:       person,
		StatusCode: http.StatusOK,
	}, nil
}

func (u *usecase) GetWithFilterAndPagination(req *domain.FilterWithPagination) (*domain.Response, error) {
	return nil, nil
}

func (u *usecase) Update(req *domain.Person) (*domain.Response, error) {
	return nil, nil
}

func (u *usecase) Delete(id int) (*domain.Response, error) {
	aff, err := u.personRepository.Delete(id)
	if err != nil {
		return nil, errors.Wrap(err, "Delete #1")
	}

	if *aff == 0 {
		return &domain.Response{
			Error:      &personNotFoundErr,
			StatusCode: http.StatusNotFound,
		}, nil
	}

	return &domain.Response{
		Data:       "record deleted successfully",
		StatusCode: http.StatusOK,
	}, nil
}
