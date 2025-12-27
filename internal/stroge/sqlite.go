package stroge

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func NewSQLite(path string) *sql.DB {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		panic(err)
	}
	return db
}

func InitDB(db *sql.DB) {

	// Создаём таблицы, если их нет
	db.Exec(`
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
	admin INT DEFAULT 1 CHECK (admin IN (1,2))
);

CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    anons TEXT NOT NULL,
    full_text TEXT NOT NULL,
    image TEXT,
    category_id INTEGER,
	user_id	INTEGER,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);
`)
	// если нет котегорий добавим их
	db.Exec(`INSERT INTO categories (name) SELECT 'Новости' WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name='Новости');`)
	db.Exec(`INSERT INTO categories (name) SELECT 'Статьи' WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name='Статьи');`)
	db.Exec(`INSERT INTO categories (name) SELECT 'Игры' WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name='Игры');`)

	log.Println("SQLite (modernc) инициализирован")
}
