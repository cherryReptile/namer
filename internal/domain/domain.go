package domain

import "time"

type Person struct {
	ID         int        `json:"id"`
	Name       string     `json:"name" form:"name"`
	Surname    string     `json:"surname" form:"surname"`
	Patronymic *string    `json:"patronymic" form:"patronymic"`
	Age        *int       `json:"age"`
	Gender     *string    `json:"gender"`
	Nation     *string    `json:"nation"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type FilterWithPagination struct {
	Filter     []Filter    `form:"filter"`
	Pagination *Pagination `form:"pagination"`
}

type Filter struct {
	Field string `form:"field"`
	Value string `form:"value"`
}

type Pagination struct {
	Limit int `form:"limit"`
	Page  int `form:"page"`
}

type Response struct {
	Data  any     `json:"data,omitempty"`
	Error *string `json:"error,omitempty"`
	Meta  *struct {
		AllRowCount int `json:"all_row_count"`
	} `json:"meta,omitempty"`
	StatusCode int `json:"-"`
}
