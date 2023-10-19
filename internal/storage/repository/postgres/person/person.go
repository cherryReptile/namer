package person

import (
	"database/sql"
	"namer/internal/domain"
)

type Repository interface {
	Create(req *domain.Person) (*domain.Person, error)
	GetByID(id int) (*domain.Person, error)
	GetWithFilterAndPagination(req *domain.FilterWithPagination) ([]domain.Person, error)
	Update(req *domain.Person) (*domain.Person, error)
	Delete(id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(req *domain.Person) (*domain.Person, error) {
	return nil, nil
}

func (r *repository) GetByID(id int) (*domain.Person, error) {
	return nil, nil
}

func (r *repository) GetWithFilterAndPagination(req *domain.FilterWithPagination) ([]domain.Person, error) {
	return nil, nil
}

func (r *repository) Update(req *domain.Person) (*domain.Person, error) {
	return nil, nil
}

func (r *repository) Delete(id int) error {
	return nil
}
