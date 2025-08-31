package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thelostleo/CTF-backend/database"
	"github.com/thelostleo/CTF-backend/models"
	"github.com/thelostleo/CTF-backend/utils"
)

type UserController struct{}

// hashPassword creates a SHA-256 hash of the password
func (uc *UserController) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// verifyPassword compares a plain password with a hashed password
func (uc *UserController) verifyPassword(hashedPassword, password string) bool {
	return uc.hashPassword(password) == hashedPassword
}

// Register handles user registration
func (uc *UserController) Register(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Check if username already exists
	var existingUser models.User
	if err := database.DB.Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Username already exists",
		})
		return
	}

	// Check if email already exists
	if err := database.DB.Where("email = ?", request.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email already exists",
		})
		return
	}

	// Create new user
	user := models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: uc.hashPassword(request.Password),
		Score:    0,
		IsAdmin:  false, // Default to regular user
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"is_admin": user.IsAdmin,
		},
	})
}

// Login handles user authentication
func (uc *UserController) Login(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Find user by username
	var user models.User
	if err := database.DB.Where("username = ?", request.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Verify password
	if !uc.verifyPassword(user.Password, request.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"is_admin": user.IsAdmin,
			"score":    user.Score,
		},
	})
}

// GetProfile handles getting user profile
func (uc *UserController) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Find user
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"is_admin": user.IsAdmin,
			"score":    user.Score,
		},
	})
}

// GetLeaderboard handles getting the leaderboard
func (uc *UserController) GetLeaderboard(c *gin.Context) {
	var users []models.User

	// Get top 10 users by score
	if err := database.DB.Select("id, username, score").
		Order("score DESC").
		Limit(10).
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch leaderboard",
		})
		return
	}

	// Create leaderboard response
	leaderboard := make([]gin.H, len(users))
	for i, user := range users {
		leaderboard[i] = gin.H{
			"rank":     i + 1,
			"username": user.Username,
			"score":    user.Score,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"leaderboard": leaderboard,
	})
}

// RefreshToken handles JWT token refresh
func (uc *UserController) RefreshToken(c *gin.Context) {
	// Get the current token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization header required",
		})
		return
	}

	// Extract token from "Bearer TOKEN" format
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authorization format",
		})
		return
	}

	oldToken := authHeader[7:]

	// Generate a new token
	newToken, err := utils.RefreshJWTToken(oldToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token or token expired",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"token":   newToken,
	})
}
