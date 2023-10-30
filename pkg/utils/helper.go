package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"namer/internal/domain"
	"strings"
)

func StringToPtr(str string) *string {
	return &str
}

func PrepareRequest(req *domain.Person) {
	req.Name = strings.ReplaceAll(strings.TrimSpace(req.Name), " ", "")
	req.Surname = strings.ReplaceAll(strings.TrimSpace(req.Surname), " ", "")

	if req.Patronymic != nil {
		p := strings.ReplaceAll(strings.TrimSpace(*req.Patronymic), " ", "")
		req.Patronymic = &p
	}
}

func GetFilterAndPagination(req *domain.FilterWithPagination, alias string) ([]string, error) {
	var filter string

	if req.Filter != nil {
		for i := range req.Filter {
			if strings.Contains(req.Filter[i].Field, " ") {
				return nil, errors.New("invalid field")
			}

			if req.Filter[i].Field == "" {
				return nil, errors.New("empty field")
			}

			if req.Filter[i].Value == "" {
				return nil, errors.New("empty value")
			}

			if i == 0 {
				filter = fmt.Sprintf(
					"where %s.%s ilike '%s%s%s'",
					alias,
					req.Filter[i].Field,
					"%",
					req.Filter[i].Value,
					"%",
				)

				continue
			}

			filter = fmt.Sprintf(
				"%s and %s.%s ilike '%s%s%s'",
				filter,
				alias,
				req.Filter[i].Field,
				"%",
				req.Filter[i].Value,
				"%",
			)
		}
	}

	if req.Pagination == nil {
		req.Pagination = &domain.Pagination{
			Limit: 5,
			Page:  1,
		}
	}

	var pagination string

	if req.Pagination.Page < 1 {
		return nil, errors.New("invalid page")
	}

	if req.Pagination.Limit < 1 {
		return nil, errors.New("invalid limit")
	}

	if req.Pagination.Page == 1 {
		pagination = fmt.Sprintf(
			"limit %d",
			req.Pagination.Limit,
		)
	} else {
		pagination = fmt.Sprintf(
			"limit %d offset %d",
			req.Pagination.Limit,
			(req.Pagination.Page-1)*req.Pagination.Limit,
		)
	}

	return []string{filter, pagination}, nil
}
