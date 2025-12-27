package user_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "modernc.org/sqlite"

	"newww/internal/apps/user"
	"newww/internal/middleware/jwt"
	"newww/internal/model"
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
	INSERT OR IGNORE INTO users (id, username, email, password_hash, admin)
	VALUES (1,'test','t@test.com','hash',2);
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
func TestGetUser(t *testing.T) {
	db := setupTestDB(t, "user_test")
	repo := user.NewRepository(db)
	service := user.NewService(repo)
	handler := user.NewUserHandler(service)

	req := httptest.NewRequest("GET", "/api/users", nil)
	req = authContext(req)
	w := httptest.NewRecorder()

	handler.Users(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var u model.User
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		t.Fatalf("cannot decode response: %v", err)
	}

	if u.ID != 1 {
		t.Errorf("expected ID=1, got %d", u.ID)
	}
	if u.Username != "test" {
		t.Errorf("expected username 'test', got '%s'", u.Username)
	}
}
