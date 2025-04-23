package handlers

import (
	"avito-test/internal/models"
	"avito-test/internal/storage"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var allowedCities = map[string]bool{
	"Москва":          true,
	"Казань":          true,
	"Санкт-Петербург": true,
}

func CreatePVZ(c *gin.Context) {
	var req struct {
		City string `json:"city"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || !allowedCities[req.City] {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid city"})
		return
	}

	id := uuid.New().String()
	now := time.Now().Format(time.RFC3339)

	_, err := storage.DB.Exec(
		"INSERT INTO pvz (id, registration_date, city) VALUES ($1, $2, $3)",
		id, now, req.City,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "db error"})
		return
	}

	c.JSON(http.StatusCreated, models.PVZ{
		ID:               id,
		RegistrationDate: now,
		City:             req.City,
	})
}
func GetPVZList(c *gin.Context) {
	start := c.DefaultQuery("startDate", "")
	end := c.DefaultQuery("endDate", "")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	// Собираем основной список ПВЗ
	rows, err := storage.DB.Query(
		"SELECT id, registration_date, city FROM pvz ORDER BY registration_date DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "db error"})
		return
	}
	defer rows.Close()

	var result []models.PVZWithReceptions

	for rows.Next() {
		var pvz models.PVZ
		if err := rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
			continue
		}

		// Получаем приёмки
		receptionQuery := "SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = $1"
		args := []interface{}{pvz.ID}

		if start != "" {
			receptionQuery += " AND date_time >= $2"
			args = append(args, start)
		}
		if end != "" {
			receptionQuery += " AND date_time <= $3"
			args = append(args, end)
		}

		receptionRows, err := storage.DB.Query(receptionQuery, args...)
		if err != nil {
			continue
		}

		var receptions []models.ReceptionWithItems

		for receptionRows.Next() {
			var r models.Reception
			if err := receptionRows.Scan(&r.ID, &r.DateTime, &r.PVZID, &r.Status); err != nil {
				continue
			}

			// Получаем товары по приёмке
			products := []models.Product{}
			productRows, err := storage.DB.Query(
				"SELECT id, date_time, type, reception_id FROM products WHERE reception_id = $1",
				r.ID,
			)
			if err == nil {
				for productRows.Next() {
					var p models.Product
					_ = productRows.Scan(&p.ID, &p.DateTime, &p.Type, &p.ReceptionID)
					products = append(products, p)
				}
				productRows.Close()
			}

			receptions = append(receptions, models.ReceptionWithItems{
				Reception: r,
				Products:  products,
			})
		}
		receptionRows.Close()

		result = append(result, models.PVZWithReceptions{
			PVZ:        pvz,
			Receptions: receptions,
		})
	}

	c.JSON(http.StatusOK, result)
}
