package article

import (
	"database/sql"
	"newww/internal/model"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]*model.Article, error) {
	rows, err := r.db.Query(`
	 SELECT id, title, anons, full_text, image, category_id, user_id
	 FROM articles
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*model.Article
	for rows.Next() {
		a := &model.Article{} // 🔹 создаём новый объект для каждой строки
		err := rows.Scan(&a.ID, &a.Title, &a.Anons, &a.FullText, &a.Image, &a.CategoryID, &a.UserID)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a) // 🔹 добавляем в слайс
	}

	return articles, nil
}

func (r *Repository) GetByUser(userID int) ([]*model.Article, error) {
	rows, err := r.db.Query(`
  SELECT id, title, anons, full_text, image, category_id, user_id
  FROM articles WHERE user_id = ?
 `, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*model.Article
	for rows.Next() {
		a := &model.Article{}
		rows.Scan(
			&a.ID,
			&a.Title,
			&a.Anons,
			&a.FullText,
			&a.Image,
			&a.CategoryID,
			&a.UserID,
		)
		articles = append(articles, a)
	}

	return articles, nil
}

func (r *Repository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM articles WHERE id = ?", id)
	return err
}

func (r *Repository) GetByID(id int) (*model.Article, error) {
	a := &model.Article{}

	row := r.db.QueryRow(`
	 SELECT id, title, anons, full_text, image, category_id, user_id
	 FROM articles
	 WHERE id = ?
	`, id)

	err := row.Scan(
		&a.ID,
		&a.Title,
		&a.Anons,
		&a.FullText,
		&a.Image,
		&a.CategoryID,
		&a.UserID,
	)
	if err != nil {
		return nil, err
	}

	return a, nil
}
func (r *Repository) Insert(a *model.Article) (int64, error) {
	var category sql.NullInt64
	if a.CategoryID != nil {
		category = sql.NullInt64{Int64: int64(*a.CategoryID), Valid: true}
	} else {
		category = sql.NullInt64{Valid: false}
	}

	res, err := r.db.Exec(`
	 INSERT INTO articles (title, anons, full_text, image, category_id, user_id)
	 VALUES (?, ?, ?, ?, ?, ?)
	`, a.Title, a.Anons, a.FullText, a.Image, category, a.UserID)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (r *Repository) GetByCategory(categoryID int) ([]*model.Article, error) {
	rows, err := r.db.Query(`
	 SELECT id, title, anons, full_text, image, category_id, user_id
	 FROM articles WHERE category_id = ?
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*model.Article
	for rows.Next() {
		a := &model.Article{}
		rows.Scan(&a.ID, &a.Title, &a.Anons, &a.FullText, &a.Image, &a.CategoryID, &a.UserID)
		articles = append(articles, a)
	}
	return articles, nil
}

func (r *Repository) Update(a *model.Article) error {
	var category sql.NullInt64
	if a.CategoryID != nil {
		category = sql.NullInt64{Int64: int64(*a.CategoryID), Valid: true}
	} else {
		category = sql.NullInt64{Valid: false}
	}

	_, err := r.db.Exec(`
	 UPDATE articles
	 SET title=?, anons=?, full_text=?, image=?, category_id=?
	 WHERE id=?
	`, a.Title, a.Anons, a.FullText, a.Image, category, a.ID)

	return err
}
func (r *Repository) ListPublic(categoryID *int) ([]model.Article, error) {
	query := `SELECT id, title, anons, image, category_id, user_id FROM articles`
	args := []any{}

	if categoryID != nil {
		query += " WHERE category_id = ?"
		args = append(args, *categoryID)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var a model.Article
		var cid sql.NullInt64

		if err := rows.Scan(
			&a.ID, &a.Title, &a.Anons, &a.Image, &cid, &a.UserID,
		); err != nil {
			return nil, err
		}
		if cid.Valid {
			c := int(cid.Int64)
			a.CategoryID = &c
		} else {
			a.CategoryID = nil
		}
		articles = append(articles, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return articles, nil
}
func (r *Repository) ListDashboard(userID int, admin bool) ([]model.Article, error) {
	query := `SELECT id, title, anons, full_text, image, category_id, user_id FROM articles`
	args := []any{}

	if !admin {
		query += " WHERE user_id = ?"
		args = append(args, userID)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var a model.Article
		var cid sql.NullInt64

		if err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Anons,
			&a.FullText,
			&a.Image,
			&cid,
			&a.UserID,
		); err != nil {
			return nil, err
		}

		if cid.Valid {
			c := int(cid.Int64)
			a.CategoryID = &c
		}

		articles = append(articles, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *Repository) Categories() ([]model.Category, error) {
	rows, err := r.db.Query(`SELECT id,name FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []model.Category
	for rows.Next() {
		var c model.Category
		rows.Scan(&c.ID, &c.Name)
		categories = append(categories, c)
	}
	return categories, nil
}
