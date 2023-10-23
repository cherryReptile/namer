package api

type API interface {
	GetAge() (*int, error)
	GetGender() (*string, error)
	GetNation() (*string, error)
}

type api struct{}

func NewAPI() API {
	return &api{}
}

func (a *api) GetAge() (*int, error) {
	return nil, nil
}

func (a *api) GetGender() (*string, error) {
	return nil, nil
}

func (a *api) GetNation() (*string, error) {
	return nil, nil
}
