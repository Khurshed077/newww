package article

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "modernc.org/sqlite"

	"newww/internal/apps/article"
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
		id INTEGER PRIMARY KEY,
		username TEXT,
		email TEXT,
		password_hash TEXT,
		admin INTEGER
	);
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY,
		name TEXT
	);
	CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		anons TEXT,
		full_text TEXT,
		image TEXT,
		category_id INTEGER,
		user_id INTEGER
	);
	INSERT OR IGNORE INTO users (id, username, email, password_hash, admin)
	VALUES (1,'test','t@test.com','hash',2);
	INSERT OR IGNORE INTO categories (id, name) VALUES (1,'Tech');
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
func TestCreateArticle(t *testing.T) {
	db := setupTestDB(t, "create_test")
	repo := article.NewRepository(db)
	service := article.NewService(repo)
	handler := article.NewHandler(service)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "Title")
	writer.WriteField("anons", "Anons")
	writer.WriteField("full_text", "Full text")
	writer.WriteField("category_id", "1")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/articles/create", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = authContext(req)

	w := httptest.NewRecorder()
	handler.Create(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var art model.Article
	if err := json.NewDecoder(resp.Body).Decode(&art); err != nil {
		t.Fatalf("cannot decode response: %v", err)
	}

	if art.ID == 0 {
		t.Errorf("expected article ID > 0")
	}
	if art.UserID != 1 {
		t.Errorf("expected UserID=1, got %d", art.UserID)
	}
}

func TestListArticles(t *testing.T) {
	db := setupTestDB(t, "list_test")
	repo := article.NewRepository(db)
	service := article.NewService(repo)
	handler := article.NewHandler(service)

	// создаём статью для проверки
	service.Create(&model.Article{Title: "List", Anons: "A", FullText: "F", UserID: 1, CategoryID: nil})

	req := httptest.NewRequest("GET", "/api/articles", nil)
	req = authContext(req)
	w := httptest.NewRecorder()
	handler.List(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var articles []*model.Article
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		t.Fatalf("cannot decode response: %v", err)
	}

	if len(articles) == 0 {
		t.Errorf("expected at least 1 article")
	}
}

func TestDetailArticle(t *testing.T) {
	db := setupTestDB(t, "detail_test")
	repo := article.NewRepository(db)
	service := article.NewService(repo)
	handler := article.NewHandler(service)

	// создаём статью для проверки
	service.Create(&model.Article{Title: "Detail", Anons: "A", FullText: "F", UserID: 1, CategoryID: nil})

	req := httptest.NewRequest("GET", "/api/articles/detail?id=1", nil)
	w := httptest.NewRecorder()
	handler.Detail(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var result model.Article
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("cannot decode response: %v", err)
	}

	if result.ID != 1 {
		t.Errorf("expected ID=1, got %d", result.ID)
	}
}

func TestDeleteArticle(t *testing.T) {
	db := setupTestDB(t, "delete")
	repo := article.NewRepository(db)
	service := article.NewService(repo)
	handler := article.NewHandler(service)

	// создаём статью для удаления
	service.Create(&model.Article{Title: "Delete", Anons: "A", FullText: "F", UserID: 1, CategoryID: nil})

	req := httptest.NewRequest("DELETE", "/api/articles/delete?id=1", nil)
	req = authContext(req)
	w := httptest.NewRecorder()
	handler.Delete(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d", resp.StatusCode)
	}
}
