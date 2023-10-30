package person

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"namer/internal/domain"
	"namer/pkg/utils"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func connectToDB() (*sql.DB, error) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if err = godotenv.Load(filepath.Join(filepath.Dir(wd), "../../../..", ".env")); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST_TEST"),
			os.Getenv("DB_PORT_TEST"),
			os.Getenv("DB_USER_TEST"),
			os.Getenv("DB_PASSWORD_TEST"),
			os.Getenv("DB_NAME_TEST"),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err = db.PingContext(timeout); err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func TestCreate(t *testing.T) {
	db, err := connectToDB()

	require.NoError(t, err)
	require.NotNil(t, db)

	repo := NewRepository(db)

	t.Run("success", func(t *testing.T) {
		req := domain.Person{
			Name:    "Test",
			Surname: "Test",
		}

		createPerson(t, repo, &req)

		deletePerson(t, repo, &req)
	})

	t.Run("error", func(t *testing.T) {
		bad := domain.Person{
			Gender: utils.StringToPtr("Test"),
		}

		assert.Error(t, repo.Create(&bad))
	})
}

func TestGetByID(t *testing.T) {
	db, err := connectToDB()

	require.NoError(t, err)
	require.NotNil(t, db)

	repo := NewRepository(db)

	t.Run("success", func(t *testing.T) {
		req := domain.Person{
			Name:    "Test",
			Surname: "Test",
		}

		createPerson(t, repo, &req)

		p, err := repo.GetByID(req.ID)

		require.NoError(t, err)

		assert.Equal(t, p, &req)

		deletePerson(t, repo, p)
	})

	t.Run("error", func(t *testing.T) {
		badP, err := repo.GetByID(0)

		assert.Nil(t, badP)
		if assert.Error(t, err) {
			assert.Equal(t, sql.ErrNoRows, errors.Cause(err))
		}
	})
}

func TestGetWithFilterAndPagination(t *testing.T) {
	db, err := connectToDB()

	require.NoError(t, err)
	require.NotNil(t, db)

	repo := NewRepository(db)

	t.Run("success", func(t *testing.T) {
		var persons []domain.Person

		for i := 0; i < 5; i++ {
			person := domain.Person{
				Name:    fmt.Sprintf("Test%d", i),
				Surname: fmt.Sprintf("Test%d", i),
			}

			createPerson(t, repo, &person)

			persons = append(persons, person)
		}

		req := domain.FilterWithPagination{
			Filter: []domain.Filter{
				{"name", "te"},
				{"surname", "t"},
			},
			Pagination: &domain.Pagination{
				Page:  1,
				Limit: 5,
			},
		}

		result, err := utils.GetFilterAndPagination(&req, "pt")
		assert.NoError(t, err)
		require.Equal(t, len(result), 2)

		b, err := repo.GetWithFilterAndPagination(result[0], result[1])
		assert.NoError(t, err)
		assert.NotNil(t, b)

		resp := struct {
			Data []domain.Person `json:"data"`
			Meta *struct {
				AllRowCount int `json:"all_row_count"`
			}
		}{}

		assert.Nil(t, json.Unmarshal(b, &resp))

		if assert.NotNil(t, resp.Data) {
			for i := range resp.Data {
				if i == 0 {
					assert.Equal(t, resp.Data[i].Name, persons[len(persons)-1].Name)
					assert.Equal(t, resp.Data[i].Surname, persons[len(persons)-1].Surname)
				} else {
					assert.Equal(t, resp.Data[i].Name, persons[len(persons)-i-1].Name)
					assert.Equal(t, resp.Data[i].Surname, persons[len(persons)-i-1].Surname)
				}

				assert.Equal(t, resp.Meta.AllRowCount, len(persons))

				deletePerson(t, repo, &resp.Data[i])
			}
		}
	})

	t.Run("error", func(t *testing.T) {
		b, err := repo.GetWithFilterAndPagination("test", "test")
		assert.Nil(t, b)
		assert.Error(t, err)
	})
}

func TestUpdate(t *testing.T) {
	db, err := connectToDB()

	require.NoError(t, err)
	require.NotNil(t, db)

	repo := NewRepository(db)

	t.Run("success", func(t *testing.T) {
		req := domain.Person{
			Name:    "Test",
			Surname: "Test",
		}

		createPerson(t, repo, &req)

		req.Name, req.Surname, req.Gender = "Test1", "Test1", utils.StringToPtr("male")

		assert.Nil(t, repo.Update(&req))
		assert.Equal(t, req.Name, "Test1")
		assert.Equal(t, req.Surname, "Test1")
		assert.Equal(t, req.Gender, utils.StringToPtr("male"))
		assert.NotNil(t, req.UpdatedAt)

		deletePerson(t, repo, &req)
	})

	t.Run("error", func(t *testing.T) {
		err = repo.Update(&domain.Person{})
		if assert.Error(t, err) {
			assert.Equal(t, sql.ErrNoRows, errors.Cause(err))
		}
	})
}

func TestDelete(t *testing.T) {
	db, err := connectToDB()

	require.NoError(t, err)
	require.NotNil(t, db)

	repo := NewRepository(db)

	t.Run("success", func(t *testing.T) {
		req := domain.Person{
			Name:    "Test",
			Surname: "Test",
		}

		createPerson(t, repo, &req)

		deletePerson(t, repo, &req)
	})

	t.Run("error", func(t *testing.T) {
		aff, err := repo.Delete(0)

		assert.NoError(t, err)
		assert.Equal(t, *aff, int64(0))
	})
}

func createPerson(t *testing.T, repo *PersonRepository, person *domain.Person) {
	require.NoError(t, repo.Create(person))

	assert.NotEqual(t, person.ID, 0)
	assert.Equal(t, person.Name, person.Name)
	assert.Equal(t, person.Surname, person.Surname)
	assert.NotNil(t, person.CreatedAt)
}

func deletePerson(t *testing.T, repo *PersonRepository, person *domain.Person) {
	aff, err := repo.Delete(person.ID)

	require.NoError(t, err)
	require.NotEqual(t, *aff, 0)
}
