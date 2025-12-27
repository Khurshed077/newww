package auth

import (
	"database/sql"
	"newww/internal/model"
)

type AuthRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) GetByEmail(email string) (*model.User, error) {
	u := &model.User{}

	err := r.db.QueryRow(`
		SELECT id, username, email, password_hash, admin
		FROM users
		WHERE email = $1
	`, email).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.Admin,
	)

	if err != nil {
		return nil, err
	}

	return u, nil
}

// ВОТ ЧЕГО НЕ ХВАТАЛО
func (r *AuthRepository) Create(u *model.User) error {
	return r.db.QueryRow(`
		INSERT INTO users (username, email, password_hash, admin)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, u.Username, u.Email, u.PasswordHash, u.Admin).Scan(&u.ID)
}
