package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
	"meal-planner/internal/config"
)

// Database инкапсулирует подключение к БД
type Database struct {
	db *sql.DB
}

// New создает новое подключение к БД
func New(cfg *config.Config) (*Database, error) {
	dbURL := cfg.GetDatabaseURL()

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем подключение
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✓ Connected to database")

	d := &Database{db: conn}

	// Выполняем миграции
	if err := d.runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return d, nil
}

// runMigrations выполняет все SQL миграции из папки migrations/
func (d *Database) runMigrations() error {
	// Читаем миграции из папки
	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			filePath := filepath.Join("migrations", f.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", filePath, err)
			}

			// Выполняем SQL
			if _, err := d.db.Exec(string(content)); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", filePath, err)
			}

			log.Printf("✓ Executed migration: %s", f.Name())
		}
	}

	return nil
}

// Close закрывает подключение к БД
func (d *Database) Close() error {
	return d.db.Close()
}

// GetDB возвращает само подключение (для запросов)
func (d *Database) GetDB() *sql.DB {
	return d.db
}
