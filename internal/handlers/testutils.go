package handlers

import (
	"avito-test/internal/storage"
	"database/sql"
	"log"

	"testing"

	_ "modernc.org/sqlite" // SQLite driver
)

func setupInMemoryDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory DB: %v", err)
	}

	// Таблица PVZ
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS pvz (
		id TEXT PRIMARY KEY,
		registration_date TEXT NOT NULL,
		city TEXT NOT NULL
	);`)
	if err != nil {
		t.Fatalf("failed to create pvz table: %v", err)
	}

	log.Println("Using in-memory SQLite DB for unit test")
	storage.DB = db
	return db
}
func setupReceptionDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory DB: %v", err)
	}

	// Создание таблиц
	_, _ = db.Exec(`CREATE TABLE pvz (
		id TEXT PRIMARY KEY,
		registration_date TEXT NOT NULL,
		city TEXT NOT NULL
	);`)

	_, _ = db.Exec(`CREATE TABLE receptions (
		id TEXT PRIMARY KEY,
		date_time TEXT NOT NULL,
		pvz_id TEXT NOT NULL,
		status TEXT NOT NULL
	);`)

	storage.DB = db

	// Подготовим ПВЗ
	_, _ = db.Exec(`INSERT INTO pvz (id, registration_date, city) VALUES ('pvz-1', '2025-01-01T00:00:00Z', 'Казань');`)

	log.Println("Using in-memory SQLite DB for receptions")
	return db
}
