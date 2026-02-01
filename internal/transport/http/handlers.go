package http

type Service interface {

}

type Handlers struct {
	s Service
}

func NewHandlers(s Service) *Handlers {
	return &Handlers{
		s: s,
	}
}