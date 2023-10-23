package person

import (
	"namer/internal/domain"
	"strings"
)

func prepareRequest(req *domain.Person) {
	req.Name = strings.ReplaceAll(strings.TrimSpace(req.Name), " ", "")
	req.Surname = strings.ReplaceAll(strings.TrimSpace(req.Surname), " ", "")

	p := strings.ReplaceAll(strings.TrimSpace(*req.Patronymic), " ", "")
	req.Patronymic = &p
}
