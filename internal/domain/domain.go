package domain

import "time"

type Person struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Surname    string     `json:"surname"`
	Patronymic *string    `json:"patronymic"`
	Age        *int       `json:"age"`
	Gender     *string    `json:"gender"`
	Nation     *string    `json:"nation"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type FilterWithPagination struct {
	Filter     []Filter    `json:"filter"`
	Pagination *Pagination `json:"pagination"`
}

type Filter struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type Response struct {
	Data       any     `json:"data,omitempty"`
	Error      *string `json:"error,omitempty"`
	StatusCode int     `json:"-"`
}
