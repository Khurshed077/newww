package article

import (
	"errors"
	"newww/internal/model"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Delete(
	articleID int,
	userID int,
	admin int,
) error {

	article, err := s.repo.GetByID(articleID)
	if err != nil {
		return err
	}

	if !CanDelete(userID, admin, article) {
		return errors.New("forbidden")
	}

	return s.repo.Delete(articleID)
}
func (s *Service) Create(a *model.Article) error {
	id, err := s.repo.Insert(a)
	if err != nil {
		return err
	}
	a.ID = int(id) // 🔹 присваиваем ID вставленной статьи
	return nil
}

func (s *Service) List(userID int, admin int, categoryID *int) ([]*model.Article, error) {
	var articles []*model.Article
	var err error

	if categoryID != nil {
		articles, err = s.repo.GetByCategory(*categoryID)
	} else if admin == 2 {
		articles, err = s.repo.GetAll()
	} else {
		articles, err = s.repo.GetByUser(userID)
	}

	return articles, err
}

func (s *Service) Update(a *model.Article, userID int, admin int) error {
	if !CanEdit(userID, admin, a) {
		return errors.New("forbidden")
	}
	return s.repo.Update(a)
}

func CanEdit(userID, admin int, a *model.Article) bool {
	if admin == 2 {
		return true
	}
	return a.UserID == userID
}

func (s *Service) GetByID(id int) (*model.Article, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Public(categoryID *int) ([]model.Article, error) {
	return s.repo.ListPublic(categoryID)
}

func (s *Service) Dashboard(userID int, admin bool) ([]model.Article, error) {
	// НИКАКИХ ПРОВЕРОК НА admin == true
	return s.repo.ListDashboard(userID, admin)
}
func (s *Service) ListDashboard(userID int, admin bool) ([]model.Article, error) {
	return s.repo.ListDashboard(userID, admin)
}

func (s *Service) Categories() ([]model.Category, error) {
	return s.repo.Categories()
}
