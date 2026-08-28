package service

type Service struct {
	storer Storer
}

func NewService(storer Storer) *Service {
	return &Service{storer: storer}
}
