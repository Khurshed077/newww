package auth

import (
	"database/sql"
	"errors"
	jj "newww/internal/middleware/jwt"
	"newww/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Register(username, email, password string) (*model.User, error) {
	// 1. Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 2. Вставляем пользователя в базу (SQLite)
	res, err := s.db.Exec(
		"INSERT INTO users(username, email, password_hash, admin) VALUES(?, ?, ?, ?)",
		username, email, string(hashedPassword), 1,
	)
	if err != nil {
		return nil, err
	}

	// 3. Получаем ID вставленного пользователя
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 4. Возвращаем объект пользователя
	return &model.User{
		ID:           int(id),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Admin:        1,
	}, nil
}

func (s *AuthService) Login(email, password string) (accessToken string, refreshToken string, err error) {
	var user model.User
	err = s.db.QueryRow(
		"SELECT id, username, email, password_hash, admin FROM users WHERE email=?",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Admin)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Генерация токенов
	accessToken, err = jj.GenerateAccessToken(user.ID, user.Admin)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = jj.GenerateRefreshToken(user.ID, user.Admin)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) GetUserByID(userID int) (*model.User, error) {
	user := &model.User{}
	err := s.db.QueryRow("SELECT id, username, email, admin FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Username, &user.Email, &user.Admin)
	if err != nil {
		return nil, err
	}
	return user, nil
}
