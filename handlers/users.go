package handlers

import (
	"book-management-api/config"
	"book-management-api/models"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	DB  *sql.DB
	Cfg *config.Config
}

func NewUserHandler(db *sql.DB, cfg *config.Config) *UserHandler {
	return &UserHandler{
		DB:  db,
		Cfg: cfg,
	}
}

func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Log the login attempt for debugging
	c.Header("Content-Type", "application/json")

	// Get user from database
	var user models.User
	query := `SELECT id, username, password, created_at, created_by, modified_at, modified_by 
			  FROM users WHERE username = $1`
	
	err := h.DB.QueryRow(query, req.Username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.ModifiedAt,
		&user.ModifiedBy,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Invalid credentials",
			Error:   "user not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database error",
			Error:   err.Error(),
		})
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Invalid credentials",
			Error:   "invalid password",
		})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.Cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate token",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data: models.LoginResponse{
			Token: tokenString,
			User:  user,
		},
	})
}

// SeedAdmin creates default admin user if not exists
func (h *UserHandler) SeedAdmin(c *gin.Context) {
	// Check if admin user already exists
	var count int
	checkQuery := `SELECT COUNT(*) FROM users WHERE username = 'admin'`
	err := h.DB.QueryRow(checkQuery).Scan(&count)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database error",
			Error:   err.Error(),
		})
		return
	}
	
	if count > 0 {
		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Message: "Admin user already exists",
			Data:    gin.H{"username": "admin"},
		})
		return
	}
	
	// Hash password "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to hash password",
			Error:   err.Error(),
		})
		return
	}
	
	// Insert admin user
	insertQuery := `INSERT INTO users (username, password, created_by) VALUES ($1, $2, $3)`
	_, err = h.DB.Exec(insertQuery, "admin", string(hashedPassword), "system")
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create admin user",
			Error:   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Admin user created successfully",
		Data:    gin.H{"username": "admin", "password": "admin123"},
	})
}

// ResetAdminPassword resets admin password to "admin123"
func (h *UserHandler) ResetAdminPassword(c *gin.Context) {
	// Hash password "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to hash password",
			Error:   err.Error(),
		})
		return
	}
	
	// Update admin user password
	_, err = h.DB.Exec(`UPDATE users SET password = $1, modified_at = CURRENT_TIMESTAMP WHERE username = 'admin'`, string(hashedPassword))
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to reset admin password",
			Error:   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Admin password reset successfully",
		Data:    gin.H{"username": "admin", "password": "admin123"},
	})
}