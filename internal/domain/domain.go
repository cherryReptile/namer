package domain

import "time"

type Person struct {
	ID         int        `json:"id"`
	Name       string     `form:"name"`
	Surname    string     `form:"surname"`
	Patronymic *string    `form:"patronymic"`
	Age        *int       `json:"age"`
	Gender     *string    `json:"gender"`
	Nation     *string    `json:"nation"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type FilterWithPagination struct {
	Filter []struct {
		Field string `form:"field"`
		Value string `form:"field"`
	} `form:"filter"`

	Pagination *struct {
		Limit int `form:"limit"`
		Page  int `form:"page"`
	} `form:"pagination"`
}

type Response struct {
	Data       any     `json:"data,omitempty"`
	Error      *string `json:"error,omitempty"`
	StatusCode int     `json:"-"`
}
