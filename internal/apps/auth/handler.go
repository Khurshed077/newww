package auth

import (
	"encoding/json"
	"net/http"
	"newww/internal/middleware/jwt"
)

func NewAuthHandler(s *AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// @Summary Register a new user
// @Description Создать нового пользователя и получить токены
// @Tags Auth
// @Accept json
// @Produce json
// @Param register body auth.RegisterRequest true "User registration info"
// @Success 200 {object} auth.UserResponse
// @Failure 400 {string} string "invalid json"
// @Router /api/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Регистрируем пользователя
	user, err := h.service.Register(req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ---------- Генерируем токены ----------
	accessToken, refreshToken, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "login after registration failed", http.StatusInternalServerError)
		return
	}

	// Отправляем cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/",
	})

	// Ответ JSON с данными пользователя
	resp := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Admin:    user.Admin,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary Login user
// @Description Вход пользователя и установка access/refresh токенов в cookie
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body auth.LoginRequest true "User login info"
// @Success 200 {object} map[string]string
// @Failure 401 {string} string "invalid credentials"
// @Router /api/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// @Summary Logout user
// @Description Удалить access и refresh токены
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Удаляем access token
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	// Удаляем refresh token
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "logged out"})
}

// @Summary Refresh access token
// @Description Получить новый access token используя refresh token
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 401 {string} string "invalid refresh token"
// @Router /api/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "refresh token missing", http.StatusUnauthorized)
		return
	}

	claims, err := jwt.ParseRefreshToken(cookie.Value)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetUserByID(claims.UserID)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	newAccessToken, err := jwt.GenerateAccessToken(user.ID, user.Admin)
	if err != nil {
		http.Error(w, "could not generate access token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "access token refreshed",
		"user":    user.Username,
		"user_id": user.ID,
		"admin":   user.Admin,
	})
}
