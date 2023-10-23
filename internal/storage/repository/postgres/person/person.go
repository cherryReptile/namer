package person

import (
	"database/sql"
	"github.com/pkg/errors"
	"namer/internal/domain"
)

type Repository interface {
	Create(req *domain.Person) error
	GetByID(id int) (*domain.Person, error)
	GetWithFilterAndPagination(req *domain.FilterWithPagination) ([]domain.Person, error)
	Update(req *domain.Person) (*int64, error)
	Delete(id int) (*int64, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(req *domain.Person) error {
	query := `
		insert into persons.persons_table
		(
		 name,
		 surname,
		 patronymic,
		 age,
		 gender,
		 nation
		)
		values ($1, $2, $3, $4, $5, $6)
		returning id, created_at
	`

	if err := r.db.QueryRow(
		query,
		req.Name,
		req.Surname,
		req.Patronymic,
		req.Age,
		req.Gender,
		req.Nation,
	).Scan(
		&req.ID,
		&req.CreatedAt,
	); err != nil {
		return errors.Wrap(err, "Create #1")
	}

	return nil
}

func (r *repository) GetByID(id int) (*domain.Person, error) {
	query := `
		select
		    p.id, 
		    p.name, 
		    p.surname, 
		    p.patronymic, 
		    p.age, 
		    p.gender, 
		    p.nation, 
		    p.created_at, 
		    p.updated_at
		from persons.persons_table p 
		where id = $1
	`

	var person domain.Person

	if err := r.db.QueryRow(query, id).Scan(
		&person.ID,
		&person.Name,
		&person.Surname,
		&person.Patronymic,
		&person.Age,
		&person.Gender,
		&person.Nation,
		&person.CreatedAt,
		&person.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "GetByID #1")
	}

	return &person, nil
}

func (r *repository) GetWithFilterAndPagination(req *domain.FilterWithPagination) ([]domain.Person, error) {
	return nil, nil
}

func (r *repository) Update(req *domain.Person) (*int64, error) {
	return nil, nil
}

func (r *repository) Delete(id int) (*int64, error) {
	query := `
		delete
		from persons.persons_table
		where id = $1
	`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return nil, errors.Wrap(err, "Delete #1")
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return nil, errors.Wrap(err, "Delete #2")
	}

	return &aff, nil
}
