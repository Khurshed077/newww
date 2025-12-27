package category

import "newww/internal/model"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List() ([]model.Category, error) {
	return s.repo.List()
}
