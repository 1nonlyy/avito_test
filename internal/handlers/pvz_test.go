package handlers

import (
	"avito-test/internal/storage"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreatePVZ_ValidCity(t *testing.T) {
	// Arrange
	router := gin.New()
	router.POST("/pvz", CreatePVZ)

	// Примонтировать фейковое подключение
	storage.DB = setupInMemoryDB(t)

	body := []byte(`{"city": "Москва"}`)
	req, _ := http.NewRequest("POST", "/pvz", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", w.Code)
	}
}

func TestCreatePVZ_InvalidCity(t *testing.T) {
	router := gin.New()
	router.POST("/pvz", CreatePVZ)

	storage.DB = setupInMemoryDB(t)

	body := []byte(`{"city": "Алмата"}`)
	req, _ := http.NewRequest("POST", "/pvz", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", w.Code)
	}
}
