package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"yourusername/gym-planner/internal/auth"
	"yourusername/gym-planner/internal/database"
	"yourusername/gym-planner/internal/models"
)

type UserHandler struct {
	db     *database.DB
	jwtMgr *auth.JWTManager
}

func NewUserHandler(db *database.DB, jwtMgr *auth.JWTManager) *UserHandler {
	return &UserHandler{db: db, jwtMgr: jwtMgr}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error processing password"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	var id int64
	err = h.db.QueryRow(
		"INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id",
		user.Username, user.Password,
	).Scan(&id)

	if err != nil {
		c.JSON(500, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(201, gin.H{"id": id})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	var user models.User
	err := h.db.Get(&user, "SELECT * FROM users WHERE username = $1", req.Username)
	if err == sql.ErrNoRows {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := h.jwtMgr.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (h *UserHandler) Logout(c *gin.Context) {
	// Since we're using JWT, we don't need to do anything server-side
	// The client should remove the token
	c.Status(200)
} 