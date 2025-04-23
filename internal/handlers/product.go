package handlers

import (
	"avito-test/internal/models"
	"avito-test/internal/storage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddProduct(c *gin.Context) {
	var req struct {
		Type  string `json:"type"`
		PVZID string `json:"pvzId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || (req.Type != "электроника" && req.Type != "одежда" && req.Type != "обувь") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid type or request"})
		return
	}

	// Найти открытую приёмку для указанного ПВЗ
	var receptionID string
	err := storage.DB.QueryRow(
		"SELECT id FROM receptions WHERE pvz_id = $1 AND status = 'in_progress' ORDER BY date_time DESC LIMIT 1",
		req.PVZID,
	).Scan(&receptionID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no open reception"})
		return
	}

	// Добавить товар
	id := uuid.New().String()
	now := time.Now().Format(time.RFC3339)

	_, err = storage.DB.Exec(
		"INSERT INTO products (id, date_time, type, reception_id) VALUES ($1, $2, $3, $4)",
		id, now, req.Type, receptionID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "db error"})
		return
	}

	c.JSON(http.StatusCreated, models.Product{
		ID:          id,
		DateTime:    now,
		Type:        req.Type,
		ReceptionID: receptionID,
	})
}
func DeleteLastProduct(c *gin.Context) {
	pvzID := c.Param("pvzId")

	// Найти открытую приёмку
	var receptionID string
	err := storage.DB.QueryRow(
		"SELECT id FROM receptions WHERE pvz_id = $1 AND status = 'in_progress' ORDER BY date_time DESC LIMIT 1",
		pvzID,
	).Scan(&receptionID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no open reception"})
		return
	}

	// Найти последний добавленный товар (по дате или ID)
	var productID string
	err = storage.DB.QueryRow(
		"SELECT id FROM products WHERE reception_id = $1 ORDER BY date_time DESC LIMIT 1",
		receptionID,
	).Scan(&productID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no products to delete"})
		return
	}

	// Удалить товар
	_, err = storage.DB.Exec("DELETE FROM products WHERE id = $1", productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted", "id": productID})
}
