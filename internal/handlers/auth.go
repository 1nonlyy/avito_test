package handlers

import (
	"avito-test/internal/utils"
	"net/http"

	"avito-test/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type DummyLoginRequest struct {
	Role string `json:"role"`
}
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // client или moderator
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func DummyLogin(c *gin.Context) {
	var req DummyLoginRequest
	if err := c.BindJSON(&req); err != nil || (req.Role != "employee" && req.Role != "moderator") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid role"})
		return
	}
	userID := uuid.New().String()
	token, err := utils.GenerateToken(userID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "token generation failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	if req.Role != "client" && req.Role != "moderator" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid role"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot hash password"})
		return
	}

	userID := uuid.New().String()
	_, err = storage.DB.Exec(
		"INSERT INTO users (id, email, password_hash, role) VALUES ($1, $2, $3, $4)",
		userID, req.Email, string(hash), req.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "user already exists or db error"})
		return
	}

	token, _ := utils.GenerateToken(userID, req.Role)
	c.JSON(http.StatusCreated, gin.H{"token": token})
}
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	var id, hash, role string
	err := storage.DB.QueryRow(
		"SELECT id, password_hash, role FROM users WHERE email = $1", req.Email,
	).Scan(&id, &hash, &role)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	token, _ := utils.GenerateToken(id, role)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
