package user

import (
	"database/sql"
	"newww/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(
		"SELECT id, username, email, password_hash, admin FROM users WHERE email = ?",
		email,
	).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Admin)

	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) Create(u *model.User) error {
	res, err := r.db.Exec(
		"INSERT INTO users (username, email, password_hash, admin) VALUES (?, ?, ?, ?)",
		u.Username, u.Email, u.PasswordHash, u.Admin,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err == nil {
		u.ID = int(id)
	}
	return nil
}

func (r *UserRepository) GetAll() ([]*model.User, error) {
	rows, err := r.db.Query(
		"SELECT id, username, email, password_hash, admin FROM users",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.PasswordHash,
			&u.Admin,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) GetByID(id int) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(
		"SELECT id, username, email, password_hash, admin FROM users WHERE id = ?",
		id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Admin)

	if err != nil {
		return nil, err
	}
	return u, nil
}
