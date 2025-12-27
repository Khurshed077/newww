package user

import "newww/internal/model"

type UserService struct {
	userRepo *UserRepository
}

func NewService(repo *UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) GetAllUsers() ([]*model.User, error) {
	return s.userRepo.GetAll()
}
func (s *UserService) GetByID(id int) (*model.User, error) {
	return s.userRepo.GetByID(id)
}
