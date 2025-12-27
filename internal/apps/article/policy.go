package article

import "newww/internal/model"

func CanDelete(userID int, admin int, article *model.Article) bool {
	if admin == 2 {
		return true
	}
	return article.UserID == userID
}
func CanEditArticle(userID, admin int, a *model.Article) bool {
	if admin == 2 {
		return true
	}
	return a.UserID == userID
}
