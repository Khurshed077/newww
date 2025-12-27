package model

type Article struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Anons      string `json:"anons"`
	FullText   string `json:"full_text"`
	Image      string `json:"image"`
	CategoryID *int   `json:"category_id"`
	UserID     int    `json:"user_id"`
}

func CanEdit(userID, admin int, a *Article) bool {
	if admin == 2 {
		return true
	}
	return a.UserID == userID
}
