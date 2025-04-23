package main

import (
	"avito-test/internal/handlers"
	"avito-test/internal/middleware"
	"avito-test/internal/storage"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := storage.InitDB(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	r := gin.Default()

	r.POST("/dummyLogin", handlers.DummyLogin)
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.POST("/pvz", middleware.RequireRole("moderator"), handlers.CreatePVZ)
	r.POST("/receptions", middleware.RequireRole("employee"), handlers.CreateReception)
	r.POST("/products", middleware.RequireRole("employee"), handlers.AddProduct)
	r.POST("/pvz/:pvzId/delete_last_product", middleware.RequireRole("employee"), handlers.DeleteLastProduct)
	r.POST("/pvz/:pvzId/close_last_reception", middleware.RequireRole("employee"), handlers.CloseLastReception)
	r.GET("/pvz", middleware.RequireRole("moderator", "employee"), handlers.GetPVZList)

	r.Run(":8080")
}
