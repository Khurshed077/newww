package article

import (
	"encoding/json"
	"io"
	"net/http"
	"newww/internal/middleware/jwt"
	"newww/internal/model"
	"os"
	"strconv"
	"strings"
)

// Handler struct
type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// ErrorResponse стандарт для ошибок
type ErrorResponse struct {
	Error string `json:"error"`
}

// ------------------ Handlers ------------------

// @Summary Create new article
// @Description Создать статью с изображением
// @Tags Articles
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Title"
// @Param anons formData string true "Anons"
// @Param full_text formData string true "Full text"
// @Param category_id formData int false "Category ID"
// @Param image formData file false "Image"
// @Success 200 {object} model.Article
// @Failure 400 {object} article.ErrorResponse
// @Failure 401 {object} article.ErrorResponse
// @Security BearerAuth
// @Router /api/articles/create [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(jwt.ClaimsKey).(*jwt.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		categories, err := h.service.Categories()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to load categories"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categories)
		return
	}

	contentType := r.Header.Get("Content-Type")
	var article model.Article

	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid json"})
			return
		}
	} else {
		article.Title = r.FormValue("title")
		article.Anons = r.FormValue("anons")
		article.FullText = r.FormValue("full_text")

		if cid := r.FormValue("category_id"); cid != "" {
			if id, err := strconv.Atoi(cid); err == nil {
				article.CategoryID = &id
			}
		}

		file, header, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			os.MkdirAll("uploads", 0755)
			dst, _ := os.Create("uploads/" + header.Filename)
			defer dst.Close()
			io.Copy(dst, file)
			article.Image = header.Filename
		}
	}

	article.UserID = claims.UserID
	if err := h.service.Create(&article); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// List возвращает список статей (JSON)
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(jwt.ClaimsKey).(*jwt.Claims)
	var categoryID *int
	if cid := r.URL.Query().Get("category"); cid != "" {
		if id, err := strconv.Atoi(cid); err == nil {
			categoryID = &id
		}
	}

	userID := 0
	admin := 0
	if claims != nil {
		userID = claims.UserID
		admin = claims.Admin
	}

	articles, err := h.service.List(userID, admin, categoryID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to load articles"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

// @Summary Get article details
// @Description Get article by ID
// @Tags Articles
// @Produce json
// @Param id query int true "Article ID"
// @Success 200 {object} model.Article
// @Failure 404 {string} string "article not found"
// @Router /api/articles/detail [get]
func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid id"})
		return
	}

	article, err := h.service.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "article not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// @Summary Edit article
// @Description Обновить статью
// @Tags Articles
// @Accept json
// @Produce json
// @Param id query int true "Article ID"
// @Param edit body model.Article true "Данные статьи"
// @Success 200 {object} model.Article
// @Failure 400 {object} article.ErrorResponse
// @Failure 403 {object} article.ErrorResponse
// @Security BearerAuth
// @Router /api/articles/edit [put]
func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(jwt.ClaimsKey).(*jwt.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	article, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "article not found", http.StatusNotFound)
		return
	}

	if !CanEdit(claims.UserID, claims.Admin, article) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// Определяем, пришёл ли запрос через JSON API
	isAPI := strings.HasPrefix(r.URL.Path, "/api/")

	if isAPI && r.Header.Get("Content-Type") == "application/json" {
		// ---------- JSON API ----------
		var input struct {
			Title      string `json:"title"`
			Anons      string `json:"anons"`
			FullText   string `json:"full_text"`
			CategoryID *int   `json:"category_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		article.Title = input.Title
		article.Anons = input.Anons
		article.FullText = input.FullText
		article.CategoryID = input.CategoryID

		if err := h.service.Update(article, claims.UserID, claims.Admin); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(article)
		return
	}

	// ---------- FORM (HTML / FormData) ----------
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	article.Title = r.FormValue("title")
	article.Anons = r.FormValue("anons")
	article.FullText = r.FormValue("full_text")

	if cid := r.FormValue("category_id"); cid != "" {
		cidInt, _ := strconv.Atoi(cid)
		article.CategoryID = &cidInt
	}

	// Работа с файлом (если есть)
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		os.MkdirAll("uploads", 0755)
		dst, _ := os.Create("uploads/" + header.Filename)
		defer dst.Close()
		io.Copy(dst, file)
		article.Image = header.Filename
	}

	if err := h.service.Update(article, claims.UserID, claims.Admin); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Перенаправление на dashboard после редактирования
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// @Summary Delete article
// @Description Удалить статью по ID
// @Tags Articles
// @Produce json
// @Param id query int true "Article ID"
// @Success 204 {object} map[string]string "deleted"
// @Failure 403 {object} article.ErrorResponse
// @Failure 404 {object} article.ErrorResponse
// @Security BearerAuth
// @Router /api/articles/delete [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(jwt.ClaimsKey).(*jwt.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid id"})
		return
	}

	if err := h.service.Delete(id, claims.UserID, claims.Admin); err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "forbidden"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary List articles
// @Description Get all articles or filter by category
// @Tags Articles
// @Produce json
// @Param category_id query int false "Category ID for filtering"
// @Success 200 {object} []model.Article
// @Router /api/articles [get]
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	var categoryID *int
	if c := r.URL.Query().Get("category_id"); c != "" {
		if id, err := strconv.Atoi(c); err == nil {
			categoryID = &id
		}
	}

	articles, err := h.service.Public(categoryID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to load articles"})
		return
	}

	categories, err := h.service.Categories()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to load categories"})
		return
	}

	_, authorized := r.Context().Value(jwt.ClaimsKey).(*jwt.Claims)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"articles":   articles,
		"categories": categories,
		"authorized": authorized,
	})
}

// @Summary Dashboard
// @Description Получить статьи текущего пользователя или все, если Admin
// @Tags Dashboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} article.ErrorResponse
// @Security BearerAuth
// @Router /dashboard [get]
func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Получаем данные пользователя из JWT
	claims, ok := r.Context().Value(jwt.ClaimsKey).(*jwt.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	isAdmin := claims.Admin == 2

	// Получаем статьи
	articles, err := h.service.ListDashboard(claims.UserID, isAdmin)
	if err != nil {
		http.Error(w, "failed to load articles", http.StatusInternalServerError)
		return
	}

	// Формируем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":  claims.UserID,
		"admin":    claims.Admin,
		"articles": articles,
	})
}
