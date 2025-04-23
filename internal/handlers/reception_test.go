package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateReception_Success(t *testing.T) {
	db := setupReceptionDB(t)

	router := gin.Default()
	router.POST("/receptions", CreateReception)

	body := []byte(`{"pvzId":"pvz-1"}`)
	req, _ := http.NewRequest("POST", "/receptions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", w.Code)
	}

	_ = db.Close()
}

func TestCreateReception_AlreadyInProgress(t *testing.T) {
	db := setupReceptionDB(t)

	// Добавим приёмку in_progress вручную
	_, _ = db.Exec(`INSERT INTO receptions (id, date_time, pvz_id, status)
	                VALUES ('r-1', '2025-01-01T01:00:00Z', 'pvz-1', 'in_progress');`)

	router := gin.Default()
	router.POST("/receptions", CreateReception)

	body := []byte(`{"pvzId":"pvz-1"}`)
	req, _ := http.NewRequest("POST", "/receptions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", w.Code)
	}

	_ = db.Close()
}
