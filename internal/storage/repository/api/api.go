package api

type API interface {
}

type api struct{}

func NewAPI() API {
	return &api{}
}
