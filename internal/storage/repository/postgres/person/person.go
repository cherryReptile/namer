package person

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"namer/internal/domain"
	"time"
)

type PersonRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PersonRepository {
	return &PersonRepository{
		db: db,
	}
}

func (r *PersonRepository) Create(req *domain.Person) error {
	query := `
		insert into persons.persons_table
		(
		 name,
		 surname,
		 patronymic,
		 age,
		 gender,
		 nation,
		 created_at
		)
		values ($1, $2, $3, $4, $5, $6, $7)
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
		time.Now(),
	).Scan(
		&req.ID,
		&req.CreatedAt,
	); err != nil {
		return errors.Wrap(err, "Create #1")
	}

	return nil
}

func (r *PersonRepository) GetByID(id int) (*domain.Person, error) {
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

func (r *PersonRepository) GetWithFilterAndPagination(filter, pagination string) ([]byte, error) {
	query := fmt.Sprintf(`
		select jsonb_build_object(
					   'data',
					   jsonb_agg(
							   jsonb_build_object(
									   'id', p.id,
									   'name', p.name,
									   'surname', p.surname,
									   'patronymic', p.patronymic,
									   'age', p.age,
									   'gender', p.gender,
									   'nation', p.nation,
									   'created_at', to_char(p.created_at, '2006-01-02T15:04:05.999999-07:00'),
									   'updated_at', p.updated_at
								   )
						   ),
					   'meta', jsonb_build_object(
							   'all_row_count', (select count(*) from persons.persons_table)
						   )
				   )
		from (select pt.id,
					 pt.name,
					 pt.surname,
					 pt.patronymic,
					 pt.age,
					 pt.gender,
					 pt.nation,
					 pt.created_at,
					 pt.updated_at
			  from persons.persons_table as pt %s
			  order by pt.id desc %s
			  ) as p;
	`, filter, pagination)

	var b []byte

	if err := r.db.QueryRow(query).Scan(&b); err != nil {
		return nil, errors.Wrap(err, "GetWithFilterAndPagination #1")
	}

	return b, nil
}

func (r *PersonRepository) Update(req *domain.Person) error {
	query := `
		update persons.persons_table
		set 
		    name = case 
		        		when $1 != ''
		        			then $1
		        		else persons_table.name
		           end,
		    surname = case 
		        		when $2 != ''
		        			then $2
						else persons_table.surname
		              end,
		    patronymic = $3,
		    age = $4,
		    gender = $5,
		    nation = $6,
		    updated_at = $7
		where id = $8
		returning name, surname, patronymic, age, gender, nation, created_at, updated_at;
	`

	if err := r.db.QueryRow(
		query,
		req.Name,
		req.Surname,
		req.Patronymic,
		req.Age,
		req.Gender,
		req.Nation,
		time.Now(),
		req.ID,
	).Scan(
		&req.Name,
		&req.Surname,
		&req.Patronymic,
		&req.Age,
		&req.Gender,
		&req.Nation,
		&req.CreatedAt,
		&req.UpdatedAt,
	); err != nil {
		return errors.Wrap(err, "Update #1")
	}

	return nil
}

func (r *PersonRepository) Delete(id int) (*int64, error) {
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
