package person

import (
	"database/sql"
	"namer/internal/domain"
	"namer/internal/storage/repository/postgres/person"
)

type Usecase interface {
	NewPerson(req *domain.Person) (*domain.Person, error)
	GetByID(id int) (*domain.Person, error)
	GetWithFilterAndPagination(req *domain.FilterWithPagination) ([]domain.Person, error)
	Update(req *domain.Person) (*domain.Person, error)
	Delete(id int) error
}

type usecase struct {
	personRepository person.Repository
}

func NewUsecase(db *sql.DB) Usecase {
	return &usecase{
		personRepository: person.NewRepository(db),
	}
}

func (u *usecase) NewPerson(req *domain.Person) (*domain.Person, error) {
	return nil, nil
}

func (u *usecase) GetByID(id int) (*domain.Person, error) {
	return nil, nil
}

func (u *usecase) GetWithFilterAndPagination(req *domain.FilterWithPagination) ([]domain.Person, error) {
	return nil, nil
}

func (u *usecase) Update(req *domain.Person) (*domain.Person, error) {
	return nil, nil
}

func (u *usecase) Delete(id int) error {
	return nil
}
