package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "modernc.org/sqlite"

	"newww/internal/apps/auth"
	"newww/internal/middleware/jwt"
)

// ------------------ Helpers ------------------
func setupTestDB(t *testing.T, name string) *sql.DB {
	db, err := sql.Open("sqlite", "file:"+name+"?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		email TEXT,
		password_hash TEXT,
		admin INTEGER
	);
	`
	if _, err := db.Exec(sqlStmt); err != nil {
		t.Fatal(err)
	}

	return db
}

func authContext(req *http.Request) *http.Request {
	claims := &jwt.Claims{UserID: 1, Admin: 2}
	return req.WithContext(jwt.SetClaimsToContext(req.Context(), claims))
}

// ------------------ Tests ------------------
func TestRegisterLoginLogout(t *testing.T) {
	db := setupTestDB(t, "auth_test")
	service := auth.NewAuthService(db)
	handler := auth.NewAuthHandler(service)

	// ---------- Register ----------
	regBody := map[string]string{
		"username": "testuser",
		"email":    "test@test.com",
		"password": "password123",
	}
	bodyBytes, _ := json.Marshal(regBody)
	req := httptest.NewRequest("POST", "/api/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Register(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for register, got %d", resp.StatusCode)
	}

	var userResp auth.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		t.Fatalf("cannot decode register response: %v", err)
	}
	if userResp.ID == 0 || userResp.Email != "test@test.com" {
		t.Fatalf("unexpected register response: %+v", userResp)
	}

	// ---------- Login ----------
	loginBody := map[string]string{
		"email":    "test@test.com",
		"password": "password123",
	}
	bodyBytes, _ = json.Marshal(loginBody)
	req = httptest.NewRequest("POST", "/api/login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.Login(w, req)

	resp = w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for login, got %d", resp.StatusCode)
	}

	// ---------- Refresh Token ----------
	// Для теста ставим refresh_token вручную
	refreshToken, _ := jwt.GenerateRefreshToken(userResp.ID, userResp.Admin)
	req = httptest.NewRequest("POST", "/api/refresh", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken})
	w = httptest.NewRecorder()
	handler.RefreshToken(w, req)

	resp = w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for refresh token, got %d", resp.StatusCode)
	}

	// ---------- Logout ----------
	req = httptest.NewRequest("POST", "/api/logout", nil)
	w = httptest.NewRecorder()
	handler.Logout(w, req)

	resp = w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for logout, got %d", resp.StatusCode)
	}
}
