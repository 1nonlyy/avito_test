package handlers

import (
	"avito-test/internal/models"
	"avito-test/internal/storage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateReception(c *gin.Context) {
	var req struct {
		PVZID string `json:"pvzId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	// Проверка: есть ли уже открытая приёмка
	var count int
	err := storage.DB.QueryRow(
		"SELECT COUNT(*) FROM receptions WHERE pvz_id = $1 AND status = 'in_progress'",
		req.PVZID,
	).Scan(&count)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "db error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "already have open reception"})
		return
	}

	// Создание новой приёмки
	id := uuid.New().String()
	now := time.Now().Format(time.RFC3339)
	status := "in_progress"

	_, err = storage.DB.Exec(
		"INSERT INTO receptions (id, date_time, pvz_id, status) VALUES ($1, $2, $3, $4)",
		id, now, req.PVZID, status,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "insert error"})
		return
	}

	c.JSON(http.StatusCreated, models.Reception{
		ID:       id,
		DateTime: now,
		PVZID:    req.PVZID,
		Status:   status,
	})
}
func CloseLastReception(c *gin.Context) {
	pvzID := c.Param("pvzId")

	// Найти последнюю открытую приёмку
	var reception models.Reception
	err := storage.DB.QueryRow(
		"SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = $1 AND status = 'in_progress' ORDER BY date_time DESC LIMIT 1",
		pvzID,
	).Scan(&reception.ID, &reception.DateTime, &reception.PVZID, &reception.Status)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no open reception"})
		return
	}

	// Обновить статус приёмки
	_, err = storage.DB.Exec(
		"UPDATE receptions SET status = 'close' WHERE id = $1",
		reception.ID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to close reception"})
		return
	}

	reception.Status = "close"
	c.JSON(http.StatusOK, reception)
}
