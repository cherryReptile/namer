package person

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"namer/internal/domain/external"
	"net/http"
	"time"
)

type APIRepository struct {
	client    *http.Client
	ageURL    string
	genderURL string
	nationURL string
}

func NewRepository() *APIRepository {
	return &APIRepository{
		client: &http.Client{
			Timeout: time.Second * 15,
		},
		ageURL:    "https://api.agify.io/?name=%s",
		genderURL: "https://api.genderize.io/?name=%s",
		nationURL: "https://api.nationalize.io/?name=%s",
	}
}

func (a *APIRepository) GetNameInfo(name string) (*external.ExternalResponse, error) {
	var (
		res external.ExternalResponse
		err error
	)

	if res.Agify, err = a.GetAge(name); err != nil {
		return nil, errors.Wrap(err, "GetNameInfo #1")
	}

	if res.Agify.Error != nil {
		res.Error, res.StatusCode = res.Agify.Error, res.Agify.StatusCode

		return &res, nil
	}

	if res.Genderize, err = a.GetGender(name); err != nil {
		return nil, errors.Wrap(err, "GetNameInfo #2")
	}

	if res.Genderize.Error != nil {
		res.Error, res.StatusCode = res.Genderize.Error, res.Genderize.StatusCode

		return &res, nil
	}

	if res.Nationalize, err = a.GetNation(name); err != nil {
		return nil, errors.Wrap(err, "GetNameInfo #3")
	}

	if res.Nationalize.Error != nil {
		res.Error, res.StatusCode = res.Nationalize.Error, res.Nationalize.StatusCode

		return &res, nil
	}

	return &res, nil
}

func (a *APIRepository) GetAge(name string) (*external.AgifyResponse, error) {
	res, err := http.Get(
		fmt.Sprintf(a.ageURL, name),
	)

	if err != nil {
		return nil, errors.Wrap(err, "GetAge #1")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "GetAge #2")
	}

	var agify external.AgifyResponse

	if err = json.Unmarshal(body, &agify); err != nil {
		return nil, errors.Wrap(err, "GetAge #3")
	}

	agify.StatusCode = res.StatusCode

	return &agify, nil
}

func (a *APIRepository) GetGender(name string) (*external.GenderizeResponse, error) {
	res, err := http.Get(
		fmt.Sprintf(a.genderURL, name),
	)

	if err != nil {
		return nil, errors.Wrap(err, "GetGender #1")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "GetGender #2")
	}

	var genderize external.GenderizeResponse

	if err = json.Unmarshal(body, &genderize); err != nil {
		return nil, errors.Wrap(err, "GetGender #3")
	}

	genderize.StatusCode = res.StatusCode

	return &genderize, nil
}

func (a *APIRepository) GetNation(name string) (*external.NationalizeResponse, error) {
	res, err := http.Get(
		fmt.Sprintf(a.nationURL, name),
	)

	if err != nil {
		return nil, errors.Wrap(err, "GetNation #1")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "GetNation #2")
	}

	var nationalize external.NationalizeResponse

	if err = json.Unmarshal(body, &nationalize); err != nil {
		return nil, errors.Wrap(err, "GetNation #3")
	}

	nationalize.StatusCode = res.StatusCode

	return &nationalize, nil
}
