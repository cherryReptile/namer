package external

type ExternalResponse struct {
	Agify       *AgifyResponse
	Genderize   *GenderizeResponse
	Nationalize *NationalizeResponse
	Error       *string
	StatusCode  int
}

type AgifyResponse struct {
	Count      *int    `json:"count"`
	Name       *string `json:"name"`
	Age        *int    `json:"age"`
	Error      *string `json:"error"`
	StatusCode int     `json:"-"`
}

type GenderizeResponse struct {
	Count       *int     `json:"count"`
	Name        *string  `json:"name"`
	Gender      *string  `json:"gender"`
	Probability *float64 `json:"probability"`
	Error       *string  `json:"error"`
	StatusCode  int      `json:"-"`
}

type NationalizeResponse struct {
	Count   *int    `json:"count"`
	Name    *string `json:"name"`
	Country []struct {
		CountryId   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
	Error      *string `json:"error"`
	StatusCode int     `json:"-"`
}
